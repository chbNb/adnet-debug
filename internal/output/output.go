package output

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/gogo/protobuf/proto"
	"gitlab.mobvista.com/ADN/adnet/internal/extractor"
	"gitlab.mobvista.com/ADN/adnet/internal/mvconst"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	"gitlab.mobvista.com/ADN/adnet/internal/protobuf"
	"gitlab.mobvista.com/ADN/adnet/internal/redis"
	"gitlab.mobvista.com/ADN/adnet/internal/watcher"
	smodel "gitlab.mobvista.com/ADN/structs/model"
	"gitlab.mobvista.com/adserver/recommend_protocols/go/ad_server"
)

const (
	TAB_CAMP     = "campaign"
	TimeoutMongo = 200
)

// AdData 定义最终接口返回的数据
type AdData struct {
	SessionID         string `json:"session_id"`
	ParentSessionID   string `json:"parent_session_id"`
	AdType            uint8  `json:"ad_type"`
	Template          uint8  `json:"template"`
	UnitSize          string `json:"unit_size"`
	HTMLURL           string `json:"html_url"`
	OnlyImpressionURL string `json:"only_impression_url"`
	Ads               []Ad   `json:"ads"`
}

func GetCampaignsInfo(campaignIds []int64) (map[int64]*smodel.CampaignInfo, error) {
	if len(campaignIds) == 0 {
		return nil, errors.New("params campaignIds is empty error")
	}
	result := make(map[int64]*smodel.CampaignInfo, len(campaignIds))
	for _, camId := range campaignIds {
		watcher.AddWatchValue("get_campaign_data_count", float64(1))
		if camInfo, ifind := extractor.GetCampaignInfo(camId); ifind {
			result[camId] = camInfo
		} else {
			watcher.AddWatchValue("get_campaign_data_error_count", float64(1))
		}
	}
	return result, nil
}

func GetCreativePbInfos(params1, params2 map[int64]map[string]map[ad_server.CreativeType]int64, r *mvutil.RequestParams, campaigns map[int64]*smodel.CampaignInfo) (map[int64]map[ad_server.CreativeType]*protobuf.Creative, error) {
	var creative1, creative2 map[int64]map[ad_server.CreativeType]*protobuf.Creative
	var err1, err2 error
	// appwall 才从　campaign mongo 中取
	if mvutil.IsAppwallOrMoreOffer(r.Param.AdType) &&
		r.Param.Scenario == mvconst.SCENARIO_OPENAPI {
		creative2, err1 = GetCreativePbNormalWall(params2, campaigns, r)
		creative1, err2 = GetCreativePbNormalWall(params1, campaigns, r)
	} else {
		creative2, err1 = GetCreativePbPipeline(params2)
		creative1, err2 = GetCreativePbPipeline(params1)
	}

	if err1 != nil && err2 != nil {
		return nil, err1
	}
	if err2 != nil {
		return creative1, nil
	}
	for c, cMap := range creative1 {
		if c2, ok := creative2[c]; ok {
			for k, v := range c2 {
				cMap[k] = v
			}
		}
	}
	return creative1, nil
}

func GetCreativePbInfosV2(crParams map[int64]map[string]map[ad_server.CreativeType]int64,
	r *mvutil.RequestParams,
	campaigns map[int64]*smodel.CampaignInfo) (map[int64]map[ad_server.CreativeType]*protobuf.Creative, error) {
	var creativeData map[int64]map[ad_server.CreativeType]*protobuf.Creative
	var err error
	// appwall 才从　campaign mongo 中取
	if mvutil.IsAppwallOrMoreOffer(r.Param.AdType) &&
		r.Param.Scenario == mvconst.SCENARIO_OPENAPI {
		creativeData, err = GetCreativePbNormalWall(crParams, campaigns, r)
	} else {
		creativeData, err = GetCreativePbPipeline(crParams)
	}
	if err != nil {
		return nil, err
	}
	return creativeData, nil
}

// 对 appwall 类型的creative读取逻辑
func GetCreativePbNormalWall(params map[int64]map[string]map[ad_server.CreativeType]int64,
	campaigns map[int64]*smodel.CampaignInfo,
	r *mvutil.RequestParams) (map[int64]map[ad_server.CreativeType]*protobuf.Creative, error) {
	if len(params) == 0 {
		return nil, errors.New("GetCreativePbPipeline Params Invalidate")
	}

	paramsRedis := make(map[int64]map[string]map[ad_server.CreativeType]int64)
	paramsMongo := make(map[int64]map[string]map[ad_server.CreativeType]int64)
	for camId, v := range params {
		cam, ifFind := campaigns[camId]
		if !ifFind {
			continue
		}
		// 请求切量，还要根据readCreative 来决定最终是否切。
		// readCreative=1走原逻辑，=2才走campaign
		if cam.ReadCreative == 2 {
			paramsMongo[camId] = v
		} else {
			paramsRedis[camId] = v
		}
	}
	resultRedis := make(map[int64]map[ad_server.CreativeType]*protobuf.Creative)
	resultMongo := make(map[int64]map[ad_server.CreativeType]*protobuf.Creative)
	var err error
	// 先判断paramsXXX 长度，避免直接调用对应函数产生报错
	if len(paramsRedis) > 0 {
		resultRedis, err = GetCreativePbPipeline(paramsRedis)
		if err != nil {
			watcher.AddWatchValue("appwall_cr_redis_err", 1)
			paramLog, _ := json.Marshal(paramsRedis)
			mvutil.Logger.Runtime.Warnf("appwall_creative_redis error, paramRedis: " + string(paramLog))
			return nil, err
		}
	}
	if len(paramsMongo) > 0 {
		resultMongo, err = GetCreativePbMongo(paramsMongo, campaigns)
		if err != nil {
			watcher.AddWatchValue("appwall_cr_mongo_err", 1)
			paramLog, _ := json.Marshal(paramsMongo)
			mvutil.Logger.Runtime.Warnf("appwall_creative_mongo error, paramMongo: " + string(paramLog))
			return nil, err
		}
	}
	for k, v := range resultRedis {
		resultMongo[k] = v
	}

	return resultMongo, nil

}

// CreativeTypeRawString 返回原数字的字符串 eg: CreativeType_APP_NAME -> 401 -> "401"
func CreativeTypeRawString(p ad_server.CreativeType) string {
	return strconv.FormatInt(int64(p), 10)
}

// GetCreativePbMongo 从 mongo campaign 表中获取creative
// http://confluence.mobvista.com/pages/viewpage.action?pageId=8685701
func GetCreativePbMongo(params map[int64]map[string]map[ad_server.CreativeType]int64,
	campaigns map[int64]*smodel.CampaignInfo) (map[int64]map[ad_server.CreativeType]*protobuf.Creative, error) {
	if len(params) == 0 {
		return nil, errors.New("GetCreativePbPipeline Params Invalidate")
	}
	result := make(map[int64]map[ad_server.CreativeType]*protobuf.Creative, len(params))
	for camId, v := range params { // c: campaignid
		campaign, ifFind := campaigns[camId]
		if !ifFind {
			continue
		}
		tmpCreatives := make(map[ad_server.CreativeType]*protobuf.Creative) // crIdInt -> cr
		for _, vv := range v {
			for cType, _ := range vv {
				creative := new(protobuf.Creative)
				switch cType {
				case ad_server.CreativeType_APP_NAME:
					if campaign.BasicCrList == nil {
						continue
					}
					creative.SValue = campaign.BasicCrList.AppName
				case ad_server.CreativeType_APP_DESC:
					if campaign.BasicCrList == nil {
						continue
					}
					creative.SValue = campaign.BasicCrList.AppDesc
				case ad_server.CreativeType_ICON:
					if campaign.BasicCrList == nil {
						continue
					}
					creative.SValue = campaign.BasicCrList.AppIcon
				case ad_server.CreativeType_APP_RATE:
					if campaign.BasicCrList == nil {
						continue
					}
					creative.FValue = campaign.BasicCrList.AppRate
				case ad_server.CreativeType_COMMENT:
					if campaign.BasicCrList == nil {
						continue
					}
					creative.IValue = int32(campaign.BasicCrList.NumRating)
				case ad_server.CreativeType_CTA_BUTTON:
					creative.SValue = campaign.BasicCrList.CtaButton
				}
				tmpCreatives[cType] = creative
			}
		}
		result[camId] = tmpCreatives
	}
	return result, nil
}

//func GetCreativePbPipeline(params map[int64]map[string]map[ad_server.CreativeType]int64) (map[int64]map[int64]*protobuf.Creative, error) {
//	if len(params) == 0 {
//		return nil, errors.New("GetCreativePbPipeline Params Invalidate")
//	}
//	args := make(map[string][]string, len(params))
//	creativeIdList := make(map[int64]map[string][]int64, len(params))
//	// creativeIdList:  campaignId: creativeIdStr:[]creativeId
//	for c, v := range params {
//		tmpCreativeIdMap := make(map[string][]int64, len(v))
//		for k, vv := range v {
//			key := "creative:" + k
//			tmpArgs := make([]string, 0, len(vv))
//			tmpCreativeIdList := make([]int64, 0, len(vv))
//			for _, vvv := range vv {
//				tmpArgs = append(tmpArgs, strconv.FormatInt(vvv, 10))
//				tmpCreativeIdList = append(tmpCreativeIdList, vvv)
//			}
//			args[key] = tmpArgs
//			tmpCreativeIdMap[key] = tmpCreativeIdList
//		}
//		creativeIdList[c] = tmpCreativeIdMap
//	}
//	creativeInfos, err := redis.LocalRedisHMGetPipeline(args)
//	if err != nil {
//		return nil, err
//	}
//	ret := make(map[int64]map[int64]*protobuf.Creative, len(params))
//	var tmpPb map[int64]*protobuf.Creative
//	for c, vmap := range creativeIdList {
//		tmpPb = make(map[int64]*protobuf.Creative, 9)
//		for k, vv := range vmap {
//			if vc, ok := creativeInfos[k]; ok {
//				for kk, vvv := range vv {
//					creativeStr := vc[kk]
//					if len(creativeStr) == 0 {
//						continue
//					}
//					creativeInfo := &protobuf.Creative{}
//					err := proto.Unmarshal([]byte(creativeStr), creativeInfo)
//					if err != nil {
//						continue
//					}
//					tmpPb[vvv] = creativeInfo
//				}
//			}
//		}
//		ret[c] = tmpPb
//	}
//	return ret, nil
//}

func GetCreativePbPipeline(params map[int64]map[string]map[ad_server.CreativeType]int64) (map[int64]map[ad_server.CreativeType]*protobuf.Creative, error) {
	if len(params) == 0 {
		return nil, errors.New("GetCreativePbPipeline Params Invalidate")
	}
	args := make(map[string][]string, len(params))
	creativeIdList := make(map[int64]map[string][]ad_server.CreativeType, len(params))
	// creativeIdList:  campaignId: creativeIdStr:[]creativeId
	for c, v := range params {
		tmpCreativeIdMap := make(map[string][]ad_server.CreativeType, len(v))
		for k, vv := range v {
			key := "creative:" + k
			tmpArgs := make([]string, 0, len(vv))
			tmpCreativeIdList := make([]ad_server.CreativeType, 0, len(vv))
			for crType, vvv := range vv {
				tmpArgs = append(tmpArgs, strconv.FormatInt(vvv, 10))
				tmpCreativeIdList = append(tmpCreativeIdList, crType)
			}
			args[key] = tmpArgs
			tmpCreativeIdMap[key] = tmpCreativeIdList
		}
		creativeIdList[c] = tmpCreativeIdMap
	}
	creativeInfos, err := redis.LocalRedisHMGetPipeline(args)
	if err != nil {
		mvutil.Logger.Runtime.Warnf("LocalRedisHMGetPipeline error=[%s]", err.Error())
		return nil, err
	}
	ret := make(map[int64]map[ad_server.CreativeType]*protobuf.Creative, len(params))
	var tmpPb map[ad_server.CreativeType]*protobuf.Creative
	for c, vmap := range creativeIdList {
		tmpPb = make(map[ad_server.CreativeType]*protobuf.Creative, 9)
		for k, vv := range vmap {
			if vc, ok := creativeInfos[k]; ok {
				for kk, vvv := range vv {
					creativeStr := vc[kk]
					if len(creativeStr) == 0 {
						continue
					}
					creativeInfo := &protobuf.Creative{}
					err := proto.Unmarshal([]byte(creativeStr), creativeInfo)
					if err != nil {
						continue
					}
					tmpPb[vvv] = creativeInfo
				}
			}
		}
		ret[c] = tmpPb
	}
	return ret, nil
}

func GetDcoMaterialData(dcoMaterialMIds map[int64]map[ad_server.CreativeType]int64) map[int64]map[ad_server.CreativeType]*protobuf.Material {
	if len(dcoMaterialMIds) == 0 {
		return nil
	}
	materialData, err := getMaterialData(dcoMaterialMIds)
	if err != nil {
		return nil
	}
	return materialData
}

func GetDcoOfferMaterialData(dcoOfferMaterialOmIds map[int64]map[ad_server.CreativeType]int64) map[int64]map[ad_server.CreativeType]*protobuf.OfferMaterial {
	if len(dcoOfferMaterialOmIds) == 0 {
		return nil
	}
	offerMaterialData, err := getOfferMaterialData(dcoOfferMaterialOmIds)
	if err != nil {
		return nil
	}
	return offerMaterialData
}

func getMaterialData(dcoMaterialMIds map[int64]map[ad_server.CreativeType]int64) (map[int64]map[ad_server.CreativeType]*protobuf.Material, error) {
	// 获取mId list
	mIds := make([]string, 0)
	for _, val := range dcoMaterialMIds {
		for _, mId := range val {
			mIds = append(mIds, strconv.FormatInt(mId, 10))
		}
	}
	// 根据mid查询material
	materialInfos, err := redis.LocalRedisHMGet(mvconst.Material, mIds)
	if err != nil {
		mvutil.Logger.Runtime.Warnf("getMaterialData mIds=[%s] error=[%s]", mIds, err.Error())
		watcher.AddWatchValue("material_hmget_err", 1)
		return nil, err
	}
	// 按照顺序获取mid对应的素材结果
	materialMIdInfos := make(map[int64]*protobuf.Material)
	for k, materialStr := range materialInfos {
		if len(materialStr) == 0 {
			mvutil.Logger.Runtime.Warnf("material_empty mId=[%s]", mIds[k])
			watcher.AddWatchValue("material_empty", 1)
			continue
		}
		materialInfo := &protobuf.Material{}
		err := proto.Unmarshal([]byte(materialStr), materialInfo)
		if err != nil {
			mvutil.Logger.Runtime.Warnf("material_unmarshal_err mId=[%s]", mIds[k])
			watcher.AddWatchValue("material_unmarshal_err", 1)
			continue
		}

		mId, err := strconv.ParseInt(mIds[k], 10, 64)
		if err == nil {
			materialMIdInfos[mId] = materialInfo
		}

	}

	result := make(map[int64]map[ad_server.CreativeType]*protobuf.Material)
	// 根据campaignId & creative type分配查询结果
	for camId, val := range dcoMaterialMIds {
		crTypeMaterialInfo := make(map[ad_server.CreativeType]*protobuf.Material)
		for crType, mId := range val {
			if materialInfo, ok := materialMIdInfos[mId]; ok {
				crTypeMaterialInfo[crType] = materialInfo
			}
		}
		result[camId] = crTypeMaterialInfo
	}
	return result, nil
}

func getOfferMaterialData(dcoOfferMaterialOmIds map[int64]map[ad_server.CreativeType]int64) (map[int64]map[ad_server.CreativeType]*protobuf.OfferMaterial, error) {
	// 获取omId list
	omIds := make([]string, 0)
	for _, val := range dcoOfferMaterialOmIds {
		for _, omId := range val {
			omIds = append(omIds, strconv.FormatInt(omId, 10))
		}
	}
	// 根据omid查询offer_material
	offerMaterialInfos, err := redis.LocalRedisHMGet(mvconst.OfferMaterial, omIds)
	if err != nil {
		mvutil.Logger.Runtime.Warnf("getOfferMaterialData mIds=[%s] error=[%s]", omIds, err.Error())
		watcher.AddWatchValue("offer_material_hmget_err", 1)
		return nil, err
	}
	// 按照顺序获取omid对应的素材关系结果
	offerMaterialOmIdInfos := make(map[int64]*protobuf.OfferMaterial)
	for k, offerMaterialStr := range offerMaterialInfos {
		if len(offerMaterialStr) == 0 {
			mvutil.Logger.Runtime.Warnf("offer_material_empty omId=[%s]", omIds[k])
			watcher.AddWatchValue("offer_material_empty", 1)
			continue
		}
		offerMaterialInfo := &protobuf.OfferMaterial{}
		err := proto.Unmarshal([]byte(offerMaterialStr), offerMaterialInfo)
		if err != nil {
			mvutil.Logger.Runtime.Warnf("material_unmarshal_err omId=[%s]", omIds[k])
			watcher.AddWatchValue("material_unmarshal_err", 1)
			continue
		}

		omId, err := strconv.ParseInt(omIds[k], 10, 64)
		if err == nil {
			offerMaterialOmIdInfos[omId] = offerMaterialInfo
		}

	}

	result := make(map[int64]map[ad_server.CreativeType]*protobuf.OfferMaterial)
	// 根据campaignId & creative type分配查询结果
	for camId, val := range dcoOfferMaterialOmIds {
		crTypeOfferMaterialInfo := make(map[ad_server.CreativeType]*protobuf.OfferMaterial)
		for crType, omId := range val {
			if materialInfo, ok := offerMaterialOmIdInfos[omId]; ok {
				// redis中并不会同步omid，as会返回，因此使用as返回的adn creative id
				materialInfo.OmID = omId
				crTypeOfferMaterialInfo[crType] = materialInfo
			}
		}
		result[camId] = crTypeOfferMaterialInfo
	}
	return result, nil
}

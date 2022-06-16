package backend

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
	"time"

	"gitlab.mobvista.com/ADN/chasm/module/demand"

	jsoniter "github.com/json-iterator/go"
	"gitlab.mobvista.com/ADN/adnet/internal/corsair_proto"
	"gitlab.mobvista.com/ADN/adnet/internal/extractor"
	"gitlab.mobvista.com/ADN/adnet/internal/mvconst"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	"gitlab.mobvista.com/ADN/adnet/internal/req_context"
	"gitlab.mobvista.com/adserver/recommend_protocols/go/ad_server"
)

type AdServer struct {
	Backend
}

func (backend AdServer) getRequestNode() string {
	return ""
}

func (backend AdServer) filterBackend(reqCtx *mvutil.ReqCtx) int {
	return mvconst.BackendOK
}

func ConstructASRequest(reqCtx *mvutil.ReqCtx) (*ad_server.QueryParam, error) {
	return composeAdServerRequest(reqCtx, mvconst.MAdx)
}

func composeAdServerRequest(reqCtx *mvutil.ReqCtx, bid int) (*ad_server.QueryParam, error) {
	if reqCtx == nil || reqCtx.ReqParams == nil {
		return nil, errors.New("composeAdServerRequest params is invalidate")
	}

	var req ad_server.QueryParam

	// fill request
	req.Timestamp = time.Now().UTC().Unix()
	req.AppId = &reqCtx.ReqParams.Param.AppID
	req.UnitId = &reqCtx.ReqParams.Param.UnitID
	req.Scenario = &reqCtx.ReqParams.Param.Scenario
	adTypeStr := mvutil.GetAdTypeStr(reqCtx.ReqParams.Param.AdType)
	if reqCtx.ReqParams.Param.RequestPath == mvconst.PATHREQUEST && len(adTypeStr) <= 0 {
		adTypeStr = reqCtx.ReqParams.Param.AdTypeStr
	}
	req.AdTypeStr = &adTypeStr

	req.ExcludeIdSet = renderExcludeIDs(reqCtx.ReqParams.Param.ExcludeIDS)
	req.AdNum = reqCtx.ReqParams.Param.AdNum
	req.ImageSizeId = ad_server.ImageSizeEnum(reqCtx.ReqParams.Param.ImageSizeID)
	req.RequestType = ad_server.RequestType(reqCtx.ReqParams.Param.RequestType)
	req.UnitSize = &reqCtx.ReqParams.Param.UnitSize
	req.Category = ad_server.Category(reqCtx.ReqParams.Param.Category)
	recallAdnOffer := int32(1)
	req.RecallADNOffer = &recallAdnOffer
	req.Platform = ad_server.Platform(reqCtx.ReqParams.Param.Platform)
	req.OsVersion = reqCtx.ReqParams.Param.OSVersion
	req.SdkVersion = reqCtx.ReqParams.Param.SDKVersion
	req.PackageName = reqCtx.ReqParams.Param.AppPackageName
	req.AppVersionName = reqCtx.ReqParams.Param.AppVersionName
	req.AppVersionCode = reqCtx.ReqParams.Param.AppVersionCode
	req.Imei = reqCtx.ReqParams.Param.IMEI
	req.Mac = reqCtx.ReqParams.Param.MAC
	req.DevId = reqCtx.ReqParams.Param.AndroidID
	req.DeviceModel = reqCtx.ReqParams.Param.Model
	req.UnifiedDevModel = &reqCtx.ReqParams.Param.ReplaceModel
	req.ScreenSize = reqCtx.ReqParams.Param.ScreenSize
	req.Orientation = ad_server.Orientation(reqCtx.ReqParams.Param.Orientation)
	req.Mnc = reqCtx.ReqParams.Param.MNC
	req.Mcc = reqCtx.ReqParams.Param.MCC
	req.NetworkType = ad_server.NetworkType(reqCtx.ReqParams.Param.NetworkType)
	req.Language = reqCtx.ReqParams.Param.Language
	req.IP = reqCtx.ReqParams.Param.ClientIP
	req.AdnServerIp = req_context.GetInstance().ServerIP()
	req.CountryCode = reqCtx.ReqParams.Param.CountryCode
	sessionId := mvutil.GetRequestID()
	req.SessionId = &sessionId
	pSessionId := mvutil.GetRequestID()
	req.ParentSessionId = &pSessionId
	req.Timezone = &reqCtx.ReqParams.Param.TimeZone
	req.GPVersion = reqCtx.ReqParams.Param.GPVersion
	req.AdSourceList = renderAdSource(reqCtx.ReqParams)
	req.GoogleAdvertisingId = &reqCtx.ReqParams.Param.GAID
	req.DEPRECATEDPublisherType = ad_server.PublisherType_M
	req.NetworkId = &reqCtx.ReqParams.Param.Network
	req.RequestId = &reqCtx.ReqParams.Param.RequestID
	req.UnitSize = &reqCtx.ReqParams.Param.UnitSize
	req.Offset = reqCtx.ReqParams.Param.Offset

	if len(reqCtx.ReqParams.Param.NativeInfoList) > 0 {
		for _, v := range reqCtx.ReqParams.Param.NativeInfoList {
			nativeInfo := ad_server.NewNativeInfo()
			nativeInfo.AdTemplate = ad_server.ADTemplate(v.AdTemplate)
			nativeInfo.RequireNum = int64(reqCtx.ReqParams.Param.AdNum)
			req.NativeInfoList = append(req.NativeInfoList, nativeInfo)
		}
	}

	req.ShowedCampaignIdList = reqCtx.ReqParams.Param.DisplayCamIds
	i64FrameNum := int64(reqCtx.ReqParams.Param.FrameNum)
	req.FrameNum = &i64FrameNum
	req.Idfa = &reqCtx.ReqParams.Param.IDFA
	tnum := int64(reqCtx.ReqParams.Param.TNum)
	// 达到展示控cap的条件，trueNum传0给AdServer
	if reqCtx.ReqParams.Param.IsBlockByImpCap {
		tnum = int64(0)
	}
	req.TrueNum = &tnum
	req.InstallIdSet = renderExcludeIDs(reqCtx.ReqParams.Param.InstallIDS)
	req.DeviceBrand = &reqCtx.ReqParams.Param.Brand
	req.UnifiedDevBrand = &reqCtx.ReqParams.Param.ReplaceBrand
	req.CityCode = &reqCtx.ReqParams.Param.CityCode
	req.PriceFloor = &reqCtx.ReqParams.Param.PriceFloor

	req.RealAppId = &reqCtx.ReqParams.Param.RealAppID
	req.UnSupportSdkTruenum = &reqCtx.ReqParams.Param.UnSupportSdkTrueNum
	req.Ua = &reqCtx.ReqParams.Param.UserAgent
	req.OsVersionCodeV2 = &reqCtx.ReqParams.Param.OSVersionCode
	IfSupportSeperateCreative := int32(2)
	// 是否走素材三期逻辑
	if reqCtx.ReqParams.Param.NewCreativeFlag {
		IfSupportSeperateCreative = int32(3)
	}
	req.IfSupportSeperateCreative = &IfSupportSeperateCreative
	req.VideoVersion = &reqCtx.ReqParams.Param.VideoVersion

	req.Videoh = &reqCtx.ReqParams.Param.VideoH
	req.Videow = &reqCtx.ReqParams.Param.VideoW
	req.RankerInfo = &reqCtx.ReqParams.Param.RankerInfo
	req.DebugMode = reqCtx.ReqParams.Param.DebugMode
	req.ApiVersion = &reqCtx.ReqParams.Param.ApiVersionCode
	rst := ad_server.InteractiveResourceType(reqCtx.ReqParams.Param.IARst)
	req.ResourceType = &rst
	req.LowDevice = reqCtx.ReqParams.Param.LowDevice

	// bidRequest
	if reqCtx.ReqParams.IsBidRequest {
		req.BidFloor = &reqCtx.ReqParams.Param.BidFloor
	}

	trafficTag := "adnet-nothb"
	if bid == mvconst.MAdx {
		req.BidFloor = &reqCtx.ReqParams.Param.AdxBidFloor
		trafficTag = "adx-nothb"
	}
	if reqCtx.ReqParams.IsHBRequest {
		extra := reqCtx.ReqParams.Param.Algorithm + "-adserver"
		reqCtx.ReqParams.Param.Algorithm = extra
		reqCtx.ReqParams.Param.Extra = extra
		trafficTag = "aladdin-hb"
	}
	req.CampaignKind = &trafficTag
	// as 传递包名，只召回对应包名的单子
	req.TargetPackageNameSet = renderPackageList(reqCtx.ReqParams.Param.PlPkg)
	req.DEPRECATEDExcludePackageNameSet = reqCtx.ReqParams.Param.ExcludePackageNames

	// more offer
	if reqCtx.ReqParams.Param.Mof == 1 {
		// 固定unit和拆分unit都需要传parent_unit
		if reqCtx.ReqParams.Param.UcParentUnitId > 0 {
			req.ParentUnitId = &reqCtx.ReqParams.Param.UcParentUnitId
		} else {
			req.ParentUnitId = &reqCtx.ReqParams.Param.ParentUnitId
		}
		req.MofData = &reqCtx.ReqParams.Param.MofData
	}
	req.SystemId = &reqCtx.ReqParams.Param.SysId
	req.SysbkupId = &reqCtx.ReqParams.Param.BkupId
	req.Ruid = &reqCtx.ReqParams.Param.RuId
	// algo experiment
	if reqCtx.ReqParams.Param.PubFlowExpectPrice > 0.0 {
		req.PubFlowExpectPrice = &reqCtx.ReqParams.Param.PubFlowExpectPrice
	}
	if reqCtx.ReqParams.Param.FillEcpmFloor > 0.0 {
		req.FillEcpmFloor = &reqCtx.ReqParams.Param.FillEcpmFloor
	}
	req.FillEcpmFloor = &reqCtx.ReqParams.Param.FillEcpmFloor
	req.TestMode = reqCtx.ReqParams.Param.AsTestMode

	// 传递设备id md5值
	req.ImeiMd5 = &reqCtx.ReqParams.Param.ImeiMd5
	req.DevIdMd5 = &reqCtx.ReqParams.Param.AndroidIDMd5
	req.GaidMd5 = &reqCtx.ReqParams.Param.GAIDMd5
	req.IdfaMd5 = &reqCtx.ReqParams.Param.IDFAMd5
	req.IfSupportDco = &reqCtx.ReqParams.Param.IfSupDco

	req.IfSupportRewardPlus = &reqCtx.ReqParams.Param.RwPlus

	if reqCtx.ReqParams.Param.Skadnetwork != nil {
		if demand.IsSupMtgSKAdnetwork(reqCtx.ReqParams.AppInfo.RealPackageName, reqCtx.ReqParams.Param.Skadnetwork.Adnetids, reqCtx.ReqParams.Param.Skadnetwork.Ver) {
			supportSkadn := true
			req.IfSupportSkadn = &supportSkadn
		}
	}

	// SDK版本 是否支持deeplink —— AS也需要传这个字段
	param := reqCtx.ReqParams.Param
	if onlineAPISupport, _ := onlineApiSupportDeepLink(reqCtx); (param.Platform == mvconst.PlatformIOS && param.FormatSDKVersion.SDKVersionCode >= 40700) ||
		(param.Platform == mvconst.PlatformAndroid && param.FormatSDKVersion.SDKVersionCode >= 90300) ||
		onlineAPISupport {
		isDeeplink := true
		req.IfSupportDeepLink = &isDeeplink
	}

	// 传递大模板标记
	if reqCtx.ReqParams.Param.BigTemplateFlag {
		IfSupportBigTemplate := int32(1)
		req.IfSupportBigTemplate = &IfSupportBigTemplate
	}
	// 传递placement id给as，rs
	req.PlacementId = &reqCtx.ReqParams.Param.FinalPlacementId
	if reqCtx.ReqParams.Param.RandNum > 0 {
		req.RandNum = reqCtx.ReqParams.Param.RandNum
	}
	req.PassthroughData = RenderPassthroughData(&reqCtx.ReqParams.Param)

	// 判断是否切量polaris
	if reqCtx.ReqParams.Param.PolarisFlag {
		supportPolaris := int32(1)
		req.IfSupportPolaris = &supportPolaris
	}
	req.WebEnv = RenderWebEnv(&reqCtx.ReqParams.Param)

	// 改为as记录需要as传给rs的素材信息的单子
	if len(reqCtx.ReqParams.Param.NeedCreativeDataCIds) > 0 {
		req.IdCheck = RenderIdCheck(&reqCtx.ReqParams.Param)
	}

	req.DebugCampaignIdList = renderExcludeIDslice(reqCtx.ReqParams.Param.TargetIds)
	req.DspJunoRes = &reqCtx.ReqParams.Param.DspMoreOfferInfo

	if mvutil.Config.CommonConfig.LogConfig.OutputFullReqRes || reqCtx.ReqParams.Param.Debug > 0 || reqCtx.ReqParams.Param.DebugMode {
		rStr, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(req)
		if reqCtx.ReqParams.Param.Debug > 0 {
			reqCtx.ReqParams.DebugInfo += "request_ad_server_data@@@" + string(rStr) + "<\br>"
		} else if reqCtx.ReqParams.Param.DebugMode {
			reqCtx.DebugModeInfo = append(reqCtx.DebugModeInfo, req)
		} else {
			mvutil.Logger.Runtime.Debugf("request_id=[%s] ad_server request data=%s",
				reqCtx.ReqParams.Param.RequestID, rStr)
		}
		reqCtx.ReqParams.Param.AsDebugParam = &req
	}
	return &req, nil
}

func FillMoreAsAds(reqCtx *mvutil.ReqCtx, ads *corsair_proto.BackendAds, r *ad_server.QueryResult_) (int, error) {
	return doFillAdServerAds(reqCtx, ads, r, true)
}

func FillAsAds(reqCtx *mvutil.ReqCtx, backendCtx *mvutil.BackendCtx, r *ad_server.QueryResult_) (int, error) {
	return doFillAdServerAds(reqCtx, backendCtx.Ads, r, false)
}

func doFillAdServerAds(reqCtx *mvutil.ReqCtx, ads *corsair_proto.BackendAds, r *ad_server.QueryResult_, skip bool) (int, error) {
	if reqCtx == nil || r == nil || ads == nil {
		return 0, errors.New("fillAdServerAds params invalidate")
	}
	// backendCtx.Ads.RequestKey = backendCtx.AdReqKeyName
	ads.Strategy = r.Strategy
	if r.IsSetRunTimeVariables() {
		var m corsair_proto.RunTimeVariable
		m.GetCampaignDetailTime = r.GetRunTimeVariables().GetCampaignDetailTime
		m.GetCampaignIdTime = r.GetRunTimeVariables().GetCampaignIdTime
		m.GetCampaignInfoTime = r.GetRunTimeVariables().GetCampaignInfoTime
		m.GetCampaignPositionInfoTime = r.GetRunTimeVariables().GetCampaignPositionInfoTime
		m.NumRecalled = r.GetRunTimeVariables().NumRecalled
		m.PublisherId = r.GetRunTimeVariables().PublisherId
		m.RankTime = r.GetRunTimeVariables().RankTime
		m.RecallSize = r.GetRunTimeVariables().RecallSize
		m.RecallSizeAll = r.GetRunTimeVariables().RecallSizeAll
		ads.RunTimeVariables = &m
	}
	ads.AdTemplate = r.AdTemplate
	if r.IsSetNoneResultReason() {
		// var m corsair_proto.NoneResultReason
		m := corsair_proto.NoneResultReason(int64(r.GetNoneResultReason()))
		ads.NoneResultReason = &m
	}
	for _, reason := range r.FilterReason {
		var m corsair_proto.FilterReason
		m.CampaignId = reason.GetCampaignId()
		m.Reason = reason.GetReason()
		ads.FilterReason = append(ads.FilterReason, &m)
	}
	ads.AlgoFeatInfo = r.AlgoFeatInfo
	ads.IfLowerImp = r.IfLowerImp
	ads.ResourceType = r.ResourceType
	ads.EndScreenTemplateId = r.EndScreenTemplateId
	ads.ExtAdxAlgo = r.GetExtAdxAlgo()
	ads.EcpmFloor = r.GetEcpmFloor()
	ads.DspId = mvconst.FakeAdserverDsp
	if r.BigTemplateInfo != nil {
		var bigTempalte corsair_proto.BigTemplate
		bigTempalte.BigTemplateId = r.BigTemplateInfo.BigTemplateId
		bigTempalte.SlotIndexCampaignIdMap = r.BigTemplateInfo.SlotIndexCampaignIdMap
		ads.BigTemplateInfo = &bigTempalte
	}
	bidsPrice := make([]string, 0, len(r.CampaignList))
	index := 0

	for _, rad := range r.CampaignList {
		var ad corsair_proto.Campaign
		if rad == nil {
			continue
		}
		var bidPrice float64
		if rad.GetBidPrice() > 0 {
			bidPrice = rad.GetBidPrice()
		} else {
			bidPrice = r.GetBidPrice()
		}
		bidsPrice = append(bidsPrice, strconv.FormatInt(rad.GetCampaignId(), 10)+":"+strconv.FormatFloat(bidPrice, 'f', 2, 64))

		fillCampaign(&ad, rad)

		index++
		if index == 1 && skip {
			continue
		}
		ads.CampaignList = append(ads.CampaignList, &ad)
	}
	reqCtx.ReqParams.Param.BidsPrice = strings.Join(bidsPrice, ",")

	asABTestRes, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(r.GetUsedAbKeys())
	reqCtx.ReqParams.Param.AsABTestResTag = string(asABTestRes)

	reqCtx.ReqParams.Param.JunoCommonLogInfoJson = r.GetJunoCommonLogInfoJson()

	return ERR_OK, nil
}

func fillCampaign(ad *corsair_proto.Campaign, rad *ad_server.Campaign) {
	if ad == nil || rad == nil {
		return
	}

	ad.CampaignId = strconv.FormatInt(rad.CampaignId, 10)

	ad.AdSource = rad.AdSource
	ad.AdTemplate = rad.AdTemplate
	ad.ImageSizeId = rad.ImageSizeId
	ad.OfferType = rad.OfferType
	ad.BtType = rad.BtType
	ad.CreativeId = rad.CreativeId
	ad.CreativeTypeIdMap = rad.CreativeTypeIdMap
	ad.CreativeId2 = rad.CreativeId2
	ad.CreativeTypeIdMap2 = rad.CreativeTypeIdMap2

	// 封装CreativeDataList
	renderCreativeDataList(ad, rad.CreativeDataList)

	ad.DynamicCreative = rad.DynamicCreative
	ad.AdElementTemplate = rad.AdElementTemplate

	// playable adserver控制
	ad.Playable = rad.Playable
	ad.EndcardUrl = rad.EndcardUrl
	ad.ExtPlayable = rad.ExtPlayable
	ad.VideoEndTypeAs = rad.VideoEndType
	ad.Orientation = rad.Orientation
	ad.UsageVideo = rad.UsageVideo
	ad.TemplateGroup = rad.TemplateGroup
	ad.VideoTemplateId = rad.VideoTemplateId
	ad.EndCardTemplateId = rad.EndCardTemplateId
	ad.MiniCardTemplateId = rad.MiniCardTemplateId
	ad.BidPrice = rad.BidPrice
	ad.UseAlgoPrice = rad.UseAlgoPrice
	ad.AlgoPriceIn = rad.AlgoPriceIn
	ad.AlgoPriceOut = rad.AlgoPriceOut
	ad.AsABTestResTag = rad.AsAbTestTag
}

func renderCreativeDataList(ad *corsair_proto.Campaign, creativeData []*ad_server.CreativeData) {
	for _, v := range creativeData {
		var crData corsair_proto.CreativeData
		crData.DocId = v.GetDocId()
		for _, val := range v.CreativeTypeIdList {
			var crTypeIdList corsair_proto.CreativeTypeId
			crTypeIdList.CreativeId = val.GetCreativeId()
			crTypeIdList.Type = val.GetType()
			crTypeIdList.MId = val.MId
			crTypeIdList.CpdId = val.CpdId
			crTypeIdList.OmId = val.OmId
			crData.CreativeTypeIdList = append(crData.CreativeTypeIdList, &crTypeIdList)
		}
		ad.CreativeDataList = append(ad.CreativeDataList, &crData)
	}
}

func RenderPassthroughData(params *mvutil.Params) map[string]string {
	res := map[string]string{}
	paramConf := extractor.GetPassthroughData()
	if len(paramConf) == 0 {
		return res
	}
	value := reflect.ValueOf(params).Elem()
	for _, paramName := range paramConf {
		res[paramName] = mvutil.GetParamValInString(value.FieldByName(paramName))
	}
	return res
}

func RenderWebEnv(params *mvutil.Params) map[string]string {
	res := map[string]string{}
	if params.WebEnvData.Webgl != nil {
		res["webgl"] = strconv.Itoa(*params.WebEnvData.Webgl)
	}
	return res
}

func RenderIdCheck(params *mvutil.Params) map[int64]bool {
	return renderExcludeIDs(params.NeedCreativeDataCIds)
}

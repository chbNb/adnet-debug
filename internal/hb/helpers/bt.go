package helpers

import (
	"hash/crc32"
	"sort"
	"strconv"

	"gitlab.mobvista.com/ADN/adnet/internal/extractor"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/constant"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/params"
	smodel "gitlab.mobvista.com/ADN/structs/model"
)

const MtgIdDefault int64 = 1700785270

func BlendTraffic(ad *params.Ad, campaign *smodel.CampaignInfo, params *params.LoadReqData) int64 {
	appId := params.AppId
	if ad == nil || campaign == nil || campaign.Ctype == 0 {
		return appId
	}

	if !InArray(int(campaign.Ctype), constant.InstallPayTypes) {
		return appId
	}
	// btV3
	return blendTrafficV3(campaign, params)
}

func blendTrafficV3(campaign *smodel.CampaignInfo, params *params.LoadReqData) int64 {
	appId := params.AppId
	btClass := getBTClass(params)
	// 获取offer维度BT规则
	if campaign.BtV4 == nil || btClass == 0 {
		return appId
	}
	btV3 := *(campaign.BtV4)
	// 记录ext_btclass
	params.Extbtclass = btClass
	// 替换subId
	// 兜底
	subIds := btV3.SubIds
	if len(subIds) <= 0 {
		apps, ifFind := extractor.GetReplenishApp()
		if ifFind && len(apps) > 0 {
			return RandInt64Arr(apps)
		}
		return appId
	}
	// rand
	randMap := make(map[int]int)
	for k, v := range subIds {
		if v.Rate <= 0 {
			continue
		}
		kInt, _ := strconv.Atoi(k)
		randMap[kInt] = v.Rate
	}
	if len(randMap) <= 0 {
		apps, ifFind := extractor.GetReplenishApp()
		if ifFind && len(apps) > 0 {
			return RandInt64Arr(apps)
		}
		return appId
	}
	randSubId := RandByRate(randMap)
	if randSubId == 0 {
		return appId
	}
	randSubIdStr := strconv.Itoa(randSubId)
	subInfo := subIds[randSubIdStr]
	// 赋值包名
	params.ExtfinalPackageName = subInfo.PackageName
	// dsp
	if len(subInfo.DspSubIds) > 0 {
		dspRandMap := make(map[int]int)
		for k, v := range subInfo.DspSubIds {
			if v.Rate <= 0 {
				continue
			}
			kInt, err := strconv.Atoi(k)
			if err != nil {
				continue
			}
			dspRandMap[kInt] = v.Rate
		}
		if len(dspRandMap) > 0 {
			dspRandKey := RandByRate(dspRandMap)
			dspRandKeyStr := strconv.Itoa(dspRandKey)
			params.Extfinalsubid = int64(dspRandKey)
			params.ExtfinalPackageName = subInfo.DspSubIds[dspRandKeyStr].PackageName
		}
	}
	if randSubId <= 0 {
		return appId
	} else {
		return int64(randSubId)
	}
}

func getBTClass(params *params.LoadReqData) int {
	publisherType := params.PublisherType
	if publisherType == constant.PublisherTypeM {
		return params.UnitBtClass
	} else {
		return params.AppBtClass
	}
}

// 广告请求类型渠道信息透明化
func RenderRequestPackage(campaign *smodel.CampaignInfo, params *params.LoadReqData) {
	renderRequestParams(params)
	// 替换newsubid
	renderNewSubId(campaign, params)
	// 是否传递包名
	needPackageName := needPackageName(campaign, params)
	if !needPackageName {
		params.ExtfinalPackageName = ""
	}
	// mtgId
	params.ExtMtgId = extractor.GetAppPackageMtgID(params.ExtfinalPackageName)
	if params.ExtMtgId == MtgIdDefault {
		params.ExtfinalPackageName = ""
	}
}

// 整理广告请求类型参数
func renderRequestParams(params *params.LoadReqData) {
	if params.Extfinalsubid <= int64(0) {
		if params.Extra14 > int64(0) {
			params.Extfinalsubid = params.Extra14
		} else {
			params.Extfinalsubid = params.AppId
		}
	}
}

// 替换newsubid
func renderNewSubId(campaign *smodel.CampaignInfo, params *params.LoadReqData) {
	conf := getNewSubIdConf(campaign, params.Extfinalsubid, params)
	llen := len(conf)
	if llen <= 0 {
		return
	}
	// 如果是对1
	if llen == 1 {
		for k, v := range conf {
			kInt, _ := strconv.ParseInt(k, 10, 64)
			params.Extfinalsubid = kInt
			params.ExtfinalPackageName = v
		}
		return
	}
	// 如果是对多
	subId := randNewSubId(conf, params.Extfinalsubid)
	subIdStr := strconv.FormatInt(subId, 10)
	if subId > int64(0) {
		params.Extfinalsubid = subId
		params.ExtfinalPackageName = conf[subIdStr]
	}
}

// 一个subid只能替换成对应的newsubid
func randNewSubId(subArr map[string]string, finalSubId int64) int64 {
	llen := len(subArr)
	subIdStr := strconv.FormatInt(finalSubId, 10)
	rand := int(crc32.ChecksumIEEE([]byte(subIdStr))) % llen
	// 按appid排序
	subIdList := make([]int, 0, llen)
	for k := range subArr {
		// ////////////todo
		kInt, _ := strconv.Atoi(k)
		subIdList = append(subIdList, kInt)
	}
	sort.Ints(subIdList)
	return int64(subIdList[rand])
	// i := 0
	// for k, _ := range subArr {
	// 	if i == rand {
	// 		kInt, _ := strconv.ParseInt(k, 10, 64)
	// 		return kInt
	// 	}
	// 	i++
	// }
	// return int64(0)
}

// 获取newsubid配置
func getNewSubIdConf(campaign *smodel.CampaignInfo, subIdInt64 int64, params *params.LoadReqData) map[string]string {
	subId := strconv.FormatInt(subIdInt64, 10)
	// campaign维度
	confs := campaign.BlackSubidListV2
	if len(confs[subId]) > 0 {
		return confs[subId]
	}

	cf := getAddSubIdConf(params, confs)
	if cf != nil {
		return cf
	}

	// advertiser维度
	res := make(map[string]string)
	config, _ := extractor.GetAdvBlackSubIdList()
	if campaign.AdvertiserId == 0 {
		return res
	}

	advId := strconv.FormatInt(int64(campaign.AdvertiserId), 10)
	if len(config[advId]) <= 0 {
		return res
	}
	confs = config[advId]
	if len(confs[subId]) > 0 {
		return confs[subId]
	}

	cf = getAddSubIdConf(params, confs)
	if cf != nil {
		return cf
	}

	return res
}

func getAddSubIdConf(params *params.LoadReqData, confs map[string]map[string]string) map[string]string {
	// M traffic 维度
	if params.PublisherType == constant.PublisherTypeM && len(confs["-1"]) > 0 {
		return confs["-1"]
	}
	// DSP traffic 维度
	if params.PublisherId == constant.DspPublisherID && len(confs["-2"]) > 0 {
		return confs["-2"]
	}
	// 默认设置
	if len(confs["0"]) > 0 {
		return confs["0"]
	}
	return nil
}

func needPackageName(campaign *smodel.CampaignInfo, params *params.LoadReqData) bool {
	appId := params.AppId
	if params.PublisherId == constant.DspPublisherID {
		appId = params.ExtdspRealAppid
	}
	appStr := strconv.FormatInt(appId, 10)
	var conf *smodel.AppPostList
	// campaign维度
	if campaign.AppPostList != nil {
		conf = campaign.AppPostList
	} else {
		if campaign.AdvertiserId == 0 {
			return false
		}
		// advertiser维度
		confs, ifFind := extractor.GetAppPostList()
		if !ifFind {
			return false
		}

		conf = confs[campaign.AdvertiserId]
	}

	if conf == nil {
		return false
	}

	if len(conf.Exclude) > 0 && InStrArray(appStr, conf.Exclude) {
		return false
	}
	if len(conf.Include) > 0 {
		if InStrArray(appStr, conf.Include) || InStrArray("ALL", conf.Include) {
			return true
		}
	}
	return false
}

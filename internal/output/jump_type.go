package output

import (
	supply_mvutil "gitlab.mobvista.com/ADN/chasm/module/supply/mvutil"
	"strconv"
	"strings"

	"gitlab.mobvista.com/ADN/adnet/internal/extractor"
	"gitlab.mobvista.com/ADN/adnet/internal/mvconst"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	smodel "gitlab.mobvista.com/ADN/structs/model"
)

//func RandJumpType(r *mvutil.RequestParams, campaign *smodel.CampaignInfo, params *mvutil.Params) string {
//	// 如果是onlineapi，则返回0
//	if params.RequestType == mvconst.REQUEST_TYPE_OPENAPI_AD {
//		// 若是可以做click_mode=6的开发者，且offer，publisher，app，unit维度可做6，且传的gaid或idfa不为空，则返回clickmode=6
//		res := oldOnlineClickMode(r, campaign, params)
//		if res == mvconst.JUMP_TYPE_CLIENT_DO_ALL {
//			return res
//		}
//		// 对于online api，使用pingmode模式则不能切到clickmode=13，clickmode都为0
//		if params.PingMode == 1 {
//			return mvconst.JUMP_TYPE_NORMAL
//		}
//
//		// 指定三方空设备只能走click mode 0
//		onlineEmptyDeviceNoServerJump, _ := extractor.GetOnlineEmptyDeviceNoServerJump()
//		if onlineEmptyDeviceNoServerJump.Status &&
//			mvutil.IsDevidEmpty(params) {
//			if mvutil.InStrArray(campaign.ThirdParty, onlineEmptyDeviceNoServerJump.ThirdParty) {
//				return mvconst.JUMP_TYPE_NORMAL
//			}
//
//			if len(onlineEmptyDeviceNoServerJump.BlackThirdParty) > 0 &&
//				!mvutil.InStrArray(campaign.ThirdParty, onlineEmptyDeviceNoServerJump.BlackThirdParty) {
//				return mvconst.JUMP_TYPE_NORMAL
//			}
//		}
//
//		// online & dsp clickmode=13切量
//		conf := getJumpTypeConfV2(*r, *campaign)
//		res = getJumpTypeByConf(conf)
//
//		emptyIPUA, _ := extractor.GetOnlineEmptyDeviceIPUAABTest()
//		// 临时写死做实验
//		if emptyIPUA.Status && mvutil.IsDevidEmpty(params) && campaign.ThirdParty == "Adjust" {
//			if res == mvconst.JUMP_TYPE_ONLINE_DSP_AJUMP_SERVER {
//				// 服务端跳转
//				params.ExtDataInit.AdjustS2S = "1"
//			} else if res == mvconst.JUMP_TYPE_NORMAL {
//				// 客户端跳转
//				params.ExtDataInit.AdjustS2S = "0"
//			}
//		}
//		return res
//	}
//	// 如果是jm_icon，返回5
//	if params.AdType == mvconst.ADTypeJMIcon {
//		return mvconst.JUMP_TYPE_CLIENT_SEND_DEVID
//	}
//
//	conf := getJumpTypeArr(r, campaign, params)
//	rateSum := 0
//	randMap := make(map[int]int)
//	for jumpType, rate := range conf {
//		jumpTypeInt, err := strconv.Atoi(jumpType)
//		if err != nil {
//			continue
//		}
//		rateSum = rateSum + int(rate)
//		randMap[jumpTypeInt] = int(rate)
//	}
//	if rateSum <= 0 || len(randMap) <= 0 {
//		return mvconst.JUMP_TYPE_NORMAL
//	}
//	resInt := mvutil.RandByRate(randMap)
//	return strconv.Itoa(resInt)
//}

// 获取jumpType配置
func getJumpTypeArr(r *mvutil.RequestParams, campaign *smodel.CampaignInfo, params *mvutil.Params) map[string]int32 {
	// 分别获取Android ios配置
	if r.Param.Platform == mvconst.PlatformAndroid {
		return getJumpTypeAndroid(r, campaign, params)
	} else if r.Param.Platform == mvconst.PlatformIOS {
		return getJumpTypeIOS(r, campaign, params)
	}
	return map[string]int32{mvconst.JUMP_TYPE_NORMAL: int32(100)}
}

func getJumpTypeConf(r *mvutil.RequestParams, campaign *smodel.CampaignInfo) map[string]int32 {
	var conf map[string]int32
	// 优先获取campaign的配置
	if len(campaign.JumpTypeConfig) > 0 {
		conf = campaign.JumpTypeConfig
	}
	// unit维度
	if len(conf) <= 0 {
		conf = r.UnitInfo.JumptypeConfig
	}
	// app维度
	if len(conf) <= 0 {
		conf = r.AppInfo.JumptypeConfig
	}
	// publisher维度
	if len(conf) <= 0 {
		conf = r.PublisherInfo.JumptypeConfig
	}
	//三方维度
	if len(conf) <= 0 {
		confTP, _ := extractor.GetJUMP_TYPE_CONFIG_THIRD_PARTY()
		platformStr := strconv.Itoa(r.Param.Platform)
		confs, ok := confTP[platformStr]
		thirdParty := strings.ToLower(campaign.ThirdParty)
		if ok && len(confs) > 0 && len(thirdParty) > 0 {
			tpConf, ok := confs[thirdParty]
			if ok {
				conf = tpConf
			}
		}
	}
	// adv维度
	if len(conf) <= 0 {
		confAdv, _ := extractor.GetJUMP_TYPE_CONFIG_ADV()
		platformStr := strconv.Itoa(r.Param.Platform)
		confs, ok := confAdv[platformStr]
		if ok && len(confs) > 0 && campaign.AdvertiserId != 0 {
			advStr := strconv.FormatInt(int64(campaign.AdvertiserId), 10)
			confTmp, ok := confs[advStr]
			if ok {
				conf = confTmp
			}
		}
	}
	// 兜底
	if len(conf) <= 0 {
		platform := r.Param.Platform
		var ifFind bool
		if platform == mvconst.PlatformAndroid {
			conf, ifFind = extractor.GetJUMP_TYPE_CONFIG()
			if !ifFind {
				conf = map[string]int32{"5": int32(1)}
			}
			return conf
		}
		if platform == mvconst.PlatformIOS {
			conf, ifFind = extractor.GetJUMP_TYPE_CONFIG_IOS()
			if !ifFind {
				conf = map[string]int32{"6": int32(1)}
			}
			return conf
		}
	}
	return conf
}

func getJumpTypeAndroid(r *mvutil.RequestParams, campaign *smodel.CampaignInfo, params *mvutil.Params) map[string]int32 {
	conf := getJumpTypeConf(r, campaign)
	res := make(map[string]int32)

	// 当做clickmode=0方式的跳转的  17-10-14号
	// 1）campaignType 不为1且为2 即不为App Store 且GooglePlay
	// 或者
	// 2） 没有包名的offer
	if isClickMode0(campaign, params) {
		res[mvconst.JUMP_TYPE_NORMAL] = getJumpTypeVal(conf, mvconst.JUMP_TYPE_NORMAL)
		return res
	}
	// 支持clickmode=1跳转，判断逻辑改为满足如下两个条件
	// 针对android平台 只判断gaid;
	// 支持clickmode=4跳转，返回的clickurl中：直接使用offer的包名，拼接market地址
	// clickmode=4的判断逻辑改为满足如下两个条件
	// 针对android平台  判断逻辑改为canVTA的方法 && 只判断gaid;
	if canClickMode1(campaign, params) {
		res[mvconst.JUMP_TYPE_SERVER_JUMP] = getJumpTypeVal(conf, mvconst.JUMP_TYPE_SERVER_JUMP)
		res[mvconst.JUMP_TYPE_SDK_TO_MARKET] = getJumpTypeVal(conf, mvconst.JUMP_TYPE_SDK_TO_MARKET)
	}
	// 不再支持clickmode=3跳转

	// 支持clickmode=5跳转
	// 针对ios，idfa为空值也走00000000-0000-0000-0000-000000000000和idfa情况下的逻辑，也能走5,6,11,12。
	adnConf, _ := extractor.GetADNET_SWITCHS()
	// 默认关闭，若开启，则gaid为空值的情况则走0
	gaidEmptyCanJump := true
	if isCloseGaid, ok := adnConf["icGaid"]; ok && isCloseGaid == 1 {
		gaidEmptyCanJump = false
	}
	// sdk版本限制+gaid
	if compareSDKVersionForJumpType(params, mvconst.JUMP_TYPE_CLIENT_SEND_DEVID) && (len(params.GAID) > 0 || gaidEmptyCanJump) ||
		params.RequestPath == mvconst.PATHJssdkApi {
		res[mvconst.JUMP_TYPE_CLIENT_SEND_DEVID] = getJumpTypeVal(conf, mvconst.JUMP_TYPE_CLIENT_SEND_DEVID)
		if len(campaign.DirectUrl) != 0 {
			res[mvconst.JUMP_TYPE_CLIENT_SEND_DEVID_PING_SERVER] = getJumpTypeVal(conf, mvconst.JUMP_TYPE_CLIENT_SEND_DEVID_PING_SERVER)
		} else {
			// 针对mtg tracking无法走clickmodel=11且存在device id的情况下，兜底方案为clickmodel=5
			res[mvconst.JUMP_TYPE_CLIENT_SEND_DEVID] = res[mvconst.JUMP_TYPE_CLIENT_SEND_DEVID] + getJumpTypeVal(conf, mvconst.JUMP_TYPE_CLIENT_SEND_DEVID_PING_SERVER)
		}
	}
	// 支持clickmode=6跳转
	// sdk版本限制+gaid
	if compareSDKVersionForJumpType(params, mvconst.JUMP_TYPE_CLIENT_DO_ALL) && (len(params.GAID) > 0 || gaidEmptyCanJump) ||
		params.RequestPath == mvconst.PATHJssdkApi {
		res[mvconst.JUMP_TYPE_CLIENT_DO_ALL] = getJumpTypeVal(conf, mvconst.JUMP_TYPE_CLIENT_DO_ALL)
		if len(campaign.DirectUrl) != 0 {
			res[mvconst.JUMP_TYPE_CLIENT_DO_ALL_PING_SERVER] = getJumpTypeVal(conf, mvconst.JUMP_TYPE_CLIENT_DO_ALL_PING_SERVER)
		} else {
			// 针对mtg tracking无法走clickmodel=12且存在device id的情况下，兜底方案为clickmodel=6
			res[mvconst.JUMP_TYPE_CLIENT_DO_ALL] = res[mvconst.JUMP_TYPE_CLIENT_DO_ALL] + getJumpTypeVal(conf, mvconst.JUMP_TYPE_CLIENT_DO_ALL_PING_SERVER)
		}
	}
	// 暂不支持clickmode = 8跳转方式（目前还不确定sdk的版本等过滤逻辑）
	// 正常跳转
	res[mvconst.JUMP_TYPE_NORMAL] = getJumpTypeVal(conf, mvconst.JUMP_TYPE_NORMAL)
	return res
}

func getJumpTypeIOS(r *mvutil.RequestParams, campaign *smodel.CampaignInfo, params *mvutil.Params) map[string]int32 {
	conf := getJumpTypeConf(r, campaign)
	res := make(map[string]int32)
	// 当做clickmode=0方式的跳转的  10-14号
	// 1）campaignType 不为1且为2 即不为App Store 且GooglePlay
	// 或者
	// 2） 没有包名的offer
	if isClickMode0(campaign, params) {
		res[mvconst.JUMP_TYPE_NORMAL] = getJumpTypeVal(conf, mvconst.JUMP_TYPE_NORMAL)
		return res
	}
	// 支持clickmode=1跳转，判断逻辑改为满足如下两个条件
	// 针对ios平台 判断idfa（IOS本次增加支持clickmode=1的跳转）
	// 支持clickmode=4跳转，返回的clickurl中：直接使用offer的包名，拼接market地址
	// clickmode=4的判断逻辑改为满足如下两个条件
	// ios平台 判断逻辑改为canVTA的方法 && 判断idfa（所有的click_mode的逻辑中，如果idfa为00000000-0000-0000-0000-000000000000的情况，认为是有设置idfa）
	if canClickMode1(campaign, params) {
		res[mvconst.JUMP_TYPE_SERVER_JUMP] = getJumpTypeVal(conf, mvconst.JUMP_TYPE_SERVER_JUMP)
		res[mvconst.JUMP_TYPE_SDK_TO_MARKET] = getJumpTypeVal(conf, mvconst.JUMP_TYPE_SDK_TO_MARKET)
	}
	// 支持clickmode=5跳转
	// 针对ios，idfa为空值也走00000000-0000-0000-0000-000000000000和idfa情况下的逻辑，也能走5,6,11,12。
	adnConf, _ := extractor.GetADNET_SWITCHS()
	// 默认关闭，若开启，则idfa为空值的情况则走0
	idfaEmptyCanJump := true
	if isCloseIdfa, ok := adnConf["icIdfa"]; ok && isCloseIdfa == 1 {
		idfaEmptyCanJump = false
	}
	// sdk版本限制+idfa
	if compareSDKVersionForJumpType(params, mvconst.JUMP_TYPE_CLIENT_SEND_DEVID) && (len(params.IDFA) > 0 || idfaEmptyCanJump) ||
		params.RequestPath == mvconst.PATHJssdkApi {
		res[mvconst.JUMP_TYPE_CLIENT_SEND_DEVID] = getJumpTypeVal(conf, mvconst.JUMP_TYPE_CLIENT_SEND_DEVID)
		if len(campaign.DirectUrl) != 0 {
			res[mvconst.JUMP_TYPE_CLIENT_SEND_DEVID_PING_SERVER] = getJumpTypeVal(conf, mvconst.JUMP_TYPE_CLIENT_SEND_DEVID_PING_SERVER)
		} else {
			// 针对mtg tracking无法走clickmodel=11且存在device id的情况下，兜底方案为clickmodel=5
			res[mvconst.JUMP_TYPE_CLIENT_SEND_DEVID] = res[mvconst.JUMP_TYPE_CLIENT_SEND_DEVID] + getJumpTypeVal(conf, mvconst.JUMP_TYPE_CLIENT_SEND_DEVID_PING_SERVER)
		}
	}
	// 支持clickmode=6跳转
	// sdk版本限制+idfa
	if compareSDKVersionForJumpType(params, mvconst.JUMP_TYPE_CLIENT_DO_ALL) && (len(params.IDFA) > 0 || idfaEmptyCanJump) ||
		params.RequestPath == mvconst.PATHJssdkApi {
		res[mvconst.JUMP_TYPE_CLIENT_DO_ALL] = getJumpTypeVal(conf, mvconst.JUMP_TYPE_CLIENT_DO_ALL)
		if len(campaign.DirectUrl) != 0 {
			res[mvconst.JUMP_TYPE_CLIENT_DO_ALL_PING_SERVER] = getJumpTypeVal(conf, mvconst.JUMP_TYPE_CLIENT_DO_ALL_PING_SERVER)
		} else {
			// 针对mtg tracking无法走clickmodel=12且存在device id的情况下，兜底方案为clickmodel=6
			res[mvconst.JUMP_TYPE_CLIENT_DO_ALL] = res[mvconst.JUMP_TYPE_CLIENT_DO_ALL] + getJumpTypeVal(conf, mvconst.JUMP_TYPE_CLIENT_DO_ALL_PING_SERVER)
		}
	}
	// 暂不支持clickmode = 8跳转方式（目前还不确定sdk的版本等过滤逻辑）
	// 正常跳转
	res[mvconst.JUMP_TYPE_NORMAL] = getJumpTypeVal(conf, mvconst.JUMP_TYPE_NORMAL)
	return res
}

func getJumpTypeVal(conf map[string]int32, jumpType string) int32 {
	val, ok := conf[jumpType]
	if ok {
		return val
	}
	return int32(0)
}

func isClickMode0(campaign *smodel.CampaignInfo, params *mvutil.Params) bool {
	if campaign.PackageName == "" {
		return true
	}
	// if params.LinkType != 1 && params.LinkType != 2 {
	// 	return true
	// }
	return false
}

func canClickMode1(campaign *smodel.CampaignInfo, params *mvutil.Params) bool {
	platform := params.Platform
	if platform == mvconst.PlatformAndroid && canVTA(campaign) && len(params.GAID) > 0 {
		return true
	}
	if platform == mvconst.PlatformIOS && canVTA(campaign) && len(params.IDFA) > 0 {
		return true
	}
	return false
}

func canVTA(campaign *smodel.CampaignInfo) bool {
	if campaign.JumpType == 0 {
		return false
	}

	jumpType := campaign.JumpType
	if jumpType == int32(1) && len(campaign.DirectPackageName) > 0 {
		return true
	}

	if jumpType == int32(3) || jumpType == int32(7) {
		return true
	}
	return false
}

func compareSDKVersionForJumpType(params *mvutil.Params, jumpType string) bool {
	confs, ifFind := extractor.GetJUMPTYPE_SDKVERSION()
	if !ifFind {
		return false
	}
	conf, ok := confs[jumpType]
	if !ok {
		return false
	}
	sdkversionItem := supply_mvutil.RenderSDKVersion(params.SDKVersion)
	if len(sdkversionItem.SDKType) <= 0 {
		return false
	}
	confVersion, ok := conf[sdkversionItem.SDKType]
	if !ok {
		return false
	}
	confVersionCode := mvutil.GetVersionCode(confVersion)
	// if sdkversionItem.SDKVersionCode >= confVersionCode {
	// 	return true
	// }
	// return false
	return sdkversionItem.SDKVersionCode >= confVersionCode
}

func oldOnlineClickMode(r *mvutil.RequestParams, campaign *smodel.CampaignInfo, params *mvutil.Params) string {
	conf := getJumpTypeConf(r, campaign)
	CM6Publishers, _ := extractor.GetCAN_CLICK_MODE_SIX_PUBLISHER()
	if _, ok := conf[mvconst.JUMP_TYPE_CLIENT_DO_ALL]; mvutil.InInt64Arr(r.Param.PublisherID, CM6Publishers) && ok && !mvutil.IsDevidEmpty(params) {
		return mvconst.JUMP_TYPE_CLIENT_DO_ALL
	}
	return mvconst.JUMP_TYPE_NORMAL
}

func getJumpTypeByConf(conf map[string]int32) string {
	rateSum := 0
	randMap := make(map[int]int)
	for jumpType, rate := range conf {
		jumpTypeInt, err := strconv.Atoi(jumpType)
		if err != nil {
			continue
		}
		rateSum = rateSum + int(rate)
		randMap[jumpTypeInt] = int(rate)
	}
	if rateSum <= 0 || len(randMap) <= 0 {
		return mvconst.JUMP_TYPE_NORMAL
	}
	resInt := mvutil.RandByRate(randMap)
	return strconv.Itoa(resInt)
}

func getJumpTypeConfV2(r mvutil.RequestParams, campaign smodel.CampaignInfo) map[string]int32 {
	var conf map[string]int32
	// 优先获取campaign的配置
	if campaign.JumpTypeConfig2 != nil && len(campaign.JumpTypeConfig2) > 0 {
		conf = campaign.JumpTypeConfig2
	}
	// unit维度
	if len(conf) <= 0 {
		conf = r.UnitInfo.JumptypeConfigV2
	}
	// app维度
	if len(conf) <= 0 {
		conf = r.AppInfo.JumptypeConfigV2
	}
	// publisher维度
	if len(conf) <= 0 {
		conf = r.PublisherInfo.JumptypeConfigV2
	}
	//三方维度
	if len(conf) <= 0 {
		confTP, _ := extractor.GetJUMP_TYPE_CONFIG_THIRD_PARTY()
		platformStr := "androidV2"
		if r.Param.Platform == mvconst.PlatformIOS {
			platformStr = "iosV2"
		}
		confs, ok := confTP[platformStr]
		thirdParty := strings.ToLower(campaign.ThirdParty)
		if ok && len(confs) > 0 && len(thirdParty) > 0 {
			tpConf, ok := confs[thirdParty]
			if ok {
				conf = tpConf
			}
		}
	}
	// adv维度
	if len(conf) <= 0 {
		confAdv, _ := extractor.GetJUMP_TYPE_CONFIG_ADV()
		platformStr := "androidV2"
		if r.Param.Platform == mvconst.PlatformIOS {
			platformStr = "iosV2"
		}
		confs, ok := confAdv[platformStr]
		if ok && len(confs) > 0 && campaign.AdvertiserId != 0 {
			advStr := strconv.FormatInt(int64(campaign.AdvertiserId), 10)
			confTmp, ok := confs[advStr]
			if ok {
				conf = confTmp
			}
		}
	}
	// 兜底
	if len(conf) <= 0 {
		// 兜底为0
		conf = map[string]int32{"0": int32(1)}
	}
	return conf
}

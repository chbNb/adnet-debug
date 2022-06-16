package process_pipeline

import (
	"errors"
	"math/rand"
	"strings"

	"gitlab.mobvista.com/ADN/adnet/internal/backend"
	"gitlab.mobvista.com/ADN/adnet/internal/extractor"
	"gitlab.mobvista.com/ADN/adnet/internal/mvconst"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	"gitlab.mobvista.com/ADN/adnet/internal/watcher"
	"gitlab.mobvista.com/adserver/recommend_protocols/go/ad_server"
)

const AdxFactor = 1.00

type TrafficAllotFilter struct {
}

func (taf *TrafficAllotFilter) Process(data interface{}) (interface{}, error) {
	in, ok := data.(*mvutil.ReqCtx)
	if !ok {
		return nil, errors.New("TrafficAllotFilter input type should be *mvutil.ReqCtx")
	}
	in.FlowTagID = mvconst.FlowTagDefault
	isProxyVer := isProxyVersion(in.ReqParams)
	orderedbackends := make([]int, 0)
	addMobvista := true
	// M系统的流量才走流量分配，或者说抽样模块
	if !isProxyVer {
		// M系统的默认流量flow_tag_id=1
		in.FlowTagID = mvconst.FlowTagDefaultM

		doSample := true
		// 新增JSSDK/appwall/offerwall不切量
		if in.ReqParams.Param.Scenario == mvconst.PATHJssdkApi ||
			in.ReqParams.Param.AdType == mvconst.ADTypeAppwall ||
			in.ReqParams.Param.AdType == mvconst.ADTypeOfferWall ||
			(in.ReqParams.Param.Scenario == mvconst.SCENARIO_OPENAPI && in.ReqParams.Param.IsLowFlowUnitReq) ||
			(in.ReqParams.Param.RequestType == mvconst.REQUEST_TYPE_OPENAPI_V3 &&
				backend.FilterAdTracking(in.ReqParams.Param.FormatAdType,
					in.ReqParams.Param.Platform,
					mvconst.AdTracking_Both_Click_Imp,
					in.ReqParams.Param.FormatSDKVersion)) {
			doSample = false
		}
		// myoffer优先
		for _, adSource := range renderAdSource(in.ReqParams) {
			if mvconst.ADSourceMYOffer == int(adSource) {
				// mvutil.Logger.Runtime.Infof("request_id=[%s] is myoffer priority",
				// 	in.ReqParams.param.RequestID)
				doSample = false
				break
			}
		}
		if in.ReqParams.IsTopon {
			doSample = true
		}

		if doSample {
			watcher.AddWatchValue("sample", float64(1))
			// sample
			if distributionData, hit, withAsInfo := sample(in); hit {
				addMobvista = false
				in.FlowTagID = distributionData.FlowTagId
				orderedbackends = distributionData.AdBackends
				// native必须 appid或者 unitid 与sdkversion都在白名单才支持 adchoice，其他 adtype 都支持
				in.ReqParams.Param.SupportAdChoice = true
				if mvutil.InArray(mvconst.MAdx, distributionData.AdBackends) && in.ReqParams.Param.AdType == mvconst.ADTypeNative {
					in.ReqParams.Param.SupportAdChoice = false
					adChoiceConfig := extractor.GetAdChoiceConfigData()
					if adChoiceConfig != nil {
						if (in.ReqParams.Param.Platform == mvconst.PlatformAndroid && adChoiceConfig.AndroidSDKVersions != nil &&
							!filterSDKVersion(adChoiceConfig.AndroidSDKVersions, in.ReqParams.Param.FormatSDKVersion.SDKVersionCode)) ||
							(in.ReqParams.Param.Platform == mvconst.PlatformIOS && adChoiceConfig.IOSSDKVersions != nil &&
								!filterSDKVersion(adChoiceConfig.IOSSDKVersions, in.ReqParams.Param.FormatSDKVersion.SDKVersionCode)) {
							if mvutil.InInt64Arr(in.ReqParams.Param.UnitID, adChoiceConfig.UnitIds) || mvutil.InInt64Arr(in.ReqParams.Param.AppID, adChoiceConfig.AppIds) {
								in.ReqParams.Param.SupportAdChoice = true
							}
						}
					}
				}
				// moat tag
				if mvutil.InArray(mvconst.MAdx, distributionData.AdBackends) {
					in.ReqParams.Param.SupportMoattag = false
					moattagConfig := extractor.GetMoattagConfig()
					if moattagConfig != nil {
						if (in.ReqParams.Param.Platform == mvconst.PlatformAndroid && moattagConfig.AndroidSDKVersions != nil &&
							!filterSDKVersion(moattagConfig.AndroidSDKVersions, in.ReqParams.Param.FormatSDKVersion.SDKVersionCode)) ||
							(in.ReqParams.Param.Platform == mvconst.PlatformIOS && moattagConfig.IOSSDKVersions != nil &&
								!filterSDKVersion(moattagConfig.IOSSDKVersions, in.ReqParams.Param.FormatSDKVersion.SDKVersionCode)) {
							in.ReqParams.Param.SupportMoattag = true
						}
					}
				}

				if mvutil.InArray(mvconst.MAdx, distributionData.AdBackends) {
					// adBackendConfigInfo, ifFind := extractor.GetAdBackendConfigInfo(int64(mvconst.MAdx))

					// if ifFind {
					adxBackendCtx := mvutil.NewBackendCtx("", "", "", AdxFactor,
						1, []string{"ALL"})
					if in.ReqParams.IsTopon { // topon只打给mas
						adxBackendCtx.IsBidMAS = true
						adxBackendCtx.IsBidAdServer = false
					} else {
						adxBackendCtx.IsBidAdServer = true
						adxMedia, ifFind := extractor.GetAdxTrafficMediaConfig(in.ReqParams.Param.UnitID, in.ReqParams.Param.FormatAdType, in.ReqParams.Param.CountryCode)
						if ifFind {
							for _, dspId := range adxMedia.DspWhiteList {
								if dspId == mvconst.FakeAdserverDsp {
									adxBackendCtx.IsBidAdServer = true
								}
								if dspId == mvconst.MAS {
									adxBackendCtx.IsBidMAS = true
								}
							}
						} else if withAsInfo { //找不到配置才需要使用返回值
							adxBackendCtx.IsBidMAS = true
							adxBackendCtx.IsBidAdServer = true
						}
					}
					// 默认都可以切给 asDsp
					in.Backends[mvconst.MAdx] = adxBackendCtx
				}
				// }
			} else {
				watcher.AddWatchValue("sample_0", float64(1))
			}
		}
	}

	// more_offer，appwall 切量到pioneer
	// mp 流量也切量pioneer
	if mvutil.IsAppwallOrMoreOffer(in.ReqParams.Param.AdType) {
		if in.ReqParams.Param.ExtDataInit.MoreofferAndAppwallMvToPioneerTag == mvconst.REQUEST_PIONEER {
			in.Backends[mvconst.Pioneer] = mvutil.NewMobvistaCtx()
			orderedbackends = append(orderedbackends, mvconst.Pioneer)
		} else {
			if _, en := in.Backends[mvconst.Mobvista]; !en {
				in.Backends[mvconst.Mobvista] = mvutil.NewMobvistaCtx()
				orderedbackends = append(orderedbackends, mvconst.Mobvista)
			}
		}
	}

	// sdk banner不做as兜底。more_offer和appwall也不需要兜底，在上面已经进行了兜底
	if !mvutil.IsBannerOrSplashOrNativeH5(in.ReqParams.Param.AdType) && !mvutil.IsAppwallOrMoreOffer(in.ReqParams.Param.AdType) && addMobvista {
		if _, en := in.Backends[mvconst.Mobvista]; !en {
			in.Backends[mvconst.Mobvista] = mvutil.NewMobvistaCtx()
			orderedbackends = append(orderedbackends, mvconst.Mobvista)
		}
	}

	in.ReqParams.FlowTagID = in.FlowTagID
	in.OrderedBackends = orderedbackends
	if rand.Intn(100) <= mvutil.Config.AreaConfig.HttpConfig.RuntimeLogRate {
		mvutil.Logger.Runtime.Infof("request_id=[%s] is hit sample flaw_tag_id=[%d], ip=[%s], countryCode=[%s], backends %v, isProxyVersion=[%t], ASTest=[%t]",
			in.ReqParams.Param.RequestID, in.FlowTagID, in.ReqParams.Param.ClientIP,
			in.ReqParams.Param.CountryCode, in.OrderedBackends, isProxyVer, in.AdsTest)
	}

	return in, nil
}

func renderAdSource(r *mvutil.RequestParams) []ad_server.ADSource {
	var sourceList []ad_server.ADSource
	sourceList = append(sourceList, ad_server.ADSource(r.Param.AdSourceID))
	if r.Param.PublisherID != int64(5488) || r.Param.AdSourceID != mvconst.ADSourceMYOffer {
		return sourceList
	}
	sList := r.UnitInfo.GetAdSourceIDList(r.Param.CountryCode)
	for _, v := range sList {
		if r.Param.AdSourceID == v {
			return []ad_server.ADSource{ad_server.ADSource_APIOFFER, ad_server.ADSource_MYOFFER}
		}
	}
	return sourceList
}

func isProxyVersion(paramList *mvutil.RequestParams) bool {
	sdkVersion := paramList.Param.SDKVersion
	prefix := "empty"
	strVer := ""
	if paramList.IsHBRequest {
		return false
	}
	if paramList.Param.RequestType == int(ad_server.RequestType_ONLINE_API) {
		return false
	}
	// 针对jssdk，不做过滤
	if paramList.Param.RequestPath == mvconst.PATHJssdkApi {
		return false
	}
	if paramList.Param.RequestType != int(ad_server.RequestType_OPENAPI_V3) {
		return true
	}
	if paramList.Param.PublisherType != mvconst.PublisherTypeM {
		return true
	}
	if paramList.IsTopon {
		return false
	}
	if len(sdkVersion) > 0 {
		if strings.Contains(sdkVersion, "_") {
			sdkVerList := strings.Split(sdkVersion, "_")
			prefix = strings.ToLower(strings.TrimSpace(sdkVerList[0]))
			strVer = strings.TrimSpace(sdkVerList[1])
		} else {
			strVer = strings.TrimSpace(sdkVersion)
		}
		intVer, err := mvutil.IntVer(strVer)
		if err != nil {
			mvutil.Logger.Runtime.Warnf("request_id=[%s] isProxyVersion() Parse Version error:%s",
				paramList.Param.RequestID, err.Error())
			return true
		}
		if configIntVer, ok := mvutil.SDKVersions[prefix]; ok {
			if configIntVer <= intVer {
				return false
			}
		}
	}
	return true
}

func sample(reqCtx *mvutil.ReqCtx) (distributionInfo *mvutil.DistributionData, ifFind bool, withAsInfo bool) {

	if reqCtx.ReqParams.Param.IsTowardAdx { //debug调试, 通过参数强制调用ADX，生成勿用
		var distributionInfo mvutil.DistributionData
		distributionInfo.FlowTagId = mvconst.FlowTagMedia
		distributionInfo.AdBackends = []int{mvconst.MAdx}
		distributionInfo.AdReqKeys = map[int]string{
			17: "ADX",
		}
		return &distributionInfo, true, true
	}

	// sdkbanner默认切给adx
	// IOS 小于5.8.6的不切
	if reqCtx.ReqParams.IsHBRequest {
		var distributionInfo mvutil.DistributionData
		distributionInfo.FlowTagId = 101 // TODO refactor hb and adnet
		distributionInfo.AdBackends = []int{mvconst.MAdx}
		distributionInfo.AdReqKeys = map[int]string{
			17: "ADX",
		}
		return &distributionInfo, true, false
	}

	// Topon 只打给madx
	if reqCtx.ReqParams.IsTopon {
		var distributionInfo mvutil.DistributionData
		distributionInfo.FlowTagId = 101
		distributionInfo.AdBackends = []int{mvconst.MAdx}
		distributionInfo.AdReqKeys = map[int]string{
			17: "ADX",
		}
		return &distributionInfo, true, false
	}

	if mvutil.IsBannerOrSplashOrNativeH5(reqCtx.ReqParams.Param.AdType) {
		var distributionInfo mvutil.DistributionData
		distributionInfo.FlowTagId = mvconst.FlowTagSdkBanner
		distributionInfo.AdBackends = []int{mvconst.MAdx}
		distributionInfo.AdReqKeys = map[int]string{
			17: "ADX",
		}
		return &distributionInfo, true, false
	}

	// 只要有unit 或者 ad_type 维度配置在 adx_traffic_media_config则切给 adx
	adxCfg, mediaHit := extractor.GetAdxTrafficMediaConfig(reqCtx.ReqParams.Param.UnitID, reqCtx.ReqParams.Param.FormatAdType, reqCtx.ReqParams.Param.CountryCode)
	if mediaHit && (adxCfg.UnitId > 1 || (adxCfg.AdType > 0 && reqCtx.ReqParams.Param.RequestType == mvconst.REQUEST_TYPE_OPENAPI_V3)) {
		var distributionInfo mvutil.DistributionData
		distributionInfo.FlowTagId = mvconst.FlowTagMedia
		distributionInfo.AdBackends = []int{mvconst.MAdx}
		distributionInfo.AdReqKeys = map[int]string{
			17: "ADX",
		}
		return &distributionInfo, true, false
	}

	//onlineApi按切量切到ADX
	param := &reqCtx.ReqParams.Param
	if param.RequestType == mvconst.REQUEST_TYPE_OPENAPI_AD && rand.Intn(1000) < getOlineApiUseAdxRate(
		param.PublisherID, param.AppID, param.UnitID, int64(param.AdType), param.DebugDspId) {
		mvutil.Logger.Runtime.Infof("onlineApi use adx, requestID：%s, pid: %d, aid:%d, uid: %d",
			param.RequestID, param.PublisherID, param.AppID, param.UnitID)
		var distributionInfo mvutil.DistributionData
		distributionInfo.FlowTagId = mvconst.FlowTagMedia
		distributionInfo.AdBackends = []int{mvconst.MAdx}
		distributionInfo.AdReqKeys = map[int]string{
			17: "ADX",
		}
		return &distributionInfo, true, true
	}

	return nil, false, false
}

func getOlineApiUseAdxRate(publisherId, appId, unitId, adType int64, debugDspId int) int {
	conf := extractor.GetONLINE_API_USE_ADX()

	// debug
	if debugDspId == 6 || debugDspId == 13 { //6，13 一定走ADX
		return 1000
	} else if debugDspId == -1 { //一定不走ADX
		return 0
	}

	if subConf, ok := conf["unit"]; ok && unitId > 0 {
		if r, ok := subConf[unitId]; ok {
			return r
		}
	}
	if subConf, ok := conf["app"]; ok && appId > 0 {
		if r, ok := subConf[appId]; ok {
			return r
		}
	}
	if subConf, ok := conf["publisher"]; ok && publisherId > 0 {
		if r, ok := subConf[publisherId]; ok {
			return r
		}
	}
	if subConf, ok := conf["adtype"]; ok && adType > 0 {
		if r, ok := subConf[adType]; ok {
			return r
		}
	}

	return 0
}

// func anyExists(data map[string]bool, keys ...string) bool {
// 	if len(data) == 0 {
// 		return false
// 	}
// 	for _, k := range keys {
// 		if len(k) == 0 {
// 			continue
// 		}
// 		if _, ok := data[k]; ok {
// 			return true
// 		}
// 	}
// 	return false
// }

func filterSDKVersion(versionData *mvutil.DataSDKVersion, versionCode int32) bool {
	// 有一个配置排掉了，就丢掉
	for _, item := range versionData.Exclude {
		if item.Min <= versionCode && item.Max >= versionCode {
			return true
		}
	}
	// 所有配置都没有排掉，没有配置白名单, 就使用
	if len(versionData.Include) == 0 {
		return false
	}
	// 所有配置都没有排掉，有一个要的就使用
	for _, item := range versionData.Include {
		if item.Min <= versionCode && item.Max >= versionCode {
			return false
		}
	}

	// 所有配置没有排掉，也没有一个要的，就不要了
	return true
}

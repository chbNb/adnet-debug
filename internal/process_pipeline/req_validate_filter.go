package process_pipeline

import (
	"errors"
	"fmt"
	"gitlab.mobvista.com/ADN/chasm/module/demand"
	"math/rand"
	"strconv"
	"strings"

	jsoniter "github.com/json-iterator/go"
	"gitlab.mobvista.com/ADN/adnet/internal/errorcode"
	"gitlab.mobvista.com/ADN/adnet/internal/extractor"
	"gitlab.mobvista.com/ADN/adnet/internal/mvconst"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	"gitlab.mobvista.com/ADN/exporter/metrics"
)

type ReqValidateFilter struct {
}

type NativeInfo struct {
	Id    int `json:"id"`
	AdNum int `json:"ad_num"`
}

// 参数校验
func (rvf *ReqValidateFilter) Process(data interface{}) (interface{}, error) {
	in, ok := data.(*mvutil.RequestParams)
	if !ok {
		return nil, errors.New("ReqValidateFilter input type should be *params.RequestParams")
	}

	//filter-openSource
	if in.Param.ApiVersion >= mvconst.API_VERSION_2_1 && in.Param.Open == 1 {
		if in.PublisherInfo.Publisher.OpenSource != 1 {
			return nil, errorcode.EXCEPTION_OP_FORBIDDEN
		}
		openSourceSign := mvutil.Md5(in.Param.Sign + in.Param.TS)
		if len(in.Param.SignWithTimeStamp) == 0 || openSourceSign != in.Param.SignWithTimeStamp {
			return nil, errorcode.ExCEPTION_OP_SIGN_CHECK_ERROR
		}
	}

	if in.Param.RequestPath != mvconst.PATHOpenApiV2 && in.Param.RequestPath != mvconst.PATHREQUEST {
		if in.UnitInfo.Unit.Status != 1 {
			mvutil.Logger.Runtime.Warnf("request_id=[%s] has filter by Unit entry is not active", in.Param.RequestID)
			return nil, errorcode.EXCEPTION_UNIT_NOT_ACTIVE
		}

		if in.UnitInfo.AppId != in.Param.AppID {
			mvutil.Logger.Runtime.Warnf("request_id=[%s] has filter by appID is validate", in.Param.RequestID)
			return nil, errorcode.EXCEPTION_UNIT_NOT_FOUND_IN_APP
		}
	}
	// check adSourceID
	adSourceID := in.UnitInfo.GetAdSourceID(in.Param.CountryCode, in.Param.AdSourceID)
	if adSourceID == 0 {
		mvutil.Logger.Runtime.Warnf("request_id=[%s] has filter by adsourceid is validate. ad_source_id=[%d] checkresult=[%d]", in.Param.RequestID, in.Param.AdSourceID, adSourceID)
		return nil, errorcode.EXCEPTION_SERVICE_REQUEST_AD_SOURCE_CLOSED
	}
	in.Param.AdSourceID = adSourceID

	if in.Param.Platform != mvconst.PlatformAndroid && in.Param.Platform != mvconst.PlatformIOS {
		mvutil.Logger.Runtime.Warnf("request_id=[%s] has filter by platform is validate. platform=[%d]", in.Param.RequestID, in.Param.Platform)
		return nil, errorcode.EXCEPTION_APP_PLATFORM_ERROR
	}

	if in.Param.RequestPath == mvconst.PATHJssdkApi {
		// 对js_native_video(id=285)，js_banner_video(id=286)不再进行横竖屏校验
		if !mvutil.IsJsVideo(in.Param.AdType) && !mvutil.InArray(in.Param.Orientation, []int{mvconst.ORIENTATION_PORTRAIT, mvconst.ORIENTATION_LANDSCAPE}) {
			mvutil.Logger.Runtime.Warnf("request_id=[%s] has filter by orientation is validate", in.Param.RequestID)
			return nil, errorcode.EXCEPTION_PARAMS_ERROR
		}
		// 主域名校验
		if in.AppInfo.App.DomainVerify == 1 && in.AppInfo.App.Domain != in.Param.MainDomain {
			mvutil.Logger.Runtime.Warnf("request_id=[%s] has filter by domain is validate", in.Param.RequestID)
			return nil, errorcode.EXCEPTION_DOMAIN_ERROR
		}
		// 针对appwall,wx小程序adtype不做networktype校验
		if !mvutil.IsAppwallOrMoreOffer(in.Param.AdType) && !mvutil.IsWxAdType(in.Param.AdType) {
			// network校验
			var RecallNetArr []int
			for _, RecallNet := range strings.Split(in.UnitInfo.Unit.RecallNet, ";") {
				RecallNet, err := strconv.Atoi(RecallNet)
				if err == nil {
					RecallNetArr = append(RecallNetArr, RecallNet)
				}
			}
			if !mvutil.InArray(in.Param.NetworkType, RecallNetArr) {
				mvutil.Logger.Runtime.Warnf("request_id=[%s] has filter by recallnet is validate", in.Param.RequestID)
				return nil, errorcode.EXCEPTION_NETWORK_ERROR
			}
		}
	}
	if in.Param.RequestPath == mvconst.PATHJssdkApi {
		// 若为jssdk则使用setting的值做校验（app维度的platform值均为3）
		if in.UnitInfo.Setting.AdPlatform != 0 && in.Param.Platform != in.UnitInfo.Setting.AdPlatform {
			mvutil.Logger.Runtime.Warnf("request_id=[%s] has filter by platform is validate", in.Param.RequestID)
			return nil, errorcode.EXCEPTION_PARAMS_ERROR
		}
	} else if in.AppInfo.App.Platform != in.Param.Platform {
		mvutil.Logger.Runtime.Warnf("request_id=[%s] has filter by platform is validate", in.Param.RequestID)
		return nil, errorcode.EXCEPTION_APP_PLATFORM_ERROR
	}

	// category校验
	if in.Param.Category > 9999 {
		mvutil.Logger.Runtime.Warnf("request_id=[%s] has filter by category is validate", in.Param.RequestID)
		return nil, errorcode.EXCEPTION_CATEGORY_ERROR
	}

	// 填充deviceType
	if in.Param.Platform == mvconst.PlatformAndroid {
		in.Param.DeviceType = mvconst.DevicePhone
		in.Param.ReqDeviceType = 1
	} else {
		in.Param.DeviceType = mvconst.DeviceUnknown
		if len(in.Param.RModel) > 0 && strings.HasPrefix(in.Param.RModel, "iPhone") {
			in.Param.DeviceType = mvconst.DevicePhone
			in.Param.ReqDeviceType = 1
		}
		if len(in.Param.RModel) > 0 && strings.HasPrefix(in.Param.RModel, "iPad") {
			in.Param.DeviceType = mvconst.DeviceTablet
			in.Param.ReqDeviceType = 2
		}
	}

	if needCheckSign(in) {
		rawSign := fmt.Sprintf("%d%s", in.Param.AppID, in.PublisherInfo.Publisher.Apikey)
		if in.Param.RequestPath == mvconst.PATHREQUEST {
			rawSign = fmt.Sprintf("%d%s%s%s", in.Param.AppID, in.PublisherInfo.Publisher.Apikey, in.Param.Scenario, in.Param.TimeStamp)
		}
		sign := mvutil.Md5(rawSign)
		if sign != in.Param.Sign {
			mvutil.Logger.Runtime.Warnf("request_id=[%s] has filter by checkSign is validate appID=[%d],apikey=[%s],sign=[%s], param.sign=[%s]",
				in.Param.RequestID, in.Param.AppID, in.PublisherInfo.Publisher.Apikey, sign, in.Param.Sign)
			return nil, errorcode.EXCEPTION_SIGN_ERROR
		}
	}

	if in.Param.Debug > 0 {
		// 初始化，写入
		var json = jsoniter.ConfigCompatibleWithStandardLibrary
		paramStr, _ := json.Marshal(*in)
		in.DebugInfo += "after_params_invalidate@@@" + string(paramStr) + "<\br>"
	}

	// request接口判断是否为sa ban掉的请求
	if in.Param.RequestPath == mvconst.PATHREQUEST {
		BlackConfs, _ := extractor.GetREQUEST_BLACKLIST()
		if BlackConfs.AppIds != nil && mvutil.InInt64Arr(in.Param.AppID, *BlackConfs.AppIds) {
			return nil, errorcode.EXCEPTION_RETURN_EMPTY
		}
		if BlackConfs.DeviceModels != nil && mvutil.InStrArray(in.Param.Model, *BlackConfs.DeviceModels) {
			return nil, errorcode.EXCEPTION_RETURN_EMPTY
		}
		if BlackConfs.Countries != nil && mvutil.InStrArray(in.Param.CountryCode, *BlackConfs.Countries) {
			return nil, errorcode.EXCEPTION_RETURN_EMPTY
		}
	}

	// check system
	if checkSystemErr(in) {
		mvutil.Logger.Runtime.Warnf("request_id=[%s] has filter by system is validate", in.Param.RequestID)
		return nil, errorcode.EXCEPTION_404_NOT_FOUND
	}

	// check deviceid
	if mvutil.GetRandConsiderZero(in.Param.GAID, in.Param.IDFA, mvconst.SALT_CHECK_DEVID, 100) == -1 {
		if !canRequestWithoutDevId(in) {
			if in.Param.Platform == mvconst.PlatformAndroid {
				mvutil.Logger.Runtime.Warnf("request_id=[%s] has filter by deviceId is validate", in.Param.RequestID)
				return nil, errorcode.EXCEPTION_GAID_EMPTY
				// return mvconst.EXCEPTION_GAID_EMPTY, fmt.Errorf("EXCEPTION_GAID_EMPTY")
			} else if in.Param.Platform == mvconst.PlatformIOS {
				mvutil.Logger.Runtime.Warnf("request_id=[%s] has filter by deviceId is validate", in.Param.RequestID)
				return nil, errorcode.EXCEPTION_IDFA_EMPTY
				// return mvconst.EXCEPTION_IDFA_EMPTY, fmt.Errorf("EXCEPTION_IDFA_EMPTY")
			}
		}
	}
	if (in.Param.AdType == mvconst.ADTypeAppwall || in.Param.AdType == mvconst.ADTypeMoreOffer) && in.Param.Platform == mvconst.PlatformAndroid &&
		in.Param.FormatSDKVersion.SDKVersionCode < mvconst.MoreOfferBlock {
		blockMoreOfferConf, _ := extractor.GetBLOCK_MORE_OFFER_CONFIG()
		if blockMoreOfferConf.Status && mvutil.InInt64Arr(in.Param.UnitID, blockMoreOfferConf.UnitIds) {
			mvutil.Logger.Runtime.Warnf("unit_id=[%d], sdk_version=[%s] more offer filter by sdk_version", in.Param.UnitID, in.Param.SDKVersion)
			return nil, errorcode.EXCEPTION_RETURN_EMPTY
		}
	}
	// 过滤有问题的流量，不做召回，也不记录请求日志。
	if filterBadRequest(in) {
		in.Param.Extreject = errorcode.EXCEPTION_TC_FILTER_BY_BAD_REQUEST.String()
		mvutil.StatLossRequestLog(&in.Param)
		return nil, errorcode.EXCEPTION_RETURN_EMPTY
	}
	// 针对sdk banner进行没传unit_size的过滤
	if in.Param.AdType == mvconst.ADTypeSdkBanner && (in.Param.SdkBannerUnitWidth == 0 || in.Param.SdkBannerUnitHeight == 0) {
		mvutil.Logger.Runtime.Warnf("unit_id=[%d] sdk banner filter by empty unit_size", in.Param.UnitID)
		return nil, errorcode.EXCEPTION_BANNER_UNIT_SIZE_EMPTY
	}

	// 过滤低于谷歌合规版本SDK的广告请求
	if filterIllegalGPSdkVersion(in) {
		in.Param.Extreject = errorcode.EXCEPTION_FILTER_BY_ILLEGAL_GP_SDK_VERSION.String()
		mvutil.StatLossRequestLog(&in.Param)
		return nil, errorcode.EXCEPTION_RETURN_EMPTY
	}

	// 过滤sdk传入的placement id和unit配置中不一致的请求
	if filterByPlacementId(in) {
		mvutil.Logger.Runtime.Warnf("unitId=[%d] has filter by placement_id inconsistent", in.Param.UnitID)
		return nil, errorcode.EXCEPTION_FILTER_BY_PLACEMENTID_INCONSISTENT
	}

	if in.Param.Scenario == mvconst.SCENARIO_OPENAPI && in.Param.IsLowFlowUnitReq {
		mvutil.StatRequestLog(&in.Param)
		return nil, errors.New("isLowFlowUnit unitid:" + strconv.FormatInt(in.Param.UnitID, 10))
	}

	// 是否命中了ip黑名单
	if demand.BlockByIPBlacklist(in.Param.ClientIP, extractor.DemandDao) {
		in.Param.FilterRequestReason = mvconst.FilterRequestByIpBlacklist
		mvutil.StatRequestLog(&in.Param)
		metrics.IncCounterWithLabelValues(26, in.Param.FilterRequestReason)
		return nil, errors.New("filter by ip blacklist. unit_id=" + strconv.FormatInt(in.Param.UnitID, 10) + ",ip=" + in.Param.ClientIP)
	}

	if FilterRequestByStack(in) {
		in.Param.FilterRequestReason = mvconst.FilterRequestByStack
		mvutil.StatRequestLog(&in.Param)
		metrics.IncCounterWithLabelValues(26, in.Param.FilterRequestReason)
		return nil, errors.New("filter by stack. unitid:" + strconv.FormatInt(in.Param.UnitID, 10))
	}

	// 剩余磁盘空间过滤
	if FilterRequestByResidualMemory(in) {
		in.Param.FilterRequestReason = mvconst.FilterRequestByResidualMemory
		mvutil.StatLossRequestLog(&in.Param)
		return nil, errors.New("filter by residualMemory. unitid:" + strconv.FormatInt(in.Param.UnitID, 10))
	}

	reqCtx := mvutil.NewReqCtx()
	reqCtx.ReqParams = in
	return reqCtx, nil
}

func FilterRequestByResidualMemory(r *mvutil.RequestParams) bool {
	if len(r.Param.ResidualMemory) == 0 {
		return false
	}
	if !mvutil.IsHbOrV3OrV5Request(r.Param.RequestPath) {
		return false
	}
	residualMemory, err := strconv.ParseFloat(r.Param.ResidualMemory, 64)
	if err != nil {
		return false
	}
	conf := extractor.GetResidualMemoryFilterConf()
	if conf == nil {
		return false
	}
	if len(conf.AdTypeList) > 0 && !mvutil.InInt32Arr(r.Param.AdType, conf.AdTypeList) {
		return false
	}
	if len(conf.PlatformList) > 0 && !mvutil.InArray(r.Param.Platform, conf.PlatformList) {
		return false
	}
	if residualMemory > conf.ResidualMemory {
		return false
	}
	if conf.Rate == 100 {
		return true
	}
	randVal := rand.Intn(100)
	if conf.Rate > randVal {
		return true
	}
	return false
}

func FilterRequestByStack(r *mvutil.RequestParams) bool {
	if len(r.Param.TrafficInfo) == 0 {
		return false
	}
	conf := extractor.GetFilterByStackConf()
	if conf == nil {
		return false
	}
	var rate int
	var ok bool
	// 安卓只会上报堆栈信息
	if r.Param.PlatformName == mvconst.PlatformNameAndroid {
		if len(r.Param.StackList) == 0 {
			return false
		}
		stackList := strings.Split(r.Param.StackList, "|")
		rate, ok = getFilterByStackConf(r, stackList, conf)
		if !ok {
			return false
		}
	} else {
		// 优先级：类名>协议>堆栈
		if len(r.Param.ClassNameList) > 0 {
			classNameList := strings.Split(r.Param.ClassNameList, "|")
			rate, ok = getFilterByStackConf(r, classNameList, conf)
		}

		if !ok && len(r.Param.ProtocolList) > 0 {
			protocolList := strings.Split(r.Param.ProtocolList, "|")
			rate, ok = getFilterByStackConf(r, protocolList, conf)
		}

		if !ok && len(r.Param.StackList) > 0 {
			stackList := strings.Split(r.Param.StackList, "|")
			rate, ok = getFilterByStackConf(r, stackList, conf)
		}
		if !ok {
			return false
		}
	}
	randVal := rand.Intn(100)
	if rate > randVal {
		return true
	}
	return false
}

func getFilterByStackConf(r *mvutil.RequestParams, stackList []string, confs []*mvutil.FilterByStackConf) (int, bool) {
	for _, v := range stackList {
		for _, stackConf := range confs {
			if stackConf == nil {
				continue
			}
			if !strings.Contains(v, stackConf.StackName) {
				continue
			}
			// app
			appRate, ok := stackConf.AppConf[strconv.FormatInt(r.Param.AppID, 10)]
			if ok {
				return appRate, ok
			}
			// pub
			pubRate, ok := stackConf.PubConf[strconv.FormatInt(r.Param.PublisherID, 10)]
			if ok {
				return pubRate, ok
			}
			// 整体切量
			if stackConf.TotalConf > 0 {
				return stackConf.TotalConf, true
			}
		}
	}
	return 0, false
}

func filterByPlacementId(r *mvutil.RequestParams) bool {
	// sdk有传placementid及unit维度有配置placementid才做校验
	if r.Param.PlacementId == 0 || r.UnitInfo.Unit.PlacementId == 0 {
		return false
	}
	if r.Param.PlacementId == r.UnitInfo.Unit.PlacementId {
		return false
	}
	return true
}

func filterBadRequest(r *mvutil.RequestParams) bool {
	badRequestFilterConf := extractor.GetBAD_REQUEST_FILTER_CONF()
	if len(badRequestFilterConf) == 0 {
		return false
	}
	for _, v := range badRequestFilterConf {
		if len(v.PubWList) > 0 && !mvutil.InInt64Arr(r.Param.PublisherID, v.PubWList) {
			continue
		}
		if len(v.AppWList) > 0 && !mvutil.InInt64Arr(r.Param.AppID, v.AppWList) {
			continue
		}
		if len(v.UnitWList) > 0 && !mvutil.InInt64Arr(r.Param.UnitID, v.UnitWList) {
			continue
		}
		if len(v.AdTypeWList) > 0 && !mvutil.InInt32Arr(r.Param.AdType, v.AdTypeWList) {
			continue
		}
		if len(v.CountryWList) > 0 && !mvutil.InStrArray(r.Param.CountryCode, v.CountryWList) {
			continue
		}
		if len(v.ApiVerWList) > 0 && !mvutil.InInt32Arr(r.Param.ApiVersionCode, v.ApiVerWList) {
			continue
		}
		if len(v.PlatformWList) > 0 && !mvutil.InArray(r.Param.Platform, v.PlatformWList) {
			continue
		}
		if len(v.DeviceTypeWList) > 0 && !mvutil.InArray(r.Param.DeviceType, v.DeviceTypeWList) {
			continue
		}
		if len(v.ModelWList) > 0 && !mvutil.InStrArray(r.Param.Model, v.ModelWList) {
			continue
		}
		if v.MinSdkVer > 0 && r.Param.FormatSDKVersion.SDKVersionCode < v.MinSdkVer {
			continue
		}
		if v.MaxSdkVer > 0 && r.Param.FormatSDKVersion.SDKVersionCode > v.MaxSdkVer {
			continue
		}
		if v.MinOsVer > 0 && r.Param.OSVersionCode < v.MinOsVer {
			continue
		}
		if v.MaxOsVer > 0 && r.Param.OSVersionCode > v.MaxOsVer {
			continue
		}

		if len(v.AppVersionCodeList) > 0 && !mvutil.InStrArray(r.Param.AppVersionCode, v.AppVersionCodeList) {
			continue
		}

		if len(v.AppVersionNameList) > 0 && !mvutil.InStrArray(r.Param.AppVersionName, v.AppVersionNameList) {
			continue
		}

		if len(v.HModelWList) > 0 && !mvutil.InStrArray(r.Param.HardwareModel, v.HModelWList) {
			continue
		}

		// 模糊匹配
		if len(v.HModelMatchingWList) > 0 {
			isContain := false
			for _, v := range v.HModelMatchingWList {
				if strings.Contains(r.Param.HardwareModel, v) {
					isContain = true
					break
				}
			}
			if !isContain {
				continue
			}
		}

		if len(v.PubWList) > 0 || len(v.AppWList) > 0 || len(v.UnitWList) > 0 || len(v.AdTypeWList) > 0 || len(v.CountryWList) > 0 ||
			len(v.ApiVerWList) > 0 || len(v.PlatformWList) > 0 || len(v.DeviceTypeWList) > 0 || v.MinSdkVer > 0 || v.MaxSdkVer > 0 ||
			v.MinOsVer > 0 || v.MaxOsVer > 0 || len(v.ModelWList) > 0 || len(v.AppVersionCodeList) > 0 || len(v.AppVersionNameList) > 0 ||
			len(v.HModelWList) > 0 || len(v.HModelMatchingWList) > 0 {
			return true
		}
	}
	return false
}

func needCheckSign(r *mvutil.RequestParams) bool {
	if r.Param.Sign == mvconst.NO_CHECK_SIGN || r.IsHBRequest || r.IsTopon {
		return false
	}
	noCheckAppList, ifFind := extractor.GetSIGN_NO_CHECK_APPS()
	if ifFind {
		return !mvutil.InInt64Arr(r.Param.AppID, noCheckAppList)
	}
	return true
}

func checkSystemErr(r *mvutil.RequestParams) bool {
	pSystem := r.PublisherInfo.Publisher.Type
	if pSystem == 0 {
		return false
	}
	system := extractor.GetSYSTEM()
	if system == mvconst.SERVER_SYSTEM_TEST {
		return false
	}
	var allowType []int
	switch system {
	case mvconst.SERVER_SYSTEM_M:
		allowType = []int{mvconst.PublisherTypeADN, mvconst.PublisherTypeMediabuy, mvconst.PublisherTypeM, mvconst.PublisherTypeDSP}
	case mvconst.SERVER_SYSTEM_SA:
		allowType = []int{mvconst.PublisherTypeSuperads}
	case mvconst.SERVER_SYSTEM_MOBPOWER:
		allowType = []int{mvconst.PublisherTypeMobpower, mvconst.PublisherTypeM}
	default:
		return false
	}
	return !mvutil.InArray(pSystem, allowType)
}

// 判断是否能无deviceid请求广告
func canRequestWithoutDevId(r *mvutil.RequestParams) bool {
	// 优先取unit维度的配置，若unit维度无配置则去app维度配置，默认为允许召回
	isCan := r.UnitInfo.Unit.DevIdAllowNull
	if isCan == 0 {
		isCan = r.AppInfo.App.DevIdAllowNull
	}
	if isCan == 2 {
		return false
	} else if isCan == 1 {
		return true
	}
	return true
}

func filterIllegalGPSdkVersion(r *mvutil.RequestParams) bool {
	// 限制sdk 安卓的流量
	if !mvutil.IsHbOrV3OrV5Request(r.Param.RequestPath) || r.Param.Platform != mvconst.PlatformAndroid {
		return false
	}
	// 是否开启过滤逻辑
	adnConf, _ := extractor.GetADNET_SWITCHS()
	if filterIllegalGPSdkVersion, ok := adnConf["filterIllegalGPSdkVersion"]; !ok || filterIllegalGPSdkVersion != 1 {
		return false
	}
	// app白名单配置
	adnConfList := extractor.GetADNET_CONF_LIST()
	if illegalGPSdkWhiteAppList, ok := adnConfList["illegalGPSdkWhiteAppList"]; ok && mvutil.InInt64Arr(r.Param.AppID, illegalGPSdkWhiteAppList) {
		return false
	}
	// 限制sdkversion
	if r.AppInfo.App.LivePlatform == 1 && (r.Param.FormatSDKVersion.SDKVersionCode < mvconst.FilterIllegalGPSdkVersion ||
		(r.Param.FormatSDKVersion.SDKVersionCode >= mvconst.FilterIllegalGPSdkVersionAndroidXMin &&
			r.Param.FormatSDKVersion.SDKVersionCode < mvconst.FilterIllegalGPSdkVersionAndroidXMax)) {
		return true
	}
	return false
}

func filterBidByDeviceGEOCountry(r *mvutil.RequestParams) (flag bool) {
	confs := extractor.GetHBCheckDeviceGEOCCConf().HBFilterDeviceGEOConfigs
	for _, conf := range confs {
		if selectorHBCheckDeviceGEOCC(r, conf) {
			flag = true
			break
		}
	}
	return
}

func selectorHBCheckDeviceGEOCC(r *mvutil.RequestParams, conf *mvutil.HBFilterDeviceGEOConfig) bool {
	if len(conf.MediationName) > 0 && !mvutil.InStrArray(strings.ToUpper(r.Param.MediationName), conf.MediationName) {
		return false
	}
	if len(conf.CountryCode) > 0 && !mvutil.InStrArray(r.Param.BidDeviceGEOCC, conf.CountryCode) {
		return false
	}
	if len(conf.PublisherID) > 0 && !mvutil.InInt64Arr(r.Param.PublisherID, conf.PublisherID) {
		return false
	}
	if len(conf.AppID) > 0 && !mvutil.InInt64Arr(r.Param.AppID, conf.AppID) {
		return false
	}
	return true
}

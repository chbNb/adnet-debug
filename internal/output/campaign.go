package output

import (
	"bytes"
	"errors"
	"fmt"
	"hash/crc32"
	"math"
	"math/rand"
	"net/url"
	"reflect"
	"sort"
	"strconv"
	"strings"

	supply_mvutil "gitlab.mobvista.com/ADN/chasm/module/supply/mvutil"

	jsoniter "github.com/json-iterator/go"
	"gitlab.mobvista.com/ADN/adnet/internal/watcher"
	"gitlab.mobvista.com/ADN/structs/constant"

	"gitlab.mobvista.com/ADN/adnet/internal/corsair_proto"
	"gitlab.mobvista.com/ADN/adnet/internal/extractor"
	"gitlab.mobvista.com/ADN/adnet/internal/mvconst"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	"gitlab.mobvista.com/ADN/adnet/internal/protobuf"
	"gitlab.mobvista.com/ADN/chasm/module/demand"
	smodel "gitlab.mobvista.com/ADN/structs/model"
	"gitlab.mobvista.com/adserver/recommend_protocols/go/ad_server"
	"gitlab.mobvista.com/voyager/clickmode/mclickmode"
)

const (
	APP_DESC_LEN = 75
)

func RenderCampaignWithCreative(r *mvutil.RequestParams, corsairCampaign corsair_proto.Campaign, campaign *smodel.CampaignInfo, cContent map[ad_server.CreativeType]*protobuf.Creative) (Ad, error) {
	// 声明
	var ad Ad
	params := r.Param
	// 默认价格
	if campaign.Price != 0 {
		params.PriceOut = campaign.Price
	}

	if campaign.OriPrice != 0 {
		params.PriceIn = campaign.OriPrice
	}
	// 素材压缩abtest
	creativeCompressABTestV2(&params, &corsairCampaign, cContent)

	cleanEmptyDeviceIdABTest(&params, campaign)
	// 渲染算法价格
	renderAlgoPrice(&params, &corsairCampaign)

	if len(cContent) > 0 {
		renderCreativePb(&ad, &params, corsairCampaign, cContent, campaign)
	}

	// 分渠道价格处理 PriceOut
	getChannelPriceOut(&params, campaign)
	// BT
	renderBT(*r, &params, corsairCampaign.BtType, campaign)
	// 分渠道价格处理 PriceIn
	getChannelPriceIn(&params, campaign)
	// 封装ad
	renderCampaignInfo(&ad, r, campaign, &params, corsairCampaign)
	// 封装mpad
	if mvutil.IsMP(r.Param.RequestPath) {
		RenderMPInfo(&ad)
	}

	// 无idfa实验，针对af，903，ios流量
	// 逻辑迁移到clickmode分配逻辑内。
	// canUseAfWhiteList(&params, campaign, &ad)
	// jssdkWhiteList(&params, campaign, &ad)
	// 若走了af白名单通道实验，则不走adtracking和h5点击逻辑。
	if params.ExtDataInit.ClickInServer != 1 && params.ExtDataInit.ClickInServer != 3 && params.Extra10 != mvconst.JUMP_TYPE_TRACKING_PING_THIRDPARTY_CLICK_URL {
		// 判断是否走click_url替换adtracking.click逻辑。因需要记录标记分析数据，因此需要在renderurl前
		isAdClickReplaceByClickUrl(&params, campaign, &ad)
		// 判断是否需要走h5发点击上报的逻辑
		isH5ToClick(&params)
		// AdClickReplaceclickByClick 回退逻辑
		adClickBack(&params, &ad)
	}

	// 服务端点击设备去重时间窗abtest
	serverUniqClickTimeABTest(&params, campaign)

	// abtest 框架处理
	adParamABTest(&params, campaign, &ad, mvconst.FakeAdserverDsp)

	// 封装abtest 结果参数
	renderABTestParams(&params)

	//ios 13 storekit time 下毒逻辑
	prisonStorekitTime(&params, &ad)

	// 针对903，三方为af的单子,rv,iv流量的情况下，添加token宏
	addAppsflyerTokenParam(&params, campaign)

	// online api 返回出价
	renderOnlineApiBidPrice(&params, &ad, r, corsairCampaign)

	// 限制算法出价过高的情况，避免算法出价过高导致的爆量
	if filterRequestByHighBidPrice(&params, &ad) {
		watcher.AddWatchValue("bid_price_out_of_limit", float64(1))
		mvutil.Logger.Runtime.Errorf("online api bid price out of limit. bid price is:[%s],requestid is:[%s]", params.OnlineApiBidPrice, params.RequestID)
		return ad, errors.New("online api bid price out of limit")
	}

	// support smart vba
	renderSupportSmartVBA(&params, &ad, campaign)

	// 封装params
	renderParams(ad, &params, campaign, corsairCampaign)

	// demand lib abtest
	demandSideLibABTest(&params, campaign, r) // 会更新价格

	renderOfferSkadnetwork(corsairCampaign, &ad, &params)

	if params.PriceIn == 0 && campaign.OriPrice != 0 {
		params.PriceIn = campaign.OriPrice
	}

	if params.PriceOut == 0 && campaign.Price != 0 {
		params.PriceOut = campaign.Price
	}
	// priceout x bt 打折系数
	if params.BtPriceOutPercent != nil {
		params.PriceOut = params.PriceOut * *params.BtPriceOutPercent
	}

	// 3s单子情况
	if params.LocalCurrency == 0 {
		params.LocalCurrency = constant.USD
		params.LocalChannelPriceIn = params.PriceIn
	}

	// ext_bp
	params.Extbp = RenderExtBp(&params, campaign)

	if r.IsBidRequest {
		params.BidPrice = corsairCampaign.GetBidPrice()
	}
	// 记录每个offer对应的k值
	ad.ParamK = params.RequestID

	// 记录extra20
	r.Param.Extra20List = params.Extra20List
	// 记录request的playable
	r.Param.ExtPlayableList = params.ExtPlayableList
	r.Param.Extrvtemplate = params.Extrvtemplate
	r.Param.ImageCreativeid = params.ImageCreativeid
	r.Param.VideoCreativeid = params.VideoCreativeid
	r.Param.ExtABTestList = params.ExtABTestList
	r.Param.ExtBigTplOfferDataList = params.ExtBigTplOfferDataList

	// 处理appwall的rs参数
	if (mvutil.IsAppwallOrMoreOffer(r.Param.AdType) && r.Param.MofVersion < 2) || params.AdType == mvconst.ADTypeOfferWall || params.AdType == mvconst.ADTypeWXAppwall {
		r.Param.QueryRsList = append(params.QueryRsList, params.QueryCRs)
	}
	if r.Param.AdType == mvconst.ADTypeInteractive {
		r.Param.IAOrientation = params.IAOrientation
		r.Param.IAPlayableUrl = params.IAPlayableUrl
	}

	// for V4, before render P & Q, set RequestType Back to V4
	if r.Param.RequestPath == mvconst.PATHOpenApiV4 {
		r.Param.RequestType = mvconst.REQUEST_TYPE_OPENAPI_V4
	}

	RenderThirdDemandAKS(corsairCampaign, &ad)
	if strings.HasSuffix(params.RequestID, "v") {
		r.Param.CNTrackingDomainTag = true
		params.Domain = extractor.GetTrackingCNABTestConf().Domain
	}
	// render new url
	RenderNewUrls(r, &ad, &params, campaign, 0, 0, false)

	// 生成Url
	RenderUrls(r, &ad, &params, campaign)
	if corsairCampaign.GetAdTracking() != nil {
		adTracking := corsairCampaign.GetAdTracking()
		if ad.AdTrackingPoint == nil {
			ad.AdTrackingPoint = &CAdTracking{}
		}
		ad.AdTrackingPoint.Click = append(ad.AdTrackingPoint.Click, adTracking.Click...)
		ad.AdTrackingPoint.Impression = append(ad.AdTrackingPoint.Impression, adTracking.Impression...)
		if ad.AdTrackingPoint.Play_percentage == nil && len(adTracking.PlayPerct) > 0 {
			ad.AdTrackingPoint.Play_percentage = make([]CPlayTracking, 0, len(adTracking.PlayPerct))
		}

		for _, playPerct := range adTracking.PlayPerct {
			ad.AdTrackingPoint.Play_percentage = append(ad.AdTrackingPoint.Play_percentage, CPlayTracking{
				Rate: int(playPerct.Rate),
				Url:  playPerct.URL,
			})
		}
	}

	//chetconfig-link
	if link, ok := hitChetLinkConfig(&params); ok {
		if ad.AdTrackingPoint == nil {
			ad.AdTrackingPoint = &CAdTracking{}
		}
		ad.AdTrackingPoint.Impression = append(ad.AdTrackingPoint.Impression, replaceChetLink(link, &params))
	}

	// 针对sdk有问题版本下毒。541,542,550,551，这些版本无法解析adtracking.imp给h5
	if ReplaceImpt2withImp(&params) && ad.AdTrackingPoint != nil && len(ad.AdTrackingPoint.Impression) > 0 {
		for _, imp := range ad.AdTrackingPoint.Impression {
			ad.AdTrackingPoint.Play_percentage = append(ad.AdTrackingPoint.Play_percentage, CPlayTracking{Rate: 0, Url: imp})
		}
		//	ad.AdTrackingPoint.Impression_t2 = append(ad.AdTrackingPoint.Impression_t2, ad.AdTrackingPoint.Impression...)
		ad.AdTrackingPoint.Impression = []string{}
	}

	// 渲染adchoice
	if mvutil.IsHbOrV3OrV5Request(params.RequestPath) {
		renderAdchoice4Ad(&ad, campaign)
	}

	ad.ApkFMd5 = renderApkFMd5(&params, campaign.ApkUrl, &ad, campaign.TrackingUrl)
	ad.Ntbarpt = r.Param.Ntbarpt
	ad.Ntbarpasbl = r.Param.Ntbarpasbl
	ad.AtatType = r.Param.AtatType

	// reward plus
	RenderRewardPlus(r, &ad)

	renderApkInfo(&ad, campaign, &params)

	// 默认maitve为0
	ad.Maitve = 0
	if ad.Maitve == 1 {
		ad.MaitveSrc = "Mtg"
	}

	if strings.HasSuffix(params.RequestID, "v") {
		// 因为存在offer侧的切量，因此命中且api_version大于等于1.4的情况，不需要sdk替换url里的域名宏
		RepalceDoMacro(extractor.GetTrackingCNABTestConf().Domain, &ad, &params)
	}

	return ad, nil
}

func renderApkInfo(ad *Ad, campaign *smodel.CampaignInfo, params *mvutil.Params) {
	// 只需针对国内流量（流量country code为CN）&apk广告（link_type为apk）&新版本sdk时返回
	if !IsCNApkTraffic(ad, params) {
		return
	}

	if len(campaign.NetworkCid) == 0 {
		return
	}
	advOfferInfo, ok := extractor.GetAdvOffer(campaign.NetworkCid)
	if !ok {
		return
	}

	if len(advOfferInfo.AppName) == 0 && len(advOfferInfo.SensitivePermission) == 0 && len(advOfferInfo.OriginSensitivePermission) == 0 &&
		len(advOfferInfo.PrivacyUrl) == 0 && len(advOfferInfo.AppVersionUpdateTime) == 0 && len(advOfferInfo.AppVersion) == 0 &&
		len(advOfferInfo.AppVersion) == 0 && len(advOfferInfo.DeveloperName) == 0 {
		return
	}

	var apkInfo ApkInfo
	apkInfo.AppName = advOfferInfo.AppName
	apkInfo.SensitivePermission = advOfferInfo.SensitivePermission
	apkInfo.OriginSensitivePermission = advOfferInfo.OriginSensitivePermission
	apkInfo.PrivacyUrl = advOfferInfo.PrivacyUrl
	apkInfo.AppVersionUpdateTime = advOfferInfo.AppVersionUpdateTime
	apkInfo.AppVersion = advOfferInfo.AppVersion
	apkInfo.DeveloperName = advOfferInfo.DeveloperName

	ad.ApkInfo = &apkInfo
}

func IsCNApkTraffic(ad *Ad, params *mvutil.Params) bool {
	return ad.CampaignType == mvconst.LINK_TYPE_APK && params.CountryCode == "CN" && params.FormatSDKVersion.SDKVersionCode >= mvconst.AndroidSupportApkInfoVersion
}

func RenderRewardPlus(r *mvutil.RequestParams, ad *Ad) {
	if !mvutil.IsHbOrV3OrV5Request(r.Param.RequestPath) {
		return
	}
	// 限制rv广告位
	if r.Param.AdType != mvconst.ADTypeRewardVideo {
		return
	}
	// 针对支持的sdk版本下发
	if r.Param.PlatformName == constant.PLATFORM_ANDROID_NAME && r.Param.FormatSDKVersion.SDKVersionCode < mvconst.AndroidSupportApkInfoVersion {
		return
	}
	if r.Param.PlatformName == constant.PLATFORM_IOS_NAME && r.Param.FormatSDKVersion.SDKVersionCode < mvconst.IosSupportOfferRewardPlusVersion {
		return
	}

	var rewardPlus RewardPlus
	rewardPlus.CurrencyId = r.UnitInfo.Setting.CurrencyId
	rewardPlus.CurrencyDesc = r.UnitInfo.Setting.CurrencyDesc
	if len(rewardPlus.CurrencyDesc) == 0 {
		rewardPlus.CurrencyDesc = "DefaultName"
	}
	currencyName, ok := r.UnitInfo.Setting.CurrencyName[r.Param.Language]
	var languageArr []string
	if ok {
		rewardPlus.CurrencyName = currencyName
	} else {
		languageArr = strings.Split(r.Param.Language, "-")
	}
	// 前两位取配置
	if len(rewardPlus.CurrencyName) == 0 && len(languageArr) >= 2 {
		currencyName, ok := r.UnitInfo.Setting.CurrencyName[languageArr[0]+"-"+languageArr[1]]
		if ok {
			rewardPlus.CurrencyName = currencyName
		}
	}
	// 首位取配置
	if len(rewardPlus.CurrencyName) == 0 && len(languageArr) >= 1 {
		currencyName, ok := r.UnitInfo.Setting.CurrencyName[languageArr[0]]
		if ok {
			rewardPlus.CurrencyName = currencyName
		}
	}

	rewardPlus.CurrencyReward = r.UnitInfo.Setting.CurrencyReward
	if rewardPlus.CurrencyReward == 0 {
		rewardPlus.CurrencyReward = 1
	}
	rewardPlus.CurrencyRewardPlus = r.UnitInfo.Setting.CurrencyRewardPlus
	rewardPlus.CurrencyCbType = r.UnitInfo.Setting.CurrencyCbType
	rewardPlus.CurrencyIcon = r.UnitInfo.Setting.CurrencyIcon

	ad.RewardPlus = &rewardPlus
}

func renderSupportSmartVBA(params *mvutil.Params, ad *Ad, campaign *smodel.CampaignInfo) {
	cfg := extractor.GetSupportSmartVBAConfig()
	if !cfg.Status {
		return
	}

	for _, item := range cfg.Items {
		if renderSupportSmartWTike(params, ad, campaign, item.WTick) {
			break
		}
	}

	for _, item := range cfg.Items {
		if renderSupportSmartReplacePackageName(params, ad, campaign, item.ReplacePackageName) {
			break
		}
	}
}

func selectorFilter(params *mvutil.Params, campaign *smodel.CampaignInfo, selector *mvutil.SmartVBASelector) bool {
	if len(selector.IncCampaignPackageName) == 0 || (!mvutil.InStrArray(campaign.PackageName, selector.IncCampaignPackageName) &&
		!mvutil.InStrArray("-1", selector.IncCampaignPackageName)) {
		return false
	}

	if len(selector.ExcCampaignPackageName) > 0 && mvutil.InStrArray(campaign.PackageName, selector.ExcCampaignPackageName) {
		return false
	}

	if len(selector.IncCampaignIds) == 0 || (!mvutil.InInt64Arr(campaign.CampaignId, selector.IncCampaignIds) &&
		!mvutil.InInt64Arr(-1, selector.IncCampaignIds)) {
		return false
	}

	if len(selector.ExcCampaignIds) > 0 && mvutil.InInt64Arr(campaign.CampaignId, selector.ExcCampaignIds) {
		return false
	}

	wtick := params.ExtDataInit.WTick
	packageReplace := params.ExtDataInit.PackageReplace
	if len(selector.IncWtick) == 0 || (!mvutil.InArray(wtick, selector.IncWtick) &&
		!mvutil.InArray(-1, selector.IncWtick)) {
		return false
	}

	if len(selector.ExcWtick) > 0 && mvutil.InArray(wtick, selector.ExcWtick) {
		return false
	}

	if len(selector.IncReplacePackage) == 0 || (!mvutil.InArray(packageReplace, selector.IncReplacePackage) &&
		!mvutil.InArray(-1, selector.IncReplacePackage)) {
		return false
	}

	if len(selector.ExcReplacePackage) > 0 && mvutil.InArray(packageReplace, selector.ExcReplacePackage) {
		return false
	}

	if len(selector.IncRequestTypes) == 0 || (!mvutil.InArray(params.RequestType, selector.IncRequestTypes) &&
		!mvutil.InArray(-1, selector.IncRequestTypes)) {
		return false
	}

	if len(selector.ExcRequestTypes) > 0 && mvutil.InArray(params.RequestType, selector.ExcRequestTypes) {
		return false
	}

	if len(selector.IncPackageNames) == 0 || (!mvutil.InStrArray(params.ExtfinalPackageName, selector.IncPackageNames) &&
		!mvutil.InStrArray("-1", selector.IncPackageNames)) {
		return false
	}

	if len(selector.ExcPackageNames) > 0 && mvutil.InStrArray(params.ExtfinalPackageName, selector.ExcPackageNames) {
		return false
	}

	if len(selector.IncAdTypes) == 0 || (!mvutil.InArray(params.CreativeAdType, selector.IncAdTypes) &&
		!mvutil.InArray(-1, selector.IncAdTypes)) {
		return false
	}

	if len(selector.ExcAdTypes) > 0 && mvutil.InArray(params.CreativeAdType, selector.ExcAdTypes) {
		return false
	}

	if len(selector.IncPublisherIds) == 0 || (!mvutil.InInt64Arr(params.PublisherID, selector.IncPublisherIds) &&
		!mvutil.InInt64Arr(-1, selector.IncPublisherIds)) {
		return false
	}

	if len(selector.ExcPublisherIds) > 0 && mvutil.InInt64Arr(params.PublisherID, selector.ExcPublisherIds) {
		return false
	}

	if len(selector.IncAppIds) == 0 || (!mvutil.InInt64Arr(params.AppID, selector.IncAppIds) &&
		!mvutil.InInt64Arr(-1, selector.IncAppIds)) {
		return false
	}

	if len(selector.ExcAppIds) > 0 && mvutil.InInt64Arr(params.AppID, selector.ExcAppIds) {
		return false
	}

	if len(selector.IncUnitIds) == 0 || (!mvutil.InInt64Arr(params.UnitID, selector.IncUnitIds) &&
		!mvutil.InInt64Arr(-1, selector.IncUnitIds)) {
		return false
	}

	if len(selector.ExcUnitIds) > 0 && mvutil.InInt64Arr(params.UnitID, selector.ExcUnitIds) {
		return false
	}

	return true
}

func renderSupportSmartWTike(params *mvutil.Params, ad *Ad, campaign *smodel.CampaignInfo, cfg *mvutil.SmartVBAWTick) bool {
	if cfg == nil {
		return false
	}

	if cfg.Selector == nil {
		return false
	}

	if !selectorFilter(params, campaign, cfg.Selector) {
		return false
	}

	if mvutil.GetRandByGlobalTagId(params, mvconst.SALT_WTICK, 100) < cfg.Rate {
		params.ExtDataInit.WTick = 1
		ad.WithOutInstallCheck = 1
	} else {
		params.ExtDataInit.WTick = 2
	}
	return true
}

func renderSupportSmartReplacePackageName(params *mvutil.Params, ad *Ad, campaign *smodel.CampaignInfo, cfg *mvutil.SmartVBAWReplacePackageName) bool {
	if cfg == nil {
		return false
	}

	if cfg.Selector == nil {
		return false
	}

	if !selectorFilter(params, campaign, cfg.Selector) {
		return false
	}

	if mvutil.GetRandByGlobalTagId(params, mvconst.SALT_REPLACE_PACKCAGE, 100) < cfg.Rate {
		params.ExtDataInit.PackageReplace = 1
		ad.PackageName = cfg.ReplacePackageName
	} else {
		params.ExtDataInit.PackageReplace = 2
	}
	return true
}

func renderApkFMd5(params *mvutil.Params, url string, ad *Ad, trackingUrl string) string {
	if params.Platform == mvconst.PlatformAndroid && ad.CampaignType == mvconst.LINK_TYPE_APK {
		if len(url) > 0 {
			return mvutil.Md5(url)
		}
		// 当apk单子，没有apk url的时候，使用tracking url兜底
		if len(trackingUrl) > 0 {
			return mvutil.Md5(trackingUrl)
		}
	}
	return ""
}

func filterRequestByHighBidPrice(params *mvutil.Params, ad *Ad) bool {
	if params.RequestType != mvconst.REQUEST_TYPE_OPENAPI_AD || ad.OnlineApiBidPrice == 0 {
		return false
	}
	conf := extractor.GetONLINE_API_MAX_BID_PRICE()
	// unit>app>pub>total
	unitStr := strconv.FormatInt(params.UnitID, 10)
	if maxBidPrice, ok := conf.UnitBidPrice[unitStr]; ok && maxBidPrice > 0 {
		// 出价比配置的max 出价还高，则过滤
		return ad.OnlineApiBidPrice > maxBidPrice
	}
	appStr := strconv.FormatInt(params.AppID, 10)
	if maxBidPrice, ok := conf.AppBidPrice[appStr]; ok && maxBidPrice > 0 {
		return ad.OnlineApiBidPrice > maxBidPrice
	}
	pubStr := strconv.FormatInt(params.PublisherID, 10)
	if maxBidPrice, ok := conf.PubBidPrice[pubStr]; ok && maxBidPrice > 0 {
		return ad.OnlineApiBidPrice > maxBidPrice
	}
	if conf.TotalBidPrice > 0 && ad.OnlineApiBidPrice > conf.TotalBidPrice {
		return true
	}
	return false
}

func renderOnlineApiBidPrice(params *mvutil.Params, ad *Ad, r *mvutil.RequestParams, corsairCampaign corsair_proto.Campaign) {
	// 限定online api 流量
	if params.RequestType != mvconst.REQUEST_TYPE_OPENAPI_AD {
		return
	}
	value, ifFind := extractor.GetSspProfitDistributionRuleByUnitIdAndCountryCode(params.UnitID, params.CountryCode)
	if !ifFind || value == nil {
		return
	}
	if value.Type != mvconst.SspProfitDistributionRuleOnlineApiEcpm {
		return
	}

	// 对于fixedecpm配置的值为0的情况（遵循算法出价），则使用算法出价，对于有值的情况下，表示运营希望固定出价
	if value.FixedEcpm > 0 {
		ad.OnlineApiBidPrice = value.FixedEcpm
	}
	// 支持debug bidprice
	renderDebugBidPrice(ad, params)

	// dspext有值表示切量到adx
	_, err := r.GetDspExt()
	if err == nil && ad.OnlineApiBidPrice == 0 {
		// dsp ext里的单位是美分，需要转成美元
		// 统一获取dspbidprice，保证获取到的精度，价格一致
		if dspPriceRft, err := GetAdxPrice(r, &corsairCampaign); err == nil {
			ad.OnlineApiBidPrice = dspPriceRft
		}
	}

	// dspext没有值则表示直接请求as
	bidPrice := corsairCampaign.GetBidPrice()
	if bidPrice > 0 && ad.OnlineApiBidPrice == 0 {
		ad.OnlineApiBidPrice = bidPrice
	}
	if ad.OnlineApiBidPrice > 0 {
		// 控制bid price 精度
		ad.OnlineApiBidPrice, _ = strconv.ParseFloat(fmt.Sprintf("%.8f", ad.OnlineApiBidPrice), 64)
		// 记录bid_price，tracking使用
		params.OnlineApiBidPrice = strconv.FormatFloat(ad.OnlineApiBidPrice, 'f', 8, 64)
	}
}

func renderDebugBidPrice(ad *Ad, params *mvutil.Params) {
	debugBidFloorAndBidPriceConf := extractor.GetDEBUG_BID_FLOOR_AND_BID_PRICE_CONF()
	if len(debugBidFloorAndBidPriceConf) == 0 {
		return
	}
	key := mvutil.GetDebugBidFloorAndBidPriceKey(params)
	if conf, ok := debugBidFloorAndBidPriceConf[key]; ok && conf != nil {
		ad.OnlineApiBidPrice = conf.DebugBidPrice
	}
}

func IsMtgPid(url string) bool {
	val := "pid=mintegral_int"
	// 配置化
	shareitAfWhiteListConf := extractor.GetSHAREIT_AF_WHITE_LIST_CONF()
	if shareitAfWhiteListConf != nil && len(shareitAfWhiteListConf.IsMtgPidValue) > 0 {
		val = shareitAfWhiteListConf.IsMtgPidValue
	}
	return strings.Contains(url, val)
}

func RenderThirdDemandAKS(corsairCampaign corsair_proto.Campaign, ad *Ad) {
	if corsairCampaign.AKS == nil || len(corsairCampaign.AKS) == 0 {
		return
	}

	if ad.AKS == nil {
		ad.AKS = make(map[string]string)
	}

	for k, v := range corsairCampaign.AKS {
		if !strings.HasPrefix(k, "__") || !strings.HasSuffix(k, "__") ||
			strings.Contains(k, "{") || strings.Contains(k, "}") {
			// 特殊限制
			continue
		}
		ad.AKS[k] = v
	}
}

var ImageCreativeType = []int64{
	int64(ad_server.CreativeType_SIZE_320x50),
	int64(ad_server.CreativeType_SIZE_300x250),
	int64(ad_server.CreativeType_SIZE_480x320),
	int64(ad_server.CreativeType_SIZE_320x480),
	int64(ad_server.CreativeType_SIZE_300x300),
	int64(ad_server.CreativeType_SIZE_1200x627),
	int64(ad_server.CreativeType_JS_TAG_320x50),
	int64(ad_server.CreativeType_JS_TAG_300x250),
	int64(ad_server.CreativeType_JS_TAG_480x320),
	int64(ad_server.CreativeType_JS_TAG_320x480),
	int64(ad_server.CreativeType_JS_TAG_300x300),
	int64(ad_server.CreativeType_JS_TAG_1200x627),
	int64(ad_server.CreativeType_SIZE_1125x310),
	int64(ad_server.CreativeType_SIZE_1080x540),
}

var VideoCreativeType = []int64{
	int64(ad_server.CreativeType_VIDEO),
	int64(ad_server.CreativeType_JS_TAG),
}

func creativeCompressABTestV2(params *mvutil.Params, campaign *corsair_proto.Campaign, content map[ad_server.CreativeType]*protobuf.Creative) {
	// 获取整体切量配置
	abtestConfs := extractor.GetADNET_CREATIVE_COMPRESS_ABTEST_CONF_V3()
	params.CreativeCompressData = make(map[ad_server.CreativeType]mvutil.CreativeCompressDataMap)
	for crType, cr := range content {
		renderVideoCompressABTest(abtestConfs, cr, params, campaign, crType)
	}
}

func renderVideoCompressABTest(abtestConfs map[string]*mvutil.CreativeCompressABTestV3Data, cr *protobuf.Creative, params *mvutil.Params, campaign *corsair_proto.Campaign, crType ad_server.CreativeType) {
	var conf *mvutil.CreativeCompressABTestV3Data
	if videoConf, ok := abtestConfs["video"]; ok && videoConf != nil && mvutil.InInt64Arr(int64(crType), VideoCreativeType) {
		conf = videoConf
	} else if imageConf, ok := abtestConfs["image"]; ok && imageConf != nil && mvutil.InInt64Arr(int64(crType), ImageCreativeType) {
		conf = imageConf
	} else if iconConf, ok := abtestConfs["icon"]; ok && iconConf != nil && crType == ad_server.CreativeType_ICON {
		conf = iconConf
	}

	if conf == nil {
		return
	}

	// 判断是否需要进入实验
	if crCompressBlackList(conf, params, campaign) {
		return
	}
	// 判断切量
	resKey := getCrCpmpressABTestResByConf(conf.RateMap, params)
	if len(resKey) == 0 {
		return
	}

	var creativeCompressDataMap mvutil.CreativeCompressDataMap
	adnCreativeIdStr := strconv.FormatInt(cr.CreativeId, 10)
	// 需要在请求日志中记录切量标记，以分析不同分组下load失败的结果。
	if params.ABTestTags == nil {
		params.ABTestTags = make(map[string]int)
	}
	// 根据crType替换值
	if mvutil.InInt64Arr(int64(crType), VideoCreativeType) {
		// 获取config表，查看此视频素材有无video size，若没有，则不需要做abtest
		videoCompressABTestMap := extractor.GetCREATIVE_VIDEO_COMPRESS_ABTEST()
		videoSize, ok := videoCompressABTestMap[adnCreativeIdStr]
		if !ok {
			return
		}

		// 记录切量标记
		params.ExtDataInit.VideoCompressAbtestTag, _ = strconv.Atoi(resKey)
		params.ABTestTags[mvconst.ABTestTagVideoCompress] = params.ExtDataInit.VideoCompressAbtestTag
		// 若为实验组，则替换cdn 地址以及 video size，并将fmd5清空
		if resKey != mvconst.CreativeCompress {
			return
		}

		creativeCompressDataMap.Fmd5 = ""
		creativeCompressDataMap.VideoSize = videoSize
		creativeCompressDataMap.Url = strings.Replace(cr.Url, mvconst.BeforeCompressPath, mvconst.AfterCompressPath, 1)
	} else if mvutil.InInt64Arr(int64(crType), ImageCreativeType) {
		//先判断此素材有无切量
		imageCompressABTestMap := extractor.GetCREATIVE_IMG_COMPRESS_ABTEST()
		if _, ok := imageCompressABTestMap[adnCreativeIdStr]; !ok {
			return
		}
		// 记录切量标记
		params.ExtDataInit.ImageCompressAbtestTag, _ = strconv.Atoi(resKey)
		params.ABTestTags[mvconst.ABTestTagImageCompress] = params.ExtDataInit.ImageCompressAbtestTag

		if resKey != mvconst.CreativeCompress {
			return
		}

		creativeCompressDataMap.Url = strings.Replace(cr.Url, mvconst.BeforeCompressPath, mvconst.AfterCompressPath, 1)
	} else if crType == ad_server.CreativeType_ICON {
		//先判断此素材有无切量
		iconCompressABTestMap := extractor.GetCREATIVE_ICON_COMPRESS_ABTEST()
		if _, ok := iconCompressABTestMap[adnCreativeIdStr]; !ok {
			return
		}
		// 记录切量标记
		params.ExtDataInit.IconCompressAbtestTag, _ = strconv.Atoi(resKey)
		params.ABTestTags[mvconst.ABTestTagIconCompress] = params.ExtDataInit.IconCompressAbtestTag

		if resKey != mvconst.CreativeCompress {
			return
		}

		creativeCompressDataMap.Url = strings.Replace(cr.SValue, mvconst.BeforeCompressPath, mvconst.AfterCompressPath, 1)
	}

	params.CreativeCompressData[crType] = creativeCompressDataMap
}

func crCompressBlackList(conf *mvutil.CreativeCompressABTestV3Data, params *mvutil.Params, campaign *corsair_proto.Campaign) bool {
	camId, _ := strconv.ParseInt(campaign.CampaignId, 10, 64)
	if mvutil.InInt64Arr(camId, conf.CampaignBList) {
		return true
	}
	if mvutil.InInt64Arr(params.UnitID, conf.UnitBList) {
		return true
	}
	if mvutil.InInt64Arr(params.AppID, conf.AppBList) {
		return true
	}
	if mvutil.InInt64Arr(params.PublisherID, conf.PubBList) {
		return true
	}
	return false
}

func getCrCpmpressABTestResByConf(weightMap map[string]int, params *mvutil.Params) string {
	adnConf, _ := extractor.GetADNET_SWITCHS()
	var resKey string
	if crCompressABTestRandByDev, ok := adnConf["crCompressABTestRandByDev"]; ok && crCompressABTestRandByDev == 1 {
		resKey = mvutil.RandByDeviceRate(weightMap, params)
	} else {
		resKey = mvutil.RandByRate3(weightMap)
	}
	return resKey
}

func renderAdchoice4Ad(ad *Ad, campaign *smodel.CampaignInfo) {
	advID := int64(campaign.AdvertiserId)
	if advID == 0 {
		return
	}

	advertiserInfo, ok := extractor.GetAdvertiserInfo(advID)
	if !ok {
		return
	}

	if !advertiserInfo.Advertiser.IsShowAdChoice {
		return
	}

	adchoice := &AdChoice{
		AdLogolink:   advertiserInfo.Advertiser.AdLogolink,
		AdchoiceIcon: advertiserInfo.Advertiser.AdchoiceIcon,
		AdchoiceLink: advertiserInfo.Advertiser.AdchoiceLink,
		AdchoiceSize: advertiserInfo.Advertiser.AdchoiceSize,
		AdvLogo:      advertiserInfo.Advertiser.AdvLogo,
		AdvName:      advertiserInfo.Advertiser.AdvName,
		PlatformLogo: advertiserInfo.Advertiser.PlatformLogo,
		PlatformName: advertiserInfo.Advertiser.PlatformName,
	}

	ad.AdChoice = adchoice
	return
}

func renderParams(ad Ad, params *mvutil.Params, campaign *smodel.CampaignInfo, corsairCampaign corsair_proto.Campaign) {
	if !campaign.IsSSPlatform() {
		params.RequestID = mvutil.GetGoTkClickID()
	}

	if campaign.PublisherId == int64(0) {
		params.Extra8 = -1
	} else {
		params.Extra8 = int(campaign.AdvertiserId)
	}
	if corsairCampaign.AdTemplate == nil {
		params.Extra16 = 0
	} else {
		params.Extra16 = int(*(corsairCampaign.AdTemplate))
	}
	// creative r TODO

	// ext_algo rank信息收集
	params.Extalgo = params.AlgoMap[campaign.CampaignId]
	// 其他
	if campaign.AdvertiserId != 0 {
		params.AdvertiserID = campaign.AdvertiserId
	}
	params.CampaignID = campaign.CampaignId
	if campaign.AdSourceId != 0 {
		params.Extra13 = campaign.AdSourceId
	}
	if campaign.Ctype != 0 {
		params.Extctype = campaign.Ctype
	}

	params.Extattr = renderExtAttr(campaign)

	// 标记vtatag和ext_installFrom
	vtaLink := mvutil.GetVTALink(campaign)
	if len(vtaLink) > 0 {
		params.ExtinstallFrom = mvconst.INSTALL_FROM_CLICK
	}
	// jumptype TODO
	renderNativeVideoFlag(ad, params)
	// nativex
	renderNativex(params, campaign)

	if campaign.Source != 0 {
		params.Extsource = campaign.Source
	}
	if corsairCampaign.OfferType != nil {
		params.OfferType = *(corsairCampaign.OfferType)
	}
	// extra20
	params.Extra20List = renderExtra20(params, campaign)
	// request返回ext_playable
	params.ExtPlayableList = renderExtPlayable(params, campaign)

	params.ExtABTestList = renderExtABTestList(params, campaign)

	if params.ExtCampaignTagList != nil {
		if info, ok := params.ExtCampaignTagList[campaign.CampaignId]; ok {
			info.CDNAbTest = params.ExtCDNAbTest
			info.VideoCreativeid = params.VideoCreativeid
		}
	}

	// 记录slot_id到q参数中
	if ad.OfferExtData != nil {
		params.ExtSlotId = strconv.Itoa(int(ad.OfferExtData.SlotId))
	}
	// 整理算法需要的大模板信息
	params.ExtBigTplOfferDataList = renderExtBigTplOfferDataList(params, &ad)

	if params.FlowTagID > 0 {
		RenderMWParams(params, &ad, false)
	}
	// 固定unit的parent_unit已直接赋值，拆分需使用params.ParentUnitId赋值
	if params.ExtDataInit.ParentUnitId == 0 {
		// 记录more offer的parent_id
		params.ExtDataInit.ParentUnitId = params.ParentUnitId
	}
	// 记录h5 type 1表示endcard；2表示playable
	params.ExtDataInit.H5Type = params.H5Type
	// tracking记录req_type aabtest结果
	params.ExtDataInit.ReqTypeTest = params.ReqTypeAABTest
	params.ExtDataInit.CleanDeviceTest = params.CleanDeviceTest
	params.ExtDataInit.DisplayCampaignABTest = params.DisplayCampaignABTest
	params.ExtDataInit.VcnABTest = params.VcnABTest
	//params.ExtDataInit.CtnSizeTest = params.CtnSizeTag
	params.ExtDataInit.IsReplaceAdClick = params.AdClickReplaceclickByClickTag
	params.ExtDataInit.BandWidth = params.BandWidth
	if isClickmode6NotInGpAndAppstore(params, &ad) {
		// 记录标记
		params.ExtDataInit.ClickMode6NotInGpAndAppstore = 1
	}
	// 处理offer维度的extData
	params.ExtData2 = renderExtData(params)

	// 插入adn lib 的abtest标记
	mergeAdnLibABTestTags(params)
}

func mergeAdnLibABTestTags(params *mvutil.Params) {
	if len(params.AdnLibABTestTags) == 0 {
		return
	}
	params.UnmarshalExtData2()
	for key, value := range params.AdnLibABTestTags {
		params.SetExtData2(key, value)
	}
	params.MarshalExtData2()
}

func renderExtBigTplOfferDataList(params *mvutil.Params, ad *Ad) string {
	// 切量部分才会记录
	if params.BigTemplateFlag {
		var buffer bytes.Buffer
		buffer.WriteString(strconv.FormatInt(ad.CampaignID, 10))
		buffer.WriteString(":")
		buffer.WriteString(strconv.FormatInt(params.ImageCreativeid, 10))
		buffer.WriteString(",")
		buffer.WriteString(strconv.FormatInt(params.VideoCreativeid, 10))
		buffer.WriteString(",")
		buffer.WriteString(strconv.FormatInt(params.EndcardCreativeID, 10))
		buffer.WriteString(",")
		buffer.WriteString(strconv.Itoa(params.Extrvtemplate))
		buffer.WriteString(",")
		buffer.WriteString(params.Extendcard)
		buffer.WriteString(",")
		buffer.WriteString(params.ExtSlotId)
		return buffer.String()
	}
	return ""
}

func renderDeleteDevID(ad *Ad, params *mvutil.Params, campaign *smodel.CampaignInfo) {
	rate := int(campaign.SendDeviceidRate)
	if rate < 0 || rate >= 100 {
		return
	}
	randV := rand.Intn(100)
	if randV >= rate {
		params.ExtdeleteDevid = mvconst.DELETE_DIVICEID_TRUE
		ad.ClickMode, _ = strconv.Atoi(mvconst.JUMP_TYPE_NORMAL)
		params.Extra10 = mvconst.JUMP_TYPE_NORMAL
	} else {
		params.ExtdeleteDevid = mvconst.DELETE_DIVICEID_FALSE
		ad.ClickMode, _ = strconv.Atoi(mvconst.JUMP_TYPE_CLIENT_SEND_DEVID)
		params.Extra10 = mvconst.JUMP_TYPE_CLIENT_SEND_DEVID
	}
}

func renderExtra20(params *mvutil.Params, campaign *smodel.CampaignInfo) []string {
	s1 := strconv.FormatInt(params.CampaignID, 10)
	s2 := strconv.FormatInt(int64(params.AdvertiserID), 10)
	s3 := strconv.FormatInt(int64(params.OfferType), 10)
	s4 := "0"
	s5 := strconv.Itoa(mvutil.GetVTATag(campaign))
	s6 := strconv.Itoa(params.CUA)
	s7 := ""
	// appwall,offerwall,online api不返回requestid
	if !mvutil.InArray(int(params.AdType), []int{mvconst.ADTypeAppwall, mvconst.ADTypeOfferWall, mvconst.ADTypeOnlineVideo, mvconst.ADTypeWXAppwall, mvconst.ADTypeMoreOffer}) {
		s7 = params.RequestID
	}
	list := []string{s1, s2, s3, s4, s5, s6, s7}
	for k, v := range list {
		if v == "0" {
			list[k] = ""
		}
	}
	return list
}

type StringArray []string

// 分渠道出价
func RenderExtBp(params *mvutil.Params, campaign *smodel.CampaignInfo) string {
	var arr StringArray
	priceIn := params.PriceIn
	priceOut := params.PriceOut
	if params.UseAlgoPrice {
		if params.AlgoPriceIn > 0 {
			priceIn = params.AlgoPriceIn
		}

		if params.AlgoPriceOut > 0 {
			priceOut = params.AlgoPriceOut
		}
	}
	s1 := mvutil.FormatFloat64(priceIn)
	s2 := mvutil.FormatFloat64(priceOut)
	s3 := strconv.FormatInt(int64(campaign.Ctype), 10)
	s4 := s3
	if campaign.CostType != 0 {
		s4 = strconv.FormatInt(int64(campaign.CostType), 10)
	}
	s5 := strconv.Itoa(params.LocalCurrency)
	s6 := mvutil.FormatFloat64(params.LocalChannelPriceIn)
	arr = append(arr, s1, s2, s3, s4, s5, s6)
	str, err := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(arr)
	if err != nil {
		return ""
	}
	return string(str)
}

func renderNativex(params *mvutil.Params, campaign *smodel.CampaignInfo) {
	params.Extnativex = 2
	if campaign.BelongType == 1 {
		params.Extnativex = 1
	}
	return
	// 不再兜底
}

func renderNativeVideoFlag(ad Ad, params *mvutil.Params) {
	params.ExtnativeVideo = 0
	if params.ApiVersion >= mvconst.API_VERSION_1_0 {
		return
	}
	if params.AdType == mvconst.ADTypeNative && len(params.VideoVersion) > 0 && params.VideoVersion != "0" && len(ad.VideoURL) > 0 {
		params.ExtnativeVideo = 1
		return
	}
}

func renderBT(r mvutil.RequestParams, params *mvutil.Params, btType ad_server.BtType, campaign *smodel.CampaignInfo) {
	// 统一赋值包名
	if r.Param.DspMof == 1 {
		params.ExtfinalPackageName = r.Param.PackageName
	} else {
		params.ExtfinalPackageName = r.AppInfo.RealPackageName
	}
	// D级流量统一
	subId := HandleGradeD(r)
	params.Extra14 = subId
	if params.AppID == params.Extra14 {
		// 获取adserver返回的bttype
		if btType == ad_server.BtType_BTOFFER {
			// bt
			subId := BlendTraffic(r, campaign, params)
			params.Extra14 = subId
		}
	}
	// 渠道信息透明化
	RenderRequestPackage(r, *campaign, params)
}

func renderCampaignInfo(ad *Ad, r *mvutil.RequestParams, campaign *smodel.CampaignInfo, params *mvutil.Params, corsairCampaign corsair_proto.Campaign) {
	ad.CampaignID = campaign.CampaignId
	// AdSourceID
	ad.AdSourceID = r.Param.AdSourceID
	if campaign.PackageName != "" {
		ad.PackageName = campaign.PackageName
	}
	if campaign.Ctype != 0 {
		ad.CType = int(campaign.Ctype)
	}
	if corsairCampaign.OfferType != nil {
		ad.OfferType = int(*(corsairCampaign.OfferType))
	}
	// 默认返回advId
	if campaign.AdvertiserId != 0 {
		ad.AdvID = int(campaign.AdvertiserId)
	}
	// 处理numberRating
	ad.NumberRating = HandleNumberRating(ad.NumberRating)
	ad.Price = float32(params.PriceOut)
	// 获取下毒的camid 及fca的默认值
	camIds, _ := extractor.GetFCA_CAMIDS()
	newDefaultFca, _ := extractor.GetNEW_DEFAULT_FCA()
	ad.FCA = mvutil.GetFCA(r, campaign, camIds, newDefaultFca)
	ad.FCB = mvutil.GetFCB(params.CreativeDescSource)
	if len(ad.CtaText) <= 0 {
		if campaign.CampaignType == 0 {
			ad.CtaText = "install"
		} else {
			ad.CtaText = handleCTAButton(campaign.CampaignType, params)
		}
	}
	if corsairCampaign.AdTemplate != nil {
		ad.Template = int(*(corsairCampaign.AdTemplate))
	}
	// 封装adUrlList advImp，要在bt之后
	renderImpUrl(ad, r, campaign, params)
	// offerwall，要在bt之后
	if corsairCampaign.OfferType == nil {
		tmpOfferType := int32(0)
		corsairCampaign.OfferType = &tmpOfferType
	}

	offerType := int32(0)
	if corsairCampaign.OfferType != nil {
		offerType = *(corsairCampaign.OfferType)
	}
	// 获取templategroupid
	if corsairCampaign.TemplateGroup != nil {
		params.TemplateGroupId = int(*corsairCampaign.TemplateGroup)
	}

	renderOfferwall(ad, r, offerType, params.PriceOut)

	// linktype
	renderLinkType(ad, campaign, params)

	if IsRvNvNewCreative(r.Param.NewCreativeFlag, r.Param.AdType) {
		// rv
		renderRewardVideoTemplateNew(ad, params, corsairCampaign)
		// endcard
		renderEndcardNew(ad, params, corsairCampaign)
		// videoEndType
		if corsairCampaign.VideoEndTypeAs != nil {
			ad.VideoEndType = int(*corsairCampaign.VideoEndTypeAs)
		}
		// 记录templategroupid，由more offer逻辑下发templategroupid使用
		r.Param.TplGroupId = int(*corsairCampaign.TemplateGroup)
	}

	// nvt2 在endcard之后
	if params.ApiVersion >= mvconst.API_VERSION_1_3 && params.AdType == mvconst.ADTypeNative {
		ad.NVT2 = r.UnitInfo.Unit.NVTemplate
		params.Extnvt2 = r.UnitInfo.Unit.NVTemplate
		if ad.NVT2 != int32(3) {
			ad.EndcardUrl = ""
		}
	}

	// 处理endcard_url，服务端控制endcard 逻辑（自动storekit、全局可点、猜你喜欢）
	renderEndcardProperty(ad, campaign, params, r)

	// 不使用ad_server的返回值，也需记录ad_server返回的分配情况
	recordPlayableTag(corsairCampaign, params)

	// jumptype

	//renderJumpTypeInfo(ad, r, campaign, params)
	// 目前903都是ss单了
	// 这样context.CampaignID才会被赋值
	params.CampaignID = campaign.CampaignId
	params.AdvertiserID = campaign.AdvertiserId
	ctx := NewDemandContext(params, r)
	previewUrl := demand.RenderThirdPartyLink(ctx, extractor.DemandDao, demand.GetCampaignNormalLandingPage(campaign))
	if previewUrl != "" {
		ad.PreviewUrl = previewUrl
	}
	// 获取clickmode
	clickmodeContext := NewClickmodeContext(params, r)
	jumpType, adClickMode := mclickmode.RenderClickMode(clickmodeContext, extractor.DemandDao)
	if clickmodeContext != nil {
		// 记录切量维度的标记
		params.ExtDataInit.ClickmodeGlobalConfigTag = clickmodeContext.ADNABTestData["cmgc_t"]
		params.ExtDataInit.ClickmodeGroupTag = clickmodeContext.ClickModeGroupTag
	}
	params.Extra10 = jumpType
	ad.ClickMode, _ = strconv.Atoi(adClickMode)
	// 新增标记看灰度效果
	params.ExtDataInit.NewClickmodeTag = 1
	renderNewClickMode(params, ad)

	// 去devid逻辑
	renderDeleteDevID(ad, params, campaign)

	renderSettingInfo(ad, params, campaign)

	// jssdk CDN域名处理
	RenderJssdkCDN(ad, *params)
	// render https
	RenderHttpsUrls(ad, *params)

	//根据配置切换CDN域名
	RenderCreativeUrls(ad, r, params)

	// 根据app来切量
	toNewCDN(ad, params)
	ad.RetargetOffer = 2
	if campaign.RetargetOffer > 0 {
		ad.RetargetOffer = int(campaign.RetargetOffer)
	}

	// TODO impua cua
	renderUa(ad, campaign, params)
	// loopback
	renderLoopback(ad, campaign)
	// 整理videoEndType
	renderVideoInfo(ad, params, r)
	// 处理playableAdsWithoutVideo
	renderPlayableAdsWithoutVideo(ad, corsairCampaign, r)
	// 素材三期记录日志
	if IsRvNvNewCreative(r.Param.NewCreativeFlag, r.Param.AdType) {
		renderExtCreativeNew(ad, params, r, corsairCampaign)
	}

	if len(campaign.SubCategoryName) > 0 {
		ad.SubCategoryName = campaign.SubCategoryName
	}

	// 降展示
	if params.ExtifLowerImp == int32(1) {
		ad.ImageURL = ""
		ad.VideoURL = ""
	}
	if campaign.AppSize != "" {
		ad.AppSize = campaign.AppSize
	}

	//storekit_time
	ad.StoreKitTime = 1
	if r.AppInfo.App.StorekitLoading == 2 {
		ad.StoreKitTime = 2
	}
	// c_toi
	cToi, _ := extractor.GetC_TOI()
	ad.CToi = cToi

	// jm_icon
	renderJM(ad, *params)

	// 针对360max返回ad.category
	if params.RequestPath == mvconst.PATHMAXADX {
		if campaign.Category != 0 {
			ad.Category = campaign.Category
		} else {
			ad.Category = 2
		}

	}

	if len(campaign.WxAppId) > 0 {
		ad.WxAppId = campaign.WxAppId
	}
	if len(campaign.WxPath) > 0 {
		ad.WxPath = mvutil.Base64Encode(campaign.WxPath)
	}
	if len(campaign.BindId) > 0 {
		ad.BindId = campaign.BindId
	}
	if len(campaign.DeepLink) > 0 {
		ad.DeepLink = campaign.DeepLink
	}
	// 360max使用
	if len(campaign.ApkVersion) > 0 {
		ad.ApkVersion = campaign.ApkVersion
	}
	if len(campaign.ApkMd5) > 0 {
		ad.ApkMd5 = campaign.ApkMd5
	}
	// 拉活offer
	if campaign.UserActivation == 1 {
		ad.UserActivation = true
	}
	// request接口使用
	if r.Param.RequestType == mvconst.REQUEST_TYPE_SDK && len(campaign.ApkUrl) > 0 {
		ad.ApkUrl = campaign.ApkUrl
	}
	// webview下毒
	if r.Param.NeedWebviewPrison {
		ad.EndcardUrl = ""
		ad.RvPoint = nil
	}
	// 获取adsource维度配置的有效缓存时间，备用缓存时间
	renderPlct(ad, params, mvconst.Mobvista, 0)

	// 针对安卓，MD5_file下毒
	//if r.Param.Platform == mvconst.PlatformAndroid {
	//	pubsConf, _ := extractor.GetMD5_FILE_PRISON_PUB()
	//	if mvutil.InInt64Arr(r.Param.PublisherID, pubsConf) {
	//		ad.Md5File = ""
	//	}
	//}
	// 设置readyRate
	renderReadyRate(ad)

	// 封装ext_data
	renderOfferExtData(ad, params, corsairCampaign)

	//APK alert
	renderOfferApkAlt(r, ad)

	// 获取小米deeplink链接
	renderXiaomiDeeplinkUrl(params, ad, ctx)

	// 当投放lazada deeplink单子的情况下，目前对于online api的流量，会单独建立单链单子来跑（deeplink 字段为空值）
	// 而给sdk跑的单子是需要有deeplink的，运营同学现希望lazada或以后的deeplink单子，也使用sdk投放的单子来跑，减少与广告主对数的成本。
	// 因此需要把online api流量，deeplink 双链的单子改为单链来投放。
	changeOnlineDeeplinkWay(r, ad)
}

func renderXiaomiDeeplinkUrl(params *mvutil.Params, ad *Ad, ctx *demand.Context) {
	// 原本单子有deeplink，不参与，以免影响效果
	if len(ad.DeepLink) > 0 {
		return
	}
	ad.DeepLink = demand.RenderXiaoMiDeeplink(ctx, extractor.DemandDao)
	// 记录切量标记
	params.AdnLibABTestTags = ctx.AdnLibABTestTags
}

func renderNewClickMode(params *mvutil.Params, ad *Ad) {
	// click mode 14,返回给sdk的clickmode需要为5
	if params.Extra10 == mvconst.JUMP_TYPE_TRACKING_PING_THIRDPARTY_CLICK_URL {
		ad.ClickMode = 5
	} else if params.Extra10 == mvconst.JUMP_TYPE_ONLINE_DSP_AJUMP_SERVER || params.Extra10 == mvconst.JUMP_TYPE_NORMAL {
		ad.ClickMode = 0
		params.PingMode = 0
	}
}

func RenderTemplateCreativeDomainMacro(tpl string) string {
	tplCreativeDomainConf := extractor.GetTPL_CREATIVE_DOMAIN_CONF()
	if tplCreativeDomainConf == nil {
		return tpl
	}
	for domainMacro, domainWeightMaps := range tplCreativeDomainConf {
		if len(domainWeightMaps) == 0 {
			continue
		}
		// 选择宏需要替换的域名
		domain, _ := getDomainByWeight(domainWeightMaps)
		// 记录标记
		tpl = strings.ReplaceAll(tpl, "__"+domainMacro+"__", domain)
	}
	return tpl
}

func changeOnlineDeeplinkWay(r *mvutil.RequestParams, ad *Ad) {
	// 必须为仅支持online api单链的流量
	if r.Param.RequestType != mvconst.REQUEST_TYPE_OPENAPI_AD {
		return
	}
	// 限定deeplink单子
	if len(ad.DeepLink) == 0 {
		return
	}
	// 限制offer包名
	changeOnlineDeeplinkWayPackageList := extractor.GetCHANGE_ONLINE_DEEPLINK_WAY_PACKAGE_LIST()
	if !mvutil.InStrArray(ad.PackageName, changeOnlineDeeplinkWayPackageList) {
		return
	}
	// 不支持双链的流量
	adnConfList := extractor.GetADNET_CONF_LIST()
	if onlineUnsupportTwoLinkPubList, ok := adnConfList["onlineUnsupportTwoLinkPubList"]; ok && mvutil.InInt64Arr(r.Param.PublisherID, onlineUnsupportTwoLinkPubList) {
		ad.DeepLink = ""
	}
}

func renderOfferApkAlt(r *mvutil.RequestParams, ad *Ad) {
	if IsCNApkTraffic(ad, &r.Param) && mvutil.InInt32Arr(r.Param.AdType, r.AppInfo.App.ApkAltAdtype) {
		ad.ApkAlt = 1
		return
	}
}

func renderOfferSkadnetwork(corsairCampaign corsair_proto.Campaign, ad *Ad, params *mvutil.Params) {
	if corsairCampaign.Skad != nil && len(corsairCampaign.Skad.Version) > 0 {
		ad.Skadnetwork = &Skadnetwork{
			Version:         corsairCampaign.Skad.Version,
			Network:         corsairCampaign.Skad.Network,
			AppleCampaignId: corsairCampaign.Skad.AppleCampaignId,
			Targetid:        corsairCampaign.Skad.Targetid,
			Nonce:           corsairCampaign.Skad.Nonce,
			Sourceid:        corsairCampaign.Skad.Sourceid,
			Timestamp:       corsairCampaign.Skad.Timestamp,
			Sign:            corsairCampaign.Skad.Sign,
			Need:            int(corsairCampaign.Skad.Need),
		}
		return
	}
	if len(params.SkSign) > 0 {
		var skadnetworkVer string
		if params.Skadnetwork != nil {
			skadnetworkVer = params.Skadnetwork.Ver
		}
		ad.Skadnetwork = &Skadnetwork{
			Version:         skadnetworkVer,
			Network:         params.SkAdNetworkId,
			AppleCampaignId: params.SkCid,
			Targetid:        params.SkTargetId,
			Nonce:           params.SkNonce,
			Sourceid:        params.SkSourceId,
			Timestamp:       params.SkTmp,
			Sign:            params.SkSign,
			Need:            params.SkNeed,
		}
		if len(params.SkViewSign) > 0 {
			ad.SkImp = &SkImp{
				ViewSign: params.SkViewSign,
			}
		}
	}
}

func renderOfferExtData(ad *Ad, params *mvutil.Params, campaign corsair_proto.Campaign) {
	var offerExtData OfferExtData
	// 优先获取campaign维度的slotid。（新的取值方式）
	if campaign.SlotId != nil {
		offerExtData.SlotId = *campaign.SlotId
	} else {
		// 记录slot id
		for slotId, campaignId := range params.BigTemplateSlotMap {
			if campaignId != ad.CampaignID {
				continue
			}
			offerExtData.SlotId = slotId
		}
	}

	ad.OfferExtData = &offerExtData
}

func renderReadyRate(ad *Ad) {
	adnConf, _ := extractor.GetADNET_SWITCHS()
	if readyRateVal, ok := adnConf["readyRate"]; ok {
		ad.ReadyRate = readyRateVal
	} else {
		ad.ReadyRate = 100
	}
}

func renderJM(ad *Ad, params mvutil.Params) {
	if params.AdType != mvconst.ADTypeJMIcon {
		return
	}
	if ad.CampaignType == 3 {
		pd := 2
		ad.JMPD = &pd
	}
}

func renderVideoInfo(ad *Ad, params *mvutil.Params, r *mvutil.RequestParams) {
	if ad.VideoEndType <= 0 {
		// 整理videoEndType
		videoEndType := 2
		if r.UnitInfo.Unit.VideoEndType > 0 {
			videoEndType = r.UnitInfo.Unit.VideoEndType
		}
		if params.ApiVersion >= mvconst.API_VERSION_1_2 {

		} else {
			if videoEndType <= 0 || videoEndType > 5 {
				videoEndType = 2
			}
		}
		ad.VideoEndType = videoEndType
	}
	// 整理storekit
	if params.Platform == mvconst.PlatformIOS {
		storekit := 0
		if params.ApiVersion >= mvconst.API_VERSION_1_3 && params.AdType == mvconst.ADTypeNative && params.Extnvt2 == int32(5) {
			storekit = mvconst.StorekitLoad
		}
		// 素材三期后v4模块id为401,402，此前为4
		if ad.Rv.VideoTemplate == 4 || ad.Rv.VideoTemplate == 401 ||
			ad.Rv.VideoTemplate == 402 || ad.VideoEndType == 6 {
			storekit = mvconst.StorekitLoad
		}
		if mvutil.IsAppwallOrMoreOffer(params.AdType) {
			storekit = mvconst.StorekitNotLoad
		}
		// 针对开发者下毒，不出storekit。设置开关，默认Wie关闭，1则为关闭下毒
		adnConf, _ := extractor.GetADNET_SWITCHS()
		closeSkPoison, ok := adnConf["closeSkPoison"]
		if (!ok || closeSkPoison != 1) && params.PublisherID == 14228 {
			storekit = mvconst.StorekitNotLoad
		}
		// ios sdk 问题版本storekit下毒
		if params.IosStorekitPoisonFlag {
			storekit = mvconst.StorekitNotLoad
		}
		ad.Storekit = storekit

	}

}

func renderPlayableAdsWithoutVideo(ad *Ad, corsairCampaign corsair_proto.Campaign, r *mvutil.RequestParams) {
	if IsRvNvNewCreative(r.Param.NewCreativeFlag, r.Param.AdType) {
		// 素材三期，若ad_server返回的UsageVideo为false则返回2。PlayableAdsWithoutVideo不为0代表走了playable逻辑。
		if corsairCampaign.UsageVideo != nil && *(corsairCampaign.UsageVideo) == false {
			ad.PlayableAdsWithoutVideo = 2
		}
	}

	if ad.PlayableAdsWithoutVideo == 0 {
		ad.PlayableAdsWithoutVideo = 1
	}
}

func renderLoopback(ad *Ad, campaign *smodel.CampaignInfo) {
	if campaign.Loopback == nil {
		return
	}

	loopback := campaign.Loopback
	if loopback.Rate == 0 {
		return
	}

	randInt := rand.Intn(100)
	rate := int(loopback.Rate)
	if randInt < rate {
		ad.LoopBack = map[string]string{
			"domain": loopback.Domain,
			"key":    loopback.Key,
			"value":  loopback.Value,
		}
	}
}

func renderSettingInfo(ad *Ad, params *mvutil.Params, campaign *smodel.CampaignInfo) {
	// campaign or 三方维度 abtest
	cctVal, abtestOk := cctAbtest(ad, params, campaign)
	if abtestOk {
		ad.ClickCacheTime = cctVal
		params.ExtDataInit.CctAbTest = &cctVal
		return
	}

	settingConfs, ifFind := extractor.GetSETTING_CONFIG()
	if !ifFind {
		return
	}
	if params.Platform == mvconst.PlatformAndroid {
		ad.ClickCacheTime = settingConfs.ACCT
		return
	}
	ad.ClickCacheTime = settingConfs.CCT
}

func fixAndroidJumpType(jumpType string, r *mvutil.RequestParams, campaign *smodel.CampaignInfo, params *mvutil.Params) string {
	if r.Param.Platform != mvconst.PlatformAndroid || r.Param.RequestType != mvconst.REQUEST_TYPE_OPENAPI_V3 {
		return jumpType
	}

	// android sdk click mode 6\12 修正
	if jumpType != mvconst.JUMP_TYPE_CLIENT_DO_ALL && jumpType != mvconst.JUMP_TYPE_CLIENT_DO_ALL_PING_SERVER {
		return jumpType
	}

	adnConf, _ := extractor.GetADNET_SWITCHS()
	// 默认关闭，若开启，则gaid为空值的情况则走0
	gaidEmptyCanJump := true
	if isCloseGaid, ok := adnConf["icGaid"]; ok && isCloseGaid == 1 {
		gaidEmptyCanJump = false
	}
	if params.GAID == "" && params.AndroidID == "" && !gaidEmptyCanJump {
		// 无设备ID走0
		return mvconst.JUMP_TYPE_NORMAL
	}

	var linkType int32
	if campaign.OpenType != 0 {
		linkType = campaign.OpenType
	} else if campaign.CampaignType != 0 {
		linkType = campaign.CampaignType
	}

	if linkType != mvconst.GP {
		if len(campaign.DirectUrl) == 0 {
			// 非GP且direct URL为空则走5
			return mvconst.JUMP_TYPE_CLIENT_SEND_DEVID
		}

		if jumpType == mvconst.JUMP_TYPE_CLIENT_DO_ALL {
			// 非GP配置了6，directURL不为空也会走5
			return mvconst.JUMP_TYPE_CLIENT_SEND_DEVID
		}
		// 非GP配置了12，directURL不为空走11
		return mvconst.JUMP_TYPE_CLIENT_SEND_DEVID_PING_SERVER
	}

	if jumpType == mvconst.JUMP_TYPE_CLIENT_DO_ALL {
		// 配置了6满足条件则以配置优先
		return jumpType
	}

	if len(campaign.DirectUrl) == 0 {
		// GP 配置12，DirectURL为空，走6
		return mvconst.JUMP_TYPE_CLIENT_DO_ALL
	}
	return mvconst.JUMP_TYPE_CLIENT_DO_ALL_PING_SERVER
}

//func renderJumpTypeInfo(ad *Ad, r *mvutil.RequestParams, campaign *smodel.CampaignInfo, params *mvutil.Params) {
//	jumpType := RandJumpType(r, campaign, params)
//
//	jumpType = fixAndroidJumpType(jumpType, r, campaign, params)
//	if jumpType == mvconst.JUMP_TYPE_CLIENT_DO_ALL || jumpType == mvconst.JUMP_TYPE_CLIENT_DO_ALL_PING_SERVER {
//		if len(campaign.SdkPackageName) > 0 {
//			ad.PackageName = campaign.SdkPackageName
//		}
//	}
//
//	// 若是onlineapi或者vast 则为0 TODO
//	// 赋值
//	params.Extra10 = jumpType
//
//	// 如果最终返回jumptype=11,则返回给sdk时需要置为5，如果jumptype=12,则返回给sdk时置为6 by 2018.01.30 jj
//	if jumpType == mvconst.JUMP_TYPE_CLIENT_SEND_DEVID_PING_SERVER {
//		jumpType = mvconst.JUMP_TYPE_CLIENT_SEND_DEVID
//	} else if jumpType == mvconst.JUMP_TYPE_CLIENT_DO_ALL_PING_SERVER {
//		jumpType = mvconst.JUMP_TYPE_CLIENT_DO_ALL
//	}
//	params.JumpType = jumpType
//	ad.ClickMode, _ = strconv.Atoi(jumpType)
//}

// renderPlayable3 adserver 的召回逻辑 adserver控制playable
// func renderPlayable3(ad *Ad, corsairCampaign corsair_proto.Campaign, params *mvutil.Params, r *mvutil.RequestParams) {
// 	if !corsairCampaign.Playable { // !corsairCampaign.Playable 只打log,不改ad的属性
// 		// corsairCampaign.ExtPlayable >0 的才打log
// 		if corsairCampaign.ExtPlayable == nil || *corsairCampaign.ExtPlayable == 0 {
// 			return
// 		}
// 		extPlayable := 20000 // 2AABB
// 		switch *corsairCampaign.ExtPlayable {
// 		case int32(1):
// 			extPlayable = extPlayable + 100
// 		case int32(2):
// 			extPlayable = extPlayable + 200
// 		case int32(3):
// 			extPlayable = extPlayable + 300
// 		}
// 		if corsairCampaign.VideoEndTypeAs != nil {
// 			extPlayable = extPlayable + int(*corsairCampaign.VideoEndTypeAs)
// 		}
// 		params.Extplayable = extPlayable
// 		return
// 	}
// 	extPlayable := 10000
// 	playableUrl := ""
// 	if corsairCampaign.Playable && corsairCampaign.EndcardUrl != nil {
// 		playableUrl = *corsairCampaign.EndcardUrl
// 		// handle scheme
// 		if !strings.HasPrefix(playableUrl, "http") {
// 			playableUrl = GetUrlScheme(params) + "://" + playableUrl
// 		}
// 	}

// 	if corsairCampaign.ExtPlayable != nil {
// 		// 1、返回offer的endcard url为空，正常返回video url和endscreen url
// 		// 2、返回offer的endcard url指向playable ads页面，正常返回video url和endscreen url
// 		// 3、返回offer的endcard url指向playable ads页面，正常返回endscreen url，但video url为空，且is_playable_without_video值为true
// 		switch *corsairCampaign.ExtPlayable {
// 		case int32(1):
// 			ad.EndcardUrl = ""
// 			extPlayable = extPlayable + 100
// 		case int32(2):
// 			ad.EndcardUrl = playableUrl
// 			extPlayable = extPlayable + 200
// 		case int32(3):
// 			ad.EndcardUrl = playableUrl
// 			ad.PlayableAdsWithoutVideo = 2
// 			ad.VideoURL = ""
// 			extPlayable = extPlayable + 300
// 		}
// 	}
// 	if corsairCampaign.VideoEndTypeAs != nil {
// 		ad.VideoEndType = int(*corsairCampaign.VideoEndTypeAs)
// 		extPlayable = extPlayable + ad.VideoEndType
// 	}
// 	params.Extplayable = extPlayable
// }

// renderPlayableNew Adnet的召回逻辑
// func renderPlayableNew(ad *Ad, campaign *smodel.CampaignInfo, params *mvutil.Params, r *mvutil.RequestParams) {
// 	if params.AdType != mvconst.ADTypeRewardVideo && params.AdType != mvconst.ADTypeInterstitialVideo {
// 		return
// 	}
// 	if campaign.Endcard == nil {
// 		return
// 	}
// 	// version compare
// 	if !Compare(r, "playable") {
// 		return
// 	}
// 	var randInt int
// 	// 中国安卓流量切gaid为空情况下，取ip做随机值
// 	if r.Param.CountryCode == "CN" && len(r.Param.GAID) == 0 && r.Param.Platform == mvconst.PlatformAndroid {
// 		randInt = mvutil.GetRandConsiderZero(r.Param.ClientIP, r.Param.IDFA, mvconst.SALT_PLAYABLE, 100)
// 	} else {
// 		randInt = mvutil.GetRandConsiderZero(r.Param.GAID, r.Param.IDFA, mvconst.SALT_PLAYABLE, 100)
// 	}
// 	if randInt == -1 {
// 		return
// 	}
// 	// 获取conf
// 	confs := campaign.Endcard
// 	endcard := getEndcard(confs, r)

// 	endcardUrl := ""
// 	if params.ApiVersion >= mvconst.API_VERSION_1_1 {
// 		endcardUrl = endcard.UrlV2
// 	}
// 	if len(endcardUrl) <= 0 {
// 		endcardUrl = endcard.Url
// 	}
// 	if len(endcardUrl) <= 0 {
// 		return
// 	}

// 	i := 0
// 	params.Extplayable = 0
// 	for k, v := range endcard.EndcardRate {
// 		if randInt < v+i {
// 			ik, _ := strconv.Atoi(k)
// 			params.Extplayable = ik
// 			break
// 		}
// 		i = i + v
// 	}
// 	// 如果url没有scheme，则要加上
// 	playableUrl := ""
// 	switch endcard.EndcardProtocal {
// 	case 1:
// 		playableUrl = "http://" + endcardUrl
// 	case 2:
// 		playableUrl = "https://" + endcardUrl
// 	default:
// 		playableUrl = GetUrlScheme(params) + "://" + endcardUrl
// 	}

// 	// 1、返回offer的endcard url为空，正常返回video url和endscreen url
// 	// 2、返回offer的endcard url指向playable ads页面，正常返回video url和endscreen url
// 	// 3、返回offer的endcard url指向playable ads页面，正常返回endscreen url，但video url为空，且is_playable_without_video值为true
// 	switch params.Extplayable {
// 	case 1:
// 		ad.EndcardUrl = ""
// 	case 2:
// 		ad.EndcardUrl = playableUrl
// 	case 3:
// 		ad.EndcardUrl = playableUrl
// 		ad.PlayableAdsWithoutVideo = 2
// 		ad.VideoURL = ""
// 	}
// }

// renderPlayable 已弃用， 改用另两个函数
// func renderPlayable(ad *Ad, campaign *smodel.CampaignInfo, params *mvutil.Params, r *mvutil.RequestParams) {
// 	if params.AdType != mvconst.ADTypeRewardVideo {
// 		return
// 	}
// 	// version compare
// 	if !Compare(r, "playable") {
// 		return
// 	}
// 	confs, ifFind := extractor.GetPlayableTest()
// 	if !ifFind {
// 		return
// 	}
// 	conf, ok := confs[campaign.CampaignId]
// 	if !ok {
// 		return
// 	}
// 	randInt := mvutil.GetRandConsiderZero(r.Param.GAID, r.Param.IDFA, mvconst.SALT_PLAYABLE, 100)
// 	if randInt == -1 {
// 		return
// 	}
// 	i := 0
// 	params.Extplayable = 0
// 	for k, v := range conf.Rate {
// 		if randInt < v+i {
// 			params.Extplayable = k
// 			break
// 		}
// 		i = i + v
// 	}
// 	// 如果url没有scheme，则要加上
// 	if !strings.Contains(conf.Url, "http") {
// 		conf.Url = GetUrlScheme(params) + "://" + conf.Url
// 	}
// 	// 1、返回offer的endcard url为空，正常返回video url和endscreen url
// 	// 2、返回offer的endcard url指向playable ads页面，正常返回video url和endscreen url
// 	// 3、返回offer的endcard url指向playable ads页面，正常返回endscreen url，但video url为空，且is_playable_without_video值为true
// 	switch params.Extplayable {
// 	case 1:
// 		ad.EndcardUrl = ""
// 	case 2:
// 		ad.EndcardUrl = conf.Url
// 	case 3:
// 		ad.EndcardUrl = conf.Url
// 		ad.PlayableAdsWithoutVideo = 2
// 		ad.VideoURL = ""
// 	}
// }

func renderLinkType(ad *Ad, campaign *smodel.CampaignInfo, params *mvutil.Params) {
	//替换神回避的linkType
	unitId := strconv.FormatInt(params.UnitID, 10)
	lt, Ok := extractor.GetLINKTYPE_UNITID()
	if Ok {
		if v, ok := lt[unitId]; ok {
			params.LinkType = v
			ad.CampaignType = v
			return
		}
	}

	if campaign.OpenType != 0 {
		openType := campaign.OpenType
		if (!mvutil.IsMpad(params.RequestPath) && params.ApiVersion < mvconst.API_VERSION_1_3) && (openType == int32(8) || openType == int32(9)) ||
			(params.Platform == mvconst.PlatformIOS && len(mvutil.GetIdfaString(params.IDFA)) == 0 && strings.ToLower(campaign.ThirdParty) == mvconst.THIRD_PARTY_S2S) {
			openType = int32(4)
		}
		// 针对mp，sdk<4.xx且linktype=8/9则强制转为4
		if mvutil.GetMPSdkVersionCompare(params.SDKVersion) && (openType == int32(8) || openType == int32(9)) {
			openType = int32(4)
		}
		// 针对mp mtg请求，linktype=3，sdk version为4.1以下的版本linktype返回4,4.1及以上的按原样返回
		openType = renderMpOpenType(*params, openType)
		ad.CampaignType = int(openType)
		params.LinkType = int(openType)
		return
	}

	linkType := int(campaign.CampaignType)
	// 针对mp，sdk<4.xx且linktype=8/9则强制转为4
	if mvutil.InArray(linkType, []int{5, 6, 7}) ||
		(params.Platform == mvconst.PlatformIOS && len(mvutil.GetIdfaString(params.IDFA)) == 0 && strings.ToLower(campaign.ThirdParty) == mvconst.THIRD_PARTY_S2S) ||
		(mvutil.GetMPSdkVersionCompare(params.SDKVersion) && (linkType == 8 || linkType == 9)) {
		linkType = 4
	}
	// 针对mp mtg请求，linktype=3，sdk version为4.1以下的版本linktype返回4,4.1及以上的按原样返回
	linkTypeInt32 := renderMpOpenType(*params, int32(linkType))
	linkType = int(linkTypeInt32)
	ad.CampaignType = linkType
	params.LinkType = linkType
}

// 必须在bt后，用的价格是bt后的价格
func renderOfferwall(ad *Ad, r *mvutil.RequestParams, offerType int32, price float64) {
	if r.Param.AdType != int32(mvconst.ADTypeOfferWall) {
		return
	}
	// guidelines
	conf, ifFind := extractor.GetOfferwallGuidelines()
	if ifFind {
		ad.Guidelines = conf[offerType]
	}
	// reward
	reward := r.UnitInfo.VirtualReward
	ad.RewardName = reward.Name
	if offerType == int32(1) {
		ad.RewardAmount = reward.StaticReward
	} else if offerType == int32(2) {
		amount := float64(reward.ExchangeRate) * price
		ad.RewardAmount = int(math.Floor(amount))
	}
}

func renderImpUrl(ad *Ad, r *mvutil.RequestParams, campaign *smodel.CampaignInfo, params *mvutil.Params) {
	ad.AdvImp = []CAdvImp{}
	for _, v := range campaign.AdvImp {
		var cAdvImp CAdvImp
		cAdvImp.Sec = int(v.Sec)
		if v.Url != "" {
			cAdvImp.Url = RenderImpUrl(v.Url, params, r)
		}
		ad.AdvImp = append(ad.AdvImp, cAdvImp)
	}

	ad.AdURLList = []string{}
	for _, v := range campaign.AdUrlList {
		url := RenderImpUrl(v, params, r)
		ad.AdURLList = append(ad.AdURLList, url)
	}
}

func handleCTAButton(linkType int32, params *mvutil.Params) string {
	language := handleLanguage(params.Language)
	/**
	  * campaignType枚举值
	     1 => 'AppStore',
	     2 => 'GooglePlay',
	     3 => 'APK',
	     5 => 'IPA',
	     6 => 'Subscription',
	     7 => 'Website',
	     4 => 'Other',
	*/
	if linkType == 1 || linkType == 2 || linkType == 3 {
		cta := mvconst.CTALang[language]
		if len(cta) > 0 {
			return cta
		}
	} else if linkType == 4 || linkType == 6 || linkType == 7 || linkType == 9 {
		cta := mvconst.CTAViewLang[language]
		if len(cta) > 0 {
			return cta
		}
	}
	return "install"
}

func handleLanguage(language string) string {
	// 去掉空格
	language = strings.Replace(language, " ", "", -1)
	// 如果没有-，则直接返回
	if !strings.Contains(language, "-") {
		return language
	}
	// 如果是中文
	if strings.Contains(language, "zh-Hant") {
		return "zh-Hant"
	}
	if strings.Contains(language, "zh-Hans") {
		return "zh-Hans"
	}
	if strings.Contains(language, "zh-TW") {
		return "zh-Hant"
	}
	arr := strings.Split(language, "-")
	return arr[0]
}

// [10000, 50000]
func HandleNumberRating(numberRating int) int {
	if numberRating > 10000 {
		return numberRating
	}
	return rand.Intn(40001) + 10000
}

func renderGInfo(crType ad_server.CreativeType, cr *protobuf.Creative) string {
	gInfoList := []string{
		strconv.FormatInt(int64(crType), 10),
		strconv.FormatInt(cr.CreativeId, 10),
		cr.AdvCreativeId,
		cr.VideoResolution,
		"",
		strconv.FormatInt(cr.UniqCId, 10),
	}
	return strings.Join(gInfoList, ",")
}

func renderCreativePb(ad *Ad, params *mvutil.Params, corsairCampaign corsair_proto.Campaign, content map[ad_server.CreativeType]*protobuf.Creative, campaign *smodel.CampaignInfo) {
	params.CreativeId = int64(0)
	var advCrMap = map[string]*mvutil.AdvCreativeMap{}

	creativeDco := 0
	gInfoList := make([]string, 0, len(corsairCampaign.CreativeTypeIdMap))
	gidList := make([]int, 0, len(corsairCampaign.CreativeTypeIdMap))
	gInfoForAppwall := make([]string, 0, 1)

	for crType, cr := range content {
		gidList = append(gidList, int(cr.CreativeId))
		gInfoList = append(gInfoList, renderGInfo(crType, cr))
		//crTypeStr := crType.String()
		switch crType {
		case ad_server.CreativeType_APP_NAME:
			// appName
			//appName, _ := cr.Value.(string)
			ad.AppName = cr.SValue
		case ad_server.CreativeType_APP_DESC:
			// appDesc
			//appDesc := cr.SValue
			ad.AppDesc = mvutil.SubUtf8Str(cr.SValue, APP_DESC_LEN)
			params.CreativeDescSource = cr.Source
		case ad_server.CreativeType_ICON:
			// iconUrl
			//iconUrl, _ := cr.Value.(string)
			if crCompressData, ok := params.CreativeCompressData[crType]; ok {
				cr.SValue = crCompressData.Url
			}
			ad.IconURL = cr.SValue
			ad.IconMime = cr.Mime
			ad.IconResolution = cr.Resolution
			if cr.AdvCreativeId != "" && cr.AdvCreativeId != "0" {
				advCrMap["icon"] = &mvutil.AdvCreativeMap{
					AdvCreativeId:      cr.AdvCreativeId,
					AdvCreativeName:    cr.Cname,
					AdvCreativeGroupId: cr.CsetId,
				}
			}
			// 封装appwall的glist
			if mvutil.IsAppwallOrMoreOffer(params.AdType) || params.AdType == mvconst.ADTypeOfferWall || params.AdType == mvconst.ADTypeWXAppwall {
				gInfoForAppwall = []string{renderGInfo(crType, cr)}
			}
		case ad_server.CreativeType_APP_RATE:
			// appScore
			// rating, _ := cr.Value.(float64)
			ad.Rating = float32(mvutil.NumFormat(cr.FValue, 1))
		case ad_server.CreativeType_CTA_BUTTON:
			// ctatext
			//CtaText, _ := cr.Value.(string)
			ad.CtaText = cr.SValue
		case ad_server.CreativeType_SIZE_320x50,
			ad_server.CreativeType_SIZE_300x250,
			ad_server.CreativeType_SIZE_480x320,
			ad_server.CreativeType_SIZE_320x480,
			ad_server.CreativeType_SIZE_300x300,
			ad_server.CreativeType_SIZE_1200x627,
			ad_server.CreativeType_JS_TAG_320x50,
			ad_server.CreativeType_JS_TAG_300x250,
			ad_server.CreativeType_JS_TAG_480x320,
			ad_server.CreativeType_JS_TAG_320x480,
			ad_server.CreativeType_JS_TAG_300x300,
			ad_server.CreativeType_JS_TAG_1200x627,
			ad_server.CreativeType_SIZE_1125x310,
			ad_server.CreativeType_SIZE_1080x540,
			ad_server.CreativeType_SIZE_640x640,
			ad_server.CreativeType_SIZE_720x576:
			// image
			// 素材压缩abtest修改返回值
			if crCompressData, ok := params.CreativeCompressData[crType]; ok {
				cr.Url = crCompressData.Url
			}
			renderImageInfo(ad, params, cr)
			if cr.AdvCreativeId != "" && cr.AdvCreativeId != "0" {
				advCrMap["image"] = &mvutil.AdvCreativeMap{
					AdvCreativeId:      cr.AdvCreativeId,
					AdvCreativeName:    cr.Cname,
					AdvCreativeGroupId: cr.CsetId,
				}
			}
		case ad_server.CreativeType_GIF_360x640,
			ad_server.CreativeType_GIF_640x360:
			// 目前仅小程序rv广告位使用
			if params.AdType == mvconst.ADTypeWXRewardImg {
				ad.GifURL = cr.Url
				if cr.AdvCreativeId != "" && cr.AdvCreativeId != "0" {
					advCrMap["image"] = &mvutil.AdvCreativeMap{
						AdvCreativeId:      cr.AdvCreativeId,
						AdvCreativeName:    cr.Cname,
						AdvCreativeGroupId: cr.CsetId,
					}
				}
			}
		case ad_server.CreativeType_GIF_1200x627:
			ad.GifURL = cr.Url
			if cr.AdvCreativeId != "" && cr.AdvCreativeId != "0" {
				advCrMap["image"] = &mvutil.AdvCreativeMap{
					AdvCreativeId:      cr.AdvCreativeId,
					AdvCreativeName:    cr.Cname,
					AdvCreativeGroupId: cr.CsetId,
				}
			}
		case ad_server.CreativeType_GIF_INDUCED:
			// 只有rv，iv有诱导gif素材
			if mvutil.IsIvOrRv(params.AdType) {
				imageArr := Image{
					IdcdImg: []string{cr.Url},
				}
				ad.Rv.Image = &imageArr
				if cr.AdvCreativeId != "" && cr.AdvCreativeId != "0" {
					advCrMap["image"] = &mvutil.AdvCreativeMap{
						AdvCreativeId:      cr.AdvCreativeId,
						AdvCreativeName:    cr.Cname,
						AdvCreativeGroupId: cr.CsetId,
					}
				}
			}
		case ad_server.CreativeType_VIDEO, ad_server.CreativeType_JS_TAG:
			// 素材压缩abtest修改返回值
			if crCompressData, ok := params.CreativeCompressData[crType]; ok {
				cr.Url = crCompressData.Url
				cr.VideoSize = crCompressData.VideoSize
				cr.FMd5 = crCompressData.Fmd5
			}
			// video
			// videourl要base64encode
			ad.VideoURL = mvutil.Base64Encode(cr.Url)
			ad.VideoLength = int(cr.VideoLength)
			ad.VideoSize = int(cr.VideoSize)
			ad.VideoResolution = cr.VideoResolution
			watchMile := float64(cr.VideoLength * cr.WatchMile / 100)
			ad.WatchMile = int(watchMile)
			ad.VideoWidth = cr.Width
			ad.VideoHeight = cr.Height
			ad.Bitrate = cr.BitRate
			ad.VideoMime = cr.Mime
			//ad.Md5File = cr.FMd5
			params.VideoFmd5 = cr.FMd5
			params.CreativeId = cr.CreativeId
			params.VideoCreativeid = cr.CreativeId
			if cr.AdvCreativeId != "" && cr.AdvCreativeId != "0" {
				advCrMap["video"] = &mvutil.AdvCreativeMap{
					AdvCreativeId:      cr.AdvCreativeId,
					AdvCreativeName:    cr.Cname,
					AdvCreativeGroupId: cr.CsetId,
				}
			}
		case ad_server.CreativeType_ENDCARD,
			ad_server.CreativeType_ENDCARD_ZIP:
			// endcard
			ad.EndcardUrl = cr.Url
			params.EndcardCreativeID = cr.CreativeId
			if cr.AdvCreativeId != "" && cr.AdvCreativeId != "0" {
				advCrMap["playable"] = &mvutil.AdvCreativeMap{
					AdvCreativeId:      cr.AdvCreativeId,
					AdvCreativeName:    cr.Cname,
					AdvCreativeGroupId: cr.CsetId,
				}
			}
		case ad_server.CreativeType_COMMENT:
			// rating
			//comment, _ := cr.IValue
			ad.NumberRating = int(cr.IValue)
		case ad_server.CreativeType_PLAYABLE_URL,
			ad_server.CreativeType_PLAYABLE_ZIP,
			ad_server.CreativeType_AR_URL,
			ad_server.CreativeType_AR_ZIP,
			ad_server.CreativeType_FULLVIEW_PLAYABLE_URL,
			ad_server.CreativeType_FULLVIEW_PLAYABLE_ZIP:
			// playable
			playableUrl := GetPlayableUrl(cr.Url, int(cr.Protocal), params)
			ad.EndcardUrl = playableUrl
			params.EndcardCreativeID = cr.CreativeId
			params.IAPlayableUrl = playableUrl
			params.IAOrientation = int(cr.Orientation)
			if cr.AdvCreativeId != "" && cr.AdvCreativeId != "0" {
				advCrMap["playable"] = &mvutil.AdvCreativeMap{
					AdvCreativeId:      cr.AdvCreativeId,
					AdvCreativeName:    cr.Cname,
					AdvCreativeGroupId: cr.CsetId,
				}
			}
		case ad_server.CreativeType_SIZE_720x1280,
			ad_server.CreativeType_SIZE_1280x720:
			// 小程序appwall二维码大图使用
			if params.AdType == mvconst.ADTypeWXNative || params.AdType == mvconst.ADTypeWXBanner || params.AdType == mvconst.ADTypeWXRewardImg {
				ad.ExtImg = cr.Url
			} else if params.AdType == mvconst.ADTypeWXAppwall {
				// image
				ad.ImageURL = cr.Url
				ad.ImageMime = cr.Mime
				ad.ImageResolution = cr.Resolution
				params.CreativeId = cr.CreativeId
				params.ImageCreativeid = cr.CreativeId
				if cr.AdvCreativeId != "" && cr.AdvCreativeId != "0" {
					advCrMap["image"] = &mvutil.AdvCreativeMap{
						AdvCreativeId:      cr.AdvCreativeId,
						AdvCreativeName:    cr.Cname,
						AdvCreativeGroupId: cr.CsetId,
					}
				}
			} else if params.RequestType == mvconst.REQUEST_TYPE_OPENAPI_AD {
				renderImageInfo(ad, params, cr)
				if cr.AdvCreativeId != "" && cr.AdvCreativeId != "0" {
					advCrMap["image"] = &mvutil.AdvCreativeMap{
						AdvCreativeId:      cr.AdvCreativeId,
						AdvCreativeName:    cr.Cname,
						AdvCreativeGroupId: cr.CsetId,
					}
				}
			}
		}
		//素材dco
		if dcoTitle, ok := corsairCampaign.DynamicCreative[crType]; ok {
			if crType == ad_server.CreativeType_APP_NAME {
				ad.AppName = dcoTitle
				creativeDco = 1
			} else if crType == ad_server.CreativeType_APP_DESC {
				ad.AppDesc = dcoTitle
				creativeDco = 1
			}
		}
	}

	// create query r
	var qr mvutil.QueryR
	qr.Gid = getCreativeGroupID(gidList)
	if corsairCampaign.AdElementTemplate != nil {
		qr.Tpid = int(*(corsairCampaign.AdElementTemplate))
	}
	qr.Crat = getQueryRAdType(*params, *ad)
	// 针对more offer情况进行处理。因为目前广告主仍未识别more_offer这种广告形式，因此，返回映射值11会被归为other，导致数据对不齐。
	// 对于自身流量，又需要做ad_type的区分，因此针对上报的ad_type，沿用appwall
	qr.Crat = getMofCrat(params, qr.Crat)

	mainAdvCreativeMap := getMainAdvCreativeMap(advCrMap)
	if mainAdvCreativeMap != nil {
		advCrid, _ := strconv.Atoi(mainAdvCreativeMap.AdvCreativeId)
		qr.AdvCrid = advCrid
		qr.Cname = mainAdvCreativeMap.AdvCreativeName
		qr.CsetName = getCsetName(mainAdvCreativeMap.AdvCreativeGroupId)
	}

	if campaign.IsCampaignCreative != 0 {
		qr.Icc = int(campaign.IsCampaignCreative)
	}

	qr.Pi = params.PriceIn
	qr.Po = params.PriceOut
	qr.Dco = creativeDco
	qr.Glist = strings.Join(gInfoList, "|")
	// 记录pcdIds
	qr.CpdIds = strings.Join(corsairCampaign.CpdIds, ",")
	queryR, err := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(qr)
	if err == nil {
		params.QueryR = mvutil.Base64(queryR)
	}
	//针对appwall处理only_impression的rs参数
	if (mvutil.IsAppwallOrMoreOffer(params.AdType) && params.MofVersion < 2) || params.AdType == mvconst.ADTypeOfferWall || params.AdType == mvconst.ADTypeWXAppwall {
		params.QueryCRs = renderAWParamRs(qr, gInfoForAppwall, campaign.CampaignId)
	}

	// url replace
	params.CreativeAdType = qr.Crat
	params.AdvCreativeID = qr.AdvCrid
	params.CreativeName = qr.Cname
	params.CsetName = qr.CsetName

	// 其他
	if params.RequestType == mvconst.REQUEST_TYPE_OPENAPI_AD && corsairCampaign.ImageSize != nil &&
		len(*corsairCampaign.ImageSize) > 0 &&
		corsairCampaign.ImageSizeId == ad_server.ImageSizeEnum_UNKNOWN {
		ad.ImageSize = *corsairCampaign.ImageSize
	} else {
		ad.ImageSize = mvconst.GetImageSizeByID(int(corsairCampaign.ImageSizeId))
	}
	// 针对小程序的banner广告位，固定返回1125x310(wx_banner只有这种尺寸的大图)
	if params.AdType == mvconst.ADTypeWXBanner {
		ad.ImageSize = "1125x310"
	}
	params.ImageSize = ad.ImageSize

	// REQUEST_TYPE_SDK，imageUrl为空则使用iconUrl
	if params.RequestPath == mvconst.PATHREQUEST && len(ad.ImageURL) == 0 {
		ad.ImageURL = ad.IconURL
	}
}

func getCsetName(csetId int64) string {
	creativePackage, ifFind := extractor.GetCreativePackage(csetId)
	if !ifFind {
		return ""
	}
	return creativePackage.CsetName
}

func renderImageInfo(ad *Ad, params *mvutil.Params, cr *protobuf.Creative) {
	ad.ImageURL = cr.Url
	ad.ImageMime = cr.Mime
	ad.ImageResolution = cr.Resolution
	params.CreativeId = cr.CreativeId
	params.ImageCreativeid = cr.CreativeId
	// 图片md5
	params.ImgFMD5 = cr.FMd5
}

func getMofCrat(params *mvutil.Params, cadtype int) int {
	adnConf, _ := extractor.GetADNET_SWITCHS()
	if closeCAdType, ok := adnConf["closeChangeMofType"]; ok && closeCAdType == 1 {
		return cadtype
	}
	if params.AdType == mvconst.ADTypeMoreOffer {
		cadtype = mvconst.CREATIVE_AD_TYPE_APPWALL
	}
	return cadtype
}

// func renderCreative(ad *Ad, params *mvutil.Params, corsairCampaign corsair_proto.Campaign, content map[int64]mvutil.Content, campaign *smodel.CampaignInfo) {
// 	params.CreativeId = int64(0)
// 	// 封装r参数
// 	var gidList []int
// 	var advCrIDMap = map[string]string{
// 		"image": "",
// 		"icon":  "",
// 		"video": "",
// 	}
// 	creativeDco := 0
// 	//var gInfoList []string
// 	gInfoList := make([]string, 0, len(corsairCampaign.CreativeTypeIdMap))
// 	for crType, crId := range corsairCampaign.CreativeTypeIdMap {
// 		cr := content[crId]
// 		gidList = append(gidList, int(cr.CreativeId))
// 		gInfoList = append(gInfoList, renderGInfo(crType, cr))
// 		//crTypeStr := crType.String()
// 		switch crType {
// 		case ad_server.CreativeType_APP_NAME:
// 			// appName
// 			appName, _ := cr.Value.(string)
// 			ad.AppName = appName
// 		case ad_server.CreativeType_APP_DESC:
// 			// appDesc
// 			appDesc, _ := cr.Value.(string)
// 			ad.AppDesc = mvutil.SubUtf8Str(appDesc, APP_DESC_LEN)
// 			params.CreativeDescSource = cr.Source
// 		case ad_server.CreativeType_ICON:
// 			// iconUrl
// 			iconUrl, _ := cr.Value.(string)
// 			ad.IconURL = iconUrl
// 			ad.IconMime = cr.Mime
// 			ad.IconResolution = cr.Resolution
// 			advCrIDMap["icon"] = cr.AdvCreativeId
// 		case ad_server.CreativeType_APP_RATE:
// 			// appScore
// 			rating, _ := cr.Value.(float64)
// 			ad.Rating = float32(mvutil.NumFormat(rating, 1))
// 		case ad_server.CreativeType_CTA_BUTTON:
// 			// ctatext
// 			CtaText, _ := cr.Value.(string)
// 			ad.CtaText = CtaText
// 		case ad_server.CreativeType_SIZE_320x50,
// 			ad_server.CreativeType_SIZE_300x250,
// 			ad_server.CreativeType_SIZE_480x320,
// 			ad_server.CreativeType_SIZE_320x480,
// 			ad_server.CreativeType_SIZE_300x300,
// 			ad_server.CreativeType_SIZE_1200x627,
// 			ad_server.CreativeType_JS_TAG_320x50,
// 			ad_server.CreativeType_JS_TAG_300x250,
// 			ad_server.CreativeType_JS_TAG_480x320,
// 			ad_server.CreativeType_JS_TAG_320x480,
// 			ad_server.CreativeType_JS_TAG_300x300,
// 			ad_server.CreativeType_JS_TAG_1200x627:
// 			// image
// 			ad.ImageURL = cr.Url
// 			ad.ImageMime = cr.Mime
// 			ad.ImageResolution = cr.Resolution
// 			params.CreativeId = cr.CreativeId
// 			advCrIDMap["image"] = cr.AdvCreativeId
// 		case ad_server.CreativeType_VIDEO, ad_server.CreativeType_JS_TAG:
// 			// video
// 			// videourl要base64encode
// 			ad.VideoURL = mvutil.Base64Encode(cr.Url)
// 			ad.VideoLength = int(cr.VideoLength)
// 			ad.VideoSize = int(cr.VideoSize)
// 			ad.VideoResolution = cr.VideoResolution
// 			watchMile := math.Ceil(float64(cr.VideoLength * cr.WatchMile / 100))
// 			ad.WatchMile = int(watchMile)
// 			ad.VideoWidth = cr.Width
// 			ad.VideoHeight = cr.Height
// 			ad.Bitrate = cr.BitRate
// 			ad.VideoMime = cr.Mime
// 			ad.Md5File = &cr.FMd5
// 			params.CreativeId = cr.CreativeId
// 			advCrIDMap["video"] = cr.AdvCreativeId
// 		case ad_server.CreativeType_ENDCARD:
// 			// endcard
// 			ad.EndcardUrl = cr.Url
// 			params.EndcardCreativeID = cr.CreativeId
// 		case ad_server.CreativeType_COMMENT:
// 			// rating
// 			comment, _ := cr.Value.(int)
// 			ad.NumberRating = comment
// 		case ad_server.CreativeType_PLAYABLE_URL, ad_server.CreativeType_PLAYABLE_ZIP:
// 			// playable
// 			params.IAPlayableUrl = GetPlayableUrl(cr.Url, cr.Protocal, *params)
// 			params.IAOrientation = cr.Orientation
// 		}
// 		// if crType == ad_server.CreativeType_APP_NAME {

// 		// }
// 		// if mvutil.InStrArray(crTypeStr, mvutil.CreativeAppName) {

// 		// } else if mvutil.InStrArray(crTypeStr, mvutil.CreativeAppDesc) {

// 		// } else if mvutil.InStrArray(crTypeStr, mvutil.CreativeIcon) {

// 		// } else if mvutil.InStrArray(crTypeStr, mvutil.CreativeAppRATE) {

// 		// } else if mvutil.InStrArray(crTypeStr, mvutil.CreativeCtaButton) {

// 		// } else if mvutil.InStrArray(crTypeStr, mvutil.CreativeImageV2) {

// 		// } else if mvutil.InStrArray(crTypeStr, mvutil.CreativeVideo) {

// 		// } else if mvutil.InStrArray(crTypeStr, mvutil.CreativeEndcard) {

// 		// } else if mvutil.InStrArray(crTypeStr, mvutil.CreativeRating) {

// 		// } else if mvutil.InStrArray(crTypeStr, mvutil.Playable) {

// 		// }

// 		//素材dco
// 		if dcoTitle, ok := corsairCampaign.DynamicCreative[crType]; ok {
// 			if crType == ad_server.CreativeType_APP_NAME {
// 				ad.AppName = dcoTitle
// 				creativeDco = 1
// 			} else if crType == ad_server.CreativeType_APP_DESC {
// 				ad.AppDesc = dcoTitle
// 				creativeDco = 1
// 			}
// 		}
// 	}

// 	// create query r
// 	var qr mvutil.QueryR
// 	qr.Gid = getCreativeGroupID(gidList)
// 	if corsairCampaign.AdElementTemplate != nil {
// 		qr.Tpid = int(*(corsairCampaign.AdElementTemplate))
// 	}
// 	qr.Crat = getQueryRAdType(*params, *ad)
// 	advCrid, _ := strconv.Atoi(getQueryRAdvID(qr.Crat, advCrIDMap))
// 	qr.AdvCrid = advCrid
// 	if campaign.IsCampaignCreative != nil {
// 		qr.Icc = int(*(campaign.IsCampaignCreative))
// 	}
// 	if campaign.OriPrice != nil {
// 		qr.Pi = *(campaign.OriPrice)
// 	}
// 	if campaign.Price != nil {
// 		qr.Po = *(campaign.Price)
// 	}
// 	qr.Dco = creativeDco
// 	qr.Glist = strings.Join(gInfoList, "|")
// 	var json = jsoniter.ConfigCompatibleWithStandardLibrary
// 	queryR, err := json.Marshal(qr)
// 	if err == nil {
// 		params.QueryR = mvutil.Base64(queryR)
// 	}
// 	// url replace
// 	params.CreativeAdType = qr.Crat
// 	params.AdvCreativeID = qr.AdvCrid

// 	// 其他
// 	ad.ImageSize = mvconst.GetImageSizeByID(int(corsairCampaign.ImageSizeId))
// 	params.ImageSize = ad.ImageSize
// }

func getMainAdvCreativeMap(advCrMap map[string]*mvutil.AdvCreativeMap) *mvutil.AdvCreativeMap {
	// 按素材优先级上报，不区分adtype
	// 素材优先级为video->playable->image->icon
	if val, ok := advCrMap["video"]; ok && val != nil {
		return val
	}

	if val, ok := advCrMap["playable"]; ok && val != nil {
		return val
	}

	if val, ok := advCrMap["image"]; ok && val != nil {
		return val
	}

	if val, ok := advCrMap["icon"]; ok && val != nil {
		return val
	}

	return nil
}

func getQueryRAdType(params mvutil.Params, ad Ad) int {
	switch params.AdType {
	case mvconst.ADTypeBanner,
		mvconst.ADTypeWXBanner,
		mvconst.ADTypeSdkBanner:
		return mvconst.CREATIVE_AD_TYPE_BANNER
	case mvconst.ADTypeInterstitial,
		mvconst.ADTypeInterstitialSdk:
		return mvconst.CREATIVE_AD_TYPE_INTERSTITIAL
	case mvconst.ADTypeNative,
		mvconst.ADTypeWXNative:
		if len(params.VideoVersion) > 0 && len(ad.VideoURL) > 0 {
			return mvconst.CREATIVE_AD_TYPE_NATIVE_VIDEO
		}
		return mvconst.CREATIVE_AD_TYPE_NATIVE
	case mvconst.ADTypeOnlineVideo,
		mvconst.ADTypeJSBannerVideo,
		mvconst.ADTypeJSNativeVideo:
		return mvconst.CREATIVE_AD_TYPE_NATIVE_VIDEO
	case mvconst.ADTypeAppwall,
		mvconst.ADTypeWXAppwall:
		return mvconst.CREATIVE_AD_TYPE_APPWALL
	case mvconst.ADTypeOfferWall:
		return mvconst.CREATIVE_AD_TYPE_OFFERWALL
	case mvconst.ADTypeRewardVideo,
		mvconst.ADTypeFullScreen:
		return mvconst.CREATIVE_AD_TYPE_REWARDED_VIDEO
	case mvconst.ADTypeInterstitialVideo:
		return mvconst.CREATIVE_AD_TYPE_INTERSTITIAL_VIDEO
	case mvconst.ADTypeMoreOffer:
		return mvconst.CREATIVE_AD_TYPE_MORE_OFFER
	case mvconst.ADTypeSplash:
		return mvconst.CREATIVE_AD_TYPE_SPLASH
	}
	return 0
}

func getCreativeGroupID(gidList []int) string {
	if len(gidList) <= 0 {
		return ""
	}
	// sort
	sort.Ints(gidList)
	// join
	var gidStrList []string
	for _, v := range gidList {
		gidStrList = append(gidStrList, strconv.Itoa(v))
	}
	gidStr := strings.Join(gidStrList, ",")
	// md5
	return mvutil.Md5(gidStr)
}

// func renderRewardVideoTemplate(ad *Ad, params *mvutil.Params, r *mvutil.RequestParams, campaign *smodel.CampaignInfo) {
// 	if r.Param.AdType != mvconst.ADTypeRewardVideo && r.Param.AdType != mvconst.ADTypeInterstitialVideo {
// 		return
// 	}

// 	//rv endcard_click_result赋值
// 	ad.EndcardClickResult = 1
// 	if ad.CampaignType == 3 && params.TemplateGroupId != 3 && params.TemplateGroupId != 4 {
// 		ad.EndcardClickResult = 2
// 	}

// 	if len(ad.VideoURL) <= 0 || len(ad.VideoResolution) <= 0 {
// 		return
// 	}
// 	// template := getRVTemplate(r, campaign)
// 	// if template.URL == nil || len(*(template.URL)) <= 0 {
// 	// 	return
// 	// }
// 	template, err := renderRvTemplate(r, campaign)
// 	if err != nil {
// 		return
// 	}
// 	handleRVTemplate(ad, params, r, template)
// 	if len(ad.Rv.TemplateUrl) <= 0 {
// 		return
// 	}
// 	params.Extrvtemplate = ad.Rv.VideoTemplate
// 	ad.Rv.TemplateUrl = renderUrlAttachSchema(params.HTTPReq, ad.Rv.TemplateUrl)
// 	// ios rv,iv tpl模板修复
// 	ad.Rv.TemplateUrl = iosIvTemplate(ad.Rv.TemplateUrl, params.Platform, params.AdType)
// 	ad.Rv.PausedUrl = renderUrlAttachSchema(params.HTTPReq, ad.Rv.PausedUrl)
// 	ad.RvPoint = &(ad.Rv)
// }

func renderRewardVideoTemplateNew(ad *Ad, params *mvutil.Params, corsairCampaign corsair_proto.Campaign) {
	if params.AdType != mvconst.ADTypeRewardVideo && params.AdType != mvconst.ADTypeInterstitialVideo {
		return
	}
	//rv endcard_click_result赋值
	ad.EndcardClickResult = 1
	if corsairCampaign.VideoTemplateId != nil {
		ad.Rv.VideoTemplate = int(*(corsairCampaign.VideoTemplateId))
	}
	if len(ad.VideoURL) <= 0 || len(ad.VideoResolution) <= 0 {
		return
	}
	if corsairCampaign.UsageVideo != nil && *(corsairCampaign.UsageVideo) == false {
		return
	}

	var minicardId int
	if corsairCampaign.MiniCardTemplateId != nil {
		minicardId = int(*(corsairCampaign.MiniCardTemplateId))
	}
	if corsairCampaign.Orientation != nil {
		ad.Rv.Orientation = int(*(corsairCampaign.Orientation))
	}
	params.Extrvtemplate = ad.Rv.VideoTemplate
	// 获取视频模板
	templateMapConf, ifFind := extractor.GetTEMPLATE_MAP()
	if ifFind {
		templateIdStr := strconv.Itoa(ad.Rv.VideoTemplate)
		minicardIdStr := strconv.Itoa(minicardId)
		if templateUrl, ok := templateMapConf.Video[templateIdStr]; ok {
			ad.Rv.TemplateUrl = templateUrl
		}
		if minicardUrl, ok := templateMapConf.MiniCard[minicardIdStr]; ok {
			ad.Rv.PausedUrl = minicardUrl
		}
	}
	ad.Rv.TemplateUrl = renderUrlAttachSchema(params.HTTPReq, ad.Rv.TemplateUrl)
	// ios rv,iv tpl模板修复
	ad.Rv.TemplateUrl = iosIvTemplate(ad.Rv.TemplateUrl, params.Platform, params.AdType)
	// 替换api_version<1.4的pl，视频模版的资源域名的宏
	ad.Rv.TemplateUrl = RenderTemplateCreativeDomainMacro(ad.Rv.TemplateUrl)
	ad.Rv.PausedUrl = renderUrlAttachSchema(params.HTTPReq, ad.Rv.PausedUrl)
	ad.RvPoint = &(ad.Rv)
	return
}

func iosIvTemplate(tpl string, platform int, adType int32) string {
	if len(tpl) > 0 && platform == mvconst.PlatformIOS && adType == mvconst.ADTypeInterstitialVideo {
		tpl += "&ad_type=287"
	}
	return tpl
}

// func handleRVTemplate(ad *Ad, params *mvutil.Params, r *mvutil.RequestParams, template *smodel.VideoTemplateUrlItem) {
// 	// 获取defConf
// 	defConf, _ := extractor.GetDefRVTemplate()
// 	if r.Param.ApiVersion >= mvconst.API_VERSION_1_2 {
// 		if len(defConf.PausedURLZip) > 0 {
// 			defConf.PausedURL = defConf.PausedURLZip
// 		}
// 		if len(defConf.URLZip) > 0 {
// 			defConf.URL = defConf.URLZip
// 		}
// 	}
// 	if template.ID != 0 {
// 		ad.Rv.VideoTemplate = int(template.ID)
// 	}
// 	if template.URL != "" {
// 		ad.Rv.TemplateUrl = template.URL
// 	}
// 	ad.Rv.Orientation = handleOrientation(ad.Rv.VideoTemplate)
// 	if defConf.PausedURL != "" {
// 		ad.Rv.PausedUrl = defConf.PausedURL
// 	}
// 	params.IsRVBack = false
// 	// 设备为空时，用默认模板
// 	if mvutil.IsDevidEmpty(params) {
// 		ad.Rv.VideoTemplate = 0
// 		if defConf.URL != "" {
// 			ad.Rv.TemplateUrl = defConf.URL
// 		}
// 		ad.Rv.Orientation = mvconst.ORIENTATION_BOTH
// 		params.IsRVBack = true
// 		return
// 	}
// 	// rv模板回退逻辑
// 	RvIsBack(ad, params, defConf)
// }

// func checkResolution(resolution string) bool {
// 	arr := strings.Split(resolution, "x")
// 	w, _ := strconv.ParseFloat(arr[0], 64)
// 	h, _ := strconv.ParseFloat(arr[1], 64)
// 	b := math.Floor(w / h)
// 	b2 := math.Floor(16 / 9)
// 	return b == b2
// }

// func isOriPortrait(ori int) bool {
// 	if ori == mvconst.ORIENTATION_PORTRAIT || ori == mvconst.ORIENTATION_BOTH {
// 		return true
// 	}
// 	return false
// }

// func isUnitIOS(platform int) bool {
// 	return platform == mvconst.PlatformIOS
// }

// func handleOrientation(tid int) int {
// 	if tid == 3 || tid == 4 {
// 		return mvconst.ORIENTATION_PORTRAIT
// 	}
// 	return mvconst.ORIENTATION_BOTH
// }

// func getRVTemplate(r mvutil.RequestParams, campaign *smodel.CampaignInfo) *smodel.VideoTemplateUrlItem {
// 	template := &smodel.VideoTemplateUrlItem{}
// 	if len(campaign.Endcard) > 0 {
// 		// offer+unit维度
// 		unitId := r.Param.UnitID
// 		unitIdStr := strconv.FormatInt(unitId, 10)
// 		conf := campaign.Endcard[unitIdStr]
// 		template = randRVTemplate(conf, r)
// 		if len(template.URL) > 0 {
// 			return template
// 		}
// 		// offer维度
// 		conf = campaign.Endcard["ALL"]
// 		template = randRVTemplate(conf, r)
// 		if len(template.URL) > 0 {
// 			return template
// 		}
// 	}
// 	// unit维度
// 	if r.UnitInfo.Endcard != nil {
// 		template = randRVTemplate(r.UnitInfo.Endcard, r)
// 		if len(template.URL) > 0 {
// 			return template
// 		}
// 	}
// 	// config维度
// 	conf, ifFind := extractor.GetEndcard()
// 	if !ifFind {
// 		return template
// 	}
// 	template = randRVTemplate(&conf, r)
// 	if len(template.URL) > 0 {
// 		return template
// 	}
// 	return template
// }

// func renderRvTemplate(r *mvutil.RequestParams, campaign *smodel.CampaignInfo) (*smodel.VideoTemplateUrlItem, error) {
// 	if len(campaign.Endcard) > 0 {
// 		confs := campaign.Endcard
// 		if cEndCard, ok := confs[strconv.FormatInt(r.Param.UnitID, 10)]; ok {
// 			rvItem, err := renderRvTemplateItem(cEndCard, &r.Param)
// 			if err == nil {
// 				return rvItem, err
// 			}
// 		}
// 		rvItem, err := renderRvTemplateItem(confs["ALL"], &r.Param)
// 		if err == nil {
// 			return rvItem, err
// 		}
// 	}
// 	//unit
// 	if r.UnitInfo.Endcard != nil {
// 		rvItem, err := renderRvTemplateItem(r.UnitInfo.Endcard, &r.Param)
// 		if err == nil {
// 			return rvItem, err
// 		}
// 	}
// 	//config
// 	conf, ifFind := extractor.GetEndcard()
// 	if !ifFind {
// 		return &smodel.VideoTemplateUrlItem{}, errors.New("get config EndCard error")
// 	}
// 	return renderRvTemplateItem(&conf, &r.Param)
// }

// func renderRvTemplateItem(endcard *smodel.EndCard, param *mvutil.Params) (rvItem *smodel.VideoTemplateUrlItem, err error) {
// 	if endcard == nil || endcard.Status != 1 {
// 		err = errors.New("endcard Item status invalidate")
// 		return
// 	}

// 	if endcard.VideoTemplateUrl == nil || len(endcard.VideoTemplateUrl) == 0 {
// 		err = errors.New("endcard videoTemplateUrls is empty")
// 		return
// 	}

// 	vTempList := endcard.VideoTemplateUrl
// 	weightArr := make(map[int]int, len(vTempList))
// 	sumWeight := 0
// 	var temListKey []int
// 	for _, v := range vTempList {
// 		weight := 1
// 		if v.ID == 0 {
// 			continue
// 		}

// 		id := int(v.ID)
// 		if v.Weight != 0 {
// 			weight = int(v.Weight)
// 		}
// 		weightArr[id] = weight
// 		sumWeight += weight
// 		temListKey = append(temListKey, id)
// 	}
// 	sort.Ints(temListKey)

// 	randVal := mvutil.GetRandConsiderZero(param.GAID, param.IDFA, mvconst.SALT_RV_TEMPLATE, sumWeight)
// 	//randVal := mvutil.GetPureRand(sumWeight) // 使用纯随机代替按设备随机
// 	i := 0
// 	key := 0
// 	for _, v := range temListKey {
// 		i = i + weightArr[v]
// 		if randVal < i {
// 			key = v
// 			break
// 		}
// 	}

// 	rvTmp := &smodel.VideoTemplateUrlItem{}
// 	for _, v := range vTempList {
// 		if int(v.ID) == key {
// 			rvTmp = v
// 			break
// 		}
// 	}
// 	if param.ApiVersion >= mvconst.API_VERSION_1_2 && len(rvTmp.URLZip) > 0 {
// 		rvTmp.URL = rvTmp.URLZip
// 	}
// 	if len(rvTmp.URL) == 0 {
// 		err = errors.New("rv Template Url is empty")
// 		return
// 	}
// 	rvItem = rvTmp
// 	return
// }

// func randRVTemplate(endcard *smodel.EndCard, r mvutil.RequestParams) *smodel.VideoTemplateUrlItem {
// 	template := &smodel.VideoTemplateUrlItem{}
// 	if endcard.Status != 1 {
// 		return template
// 	}
// 	if len(endcard.VideoTemplateUrl) == 0 {
// 		return template
// 	}

// 	list := endcard.VideoTemplateUrl
// 	randArr := make(map[int]int, len(list))
// 	for k, v := range list {
// 		weight := 1
// 		if v.Weight != 0 {
// 			weight = int(v.Weight)
// 		}
// 		randArr[k] = weight
// 	}
// 	key := mvutil.RandArr(randArr, r.Param.GAID, r.Param.IDFA, mvconst.SALT_RV_TEMPLATE)
// 	template = list[key]
// 	if r.Param.ApiVersion >= mvconst.API_VERSION_1_2 && len(template.URLZip) > 0 {
// 		template.URL = template.URLZip
// 	}
// 	return template
// }

// func renderEndcard(ad *Ad, params *mvutil.Params, r *mvutil.RequestParams, campaign *smodel.CampaignInfo) {
// 	// 素材二期，优先使用adserver返回的endcard
// 	if params.EndcardCreativeID > int64(0) && len(ad.EndcardUrl) > 0 {
// 		params.Extendcard = strconv.FormatInt(params.EndcardCreativeID, 10)
// 		ad.EndcardUrl = RenderEndcardUrl(params, ad.EndcardUrl)
// 		return
// 	}
// 	// 版本判断
// 	if !IsReturnEndcard(r) {
// 		return
// 	}
// 	// 获取conf
// 	if len(campaign.Endcard) == 0 {
// 		return
// 	}

// 	endcard := getEndcard(campaign.Endcard, r)
// 	// 新版本使用url_v2，旧版本使用url
// 	endcardFlag := "B"
// 	if params.ApiVersion >= mvconst.API_VERSION_1_1 && len(endcard.UrlV2) > 0 {
// 		endcard.Url = endcard.UrlV2
// 		endcardFlag = "A"
// 	}
// 	// 校验orientation
// 	unitOri := r.Param.FormatOrientation
// 	endcardOri := int(endcard.Orientation)
// 	if len(ad.Rv.TemplateUrl) > 0 {
// 		if !params.IsRVBack {
// 			rvOri := ad.Rv.Orientation
// 			offerOri := mvutil.Max(rvOri, unitOri)
// 			if checkOrientation(offerOri, endcardOri) {
// 				ad.Rv.Orientation = mvutil.Max(offerOri, endcardOri)
// 			} else {
// 				endcard.Url = ""
// 				endcard.ID = int32(0)
// 			}
// 		} else {
// 			if checkOrientation(unitOri, endcardOri) {
// 				ad.Rv.Orientation = mvutil.Max(unitOri, endcardOri)
// 			} else {
// 				ad.Rv.Orientation = mvconst.ORIENTATION_BOTH
// 				endcard.Url = ""
// 				endcard.ID = int32(0)
// 			}
// 		}
// 	} else {
// 		if !checkOrientation(unitOri, endcardOri) {
// 			endcard.Url = ""
// 			endcard.ID = int32(0)
// 		}
// 	}

// 	// 封装
// 	ad.EndcardUrl = RenderEndcardUrl(params, endcard.Url)
// 	params.Extendcard = endcardFlag + strconv.FormatInt(int64(endcard.ID), 10)
// }

func renderEndcardNew(ad *Ad, params *mvutil.Params, corsairCampaign corsair_proto.Campaign) {
	// 获取视频模板
	templateMapConf, _ := extractor.GetTEMPLATE_MAP()
	// 素材三期
	if corsairCampaign.EndCardTemplateId != nil {
		endcardId := int(*(corsairCampaign.EndCardTemplateId))
		// 记录EndCardTemplateId为0的情况。
		if endcardId == 0 {
			params.Extendcard = "0"
		}
		endcardIdStr := strconv.Itoa(endcardId)
		if endcardUrl, ok := templateMapConf.EndScreen[endcardIdStr]; ok {
			params.Extendcard = endcardIdStr
			ad.EndcardUrl = RenderEndcardUrl(params, endcardUrl)
			return
		}
		// 若获取不到对应url，则认为是定制模板
		if params.EndcardCreativeID > int64(0) && len(ad.EndcardUrl) > 0 {
			params.Extendcard = endcardIdStr
			ad.EndcardUrl = RenderEndcardUrl(params, ad.EndcardUrl)
			return
		}
	} else {
		// 若as没有传EndCardTemplateId，则记录为空字符串
		params.Extendcard = ""
	}
	return
}

// func getEndcard(confs map[string]*smodel.EndCard, r *mvutil.RequestParams) mvutil.EndcardItem {
// 	// unit+campaign维度
// 	unitId := r.Param.UnitID
// 	unitIdStr := strconv.FormatInt(unitId, 10)
// 	conf := confs[unitIdStr]
// 	endcard := RandEndcard(conf, r)
// 	if len(endcard.Url) > 0 {
// 		return endcard
// 	}
// 	// campaign维度
// 	conf = confs["ALL"]
// 	endcard = RandEndcard(conf, r)
// 	return endcard
// }

func RandEndcard(conf *smodel.EndCard, r *mvutil.RequestParams) mvutil.EndcardItem {
	var endcard mvutil.EndcardItem
	if conf == nil || conf.Status != 1 {
		return endcard
	}

	if len(conf.Urls) == 0 {
		return endcard
	}

	randArr := make(map[int]int)
	for k, v := range conf.Urls {
		weight := 1
		if v.Weight != 0 {
			weight = int(v.Weight)
		}
		randArr[k] = weight
	}

	key := mvutil.RandArr(randArr, r.Param.GAID, r.Param.IDFA, mvconst.SALT_ENDCARD)
	url := conf.Urls[key]
	if url.Url != "" {
		endcard.Url = url.Url
	}

	if url.Id != 0 {
		endcard.ID = url.Id
	}

	if conf.Orientation != 0 {
		endcard.Orientation = conf.Orientation
	}

	if url.UrlV2 != "" {
		endcard.UrlV2 = url.UrlV2
	}

	if conf.EndcardProtocal != 0 {
		endcard.EndcardProtocal = conf.EndcardProtocal
	}

	if len(conf.EndcardRate) != 0 {
		endcard.EndcardRate = conf.EndcardRate
	}
	return endcard
}

func checkOrientation(offerOri, unitOri int) bool {
	if offerOri == mvconst.ORIENTATION_BOTH || unitOri == mvconst.ORIENTATION_BOTH {
		return true
	}
	return offerOri == unitOri
}

func renderExtAttr(campaign *smodel.CampaignInfo) int {
	attr := 0
	if mvutil.IsBrandOffer(campaign) {
		attr = attr + mvutil.ATTR_BRAND_OFFER
	}
	if mvutil.IsVTA(campaign) {
		attr = attr + mvutil.ATTR_VTA_OFFER
	}
	if mvutil.IsCityOffer(campaign) {
		attr = attr + mvutil.ATTR_CITY_OFFER
	}
	return attr
}

func renderExtPlayable(p *mvutil.Params, c *smodel.CampaignInfo) string {
	Campaign := strconv.FormatInt(c.CampaignId, 10)
	ExtPlayable := strconv.Itoa(p.Extplayable)
	RequestExtPlayable := Campaign + ":" + ExtPlayable
	return RequestExtPlayable
}

// func renderExtTag(p *mvutil.Params, c *smodel.CampaignInfo) string {
// 	Campaign := strconv.FormatInt(c.CampaignId, 10)
// 	ExtCDNAbTest := strconv.Itoa(p.ExtCDNAbTest)
// 	RequestExtTag := Campaign + ":" + ExtCDNAbTest
// 	return RequestExtTag
// }

func renderUa(ad *Ad, c *smodel.CampaignInfo, p *mvutil.Params) {
	// 因ad_server没有传imp_ua和c_ua，所以现默认返回1 （1=webview ua 2=system ua）
	ad.ImpUa = 1
	ad.CUA = 1
	UAAbTest(ad, p, c)
	// 针对appsflyer，imp_ua和c_ua均使用默认的ua（system ua）
	if strings.ToLower(c.ThirdParty) == mvconst.THIRD_PARTY_APPSFLYER {
		ad.ImpUa = 2
		ad.CUA = 2
		RenderAppsflyerUaABTest(ad, p)
		// 记录abtest标记。展示点击的值都一样
		p.ExtDataInit.AppsflyerUaABTestTag = ad.CUA
	}
	p.CUA = ad.CUA
}

func RenderAppsflyerUaABTest(ad *Ad, params *mvutil.Params) {
	randVal := int(crc32.ChecksumIEEE([]byte(mvconst.SALT_APPSFLYER_UA_ABTEST+"_"+mvutil.GetGlobalUniqDeviceTag(params))) % 100) // 64 bit sys
	adnConf, _ := extractor.GetADNET_SWITCHS()
	if appsflyerUaABTestRate, ok := adnConf["appsflyerUaABTestRate"]; ok && appsflyerUaABTestRate > randVal {
		ad.ImpUa = 1
		ad.CUA = 1
	}
}

func canUAAbTest(param *mvutil.Params, c *smodel.CampaignInfo) bool {
	campaignConfs, _ := extractor.GetUA_AB_TEST_CAMPAIGN_CONFIG()
	thirdPartyConfs, _ := extractor.GetUA_AB_TEST_THIRD_PARTY_CONFIG()
	// 优先判断campaign维度没有配置且三方维度没有配置则不做实验
	if !mvutil.InInt64Arr(c.CampaignId, campaignConfs) {
		// 目前只针对appflyer
		if !mvutil.InStrArray(strings.ToLower(c.ThirdParty), thirdPartyConfs) {
			return false
		}
	}

	//device id为空的情况，不做实验
	if param.Platform == mvconst.PlatformAndroid && len(param.GAID) <= 0 && len(param.AndroidID) <= 0 {
		return false
	}
	if param.Platform == mvconst.PlatformIOS && len(param.IDFA) <= 0 {
		return false
	}
	// 处理sdkversion字符串
	sdkData := supply_mvutil.RenderSDKVersion(param.SDKVersion)
	if sdkData.SDKType == "mp" {
		return false
	}
	sdkVersionNum := sdkData.SDKNumber
	var campareValue int
	// 安卓需大于等于8.7.4，ios需大于等于3.3.5才能做abtest
	sdkOs, ifFind := extractor.GetUA_AB_TEST_SDK_OS_CONFIG()
	platformStr := strconv.Itoa(param.Platform)
	if ifFind && len(sdkOs[platformStr]) > 0 {
		campareValue, _ = mvutil.VerCampare(sdkVersionNum, sdkOs[platformStr])
	}
	//if param.Platform == 1 {
	//	campareValue, _ = mvutil.VerCampare(sdkVersionNum, "8.7.4")
	//} else if param.Platform == 2 {
	//	campareValue, _ = mvutil.VerCampare(sdkVersionNum, "3.3.5")
	//}
	if campareValue == 1 || campareValue == 0 {
		return true
	}
	return false
}

func UAAbTest(ad *Ad, p *mvutil.Params, c *smodel.CampaignInfo) {
	p.ExtappearUa = 0
	if canUAAbTest(p, c) {
		testRand, ifFind := extractor.GetUA_AB_TEST_CONFIG()
		if ifFind && len(testRand) > 0 {
			testRandList := make(map[int]int)
			for k, v := range testRand {
				kInt, err := strconv.Atoi(k)
				if err != nil {
					continue
				}
				testRandList[kInt] = v
			}
			if len(testRandList) > 0 {
				testRes := mvutil.RandByRate(testRandList)
				if testRes == 0 {
					return
				}
				switch testRes {
				case 1:
					ad.ImpUa = testRes
					ad.CUA = testRes

				case 2:
					ad.ImpUa = testRes
					ad.CUA = testRes

				case 3:
					ad.ImpUa = 1
					ad.CUA = 1
					p.ExtdeleteDevid = mvconst.DELETE_DIVICEID_BUT_NOT_IMPRESSION
				case 4:
					ad.ImpUa = 2
					ad.CUA = 2
					p.ExtdeleteDevid = mvconst.DELETE_DIVICEID_BUT_NOT_IMPRESSION
				}
				p.ExtappearUa = testRes
			}
		}
	}
}

func cleanEmptyDeviceIdABTest(params *mvutil.Params, campaign *smodel.CampaignInfo) {
	if params.Scenario != mvconst.SCENARIO_OPENAPI || params.RequestType != mvconst.REQUEST_TYPE_OPENAPI_V3 {
		return
	}

	abtest, _ := extractor.GetCleanEmptyDeviceABTest()
	if !abtest.Status {
		return
	}

	if !canDeviceABTest(params) {
		return
	}

	if len(abtest.BlackThirdParty) > 0 && mvutil.InStrArray(campaign.ThirdParty, abtest.BlackThirdParty) {
		// 某些三方不做切量
		return
	}

	if len(abtest.WhiteThirdParty) > 0 && !mvutil.InStrArray(campaign.ThirdParty, abtest.WhiteThirdParty) {
		// 限制三方切量
		return
	}

	if thirdPartyInject, ok := abtest.ThirdPartyInjects[campaign.ThirdParty]; ok && thirdPartyInject.CleanDevId {
		// 按三方配置清空空设备ID
		params.CleanDeviceTest = 2
		return
	}

	if abtest.Rate == 0 {
		return
	}

	if rand.Intn(10000) < abtest.Rate {
		params.CleanDeviceTest = 2
		if thirdPartyInject, ok := abtest.ThirdPartyInjects[campaign.ThirdParty]; ok {
			params.ThirdPartyInjectParams = thirdPartyInject.InjectParams
			if rand.Intn(10000) < thirdPartyInject.SubRate["7"] {
				params.CleanDeviceTest = 7
			} else {
				params.CleanDeviceTest = 6
			}
		}
	} else {
		params.CleanDeviceTest = 1
		if thirdPartyInject, ok := abtest.ThirdPartyInjects[campaign.ThirdParty]; ok {
			params.ThirdPartyInjectParams = thirdPartyInject.InjectParams
			if rand.Intn(10000) < thirdPartyInject.SubRate["5"] {
				params.CleanDeviceTest = 5
			} else {
				params.CleanDeviceTest = 4
			}
		}
	}
	return
}

func canDeviceABTest(params *mvutil.Params) bool {
	if !mvutil.IsDevidEmpty(params) {
		return false
	}

	if len(params.GAID) > 0 || len(params.IDFA) > 0 {
		return true
	}
	return false
}

// 处理mv=>mv 映射值
func RenderExtMpNormalMap(r *mvutil.RequestParams) string {
	MpNormalMap := ""
	if r.UnitInfo.MVToMP == nil {
		return MpNormalMap
	}
	mvmp := r.UnitInfo.MVToMP
	uStr := strconv.FormatInt(mvmp.UnitId, 10)
	aStr := strconv.FormatInt(mvmp.AppId, 10)
	pStr := strconv.FormatInt(mvmp.PublisherId, 10)
	if mvmp.UnitId > 0 && mvmp.AppId > 0 && mvmp.PublisherId > 0 {
		MpNormalMap = uStr + "," + aStr + "," + pStr
	}
	return MpNormalMap
}

// getPriceOut 获取price out
func getPriceOut(params *mvutil.Params, campaign *smodel.CampaignInfo) float64 {
	countryCodeUpper := strings.ToUpper(params.CountryCode)
	key := fmt.Sprintf("%s_%d", countryCodeUpper, params.AppID)
	if price, ok := campaign.MCountryChanlPrice[key]; ok && price > 0 {
		return price
	}

	if price, ok := campaign.MCountryChanlPrice[countryCodeUpper]; ok && price > 0 {
		return price
	}

	if params.PriceOut > 0 {
		return params.PriceOut
	}

	if campaign.Price != 0 {
		return campaign.Price
	}

	return 0
}

// getPriceIn 获取price in
func getPriceIn(params *mvutil.Params, campaign *smodel.CampaignInfo) float64 {
	countryCodeUpper := strings.ToUpper(params.CountryCode)
	subID := getSubId(params)
	key := fmt.Sprintf("%s_%d", countryCodeUpper, subID)
	if price, ok := campaign.MCountryChanlOriPrice[key]; ok && price > 0 {
		return price
	}

	if price, ok := campaign.MCountryChanlOriPrice[countryCodeUpper]; ok && price > 0 {
		return price
	}

	if params.PriceIn > 0 {
		return params.PriceIn
	}

	if campaign.OriPrice != 0 {
		return campaign.OriPrice
	}
	return 0
}

func getChannelPriceIn(params *mvutil.Params, campaign *smodel.CampaignInfo) {
	// 处理 price in
	params.PriceIn = getPriceIn(params, campaign)
}

// 获取渠道维度单子出价
func getChannelPriceOut(params *mvutil.Params, campaign *smodel.CampaignInfo) {
	if params.PublisherType != mvconst.PublisherTypeM {
		return
	}

	// 新逻辑 price out
	params.PriceOut = getPriceOut(params, campaign)
}

func renderMpOpenType(params mvutil.Params, openType int32) int32 {
	// 逻辑只针对mp流量
	if !mvutil.IsMpad(params.RequestPath) {
		return openType
	}
	// 仅针对apk单子
	if openType != int32(3) {
		return openType
	}
	sdkData := supply_mvutil.RenderSDKVersion(params.SDKVersion)
	// sdk_version < 4.1.0 openType强制改为4
	campareValue, _ := mvutil.VerCampare(sdkData.SDKNumber, "4.1.0")
	if campareValue == -1 {
		openType = int32(4)
	}
	return openType
}

func renderAWParamRs(qr mvutil.QueryR, gInfoForAppwall []string, campaignId int64) mvutil.QueryR {
	qr.Glist = strings.Join(gInfoForAppwall, "|")
	qr.Cid = campaignId
	return qr
}

func renderExtCreativeNew(ad *Ad, params *mvutil.Params, r *mvutil.RequestParams, corsairCampaign corsair_proto.Campaign) {
	var extCreativeNew mvutil.ExtCreativeNew
	extCreativeNew.PlayWithoutVideo = ad.PlayableAdsWithoutVideo
	extCreativeNew.VideoEndType = ad.VideoEndType
	if corsairCampaign.TemplateGroup != nil {
		tGroup := int(*corsairCampaign.TemplateGroup)
		extCreativeNew.TemplateGroupId = &tGroup
		params.TemplateGroupId = tGroup
	}
	extCreativeNew.EndScreenId = r.Param.Extendcard
	str, err := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(extCreativeNew)
	if err == nil {
		params.ExtCreativeNew = string(str)
	}
	return
}

func recordPlayableTag(corsairCampaign corsair_proto.Campaign, params *mvutil.Params) {
	if params.AdType != mvconst.ADTypeRewardVideo && params.AdType != mvconst.ADTypeInterstitialVideo {
		return
	}
	if corsairCampaign.ExtPlayable == nil || *corsairCampaign.ExtPlayable == 0 {
		return
	}
	extPlayable := 10000
	if !corsairCampaign.Playable {
		extPlayable = 20000
	}
	switch *corsairCampaign.ExtPlayable {
	case int32(1):
		extPlayable = extPlayable + 100
	case int32(2):
		extPlayable = extPlayable + 200
	case int32(3):
		extPlayable = extPlayable + 300
	}
	if corsairCampaign.VideoEndTypeAs != nil {
		extPlayable = extPlayable + int(*corsairCampaign.VideoEndTypeAs)
	}
	params.Extplayable2 = extPlayable
}

func renderAlgoPrice(params *mvutil.Params, corsairCampaign *corsair_proto.Campaign) {
	if corsairCampaign.UseAlgoPrice != nil {
		params.UseAlgoPrice = *corsairCampaign.UseAlgoPrice
	}

	if corsairCampaign.AlgoPriceIn != nil {
		params.AlgoPriceIn = *corsairCampaign.AlgoPriceIn
	}

	if corsairCampaign.AlgoPriceOut != nil {
		params.AlgoPriceOut = *corsairCampaign.AlgoPriceOut
	}
}

func IsRvNvNewCreative(flag bool, adType int32) bool {
	if flag && (adType == mvconst.ADTypeRewardVideo || adType == mvconst.ADTypeInterstitialVideo) {
		return true
	}
	return false
}

func renderPlct(ad *Ad, params *mvutil.Params, backendId int, dspId int64) {
	adSourceConf, _ := extractor.GetCONFIG_OFFER_PLCT()
	if len(adSourceConf) == 0 {
		return
	}
	backendIdStr := strconv.Itoa(backendId)
	dspIdStr := strconv.FormatInt(dspId, 10)
	// 配置的key为广告来源id和dspid和ruleType的组合，使用_分割。ruleType：1表示ALL,2表示specific app
	// 优先specific
	if specificConf, ok := adSourceConf[backendIdStr+"_"+dspIdStr+"_2"]; ok {
		for _, v := range specificConf {
			if mvutil.InInt64Arr(params.AppID, v.AppIds) {
				ad.Plct = v.Plct
				ad.Plctb = v.Plctb
				return
			}
		}
	}
	// All维度配置
	if allConf, ok := adSourceConf[backendIdStr+"_"+dspIdStr+"_1"]; ok {
		for _, v := range allConf {
			ad.Plct = v.Plct
			ad.Plctb = v.Plctb
		}
	}
}

func cctAbtest(ad *Ad, params *mvutil.Params, campaign *smodel.CampaignInfo) (int, bool) {
	if !mvutil.IsHbOrV3OrV5Request(params.RequestPath) {
		return 0, false
	}
	// 存在设备信息才做abtest，否则不做实验
	if len(params.GAID) > 0 || len(params.IDFA) > 0 || len(params.AndroidID) > 0 {
		cctConf, _ := extractor.GetCCT_ABTEST_CONF()
		// 对于国内流量，若无gaid，则使用Androidid随机
		gaidOrAndroid := mvutil.GetGaidOrAndroid(params.GAID, params.AndroidID)
		// 随机值
		randVal := mvutil.GetRandConsiderZero(gaidOrAndroid, params.IDFA, mvconst.SALT_CCTABTEST, 100)
		// 优先campaign维度
		camIdStr := strconv.FormatInt(ad.CampaignID, 10)
		if camConf, ok := cctConf[camIdStr]; ok {
			if len(camConf) > 0 {
				cctVal, randOK := mvutil.RandByRateInMap(camConf, randVal)
				if randOK {
					return cctVal, true
				}
			}
		}
		if thirdPartyConf, ok := cctConf[campaign.ThirdParty]; ok {
			if len(thirdPartyConf) > 0 {
				cctVal, randOK := mvutil.RandByRateInMap(thirdPartyConf, randVal)
				if randOK {
					return cctVal, true
				}
			}
		}
	}
	return 0, false
}

func renderExtData(params *mvutil.Params) string {
	if params.UseAlgoPrice {
		params.ExtDataInit.UseAlgoPrice = params.UseAlgoPrice
		params.ExtDataInit.AlgoPriceIn = params.AlgoPriceIn
		params.ExtDataInit.AlgoPriceOut = params.AlgoPriceOut
		params.ExtDataInit.PriceIn = params.PriceIn
		params.ExtDataInit.PriceOut = params.PriceOut
	}

	params.ExtDataInit.RwPlus = params.RwPlus
	params.ExtDataInit.SDKOpen = params.Open

	extDataJson, err := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(params.ExtDataInit)
	if err == nil {
		extDataStr := string(extDataJson)
		return extDataStr
	}
	return ""
}

func renderEndcardProperty(ad *Ad, campaign *smodel.CampaignInfo, params *mvutil.Params, r *mvutil.RequestParams) {
	// endcard_url为空值则不处理
	if len(ad.EndcardUrl) == 0 && len(params.IAPlayableUrl) == 0 {
		return
	}
	// 仅针对iv，rv，interactive流量
	if params.AdType != mvconst.ADTypeInterstitialVideo && params.AdType != mvconst.ADTypeRewardVideo &&
		params.AdType != mvconst.ADTypeInteractive {
		return
	}

	endcardParams := ""
	// 判断是否允许自动出storekit
	if judgeByUnitAndOffer(r.UnitInfo.Unit.Alac, campaign.AlacRate) {
		canAlac := true
		prisonConf, _ := extractor.GetALAC_PRISON_CONFIG()
		if pubConf, ok := prisonConf[strconv.FormatInt(params.PublisherID, 10)]; ok {
			if mvutil.InArray(params.TemplateGroupId, pubConf) {
				canAlac = false
			}
		}
		if canAlac {
			endcardParams = endcardParams + "&alac=1"
			params.ExtDataInit.Alac = 1
		}
	} else {
		endcardParams += renderClsdly(params)
	}
	// 判断是否允许endcard全局可点
	if judgeByUnitAndOffer(r.UnitInfo.Unit.Alecfc, campaign.AlecfcRate) {
		endcardParams = endcardParams + "&alecfc=1"
		params.ExtDataInit.Alecfc = 1
	}
	// 判断是否允许展示更多offer
	// unit维度默认开启
	mof := 1
	// offer维度默认开启
	oMof := 1
	if r.UnitInfo.Unit.Mof != nil {
		mof = *r.UnitInfo.Unit.Mof
	}
	if campaign.Mof != 0 {
		oMof = campaign.Mof
	}
	// unit及offer均为开启则允许展示更多offer
	if mof == 1 && oMof == 1 {
		endcardParams = endcardParams + "&mof=1"
		params.ExtDataInit.Mof = 1
		if r.Param.MofAbFlag {
			endcardParams += "&mof_ab=1"
		}
	}
	// 下发more offer或关闭场景广告才下发主offer信息
	if (mof == 1 && oMof == 1) || len(params.ExtDataInit.CloseAdTag) > 0 {
		if r.Param.NewMoreOfferFlag {
			endcardParams += "&mof_uid=" + strconv.FormatInt(r.UnitInfo.Unit.MofUnitId, 10)
		}
		// 改为只要出more offer，就下发主offer相关信息，不区分固定unit和拆分unit
		r.Param.NewMoreOfferParams = addNewMofParams(ad, params, r)
		endcardParams += "&ec_id=" + params.Extendcard + r.Param.NewMoreOfferParams
	}

	if r.Param.NewMofImpFlag {
		endcardParams += "&n_imp=1"
	}
	// 关闭场景广告参数下发
	if len(params.ExtDataInit.CloseAdTag) > 0 {
		endcardParams += "&clsad=" + params.ExtDataInit.CloseAdTag
	}
	// 返回流量对应的包名
	if len(params.PackageName) > 0 {
		endcardParams += "&mof_pkg=" + params.PackageName
	}

	// 下发自有id
	if len(params.ExtSysId) > 0 && params.ExtSysId != "," && params.ApiVersion < mvconst.API_VERSION_2_2 {
		endcardParams += "&n_bc=" + url.QueryEscape(params.ExtSysId)
	}

	// 下发region，pl 没有分region，需要adnet和pioneer返回
	if len(params.ExtcdnType) > 0 {
		endcardParams += "&n_region=" + params.ExtcdnType
	}

	if params.AdType == mvconst.ADTypeInteractive {
		params.IAPlayableUrl = params.IAPlayableUrl + endcardParams
	} else {
		ad.EndcardUrl = ad.EndcardUrl + endcardParams
	}
	// endcard_url 下发参数abtest
	renderEndcardUrlParamABTest(ad, campaign, params, mvconst.FakeAdserverDsp)

	// 替换api_version<1.4的pl，视频模版的资源域名的宏
	ad.EndcardUrl = RenderTemplateCreativeDomainMacro(ad.EndcardUrl)
}

func renderClsdly(params *mvutil.Params) string {
	adnConf, _ := extractor.GetADNET_SWITCHS()
	adnConfList := extractor.GetADNET_CONF_LIST()
	var clsdlyStr string
	if clsdlyPublisherBlackList, ok := adnConfList["clsdlyPubBList"]; ok && len(clsdlyPublisherBlackList) > 0 {
		if mvutil.InInt64Arr(params.PublisherID, clsdlyPublisherBlackList) {
			return clsdlyStr
		}
	}
	clsdlyTime := 3
	if clsdly, ok := adnConf["clsdly"]; ok {
		clsdlyTime = clsdly
	}
	clsdlyStr = "&clsdly=" + strconv.Itoa(clsdlyTime)
	return clsdlyStr
}

func addNewMofParams(ad *Ad, params *mvutil.Params, r *mvutil.RequestParams) string {
	var str string
	str = "&rv_tid=" + strconv.Itoa(ad.Rv.VideoTemplate) + "&tplgp=" + strconv.Itoa(r.Param.TplGroupId) +
		"&v_fmd5=" + params.VideoFmd5 + "&i_fmd5=" + params.ImgFMD5 + "&mcc=" + r.Param.MCC + "&mnc=" + r.Param.MNC
	return str
}

func judgeByUnitAndOffer(unitSwitch *int, rate int) bool {
	// unit维度默认开启允许H5自动出storekit和允许endcard全局可点和允许展示更多offer
	uswitch := 1
	if unitSwitch != nil {
		uswitch = *unitSwitch
	}
	if uswitch != 1 {
		return false
	}
	// offer维度配置0则不下发
	if rate == 0 {
		return false
	}
	randV := rand.Intn(100)
	if rate > randV {
		return true
	}
	return false
}

func isAdClickReplaceByClickUrl(params *mvutil.Params, campaign *smodel.CampaignInfo, ad *Ad) {
	if params.Scenario != mvconst.SCENARIO_OPENAPI || !mvutil.IsHbOrV3OrV5Request(params.RequestPath) {
		return
	}
	if params.AdType != mvconst.ADTypeRewardVideo && params.AdType != mvconst.ADTypeInterstitialVideo &&
		params.AdType != mvconst.ADTypeNative {
		return
	}
	if !iosSdkSupportAdClick(params.FormatSDKVersion, params.AdType, params.Platform) {
		return
	}
	adnConf, _ := extractor.GetADNET_SWITCHS()

	plSwitch, plOk := adnConf["closePlSwitch"]
	if plOk && plSwitch != 1 {
		if isPlayable(params.TemplateGroupId) {
			return
		}
	}
	advId := int(campaign.AdvertiserId)

	if adClickSwitch, ok := adnConf["rpAdClick"]; ok {
		if adClickSwitch == 1 {
			if advId == 903 && strings.ToLower(campaign.ThirdParty) == mvconst.THIRD_PARTY_APPSFLYER &&
				(params.Extra10 == mvconst.JUMP_TYPE_CLIENT_DO_ALL_PING_SERVER ||
					mvutil.IsSSOfferAndCM6(campaign, params)) {
				if mvutil.IsDevidEmpty(params) {
					// 表示无idfa的情况
					params.ExtDataInit.ClickWithUaLangTag = 1
					params.AdClickReplaceclickByClickTag = true
					// af不参与下面的clickInAdTrackingTest实验
					return
				}
				// af 安卓有设备id的流量保留实验
				if params.Platform == mvconst.PlatformAndroid {
					hasDeviceIdABTest(adnConf, params)
					return
				}
			}
		}
	}
	// 放开903和clickmode的开关，若有问题则不放开。默认都放开
	if advAndCMSwitch, ok := adnConf["advAndCMSwitch"]; ok && advAndCMSwitch == 1 {
		// plan c。使用adtracking.click跳转。
		if advId == 903 && (params.Extra10 == mvconst.JUMP_TYPE_CLIENT_DO_ALL_PING_SERVER || mvutil.IsSSOfferAndCM6(campaign, params)) &&
			params.Platform == mvconst.PlatformIOS && ad.CampaignType == 1 && len(ad.PackageName) > 0 && ad.CType != mvconst.PAY_TYPE_CPC {
			clickInAdTrackingTest(params, campaign.ThirdParty)
		}
	} else if advId == 903 && params.Platform == mvconst.PlatformIOS && ad.CampaignType == 1 && len(ad.PackageName) > 0 &&
		(params.Extra10 == mvconst.JUMP_TYPE_CLIENT_DO_ALL || params.Extra10 == mvconst.JUMP_TYPE_CLIENT_DO_ALL_PING_SERVER) && ad.CType != mvconst.PAY_TYPE_CPC {
		// af 为6的情况不能做，跳到af时会有landing page，所以af为6不能切
		if params.Extra10 == mvconst.JUMP_TYPE_CLIENT_DO_ALL && strings.ToLower(campaign.ThirdParty) == mvconst.THIRD_PARTY_APPSFLYER &&
			!mvutil.IsSSOfferAndCM6(campaign, params) {
			return
		}
		clickInAdTrackingTest(params, campaign.ThirdParty)
	}

}

func iosSdkSupportAdClick(sdkVerion supply_mvutil.SDKVersionItem, adType int32, platform int) bool {
	if sdkVerion.SDKType != "mi" && sdkVerion.SDKType != "mal" {
		return false
	}
	//iOS 排除 390和 391
	// 因390及以上ios sdk版本支持返回campaign信息给到h5，因此可以放开这两个下毒版本。
	//if sdkVerion.SDKVersionCode == mvconst.AdTrackingIOSExpV1 || sdkVerion.SDKVersionCode == mvconst.AdTrackingIOSExpV2 {
	//	return false
	//}
	if platform == mvconst.PlatformIOS {
		// ios rv 展示和点击 2.8.0 及以上版本
		if adType == mvconst.ADTypeRewardVideo && sdkVerion.SDKVersionCode < mvconst.AdTrackingIOSRV {
			return false
		}
		// ios iv 展示和点击 3.6.0及以上版本
		if adType == mvconst.ADTypeInterstitialVideo && sdkVerion.SDKVersionCode < mvconst.AdTrackingIOSIV {
			return false
		}
		// 针对nv
		if adType == mvconst.ADTypeNative && sdkVerion.SDKVersionCode < mvconst.AdTrackingIOSNativeClick {
			return false
		}
	} else {
		if sdkVerion.SDKVersionCode < mvconst.AdTrackingAndroidNativePicClick {
			return false
		}
	}

	return true
}

func clickInAdTrackingTest(params *mvutil.Params, thirdParty string) {
	adnConf, _ := extractor.GetADNET_SWITCHS()
	thirdPartyConf, _ := extractor.GetCLICK_IN_ADTRACKING_THIRD_PARTY_CONFIG()
	clickInAdTrackingRate, ok := adnConf["ciatRate"]
	// 开关都关闭则不做处理
	if len(thirdPartyConf) == 0 && !ok {
		return
	}
	finalRate := clickInAdTrackingRate
	if thirdPartyRate, ok := thirdPartyConf[thirdParty]; ok {
		finalRate = thirdPartyRate
	}
	rateRand := rand.Intn(100)
	if finalRate > rateRand {
		params.ExtDataInit.ClickInAdTracking = 1
		params.AdClickReplaceclickByClickTag = true
		return
	}
	params.ExtDataInit.ClickInAdTracking = 2
}

func hasDeviceIdABTest(conf map[string]int, params *mvutil.Params) {
	// ios 有idfa情况abtest
	if iosClickWithUaLangRate, ok := conf["iosCwulTest"]; ok && params.Platform == mvconst.PlatformIOS {
		rateRand := rand.Intn(100)
		if iosClickWithUaLangRate > rateRand {
			params.ExtDataInit.ClickWithUaLangTag = 2
			params.AdClickReplaceclickByClickTag = true
			return
		}
	}
	// 安卓由gaid情况下abtest
	if androidClickWithUaLangRate, ok := conf["adCwulTest"]; ok && params.Platform == mvconst.PlatformAndroid {
		rateRand := rand.Intn(100)
		if androidClickWithUaLangRate > rateRand {
			params.ExtDataInit.ClickWithUaLangTag = 2
			params.AdClickReplaceclickByClickTag = true
			return
		}
	}
	params.ExtDataInit.ClickWithUaLangTag = 3
}

func isH5ToClick(params *mvutil.Params) {
	// 只针对原来的使用click_url赋值给adtracking.click，并把click_url置空的流量生效
	if !params.AdClickReplaceclickByClickTag {
		return
	}
	if params.NeedWebviewPrison {
		return
	}
	// 只针对iv，rv流量
	if params.AdType != mvconst.ADTypeRewardVideo && params.AdType != mvconst.ADTypeInterstitialVideo {
		return
	}
	// 判断支持h5发点击上报的版本
	if !enableSdkVersion(params) {
		return
	}
	adnConf, _ := extractor.GetADNET_SWITCHS()
	// 安卓开关
	androidSwitch, adOk := adnConf["adClickSwitch"]
	if !adOk || androidSwitch != 1 && params.Platform == mvconst.PlatformAndroid {
		return
	}
	// ios开关
	iosSwitch, iosOk := adnConf["iosClickSwitch"]
	if !iosOk || iosSwitch != 1 && params.Platform == mvconst.PlatformIOS {
		return
	}

	// endcard控制开关
	if endcardSwitchRate, ok := adnConf["ecClickSwitch"]; ok && params.TemplateGroupId != 4 {
		rateRand := rand.Intn(100)
		if endcardSwitchRate > rateRand {
			params.ExtDataInit.H5Handle = 1
		}
	}
	// playable控制开关
	if playableSwitchRate, ok := adnConf["plClickSwitch"]; ok && params.TemplateGroupId == 4 {
		rateRand := rand.Intn(100)
		if playableSwitchRate > rateRand {
			params.ExtDataInit.H5Handle = 1
		}
	}

}

func enableSdkVersion(params *mvutil.Params) bool {
	if params.Platform == mvconst.PlatformIOS {
		if params.FormatSDKVersion.SDKVersionCode >= mvconst.AdTrackingIOSExpV1 {
			return true
		}
	} else {
		if params.FormatSDKVersion.SDKVersionCode >= mvconst.AdTrackingAndroidNativePicClick {
			return true
		}
	}
	return false
}

func adClickBack(params *mvutil.Params, ad *Ad) {
	// AdClickReplaceclickByClick 回退逻辑
	if params.ExtDataInit.H5Handle != 1 && params.Platform == mvconst.PlatformIOS && (params.FormatSDKVersion.SDKVersionCode == mvconst.AdTrackingIOSExpV1 ||
		params.FormatSDKVersion.SDKVersionCode == mvconst.AdTrackingIOSExpV2) && params.AdClickReplaceclickByClickTag == true {
		params.AdClickReplaceclickByClickTag = false
	}
	//  ios 限定link_type=1且存在包名才能放置到ad_tracking
	if params.Platform == mvconst.PlatformIOS && params.AdClickReplaceclickByClickTag == true && (ad.CampaignType != 1 || len(ad.PackageName) == 0) {
		params.AdClickReplaceclickByClickTag = false
	}
	// 对于安卓不支持h5发的情况，回退到最初的逻辑，不使用click_url给adtracking.click赋值及不将click_url置空
	if params.ExtDataInit.H5Handle != 1 && params.Platform == mvconst.PlatformAndroid && params.AdClickReplaceclickByClickTag == true {
		params.AdClickReplaceclickByClickTag = false
	}
}

func isPlayable(tgroupId int) bool {
	return tgroupId == 3 || tgroupId == 4
}

func isSSNeedNotice3s(abtestTag int) bool {
	return abtestTag == 1 || abtestTag == 2
}

func isSSABTestCreateSourse4(abtestTag int) bool {
	return abtestTag == 2
}

func isDirectUrlEmpty(directUrl string, abtestTag int) bool {
	return len(directUrl) == 0 || isSSABTestCreateSourse4(abtestTag)
}

// 服务端点击设备去重实验，寻找最优的窗口，避免误伤
func serverUniqClickTimeABTest(params *mvutil.Params, campaign *smodel.CampaignInfo) {
	if params.Extra10 != mvconst.JUMP_TYPE_ONLINE_DSP_AJUMP_SERVER && params.Extra10 != mvconst.JUMP_TYPE_TRACKING_PING_THIRDPARTY_CLICK_URL {
		return
	}

	if params.Scenario != mvconst.SCENARIO_OPENAPI || params.RequestType != mvconst.REQUEST_TYPE_OPENAPI_V3 {
		return
	}

	cfg := extractor.GetNoticeClickUniq()
	if !cfg.ServerClickUniq || len(cfg.ServerClickUniqTime) == 0 {
		return
	}

	intervalStr, _ := mvutil.RandByRateInMapInString(cfg.ServerClickUniqTime, int(crc32.ChecksumIEEE([]byte(strings.ToLower(mvutil.GetDeviceTag(params))))%10000))
	interval, err := strconv.Atoi(intervalStr)
	if err != nil {
		return
	}

	params.ExtDataInit.ServerUniqClickTime = interval
	return
}

func demandSideLibABTest(params *mvutil.Params, campaign *smodel.CampaignInfo, r *mvutil.RequestParams) {
	// demand side lib abtest
	if !campaign.IsSSPlatform() {
		return
	}
	renderDemandContext(params, r)
	// tracking cn abtest flag

	return
}

func renderDemandContext(params *mvutil.Params, r *mvutil.RequestParams) {
	ctx := NewDemandContext(params, r)
	demand.RenderContext(ctx, extractor.DemandDao)
	ParseDemandContext4Params(ctx, params)
	params.DemandContext = ctx
}

func jssdkWhiteList(params *mvutil.Params, campaign *smodel.CampaignInfo, ad *Ad) {
	// 限制jssdk流量
	if params.Scenario != mvconst.SCENARIO_OPENAPI || params.RequestPath != mvconst.PATHJssdkApi {
		return
	}

	// clickmode 14无需走此逻辑
	if params.Extra10 == mvconst.JUMP_TYPE_TRACKING_PING_THIRDPARTY_CLICK_URL {
		return
	}

	advId := int(campaign.AdvertiserId)

	if advId != mvconst.NormalAgency {
		return
	}

	if !mvutil.IsDevidEmpty(params) {
		return
	}

	// CN 空设备判断
	if params.Platform == mvconst.PlatformAndroid && params.CountryCode == "CN" && (len(params.AndroidID) > 0 || len(params.IMEI) > 0) {
		return
	}

	cfg := extractor.GetThirdPartyWhiteList()
	if !cfg.Status {
		return
	}

	if len(cfg.PublisherIds) > 0 && !mvutil.InInt64Arr(params.PublisherID, cfg.PublisherIds) {
		return
	}

	if len(cfg.AppIds) > 0 && !mvutil.InInt64Arr(params.AppID, cfg.AppIds) {
		return
	}

	if len(cfg.UnitIds) > 0 && !mvutil.InInt64Arr(params.UnitID, cfg.UnitIds) {
		return
	}

	if len(cfg.Platform) > 0 && !mvutil.InStrArray(strings.ToLower(params.PlatformName), cfg.Platform) {
		return
	}

	if len(cfg.CountryCode) > 0 && !mvutil.InStrArray(strings.ToUpper(params.CountryCode), cfg.CountryCode) {
		return
	}

	thirdParty := strings.ToLower(campaign.ThirdParty)
	thirdPartyCfg, ok := cfg.JSSDKRate[thirdParty]
	if !ok {
		return
	}

	if len(thirdPartyCfg.CampaignIds) > 0 && !mvutil.InInt64Arr(campaign.CampaignId, thirdPartyCfg.CampaignIds) {
		return
	}

	if thirdPartyCfg.Rate > rand.Intn(10000) {
		params.ExtDataInit.ClickInServer = 1
		// 对于此情况，clickmode要改为13，tracking根据此值做跳转，返回给sdk的clickmode不变。
		params.ExtDataInit.OldClickMode = params.Extra10
		params.Extra10 = mvconst.JUMP_TYPE_ONLINE_DSP_AJUMP_SERVER
		return
	}
	params.ExtDataInit.ClickInServer = 2
}

func canUseAfWhiteList(params *mvutil.Params, campaign *smodel.CampaignInfo, ad *Ad) {
	// 限制sdk流量
	if params.Scenario != mvconst.SCENARIO_OPENAPI || !mvutil.IsHbOrV3OrV5Request(params.RequestPath) {
		return
	}
	// clickmode 14无需走此逻辑
	// 0/13也切量,确保配置了0，13的情况下，不会被clickinserver逻辑覆盖。后续会把逻辑迁移到adn lib里面
	if params.Extra10 == mvconst.JUMP_TYPE_TRACKING_PING_THIRDPARTY_CLICK_URL || params.Extra10 == mvconst.JUMP_TYPE_NORMAL ||
		params.Extra10 == mvconst.JUMP_TYPE_ONLINE_DSP_AJUMP_SERVER {
		return
	}
	// 需要ios，设备为空值，offer的advid为903，三方为af。
	advId := int(campaign.AdvertiserId)
	// 只针对903单子
	if advId != 903 {
		return
	}
	// 有device id和无device id分开实验。
	isDeviceEmpty := mvutil.CisIsDevidEmpty(params)

	// 开放有devid，及其他三方进行服务端点击处理。有无devid + platform + 三方 分开切量
	rate := getClickInServerRate(params, campaign, isDeviceEmpty)
	if rate > 0 {
		// 新增枚举值以区分是否有设备id,偏移2位
		offset := 0
		if !isDeviceEmpty {
			offset = 2
		}

		var rateRand int
		if extractor.GetCLICK_IN_SERVER_CONF_NEW().RandDev {
			rateRand = int(crc32.ChecksumIEEE([]byte(mvconst.SALT_CLICK_IN_SERVER__ABTEST+"_"+mvutil.GetGlobalUniqDeviceTag(params))) % 100) // 64 bit sys
		} else {
			rateRand = rand.Intn(100)
		}

		if rate > rateRand {
			params.ExtDataInit.ClickInServer = 1 + offset
			// 记录原有的clickmode
			params.ExtDataInit.OldClickMode = params.Extra10
			// 对于此情况，clickmode要改为13，tracking根据此值做跳转，返回给sdk的clickmode不变。
			params.Extra10 = mvconst.JUMP_TYPE_ONLINE_DSP_AJUMP_SERVER
			// 走白名单通道，因为af带的ua为system ua，导致无法ipua归因
			if strings.ToLower(campaign.ThirdParty) == mvconst.THIRD_PARTY_APPSFLYER {
				ad.CUA = 1
			}
			return
		}
		params.ExtDataInit.ClickInServer = 2 + offset
	}
}

func getIpuaWhiteListRate(params *mvutil.Params, campaign *smodel.CampaignInfo, rateConf []mvutil.ClickInServerRateConf) int {
	region := extractor.GetSYSTEM_AREA()
	for _, v := range rateConf {
		// 排除掉黑名单内的流量
		if mvutil.InStrArray(region, v.BlackRegion) {
			continue
		}
		if mvutil.InStrArray(params.CountryCode, v.BlackCountry) {
			continue
		}
		if mvutil.InStrArray(campaign.ThirdParty, v.BlackThirdParty) {
			continue
		}
		// 白名单内流量才能获取切量rate
		if len(v.Platform) > 0 && !mvutil.InArray(params.Platform, v.Platform) {
			continue
		}
		if len(v.Region) > 0 && !mvutil.InStrArray(region, v.Region) {
			continue
		}
		if len(v.Country) > 0 && !mvutil.InStrArray(params.CountryCode, v.Country) {
			continue
		}
		if len(v.ThirdParty) > 0 && !mvutil.InStrArray(campaign.ThirdParty, v.ThirdParty) {
			continue
		}
		return v.Rate
	}
	return 0
}

func getClickInServerRate(params *mvutil.Params, campaign *smodel.CampaignInfo, isDeviceEmpty bool) int {
	rateConf := extractor.GetCLICK_IN_SERVER_CONF_NEW()
	var finalRate int
	// 先过滤掉不能做服务端点击上报的情况
	if mvutil.InInt64Arr(campaign.CampaignId, rateConf.OfferBList) {
		return finalRate
	}
	if mvutil.InInt64Arr(params.PublisherID, rateConf.PubBList) {
		return finalRate
	}
	if mvutil.InInt64Arr(params.AppID, rateConf.AppBList) {
		return finalRate
	}
	if mvutil.InInt64Arr(params.UnitID, rateConf.UnitBList) {
		return finalRate
	}
	if mvutil.InStrArray(params.Extra10, rateConf.ClickModeBList) {
		return finalRate
	}

	if mvutil.InArray(params.LinkType, rateConf.OpenTypeBList) {
		return finalRate
	}

	// 获取对应配置
	if isDeviceEmpty {
		return getIpuaWhiteListRate(params, campaign, rateConf.NoDevIdConf)
	} else {
		return getIpuaWhiteListRate(params, campaign, rateConf.HasDevIdConf)
	}
	return finalRate
}

// true 切量 false 没有切量或者对照组
// func adjustPostbackABTest(params *mvutil.Params, campaign *smodel.CampaignInfo) (btest bool) {
// 	if !campaign.IsSSPlatform() || strings.ToLower(campaign.ThirdParty) != mvconst.THIRD_PARTY_ADJUST {
// 		return
// 	}

// 	cfg := extractor.GetAdjustPostbackABTest()
// 	if !cfg.Status {
// 		return
// 	}

// 	if len(cfg.Scenario) > 0 && !mvutil.InStrArray(params.Scenario, cfg.Scenario) {
// 		return
// 	}

// 	if len(cfg.BListCampaign) > 0 && mvutil.InInt64Arr(campaign.CampaignId, cfg.BListCampaign) {
// 		return
// 	}

// 	if len(cfg.WListCampaign) > 0 && !mvutil.InInt64Arr(campaign.CampaignId, cfg.WListCampaign) {
// 		return
// 	}

// 	if rand.Intn(10000) < cfg.Rate {
// 		// 切量
// 		params.ExtDataInit.AdjustPostbackABTest = 2
// 		return true
// 	}
// 	params.ExtDataInit.AdjustPostbackABTest = 1
// 	return
// }

func adParamABTest(params *mvutil.Params, campaign *smodel.CampaignInfo, ad *Ad, dspId int64) {
	if !mvutil.IsHbOrV3OrV5Request(params.RequestPath) {
		return
	}
	if params.ABTestTags == nil {
		params.ABTestTags = make(map[string]int)
	}
	value := reflect.ValueOf(ad).Elem()
	abtestFieldsConf, _ := extractor.GetABTEST_FIELDS()
	abtestConfs, _ := extractor.GetABTEST_CONFS()
	// 遍历ad结构体
	for i := 0; i < value.NumField(); i++ {
		// 获取ad下字段对应的变量名
		varName := value.Type().Field(i).Name
		// 在配置的fields里的变量才需要做实验
		if mvutil.InStrArray(varName, abtestFieldsConf.AdParams) {
			// 获取字段对应的配置
			if abtestConf, ok := abtestConfs[varName]; ok {
				// 根据流量条件，获取abtest配置
				conf, randType := GetABTestConf(params, campaign, abtestConf, dspId)
				if len(conf) == 0 {
					continue
				}
				// 获取abtest结果
				finalVal, randOk := GetABTestRes(conf, randType, params, varName)
				if !randOk {
					continue
				}
				v := value.FieldByName(varName)
				// 修改abtest字段的值
				mvutil.SetFieldValueForABTest(v, finalVal)
				// abtest结果记录到日志中
				params.ABTestTags[varName] = finalVal
			}
		}
	}
}

func GetABTestRes(conf map[string]int, randType int32, params *mvutil.Params, varName string) (int, bool) {
	var randVal int
	if randType == mvconst.RANDOM_BY_DEVICE {
		if params.CountryCode == "CN" && len(params.GAID) == 0 && params.Platform == mvconst.PlatformAndroid {
			randVal = mvutil.GetRandConsiderZero(params.ClientIP+params.Model, params.IDFA, varName, 100)
		} else {
			randVal = mvutil.GetRandConsiderZero(params.GAID, params.IDFA, varName, 100)
		}
		// 设备id为空值，则不做abtest
		if randVal == -1 {
			return 0, false
		}
	} else {
		randVal = rand.Intn(100)
	}
	aabTestVal, randOK := mvutil.RandByRateInMap(conf, randVal)
	if !randOK {
		return 0, false
	}
	return aabTestVal, true
}

func GetABTestConf(params *mvutil.Params, campaign *smodel.CampaignInfo, abtestConfs []mvutil.ABTEST_CONF, dspId int64) (map[string]int, int32) {
	var defaultConf map[string]int

	for _, v := range abtestConfs {
		// 判断是否在黑名单内
		if blockByBlackList(v, params, campaign, dspId) {
			return defaultConf, 0
		}

		// 判断是否在白名单内,若存在，则返回切量配置内容
		if isInWhiteList(v, params, campaign, dspId) {
			return v.RateMap, v.RandType
		}
	}

	return defaultConf, 0
}

func blockByBlackList(abtestConf mvutil.ABTEST_CONF, params *mvutil.Params, campaign *smodel.CampaignInfo, dspId int64) bool {
	if mvutil.InInt64Arr(params.UnitID, abtestConf.UnitBList) {
		return true
	}
	if mvutil.InInt64Arr(params.AppID, abtestConf.AppBList) {
		return true
	}
	if mvutil.InInt64Arr(params.PublisherID, abtestConf.PubBList) {
		return true
	}
	if campaign != nil && mvutil.InInt64Arr(campaign.CampaignId, abtestConf.CIdBList) {
		return true
	}
	if mvutil.InInt32Arr(params.FormatSDKVersion.SDKVersionCode, abtestConf.SdkVerBList) {
		return true
	}
	if mvutil.InStrArray(params.CountryCode, abtestConf.CountryBList) {
		return true
	}
	if mvutil.InInt32Arr(params.AdType, abtestConf.AdTypeBList) {
		return true
	}
	if mvutil.InArray(params.Platform, abtestConf.PlatformBList) {
		return true
	}
	if mvutil.InInt32Arr(params.ApiVersionCode, abtestConf.ApiVerBList) {
		return true
	}
	if campaign != nil && mvutil.InStrArray(campaign.ThirdParty, abtestConf.ThirdPartyBList) {
		return true
	}
	if mvutil.InInt32Arr(params.OSVersionCode, abtestConf.OsVerBList) {
		return true
	}
	if mvutil.InInt64Arr(dspId, abtestConf.DspIdBList) {
		return true
	}
	return false
}

func isInWhiteList(abtestConf mvutil.ABTEST_CONF, params *mvutil.Params, campaign *smodel.CampaignInfo, dspId int64) bool {
	if len(abtestConf.UnitWList) > 0 && !mvutil.InInt64Arr(params.UnitID, abtestConf.UnitWList) {
		return false
	}

	if len(abtestConf.AppWList) > 0 && !mvutil.InInt64Arr(params.AppID, abtestConf.AppWList) {
		return false
	}

	if len(abtestConf.PubWList) > 0 && !mvutil.InInt64Arr(params.PublisherID, abtestConf.PubWList) {
		return false
	}

	if campaign != nil && len(abtestConf.CIdWList) > 0 && !mvutil.InInt64Arr(campaign.CampaignId, abtestConf.CIdWList) {
		return false
	}

	if abtestConf.MinSdkVer > 0 && params.FormatSDKVersion.SDKVersionCode < abtestConf.MinSdkVer {
		return false
	}

	if abtestConf.MaxSdkVer > 0 && params.FormatSDKVersion.SDKVersionCode > abtestConf.MaxSdkVer {
		return false
	}

	if len(abtestConf.CountryWList) > 0 && !mvutil.InStrArray(params.CountryCode, abtestConf.CountryWList) {
		return false
	}

	if len(abtestConf.AdTypeWList) > 0 && !mvutil.InInt32Arr(params.AdType, abtestConf.AdTypeWList) {
		return false
	}

	if len(abtestConf.PlatformWList) > 0 && !mvutil.InArray(params.Platform, abtestConf.PlatformWList) {
		return false
	}

	if len(abtestConf.ApiVerWList) > 0 && !mvutil.InInt32Arr(params.ApiVersionCode, abtestConf.ApiVerWList) {
		return false
	}

	if campaign != nil && len(abtestConf.ThirdPartyWList) > 0 && !mvutil.InStrArray(campaign.ThirdParty, abtestConf.ThirdPartyWList) {
		return false
	}

	if abtestConf.MinOsVer > 0 && params.OSVersionCode < abtestConf.MinOsVer {
		return false
	}

	if abtestConf.MaxOsVer > 0 && params.OSVersionCode > abtestConf.MaxOsVer {
		return false
	}

	if len(abtestConf.DspIdWList) > 0 && !mvutil.InInt64Arr(dspId, abtestConf.DspIdWList) {
		return false
	}

	// 如果均没有配置白名单，则判断整体切量概率配置
	if len(abtestConf.UnitWList) == 0 && len(abtestConf.AppWList) == 0 && len(abtestConf.PubWList) == 0 && len(abtestConf.CIdWList) == 0 &&
		abtestConf.MinSdkVer == 0 && abtestConf.MaxSdkVer == 0 && len(abtestConf.CountryWList) == 0 && len(abtestConf.AdTypeWList) == 0 &&
		len(abtestConf.PlatformWList) == 0 && len(abtestConf.ApiVerWList) == 0 && len(abtestConf.ThirdPartyWList) == 0 && abtestConf.MinOsVer == 0 &&
		abtestConf.MaxOsVer == 0 && len(abtestConf.DspIdWList) == 0 {
		if abtestConf.TotalRate > rand.Intn(100) {
			return true
		} else {
			return false
		}
	}

	// 符合条件的流量，根据total rate进行切量
	if abtestConf.TotalRate > rand.Intn(100) {
		return true
	}

	return false
}

func renderExtABTestList(p *mvutil.Params, c *smodel.CampaignInfo) string {
	if len(p.ABTestTagStr) > 0 {
		Campaign := strconv.FormatInt(c.CampaignId, 10)
		return Campaign + ":" + p.ABTestTagStr
	}
	return ""
}

func renderEndcardUrlParamABTest(ad *Ad, campaign *smodel.CampaignInfo, params *mvutil.Params, dspId int64) {
	abtestFieldsConf, _ := extractor.GetABTEST_FIELDS()
	abtestConfs, _ := extractor.GetABTEST_CONFS()

	if len(abtestFieldsConf.EndcardUrlParams) == 0 {
		return
	}
	if params.ABTestTags == nil {
		params.ABTestTags = make(map[string]int)
	}
	var abtestParams string
	for _, v := range abtestFieldsConf.EndcardUrlParams {
		if abtestConf, ok := abtestConfs[v]; ok {
			conf, randType := GetABTestConf(params, campaign, abtestConf, dspId)
			if len(conf) == 0 {
				continue
			}
			finalVal, randOk := GetABTestRes(conf, randType, params, v)
			if !randOk {
				continue
			}
			params.ABTestTags[v] = finalVal
			// 拼装endcard_url参数
			abtestParams += "&" + v + "=" + strconv.Itoa(finalVal)
		}
	}
	// ia 的encardurl也需兼容
	if params.AdType == mvconst.ADTypeInteractive {
		params.IAPlayableUrl += abtestParams
	} else {
		ad.EndcardUrl += abtestParams
	}
}

func renderABTestParams(params *mvutil.Params) {
	if params.ABTestTags == nil || len(params.ABTestTags) == 0 {
		return
	}
	str, err := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(params.ABTestTags)
	if err == nil {
		params.ABTestTagStr = string(str)
	}
}

func prisonStorekitTime(params *mvutil.Params, ad *Ad) {
	// 针对ios sdk的流量做处理
	if !mvutil.IsHbOrV3OrV5Request(params.RequestPath) || params.Platform != mvconst.PlatformIOS {
		return
	}
	// ios设备切系统版本大于等于13，sdk版本处于 [520, 561] 时，返回广告的 storekit_time强制设置为1
	skTimeConf := extractor.GetSTOREKIT_TIME_PRISON_CONF()
	if skTimeConf.MinOsVer > 0 && params.OSVersionCode < skTimeConf.MinOsVer {
		return
	}
	if skTimeConf.MaxOsVer > 0 && params.OSVersionCode > skTimeConf.MaxOsVer {
		return
	}
	if skTimeConf.MinSdkVer > 0 && params.FormatSDKVersion.SDKVersionCode < skTimeConf.MinSdkVer {
		return
	}
	if skTimeConf.MaxSdkVer > 0 && params.FormatSDKVersion.SDKVersionCode > skTimeConf.MaxSdkVer {
		return
	}
	// 若均无配置，则不做下毒处理
	if skTimeConf.MinOsVer == 0 && skTimeConf.MaxOsVer == 0 && skTimeConf.MinSdkVer == 0 && skTimeConf.MaxSdkVer == 0 {
		return
	}
	skTime := int32(1)
	if skTimeConf.StorekitTimeVal > 0 {
		skTime = skTimeConf.StorekitTimeVal
	}
	ad.StoreKitTime = skTime
}

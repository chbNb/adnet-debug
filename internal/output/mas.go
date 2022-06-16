package output

import (
	"gitlab.mobvista.com/ADN/exporter/metrics"
	"net/url"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"gitlab.mobvista.com/ADN/adnet/internal/corsair_proto"
	"gitlab.mobvista.com/ADN/adnet/internal/extractor"
	"gitlab.mobvista.com/ADN/adnet/internal/mvconst"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	"gitlab.mobvista.com/ADN/adnet/internal/watcher"
	"gitlab.mobvista.com/ADN/mtg_openrtb/pkg/mtgrtb"
	"gitlab.mobvista.com/adserver/recommend_protocols/go/ad_server"
)

// renderMasCampaigns   mas(ad server)其实是pioneer的别称
// isMoreAds：　如果true, 需要重置dspext
func renderMasCampaigns(r *mvutil.RequestParams, ads *corsair_proto.BackendAds, isMoreAds bool) []Ad {
	asResp := r.AsResp
	r.Param.Extra2 = asResp.GetExtra2()
	extra2Arr := strings.Split(r.Param.Extra2, ",")
	r.Param.Extra6 = len(extra2Arr)
	r.Param.Extra3 = asResp.GetExtra3()
	r.Param.Extra20 = asResp.GetExtra20()
	r.Param.ExtfinalPackageName = asResp.GetExtFinalPackageName()
	r.Param.Extb2t, _ = strconv.Atoi(asResp.GetExtB2T())
	r.Param.Extalgo = asResp.GetExtAlgo()
	extifLowerImp, _ := strconv.Atoi(asResp.GetExtIfLowerImp())
	r.Param.ExtifLowerImp = int32(extifLowerImp)
	r.Param.ExtTagArrStr = asResp.GetExtTagList()
	r.Param.BidsPrice = asResp.GetExtBidsPrice()
	r.Param.ExtPlayableArr = asResp.GetExtPlayable()
	r.Param.EndcardUrl = asResp.GetEndscreenUrl()
	r.Param.Extendcard = strconv.FormatInt(int64(asResp.GetEndscreenId()), 10)
	r.Param.BigTemplateId = asResp.GetBigTemplateId()
	r.Param.ExtBigTemId = strconv.FormatInt(r.Param.BigTemplateId, 10)
	r.Param.BigTemplateUrl = asResp.GetBigTemplateUrl()
	r.Param.ExtBigTplOfferData = asResp.GetBigTemplateOfferData()
	r.Param.PolarisTplData = asResp.GetPolarisTplData()
	r.Param.PolarisCreativeData = asResp.GetPolarisCreativeData()
	r.Param.TplGrayTag = asResp.GetTplGrayTag()
	r.Param.ABTestTagStr = asResp.GetAbtestResTag()
	r.Param.AsABTestResTag = asResp.GetAsAbtestResTag()
	r.Param.JunoCommonLogInfoJson = asResp.GetJunoLogInfo()
	r.Param.AdspaceType = asResp.GetAdspaceType()
	r.Param.MaterialType = asResp.GetMaterialType()
	r.Param.PioneerExtdataInfo = asResp.GetPioneerExtdataInfo()
	r.Param.PioneerOfferExtdataInfo = asResp.GetPioneerOfferExtdataInfo()

	r.Param.BackendID = ads.BackendId
	// backendid 记录为mobvista。兼容报表逻辑
	if mvutil.IsRequestPioneerDirectly(&r.Param) {
		r.Param.BackendID = mvconst.Mobvista
	}
	r.Param.RequestKey = ads.RequestKey

	// tokenrule headerbidding流量，如果1表示该广告可以被应用在其他token中，比如其他token如果load，本地已有标记为1的有效缓存广告可直接回调load成功。不下发或2表示token各自独立。默认不下发。需支持abtest。{encrypt_p}
	r.Param.TokenRule = asResp.GetTokenRule()

	var adList []Ad
	if len(ads.CampaignList) <= 0 {
		return adList
	}
	for k, v := range ads.CampaignList {
		if v == nil {
			continue
		}
		if len(v.CampaignId) == 0 {
			continue
		}
		//降填充逻辑
		if isReduceFill(r, v, ads.BackendId) {
			continue
		}

		if filterFloor(r, v, ads.BackendId) {
			continue
		}

		ad, err := RenderMasCampaign(r, *v, ads.BackendId, k, isMoreAds, asResp.GetUrlParam())
		if err != nil {
			continue
		}
		adList = append(adList, ad)
	}

	if len(adList) > 0 {
		r.Param.FillBackendId = append(r.Param.FillBackendId, ads.BackendId)
	}
	// pv的adbeckend 由原来记录请求的backendid改为有填充的backendid
	r.Param.BackendList = append(r.Param.BackendList, ads.BackendId)
	return adList
}

func RenderMasCampaign(r *mvutil.RequestParams, corsairCampaign corsair_proto.Campaign, backendId int32, k int, isMoreAds bool, urlParams []*mtgrtb.BidResponse_UrlParam) (Ad, error) {
	// 声明
	var ad Ad

	params := r.Param

	// 根据map赋值
	campaignID, _ := strconv.ParseInt(corsairCampaign.CampaignId, 10, 64)
	ad.CampaignID = campaignID

	if corsairCampaign.AppName != nil {
		ad.AppName = *(corsairCampaign.AppName)
	}
	if corsairCampaign.AppDesc != nil {
		ad.AppDesc = *(corsairCampaign.AppDesc)
	}
	if corsairCampaign.PackageName != nil {
		ad.PackageName = *(corsairCampaign.PackageName)
	}
	if corsairCampaign.WTick != nil {
		ad.WithOutInstallCheck = int(*(corsairCampaign.WTick))
	}
	if corsairCampaign.Floatball != nil {
		ad.Floatball = *(corsairCampaign.Floatball)
	}
	if corsairCampaign.FloatballSkipTime != nil {
		ad.FloatballSkipTime = *(corsairCampaign.FloatballSkipTime)
	}
	if corsairCampaign.ImageSize != nil {
		ad.ImageSize = *(corsairCampaign.ImageSize)
	}
	if corsairCampaign.VideoLength != nil {
		ad.VideoLength = int(*(corsairCampaign.VideoLength))
	}
	if corsairCampaign.VideoSize != nil {
		ad.VideoSize = int(*(corsairCampaign.VideoSize))
	}
	if corsairCampaign.BitRate != nil {
		ad.Bitrate = *(corsairCampaign.BitRate)
	}
	if corsairCampaign.VideoResolution != nil {
		ad.VideoResolution = *(corsairCampaign.VideoResolution)
	}
	if corsairCampaign.AdTemplate != nil {
		ad.Template = int(*(corsairCampaign.AdTemplate))
	}

	ad.PlayableAdsWithoutVideo = 1

	if corsairCampaign.IconURL != nil {
		ad.IconURL = *(corsairCampaign.IconURL)
	}
	if corsairCampaign.ImageURL != nil {
		ad.ImageURL = *(corsairCampaign.ImageURL)
	}
	if corsairCampaign.ImageResolution != nil {
		ad.ImageResolution = *(corsairCampaign.ImageResolution)
	}
	if corsairCampaign.ImageMime != nil {
		ad.ImageMime = *(corsairCampaign.ImageMime)
	}
	if corsairCampaign.ClickURL != nil {
		ad.ClickURL = *(corsairCampaign.ClickURL)
	}
	if corsairCampaign.NoticeURL != nil {
		ad.NoticeURL = *(corsairCampaign.NoticeURL)
	}
	if corsairCampaign.ClickMode != nil {
		ad.ClickMode = int(*(corsairCampaign.ClickMode))
	}
	if corsairCampaign.LinkType != nil {
		ad.CampaignType = int(*(corsairCampaign.LinkType))
	}
	if corsairCampaign.AdTemplate != nil {
		ad.Template = int(*(corsairCampaign.AdTemplate))
	}
	if corsairCampaign.VideoEndType != nil {
		ad.VideoEndType = int(*corsairCampaign.VideoEndType)
	}
	ad.AdSourceID = int(corsairCampaign.AdSource)
	if corsairCampaign.FCA != nil {
		ad.FCA = int(*(corsairCampaign.FCA))
	}
	if corsairCampaign.FCB != nil {
		ad.FCB = int(*(corsairCampaign.FCB))
	}
	if corsairCampaign.VideoURL != nil {
		ad.VideoURL = *(corsairCampaign.VideoURL)
	}
	if corsairCampaign.CType != nil {
		ad.CType = int(*(corsairCampaign.CType))
	}
	if corsairCampaign.CreativeId != nil {
		creIdInt, _ := strconv.ParseInt(*(corsairCampaign.CreativeId), 10, 64)
		ad.CreativeId = creIdInt
	}
	if corsairCampaign.OfferType != nil {
		ad.OfferType = int(*corsairCampaign.OfferType)
	}
	ad.AdvImp = []CAdvImp{}
	if len(corsairCampaign.AdvImpList) > 0 {
		for _, v := range corsairCampaign.AdvImpList {
			if v == nil {
				continue
			}
			var advImp CAdvImp
			advImp.Sec = int(v.Second)
			advImp.Url = v.URL
			ad.AdvImp = append(ad.AdvImp, advImp)
		}
	}
	ad.AdURLList = []string{}
	if len(corsairCampaign.AdURLList) > 0 {
		ad.AdURLList = corsairCampaign.AdURLList
	}
	ad.AdvID = 0
	if corsairCampaign.GuideLines != nil {
		ad.Guidelines = *(corsairCampaign.GuideLines)
	}
	if corsairCampaign.RetargetOffer != nil {
		ad.RetargetOffer = int(*(corsairCampaign.RetargetOffer))
	}
	if corsairCampaign.StatsURL != nil {
		ad.StatsURL = *(corsairCampaign.StatsURL)
	}
	ad.LoopBack = corsairCampaign.LoopBack
	if corsairCampaign.EndcardURL != nil {
		ad.EndcardUrl = renderUrlAttachSchema(params.HTTPReq, *corsairCampaign.EndcardURL)
	}
	if corsairCampaign.OfferName != nil {
		ad.OfferName = *(corsairCampaign.OfferName)
	}

	if r.Param.AdType == mvconst.ADTypeRewardVideo || r.Param.AdType == mvconst.ADTypeInterstitialVideo {
		ad.EndcardClickResult = 1
	}

	if corsairCampaign.WatchMile != nil {
		ad.WatchMile = int(corsairCampaign.GetWatchMile())
	}

	if corsairCampaign.AdvertiserId != nil {
		ad.AdvID = int(corsairCampaign.GetAdvertiserId())
	}

	if corsairCampaign.AppSize != nil {
		ad.AppSize = corsairCampaign.GetAppSize()
	}

	if corsairCampaign.Rating != nil {
		ad.Rating = float32(corsairCampaign.GetRating())
	}
	if corsairCampaign.NumberRating != nil {
		ad.NumberRating = int(corsairCampaign.GetNumberRating())
	}
	if ad.NumberRating == 0 {
		ad.NumberRating = HandleNumberRating(0)
	}

	if corsairCampaign.CtaText != nil {
		ad.CtaText = corsairCampaign.GetCtaText()
	}

	if params.ApiVersion >= mvconst.API_VERSION_2_0 {
		if corsairCampaign.UrlTemplate != nil {
			ad.CamTplUrl = *corsairCampaign.UrlTemplate
		}

		if corsairCampaign.HtmlTemplate != nil {
			ad.CamHtml = *corsairCampaign.HtmlTemplate
		}
	}

	// c_ct 字段
	settingConfs, ifFind := extractor.GetSETTING_CONFIG()
	if ifFind {
		if params.Platform == mvconst.PlatformAndroid {
			ad.ClickCacheTime = settingConfs.ACCT
		}
		ad.ClickCacheTime = settingConfs.CCT
	}

	// 其他
	if len(*corsairCampaign.ImageSize) > 0 && corsairCampaign.ImageSizeId == ad_server.ImageSizeEnum_UNKNOWN && r.Param.RequestType == mvconst.REQUEST_TYPE_OPENAPI_AD {
		ad.ImageSize = *corsairCampaign.ImageSize
	} else {
		ad.ImageSize = mvconst.GetImageSizeByID(int(corsairCampaign.ImageSizeId))
	}
	// 针对小程序的banner广告位，固定返回1125x310(wx_banner只有这种尺寸的大图)
	if params.AdType == mvconst.ADTypeWXBanner {
		ad.ImageSize = "1125x310"
	}
	params.ImageSize = ad.ImageSize

	// render
	// 截取desc
	ad.AppDesc = mvutil.SubUtf8Str(ad.AppDesc, 38)

	if params.AdType == mvconst.ADTypeRewardVideo || params.AdType == mvconst.ADTypeInterstitialVideo {
		var rv RV
		if corsairCampaign.Rv != nil {
			rv.Orientation = int(corsairCampaign.Rv.GetOrientation())
			rv.TemplateUrl = renderUrlAttachSchema(params.HTTPReq, corsairCampaign.Rv.GetTemplateURL())
			rv.PausedUrl = renderUrlAttachSchema(params.HTTPReq, corsairCampaign.Rv.GetPausedURL())
			rv.VideoTemplate = int(corsairCampaign.Rv.GetTemplate())
			r.Param.Extrvtemplate = rv.VideoTemplate
		}
		ad.Rv = rv
		// 根据pub，app，unit对orientaton下毒。
		if OrientationPoison(&params) {
			ad.Rv.Orientation = 0
		}
		ad.RvPoint = &ad.Rv
	}
	params.DspExt = r.DspExt
	strBackend := strconv.FormatInt(int64(backendId), 10)
	priceFactorObj, ifFind := r.PriceFactor.Load(strBackend)
	if ifFind {
		priceFactor, ok := priceFactorObj.(string)
		if ok {
			r.Param.PriceFactor = priceFactor
		}
	}
	params.IsHBRequest = r.IsHBRequest
	//// 针对mv dsp的sdk banner流量的处理逻辑
	//if params.AdType == mvconst.ADTypeSdkBanner {
	//	renderClickmode(&ad, &params, corsairCampaign, r)
	//}
	//renderABTestRes(&params)
	RenderMWParams(&params, &ad, false)
	// tracking cn abtest flag
	var urlParam *mtgrtb.BidResponse_UrlParam
	if len(urlParams) > 0 && len(urlParams) > k {
		urlParam = urlParams[k]
		if strings.HasSuffix(urlParam.GetK(), "v") {
			r.Param.CNTrackingDomainTag = true
			// api_version小于1.4,会直接使用params.Domain生成链接
			params.Domain = extractor.GetTrackingCNABTestConf().Domain
		}
	}
	//渲染URL
	RenderNewUrls(r, &ad, &params, nil, campaignID, k, isMoreAds)
	ad.ImpressionURL = createImpressionUrl(&params)
	// videourl encode
	if len(ad.VideoURL) > 0 && !mvutil.IsBannerOrSplashOrDI(params.AdType) {
		ad.VideoURL = mvutil.Base64Encode(ad.VideoURL)
	}
	//storekit_time
	ad.StoreKitTime = 1
	if r.AppInfo.App.StorekitLoading == 2 {
		ad.StoreKitTime = 2
	}
	sdkParams := r.AsResp.GetSdkParam()
	if len(sdkParams) > 0 && len(sdkParams) > k {
		sdkParam := sdkParams[k]
		r.Param.ExtDataInit.ReturnAfTokenParam = int(sdkParam.GetReturnAfToken())
		r.Param.ExtDataInit.ClickInServer = int(sdkParam.GetClickInServer())
		// param里需要使用ClickInServer信息
		params.ExtDataInit.ClickInServer = int(sdkParam.GetClickInServer())
		//onlineAPI
		params.Extra10 = strconv.Itoa(int(sdkParam.GetClickMode()))
	}
	if r.Param.RequestType != mvconst.REQUEST_TYPE_OPENAPI_AD { //onlineAPI流量的Notice不需要赋值
		ad.NoticeURL = createClickUrl(&params, true)
	}

	// click mode 14情况下，notice url新增&redirect=1，客户端无需跳落地页
	// 返回给sdk 的clickmode 固定为5
	if strconv.Itoa(ad.ClickMode) == mvconst.JUMP_TYPE_TRACKING_PING_THIRDPARTY_CLICK_URL && len(ad.NoticeURL) > 0 {
		ad.NoticeURL += mvconst.REDIRECT
		ad.ClickMode = 5
	} else if mvutil.InStrArray(strconv.Itoa(ad.ClickMode), []string{mvconst.JUMP_TYPE_NORMAL, mvconst.JUMP_TYPE_ONLINE_DSP_AJUMP_SERVER}) &&
		len(ad.NoticeURL) > 0 && r.Param.RequestType == mvconst.REQUEST_TYPE_OPENAPI_V3 {
		ad.NoticeURL += mvconst.USELESS_NOTICE
		// mp pingmode 为0的情况下，不需要返回notice url,节省流量成本
		if mvutil.IsMpPingmode0(r) {
			ad.NoticeURL = ""
		}
		// 处理api_version小于1.3的click_url（more_offer）
		if !params.IsNewUrl {
			ad.ClickURL = replaceUrlMacro(ad.ClickURL, &params, GetUrlScheme(&params), params.Domain)
		}
		ad.ClickMode = 0
	}

	lazadaUrlSubStr := ""
	if len(sdkParams) > 0 && len(sdkParams) > k {
		sdkParam := sdkParams[k]
		ad.SubCategoryName = sdkParam.SubCategoryName
		//ad.Md5File = sdkParam.GetMd5File()
		ad.Storekit = int(sdkParam.GetStorekit())
		ad.StoreKitTime = sdkParam.GetStorekitTime()
		// 由adnet固定下发为1，因此不需要使用pioneer返回的值
		//ad.EndcardClickResult = sdkParam.GetEndcardClickResult()
		ad.ImpUa = int(sdkParam.GetImpUa())
		ad.CUA = int(sdkParam.GetCUa())
		ad.DeepLink = sdkParam.GetDeepLink()
		ad.NVT2 = sdkParam.GetNvT2()
		if mvutil.IsHbOrV3OrV5Request(r.Param.RequestPath) {
			sdkParam.GetAdchoice()
			//ad.AdChoice = sdkParam.Ad
			renderAdchoice4Mas(&ad, sdkParam.GetAdchoice())
		}

		var offerExtData OfferExtData
		offerExtData.SlotId = sdkParam.GetSlotId()
		if corsairCampaign.GifWidgets != nil {
			offerExtData.GifWidgets = *corsairCampaign.GifWidgets
		}
		ad.OfferExtData = &offerExtData

		if sdkParam.IdcdImgUrl != nil && len(*sdkParam.IdcdImgUrl) > 0 {
			if ad.Rv.Image == nil {
				ad.Rv.Image = &Image{IdcdImg: make([]string, 0)}
			}
			ad.Rv.Image.IdcdImg = append(ad.Rv.Image.IdcdImg, *sdkParam.IdcdImgUrl)
		}
		ad.GifURL = sdkParam.GetGifUrl()
		renderMasOfferSkadnetwork(&ad, sdkParam.GetSkAdnetwork())

		if len(sdkParam.GetLzdToken()) > 0 {
			lazadaUrlSubStr += "&lzd_token=" + sdkParam.GetLzdToken()
		}
		if len(sdkParam.GetLzdRtaid()) > 0 {
			lazadaUrlSubStr += "&lzd_rtaid=" + sdkParam.GetLzdRtaid()
		}

		// pioneer 支持lzd
		if len(ad.ImpressionURL) > 0 && len(lazadaUrlSubStr) > 0 {
			ad.ImpressionURL += lazadaUrlSubStr
			//if len(sdkParam.GetLzdToken()) > 0 {
			//	ad.ImpressionURL += "&lzd_token=" + sdkParam.GetLzdToken()
			//}
			//if len(sdkParam.GetLzdRtaid()) > 0 {
			//	ad.ImpressionURL += "&lzd_rtaid=" + sdkParam.GetLzdRtaid()
			//}
		}
		if len(ad.NoticeURL) > 0 {
			if sdkParam.GetReturnAfToken() == 1 {
				ad.NoticeURL += mvconst.APPSFLYER_TOKEN_PARAMS
			}
			if len(lazadaUrlSubStr) > 0 {
				ad.NoticeURL += lazadaUrlSubStr
			}
			//if len(sdkParam.GetLzdToken()) > 0 {
			//	ad.NoticeURL += "&lzd_token=" + sdkParam.GetLzdToken()
			//}
			//if len(sdkParam.GetLzdRtaid()) > 0 {
			//	ad.NoticeURL += "&lzd_rtaid=" + sdkParam.GetLzdRtaid()
			//}
		}
		ad.UserActivation = sdkParam.GetUserActivation()
		ad.PreviewUrl = sdkParam.GetPreviewUrl()

		if len(sdkParam.GetApkFmd5()) > 0 {
			ad.ApkFMd5 = sdkParam.GetApkFmd5()
		}

		ad.Ntbarpt = int(sdkParam.GetNtbarpt())
		ad.Ntbarpasbl = int(sdkParam.GetNtbarpasbl())
		ad.AtatType = int(sdkParam.GetAtatType())
		// 接收mas返回的apk_info
		RenderMasApkInfo(&ad, sdkParam.GetApkInfo(), &params)

		// lazada maitve
		renderMaitveInfo(&ad, sdkParam.GetMaitve())
		ad.VideoCtnType = sdkParam.GetVideoCtnType()
		ad.VideoCheckType = sdkParam.GetVideoCheckType()
		ad.RsIgnoreCheckRule = sdkParam.GetRsIgnoreCheckRule()
	}

	// 当投放lazada deeplink单子的情况下，目前对于online api的流量，会单独建立单链单子来跑（deeplink 字段为空值）
	// 而给sdk跑的单子是需要有deeplink的，运营同学现希望lazada或以后的deeplink单子，也使用sdk投放的单子来跑，减少与广告主对数的成本。
	// 因此需要把online api流量，deeplink 双链的单子改为单链来投放。
	changeOnlineDeeplinkWay(r, &ad)

	RenderMasAdtracking(&ad, corsairCampaign, &params, r, lazadaUrlSubStr)

	prisonStorekitTime(&params, &ad)
	adnConf, _ := extractor.GetADNET_SWITCHS()
	if params.Platform == mvconst.PlatformIOS {
		// 针对开发者下毒，不出storekit。设置开关，1则为关闭下毒
		closeSkPoison, ok := adnConf["closeSkPoison"]
		if (!ok || closeSkPoison != 1) && params.PublisherID == 14228 {
			ad.Storekit = mvconst.StorekitNotLoad
		}
		// ios sdk 问题版本storekit下毒
		if params.IosStorekitPoisonFlag {
			ad.Storekit = mvconst.StorekitNotLoad
		}
	}

	// topon adx下毒。返回Storekit 固定为1。因为topon adx使用的sdk插件，没有请求setting，无法获得setting的storekit值。
	if (r.IsTopon || params.PublisherID == mvconst.TOPONADXPublisherID) && r.Param.Platform == mvconst.PlatformIOS {
		ad.Storekit = mvconst.StorekitLoad
	}

	// 针对安卓，MD5_file下毒
	//if r.Param.Platform == mvconst.PlatformAndroid {
	//	pubsConf, _ := extractor.GetMD5_FILE_PRISON_PUB()
	//	if mvutil.InInt64Arr(r.Param.PublisherID, pubsConf) {
	//		ad.Md5File = ""
	//	}
	//}

	// notice url 下发appsflyer token宏

	// c_toi
	cToi, ifFind := extractor.GetC_TOI()
	if ifFind {
		ad.CToi = cToi
	}

	// webview下毒
	if r.Param.NeedWebviewPrison {
		ad.EndcardUrl = ""
		ad.RvPoint = nil
	}

	RenderThirdDemandAKS(corsairCampaign, &ad)
	renderPlct(&ad, &params, int(backendId), mvconst.MAS)

	ad.PlayableAdsWithoutVideo = int(corsairCampaign.GetPlayableAdsWithoutVideo())
	// 设置readyRate
	renderReadyRate(&ad)
	// 降展示
	if params.ExtifLowerImp == int32(1) {
		ad.ImageURL = ""
		ad.VideoURL = ""
	}
	//APK alert
	renderOfferApkAlt(r, &ad)

	// 控制是否要用pioneer生成的click_url。不配置，或配置的值不为1，则使用pioneer生成的。
	// deeplinkClickUrlSwitch, ok := adnConf["deeplinkClickUrlSwitch"]
	// 和小米storekit拆开上
	if ad.DeepLink != "" && ad.NoticeURL != "" {
		// clickmode 14，deeplink单子不添加&redirect=1，避免加上了&redirect=1导致无法兜底跳落地页的情况
		// 删除notice url的redirect=1
		ad.ClickURL = mvutil.DelSubStr(ad.NoticeURL, mvconst.REDIRECT) + mvconst.FORWARD // 只是记录点击
	}

	// 处理online api切量到 Mas 返回的bid_price
	renderThirdPartyOnlineApiBidPrice(&params, r, &ad)

	// 对于topon流量，新增出价到impression url上
	addToponBidPriceOnImpressionUrl(&ad, r)

	addUnitIdOnTrackingUrl(&ad, &params)
	renderOnlinePubClickUrl(&ad, &params)

	// 限制算法出价过高的情况，避免算法出价过高导致的爆量
	if filterRequestByHighBidPrice(&params, &ad) {
		watcher.AddWatchValue("bid_price_out_of_limit_max", float64(1))
		mvutil.Logger.Runtime.Errorf("online api bid price out of limit[mas]. bid price is:[%s],requestid is:[%s]", params.OnlineApiBidPrice, params.RequestID)
		return ad, errors.New("online api bid price out of limit[mas]")
	}

	if r.Param.HtmlSupport == 1 && corsairCampaign.AdHtml != nil {
		ad.AdHtml = *corsairCampaign.AdHtml
	}

	// tracking cn abtest flag
	if strings.HasSuffix(urlParam.GetK(), "v") {
		// 因为存在offer侧的切量，因此命中且api_version大于等于1.4的情况，不需要sdk替换url里的域名宏
		RepalceDoMacro(extractor.GetTrackingCNABTestConf().Domain, &ad, &params)
	}

	//onlineAPI的流量统一替换{*}
	replaceOnlineApiUrl(&ad, &params)

	// topon 记录link_type
	if r.IsTopon {
		r.Param.LinkType = ad.CampaignType
	}
	// reward plus
	RenderRewardPlus(r, &ad)

	ad.ViewCompletedTime = r.UnitInfo.Unit.ViewCompletedTime

	// 新插屏广告处理
	RenderNewInterstitialResponse(r, &ad)

	if mvutil.IsMP(r.Param.RequestPath) {
		RenderMPInfo(&ad)
	}

	// endcard_url,rv 模版合法性校验
	checkTemplateUrl(&ad, mvconst.MAS, params.RequestID)

	return ad, nil
}

func checkTemplateUrl(ad *Ad, dspId int, requestid string) {
	if len(ad.Rv.TemplateUrl) > 0 {
		if !isLegalTemplateUrl(ad.Rv.TemplateUrl) {
			TplErrCountAndLog(ad.Rv.TemplateUrl, mvconst.TEMPLATE_TYPE_TYPE_VIDEO, dspId, requestid)
		}
	}
	if len(ad.Rv.PausedUrl) > 0 {
		if !isLegalTemplateUrl(ad.Rv.PausedUrl) {
			TplErrCountAndLog(ad.Rv.PausedUrl, mvconst.TEMPLATE_TYPE_TYPE_PAUSE, dspId, requestid)
		}
	}
	if len(ad.EndcardUrl) > 0 {
		if !isLegalTemplateUrl(ad.EndcardUrl) {
			TplErrCountAndLog(ad.EndcardUrl, mvconst.TEMPLATE_TYPE_TYPE_ENDCARD, dspId, requestid)
		}
	}
	if len(ad.CamTplUrl) > 0 {
		if !isLegalTemplateUrl(ad.CamTplUrl) {
			TplErrCountAndLog(ad.CamTplUrl, mvconst.TEMPLATE_TYPE_TYPE_CAMTPL, dspId, requestid)
		}
	}
}

func TplErrCountAndLog(tplUrl string, templateType string, dspId int, requestid string) {
	mvutil.Logger.Runtime.Errorf("template url err.url is [%s], type is [%s], dspid is [%d], requestid is [%s]", tplUrl, templateType, dspId, requestid)
	metrics.IncCounterWithLabelValues(36, templateType, strconv.Itoa(dspId))
}

func isLegalTemplateUrl(tplUrl string) bool {
	urlData, err := url.ParseRequestURI(tplUrl)
	return err == nil && len(urlData.Host) > 0 && len(urlData.Path) > 0
}

func RenderNewInterstitialResponse(r *mvutil.RequestParams, ad *Ad) {
	// 限制api_version
	if !mvutil.IsNewIv(&r.Param) {
		return
	}
	// 1＝全屏 2＝半屏
	ad.AdspaceType = r.Param.AdspaceType
	// materialType 不为1，则为插屏视频
	// 除了全屏视频外需要查询unit配置，其余赋值为0
	closeButtonDelay := int32(0)
	videoSkipTime := int32(0)
	if r.Param.AdspaceType == 1 && r.Param.MaterialType != 1 {
		closeButtonDelay = int32(r.UnitInfo.Setting.CloseButtonDelay)
		videoSkipTime = int32(r.UnitInfo.Setting.VideoSkipTime)
	}
	ad.CloseButtonDelay = &closeButtonDelay
	ad.VideoSkipTime = &videoSkipTime

	// 如果是半屏，如果有end_screen_url，就会展示end_screen_url的html，会出现适配问题（各种按钮icon不匹配的问题）
	// 因此服务端针对半屏，不下发end_screen_url
	if r.Param.AdspaceType == 2 {
		r.Param.NewIVClearEndScreenUrl = true
	}
}

func renderMaitveInfo(ad *Ad, maitve int32) {
	ad.Maitve = maitve
	if maitve == 1 {
		ad.MaitveSrc = "Mtg"
	}
}

func OrientationPoison(params *mvutil.Params) bool {
	orientationPoisonConf := extractor.GetOrientationPoisonConf()
	if len(orientationPoisonConf.UnitBlackList) > 0 && mvutil.InInt64Arr(params.UnitID, orientationPoisonConf.UnitBlackList) {
		return false
	}
	if len(orientationPoisonConf.AppBlackList) > 0 && mvutil.InInt64Arr(params.AppID, orientationPoisonConf.AppBlackList) {
		return false
	}
	if len(orientationPoisonConf.UnitList) > 0 && mvutil.InInt64Arr(params.UnitID, orientationPoisonConf.UnitList) {
		return true
	}
	if len(orientationPoisonConf.AppList) > 0 && mvutil.InInt64Arr(params.AppID, orientationPoisonConf.AppList) {
		return true
	}
	if len(orientationPoisonConf.PubList) > 0 && mvutil.InInt64Arr(params.PublisherID, orientationPoisonConf.PubList) {
		return true
	}
	return false
}

func RenderMasApkInfo(ad *Ad, apkInfo *mtgrtb.BidResponse_ApkInfo, params *mvutil.Params) {
	// 只需针对国内流量（流量country code为CN）&apk广告（link_type为apk）&新版本sdk时返回
	if !IsCNApkTraffic(ad, params) {
		return
	}
	if len(apkInfo.GetAppName()) == 0 && len(apkInfo.GetSensitivePermission()) == 0 && len(apkInfo.GetOriginSensitivePermission()) == 0 &&
		len(apkInfo.GetAppVersionUpdateTime()) == 0 && len(apkInfo.GetAppVersion()) == 0 && len(apkInfo.GetDeveloperName()) == 0 &&
		len(apkInfo.GetPrivacyUrl()) == 0 {
		return
	}
	var adApkInfo ApkInfo
	adApkInfo.AppName = apkInfo.GetAppName()
	adApkInfo.SensitivePermission = apkInfo.GetSensitivePermission()
	adApkInfo.OriginSensitivePermission = apkInfo.GetOriginSensitivePermission()
	adApkInfo.PrivacyUrl = apkInfo.GetPrivacyUrl()
	adApkInfo.AppVersionUpdateTime = apkInfo.GetAppVersionUpdateTime()
	adApkInfo.AppVersion = apkInfo.GetAppVersion()
	adApkInfo.DeveloperName = apkInfo.GetDeveloperName()
	ad.ApkInfo = &adApkInfo
}

func renderMasOfferSkadnetwork(ad *Ad, skadnetwork *mtgrtb.BidResponse_SkAdNetwork) {
	if skadnetwork != nil && len(skadnetwork.GetVersion()) > 0 {
		ad.Skadnetwork = &Skadnetwork{
			Version:         skadnetwork.GetVersion(),
			Network:         skadnetwork.GetNetwork(),
			AppleCampaignId: skadnetwork.GetCampaign(),
			Targetid:        skadnetwork.GetItunesitem(),
			Nonce:           skadnetwork.GetNonce(),
			Sourceid:        skadnetwork.GetSourceapp(),
			Timestamp:       skadnetwork.GetTimestamp(),
			Sign:            skadnetwork.GetSignature(),
			Need:            int(skadnetwork.GetSkneed()),
		}
		if len(skadnetwork.GetViewSignature()) > 0 {
			ad.SkImp = &SkImp{
				ViewSign: skadnetwork.GetViewSignature(),
			}
		}
	}
}

func RenderMasAdtracking(ad *Ad,
	corsairCampaign corsair_proto.Campaign,
	params *mvutil.Params,
	r *mvutil.RequestParams,
	lazadaUrlSubStr string) {
	if corsairCampaign.AdTracking == nil {
		return
	}

	var adtracking CAdTracking
	csAdtracking := *(corsairCampaign.AdTracking)
	// 今日头条
	if len(csAdtracking.PlayPerct) > 0 {
		adtracking.Play_percentage = []CPlayTracking{}
		for _, v := range csAdtracking.PlayPerct {
			if v == nil {
				continue
			}
			var playPct CPlayTracking
			tmpPlay := *v
			playPct.Rate = int(tmpPlay.Rate)
			playPct.Url = tmpPlay.URL
			adtracking.Play_percentage = append(adtracking.Play_percentage, playPct)
		}
	}

	// 赋值
	adtracking.Start = csAdtracking.Start
	adtracking.First_quartile = csAdtracking.FirstQuart
	adtracking.Midpoint = csAdtracking.Mid
	adtracking.Third_quartile = csAdtracking.ThirdQuart
	adtracking.Complete = csAdtracking.Complete
	adtracking.Mute = csAdtracking.Mute
	adtracking.Unmute = csAdtracking.Unmute
	adtracking.Impression = csAdtracking.Impression
	adtracking.Click = csAdtracking.Click
	adtracking.Endcard_show = csAdtracking.EndcardShow
	adtracking.Close = csAdtracking.Close
	adtracking.Pause = csAdtracking.Pause
	adtracking.ApkDownloadStart = csAdtracking.ApkDownloadStart
	adtracking.ApkDownloadEnd = csAdtracking.ApkDownloadEnd
	adtracking.ApkInstall = csAdtracking.ApkInstall
	adtracking.PubImp = csAdtracking.PubImp
	adtracking.Video_Click = csAdtracking.VideoClick
	adtracking.Impression_t2 = csAdtracking.ImpressionT2

	// // 针对ads接口切第三方广告主的情况
	// 不需要处理了，因为针对REQUEST_TYPE_OPENAPI_AD不会给ad.NoticeURL
	// if params.RequestType == mvconst.REQUEST_TYPE_OPENAPI_AD && !params.IsVast && len(ad.NoticeURL) > 0 &&
	// 	len(ad.NoticeURL) > 0 {
	// 	adtracking.Click = append(adtracking.Click, ad.NoticeURL)
	// }

	if params.RequestType == mvconst.REQUEST_TYPE_OPENAPI_AD && ad.DeepLink != "" && ad.ClickURL != "" {
		// 控制是否要用pioneer生成的click_url。不配置，或配置的值不为1，则使用pioneer生成的。
		//adnConf, _ := extractor.GetADNET_SWITCHS()
		//onlineApideeplinkClickUrlSwitch, ok := adnConf["onlineApideeplinkClickUrlSwitch"]
		//if ok && onlineApideeplinkClickUrlSwitch == 1 {
		adtracking.Click = append(adtracking.Click, ad.ClickURL)
		//fallback
		deepLinkFallbackUrl := createClickUrl(params, false) + mvconst.FORWARD
		if len(lazadaUrlSubStr) > 0 {
			deepLinkFallbackUrl += lazadaUrlSubStr
		}
		ad.ClickURL = deepLinkFallbackUrl
		ad.PreviewUrl = deepLinkFallbackUrl
		//} else {
		//	ad.PreviewUrl = ad.ClickURL
		//}
	}

	// clickmode 13+ deeplink的情况下，clickurl 加上redirect=1,避免再跳到落地页的情况。
	// 这个逻辑做到pioneer更合适，后面和forward=1参数一起迁移到pioneer
	var originClickmodeStr string
	if corsairCampaign.ClickMode != nil {
		originClickmodeStr = strconv.Itoa(int(*(corsairCampaign.ClickMode)))
	}
	if originClickmodeStr == mvconst.JUMP_TYPE_ONLINE_DSP_AJUMP_SERVER && len(ad.DeepLink) > 0 {
		var clickList []string
		for _, clickUrl := range adtracking.Click {
			// 只对adn trackingurl处理。后面一定要迁移到pioneer处理。
			if strings.Contains(clickUrl, "/click?k=") && !strings.Contains(clickUrl, mvconst.REDIRECT) {
				clickUrl += mvconst.REDIRECT
			}
			clickList = append(clickList, clickUrl)
		}
		adtracking.Click = clickList
	}

	//埋chetconfig-link
	if link, ok := hitChetLinkConfig(params); ok {
		adtracking.Impression = append(adtracking.Impression, replaceChetLink(link, params))
	}

	// 针对sdk有问题版本下毒。541,542,550,551,这些版本无法解析adtracking.imp给h5
	if ReplaceImpt2withImp(params) {
		for _, imp := range adtracking.Impression {
			adtracking.Play_percentage = append(adtracking.Play_percentage, CPlayTracking{Rate: 0, Url: imp})
		}
		//adtracking.Impression_t2 = adtracking.Impression
		adtracking.Impression = []string{}
	}

	// 对于online api 提前替换domain
	if params.UseCdnTrackingDomain == 1 {
		ReplaceAdTrackingDoMacroExcludeImpAndClick(&adtracking, params.TrackingCdnDomain)
	}

	replace4OnlineApi(params, &adtracking, r.Param.RequestType)

	ad.AdTracking = adtracking
	ad.AdTrackingPoint = &adtracking
}

//onlineAPI 的adntracking需要替换Mas的链接
func replace4OnlineApi(params *mvutil.Params, adtracking *CAdTracking, reqeustType int) {
	sh := GetUrlScheme(params)
	do := params.Domain
	if reqeustType == mvconst.REQUEST_TYPE_OPENAPI_AD {
		adtracking.Start = replaceUrls(adtracking.Start, params, sh, do)
		adtracking.First_quartile = replaceUrls(adtracking.First_quartile, params, sh, do)
		adtracking.Midpoint = replaceUrls(adtracking.Midpoint, params, sh, do)
		adtracking.Third_quartile = replaceUrls(adtracking.Third_quartile, params, sh, do)
		adtracking.Complete = replaceUrls(adtracking.Complete, params, sh, do)
		adtracking.Mute = replaceUrls(adtracking.Mute, params, sh, do)
		adtracking.Unmute = replaceUrls(adtracking.Unmute, params, sh, do)
		adtracking.Impression = replaceUrls(adtracking.Impression, params, sh, do)
		adtracking.Click = replaceUrls(adtracking.Click, params, sh, do)
		adtracking.Endcard_show = replaceUrls(adtracking.Endcard_show, params, sh, do)
		adtracking.Close = replaceUrls(adtracking.Close, params, sh, do)
		adtracking.Pause = replaceUrls(adtracking.Pause, params, sh, do)
		adtracking.Video_Click = replaceUrls(adtracking.Video_Click, params, sh, do)
		adtracking.Impression_t2 = replaceUrls(adtracking.Impression_t2, params, sh, do)
		adtracking.ApkDownloadStart = replaceUrls(adtracking.ApkDownloadStart, params, sh, do)
		adtracking.ApkDownloadEnd = replaceUrls(adtracking.ApkDownloadEnd, params, sh, do)
		adtracking.ApkInstall = replaceUrls(adtracking.ApkInstall, params, sh, do)
		adtracking.Dropout = replaceUrls(adtracking.Dropout, params, sh, do)
		adtracking.Plycmpt = replaceUrls(adtracking.Plycmpt, params, sh, do)
		adtracking.PubImp = replaceUrls(adtracking.PubImp, params, sh, do)
		adtracking.ExaClick = replaceUrls(adtracking.ExaClick, params, sh, do)
		adtracking.ExaImp = replaceUrls(adtracking.ExaImp, params, sh, do)
		for k, v := range adtracking.Play_percentage {
			adtracking.Play_percentage[k].Url = replaceUrls([]string{v.Url}, params, sh, do)[0]
		}
	}
}

func renderAdchoice4Mas(ad *Ad, adchoice *mtgrtb.BidResponse_AdChoice) {
	if adchoice != nil {
		ad.AdChoice = &AdChoice{
			AdLogolink:   adchoice.GetAdLogoLink(),
			AdchoiceIcon: adchoice.GetAdchoiceIcon(),
			AdchoiceLink: adchoice.GetAdchoiceLink(),
			AdchoiceSize: adchoice.GetAdchoiceSize(),
			AdvLogo:      adchoice.GetAdvLogo(),
			AdvName:      adchoice.GetAdvName(),
			PlatformLogo: adchoice.GetPlatformLogo(),
			PlatformName: adchoice.GetPlatformName(),
		}
	}

}

package output

import (
	"errors"
	"fmt"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/helpers"
	rtb "gitlab.mobvista.com/ADN/mtg_openrtb/pkg/mtgrtb"
	"math/rand"
	"net/url"
	"strconv"
	"strings"

	supply_mvutil "gitlab.mobvista.com/ADN/chasm/module/supply/mvutil"

	jsoniter "github.com/json-iterator/go"
	"gitlab.mobvista.com/ADN/adnet/internal/corsair_proto"
	"gitlab.mobvista.com/ADN/adnet/internal/extractor"
	"gitlab.mobvista.com/ADN/adnet/internal/mvconst"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	smodel "gitlab.mobvista.com/ADN/structs/model"
)

func RenderThirdPartyCampaign(r *mvutil.RequestParams, corsairCampaign corsair_proto.Campaign, backendId int32, i int) Ad {
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
	if corsairCampaign.ImageSize != nil {
		ad.ImageSize = *(corsairCampaign.ImageSize)
	}
	if corsairCampaign.VideoLength != nil {
		ad.VideoLength = int(*(corsairCampaign.VideoLength))
	}
	if corsairCampaign.VideoSize != nil {
		ad.VideoSize = int(*(corsairCampaign.VideoSize))
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
		// 对于RV和IV， 从M ADX取回的素材返回给SDK时，视频模板置空 （ADNET-151）
		//if backendId == mvconst.MAdx &&
		//	(r.Param.AdType == constant.AD_TYPE_RV || r.Param.AdType == constant.AD_TYPE_IV) &&
		//	r.Param.RequestPath == mvconst.PATHOpenApiV3 {
		//	ad.Template = 0
		//} else {
		ad.Template = int(*(corsairCampaign.AdTemplate))
		//}
	}
	if r.Param.AdType == mvconst.ADTypeRewardVideo || r.Param.AdType == mvconst.ADTypeInterstitialVideo {
		ad.VideoEndType = int(*corsairCampaign.VideoEndType)
	}
	ad.AdSourceID = int(corsairCampaign.AdSource)
	if corsairCampaign.FCA != nil {
		ad.FCA = int(*(corsairCampaign.FCA))
	}
	if corsairCampaign.FCB != nil {
		ad.FCB = int(*(corsairCampaign.FCB))
	}
	if corsairCampaign.Rating != nil {
		ad.Rating = float32(*(corsairCampaign.Rating))
	}
	if corsairCampaign.CtaText != nil {
		ad.CtaText = *(corsairCampaign.CtaText)
	}

	// 开屏cta text
	// 优先级高于三方dsp返回的cta text，主要为了合规性
	renderSplashCtaText(&ad, &params)

	if len(ad.CtaText) == 0 {
		ad.CtaText = handleCTAButton(int32(ad.CampaignType), &params)
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
	if corsairCampaign.AdHtml != nil {
		ad.AdHtml = *corsairCampaign.AdHtml
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

	if mvutil.IsNewIv(&params) {
		r.Param.AdspaceType = r.UnitInfo.Unit.AdSpaceType
		r.Param.MaterialType = r.UnitInfo.Unit.MaterialType
	}

	var dspId int64
	dspExt, err := r.GetDspExt()
	// 获取dspId，用于plct，plctb取值
	if err == nil && dspExt != nil {
		dspId = dspExt.DspId
	}

	// 限制ad_type，以免和
	if (params.ApiVersion >= mvconst.API_VERSION_2_0 || extractor.IsVastBannerDsp(dspId)) &&
		mvutil.IsBannerOrSplashOrNativeH5(params.FormatAdType) {
		if corsairCampaign.UrlTemplate != nil {
			ad.CamTplUrl = renderBannerHtml(*corsairCampaign.UrlTemplate, dspId)
		}

		if corsairCampaign.HtmlTemplate != nil {
			ad.CamHtml = renderBannerHtml(*corsairCampaign.HtmlTemplate, dspId)
		}
		// sdk_banner，splash,interstital的模版改为在adnet获取下发
		renderCamTplUrl(&ad, dspId, &params, corsairCampaign)
	}
	// omsdk
	ad.OmSDK = r.OmSDK
	if corsairCampaign.OfferName != nil {
		ad.OfferName = *(corsairCampaign.OfferName)
	}

	if r.Param.AdType == mvconst.ADTypeRewardVideo || r.Param.AdType == mvconst.ADTypeInterstitialVideo {
		ad.EndcardClickResult = 1
	}

	if corsairCampaign.DeepLink != nil {
		ad.DeepLink = *(corsairCampaign.DeepLink)
	}
	// 当投放lazada deeplink单子的情况下，目前对于online api的流量，会单独建立单链单子来跑（deeplink 字段为空值）
	// 而给sdk跑的单子是需要有deeplink的，运营同学现希望lazada或以后的deeplink单子，也使用sdk投放的单子来跑，减少与广告主对数的成本。
	// 因此需要把online api流量，deeplink 双链的单子改为单链来投放。
	changeOnlineDeeplinkWay(r, &ad)

	// render
	// 截取desc
	ad.AppDesc = mvutil.SubUtf8Str(ad.AppDesc, 38)

	if r.Adchoice != nil {
		ad.AdChoice = &AdChoice{
			AdchoiceIcon: r.Adchoice.AdchoiceIcon,
			AdchoiceLink: r.Adchoice.AdchoiceLink,
		}
	}

	if mvutil.IsIvOrRv(params.AdType) {
		if len(corsairCampaign.GetHtmlTemplate()) > 0 {
			ad.Mraid = renderMraidByDsp(corsairCampaign.GetHtmlTemplate(), dspId)
		}

		if len(corsairCampaign.GetUrlTemplate()) > 0 {
			ad.Mraid = renderMraidByDsp(corsairCampaign.GetUrlTemplate(), dspId)
		}

		if len(ad.VideoURL) == 0 && (len(corsairCampaign.GetHtmlTemplate()) > 0 || len(corsairCampaign.GetUrlTemplate()) > 0) {
			ad.PlayableAdsWithoutVideo = 2
		}
	}

	// rv
	if (params.AdType == mvconst.ADTypeRewardVideo || params.AdType == mvconst.ADTypeInterstitialVideo) && isReturnRv(r.Param.UnitID, dspId) {
		var rv RV
		if corsairCampaign.Rv != nil {
			rv.Orientation = int(corsairCampaign.Rv.GetOrientation())
			rv.TemplateUrl = renderUrlAttachSchema(params.HTTPReq, corsairCampaign.Rv.GetTemplateURL())
			rv.PausedUrl = renderUrlAttachSchema(params.HTTPReq, corsairCampaign.Rv.GetPausedURL())
			rv.VideoTemplate = int(corsairCampaign.Rv.GetTemplate())
		}
		// 对于第三方的DSP, 固定使用902模板
		if mvutil.IsThirdDspWithRtdsp(dspId) {
			templateMapConf, ifFind := extractor.GetTEMPLATE_MAP()
			if ifFind {
				// 902: 视频模板V9
				videoTplId := 902

				// 三方dsp支持切量全局可点
				if hit := IsAllowSDKVideoFullClick(dspId, r); hit {
					// 命中切量
					videoTplId = 904
					// Log
					params.ThirdPartyABTestRes.SDKVideoFullClick = 1
				}

				// 判断有无走abtest逻辑
				if r.Param.VideoDspTplABTestId > 0 {
					videoTplId = r.Param.VideoDspTplABTestId
				}
				if templateUrl, ok := templateMapConf.Video[strconv.Itoa(videoTplId)]; ok {
					rv.TemplateUrl = renderUrlAttachSchema(params.HTTPReq, templateUrl)
					rv.VideoTemplate = videoTplId
				}

				// 1102: 第三方DSP专用endcard模板
				r.Param.EndcardTplId = "1102"
				if r.Param.EndcardDspTplABTestId > 0 {
					r.Param.EndcardTplId = strconv.Itoa(r.Param.EndcardDspTplABTestId)
				}
				if endcardUrl, ok := templateMapConf.EndScreen[r.Param.EndcardTplId]; ok {
					ad.EndcardUrl = renderUrlAttachSchema(params.HTTPReq, endcardUrl)
				}

				// 针对新插屏的视频，ec模版
				renderNewIvEndcardAndVideoTemplate(&rv, r, &ad)

				// 新旧模版切量abtest
				renderRvTemplateUrlAbTest(&rv, &params, videoTplId)

				renderEndcardTemplateUrlAbTest(&params, r.Param.EndcardTplId, &ad)

				// 针对模版上的宏参数进行替换
				renderCretiveDomainAbTest(&rv, &params, &ad)
			}
		}

		// default维度
		conf, ifFind := extractor.GetDefRVTemplate()
		if r.Param.ApiVersion >= mvconst.API_VERSION_1_2 && ifFind {
			if len(conf.PausedURLZip) > 0 {
				conf.PausedURL = conf.PausedURLZip
			}
		}
		ad.Rv = rv

		// 模板回退逻辑
		//RvIsBack(&ad, &params, conf)
		ad.Rv.TemplateUrl = renderUrlAttachSchema(params.HTTPReq, ad.Rv.TemplateUrl)
		if len(ad.Rv.PausedUrl) == 0 && len(conf.PausedURL) > 0 {
			ad.Rv.PausedUrl = conf.PausedURL
		}
		ad.Rv.PausedUrl = renderUrlAttachSchema(params.HTTPReq, ad.Rv.PausedUrl)
		ad.Rv.Orientation = r.Param.FormatOrientation
		// 根据pub，app，unit对orientaton下毒。
		if OrientationPoison(&params) {
			ad.Rv.Orientation = 0
		}
		ad.RvPoint = &ad.Rv
	}
	renderThirdPartyEndcardProperty(&ad, r, dspId)

	// 开屏模版
	renderSpecalSplashProperty(r.UnitInfo, &r.Param, dspId, &ad)

	if backendId == mvconst.MAdx {
		params.DspExt = r.DspExt
		//params.PriceFactor = r.PriceFactor
	}

	strBackend := strconv.FormatInt(int64(backendId), 10)
	priceFactorObj, ifFind := r.PriceFactor.Load(strBackend)
	if ifFind {
		priceFactor, ok := priceFactorObj.(string)
		if ok {
			params.PriceFactor = priceFactor
		}
	}
	// 针对mv dsp的sdk banner流量的处理逻辑
	if mvutil.IsBannerOrSplashOrDI(params.AdType) {
		renderClickmode(&ad, &params, corsairCampaign, r)
	}
	renderABTestRes(&params, r, &ad)

	// 处理online api切量到三方dsp返回的bid_price
	renderThirdPartyOnlineApiBidPrice(&params, r, &ad)

	// 对于topon流量，新增出价到impression url上
	addToponBidPriceOnImpressionUrl(&ad, r)

	//params.DPrice = r.DPrice

	if r.IsTopon {
		params.IsNewUrl = false
		params.HTTPReq = 2
	}
	RenderThirdPartyUrls(&params, &ad, r)
	// topon 获取mp tracking url
	if r.IsTopon {
		r.Param.ToponThirdPartyImpUrl = ad.ImpressionURL
		r.Param.ToponThirdPartyNoticeUrl = ad.NoticeURL
		// topon 记录link_type
		r.Param.LinkType = ad.CampaignType
	}

	// adTracking
	RenderThirdAdtracking(&ad, corsairCampaign, &params, backendId, r)
	// videourl encode
	if len(ad.VideoURL) > 0 {
		ad.VideoURL = mvutil.Base64Encode(ad.VideoURL)
	}
	ad.FCB = 2
	tmpRating := strconv.FormatFloat(float64(ad.Rating), 'f', 0, 32)
	tmpRatingF, err := strconv.ParseFloat(tmpRating, 32)
	if err == nil {
		ad.Rating = float32(tmpRatingF)
	}
	//storekit_time
	ad.StoreKitTime = 1
	if r.AppInfo.App.StorekitLoading == 2 {
		ad.StoreKitTime = 2
	}
	prisonStorekitTime(&params, &ad)
	if params.Platform == mvconst.PlatformIOS {
		// 针对开发者下毒，不出storekit。设置开关，1则为关闭下毒
		adnConf, _ := extractor.GetADNET_SWITCHS()
		closeSkPoison, ok := adnConf["closeSkPoison"]
		if (!ok || closeSkPoison != 1) && params.PublisherID == 14228 {
			ad.Storekit = mvconst.StorekitNotLoad
		}
		// ios sdk 问题版本storekit下毒
		if params.IosStorekitPoisonFlag {
			ad.Storekit = mvconst.StorekitNotLoad
		}
	}

	// c_toi
	cToi, ifFind := extractor.GetC_TOI()
	if ifFind {
		ad.CToi = cToi
	}
	if ad.NumberRating == 0 {
		ad.NumberRating = HandleNumberRating(0)
	}

	// webview下毒
	if r.Param.NeedWebviewPrison {
		ad.EndcardUrl = ""
		ad.RvPoint = nil
	}
	renderPlct(&ad, &params, int(backendId), dspId)
	// abtest 框架处理(主要用于下毒逻辑处理及全量配置。abtest主要以mv结果为准)
	var campaign smodel.CampaignInfo
	adParamABTest(&params, &campaign, &ad, dspId)

	//APK alert
	renderOfferApkAlt(r, &ad)

	renderOfferSkadnetwork(corsairCampaign, &ad, &params)

	if corsairCampaign.AKS != nil {
		if ad.AKS == nil {
			ad.AKS = corsairCampaign.AKS
		} else {
			for k, v := range corsairCampaign.AKS {
				ad.AKS[k] = v
			}
		}
	}

	// 三方dsp没有apk url，使用click_url代替
	ad.ApkFMd5 = renderApkFMd5(&params, ad.ClickURL, &ad, "")
	ad.Ntbarpt = r.Param.Ntbarpt
	ad.Ntbarpasbl = r.Param.Ntbarpasbl
	ad.AtatType = r.Param.AtatType

	RenderNewInterstitialResponse(r, &ad)

	// reward plus
	RenderRewardPlus(r, &ad)

	ad.ViewCompletedTime = r.UnitInfo.Unit.ViewCompletedTime

	checkTemplateUrl(&ad, int(dspId), params.RequestID)

	ad.ThirdPartyOffer = 1
	ad.FilterAutoClick = renderFilterAutoClick(&params, dspId)

	return ad
}

func renderFilterAutoClick(params *mvutil.Params, dspId int64) (isFilter int) {
	conf := extractor.GetFilterAutoClickConf()
	if conf == nil {
		return
	}
	// dsp,pub任一维度下发开启，则认为需要过滤点击
	if len(conf.DspList) > 0 && mvutil.InInt64Arr(dspId, conf.DspList) {
		return 1
	}
	if len(conf.PubList) > 0 && mvutil.InInt64Arr(params.PublisherID, conf.PubList) {
		return 1
	}

	// 整体下发了开启的情况下，dsp，pub的黑名单
	if len(conf.DspBlacklist) > 0 && mvutil.InInt64Arr(dspId, conf.DspBlacklist) {
		return
	}
	if len(conf.PubBlacklist) > 0 && mvutil.InInt64Arr(params.PublisherID, conf.PubBlacklist) {
		return
	}
	if conf.TotalStatus {
		return 1
	}
	return
}

func renderSplashCtaText(ad *Ad, params *mvutil.Params) {
	if params.AdType != mvconst.ADTypeSplash {
		return
	}
	language := handleLanguage(params.Language)
	if IsViewCampaignType(int32(ad.CampaignType)) {
		ad.CtaText = mvconst.SplashViewLang[language]
		if len(ad.CtaText) == 0 {
			ad.CtaText = mvconst.View
		}
	}
	if IsInstallCampaignType(int32(ad.CampaignType)) {
		ad.CtaText = mvconst.SplashInstallLang[language]
		if len(ad.CtaText) == 0 {
			ad.CtaText = mvconst.Install
		}
	}
}

func IsViewCampaignType(campaignType int32) bool {
	return campaignType == int32(rtb.LinkType_WEBVIEW) || campaignType == int32(rtb.LinkType_DEFAULTBROWSER)
}

func IsInstallCampaignType(campaignType int32) bool {
	return campaignType == int32(rtb.LinkType_APPSTORE) || campaignType == int32(rtb.LinkType_GOOGLEPLAY) || campaignType == int32(rtb.LinkType_APKDOWNLOAD)
}

func renderNewIvEndcardAndVideoTemplate(rv *RV, r *mvutil.RequestParams, ad *Ad) {
	// 限制新插屏
	if !mvutil.IsNewIv(&r.Param) {
		return
	}
	// 全屏视频+图片 以及 全屏视频 走旧逻辑 （视频模版为902/904，endcard模版为1302）
	if (r.Param.AdspaceType == mvconst.AdSpaceTypeFullScreen) &&
		(r.Param.MaterialType == mvconst.MaterialTypeVideoEndcard || r.Param.MaterialType == mvconst.MaterialTypeVideo) {
		return
	}
	templateMap, _ := extractor.GetTEMPLATE_MAP()
	adnetDefaultValue := extractor.GetADNET_DEFAULT_VALUE()
	// 全屏图片 只下发ec模版
	if r.Param.AdspaceType == mvconst.AdSpaceTypeFullScreen && r.Param.MaterialType == mvconst.MaterialTypeEndcard {
		r.Param.EndcardTplId = "404"
		if fullScreenEndcardTplId, ok := adnetDefaultValue["fullScreenEndcardTplId"]; ok && len(fullScreenEndcardTplId) > 0 {
			r.Param.EndcardTplId = fullScreenEndcardTplId
		}
		subTemplateMap, ok := templateMap.DiverseEndScreen[r.Param.EndcardTplId]
		if !ok {
			return
		}
		ad.EndcardUrl = renderUrlAttachSchema(r.Param.HTTPReq, subTemplateMap.FullScreen)
		// 清空视频模版
		rv.TemplateUrl = ""
		rv.VideoTemplate = 0
		// 如果三方dsp返回了视频，需要遵循MaterialType配置，只返回图片
		ad.VideoURL = ""
		// 纯图片需要返回，置为2，不然sdk 会渲染失败
		ad.PlayableAdsWithoutVideo = 2
		return
	}

	// 半屏视频+ec 或 半屏视频。只下发视频模版
	if r.Param.AdspaceType == mvconst.AdSpaceTypeHalfScreen &&
		(r.Param.MaterialType == mvconst.MaterialTypeVideo || r.Param.MaterialType == mvconst.MaterialTypeVideoEndcard) {
		videoTemplateId := "102"
		if halfScreenVideoTplId, ok := adnetDefaultValue["halfScreenVideoTplId"]; ok && len(halfScreenVideoTplId) > 0 {
			videoTemplateId = halfScreenVideoTplId
		}
		subTemplateMap, ok := templateMap.DiverseVideo[videoTemplateId]
		if !ok {
			return
		}
		rv.VideoTemplate, _ = strconv.Atoi(videoTemplateId)
		rv.TemplateUrl = renderUrlAttachSchema(r.Param.HTTPReq, subTemplateMap.HalfScreen)
		// 清空endcard模版
		r.Param.EndcardTplId = ""
		ad.EndcardUrl = ""

		// 如果三方dsp返回了图片，需要遵循MaterialType配置，只返回视频
		if r.Param.MaterialType == mvconst.MaterialTypeVideo {
			ad.ImageURL = ""
		}
		return
	}

	// 半屏图片 只下发ec模版
	if r.Param.AdspaceType == mvconst.AdSpaceTypeHalfScreen && r.Param.MaterialType == mvconst.MaterialTypeEndcard {
		r.Param.EndcardTplId = "404"
		if halfScreenEndcardTplId, ok := adnetDefaultValue["halfScreenEndcardTplId"]; ok && len(halfScreenEndcardTplId) > 0 {
			r.Param.EndcardTplId = halfScreenEndcardTplId
		}
		subTemplateMap, ok := templateMap.DiverseEndScreen[r.Param.EndcardTplId]
		if !ok {
			return
		}
		ad.EndcardUrl = renderUrlAttachSchema(r.Param.HTTPReq, subTemplateMap.HalfScreen)
		// 清空视频模版
		rv.TemplateUrl = ""
		rv.VideoTemplate = 0

		// 如果三方dsp返回了视频，需要遵循MaterialType配置，只返回图片
		ad.VideoURL = ""
		// 纯图片需要返回，置为2，不然sdk 会渲染失败
		ad.PlayableAdsWithoutVideo = 2
		return
	}

}

func renderCamTplUrl(ad *Ad, dspId int64, params *mvutil.Params, corsairCampaign corsair_proto.Campaign) {
	// 切量开关
	adnetConf, _ := extractor.GetADNET_SWITCHS()
	tpDspGetTplABTestRate, ok := adnetConf["tpDspGetTplABTestRate"]
	if !ok || tpDspGetTplABTestRate == 0 {
		return
	}
	randVal := mvutil.GetRandByGlobalTagId(params, mvconst.SALT_THIRDPARTY_DSP_TPL_ABTEST, 100)
	if tpDspGetTplABTestRate > randVal {
		// 记录切量标记
		params.ThirdPartyABTestRes.ThirdPartyDspTplTag = 1
	} else {
		params.ThirdPartyABTestRes.ThirdPartyDspTplTag = 2
		return
	}

	var imageWidth, imageHeight int
	if corsairCampaign.ImageWidth != nil {
		imageWidth = *corsairCampaign.ImageWidth
	}
	if corsairCampaign.ImageHeight != nil {
		imageHeight = *corsairCampaign.ImageHeight
	}
	templateV2Conf := extractor.GetTEMPLATEMAPV2()
	adnetDefaultConf := extractor.GetADNET_DEFAULT_VALUE()
	switch params.AdType {
	case mvconst.ADTypeSdkBanner:
		if dspId == mvconst.GDTYLH {
			if mvutil.IsEqualProportion(imageWidth, imageHeight, 16, 9) {
				ad.CamTplUrl = GetTemplateUrlFromTemplateV2Conf(mvconst.Banner1002006, mvconst.ApiFrameworkURL, templateV2Conf)
			} else if mvutil.IsEqualProportion(imageWidth, imageHeight, 9, 16) {
				ad.CamTplUrl = GetTemplateUrlFromTemplateV2Conf(mvconst.Banner1002003, mvconst.ApiFrameworkURL, templateV2Conf)
			} else {
				ad.CamTplUrl = GetTemplateUrlFromTemplateV2Conf(mvconst.Banner1002004, mvconst.ApiFrameworkURL, templateV2Conf)
			}
		}

		// 1002005
		if dspId == mvconst.JD || dspId == mvconst.XunfeiDsp || dspId == mvconst.BaiduDsp {
			ad.CamTplUrl = GetTemplateUrlFromTemplateV2Conf(mvconst.Banner1002005, mvconst.ApiFrameworkURL, templateV2Conf)
		}

		if dspId == mvconst.Tanx || dspId == mvconst.InmobiDSP {
			ad.CamTplUrl = GetTemplateUrlFromTemplateV2Conf(mvconst.Banner1002006, mvconst.ApiFrameworkURL, templateV2Conf)
		}
	case mvconst.ADTypeSplash:
		if dspId == mvconst.Tanx {
			if params.FormatOrientation == mvconst.ORIENTATION_PORTRAIT && imageWidth >= imageHeight {
				// 横屏 （暂时没有id）
				splash2TplUrl, ok := adnetDefaultConf["splash2TplUrl"]
				if ok {
					ad.CamTplUrl = splash2TplUrl
				}
			} else {
				// 竖屏 10002001
				ad.CamTplUrl = GetTemplateUrlFromTemplateV2Conf(mvconst.Splash10002001, mvconst.ApiFrameworkZIP, templateV2Conf)
			}
		}
		if dspId == mvconst.JD || dspId == mvconst.XunfeiDsp || dspId == mvconst.BaiduDsp || dspId == mvconst.InmobiDSP {
			ad.CamTplUrl = GetTemplateUrlFromTemplateV2Conf(mvconst.Splash10002001, mvconst.ApiFrameworkZIP, templateV2Conf)
		}
	case mvconst.ADTypeNativeH5:
		if dspId == mvconst.BaiduDsp || dspId == mvconst.Tanx {
			ad.CamTplUrl = GetTemplateUrlFromTemplateV2Conf(mvconst.NativeH59002002, mvconst.ApiFrameworkZIP, templateV2Conf)
		}
	case mvconst.ADTypeInterstitialSdk:
		if dspId == mvconst.BaiduDsp || dspId == mvconst.Tanx || dspId == mvconst.InmobiDSP {
			interstitialTplUrl, ok := adnetDefaultConf["interstitialTplUrl"]
			if ok {
				ad.CamTplUrl = interstitialTplUrl
			}
		}
	}
	// 加上协议头, 都用https
	ad.CamTplUrl = helpers.GenFullUrl(2, ad.CamTplUrl)
}

func GetTemplateUrlFromTemplateV2Conf(tmpId string, apiFramework string, conf map[string]map[string][]*mvutil.TemplateWeightMap) string {
	var templateUrl string
	tmpIdConf, ok := conf[tmpId]
	if !ok {
		return templateUrl
	}
	apiFrameworkConf, ok := tmpIdConf[apiFramework]
	if !ok {
		return templateUrl
	}
	templateUrl, _ = getUrlByWeight(apiFrameworkConf)
	return templateUrl
}

func renderCretiveDomainAbTest(rv *RV, params *mvutil.Params, ad *Ad) {
	tplCreativeDomainConf := extractor.GetTPL_CREATIVE_DOMAIN_CONF()
	if tplCreativeDomainConf == nil {
		return
	}
	abtestTagMap := make(map[string]int)
	for domainMacro, domainWeightMaps := range tplCreativeDomainConf {
		if len(domainWeightMaps) == 0 {
			continue
		}
		// 选择宏需要替换的域名
		domain, id := getDomainByWeight(domainWeightMaps)
		// 记录标记
		abtestTagMap[domainMacro] = id
		rv.TemplateUrl = strings.ReplaceAll(rv.TemplateUrl, "__"+domainMacro+"__", domain)
		ad.EndcardUrl = strings.ReplaceAll(ad.EndcardUrl, "__"+domainMacro+"__", domain)
	}
	// 记录到ThirdPartyABTestRes 中
	if len(abtestTagMap) > 0 {
		jsonByte, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(abtestTagMap)
		params.ThirdPartyABTestRes.TplCreativeDomainTag = string(jsonByte)
	}
}

func renderRvTemplateUrlAbTest(rv *RV, params *mvutil.Params, tplId int) {
	replaceTplconf := extractor.GetREPLACE_TEMPLATE_URL_CONF()
	if replaceTplconf == nil {
		return
	}
	if replaceTplconf.RvTemplate == nil {
		return
	}
	// 获取对应模版id的配置
	tplIdConf, ok := replaceTplconf.RvTemplate[strconv.Itoa(tplId)]
	if !ok || len(tplIdConf) == 0 {
		return
	}
	// 根据权重选择url
	replaceUrl, id := getUrlByWeight(tplIdConf)
	// 记录切量标记
	params.ThirdPartyABTestRes.RvTplId = id
	// 为空值表示对照组
	if len(replaceUrl) == 0 {
		return
	}
	// 替换url
	rv.TemplateUrl = renderUrlAttachSchema(params.HTTPReq, replaceUrl)
}

func renderEndcardTemplateUrlAbTest(params *mvutil.Params, tplId string, ad *Ad) {
	replaceTplconf := extractor.GetREPLACE_TEMPLATE_URL_CONF()
	if replaceTplconf == nil {
		return
	}
	if replaceTplconf.Endcard == nil {
		return
	}
	// 获取对应模版id的配置
	tplIdConf, ok := replaceTplconf.Endcard[tplId]
	if !ok || len(tplIdConf) == 0 {
		return
	}
	// 根据权重选择url
	replaceUrl, id := getUrlByWeight(tplIdConf)
	// 记录切量标记
	params.ThirdPartyABTestRes.EndcardTplId = id
	// 为空值表示对照组
	if len(replaceUrl) == 0 {
		return
	}
	// 替换url
	ad.EndcardUrl = renderUrlAttachSchema(params.HTTPReq, replaceUrl)
}

func getUrlByWeight(confs []*mvutil.TemplateWeightMap) (string, int) {
	urlMap := make(map[string]int)
	IdMap := make(map[string]int)
	for _, tplMap := range confs {
		if tplMap == nil {
			continue
		}
		urlMap[tplMap.Url] = tplMap.Weight
		IdMap[tplMap.Url] = tplMap.Id
	}
	newUrl := mvutil.RandByRate2(urlMap)
	id := IdMap[newUrl]
	return newUrl, id
}

func getDomainByWeight(confs []*mvutil.TemplateCreativeDomainMap) (string, int) {
	domainMap := make(map[string]int)
	IdMap := make(map[string]int)

	for _, domainConf := range confs {
		if domainConf == nil {
			continue
		}
		domainMap[domainConf.Domain] = domainConf.Weight
		IdMap[domainConf.Domain] = domainConf.Id
	}
	domain := mvutil.RandByRate2(domainMap)
	id := IdMap[domain]
	return domain, id
}

// func renderTtAlgo(r *mvutil.RequestParams, corsairCampaign corsair_proto.Campaign) (camInfo *smodel.CampaignInfo, err error) {
// 	//var campaignInfo smodel.CampaignInfo
// 	if r.Param.BackendID != mvconst.TouTiao {
// 		return nil, errors.New("params error")
// 	}
// 	testMap, _ := extractor.GetTOUTIAO_ALGO()
// 	cidInt64, _ := strconv.ParseInt(corsairCampaign.CampaignId, 10, 64)
// 	mvList, ok := testMap[cidInt64]
// 	if !ok || len(mvList) <= 0 {
// 		return nil, errors.New("params error")
// 	}
// 	campaigns, err := GetCampaignInfo(mvList)
// 	if err != nil {
// 		return nil, err
// 	}
// 	// 拿第一个单子
// 	mvcid := int64(0)
// 	for _, v := range mvList {
// 		campaign, ok := campaigns[v]
// 		if ok {
// 			mvcid = v
// 			camInfo = campaign
// 			break
// 		}
// 	}
// 	if mvcid == int64(0) {
// 		return nil, errors.New("not found")
// 	}
// 	// rand
// 	randInt := mvutil.GetRandConsiderZero(r.Param.GAID, r.Param.IDFA, mvconst.SALT_TTALGO, 100)
// 	if randInt == -1 {
// 		return nil, errors.New("rand not hit")
// 	}
// 	// 如果是对照组
// 	if randInt >= 10 {
// 		r.Param.Extra = "adnet_tt1"
// 	} else {
// 		r.Param.Extra = "adnet_tt2"
// 	}
// 	return camInfo, nil
// }

// func renderTtAlgoInfo(ad *Ad, campaign *smodel.CampaignInfo, r *mvutil.RequestParams) {
// 	if r.Param.Extra != "adnet_tt2" {
// 		return
// 	}
// 	var corsairCampaign corsair_proto.Campaign
// 	corsairCampaign.CampaignId = strconv.FormatInt(campaign.CampaignId, 10)
// 	corsairCampaign.AdSource = enum.ADSource_APIOFFER
// 	mvad, _ := RenderCampaign(r, corsairCampaign, campaign)
// 	ad.ImpressionURL = mvad.ImpressionURL
// 	ad.ClickURL = mvad.ClickURL
// 	ad.NoticeURL = mvad.NoticeURL
// 	ad.AdTrackingPoint = mvad.AdTrackingPoint
// 	ad.AdvImp = mvad.AdvImp
// 	ad.AdURLList = mvad.AdURLList
// 	ad.ClickMode = mvad.ClickMode

// 	r.Param.ExtthirdCid = corsairCampaign.CampaignId
// }

func renderUrlAttachSchema(httpReq int32, query string) string {
	if len(query) == 0 {
		return ""
	}

	if strings.HasPrefix(query, "http") {
		return query
	}
	schema := "http"
	if httpReq == int32(2) {
		schema = "https"
	}
	return schema + "://" + query
}

func RenderThirdAdtracking(ad *Ad, corsairCampaign corsair_proto.Campaign, params *mvutil.Params, backendId int32, r *mvutil.RequestParams) {
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

	// 针对ads接口切第三方广告主的情况
	if params.RequestType == mvconst.REQUEST_TYPE_OPENAPI_AD && !params.IsVast {
		adtracking.Click = append(adtracking.Click, ad.NoticeURL)
	}

	//广点通安卓使用adtracing.fcb==3 加广点通的图标
	// if (backendId == mvconst.Gdt || backendId == mvconst.JinShanYun) && params.Platform == mvconst.PlatformAndroid {
	// 	adtracking.Fcb = 3
	// }

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

	ad.AdTracking = adtracking
	ad.AdTrackingPoint = &adtracking
}

func renderThirdPartyOnlineApiBidPrice(params *mvutil.Params, r *mvutil.RequestParams, ad *Ad) {
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
		if dspPriceRft, err := GetAdxPrice(r, nil); err == nil {
			ad.OnlineApiBidPrice = dspPriceRft
		}
	}
	if ad.OnlineApiBidPrice > 0 {
		// 控制bid price 精度
		ad.OnlineApiBidPrice, _ = strconv.ParseFloat(fmt.Sprintf("%.8f", ad.OnlineApiBidPrice), 64)
		// 记录bid_price，tracking使用
		params.OnlineApiBidPrice = strconv.FormatFloat(ad.OnlineApiBidPrice, 'f', 8, 64)
	}
}

func addToponBidPriceOnImpressionUrl(ad *Ad, r *mvutil.RequestParams) {
	if r.IsTopon && len(r.Param.OnlineApiBidPrice) > 0 && r.Param.OnlineApiBidPrice != "0.00000000" {
		ad.ImpressionURL += "&bid_p=" + r.Param.OnlineApiBidPrice
	}
}

func RenderThirdPartyUrls(params *mvutil.Params, ad *Ad, r *mvutil.RequestParams) {
	RenderMWParams(params, ad, false)
	mwCreative := make([]int32, 0)
	mwCreative = append(mwCreative, int32(ad.Rv.VideoTemplate))
	// note 三方胜出处理，序列化中没有extId，需要添加extId			   会带给tracking的字段，新增字段要放在尾巴，不然tracking那边也要改很多
	queryMP := mvutil.SerializeMP(params, mvutil.SerializeMwCreative(mwCreative, params.RequestID))
	queryK := params.RequestID

	if params.IsNewUrl {
		ad.ImpressionURL = "{sh}://{do}/impression?k={k}&mp={mp}"
		if len(r.Param.ReplacedImpTrackDomain) > 0 {
			ad.ImpressionURL = "{sh}://{ido}/impression?k={k}&mp={mp}"
		}
		ad.NoticeURL = "{sh}://{do}/click?k={k}&mp={mp}"
		if len(r.Param.ReplacedClickTrackDomain) > 0 {
			ad.NoticeURL = "{sh}://{cdo}/click?k={k}&mp={mp}"
		}
		if ad.AKS == nil {
			ad.AKS = make(map[string]string)
		}

		ad.AKS["mp"] = mvutil.UrlEncode(queryMP)
		ad.AKS["k"] = queryK
		return
	}

	scheme := GetUrlScheme(params)
	urlQuery := "k=" + queryK + "&mp=" + mvutil.UrlEncode(queryMP)
	// impression
	ad.ImpressionURL = scheme + "://" + params.Domain + "/impression?" + urlQuery
	// // 若为pubnative平台，则把返回的noticeURL放入h参数
	// if params.BackendID == mvconst.PubNative {
	// 	urlQuery = urlQuery + "&h=" + mvutil.UrlEncode(ad.NoticeURL)
	// }
	ad.NoticeURL = scheme + "://" + params.Domain + "/click?" + urlQuery
}

func RenderPVQuery(params *mvutil.Params) string {
	RenderMWParamsForPv(params)
	queryMP := mvutil.SerializeMP(params, mvutil.SerializeMwCreative(params.MwCreatvie, params.RequestID))
	return "k=" + params.RequestID + "&mp=" + url.QueryEscape(queryMP)
}

func RenderMWParamsForPv(params *mvutil.Params) {
	if len(params.BackendList) > 0 {
		var bList []string
		for _, v := range params.BackendList {
			bList = append(bList, strconv.FormatInt(int64(v), 10))
		}
		params.MWadBackend = strings.Join(bList, ",")
	} else {
		params.MWadBackend = "0"
	}
	params.MWadBackendData = "0"
	params.MWbackendConfig = params.AdBackendConfig
}

func RenderMWParams(params *mvutil.Params, ad *Ad, isPv bool) {
	if isPv {
		RenderMWParamsForPv(params)
	} else {
		params.MWadBackend = strconv.Itoa(int(params.BackendID))
		cidStr := strconv.FormatInt(ad.CampaignID, 10)
		params.MWadBackendData = params.MWadBackend + ":" + cidStr + ":1"
		content := "2"
		if len(ad.VideoURL) > 0 {
			content = "3"
		}
		params.MWbackendConfig = params.MWadBackend + ":" + params.RequestKey + ":" + content
	}
}

func isReturnRv(unitId int64, dspId int64) bool {
	noRvUnits, ifFind := extractor.GetTP_WITHOUT_RV()
	if ifFind {
		if mvutil.InInt64Arr(unitId, noRvUnits) {
			return false
		}
	}
	// 如果是自有DSP，不返回rv
	if dspId == mvconst.MVDSP {
		return false
	}
	return true
}

func renderThirdPartyEndcardProperty(ad *Ad, r *mvutil.RequestParams, dspId int64) {
	// 因是第三方的offer，alac，alecfc相对不可控，都写死为关闭，mof获取unit维度配置，默认开启
	if len(ad.EndcardUrl) == 0 {
		return
	}

	//more offer
	endcardUrl, err := renderMoreOfferProperty(r.UnitInfo, &r.Param, ad.EndcardUrl, dspId)
	if err == nil {
		ad.EndcardUrl = endcardUrl
	}
	//specail ec support
	endcardUrl, err = renderSpecalEcProperty(r.UnitInfo, &r.Param, ad.EndcardUrl, dspId, ad.PackageName)
	if err == nil {
		ad.EndcardUrl = endcardUrl
	}
	// endcard_url 下发参数abtest(主要用于下毒功能，及全量上线。abtest还是以mv广告为准)
	var campaign smodel.CampaignInfo
	renderEndcardUrlParamABTest(ad, &campaign, &r.Param, dspId)
}

func renderSpecalSplashProperty(unitInfo *smodel.UnitInfo, param *mvutil.Params, dspId int64, ad *Ad) {
	if unitInfo == nil || param == nil || len(ad.CamTplUrl) == 0 || dspId <= 0 || param.AdType != mvconst.ADTypeSplash {
		return
	}
	dspCfg, ok := extractor.GetDspConfig(dspId)
	if !ok {
		return
	}
	// 关闭全局可点
	var CloseAlecfc bool
	// mongop 里的含义：1表示开 0表示关
	// 对于h5 【splash 】上alecfc=0为不开启全局可点，alecfc=其余值或alecfc不存在均为全局可点
	if unitInfo.Unit.Aladfc != nil && *unitInfo.Unit.Aladfc == 0 {
		CloseAlecfc = true
	}
	if rate, ok := dspCfg.Target.SDKSplashTemplate[mvconst.SplashFullClick]; ok {
		randRate := rand.Intn(100)
		if randRate >= rate {
			CloseAlecfc = true
		}
	}

	if CloseAlecfc {
		ad.CamTplUrl += "&alecfc=0"
	}

	// mongop 里的含义：1表示隐藏按钮 0表示显示按钮
	// 对于H5 hdbtn=1表示隐藏按钮，hdbtn=0或者其余值或alecfc不存在均为显示按钮
	// 当不开启全局可点的时候，必须保证按钮常驻展示，避免页面没有可点击区域
	if unitInfo.Unit.Hdbtn != nil && *unitInfo.Unit.Hdbtn == 1 && !CloseAlecfc {
		ad.CamTplUrl += "&hdbtn=1"
	}

}

func renderSpecalEcProperty(unitInfo *smodel.UnitInfo, param *mvutil.Params, endCardUrl string, dspId int64, packageName string) (string, error) {
	if unitInfo == nil || param == nil || len(endCardUrl) == 0 || dspId <= 0 {
		return "", errors.New("renderMoreOfferProperty params error")
	}
	dspCfg, ok := extractor.GetDspConfig(dspId)
	if !ok {
		return endCardUrl, nil
	}
	if dspCfg.Target == nil || len(dspCfg.Target.SDKEcTemplate) == 0 {
		return endCardUrl, nil
	}

	//fullclick
	if unitInfo.Unit.Alecfc != nil && *unitInfo.Unit.Alecfc == 1 {
		if rate, ok := dspCfg.Target.SDKEcTemplate[mvconst.FullClick]; ok {
			randRate := rand.Intn(100)
			if rate > randRate {
				endCardUrl += "&alecfc=1"
			}
		}
	}

	if param.Platform != mvconst.PlatformIOS || len(packageName) == 0 {
		return endCardUrl, nil
	}
	//autoJump
	if unitInfo.Unit.Alac != nil && *unitInfo.Unit.Alac == 1 {
		if rate, ok := dspCfg.Target.SDKEcTemplate[mvconst.AutuJump]; ok {
			randRate := rand.Intn(100)
			if rate > randRate {
				endCardUrl += "&alac=1"
			}
		}
	}

	return endCardUrl, nil
}

func renderMoreOfferProperty(unitInfo *smodel.UnitInfo, param *mvutil.Params, endCardUrl string, dspId int64) (string, error) {
	if unitInfo == nil || param == nil || len(endCardUrl) == 0 {
		return "", errors.New("renderMoreOfferProperty params error")
	}
	dspCfg, ok := extractor.GetDspConfig(dspId)
	if !ok {
		return endCardUrl, nil
	}
	if dspCfg.Target == nil || len(dspCfg.Target.SDKEcTemplate) == 0 {
		return endCardUrl, nil
	}

	mof := 1
	if unitInfo.Unit.Mof != nil {
		mof = *unitInfo.Unit.Mof
	}

	if moreOfferSwitch, mofOk := dspCfg.Target.SDKEcTemplate[mvconst.MoreOffer]; mofOk && mof == 1 && moreOfferSwitch == 1 {
		// dsp配置开启，及unit维度配置开启，才下发moreoffer标识
		endCardUrl = endCardUrl + "&mof=1"

		if param.MofAbFlag {
			endCardUrl += "&mof_ab=1"
		}
		if len(param.CtnSizeTag) > 0 {
			endCardUrl += "&ctnsize=" + param.CtnSizeTag
		}
	}
	if param.NewMoreOfferFlag {
		endCardUrl += "&mof_uid=" + strconv.FormatInt(unitInfo.Unit.MofUnitId, 10)
	}
	if param.NewMofImpFlag {
		endCardUrl += "&n_imp=1"
	}

	// dsp配置开启，及unit维度配置开启，才下发关闭场景广告标识
	if clsdSwitch, clsdOk := dspCfg.Target.SDKEcTemplate[mvconst.CloseButtonAd]; clsdOk && len(param.ExtDataInit.CloseAdTag) > 0 && clsdSwitch == 1 {
		endCardUrl += "&clsad=" + param.ExtDataInit.CloseAdTag
	}
	// 返回流量对应的包名
	if len(param.PackageName) > 0 {
		endCardUrl += "&mof_pkg=" + param.PackageName
	}
	// 下发自有id
	if len(param.ExtSysId) > 0 && param.ExtSysId != "," && param.ApiVersion < mvconst.API_VERSION_2_2 {
		endCardUrl += "&n_bc=" + url.QueryEscape(param.ExtSysId)
	}
	// 下发region，pl 没有分region，需要adnet和pioneer返回
	if len(param.ExtcdnType) > 0 {
		endCardUrl += "&n_region=" + param.ExtcdnType
	}
	return endCardUrl, nil
}

func renderClickmode(ad *Ad, params *mvutil.Params, corsairCampaign corsair_proto.Campaign, r *mvutil.RequestParams) {
	// 获取原始的campaignid
	if corsairCampaign.OriCampaignId == nil {
		return
	}
	campaignId, err := strconv.ParseInt(*(corsairCampaign.OriCampaignId), 10, 64)
	if err != nil {
		return
	}
	// 获取此campaignId的信息
	camInfo, ifind := extractor.GetCampaignInfo(campaignId)
	if !ifind {
		return
	}
	clickmodeConf := getJumpTypeArr(r, camInfo, params)
	var canCM6 bool
	for jumpType, rate := range clickmodeConf {
		// 无论安卓和ios，若能做clickmode 6或者12，则认为可以走clickmode6的方式
		if (jumpType == mvconst.JUMP_TYPE_CLIENT_DO_ALL || jumpType == mvconst.JUMP_TYPE_CLIENT_DO_ALL_PING_SERVER) && rate != 0 {
			canCM6 = true
			continue
		}
	}
	// 若不能做,则不参与abtest
	if !canCM6 {
		return
	}
	// 判断abtest 切量
	adnConf, _ := extractor.GetADNET_SWITCHS()
	if clickmodeTestRate, ok := adnConf["cmTestRate"]; ok {
		if clickmodeTestRate > rand.Intn(100) {
			ad.ClickMode, _ = strconv.Atoi(mvconst.JUMP_TYPE_CLIENT_DO_ALL)
			// 记录实验标记
			params.ThirdPartyABTestRes.ClickmodeRes = 1
		} else {
			// 未命中标记
			params.ThirdPartyABTestRes.ClickmodeRes = 2
		}

	}
}

func renderABTestRes(params *mvutil.Params, r *mvutil.RequestParams, ad *Ad) {
	params.ThirdPartyABTestRes.CNTrackDomain = r.Param.ExtDataInit.CNTrackDomain
	params.ThirdPartyABTestRes.VideoDspTplABTest = r.Param.VideoDspTplABTestId
	params.ThirdPartyABTestRes.EndcardDspTplABTest = r.Param.EndcardDspTplABTestId
	//赋值给MP参数
	params.ThirdPartyABTestRes.PriceFactor = r.Param.ExtDataInit.PriceFactor
	params.ThirdPartyABTestRes.PriceFactorGroupName = r.Param.ExtDataInit.PriceFactorGroupName
	params.ThirdPartyABTestRes.PriceFactorTag = r.Param.ExtDataInit.PriceFactorTag
	params.ThirdPartyABTestRes.PriceFactorFreq = r.Param.ExtDataInit.PriceFactorFreq
	params.ThirdPartyABTestRes.PriceFactorHit = r.Param.ExtDataInit.PriceFactorHit
	//
	params.ThirdPartyABTestRes.ImpressionCap = r.Param.ExtDataInit.ImpressionCap
	params.ThirdPartyABTestRes.ImpressionCapTime = r.Param.ExtDataInit.ImpressionCapTime
	// adnet生成的三方dsp 视频素材id
	params.ThirdPartyABTestRes.TPDspVideoCrId = r.Param.ThirdPartyDspVideoCreativeId
	params.ThirdPartyABTestRes.YLHHit = r.Param.YLHHit
	//v5
	params.ThirdPartyABTestRes.V5AbtestTag = r.Param.ExtDataInit.V5AbtestTag

	params.ThirdPartyABTestRes.BandWidth = r.Param.BandWidth
	params.ThirdPartyABTestRes.TKSysTag = r.Param.ExtDataInit.TKSysTag
	params.ThirdPartyABTestRes.CdnTrackingDomainABTestTag = r.Param.ExtDataInit.CdnTrackingDomainABTestTag

	// 限制版本
	if mvutil.IsNewIv(params) {
		params.ThirdPartyABTestRes.AdspaceType = r.UnitInfo.Unit.AdSpaceType
		params.ThirdPartyABTestRes.MaterialType = r.UnitInfo.Unit.MaterialType

		//creative.linear.MediaFiles 有值表示视频。
		//creative.companion.companionad 里，image url或HTMLResource有值，则表示endcard。
		//endcard里：
		//staticResource -> 图片
		//htmlResource -> html
		if videoAndImage(ad) {
			// 图片+视频
			params.ThirdPartyABTestRes.TemplateType = 1
		} else if OnlyImage(ad) {
			params.ThirdPartyABTestRes.TemplateType = 30
		} else if OnlyVideo(ad) {
			params.ThirdPartyABTestRes.TemplateType = 11
		}
	}

	extDataJson, err := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(params.ThirdPartyABTestRes)
	if err != nil {
		params.ThirdPartyABTestStr = ""
	}
	params.ThirdPartyABTestStr = string(extDataJson)
}

func OnlyVideo(ad *Ad) bool {
	return len(ad.VideoURL) > 0 && len(ad.ImageURL) == 0
}

func OnlyImage(ad *Ad) bool {
	return (len(ad.ImageURL) > 0 && len(ad.VideoURL) == 0) || len(ad.Mraid) > 0
}

func videoAndImage(ad *Ad) bool {
	return len(ad.VideoURL) > 0 && len(ad.ImageURL) > 0
}

// see https://confluence.mobvista.com/pages/viewpage.action?pageId=62739115
// 三方dsp支持切量全局可点
func IsAllowSDKVideoFullClick(dspId int64, r *mvutil.RequestParams) bool {

	//
	if isSupport, err := IsSupportSDKV9Template(r); !isSupport {
		// 不支持 V9
		msg := ""
		if err != nil {
			msg = err.Error()
		}
		mvutil.Logger.Runtime.Warnf("not support V9 template for reason [%s]", msg)
		return false
	}

	// 如果配置了 V9 则走 dsp 配置逻辑
	dspConf, find := extractor.GetDspConfig(dspId)
	if !find {
		// not configuration
		return false
	}
	if dspConf.Target == nil || len(dspConf.Target.SDKVideoTemplate) == 0 {
		// data empty
		return false
	}
	//
	v, ok := dspConf.Target.SDKVideoTemplate[string(mvconst.VideoFullClick)]
	if !ok {
		// data not exist
		return false
	}

	randRate := rand.Intn(100)
	if randRate >= v {
		//
		return false
	}

	// 切量
	return true
}

// adnet优先判断manage unit template的配置
// unitID+orientation+sdkversion+os version维度配置了视频模板并且没有配置全屏可点的话就不支持，其他的情况都支持全屏可点
func IsSupportSDKV9Template(r *mvutil.RequestParams) (bool, error) {

	//
	unitTemplateInfo := r.UnitInfo.TemplateConf

	// portal没有配置的情况下走dsp切量
	if len(unitTemplateInfo) == 0 {
		return true, nil
	}

	// 只有匹配维度都命中, 且配置了模版(不配置表示都支持), 且模版中不包含V9, 才将 hasConfigV9 置为false
	// 因此默认兜底设置为 true
	supportV9 := true

	for _, templateConf := range unitTemplateInfo {

		// 匹配 横竖 屏
		if templateConf.Orientation != r.Param.FormatOrientation {
			// orientation 不匹配
			continue
		}

		// 匹配 OS版本
		if templateConf.OSMin != 0 && templateConf.OSMin > int(r.Param.OSVersionCode) {
			//
			continue
		}
		if templateConf.OSMax != 0 && templateConf.OSMax < int(r.Param.OSVersionCode) {
			//
			continue
		}

		// 匹配 Sdk 版本
		if includeSdkVersion, ok := templateConf.SdkVersion["include"]; ok {
			if exist, _ := IsInSdkVersion(r.Param.FormatSDKVersion, includeSdkVersion); !exist {
				continue
			}
		}
		//
		if excludeSdkVersion, ok := templateConf.SdkVersion["exclude"]; ok {
			if exist, _ := IsInSdkVersion(r.Param.FormatSDKVersion, excludeSdkVersion); exist {
				continue
			}
		}

		// 没配置模版, 表示都支持
		if len(templateConf.VideoTemplate) == 0 {
			continue
		}

		// 配置了模版, 则寻找是否有 v9 模版
		found := false
		for _, videoTemplate := range templateConf.VideoTemplate {
			if videoTemplate.Type == smodel.VideoTemplateTypeV9 {
				// 有 V9 模版, 走下面逻辑
				found = true
				break
			}
		}
		if !found {
			// found
			supportV9 = false
			break
		}
	}

	// 如果没有配置 V9, 则不切
	if supportV9 {
		return true, nil
	} else {
		return false, errors.New(fmt.Sprintf("unit id [%d] has block v9 template for formatOrientation: [%d], os versio code: [%d], sdk version: [%s]",
			r.Param.UnitID,
			r.Param.FormatOrientation,
			r.Param.OSVersionCode,
			r.Param.SDKVersion,
		))
	}
}

// 是否在 Include Sdk Version 中
func IsInSdkVersion(target supply_mvutil.SDKVersionItem, rules []smodel.SdkVersionRule) (bool, error) {
	// 如果没有配置, 则默认为所有
	if len(rules) == 0 {
		return true, nil
	}

	sb := strings.Builder{}

	for _, rule := range rules {
		// 如果配置了, 则在任意一个都算
		MinSdkVersionItem := supply_mvutil.RenderSDKVersion(rule.Min)
		MaxSdkVersionItem := supply_mvutil.RenderSDKVersion(rule.Max)
		isMatch, err := MatchSdkVersion(target, MinSdkVersionItem, MaxSdkVersionItem)
		if err != nil {
			sb.WriteString(err.Error())
		}

		//
		if isMatch {
			return true, nil
		}
	}

	return false, errors.New(fmt.Sprintf("Not InSdkVersions, Error:[%s]", sb.String()))
}

// 例如判断 mi_7.0.3 是否在 [mi_7.0.0, mi_8.0.3] 中
func MatchSdkVersion(target supply_mvutil.SDKVersionItem, min supply_mvutil.SDKVersionItem, max supply_mvutil.SDKVersionItem) (bool, error) {

	// 不能同时为空
	if min.SDKVersionCode == 0 && max.SDKVersionCode == 0 {
		return false, errors.New(fmt.Sprintf("sdk min max both empty, target verion : [%s_%s], min version: [%s_%s], max version: [%s_%s]",
			target.SDKType,
			target.SDKNumber,
			min.SDKType,
			min.SDKNumber,
			max.SDKType,
			max.SDKNumber,
		))
	}
	if min.SDKVersionCode != 0 && target.SDKType != min.SDKType {
		// 前缀不匹配 如 mi != mal
		return false, errors.New(fmt.Sprintf("prefix not match target: [%s_%s], min: [%s_%s]", target.SDKType, target.SDKNumber, min.SDKType, min.SDKNumber))
	}
	if max.SDKVersionCode != 0 && target.SDKType != max.SDKType {
		// 前缀不匹配 如 mi != mal
		return false, errors.New(fmt.Sprintf("prefix not match target: [%s_%s], max: [%s_%s]", target.SDKType, target.SDKNumber, max.SDKType, max.SDKNumber))
	}
	if min.SDKVersionCode != 0 && target.SDKVersionCode < min.SDKVersionCode {
		// min 存在
		return false, errors.New(fmt.Sprintf("target less then min, target: [%d], min [%d]", target.SDKVersionCode, min.SDKVersionCode))
	}
	if max.SDKVersionCode != 0 && target.SDKVersionCode > max.SDKVersionCode {
		// max 存在
		return false, errors.New(fmt.Sprintf("target greater then max, target: [%d], max [%d]", target.SDKVersionCode, max.SDKVersionCode))
	}
	// success
	return true, nil
}

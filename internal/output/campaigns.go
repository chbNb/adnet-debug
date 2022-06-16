package output

import (
	"errors"
	"hash/crc32"
	"math/rand"
	"net/url"
	"os"
	"strconv"
	"strings"

	supply_mvutil "gitlab.mobvista.com/ADN/chasm/module/supply/mvutil"

	jsoniter "github.com/json-iterator/go"
	"gitlab.mobvista.com/ADN/adnet/internal/corsair_proto"
	"gitlab.mobvista.com/ADN/adnet/internal/extractor"
	"gitlab.mobvista.com/ADN/adnet/internal/mvconst"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	"gitlab.mobvista.com/ADN/adnet/internal/protobuf"
	"gitlab.mobvista.com/ADN/adnet/internal/redis"
	"gitlab.mobvista.com/ADN/adnet/internal/uuid"
	"gitlab.mobvista.com/ADN/adnet/internal/watcher"
	"gitlab.mobvista.com/adserver/recommend_protocols/go/ad_server"
	mydecimal "gitlab.mobvista.com/mae/go-kit/decimal"
)

func RenderOutput(r *mvutil.RequestParams, res *corsair_proto.QueryResult_) (*MobvistaResult, error) {
	if r == nil || res == nil {
		return nil, errors.New("RenderOutput Params error")
	}
	var mr MobvistaResult
	mr.Status = 1
	mr.Msg = "success"
	if r.Param.ApiVersion >= mvconst.API_VERSION_2_2 {
		mr.Data.NewVersionSessionID = r.Param.SessionID
	} else {
		mr.Data.SessionID = r.Param.SessionID
	}
	// 判断一定时间范围内，判断是否需要更新EncryptedId
	renderNewEncryptedId(r)
	// 若id无变化，则不需要返回给sdk
	if r.Param.NewEncryptedSysId != r.Param.EncryptedSysId {
		mr.Data.EncryptedSysId = r.Param.NewEncryptedSysId
	}
	if r.Param.NewEncryptedBkupId != r.Param.EncryptedBkupId {
		mr.Data.EncryptedBkupId = r.Param.NewEncryptedBkupId
	}
	mr.Data.ParentSessionID = r.Param.ParentSessionID
	mr.Data.AdType = int(r.Param.AdType)

	hasThird := false
	hasMv := false
	thirdDspId := int64(0)
	if len(r.DspExt) > 0 {
		dspExt, _ := r.GetDspExt()
		thirdDspId = dspExt.DspId
	}
	// render params
	r.Param.Extra4 = res.LogId
	r.Param.Extra7 = int(r.Param.AdNum)
	r.Param.Extra15 = r.Param.TNum
	r.Param.Extra9 = mvutil.RawUrlEncode(r.Param.UserAgent)
	r.Param.MWRandValue = int(res.RandValue)
	r.Param.MWFlowTagID = int(res.FlowTagId)
	r.Param.FlowTagID = int(res.FlowTagId)
	r.Param.ExtflowTagId = r.Param.FlowTagID
	r.Param.AdBackendConfig = res.AdBackendConfig
	r.Param.Extra5 = r.Param.RequestID
	r.Param.ExtpackageName = r.Param.PackageName
	r.Param.Extra13 = int32(r.Param.AdSourceID)
	r.Param.Extra16 = r.Param.Template
	if r.Param.DspMof == 1 {
		r.Param.ExtfinalPackageName = r.Param.PackageName
	} else {
		r.Param.ExtfinalPackageName = r.AppInfo.RealPackageName
	}
	r.Param.Extra = getCloudExtra()
	if r.IsHBRequest {
		r.Param.Extra = r.Param.Algorithm
		cfg := extractor.GetHBCachedAdNumConfig()
		if v, ok := cfg[strings.ToUpper(r.Param.MediationName)]; ok {
			mr.Data.Vcn = v
		} else {
			mr.Data.Vcn = cfg["default"]
		}
	}

	// 对于关闭场景广告，加上后缀方便在报表区分数据
	if mvutil.IsRequestPioneerDirectly(&r.Param) {
		r.Param.Extra += "_pioneer"
		if r.Param.MofType == mvconst.MOF_TYPE_CLOSE_BUTTON_AD ||
			r.Param.MofType == mvconst.MOF_TYPE_CLOSE_BUTTON_AD_MORE_OFFER {
			r.Param.Extra += "_clsad"
		}
	}

	// 判断能否做ctn_size test,只针对adx流量。
	renderCtnSizeABTest(r, res)
	// 三方广告源endcard，视频模版abtest
	renderDspTplABTest(r, thirdDspId)

	// 素材三期打标记
	r.Param.ExtCreativeNew = renderReExtCreativeNewTag(r.Param.NewCreativeFlag, r.Param.AdType)

	// 替换adywind的tracking域名
	replaceAdywindDomain(r)
	// 替换jssdk 新tracking域名
	replaceJssdkNewDomain(r)
	// 替换tracking 域名
	replaceTrackingDomain(r)
	// 处理sdk传来的has_wx（是否安装微信）
	r.Param.ExtData = RenderExtData(r)
	// 3s cn line
	Is3SCNLine(r)
	// webview下毒
	needPrisonForAndroidWebview(r)

	// 返回more_offer的请求域名
	renderMoreOfferRequestDomain(r)

	// 结果记录到extdev里
	r.Param.ExtDeviceId = RenderExtDeviceId(r)

	if thirdDspId != mvconst.MAS && IsRvNvNewCreative(r.Param.NewCreativeFlag, r.Param.AdType) {
		var ads *corsair_proto.BackendAds
		for _, v := range res.AdsByBackend {
			if v.BackendId == mvconst.Mobvista || (v.BackendId == mvconst.MAdx && r.IsFakeAs) {
				ads = v
				break
			}
		}
		RenderUnitEndcardNew(r, ads)
	}

	// 兼容对于切到adx，adx有填充，且ad_type为iv,不使用as返回的情况
	if len(r.Param.EndcardUrl) == 0 {
		RenderUnitEndcard(r)
	}

	var ads []Ad
	rks := make(map[string]string)
	for _, v := range res.AdsByBackend {
		var oneAds []Ad
		// 传递大模板slot map
		if r.Param.BigTemplateFlag && v.BigTemplateInfo != nil {
			r.Param.BigTemplateId = int64(v.BigTemplateInfo.BigTemplateId)
			r.Param.ExtBigTemId = strconv.FormatInt(r.Param.BigTemplateId, 10)
			r.Param.BigTemplateSlotMap = v.BigTemplateInfo.SlotIndexCampaignIdMap
			// 处理大模板abtest框架逻辑，将标记带入到offer维度内。
			renderBigTempalteABTestTag(&r.Param, thirdDspId)
		}

		if v.BackendId == mvconst.Mobvista {
			oneAds = renderMobvistaCampaigns(r, v) // adnet -> as
			if len(oneAds) > 0 {
				hasMv = true
			}
		} else if v.BackendId == mvconst.MAdx && r.IsFakeAs {
			oneAds = renderMobvistaCampaigns(r, v) // adnet -> adx -> as
			r.Param.IsFakeAs = r.IsFakeAs
			if len(oneAds) > 0 {

			}
			// 流程主要看两个打断点的地方
		} else if (v.BackendId == mvconst.MAdx && (thirdDspId == mvconst.MAS || v.DspId == mvconst.MAS)) ||
			(mvutil.IsRequestPioneerDirectly(&r.Param) && v.BackendId == mvconst.Pioneer) {
			// thirdDspId 是这次请求正式返回的dsp id, v.DspId 是mas竞价失败时返回的more ads里的
			isMoreAds := thirdDspId != mvconst.MAS && v.DspId == mvconst.MAS
			oneAds = renderMasCampaigns(r, v, isMoreAds) // adnet -> adx -> pioneer -> as
			renderBannerTemplate(&mr, v, r, thirdDspId)
			renderRequestExtData(&mr, r)
			r.Param.DspExt = r.DspExt
		} else {
			// adnet处理针对三方dsp的cn tracking切量
			thirdPartyCnTrackingABTest(r)
			oneAds = renderThirdPartyCampaigns(r, v) // adnet -> adx -> 外部dsp
			if v.BackendId == mvconst.MAdx {
				r.Param.DspExt = r.DspExt
				// r.Param.PriceFactor = r.PriceFactor
			}
			if len(oneAds) > 0 {
				hasThird = true
			}
			// sdk banner/splash 返回bannerHtml或bannerUrl.splash 返回ad_tpl_url及ad_html
			if mvutil.IsBannerOrSplashOrDI(r.Param.AdType) {
				renderBannerTemplate(&mr, v, r, thirdDspId)
			} else if (v.BannerUrl != nil || v.BannerHtml != nil) && len(oneAds) > 0 {
				// 不是banner还有bannerUrl、bannerHtml的场景： video中支持只返回一个mraid,不返回视频素材
				// http://confluence.mobvista.com/pages/viewpage.action?pageId=24527937
				if v.BannerUrl != nil {
					oneAds[0].Mraid = renderMraidByDsp(v.GetBannerUrl(), thirdDspId)
				}
				if v.BannerHtml != nil {
					oneAds[0].Mraid = renderMraidByDsp(v.GetBannerHtml(), thirdDspId)
				}
				oneAds[0].PlayableAdsWithoutVideo = 2 // 没有视频只有mraid时， 不播放视频
			}
		}

		for k, v := range v.RKS {
			if !strings.HasPrefix(k, "__") || !strings.HasSuffix(k, "__") ||
				strings.Contains(k, "{") || strings.Contains(k, "}") {
				// 特殊限制
				continue
			}

			rks[k] = v // 如果多个demand返回的占位符一致，则可能会被覆盖.目前我们一份流量只会售卖给一个demand，故不存在多个demand之间相互覆盖
		}
		// adsMap[v.BackendId] = oneAds
		ads = append(ads, oneAds...)
		// pv的adbeckend 由原来记录请求的backendid改为有填充的backendid
		// r.Param.BackendList = append(r.Param.BackendList, v.BackendId)
	}
	// 写入降填充标记
	renderReduceFillFlag(&r.Param)

	if len(ads) == 0 {
		return nil, errors.New("EXCEPTION_RETURN_EMPTY")
	}
	mr.Data.Ads = ads
	mr.Data.UnitSize = r.Param.UnitSize
	mr.Data.Template = r.Param.Template

	// 处理大模板下发参数
	renderBigTemplateParams(&mr, r, res, thirdDspId)

	// new url
	if r.Param.IsNewUrl {
		sh := GetUrlScheme(&r.Param)
		rks["sh"] = sh
		rks["do"] = r.Param.Domain

		if len(r.Param.ReplacedImpTrackDomain) > 0 {
			rks["ido"] = r.Param.ReplacedImpTrackDomain
		}
		if len(r.Param.ReplacedClickTrackDomain) > 0 {
			rks["cdo"] = r.Param.ReplacedClickTrackDomain
		}
		queryZ := r.Param.QueryZ
		queryZ = UrlReplace1(queryZ)
		rks["z"] = queryZ
		mr.Data.RKS = rks
	}

	// 因为pv_urls不支持宏替换，所以没法兼容大模版的v5的返回。因此服务端对pv_urls做宏替换处理逻辑        宏替换，就是将{xx}替换成真正的字符串
	replacePvUrlsMacro(&mr)

	// more offer 参数下发, mas不参与，直接使用mas返回的
	if thirdDspId != mvconst.MAS {
		RenderEndscreenProperty(r)
	}

	if len(r.Param.CtnSizeTag) > 0 {
		AddCtnSizeToEndscreen(r)
	}
	// htmlurl endscreenurl onlyimpressionurl
	mr.Data.OnlyImpressionURL = CreateOnlyImpressionUrl(r.Param, r)
	if thirdDspId != mvconst.MAS {
		RenderEndscreenUrlWithInfo(&mr, r, &r.Param)
	} else {
		mr.Data.EndScreenURL = r.Param.EndcardUrl //
	}
	if len(r.Param.ThirdCidList) > 0 && len(r.Param.ExtthirdCid) <= 0 {
		r.Param.ExtthirdCid = strings.Join(r.Param.ThirdCidList, ",")
	}
	if len(r.Param.Extra20) == 0 {
		r.Param.Extra20 = ""
	}
	// 如果有第三方单子，没有adn单子，则置空extra3
	if hasThird && !hasMv {
		//	r.Param.Extra3 = ""
	}

	// interactive
	if r.Param.AdType == mvconst.ADTypeInteractive {
		renderIA(&mr, *r)
	}
	// webview下毒
	if r.Param.NeedWebviewPrison || r.Param.NewIVClearEndScreenUrl {
		mr.Data.EndScreenURL = ""
	}

	if r.Param.TokenRule > 0 && r.IsHBRequest {
		mr.Data.TokenRule = r.Param.TokenRule
	}
	// 下发加密的出价
	if r.IsHBRequest && r.Price > 0 {
		mr.Data.EncryptPrice = mvutil.Base64EncodeWithURLEncoding(strconv.FormatFloat(r.Price, 'f', -1, 64))
	}

	return &mr, nil
}

func replacePvUrlsMacro(mr *MobvistaResult) {
	if len(mr.Data.PvUrls) == 0 {
		return
	}

	var newReplaceUrlList []string
	for _, pvUrl := range mr.Data.PvUrls {
		replacedUrl := replaceReqUrlMacro(pvUrl, mr.Data.RKS)
		newReplaceUrlList = append(newReplaceUrlList, replacedUrl)
	}
	mr.Data.PvUrls = newReplaceUrlList
}

func thirdPartyCnTrackingABTest(r *mvutil.RequestParams) {
	conf := extractor.GetTrackingCNABTestConf()
	if conf.Enable && r.Param.CountryCode == "CN" {
		randVal := int(crc32.ChecksumIEEE([]byte(mvconst.SALT_TRACK_CN_ABTEST+"_"+mvutil.GetGlobalUniqDeviceTag(&r.Param))) % 100)
		if IsTrackingCNABTest(r, conf, randVal) {
			r.Param.ThirdPartyABTestRes.TKCNABTestTag = 1
			r.Param.Domain = conf.Domain
		} else {
			r.Param.ThirdPartyABTestRes.TKCNABTestTag = 2
		}
	}
}

func IsTrackingCNABTest(r *mvutil.RequestParams, c *mvutil.TrackingCNABTestConf, rand int) bool {
	if mvutil.InInt64Arr(r.Param.UnitID, c.Unit.TagList) {
		return rand < c.Unit.Rate
	}
	if mvutil.InInt64Arr(r.Param.AppID, c.App.TagList) {
		return rand < c.App.Rate
	}
	if mvutil.InInt64Arr(r.Param.PublisherID, c.Pub.TagList) {
		return rand < c.Pub.Rate
	}
	if mvutil.InInt64Arr(int64(r.Param.AdType), c.AdType.TagList) {
		return rand < c.AdType.Rate
	}
	if mvutil.InInt64Arr(int64(r.Param.Platform), c.Platform.TagList) {
		return rand < c.Platform.Rate
	}
	return rand < c.Rate
}

func renderMoreOfferRequestDomain(r *mvutil.RequestParams) {
	// 限制走more_offer cache的量才下发域名
	if r.UnitInfo.Unit.MofUnitId == 0 {
		return
	}
	conf := extractor.GetMORE_OFFER_REQUEST_DOMAIN()
	if domain, ok := conf[mvutil.Cloud()+"-"+mvutil.Region()]; ok && len(domain) > 0 {
		r.Param.MofRequestDomain = domain
	}
	// 限制全链路灰度返回 moreoffer 的请求域名是灰度访问地址
	if os.Getenv("FORCE_ADNET_ENDPIONT_FROM_ENV") == "1" && len(os.Getenv("ADNET_ENDPIONT")) > 0 {
		r.Param.MofRequestDomain = os.Getenv("ADNET_ENDPIONT")
	}
	// 针对cc控制。针对埃及屏蔽rayjump.com以及后面域名替换的情况。
	byCCrequestDomainConf := extractor.GetMoreOfferRequestDomainByCountryCodeConf()
	urlInfo, err := url.Parse(r.Param.MofRequestDomain)
	if err != nil {
		return
	}
	hostNameKey := strings.Replace(urlInfo.Hostname(), ".", "_", -1)
	domainConf, ok := byCCrequestDomainConf[hostNameKey]
	if !ok {
		return
	}
	newDomain, ok := domainConf[r.Param.CountryCode]
	if !ok {
		return
	}
	r.Param.MofRequestDomain = strings.Replace(r.Param.MofRequestDomain, urlInfo.Hostname(), newDomain, -1)
}

func renderRequestExtData(mr *MobvistaResult, r *mvutil.RequestParams) {
	var requestExtData RequestExtData
	requestExtData.ParentId = r.Param.ParentId
	requestExtData.MofRequestDomain = r.Param.MofRequestDomain
	mr.Data.RequestExtData = &requestExtData
}

func RenderExtDeviceId(r *mvutil.RequestParams) string {
	// 有ruid才记录
	if len(r.Param.RuId) == 0 && len(r.Param.MappingIdfa) == 0 {
		return ""
	}

	var extDeviceId mvutil.ExtDeviceId
	extDeviceId = mvutil.ExtDeviceId{
		Ruid:        r.Param.RuId,
		MappingIdfa: r.Param.MappingIdfa,
	}
	extDeviceIdStr, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(extDeviceId)
	return string(extDeviceIdStr)
}

func renderDspTplABTest(r *mvutil.RequestParams, thirdDspId int64) {
	// 仅针对sdk，rv，iv走第三方广告源的流量进行abtest
	if !mvutil.IsHbOrV3OrV5Request(r.Param.RequestPath) {
		return
	}
	if r.Param.AdType != mvconst.ADTypeRewardVideo && r.Param.AdType != mvconst.ADTypeInterstitialVideo {
		return
	}
	if !mvutil.IsThirdDspWithRtdsp(thirdDspId) {
		return
	}
	dspTplABTestConf := extractor.GetDSP_TPL_ABTEST_CONF()
	if len(dspTplABTestConf) == 0 {
		return
	}
	// 视频模版abtest
	if videoDspTplABTestConf, ok := dspTplABTestConf["video"]; ok && len(videoDspTplABTestConf) > 0 {
		r.Param.VideoDspTplABTestId = getDspTplABTestId(videoDspTplABTestConf)
	}
	// endcard模版abtest
	if endcardDspTplABTestConf, ok := dspTplABTestConf["endcard"]; ok && len(endcardDspTplABTestConf) > 0 {
		r.Param.EndcardDspTplABTestId = getDspTplABTestId(endcardDspTplABTestConf)
	}
}

func getDspTplABTestId(tplConfs []*mvutil.DspTplAbtest) int {
	tplAbtestMap := make(map[string]int)
	for _, tplConf := range tplConfs {
		tplAbtestMap[strconv.Itoa(tplConf.Id)] = tplConf.Weight
	}
	randVal := rand.Intn(100)
	tplId, _ := mvutil.RandByRateInMap(tplAbtestMap, randVal)
	return tplId
}

func renderBigTemplateParams(mr *MobvistaResult, r *mvutil.RequestParams, res *corsair_proto.QueryResult_, thirdDspId int64) {
	for _, v := range res.AdsByBackend {
		// 仅处理as和mas的情况
		if v.BackendId != mvconst.MAdx && v.BackendId != mvconst.Mobvista {
			continue
		}
		if v.BackendId == mvconst.MAdx && thirdDspId != mvconst.MAS && thirdDspId != mvconst.FakeAdserverDsp {
			continue
		}
		// 若无广告返回，则过滤
		if v.CampaignList == nil || len(v.CampaignList) == 0 {
			continue
		}
		// 处理nscpt。用于控制当前返回的所有广告是否都要求sdk必须加载完素材才可展示
		// 默认为1。1表示需要所有广告都达到ready_rate条件才能展示，0表示任意广告达到ready_rate就可以展示。
		if r.Param.BigTemplateFlag {
			adnConf, _ := extractor.GetADNET_SWITCHS()
			mr.Data.Nscpt = 1
			if closeNscpt, ok := adnConf["closeNscpt"]; ok && closeNscpt == 1 {
				mr.Data.Nscpt = 0
			}
			// 获取大模板id
			mr.Data.BigTemplateId = r.Param.BigTemplateId
			mr.Data.PvUrls = r.Param.BigTempalteAdxPvUrl
		}
		mr.Data.BigTemplateUrl = renderBigTplUrl(v, r, thirdDspId)
	}
}

func renderBigTplUrl(backend *corsair_proto.BackendAds, r *mvutil.RequestParams, thirdDspId int64) string {
	// mas使用返回的大模版url
	if thirdDspId == mvconst.MAS {
		return r.Param.BigTemplateUrl
	}
	templateMapConf, _ := extractor.GetTEMPLATE_MAP()
	// 判断101或为0是否
	if !mvutil.IsBigTemplate(r.Param.BigTemplateId) {
		// 判断是否为v11模版，若是，则需要使用fake的大模版
		if fakeBigTpl, ok := templateMapConf.BigTempalte[mvconst.BigTemplateUrlFake]; ok && len(fakeBigTpl) > 0 && withVideoTemplateV11(backend) {
			// 大模板abtest追加参数
			return GetUrlScheme(&r.Param) + "://" + fakeBigTpl + r.Param.BigTplABTestParamsStr
		}
	}
	if bigTplUrl, ok := templateMapConf.BigTempalte[r.Param.ExtBigTemId]; ok && len(bigTplUrl) > 0 {
		// 大模板abtest追加参数
		return GetUrlScheme(&r.Param) + "://" + bigTplUrl + r.Param.BigTplABTestParamsStr
	}
	return ""
}

func withVideoTemplateV11(backend *corsair_proto.BackendAds) bool {
	for _, v := range backend.CampaignList {
		if v.VideoTemplateId != nil && *v.VideoTemplateId == ad_server.VideoTemplateId_V11_ZIP {
			return true
		}
	}
	return false
}

func renderReduceFillFlag(params *mvutil.Params) {

	params.ExtReduceFillResp = 0
	params.ExtReduceFillReq = 0

	if len(params.ReqBackends) > 0 {
		// 取第一个请求的backendID
		firstBackendId, _ := strconv.Atoi(params.ReqBackends[0])
		firstFillBackendId := 0
		// 返回的填充以第一个为准
		if len(params.FillBackendId) > 0 {
			firstFillBackendId = int(params.FillBackendId[0])
		}

		if len(params.ReqBackends) > 1 {
			// 若请求了多个backend
			//	1.第1个backend 为非online API 则请求计入, 返回看填充情况
			//  2.第1个backend 为online API
			//		a.online API 有填充 , 则请求与返回都不不计入
			//		b.online API 没有填充, 则请求计入, 返回看填充情况
			if mvconst.Mobvista == firstBackendId || mvconst.MAdx == firstBackendId {
				params.ExtReduceFillReq = 1
				// 现在场景不会存在ReqBackends大于1的情况，保险点还是加上直接请求pioneer的逻辑
				if mvconst.Mobvista == firstFillBackendId || mvconst.MAdx == firstFillBackendId || mvconst.Pioneer == firstFillBackendId {
					params.ExtReduceFillResp = 1
				}
			} else { // 如果第一个是online API
				if firstFillBackendId == firstBackendId {
					params.ExtReduceFillResp = 0
					params.ExtReduceFillReq = 0
				} else {
					params.ExtReduceFillReq = 1
					if len(params.FillBackendId) > 0 {
						params.ExtReduceFillResp = 1
					}
				}
			}

		} else if len(params.ReqBackends) == 1 {
			// 若只请求了1个backend
			//	1.为非online API 则请求计入, 返回看填充情况
			//	2.为online API则请求与返回都 不计入
			if mvconst.Mobvista == firstBackendId || mvconst.MAdx == firstBackendId {
				params.ExtReduceFillReq = 1
				if mvconst.Mobvista == firstFillBackendId || mvconst.MAdx == firstFillBackendId || mvconst.Pioneer == firstFillBackendId {
					params.ExtReduceFillResp = 1
				}
			} else {
				params.ExtReduceFillReq = 0
				params.ExtReduceFillResp = 0
			}
		}
	}

	if !extractor.GetONLY_REQUEST_THIRD_DSP_SWITCH() { // 开闭开关是需要Adnet获取SspProfit
		// sspDailyProfitCap
		if value, found := extractor.GetSspProfitDistributionRuleByUnitIdAndCountryCode(params.UnitID, params.CountryCode); found {
			// 限制type 为1的情况才处理。
			if value.Type != mvconst.SspProfitDistributionRuleFixedEcpm {
				return
			}
			v, ok := value.DailyDeficitCap[params.CountryCode]
			if ok && v == 1 {
				params.ExtReduceFillReq = 0
				params.ExtReduceFillResp = 0
			}
		}
	}
}

func renderIA(mr *MobvistaResult, r mvutil.RequestParams) {
	mr.Data.IARst = &(r.Param.IAADRst)
	mr.Data.IAOri = &(r.Param.IAOrientation)

	// ia icon
	rstStr := strconv.Itoa(r.Param.IAADRst)
	iconMap := r.UnitInfo.Unit.EntranceImg
	iconUrl, ok := iconMap[rstStr]
	if ok {
		if NeedSchemeHttps(r.Param.HTTPReq) {
			iconUrl = renderCDNUrl2Https(iconUrl)
		}
	}
	mr.Data.IAIcon = &(iconUrl)
	// ia url
	if r.Param.IAADRst == mvconst.IARST_APPWALL {
		// 如果是appwall
		iaUrl, _ := extractor.GetIA_APPWALL()
		iaUrl = GetUrlScheme(&r.Param) + "://" + iaUrl
		mr.Data.IAUrl = &(iaUrl)
	} else if r.Param.IAADRst == mvconst.IARST_PLAYABLE {
		// 如果是playable
		mr.Data.IAUrl = &(r.Param.IAPlayableUrl)
	}
}

func renderThirdPartyCampaigns(r *mvutil.RequestParams, ads *corsair_proto.BackendAds) []Ad {
	var adList []Ad
	r.Param.BackendID = ads.BackendId
	r.Param.RequestKey = ads.RequestKey

	if len(ads.CampaignList) <= 0 {
		return adList
	}

	for i, v := range ads.CampaignList {
		if v == nil {
			continue
		}

		// 降填充逻辑
		if isReduceFill(r, v, ads.BackendId) {
			continue
		}

		// floor过滤逻辑
		if filterFloor(r, v, ads.BackendId) {
			continue
		}

		ad := RenderThirdPartyCampaign(r, *v, ads.BackendId, i)
		adList = append(adList, ad)
		r.Param.ThirdCidList = append(r.Param.ThirdCidList, v.CampaignId)
	}

	if len(adList) > 0 {
		r.Param.FillBackendId = append(r.Param.FillBackendId, ads.BackendId)
	}
	// pv的adbeckend 由原来记录请求的backendid改为有填充的backendid
	r.Param.BackendList = append(r.Param.BackendList, ads.BackendId)
	return adList
}

type Stringss [][]string

func renderMobvistaCampaigns(r *mvutil.RequestParams, ads *corsair_proto.BackendAds) []Ad {
	// ss 切量
	ssTestRuleByAds(ads)

	// render params
	if ads != nil {
		renderAllParams(r, *ads)
	}
	if r == nil || ads == nil || ads.CampaignList == nil {
		return []Ad{}
	}

	// 修复因为abtest导致algo值缺失问题
	fixSSABTestAlgo(r, ads)

	// var adList []Ad
	campaignList := ads.CampaignList
	listLen := len(campaignList)
	adList := make([]Ad, 0, listLen)
	// 批量获取campaign信息
	campaignIds := make([]int64, 0, listLen)
	extra2List := make([]string, 0, listLen)
	creative1 := make(map[int64]map[string]map[ad_server.CreativeType]int64, listLen)
	creative2 := make(map[int64]map[string]map[ad_server.CreativeType]int64, listLen)
	creativeNew := make(map[int64]map[string]map[ad_server.CreativeType]int64, listLen)
	dcoMaterialMIds := make(map[int64]map[ad_server.CreativeType]int64)
	dcoOfferMaterialOmId := make(map[int64]map[ad_server.CreativeType]int64)
	if r.Param.ExtCampaignTagList == nil {
		r.Param.ExtCampaignTagList = make(map[int64]*mvutil.CampaignTagInfo)
	}
	for _, v := range campaignList {
		campaignID, _ := strconv.ParseInt(v.CampaignId, 10, 64)
		extra2List = append(extra2List, v.CampaignId)
		campaignIds = append(campaignIds, campaignID)

		r.Param.ExtCampaignTagList[campaignID] = &mvutil.CampaignTagInfo{
			CampaignID:        campaignID,
			TemplateGroup:     int64(v.GetTemplateGroup()),
			EndCardTemplateId: int64(v.GetEndCardTemplateId()),
			VideoTemplateId:   int64(v.GetVideoTemplateId()),
			IsReduceFill:      0,
		}
		// 优先判断是否为新结构下发素材，是则从新素材中获取
		if len(v.CreativeDataList) > 0 {
			ctmpMap := make(map[string]map[ad_server.CreativeType]int64, len(v.CreativeDataList))
			crTypeIdMap := make(map[ad_server.CreativeType]int64, 9)
			dcoMIds := make(map[ad_server.CreativeType]int64)
			dcoOmIds := make(map[ad_server.CreativeType]int64)
			for _, val := range v.CreativeDataList {
				crMap := make(map[ad_server.CreativeType]int64, len(v.CreativeDataList))
				for _, crTypeData := range val.CreativeTypeIdList {
					crMap[crTypeData.Type] = crTypeData.CreativeId
					// 封装CreativeTypeIdMap,用于渲染素材信息
					crTypeIdMap[crTypeData.Type] = crTypeData.CreativeId
					// 处理dco素材id,为查询redis做准备
					if crTypeData.MId != nil {
						dcoMIds[crTypeData.Type] = *crTypeData.MId
					}
					if crTypeData.OmId != nil {
						dcoOmIds[crTypeData.Type] = *crTypeData.OmId
					}
					if crTypeData.CpdId != nil {
						// 获取pcdId，记录到r参数中
						v.CpdIds = append(v.CpdIds, strconv.FormatInt(*crTypeData.CpdId, 10))
					}
				}
				ctmpMap[val.DocId] = crMap
				// 封装参数，用于查询素材信息
				creativeNew[campaignID] = ctmpMap
			}
			// dco
			if len(dcoMIds) > 0 {
				dcoMaterialMIds[campaignID] = dcoMIds
			}
			if len(dcoOmIds) > 0 {
				dcoOfferMaterialOmId[campaignID] = dcoOmIds
			}
			v.CreativeTypeIdMap = crTypeIdMap
		} else {
			if v.CreativeId != nil && len(v.CreativeTypeIdMap) > 0 {
				ctmpMap := make(map[string]map[ad_server.CreativeType]int64, 1)
				ctmpMap[*v.CreativeId] = v.CreativeTypeIdMap
				creative1[campaignID] = ctmpMap
			}
			if v.CreativeId2 != nil && len(v.CreativeTypeIdMap2) > 0 {
				ctmpMap := make(map[string]map[ad_server.CreativeType]int64, 1)
				ctmpMap[*v.CreativeId2] = v.CreativeTypeIdMap2
				creative2[campaignID] = ctmpMap
				for k, info := range v.CreativeTypeIdMap2 {
					v.CreativeTypeIdMap[k] = info
				}
			}
		}
	}
	campaigns, err := GetCampaignsInfo(campaignIds)
	if err != nil {
		watcher.AddWatchValue("get_campaign_error_count", float64(1))
		return []Ad{}
	}

	var creativeInfos map[int64]map[ad_server.CreativeType]*protobuf.Creative
	var dcoMaterialInfos map[int64]map[ad_server.CreativeType]*protobuf.Material
	var dcoOfferMaterialInfos map[int64]map[ad_server.CreativeType]*protobuf.OfferMaterial
	if ads.BackendId == mvconst.Mobvista || (ads.BackendId == mvconst.MAdx && r.IsFakeAs) {
		if len(creativeNew) > 0 {
			creativeInfos, err = GetCreativePbInfosV2(creativeNew, r, campaigns)
			// dco获取素材信息
			dcoMaterialInfos = GetDcoMaterialData(dcoMaterialMIds)
			dcoOfferMaterialInfos = GetDcoOfferMaterialData(dcoOfferMaterialOmId)
		} else {
			creativeInfos, err = GetCreativePbInfos(creative1, creative2, r, campaigns)
		}
		if err != nil {
			watcher.AddWatchValue("get_creative_error", float64(1))
			// mvutil.Logger.Runtime.Warnf("request_id=[%s] GetCreativePbInfos error=[%s]", r.Param.RequestID, err.Error())
		}
	}
	adMap := make(map[int64]Ad, listLen)
	var extra20Arr Stringss
	var ExtABTestArr, ExtPlayableArr, ExtBigTplOfferDataArr []string
	for _, v := range campaignList {
		if v == nil {
			continue
		}

		// 降填充逻辑
		if isReduceFill(r, v, ads.BackendId) {
			campaignID, _ := strconv.ParseInt(v.CampaignId, 10, 64)
			r.Param.ExtCampaignTagList[campaignID].IsReduceFill = 1
			continue
		}

		if filterFloor(r, v, ads.BackendId) {
			continue
		}
		campaignID, _ := strconv.ParseInt(v.CampaignId, 10, 64)
		campaign, ok := campaigns[campaignID]
		if !ok {
			continue
		}
		if campaign.CampaignId == int64(0) {
			continue
		}

		var crContent map[ad_server.CreativeType]*protobuf.Creative
		var materialContent map[ad_server.CreativeType]*protobuf.Material
		var offerMaterialContent map[ad_server.CreativeType]*protobuf.OfferMaterial
		if creativeInfos != nil {
			if content, ok := creativeInfos[campaignID]; ok {
				crContent = content
			}
		}
		if dcoMaterialInfos != nil {
			if dcoMaterialInfo, ok := dcoMaterialInfos[campaignID]; ok {
				materialContent = dcoMaterialInfo
			}
		}
		if dcoOfferMaterialInfos != nil {
			if dcoOfferMaterialInfo, ok := dcoOfferMaterialInfos[campaignID]; ok {
				offerMaterialContent = dcoOfferMaterialInfo
			}
		}
		// 将dco 素材信息合到原来的crContent结构中
		crContent = MergeDcoCreativeInfo(crContent, materialContent, offerMaterialContent)
		var ad Ad
		ad, err = RenderCampaignWithCreative(r, *v, campaign, crContent)
		if err != nil {
			continue
		}
		adMap[campaignID] = ad

		extra20Arr = append(extra20Arr, r.Param.Extra20List)
		// 记录请求的playable
		ExtPlayableArr = append(ExtPlayableArr, r.Param.ExtPlayableList)
		if len(r.Param.ExtABTestList) > 0 {
			ExtABTestArr = append(ExtABTestArr, r.Param.ExtABTestList)
		}
		if len(r.Param.ExtBigTplOfferDataList) > 0 {
			ExtBigTplOfferDataArr = append(ExtBigTplOfferDataArr, r.Param.ExtBigTplOfferDataList)
		}
	}
	// 整理extra20
	if len(extra20Arr) > 0 {
		extra20Byte, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(extra20Arr)
		r.Param.Extra20 = string(extra20Byte)

		r.Param.FillBackendId = append(r.Param.FillBackendId, ads.BackendId)
	}
	// 整理request维度的ext_playable
	extPlayableJson, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(ExtPlayableArr)
	r.Param.ExtPlayableArr = string(extPlayableJson)
	// 整理request维度的标记值
	extTagJson, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(renderExtTagList(&r.Param))
	r.Param.ExtTagArrStr = string(extTagJson)
	// 请求日志存储 abtest标记
	if len(ExtABTestArr) > 0 {
		extABTestJson, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(ExtABTestArr)
		r.Param.ABTestTagStr = string(extABTestJson)
	}
	// 整理extra2
	r.Param.Extra2 = strings.Join(extra2List, ",")
	r.Param.Extra6 = len(extra2List)
	// pv的adbeckend 由原来记录请求的backendid改为有填充的backendid
	r.Param.BackendList = append(r.Param.BackendList, ads.BackendId)
	// 记录算法需要的大模板offer侧信息
	r.Param.ExtBigTplOfferData = strings.Join(ExtBigTplOfferDataArr, ";")

	// 按照adserver返回顺序排序
	for _, v := range campaignIds {
		ad, ok := adMap[v]
		if ok {
			adList = append(adList, ad)
		}
	}

	return adList
}

func IsOnlineEcpmUnit(params *mvutil.Params) bool {
	if params.RequestType != mvconst.REQUEST_TYPE_OPENAPI_AD {
		return false
	}
	value, ifFind := extractor.GetSspProfitDistributionRuleByUnitIdAndCountryCode(params.UnitID, params.CountryCode)
	if !ifFind || value == nil {
		return false
	}
	return value.Type == mvconst.SspProfitDistributionRuleOnlineApiEcpm
}

func MergeDcoCreativeInfo(crContent map[ad_server.CreativeType]*protobuf.Creative, materialContent map[ad_server.CreativeType]*protobuf.Material, offerMaterialContent map[ad_server.CreativeType]*protobuf.OfferMaterial) map[ad_server.CreativeType]*protobuf.Creative {
	// material的素材信息赋值到原有的creative结构中
	if materialContent != nil {
		for crType, cr := range materialContent {
			crContent[crType] = &protobuf.Creative{
				Url:             cr.Url,
				VideoLength:     cr.VideoLength,
				VideoSize:       cr.VideoSize,
				VideoResolution: cr.Resolution,
				Width:           cr.Width,
				Height:          cr.Height,
				WatchMile:       100,
				BitRate:         cr.BitRate,
				IValue:          cr.IValue,
				SValue:          cr.SValue,
				FValue:          cr.FValue,
				Resolution:      cr.Resolution,
				Mime:            cr.Mime,
				FMd5:            cr.FMd5,
				Orientation:     cr.Orientation,
				Protocal:        cr.Protocol,
			}
		}
	}
	// offer_material的素材关系信息赋值到原有的creative结构中
	if offerMaterialContent != nil {
		for crType, cr := range offerMaterialContent {
			crContent[crType].Source = cr.Source
			crContent[crType].AdvCreativeId = cr.AdvCreativeId
			crContent[crType].CreativeId = cr.OmID
			crContent[crType].Cname = cr.AdvCname
			crContent[crType].CsetId = cr.CsetId
		}
	}
	return crContent
}

func renderExtTagList(param *mvutil.Params) []string {
	var item []string
	for _, info := range param.ExtCampaignTagList {
		item = append(item,
			strings.Join([]string{
				strconv.FormatInt(info.CampaignID, 10),
				strconv.FormatInt(int64(info.CDNAbTest), 10),
				strconv.FormatInt(info.TemplateGroup, 10),
				strconv.FormatInt(info.EndCardTemplateId, 10),
				strconv.FormatInt(info.VideoCreativeid, 10),
				strconv.FormatInt(info.VideoTemplateId, 10),
				strconv.FormatInt(info.IsReduceFill, 10),
			}, ":"),
		)
	}

	return item

}

func renderAllParams(r *mvutil.RequestParams, ads corsair_proto.BackendAds) {
	if !strings.Contains(r.Param.Extra, "_tt") {
		if r.Param.FBFlag == 1 {
			r.Param.Extra = r.Param.Extra + "_fb1"
		} else if r.Param.FBFlag == 2 {
			r.Param.Extra = r.Param.Extra + "_fb2"
		} else if r.Param.FBFlag == 3 {
			r.Param.Extra = r.Param.Extra + "_fb3"
		} else if r.Param.FBFlag == 4 {
			r.Param.Extra = r.Param.Extra + "_fb4"
		} else if r.IsHBRequest {
			r.Param.Extra = r.Param.Algorithm // hb-${region}-adserver/pioneer
		} else {
			r.Param.Extra = r.Param.Extra + "_adserver"
		}
	}

	// 对于关闭场景广告，加上后缀方便在报表区分数据
	if mvutil.IsAppwallOrMoreOffer(r.Param.AdType) && (r.Param.MofType == mvconst.MOF_TYPE_CLOSE_BUTTON_AD ||
		r.Param.MofType == mvconst.MOF_TYPE_CLOSE_BUTTON_AD_MORE_OFFER) {
		r.Param.Extra += "_clsad"
	}

	if ads.AlgoFeatInfo != nil {
		r.Param.AlgoMap = renderAlgo(*(ads.AlgoFeatInfo))
		r.Param.Extalgo = *(ads.AlgoFeatInfo)
	}
	if len(ads.ExtAdxAlgo) > 0 {
		r.Param.ExtAdxAlgo = ads.ExtAdxAlgo
	}

	r.Param.RespFillEcpmFloor = ads.EcpmFloor

	r.Param.Extra5 = r.Param.RequestID
	r.Param.ExtflowTagId = r.Param.FlowTagID
	r.Param.Extra3 = mvutil.CheckParam(ads.Strategy)
	if ads.IfLowerImp != nil {
		r.Param.ExtifLowerImp = *(ads.IfLowerImp)
	}
	r.Param.BackendID = ads.BackendId
	r.Param.RequestKey = ads.RequestKey
	// interactive
	if ads.ResourceType != nil {
		r.Param.IAADRst = int(*(ads.ResourceType))
		r.Param.Extb2t = r.Param.IAADRst
	}
}

func renderAlgo(algo string) map[int64]string {
	algoMap := make(map[int64]string)
	if len(algo) <= 0 {
		return algoMap
	}
	algoList := strings.Split(algo, "#")
	if len(algoList) <= 0 {
		return algoMap
	}
	for _, v := range algoList {
		algoData := strings.Split(v, ":")
		if len(algoData) != 2 {
			continue
		}
		if len(algoData[0]) <= 0 || len(algoData[1]) <= 0 {
			continue
		}
		campaignID, _ := strconv.ParseInt(algoData[0], 10, 64)
		if campaignID <= int64(0) {
			continue
		}
		algoMap[campaignID] = algoData[1]
	}
	return algoMap
}

// fixSSABTestAlgo 修正因为ss切量导致算法值缺失
func fixSSABTestAlgo(r *mvutil.RequestParams, ads *corsair_proto.BackendAds) {
	campaignList := ads.CampaignList
	if len(r.Param.AlgoMap) == 0 || len(campaignList) == 0 {
		return
	}

	for _, c := range campaignList {
		if c.CampaignId == c.GetRawCampaignId() || c.GetRawCampaignId() == "" {
			continue
		}

		rawCampaignID, _ := strconv.ParseInt(c.GetRawCampaignId(), 10, 64)
		if algoValue, ok := r.Param.AlgoMap[rawCampaignID]; ok {
			campaignID, _ := strconv.ParseInt(c.CampaignId, 10, 64)
			r.Param.AlgoMap[campaignID] = algoValue
		}
	}
}

func replaceAdywindDomain(r *mvutil.RequestParams) {
	// 若为adywind 则替换其tracking域名
	if mvutil.IsMpad(r.Param.RequestPath) {
		// 请求域名与tracking域名的映射关系，不单止adywind
		adywindDomainConf := extractor.GetMP_DOMAIN_CONF()
		if len(adywindDomainConf) <= 0 {
			return
		}
		// 查看此域名是否有配置
		for _, v := range adywindDomainConf {
			if v.SearchDomain != nil && v.ReplaceDomain != nil && strings.Contains(r.Param.MpReqDomain, *v.SearchDomain) {
				r.Param.Domain = *v.ReplaceDomain
				return
			}
		}
	}
}

func renderReExtCreativeNewTag(flag bool, adType int32) string {
	creativeNewTag := ""
	if IsRvNvNewCreative(flag, adType) {
		var extCreativeNew mvutil.ExtCreativeNew
		extCreativeNew.IsCreativeNew = true
		str, err := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(extCreativeNew)
		if err == nil {
			creativeNewTag = string(str)
		}
	}
	return creativeNewTag
}

// android webview >= 72 下毒
func needPrisonForAndroidWebview(r *mvutil.RequestParams) {

	if r.Param.AdType != mvconst.ADTypeRewardVideo && r.Param.AdType != mvconst.ADTypeInterstitialVideo {
		return
	}
	// sdk 流量
	if !mvutil.IsHbOrV3OrV5Request(r.Param.RequestPath) {
		return
	}
	// 19618这个开发者load和show之间必须要求1秒以内，加载模板会耽误渲染时间
	adnetListConf := extractor.GetADNET_CONF_LIST()
	if temPoisonPubList, ok := adnetListConf["temPoisonPubList"]; ok && mvutil.InInt64Arr(r.Param.PublisherID, temPoisonPubList) {
		r.Param.NeedWebviewPrison = true
		return
	}

	if r.Param.Platform != mvconst.PlatformAndroid {
		return
	}

	// 另外一种下毒，这个问题是在部分机型，c++层空指针崩溃，崩溃到了系统chrome内核里面，崩溃率2.5%左右，分析是chromium的bug，java层暂时没有好方案修复
	// 主要是开发者在bugly上统计到的 有关国内定制ROM厂商品牌的手机上 出现的， 集中在Android 5.x和Android 6。
	pubAndOsvConf, _ := extractor.GetANR_PRISION_BY_PUB_AND_OSV()
	if len(pubAndOsvConf.PubIds) > 0 && len(pubAndOsvConf.Osvs) > 0 {
		if mvutil.InInt64Arr(r.Param.PublisherID, pubAndOsvConf.PubIds) && mvutil.InStrArray(r.Param.OSVersion, pubAndOsvConf.Osvs) {
			r.Param.NeedWebviewPrison = true
			return
		}
	}

	// 获取webview版本
	ua := mvutil.UaParser.ParseUserAgent(r.Param.UserAgent)
	if len(ua.Major) == 0 {
		return
	}
	majorInt, _ := strconv.Atoi(ua.Major)
	// >72则不下毒
	if majorInt > 72 {
		return
	} else if majorInt == 72 {
		r.Param.NeedWebviewPrison = true
		return
	}

	// 获取sdk版本信息
	sdkVersionData := supply_mvutil.RenderSDKVersion(r.Param.SDKVersion)
	campareValue, _ := mvutil.VerCampare(sdkVersionData.SDKNumber, "9.9.0")

	if campareValue == -1 {
		// 根据app or pub下毒
		prisonConf, _ := extractor.GetWEBVIEW_PRISON_PUBIDS_APPIDS()
		if appIdConf, ok := prisonConf["appId"]; ok {
			if mvutil.InInt64Arr(r.Param.AppID, appIdConf) {
				r.Param.NeedWebviewPrison = true
				return
			}
		}
		if pubIdConf, ok := prisonConf["pubId"]; ok {
			if mvutil.InInt64Arr(r.Param.PublisherID, pubIdConf) {
				r.Param.NeedWebviewPrison = true
				return
			}
		}
	}

	return
}

func ssTestRuleByAds(ads *corsair_proto.BackendAds) {
	if ads == nil {
		return
	}

	for _, ad := range ads.CampaignList {
		ssTestRule(ad)
	}
}

func ssTestRule(ad *corsair_proto.Campaign) {
	campaignRules, ok := extractor.GetSS_ABTEST_CAMPAIGN()
	if !ok {
		return
	}

	rule, ok := campaignRules[ad.CampaignId]
	if !ok {
		return
	}

	if rule.CampaignId == 0 {
		return
	}

	rawCampaignID := ad.CampaignId
	ad.RawCampaignId = &rawCampaignID
	if rand.Intn(10000) < rule.ABRatio {
		// AB Test
		ad.CampaignId = strconv.Itoa(rule.CampaignId)
		hitSSABTest := true
		ad.SSABTest = &hitSSABTest
		return
	}

	if rand.Intn(10000) < rule.AARatio {
		// AA Test
		hitSSAATest := true
		ad.SSAATest = &hitSSAATest
		return
	}
	return
}

func replaceJssdkNewDomain(r *mvutil.RequestParams) {
	if mvutil.NeedNewJssdkDomain(r.Param.RequestPath, r.Param.Ndm) {
		if r.Param.ExtDataInit.TKCNABTestTag == 1 {
			conf := extractor.GetTrackingCNABTestConf()
			if len(conf.JsDomain) > 0 {
				r.Param.Domain = conf.JsDomain
			}
		} else {
			trackDomainConf := extractor.GetJSSDK_DOMAIN_TRACK()
			if len(trackDomainConf) > 0 {
				r.Param.Domain = trackDomainConf
			}
		}
	}
}

func filterFloor(r *mvutil.RequestParams, v *corsair_proto.Campaign, backendID int32) bool {
	var bidPrice float64 = v.GetBidPrice()
	if backendID == mvconst.MAdx {
		if dspPriceRft, err := GetAdxPrice(r, v); err == nil {
			bidPrice = dspPriceRft
		}
	}
	if r.Param.ReduceFillConfig != nil && r.Param.ReduceFillConfig.ControlMode == mvconst.EcpmFloor &&
		r.Param.BidFloor > 0 {
		if bidPrice < r.Param.BidFloor {
			r.Param.ReduceFill = "floor"
			return true
		}
	}

	// online api 过滤掉低于底价的请求
	if IsOnlineEcpmUnit(&r.Param) {
		// 做一个开发者的黑名单，对于这些开发者，可以不做过滤
		adnConfList := extractor.GetADNET_CONF_LIST()
		if onlineFilterByBidFloorBlackPubList, ok := adnConfList["onlineFilterByBidFloorBlackPubList"]; ok && mvutil.InInt64Arr(r.Param.PublisherID, onlineFilterByBidFloorBlackPubList) {
			return false
		}
		if bidPrice < r.Param.BidFloor {
			r.Param.OnlineFilterByBidFloor = true
			return true
		}
	}
	return false
}

// isReduceFill 降填充逻辑
func isReduceFill(r *mvutil.RequestParams, v *corsair_proto.Campaign, backendID int32) bool {
	if r.IsHBRequest {
		return false
	}

	if len(r.Param.ExtReduceFillValue) > 0 && (backendID == mvconst.MAdx || backendID == mvconst.Mobvista) {
		var log mvutil.ReduceFillLog

		log.IsReduceFill = false
		log.BackendID = backendID
		log.CampaignID = v.CampaignId

		defer func() {
			r.Param.ReduceFillList = append(r.Param.ReduceFillList, renderReduceFillLog(&log))
		}()

		ecpmFloorPrice := r.Param.FillEcpmFloor // 前置到 param_render_filter 处理 redis 的降填充底价（算法负责写入）查询
		// 请求级别有拿到降填充的门限使用请求级别的
		if r.Param.RespFillEcpmFloor > 0 {
			ecpmFloorPrice = r.Param.RespFillEcpmFloor
		}
		// 算法计算生成的降填充底价范围为 [0.00000000001, 100000] 所以不在这个范围的为异常数据
		if ecpmFloorPrice < mvconst.REDUCE_FILL_ECPMFLOOR_MIN || ecpmFloorPrice > mvconst.REDUCE_FILL_ECPMFLOOR_MAX {
			mvutil.Logger.Runtime.Warnf("key[%s] reduce fill ecpm_floor[%s] outside range value", r.Param.FillEcpmFloorKey, r.Param.ExtReduceFillValue)
			return false
		}

		var bidPrice float64 = v.GetBidPrice()
		if backendID == mvconst.MAdx {
			if dspPriceRft, err := GetAdxPrice(r, v); err == nil {
				bidPrice = dspPriceRft
			}
		}
		log.FillPrice = bidPrice
		log.ReduceEcpmFloor = ecpmFloorPrice
		log.ReduceFillKey = r.Param.FillEcpmFloorKey
		log.Version = r.Param.FillEcpmFloorVer

		// TODO
		// 填充率是否正常
		// check ecpmFloorPriceVer && getClickHouseRealTimeFillRate by key
		// 降填充兜底策略
		// 落监控日志触发报警

		// 算法实验
		algoStrategy := strings.Split(r.Param.Extra3, ";")
		if len(algoStrategy) > 9 && algoStrategy[9] == "1" {
			log.IsAlgoExperiment = true
			return false
		}

		// 设备白名单不做降填充
		var deviceKey string
		if r.Param.Platform == mvutil.IOSPLATFORM {
			deviceKey = r.Param.IDFA + "_2"
		} else if r.Param.Platform == mvutil.ANDROIDPLATFORM {
			deviceKey = r.Param.GAID + "_1"
		}
		findDevice, _ := redis.LocalRedisAlgoHExists("dev_white_list", deviceKey)
		if findDevice {
			log.IsWhiteListDev = true
			return false
		}

		// 出价小于 fillrate ecpmFloor 不填充广告
		// config fillrate to 0, ecpmFloorPrice is 100000.
		if bidPrice < ecpmFloorPrice {
			log.IsReduceFill = true
			r.Param.ReduceFill = "fillrate"
			if r.Param.Debug > 0 {
				r.DebugInfo += "is_reduce_fill@@@" + r.Param.FillEcpmFloorKey + ":" + r.Param.ExtReduceFillValue + ":" + strconv.FormatFloat(log.FillPrice, 'E', 6, 64) + "<br>"
			}
			return true
		}
	}

	return false
}

func GetReduceFillSwitch(key string) bool {
	sw, _ := extractor.GetREDUCE_FILL_SWITCH()
	if sw.Total != 0 {
		return true
	}
	// 判断白名单
	if _, ok := sw.WhiteList[key]; ok {
		return true
	}

	return false

}

func GetReduceFillFlinkSwitch(key string) bool {
	sw, _ := extractor.GetREDUCE_FILL_FLINK_SWITCH()
	if sw.Total != 0 {
		return true
	}
	// 判断白名单
	if _, ok := sw.WhiteList[key]; ok {
		return true
	}

	return false

}

func renderReduceFillLog(log *mvutil.ReduceFillLog) string {
	var reduce string = "0"
	if log.IsReduceFill {
		reduce = "1"
	}
	var whiteListDev string = "0"
	if log.IsWhiteListDev {
		whiteListDev = "1"
	}
	var algoExperiment string = "0"
	if log.IsAlgoExperiment {
		algoExperiment = "1"
	}

	attr := strings.Join([]string{
		log.CampaignID,
		strconv.FormatInt(int64(log.BackendID), 10),
		log.ReduceFillKey,
		reduce,
		whiteListDev,
		algoExperiment,
		strconv.FormatFloat(log.FillPrice, 'E', 6, 64),
		strconv.FormatFloat(log.ReduceEcpmFloor, 'E', 6, 64),
		log.Version,
	}, ":")

	return attr
}

func getCloudExtra() string {
	cloud := extractor.GetCLOUD_NAME()
	if len(os.Getenv("POD_NAME")) > 0 {
		cloud = cloud + "k8s"
	}
	return mvutil.GetCloudExtra(cloud)
}

func GetAdxPrice(r *mvutil.RequestParams, v *corsair_proto.Campaign) (float64, error) {
	dspExt, err := r.GetDspExt()

	d1 := mydecimal.NewMDecimal()
	d1.FromFloat64(dspExt.PriceOut) // ADX 返回的税后价格（美分）
	var divisor int64 = 100
	d2 := mydecimal.NewMDecimal()
	d2.FromInt(divisor)
	d3 := mydecimal.NewMDecimal()
	err = mydecimal.Div(d1, d2, d3, 2) // 美分转美元
	if err != nil {
		return 0.0, err
	}

	return d3.ToFloat64()
}

func renderCtnSizeABTest(r *mvutil.RequestParams, res *corsair_proto.QueryResult_) {
	if !mvutil.IsHbOrV3OrV5Request(r.Param.RequestPath) {
		return
	}
	if r.Param.AdType != mvconst.ADTypeRewardVideo && r.Param.AdType != mvconst.ADTypeInterstitialVideo {
		return
	}
	ctnSizeABTestConf, _ := extractor.GetCTN_SIZE_ABTEST()
	if len(ctnSizeABTestConf) == 0 {
		return
	}
	for _, v := range res.AdsByBackend {
		if dspExt, _ := r.GetDspExt(); v.BackendId == mvconst.MAdx && !r.IsFakeAs && !(dspExt != nil && dspExt.DspId == mvconst.MAS) {
			// 只针对online api及adx做切量,只选择第一个backend, 不走as和mas
			dspId := strconv.FormatInt(dspExt.DspId, 10)

			if dspUnitConf, ok := ctnSizeABTestConf[dspId]; ok {
				if dspRate, ok := dspUnitConf[v.RequestKey]; ok {
					dspRateRand := rand.Intn(100)
					if dspRate > int32(dspRateRand) {
						r.Param.CtnSizeTag = "2"
					}
				}
			}
			return
		}
	}
}

func renderBannerHtml(bannerHtml string, dspId int64) string {
	// 对于pioneer返回的html模版，不需要另外加html标签
	if dspId == mvconst.MAS || extractor.IsVastBannerDsp(dspId) {
		return bannerHtml
	}
	return addHtmlCodeOnMraid(bannerHtml)
}

func renderMraidByDsp(mraid string, dspId int64) string {
	adnConfList := extractor.GetADNET_CONF_LIST()
	if needHtmlCodeDspList, ok := adnConfList["needHtmlCodeDspList"]; ok && mvutil.InInt64Arr(dspId, needHtmlCodeDspList) {
		return addHtmlCodeOnMraid(mraid)
	}
	return mraid
}

func addHtmlCodeOnMraid(mraid string) string {
	htmlStrConf, _ := extractor.GetBANNER_HTML_STR()
	htmlStr := "<!DOCTYPE html><html lang=\"en\"><head><meta charset=\"UTF-8\"><meta name=\"viewport\"content=\"width=device-width, initial-scale=1.0\"><meta http-equiv=\"X-UA-Compatible\"content=\"ie=edge\"><title>MTG</title></head><body><div id=\"MTGWRAPPER\">{mv_adn_replace}</div></body></html>"
	if len(htmlStrConf) > 0 {
		htmlStr = htmlStrConf
	}
	htmlStr = strings.Replace(htmlStr, "{mv_adn_replace}", mraid, -1)
	return htmlStr
}

func RenderExtData(r *mvutil.RequestParams) string {
	var extData mvutil.ExtData
	if r.Param.NewMoreOfferFlag {
		extData.MofUnitId = r.UnitInfo.Unit.MofUnitId
	}
	// 判断主offer 是否为online api返回的offer
	if r.Param.Mof == 1 {
		// 即无ifmd5也无vfmd5则判定主offer为三方广告主返回的单子
		extData.IsThirdParty = r.Param.IsThirdPartyMoreoffer
		// 记录主offer的request_id
		extData.CrtRid = r.Param.ExtDataInit.CrtRid
		// 请求日志中区分more offer请求还是appwall 请求
		if r.Param.ParentUnitId > 0 {
			extData.IsMoreOffer = 1
		}
	}
	if r.Param.HasWX == true {
		extData.HasWX = r.Param.HasWX
	}
	if r.Param.ReqTypeAABTest > 0 {
		extData.ReqTypeTest = r.Param.ReqTypeAABTest
	}

	if r.Param.DisplayCampaignABTest > 0 {
		extData.DisplayCampaignABTest = r.Param.DisplayCampaignABTest
	}

	if r.Param.CleanDeviceTest > 0 {
		extData.CleanDeviceTest = r.Param.CleanDeviceTest
	}

	if r.Param.VcnABTest > 0 {
		extData.VcnABTest = r.Param.VcnABTest
	}
	if len(r.Param.CtnSizeTag) > 0 {
		extData.CtnSizeTest = r.Param.CtnSizeTag
	}
	// 记录parent_unit id,UcParentUnitId为固定unit，ParentUnitId为拆分unit
	if r.Param.UcParentUnitId > 0 {
		extData.ParentUnitId = r.Param.UcParentUnitId
		// 固定unit的parentunit也记录下来，方便算法做数据统计
		r.Param.ExtDataInit.ParentUnitId = r.Param.UcParentUnitId
	}
	if r.Param.ParentUnitId > 0 {
		extData.ParentUnitId = r.Param.ParentUnitId
	}
	// 记录mof_type以区分more offer还是close button ad
	if r.Param.MofType > 0 {
		extData.MofType = r.Param.MofType
		r.Param.ExtDataInit.MofType = r.Param.MofType
	}
	// 请求时记录clsad标志
	if len(r.Param.ExtDataInit.CloseAdTag) > 0 {
		extData.CloseAdTag = r.Param.ExtDataInit.CloseAdTag
	}
	// 记录h5_type
	if r.Param.H5Type > 0 {
		extData.H5Type = r.Param.H5Type
	}

	// 记录H5上报的标记数据
	h5Data := checkH5Data(r.Param.H5Data, &r.Param)
	extData.H5Data = h5Data
	r.Param.ExtDataInit.H5Data = h5Data

	// 记录中国专线切量标记
	extData.MvLine = r.Param.MvLine
	// tracking日志也需记录
	r.Param.ExtDataInit.MvLine = r.Param.MvLine

	// 记录tracking中国专线域名
	extData.CNTrackDomain = r.Param.ExtDataInit.CNTrackDomain

	// tk 域名abtest
	extData.TKSysTag = r.Param.ExtDataInit.TKSysTag

	// 请求日志记录dmp切量标记
	extData.GaidTag = r.Param.ExtDataInit.GaidTag
	extData.IdfaTag = r.Param.ExtDataInit.IdfaTag
	extData.AndroidIdTag = r.Param.ExtDataInit.AndroidIdTag
	extData.ImeiTag = r.Param.ExtDataInit.ImeiTag
	extData.ImeiMd5Tag = r.Param.ExtDataInit.ImeiMd5Tag

	extData.ExcludePkg = r.Param.ExtDataInit.ExcludePkg
	extData.ExcludePsbPkg = r.Param.ExtDataInit.ExcludePsbPkg
	extData.ImpExcludePkg = r.Param.ExtDataInit.ImpExcludePkg
	// reduce fill mode flag
	if r.Param.ReduceFillConfig != nil {
		extData.ReduceFillMode = r.Param.ReduceFillConfig.ControlMode
	}
	// 记录dmt，dmf，cpu type
	extData.Dmt = r.Param.Dmt
	extData.Dmf = r.Param.Dmf
	extData.CpuType = r.Param.Ct
	// 根据开关决定是否需要记录到tracking日志中
	adnConf, _ := extractor.GetADNET_SWITCHS()
	if dmTrackSwitch, ok := adnConf["dmTrackSwitch"]; ok && dmTrackSwitch == 1 {
		r.Param.ExtDataInit.Dmt = r.Param.Dmt
		r.Param.ExtDataInit.Dmf = r.Param.Dmf
		r.Param.ExtDataInit.CpuType = r.Param.Ct
	}
	// 记录adnet和adx的灰度标记
	renderStartMode(&extData, r)
	// 频次控制
	extData.PriceFactor = r.Param.ExtDataInit.PriceFactor
	extData.PriceFactorGroupName = r.Param.ExtDataInit.PriceFactorGroupName
	extData.PriceFactorTag = r.Param.ExtDataInit.PriceFactorTag
	extData.PriceFactorFreq = r.Param.ExtDataInit.PriceFactorFreq
	extData.PriceFactorHit = r.Param.ExtDataInit.PriceFactorHit
	// placement维度的频控
	extData.ImpressionCap = r.Param.ExtDataInit.ImpressionCap
	extData.ImpressionCapTime = r.Param.ExtDataInit.ImpressionCapTime

	// Aerospike 是否使用Gzip压缩
	// extData.AerospikeGzipEnable = r.Param.ExtDataInit.AerospikeGzipEnable
	// Aerospike 是否移除冗余
	// extData.AerospikeRemoveRedundancyEnable = r.Param.ExtDataInit.AerospikeRemoveRedundancyEnable
	// 当前Go版本
	// extData.GoVersion = runtime.Version()
	// 是否请求bid server
	// extData.RequestBidServer = r.Param.ExtDataInit.RequestBidServer
	// 三方广告源视频，endcard模版abest标记
	extData.VideoDspTplABTest = r.Param.VideoDspTplABTestId
	extData.EndcardDspTplABTest = r.Param.EndcardDspTplABTestId

	//接入treasure box的标记位
	extData.TreasureBoxAbtestTag = r.Param.ExtDataInit.TreasureBoxAbtestTag

	// 大模板召回开关
	extData.RwPlus = r.Param.RwPlus
	// V5 接口
	extData.V5AbtestTag = r.Param.ExtDataInit.V5AbtestTag

	extData.ReplacedImpTrackDomainId = r.Param.ExtDataInit.ReplacedImpTrackDomainId
	extData.ReplacedClickTrackDomainId = r.Param.ExtDataInit.ReplacedClickTrackDomainId

	// 带宽
	extData.BandWidth = r.Param.BandWidth

	extData.AdxAbTest = r.Param.ExtDataInit.AdxAbTest

	extData.SqsCollect = r.Param.ExtDataInit.SqsCollect
	if r.Param.BigTemplateFlag {
		extData.BigTemplateTag = 1
	} else {
		extData.BigTemplateTag = 2
	}

	// 记录频次控制传递给as的包名
	if len(r.Param.ExcludePackageNames) > 0 {
		extData.FreExcludePkgList = renderFreExcludePkgList(r.Param.ExcludePackageNames)
	}
	extData.DeeplinkType = r.Param.DeeplinkType
	extData.HtmlSupport = r.Param.HtmlSupport

	extData.ToponRequestId = r.Param.ExtDataInit.ToponRequestId
	extData.OsvUpTime = r.Param.OsvUpTime
	extData.UpTime = r.Param.UpTime
	extData.NewId = r.Param.NewId
	extData.OldId = r.Param.OldId
	extData.ImeiABTest = r.Param.ExtDataInit.ImeiABTest

	r.Param.ExtDataInit.IsReturnWtickTag = r.Param.IsReturnWtick
	extData.IsReturnWtickTag = r.Param.IsReturnWtick

	// 记录ios 设备信息。
	r.Param.ExtDataInit.UpTime = r.Param.UpTime
	r.Param.ExtDataInit.Brt = r.Param.Brt
	r.Param.ExtDataInit.Vol = r.Param.Vol
	r.Param.ExtDataInit.Lpm = r.Param.Lpm
	r.Param.ExtDataInit.Font = r.Param.Font

	extData.LimitTrk = r.Param.LimitTrk
	extData.Att = r.Param.Att
	extData.Brt = r.Param.Brt
	extData.Vol = r.Param.Vol
	extData.Lpm = r.Param.Lpm
	extData.Font = r.Param.Font

	extData.FwType = r.Param.FwType
	extData.HardwareModel = r.Param.HardwareModel
	extData.Ntbarpt = r.Param.Ntbarpt
	r.Param.ExtDataInit.Ntbarpt = r.Param.Ntbarpt
	extData.Ntbarpasbl = r.Param.Ntbarpasbl
	r.Param.ExtDataInit.Ntbarpasbl = r.Param.Ntbarpasbl
	extData.AtatType = r.Param.AtatType
	r.Param.ExtDataInit.AtatType = r.Param.AtatType

	extData.MappingIdfaTag = r.Param.ExtDataInit.MappingIdfaTag
	extData.ExcludeAopPkg = r.Param.ExtDataInit.ExcludeAopPkg
	extData.GdprConsent = r.Param.GdprConsent

	extData.TKCNABTestTag = r.Param.ExtDataInit.TKCNABTestTag
	extData.TKCNABTestAATag = r.Param.ExtDataInit.TKCNABTestAATag
	extData.MappingIdfaCoverIdfaTag = r.Param.ExtDataInit.MappingIdfaCoverIdfaTag

	extData.AppSettingId = r.Param.AppSettingId
	extData.UnitSettingId = r.Param.UnitSettingId
	extData.RewardSettingId = r.Param.RewardSettingId
	extData.MiskSpt = r.Param.MiSkSpt
	extData.MoreofferAndAppwallMvToPioneerTag = r.Param.ExtDataInit.MoreofferAndAppwallMvToPioneerTag
	extData.PioneerHttpCode = r.Param.ExtDataInit.PioneerHttpCode
	extData.OnlineApiNeedOfferBidPrice = r.Param.OnlineApiNeedOfferBidPrice
	extData.TrackDomainByCountryCode = r.Param.ExtDataInit.TrackDomainByCountryCode
	extData.DecryptHarmonyInfo = r.Param.DecryptHarmonyInfo
	extData.LoadCDNTag = r.Param.ExtDataInit.LoadCDNTag
	extData.TmaxABTestTag = r.Param.ExtDataInit.TmaxABTestTag
	extData.UseDynamicTmax = r.Param.ExtDataInit.UseDynamicTmax
	extData.ReqTimeout = r.Param.ExtDataInit.ReqTimeout
	extData.MasResponseTime = r.Param.MasResponseTime
	extData.Dnt = r.Param.Dnt
	extData.MultiVcn = r.Param.ExtDataInit.MultiVcn
	extData.VcnCampaigns = r.Param.ExtDataInit.VcnCampaigns
	extData.ExpIds = r.Param.ExtDataInit.ExpIds
	extData.MpToPioneerTag = r.Param.ExtDataInit.MpToPioneerTag
	extData.GetAdsErrbackendList = r.Param.GetAdsErrbackendList
	extData.GetAdsErr = r.Param.GetAdsErr
	extData.DeviceGEOCCMatch = r.Param.ExtDataInit.DeviceGEOCCMatch
	extData.ThreeLetterCountry = r.Param.ExtDataInit.ThreeLetterCountry
	extData.HBSubsidyType = r.Param.ExtDataInit.HBSubsidyType

	// 渲染 bid server tag
	renderBidServerTag(&extData, r)
	// 渲染 rs 精准出价
	extData.BidServerRsPrice = r.Param.ExtDataInit.BidServerRsPrice
	// hb load 过滤
	extData.LoadRejectCode = r.LoadRejectCode
	extData.CdnTrackingDomainABTestTag = r.Param.ExtDataInit.CdnTrackingDomainABTestTag

	if len(r.Param.DecryptTrafficInfoStr) > 0 {
		// 抽样记录下来
		randVal := rand.Intn(10000)
		if recordTrafficInfoRate, ok := adnConf["recordTrafficInfoRate"]; ok && recordTrafficInfoRate > randVal {
			extData.DecryptTrafficInfoStr = r.Param.DecryptTrafficInfoStr
		}
	}

	extDataJson, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(extData)
	extDataStr := string(extDataJson)
	return extDataStr
}

func renderBidServerTag(extData *mvutil.ExtData, r *mvutil.RequestParams) {
	if r.Param.IsHitRequestBidServer == 1 {
		// 说明命中了切量
		if r.Param.BidServerCtx != nil {
			// 该字段不为空说明为 bid server 出价的情况且胜出
			extData.BidServerTag = 1
			return
		} else if dspExt, err := r.GetDspExt(); err == nil && dspExt.DspId == mvconst.MAS {
			// 如果有广告且dspId是Mas说明是 bid server 回退的情况, 且胜出了
			extData.BidServerTag = 2
			return
		} else {
			// 这种情况说明命中了 bid server, 但是胜出的不是bid server
			extData.BidServerTag = 3
			return
		}
	}
	return
}

func renderFreExcludePkgList(pkgList map[string]bool) []string {
	var excludePkgList []string
	for pkg, _ := range pkgList {
		excludePkgList = append(excludePkgList, pkg)
	}
	return excludePkgList
}

func renderStartMode(extData *mvutil.ExtData, r *mvutil.RequestParams) {
	extData.AdnetStartMode = r.Param.StartMode
	if r.Param.StartModeTags == nil {
		r.Param.StartModeTags = make(map[string]string)
	}
	r.Param.StartModeTags[mvconst.AdnetStartModeTag] = r.Param.StartMode
	if mode, ok := r.Param.AdxStartMode[mvconst.AdxStartModeTag]; ok {
		extData.AdxStartMode = mode
		r.Param.StartModeTags[mvconst.AdxStartModeTag] = mode
	}
	// 封装带到tracking的start mode 标记
	startModeTags, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(r.Param.StartModeTags)
	r.Param.StartModeTagsStr = string(startModeTags)
}

func checkH5Data(h5Data string, params *mvutil.Params) string {
	adnetConf, _ := extractor.GetADNET_SWITCHS()
	size := 1000
	if h5DataSize, ok := adnetConf["h5DataSize"]; ok {
		size = h5DataSize
	}
	//  过长则不记录
	if len(h5Data) > size {
		mvutil.Logger.Runtime.Warnf("request_id=[%s],h5Data=[%s] h5Data out of size", params.RequestID, h5Data)
		return ""
	}
	// 对传入内容进行清洗
	return mvutil.CleanH5Data(h5Data)
}

func renderBigTempalteABTestTag(params *mvutil.Params, dspId int64) {
	// 真正下发大模板url才可以做abtest
	if !mvutil.IsBigTemplate(params.BigTemplateId) {
		return
	}
	abtestFieldsConf, _ := extractor.GetABTEST_FIELDS()
	abtestConfs, _ := extractor.GetABTEST_CONFS()
	if len(abtestFieldsConf.BigTemplateParams) == 0 {
		return
	}
	if params.ABTestTags == nil {
		params.ABTestTags = make(map[string]int)
	}
	for _, v := range abtestFieldsConf.BigTemplateParams {
		if abtestConf, ok := abtestConfs[v]; ok {
			conf, randType := GetABTestConf(params, nil, abtestConf, dspId)
			if len(conf) == 0 {
				continue
			}
			finalVal, randOk := GetABTestRes(conf, randType, params, v)
			if !randOk {
				continue
			}
			params.ABTestTags[v] = finalVal
			// 拼装大模板abtest下发参数
			params.BigTplABTestParamsStr += "&" + v + "=" + strconv.Itoa(finalVal)
		}
	}
}

// func renderCrLangABTest(campaignId int64, r *mvutil.RequestParams, crMap, crTypeIdMap *map[ad_server.CreativeType]int64) {
// 	langCreativeABTestConf := extractor.GetAdnetLangCreativeABTestConf()
// 	adnetConfList := extractor.GetADNET_CONF_LIST()
// 	if len(langCreativeABTestConf) == 0 {
// 		return
// 	}
// 	abTestCreativeTypeList, ok := adnetConfList["abTestCreativeTypeList"]
// 	if !ok || len(abTestCreativeTypeList) == 0 {
// 		return
// 	}
// 	if camCrInfo, ok := langCreativeABTestConf[strconv.FormatInt(campaignId, 10)+"_"+r.Param.CountryCode]; ok && len(camCrInfo) > 0 {
// 		// 切量标记
// 		r.Param.CrLangABTestTag = true
// 		// 随机获取素材组
// 		randVal := rand.Intn(len(camCrInfo))
// 		// 获取素材组
// 		replaceCreatives := camCrInfo[randVal]
// 		// 替换素材组
// 		for _, creativeType := range abTestCreativeTypeList {
// 			(*crMap)[ad_server.CreativeType(creativeType)] = replaceCreatives[strconv.FormatInt(creativeType, 10)]
// 			(*crTypeIdMap)[ad_server.CreativeType(creativeType)] = replaceCreatives[strconv.FormatInt(creativeType, 10)]
// 		}
// 	}
// 	return
// }

func renderBannerTemplate(mr *MobvistaResult, ads *corsair_proto.BackendAds, r *mvutil.RequestParams, dspId int64) {
	if ads.BannerUrl != nil {
		if r.Param.AdType == mvconst.ADTypeSplash {
			mr.Data.AdTplUrl = *ads.BannerUrl
		} else {
			mr.Data.BannerUrl = *ads.BannerUrl
		}
	}
	if ads.BannerHtml != nil {
		if r.Param.AdType == mvconst.ADTypeSplash {
			mr.Data.AdHtml = renderBannerHtml(*ads.BannerHtml, dspId)
		} else {
			mr.Data.BannerHtml = renderBannerHtml(*ads.BannerHtml, dspId)
		}
	}
}

func replaceTrackingDomain(r *mvutil.RequestParams) {
	// 排出掉mp流量
	if mvutil.IsMP(r.Param.RequestPath) {
		return
	}
	replaceTrackingDomainConfs := extractor.GetREPLACE_TRACKING_DOMAIN_CONF()
	if replaceTrackingDomainConfs == nil {
		return
	}
	region := mvutil.Config.AreaConfig.HttpConfig.RegionName
	// 优先级：unit>app>pub>total
	if replaceTrackingDomainConfs.UnitConfs != nil {
		if regionConf, ok := replaceTrackingDomainConfs.UnitConfs[region]; ok && regionConf != nil {
			if unitConf, ok := regionConf[strconv.FormatInt(r.Param.UnitID, 10)]; ok && unitConf != nil {
				chooseTrackingDomainByAction(r, unitConf)
				return
			}
		}
	}

	if replaceTrackingDomainConfs.AppConfs != nil {
		if regionConf, ok := replaceTrackingDomainConfs.AppConfs[region]; ok && regionConf != nil {
			if appConf, ok := regionConf[strconv.FormatInt(r.Param.AppID, 10)]; ok && appConf != nil {
				chooseTrackingDomainByAction(r, appConf)
				return
			}
		}
	}

	if replaceTrackingDomainConfs.PubConfs != nil {
		if regionConf, ok := replaceTrackingDomainConfs.PubConfs[region]; ok && regionConf != nil {
			if pubConf, ok := regionConf[strconv.FormatInt(r.Param.PublisherID, 10)]; ok && pubConf != nil {
				chooseTrackingDomainByAction(r, pubConf)
				return
			}
		}
	}

	// 整体维度
	if regionConf, ok := replaceTrackingDomainConfs.TotalConf[region]; ok && regionConf != nil {
		chooseTrackingDomainByAction(r, regionConf)
	}
}

func chooseTrackingDomainByAction(r *mvutil.RequestParams, actionConfs *mvutil.TrackingDomainActionConf) {
	// 替换展示tracking domain
	if actionConfs.ImpressionTrackConf != nil {
		r.Param.ExtDataInit.ReplacedImpTrackDomainId, r.Param.ReplacedImpTrackDomain = chooseTrackingDomain(actionConfs.ImpressionTrackConf)
	}
	// 替换点击tracking domain
	if actionConfs.ClickTrackConf != nil {
		r.Param.ExtDataInit.ReplacedClickTrackDomainId, r.Param.ReplacedClickTrackDomain = chooseTrackingDomain(actionConfs.ClickTrackConf)
	}
}

func chooseTrackingDomain(Confs []*mvutil.TrackingDomainWeightConf) (int, string) {
	trackDomainWeightMap := make(map[string]int)
	trackDomainIdMap := make(map[string]int)
	for _, conf := range Confs {
		trackDomainWeightMap[conf.Domain] = conf.Weight
		trackDomainIdMap[conf.Domain] = conf.Id
	}
	var newTrackingDomain string
	var trackId int
	newTrackingDomain = mvutil.RandByRate2(trackDomainWeightMap)
	if id, ok := trackDomainIdMap[newTrackingDomain]; ok {
		trackId = id
	}
	return trackId, newTrackingDomain
}

func renderNewEncryptedId(r *mvutil.RequestParams) {
	renderNewEncryptedSysId(r)

	RenderNewEncryptedBkupId(r)
}

func RenderNewEncryptedBkupId(r *mvutil.RequestParams) {
	// 新版本才处理
	if r.Param.ApiVersion < mvconst.API_VERSION_2_2 {
		return
	}
	if !mvutil.IsHbOrV3OrV5Request(r.Param.RequestPath) {
		return
	}

	// 如果此前已经生成了NewEncryptedBkupId，则不再生成。
	if len(r.Param.NewEncryptedBkupId) > 0 {
		return
	}

	// 新session或超过设置时间范围，才做重新加密bkupid的逻辑
	adnConf, _ := extractor.GetADNET_SWITCHS()
	bkupIdUpInterval, ok := adnConf["bkupIdUpInterval"]
	if (!ok || r.Param.RequestTime-r.Param.EncryptedBkupIdTimestamp <= int64(bkupIdUpInterval)) && !r.Param.IsNewSession {
		return
	}
	// 当需要替换的时候，请求mapping server生成加密后的ruid
	if getBkupIdFromMpServerSwitch, ok := adnConf["getBkupIdFromMpServerSwitch"]; ok && getBkupIdFromMpServerSwitch == 1 {
		if r.Param.PlatformName == mvconst.PlatformNameIOS {
			RenderRuid(r, "encryption")
			if len(r.Param.EncryptedRuid) > 0 {
				r.Param.BkupId = r.Param.EncryptedRuid
			}
		} else {
			// 对于安卓新版本，目前还是使用setting 旧逻辑生成bkupid，等mapping server支持安卓的时候，再请求mapping server
			// 对于已有bkupid，则不需要重新生成了。
			if len(r.Param.BkupId) == 0 {
				v4, _ := uuid.NewV4()
				r.Param.BkupId = v4.String()
			}
		}
	}
	r.Param.NewEncryptedBkupId = NewEncryptDevId(r, r.Param.BkupId)
}

func renderNewEncryptedSysId(r *mvutil.RequestParams) {
	// 新版本才处理
	if r.Param.ApiVersion < mvconst.API_VERSION_2_2 {
		return
	}
	if !mvutil.IsHbOrV3OrV5Request(r.Param.RequestPath) {
		return
	}
	if len(r.Param.SysId) == 0 {
		return
	}
	// 新session则重新生成
	if r.Param.IsNewSession {
		r.Param.NewEncryptedSysId = NewEncryptDevId(r, r.Param.SysId)
		return
	}
	// sysid
	adnConf, _ := extractor.GetADNET_SWITCHS()
	sysIdUpInterval, ok := adnConf["sysIdUpInterval"]
	if !ok {
		return
	}
	if r.Param.RequestTime-r.Param.EncryptedSysIdTimestamp <= int64(sysIdUpInterval) {
		return
	}
	r.Param.NewEncryptedSysId = NewEncryptDevId(r, r.Param.SysId)
}

func NewEncryptDevId(r *mvutil.RequestParams, devId string) string {
	newEncryptedDevId, err := mvutil.DevIdCbcEncrypt(strconv.FormatInt(r.Param.RequestTime, 10) + "_" + devId)
	if err != nil {
		mvutil.Logger.Runtime.Errorf("str=[%s] encrypt id error. error:%s", devId, err.Error())
		return ""
	}
	return newEncryptedDevId
}

func RenderHBRequestLogParam(r *mvutil.RequestParams) {
	r.Param.Extra = getCloudExtra()
	if r.IsHBRequest {
		r.Param.Extra = r.Param.Algorithm
	}
	// 处理sdk传来的has_wx（是否安装微信）
	r.Param.ExtData = RenderExtData(r)
	// 结果记录到extdev里
	r.Param.ExtDeviceId = RenderExtDeviceId(r)
	if r.AsResp != nil {
		r.Param.AsABTestResTag = r.AsResp.GetAsAbtestResTag()
	}
	r.Param.ABTestTagStr = r.AsResp.GetAbtestResTag()
	r.Param.JunoCommonLogInfoJson = r.AsResp.GetJunoLogInfo()
	r.Param.PioneerExtdataInfo = r.AsResp.GetPioneerExtdataInfo()
	r.Param.PioneerOfferExtdataInfo = r.AsResp.GetPioneerOfferExtdataInfo()
}

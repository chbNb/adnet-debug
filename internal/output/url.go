package output

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
	"gitlab.mobvista.com/ADN/adx_common/model"

	"gitlab.mobvista.com/ADN/adnet/internal/corsair_proto"
	"gitlab.mobvista.com/ADN/adnet/internal/extractor"
	"gitlab.mobvista.com/ADN/adnet/internal/mvconst"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	"gitlab.mobvista.com/ADN/chasm/module/demand"
	smodel "gitlab.mobvista.com/ADN/structs/model"
)

func RenderImpUrl(queryURL string, params *mvutil.Params, r *mvutil.RequestParams) string {
	// gaid
	queryURL = strings.Replace(queryURL, "{gaid}", params.GAID, -1)
	// idfa
	queryURL = strings.Replace(queryURL, "{idfa}", params.IDFA, -1)
	// devid
	devId := getDevId(r, params)
	queryURL = strings.Replace(queryURL, "{devId}", devId, -1)
	// imei
	queryURL = strings.Replace(queryURL, "{imei}", params.IMEI, -1)
	// oaid
	queryURL = strings.Replace(queryURL, "{oaid}", url.QueryEscape(params.OAID), -1)
	// mac
	queryURL = strings.Replace(queryURL, "{mac}", params.MAC, -1)
	idfa := params.IDFA
	var idfaMd5 string
	if len(params.IDFAMd5) == 0 {
		idfaMd5 = getMd5(idfa)
	} else {
		idfaMd5 = strings.ToUpper(params.IDFAMd5)
	}
	idfaSha1 := getSha1(idfa)

	gaid := params.GAID
	var gaidMd5 string
	if len(params.GAIDMd5) == 0 {
		gaidMd5 = getMd5(gaid)
	} else {
		gaidMd5 = strings.ToUpper(params.GAIDMd5)
	}
	gaidSha1 := getSha1(gaid)

	devID := params.AndroidID
	var devIdMd5 string
	if len(params.AndroidIDMd5) == 0 {
		devIdMd5 = getMd5(devID)
	} else {
		devIdMd5 = strings.ToUpper(params.AndroidIDMd5)
	}

	devIdSha1 := getSha1(devID)

	imei := params.IMEI
	var imeiMd5 string
	if len(params.ImeiMd5) == 0 {
		imeiMd5 = getMd5Lower(imei)
	} else {
		imeiMd5 = strings.ToLower(params.ImeiMd5)
	}
	imeiSha1 := getSha1(imei)

	mac := params.MAC
	macMd5 := getMd5(mac)
	macSha1 := getSha1(mac)

	queryURL = strings.Replace(queryURL, "{gaidMd5}", gaidMd5, -1)
	queryURL = strings.Replace(queryURL, "{gaidSha1}", gaidSha1, -1)
	queryURL = strings.Replace(queryURL, "{imeiMd5}", imeiMd5, -1)
	queryURL = strings.Replace(queryURL, "{imeiSha1}", imeiSha1, -1)
	queryURL = strings.Replace(queryURL, "{idfaMd5}", idfaMd5, -1)
	queryURL = strings.Replace(queryURL, "{idfaSha1}", idfaSha1, -1)
	queryURL = strings.Replace(queryURL, "{macMd5}", macMd5, -1)
	queryURL = strings.Replace(queryURL, "{macSha1}", macSha1, -1)
	queryURL = strings.Replace(queryURL, "{devIdMd5}", devIdMd5, -1)
	queryURL = strings.Replace(queryURL, "{devIdSha1}", devIdSha1, -1)

	// network
	queryURL = strings.Replace(queryURL, "{network}", params.NetworkTypeName, -1)

	// packageName
	queryURL = strings.Replace(queryURL, "{package_name}", params.ExtfinalPackageName, -1)
	// subId
	subId := getSubId(params)
	subIdStr := strconv.FormatInt(subId, 10)
	queryURL = strings.Replace(queryURL, "{subId}", subIdStr, -1)
	// mtgId
	mtgIdStr := "mtg" + strconv.FormatInt(params.ExtMtgId, 10)
	queryURL = strings.Replace(queryURL, "{mtgId}", mtgIdStr, -1)

	queryURL = renderEncodeURIComponent(queryURL)
	return queryURL
}

func getDevId(r *mvutil.RequestParams, params *mvutil.Params) string {
	devId := params.AndroidID
	devinfoEncrypt := smodel.IsDevinfoEncrypt(r.AppInfo)
	if devinfoEncrypt {
		devId = ""
	}
	return devId
}

func getSubId(params *mvutil.Params) int64 {
	if params.Extfinalsubid > int64(0) {
		return params.Extfinalsubid
	}
	if params.Extra14 > int64(0) {
		return params.Extra14
	}
	return params.AppID
}

func GetPlayableUrl(playableUrl string, protocal int, params *mvutil.Params) string {
	switch protocal {
	case 1:
		playableUrl = "http://" + playableUrl
	case 2:
		playableUrl = "https://" + playableUrl
	default:
		playableUrl = GetUrlScheme(params) + "://" + playableUrl
	}
	return playableUrl
}

func GetUrlScheme(params *mvutil.Params) string {
	if params.HTTPReq == int32(2) {
		return "https"
	}
	return "http"
}

func NeedSchemeHttps(httpReq int32) bool {
	return httpReq == int32(2)
}

func RenderEndcardUrl(params *mvutil.Params, endcardUrl string) string {
	if len(endcardUrl) <= 0 {
		return endcardUrl
	}
	if len(endcardUrl) >= 4 && endcardUrl[0:4] == "http" {
		return endcardUrl
	}
	scheme := GetUrlScheme(params)
	endcardUrl = scheme + "://" + endcardUrl
	return endcardUrl
}

// RenderNewUrls
// isMoreAds: 如果是true, 需要重置dspExt
func RenderNewUrls(
	r *mvutil.RequestParams, ad *Ad,
	params *mvutil.Params,
	campaign *smodel.CampaignInfo,
	campaignId int64, //campaign & campaignId 二选一
	key int,
	isMoreAds bool,
) {
	if len(r.Param.QueryP) == 0 {
		r.Param.QueryP = UrlReplace1(mvutil.SerializeP(params))
	}
	if ad.AKS == nil {
		ad.AKS = make(map[string]string)
	}

	dspExt, err := r.GetDspExt()
	if (err == nil && dspExt != nil &&
		(dspExt.DspId == mvconst.MAS ||
			//如果是onlineAPI的流量第二个offer需要进入逻辑
			((r.Param.AdType == mvconst.ADTypeNative || r.Param.RequestType == mvconst.REQUEST_TYPE_OPENAPI_AD) && key > 0)) ||
		mvutil.IsRequestPioneerDirectly(&r.Param)) &&
		r.AsResp != nil { // 走 MAS
		if urlParams := r.AsResp.GetUrlParam(); len(urlParams) > 0 {
			urlParam := urlParams[key]
			params.QueryQ = urlParam.GetQ()
			params.QueryQ2 = UrlDeReplace3(params.QueryQ) // 和下面 else 中的逻辑保持一致
			params.QueryR = urlParam.GetR()
			params.QueryAL = urlParam.GetAl()
			params.QueryZ = urlParam.GetZ()
			r.Param.QueryZ = params.QueryZ
			k := urlParam.GetK()

			// onlineAPi流量如果接口有返回则使用接口返回的
			// more_offer的api version为1.3,也不支持宏替换
			if r.Param.RequestType == mvconst.REQUEST_TYPE_OPENAPI_AD || (mvutil.IsRequestPioneerDirectly(&r.Param) && r.Param.ApiVersion < mvconst.API_VERSION_1_4) { //不需要aks
				if len(urlParam.GetP()) > 0 {
					params.QueryP = urlParam.GetP()
					r.Param.QueryP = params.QueryP
				}
				if len(k) > 0 {
					params.RequestID = k //使用mas返回的k
				}
			} else {
				ad.AKS["k"] = k
			}
		}
	} else {
		if campaign == nil && campaignId > 0 { //兼容模式
			var ifFind bool
			campaign, ifFind = extractor.GetCampaignInfo(campaignId)
			if !ifFind {
				return
			}
		}
		params.QueryP = r.Param.QueryP
		params.QueryQ2 = mvutil.SerializeQ(params, campaign)
		params.QueryQ = UrlReplace3(params.QueryQ2)
		params.QueryR = UrlReplace1(params.QueryR)
		params.QueryAL = url.QueryEscape(params.Extalgo)
	}
	if isMoreAds {
		// moreAds 时， 实际出的单子是 mas通过 VastMas返回的， 但此时的dspExt还是原来adx里返回时生成的，因此需要重新生成
		newDspExt := &model.DspExt{
			ChannelId: dspExt.ChannelId,
			DspId:     mvconst.MAS,
			PriceIn:   dspExt.PriceIn,
			PriceOut:  dspExt.PriceOut,
			IsHB:      dspExt.IsHB,
		}
		dspExt, _ := jsoniter.ConfigFastest.Marshal(newDspExt)
		params.QueryCSP = UrlReplace3(CreateCsp(params, string(dspExt)))
	} else {
		params.QueryCSP = UrlReplace3(CreateCsp(params, r.DspExt))
	}

	if !r.Param.IsNewUrl {
		return
	}
	if len(r.Param.QueryZ) == 0 {
		r.Param.QueryZ = mvutil.SerializeZ(params)
	}

	// mas 的在前面已经赋值
	if ad.AKS["k"] == "" {
		ad.AKS["k"] = params.RequestID
	}

	ad.AKS["q"] = params.QueryQ
	ad.AKS["r"] = params.QueryR
	ad.AKS["al"] = params.QueryAL
	ad.AKS["csp"] = params.QueryCSP
}

func RenderUrls(r *mvutil.RequestParams, ad *Ad, params *mvutil.Params, campaign *smodel.CampaignInfo) {
	// click notice url
	if params.PingMode == 1 {
		ad.ClickURL = GetJumpUrl(r, params, campaign, params.QueryQ2, ad)
		ad.NoticeURL = createClickUrl(params, true)
	} else {
		ad.ClickURL = createClickUrl(params, false)
	}

	if params.PingMode != 1 {
		// 排查joox问题，下发埋点 点击trackingurl 来分析
		testUnitConf, _ := extractor.GetCHET_URL_UNIT()
		if len(testUnitConf) > 0 {
			if mvutil.InInt64Arr(r.UnitInfo.UnitId, testUnitConf) {
				ad.ClickURL = ad.ClickURL + "&chet=1" + "&uid=" + strconv.FormatInt(params.UnitID, 10)
				ad.NoticeURL = createChetUrl(params)
			}
		}
		// 对于sdk的流量，必须得有notice_url，埋点上报的rid会从notice_url中的k参数获取
		// 加上useless_notice=1，单独记录一份日志
		if r.Param.RequestType == mvconst.REQUEST_TYPE_OPENAPI_V3 {
			ad.NoticeURL = createClickUrl(params, true) + mvconst.USELESS_NOTICE
		}

		// mp pingmode 为0的情况下，不需要返回notice url,节省流量成本
		// 使用r.Param.PingMode,按照流量的pingmode来确定是否清除noticeurl
		if mvutil.IsMpPingmode0(r) {
			ad.NoticeURL = ""
		}
	}
	// notice url 下发appsflyer token宏
	if params.ExtDataInit.ReturnAfTokenParam == 1 {
		ad.NoticeURL += mvconst.APPSFLYER_TOKEN_PARAMS
	}

	// impression url
	ad.ImpressionURL = createImpressionUrl(params)

	// 针对online api，给impression_url和click_url添加unit_id，以方便对数
	addUnitIdOnTrackingUrl(ad, params)

	// adtracking
	renderAdTracking(r, ad, params)

	// 对于clickmode为6的非gp，appstore单子。三方点击url放到ad_tracking.click中
	renderClickmode6NotInGpAndAppstore(params, ad, campaign, r)

	// 对于切给af 白名单请求的情况，需把click_url改为market地址
	renderThirdpartyWhiteListClickUrl(ad, params, campaign, r)
	// 针对三星应用商店在印度的代理商开发者做兼容，对于此开发者，click_url不做302，由服务端ping三方
	renderOnlinePubClickUrl(ad, params)
	// 淘宝单子因为SDK打开deeplink失败时才会调用clickUrl，所以需要把归因的tracking url需放到adtracking.click中由SDK点击，否则不会归因
	// 删除掉原来taobaooffer的逻辑
	renderDeepLinkOffer(ad, params)

	// 19-07-05 因appflyer ip+ua归因较低，后发现点击跳转时，因下发的c_ua为2，导致会调用header带的ua是无法做ipua归因的，且没有带language。
	// 而adtracking里的click上报是不受影响的。因此需要对于ios，903广告主，三方为appflyer，idfa为空值，clickmode为6的情况，使用click_url赋值到adtracking里的click内
	renderAdClickAndClickUrl(ad, params)

	if params.DemandContext != nil {
		ad.DeepLink = demand.RenderThirdPartyLink(params.DemandContext, extractor.DemandDao, ad.DeepLink)
	}
}

func isClickmode6NotInGpAndAppstore(params *mvutil.Params, ad *Ad) bool {
	return params.Extra10 == mvconst.JUMP_TYPE_CLIENT_DO_ALL && ad.CampaignType != mvconst.LINK_TYPE_APPSTORE && ad.CampaignType != mvconst.LINK_TYPE_GOOGLEPLAY &&
		params.AdType != mvconst.ADTypeAppwall
}

func renderClickmode6NotInGpAndAppstore(params *mvutil.Params, ad *Ad, campaign *smodel.CampaignInfo, r *mvutil.RequestParams) {
	if isClickmode6NotInGpAndAppstore(params, ad) {
		// 排除appwall不支持ad_tracking.click的量。
		ad.AdTracking.Click = append(ad.AdTracking.Click, ad.ClickURL)
		ad.AdTrackingPoint = &ad.AdTracking
		ctx := NewDemandContext(params, r)
		ad.ClickURL = demand.RenderThirdPartyLink(ctx, extractor.DemandDao, demand.GetCampaignNormalLandingPage(campaign))
	}
}

func renderDeepLinkOffer(ad *Ad, params *mvutil.Params) {
	if len(ad.DeepLink) > 0 {
		// deeplink + clickmode 13,不需要302到落地页
		if params.Extra10 == mvconst.JUMP_TYPE_ONLINE_DSP_AJUMP_SERVER {
			ad.ClickURL += mvconst.REDIRECT
		}
		if params.ExtDataInit.ClickMode6NotInGpAndAppstore == 1 {
			ad.ClickURL = ad.NoticeURL + mvconst.FORWARD // 只是记录点击
		} else {
			// 对于online api,没有noticeurl，因此需要用click_url代替
			if params.RequestType == mvconst.REQUEST_TYPE_OPENAPI_AD {
				// online api兜底链接需要删掉redirect，不然无法实现兜底逻辑
				ad.PreviewUrl = mvutil.DelSubStr(ad.ClickURL, mvconst.REDIRECT) + mvconst.FORWARD // 只是记录点击
			} else {
				ad.PreviewUrl = ad.NoticeURL + mvconst.FORWARD // 只是记录点击
			}

			params.AdClickReplaceclickByClickTag = true
		}
	}
}

func renderOnlinePubClickUrl(ad *Ad, params *mvutil.Params) {
	// 限制online api流量及clickmode=13才能做此逻辑
	if params.RequestType != mvconst.REQUEST_TYPE_OPENAPI_AD || params.Extra10 != mvconst.JUMP_TYPE_ONLINE_DSP_AJUMP_SERVER {
		return
	}
	adnetListConf := extractor.GetADNET_CONF_LIST()
	if noDirectPubs, ok := adnetListConf["noDirectPubs"]; ok && mvutil.InInt64Arr(params.PublisherID, noDirectPubs) {
		ad.ClickURL = ad.ClickURL + mvconst.REDIRECT
		// 改为notice模式
		if !strings.Contains(ad.ClickURL, "&notice=1") {
			ad.ClickURL += "&notice=1"
		}
	}
}

func addUnitIdOnTrackingUrl(ad *Ad, params *mvutil.Params) {
	if params.RequestType != mvconst.REQUEST_TYPE_OPENAPI_AD {
		return
	}
	units, _ := extractor.GetADD_UNIT_ID_ON_TRACKING_URL_UNIT()
	if len(units) > 0 && mvutil.InInt64Arr(params.UnitID, units) {
		unitIdStr := strconv.FormatInt(params.UnitID, 10)
		if len(ad.ClickURL) > 0 {
			ad.ClickURL += "&unit_id=" + unitIdStr
		}
		if len(ad.ImpressionURL) > 0 {
			ad.ImpressionURL += "&unit_id=" + unitIdStr
		}
	}
}

func renderAdTracking(r *mvutil.RequestParams, ad *Ad, params *mvutil.Params) {
	attype := ""
	if params.AdType == mvconst.ADTypeNative && len(params.VideoVersion) > 0 && len(ad.VideoURL) > 0 {
		attype = mvconst.ADTRACKING_TYPE_NATIVEVIDEO
	} else if params.AdType == mvconst.ADTypeRewardVideo && params.VersionFlag >= 1 {
		// TODO versionflag 解析
		attype = mvconst.ADTRACKING_TYPE_REWARDVIDEO
	} else if params.AdType == mvconst.ADTypeInterstitialVideo {
		attype = mvconst.ADTRACKING_TYPE_INTERSTITIALVIDEO
	} else if params.AdType == mvconst.ADTypeOnlineVideo {
		attype = mvconst.ADTRACKING_TYPE_ONLINEVIDEO
	} else if mvutil.IsJsVideo(r.Param.AdType) {
		if params.AdType == mvconst.ADTypeJSBannerVideo {
			attype = mvconst.ADTRACKING_TYPE_JSBANNERVIDEO
		} else if params.AdType == mvconst.ADTypeJSNativeVideo {
			attype = mvconst.ADTRACKING_TYPE_JSNATIVEVIDEO
		}
	} else if params.AdType == mvconst.ADTypeInteractive && r.Param.IAADRst == mvconst.IARST_PLAYABLE {
		attype = mvconst.ADTRACKING_TYPE_INTERACTIVE
	}
	if attype == "" && !isAdsense(params) && !demand.IsXiaoMiDeepLink(ad.DeepLink) {
		return
	}

	subUrl := renderSubUrl(params, attype)

	if attype == mvconst.ADTRACKING_TYPE_REWARDVIDEO || attype == mvconst.ADTRACKING_TYPE_INTERSTITIALVIDEO || attype == mvconst.ADTRACKING_TYPE_ONLINEVIDEO {
		renderRewardVideoAdTracking(ad, subUrl)
	} else if attype == mvconst.ADTRACKING_TYPE_NATIVEVIDEO || mvutil.IsJsVideo(r.Param.AdType) {
		renderNativeVideoAdTracking(ad, subUrl, params)
	} else if attype == mvconst.ADTRACKING_TYPE_INTERACTIVE {
		renderInteractiveAdTracking(ad, subUrl)
	}
	// 新pv统计逻辑,仅针对iv，rv下发
	if params.ApiVersion >= mvconst.API_VERSION_1_5 && params.ApiVersion < mvconst.API_VERSION_1_9 && (attype == mvconst.ADTRACKING_TYPE_REWARDVIDEO || attype == mvconst.ADTRACKING_TYPE_INTERSTITIALVIDEO) {
		ad.AdTracking.PubImp = getAdTrackingUrlList("pub_imp", subUrl, false, 0)
	}
	// adsense 下发exa_click
	if isAdsense(params) {
		subUrl = renderSubUrl(params, mvconst.ADTRACKING_TYPE_ADSENSE)
		ad.AdTracking.ExaClick = getAdTrackingUrlList("exa_click", subUrl, false, 0)
		ad.AdTracking.ExaImp = getAdTrackingUrlList("exa_imp", subUrl, false, 0)
	}

	// 只要是小米storekit的量，就需要下发apk下载安装进度埋点
	if len(ad.AdTracking.ApkDownloadStart) == 0 && len(ad.AdTracking.ApkDownloadEnd) == 0 && len(ad.AdTracking.ApkInstall) == 0 &&
		demand.IsXiaoMiDeepLink(ad.DeepLink) {
		renderApkUrls(ad, subUrl)
	}

	ad.AdTrackingPoint = &ad.AdTracking
}

func isAdsense(params *mvutil.Params) bool {
	return mvutil.IsMpad(params.RequestPath) && params.AdSense == 1
}

func renderSubUrl(params *mvutil.Params, attype string) string {
	subUrl := ""
	if params.IsNewUrl {
		subUrl = "{sh}://{do}/trackv2?z={z}&q={q}&type=" + attype + "&r={r}&c={c}&csp={csp}"
	} else {
		subUrl = GetUrlScheme(params) + "://" + params.Domain + "/trackv2?" + "p=" + params.QueryP + "&q=" + params.QueryQ + "&type=" + attype + "&r=" + params.QueryR + "&csp=" + params.QueryCSP
	}
	return subUrl
}

func renderRewardVideoAdTracking(ad *Ad, subUrl string) {
	ad.AdTracking.Mute = getAdTrackingUrlList("mute", subUrl, false, 0)
	ad.AdTracking.Unmute = getAdTrackingUrlList("unmute", subUrl, false, 0)
	ad.AdTracking.Endcard_show = getAdTrackingUrlList("endcard_show", subUrl, false, 0)
	ad.AdTracking.Close = getAdTrackingUrlList("close", subUrl, false, 0)
	ad.AdTracking.Pause = getAdTrackingUrlList("pause", subUrl, false, 0)
	renderApkUrls(ad, subUrl)
	renderPlayPercentage(ad, subUrl)
}

func renderApkUrls(ad *Ad, subUrl string) {
	if ad.CampaignType == 3 || demand.IsXiaoMiDeepLink(ad.DeepLink) {
		//apk
		ad.AdTracking.ApkDownloadStart = getAdTrackingUrlList("apk_download_start", subUrl, false, 0)
		ad.AdTracking.ApkDownloadEnd = getAdTrackingUrlList("apk_download_end", subUrl, false, 0)
		ad.AdTracking.ApkInstall = getAdTrackingUrlList("apk_install", subUrl, false, 0)
	}
}

func renderPlayPercentage(ad *Ad, subUrl string) {
	var ppList []CPlayTracking
	percentageList := []int{0, 25, 50, 75, 100}
	for _, percentage := range percentageList {
		var pp CPlayTracking
		pp.Rate = percentage
		pp.Url = getAdTrackingUrl("play_percentage", subUrl, true, percentage)
		ppList = append(ppList, pp)
	}
	if ad.AdTracking.Play_percentage == nil {
		ad.AdTracking.Play_percentage = ppList
		return
	}
	ad.AdTracking.Play_percentage = append(ad.AdTracking.Play_percentage, ppList...)
}

func renderNativeVideoAdTracking(ad *Ad, subUrl string, params *mvutil.Params) {
	if params.ApiVersion >= mvconst.API_VERSION_1_3 {
		renderPlayPercentage(ad, subUrl)
		if params.Extnvt2 > int32(0) && params.Extnvt2 != int32(1) {
			renderImpressionT2(ad, subUrl, params.Extnvt2)
		}
	} else if params.ApiVersion >= mvconst.API_VERSION_1_0 {
		renderPlayPercentage(ad, subUrl)
	} else {
		renderNVDefAdTracking(ad, subUrl)
	}

	if params.Extnvt2 == int32(3) {
		ad.AdTracking.Endcard_show = getAdTrackingUrlList("endcard_show", subUrl, false, 0)
	}
	if mvutil.InInt32Arr(params.AdType, []int32{mvconst.ADTypeJSBannerVideo, mvconst.ADTypeJSNativeVideo}) {
		ad.AdTracking.Mute = getAdTrackingUrlList("mute", subUrl, false, 0)
		ad.AdTracking.Unmute = getAdTrackingUrlList("unmute", subUrl, false, 0)
		ad.AdTracking.Close = getAdTrackingUrlList("close", subUrl, false, 0)
		ad.AdTracking.Pause = getAdTrackingUrlList("pause", subUrl, false, 0)
		ad.AdTracking.Impression = getAdTrackingUrlList("impression", subUrl, false, 0)
		ad.AdTracking.Click = getAdTrackingUrlList("click", subUrl, false, 0)
	}
	ad.AdTracking.Video_Click = getAdTrackingUrlList("video_click", subUrl, false, 0)
	renderApkUrls(ad, subUrl)
}

func renderNVDefAdTracking(ad *Ad, subUrl string) {
	ad.AdTracking.Start = getAdTrackingUrlList("play_percentage", subUrl, true, 0)
	ad.AdTracking.First_quartile = getAdTrackingUrlList("play_percentage", subUrl, true, 25)
	ad.AdTracking.Midpoint = getAdTrackingUrlList("play_percentage", subUrl, true, 50)
	ad.AdTracking.Third_quartile = getAdTrackingUrlList("play_percentage", subUrl, true, 75)
	ad.AdTracking.Complete = getAdTrackingUrlList("play_percentage", subUrl, true, 100)
}

func renderImpressionT2(ad *Ad, subUrl string, nvt2 int32) {
	var list []string
	nvt2Str := strconv.FormatInt(int64(nvt2), 10)
	aturl := subUrl + fmt.Sprintf("&key=%s", "impression_t2") + fmt.Sprintf("&nv_t2=%s", nvt2Str)
	list = append(list, aturl)
	ad.AdTracking.Impression_t2 = list
}

func getAdTrackingUrlList(key string, subUrl string, isPercentage bool, rate int) []string {
	var list []string
	aturl := getAdTrackingUrl(key, subUrl, isPercentage, rate)
	list = append(list, aturl)
	return list
}

func getAdTrackingUrl(key string, subUrl string, isPercentage bool, rate int) string {
	subUrl = subUrl + fmt.Sprintf("&key=%s", key)
	if isPercentage {
		subUrl = subUrl + fmt.Sprintf("&rate=%d", rate)
	}
	return subUrl
}

func GetJumpUrl(r *mvutil.RequestParams, params *mvutil.Params, campaign *smodel.CampaignInfo, queryQ string, ad *Ad) string {
	if params.Platform == 0 {
		params.Platform = mvconst.PlatformAndroid
	}
	jumpUrl := ""
	switch params.Extra10 {
	case mvconst.JUMP_TYPE_SDK_TO_MARKET:
		jumpUrl = getMarketUrl(params, campaign)
	case mvconst.JUMP_TYPE_CLIENT_SEND_DEVID_PING_SERVER:
		jumpUrl = createDirectUrl(r, params, campaign)
	case mvconst.JUMP_TYPE_CLIENT_DO_ALL_PING_SERVER:
		jumpUrl = createDirectUrl(r, params, campaign)
	case mvconst.JUMP_TYPE_TRACKING_PING_THIRDPARTY_CLICK_URL:
		jumpUrl = getMarketUrlByLinkType(params, campaign, ad, r)
	default:
		jumpUrl = CreateTrackUrl(r, params, campaign, "", queryQ)
	}
	return jumpUrl
}

func getMd5(str string) string {
	if len(str) <= 0 {
		return ""
	}
	str = strings.ToUpper(str)
	str = mvutil.Md5(str)
	str = strings.ToUpper(str)
	return str
}

func getMd5Lower(str string) string {
	if len(str) <= 0 {
		return ""
	}
	str = strings.ToLower(str)
	str = mvutil.Md5(str)
	return str
}

func getSha1(str string) string {
	if len(str) <= 0 {
		return ""
	}
	return mvutil.Sha1(str)
}

func handleUaParam(str string) string {
	return mvutil.UrlEncode(strings.ToLower(str))
}

func getAdvSubid(appID int64, chnID int) string {
	str := strconv.Itoa(chnID) + "_" + strconv.FormatInt(appID, 10)
	str = mvutil.Md5(str)
	str = mvutil.SubString(str, 0, 16)
	return str
}

var (
	appsflyerReg, _          = regexp.Compile("{af_aes_{\\S*?}}")
	matReg, _                = regexp.Compile("{mat_aes_{\\S*?}}")
	md5SignReg, _            = regexp.Compile("{md5{\\S*?}}")
	encodeURIComponentReg, _ = regexp.Compile("{encode_uri_component_{[\\s\\S]*?}}")

	appsflyerEncrypt = mvutil.NewAESECBEncrypt([]byte("wMOs8JrslMwnbK44"), 16)

	matPrivateKey = "15d09e2d839ae599992bcfd48a20e3d0"[:16]
	matKey        = "b4aa405de27c5c8dd0fea55d482ddcca"
	matEncrypt    = mvutil.NewAESCBCEncrypt([]byte(matKey), []byte(matPrivateKey))
)

// hardcode encrypt
func renderAppsflyerAESEncrypt(durl string) string {
	return renderReg(appsflyerReg, "{af_aes_{", appsflyerEncrypt.Encrypt, durl)
}

// hardcode encrypt
func renderMatAESEncrypt(durl string) string {
	return renderReg(matReg, "{mat_aes_{", matEncrypt.Encrypt, durl)
}

func renderSignMD5(durl string) string {
	return renderReg(md5SignReg, "{md5{", func(s string) (string, error) {
		return mvutil.Md5(s), nil
	}, durl)
}

func renderEncodeURIComponent(durl string) string {
	return renderReg(encodeURIComponentReg, "{encode_uri_component_{", func(s string) (string, error) {
		return url.QueryEscape(s), nil
	}, durl)
}

func renderReg(reg *regexp.Regexp, prefix string, renderFunc func(string) (string, error), durl string) string {
	if !strings.Contains(durl, prefix) {
		return durl
	}

	count := strings.Count(durl, prefix)
	for i := 0; i < count; i++ {
		matchStr := reg.FindString(durl)
		if matchStr == "" {
			continue
		}

		repStr := strings.Replace(matchStr, prefix, "", 1)
		if len(repStr) > 2 {
			repStr = repStr[:len(repStr)-2]
		}

		if repStr == "" {
			continue
		}

		encrypt, err := renderFunc(repStr)
		if err != nil {
			mvutil.Logger.Runtime.Warnf("%s raw:%v Encrypt:%v err:%v", prefix, matchStr, repStr, encrypt)
		}
		durl = strings.Replace(durl, matchStr, encrypt, 1)
	}
	return durl
}

func renderDirectUrl(dUrl string, r *mvutil.RequestParams, params *mvutil.Params, campaign *smodel.CampaignInfo) string {
	//-------- 清空设备切量- 11、12  +ss单子
	if dUrl[len(dUrl)-1] == '&' {
		dUrl = dUrl[:len(dUrl)-1]
	}

	if mvutil.InArray(params.CleanDeviceTest, []int{2, 6, 7}) {
		params.GAID = ""
		params.IDFA = ""
	}

	if params.CleanDeviceTest == 5 || params.CleanDeviceTest == 7 {
		dUrl += params.ThirdPartyInjectParams
	}

	// ---------

	if params.ExtdeleteDevid == mvconst.DELETE_DIVICEID_TRUE || params.ExtdeleteDevid == mvconst.DELETE_DIVICEID_BUT_NOT_IMPRESSION {
		params.GAID = ""
		params.IDFA = ""
		params.AndroidID = ""
	}

	switchs, _ := extractor.GetADNET_SWITCHS()
	if params.ExtDataInit.DemandLibABTest == 1 || switchs["adnetDemandLibABtest"] >= 100 {
		dUrl = renderSSOfferMacroByDemandLib(params, campaign, dUrl)
	} else {
		dUrl = renderThirdPartyUrl(params, campaign, dUrl)
	}

	//dUrl = removeAdjustCallback(dUrl, params, campaign)
	// 对于命中的情况，针对903，三方为af的单子,rv,iv流量的情况下，添加token宏
	if params.ExtDataInit.ReturnAfTokenParam == 1 {
		dUrl += mvconst.APPSFLYER_TOKEN_PARAMS
	}
	return dUrl
}

func renderSSOfferMacroByDemandLib(params *mvutil.Params, campaign *smodel.CampaignInfo, trackURL string) string {
	if params.DemandContext == nil {
		return renderThirdPartyUrl(params, campaign, trackURL)
	}
	// 清空设备实验遗留物。后续需要干掉这个实验 TODO
	params.DemandContext.GAID = params.GAID
	params.DemandContext.IDFA = params.IDFA
	params.DemandContext.AndroidId = params.AndroidID
	return demand.RenderThirdPartyLink(params.DemandContext, extractor.DemandDao, trackURL)
}

func renderThirdPartyUrl(params *mvutil.Params, campaign *smodel.CampaignInfo, dUrl string) string {
	gaid := params.GAID
	idfa := params.IDFA
	devId := params.AndroidID
	var idfaMd5 string
	if len(params.IDFAMd5) == 0 {
		idfaMd5 = getMd5(idfa)
	} else {
		idfaMd5 = strings.ToUpper(params.IDFAMd5)
	}
	idfaSha1 := getSha1(idfa)
	var gaidMd5 string
	if len(params.GAIDMd5) == 0 {
		gaidMd5 = getMd5(gaid)
	} else {
		gaidMd5 = strings.ToUpper(params.GAIDMd5)
	}
	gaidSha1 := getSha1(gaid)
	var devIdMd5 string
	if len(params.AndroidIDMd5) == 0 {
		devIdMd5 = getMd5(devId)
	} else {
		devIdMd5 = strings.ToUpper(params.AndroidIDMd5)
	}
	devIdSha1 := getSha1(devId)
	imei := params.IMEI
	var imeiMd5 string
	if len(params.ImeiMd5) == 0 {
		imeiMd5 = getMd5Lower(imei)
	} else {
		imeiMd5 = strings.ToLower(params.ImeiMd5)
	}
	imeiSha1 := getSha1(imei)
	mac := params.MAC
	macMd5 := getMd5(mac)
	macSha1 := getSha1(mac)
	gaidDevId := gaid
	if len(gaid) <= 0 {
		gaidDevId = devId
	}
	ip := params.ClientIP
	countryCode := strings.ToLower(params.CountryCode)
	ua := params.UserAgent
	uaOsPlatform := handleUaParam(params.PlatformName)
	uaOsVersion := handleUaParam(params.OSVersion)
	uaDeviceModel := handleUaParam(params.Model)
	uaInfo := mvutil.UaParser.Parse(ua)
	uaOs := mvutil.UrlEncode(uaInfo.Os.Family + " " + uaInfo.Os.ToVersionString())
	ua = mvutil.UrlEncode(ua)
	timestamp := time.Now().Unix()
	microtime := time.Now().UnixNano() / 1e6
	cityString := mvutil.UrlEncode(params.CityString)
	mgid := "mtg" + params.RequestID
	//appId := params.AppID
	subId := getSubId(params)
	chnId := campaign.ChnID
	advSubId := getAdvSubid(subId, chnId)
	mbSubId := "mob" + advSubId
	cMbSubId := "c_" + advSubId
	packageName := mvutil.UrlEncode(params.ExtfinalPackageName)
	packageMbSubId := packageName
	if len(packageMbSubId) <= 0 {
		packageMbSubId = mvutil.UrlEncode(mbSubId)
	}
	// 大写的gaid或idfa
	gaidOrIdfa := ""
	if params.Platform == mvconst.PlatformAndroid {
		gaidOrIdfa = gaid
	} else if params.Platform == mvconst.PlatformIOS {
		gaidOrIdfa = idfa
	}
	gaidOrIdfa = strings.ToUpper(gaidOrIdfa)
	// 大写的countrycode
	upperCountryCode := strings.ToUpper(params.CountryCode)
	// 素材独享逻辑 TODO
	crat := strconv.Itoa(params.CreativeAdType)
	// advCreativeid
	advCrid := strconv.Itoa(params.AdvCreativeID)
	priceIn := mvutil.FormatFloat64(params.PriceIn)

	dUrl = strings.Replace(dUrl, "{mgid}", mgid, -1)
	dUrl = strings.Replace(dUrl, "{mbSubId}", mbSubId, -1)
	dUrl = strings.Replace(dUrl, "{ip}", ip, -1)
	dUrl = strings.Replace(dUrl, "{package_name}", packageName, -1)
	dUrl = strings.Replace(dUrl, "{gaid_devId}", gaidDevId, -1)
	dUrl = strings.Replace(dUrl, "{ua}", ua, -1)
	dUrl = strings.Replace(dUrl, "{microtime}", strconv.FormatInt(microtime, 10), -1)
	dUrl = strings.Replace(dUrl, "{countryCode}", countryCode, -1)
	dUrl = strings.Replace(dUrl, "{package_mbSubId}", packageMbSubId, -1)
	dUrl = strings.Replace(dUrl, "{c_mbSubId}", cMbSubId, -1)
	dUrl = strings.Replace(dUrl, "{idfa}", idfa, -1)
	dUrl = strings.Replace(dUrl, "{idfaMd5}", idfaMd5, -1)
	dUrl = strings.Replace(dUrl, "{idfaSha1}", idfaSha1, -1)
	dUrl = strings.Replace(dUrl, "{gaid}", gaid, -1)
	dUrl = strings.Replace(dUrl, "{gaidMd5}", gaidMd5, -1)
	dUrl = strings.Replace(dUrl, "{gaidSha1}", gaidSha1, -1)
	dUrl = strings.Replace(dUrl, "{devId}", devId, -1)
	dUrl = strings.Replace(dUrl, "{devIdMd5}", devIdMd5, -1)
	dUrl = strings.Replace(dUrl, "{devIdSha1}", devIdSha1, -1)
	dUrl = strings.Replace(dUrl, "{imei}", imei, -1)
	dUrl = strings.Replace(dUrl, "{oaid}", url.QueryEscape(params.OAID), -1)
	dUrl = strings.Replace(dUrl, "{imeiMd5}", imeiMd5, -1)
	dUrl = strings.Replace(dUrl, "{imeiSha1}", imeiSha1, -1)
	dUrl = strings.Replace(dUrl, "{mac}", mac, -1)
	dUrl = strings.Replace(dUrl, "{macMd5}", macMd5, -1)
	dUrl = strings.Replace(dUrl, "{macSha1}", macSha1, -1)
	dUrl = strings.Replace(dUrl, "{uaDevice}", uaDeviceModel, -1)
	dUrl = strings.Replace(dUrl, "{uaOsPlatform}", uaOsPlatform, -1)
	dUrl = strings.Replace(dUrl, "{uaOsVersion}", uaOsVersion, -1)
	dUrl = strings.Replace(dUrl, "{timestamp}", strconv.FormatInt(timestamp, 10), -1)
	dUrl = strings.Replace(dUrl, "{city}", cityString, -1)
	dUrl = strings.Replace(dUrl, "{uaOs}", uaOs, -1)
	dUrl = strings.Replace(dUrl, "{gaid_idfa}", gaidOrIdfa, -1)
	dUrl = strings.Replace(dUrl, "{upperCountryCode}", upperCountryCode, -1)
	dUrl = strings.Replace(dUrl, "{creativeId}", advCrid, -1)
	dUrl = strings.Replace(dUrl, "{adType}", crat, -1)
	dUrl = strings.Replace(dUrl, "{creativeName}", params.CreativeName, -1)
	dUrl = strings.Replace(dUrl, "{price}", priceIn, -1)
	dUrl = strings.Replace(dUrl, "{lang}", url.QueryEscape(params.Language), -1)
	dUrl = strings.Replace(dUrl, "{network}", params.NetworkTypeName, -1)

	// mtgId
	mtgIdStr := "mtg" + strconv.FormatInt(params.ExtMtgId, 10)
	dUrl = strings.Replace(dUrl, "{mtgId}", mtgIdStr, -1)

	// 顶级宏替换
	dUrl = renderAppsflyerAESEncrypt(dUrl)
	dUrl = renderMatAESEncrypt(dUrl)
	dUrl = renderSignMD5(dUrl)
	dUrl = renderEncodeURIComponent(dUrl)
	return dUrl
}

func addAppsflyerTokenParam(params *mvutil.Params, campaign *smodel.CampaignInfo) {
	advId := int(campaign.AdvertiserId)
	//仅针对mv sdk,rv,iv,appsflyer,903的流量
	if !mvutil.IsHbOrV3OrV5Request(params.RequestPath) || !mvutil.IsIvOrRv(params.AdType) ||
		strings.ToLower(campaign.ThirdParty) != mvconst.THIRD_PARTY_APPSFLYER || advId != 903 {
		return
	}
	adnetConf, _ := extractor.GetADNET_SWITCHS()
	iosMinVer := 99000000
	androidMinVer := 99000000
	// 获取ios，Android 对应的os_version
	if iosMinVerConf, ok := adnetConf["iosMinVer"]; ok {
		iosMinVer = iosMinVerConf
	}
	if androidMinVerConf, ok := adnetConf["androidMinVer"]; ok {
		androidMinVer = androidMinVerConf
	}
	// 安卓情况下，9.3.0不支持传campaign信息给sdk，os_version在4.3.1以下不支持af_token，下发af_token不支持放置到adtracking中。所以这部分流量不能做此逻辑
	if params.Platform == mvconst.PlatformAndroid && (params.FormatSDKVersion.SDKVersionCode < mvconst.MoreOfferBlock || params.OSVersionCode < int32(androidMinVer) || params.AdClickReplaceclickByClickTag) {
		return
	}
	// ios 情况下，os_version在11.4以下不支持af_token
	if params.Platform == mvconst.PlatformIOS && params.OSVersionCode < int32(iosMinVer) {
		return
	}
	// ss单子，或者3s单子clickmode为11,12才能做切量
	if campaign.IsSSPlatform() ||
		params.Extra10 == mvconst.JUMP_TYPE_CLIENT_SEND_DEVID_PING_SERVER || params.Extra10 == mvconst.JUMP_TYPE_CLIENT_DO_ALL_PING_SERVER {
		if returnAfTokenParamRate, ok := adnetConf["rtAfTokenRate"]; ok {
			randRate := rand.Intn(100)
			if returnAfTokenParamRate > randRate {
				params.ExtDataInit.ReturnAfTokenParam = 1
			} else {
				params.ExtDataInit.ReturnAfTokenParam = 2
			}
		}
	}
}

// 跳转链路优化 TODO
func createDirectUrl(r *mvutil.RequestParams, params *mvutil.Params, campaign *smodel.CampaignInfo) string {
	if len(campaign.DirectUrl) == 0 {
		return ""
	}

	return renderDirectUrl(campaign.DirectUrl, r, params, campaign)
}

func getMarketUrl(params *mvutil.Params, campaign *smodel.CampaignInfo) string {
	marketUrl := ""
	conf, ifFind := extractor.GetADSTACKING()
	if !ifFind {
		return marketUrl
	}
	if params.Platform == mvconst.PlatformAndroid {
		marketUrl = conf.Android
	} else if params.Platform == mvconst.PlatformIOS {
		marketUrl = conf.IOS
	}
	if len(marketUrl) <= 0 {
		return ""
	}
	if campaign.PackageName == "" {
		return ""
	}
	marketUrl = strings.Replace(marketUrl, "{package_name}", campaign.PackageName, -1)
	return marketUrl
}

func getMarketUrlByLinkType(params *mvutil.Params, campaign *smodel.CampaignInfo, ad *Ad, r *mvutil.RequestParams) string {
	if mvutil.InArray(ad.CampaignType, []int{mvconst.LINK_TYPE_APPSTORE, mvconst.LINK_TYPE_GOOGLEPLAY}) {
		return getMarketUrl(params, campaign)
	}

	ctx := NewDemandContext(params, r)
	if ad.CampaignType == mvconst.LINK_TYPE_APK && campaign.ApkUrl != "" {
		adnetSwitchConf, _ := extractor.GetADNET_SWITCHS()
		useApkUrl, ok := adnetSwitchConf["useApkUrl"]
		if ok && useApkUrl == 1 {
			return campaign.ApkUrl
		}
	}

	previewUrl := demand.GetCampaignNormalLandingPage(campaign)
	if previewUrl != "" {
		return demand.RenderThirdPartyLink(ctx, extractor.DemandDao, previewUrl)
	}
	return ""
}

func CreateTrackUrl(r *mvutil.RequestParams, params *mvutil.Params, campaign *smodel.CampaignInfo, trackUrl string, queryQ string) string {
	if len(trackUrl) <= 0 {
		if len(campaign.TrackingUrl) == 0 {
			return ""
		}
		trackUrl = campaign.TrackingUrl

		// android https
		if params.Platform == mvconst.PlatformAndroid && NeedSchemeHttps(params.HTTPReq) && len(campaign.TrackingUrlHttps) > 0 {
			trackUrl = campaign.TrackingUrlHttps
		}
	}

	if campaign.IsSSPlatform() {
		trackUrl = renderDirectUrl(trackUrl, r, params, campaign)
		return trackUrl
	}

	// append url
	if campaign.Network != 0 {
		network := campaign.Network
		confs, ifFind := extractor.GetTRACK_URL_CONFIG_NEW()
		if ifFind {
			conf, ok := confs[network]
			if ok {
				urlAppend := ""
				if params.Platform == mvconst.PlatformAndroid {
					urlAppend = conf.Android
				} else if params.Platform == mvconst.PlatformIOS {
					urlAppend = conf.IOS
				}
				// 针对小程序做处理
				if network == mvconst.NETWORK_SMALL_ROUTINE {
					// 小程序相关替换规则
					trackUrl = renderSmallRoutineUrl(trackUrl, params, campaign.CampaignId, urlAppend)
				} else {
					trackUrl = trackUrl + urlAppend
				}
			}
		}
	}

	if params.ExtdeleteDevid == mvconst.DELETE_DIVICEID_TRUE || params.ExtdeleteDevid == mvconst.DELETE_DIVICEID_BUT_NOT_IMPRESSION {
		params.GAID = ""
		params.IDFA = ""
		params.AndroidID = ""
	}

	// gaid
	trackUrl = strings.Replace(trackUrl, "{gaid}", params.GAID, -1)
	// idfa
	trackUrl = strings.Replace(trackUrl, "{idfa}", params.IDFA, -1)
	// devid
	devId := getDevId(r, params)
	trackUrl = strings.Replace(trackUrl, "{devId}", devId, -1)
	// imei
	trackUrl = strings.Replace(trackUrl, "{imei}", params.IMEI, -1)
	// oaid
	trackUrl = strings.Replace(trackUrl, "{oaid}", url.QueryEscape(params.OAID), -1)
	// mac
	trackUrl = strings.Replace(trackUrl, "{mac}", params.MAC, -1)

	idfa := params.IDFA
	var idfaMd5 string
	if len(params.IDFAMd5) == 0 {
		idfaMd5 = getMd5(idfa)
	} else {
		idfaMd5 = strings.ToUpper(params.IDFAMd5)
	}
	idfaSha1 := getSha1(idfa)

	gaid := params.GAID
	var gaidMd5 string
	if len(params.GAIDMd5) == 0 {
		gaidMd5 = getMd5(gaid)
	} else {
		gaidMd5 = strings.ToUpper(params.GAIDMd5)
	}
	gaidSha1 := getSha1(gaid)

	devID := params.AndroidID
	var devIdMd5 string
	if len(params.AndroidIDMd5) == 0 {
		devIdMd5 = getMd5(devID)
	} else {
		devIdMd5 = strings.ToUpper(params.AndroidIDMd5)
	}
	devIdSha1 := getSha1(devID)

	imei := params.IMEI
	var imeiMd5 string
	if len(params.ImeiMd5) == 0 {
		imeiMd5 = getMd5Lower(imei)
	} else {
		imeiMd5 = strings.ToLower(params.ImeiMd5)
	}
	imeiSha1 := getSha1(imei)

	mac := params.MAC
	macMd5 := getMd5(mac)
	macSha1 := getSha1(mac)

	trackUrl = strings.Replace(trackUrl, "{gaidMd5}", gaidMd5, -1)
	trackUrl = strings.Replace(trackUrl, "{gaidSha1}", gaidSha1, -1)
	trackUrl = strings.Replace(trackUrl, "{imeiMd5}", imeiMd5, -1)
	trackUrl = strings.Replace(trackUrl, "{imeiSha1}", imeiSha1, -1)
	trackUrl = strings.Replace(trackUrl, "{idfaMd5}", idfaMd5, -1)
	trackUrl = strings.Replace(trackUrl, "{idfaSha1}", idfaSha1, -1)
	trackUrl = strings.Replace(trackUrl, "{macMd5}", macMd5, -1)
	trackUrl = strings.Replace(trackUrl, "{macSha1}", macSha1, -1)
	trackUrl = strings.Replace(trackUrl, "{devIdMd5}", devIdMd5, -1)
	trackUrl = strings.Replace(trackUrl, "{devIdSha1}", devIdSha1, -1)

	trackUrl = strings.Replace(trackUrl, "{network}", params.NetworkTypeName, -1)
	// packageName
	trackUrl = strings.Replace(trackUrl, "{package_name}", params.ExtfinalPackageName, -1)
	// subId
	subId := getSubId(params)
	subIdStr := strconv.FormatInt(subId, 10)
	trackUrl = strings.Replace(trackUrl, "{subId}", subIdStr, -1)
	// mtgId
	mtgIdStr := "mtg" + strconv.FormatInt(params.ExtMtgId, 10)
	trackUrl = strings.Replace(trackUrl, "{mtgId}", mtgIdStr, -1)
	// clickId
	if len(queryQ) <= 0 {
		queryQ = mvutil.SerializeQ(params, campaign)
	}
	queryQ = mvutil.UrlEncode(queryQ)
	trackUrl = strings.Replace(trackUrl, "{clickId}", queryQ, -1)
	// 素材独享逻辑 TODO
	crat := strconv.Itoa(params.CreativeAdType)
	advCrid := strconv.Itoa(params.AdvCreativeID)
	trackUrl = strings.Replace(trackUrl, "{creativeId}", advCrid, -1)
	trackUrl = strings.Replace(trackUrl, "{adType}", crat, -1)
	trackUrl = strings.Replace(trackUrl, "{creativeName}", params.CreativeName, -1)
	trackUrl = strings.Replace(trackUrl, "{lang}", url.QueryEscape(params.Language), -1)
	trackUrl = renderEncodeURIComponent(trackUrl)

	// 系统是SA，替换域名
	system := extractor.GetSYSTEM()
	if system == mvconst.SERVER_SYSTEM_SA {
		domain := extractor.GetDOMAIN()
		domainTrack := extractor.GetDOMAIN_TRACK()
		trackUrl = strings.Replace(trackUrl, "http://net.rayjump.com", "http://"+domain, -1)
		trackUrl = strings.Replace(trackUrl, "http://tknet.rayjump.com", "http://"+domainTrack, -1)
	}
	// 替换agent/click的https
	if params.HTTPReq == 2 {
		trackUrl = strings.Replace(trackUrl, "http://net.rayjump.com", "https://net.rayjump.com", -1)
		trackUrl = strings.Replace(trackUrl, "http://tknet.rayjump.com", "https://tknet.rayjump.com", -1)
	}

	// 3s CN domain
	trackUrl = change3SCNDomain(trackUrl, params)
	// dsp domain
	// changeDspDomain

	return trackUrl
}

// func changeDspDomain(trackUrl *string, params mvutil.Params) {
// 	if params.PublisherID != mvconst.DspPublisherID {
// 		return
// 	}
// 	if params.
// }

func change3SCNDomain(trackUrl string, params *mvutil.Params) string {
	conf, ifFind := extractor.Get3S_CHINA_DOMAIN()
	if !ifFind {
		return trackUrl
	}
	countrys := conf.Countrys
	if !mvutil.InStrArray(params.CountryCode, countrys) {
		return trackUrl
	}
	domains := conf.Domains
	cnDomain := conf.CNDomain
	if len(domains) <= 0 || len(cnDomain) <= 0 {
		return trackUrl
	}
	for _, v := range domains {
		if strings.Contains(trackUrl, v) {
			if params.CNDomainTest == 2 {
				cnDomain = conf.CNLineDo
			}
			trackUrl = strings.Replace(trackUrl, v, cnDomain, -1)
			trackUrl = trackUrl + "&mb_trackingcn=1"
			return trackUrl
		}
	}
	return trackUrl
}

func Is3SCNLine(r *mvutil.RequestParams) {
	conf, ifFind := extractor.Get3S_CHINA_DOMAIN()
	if !ifFind {
		return
	}
	countrys := conf.Countrys
	if !mvutil.InStrArray(r.Param.CountryCode, countrys) {
		return
	}
	uid := strconv.FormatInt(r.Param.UnitID, 10)
	if !mvutil.InStrArray(uid, conf.Units) && !mvutil.InStrArray("ALL", conf.Units) {
		return
	}
	// rand
	randInt := mvutil.GetRandConsiderZero(r.Param.GAID, r.Param.IDFA, mvconst.SALT_3SLINE, 100)
	if randInt == -1 {
		return
	}
	if randInt < conf.Rate {
		r.Param.CNDomainTest = 2
	} else {
		r.Param.CNDomainTest = 3
	}
}

func createClickUrl(params *mvutil.Params, isNotice bool) string {
	res := ""
	// 如果走af 白名单通道，则不加notice参数,tracking根据clickmode=13ping 3s和三方
	if params.ExtDataInit.ClickInServer == 1 || params.ExtDataInit.ClickInServer == 3 || params.Extra10 == mvconst.JUMP_TYPE_TRACKING_PING_THIRDPARTY_CLICK_URL {
		isNotice = false
	}
	// new url
	if params.IsNewUrl {
		res = "{sh}://{do}/click?k={k}&z={z}&q={q}&r={r}&al={al}&csp={csp}&c={c}"
		if len(params.ReplacedClickTrackDomain) > 0 {
			res = "{sh}://{cdo}/click?k={k}&z={z}&q={q}&r={r}&al={al}&csp={csp}&c={c}"
		}
		if isNotice {
			res = res + "&notice=1"
		}
		if params.PublisherID == mvconst.TOPONADXPublisherID && params.IsHBRequest {
			res += "&tpua={tpua}"
		}
		return res
	}

	queryStr := "k=" + params.RequestID
	domain := params.Domain

	// 替换tracking domain
	if len(params.ReplacedClickTrackDomain) > 0 {
		domain = params.ReplacedClickTrackDomain
	}

	queryStr += "&p=" + params.QueryP + "&q=" + params.QueryQ + "&r=" + params.QueryR + "&al=" + params.QueryAL + "&csp=" + params.QueryCSP
	if isNotice {
		queryStr = queryStr + "&notice=1"
	}

	res = GetUrlScheme(params) + "://" + domain + "/click?" + queryStr
	return res
}

func RepalceDoMacro(domain string, ad *Ad, params *mvutil.Params) {
	if params.ApiVersion < mvconst.API_VERSION_1_4 {
		return
	}
	if !mvutil.IsHbOrV3OrV5Request(params.RequestPath) {
		return
	}
	if len(domain) == 0 {
		return
	}
	if len(ad.ClickURL) > 0 {
		ad.ClickURL = strings.Replace(ad.ClickURL, "{do}", domain, -1)
	}
	if len(ad.NoticeURL) > 0 {
		ad.NoticeURL = strings.Replace(ad.NoticeURL, "{do}", domain, -1)
	}
	if len(ad.PreviewUrl) > 0 {
		ad.PreviewUrl = strings.Replace(ad.PreviewUrl, "{do}", domain, -1)

	}
	if len(ad.ImpressionURL) > 0 {
		ad.ImpressionURL = strings.Replace(ad.ImpressionURL, "{do}", domain, -1)
	}
	if len(ad.DeepLink) > 0 {
		ad.DeepLink = strings.Replace(ad.DeepLink, "{do}", domain, -1)
	}
	ReplaceAdTrackingDoMacro(ad.AdTrackingPoint, domain)
}

func ReplaceAdTrackingDoMacro(adTracking *CAdTracking, domain string) {
	if adTracking == nil {
		return
	}
	adTracking.Impression = ReplaceUrlsDoMacro(adTracking.Impression, domain)
	adTracking.Click = ReplaceUrlsDoMacro(adTracking.Click, domain)
	ReplaceAdTrackingDoMacroExcludeImpAndClick(adTracking, domain)
}

func ReplaceAdTrackingDoMacroExcludeImpAndClick(adTracking *CAdTracking, domain string) {
	if adTracking == nil {
		return
	}
	adTracking.Start = ReplaceUrlsDoMacro(adTracking.Start, domain)
	adTracking.First_quartile = ReplaceUrlsDoMacro(adTracking.First_quartile, domain)
	adTracking.Midpoint = ReplaceUrlsDoMacro(adTracking.Midpoint, domain)
	adTracking.Third_quartile = ReplaceUrlsDoMacro(adTracking.Third_quartile, domain)
	adTracking.Complete = ReplaceUrlsDoMacro(adTracking.Complete, domain)
	adTracking.Mute = ReplaceUrlsDoMacro(adTracking.Mute, domain)
	adTracking.Unmute = ReplaceUrlsDoMacro(adTracking.Unmute, domain)
	adTracking.Endcard_show = ReplaceUrlsDoMacro(adTracking.Endcard_show, domain)
	adTracking.Close = ReplaceUrlsDoMacro(adTracking.Close, domain)
	for k, v := range adTracking.Play_percentage {
		adTracking.Play_percentage[k].Url = ReplaceUrlsDoMacro([]string{v.Url}, domain)[0]
	}
	adTracking.Pause = ReplaceUrlsDoMacro(adTracking.Pause, domain)
	adTracking.Video_Click = ReplaceUrlsDoMacro(adTracking.Video_Click, domain)
	adTracking.Impression_t2 = ReplaceUrlsDoMacro(adTracking.Impression_t2, domain)
	adTracking.ApkDownloadStart = ReplaceUrlsDoMacro(adTracking.ApkDownloadStart, domain)
	adTracking.ApkDownloadEnd = ReplaceUrlsDoMacro(adTracking.ApkDownloadEnd, domain)
	adTracking.ApkInstall = ReplaceUrlsDoMacro(adTracking.ApkInstall, domain)
	adTracking.Dropout = ReplaceUrlsDoMacro(adTracking.Dropout, domain)
	adTracking.Plycmpt = ReplaceUrlsDoMacro(adTracking.Plycmpt, domain)
	adTracking.PubImp = ReplaceUrlsDoMacro(adTracking.PubImp, domain)
	adTracking.ExaClick = ReplaceUrlsDoMacro(adTracking.ExaClick, domain)
	adTracking.ExaImp = ReplaceUrlsDoMacro(adTracking.ExaImp, domain)
}

func ReplaceUrlsDoMacro(urls []string, domain string) []string {
	if len(urls) == 0 {
		return []string{}
	}
	for key, str := range urls {
		str = strings.Replace(str, "{do}", domain, -1)
		urls[key] = str
	}
	return urls
}

func CreateClickUrl(params *mvutil.Params, isNotice bool) string {
	return createClickUrl(params, isNotice)
}

func createImpressionUrl(params *mvutil.Params) string {
	if params.IAADRst == mvconst.IARST_APPWALL {
		return ""
	}
	x := strconv.Itoa(params.ImpressionImage)
	res := ""
	if params.IsNewUrl {
		res = "{sh}://{do}/impression?k={k}&z={z}&q={q}&x=" + x + "&r={r}&al={al}&csp={csp}&c={c}"
		if len(params.ReplacedImpTrackDomain) > 0 {
			res = "{sh}://{ido}/impression?k={k}&z={z}&q={q}&x=" + x + "&r={r}&al={al}&csp={csp}&c={c}"
		}
		if params.PublisherID == mvconst.TOPONADXPublisherID && params.IsHBRequest {
			res += "&tpua={tpua}"
		}
		// token_r值为1则在展示点击url中下发{encrypt_p},{irlfa}的宏
		if params.TokenRule == 1 && params.IsHBRequest {
			res += "&encrypt_p={encrypt_p}&irlfa={irlfa}"
		}
		return res
	}

	queryStr := "k=" + params.RequestID
	domain := params.Domain

	if len(params.ReplacedImpTrackDomain) > 0 {
		domain = params.ReplacedImpTrackDomain
	}

	queryStr += "&p=" + params.QueryP + "&q=" + params.QueryQ + "&x=" + x + "&r=" + params.QueryR + "&al=" + params.QueryAL + "&csp=" + params.QueryCSP
	res = GetUrlScheme(params) + "://" + domain + "/impression?" + queryStr
	return res
}

func CreateImpressionUrl(params *mvutil.Params) string {
	return createImpressionUrl(params)
}

func CreateCsp(params *mvutil.Params, dspExt string) string {
	if params.MWFlowTagID <= 0 || params.Extra == "adnet_tt2" {
		return ""
	}
	// backend id为1的情况无dspext。
	if params.BackendID == mvconst.Mobvista {
		dspExt = ""
	}
	return mvutil.SerializeCSP(params, dspExt)
}

func RenderHttpsUrls(ad *Ad, params mvutil.Params) {
	if !NeedSchemeHttps(params.HTTPReq) {
		return
	}

	if len(ad.IconURL) > 0 {
		ad.IconURL = renderCDNUrl2Https(ad.IconURL)
	}
	if len(ad.ImageURL) > 0 {
		ad.ImageURL = renderCDNUrl2Https(ad.ImageURL)
	}
	if len(ad.VideoURL) > 0 {
		videoUrl := mvutil.Base64Decode(ad.VideoURL)
		if len(videoUrl) > 0 {
			videoUrl = renderCDNUrl2Https(videoUrl)
			ad.VideoURL = mvutil.Base64Encode(videoUrl)
		}
	}
	if len(ad.ExtImg) > 0 {
		ad.ExtImg = renderCDNUrl2Https(ad.ExtImg)
	}
	if len(ad.GifURL) > 0 {
		ad.GifURL = renderCDNUrl2Https(ad.GifURL)
	}
	if ad.Rv.Image != nil && len(ad.Rv.Image.IdcdImg) > 0 {
		for k, v := range ad.Rv.Image.IdcdImg {
			ad.Rv.Image.IdcdImg[k] = renderCDNUrl2Https(v)
		}
	}
}

func RenderCreativeUrls(ad *Ad, r *mvutil.RequestParams, params *mvutil.Params) {
	// 针对小程序广告位，不做cdn域名替换
	if mvutil.IsWxAdType(r.Param.AdType) {
		return
	}
	newCDNUrl := ""
	// unit维度切量配置
	if len(r.UnitInfo.CdnSetting) > 0 {
		newCDNUrl = getNewCdnUrl(r.UnitInfo.CdnSetting, params, "")
	}
	// 判断广告主+地区维度
	if len(newCDNUrl) == 0 {
		cdnConfig, _ := extractor.GetPUB_CC_CDN()
		pubStr := strconv.FormatInt(r.Param.PublisherID, 10)
		pubConfig, ok := cdnConfig[pubStr]
		if ok {
			// countrycode维度
			newCDNUrl = pubConfig[r.Param.CountryCode]
			// ALL维度
			if len(newCDNUrl) == 0 {
				newCDNUrl = pubConfig["ALL"]
			}
		}
	}
	cdnConfs := extractor.GetNEW_CDN_TEST()
	// 整体切量配置
	if len(newCDNUrl) == 0 {
		// 只取icon url的域名来判断就好了
		oldHost := getOldHost(ad.IconURL)
		// 优先获取cdn+cc配置，无配置则获取cdn+ALL配置
		cdnConf := getCdnConf(cdnConfs, params, oldHost)
		if len(cdnConf) > 0 {
			newCDNUrl = getNewCdnUrl(cdnConf, params, oldHost)
		}
	}
	// 若newCDNUrl存在值才替换
	if len(newCDNUrl) > 0 {
		// 切换新的cdn的切量配置
		replaceCdnDomain(ad, newCDNUrl, params.HTTPReq)
	}

	// endcard 域名（playable.rayjump.com）abtest。因为ucloud太贵了，因此需要切其他CDN，因此需要做abtest
	if len(ad.EndcardUrl) > 0 {
		oldHost := getOldHost(ad.EndcardUrl)
		var newEndcardCDNUrl string
		cdnConf := getCdnConf(cdnConfs, params, oldHost)
		if len(cdnConf) > 0 {
			newEndcardCDNUrl = getEndcardNewCdnUrl(cdnConf, params)
		}
		if len(newEndcardCDNUrl) > 0 {
			ad.EndcardUrl = replaceSchemeAndHost(ad.EndcardUrl, newEndcardCDNUrl, params.HTTPReq)
		}
	}

}

func getCdnConf(cdnConfs map[string]map[string][]*smodel.CdnSetting, params *mvutil.Params, oldHost string) []*smodel.CdnSetting {
	ccConf, ok := cdnConfs[oldHost]
	if !ok {
		return []*smodel.CdnSetting{}
	}
	if cdnConf, ok := ccConf[params.CountryCode]; ok && len(cdnConf) > 0 {
		return cdnConf
	}
	if cdnConf, ok := ccConf["ALL"]; ok && len(cdnConf) > 0 {
		return cdnConf
	}
	return []*smodel.CdnSetting{}
}

// 读取配置，根据权重获取最终替换的cdn地址，并且记录cdn标记
func getEndcardNewCdnUrl(cdnSetting []*smodel.CdnSetting, params *mvutil.Params) string {
	newCDNUrl := ""
	//var cdnMap map[string]int
	cdnMap := make(map[string]int)
	cdnId := make(map[string]int)
	for _, cdnSetting := range cdnSetting {
		cdnMap[cdnSetting.Url] = cdnSetting.Weight
		cdnId[cdnSetting.Url] = cdnSetting.Id
	}
	newCDNUrl = mvutil.RandByRate2(cdnMap)
	// 记录cdn测试标记
	if id, ok := cdnId[newCDNUrl]; ok {
		params.ExtDataInit.ECCDNTag = id
	}
	return newCDNUrl
}

func getOldHost(creUrl string) string {
	oldCdn, err := url.Parse(creUrl)
	oldHost := ""
	if err == nil {
		oldHost = oldCdn.Host
	}
	// mongo的key不能有.，使用_代替.查询数据
	oldHost = strings.Replace(oldHost, ".", "_", -1)
	return oldHost
}

func toNewCDN(ad *Ad, params *mvutil.Params) {
	// 针对小程序广告位，不做cdn域名替换
	if mvutil.IsWxAdType(params.AdType) {
		return
	}
	if !NeedSchemeHttps(params.HTTPReq) {
		return
	}
	// 切换新的cdn的切量配置
	toNewCDNAppList, ifFind := extractor.GetTO_NEW_CDN_APPS()
	if !(ifFind && mvutil.InInt64Arr(params.AppID, toNewCDNAppList)) {
		return
	}
	params.ExtCDNAbTest = 1
	if rand.Intn(2) != 1 {
		return
	}
	params.ExtCDNAbTest = 2

	if len(ad.IconURL) > 0 {
		ad.IconURL = renderNewCDNUrl(ad.IconURL)
	}
	if len(ad.ImageURL) > 0 {
		ad.ImageURL = renderNewCDNUrl(ad.ImageURL)
	}
	if len(ad.VideoURL) > 0 {
		videoUrl := mvutil.Base64Decode(ad.VideoURL)
		if len(videoUrl) > 0 {
			videoUrl = renderNewCDNUrl(videoUrl)
			ad.VideoURL = mvutil.Base64Encode(videoUrl)
		}
	}
	if len(ad.ExtImg) > 0 {
		ad.ExtImg = renderNewCDNUrl(ad.ExtImg)
	}
}

func renderCDNUrl2Https(oriUrl string) string {
	oriUrl = strings.Replace(oriUrl, "http://cdn-adn.rayjump.com", "https://cdn-adn-https.rayjump.com", -1)
	oriUrl = strings.Replace(oriUrl, "http://res.rayjump.com", "https://res-https.rayjump.com", -1)
	oriUrl = strings.Replace(oriUrl, "http://d11kdtiohse1a9.cloudfront.net", "https://res-https.rayjump.com", -1)
	return oriUrl
}

func renderNewCDNUrl(url string) string {
	url = strings.Replace(url, "https://cdn-adn-https.rayjump.com", "https://cdn-adn-https-new.rayjump.com", -1)
	return url
}

func renderNewCDNUrl2(url, newUrl string) string {
	url = strings.Replace(url, "cdn-adn-https.rayjump.com", newUrl, -1)
	url = strings.Replace(url, "cdn-adn.rayjump.com", newUrl, -1)
	url = strings.Replace(url, "res.rayjump.com", newUrl, -1)
	url = strings.Replace(url, "d11kdtiohse1a9.cloudfront.net", newUrl, -1)
	url = strings.Replace(url, "res-https.rayjump.com", newUrl, -1)
	return url
}

func CreateOnlyImpressionUrl(params mvutil.Params, r *mvutil.RequestParams) string {
	if params.OnlyImpression != 1 {
		return ""
	}
	params.Extra14 = params.AppID
	params.Extalgo = ""

	if params.Mof == 1 && params.MofVersion >= 2 && params.AdType != mvconst.ADTypeMoreOffer {
		params.Extra20 = ""
	}

	var ad Ad
	subUrl := ""
	if params.MWFlowTagID > 0 {
		RenderMWParams(&params, &ad, true)
	}
	// 处理rs参数
	if len(params.QueryRsList) > 0 {
		queryRslist, err := json.Marshal(params.QueryRsList)
		var rsStr string
		if err == nil {
			rsStr = mvutil.Base64(queryRslist)
		}
		params.QueryRs = rsStr // adnet自己的赋值逻辑
	}
	dspExt, _ := r.GetDspExt()

	// 如果是第三方独占
	// native more ads 当第三方胜出且有多个单子时， 不能进这里
	if (!params.IsFakeAs && len(params.BackendList) > 0) && (dspExt == nil ||
		(dspExt.DspId != mvconst.MAS && !r.IsMoreAsAds)) && !mvutil.IsRequestPioneerDirectly(&params) {
		if !mvutil.InInt32Arr(int32(1), params.BackendList) {
			//subUrl = RenderQuery(&params, ad, true)
			subUrl = RenderPVQuery(&params)
			return GetUrlScheme(&params) + "://" + params.Domain + "/onlyImpression?" + subUrl
		}
	}
	//构造 k/p/csp
	rp := ""
	rcsp := strings.Replace(CreateCsp(&params, r.DspExt), "=", "%3D", -1)
	rcsp = strings.Replace(rcsp, "+", "%2B", -1)
	rcsp = strings.Replace(rcsp, "/", "%2F", -1)
	// note 现在的请求都是打给pionner，所以这段逻辑可以优化。。    只要是AppwallOrMoreOffer就要走这段
	if mvutil.IsRequestPioneerDirectly(&params) { // 必走
		// appwall使用pioneer返回的rs参数
		if len(r.AsResp.GetRs()) > 0 {
			params.QueryRs = UrlReplace1(r.AsResp.GetRs()) // 取pioneer的返回结果
		}
		if len(r.AsResp.GetOz()) > 0 {
			rp = UrlReplace1(r.AsResp.GetOz())
		}
	}
	// note 对于这段逻辑，可以去线上看日志，观察IsNewUrl为true的比例，如果高达999,则可以直接优化这段逻辑
	// 只要apiVersion>1.4,固定为true。  不是AppwallOrMoreOffer的，apiVersion<1.4
	if params.IsNewUrl {
		dspExt, _ := r.GetDspExt()
		if dspExt != nil && dspExt.DspId == mvconst.MAS {
			//urlParams := r.AsResp.GetUrlParam()
			//if len(urlParams) > 0 {
			//	rp = UrlReplace1(urlParams[0].GetOz()) // 干掉
			//}
			// mas oz 参数移动到AsResp内
			if len(rp) == 0 {
				// 现在是返回所有。。。       支持abtest，标记字段。上线顺序
				rp = UrlReplace1(r.AsResp.GetOz()) // pi胜出
			}
		} else {
			rp = UrlReplace1(mvutil.SerializeOZ(&params)) // 三方dsp胜出
		}
		subUrl = "{sh}://{do}/onlyImpression?k=" + params.Extra5 + "&p=" + rp + "&csp=" + rcsp + "&c={c}" + "&rs=" + params.QueryRs
		// cn tracking 切量
		if CnTrackingDomainTag(r) {
			subUrl = strings.Replace(subUrl, "{do}", extractor.GetTrackingCNABTestConf().Domain, -1)
		}
	} else {
		if !mvutil.IsRequestPioneerDirectly(&params) {
			// note 目前p参数的变量都在SerializeOImpP这个方法中处理，所以要对这个方法中的参数进行划分，之后再在AsResp这个参数中新增字段来接收划分出来的参数
			rp = strings.Replace(mvutil.SerializeOImpP(&params), "=", "%3D", -1)
		}
		subUrl = "k=" + params.Extra5 + "&p=" + rp + "&csp=" + rcsp + "&rs=" + params.QueryRs

		// cn tracking 切量		 note 保留
		if CnTrackingDomainTag(r) {
			params.Domain = extractor.GetTrackingCNABTestConf().Domain
		}
		subUrl = GetUrlScheme(&params) + "://" + params.Domain + "/onlyImpression?" + subUrl
	}

	// token_r值为1则在展示点击url中下发{encrypt_p},{irlfa}的宏
	if r.Param.TokenRule == 1 && r.IsHBRequest {
		subUrl += "&encrypt_p={encrypt_p}&irlfa={irlfa}"
	}

	return subUrl
}

func CnTrackingDomainTag(r *mvutil.RequestParams) bool {
	return r.Param.CNTrackingDomainTag && len(extractor.GetTrackingCNABTestConf().Domain) > 0
}

func RenderEndscreenUrlWithInfo(mr *MobvistaResult, r *mvutil.RequestParams, params *mvutil.Params) {
	renderEndscreenUrl(mr, r, params)
	mr.Data.HTMLURL = appendEndscreen(mr.Data.HTMLURL, params)
	mr.Data.EndScreenURL = appendEndscreen(mr.Data.EndScreenURL, params)
}

func appendEndscreen(url string, params *mvutil.Params) string {
	if len(url) <= 0 {
		return ""
	}
	flag := "?"
	if strings.Contains(url, "?") {
		flag = "&"
	}
	return url + flag + "unit_id=" + strconv.FormatInt(params.UnitID, 10) + "&sdk_version=" + params.SDKVersion
}

func renderEndscreenUrl(mr *MobvistaResult, r *mvutil.RequestParams, params *mvutil.Params) {
	adType := params.AdType
	if !mvutil.InArray(int(adType), []int{mvconst.ADTypeOfferWall, mvconst.ADTypeInterstitialSdk, mvconst.ADTypeRewardVideo, mvconst.ADTypeInterstitialVideo}) {
		return
	}
	confs, ifFind := extractor.Getofferwall_urls()
	if !ifFind {
		return
	}
	var conf mvutil.COfferwallUrls_
	if NeedSchemeHttps(params.HTTPReq) {
		conf = confs.HTTPS
	} else {
		conf = confs.HTTP
	}
	if len(conf.Rewardvideo_end_screen) <= 0 {
		return
	}
	if adType == mvconst.ADTypeInterstitialSdk {
		mr.Data.HTMLURL = conf.Interstitial_sdk
		return
	}
	if adType == mvconst.ADTypeOfferWall {
		mr.Data.HTMLURL = conf.Offerwall
		mr.Data.EndScreenURL = conf.End_screen
		if isOfferwallChange(r) {
			mr.Data.HTMLURL = conf.Interstitial_sdk
		}
		return
	}
	if adType == mvconst.ADTypeRewardVideo || adType == mvconst.ADTypeInterstitialVideo {
		mr.Data.EndScreenURL = conf.Rewardvideo_end_screen
		if len(params.EndcardUrl) > 0 {
			mr.Data.EndScreenURL = GetUrlScheme(params) + "://" + params.EndcardUrl
		}
	}
}

func RenderUnitEndcard(r *mvutil.RequestParams) {
	if r.UnitInfo.Endcard != nil {
		endcard := RandEndcard(r.UnitInfo.Endcard, r)
		if len(endcard.Url) > 0 {
			r.Param.EndcardUrl = endcard.Url
			r.Param.Extendcard = strconv.FormatInt(int64(endcard.ID), 10)
			return
		}
	}
	conf, _ := extractor.GetEndcard()
	endcard := RandEndcard(&conf, r)
	if len(endcard.Url) > 0 {
		r.Param.EndcardUrl = endcard.Url
		r.Param.Extendcard = strconv.FormatInt(int64(endcard.ID), 10)
	}
}

func RenderUnitEndcardNew(r *mvutil.RequestParams, res *corsair_proto.BackendAds) {
	if res == nil {
		return
	}
	if res.EndScreenTemplateId != nil {
		templateMapConf, ifFind := extractor.GetTEMPLATE_MAP()
		if ifFind {
			endScreenId := int(*(res.EndScreenTemplateId))
			endScreenIdStr := strconv.Itoa(endScreenId)
			if endScreenUrl, ok := templateMapConf.EndScreen[endScreenIdStr]; ok && len(endScreenUrl) > 0 {
				r.Param.EndcardUrl = endScreenUrl
				r.Param.Extendcard = strconv.Itoa(endScreenId)
				return
			}
		}
	} else if r.UnitInfo.Endcard != nil {
		// 对于三方独占情况，仍需先查询unit维度配置
		endcard := RandEndcard(r.UnitInfo.Endcard, r)
		if len(endcard.Url) > 0 {
			r.Param.EndcardUrl = endcard.Url
			r.Param.Extendcard = strconv.FormatInt(int64(endcard.ID), 10)
			return
		}
	}

	conf, _ := extractor.GetEndcard()
	endcard := RandEndcard(&conf, r)
	if len(endcard.Url) > 0 {
		r.Param.EndcardUrl = endcard.Url
		r.Param.Extendcard = strconv.FormatInt(int64(endcard.ID), 10)
	}
	return
}

func isOfferwallChange(r *mvutil.RequestParams) bool {
	landingPageVersion := r.AppInfo.LandingPageVersion
	if len(landingPageVersion) <= 0 {
		return false
	}
	return mvutil.InStrArray(r.Param.AppVersionName, landingPageVersion)
}

// func IsGoTrack(r *mvutil.RequestParams) {
// 	r.Param.Extabtest1 = 10
// 	if r.Param.Gotk {
// 		r.Param.Extabtest1 = 10
// 		return
// 	}
// 	confs, _ := extractor.GetGO_TRACK()
// 	// 判断path
// 	pathConf, ok := confs["path"]
// 	if !ok {
// 		return
// 	}
// 	_, ok = pathConf[r.Param.RequestPath]
// 	if !ok {
// 		return
// 	}

// 	//rand := mvutil.GetRandConsiderZero(r.Param.GAID, r.Param.IDFA, mvconst.SALT_GOTRACK, 10000)
// 	rateRand := rand.Intn(10000)

// 	// unit维度
// 	rate := 0
// 	area := extractor.GetSYSTEM_AREA()
// 	conf, ok := confs[area]
// 	if !ok {
// 		conf, ok = confs["ALL"]
// 	}

// 	if ok {
// 		unitStr := strconv.FormatInt(r.Param.UnitID, 10)
// 		rate, ok = conf[unitStr]
// 		if !ok {
// 			rate, ok = conf["ALL"]
// 		}
// 		if ok {
// 			// if rateRand == -1 {
// 			// 	return
// 			// }
// 			//rand = 9999 - rand
// 			if rateRand < rate {
// 				r.Param.Extabtest1 = 10
// 			} else {
// 				if rateRand < 2*rate {
// 					r.Param.Extabtest1 = 8
// 				} else {
// 					r.Param.Extabtest1 = 9
// 				}
// 			}
// 			return
// 		}
// 	}

// 	// 按照adtype切量
// 	adtypeConf, ok := confs["adtype"]
// 	if ok {
// 		adtype := strconv.FormatInt(int64(r.Param.AdType), 10)
// 		rate, ok := adtypeConf[adtype]
// 		if ok {
// 			// if rateRand == -1 {
// 			// 	return
// 			// }
// 			//rand = 9999 - rand
// 			if rateRand < rate {
// 				r.Param.Extabtest1 = 10
// 			} else {
// 				if rateRand < 2*rate {
// 					r.Param.Extabtest1 = 8
// 				} else {
// 					r.Param.Extabtest1 = 9
// 				}
// 			}
// 		}
// 		return
// 	}
// }

func UrlReplace1(str string) string {
	return strings.Replace(str, "=", "%3D", -1)
}

func UrlReplace3(str string) string {
	str = strings.Replace(str, "=", "%3D", -1)
	str = strings.Replace(str, "+", "%2B", -1)
	str = strings.Replace(str, "/", "%2F", -1)
	return str
}

// UrlDeReplace3 上面 UrlReplace3 的逆过程
func UrlDeReplace3(str string) string {
	str = strings.Replace(str, "%2F", "/", -1)
	str = strings.Replace(str, "%2B", "+", -1)
	str = strings.Replace(str, "%3D", "=", -1)
	return str
}

func renderSmallRoutineUrl(trackUrl string, params *mvutil.Params, campaignId int64, urlAppend string) string {
	urlData, err := url.Parse(trackUrl)
	if err != nil {
		mvutil.Logger.Runtime.Warnf("smallRoutine tracking url error")
		return ""
	}
	query := urlData.Query()
	pathData := query.Get("path")
	if len(pathData) > 0 {
		pathData = pathData + urlAppend
		pathData = strings.Replace(pathData, "{pubId}", strconv.FormatInt(params.PublisherID, 10), -1)
		pathData = strings.Replace(pathData, "{appId}", strconv.FormatInt(params.AppID, 10), -1)
		pathData = strings.Replace(pathData, "{unitId}", strconv.FormatInt(params.UnitID, 10), -1)
		pathData = strings.Replace(pathData, "{ip}", params.ClientIP, -1)
		pathData = strings.Replace(pathData, "{offerId}", strconv.FormatInt(campaignId, 10), -1)
		pathData = strings.Replace(pathData, "{gaid}", params.GAID, -1)
		pathData = strings.Replace(pathData, "{idfa}", params.IDFA, -1)
	}
	query.Set("path", pathData)
	urlData.RawQuery = query.Encode()
	trackUrl = urlData.String()
	return trackUrl
}

func renderInteractiveAdTracking(ad *Ad, subUrl string) {
	ad.AdTracking.Dropout = getAdTrackingUrlList("dropout_track", subUrl, false, 0)
	ad.AdTracking.Plycmpt = getAdTrackingUrlList("plycmpt_track", subUrl, false, 0)
}

func RenderJssdkCDN(ad *Ad, params mvutil.Params) {
	// 判断是否为jssdk 新域名逻辑
	if mvutil.NeedNewJssdkDomain(params.RequestPath, params.Ndm) {
		if len(params.JssdkCdnDomain) > 0 {
			// 替换协议头及域名
			replaceCdnDomain(ad, params.JssdkCdnDomain, params.HTTPReq)
		}
	}
}

func replaceCdnDomain(ad *Ad, newDomain string, httpReq int32) {
	if len(ad.IconURL) > 0 {
		ad.IconURL = replaceSchemeAndHost(ad.IconURL, newDomain, httpReq)
	}
	if len(ad.ImageURL) > 0 {
		ad.ImageURL = replaceSchemeAndHost(ad.ImageURL, newDomain, httpReq)
	}
	if len(ad.VideoURL) > 0 {
		videoUrl := mvutil.Base64Decode(ad.VideoURL)
		if len(videoUrl) > 0 {
			videoUrl = replaceSchemeAndHost(videoUrl, newDomain, httpReq)
			ad.VideoURL = mvutil.Base64Encode(videoUrl)
		}
	}
	if len(ad.ExtImg) > 0 {
		ad.ExtImg = replaceSchemeAndHost(ad.ExtImg, newDomain, httpReq)
	}
	if len(ad.GifURL) > 0 {
		ad.GifURL = replaceSchemeAndHost(ad.GifURL, newDomain, httpReq)
	}

	if ad.Rv.Image != nil && len(ad.Rv.Image.IdcdImg) > 0 {
		for k, v := range ad.Rv.Image.IdcdImg {
			ad.Rv.Image.IdcdImg[k] = replaceSchemeAndHost(v, newDomain, httpReq)
		}
	}
}

func replaceSchemeAndHost(oldUrl, newDomain string, httpReq int32) string {
	u, err := url.Parse(oldUrl)
	if err == nil {
		if httpReq == 2 {
			u.Scheme = "https"
		}
		u.Host = newDomain
		return u.String()
	}
	return oldUrl
}

// 读取配置，根据权重获取最终替换的cdn地址，并且记录cdn标记
//func getNewCdnUrl(cdnSetting []*smodel.CdnSetting, params *mvutil.Params, oldHost string) string {
//	newCDNUrl := ""
//	//var cdnMap map[string]int
//	cdnMap := make(map[string]int)
//	cdnId := make(map[string]int)
//	for _, cdnSetting := range cdnSetting {
//		cdnMap[cdnSetting.Url] = cdnSetting.Weight
//		cdnId[cdnSetting.Url] = cdnSetting.Id
//	}
//	adnConf, _ := extractor.GetADNET_SWITCHS()
//	if cdnRandByReq, ok := adnConf["cdnRandByReq"]; ok && cdnRandByReq == 1 {
//		newCDNUrl = mvutil.RandByRate2(cdnMap)
//	} else {
//		newCDNUrl = mvutil.RandByDeviceRate(cdnMap, params)
//	}
//	// 设备id为空值，则指定默认的cdn地址
//	if len(oldHost) > 0 && len(newCDNUrl) == 0 {
//		adnetDefaultConf := extractor.GetADNET_DEFAULT_VALUE()
//		if defaultCdn, ok := adnetDefaultConf[oldHost]; ok && len(defaultCdn) > 0 {
//			return defaultCdn
//		}
//	}
//	if id, ok := cdnId[newCDNUrl]; ok {
//		params.ExtCDNAbTest = id
//	}
//	return newCDNUrl
//}

// 读取配置，根据权重获取最终替换的cdn地址，并且记录cdn标记
func getNewCdnUrl(cdnSetting []*smodel.CdnSetting, params *mvutil.Params, oldHost string) string {
	newCDNUrlId := 0
	cdnWeigt := make(map[int]int)
	cdnUrl := make(map[int]string)
	for _, cdnSetting := range cdnSetting {
		cdnWeigt[cdnSetting.Id] = cdnSetting.Weight
		cdnUrl[cdnSetting.Id] = cdnSetting.Url
	}
	adnConf, _ := extractor.GetADNET_SWITCHS()
	if cdnRandByReq, ok := adnConf["cdnRandByReq"]; ok && cdnRandByReq == 1 {
		newCDNUrlId = mvutil.RandByRate3WithIntMap(cdnWeigt)
	} else {
		newCDNUrlId = mvutil.RandByDeviceRateWithIntMap(cdnWeigt, params)
	}
	// 设备id为空值，则指定默认的cdn地址
	if len(oldHost) > 0 && newCDNUrlId == 0 {
		adnetDefaultConf := extractor.GetADNET_DEFAULT_VALUE()
		if defaultCdn, ok := adnetDefaultConf[oldHost]; ok && len(defaultCdn) > 0 {
			return defaultCdn
		}
	}
	// 记录cdn测试标记
	if cdnUrl, ok := cdnUrl[newCDNUrlId]; ok {
		params.ExtCDNAbTest = newCDNUrlId
		return cdnUrl
	}
	return ""
}

func createChetUrl(params *mvutil.Params) string {
	res := ""
	queryStr := "k=" + params.RequestID
	domain := params.Domain
	// 支持把埋点url切换到其他tracking服务上
	chetDomain := extractor.GetCHET_DOMAIN()
	if len(chetDomain) > 0 {
		domain = chetDomain
	}
	queryStr += "&chet=1"
	// 加入unit以区分不同的埋点url
	queryStr += "&uid=" + strconv.FormatInt(params.UnitID, 10)
	res = GetUrlScheme(params) + "://" + domain + "/chet?" + queryStr
	return res
}

func RenderEndscreenProperty(r *mvutil.RequestParams) {
	if len(r.Param.EndcardUrl) > 0 {
		flag := "?"
		if strings.Contains(r.Param.EndcardUrl, "?") {
			flag = "&"
		}
		// 无论是固定unit还是拆分unit，都下发ec值及offer信息
		r.Param.EndcardUrl += flag + "ec_id=" + r.Param.Extendcard
		if len(r.Param.NewMoreOfferParams) > 0 {
			r.Param.EndcardUrl += r.Param.NewMoreOfferParams
		}
		if r.Param.NewMoreOfferFlag {
			// 对于返回online api单子的情况，也要下发mof_uid
			r.Param.EndcardUrl += "&mof_uid=" + strconv.FormatInt(r.UnitInfo.Unit.MofUnitId, 10)

			if r.Param.MofAbFlag {
				r.Param.EndcardUrl += "&mof_ab=1"
			}
		}
		if r.Param.NewMofImpFlag {
			r.Param.EndcardUrl += "&n_imp=1"
		}
		if len(r.Param.ExtDataInit.CloseAdTag) > 0 {
			r.Param.EndcardUrl += "&clsad=" + r.Param.ExtDataInit.CloseAdTag
		}
	}
}

func AddCtnSizeToEndscreen(r *mvutil.RequestParams) {
	if len(r.Param.EndcardUrl) > 0 {
		flag := "?"
		if strings.Contains(r.Param.EndcardUrl, "?") {
			flag = "&"
		}
		r.Param.EndcardUrl += flag + "ctnsize=" + r.Param.CtnSizeTag
	}
}

func renderAdClickAndClickUrl(ad *Ad, params *mvutil.Params) {
	// h5跳转逻辑
	if params.ExtDataInit.H5Handle == 1 {
		ad.ClickURL = mvutil.Base64Encode(ad.ClickURL)
		return
	}
	// sdk跳转逻辑
	if params.ExtDataInit.ClickMode6NotInGpAndAppstore != 1 && params.AdClickReplaceclickByClickTag {
		// 对于14 deeplink单子，不需要把click_url放到adtracking.click中，因为click_url为落地页，放到adtracking.click做点击上报没有意义。
		if params.Extra10 == mvconst.JUMP_TYPE_TRACKING_PING_THIRDPARTY_CLICK_URL {
			// 兜底的click_url不添加redirect=1，避免兜底情况无法跳转到落地页
			ad.ClickURL = mvutil.DelSubStr(ad.PreviewUrl, mvconst.REDIRECT)
		} else {
			// 给adclick赋值
			ad.AdTracking.Click = append(ad.AdTracking.Click, ad.ClickURL)
			ad.AdTrackingPoint = &ad.AdTracking
			// 清空click_url避免重复上报
			ad.ClickURL = ad.PreviewUrl
		}
	}
}

func ReplaceImpt2withImp(params *mvutil.Params) bool {
	if !mvutil.IsHbOrV3OrV5Request(params.RequestPath) {
		return false
	}
	if params.Platform != mvconst.PlatformIOS {
		return false
	}
	if params.FormatSDKVersion.SDKVersionCode == 50401 || params.FormatSDKVersion.SDKVersionCode == 50402 ||
		params.FormatSDKVersion.SDKVersionCode == 50500 || params.FormatSDKVersion.SDKVersionCode == 50501 {
		return true
	}
	return false
}

func renderThirdpartyWhiteListClickUrl(ad *Ad, params *mvutil.Params, campaign *smodel.CampaignInfo, r *mvutil.RequestParams) {
	if params.PingMode != 1 {
		return
	}

	if params.ExtDataInit.ClickInServer != 1 && params.ExtDataInit.ClickInServer != 3 && params.Extra10 != mvconst.JUMP_TYPE_TRACKING_PING_THIRDPARTY_CLICK_URL {
		return
	}

	if campaign == nil {
		return
	}

	if params.RequestPath == mvconst.PATHJssdkApi ||
		mvutil.IsHbOrV3OrV5Request(params.RequestPath) || mvutil.IsMP(params.RequestPath) {
		ad.ClickURL = getMarketUrlByLinkType(params, campaign, ad, r)
		if len(ad.NoticeURL) > 0 {
			ad.NoticeURL = ad.NoticeURL + mvconst.REDIRECT
		}
	}
}

func hitChetLinkConfig(params *mvutil.Params) (string, bool) {
	if configs, ok := extractor.GetChetLinkConfigs(); ok {
		for _, cfg := range configs {
			if cfg != nil && cfg.Condition != nil && checkChetHit(cfg.Condition, params) {
				return cfg.Link, true
			}
		}
	}
	return "", false
}

func replaceChetLink(link string, params *mvutil.Params) string {
	if params == nil {
		return link
	}

	link = strings.Replace(link, "__APPID__", strconv.FormatInt(params.AppID, 10), -1)
	link = strings.Replace(link, "__UNITID__", strconv.FormatInt(params.UnitID, 10), -1)
	link = strings.Replace(link, "__SDKVER__", url.QueryEscape(params.SDKVersion), -1)
	link = strings.Replace(link, "__CC__", url.QueryEscape(params.CountryCode), -1)
	link = strings.Replace(link, "__REQID__", url.QueryEscape(params.RequestID), -1)
	return link
}

func checkChetHit(condition *mvutil.ChetLinkCondition, params *mvutil.Params) bool {
	if condition == nil && params == nil {
		return false
	}

	// if params.Platform == mvconst.PlatformIOS && len(condition.IOSSDKVersion) > 0 && !mvutil.InInt32Arr(params.FormatSDKVersion.SDKVersionCode, condition.IOSSDKVersion) {
	// 	return false
	// }

	// if params.Platform == mvconst.PlatformAndroid && len(condition.AndroidSDKVersion) > 0 && !mvutil.InInt32Arr(params.FormatSDKVersion.SDKVersionCode, condition.AndroidSDKVersion) {
	// 	return false
	// }

	sdkVersion := condition.IOSSDKVersions
	if params.Platform == mvconst.PlatformAndroid {
		sdkVersion = condition.AndroidSDKVersions
	}
	if len(sdkVersion) > 0 && !mvutil.InInt32Arr(params.FormatSDKVersion.SDKVersionCode, sdkVersion) {
		return false
	}

	if len(condition.AppIds) > 0 && !mvutil.InInt64Arr(params.AppID, condition.AppIds) {
		return false
	}

	if len(condition.UnitIds) > 0 && !mvutil.InInt64Arr(params.UnitID, condition.UnitIds) {
		return false
	}

	if len(condition.CountryCodes) > 0 && !mvutil.InStrArray(params.CountryCode, condition.CountryCodes) {
		return false
	}

	return true
}

func replaceUrls(urls []string, params *mvutil.Params, sh, do string) []string {
	if len(sh) == 0 || len(do) == 0 {
		sh, do = getShellAndDomain(params)
	}
	for key, str := range urls {
		str = replaceUrlMacro(str, params, sh, do)
		urls[key] = str
	}
	return urls
}

func replaceUrlMacro(str string, params *mvutil.Params, sh, do string) string {
	str = strings.Replace(str, "{q}", params.QueryQ, -1)
	str = strings.Replace(str, "{k}", params.RequestID, -1)
	str = strings.Replace(str, "{p}", params.QueryP, -1)
	str = strings.Replace(str, "{q}", params.QueryQ, -1)
	str = strings.Replace(str, "{r}", params.QueryR, -1)
	str = strings.Replace(str, "{al}", params.QueryAL, -1)
	str = strings.Replace(str, "{csp}", params.QueryCSP, -1)
	str = strings.Replace(str, "{do}", do, -1)
	str = strings.Replace(str, "{sh}", sh, -1)
	return str
}

func replaceReqUrlMacro(str string, rks map[string]string) string {
	for key, val := range rks {
		str = strings.Replace(str, "{"+key+"}", val, -1)
	}
	return str
}

func getShellAndDomain(params *mvutil.Params) (string, string) {
	sh := GetUrlScheme(params)
	do := params.Domain

	return sh, do
}

func replaceOnlineApiUrl(ad *Ad, params *mvutil.Params) {
	if params.RequestType != mvconst.REQUEST_TYPE_OPENAPI_AD {
		return
	}
	sh, do := getShellAndDomain(params)
	if len(ad.ClickURL) > 0 {
		ad.ClickURL = replaceUrls([]string{ad.ClickURL}, params, sh, do)[0]
	}
	if len(ad.NoticeURL) > 0 {
		ad.NoticeURL = replaceUrls([]string{ad.NoticeURL}, params, sh, do)[0]
	}
	if len(ad.PreviewUrl) > 0 {
		ad.PreviewUrl = replaceUrls([]string{ad.PreviewUrl}, params, sh, do)[0]
	}
	if len(ad.ImpressionURL) > 0 {
		ad.ImpressionURL = replaceUrls([]string{ad.ImpressionURL}, params, sh, do)[0]
	}
	if len(ad.DeepLink) > 0 {
		ad.DeepLink = replaceUrls([]string{ad.DeepLink}, params, sh, do)[0]
	}
}

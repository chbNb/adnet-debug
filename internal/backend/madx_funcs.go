package backend

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"math/rand"
	"strconv"
	"strings"

	jsoniter "github.com/json-iterator/go"
	supply_mvutil "gitlab.mobvista.com/ADN/chasm/module/supply/mvutil"

	"gitlab.mobvista.com/ADN/adnet/internal/corsair_proto"
	"gitlab.mobvista.com/ADN/adnet/internal/extractor"
	"gitlab.mobvista.com/ADN/adnet/internal/mvconst"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	"gitlab.mobvista.com/ADN/adnet/internal/redis"
	"gitlab.mobvista.com/ADN/adnet/internal/watcher"
	adxconst "gitlab.mobvista.com/ADN/adx_common/constant"
	rtb "gitlab.mobvista.com/ADN/mtg_openrtb/pkg/mtgrtb"
	"gitlab.mobvista.com/ADN/mtg_openrtb/pkg/native"
	"gitlab.mobvista.com/ADN/mtg_openrtb/pkg/vast"
	"gitlab.mobvista.com/adserver/recommend_protocols/go/ad_server"
)

// FilterAdTracking adtracking特性的过滤和一些sdk platform/version方面过滤
func FilterAdTracking(adType int32, platform int, adTracking int, sdkVerion supply_mvutil.SDKVersionItem) bool {
	return filterAdTracking(adType, platform, adTracking, sdkVerion)
}

func decodeVast(adm string) (*vast.VAST, error) {
	var vastdata vast.VAST
	err := xml.Unmarshal([]byte(adm), &vastdata)
	if err != nil {
		// fmt.Println("decodeVast error: ", err)
		watcher.AddWatchValue("mAdx_amd_vast_unmarshal_error", float64(1))
		return nil, err
	}
	return &vastdata, nil
}

func decodeNative(adm string) (*native.Native, error) {
	var nativeObject native.Native
	err := json.Unmarshal([]byte(adm), &nativeObject)
	if err != nil {
		watcher.AddWatchValue("mAdx_amd_native_unmarshal_error", float64(1))
		return nil, err
	}
	return &nativeObject, nil
}

func constructBapp(BlackBundleList []string, reqCtx *mvutil.ReqCtx) []string {
	if reqCtx.ReqParams.UnitInfo.BlackPackageList != nil {
		BlackBundleList = append(BlackBundleList, *reqCtx.ReqParams.UnitInfo.BlackPackageList...)
		return BlackBundleList
	}
	if reqCtx.ReqParams.AppInfo.BlackPackageList != nil {
		BlackBundleList = append(BlackBundleList, *reqCtx.ReqParams.AppInfo.BlackPackageList...)
		return BlackBundleList
	}
	return BlackBundleList
}

func constructImpExt(reqCtx *mvutil.ReqCtx, nativeType int64) *rtb.BidRequest_Imp_Ext {
	ext := &rtb.BidRequest_Imp_Ext{}
	ext.BlackOfferList, ext.WhiteOfferList = constructBlackAndWhiteOfferList(reqCtx)
	ext.BlackCategory = constructBlackCategory(reqCtx)
	var forbidApk bool
	if !mvutil.InArray(mvconst.APK, reqCtx.ReqParams.AppInfo.App.OfferPreference) {
		forbidApk = true
	}
	ext.ForbidApk = &forbidApk
	ext.NativeType = &nativeType
	isCcpa := strings.ToLower(reqCtx.ReqParams.Param.RegionString) == "ca" && reqCtx.ReqParams.AppInfo.App.Ccpa == 1
	ext.IsCcpa = &isCcpa
	// 传递adnet的灰度标记
	if len(reqCtx.ReqParams.Param.StartMode) > 0 {
		ext.GrayTags = map[string]string{
			mvconst.AdnetStartModeTag: reqCtx.ReqParams.Param.StartMode,
		}
	}
	// SDK版本 是否支持deeplink
	param := reqCtx.ReqParams.Param
	if (param.Platform == mvconst.PlatformIOS && param.FormatSDKVersion.SDKVersionCode >= 40700) ||
		(param.Platform == mvconst.PlatformAndroid && param.FormatSDKVersion.SDKVersionCode >= 90300) {
		isDeeplink := true
		ext.Isdeeplink = &isDeeplink
	}
	//支持deeplinkType
	if isDeeplink, deeplinkType := onlineApiSupportDeepLink(reqCtx); isDeeplink == true {
		ext.Isdeeplink = &isDeeplink
		ext.DeeplinkType = &deeplinkType
	}
	// 支持 htmlSupport
	if reqCtx.ReqParams.Param.RequestType == mvconst.REQUEST_TYPE_OPENAPI_AD &&
		reqCtx.ReqParams.Param.AdType == mvconst.ADTypeBanner { //只有nonlineAPI & banner 的量才会不支持
		htmlNonsupport := param.HtmlSupport != 1
		ext.HtmlNonsupport = &htmlNonsupport
	}

	//skadnetwork
	if reqCtx.ReqParams.Param.Skadnetwork != nil && reqCtx.ReqParams.AppInfo.App.Platform == mvconst.PlatformIOS &&
		len(reqCtx.ReqParams.AppInfo.RealPackageName) > 0 && reqCtx.ReqParams.AppInfo.RealPackageName[0:2] == "id" {
		sourceapp := reqCtx.ReqParams.AppInfo.RealPackageName[2:]
		ext.Skadn = &rtb.BidRequest_Imp_Ext_Skadn{
			Version:    &reqCtx.ReqParams.Param.Skadnetwork.Ver,
			Sourceapp:  &sourceapp,
			Skadnetids: reqCtx.ReqParams.Param.Skadnetwork.Adnetids,
		}
	}

	if mvutil.IsNewIv(&reqCtx.ReqParams.Param) {
		ext.MaterialType = &reqCtx.ReqParams.UnitInfo.Unit.MaterialType
	}

	return ext
}

func onlineApiSupportDeepLink(reqCtx *mvutil.ReqCtx) (bool, int32) {
	Param := reqCtx.ReqParams.Param
	//接口上报优先
	if Param.RequestType == mvconst.REQUEST_TYPE_OPENAPI_AD {

		if Param.DeeplinkType > 0 && Param.DeeplinkType < 3 {
			return true, int32(Param.DeeplinkType)
		}

		//如果是OnlineAPI的流量 & 配置为1的才会返回true
		conf := extractor.GetONLINE_API_SUPPORT_DEEPLINK_V2()
		if subConf, ok := conf["unit"]; ok && Param.UnitID > 0 {
			if deeplinkType, ok := subConf[Param.UnitID]; ok && deeplinkType > 0 && deeplinkType < 3 {
				return true, deeplinkType
			}
		}
		if subConf, ok := conf["app"]; ok && Param.AppID > 0 {
			if deeplinkType, ok := subConf[Param.AppID]; ok && deeplinkType > 0 && deeplinkType < 3 {
				return true, deeplinkType
			}
		}
		if subConf, ok := conf["publisher"]; ok && Param.PublisherID > 0 {
			if deeplinkType, ok := subConf[Param.PublisherID]; ok && deeplinkType > 0 && deeplinkType < 3 {
				return true, deeplinkType
			}
		}
	}
	return false, 0
}

// getNativeType  根据 content_type native_info is_video 判断native_type
func getNativeType(reqCtx *mvutil.ReqCtx) int64 {
	adType := reqCtx.ReqParams.Param.AdType
	if adType != mvconst.ADTypeNative {
		return adxconst.NativeTypeNot
	}
	param := reqCtx.ReqParams.Param
	contentType := param.VideoAdType
	nativeInfoStr := param.NativeInfo
	nativeInfo := make([]NativeInfo, 0)
	// nativeInfo 字段为空， 默认为大图
	if len(nativeInfoStr) == 0 {
		nativeInfo = append(nativeInfo, NativeInfo{Id: 2})
	} else {
		_ = jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal([]byte(nativeInfoStr), &nativeInfo)
		if len(nativeInfo) == 0 {
			nativeInfo = append(nativeInfo, NativeInfo{Id: 2})
		}
	}
	isVideo := reqCtx.ReqParams.Param.VideoVersion != ""
	if contentType == mvconst.VideoAdTypeNOLimit && nativeInfo[0].Id == 2 && isVideo {
		return adxconst.NativeTypeDisplayVideo
	}
	if contentType == mvconst.VideoAdTypeOnlyVideo && nativeInfo[0].Id == 2 && isVideo {
		return adxconst.NativeTypeVideo
	}
	return adxconst.NativeTypeDisplay
}

func constructBlackAndWhiteOfferList(reqCtx *mvutil.ReqCtx) ([]int64, []int64) {
	var blackOfferList []int64
	var whiteOfferList []int64
	blackOfferList = mvutil.RenderStringToIntList(reqCtx.ReqParams.Param.ExcludeIDS)
	// 先判读unit维度
	if reqCtx.ReqParams.UnitInfo.DirectOfferConfig.Status == 1 && len(reqCtx.ReqParams.UnitInfo.DirectOfferConfig.OfferIds) > 0 {
		if reqCtx.ReqParams.UnitInfo.DirectOfferConfig.Type == 1 {
			whiteOfferList = reqCtx.ReqParams.UnitInfo.DirectOfferConfig.OfferIds
		} else if reqCtx.ReqParams.UnitInfo.DirectOfferConfig.Type == 2 {
			blackOfferList = append(blackOfferList, reqCtx.ReqParams.UnitInfo.DirectOfferConfig.OfferIds...)
		}
		return blackOfferList, whiteOfferList
	}
	// unit维度没配置，则判断app维度
	if reqCtx.ReqParams.AppInfo.DirectOfferConfig.Status == 1 && len(reqCtx.ReqParams.AppInfo.DirectOfferConfig.OfferIds) > 0 {
		if reqCtx.ReqParams.AppInfo.DirectOfferConfig.Type == 1 {
			whiteOfferList = reqCtx.ReqParams.AppInfo.DirectOfferConfig.OfferIds
		} else if reqCtx.ReqParams.AppInfo.DirectOfferConfig.Type == 2 {
			blackOfferList = append(blackOfferList, reqCtx.ReqParams.AppInfo.DirectOfferConfig.OfferIds...)
		}
		return blackOfferList, whiteOfferList
	}
	return blackOfferList, whiteOfferList
}

func constructBlackCategory(reqCtx *mvutil.ReqCtx) []int64 {
	var blackCategoryList []int64
	// 先读取unit维度配置
	if reqCtx.ReqParams.UnitInfo.BlackCategoryList != nil {
		blackCategoryList = *reqCtx.ReqParams.UnitInfo.BlackCategoryList
		return blackCategoryList
	}
	// 若unit维度无配置，则使用app维度配置
	if reqCtx.ReqParams.AppInfo.BlackCategoryList != nil {
		blackCategoryList = *reqCtx.ReqParams.AppInfo.BlackCategoryList
		return blackCategoryList
	}
	return blackCategoryList
}

func getCriteoImg(oldImgUrl string) (imgInfo *CriteoImg, imageKey string, err error) {
	imageKey = mvutil.Md5(oldImgUrl)
	// 查询redis
	imgInfoStr, err := redis.LocalRedisGet(imageKey)
	if err != nil {
		return nil, imageKey, err
	}
	imgInfo = new(CriteoImg)
	err = jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal([]byte(imgInfoStr), imgInfo)
	if err != nil {
		return nil, imageKey, err
	}
	return
}

type CriteoImg struct {
	Url        string `json:"url,omitempty"`
	Size       int    `json:"size,omitempty"`
	Resolution string `json:"resolution,omitempty"`
	Result     int    `json:"result,omitempty"`
}

// oldVastFill 旧格式下的对Extension的处理
func oldVastFill(ad *corsair_proto.Campaign, extension *vast.Extension) error {
	template := extension.Template
	for _, asset := range template.Asset {
		// mvutil.Logger.Runtime.Infof("asset: %+v", asset)
		if asset.AssetType == "icon" {
			if asset.Cdata != "" {
				*ad.IconURL = strings.TrimSpace(asset.Cdata)
			} else {
				*ad.IconURL = strings.TrimSpace(asset.Value)
			}
		}
		if asset.AssetType == "upperimg" || asset.AssetType == "lowerimg" {
			*ad.ImageURL = strings.TrimSpace(asset.Cdata)
			*ad.ImageSize = "VIDEO"
		}
		if asset.AssetType == "CTA" {
			*ad.CtaText = asset.Value
		}
		if asset.AssetType == "title" {
			*ad.AppName = asset.Value
		}
		if asset.AssetType == "des" {
			*ad.AppDesc = asset.Value
		}
		if asset.AssetType == "starrating" {
			// todo
			*ad.Rating, _ = strconv.ParseFloat(asset.Value, 64)
		}
		if asset.AssetType == "mainimage" && *ad.ImageURL == "" {
			if asset.Cdata != "" {
				*ad.ImageURL = strings.TrimSpace(asset.Cdata)
			} else {
				*ad.ImageURL = strings.TrimSpace(asset.Value)
			}
		}
		if asset.AssetType == "numberrating" {
			// numberrating: 新加字段，评论数
			numberRating, _ := strconv.Atoi(asset.Value)
			nr := int32(numberRating)
			ad.NumberRating = &nr
		}
	}
	return nil
}

// newVastFill  新的 Vast.Extensions填充逻辑， 见 http://confluence.mobvista.com/pages/viewpage.action?pageId=24517947
func newVastFill(ad *corsair_proto.Campaign, extension *vast.Extension, HTMLResource string, param *mvutil.Params, dspId int64) error {
	// 获取视频模板
	// templateMapConf, ifFind := extractor.GetTEMPLATE_MAP()
	// if !ifFind {
	//	return errors.New("TEMPLATE_MAP not found")
	// }
	if ad.Rv == nil {
		ad.Rv = &corsair_proto.RewardVideo{}
		if extension.Orientation != nil {
			orientation, _ := strconv.Atoi(extension.Orientation.Value)
			orientation32 := int32(orientation)
			ad.Rv.Orientation = &orientation32
		} else {
			orientation := int32(param.Orientation)
			ad.Rv.Orientation = &orientation
		}
	}
	for _, template := range extension.Templates {
		if template.Id == "" && !extractor.IsVastBannerDsp(dspId) {
			continue
		}
		value := template.Cdata
		switch template.Type {
		case "video":
			ad.Rv.TemplateURL = &value
			tempId, _ := strconv.Atoi(template.Id)
			enumTempId := ad_server.VideoTemplateId(tempId) // rv.video_template
			ad.VideoTemplateId = &enumTempId
			tempId32 := int32(tempId)
			ad.Rv.Template = &tempId32
		case "endcard":
			ad.EndcardURL = &value
		case "endscreen":
			param.EndcardUrl = value
			param.Extendcard = template.Id
		case "minicard":
			ad.Rv.PausedURL = &value
		case "group":
			groupId, _ := strconv.Atoi(template.Id)
			enumTG := ad_server.TemplateGroup(groupId)
			ad.TemplateGroup = &enumTG
		case "urltpl":
			ad.UrlTemplate = &value
		case "htmltpl":
			ad.HtmlTemplate = &value
		}
	}
	for _, asset := range extension.Asset {
		// mvutil.Logger.Runtime.Infof("asset: %+v", asset)
		if asset.AssetType == "icon" {
			if asset.Cdata != "" {
				*ad.IconURL = strings.TrimSpace(asset.Cdata)
			} else {
				*ad.IconURL = strings.TrimSpace(asset.Value)
			}
		}
		if asset.AssetType == "CTA" {
			*ad.CtaText = asset.Value
		}
		if asset.AssetType == "starrating" {
			// todo
			*ad.Rating, _ = strconv.ParseFloat(asset.Value, 64)
		}
		if asset.AssetType == "numberrating" {
			// numberrating: 新加字段，评论数
			numberRating, _ := strconv.Atoi(asset.Value)
			nr := int32(numberRating)
			ad.NumberRating = &nr
		}
		if asset.AssetType == "widgets" {
			ad.GifWidgets = &asset.Value
		}
	}

	return nil
}

func isDynamicEndcard(tempId string) bool {
	return tempId == "-201" || tempId == "-202"
}
func isPlayableEndcard(tempId string) bool {
	return tempId == "-301" || tempId == "-302"
}

type VastVideoMeta struct {
	Size int32
}

type NativeInfo struct {
	Id    int `json:"id"`
	AdNum int `json:"ad_num"`
}

type ExtData2MAS struct {
	PriceFactor                       float64 `json:"pf,omitempty"`       // 频次控制- 价格系数
	PriceFactorGroupName              string  `json:"pf_g,omitempty"`     // 频次控制- 实验组名称
	PriceFactorTag                    int     `json:"pf_t,omitempty"`     // 频次控制- 实验标签，1=A, 2=B, 3=B'
	PriceFactorFreq                   *int    `json:"pf_f,omitempty"`     // 频次控制- 获取到当前的频次
	PriceFactorHit                    int     `json:"pf_h,omitempty"`     // 频次控制- 是否能命中概率， 1=命中，2=不命中
	ImpressionCap                     int     `json:"imp_c,omitempty"`    // placement的impressionCap
	IfSupDco                          int     `json:"dco,omitempty"`      // 标记dco切量结果
	V5AbtestTag                       string  `json:"v5_t,omitempty"`     // V5的实验标记， 5_5, 5_3, 或者控
	SDKOpen                           int     `json:"sdk_open,omitempty"` // SDK是否开源版本
	BandWidth                         int64   `json:"bw,omitempty"`       // 带宽
	TKSysTag                          string  `json:"tkst,omitempty"`     // tracking 集群切量
	UpTime                            string  `json:"uptime,omitempty"`   // 系统开机时间
	Brt                               string  `json:"brt,omitempty"`      // 屏幕亮度
	Vol                               string  `json:"vol,omitempty"`      // 音量
	Lpm                               string  `json:"lpm,omitempty"`      // 是否为低电量模式
	Font                              string  `json:"font,omitempty"`     // 设备默认字体大小
	ImeiABTest                        int     `json:"imei_abt,omitempty"` // imei abtest
	IsReturnWtickTag                  int     `json:"irwt,omitempty"`     // 是否给sdk返回wtick=1
	Ntbarpt                           int     `json:"ntbarpt,omitempty"`
	Ntbarpasbl                        int     `json:"ntbarpasbl,omitempty"`
	AtatType                          int     `json:"atatType,omitempty"`
	MappingIdfaTag                    string  `json:"mp_idfa,omitempty"`     // mapping idfa abtest标记
	ExcludeAopPkg                     string  `json:"exc_aop_pkg,omitempty"` // analysis offline package 不召回实验标记
	TKCNABTestTag                     int     `json:"tkcn"`                  // tracking cn 集群切量标记
	TKCNABTestAATag                   int     `json:"tkcnaa"`                // tracking cn 集群切量AA标记
	MappingIdfaCoverIdfaTag           string  `json:"mp_ici,omitempty"`      // mapping idfa 替换idfa的abtest标记
	MiskSpt                           string  `json:"misk_spt,omitempty"`    // 是否支持小米storekit。1表示支持，0表示不支持。-1表示未安装小米商店
	MoreofferAndAppwallMvToPioneerTag string  `json:"maamtp,omitempty"`      // more_offer/appwall迁移aabtest标记
	ParentUnitId                      int64   `json:"parent_id,omitempty"`
	H5Type                            int     `json:"h5_t,omitempty"`
	MofType                           int     `json:"mof_type,omitempty"` // 区分是more offer 还是 close button ad
	CrtRid                            string  `json:"crt_rid,omitempty"`  // 主offer的request_id
	CNTrackDomain                     int     `json:"cntd,omitempty"`     // 中国专线tracking域名切量标记
	OnlineApiNeedOfferBidPrice        string  `json:"olnobp,omitempty"`   // 表示算法是否需要针对hb online api请求的每个offer单独出价 1表示需要
	TrackDomainByCountryCode          int     `json:"tdbcc,omitempty"`    // 根据country code选择的tracking域名
	TmaxABTestTag                     int     `json:"a_tmax,omitempty"`   // tmax abtest tag
	UseDynamicTmax                    string  `json:"dy_tmax,omitempty"`
	Dnt                               string  `json:"dnt,omitempty"`       // dnt值为1时，表示用户退出个性化广告
	ExpIds                            string  `json:"expIds,omitempty"`    // abtest 实验id
	MpToPioneerTag                    string  `json:"mptpi,omitempty"`     // mp 流量迁移aabtest标记
	DeviceGEOCCMatch                  int     `json:"d_geo_cc,omitempty"`  // hb bid request device.geo.country 和 IP 信息的 country 是否一致, 0: device.geo.country 空, 1: 一致, 2: 不一致
	ThreeLetterCountry                string  `json:"t_cc,omitempty"`      // 三位国家码
	HBSubsidyType                     int     `json:"subsidy_t,omitempty"` // 扶持类型。普通流量：0 扶持流量+垂类保量 ：101 扶持流量+非垂类保量 ：102 冷启动流量：200
	RequestBidServer                  int32   `json:"req_bs,omitempty"`    // 是否请求bid server 1: 请求
	ParentAdType                      string  `json:"p_ad_t,omitempty"`
	ParentExchange                    string  `json:"p_exc,omitempty"`
	CdnTrackingDomain                 int     `json:"cdn_td,omitempty"`    // 使用cdn tracking domain
	CdnTrackingDomainABTestTag        string  `json:"cdn_tdabt,omitempty"` // imp，click,only_imp 的cdn切量标记
}

func handleExtData2MAS(r *mvutil.RequestParams) {
	var extData2MAS ExtData2MAS
	fcc := extractor.GetFREQ_CONTROL_CONFIG()
	if fcc != nil && fcc.Status == 1 {
		extData2MAS.PriceFactor = r.Param.ExtDataInit.PriceFactor
		extData2MAS.PriceFactorGroupName = r.Param.ExtDataInit.PriceFactorGroupName
		extData2MAS.PriceFactorTag = r.Param.ExtDataInit.PriceFactorTag
		extData2MAS.PriceFactorFreq = r.Param.ExtDataInit.PriceFactorFreq
		extData2MAS.PriceFactorHit = r.Param.ExtDataInit.PriceFactorHit
	}

	if r.PlacementInfo != nil && r.PlacementInfo.ImpressionCap > 0 {
		extData2MAS.ImpressionCap = r.PlacementInfo.ImpressionCap
	}
	// v5
	if len(r.Param.ExtDataInit.V5AbtestTag) > 0 {
		extData2MAS.V5AbtestTag = r.Param.ExtDataInit.V5AbtestTag
	}

	if r.Param.Open == 1 {
		extData2MAS.SDKOpen = r.Param.Open
	}

	if r.Param.BandWidth > 0 {
		extData2MAS.BandWidth = r.Param.BandWidth
	}

	extData2MAS.TKSysTag = r.Param.ExtDataInit.TKSysTag
	// pion, 标记记录到track。。。z参数放的s信息。。。。。
	extData2MAS.ExpIds = r.Param.ExtDataInit.ExpIds
	// 记录ios 设备信息。
	extData2MAS.UpTime = r.Param.UpTime
	extData2MAS.Brt = r.Param.Brt
	extData2MAS.Vol = r.Param.Vol
	extData2MAS.Lpm = r.Param.Lpm
	extData2MAS.Font = r.Param.Font
	extData2MAS.ImeiABTest = r.Param.ExtDataInit.ImeiABTest

	extData2MAS.IsReturnWtickTag = r.Param.IsReturnWtick

	extData2MAS.Ntbarpt = r.Param.Ntbarpt
	extData2MAS.Ntbarpasbl = r.Param.Ntbarpasbl
	extData2MAS.AtatType = r.Param.AtatType
	extData2MAS.MappingIdfaTag = r.Param.ExtDataInit.MappingIdfaTag
	extData2MAS.ExcludeAopPkg = r.Param.ExtDataInit.ExcludeAopPkg
	extData2MAS.MappingIdfaCoverIdfaTag = r.Param.ExtDataInit.MappingIdfaCoverIdfaTag

	extData2MAS.TKCNABTestTag = r.Param.ExtDataInit.TKCNABTestTag
	extData2MAS.TKCNABTestAATag = r.Param.ExtDataInit.TKCNABTestAATag

	extData2MAS.MiskSpt = r.Param.MiSkSpt
	extData2MAS.MoreofferAndAppwallMvToPioneerTag = r.Param.ExtDataInit.MoreofferAndAppwallMvToPioneerTag
	extData2MAS.MpToPioneerTag = r.Param.ExtDataInit.MpToPioneerTag
	extData2MAS.HBSubsidyType = r.Param.ExtDataInit.HBSubsidyType

	if r.Param.UcParentUnitId > 0 {
		extData2MAS.ParentUnitId = r.Param.UcParentUnitId
	} else {
		extData2MAS.ParentUnitId = r.Param.ParentUnitId
	}
	extData2MAS.H5Type = r.Param.H5Type
	extData2MAS.MofType = r.Param.MofType
	extData2MAS.CrtRid = r.Param.ExtDataInit.CrtRid
	extData2MAS.CNTrackDomain = r.Param.ExtDataInit.CNTrackDomain
	extData2MAS.OnlineApiNeedOfferBidPrice = r.Param.OnlineApiNeedOfferBidPrice
	extData2MAS.TrackDomainByCountryCode = r.Param.ExtDataInit.TrackDomainByCountryCode
	extData2MAS.TmaxABTestTag = r.Param.ExtDataInit.TmaxABTestTag
	extData2MAS.UseDynamicTmax = r.Param.ExtDataInit.UseDynamicTmax
	extData2MAS.Dnt = r.Param.Dnt
	extData2MAS.DeviceGEOCCMatch = r.Param.ExtDataInit.DeviceGEOCCMatch
	extData2MAS.ThreeLetterCountry = r.Param.ExtDataInit.ThreeLetterCountry
	extData2MAS.RequestBidServer = r.Param.IsHitRequestBidServer // 是否请求了 bid server
	extData2MAS.ParentAdType = r.Param.ParentAdType
	extData2MAS.ParentExchange = r.Param.ParentExchange
	extData2MAS.CdnTrackingDomain = r.Param.UseCdnTrackingDomain
	extData2MAS.CdnTrackingDomainABTestTag = r.Param.ExtDataInit.CdnTrackingDomainABTestTag

	jsonvalue, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(extData2MAS)
	r.Param.ExtData2MAS = string(jsonvalue)
}

// IsSupportMraid
// http://confluence.mobvista.com/pages/viewpage.action?pageId=24527937
func IsSupportMraid(reqCtx *mvutil.ReqCtx) bool {
	param := reqCtx.ReqParams.Param
	// openapi 不支持
	if param.RequestPath == mvconst.PATHOnlineApi {
		return false
	}
	// native不支持
	if param.AdType == mvconst.ADTypeNative {
		return false
	}

	if mvutil.IsHbOrV3OrV5Request(param.RequestPath) {
		// new banner: 支持
		if mvutil.IsBannerOrSplashOrDI(param.AdType) {
			return true
		}

		// 半屏对mraid情况的支持可能存在不兼容或者效果不好的情况，针对半屏的情况，设置成不支持mraid
		if mvutil.IsNewIv(&param) && reqCtx.ReqParams.UnitInfo.Unit.AdSpaceType == mvconst.AdSpaceTypeHalfScreen {
			return false
		}

		// rv iv ia: 指定版本支持
		if param.AdType == mvconst.ADTypeInterstitialVideo ||
			param.AdType == mvconst.ADTypeRewardVideo ||
			param.AdType == mvconst.ADTypeInteractive {
			// mraid 版本支持： android 12.2.0, ios 5.5.0 以上
			if (param.Platform == mvconst.PlatformAndroid && param.FormatSDKVersion.SDKVersionCode >= 120200) ||
				(param.Platform == mvconst.PlatformIOS && param.FormatSDKVersion.SDKVersionCode >= 50500) {
				return true
			}
		}
		return false
	}
	return false
}

// doMasAbtest MAS和AS的请求放到同一个池子里，统一切量
func doMasAbtest(reqCtx *mvutil.ReqCtx, request *rtb.BidRequest, backendCtx *mvutil.BackendCtx) {
	// mas_abtest 有切as或mas时，才要切量
	if backendCtx.IsBidAdServer || backendCtx.IsBidMAS {
		goMas := false
		randValue := rand.Float64()
		masAbtest, ifFind := extractor.GetMasAbtest()
		if ifFind && masAbtest != nil {
			// 切量：先按 platform+ad_type+cc 维度切量， 分别都支持all, 再按unit维度切
			tagid := strconv.FormatInt(reqCtx.ReqParams.Param.UnitID, 10)
			platform := mvconst.GetPlatformStr(reqCtx.ReqParams.Param.Platform) // str: android/ios
			adType := strconv.FormatInt(int64(reqCtx.ReqParams.Param.FormatAdType), 10)
			cc := reqCtx.ReqParams.Param.CountryCode
			for _, uCfg := range masAbtest {
				// 如果unit在黑名单里， 不切量
				if len(uCfg.Blacklist) > 0 {
					if mvutil.InStrArray(tagid, uCfg.Blacklist) {
						goMas = false
						break
					}
				}
				// 按platform+adtype+cc切量
				if len(uCfg.Key) > 0 {
					fields := strings.Split(uCfg.Key, "_")
					if len(fields) == 3 && (fields[0] == "ALL" || fields[0] == platform) &&
						(fields[1] == "ALL" || fields[1] == adType) && (fields[2] == "ALL" || fields[2] == cc) && randValue < uCfg.Rate {
						goMas = true
						break
					}

				} else if uCfg.UnitID == tagid && randValue < uCfg.Rate {
					// 按unit维度切量
					goMas = true
					break
				}
			}
		}
		// 对于sdk banner，只能通过pioneer切给as，不走adnet->adx->as。
		if mvutil.IsBannerOrSplashOrNativeH5(reqCtx.ReqParams.Param.AdType) {
			goMas = renderBannerOrSplashGoMas(reqCtx.ReqParams.Param.AdType, reqCtx.ReqParams.Param.UnitID, reqCtx.ReqParams.IsHBRequest)
		}
		if reqCtx.ReqParams.Param.ApiVersionCode < 1040000 {
			goMas = false
		}
		//Online API 符合切量到Mas的开关
		param := reqCtx.ReqParams.Param
		if onlineDoMasAbtest(int64(param.RequestType), param.PublisherID, param.AppID, param.UnitID, int64(param.AdType), param.DebugDspId) {
			goMas = true
		}

		if reqCtx.ReqParams.IsTopon {
			goMas = true
		}
		if goMas {
			backendCtx.IsBidMAS = true
			backendCtx.IsBidAdServer = false
			mas := "mas"
			user := request.User
			if user == nil {
				user = &rtb.BidRequest_User{}
			}
			user.Customdata = &mas // 不想加多字段， 用这个无意义字段当abtest用
			request.User = user
		} else {
			backendCtx.IsBidAdServer = true
			backendCtx.IsBidMAS = false
		}

	}
}

func renderBannerOrSplashGoMas(adType int32, unitId int64, isHBRequest bool) bool {
	// native h5直接可以走mas
	if adType == mvconst.ADTypeNativeH5 {
		return true
	}
	// 切量的unit，则gomas
	toAsABTestConf := extractor.GetBANNER_TO_AS_ABTEST_CONF()
	if rate, ok := toAsABTestConf[strconv.FormatInt(unitId, 10)]; ok {
		randVal := rand.Intn(100)
		if rate > randVal {
			return true
		} else {
			return false
		}
	}
	adnConf, _ := extractor.GetADNET_SWITCHS()
	if len(adnConf) == 0 {
		return false
	}
	// banner,splash 整体按比例切量
	if adType == mvconst.ADTypeSdkBanner {
		key := "bannerCanGoMasRate"
		if isHBRequest {
			key = "bannerCanGoMasRateHB"
		}
		if bannerCanGoMasRate, ok := adnConf[key]; ok {
			randVal := rand.Intn(100)
			if bannerCanGoMasRate > randVal {
				return true
			}
		}
	}
	if adType == mvconst.ADTypeSplash {
		key := "splashCanGoMasRate"
		if isHBRequest {
			key = "splashCanGoMasRateHB"
		}
		if splashCanGoMasRate, ok := adnConf[key]; ok {
			randVal := rand.Intn(100)
			if splashCanGoMasRate > randVal {
				return true
			}
		}
	}
	if adType == mvconst.ADTypeInterstitialSdk {
		return true
	}
	return false

}

func getSkipByAdType(adType int32) int32 {
	if adType == mvconst.ADTypeRewardVideo {
		return 0
	} else if adType == mvconst.ADTypeInterstitialVideo {
		return 1
	} else {
		return 1
	}
}

func videoTemplate(platform, orientation int) []int32 {
	var tamplateIds []int32
	if platform == mvconst.PlatformAndroid || platform == mvconst.PlatformIOS {
		tamplateIds = append(tamplateIds, MIDDLE_BLACK_SCREEN)
		tamplateIds = append(tamplateIds, MIDDLE_FUR_SCREEN)
		tamplateIds = append(tamplateIds, STRETCH_SCREEN)
		// return tamplateIds
	}
	if platform == mvconst.PlatformIOS && orientation == mvconst.ORIENTATION_PORTRAIT {
		tamplateIds = append(tamplateIds, ABOVE_VIDEO)
		tamplateIds = append(tamplateIds, STOREKIT_VIODE)
		tamplateIds = append(tamplateIds, IMAGE_VIDOE)

		// return tamplateIds
	}
	if platform == mvconst.PlatformAndroid && orientation == mvconst.ORIENTATION_PORTRAIT {
		tamplateIds = append(tamplateIds, ABOVE_VIDEO)
		tamplateIds = append(tamplateIds, IMAGE_VIDOE)
		// return tamplateIds
	}
	return tamplateIds
}

func templateName(templateId int32) string {
	switch templateId {
	case MIDDLE_BLACK_SCREEN:
		return "middle_black_screen"
	case MIDDLE_FUR_SCREEN:
		return "middle_fur_screen"
	case ABOVE_VIDEO:
		return "above_video"
	case STOREKIT_VIODE:
		return "storekit_video"
	case IMAGE_VIDOE:
		return "image_video"
	case STRETCH_SCREEN:
		return "stretch_screen"
	default:
		return ""
	}
}

// getWidthAndHeight  100x200 -> 100,200
func getWidthAndHeight(imageSize string) (width, height string) {
	arr := strings.Split(imageSize, "x")
	if len(arr) >= 2 && len(arr[0]) > 0 && len(arr[1]) > 0 {
		width = arr[0]
		height = arr[1]
	}
	return
}

func getVideoLength(lenth string) (int, error) {
	result := 0
	lenArr := strings.Split(lenth, ":")
	if len(lenArr) != 3 {
		return result, errors.New("input params not in the correct format length=" + lenth)
	}
	hour, err := strconv.Atoi(lenArr[0])
	if err != nil {
		return result, errors.New("convert hour data error:" + err.Error())
	}
	result += hour * 3600
	min, err := strconv.Atoi(lenArr[1])
	if err != nil {
		return result, errors.New("convert min data error:" + err.Error())
	}
	result += min * 60
	secArr := strings.Split(lenArr[2], ".")
	if len(secArr) == 0 {
		return result, errors.New("convert sec data error not in correct format data=" + lenArr[2])
	}
	sec, err := strconv.Atoi(secArr[0])
	if err != nil {
		return result, errors.New("convert sec data error:" + err.Error())
	}
	result += sec
	// if len(secArr) > 1 {
	// 	result += 1
	// }
	return result, nil
}

func getAdTypeDesc(adType int32) string {
	switch adType {
	case mvconst.ADTypeRewardVideo:
		return "reward_video"
	case mvconst.ADTypeOnlineVideo:
		return "online_video"
	case mvconst.ADTypeInterstitialVideo:
		return "interstitial_video"
	case mvconst.ADTypeNativeVideo:
		return "native_video"
	case mvconst.ADTypeNativePic:
		return "native_image"
	default:
		return ""
	}
}

func renderGdtClickUrl(clickUrl string, req *mvutil.RequestParams) string {
	if len(clickUrl) == 0 || req == nil || req.AppInfo == nil {
		return ""
	}
	clickUrl = strings.Replace(clickUrl, "__WIDTH__", strconv.Itoa(int(req.Param.ScreenWidth)), -1)
	clickUrl = strings.Replace(clickUrl, "__HEIGHT__", strconv.Itoa(int(req.Param.ScreenHeigh)), -1)
	packageName := req.AppInfo.RealPackageName
	if req.Param.Platform == mvconst.PlatformIOS {
		packageName = req.AppInfo.App.BundleId
	}
	clickUrl = strings.Replace(clickUrl, "app_bundle_id", packageName, -1)
	return clickUrl
}

func renderVideoViewLink(videoViewLink string, ori int, videoDuration int32) string {
	strVideoDuration := strconv.FormatInt(int64(videoDuration), 10)
	videoViewLink = strings.Replace(videoViewLink, "__VIDEO_TIME__", strVideoDuration, -1)
	videoViewLink = strings.Replace(videoViewLink, "__BEGIN_TIME__", "0", -1)
	videoViewLink = strings.Replace(videoViewLink, "__END_TIME__", strVideoDuration, -1)
	videoViewLink = strings.Replace(videoViewLink, "__PLAY_FIRST_FRAME__", "1", -1)
	videoViewLink = strings.Replace(videoViewLink, "__PLAY_LAST_FRAME__", "1", -1)
	scene := "2"
	if ori == mvconst.ORIENTATION_LANDSCAPE {
		scene = "4"
	}
	videoViewLink = strings.Replace(videoViewLink, "__SCENE__", scene, -1)
	videoViewLink = strings.Replace(videoViewLink, "__TYPE__", "1", -1)
	videoViewLink = strings.Replace(videoViewLink, "__BEHAVIOR__", "2", -1)
	videoViewLink = strings.Replace(videoViewLink, "__STATUS__", "0", -1)
	return videoViewLink
}

// func renderGdtVideo(adTmp *corsair_proto.Campaign, param *mvutil.Params, qKey string, videoUrl string) error {
// 	vInfo, err := redis.LocalRedisGet(qKey)
// 	needCreativeLog := false
// 	if err != nil {
// 		needCreativeLog = true
// 		mvutil.Logger.Runtime.Warnf("request_id=[%s] GDT GetVideoInfo key=[%s] from redis error:%s",
// 			param.RequestID, qKey, err.Error())
// 	}
// 	if !needCreativeLog {
// 		var videoCacheData *mvutil.VideoCacheData
// 		var err error
// 		if param.Platform == mvconst.PlatformIOS && param.FormatSDKVersion.SDKVersionCode < mvconst.IOSSupportVideoUrlWithParams {
// 			videoCacheData, err = mvutil.VideoMetaDataWithUrl(vInfo)
// 			if err != nil {
// 				mvutil.Logger.Runtime.Warnf("request_id=[%s] GDT DecodeVinfo key=[%s] data=[%s] error:%s",
// 					param.RequestID, qKey, vInfo, err.Error())
// 			}
// 			if videoCacheData != nil && len(videoCacheData.VideoMTGCdnURL) > 0 {
// 				*adTmp.VideoURL = videoCacheData.VideoMTGCdnURL
// 			} else {
// 				// ios 视频cdn链接有问题，不返回广告
// 				return errors.New("gdt creative error")
// 			}
// 		} else {
// 			videoCacheData, err = mvutil.VideoMetaDataWithoutUrl(vInfo)
// 			if err != nil {
// 				mvutil.Logger.Runtime.Warnf("request_id=[%s] GDT DecodeVinfo key=[%s] data=[%s] error:%s",
// 					param.RequestID, qKey, vInfo, err.Error())
// 			}
// 		}
// 		if videoCacheData != nil {
// 			*adTmp.VideoSize = videoCacheData.VideoSize
// 			*adTmp.VideoLength = videoCacheData.VideoLen
// 			*adTmp.VideoResolution = videoCacheData.VideoResolution
// 		}
// 	}
// 	if needCreativeLog {
// 		if param.Platform == mvconst.PlatformIOS && param.FormatSDKVersion.SDKVersionCode < mvconst.IOSSupportVideoUrlWithParams {
// 			mvutil.Logger.Creative.Infof("%s\t3\t%s", qKey, videoUrl)
// 		} else {
// 			mvutil.Logger.Creative.Infof("%s\t2\t%s", qKey, videoUrl)
// 		}
// 		// 需要拿我们的视频数据，但是还没有的情况下，视频不返回
// 		return errors.New("gdt creative error")
// 		// *adTmp.VideoLength = int32(ad.VideoLength)
// 	}
// 	return nil
// }

func fillBannerTracking(adTracking *corsair_proto.AdTracking, bid *rtb.BidResponse_SeatBid_Bid) {
	if len(bid.GetImptracker()) > 0 {
		adTracking.Impression = append(adTracking.Impression, bid.GetImptracker())
	}
	if len(bid.GetClicktracker()) > 0 {
		adTracking.Click = append(adTracking.Click, bid.GetClicktracker())
	}
}

func onlineDoMasAbtest(requestType, publisherId, appId, unitId, adType int64, debugDspId int) bool {
	if requestType == mvconst.REQUEST_TYPE_OPENAPI_AD &&
		rand.Intn(1000) < getOlineApiUseAdxMasRate(publisherId, appId, unitId, adType, debugDspId) {
		return true
	}
	return false
}

func getOlineApiUseAdxMasRate(publisherId, appId, unitId, adType int64, debugDspId int) int {
	conf := extractor.GetONLINE_API_USE_ADX_MAS()

	// debug
	if debugDspId == 13 { //13  一定走MAS
		return 1000
	} else if debugDspId == 6 { //6  一定走AS
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

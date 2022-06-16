package backend

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"hash/crc32"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/gogo/protobuf/proto"
	jsoniter "github.com/json-iterator/go"
	uuid "github.com/satori/go.uuid"
	"github.com/valyala/fasthttp"
	"gitlab.mobvista.com/ADN/adnet/internal/corsair_proto"
	"gitlab.mobvista.com/ADN/adnet/internal/extractor"
	hbconst "gitlab.mobvista.com/ADN/adnet/internal/hb/constant"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/filter"
	hbreqctx "gitlab.mobvista.com/ADN/adnet/internal/hb/req_context"
	"gitlab.mobvista.com/ADN/adnet/internal/mvconst"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	"gitlab.mobvista.com/ADN/adnet/internal/protocol"
	"gitlab.mobvista.com/ADN/adnet/internal/redis"
	"gitlab.mobvista.com/ADN/adnet/internal/watcher"
	adxconst "gitlab.mobvista.com/ADN/adx_common/constant"
	"gitlab.mobvista.com/ADN/exporter/metrics"
	"gitlab.mobvista.com/ADN/mtg_openrtb/pkg/di/abtest"
	"gitlab.mobvista.com/ADN/mtg_openrtb/pkg/mtg_hb_rtb"
	rtb "gitlab.mobvista.com/ADN/mtg_openrtb/pkg/mtgrtb"
	"gitlab.mobvista.com/ADN/mtg_openrtb/pkg/native"
	openrtb "gitlab.mobvista.com/ADN/mtg_openrtb/pkg/openrtb_v2"
	"gitlab.mobvista.com/ADN/mtg_openrtb/pkg/vast"
	smodel "gitlab.mobvista.com/ADN/structs/model"
	"gitlab.mobvista.com/adserver/recommend_protocols/go/ad_server"
	algoab "gitlab.mobvista.com/algo-engineering/abtest-sdk-go"
	"gitlab.mobvista.com/algo-engineering/abtest-sdk-go/pkg/biz"
	"gitlab.mobvista.com/mae/go-kit/decimal"
	"gitlab.mobvista.com/voyager/common/deviceId"
	"gitlab.mobvista.com/voyager/common/deviceId/deviceId_algo"
)

type MAdxBackend struct {
	Backend
}

const (
	NON_SECURE bool = false
	SECURE     bool = true
)

const (
	SUPPORT    bool = true
	UN_SUPPORT bool = false
)

const (
	FORCE   bool = true
	NOFORCE bool = false
)

const (
	GDTCM = 1
	GDTCU = 2
	GDTCT = 3
	GDTOT = 0

	ReplaceClickMarcoAndroidVersionCode = 90800
	ReplaceClickMarcoiOSVersionCode     = 40904
)

const (
	MIDDLE_BLACK_SCREEN int32 = 1 // 视频居中播放，背景为黑屏
	MIDDLE_FUR_SCREEN   int32 = 2 // 视频居中播放，背景为毛玻璃
	ABOVE_VIDEO         int32 = 3 // 上方显示视频，下方显示icon,button
	STOREKIT_VIODE      int32 = 4 // 上方显示视频，下方显示storekit
	IMAGE_VIDOE         int32 = 5 // 上下显示图片，中间显示视频
	STRETCH_SCREEN      int32 = 6 // 视频拉伸铺满屏幕
)

const (
	ENDCARD  int32 = 0
	STOREKIT int32 = 1
)

const ImpId string = "1"

const (
	// 最大视频时长
	MAXDURATION int32 = 30
	// 码率
	MAXBITRATE int32 = 2000
)

// const (
//	ICONID int32 = iota + 1
//	VIDEOID
//	TILTEID
//	DECSID
//	RATINGID
//	CTAID
//	IMAGEID
// )

func fillrateConfig(reqCtx *mvutil.ReqCtx) (cfg *smodel.ConfigAlgorithmFillRate, ifFind bool) {
	cfg, ifFind = extractor.GetFillRateControllConfig(extractor.GetFillRateKey(reqCtx.ReqParams.Param.UnitID,
		reqCtx.ReqParams.Param.AppID, reqCtx.ReqParams.Param.Platform, reqCtx.ReqParams.Param.CountryCode))
	if ifFind {
		return cfg, ifFind
	}

	cfg, ifFind = extractor.GetFillRateControllConfig(extractor.GetFillRateKey(reqCtx.ReqParams.Param.UnitID,
		reqCtx.ReqParams.Param.AppID, reqCtx.ReqParams.Param.Platform, "ALL"))
	return cfg, ifFind
}

func (backend MAdxBackend) filterBackend(reqCtx *mvutil.ReqCtx) int {
	if reqCtx == nil || reqCtx.ReqParams == nil {
		return mvconst.ParamInvalidate
	}
	// DSP对接一期 2.1 去掉原来对于SDK的版本限制 jira: ADNET-151
	// if reqCtx.ReqParams.Param.AdType == mvconst.ADTypeRewardVideo || reqCtx.ReqParams.Param.AdType == mvconst.ADTypeInterstitialVideo {
	//	if reqCtx.ReqParams.Param.Platform == mvconst.PlatformIOS && reqCtx.ReqParams.Param.FormatSDKVersion.SDKVersionCode < mvconst.ADTRACKINGIOSADX {
	//		return true
	//	}
	//	if reqCtx.ReqParams.Param.Platform == mvconst.PlatformAndroid && reqCtx.ReqParams.Param.FormatSDKVersion.SDKVersionCode < mvconst.AdTrackingAndroidNativePicClick {
	//		return true
	//	}
	// }
	// online-api的请求不走sdk版本过滤
	// if reqCtx.ReqParams.Param.RequestType == mvconst.REQUEST_TYPE_OPENAPI_AD {
	// 	return mvconst.BackendOK
	// }

	if !reqCtx.ReqParams.IsHBRequest {
		cfg, ifFind := fillrateConfig(reqCtx)
		if ifFind && cfg != nil && cfg.ControlMode == mvconst.FillRate && cfg.Rate == 0 {
			return mvconst.BackendFillRateFilter
		}
	}

	// if filterAdTracking(reqCtx.ReqParams.Param.FormatAdType,
	// 	reqCtx.ReqParams.Param.Platform,
	// 	mvconst.AdTracking_Both_Click_Imp,
	// 	reqCtx.ReqParams.Param.FormatSDKVersion) {
	// 	return mvconst.BackendTrackingFilter
	// }
	return mvconst.BackendOK
}

func constructPublisher(publisherId int64, publisherName string) *rtb.BidRequest_Publisher {
	publisher := &rtb.BidRequest_Publisher{}
	pubId := strconv.FormatInt(publisherId, 10)
	publisher.Id = &pubId
	publisher.Name = &publisherName
	return publisher
}

func constructBanner4IA(width, height int32, orientation *rtb.BidRequest_Imp_OrientationType, isSupportMraid bool) *rtb.BidRequest_Imp_Banner {
	banner := &rtb.BidRequest_Imp_Banner{}

	banner.W = &width
	banner.H = &height
	banner.Battr = append(banner.Battr, rtb.CreativeAttribute_AD_CAN_BE_SKIPPED)

	pos := rtb.BidRequest_Imp_AD_POSITION_FULLSCREEN
	banner.Pos = &pos
	// 自由流量RV、IV要求对方返回拼好的endcard，传application/javascript
	// 横竖屏要求
	var ext rtb.BidRequest_Imp_Banner_Ext
	ext.Orientation = orientation
	// 设置banner type
	// 由原来的banner type区分改为用adtype来区分
	adtype := rtb.BidRequest_Imp_Banner_Ext_ADTYPE_INTERACTIVE_ADS
	ext.Adtype = &adtype
	banner.Ext = &ext
	if isSupportMraid {
		banner.Mimes = append(banner.Mimes,
			"application/javascript", "image/jpeg", "image/jpg", "text/html", "image/png", "text/css", "image/gif")
		banner.Api = append(banner.Api, rtb.APIFramework_MRAID_1, rtb.APIFramework_MRAID_2, rtb.APIFramework_MRAID_3)
	} else {
		banner.Btype = append(banner.Btype, rtb.BidRequest_Imp_Banner_XHTML_TEXT_AD, rtb.BidRequest_Imp_Banner_XHTML_BANNER_AD, rtb.BidRequest_Imp_Banner_JAVASCRIPT_AD, rtb.BidRequest_Imp_Banner_IFRAME)
		banner.Mimes = append(banner.Mimes, "image/jpg", "image/png", "image/jpeg")
	}
	banner.Api = append(banner.Api, rtb.APIFramework_OMID_1)
	return banner
}

func constructSdkBanner(reqCtx *mvutil.ReqCtx, isSupportMraid bool) *rtb.BidRequest_Imp_Banner {
	// 处理sdk_banner及splash
	banner := &rtb.BidRequest_Imp_Banner{}

	width := int32(reqCtx.ReqParams.Param.SdkBannerUnitWidth)
	height := int32(reqCtx.ReqParams.Param.SdkBannerUnitHeight)
	banner.W = &width
	banner.H = &height
	banner.Battr = append(banner.Battr, rtb.CreativeAttribute_ADOBE_FLASH)
	var ext rtb.BidRequest_Imp_Banner_Ext
	// banner不需传orientation
	ori := rtb.BidRequest_Imp_ORIENTATION_UNKNOWN
	// 从bannertype改为用adtype区分广告类型
	adtype := rtb.BidRequest_Imp_Banner_Ext_ADTYPE_SDK_BANNER
	if reqCtx.ReqParams.Param.AdType == mvconst.ADTypeSplash {
		adtype = rtb.BidRequest_Imp_Banner_Ext_ADTYPE_SPLASH
		// 开屏广告需要根据orientation返回大图。横屏尺寸 就召回1280x720，竖屏尺寸就召回720x1280
		if reqCtx.ReqParams.Param.FormatOrientation == mvconst.ORIENTATION_PORTRAIT {
			ori = rtb.BidRequest_Imp_PORTRAIT
		} else if reqCtx.ReqParams.Param.FormatOrientation == mvconst.ORIENTATION_LANDSCAPE {
			ori = rtb.BidRequest_Imp_LANDSCAPE
		}
	}
	if reqCtx.ReqParams.Param.AdType == mvconst.ADTypeInterstitialSdk {
		adtype = rtb.BidRequest_Imp_Banner_Ext_ADTYPE_INTERSTITIAL
	}

	if reqCtx.ReqParams.Param.AdType == mvconst.ADTypeMoreOffer {
		adtype = rtb.BidRequest_Imp_Banner_Ext_ADTYPE_MORE_OFFER
	}

	if reqCtx.ReqParams.Param.AdType == mvconst.ADTypeAppwall {
		adtype = rtb.BidRequest_Imp_Banner_Ext_ADTYPE_APPWALL
	}

	ext.Orientation = &ori
	ext.Adtype = &adtype
	banner.Ext = &ext
	if isSupportMraid {
		banner.Mimes = append(banner.Mimes,
			"application/javascript", "image/jpeg", "image/jpg", "text/html", "image/png", "text/css", "image/gif")
		banner.Api = append(banner.Api, rtb.APIFramework_MRAID_1, rtb.APIFramework_MRAID_2, rtb.APIFramework_MRAID_3)
	} else {
		banner.Btype = append(banner.Btype, rtb.BidRequest_Imp_Banner_XHTML_TEXT_AD, rtb.BidRequest_Imp_Banner_XHTML_BANNER_AD, rtb.BidRequest_Imp_Banner_JAVASCRIPT_AD, rtb.BidRequest_Imp_Banner_IFRAME)
		banner.Mimes = append(banner.Mimes, "image/jpg", "image/png", "image/jpeg")
	}
	banner.Api = append(banner.Api, rtb.APIFramework_OMID_1)
	return banner
}

// mraid 用的banner, 赋值逻辑与sdk banner基本一致
func constructMraidBanner(width, height int32, isSupportMraid bool, orientation int) *rtb.BidRequest_Imp_Banner {
	if orientation == mvconst.ORIENTATION_PORTRAIT {
		width = 720
		height = 1280
	} else {
		width = 1280
		height = 720
	}
	banner := &rtb.BidRequest_Imp_Banner{}

	banner.W = &width
	banner.H = &height
	banner.Battr = append(banner.Battr, rtb.CreativeAttribute_ADOBE_FLASH)
	var ext rtb.BidRequest_Imp_Banner_Ext
	// banner不需传orientation
	ori := rtb.BidRequest_Imp_ORIENTATION_UNKNOWN
	// 从bannertype改为用adtype区分广告类型
	adtype := rtb.BidRequest_Imp_Banner_Ext_ADTYPE_SDK_BANNER

	ext.Orientation = &ori
	ext.Adtype = &adtype
	banner.Ext = &ext
	if isSupportMraid {
		banner.Mimes = append(banner.Mimes,
			"application/javascript", "image/jpeg", "image/jpg", "text/html", "image/png", "text/css", "image/gif")
		banner.Api = append(banner.Api, rtb.APIFramework_MRAID_1, rtb.APIFramework_MRAID_2, rtb.APIFramework_MRAID_3)
	} else {
		banner.Btype = append(banner.Btype, rtb.BidRequest_Imp_Banner_XHTML_TEXT_AD, rtb.BidRequest_Imp_Banner_XHTML_BANNER_AD, rtb.BidRequest_Imp_Banner_JAVASCRIPT_AD, rtb.BidRequest_Imp_Banner_IFRAME)
		banner.Mimes = append(banner.Mimes, "image/jpg", "image/png", "image/jpeg")
	}
	banner.Api = append(banner.Api, rtb.APIFramework_OMID_1)
	return banner
}

// constructBanner4Video rv/iv请求里也会带banner
func constructBanner4Video(adType, width, height int32, orientation int, isSupportMraid bool) *rtb.BidRequest_Imp_Banner {
	banner := &rtb.BidRequest_Imp_Banner{}

	// 作为该次请求每个companion的唯一标识
	id := "1"
	banner.Id = &id
	banner.Battr = append(banner.Battr, rtb.CreativeAttribute_AD_CAN_BE_SKIPPED)

	if isSupportMraid {
		banner.Mimes = append(banner.Mimes,
			"application/javascript", "image/jpeg", "image/jpg", "text/html", "image/png", "text/css", "image/gif")
		banner.Api = append(banner.Api, rtb.APIFramework_MRAID_1, rtb.APIFramework_MRAID_2, rtb.APIFramework_MRAID_3)
	} else {
		banner.Btype = append(banner.Btype, rtb.BidRequest_Imp_Banner_XHTML_TEXT_AD, rtb.BidRequest_Imp_Banner_XHTML_BANNER_AD, rtb.BidRequest_Imp_Banner_IFRAME)
		banner.Mimes = append(banner.Mimes, "image/jpg", "image/png", "image/jpeg")
	}
	banner.Api = append(banner.Api, rtb.APIFramework_OMID_1)

	banner.W = &width
	banner.H = &height

	pos := rtb.BidRequest_Imp_AD_POSITION_FULLSCREEN
	banner.Pos = &pos

	var ext rtb.BidRequest_Imp_Banner_Ext
	enOrientation := rtb.BidRequest_Imp_OrientationType(orientation)
	ext.Orientation = &enOrientation
	videoType := getVideoType(adType)
	if videoType == rtb.BidRequest_Imp_Video_Ext_REWARDED_VIDEO {
		isReward := true
		ext.IsRewarded = &isReward
	}
	banner.Ext = &ext

	return banner
}

func getVideoType(adType int32) rtb.BidRequest_Imp_Video_Ext_VideoType {
	if adType == mvconst.ADTypeRewardVideo {
		return rtb.BidRequest_Imp_Video_Ext_REWARDED_VIDEO
	} else {
		return rtb.BidRequest_Imp_Video_Ext_INTERSTITIAL_VIDEO
	}
}

func constructVideoExt(adType, width, height, sdkVerionCode int32, platform, orientation int, path string) *rtb.BidRequest_Imp_Video_Ext {
	ext := &rtb.BidRequest_Imp_Video_Ext{}

	enOrientation := rtb.BidRequest_Imp_OrientationType(orientation)
	ext.Orientation = &enOrientation

	videoType := getVideoType(adType)
	ext.Videotype = &videoType
	if videoType == rtb.BidRequest_Imp_Video_Ext_REWARDED_VIDEO {
		isReward := true
		ext.IsRewarded = &isReward
	}

	// 非SDK流量默认不传该字段
	if mvutil.IsHbOrV3OrV5Request(path) {
		ext.Videoendtype = append(ext.Videoendtype, rtb.VideoEndType_AUTO_PLAY_ENDCARD)
		if platform == mvconst.PlatformIOS && sdkVerionCode >= mvconst.AdTrackingIOSStorekit {
			ext.Videoendtype = append(ext.Videoendtype, rtb.VideoEndType_WEBVIEW_APP_STORE)
		}
	}

	templates := videoTemplate(platform, orientation)
	for _, templateId := range templates {
		vt := constructVideoTemplate(templateId, width, height, &enOrientation)
		ext.Videotemplate = append(ext.Videotemplate, vt)
	}

	// 自有SDK流量默认传1，其他流量不传 todo---
	endcardOnly := true
	ext.Endcardonly = &endcardOnly

	return ext
}

func constructVideoTemplate(templateId, width, height int32, orientation *rtb.BidRequest_Imp_OrientationType) *rtb.BidRequest_Imp_Video_Ext_Videotemplate {
	videoTemplate := &rtb.BidRequest_Imp_Video_Ext_Videotemplate{}
	videoTemplate.Id = &templateId

	name := templateName(templateId)
	videoTemplate.Name = &name

	switch templateId {
	case MIDDLE_BLACK_SCREEN, MIDDLE_FUR_SCREEN, STRETCH_SCREEN:

		videoTemplate.Videoh = &height

		videoTemplate.Videow = &width

		videoTemplate.Videoorientation = orientation

	case ABOVE_VIDEO, STOREKIT_VIODE, IMAGE_VIDOE:
		fixOrientation := rtb.BidRequest_Imp_OrientationType(rtb.BidRequest_Imp_LANDSCAPE)
		videoTemplate.Videoorientation = &fixOrientation
	}
	return videoTemplate
}

func constructVideoWithCompainion(adType, width, height, sdkVerionCode int32, platform, orientation int, path string, isSupportMraid bool) *rtb.BidRequest_Imp_Video {
	if orientation == mvconst.ORIENTATION_PORTRAIT {
		width = 720
		height = 1280
	} else {
		width = 1280
		height = 720
	}
	video := &rtb.BidRequest_Imp_Video{}
	video.Mimes = append(video.Mimes, "video/mp4")

	maxDuration := MAXDURATION
	video.Maxduration = &maxDuration

	video.Protocols = []rtb.VideoBidResponseProtocol{rtb.VideoBidResponseProtocol_VAST_2_0, rtb.VideoBidResponseProtocol_VAST_3_0,
		rtb.VideoBidResponseProtocol_VAST_2_0_WRAPPER, rtb.VideoBidResponseProtocol_VAST_3_0_WRAPPER}

	linearity := rtb.BidRequest_Imp_Video_LINEAR
	video.Linearity = &linearity

	// 只有RV的battr=[16], IV为空。 因为RV视频不允许跳过
	if adType == mvconst.ADTypeRewardVideo {
		video.Battr = append(video.Battr, rtb.CreativeAttribute_AD_CAN_BE_SKIPPED)
	}

	video.Delivery = append(video.Delivery, rtb.BidRequest_Imp_Video_PROGRESSIVE)

	skip := getSkipByAdType(adType)
	video.Skip = &skip

	video.W = &width
	video.H = &height
	pos := rtb.BidRequest_Imp_AD_POSITION_FULLSCREEN // full screen, 见openrtb list 5.4
	video.Pos = &pos

	maxBitrate := MAXBITRATE
	video.Maxbitrate = &maxBitrate

	if isSupportMraid {
		video.Companiontype = append(video.Companiontype, rtb.BidRequest_Imp_STATIC, rtb.BidRequest_Imp_HTML, rtb.BidRequest_Imp_COMPANION_IFRAME)
		video.Api = append(video.Api, rtb.APIFramework_MRAID_1, rtb.APIFramework_MRAID_2, rtb.APIFramework_MRAID_3)
	} else {
		video.Companiontype = append(video.Companiontype, rtb.BidRequest_Imp_STATIC)
	}
	video.Api = append(video.Api, rtb.APIFramework_OMID_1)
	banner := constructBanner4Video(adType, width, height, orientation, isSupportMraid)

	video.Companionad = append(video.Companionad, banner)
	video.Ext = constructVideoExt(adType, width, height, sdkVerionCode, platform, orientation, path)
	return video
}

func (backend *MAdxBackend) getRequestNode() string {
	return backend.MAdxClient.GetNode()
}

func (backend MAdxBackend) composeHttpRequest(reqCtx *mvutil.ReqCtx, backendCtx *mvutil.BackendCtx, req *fasthttp.Request) error {
	if reqCtx == nil || backendCtx == nil || req == nil {
		return errors.New("GDTBackend composeHttpRequest params invalidate")
	}

	if reqCtx.ReqParams.IsTopon {
		// 请求新的Mas, 填充 asInfo字段
		if backendCtx.IsBidMAS {
			reqCtx.ReqParams.ToponRequest.User = &openrtb.BidRequest_User{Customdata: "mas"}
			// mas其实是pioneer的别名
			asInfo, err := ConstructMasInfo(reqCtx)
			if err != nil {
				return err
			}
			reqCtx.ReqParams.Param.MasDebugParam = asInfo
			openAsInfo := new(openrtb.AsInfo)
			b, _ := jsoniter.ConfigFastest.Marshal(asInfo)
			jsoniter.ConfigFastest.Unmarshal(b, openAsInfo)
			if reqCtx.ReqParams.IsTopon {
				reqCtx.ReqParams.ToponRequest.Mvext.AsInfo = openAsInfo
			}
		}

		// 补充ios的device.ext信息
		if reqCtx.ReqParams.ToponRequest.Device != nil && reqCtx.ReqParams.Param.Platform == mvutil.IOSPLATFORM {
			if reqCtx.ReqParams.ToponRequest.Device.Mvext == nil {
				reqCtx.ReqParams.ToponRequest.Device.Mvext = &openrtb.BidRequest_Device_Ext{}
			}
			reqCtx.ReqParams.ToponRequest.Device.Mvext.OsvUpTime = reqCtx.ReqParams.Param.OsvUpTime // 系统更新时间
			reqCtx.ReqParams.ToponRequest.Device.Mvext.Ram = reqCtx.ReqParams.Param.Ram             // 物理内存
			reqCtx.ReqParams.ToponRequest.Device.Mvext.Uptime = reqCtx.ReqParams.Param.UpTime       // 开机时间
			//reqCtx.ReqParams.ToponRequest.Device.Mvext.CountryCode = reqCtx.ReqParams.Param.CountryCode // 国家代码
			reqCtx.ReqParams.ToponRequest.Device.Mvext.TotalMemory = reqCtx.ReqParams.Param.TotalMemory // 硬盘尺寸
			reqCtx.ReqParams.ToponRequest.Device.Mvext.TimeZone = reqCtx.ReqParams.Param.TimeZone       // 时区
		}

		protoData, _ := proto.Marshal(reqCtx.ReqParams.ToponRequest)
		req.Header.SetMethod(backendCtx.Method)
		req.SetRequestURI(backendCtx.ReqPath)
		req.SetBody(protoData)
		req.Header.Set("mtg-rtb-version", "2.0")
	} else {
		bidRequest, err := renderBidRequest(reqCtx, backendCtx)
		if err != nil {
			return err
		}
		// adxRequestByte, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(bidRequest)
		// mvutil.Logger.Runtime.Debugf("=========== adx request url: %s, timeout: %d, ad_type: %s, data: %s",
		// backendCtx.ReqPath, *bidRequest.Tmax, mvutil.GetAdTypeStr(reqCtx.ReqParams.Param.AdType), string(adxRequestByte))
		protoData, err := proto.Marshal(bidRequest)
		if err != nil {
			return err
		}
		req.Header.SetMethod(backendCtx.Method)
		// for test
		// req.SetRequestURI("http://127.0.0.1:8102/hbrtb")
		req.SetRequestURI(backendCtx.ReqPath)
		req.SetBody(protoData)
	}
	return nil
}

func renderBidRequest(reqCtx *mvutil.ReqCtx, backendCtx *mvutil.BackendCtx) (*rtb.BidRequest, error) {
	bidRequest := &rtb.BidRequest{}
	bidRequest.Id = &reqCtx.ReqParams.Param.RequestID
	imp, err := constructImp(reqCtx)
	if err != nil {
		return nil, err
	}

	// abtest 透传字段逻辑
	abTestContxt := renderAbTestParam(reqCtx.ReqParams)

	// 填充dsp pkg blacklist
	for k, data := range backendCtx.BlackList {
		dspId := k
		blacklist := &rtb.BidRequest_Imp_BlackList{DspId: &dspId, BlackPkgList: data.PkgNames}
		imp.BlackList = append(imp.BlackList, blacklist)
	}
	// tmax abtest
	var (
		tmaxABTestGroup int
		tmaxABTest      int32
		timeoutConf     int32
		useDynamicTmax  string
	)
	tmaxABTestConf := extractor.GetTmaxABTestConf()
	if !reqCtx.ReqParams.IsHBRequest {
		timeoutConf = 2000
		if reqCtx.ReqParams.Param.RequestType == mvconst.REQUEST_TYPE_OPENAPI_AD {
			if v, ok := tmaxABTestConf.Timeout["online_api"]; ok { // 配置的 online api 流量 timeout
				timeoutConf = v
			}
		} else {
			if v, ok := tmaxABTestConf.Timeout["sdk"]; ok { // 配置的 sdk 非 hb 流量 timeout
				timeoutConf = v
			}
		}
	} else {
		timeoutConf = int32(backendCtx.Tmax)
		if reqCtx.ReqParams.Param.RequestType == mvconst.REQUEST_TYPE_OPENAPI_AD {
			if v, ok := tmaxABTestConf.Timeout["online_api"]; ok { // 配置的 online api 流量 timeout
				timeoutConf = v
			}
		} else if tmaxABTestConf.UseBidRequestTmax && reqCtx.ReqParams.Param.HBTmax > 0 { // 请求中给的 tmax
			timeoutConf = reqCtx.ReqParams.Param.HBTmax
		} else if v, ok := tmaxABTestConf.Timeout[strings.ToUpper(reqCtx.ReqParams.Param.MediationName)]; ok { // 配置的聚合平台 s2s 的 timeout
			timeoutConf = v
		} else {
			if v, ok = tmaxABTestConf.Timeout["default"]; ok { // 配置的默认 timeout, 包括 c2s 和没有传聚合名字的 s2s
				timeoutConf = v
			}
		}
	}
	networkCost := tmaxABTestConf.NetworkCost["hb_to_adx"]
	for _, c := range tmaxABTestConf.TmaxABTestConfigs {
		if (timeoutConf-networkCost) > 0 && selectorTmaxABTestFilter(reqCtx, c) {
			useDynamicTmax = "0"
			// 按设备切量，且保证 tmax 配置的实验组大于1组
			if mvutil.GetRandByGlobalTagId(&reqCtx.ReqParams.Param, mvconst.SALT_TMAX_ABTEST, 10000) < c.Rate && len(c.Tmax) > 1 {

				tmaxABTest = timeoutConf - networkCost // 计算后的 tmax
				// 500, [600, 700, 800, 900], 原始值 500
				// 650, [600, 700, 800, 900], 600
				// 750, [600, 700, 800, 900], 设备信息模2
				// 850, [600, 700, 800, 900], 设备信息模3
				// 950, [600, 700, 800, 900], 设备信息模4

				if tmaxABTest < c.Tmax[0] { // 排除小于配置的超时最小值
					tmaxABTestGroup = -1
					break
				}

				groups := len(c.Tmax)
				if tmaxABTest >= c.Tmax[groups-1] { // 直接判断是否大于配置的超时最大值
					group := mvutil.GetRandByGlobalTagId(&reqCtx.ReqParams.Param, mvconst.SALT_TMAX_GROUPS, groups)
					timeoutConf = c.Tmax[group]
					tmaxABTestGroup = group
					useDynamicTmax = "1"
					break
				}

				for i := 0; i < groups; i++ {
					if tmaxABTest >= c.Tmax[i] { // 找到实验组
						tmaxABTestGroup = i
					}
				}
				if tmaxABTestGroup > 0 {
					group := mvutil.GetRandByGlobalTagId(&reqCtx.ReqParams.Param, mvconst.SALT_TMAX_GROUPS, tmaxABTestGroup+1) // 这里是用 group 的 index + 1 作为取模的数来计算实验组
					timeoutConf = c.Tmax[group]
					tmaxABTestGroup = group
				} else {
					timeoutConf = c.Tmax[tmaxABTestGroup] // 直接选择了第一个配置
				}
				useDynamicTmax = "1"
				break
			}
		}
	}
	bidRequest.Tmax = &timeoutConf
	reqCtx.ReqParams.Param.DynamicTmax = timeoutConf
	reqCtx.ReqParams.Param.UseDynamicTmax = useDynamicTmax // 透传给 poineer, 0: 原来配置, 1: 使用动态 tmax
	reqCtx.ReqParams.Param.ExtDataInit.TmaxABTestTag = tmaxABTestGroup
	reqCtx.ReqParams.Param.ExtDataInit.UseDynamicTmax = useDynamicTmax

	hitRequestBidServerTest(reqCtx)
	handleExtData2MAS(reqCtx.ReqParams)
	doMasAbtest(reqCtx, bidRequest, backendCtx)
	// 请求AS填充as_req字段
	if backendCtx.IsBidAdServer {
		reqCtx.ReqParams.Param.AdxBidFloor = imp.GetBidfloor() / 100
		queryAs, err := ConstructASRequest(reqCtx)
		if err != nil {
			return nil, err
		}
		reqCtx.ReqParams.Param.AsDebugParam = queryAs
		asData, err := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(queryAs)
		if err != nil {
			return nil, err
		}
		asReq := string(asData)
		imp.AsReq = &asReq
	}
	// 请求新的Mas, 填充 asInfo字段
	if backendCtx.IsBidMAS || mvutil.IsAppwallOrMoreOffer(reqCtx.ReqParams.Param.AdType) || mvutil.IsRequestPioneerDirectly(&reqCtx.ReqParams.Param) {
		reqCtx.ReqParams.Param.AdxBidFloor = imp.GetBidfloor() / 100
		asInfo, err := ConstructMasInfo(reqCtx)
		if err != nil {
			return nil, err
		}
		bidRequest.AsInfo = asInfo // pione从这里取出ext记录
		reqCtx.ReqParams.Param.MasDebugParam = asInfo
	}

	requestType := int32(reqCtx.ReqParams.Param.RequestType)
	imp.RequestType = &requestType
	supportDownload := filterDownload(reqCtx.ReqParams.AppInfo.App.OfferPreference, reqCtx.ReqParams.Param.Platform)
	imp.SupportDownload = &supportDownload

	bidRequest.Imp = append(bidRequest.Imp, imp)

	bidRequest.App = constructApp(reqCtx)

	bidRequest.Device = constructDevice(reqCtx)

	if !reqCtx.ReqParams.IsHBRequest {
		at := rtb.BidRequest_SECOND_PRICE
		bidRequest.At = &at
	} else {
		at := rtb.BidRequest_FIRST_PRICE
		bidRequest.At = &at
		if len(reqCtx.ReqParams.Param.RTBSourceBytes) > 0 {
			var rtbSource rtb.BidRequest_Source
			if err := jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(reqCtx.ReqParams.Param.RTBSourceBytes, &rtbSource); err == nil {
				bidRequest.Source = &rtbSource
			}
		}
	}

	var BlackBundleList []string
	// 黑名单从 adx_medi_config改为从 unit.unit维度获取
	for k, v := range reqCtx.ReqParams.UnitInfo.Unit.BlackIABCategory {
		if len(v) == 0 {
			bidRequest.Bcat = append(bidRequest.Bcat, k)
			continue
		}
		bidRequest.Bcat = append(bidRequest.Bcat, v...)
	}

	bidRequest.Badv = reqCtx.ReqParams.UnitInfo.Unit.BlackDomain
	BlackBundleList = reqCtx.ReqParams.UnitInfo.Unit.BlackBundle

	// countrycode 黑名单
	if len(reqCtx.ReqParams.Param.CountryBlackPackageList) > 0 {
		BlackBundleList = append(BlackBundleList, reqCtx.ReqParams.Param.CountryBlackPackageList...)
	}
	bidRequest.Bapp = constructBapp(BlackBundleList, reqCtx)

	if reqCtx.ReqParams.AppInfo.App.Coppa == 1 {
		regs := &rtb.BidRequest_Regs{}
		coppa := true
		regs.Coppa = &coppa
		bidRequest.Regs = regs
	}
	// exchange记录为mintegral
	exchange := "mintegral"
	ext := &rtb.BidRequest_Ext{}
	ext.Exchange = &exchange
	ext.RequestBidServer = proto.Int32(reqCtx.ReqParams.Param.IsHitRequestBidServer)
	bidRequest.Ext = ext

	if bidRequest.Ext != nil {
		bidRequest.Ext.OnlyRequstThirdDsp = &reqCtx.ReqParams.Param.OnlyRequestThirdDsp
	} else {
		bidRequest.Ext = &rtb.BidRequest_Ext{
			OnlyRequstThirdDsp: &reqCtx.ReqParams.Param.OnlyRequestThirdDsp,
		}
	}

	// ab test context
	bidRequest.Ext.AbTestContext = abTestContxt

	return bidRequest, nil
}

func selectorTmaxABTestFilter(reqCtx *mvutil.ReqCtx, c *mvutil.TmaxABTestConf) bool {
	if len(c.Cloud) > 0 && !mvutil.InStrArray(reqCtx.ReqParams.Cloud, c.Cloud) {
		return false
	}
	if len(c.Region) > 0 && !mvutil.InStrArray(reqCtx.ReqParams.Region, c.Region) {
		return false
	}

	if len(c.CountryCode) > 0 && !mvutil.InStrArray(reqCtx.ReqParams.Param.CountryCode, c.CountryCode) {
		return false
	}
	var isHB int
	if reqCtx.ReqParams.IsHBRequest {
		isHB = 1
	}
	if len(c.IsHB) > 0 && !mvutil.InArray(isHB, c.IsHB) && !mvutil.InArray(-1, c.IsHB) {
		return false
	}

	if len(c.MediationName) > 0 && !mvutil.InStrArray(reqCtx.ReqParams.Param.MediationName, c.MediationName) {
		return false
	}
	if len(c.Platform) > 0 && !mvutil.InArray(reqCtx.ReqParams.Param.Platform, c.Platform) {
		return false
	}
	if len(c.RequestType) > 0 && !mvutil.InArray(reqCtx.ReqParams.Param.RequestType, c.RequestType) {
		return false
	}
	if len(c.AdType) > 0 && !mvutil.InInt32Arr(reqCtx.ReqParams.Param.AdType, c.AdType) {
		return false
	}
	if len(c.Scenario) > 0 && !mvutil.InStrArray(reqCtx.ReqParams.Param.Scenario, c.Scenario) {
		return false
	}
	if len(c.PublisherID) > 0 && !mvutil.InInt64Arr(reqCtx.ReqParams.Param.PublisherID, c.PublisherID) {
		return false
	}
	if len(c.AppID) > 0 && !mvutil.InInt64Arr(reqCtx.ReqParams.Param.AppID, c.AppID) {
		return false
	}
	if len(c.UnitID) > 0 && !mvutil.InInt64Arr(reqCtx.ReqParams.Param.UnitID, c.UnitID) {
		return false
	}

	return true
}

func getConnetType(netWorkId int) *rtb.BidRequest_Device_ConnectionType {
	connectType := rtb.BidRequest_Device_CONNECTION_UNKNOWN
	switch netWorkId {
	case mvconst.NETWORK_TYPE_2G:
		connectType = rtb.BidRequest_Device_CELL_2G
	case mvconst.NETWORK_TYPE_3G:
		connectType = rtb.BidRequest_Device_CELL_3G
	case mvconst.NETWORK_TYPE_4G:
		connectType = rtb.BidRequest_Device_CELL_4G
	case mvconst.NETWORK_TYPE_WIFI:
		connectType = rtb.BidRequest_Device_WIFI
	}
	return &connectType
}

func constructDevice(reqCtx *mvutil.ReqCtx) *rtb.BidRequest_Device {
	device := &rtb.BidRequest_Device{}
	device.Ua = &reqCtx.ReqParams.Param.UserAgent
	// 自有流量：支持coppa的开发者默认yes（app/unit配置表）
	dnt := false
	// coppa 1：表示遵守 coppa 协议，2：表示不遵守 coppa 协议
	if reqCtx.ReqParams.AppInfo.App.Coppa == 1 {
		dnt = true
	}
	// 若sdk传来的dnt为1，则也传给三方dsp。pioneer不使用device.Dnt，通过asinfo.dnt取值
	if reqCtx.ReqParams.Param.Dnt == "1" {
		dnt = true
	}
	device.Dnt = &dnt
	// 设置开关
	var carrier string
	useMccMnc := extractor.GetAdnetSwitchConf("useCarrier")
	if useMccMnc == 1 {
		carrier = reqCtx.ReqParams.Param.Carrier
	} else {
		carrier = reqCtx.ReqParams.Param.MCC + reqCtx.ReqParams.Param.MNC
	}
	device.Carrier = &carrier

	deviceType := rtb.BidRequest_Device_PHONE
	if reqCtx.ReqParams.Param.Platform == mvutil.IOSPLATFORM && reqCtx.ReqParams.Param.DeviceType == mvconst.DeviceTablet {
		deviceType = rtb.BidRequest_Device_TABLET
	}
	device.Devicetype = &deviceType
	device.Ip = &reqCtx.ReqParams.Param.ClientIP
	if strings.Contains(reqCtx.ReqParams.Param.ClientIP, ":") {
		device.Ipv6 = &reqCtx.ReqParams.Param.ClientIP
	}
	device.Make = &reqCtx.ReqParams.Param.ReplaceBrand
	device.Model = &reqCtx.ReqParams.Param.ReplaceModel
	os := mvconst.GetPlatformStr(reqCtx.ReqParams.Param.Platform)
	device.Os = &os
	device.Osv = &reqCtx.ReqParams.Param.OSVersion

	device.Language = &reqCtx.ReqParams.Param.Language

	device.Connectiontype = getConnetType(reqCtx.ReqParams.Param.NetworkType)

	if reqCtx.ReqParams.Param.Platform == mvconst.PlatformAndroid {
		device.Ifa = &reqCtx.ReqParams.Param.GAID
	}
	if reqCtx.ReqParams.Param.Platform == mvconst.PlatformIOS {
		device.Ifa = &reqCtx.ReqParams.Param.IDFA
	}
	device.Oaid = &reqCtx.ReqParams.Param.OAID
	if len(reqCtx.ReqParams.Param.IMEI) > 0 {
		device.Imei = &reqCtx.ReqParams.Param.IMEI
		didSha1 := mvutil.Sha1(reqCtx.ReqParams.Param.IMEI)
		device.Didsha1 = &didSha1
	}
	// imei md5 取值
	var didMd5 string
	if len(reqCtx.ReqParams.Param.ImeiMd5) > 0 {
		didMd5 = reqCtx.ReqParams.Param.ImeiMd5
	} else if len(reqCtx.ReqParams.Param.IMEI) > 0 {
		didMd5 = mvutil.Md5(reqCtx.ReqParams.Param.IMEI)
	}
	device.Didmd5 = &didMd5
	if len(reqCtx.ReqParams.Param.AndroidID) > 0 {
		device.AndroidId = &reqCtx.ReqParams.Param.AndroidID
		dpidsha1 := mvutil.Sha1(reqCtx.ReqParams.Param.AndroidID)
		device.Dpidsha1 = &dpidsha1
	}
	// Android id md5取值
	var dpidmd5 string
	if len(reqCtx.ReqParams.Param.AndroidIDMd5) > 0 {
		dpidmd5 = reqCtx.ReqParams.Param.AndroidIDMd5
	} else if len(reqCtx.ReqParams.Param.AndroidID) > 0 {
		dpidmd5 = mvutil.Md5(reqCtx.ReqParams.Param.AndroidID)
	}
	device.Dpidmd5 = &dpidmd5
	w, h := int32(reqCtx.ReqParams.Param.ScreenWidth), int32(reqCtx.ReqParams.Param.ScreenHeigh)
	device.W, device.H = &w, &h

	geo := &rtb.BidRequest_Geo{}
	geo.Country = &reqCtx.ReqParams.Param.CountryCode
	// city := strconv.FormatInt(reqCtx.ReqParams.Param.CityCode, 10)
	// geo.City = &city
	geo.City = &reqCtx.ReqParams.Param.CityString
	device.Geo = geo

	// 补充ios的device.ext信息
	if reqCtx.ReqParams.Param.Platform == mvutil.IOSPLATFORM {
		device.Ext = &rtb.BidRequest_Device_Ext{
			OsvUpTime: &reqCtx.ReqParams.Param.OsvUpTime, // 系统更新时间
			Ram:       &reqCtx.ReqParams.Param.Ram,       // 物理内存
			Uptime:    &reqCtx.ReqParams.Param.UpTime,    // 开机时间
			//CountryCode: &reqCtx.ReqParams.Param.CountryCode, // 国家代码
			TotalMemory: &reqCtx.ReqParams.Param.TotalMemory, // 硬盘尺寸
			TimeZone:    &reqCtx.ReqParams.Param.TimeZone,    // 时区
			Ifv:         &reqCtx.ReqParams.Param.IDFV,
		}
	}

	//
	if device.Ext == nil {
		device.Ext = &rtb.BidRequest_Device_Ext{}
	}
	// add one id in device
	device.Ext.Oneid = proto.String(renderOneId(reqCtx.ReqParams))

	return device
}

func constructApp(reqCtx *mvutil.ReqCtx) *rtb.BidRequest_App {
	app := &rtb.BidRequest_App{}

	appId := strconv.FormatInt(reqCtx.ReqParams.Param.AppID, 10)
	app.Id = &appId

	app.Bundle = &reqCtx.ReqParams.AppInfo.RealPackageName
	if mvutil.IsAppwallOrMoreOffer(reqCtx.ReqParams.Param.AdType) {
		app.Bundle = &reqCtx.ReqParams.Param.AppPackageName
	}
	appName := reqCtx.ReqParams.AppInfo.App.OfficialName
	if len(appName) == 0 {
		appName = reqCtx.ReqParams.AppInfo.App.Name
	}
	app.Name = &appName

	if reqCtx.ReqParams.Param.Platform == mvconst.PlatformIOS && len(reqCtx.ReqParams.AppInfo.App.BundleId) > 0 {
		app.Storeid = &reqCtx.ReqParams.AppInfo.App.BundleId
	}

	for k, v := range reqCtx.ReqParams.AppInfo.App.IabCategoryV2 {
		if len(v) == 0 {
			app.Cat = append(app.Cat, k)
			continue
		}
		app.Cat = append(app.Cat, v...)
	}

	app.Storeurl = &reqCtx.ReqParams.AppInfo.App.StoreUrl

	app.Ver = &reqCtx.ReqParams.Param.AppVersionName

	app.Publisher = constructPublisher(reqCtx.ReqParams.Param.PublisherID, reqCtx.ReqParams.PublisherInfo.Publisher.Username)
	app.Domain = &reqCtx.ReqParams.AppInfo.Publisher.Domain
	return app
}

func constructImp(reqCtx *mvutil.ReqCtx) (*rtb.BidRequest_Imp, error) {
	imp := &rtb.BidRequest_Imp{}
	impId := ImpId
	imp.Id = &impId

	nativeType := getNativeType(reqCtx)

	width := int32(reqCtx.ReqParams.Param.ScreenWidth)
	height := int32(reqCtx.ReqParams.Param.ScreenHeigh)
	vWidth := reqCtx.ReqParams.Param.VideoW
	vHeight := reqCtx.ReqParams.Param.VideoH
	adType := reqCtx.ReqParams.Param.FormatAdType
	sdkVersionCode := reqCtx.ReqParams.Param.FormatSDKVersion.SDKVersionCode
	platform := reqCtx.ReqParams.Param.Platform
	path := reqCtx.ReqParams.Param.RequestPath
	imageSize := reqCtx.ReqParams.Param.ImageSize
	imageSizeId := reqCtx.ReqParams.Param.ImageSizeID
	unitId := reqCtx.ReqParams.Param.UnitID
	requestType := reqCtx.ReqParams.Param.RequestType

	if reqCtx.ReqParams.Param.SupportAdChoice {
		imp.SupportAdchoice = &reqCtx.ReqParams.Param.SupportAdChoice
	}
	if reqCtx.ReqParams.Param.SupportMoattag {
		imp.SupportMoattag = &reqCtx.ReqParams.Param.SupportMoattag
	}

	// if vWidth == 0 || vHeight == 0 { // native video的尺寸如果是0， 从unitsize中读取，重新赋值
	//	size := getSize(reqCtx.ReqParams.Param.UnitSize)
	//	vWidth = size.Width
	//	vHeight = size.Height
	// }

	var orientation int
	if mvutil.IsHbOrV3OrV5Request(path) && reqCtx.ReqParams.UnitInfo != nil { // 如果请求来自SDK，读取portal orientaion
		orientation = reqCtx.ReqParams.UnitInfo.Unit.Orientation
	}
	if orientation == mvconst.ORIENTATION_BOTH {
		orientation = FixedLandscape4Both(reqCtx.ReqParams.Param.FormatOrientation)
	}
	orien := rtb.BidRequest_Imp_OrientationType(orientation)

	isSupportMraid := IsSupportMraid(reqCtx) // http://confluence.mobvista.com/pages/viewpage.action?pageId=24527937
	// sdk interactive ads fill banner
	if mvutil.IsHbOrV3OrV5Request(path) && (reqCtx.ReqParams.Param.AdType == mvconst.ADTypeInteractive) {
		imp.Banner = constructBanner4IA(width, height, &orien, isSupportMraid)
		// interactive 只投playable
		if reqCtx.ReqParams.Param.AdType == mvconst.ADTypeInteractive && imp.Banner != nil && imp.Banner.Ext != nil {
			t := true
			imp.Banner.Ext.Bannerplayableonly = &t
		}
		instl := int32(1)
		imp.Instl = &instl
	}

	if mvutil.IsHbOrV3OrV5OrOnlineAPIRequest(path, requestType) &&
		mvutil.IsBannerOrSdkBannerOrSplash(reqCtx.ReqParams.Param.AdType) {
		imp.Banner = constructSdkBanner(reqCtx, isSupportMraid)
		if reqCtx.ReqParams.Param.AdType == mvconst.ADTypeSplash {
			instl := int32(1)
			imp.Instl = &instl
		}
	}

	// 对于SDK和JS SDK，IV和RV包含该imp.video
	if (mvutil.IsHbOrV3OrV5Request(path) || path == mvconst.PATHJssdkApi) &&
		(adType == mvconst.ADTypeRewardVideo || adType == mvconst.ADTypeInterstitialVideo) {
		imp.Video = constructVideoWithCompainion(adType, width, height, sdkVersionCode, platform, orientation, path, isSupportMraid)
		var instl int32 = 1
		imp.Instl = &instl // 1 = the ad is interstitial or full screen, 0 = not interstitial.
		// rv/iv 增加一个banner
		if extractor.GetSupportVideoBanner() && (adType == mvconst.ADTypeRewardVideo || adType == mvconst.ADTypeInterstitialVideo) && isSupportMraid {
			imp.Banner = constructMraidBanner(width, height, isSupportMraid, orientation)
		}
	}

	// native广告, online video也使用native
	if adType == mvconst.ADTypeNativeVideo ||
		adType == mvconst.ADTypeNativePic ||
		adType == mvconst.ADTypeOnlineVideo ||
		adType == mvconst.ADTypeNativeH5 {
		native, err := constructNative(imageSize, path, requestType, imageSizeId, orientation, vWidth, vHeight, adType, nativeType)
		if err != nil {
			return nil, err
		}
		imp.Native = native
	}
	// onlineAPI native 或者more_offer及appwall
	if (requestType == mvconst.REQUEST_TYPE_OPENAPI_AD && adType == mvconst.ADTypeBanner) || mvutil.IsAppwallOrMoreOffer(adType) {
		imp.Banner = constructSdkBanner(reqCtx, false)
	}

	tagid := strconv.FormatInt(unitId, 10)
	imp.Tagid = &tagid

	bidFloor := 0.0
	if extractor.GetFillRateEcpmFloorSwitch() && !reqCtx.ReqParams.IsHBRequest {
		if reqCtx.ReqParams.Param.BidFloor > 0 && bidFloor <= reqCtx.ReqParams.Param.BidFloor {
			bidFloor = reqCtx.ReqParams.Param.BidFloor
		} else {
			bidFloor = reqCtx.ReqParams.Param.FillEcpmFloor
			bidFloorType := int32(1)
			imp.BidFloorType = &bidFloorType
		}
	}

	if reqCtx.ReqParams.IsHBRequest {
		isHB := true
		imp.IsHb = &isHB
		renderBidFloor(reqCtx)
		bidFloor = reqCtx.ReqParams.Param.BidFloor
	}

	// 针对app+设备对修改其底价，用于抓包排查问题
	bidFloor = changeDebugBidFloor(reqCtx, bidFloor)

	// 过滤小于0的底价
	if bidFloor < 0.0 {
		bidFloor = 0.01
	}
	reqCtx.ReqParams.Param.AdxBidFloor = bidFloor

	// 美元转美分
	bidFloorCent := bidFloor * 100
	imp.Bidfloor = &bidFloorCent

	bidFloorCur := "USD"
	imp.Bidfloorcur = &bidFloorCur

	secure := NON_SECURE
	if reqCtx.ReqParams.Param.HTTPReq == 2 {
		secure = SECURE
	}
	imp.Secure = &secure
	imp.Displaymanagerver = &reqCtx.ReqParams.Param.SDKVersion
	displaymanager := "mintegral SDK"
	imp.Displaymanager = &displaymanager
	imp.ExtChannel = &reqCtx.ReqParams.Param.Extchannel
	imp.SupportTrackingTemplate = &reqCtx.ReqParams.Param.SupportTrackingTemplate
	imp.Ext = constructImpExt(reqCtx, nativeType)

	return imp, nil
}

func changeDebugBidFloor(in *mvutil.ReqCtx, bidfloor float64) float64 {
	debugBidFloorAndBidPriceConf := extractor.GetDEBUG_BID_FLOOR_AND_BID_PRICE_CONF()
	if len(debugBidFloorAndBidPriceConf) == 0 {
		return bidfloor
	}
	key := mvutil.GetDebugBidFloorAndBidPriceKey(&in.ReqParams.Param)
	if conf, ok := debugBidFloorAndBidPriceConf[key]; ok && conf != nil {
		return conf.DebugBidFloor
	}
	return bidfloor
}

func renderBidFloor(in *mvutil.ReqCtx) {
	if in.ReqParams.Param.BidFloor > 0 {
		return
	}
	conf, ok := fillrateConfig(in)
	if !ok {
		return
	}
	if conf.ControlMode == 2 {
		in.ReqParams.Param.BidFloor = conf.EcpmFloor
	}
}

// oneId 渲染逻辑
func renderOneId(r *mvutil.RequestParams) string {
	//
	idGen := deviceId.NewIdGenerator()
	//
	return idGen.Generate(&deviceId_algo.Request{
		Platform:     r.Param.PlatformName,
		Gaid:         r.Param.GAID,
		GaidMd5:      r.Param.GAIDMd5,
		Imei:         r.Param.IMEI,
		ImeiMd5:      r.Param.ImeiMd5,
		Oaid:         r.Param.OAID,
		AndroidId:    r.Param.AndroidID,
		AndroidIdMd5: r.Param.AndroidIDMd5,
		Idfa:         r.Param.IDFA,
		IdfaMd5:      r.Param.IDFAMd5,
		Idfv:         r.Param.IDFV,
		SysId:        r.Param.SysId,
		BackupId:     r.Param.BkupId,
		IpV4:         r.Param.ClientIP,
		Ua:           r.Param.UserAgent,
	})
}

// abtest Exp 相关参数渲染逻辑
func renderAbTestParam(r *mvutil.RequestParams) *abtest.ABTestContext {

	var ipv4 string
	if !strings.Contains(r.Param.ClientIP, ":") {
		// ipv4
		ipv4 = r.Param.ClientIP
	}

	gen := biz.DeviceIdGenerator(&biz.RequestInfo{
		Platform:  r.Param.PlatformName,
		Gaid:      r.Param.GAID,
		Oaid:      r.Param.OAID,
		Imei:      r.Param.IMEI,
		DevId:     "",
		SysId:     r.Param.SysId,
		BackupId:  r.Param.BkupId,
		AndroidId: r.Param.AndroidID,
		Idfa:      r.Param.IDFA,
		Idfv:      r.Param.IDFV,
		RequestId: r.Param.RequestID,
		Ruid:      r.Param.RuId,
		IP:        ipv4, // 这里传ipv4
		UA:        r.Param.UserAgent,
	})
	resp, err := algoab.GetExperiments(gen, nil)
	if err != nil {
		mvutil.Logger.Runtime.Warnf("render ab test param error: [%s]", err.Error())
		return nil
	}
	r.Param.ExtDataInit.ExpIds = resp.Experiments.ExpIdList
	// 透传新增结构体
	abTestContext := &abtest.ABTestContext{
		AbtestId:  proto.String(gen.Id()),
		ExpidList: proto.String(resp.Experiments.ExpIdList),
	}
	// 拼装
	exp := make(map[string]*abtest.ABTestExperiment)
	for k, v := range resp.Experiments.Experiments {
		exp[k] = &abtest.ABTestExperiment{
			Expid:  proto.Int32(v.ExpId),
			Xpath:  proto.String(v.Xpath),
			Params: v.Params,
		}
	}
	abTestContext.Experiments = exp
	return abTestContext
}

func constructImageAssert(id, width, height int32, assertType rtb.NativeRequest_Asset_Image_ImageAssetType) *rtb.NativeRequest_Asset {
	required := int32(1)
	assert := &rtb.NativeRequest_Asset{
		Img: &rtb.NativeRequest_Asset_Image{
			Type:  &assertType,
			W:     &width,
			H:     &height,
			Mimes: []string{"image/jpeg", "image/jpg", "image/png"},
		},
		Id:       &id,
		Required: &required,
	}
	return assert
}

func constructDataAssert(id int32, assertType rtb.NativeRequest_Asset_Data_DataAssetType, required bool) *rtb.NativeRequest_Asset {
	var requiredI32 int32
	if required {
		requiredI32 = 1
	}
	assert := &rtb.NativeRequest_Asset{
		Id:       &id,
		Required: &requiredI32,
		Data: &rtb.NativeRequest_Asset_Data{
			Type: &assertType,
		},
	}
	return assert
}

func constructNative(imageSize, path string, requestType, imageSizeId, orientation int, vWidth, vHeight, adType int32, nativeType int64) (*rtb.BidRequest_Imp_Native, error) {
	nativeVer := "1.1"
	request, err := constructNativeRequest(nativeVer, imageSize, path, requestType, imageSizeId, orientation, vWidth, vHeight, adType, nativeType)
	if err != nil {
		return nil, err
	}
	requestData := string(request)
	native := &rtb.BidRequest_Imp_Native{
		Request: &requestData,
		Ver:     &nativeVer,
	}
	// 针对native h5处理逻辑
	if adType == mvconst.ADTypeNativeH5 {
		ext := &rtb.BidRequest_Imp_Native_Ext{}
		nativeType := rtb.BidRequest_Imp_Native_Ext_NATIVE_H5
		ext.Nativetype = &nativeType
		native.Ext = ext
		native.Battr = append(native.Battr, rtb.CreativeAttribute_ADOBE_FLASH)
		native.Api = append(native.Api, rtb.APIFramework_MRAID_1, rtb.APIFramework_MRAID_2, rtb.APIFramework_MRAID_3)
	}
	native.Api = append(native.Api, rtb.APIFramework_OMID_1)

	return native, nil
}

func constructNativeVideo(width, height int32, orientation int) *rtb.BidRequest_Imp_Video {
	enOrientation := rtb.BidRequest_Imp_OrientationType(orientation)
	video := &rtb.BidRequest_Imp_Video{
		W:     &width,
		H:     &height,
		Mimes: []string{"video/mp4"},
		Ext: &rtb.BidRequest_Imp_Video_Ext{
			Orientation: &enOrientation,
		},
	}

	return video
}

func constructNativeRequest(ver, imageSize, path string, requestType, imageSizeId, orientation int,
	vWidth, vHeight, adType int32, nativeType int64) ([]byte, error) {
	nativeReq := &rtb.NativeRequest{}
	nativeReq.Ver = &ver

	mWidth := int32(1200)
	mHeight := int32(627)
	if requestType == mvconst.REQUEST_TYPE_OPENAPI_AD && imageSizeId > 0 && len(imageSize) > 3 {
		w, h := getWidthAndHeight(imageSize)
		swidth, err := strconv.Atoi(w)
		if err == nil {
			mWidth = int32(swidth)
		}
		sheight, err := strconv.Atoi(h)
		if err == nil {
			mHeight = int32(sheight)
		}
	}
	mainImage := constructImageAssert(adxconst.IMAGEID, mWidth, mHeight, rtb.NativeRequest_Asset_Image_MAIN)
	// native h5参考使用banner的模版
	if adType == mvconst.ADTypeNativeH5 {
		mainImage.Img.Mimes = []string{
			"application/javascript", "image/jpeg", "image/jpg", "text/html", "image/png", "text/css", "image/gif",
		}
	}
	nativeReq.Assets = append(nativeReq.Assets, mainImage)

	// icon: onlineAPI流量使用128x128, nv和ni使用300x300
	if requestType == mvconst.REQUEST_TYPE_OPENAPI_AD {
		nativeReq.Assets = append(nativeReq.Assets, constructImageAssert(adxconst.ICONID, int32(128), int32(128), rtb.NativeRequest_Asset_Image_ICON))
	} else {
		nativeReq.Assets = append(nativeReq.Assets, constructImageAssert(adxconst.ICONID, int32(300), int32(300), rtb.NativeRequest_Asset_Image_ICON))
	}

	required := int32(1)
	title := &rtb.NativeRequest_Asset_Title{}
	titleLen := int32(20)
	title.Len = &titleLen
	titleAssert := &rtb.NativeRequest_Asset{}
	titleAssert.Title = title
	titleId := adxconst.TILTEID
	titleAssert.Id = &titleId
	titleAssert.Required = &required
	nativeReq.Assets = append(nativeReq.Assets, titleAssert)

	descAssert := constructDataAssert(adxconst.DECSID, rtb.NativeRequest_Asset_Data_DESC, true)
	nativeReq.Assets = append(nativeReq.Assets, descAssert)

	ratingAssert := constructDataAssert(adxconst.RATINGID, rtb.NativeRequest_Asset_Data_RATING, true)
	nativeReq.Assets = append(nativeReq.Assets, ratingAssert)

	ctaAssert := constructDataAssert(adxconst.CTAID, rtb.NativeRequest_Asset_Data_CTATEXT, true)
	nativeReq.Assets = append(nativeReq.Assets, ctaAssert)

	// 旧逻辑生成的adType==NV并不能完全过滤到支持video的,加上新版本的过滤 (等价于 nType == Video || nType == DisplayVideo || adType == onlinevideo )
	if (adType == mvconst.ADTypeNativeVideo && nativeType != adxconst.NativeTypeDisplay) || adType == mvconst.ADTypeOnlineVideo {
		if mvutil.IsHbOrV3OrV5Request(path) { // SDK， 尺寸固定为 1280x720 / 720x1280
			if orientation == mvconst.ORIENTATION_PORTRAIT {
				vWidth = 720
				vHeight = 1280
			} else {
				vWidth = 1280
				vHeight = 720
			}
		}
		// onlineapi, 如果有传videow videoh,就直接用。如果没传，按orientation传固定值
		if (requestType == mvconst.REQUEST_TYPE_OPENAPI_AD || path == mvconst.PATHBidAds) && (vHeight == 0 || vWidth == 0) {
			if orientation == mvconst.ORIENTATION_PORTRAIT {
				vWidth = 720
				vHeight = 1280
			} else {
				vWidth = 1280
				vHeight = 720
			}
		}
		videoAssert := &rtb.NativeRequest_Asset{}
		// TODO
		videoAssert.Video = constructNativeVideo(vWidth, vHeight, orientation)
		videoId := adxconst.VIDEOID
		videoAssert.Id = &videoId
		videoAssert.Required = &required
		nativeReq.Assets = append(nativeReq.Assets, videoAssert)
	}
	return jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(nativeReq)
}

func (backend *MAdxBackend) parseHttpResponse(reqCtx *mvutil.ReqCtx, backendCtx *mvutil.BackendCtx) (int, error) {
	if reqCtx == nil || backendCtx == nil {
		return ERR_Param, errors.New("MAdxBackend parseHttpResponse params invalidate")
	}
	// 调试信息
	// resp := new(rtb.BidResponse)
	// error := proto.Unmarshal(backendCtx.RespData, resp)
	// fmt.Println("respjsonlog:", error)
	// b, _ := jsoniter.ConfigFastest.Marshal(resp)
	// fmt.Println("respjson:", string(b))

	var rdata rtb.BidResponse
	if reqCtx.ReqParams.IsTopon {
		oResp := new(openrtb.BidResponse)
		err := proto.Unmarshal(backendCtx.RespData, oResp)
		if err != nil {
			return ERR_ParseRsp, err
		}
		reqCtx.ReqParams.ToponResponse = oResp
		resp, err := protocol.Openrtbv2ToMtgrtbResp(backendCtx.RespData)
		if err != nil {
			watcher.AddWatchValue("mAdx_unmarshal_error", float64(1))
			return ERR_ParseRsp, err
		}
		if len(oResp.Seatbid) > 0 && len(oResp.Seatbid[0].Bid) > 0 {
			ext := oResp.Seatbid[0].Bid[0].Mvext
			if ext != nil {
				reqCtx.ReqParams.DspExt = ext.Dataext
				reqCtx.ReqParams.Param.DspExt = ext.Dataext
			}
			// 获取online api bid price
			reqCtx.ReqParams.Param.OnlineApiBidPrice = strconv.FormatFloat(oResp.Seatbid[0].Bid[0].Price, 'f', 8, 64)
		}
		rdata = *resp
	} else {
		err := proto.Unmarshal(backendCtx.RespData, &rdata)
		if err != nil {
			watcher.AddWatchValue("mAdx_unmarshal_error", float64(1))
			return ERR_ParseRsp, err
		}
	}
	// 解析返回的协议字段
	GetInfoFromRData(rdata, reqCtx)

	// for hb bid test mode
	testAdCacheKey := strconv.Itoa(int(reqCtx.ReqParams.Param.AdType)) + ":" + strings.ToLower(reqCtx.ReqParams.Param.Os) + ":" + time.Now().Format("2006010215")
	var testAdCacheVal []byte
	var getTestAdCacheErr error
	// hb request snapshot
	hbRequestSnapshot := extractor.GetHBRequestSnapshot()
	if reqCtx.ReqParams.IsHBRequest && hbRequestSnapshot.Enable && rand.Intn(10000) < hbRequestSnapshot.Rate {
		testAdCacheVal, getTestAdCacheErr = hbreqctx.GetInstance().AdCacheClient.Get(testAdCacheKey)
		if getTestAdCacheErr != nil && len(rdata.Seatbid) > 0 {
			adBytes, _ := proto.Marshal(&rdata)
			hbreqctx.GetInstance().AdCacheClient.SetEx(testAdCacheKey, adBytes, time.Hour*1)
		}
	}
	// hb test mode bid request
	if reqCtx.ReqParams.Param.HBBidTestMode == 1 {
		testAdCacheVal, getTestAdCacheErr = hbreqctx.GetInstance().AdCacheClient.Get(testAdCacheKey)
		if getTestAdCacheErr == nil {
			proto.Unmarshal(testAdCacheVal, &rdata)
		}
	}
	if len(rdata.Seatbid) == 0 {
		watcher.AddWatchValue("mAdx_no_ads", float64(1))
		if reqCtx.ReqParams.Param.IsHitRequestBidServer == 1 {
			// 命中了请求 bid server 逻辑, 但是返回了未填充
			metrics.IncCounterWithLabelValues(32, BidBidServerWinButNoAds.Error())
		}
		return ERR_NoAds, errors.New("return no ads, result.id=" + *rdata.Id)
	}
	reqCtx.RespData = backendCtx.RespData
	err := fillAd(&rdata, reqCtx, backendCtx)
	if err != nil {
		return ERR_ParseRsp, err
	}

	// b, _ := json.Marshal(backendCtx.Ads)
	// mvutil.Logger.Runtime.Info("ads:" + string(b))
	return ERR_OK, nil
}

func GetInfoFromRData(rdata rtb.BidResponse, reqCtx *mvutil.ReqCtx) {
	// 这个是MAS返回的结构体
	if rdata.GetAsResp() != nil {
		reqCtx.ReqParams.Param.Extra3 = rdata.GetAsResp().GetExtra3()
		reqCtx.ReqParams.Param.Extalgo = rdata.GetAsResp().GetExtAlgo()
		reqCtx.ReqParams.Param.ExtAdxAlgo = rdata.GetAsResp().GetExtAdxAlgo()
		reqCtx.ReqParams.Param.RespFillEcpmFloor = rdata.GetAsResp().GetEcpmFloor()
		// reqCtx.ReqParams.Param.ExtifLowerImp = strconv.Atoi(rdata.GetAsResp().GetExtIfLowerImp())
		// debugmode情况下，记录as返回的过滤信息
		if reqCtx.ReqParams.Param.DebugMode {
			reqCtx.DebugModeInfo = append(reqCtx.DebugModeInfo, rdata.GetAsResp().GetDebugInfo())
		}
	}

	// b, _ := json.Marshal(rdata)
	// fmt.Println(string(b))
	// 这个是AS返回的结构体
	if len(rdata.GetRespAs()) > 0 {
		asResp := ad_server.QueryResult_{}
		err := jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal([]byte(rdata.GetRespAs()), &asResp)
		if err == nil {
			reqCtx.ReqParams.Param.Extalgo = asResp.GetAlgoFeatInfo()
			reqCtx.ReqParams.Param.ExtifLowerImp = asResp.GetIfLowerImp()
			reqCtx.ReqParams.Param.Extra3 = asResp.GetStrategy()
			reqCtx.ReqParams.Param.ExtAdxAlgo = asResp.GetExtAdxAlgo()
			reqCtx.ReqParams.Param.RespFillEcpmFloor = asResp.GetEcpmFloor()
			// debugmode情况下，记录as返回的过滤信息
			if reqCtx.ReqParams.Param.DebugMode {
				reqCtx.DebugModeInfo = append(reqCtx.DebugModeInfo, asResp.GetDebugInfo())
			}
		}
	}

	if rdata.GetExt() != nil && len(rdata.GetExt().GetAbtest()) > 0 { // 将返回的Abtest转成Map
		tmpObj := make(map[string]map[string]string)
		for _, testObj := range rdata.GetExt().GetAbtest() {
			if testObj.GetKey() == BidServerAdxABTestKey {
				// 如果有这个标记, 那意味着一定是bid server(非回退) 胜出给的广告 (madx的逻辑)
				resp := new(smodel.BidServerCtx)
				err := jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal([]byte(testObj.GetValue()), resp)
				if err != nil {
					mvutil.Logger.Runtime.Warnf("[bid-server] get adx resp key by unmarshal error:[%s]", err.Error())
					continue
				}
				// 响应中的 BidResponse 没用, 置为空
				reqCtx.ReqParams.Param.BidServerCtx = resp

				// 记录adx原始的响应, 用于之后Load数据补充adx相关信息
				reqCtx.ReqParams.Param.BidServerAdxResponse = &rdata
				continue
			}
			if len(testObj.GetKey()) > 0 {
				if len(tmpObj[testObj.GetKey()]) == 0 {
					tmpObj[testObj.GetKey()] = make(map[string]string)
				}
				tmpObj[testObj.GetKey()]["v"] = testObj.GetValue()
				tmpObj[testObj.GetKey()]["t"] = testObj.GetTags()
			}
		}
		reqCtx.ReqParams.Param.ExtDataInit.AdxAbTest = tmpObj
	}

	// 获取adx启动模式（是否灰度）
	if len(rdata.GetExt().GetRepGrayTags()) > 0 {
		reqCtx.ReqParams.Param.AdxStartMode = rdata.GetExt().GetRepGrayTags()
	}
	reqCtx.ReqParams.Param.RespOnlyRequestThirdDsp = rdata.GetExt().GetOnlyRequstThirdDsp()
	reqCtx.ReqParams.Param.MasResponseTime = rdata.GetTimestamp()
}

func getBidByImpId(seatbids []*rtb.BidResponse_SeatBid, impId string) (*rtb.BidResponse_SeatBid_Bid, error) {
	if len(seatbids) == 0 {
		return nil, errors.New("get BidByImpId seatbids is empty")
	}
	for _, seatbid := range seatbids {
		for _, bid := range seatbid.GetBid() {
			if bid.GetImpid() == impId {
				return bid, nil
			}
		}
	}
	return nil, errors.New("get BidByImpId has no bid match impId:" + impId)
}

func doFillAd(ad *corsair_proto.Campaign, bid *rtb.BidResponse_SeatBid_Bid, reqCtx *mvutil.ReqCtx, backendCtx *mvutil.BackendCtx, k int) error {
	if ad == nil || bid == nil || reqCtx == nil || reqCtx.ReqParams == nil || backendCtx == nil {
		return errors.New("doFillAd input is invalidate")
	}
	reqCtx.ReqParams.DspExt = bid.GetExt().GetDataext()
	dspExt, err := reqCtx.ReqParams.GetDspExt()
	if err != nil {
		return errors.New("doFillAd decodeDspExt error:" + err.Error())
	}
	// 对于第三方的返回， 会在原adid的基础上加一个大数，以免和内部的campaignId冲突
	if dspExt.DspId == mvconst.MAS {
		// 返回的adid格式为 dspId-campaignId
		adid := bid.GetAdid()
		index := strings.Index(adid, "-")
		if index <= 0 {
			return errors.New("doFillAd invalid mas adid:" + adid)
		}
		ad.CampaignId = adid[index+1:]
	} else {
		OfferIDBase := mvconst.MaxMVOfferID + (mvconst.MAdx-1)*mvconst.BACKENDLEN
		if len(bid.GetAdid()) > 0 {
			adId, err := strconv.ParseInt(bid.GetAdid(), 10, 64)
			if err != nil {
				// 非数字字符串，直接用crc32做转换
				adId = int64(crc32.ChecksumIEEE([]byte(bid.GetAdid())))
			}
			ad.CampaignId = strconv.FormatInt(int64((adId)%mvconst.BACKENDLEN)+int64(OfferIDBase), 10)
		}
	}
	*ad.PackageName = bid.GetBundle()
	*ad.Price = bid.GetPrice()
	// 获取主素材id及click_url
	if mvutil.IsBannerOrSdkBannerOrSplashOrNativeH5(reqCtx.ReqParams.Param.AdType) {
		if bid.Crid != nil {
			arr := strings.Split(*bid.Crid, "-")
			if len(arr) == 2 {
				crid := arr[1]
				ad.CreativeId = &crid
			}
		}
		if len(bid.GetExt().GetClickUrl()) > 0 {
			clickUrl := bid.GetExt().GetClickUrl()
			ad.ClickURL = &clickUrl
		}
		// 获取mv dsp原有的campaignid，用于判断clickmode取值判断
		oriCid := bid.GetCid()
		cidArr := strings.Split(oriCid, "-")
		if len(cidArr) == 2 {
			oriCid = cidArr[1]
		}
		ad.OriCampaignId = &oriCid
	}

	if dspExt.DspId == mvconst.FakeAdserverDsp {
		reqCtx.ReqParams.IsFakeAs = true
		asResp := ad_server.QueryResult_{}
		err := jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal([]byte(bid.GetAdm()), &asResp)
		if err != nil {
			return errors.New("do FillAd parse As Response error:" + err.Error())
		}
		reqCtx.ReqParams.Param.RespFillEcpmFloor = asResp.GetEcpmFloor()
		FillAsAds(reqCtx, backendCtx, &asResp)
		if len(backendCtx.Ads.CampaignList) > 0 && bid.Ext != nil && len(bid.Ext.GetAks()) > 0 {
			for index := range backendCtx.Ads.CampaignList {
				backendCtx.Ads.CampaignList[index].AKS = bid.Ext.GetAks()
			}
		}

		if backendCtx.Ads != nil && len(backendCtx.Ads.CampaignList) > 0 {
			// 限制切量大模板的情况
			if reqCtx.ReqParams.Param.BigTemplateFlag && backendCtx.Ads.BigTemplateInfo != nil {
				for _, ad := range backendCtx.Ads.CampaignList {
					bidPrice := bid.GetPrice()
					ad.Price = &bidPrice
					ad.AdTracking = corsair_proto.NewAdTracking()
					for _, rate := range bid.GetExt().Rates {
						ad.AdTracking.PlayPerct = append(ad.AdTracking.PlayPerct, &corsair_proto.PlayPercent{
							Rate: rate.GetRate(),
							URL:  rate.GetUrl(),
						})
					}
					if len(bid.GetClicktracker()) > 0 {
						ad.AdTracking.Click = append(ad.AdTracking.Click, bid.GetClicktracker())
					}
				}
				reqCtx.ReqParams.Param.BigTempalteAdxPvUrl = append(reqCtx.ReqParams.Param.BigTempalteAdxPvUrl, bid.GetImptracker())
			} else {
				ad := backendCtx.Ads.CampaignList[0]
				bidPrice := bid.GetPrice()
				ad.Price = &bidPrice
				ad.AdTracking = corsair_proto.NewAdTracking()
				for _, rate := range bid.GetExt().Rates {
					ad.AdTracking.PlayPerct = append(ad.AdTracking.PlayPerct, &corsair_proto.PlayPercent{
						Rate: rate.GetRate(),
						URL:  rate.GetUrl(),
					})
				}
				if len(bid.GetImptracker()) > 0 {
					ad.AdTracking.Impression = append(ad.AdTracking.Impression, bid.GetImptracker())
				}
				if len(bid.GetClicktracker()) > 0 {
					ad.AdTracking.Click = append(ad.AdTracking.Click, bid.GetClicktracker())
				}
			}

			reqCtx.ReqParams.Param.DspExt = reqCtx.ReqParams.DspExt
		}
		return nil
	}

	var adTracking corsair_proto.AdTracking
	if bid.Ext != nil {
		// 如果 Linktype==0 用默认值
		if bid.GetExt().GetLinktype() > 0 {
			*ad.LinkType = int32(bid.GetExt().GetLinktype())
		}

		if len(bid.Ext.Imptrackers) > 0 {
			adTracking.Impression = append(adTracking.Impression, bid.GetExt().GetImptrackers()...)
			adTracking.Click = append(adTracking.Click, bid.GetExt().GetClicktrackers()...)
		}

		if len(bid.Ext.GetRates()) > 0 {
			for _, rate := range bid.Ext.GetRates() {
				adTracking.PlayPerct = append(adTracking.PlayPerct, &corsair_proto.PlayPercent{
					Rate: rate.GetRate(),
					URL:  rate.GetUrl(),
				})
			}
		}
		deeplink := bid.Ext.GetDeeplink()
		if len(deeplink) > 0 {
			ad.DeepLink = &deeplink
		}
	}
	// 除as和mvdsp外， madx返回的第三方DSP resp
	// IOS appstore 的单子模板自动跳转
	if dspExt.DspId != mvconst.FakeAdserverDsp && dspExt.DspId != mvconst.MVDSP && dspExt.DspId != mvconst.MVDSP_Retarget && dspExt.DspId != mvconst.MAS &&
		reqCtx.ReqParams.Param.Platform == mvconst.PlatformIOS && bid.GetBundle() != "" && ad.DeepLink == nil {
		*ad.LinkType = int32(rtb.LinkType_APPSTORE)
	}
	formatAdType := reqCtx.ReqParams.Param.FormatAdType
	isAdmVast := false
	// MAS的banner和native都是VAST
	if mvutil.IsRequestPioneerDirectly(&reqCtx.ReqParams.Param) ||
		mvutil.IsHbOrV3OrV5OrOnlineAPIRequest(reqCtx.ReqParams.Param.RequestPath, reqCtx.ReqParams.Param.RequestType) &&
			(mvutil.IsIvOrRv(formatAdType) || formatAdType == mvconst.ADTypeNativeH5 ||
				(dspExt.DspId == mvconst.MAS &&
					(mvutil.IsBannerOrSdkBannerOrSplash(reqCtx.ReqParams.Param.FormatAdType) || mvutil.IsNative(formatAdType) || formatAdType == mvconst.ADTypeOnlineVideo)) ||
				(extractor.IsVastBannerDsp(dspExt.DspId) && mvutil.IsBannerOrSplashOrDI(reqCtx.ReqParams.Param.FormatAdType))) {
		adm := bid.GetAdm()
		isAdmVast = mvutil.IsAdmVast(adm)
		if isAdmVast {
			vastData, err := decodeVast(adm)
			if err != nil {
				return errors.New("doFillAd decodeVast error:" + err.Error())
			}
			if len(vastData.Ads) == 0 {
				return errors.New("doFillAd decodeVast ads is empty")
			}
			var vastAd vast.Ad
			// mas 因为大模板的原因，vast里支持多个offer
			vastAd = vastData.Ads[k]
			err = filladWithVast(ad, &adTracking, &vastAd, reqCtx, dspExt.DspId, true)
			if err != nil {
				return err
			}
			if mvutil.IsBannerOrSdkBannerOrSplashOrNativeH5(reqCtx.ReqParams.Param.FormatAdType) {
				fillAdWithBanner(&adTracking, bid, backendCtx, dspExt.DspId, ad, reqCtx)
			}
		} else {
			fillAdWithBanner(&adTracking, bid, backendCtx, dspExt.DspId, ad, reqCtx)
		}
	} else if mvutil.IsBannerOrSdkBannerOrSplashOrNativeH5(reqCtx.ReqParams.Param.FormatAdType) {
		// 若为sdkbanner 则不使用adm做任何处理
		fillAdWithBanner(&adTracking, bid, backendCtx, dspExt.DspId, ad, reqCtx)
	} else {
		// native, online video
		responseNative, err := decodeNative(bid.GetAdm())
		if err != nil {
			return errors.New("doFillAd decodeNative error:" + err.Error())
		}
		err = filladWithNative(ad, &adTracking, &responseNative.Native, reqCtx, dspExt.DspId)
		if err != nil {
			return err
		}
	}

	if reqCtx.ReqParams.Param.HtmlSupport == 1 && mvutil.IsThirdDspWithRtdsp(dspExt.DspId) {
		adm := bid.GetAdm()
		ad.AdHtml = &adm
	}

	if mvutil.IsThirdDspWithRtdsp(dspExt.DspId) { // endcard的展示方式
		if ad.ClickURL != nil && *ad.ClickURL != "" {
			*ad.VideoEndType = int32(ad_server.VideoEndType_NORMAL_ENDCARD)
		} else {
			*ad.VideoEndType = int32(ad_server.VideoEndType_CLOSE_VIDEO)
		}
	}

	if reqCtx.ReqParams.Param.Platform == mvconst.PlatformIOS &&
		dspExt.DspId == mvconst.GDTYLH && ad.GetLinkType() == 1 &&
		isHitClickModeTest(reqCtx.ReqParams.Param.IDFA, reqCtx.ReqParams.Param.IDFV, reqCtx.ReqParams.Param.UnitID, reqCtx.ReqParams.Param.AppID) {
		// 命中实验
		*ad.ClickMode = 6
		reqCtx.ReqParams.Param.YLHHit = 1
	}
	// 解析
	// if fk, ok := backendCtx.FakeKeys[int64(dspExt.DspId)]; ok {
	// 	priceFactor := strconv.Itoa(mvconst.MAdx) + ":" + strconv.FormatInt(dspExt.DspId, 10) + ":" + strconv.FormatFloat(fk.PriceFactor, 'f', 2, 64)
	// 	reqCtx.ReqParams.PriceFactor.Store(strconv.Itoa(mvconst.MAdx), priceFactor)
	// 	reqCtx.ReqParams.ReqPriceFactor = append(reqCtx.ReqParams.ReqPriceFactor, priceFactor)
	// }

	adTypeStr := getAdTypeDesc(reqCtx.ReqParams.Param.FormatAdType)

	adReqKeyName := backendCtx.AdReqKeyName
	if len(adTypeStr) > 0 {
		// add adnet playPerct urls
		extData := ""
		// if dspExt.DspId == mvconst.FakeToutiao || dspExt.DspId == mvconst.FakeGDT {
		// 	if fakeData, ok := backendCtx.FakeKeys[dspExt.DspId]; ok {
		// 		adReqKeyName = fakeData.ReqKeyName
		// 		backendCtx.AdReqKeyName = adReqKeyName
		// 	}

		// 	extData = url.QueryEscape(fakeExt)
		// 	// orietation := mvconst.ORIENTATION_LANDSCAPE
		// 	// if bid.GetImageMode() == 15 {
		// 	// 	//竖屏
		// 	// 	orietation = mvconst.ORIENTATION_PORTRAIT
		// 	// }
		// 	//template相关
		// 	// FillRVTemplate(ad, orietation, reqCtx.ReqParams.Param.Platform, mvutil.GetBackendId(dspExt.DspId), reqCtx.ReqParams.Param.FormatAdType)

		// 	if ad.Rv != nil && ad.Rv.Template != nil {
		// 		reqCtx.ReqParams.Param.MwCreatvie = append(reqCtx.ReqParams.Param.MwCreatvie, *ad.Rv.Template)
		// 	}
		// }
		if dspExt.DspId != mvconst.MAS {
			backendData := formatBackendData(mvconst.MAdx, *ad.VideoURL, adReqKeyName)
			prefix := GetAdTrackURLPrefix(reqCtx.ReqParams, mvconst.MAdx, ad.CampaignId, backendData)
			perct100 := formatPercentageURL(PER_RATE_100, prefix, adTypeStr, extData)
			if dspExt.DspId == mvconst.GDTYLH {
				*ad.ClickURL = renderGdtClickUrl(*ad.ClickURL, reqCtx.ReqParams)
				*ad.FCB = 3
			}
			// 视频才会填充进度tracking
			if formatAdType == mvconst.ADTypeOnlineVideo || reqCtx.IsNativeVideo ||
				(mvutil.IsHbOrV3OrV5Request(reqCtx.ReqParams.Param.RequestPath) &&
					(formatAdType == mvconst.ADTypeRewardVideo || formatAdType == mvconst.ADTypeInterstitialVideo)) {
				adTracking.PlayPerct = append(adTracking.PlayPerct, formatPercentageURL(PER_RATE_0, prefix, adTypeStr, extData))
				adTracking.PlayPerct = append(adTracking.PlayPerct, formatPercentageURL(PER_RATE_25, prefix, adTypeStr, ""))
				adTracking.PlayPerct = append(adTracking.PlayPerct, formatPercentageURL(PER_RATE_50, prefix, adTypeStr, ""))
				adTracking.PlayPerct = append(adTracking.PlayPerct, formatPercentageURL(PER_RATE_75, prefix, adTypeStr, ""))
				adTracking.PlayPerct = append(adTracking.PlayPerct, perct100)
			}
			// 新pv统计逻辑
			if reqCtx.ReqParams.Param.ApiVersion >= mvconst.API_VERSION_1_5 && reqCtx.ReqParams.Param.ApiVersion < mvconst.API_VERSION_1_9 && mvutil.IsIvOrRv(reqCtx.ReqParams.Param.AdType) {
				adTracking.PubImp = append(adTracking.PubImp, formatAdTrackingURL(prefix, adTypeStr, "pub_imp", ""))
			}
		}
	}

	ad.AdTracking = &adTracking
	// addBackendLog(ad, &DemandExt{BackendId: mvconst.MAdx, DspId: dspExt.DspId}, &reqCtx.ReqParams.Param)
	return nil
}

// SupportVideoFeatrureIOS
// flag=1 不需要缓存 video info
// flag=2 需要缓存 video info
// flag=3 需要缓存 video info, 且视频在mtg cdn里
func SupportVideoFeatureIOS(ad *corsair_proto.Campaign, param *mvutil.Params, videoUrl, qKey string, videoLen, videoSize int32) (int, error) {
	flag := 1
	if (param.FormatSDKVersion.SDKVersionCode < mvconst.IOSSupportVideoUrlWithParams || mvutil.InInt32Arr(param.FormatSDKVersion.SDKVersionCode,
		[]int32{mvconst.IOSSupportVideoUrlWithParamsExpV1, mvconst.IOSSupportVideoUrlWithParamsExpV2, mvconst.IOSSupportVideoUrlWithParamsExpV3, mvconst.IOSSupportVideoUrlWithParamsExpV4})) &&
		!strings.HasSuffix(strings.ToLower(videoUrl), "mp4") {
		// 不支持视频文件mp4后面带参数，走素材处理流程
		flag = 3
	}

	if flag == 1 && param.FormatSDKVersion.SDKVersionCode < mvconst.IOSSupportVideoSizeZero && videoSize == int32(0) {
		flag = 2
	}

	if flag > 1 {
		vInfo, err := redis.LocalRedisGet(qKey)
		if err != nil {
			mvutil.Logger.Creative.Infof("%s\t%d\t%s", qKey, flag, videoUrl)
			return flag, errors.New("video creative not support")
		}
		var videoCacheData *mvutil.VideoCacheData
		if flag == 3 {
			videoCacheData, err = mvutil.VideoMetaDataWithUrl(vInfo)
			if err != nil || videoCacheData == nil || len(videoCacheData.VideoMTGCdnURL) == 0 {
				return flag, errors.New("video creative not support")
			}
			*ad.VideoURL = videoCacheData.VideoMTGCdnURL
		} else if flag == 2 {
			videoCacheData, err = mvutil.VideoMetaDataWithoutUrl(vInfo)
			if err != nil || videoCacheData == nil {
				return flag, errors.New("video creative not support")
			}
		}
		*ad.VideoSize = videoCacheData.VideoSize
		*ad.VideoLength = videoCacheData.VideoLen
		*ad.VideoResolution = videoCacheData.VideoResolution
	}
	return flag, nil
}

// filadWithVast
// masMoreAds: 这个vast是否是“第三方DSP胜出了，填充第一个单子， mas未胜出，用mas的vast填充剩下的单子”的情况
// masMoreAds==true 时， 将adx返回的click/imp/视频进度监测都清空
func filladWithVast(ad *corsair_proto.Campaign, adTracking *corsair_proto.AdTracking, vad *vast.Ad, reqCtx *mvutil.ReqCtx, dspId int64, masMoreAdsWin bool) error {
	if ad == nil || adTracking == nil || vad == nil || vad.InLine == nil {
		return errors.New("filladWithVast vad.Inline is nil")
	}
	if len(ad.CampaignId) == 0 || mvutil.IsRequestPioneerDirectly(&reqCtx.ReqParams.Param) {
		ad.CampaignId = vad.ID
	}
	*ad.AppName = vad.InLine.AdTitle.CDATA
	*ad.AppDesc = vad.InLine.Description.CDATA

	if len(vad.InLine.Creatives) == 0 {
		return errors.New("filladWithVast vad.InLine.Creatives is empty")
	}
	param := &reqCtx.ReqParams.Param
	var HTMLResource string
	// 先解析extension, 以获知是否是playableT
	if vad.InLine.Extensions == nil || len(*vad.InLine.Extensions) == 0 {
		return errors.New("filladWithVast vad.Inline.Extensions is empty")
	}
	extensions := *vad.InLine.Extensions
	exists := false
	for _, extension := range extensions {
		// extension := (*vad.InLine.Extensions)[0]
		if extension.Template != nil {
			oldVastFill(ad, &extension)
			exists = true
		} else if len(extension.Asset) > 0 || len(extension.Templates) > 0 {
			newVastFill(ad, &extension, HTMLResource, param, dspId)
			exists = true
		}
		if extension.Type == "AdVerifications" { // pokkt 返回的 omsdk 字段
			if extension.Verification == nil || extension.Verification.VerificationParameters == nil {
				continue
			}
			reqCtx.ReqParams.OmSDK = []mvutil.OmSDK{
				{
					VendorKey:              extension.Verification.Vendor,
					VerificationParameters: extension.Verification.VerificationParameters.Cdata,
					EventtrackerUrl:        extension.Verification.JavaScriptResource.URI,
				},
			}
		}
		// mas 的额外字段
		if extension.Orientation != nil {
			orientation, _ := strconv.ParseInt(extension.Orientation.Value, 10, 64)
			enumOrientation := ad_server.Orientation(orientation)
			ad.Orientation = &enumOrientation
		}
		if extension.VideoEndType != nil {
			videoEndType, _ := strconv.ParseInt(extension.VideoEndType.Value, 10, 64)
			videoEndType32 := int32(videoEndType)
			ad.VideoEndType = &videoEndType32
		}
		if extension.PlayableAdsWithoutVideo != nil {
			pawv, _ := strconv.ParseInt(extension.PlayableAdsWithoutVideo.Value, 10, 64)
			pawv32 := int32(pawv)
			ad.PlayableAdsWithoutVideo = &pawv32
		}
		if extension.Fcb != nil {
			fcb, _ := strconv.ParseInt(extension.Fcb.Value, 10, 64)
			fcb32 := int32(fcb)
			ad.FCB = &fcb32
		}
		// 记录offer id
		if len(vad.ID) > 0 && (dspId == mvconst.MAS || mvutil.IsRequestPioneerDirectly(&reqCtx.ReqParams.Param)) {
			ad.CampaignId = vad.ID
		}
	}
	if !exists {
		return errors.New("filladWithVast vad.Inline.Extensions.Template/Templates/Asset is empty")
	}
	isNative := reqCtx.ReqParams.Param.AdType == mvconst.ADTypeNative
	hasEndcard := ad.EndcardURL != nil && len(*ad.EndcardURL) > 0 || hasBannerTemplate(ad)
	isOnlineApiBanner := reqCtx.ReqParams.Param.RequestType == mvconst.REQUEST_TYPE_OPENAPI_AD && reqCtx.ReqParams.Param.AdType == mvconst.ADTypeBanner
	for _, vcreative := range vad.InLine.Creatives {
		if vcreative.Linear != nil { // creative.Linear
			if vcreative.Linear.Icons != nil && len(vcreative.Linear.Icons.Icon) > 0 && len(*ad.IconURL) == 0 {
				for _, icon := range vcreative.Linear.Icons.Icon {
					if icon.StaticResource != nil {
						ad.IconURL = &icon.StaticResource.URI
					}
				}
			}
			// mediafile里的视频链接和extension里的endcard至少要有一个.
			if len(vcreative.Linear.MediaFiles) == 0 && !hasEndcard && !isNative && !isOnlineApiBanner &&
				!mvutil.IsRequestPioneerDirectly(&reqCtx.ReqParams.Param) {
				return errors.New("filladWithVast vad.InLine.Creatives video&endcard is empty")
			}
			if len(vcreative.Linear.MediaFiles) > 0 {
				duration, err := vcreative.Linear.Duration.MarshalText()
				if err == nil && len(duration) > 0 {
					videoLen, err := getVideoLength(string(duration))
					if err == nil {
						*ad.VideoLength = int32(videoLen)
					}
				}
				*ad.VideoURL = strings.TrimSpace(vcreative.Linear.MediaFiles[0].URI)
				// 记录三方dsp的视频url,用于运营分析同包名的产品素材差异
				recordThirdPartyDspCreative(dspId, ad, reqCtx)
				mediaFile := vcreative.Linear.MediaFiles[0]
				// 对于pioneer，不需要从cache里获取
				if param.Platform == mvconst.PlatformIOS && dspId != mvconst.MAS {
					qKey := fmt.Sprintf("mtg_%s_%s", ad.CampaignId, mvutil.Md5(*ad.VideoURL))
					flag, err := SupportVideoFeatureIOS(ad, param, *ad.VideoURL, qKey, ad.GetVideoLength(), int32(mediaFile.Size))
					if err != nil {
						return err
					}
					if flag == 1 {
						*ad.VideoSize = int32(mediaFile.Size)
						*ad.VideoResolution = strconv.Itoa(vcreative.Linear.MediaFiles[0].Width) + "x" + strconv.Itoa(vcreative.Linear.MediaFiles[0].Height)
					}
				} else {
					*ad.VideoSize = int32(mediaFile.Size)
					*ad.VideoResolution = strconv.Itoa(vcreative.Linear.MediaFiles[0].Width) + "x" + strconv.Itoa(vcreative.Linear.MediaFiles[0].Height)
					*ad.BitRate = int32(vcreative.Linear.MediaFiles[0].Bitrate)
				}
				if reqCtx.ReqParams.Param.RequestType == mvconst.REQUEST_TYPE_OPENAPI_AD {
					if crid, err := strconv.ParseInt(vcreative.ID, 10, 64); err == nil {
						reqCtx.ReqParams.Param.VideoCreativeid = crid
					}
				}
			}

			// 对于广点通，impression中第一个为videoViewLink
			if (dspId == mvconst.GDTYLH) && len(vad.InLine.Impressions) > 0 && len(vad.InLine.Impressions[0].URI) > 0 {
				videoViewLink := vad.InLine.Impressions[0].URI
				videoViewLink = renderVideoViewLink(videoViewLink, param.FormatOrientation, *ad.VideoLength)
				// 封装进playcomplate
				playtmp := corsair_proto.NewPlayPercent()
				playtmp.Rate = 100
				playtmp.URL = videoViewLink
				adTracking.PlayPerct = append(adTracking.PlayPerct, playtmp)
				// 删除广点通impression list中第一个url
				if len(vad.InLine.Impressions) > 1 {
					vad.InLine.Impressions = append(vad.InLine.Impressions[:0], vad.InLine.Impressions[1:]...)
				}
			}
			fillAdTrackingWithVastTracking(vcreative.Linear.TrackingEvents, adTracking)

			if vcreative.Linear.VideoClicks != nil {
				if len(vcreative.Linear.VideoClicks.ClickThroughs) > 0 {
					*ad.ClickURL = strings.TrimSpace(vcreative.Linear.VideoClicks.ClickThroughs[0].URI)
				}
				if len(vcreative.Linear.VideoClicks.ClickTrackings) > 0 {
					for _, clickTracking := range vcreative.Linear.VideoClicks.ClickTrackings {
						adTracking.Click = append(adTracking.Click, strings.TrimSpace(clickTracking.URI))
					}
				}
			}
		}
		// creative.CompanionAds
		if vcreative.CompanionAds != nil {
			if len(vcreative.CompanionAds.Companions) == 0 {
				return errors.New("filladWithVast vad.InLine.Creatives companionAds is empty")
			}
			for _, companion := range vcreative.CompanionAds.Companions {
				if companion.StaticResource != nil {
					*ad.ImageURL = strings.TrimSpace(companion.StaticResource.URI)
					var imageResolution string
					if companion.Width > 0 && companion.Height > 0 {
						imageResolution = fmt.Sprintf("%dx%d", companion.Width, companion.Height)
						*ad.ImageWidth = companion.Width
						*ad.ImageHeight = companion.Height
					}
					if reqCtx.ReqParams.Param.RequestType == mvconst.REQUEST_TYPE_OPENAPI_AD { //online专用
						ad.ImageResolution = &imageResolution
						*ad.ImageMime = strings.TrimSpace(companion.StaticResource.CreativeType)
						if crid, err := strconv.ParseInt(vcreative.ID, 10, 64); err == nil {
							reqCtx.ReqParams.Param.ImageCreativeid = crid
						}
					}
				}
				if companion.HTMLResource != nil {
					*ad.HtmlTemplate = companion.HTMLResource.HTML
				}
				//if len(companion.IFrameResource.CDATA) > 0 {
				//	*ad.EndcardURL = companion.IFrameResource.CDATA
				//}
				for _, clickTracking := range companion.CompanionClickTracking {
					if len(clickTracking.CDATA) > 0 {
						adTracking.Click = append(adTracking.Click, strings.TrimSpace(clickTracking.CDATA))
					}
				}
				if len(*ad.ClickURL) == 0 {
					// 如果 Linear为空就将companion的clickthrough填到上报页
					*ad.ClickURL = strings.TrimSpace(companion.CompanionClickThrough.CDATA)
				}
				// endcardshow 从原来的 companion.CompanionClickThrough 移到 TrackingEvent.Event creativeView
				// 兼容creative.companionsAds.companions.trackingEvents放apk埋点上报链接的情况
				fillAdTrackingWithVastTracking(companion.TrackingEvents, adTracking)
			}
		}
	}
	for _, impression := range vad.InLine.Impressions {
		if reqCtx.ReqParams.Param.BigTemplateFlag && dspId == mvconst.MAS && reqCtx.ReqParams.Param.ApiVersion >= mvconst.API_VERSION_1_9 {
			reqCtx.ReqParams.Param.BigTempalteAdxPvUrl = append(reqCtx.ReqParams.Param.BigTempalteAdxPvUrl, impression.URI)
		} else if masMoreAdsWin {
			adTracking.Impression = append(adTracking.Impression, impression.URI)
		}
	}
	if !masMoreAdsWin {
		adTracking.Click = []string{}
	}
	return nil
}

func hasBannerTemplate(ad *corsair_proto.Campaign) bool {
	return (ad.HtmlTemplate != nil && len(*ad.HtmlTemplate) > 0) || (ad.UrlTemplate != nil && len(*ad.UrlTemplate) > 0)
}

// fillAdTrackingWithVastTracking 填充vast中带过来的视频进度 （adnet自己的视频进度不在这个函数内加）
func fillAdTrackingWithVastTracking(trackingEvents []vast.Tracking, adTracking *corsair_proto.AdTracking) {
	for _, tracking := range trackingEvents {
		if !strings.HasPrefix(tracking.URI, "http") && !strings.HasPrefix(tracking.URI, "{sh}") {
			continue
		}
		switch tracking.Event {
		case "start":
			adTracking.PlayPerct = append(adTracking.PlayPerct, &corsair_proto.PlayPercent{Rate: PER_RATE_0, URL: tracking.URI})
		case "firstQuartile":
			adTracking.PlayPerct = append(adTracking.PlayPerct, &corsair_proto.PlayPercent{Rate: PER_RATE_25, URL: tracking.URI})
		case "midpoint":
			adTracking.PlayPerct = append(adTracking.PlayPerct, &corsair_proto.PlayPercent{Rate: PER_RATE_50, URL: tracking.URI})
		case "thirdQuartile":
			adTracking.PlayPerct = append(adTracking.PlayPerct, &corsair_proto.PlayPercent{Rate: PER_RATE_75, URL: tracking.URI})
		case "complete":
			adTracking.PlayPerct = append(adTracking.PlayPerct, &corsair_proto.PlayPercent{Rate: PER_RATE_100, URL: tracking.URI})
		case "mute":
			adTracking.Mute = append(adTracking.Mute, tracking.URI)
		case "unmute":
			adTracking.Unmute = append(adTracking.Unmute, tracking.URI)
		case "pause":
			adTracking.Pause = append(adTracking.Pause, tracking.URI)
		case "close":
			adTracking.Close = append(adTracking.Close, tracking.URI)
		case "endcard_show":
			adTracking.EndcardShow = append(adTracking.EndcardShow, tracking.URI)
		case "apk_install":
			adTracking.ApkInstall = append(adTracking.ApkInstall, tracking.URI)
		case "pub_imp":
			adTracking.PubImp = append(adTracking.PubImp, tracking.URI)
		case "download_start":
			adTracking.Click = append(adTracking.Click, tracking.URI)
		case "apk_download_start":
			adTracking.ApkDownloadStart = append(adTracking.ApkDownloadStart, tracking.URI)
		case "apk_download_end":
			adTracking.ApkDownloadEnd = append(adTracking.ApkDownloadEnd, tracking.URI)
		case "video_click":
			adTracking.VideoClick = append(adTracking.VideoClick, tracking.URI)
		case "impression_t2":
			adTracking.ImpressionT2 = append(adTracking.ImpressionT2, tracking.URI)
		case "creativeView":
			adTracking.EndcardShow = append(adTracking.EndcardShow, strings.TrimSpace(tracking.URI))
		}
	}
}

func filladWithNative(ad *corsair_proto.Campaign, adTracking *corsair_proto.AdTracking, nad *native.NativeItem, reqCtx *mvutil.ReqCtx, dspId int64) error {
	if ad == nil || adTracking == nil || nad == nil || len(nad.Assets) == 0 {
		return errors.New("filladWithNative params is nil")
	}
	param := &reqCtx.ReqParams.Param
	if nad.Ext != nil && nad.Ext.Adchoice != nil { // adchoice
		reqCtx.ReqParams.Adchoice = nad.Ext.Adchoice
	}
	for _, asset := range nad.Assets {
		if asset.Title != nil && len(asset.Title.Text) > 0 {
			*ad.AppName = asset.Title.Text
		}
		if asset.Img != nil {
			if len(asset.Img.Url) > 0 && asset.Id == int(adxconst.IMAGEID) {
				if dspId == mvconst.Criteo { // 和头条一样，对图片进行额外处理
					imgInfo, imgKey, err := getCriteoImg(asset.Img.Url)
					if err != nil || imgInfo == nil {
						mvutil.Logger.Runtime.Warnf("request_id=[%s] GetImageUrl key=[%s] from redis error:%s",
							param.RequestID, asset.Img.Url, err.Error())
						mvutil.Logger.Creative.Infof("%s\t4\t%s\t1200\t627", imgKey, asset.Img.Url)
						return errors.New("filladWithNative criteo img not found:" + asset.Img.Url)
					} else {
						*ad.ImageSize = imgInfo.Resolution
						*ad.ImageURL = imgInfo.Url
					}
				} else {
					*ad.ImageURL = asset.Img.Url
					*ad.ImageSize = fmt.Sprintf("%dx%d", asset.Img.W, asset.Img.H)
				}
			}
			if len(asset.Img.Url) > 0 && asset.Id == int(adxconst.ICONID) {
				*ad.IconURL = asset.Img.Url
			}
		}
		if asset.Video != nil {
			if len(asset.Video.VastTag) > 0 {
				nativeVideo, err := decodeVast(asset.Video.VastTag)
				if err != nil {
					return errors.New("filladWithNative decodeVast error:" + err.Error())
				}
				if len(nativeVideo.Ads) == 0 {
					return errors.New("filladWithNative video is empty")
				}
				nvad := nativeVideo.Ads[0]
				if nvad.InLine == nil || len(nvad.InLine.Creatives) == 0 {
					return errors.New("filladWithNative nvad.InLine.Creatives is empty")
				}

				nvCreative := nvad.InLine.Creatives[0]
				if nvCreative.Linear == nil || len(nvCreative.Linear.MediaFiles) == 0 {
					return errors.New("filladWithNative nvad.InLine.Creatives video is empty")
				}
				reqCtx.IsNativeVideo = true // 返回的native确实是video

				linear := nvCreative.Linear
				duration, _ := linear.Duration.MarshalText()
				if len(duration) > 0 {
					// todo
					videoLen, _ := getVideoLength(string(duration))
					*ad.VideoLength = int32(videoLen)
				}

				fillAdTrackingWithVastTracking(linear.TrackingEvents, adTracking)

				*ad.VideoURL = linear.MediaFiles[0].URI
				// 记录三方dsp的视频url
				recordThirdPartyDspCreative(dspId, ad, reqCtx)
				*ad.VideoSize = int32(linear.MediaFiles[0].Size)
				*ad.VideoResolution = strconv.Itoa(linear.MediaFiles[0].Width) + "x" + strconv.Itoa(linear.MediaFiles[0].Height)
				*ad.ImageSize = "VIDEO"
			}
		}
		if asset.Data != nil {
			if asset.Id == int(adxconst.DECSID) && len(asset.Data.Value) > 0 {
				*ad.AppDesc = asset.Data.Value
			}

			if asset.Id == int(adxconst.CTAID) && len(asset.Data.Value) > 0 {
				*ad.CtaText = asset.Data.Value
			}

			if asset.Id == int(adxconst.RATINGID) && len(asset.Data.Value) > 0 {
				rating, err := strconv.ParseFloat(asset.Data.Value, 64)
				if err != nil {
					// todo
				} else {
					*ad.Rating = rating
				}
			}
		}
	}

	*ad.ClickURL = nad.Link.Url
	adTracking.Click = append(adTracking.Click, nad.Link.ClickTrackers...)
	adTracking.Impression = append(adTracking.Impression, nad.ImpTrackers...)

	return nil
}

func fillAd(rdata *rtb.BidResponse, reqCtx *mvutil.ReqCtx, backendCtx *mvutil.BackendCtx) error {
	if rdata == nil || reqCtx == nil || backendCtx == nil || reqCtx.ReqParams == nil || backendCtx.Ads == nil {
		return errors.New("fillAd params error")
	}

	bid, err := getBidByImpId(rdata.GetSeatbid(), ImpId)
	if err != nil {
		return errors.New("madx fillAd getBidByImpId error:%s" + err.Error())
	}

	reqCtx.ReqParams.DspExt = bid.GetExt().GetDataext()
	dspExt, err := reqCtx.ReqParams.GetDspExt()
	// 直接请求pioneer，不会有dspext。
	if err != nil {
		return errors.New("doFillAd decodeDspExt error:" + err.Error())
	}

	if rdata.GetExt() != nil {
		backendCtx.Ads.RKS = rdata.GetExt().GetRks()
	}

	adTmp := &corsair_proto.Campaign{}
	mvutil.SetCampaignDefaultFields(adTmp)
	backendCtx.Ads.DspId = dspExt.DspId
	// mas 大模板或native 会返回多个广告,  isMasMoreAdsWin 表示最后胜出的是这种情况
	isMasMoreAdsWin := ((reqCtx.ReqParams.Param.BigTemplateFlag ||
		reqCtx.ReqParams.Param.AdType == mvconst.ADTypeNative) && dspExt.DspId == mvconst.MAS) ||
		(reqCtx.ReqParams.Param.RequestType == mvconst.REQUEST_TYPE_OPENAPI_AD && dspExt.DspId == mvconst.MAS) ||
		mvutil.IsAppwallOrMoreOffer(reqCtx.ReqParams.Param.AdType) || mvutil.IsRequestPioneerDirectly(&reqCtx.ReqParams.Param) //onlineAPI 流量也支持多个offer返回
	if isMasMoreAdsWin && rdata.GetAsResp() != nil && len(rdata.GetAsResp().GetSdkParam()) > 0 {
		for k, _ := range rdata.GetAsResp().GetSdkParam() {
			adTmp := &corsair_proto.Campaign{}
			mvutil.SetCampaignDefaultFields(adTmp)
			err = fillAdData(adTmp, bid, reqCtx, backendCtx, rdata, k)
			backendCtx.Ads.CampaignList = append(backendCtx.Ads.CampaignList, adTmp)
		}
	} else {
		err = fillAdData(adTmp, bid, reqCtx, backendCtx, rdata, 0)
	}
	if err != nil {
		// mvutil.Logger.Runtime.Debugf("=========== adx fillAd error: %s", err.Error())
		return errors.New("madx fillAd doFillAd error:%s" + err.Error())
	}
	// madx的请求未切量给as(切给mas或都未切)， 最后胜出的不是mas的大模板或native(mas video 或第三方dsp)
	if !reqCtx.ReqParams.IsFakeAs && !isMasMoreAdsWin {
		if reqCtx.ReqParams.Param.TNum > 1 && len(rdata.GetVastMas()) > 0 && rdata.GetAsResp() != nil && len(rdata.GetAsResp().SdkParam) > 1 {
			// TODO 类似于上面的FillMoreAsAds, 将vastmas 解析，然后填到asAds里
			moreAds := corsair_proto.NewBackendAds()
			moreAds.BackendId = mvconst.MAdx
			fillVastAsAds(reqCtx, rdata, moreAds, dspExt.DspId, backendCtx)
			reqCtx.ReqParams.IsMoreAsAds = true
			backendCtx.AsAds = moreAds
		} else if len(rdata.GetRespAs()) > 0 {
			asResp := ad_server.QueryResult_{}
			err := jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal([]byte(rdata.GetRespAs()), &asResp)
			if err == nil {
				if reqCtx.ReqParams.Param.TNum > 1 && len(asResp.GetCampaignList()) > 1 {
					moreAds := corsair_proto.NewBackendAds()
					moreAds.BackendId = mvconst.Mobvista
					FillMoreAsAds(reqCtx, moreAds, &asResp)
					reqCtx.ReqParams.IsMoreAsAds = true
					backendCtx.AsAds = moreAds
				} else {
					reqCtx.ReqParams.Param.Extalgo = asResp.GetAlgoFeatInfo()
					reqCtx.ReqParams.Param.ExtifLowerImp = asResp.GetIfLowerImp()
					reqCtx.ReqParams.Param.Extra3 = asResp.GetStrategy()
					reqCtx.ReqParams.Param.RespFillEcpmFloor = asResp.GetEcpmFloor()
				}
			}
		}

		for _, camp := range backendCtx.Ads.CampaignList {
			if bid.GetExt() != nil && bid.GetExt().Aks != nil {
				camp.AKS = bid.GetExt().Aks
			}
		}
		backendCtx.Ads.CampaignList = append(backendCtx.Ads.CampaignList, adTmp)
	}

	if reqCtx.ReqParams.IsHBRequest {
		reqCtx.ReqParams.BidCur = rdata.GetCur()
		reqCtx.ReqParams.Token = rdata.GetBidid()
		reqCtx.ReqParams.BidRespID = rdata.GetId()
		reqCtx.ReqParams.Price = bid.GetPrice()

		// 算法出价调控实验
		// extAlgo 实验标记的位置 0: 非实验, 1: 实验
		// 实验的部分工程不做价格系数的计算
		var isBiddingExp bool
		// 因为依赖算法的返回
		// 这里不能把对reqCtx.ReqParams.Param.ExtDataInit.PriceFactor系数的控制放到param render的流程里
		// 会存在一种情况是日志里看到 extdata 的 price factor 有系数控制, 价格却没有乘以系数
		// 可以通过在日志的 extalgo 内容按逗号分隔的57位是否是1来看是否属于算法实验控制出价
		extalgo := strings.Split(strings.Split(reqCtx.ReqParams.Param.Extalgo, ";")[0], ",")
		if len(extalgo) > 56 && extalgo[56] == "1" {
			isBiddingExp = true
		}
		if reqCtx.ReqParams.Param.ExtDataInit.Send2RS != 1 && reqCtx.ReqParams.Param.ExtDataInit.PriceFactor != 1 && reqCtx.ReqParams.Param.ExtDataInit.PriceFactor > 0 &&
			reqCtx.ReqParams.Param.ExtDataInit.PriceFactor < 10 && (dspExt.DspId == mvconst.FakeAdserverDsp || dspExt.DspId == mvconst.MAS) && !isBiddingExp { // 需要判断为Adserver || MAS才需要乘于系数
			reqCtx.ReqParams.Price = reqCtx.ReqParams.Price * reqCtx.ReqParams.Param.ExtDataInit.PriceFactor
		}
		// for mopub bid test mode
		if reqCtx.ReqParams.Param.HBBidTestMode == 1 {
			reqCtx.ReqParams.Token = uuid.NewV4().String() + "-74657374" // "test" ascii to hex
			reqCtx.ReqParams.Price = 9900.00
		}
		// bidding price can't less bid floor
		if reqCtx.ReqParams.Param.HBBidTestMode != 1 && reqCtx.ReqParams.Price < reqCtx.ReqParams.Param.BidFloor*100 {
			reqCtx.ReqParams.Nbr = hbconst.BiddingPriceError
			reqCtx.ReqParams.BidRejectCode = hbconst.BackendContentFilter
			return filter.BiddingPriceError
		}
		reqCtx.ReqParams.Nbr = hbconst.OK
		dec1 := decimal.NewMDecimal()
		err := dec1.FromFloat64(reqCtx.ReqParams.Price)
		if err == nil {
			reqCtx.ReqParams.PriceBigDecimal = string(dec1.ToString())
		}
		dec2 := decimal.NewMDecimal()
		err = dec2.FromFloat64(reqCtx.ReqParams.Param.BidFloor * 100)
		if err == nil {
			reqCtx.ReqParams.BidFloorBigDecimal = string(dec2.ToString())
		}
	}

	return nil
}

func fillAdData(adTmp *corsair_proto.Campaign, bid *rtb.BidResponse_SeatBid_Bid, reqCtx *mvutil.ReqCtx, backendCtx *mvutil.BackendCtx, rdata *rtb.BidResponse, k int) error {
	err := doFillAd(adTmp, bid, reqCtx, backendCtx, k)
	if err != nil {
		// mvutil.Logger.Runtime.Debugf("=========== adx fillAdData doFillAd error: %s", err.Error())
		// 出现渲染错误则把dspExt置空
		reqCtx.ReqParams.DspExt = ""
		return err
	}
	dspExt, err := reqCtx.ReqParams.GetDspExt()
	if err != nil {
		return errors.New("fillAdData dspExt err: " + err.Error())
	}
	if bid.GetExt() != nil {
		if len(bid.GetExt().GetAks()) > 0 {
			adTmp.AKS = bid.Ext.Aks // 参数模板化
		}
		if bid.GetExt().GetSkadn() != nil {
			version := bid.GetExt().GetSkadn().GetVersion()
			network := bid.GetExt().GetSkadn().GetNetwork()
			campaign := bid.GetExt().GetSkadn().GetCampaign()
			itunesItem := bid.GetExt().GetSkadn().GetItunesitem()
			nonce := bid.GetExt().GetSkadn().GetNonce()
			sourceApp := bid.GetExt().GetSkadn().GetSourceapp()
			timestamp := bid.GetExt().GetSkadn().GetTimestamp()
			signature := bid.GetExt().GetSkadn().GetSignature()
			need := bid.GetExt().GetSkadn().GetSkneed()
			adTmp.Skad = &corsair_proto.Skad{
				Version:         version,
				Network:         network,
				AppleCampaignId: campaign,
				Targetid:        itunesItem,
				Nonce:           nonce,
				Sourceid:        sourceApp,
				Timestamp:       timestamp,
				Sign:            signature,
				Need:            need,
			}
			if reqCtx.ReqParams.IsBidRequest {
				reqCtx.ReqParams.BidSkAdNetwork = &mtg_hb_rtb.BidResponse_SeatBid_Bid_Ext_Skadn{
					Version:    &version,
					Network:    &network,
					Campaign:   &campaign,
					Itunesitem: &itunesItem,
					Nonce:      &nonce,
					Sourceapp:  &sourceApp,
					Timestamp:  &timestamp,
					Signature:  &signature,
				}
			}
		}
	}
	if dspExt.DspId != mvconst.MAS && dspExt.DspId != mvconst.FakeAdserverDsp {
		fillDeeplink(adTmp, bid, reqCtx) // 非6和13的deeplink特殊处理
	}
	if dspExt.DspId == mvconst.MAS || mvutil.IsAppwallOrMoreOffer(reqCtx.ReqParams.Param.AdType) || mvutil.IsRequestPioneerDirectly(&reqCtx.ReqParams.Param) {
		fillMasAd(rdata, reqCtx, backendCtx, adTmp, k) // MAS
	}
	return nil
}

// fillAdWithBanner: 对所有banner类的返回进行赋值， 包括sdk_banner, IA, mraid 等
func fillAdWithBanner(adTracking *corsair_proto.AdTracking, bid *rtb.BidResponse_SeatBid_Bid, backendCtx *mvutil.BackendCtx, dspId int64, ad *corsair_proto.Campaign, reqCtx *mvutil.ReqCtx) {
	adm := bid.GetAdm()
	// native h5 和高版本banner，splash，都在offer维度返回模版。低版本banner及splash，还是在unit维度返回模版
	// api_version小于2.0且为mas
	if extractor.IsVastBannerDsp(dspId) && mvutil.IsBannerOrSplashOrNativeH5(reqCtx.ReqParams.Param.FormatAdType) {
		var bOsv int32
		if reqCtx.ReqParams.Param.Platform == mvconst.PlatformIOS {
			bOsv = mvutil.GetVersionCode(mvconst.BannerIosSDKVersion)
		} else {
			bOsv = mvutil.GetVersionCode(mvconst.BannerAndriodSDKVersion)
		}

		if reqCtx.ReqParams.Param.FormatSDKVersion.SDKVersionCode < bOsv {
			backendCtx.Ads.BannerUrl = ad.UrlTemplate
		}
	} else if dspId == mvconst.MAS && reqCtx.ReqParams.Param.ApiVersion < mvconst.API_VERSION_2_0 {
		if ad.UrlTemplate != nil {
			backendCtx.Ads.BannerUrl = ad.UrlTemplate
		} else if ad.HtmlTemplate != nil {
			backendCtx.Ads.BannerHtml = ad.HtmlTemplate
		}
		// api_version小于2.0且不为mas（mvdsp）
	} else if dspId != mvconst.MAS && reqCtx.ReqParams.Param.ApiVersion < mvconst.API_VERSION_2_0 {
		if len(adm) > 0 && strings.HasPrefix(adm, "http") {
			backendCtx.Ads.BannerUrl = &adm
		} else {
			backendCtx.Ads.BannerHtml = &adm
		}
		// api_version大于等于2.0且不为mas
	} else if dspId != mvconst.MAS && reqCtx.ReqParams.Param.ApiVersion >= mvconst.API_VERSION_2_0 {
		if len(adm) > 0 && strings.HasPrefix(adm, "http") {
			ad.UrlTemplate = &adm
		} else {
			ad.HtmlTemplate = &adm
		}
	}
	fillBannerTracking(adTracking, bid)
}

func fillVastAsAds(reqCtx *mvutil.ReqCtx, rdata *rtb.BidResponse, ads *corsair_proto.BackendAds,
	dspId int64, backendCtx *mvutil.BackendCtx) error {
	vastStr := rdata.GetVastMas()
	masResp := rdata.GetAsResp()
	if masResp == nil || len(masResp.SdkParam) == 0 {
		return errors.New("VastMoreAds: no asResp or sdkparam")
	}
	v := new(vast.VAST)
	err := xml.Unmarshal([]byte(vastStr), v)
	if err != nil {
		return err
	}
	if len(v.Ads) != len(masResp.SdkParam) {
		return errors.New("VastMoreAds: ads & sdkparam length not equal")
	}
	for i, _ := range masResp.SdkParam {
		ad := new(corsair_proto.Campaign)
		if i > 0 {
			mvutil.SetCampaignDefaultFields(ad)
			adtracking := new(corsair_proto.AdTracking)
			vad := v.Ads[i]
			filladWithVast(ad, adtracking, &vad, reqCtx, dspId, false)

			fillMasAd(rdata, reqCtx, backendCtx, ad, i)
			ad.AdTracking = adtracking
		}
		ads.CampaignList = append(ads.CampaignList, ad)
	}
	ads.DspId = mvconst.MAS
	return nil
}

// fillDeeplink 按照SDK的deeplink使用逻辑，将字段重新赋值
// deeplink失败时，会走 clickUrl, 无论是否成功，都会上报adtracking.click
// 改动： 1. 将ad.ClickUrl （即DSP的landingPage）放到adtracking.click中；2. 将fallbackUrl放到ad.Clickurl
func fillDeeplink(ad *corsair_proto.Campaign, bid *rtb.BidResponse_SeatBid_Bid, ctx *mvutil.ReqCtx) {
	if bid.Ext == nil {
		return
	}
	if len(bid.Ext.GetDeeplink()) == 0 || len(bid.Ext.GetDeeplinkfallbackurl()) == 0 || canChangeOnlineDeeplinkWay(ctx, ad) {
		return
	}

	if ad.AdTracking == nil {
		ad.AdTracking = &corsair_proto.AdTracking{}
	}
	if len(ad.GetClickURL()) > 0 {
		ad.AdTracking.Click = append(ad.AdTracking.Click, ad.GetClickURL())
	}

	ad.ClickURL = bid.Ext.Deeplinkfallbackurl
}

func canChangeOnlineDeeplinkWay(ctx *mvutil.ReqCtx, ad *corsair_proto.Campaign) bool {
	// 必须为仅支持online api单链的流量
	if ctx.ReqParams.Param.RequestType != mvconst.REQUEST_TYPE_OPENAPI_AD {
		return false
	}
	// 限制offer包名
	changeOnlineDeeplinkWayPackageList := extractor.GetCHANGE_ONLINE_DEEPLINK_WAY_PACKAGE_LIST()
	if !mvutil.InStrArray(ad.GetPackageName(), changeOnlineDeeplinkWayPackageList) {
		return false
	}
	// 不支持双链的流量
	adnConfList := extractor.GetADNET_CONF_LIST()
	if onlineUnsupportTwoLinkPubList, ok := adnConfList["onlineUnsupportTwoLinkPubList"]; ok && mvutil.InInt64Arr(ctx.ReqParams.Param.PublisherID, onlineUnsupportTwoLinkPubList) {
		return true
	}
	return false
}

func recordThirdPartyDspCreative(dspId int64, ad *corsair_proto.Campaign, reqCtx *mvutil.ReqCtx) {
	if !mvutil.IsThirdDsp(dspId) || ad.VideoURL == nil || len(*ad.VideoURL) == 0 {
		return
	}
	if ad.PackageName == nil {
		*ad.PackageName = ""
	}
	// 记录生成的三方dsp 视频素材id
	reqCtx.ReqParams.Param.ThirdPartyDspVideoCreativeId = int64(crc32.ChecksumIEEE([]byte(*ad.VideoURL)))
	ctime := time.Now().Unix()
	logStr := strings.Join([]string{
		time.Unix(ctime, 0).Format("20060102"),
		time.Unix(ctime, 0).Format("150405"),
		strconv.FormatInt(ctime, 10),
		strconv.FormatInt(dspId, 10),
		strconv.FormatInt(reqCtx.ReqParams.Param.ThirdPartyDspVideoCreativeId, 10),
		mvutil.TrimBlank(*ad.VideoURL),
		mvutil.TrimBlank(*ad.PackageName),
	}, "\t")
	// 记录三方dsp的视频素材信息
	mvutil.Logger.DspCreativeData.Infof(logStr)
}

func isHitClickModeTest(idfa, idfv string, unitId, appId int64) bool {
	configs := extractor.GetYLHClickModeTestConfig()
	if len(configs) == 0 {
		return false
	}

	index := 0
	if len(idfa) > 0 {
		index = int(crc32.ChecksumIEEE([]byte(idfa)) % 100)
	} else if len(idfv) > 0 {
		index = int(crc32.ChecksumIEEE([]byte(idfv)) % 100)
	} else {
		index = rand.Intn(100)
	}

	unitStr := "unit_" + strconv.FormatInt(unitId, 10)
	rate, ok := configs[unitStr]
	if ok {
		return index < rate
	}

	appStr := "app_" + strconv.FormatInt(appId, 10)
	rate, ok = configs[appStr]
	if ok {
		return index < rate
	}

	rate, ok = configs["all"]
	if ok {
		return index < rate
	}

	return false
}

// 请求 bid server 逻辑切量实验
func hitRequestBidServerTest(reqCtx *mvutil.ReqCtx) {
	// 实验为HB流量
	if !reqCtx.ReqParams.IsHBRequest {
		return
	}
	if reqCtx.ReqParams.Param.RequestPath == mvconst.PATHBidAds {
		// bid_ads 不切量
		return
	}

	abConfArr := extractor.GetRequestBidServerRate()
	if abConfArr == nil || abConfArr.Configs == nil {
		return
	}

	// 根据多个维度进行切量
	for _, conf := range abConfArr.Configs {
		// 命中
		if selectorBidServerABTestFilter(reqCtx, conf) {
			// 根据 device id 切量
			factor := mvutil.GetRandByGlobalTagId(&reqCtx.ReqParams.Param, mvconst.SALT_HB_REQUEST_BID_SERVER, 10000)
			if factor < conf.Rate {
				// hit
				// reqCtx.ReqParams.Param.ExtDataInit.RequestBidServer = 1
				reqCtx.ReqParams.Param.IsHitRequestBidServer = 1
				metrics.IncCounterWithLabelValues(33, "bid")
			}
			return
		}
	}

	return
}

//
func selectorBidServerABTestFilter(reqCtx *mvutil.ReqCtx, c *mvutil.HBRequestBidServerConf) bool {
	if len(c.Cloud) > 0 && !mvutil.InStrArray(mvutil.Cloud(), c.Cloud) {
		return false
	}
	if len(c.Region) > 0 && !mvutil.InStrArray(mvutil.Region(), c.Region) {
		return false
	}
	if len(c.CountryCode) > 0 && !mvutil.InStrArray(reqCtx.ReqParams.Param.CountryCode, c.CountryCode) {
		return false
	}
	if len(c.MediationName) > 0 && !mvutil.InStrArray(reqCtx.ReqParams.Param.MediationName, c.MediationName) {
		return false
	}
	if len(c.Platform) > 0 && !mvutil.InArray(reqCtx.ReqParams.Param.Platform, c.Platform) {
		return false
	}
	if len(c.RequestType) > 0 && !mvutil.InArray(reqCtx.ReqParams.Param.RequestType, c.RequestType) {
		return false
	}
	if len(c.AdType) > 0 && !mvutil.InInt32Arr(reqCtx.ReqParams.Param.AdType, c.AdType) {
		return false
	}
	if len(c.Scenario) > 0 && !mvutil.InStrArray(reqCtx.ReqParams.Param.Scenario, c.Scenario) {
		return false
	}
	if len(c.PublisherID) > 0 && !mvutil.InInt64Arr(reqCtx.ReqParams.Param.PublisherID, c.PublisherID) {
		return false
	}
	if len(c.AppID) > 0 && !mvutil.InInt64Arr(reqCtx.ReqParams.Param.AppID, c.AppID) {
		return false
	}
	if len(c.UnitID) > 0 && !mvutil.InInt64Arr(reqCtx.ReqParams.Param.UnitID, c.UnitID) {
		return false
	}

	return true
}

func GetAdTrackURLPrefix(reqParam *mvutil.RequestParams, backend int, offerID, backendConfig string) string {
	adBackend := fmt.Sprintf("%d", backend)
	adBackendData := fmt.Sprintf("%d:%s:1", backend, offerID)
	dspExt := reqParam.DspExt
	if backend != mvconst.MAdx {
		dspExt = ""
	}
	dPrice := "0"
	if dpriceObj, ifFind := reqParam.DPrice.Load(strconv.Itoa(backend)); ifFind {
		dprice, ok := dpriceObj.(int64)
		if ok {
			dPrice = strconv.FormatInt(dprice, 10)
		}
	}

	pFactor := ""
	if priceFactorObj, ifFind := reqParam.PriceFactor.Load(strconv.Itoa(backend)); ifFind {
		priceFactor, ok := priceFactorObj.(string)
		if ok {
			pFactor = priceFactor
		}
	}

	rawParam := mvutil.SerializeMPPart(reqParam, adBackend, adBackendData, backendConfig, dspExt, dPrice, pFactor)
	base64Data := mvutil.Base64Encode(rawParam)
	base64Data = strings.Replace(base64Data, "=", "%3D", -1)
	base64Data = strings.Replace(base64Data, "+", "%2B", -1)
	base64Data = strings.Replace(base64Data, "/", "%2F", -1)
	var buf bytes.Buffer
	buf.WriteString(mvutil.GetUrlScheme(int(reqParam.Param.HTTPReq)))

	mtrackDomain := mvutil.Config.CommonConfig.TrackConfig.TrackHost
	// 非归因相关埋点上报切到走cdn的域名
	if reqParam.Param.UseCdnTrackingDomain == 1 {
		cdnTrackingDomain := extractor.GetCdnTrackingDomain(mvutil.Config.CommonConfig.TrackConfig.TrackHost)
		if len(cdnTrackingDomain) > 0 {
			mtrackDomain = cdnTrackingDomain
		}
	}

	if mvutil.NeedNewJssdkDomain(reqParam.Param.RequestPath, reqParam.Param.Ndm) && len(reqParam.Param.MTrackDomain) > 0 {
		buf.WriteString(reqParam.Param.MTrackDomain + mvutil.Config.CommonConfig.TrackConfig.PlayTrackPath)
	} else {
		buf.WriteString(mtrackDomain + mvutil.Config.CommonConfig.TrackConfig.PlayTrackPath)
	}
	buf.WriteString("?k=" + reqParam.Param.RequestID)
	buf.WriteString("&mp=" + base64Data)
	return buf.String()
}

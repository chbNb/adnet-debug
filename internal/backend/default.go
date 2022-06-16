package backend

import (
	"bytes"
	supply_mvutil "gitlab.mobvista.com/ADN/chasm/module/supply/mvutil"
	"math/rand"
	"regexp"
	"strconv"

	jsoniter "github.com/json-iterator/go"
	"gitlab.mobvista.com/ADN/adnet/internal/corsair_proto"

	"gitlab.mobvista.com/ADN/adnet/internal/mvconst"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	"gitlab.mobvista.com/adserver/recommend_protocols/go/ad_server"
)

const (
	PER_RATE_0   int32 = 0
	PER_RATE_25  int32 = 25
	PER_RATE_50  int32 = 50
	PER_RATE_75  int32 = 75
	PER_RATE_100 int32 = 100
)

const (
	SEC_SPLITTER   string = "\t"
	FIELD_SPLITTER string = ":"
)

type Filter interface {
	Filter(param *mvutil.RequestParams) bool
}

type DemandExt struct {
	BackendId int
	DspId     int64
}

type Backender interface {
	GetAds(param *mvutil.RequestParams, backendCtx *mvutil.BackendCtx) (*corsair_proto.BackendAds, error)
}

func BackendFilter(filter Filter, param *mvutil.RequestParams) bool {
	return filter.Filter(param)
}

var backenders = make(map[int]Backender)

func Register(id int, backender Backender) {
	if _, exists := backenders[id]; exists {
		mvutil.Logger.Runtime.Warnf("Backender id=[%d] already registered", id)
	}
	mvutil.Logger.Runtime.Infof("Backender id=[%d] registered success", id)
	backenders[id] = backender
}

var marcoReg *regexp.Regexp
var marco1Reg *regexp.Regexp

func init() {
	marcoReg = regexp.MustCompile("\\$\\{.+?\\}")
	marco1Reg = regexp.MustCompile("\\$\\{.+?\\}")
}

type Ints []int64

type Strs []string

func renderExcludeIDs(exclude string) map[int64]bool {
	resMap := make(map[int64]bool)
	var eList Ints
	err := jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal([]byte(exclude), &eList)
	if err == nil && len(eList) > 0 {
		for _, v := range eList {
			resMap[v] = true
		}
		return resMap
	}
	var eStrList Strs
	err = jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal([]byte(exclude), &eStrList)
	if err == nil && len(eStrList) > 0 {
		for _, v := range eStrList {
			vInt, err2 := strconv.ParseInt(v, 10, 64)
			if err2 == nil {
				resMap[vInt] = true
			}
		}
		return resMap
	}
	return resMap
}

func renderExcludeIDslice(exclude string) []int64 {
	var eList Ints
	err := jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal([]byte(exclude), &eList)
	if err == nil {
		return eList
	}

	// 兼容string list的情况
	var eStrList []string
	err = jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal([]byte(exclude), &eStrList)
	if err == nil {
		var excludeList []int64
		for _, v := range eStrList {
			vInt, err2 := strconv.ParseInt(v, 10, 64)
			if err2 == nil {
				excludeList = append(excludeList, vInt)
			}
		}
		return excludeList
	}
	return eList
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

func pickItem(data []string) string {
	dataLen := len(data)
	if dataLen == 1 {
		return data[0]
	} else if dataLen > 1 {
		index := rand.Intn(dataLen)
		return data[index]
	} else {
		return ""
	}
}

func filterDownload(offerPreference []int, platform int) bool {
	if len(offerPreference) != 0 {
		if platform == mvconst.PlatformAndroid && !mvutil.InArray(mvconst.APK, offerPreference) {
			//mvutil.Logger.Runtime.Warnf("request_id=[%s] backendId=[%d] Not Support ApkDownload AppId=[%d]", requestId, backendId, appId)
			return true
		}
		if platform == mvconst.PlatformIOS && !mvutil.InArray(mvconst.APPSTORE, offerPreference) {
			//mvutil.Logger.Runtime.Warnf("request_id=[%s] backendId=[%d] Not Support AppStoreDownload AppId=[%d]", requestId, backendId, appId)
			return true
		}
	}
	return false
}

// func fillOrientation(adType int32, orientation int) int {
// 	if adType == mvconst.ADTypeNativeVideo {
// 		return mvconst.ORIENTATION_LANDSCAPE
// 	}
// 	return orientation
// }

// func fillLinkTypeWithDownload(platform int) int32 {
// 	if platform == mvconst.PlatformIOS {
// 		return mvconst.APPSTORE
// 	} else if platform == mvconst.PlatformAndroid {
// 		return mvconst.APK
// 	} else {
// 		return mvconst.OTHER
// 	}
// }

func FixedLandscape4Both(orientation int) int {
	//横竖屏,如果是both，传横屏 todo
	if orientation == mvconst.ORIENTATION_BOTH {
		orientation = mvconst.ORIENTATION_LANDSCAPE
	}
	return orientation
}

// adtracking特性过滤
func filterAdTracking(adType int32, platform int, adTracking int, sdkVerion supply_mvutil.SDKVersionItem) bool {
	if platform == mvconst.PlatformAndroid {
		if sdkVerion.SDKType != "ma" && sdkVerion.SDKType != "mal" {
			return true
		}

		//android adTracking.click 排除 930/940/950/951/960
		if len(mvutil.Config.CommonConfig.DVIConfig.AndroidExpVersion) > 0 && mvutil.InInt32Arr(sdkVerion.SDKVersionCode, mvutil.Config.CommonConfig.DVIConfig.AndroidExpVersion) &&
			(adTracking == mvconst.AdTracking_Click ||
				adTracking == mvconst.AdTracking_Both_Click_Imp) {
			return true
		}

		//android rv 展示和点击 8.3.0 及以上版本
		if adType == mvconst.ADTypeRewardVideo && (adTracking == mvconst.AdTracking_Imp || adTracking == mvconst.AdTracking_Click || adTracking == mvconst.AdTracking_Both_Click_Imp) &&
			sdkVerion.SDKVersionCode >= mvconst.AdTrackingAndroidRV {
			return false
		}

		//android native video 展示和点击 8.2.1 及以上版本
		if adType == mvconst.ADTypeNativeVideo && (adTracking == mvconst.AdTracking_Imp || adTracking == mvconst.AdTracking_Click || adTracking == mvconst.AdTracking_Both_Click_Imp) &&
			sdkVerion.SDKVersionCode >= mvconst.AdTrackingAndroidNativeVideo {
			return false
		}

		// android iv 展示和点击 8.10.0
		if adType == mvconst.ADTypeInterstitialVideo && (adTracking == mvconst.AdTracking_Imp || adTracking == mvconst.AdTracking_Click || adTracking == mvconst.AdTracking_Both_Click_Imp) &&
			sdkVerion.SDKVersionCode >= mvconst.AdTrackingAndroidIV {
			return false
		}

		// android native pic 展示 和 点击 8.6.0
		if adType == mvconst.ADTypeNativePic && (adTracking == mvconst.AdTracking_Imp || adTracking == mvconst.AdTracking_Both_Click_Imp) && sdkVerion.SDKVersionCode >= mvconst.AdTrackingAndroidNativePicImp {
			return false
		}

		// android native pic 点击 8.4.0
		if adType == mvconst.ADTypeNativePic && adTracking == mvconst.AdTracking_Click && sdkVerion.SDKVersionCode >= mvconst.AdTrackingAndroidNativePicClick {
			return false
		}
		// sdk banner 及IA DI 不做版本判断，都支持
		if mvutil.IsBannerOrSplashOrNativeH5(adType) || adType == mvconst.ADTypeInteractive {
			return false
		}

		return true
	}

	if platform == mvconst.PlatformIOS {
		if sdkVerion.SDKType != "mi" {
			return true
		}

		//iOS 排除 390和 391
		if sdkVerion.SDKVersionCode == mvconst.AdTrackingIOSExpV1 || sdkVerion.SDKVersionCode == mvconst.AdTrackingIOSExpV2 {
			return true
		}
		// ios rv 展示和点击 2.8.0 及以上版本
		if adType == mvconst.ADTypeRewardVideo && (adTracking == mvconst.AdTracking_Imp || adTracking == mvconst.AdTracking_Click || adTracking == mvconst.AdTracking_Both_Click_Imp) &&
			sdkVerion.SDKVersionCode >= mvconst.AdTrackingIOSRV {
			return false
		}

		// ios iv 展示和点击 3.6.0及以上版本
		if adType == mvconst.ADTypeInterstitialVideo && (adTracking == mvconst.AdTracking_Imp || adTracking == mvconst.AdTracking_Click || adTracking == mvconst.AdTracking_Both_Click_Imp) &&
			sdkVerion.SDKVersionCode >= mvconst.AdTrackingIOSIV {
			return false
		}

		// ios native 展示 3.1.0
		if (adType == mvconst.ADTypeNativePic || adType == mvconst.ADTypeNativeVideo) && (adTracking == mvconst.AdTracking_Imp || adTracking == mvconst.AdTracking_Both_Click_Imp) &&
			sdkVerion.SDKVersionCode >= mvconst.AdTrackingIOSNativeImp {
			return false
		}

		// ios native 点击 2.8.0
		if (adType == mvconst.ADTypeNativePic || adType == mvconst.ADTypeNativeVideo) && adTracking == mvconst.AdTracking_Click &&
			sdkVerion.SDKVersionCode >= mvconst.AdTrackingIOSNativeClick {
			return false
		}

		// sdk banner 及IA DI 不做版本判断，都支持
		if mvutil.IsBannerOrSplashOrNativeH5(adType) || adType == mvconst.ADTypeInteractive {
			return false
		}
		return true
	}
	return true
}

func formatPercentageURL(percent int32, prefix, adType, ext string) *corsair_proto.PlayPercent {
	result := corsair_proto.NewPlayPercent()
	result.Rate = percent
	urlBuf := bytes.NewBufferString(prefix)
	urlBuf.WriteString("&type=" + adType)
	urlBuf.WriteString("&key=play_percentage")
	urlBuf.WriteString("&rate=" + strconv.FormatInt(int64(percent), 10))
	if len(ext) > 0 {
		urlBuf.WriteString("&text=" + ext)
	}
	result.URL = urlBuf.String()
	return result
}

func formatAdTrackingURL(prefix, adType, key, ext string) string {
	urlBuf := bytes.NewBufferString(prefix)
	urlBuf.WriteString("&type=" + adType)
	urlBuf.WriteString("&key=" + key)
	if len(ext) > 0 {
		urlBuf.WriteString("&text=" + ext)
	}
	return urlBuf.String()
}

func formatBackendData(backendId int, videoURL, requestKey string) string {
	dataBuf := bytes.NewBufferString("")
	dataBuf.WriteString(strconv.Itoa(backendId))
	dataBuf.WriteString(FIELD_SPLITTER)
	dataBuf.WriteString(requestKey)
	dataBuf.WriteString(FIELD_SPLITTER)
	if len(videoURL) > 0 {
		dataBuf.WriteString(strconv.Itoa(mvconst.ContentVideo))
	} else {
		dataBuf.WriteString(strconv.Itoa(mvconst.ContentImg))
	}
	return dataBuf.String()
}

func renderPackageList(pkgList string) map[string]bool {
	resMap := make(map[string]bool)
	if len(pkgList) == 0 {
		return resMap
	}
	var eStrList Strs
	err := jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal([]byte(pkgList), &eStrList)
	if err == nil && len(eStrList) > 0 {
		for _, v := range eStrList {
			resMap[v] = true
		}
	}
	return resMap
}

func renderPackageListSlice(pkgList string) []string {
	var eStrList Strs
	err := jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal([]byte(pkgList), &eStrList)
	if err != nil {
		return []string{}
	}
	return []string(eStrList)
}

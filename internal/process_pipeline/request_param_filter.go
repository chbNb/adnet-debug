package process_pipeline

import (
	"errors"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
	"gitlab.mobvista.com/ADN/adnet/internal/errorcode"
	"gitlab.mobvista.com/ADN/adnet/internal/extractor"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/constant"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/filter"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/helpers"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/req_context"
	"gitlab.mobvista.com/ADN/adnet/internal/mvconst"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	"gitlab.mobvista.com/ADN/adnet/internal/watcher"
	supply_mvutil "gitlab.mobvista.com/ADN/chasm/module/supply/mvutil"
	"gitlab.mobvista.com/ADN/mtg_openrtb/pkg/mtg_hb_rtb"
)

const (
	AndroidConverVer = 90000 // android 9.0.0 做network_type的map
)

type RequestParamFilter struct {
}

var StartMode string

func checkError(requestId string, err error) {
	if err != nil {
		mvutil.Logger.Runtime.Warnf("request_id=[%s] ParseParam error:%s", requestId, err.Error())
	}
}

func requestHeader(req *http.Request, key string) string {
	if values, ok := req.Header[key]; ok && len(values) > 0 {
		return values[0]
	}
	return ""
}

func GetClientIP(req *http.Request) string {
	forwardedByClientIP := true
	if forwardedByClientIP {
		clientIP := strings.TrimSpace(requestHeader(req, "X-Real-Ip"))
		if len(clientIP) > 0 {
			return clientIP
		}
		clientIP = requestHeader(req, "X-Forwarded-For")
		if index := strings.IndexByte(clientIP, ','); index >= 0 {
			clientIP = clientIP[0:index]
		}
		clientIP = strings.TrimSpace(clientIP)
		if len(clientIP) > 0 {
			return clientIP
		}
	}
	if ip, _, err := net.SplitHostPort(strings.TrimSpace(req.RemoteAddr)); err == nil {
		return ip
	}
	return ""
}

func AndroidConvertNetWorkType(networkType int) int {
	switch networkType {
	case 13, 19:
		return mvconst.NETWORK_TYPE_4G
	case 3, 5, 6, 8, 10, 12, 14, 15:
		return mvconst.NETWORK_TYPE_3G
	case 1, 2, 4, 7, 11:
		return mvconst.NETWORK_TYPE_2G
	default:
		return mvconst.NETWORK_TYPE_UNKNOWN
	}
}

// 依然是解析参数
func (rpf *RequestParamFilter) Process(data interface{}) (interface{}, error) {
	in, ok := data.(*mvutil.RequestParams)
	if !ok {
		return nil, errors.New("RequestParamFilter input type should be mvutil.RequestParams")
	}

	in.Param.RequestID = mvutil.GetGoTkClickID() // 生成请求ID
	in.Param.StartMode = StartMode               // 判断是否为灰度机器(启动模式)
	// TODO refactor hb and adnet
	var err error
	// 竞价请求会多带一个body数据体
	if in.IsHBRequest && len(in.PostData) > 0 {
		err = renderOpenRTBRequestParam(in)
	} else {
		err = renderRequestParam(in)
	}
	if err != nil {
		return nil, err
	}
	// 网路相关   设置NetworkType
	if in.Param.NetworkType > 0 && in.Param.NetworkType != mvconst.NETWORK_TYPE_WIFI &&
		in.Param.Platform == mvconst.PlatformAndroid && in.Param.FormatSDKVersion.SDKType == "mal" &&
		in.Param.FormatSDKVersion.SDKVersionCode >= AndroidConverVer && in.Param.FormatSDKVersion.SDKVersionCode < mvconst.AndroidCorectNetworkType {
		in.Param.NetworkType = AndroidConvertNetWorkType(in.Param.NetworkType)
	}
	if len(mvutil.Cloud()) > 0 {
		in.Cloud = mvutil.Cloud()
	}
	if len(mvutil.Region()) > 0 {
		in.Region = mvutil.Region()
	}
	// hb处理
	if in.IsHBRequest {
		// in.Param.Sign = "NO_CHECK_SIGN"
		in.Nbr = -1
		in.Param.OnlyImpression = 1
		in.Param.VideoVersion = "1.0" // 客户端支持native_video与否
		in.Param.VersionFlag = 1      //标志业务版本，位操作；传1 则返回ad_tracking字段
		in.Param.TNum = 1             // 请求normal广告个数

		// 1）通知模式，1-开启，不传或者为其他值则默认不开启
		//2）如果ping_mode=1，会返回click_url和notice_url，在用户点击click_url的时候sdk需要异步通知notice_url；
		//否则只会返回click_url，用户点击click_url即可
		if in.Param.RequestPath != mvconst.PATHBidAds {
			in.Param.PingMode = 1
		}

		var sdkVersionNoBid bool
		// 类型转换，int转str
		platformStr := helpers.GetOs(in.Param.Platform)
		appIDStr := strconv.FormatInt(in.Param.AppID, 10)
		unitIDStr := strconv.FormatInt(in.Param.UnitID, 10)
		// 获取hb sdkVersion信息
		if bidSDKVersionConfig, found := extractor.GetHBBidSDKVersionData(appIDStr, unitIDStr, platformStr); found {
			sdkVersionConf := helpers.RenderSdkVersion(bidSDKVersionConfig)
			reqSDKVersion := helpers.RenderSdkVersion(in.Param.SDKVersion)
			// skip bid_ads request
			if in.Param.RequestPath == mvconst.PATHBid && reqSDKVersion.SDKVersionCode < sdkVersionConf.SDKVersionCode {
				sdkVersionNoBid = true
			}
		}

		if sdkVersionNoBid {
			return nil, filter.BidMSDKVersionTooLow
		}

		if in.Param.Platform != constant.Android && in.Param.Platform != constant.IOS {
			return in, filter.AppPlatformError
		}

		in.Param.Scenario = constant.OpenApi
		algorithm := "hb-" + req_context.GetInstance().Cloud + "-" + req_context.GetInstance().Region
		// cdn abtest flag
		if len(in.Param.MvLine) > 0 {
			algorithm = algorithm + "-" + in.Param.MvLine
		}
		in.Param.Algorithm = algorithm
		in.Param.Extra = algorithm
		// apiVersion
		if in.Param.ApiVersion <= 0 {
			apiVersion := "1.5"
			if in.Param.RequestPath == mvconst.PATHBidAds {
				apiVersion = "1.2"
			}
			// when the s2s has not api_version
			if len(in.Param.HBS2SBidID) > 0 {
				if (in.Param.Platform == constant.IOS && checkSdkVersion(in.Param.SDKVersion, constant.MinIOSSDKVersion)) ||
					(in.Param.Platform == constant.Android && checkSdkVersion(in.Param.SDKVersion, constant.MiniAndroidSDKVersion)) {
					apiVersion = "1.9"
				}
			}
			in.Param.ApiVersion, _ = strconv.ParseFloat(apiVersion, 64)
			in.Param.ApiVersionCode, _ = mvutil.IntVer(apiVersion)
		}

		if mvutil.Config != nil {
			if len(mvutil.Config.Cloud) > 0 {
				in.Cloud = mvutil.Config.Cloud
			}
			if len(mvutil.Config.Region) > 0 {
				in.Region = mvutil.Config.Region
			}
		}
	}

	// 统计rawRequest
	watcher.AddWatchValue("raw_request", float64(1))
	return in, nil
}

func checkSdkVersion(sdkVersion string, version int32) bool {
	sdkVersionData := helpers.RenderSdkVersion(sdkVersion)
	if sdkVersionData != nil && sdkVersionData.SDKVersionCode > 0 && sdkVersionData.SDKVersionCode >= version {
		return true
	}
	return false
}

// 将req的queryMap中的参数解析到req的param中去
func renderRequestParam(req *mvutil.RequestParams) error {
	reqId, err := req.QueryMap.GetString("req_id", true, "")
	checkError(req.Param.RequestID, err)
	if len(reqId) == 24 {
		req.Param.RequestID = reqId
	}

	platform, err := req.QueryMap.GetInt("platform", 0)
	checkError(req.Param.RequestID, err)
	req.Param.Platform = platform
	req.Param.Os = helpers.GetOs(req.Param.Platform)

	osVersion, err := req.QueryMap.GetString("os_version", true, "")
	checkError(req.Param.RequestID, err)
	req.Param.OSVersion = osVersion
	req.Param.OSVersionCode, _ = mvutil.IntVer(osVersion)

	packageName, err := req.QueryMap.GetString("package_name", true, "")
	checkError(req.Param.RequestID, err)
	req.Param.PackageName = packageName

	appVersionName, err := req.QueryMap.GetString("app_version_name", true, "")
	checkError(req.Param.RequestID, err)
	req.Param.AppVersionName = appVersionName

	appVersionCode, err := req.QueryMap.GetString("app_version_code", false, "")
	checkError(req.Param.RequestID, err)
	req.Param.AppVersionCode = appVersionCode

	orientation, err := req.QueryMap.GetInt("orientation", 0)
	checkError(req.Param.RequestID, err)
	req.Param.Orientation = orientation

	brand, err := req.QueryMap.GetString("brand", true, "0")
	checkError(req.Param.RequestID, err)
	req.Param.Brand = brand

	model, err := req.QueryMap.GetString("model", true, "0")
	checkError(req.Param.RequestID, err)
	// 如果多次encode，则再一次decode
	if strings.Contains(model, "%") {
		model = mvutil.UrlDecode(model)
	}
	req.Param.Model = strings.ToLower(model)
	req.Param.RModel = model

	androidID, err := req.QueryMap.GetString("android_id", true, "")
	checkError(req.Param.RequestID, err)
	req.Param.AndroidID = androidID

	iMEI, err := req.QueryMap.GetString("imei", true, "")
	checkError(req.Param.RequestID, err)
	req.Param.IMEI = iMEI

	imeiMd5, err := req.QueryMap.GetString("imei_md5", true, "")
	checkError(req.Param.RequestID, err)
	req.Param.ImeiMd5 = imeiMd5

	mac, err := req.QueryMap.GetString("mac", true, "")
	checkError(req.Param.RequestID, err)
	req.Param.MAC = mac

	gaid, err := req.QueryMap.GetString("gaid2", true, "")
	gaid = mvutil.Base64Decode(gaid)
	checkError(req.Param.RequestID, err)
	req.Param.GAID = gaid
	// gaid加密传输需求
	if len(req.Param.GAID) <= 0 {
		gaid, err := req.QueryMap.GetString("gaid", true, "")
		checkError(req.Param.RequestID, err)
		req.Param.GAID = gaid
	}

	idfa, err := req.QueryMap.GetString("idfa", true, "")
	checkError(req.Param.RequestID, err)
	req.Param.IDFA = idfa

	mnc, err := req.QueryMap.GetString("mnc", true, "")
	checkError(req.Param.RequestID, err)
	req.Param.MNC = mnc

	mcc, err := req.QueryMap.GetString("mcc", true, "")
	checkError(req.Param.RequestID, err)
	req.Param.MCC = mcc

	networkType, err := req.QueryMap.GetInt("network_type", 0)
	checkError(req.Param.RequestID, err)
	// sdk network_type=0表示没有获取，network_type=1表示获取了但是未知，服务端只有0表示未知，将1强制赋值0
	if networkType == 1 {
		networkType = 0
	}
	req.Param.NetworkType = networkType

	language, err := req.QueryMap.GetString("language", true, "")
	checkError(req.Param.RequestID, err)
	req.Param.Language = language

	timeZone, err := req.QueryMap.GetString("timezone", true, "")
	checkError(req.Param.RequestID, err)
	req.Param.TimeZone = timeZone

	sdkVersion, err := req.QueryMap.GetString("sdk_version", true, "0")
	checkError(req.Param.RequestID, err)
	// process SDKVersion
	req.Param.SDKVersion = strings.ToLower(sdkVersion)
	req.Param.FormatSDKVersion = supply_mvutil.RenderSDKVersion(req.Param.SDKVersion)

	gpVersion, err := req.QueryMap.GetString("gp_version", true, "")
	checkError(req.Param.RequestID, err)
	req.Param.GPVersion = gpVersion

	gpsv, err := req.QueryMap.GetString("gpsv", true, "")
	checkError(req.Param.RequestID, err)
	req.Param.GPSV = gpsv

	screenSize, err := req.QueryMap.GetString("screen_size", true, "")
	checkError(req.Param.RequestID, err)
	req.Param.ScreenSize = screenSize

	lat, err := req.QueryMap.GetString("lat", true, "")
	checkError(req.Param.RequestID, err)
	req.Param.LAT = lat

	lng, err := req.QueryMap.GetString("lng", true, "")
	checkError(req.Param.RequestID, err)
	req.Param.LNG = lng

	gpst, err := req.QueryMap.GetString("gpst", true, "")
	checkError(req.Param.RequestID, err)
	req.Param.GPST = gpst

	gpsAccuracy, err := req.QueryMap.GetString("gps_accuracy", true, "")
	checkError(req.Param.RequestID, err)
	req.Param.GPSAccuracy = gpsAccuracy

	gpsType, err := req.QueryMap.GetInt("gps_type", 0)
	checkError(req.Param.RequestID, err)
	req.Param.GPSType = gpsType

	d1, err := req.QueryMap.GetString("d1", true, "")
	checkError(req.Param.RequestID, err)
	req.Param.D1 = d1

	d2, err := req.QueryMap.GetString("d2", true, "")
	checkError(req.Param.RequestID, err)
	req.Param.D2 = d2

	d3, err := req.QueryMap.GetString("d3", true, "")
	checkError(req.Param.RequestID, err)
	req.Param.D3 = d3

	appID, err := req.QueryMap.GetInt64("app_id")
	checkError(req.Param.RequestID, err)
	// TODO refactor hb and adnet
	if (err != nil || appID < 0) && req.Param.RequestPath != mvconst.PATHREQUEST {
		return errorcode.EXCEPTION_APP_ID_EMPTY
	}
	req.Param.AppID = appID

	unitID, err := req.QueryMap.GetInt64("unit_id", 0)
	checkError(req.Param.RequestID, err)
	// TODO refactor hb and adnet
	if mvutil.IsHbOrV3OrV5Request(req.Param.RequestPath) && unitID <= 0 {
		mvutil.Logger.Runtime.Warnf("request_id=[%s] has filter by path=[%s] UnitID is empty", req.Param.RequestID, req.Param.RequestPath)
		return errorcode.EXCEPTION_UNIT_ID_EMPTY
	}
	req.Param.UnitID = unitID

	sign, err := req.QueryMap.GetString("sign", false, "")
	checkError(req.Param.RequestID, err)
	// sign = "NO_CHECK_SIGN"
	req.Param.Sign = sign

	category, err := req.QueryMap.GetInt("category", 0)
	checkError(req.Param.RequestID, err)
	req.Param.Category = category

	adNum, err := req.QueryMap.GetInt("ad_num", 0)
	checkError(req.Param.RequestID, err)
	// todo
	req.Param.AdNum = int32(adNum)

	pingMode, err := req.QueryMap.GetInt("ping_mode", 0)
	checkError(req.Param.RequestID, err)
	req.Param.PingMode = pingMode

	unitSize, err := req.QueryMap.GetString("unit_size", false, "")
	checkError(req.Param.RequestID, err)
	req.Param.UnitSize = unitSize

	camIdsStr, err := req.QueryMap.GetString("display_cids", false, "")
	checkError(req.Param.RequestID, err)
	var displayCamIds Ints
	err = jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal([]byte(camIdsStr), &displayCamIds)
	if err == nil {
		// decode error
		req.Param.DisplayCamIds = displayCamIds
	}
	//使用新的别名
	excludeIDS, err := req.QueryMap.GetString("e", false, "")
	if len(excludeIDS) == 0 {
		excludeIDS, err = req.QueryMap.GetString("fqci", false, "")
	}
	if len(excludeIDS) == 0 { //新字段无值时使用老子段获取
		excludeIDS, err = req.QueryMap.GetString("exclude_ids", false, "")
	}
	checkError(req.Param.RequestID, err)
	req.Param.ExcludeIDS = excludeIDS

	offset, err := req.QueryMap.GetInt("offset", 0)
	checkError(req.Param.RequestID, err)
	// todo
	req.Param.Offset = int32(offset)

	// 参数替换 https://confluence.mobvista.com/pages/viewpage.action?pageId=49639079
	sessionID, err := req.QueryMap.GetString("a", false, "")
	if len(sessionID) == 0 {
		sessionID, err = req.QueryMap.GetString("session_id", false, "")
	}
	checkError(req.Param.RequestID, err)
	req.Param.SessionID = sessionID
	if len(req.Param.SessionID) <= 0 {
		// 记录此次为新session
		req.Param.IsNewSession = true
		req.Param.SessionID = mvutil.GetRequestID()
	}

	parentSessionID, err := req.QueryMap.GetString("parent_session_id", false, "")
	checkError(req.Param.RequestID, err)
	req.Param.ParentSessionID = parentSessionID
	if len(req.Param.ParentSessionID) <= 0 {
		req.Param.ParentSessionID = mvutil.GetRequestID()
	}

	onlyImpression, err := req.QueryMap.GetInt("only_impression", 0)
	checkError(req.Param.RequestID, err)
	req.Param.OnlyImpression = onlyImpression

	netWork, err := req.QueryMap.GetString("network", true, "")
	checkError(req.Param.RequestID, err)
	req.Param.NetWork = netWork

	impressionImage, err := req.QueryMap.GetInt("impression_image", 0)
	checkError(req.Param.RequestID, err)
	req.Param.ImpressionImage = impressionImage

	adSourceID, err := req.QueryMap.GetInt("ad_source_id", 0)
	checkError(req.Param.RequestID, err)
	req.Param.AdSourceID = adSourceID

	adType, err := req.QueryMap.GetInt("ad_type", 0)
	checkError(req.Param.RequestID, err)
	// todo
	req.Param.AdType = int32(adType)

	nativeInfo, err := req.QueryMap.GetString("native_info", false, "")
	checkError(req.Param.RequestID, err)
	req.Param.NativeInfo = nativeInfo

	realAppID, err := req.QueryMap.GetInt64("real_app_id", 0)
	checkError(req.Param.RequestID, err)
	req.Param.RealAppID = realAppID

	isOffline, err := req.QueryMap.GetInt("is_offline", 0)
	checkError(req.Param.RequestID, err)
	req.Param.IsOffline = isOffline

	scenario, err := req.QueryMap.GetString("scenario", false, "")
	checkError(req.Param.RequestID, err)
	req.Param.Scenario = scenario

	tnum, err := req.QueryMap.GetInt("tnum", 0)
	checkError(req.Param.RequestID, err)
	req.Param.TNum = tnum
	if req.Param.TNum == 0 {
		req.Param.UnSupportSdkTrueNum = 1
	} else {
		req.Param.UnSupportSdkTrueNum = 0
	}

	network, err := req.QueryMap.GetInt64("network", 0)
	checkError(req.Param.RequestID, err)
	req.Param.Network = network

	imageSize, err := req.QueryMap.GetString("image_size", false, "")
	checkError(req.Param.RequestID, err)
	req.Param.ImageSize = imageSize
	req.Param.ImageSizeID = mvconst.GetImageSizeID(imageSize)
	// 如果image_size是id
	imageSizeID, err := strconv.Atoi(imageSize)
	if err == nil && imageSizeID > 0 {
		req.Param.ImageSizeID = imageSizeID
	}
	//onlineAPI 增加 image_w & image_h
	if req.Param.RequestPath == mvconst.PATHOnlineApi && req.Param.ImageSizeID == mvconst.IMAGE_SIZE_ID_UNKOWN {
		imageW, err := req.QueryMap.GetInt("image_w", 0)
		checkError(req.Param.RequestID, err)
		imageH, err := req.QueryMap.GetInt("image_h", 0)
		checkError(req.Param.RequestID, err)
		if imageH > 0 && imageW > 0 {
			imageSize = strconv.Itoa(imageW) + "x" + strconv.Itoa(imageH)
			req.Param.ImageSize = imageSize
			req.Param.ImageSizeID = mvconst.GetImageSizeID(imageSize)
		}

	}

	installIDS, err := req.QueryMap.GetString("install_ids", false, "")
	checkError(req.Param.RequestID, err)
	req.Param.InstallIDS = installIDS

	displayCIDS, err := req.QueryMap.GetString("display_cids", false, "")
	checkError(req.Param.RequestID, err)
	req.Param.DisplayCIDS = displayCIDS

	frameNum, err := req.QueryMap.GetInt("frame_num", 1)
	checkError(req.Param.RequestID, err)
	req.Param.FrameNum = frameNum

	idfv, err := req.QueryMap.GetString("idfv", false, "")
	checkError(req.Param.RequestID, err)
	req.Param.IDFV = idfv

	openIDFA, err := req.QueryMap.GetString("openidfa", false, "")
	checkError(req.Param.RequestID, err)
	req.Param.OpenIDFA = openIDFA

	priceFloor, err := req.QueryMap.GetFloat("price_floor", 0.0)
	checkError(req.Param.RequestID, err)
	req.Param.PriceFloor = priceFloor

	httpReq, err := req.QueryMap.GetInt("http_req", 0)
	checkError(req.Param.RequestID, err)
	// todo
	req.Param.HTTPReq = int32(httpReq)

	dvi, err := req.QueryMap.GetString("dvi", false, "")
	checkError(req.Param.RequestID, err)
	req.Param.DVI = dvi

	baseIDS, err := req.QueryMap.GetString("base_ids", false, "")
	checkError(req.Param.RequestID, err)
	req.Param.BaseIDS = baseIDS

	videoVersion, err := req.QueryMap.GetString("video_version", false, "")
	checkError(req.Param.RequestID, err)
	req.Param.VideoVersion = videoVersion

	versionFlag, err := req.QueryMap.GetInt("version_flag", 0)
	checkError(req.Param.RequestID, err)
	req.Param.VersionFlag = int32(versionFlag)

	videoW, err := req.QueryMap.GetInt("video_w", 0)
	checkError(req.Param.RequestID, err)
	req.Param.VideoW = int32(videoW)

	videoH, err := req.QueryMap.GetInt("video_h", 0)
	checkError(req.Param.RequestID, err)
	req.Param.VideoH = int32(videoH)

	// rankerInfo
	powerRate, err := req.QueryMap.GetInt("power_rate", 0)
	checkError(req.Param.RequestID, err)
	req.Param.PowerRate = powerRate

	charging, err := req.QueryMap.GetInt("charging", 0)
	checkError(req.Param.RequestID, err)
	req.Param.Charging = charging

	totalMemory, err := req.QueryMap.GetString("h", false, "")
	if len(totalMemory) == 0 {
		totalMemory, err = req.QueryMap.GetString("cache1", false, "")
	}
	checkError(req.Param.RequestID, err)
	req.Param.TotalMemory = totalMemory

	residualMemory, err := req.QueryMap.GetString("i", false, "")
	if len(residualMemory) == 0 {
		residualMemory, err = req.QueryMap.GetString("cache2", false, "")
	}
	checkError(req.Param.RequestID, err)
	req.Param.ResidualMemory = residualMemory

	debug, err := req.QueryMap.GetInt("debug", 0)
	checkError(req.Param.RequestID, err)
	req.Param.Debug = debug

	debugMode, err := req.QueryMap.GetBool("debugmode", false)
	checkError(req.Param.RequestID, err)
	req.Param.DebugMode = debugMode

	channel, err := req.QueryMap.GetString("channel", false, "")
	checkError(req.Param.RequestID, err)
	req.Param.Extchannel = channel

	cc, err := req.QueryMap.GetString("cc", false, "")
	checkError(req.Param.RequestID, err)
	req.Param.CC = strings.ToUpper(cc)

	apiversion, err := req.QueryMap.GetFloat("api_version", 0.0)
	checkError(req.Param.RequestID, err)
	req.Param.ApiVersion = apiversion

	// TODO refactor hb and adnet
	apiVersionStr, err := req.QueryMap.GetString("api_version", false, "")
	checkError(req.Param.RequestID, err)
	req.Param.ApiVersionCode, _ = mvutil.IntVer(apiVersionStr)

	// interactive ads
	iarst, err := req.QueryMap.GetInt("ia_rst", 0)
	checkError(req.Param.RequestID, err)
	req.Param.IARst = iarst

	// init MwCreative
	req.Param.MwCreatvie = make([]int32, 0)

	gotk, err := req.QueryMap.GetBool("gotk", false)
	checkError(req.Param.RequestID, err)
	req.Param.Gotk = gotk

	isVast, err := req.QueryMap.GetBool("is_vast", false)
	checkError(req.Param.RequestID, err)
	req.Param.IsVast = isVast

	fallback, err := req.QueryMap.GetBool("fallback", false)
	checkError(req.Param.RequestID, err)
	req.Param.Fallback = fallback

	mainDomain, err := req.QueryMap.GetString("main_domain", true, "")
	checkError(req.Param.RequestID, err)
	req.Param.MainDomain = mainDomain

	aid, err := req.QueryMap.GetString("aid", true, "")
	checkError(req.Param.RequestID, err)
	req.Param.AID = aid

	req.Param.RequestTime = time.Now().Unix()

	// 用于小程序请求ads接口不做参数强校验
	ncp, err := req.QueryMap.GetString("ncp", true, "")
	checkError(req.Param.RequestID, err)
	req.Param.NCP = ncp

	isResInHtml, err := req.QueryMap.GetBool("rih", false)
	checkError(req.Param.RequestID, err)
	req.Param.ResInHtml = isResInHtml

	hasWx, err := req.QueryMap.GetBool("has_wx", false)
	checkError(req.Param.RequestID, err)
	req.Param.HasWX = hasWx

	// 自有id
	encryptedSysId, err := req.QueryMap.GetString("b", true, "")
	req.Param.EncryptedSysId = encryptedSysId
	var sysId string
	if len(encryptedSysId) > 0 {
		sysId, req.Param.EncryptedSysIdTimestamp = mvutil.GetOriginIdFromEncryptSysId(encryptedSysId)
	}

	if len(sysId) == 0 {
		sysId, err = req.QueryMap.GetString("sys_id", true, "")
		// 兼容开发者会滚了sdk版本的情况，他们可能会传加密后的id给我们
		if isRollbackSysId(sysId) {
			sysId, req.Param.EncryptedSysIdTimestamp = mvutil.GetOriginIdFromEncryptSysId(sysId)
		}
	}
	checkError(req.Param.RequestID, err)
	req.Param.SysId = sysId

	encryptedBkupId, err := req.QueryMap.GetString("c", true, "")
	req.Param.EncryptedBkupId = encryptedBkupId
	var bkupId string
	if len(encryptedBkupId) > 0 {
		bkupId, req.Param.EncryptedBkupIdTimestamp = mvutil.GetOriginIdFromEncryptSysId(encryptedBkupId)
	}
	if len(bkupId) == 0 {
		bkupId, err = req.QueryMap.GetString("bkup_id", true, "")
		// 兼容开发者会滚了sdk版本的情况，他们可能会传加密后的id给我们
		if isRollbackSysId(bkupId) {
			bkupId, req.Param.EncryptedBkupIdTimestamp = mvutil.GetOriginIdFromEncryptSysId(bkupId)
		}
	}
	checkError(req.Param.RequestID, err)
	req.Param.BkupId = bkupId

	oaid, err := req.QueryMap.GetString("oaid", true, "")
	checkError(req.Param.RequestID, err)
	req.Param.OAID = oaid

	ndm, err := req.QueryMap.GetInt("ndm", 0)
	checkError(req.Param.RequestID, err)
	req.Param.Ndm = ndm

	// mp adsense广告请求参数
	adSense, err := req.QueryMap.GetInt("ad_s", 0)
	checkError(req.Param.RequestID, err)
	req.Param.AdSense = adSense

	ecpmFloor, err := req.QueryMap.GetFloat("ecpm_floor", float64(0))
	checkError(req.Param.RequestID, err)
	req.Param.BidFloor = ecpmFloor

	if req.IsHBRequest && req.Param.RequestPath != mvconst.PATHBidAds {
		// hb bid_floor
		bidFloor, err := req.QueryMap.GetFloat("bid_floor", 0.0)
		checkError(req.Param.RequestID, err)
		req.Param.BidFloor = bidFloor
	}

	plPkg, err := req.QueryMap.GetString("plpkg", false, "")
	checkError(req.Param.RequestID, err)
	req.Param.PlPkg = plPkg

	mof, err := req.QueryMap.GetInt("mof", 0)
	checkError(req.Param.RequestID, err)
	req.Param.Mof = mof

	mofData, err := req.QueryMap.GetString("mof_data", false, "")
	checkError(req.Param.RequestID, err)
	req.Param.MofData = mofData

	parentUnitId, err := req.QueryMap.GetInt64("parent_unit", 0)
	checkError(req.Param.RequestID, err)
	req.Param.ParentUnitId = parentUnitId

	// 固定unit对应的parent_unit
	ucParentUnitId, err := req.QueryMap.GetInt64("uc_parent_unit", 0)
	checkError(req.Param.RequestID, err)
	req.Param.UcParentUnitId = ucParentUnitId

	h5Type, err := req.QueryMap.GetInt("h5_type", 0)
	checkError(req.Param.RequestID, err)
	req.Param.H5Type = h5Type

	mofVersion, err := req.QueryMap.GetInt("mof_ver", 0)
	checkError(req.Param.RequestID, err)
	req.Param.MofVersion = mofVersion

	reqType, err := req.QueryMap.GetString("req_type", false, "")
	checkError(req.Param.RequestID, err)
	req.Param.ReqType = reqType

	newDisplayInfoStr, err := req.QueryMap.GetString("d", false, "")
	checkError(req.Param.RequestID, err)
	var newDisplayInfo mvutil.NewDisplayInfos
	err = jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal([]byte(newDisplayInfoStr), &newDisplayInfo)
	if err == nil {
		req.Param.NewDisplayInfoData = newDisplayInfo
	}

	displayInfoStr, err := req.QueryMap.GetString("display_info", false, "")
	checkError(req.Param.RequestID, err)
	var displayInfo mvutil.DisplayInfos
	err = jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal([]byte(displayInfoStr), &displayInfo)
	if err == nil {
		req.Param.DisplayInfoData = displayInfo
	}

	mofType, err := req.QueryMap.GetInt("mof_type", 0)
	checkError(req.Param.RequestID, err)
	req.Param.MofType = mofType

	asTestMode, err := req.QueryMap.GetInt("ast_mode", 0)
	checkError(req.Param.RequestID, err)
	req.Param.AsTestMode = int32(asTestMode)

	debugModeTimeout, err := req.QueryMap.GetInt("debug_timeout", 0)
	checkError(req.Param.RequestID, err)
	req.Param.DebugModeTimeout = debugModeTimeout

	h5Data, err := req.QueryMap.GetString("h5_data", false, "")
	checkError(req.Param.RequestID, err)
	req.Param.H5Data = h5Data

	androidIDMd5, err := req.QueryMap.GetString("aid_md5", false, "")
	checkError(req.Param.RequestID, err)
	req.Param.AndroidIDMd5 = androidIDMd5

	gaidMd5, err := req.QueryMap.GetString("gaid_md5", false, "")
	checkError(req.Param.RequestID, err)
	req.Param.GAIDMd5 = gaidMd5

	idfaMd5, err := req.QueryMap.GetString("idfa_md5", false, "")
	checkError(req.Param.RequestID, err)
	req.Param.IDFAMd5 = idfaMd5

	randNum, err := req.QueryMap.GetInt("rand_num", 0)
	checkError(req.Param.RequestID, err)
	req.Param.RandNum = int32(randNum)

	dmt, err := req.QueryMap.GetFloat("dmt", 0.0)
	checkError(req.Param.RequestID, err)
	req.Param.Dmt = dmt

	dmf, err := req.QueryMap.GetFloat("dmf", 0.0)
	checkError(req.Param.RequestID, err)
	req.Param.Dmf = dmf

	cpuType, err := req.QueryMap.GetString("ct", false, "")
	checkError(req.Param.RequestID, err)
	req.Param.Ct = cpuType

	channelInfo, err := req.QueryMap.GetString("ch_info", false, "")
	// 解析channelInfo
	channelInfo = mvutil.Base64Decode(channelInfo)
	checkError(req.Param.RequestID, err)
	req.Param.ChannelInfo = channelInfo

	placementId, err := req.QueryMap.GetInt64("placement_id", 0)
	checkError(req.Param.RequestID, err)
	req.Param.PlacementId = placementId

	testMode, err := req.QueryMap.GetInt("test", 0)
	checkError(req.Param.RequestID, err)
	req.Param.HBBidTestMode = int32(testMode)

	// 是否需要按照返回线性可跳过模式返回内容。不传则默认返回linear ads。
	adFormat, err := req.QueryMap.GetInt("ad_format", 0)
	checkError(req.Param.RequestID, err)
	req.Param.AdFormat = adFormat

	// rwPlus 是否开启Reward Plus的开关
	rwPlus, err := req.QueryMap.GetInt("rw_plus", 0)
	checkError(req.Param.RequestID, err)
	req.Param.RwPlus = int32(rwPlus)

	webEnvStr, err := req.QueryMap.GetString("web_env", false, "")
	checkError(req.Param.RequestID, err)
	var webEnv mvutil.WebEnv
	err = jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal([]byte(webEnvStr), &webEnv)
	if err == nil {
		req.Param.WebEnvData = webEnv
	}

	tmpIdsSt, err := req.QueryMap.GetString("g", false, "")
	if len(tmpIdsSt) == 0 {
		tmpIdsSt, err = req.QueryMap.GetString("tmp_ids", false, "")
	}
	if len(tmpIdsSt) > 0 {
		req.Param.TmpIds = make(map[string]bool)
		tmpIdsArr := []string{}
		err = jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal([]byte(tmpIdsSt), &tmpIdsArr)
		if err == nil {
			for _, v := range tmpIdsArr {
				req.Param.TmpIds[v] = true
			}
		}
	}

	open, err := req.QueryMap.GetInt("open", 0)
	checkError(req.Param.RequestID, err)
	req.Param.Open = open

	ts, err := req.QueryMap.GetString("ts", true, "")
	checkError(req.Param.RequestID, err)
	req.Param.TS = ts

	st, err := req.QueryMap.GetString("st", true, "")
	checkError(req.Param.RequestID, err)
	req.Param.SignWithTimeStamp = st

	skadStr, err := req.QueryMap.GetString("skad", false, "")
	checkError(req.Param.RequestID, err)
	var skad mvutil.Skadnetwork
	err = jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal([]byte(skadStr), &skad)
	if err == nil {
		req.Param.Skadnetwork = &skad
	}
	// 先不处理这个参数
	// tkiStr, err := req.QueryMap.GetString("tki", false, "")
	// checkError(req.Param.RequestID, err)
	// var tki mvutil.TrackingInfo
	// err = jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal([]byte(tkiStr), &tki)
	// if err == nil {
	//	req.Param.TrackingInfo = &tki
	// }

	bandWidth, err := req.QueryMap.GetInt64("band_width", 0)
	checkError(req.Param.RequestID, err)
	req.Param.BandWidth = bandWidth

	newTki, err := req.QueryMap.GetString("f", false, "")
	checkError(req.Param.RequestID, err)
	req.Param.NewTKI = newTki

	tki, err := req.QueryMap.GetString("tki", false, "")
	checkError(req.Param.RequestID, err)
	req.Param.TKI = tki

	towardAdx, err := req.QueryMap.GetInt("toward_adx", 0)
	if towardAdx == 1 {
		req.Param.IsTowardAdx = true
	}

	deeplinkType, err := req.QueryMap.GetInt("deeplink_type", 0)
	checkError(req.Param.RequestID, err)
	if deeplinkType < 0 || deeplinkType > 2 {
		deeplinkType = 0
	}
	req.Param.DeeplinkType = deeplinkType

	htmlSupport, err := req.QueryMap.GetInt("html_support", 0)
	checkError(req.Param.RequestID, err)
	if htmlSupport != 0 && htmlSupport != 1 {
		htmlSupport = 0
	}
	req.Param.HtmlSupport = htmlSupport

	// limit_trk 是否限制广告追踪
	limitTrk, err := req.QueryMap.GetString("limit_trk", false, "")
	checkError(req.Param.RequestID, err)
	req.Param.LimitTrk = limitTrk

	// att app 级别 idfa 授权状态
	att, err := req.QueryMap.GetString("att", false, "")
	checkError(req.Param.RequestID, err)
	req.Param.Att = att

	// brt  屏幕亮度
	brt, err := req.QueryMap.GetString("brt", false, "")
	checkError(req.Param.RequestID, err)
	req.Param.Brt = brt

	// vol  音量
	vol, err := req.QueryMap.GetString("vol", false, "")
	checkError(req.Param.RequestID, err)
	req.Param.Vol = vol

	// lpm  是否为低电量模式
	lpm, err := req.QueryMap.GetString("lpm", false, "")
	checkError(req.Param.RequestID, err)
	req.Param.Lpm = lpm

	// font  设备默认字体大小
	font, err := req.QueryMap.GetString("font", false, "")
	checkError(req.Param.RequestID, err)
	req.Param.Font = font

	needCreativeDataCIds, err := req.QueryMap.GetString("need_cr_data_ids", false, "")
	checkError(req.Param.RequestID, err)
	req.Param.NeedCreativeDataCIds = needCreativeDataCIds

	debugDspId, err := req.QueryMap.GetInt("debug_dsp_id", 0)
	checkError(req.Param.RequestID, err)
	req.Param.DebugDspId = debugDspId

	mpsDebug, err := req.QueryMap.GetString("mapping_server_debug", false, "off")
	checkError(req.Param.RequestID, err)
	req.Param.MappingServerDebug = mpsDebug

	pcmReportendpoint, err := req.QueryMap.GetString("pcm_rp", false, "")
	checkError(req.Param.RequestID, err)
	req.Param.PcmReportendpoint = pcmReportendpoint

	fwType, err := req.QueryMap.GetInt("fw_type", 0)
	checkError(req.Param.RequestID, err)
	req.Param.FwType = fwType

	hardwareModel, err := req.QueryMap.GetString("h_model", false, "")
	checkError(req.Param.RequestID, err)
	req.Param.HardwareModel = hardwareModel

	targetIds, err := req.QueryMap.GetString("target_ids", false, "")
	checkError(req.Param.RequestID, err)
	req.Param.TargetIds = targetIds

	j, err := req.QueryMap.GetString("j", false, "")
	checkError(req.Param.RequestID, err)
	req.Param.TrafficInfo = j

	dspMoreOfferInfo, err := req.QueryMap.GetString("dsp_mof_info", false, "")
	checkError(req.Param.RequestID, err)
	req.Param.DspMoreOfferInfo = dspMoreOfferInfo

	gdprConsent, err := req.QueryMap.GetString("gdpr_consent", false, "")
	checkError(req.Param.RequestID, err)
	req.Param.GdprConsent = gdprConsent

	dspMof, err := req.QueryMap.GetInt("dsp_mof", 0)
	checkError(req.Param.RequestID, err)
	req.Param.DspMof = dspMof

	appSettingId, err := req.QueryMap.GetString("a_stid", false, "")
	checkError(req.Param.RequestID, err)
	req.Param.AppSettingId = appSettingId

	unitSettingId, err := req.QueryMap.GetString("u_stid", false, "")
	checkError(req.Param.RequestID, err)
	req.Param.UnitSettingId = unitSettingId

	rewardSettingId, err := req.QueryMap.GetString("r_stid", false, "")
	checkError(req.Param.RequestID, err)
	req.Param.RewardSettingId = rewardSettingId

	miSkSpt, err := req.QueryMap.GetString("misk_spt", false, "")
	checkError(req.Param.RequestID, err)
	req.Param.MiSkSpt = miSkSpt

	miSkSptDet, err := req.QueryMap.GetString("misk_spt_det", false, "")
	checkError(req.Param.RequestID, err)
	req.Param.MiSkSptDet = miSkSptDet

	// more offer 才有的。pioneer根据这个查询缓存的单子
	parentId, err := req.QueryMap.GetString("mof_parent_id", false, "")
	checkError(req.Param.RequestID, err)
	req.Param.ParentId = parentId

	// hm_info
	encryptHarmonyInfo, err := req.QueryMap.GetString("hm_info", false, "")
	checkError(req.Param.RequestID, err)
	req.Param.EncryptHarmonyInfo = encryptHarmonyInfo

	// dnt值为1时，表示用户退出个性化广告
	dnt, err := req.QueryMap.GetString("dnt", false, "")
	checkError(req.Param.RequestID, err)
	req.Param.Dnt = dnt

	// more_offer的主unit的ad_type
	parentAdType, err := req.QueryMap.GetString("parent_ad_type", false, "")
	checkError(req.Param.RequestID, err)
	req.Param.ParentAdType = parentAdType

	// more_offer的主unit 的adx名字 dsp流量才会有
	parentExchange, err := req.QueryMap.GetString("parent_exchange", false, "")
	checkError(req.Param.RequestID, err)
	req.Param.ParentExchange = parentExchange

	return nil
}

func isRollbackSysId(id string) bool {
	return len(id) != 36 || !strings.Contains(id, "-")
}

// 将req的post里的参数解析到bidRequest这个数据结构中去，再讲bidRequest中的参数解析到req中的Param中
func renderOpenRTBRequestParam(req *mvutil.RequestParams) error {
	if req == nil {
		return filter.RenderBidReqDateError
	}
	bidRequest := &mtg_hb_rtb.BidRequest{}
	err := jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(req.PostData, bidRequest)
	if err != nil {
		return filter.DecodeBidDataError
	}
	req.Param.HBS2SBidID = bidRequest.GetId()
	// TODO 啥，Imp
	imp := bidRequest.GetImp()
	if len(imp) == 0 {
		return filter.BidRequestImpEmpty
	}

	impObj := imp[0]
	req.Param.ImpID = impObj.GetId()
	req.Param.SDKVersion = impObj.GetDisplaymanagerver()
	req.Param.MediationName = impObj.GetDisplaymanager()
	// rtb imp.displaymanager mapping to mtg channel
	chaIds := extractor.GetMEDIATION_CHANNEL_ID()
	if chaId, ok := chaIds[strings.ToUpper(req.Param.MediationName)]; ok {
		req.Param.Extchannel = chaId
	}
	req.Param.BidFloor = impObj.GetBidfloor()
	req.Param.HTTPReq = 2

	if req.Param.RequestPath == mvconst.PATHMopubBid {
		if impObj.GetExt() == nil || impObj.GetExt().GetNetworkids() == nil {
			return filter.BidRequestInvalidate
		}
		// replace unitId and appID
		appId, err := strconv.ParseInt(impObj.GetExt().GetNetworkids().GetAppid(), 10, 64)
		if err != nil || appId <= 0 {
			return filter.BidRequestAppInValidate
		}
		req.Param.AppID = appId
		unitId, err := strconv.ParseInt(impObj.GetExt().GetNetworkids().GetPlacementid(), 10, 64)
		if err != nil || unitId <= 0 {
			return filter.BidRequestUnitInValidate
		}
		req.Param.UnitID = unitId
	} else {
		unitId, err := strconv.ParseInt(impObj.GetTagid(), 10, 64)
		if err != nil || unitId <= 0 {
			return filter.BidRequestUnitInValidate
		}

		req.Param.UnitID = unitId
		appId, err := strconv.ParseInt(bidRequest.GetApp().GetId(), 10, 64)
		if err != nil || appId <= 0 {
			return filter.BidRequestAppInValidate
		}
		req.Param.AppID = appId
	}
	// bidReq参数赋值到req.Param
	req.Param.HBBidTestMode = bidRequest.GetTest()
	req.Param.DebugMode = bidRequest.GetDebugMode()
	req.Param.AsTestMode = bidRequest.GetAstMode()
	req.Param.AppVersionName = bidRequest.GetApp().GetVer()
	req.Param.Orientation = int(bidRequest.GetApp().GetExt().GetOrientation())

	req.Param.UserAgent = bidRequest.GetDevice().GetUa()
	clientIp := bidRequest.GetDevice().GetIp()
	if len(clientIp) <= 0 {
		clientIp = bidRequest.GetDevice().GetIpv6()
	}
	req.Param.ClientIP = clientIp
	req.Param.DeviceType = int(bidRequest.GetDevice().GetDevicetype())
	req.Param.Brand = bidRequest.GetDevice().GetMake()
	req.Param.Model = bidRequest.GetDevice().GetModel()
	req.Param.Os = strings.ToLower(bidRequest.GetDevice().GetOs())
	req.Param.Platform = helpers.GetPlatform(strings.ToLower(req.Param.Os))
	req.Param.OSVersion = bidRequest.GetDevice().GetOsv()
	req.Param.OSVersionCode, _ = mvutil.IntVer(req.Param.OSVersion)
	req.Param.Language = bidRequest.GetDevice().GetLanguage()
	req.Param.ScreenWidth = int(bidRequest.GetDevice().GetW())
	req.Param.ScreenHeigh = int(bidRequest.GetDevice().GetH())
	if bidRequest.GetDevice().GetGeo() != nil {
		req.Param.BidDeviceGEOCC = strings.ToUpper(bidRequest.GetDevice().GetGeo().GetCountry())
	}
	// impObj数据赋值到req.Param中
	if impObj.GetNative() != nil {
		if impObj.GetNative().GetExt() != nil {
			if impObj.GetNative().GetExt().GetUnitSizeH() > 0 && impObj.GetNative().GetExt().GetUnitSizeW() > 0 {
				req.Param.UnitSize = strconv.Itoa(int(impObj.GetNative().GetExt().GetUnitSizeW())) + "x" + strconv.Itoa(int(impObj.GetNative().GetExt().GetUnitSizeH()))
			}
		}
	}
	if impObj.GetBanner() != nil {
		req.Param.SdkBannerUnitWidth = int64(impObj.GetBanner().GetW())
		req.Param.SdkBannerUnitHeight = int64(impObj.GetBanner().GetH())
	}
	if impObj.GetVideo() != nil {
		req.Param.VideoW = impObj.GetVideo().GetW()
		req.Param.VideoH = impObj.GetVideo().GetH()
	}
	if impObj.GetExt() != nil && len(impObj.GetExt().GetChInfo()) > 0 {
		req.Param.ChannelInfo = impObj.GetExt().GetChInfo()
	}
	if impObj.GetExt() != nil && impObj.GetExt().GetSkadn() != nil && len(impObj.GetExt().GetSkadn().GetVersion()) > 0 {
		req.Param.Skadnetwork = &mvutil.Skadnetwork{
			Ver:      impObj.GetExt().GetSkadn().GetVersion(),
			Adnetids: impObj.GetExt().GetSkadn().GetSkadnetids(),
		}
	}

	// 用户是否退出个性化广告
	if bidRequest.GetDevice().GetDnt() > 0 {
		req.Param.Dnt = strconv.Itoa(int(bidRequest.GetDevice().GetDnt()))
	}
	req.Param.HBTmax = bidRequest.GetTmax()
	if placementId, err := strconv.ParseInt(impObj.GetPlacementId(), 10, 64); err == nil {
		req.Param.PlacementId = placementId
	}
	req.Param.NetworkType = renderMTGConnType(bidRequest.GetDevice().GetConnectiontype())
	/* IDFA，英文全称 Identifier for Advertising ，
	可以理解为广告id，苹果公司提供的用于追踪用户的广告标识符，可以用来打通不同app之间的广告。
	每个设备只有一个IDFA，不同APP在同一设备上获取IDFA的结果是一样的
	 GAID是安卓的广告标识符
	*/
	if req.Param.Platform == constant.IOS {
		req.Param.IDFA = bidRequest.GetDevice().GetIfa()
	} else {
		req.Param.GAID = bidRequest.GetDevice().GetIfa()
	}
	// MCC,Mobile Country Code,移动国家代码。它由三位数字组成,用于标识一个国家
	// MNC,Mobile Network Code,移动网络代码。它由二到三位数字组成。它和 MCC 合在一起唯一标识一个移动网络提供者,如中国移动、中国联通
	mmdata := strings.Split(bidRequest.GetDevice().GetMccmnc(), "-")
	if len(mmdata) == 2 {
		req.Param.MCC = mmdata[0]
		req.Param.MNC = mmdata[1]
	}

	req.Param.Carrier = bidRequest.GetDevice().GetCarrier()

	if bidRequest.Regs.GetExt() != nil {
		req.Param.Ccpa = bidRequest.Regs.GetExt().GetCcpa()
	}

	if bidRequest.GetSource() != nil {
		rtbSourceBytes, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(bidRequest.GetSource())
		req.Param.RTBSourceBytes = rtbSourceBytes
	}

	buyeruid := bidRequest.GetUser().GetBuyeruid()
	err = parseBuyeruid(buyeruid, req)
	if err != nil {
		return err
	}
	// process SDKVersion
	req.Param.SDKVersion = strings.ToLower(req.Param.SDKVersion)
	req.Param.FormatSDKVersion = supply_mvutil.RenderSDKVersion(req.Param.SDKVersion)

	req.Param.RequestTime = time.Now().Unix()
	return nil
}

func renderMTGConnType(rtbConnType mtg_hb_rtb.BidRequest_Device_ConnectionType) int {
	switch rtbConnType {
	case mtg_hb_rtb.BidRequest_Device_WIFI:
		return constant.TWIFI
	case mtg_hb_rtb.BidRequest_Device_CELL_4G:
		return constant.T4G
	case mtg_hb_rtb.BidRequest_Device_CELL_3G:
		return constant.T3G
	case mtg_hb_rtb.BidRequest_Device_CELL_2G:
		return constant.T2G
	default:
		return constant.TNUNKNOWN
	}
}

func parseBuyeruid(uid string, req *mvutil.RequestParams) error {
	// http://confluence.mobvista.com/pages/viewpage.action?pageId=21466844#Header-Bidding%E6%8E%A5%E5%8F%A3%E6%96%87%E6%A1%A3-2.2.5Object%EF%BC%9AUser
	if len(uid) == 0 || req == nil {
		return filter.BuyerUidEmpty
	}
	uid = strings.Replace(uid, " ", "+", -1)
	uid = strings.Replace(uid, "\\", "", -1)
	uid = helpers.Base64Decode(uid)
	uidData := strings.SplitN(uid, "|", 33)
	if len(uidData) < 11 {
		return filter.BuyerUidDataInvalidate
	}

	if len(uidData[8]) > 0 {
		req.Param.SDKVersion = uidData[8]
	}

	if len(uidData[9]) > 0 {
		req.Param.ScreenSize = uidData[9]
	}

	if len(uidData[10]) > 0 {
		req.Param.UserAgent = uidData[10]
	}

	if len(uidData[5]) > 0 {
		req.Param.Brand = uidData[5]
	}

	if len(uidData[6]) > 0 {
		req.Param.Model = uidData[6]
	}

	if len(uidData[3]) > 0 {
		req.Param.SysId = uidData[3]
		if isRollbackSysId(req.Param.SysId) {
			req.Param.SysId, req.Param.EncryptedSysIdTimestamp = mvutil.GetOriginIdFromEncryptSysId(req.Param.SysId)
		}
	}

	if len(uidData[4]) > 0 {
		req.Param.BkupId = uidData[4]
		if isRollbackSysId(req.Param.BkupId) {
			req.Param.BkupId, req.Param.EncryptedBkupIdTimestamp = mvutil.GetOriginIdFromEncryptSysId(req.Param.BkupId)
		}
	}

	if len(uidData[7]) > 0 {
		conType, err := strconv.Atoi(uidData[7])
		if err == nil {
			req.Param.NetworkType = conType
		}
	}

	if req.Param.Platform == constant.IOS {
		if len(uidData[0]) > 0 {
			req.Param.IDFA = uidData[0]
		}
		if len(uidData[1]) > 0 {
			req.Param.IDFV = uidData[1]
		}
	}

	if req.Param.Platform == constant.Android {
		if len(uidData[0]) > 0 {
			req.Param.GAID = uidData[0]
		}
		if len(uidData[1]) > 0 {
			req.Param.AndroidID = uidData[1]
		}
		if len(uidData[2]) > 0 {
			req.Param.IMEI = uidData[2]
		}
	}

	installIds := uidData[11:12]
	if len(installIds) > 0 {
		req.Param.InstallIDS = installIds[0]
	}
	excludeIds := uidData[12:13]
	if len(excludeIds) > 0 {
		req.Param.ExcludeIDS = excludeIds[0]
	}
	tokenTimeStamp := uidData[13:14]
	if len(tokenTimeStamp) > 0 {
		req.Param.TokenTimeStamp = tokenTimeStamp[0]
	}
	apiVersion := uidData[14:15]
	if len(apiVersion) > 0 {
		req.Param.ApiVersion, _ = strconv.ParseFloat(apiVersion[0], 64)
		req.Param.ApiVersionCode, _ = mvutil.IntVer(apiVersion[0])
	}
	dmt := uidData[15:16]
	if len(dmt) > 0 {
		if v, err := strconv.ParseFloat(dmt[0], 64); err == nil {
			req.Param.Dmt = v
		}
	}
	dmf := uidData[16:17]
	if len(dmf) > 0 {
		if v, err := strconv.ParseFloat(dmf[0], 64); err == nil {
			req.Param.Dmf = v
		}
	}
	ct := uidData[17:18]
	if len(ct) > 0 {
		req.Param.Ct = ct[0]
	}
	powerRate := uidData[18:19]
	if len(powerRate) > 0 {
		if v, err := strconv.Atoi(powerRate[0]); err == nil {
			req.Param.PowerRate = v
		}
	}
	charging := uidData[19:20]
	if len(charging) > 0 {
		if v, err := strconv.Atoi(charging[0]); err == nil {
			req.Param.Charging = v
		}
	}
	totalMemory := uidData[20:21]
	if len(totalMemory) > 0 {
		req.Param.TotalMemory = totalMemory[0]
	}
	residualMemory := uidData[21:22]
	if len(residualMemory) > 0 {
		req.Param.ResidualMemory = residualMemory[0]
	}

	upTime := uidData[22:23]
	if len(upTime) > 0 {
		req.Param.UpTime = upTime[0]
	}

	skVersion := uidData[23:24]
	if len(skVersion) > 0 && len(skVersion[0]) > 0 {
		if req.Param.Skadnetwork == nil {
			req.Param.Skadnetwork = &mvutil.Skadnetwork{}
		}
		req.Param.Skadnetwork.Ver = skVersion[0]
	}

	skTag := uidData[24:25]
	if len(skTag) > 0 && len(skTag[0]) > 0 {
		if req.Param.Skadnetwork == nil {
			req.Param.Skadnetwork = &mvutil.Skadnetwork{}
		}
		req.Param.Skadnetwork.Tag = skTag[0]
	}

	// 自有系统id 加密后的sys_id
	encryptedSysId := uidData[25:26]
	if len(encryptedSysId) > 0 && len(encryptedSysId[0]) > 0 {
		req.Param.EncryptedSysId = encryptedSysId[0]
		req.Param.SysId, req.Param.EncryptedSysIdTimestamp = mvutil.GetOriginIdFromEncryptSysId(encryptedSysId[0])
	}
	encryptedBkupId := uidData[26:27]
	if len(encryptedBkupId) > 0 && len(encryptedBkupId[0]) > 0 {
		req.Param.EncryptedBkupId = encryptedBkupId[0]
		req.Param.BkupId, req.Param.EncryptedBkupIdTimestamp = mvutil.GetOriginIdFromEncryptSysId(encryptedBkupId[0])
	}

	// hardward model
	hardModel := uidData[27:28]
	if len(hardModel) > 0 && len(hardModel[0]) > 0 {
		req.Param.HardwareModel = hardModel[0]
	}

	// setting id
	appSettingId := uidData[28:29]
	if len(appSettingId) > 0 && len(appSettingId[0]) > 0 {
		req.Param.AppSettingId = appSettingId[0]
	}

	miskSpt := uidData[29:30]
	if len(miskSpt) > 0 && len(miskSpt[0]) > 0 {
		req.Param.MiSkSpt = miskSpt[0]
	}

	miskSptDet := uidData[30:31]
	if len(miskSptDet) > 0 && len(miskSptDet[0]) > 0 {
		req.Param.MiSkSptDet = miskSptDet[0]
	}

	cachedCampaignIds := uidData[31:32]
	if len(cachedCampaignIds) > 0 && len(cachedCampaignIds[0]) > 0 {
		req.Param.CachedCampaignIds = cachedCampaignIds[0]
		req.Param.ExtDataInit.MultiVcn = 1
		req.Param.ExtDataInit.VcnCampaigns = cachedCampaignIds[0]
	}

	oaid := uidData[32:33]
	if len(oaid) > 0 && len(oaid[0]) > 0 {
		req.Param.OAID = oaid[0]
	}

	// fix iOS SDK bug, as follow
	// idfa|||sys_id|bkup_id|brand|model|networktype|sdk_version|screen_size|user_agent|exclude_id|token_time_stamp|install_id
	if len(apiVersion) == 0 && req.Param.Platform == constant.IOS {
		excludeIds = uidData[11:12]
		if len(excludeIds) > 0 {
			req.Param.ExcludeIDS = excludeIds[0]
		}
		tokenTimeStamp = uidData[12:13]
		if len(tokenTimeStamp) > 0 {
			req.Param.TokenTimeStamp = tokenTimeStamp[0]
		}
		installIds = uidData[13:14]
		if len(installIds) > 0 {
			req.Param.InstallIDS = installIds[0]
		}
	}
	return nil
}

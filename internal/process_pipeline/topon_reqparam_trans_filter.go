package process_pipeline

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"gitlab.mobvista.com/ADN/adnet/internal/extractor"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/constant"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/filter"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/helpers"
	"gitlab.mobvista.com/ADN/adnet/internal/mvconst"
	"gitlab.mobvista.com/ADN/adnet/internal/watcher"
	supply_mvutil "gitlab.mobvista.com/ADN/chasm/module/supply/mvutil"
	openrtb "gitlab.mobvista.com/ADN/mtg_openrtb/pkg/openrtb_v2"
	//"gitlab.mobvista.com/adserver/recommend_protocols/go/ad_server"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"time"

	"gitlab.mobvista.com/ADN/adnet/internal/errorcode"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
)

type ToponReqparamTransFilter struct {
}

func (tf *ToponReqparamTransFilter) Process(data interface{}) (interface{}, error) {
	in, ok := data.(*http.Request)
	if !ok {
		return nil, errorcode.EXCEPTION_SPECIAL_EMPTY_STRING
	}
	body, _ := ioutil.ReadAll(in.Body)
	r := mvutil.RequestParams{}
	r.PostData = body
	r.IsTopon = true
	r.Param.RequestTime = time.Now().Unix()
	req := new(openrtb.BidRequest)

	err := proto.Unmarshal(r.PostData, req)
	// TODO 解析proto
	if err != nil {
		return nil, errorcode.EXCEPTION_SPECIAL_EMPTY_STRING
	}
	renderToponReqParam(in, &r, req)
	if req.Mvext == nil {
		req.Mvext = new(openrtb.BidRequest_Ext)
	}
	r.ToponRequest = req

	if r.Param.Platform != constant.Android && r.Param.Platform != constant.IOS {
		return r, filter.AppPlatformError
	}
	return &r, nil
}

// renderToponReqParam 从请求中解析得到所需的r.Params
func renderToponReqParam(in *http.Request, r *mvutil.RequestParams, req *openrtb.BidRequest) {
	r.Param.RequestPath = in.URL.Path
	r.Param.RequestURI = in.RequestURI
	device := req.GetDevice()
	app := req.GetApp()
	imp := req.GetImp()[0]
	clientIp := device.GetIp()
	r.Param.ParamCIP = clientIp
	r.Param.ClientIP = clientIp
	if len(r.Param.ClientIP) == 0 {
		// 请求没有传client_ip，使用http头部解析的ip
		r.Param.ClientIP = GetClientIP(in)
	} else {
		ip := net.ParseIP(r.Param.ClientIP)
		if ip == nil {
			// 请求client_ip字段为非法ip
			r.Param.ClientIP = GetClientIP(in)
		}
	}
	r.Param.UserAgent = device.GetUa()
	r.Param.Extra9 = r.Param.UserAgent
	r.Param.ExtsystemUseragent = mvutil.RawUrlEncode(in.Header.Get("User-Agent"))
	r.Param.MvLine = in.Header.Get("Mv-Line")
	// request_param_filter
	r.Param.RequestID = mvutil.GetGoTkClickID()
	connType := device.GetConnectiontype()
	r.Param.NetworkType = topOnConnTypeToNetType(connType)
	r.Param.RequestID = req.GetId()
	os := device.GetOs()
	r.Param.Platform = helpers.GetPlatform(os)
	if r.Param.Platform == constant.IOS { // 写死为高版本
		r.Param.SDKVersion = "mi_99.0.0"
	} else {
		r.Param.SDKVersion = "mal_99.0.0"
	}
	imp.Displaymanagerver = r.Param.SDKVersion
	r.Param.FormatSDKVersion = supply_mvutil.RenderSDKVersion(r.Param.SDKVersion)
	r.Param.Os = os
	r.Param.OSVersion = device.GetOsv()
	r.Param.OSVersionCode, _ = mvutil.IntVer(r.Param.OSVersion)
	r.Param.PackageName = app.GetBundle()
	// TODO appVerseionName, appVersionCode 是否要加
	r.Param.Brand = device.GetMake()
	r.Param.Model = device.GetModel()
	r.Param.RModel = r.Param.Model

	// TODO 没有mnc, mcc
	r.Param.Language = device.GetLanguage()
	// TODO 没有timezone, sdkversion, gpversion
	r.Param.LAT = strconv.FormatFloat(device.GetGeo().GetLat(), 'f', 2, 64)
	r.Param.LNG = strconv.FormatFloat(device.GetGeo().GetLon(), 'f', 2, 64)
	// TODO 没有gpst, gps_accuracy, gps_type, d1 d2 d3 appid
	r.Param.UnitID, _ = strconv.ParseInt(imp.GetTagid(), 10, 64)
	r.Param.AdNum = 1
	// TODO 没有pingMode, display_cids, exclude_ids, offset
	r.Param.SessionID = mvutil.GetRequestID()
	r.Param.ParentSessionID = mvutil.GetRequestID()
	// TODO 没有 network, impression_image, ad_source_id
	r.Param.PriceFloor = imp.GetBidfloor() // 传入美元
	r.Param.CC = device.GetGeo().GetCountry()
	r.Param.PublisherID, _ = strconv.ParseInt(app.GetPublisher().GetId(), 10, 64)
	r.Param.AppID, _ = strconv.ParseInt(app.GetId(), 10, 64)
	r.Param.ImageSizeID = 1
	r.Param.SDKVersion = imp.GetDisplaymanagerver()
	if r.Param.Platform == constant.IOS {
		r.Param.IDFA = device.GetIfa()
	} else {
		r.Param.GAID = device.GetIfa()
	}
	r.Param.ToponTemplateSupportVersion = imp.GetMvext().GetTemplateSupportVersion()
	if ext := device.GetMvext(); ext != nil {
		r.Param.IMEI = ext.GetImei()
		r.Param.AndroidID = ext.GetAndroidId()
		r.Param.IDFV = ext.GetIfv()
		r.Param.SysId = ext.GetSystemId()
		r.Param.BkupId = ext.GetSysbkupId()
		r.Param.UpTime = ext.GetUptime()
		r.Param.OsvUpTime = ext.GetOsvUpTime()
		r.Param.Ram = ext.GetRam()
		r.Param.TotalMemory = ext.GetTotalMemory()
		r.Param.TimeZone = ext.GetTimeZone()
		r.Param.OAID = ext.GetOaid()
	}
	r.Param.MAC = device.GetMacmd5() // 传入的有可能是mac原文

	adType, _ := mvutil.GetOpenrtbAdType(imp)
	adnetAdType := mvutil.GetAdnetAdType(adType)
	var w, h int32
	if b := imp.GetBanner(); b != nil {
		w = b.GetW()
		h = b.GetH()
		r.Param.SdkBannerUnitWidth = int64(w)
		r.Param.SdkBannerUnitHeight = int64(h)
	}
	if v := imp.GetVideo(); v != nil {
		w = v.GetW()
		h = v.GetH()
		r.Param.VideoW = w
		r.Param.VideoH = h
	}
	r.Param.AdType = int32(adnetAdType)
	r.Param.AdxBidFloor = imp.GetBidfloor()
	r.Param.Scenario = constant.OpenApi

	if w >= h {
		r.Param.Orientation = mvconst.ORIENTATION_LANDSCAPE
	} else {
		r.Param.Orientation = mvconst.ORIENTATION_PORTRAIT
	}
	r.Param.ScreenSize = fmt.Sprintf("%dx%d", w, h)
	//r.Param.Category = int(ad_server.Category_APPLICATION)
	r.Param.GPVersion = "1"
	r.Param.Domain = extractor.GetDOMAIN_TRACK()
	r.Param.Sign = "NO_CHECK_SIGN"
	r.Nbr = -1
	r.Param.OnlyImpression = 1
	r.Param.VideoVersion = "1.0"
	r.Param.VersionFlag = 1
	r.Param.PingMode = 1
	r.Param.TNum = 1

	apiVersion := "1.5"
	r.Param.ApiVersion, _ = strconv.ParseFloat(apiVersion, 64)
	r.Param.ApiVersionCode, _ = mvutil.IntVer(apiVersion)
	// topon 固定使用https
	r.Param.HTTPReq = 2

	// 统计rawRequest
	watcher.AddWatchValue("raw_request", float64(1))

}

func topOnConnTypeToNetType(connType openrtb.BidRequest_Device_ConnectionType) int {
	switch connType {
	case openrtb.BidRequest_Device_CELL_2G:
		return mvconst.NETWORK_TYPE_2G
	case openrtb.BidRequest_Device_CELL_3G:
		return mvconst.NETWORK_TYPE_3G
	case openrtb.BidRequest_Device_CELL_4G:
		return mvconst.NETWORK_TYPE_4G
	case openrtb.BidRequest_Device_WIFI:
		return mvconst.NETWORK_TYPE_WIFI
	}
	return mvconst.NETWORK_TYPE_UNKNOWN
}

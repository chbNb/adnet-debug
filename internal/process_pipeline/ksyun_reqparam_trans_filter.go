package process_pipeline

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gogo/protobuf/proto"
	"gitlab.mobvista.com/ADN/adnet/internal/errorcode"
	"gitlab.mobvista.com/ADN/adnet/internal/extractor"
	"gitlab.mobvista.com/ADN/adnet/internal/mvconst"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	"gitlab.mobvista.com/ADN/adnet/internal/protobuf"
)

type KSYUNReqParamTransFilter struct {
}

func (ksrptf *KSYUNReqParamTransFilter) Process(data interface{}) (interface{}, error) {
	in, ok := data.(*http.Request)
	if !ok {
		return nil, errorcode.EXCEPTION_SPECIAL_EMPTY_STRING
	}
	body, _ := ioutil.ReadAll(in.Body)
	in.Body.Close()

	r := mvutil.RequestParams{}
	//var rawQuery mvutil.RequestQueryMap
	rawQuery, err := RenderKSYUNReqParam(body, &r)
	if err != nil {
		mvutil.Logger.Runtime.Warnf("ksyun render params error. error=[%s]", err)
		return nil, errorcode.EXCEPTION_PARAMS_ERROR
	}
	RenderReqParam(in, &r, rawQuery)
	return &r, nil
}

func RenderKSYUNReqParam(body []byte, r *mvutil.RequestParams) (map[string][]string, error) {
	var BidRequest protobuf.STExchangeReq
	err := proto.Unmarshal(body, &BidRequest)
	if err != nil {
		return nil, err
	}
	// 获取Request内容
	device := BidRequest.Device
	network := BidRequest.Network
	app := BidRequest.App
	gps := BidRequest.Gps
	adslot := BidRequest.Adslot
	query := make(mvutil.RequestQueryMap)
	query["useragent"] = mvutil.RenderQueryMapV(device.Ua)
	query["mac"] = mvutil.RenderQueryMapV(device.Udid.Mac)
	if &network.Ipv4 == nil {
		query["cc"] = mvutil.RenderQueryMapV("CN")
	} else {
		query["client_ip"] = mvutil.RenderQueryMapV(network.Ipv4)
	}
	// need CC in case we don't get IPV4 address
	if device.OsType == 1 {
		query["platform"] = mvutil.RenderQueryMapV("1")
		query["android_id"] = mvutil.RenderQueryMapV(device.Udid.AndroidId)
		query["imei"] = mvutil.RenderQueryMapV(device.Udid.Imei)
	} else if device.OsType == 2 {
		query["platform"] = mvutil.RenderQueryMapV("2")
		query["idfa"] = mvutil.RenderQueryMapV(device.Udid.Idfa)
	}
	query["brand"] = mvutil.RenderQueryMapV(string(device.Vendor))
	query["model"] = mvutil.RenderQueryMapV(string(device.Model))
	osversion := device.OsVersion
	osv := strconv.Itoa(int(osversion.Major)) + "." + strconv.Itoa(int(osversion.Minor)) + "." + strconv.Itoa(int(osversion.Micro))
	query["os_version"] = mvutil.RenderQueryMapV(osv)
	// Don't know whether language is necessary, keep it unsolved...
	query["package_name"] = mvutil.RenderQueryMapV(app.AppPackage)
	if app.AppVersion != nil {
		appversion := strconv.Itoa(int(app.AppVersion.Major)) + "." + strconv.Itoa(int(app.AppVersion.Minor)) + "." + strconv.Itoa(int(app.AppVersion.Micro))
		query["app_version_name"] = mvutil.RenderQueryMapV(appversion)
	}
	query["network_type"] = mvutil.RenderQueryMapV(getKSYUNNetworkType(network.ConnectionType))
	if device.ScreenSize.Width != 0 && device.ScreenSize.Height != 0 {
		resolution := strconv.Itoa(int(device.ScreenSize.Width)) + "x" + strconv.Itoa(int(device.ScreenSize.Height))
		query["screen_size"] = mvutil.RenderQueryMapV(resolution)
	}
	if device.ScreenType == 1 {
		query["orientation"] = mvutil.RenderQueryMapV("2")
	} else if device.ScreenType == 2 {
		query["orientation"] = mvutil.RenderQueryMapV("24")
	}
	query["lat"] = mvutil.RenderQueryMapV(strconv.FormatFloat(gps.Latitude, 'E', -1, 64))
	query["lng"] = mvutil.RenderQueryMapV(strconv.FormatFloat(gps.Longitude, 'E', -1, 64))
	// Start dealing with feed back...
	query["app_id"] = mvutil.RenderQueryMapV(app.AppId)
	query["unit_id"] = mvutil.RenderQueryMapV(adslot.AdslotId)
	r.Param.KSYUNAdSlotID = adslot.AdslotId
	r.Param.KSYUNSKey = BidRequest.SearchKey
	r.Param.KSYUNRequestID = BidRequest.RequestId
	// 金山云仅处理一条广告，故返回一条.
	query["ad_num"] = mvutil.RenderQueryMapV("1")
	query["sign"] = mvutil.RenderQueryMapV(mvconst.NO_CHECK_SIGN)
	if BidRequest.RequestProtocolType == 2 {
		query["http_req"] = mvutil.RenderQueryMapV("2")
	}
	// 从mongo调取底价：
	onlinePriceFloorUnits, ifFind := extractor.GetONLINE_PRICE_FLOOR_APPID()
	if !ifFind {
		return nil, errors.New("ksyunreqparam filter can't get bid floor price")
	}
	price, ok := onlinePriceFloorUnits[app.AppId]
	if !ok {
		price = 2500
	}
	if adslot.MinimumCpm >= int32(price) {
		return nil, errors.New("MinimunCpm from KSYUN is too high")
	} else {
		r.Param.KSYUNMaxCPM = int64(price)
	}
	return query, nil
}

func getKSYUNNetworkType(networkType protobuf.ConnectionType) string {
	switch networkType {
	case 1:
		return "2"
	case 2:
		return "2"
	case 3:
		return "3"
	case 4:
		return "4"
	case 5:
		return "4"
	case 100:
		return "9"
	case 101:
		return "9"
	default:
		return "0"
	}
}

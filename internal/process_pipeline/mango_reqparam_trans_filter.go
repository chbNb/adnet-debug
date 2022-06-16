package process_pipeline

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	jsoniter "github.com/json-iterator/go"
	"gitlab.mobvista.com/ADN/adnet/internal/errorcode"
	"gitlab.mobvista.com/ADN/adnet/internal/extractor"
	"gitlab.mobvista.com/ADN/adnet/internal/mvconst"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
)

type MangoReqparamTransFilter struct {
}

type BidRequsetMango struct {
	Version int         `json:"version"` //回填
	Bid     string      `json:"bid"`     //回填
	Imp     []ImpMango  `json:"imp"`
	Device  DeviceMango `json:"device"`
}

type ImpMango struct {
	SpaceId  string `json:"space_id"`      //unitId映射
	MinPrice int    `json:"min_cpm_price"` // 流量最低竞标价格
}

type DeviceMango struct {
	IMEI  string `json:"imei"`
	IDFA  string `json:"idfa"`
	ANID  string `json:"anid"`
	MAC   string `json:"mac"`
	OS    string `json:"os"`
	Brand string `json:"brand"`
	Model string `json:"model"`
	SW    int    `json:"sw"`
	SH    int    `json:"sh"`
	IP    string `json:"ip"`
	//City_Code string `json:"city_code"` Citycode没有接收，注意
	UA             string `json:"ua"`
	ConnectionType int    `json:"connectiontype"`
	//Openudid string `json:"openudid"`
}

func (this *MangoReqparamTransFilter) Process(data interface{}) (interface{}, error) {
	in, ok := data.(*http.Request)
	if !ok {
		return nil, errorcode.EXCEPTION_SPECIAL_EMPTY_STRING
	}
	body, _ := ioutil.ReadAll(in.Body)
	in.Body.Close()
	r := mvutil.RequestParams{}
	rawQuery, err := RenderMangoReqParam(body, &r)
	if err != nil {
		return nil, errorcode.EXCEPTION_PARAMS_ERROR
	}
	RenderReqParam(in, &r, rawQuery)
	return &r, nil
}

func RenderMangoReqParam(body []byte, r *mvutil.RequestParams) (map[string][]string, error) {
	var BidRequsest BidRequsetMango
	err := jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(body, &BidRequsest)
	if err != nil {
		return nil, err
	}
	query := make(mvutil.RequestQueryMap, 0)
	query["useragent"] = mvutil.RenderQueryMapV(BidRequsest.Device.UA)
	query["client_ip"] = mvutil.RenderQueryMapV(BidRequsest.Device.IP)
	// 芒果在19年11月改为传imei md5值
	query["imei_md5"] = mvutil.RenderQueryMapV(BidRequsest.Device.IMEI)
	query["brand"] = mvutil.RenderQueryMapV(BidRequsest.Device.Brand)
	query["model"] = mvutil.RenderQueryMapV(BidRequsest.Device.Model)
	os, osv := getMangoOS(BidRequsest.Device.OS)
	if os == mvconst.PlatformNameAndroid {
		query["platform"] = mvutil.RenderQueryMapV("1")
		query["android_id"] = mvutil.RenderQueryMapV(BidRequsest.Device.ANID)
	} else if os == mvconst.PlatformNameIOS {
		query["platform"] = mvutil.RenderQueryMapV("2")
		query["idfa"] = mvutil.RenderQueryMapV(BidRequsest.Device.IDFA)
	}
	query["os_version"] = mvutil.RenderQueryMapV(osv)
	query["mac"] = mvutil.RenderQueryMapV(BidRequsest.Device.MAC)
	if BidRequsest.Device.SH != 0 && BidRequsest.Device.SW != 0 {
		resolution := strconv.Itoa(BidRequsest.Device.SW) + "x" + strconv.Itoa(BidRequsest.Device.SH)
		query["screen_size"] = mvutil.RenderQueryMapV(resolution)
	}
	query["network_type"] = mvutil.RenderQueryMapV(getMangoNetworkType(BidRequsest.Device.ConnectionType))
	//特别地，将芒果需要回传的值传递至param
	r.Param.MangoBid = BidRequsest.Bid
	r.Param.MangoVersion = BidRequsest.Version
	// 处理app_Id等内容
	config := extractor.GetMANGO_APPID_AND_UNITID()
	if len(config) == 0 {
		return nil, errorcode.EXCEPTION_PARAMS_ERROR
	}
	if len(BidRequsest.Imp) == 0 {
		return nil, errorcode.EXCEPTION_PARAMS_ERROR
	}
	// 最低竞价价格
	r.Param.MangoMinPrice = BidRequsest.Imp[0].MinPrice
	spaceId := BidRequsest.Imp[0].SpaceId
	if platformConf, ok := config[os]; ok {
		if spaceIdConf, ok := platformConf[spaceId]; ok {
			query["app_id"] = mvutil.RenderQueryMapV(strconv.FormatInt(spaceIdConf.AppId, 10))
			query["unit_id"] = mvutil.RenderQueryMapV(strconv.FormatInt(spaceIdConf.UnitId, 10))
		}
	}

	query["sign"] = mvutil.RenderQueryMapV("NO_CHECK_SIGN")
	query["ad_num"] = mvutil.RenderQueryMapV("1")
	query["orientation"] = mvutil.RenderQueryMapV("2")
	return query, nil
}

func getMangoNetworkType(networkType int) string {
	switch networkType {
	case 0:
		return "0"
	case 1:
		return "9"
	case 2:
		return "2"
	case 3:
		return "3"
	case 4:
		return "4"
	default:
		return "0"
	}
}

func getMangoOS(request string) (os string, osv string) {
	request = strings.ToLower(request)
	if len(request) == 0 {
		return "", ""
	}
	osArr := strings.Split(request, "_")
	if len(osArr) < 2 {
		return "", ""
	}
	os = osArr[0]
	osv = osArr[1]
	return os, osv
}

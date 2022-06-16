package process_pipeline

import (
	"io/ioutil"
	"net/http"
	"strconv"

	jsoniter "github.com/json-iterator/go"
	"gitlab.mobvista.com/ADN/adnet/internal/errorcode"
	"gitlab.mobvista.com/ADN/adnet/internal/extractor"
	"gitlab.mobvista.com/ADN/adnet/internal/mvconst"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
)

type CMReqparamTransFilter struct {
}

type CM struct {
	Id     string   `json:"id"`
	Imp    []ImpCM  `json:"imp"`
	App    AppCM    `json:"app"`
	Device DeviceCM `json:"device"`
}

type ImpCM struct {
	Id    string `json:"id"`    // 回填
	TagId string `json:"tagid"` // 必填
	//视频对象不需要
}

type AppCM struct {
	Id     string `json:"id"` // not used
	Bundle string `json:"bundle"`
}

type DeviceCM struct {
	UA string `json:"ua"`
	// GEO is not supported yet...
	IP             string `json:"ip"`
	Make           string `json:"make"`
	Model          string `json:"model"`
	OS             string `json:"os"`
	OSV            string `json:"osv"`
	W              int    `json:"w"`
	H              int    `json:"h"`
	ConnectionType int    `json:"connectiontype"`
	//DpidMD5 string `json:"dpidmd5"` // MD5 value of ADID or IDFA
	IFA      string `json:"ifa"` // IDFA or ADID
	IMEI     string `json:"imei"`
	Language string `json:"language"`
}

func (cmrtf *CMReqparamTransFilter) Process(data interface{}) (interface{}, error) {
	in, ok := data.(*http.Request)
	if !ok {
		return nil, errorcode.EXCEPTION_SPECIAL_EMPTY_STRING
	}
	body, _ := ioutil.ReadAll(in.Body)
	in.Body.Close()
	r := mvutil.RequestParams{}
	//var rawQuery mvutil.RequestQueryMap
	rawQuery, err := RenderCMReqParam(body, &r)
	if err != nil {
		return nil, errorcode.EXCEPTION_PARAMS_ERROR
	}
	RenderReqParam(in, &r, rawQuery)
	return &r, nil
}

func RenderCMReqParam(body []byte, r *mvutil.RequestParams) (map[string][]string, error) {
	var BidRequest CM
	err := jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(body, &BidRequest)
	if err != nil {
		return nil, err
	}
	query := make(mvutil.RequestQueryMap)
	query["useragent"] = mvutil.RenderQueryMapV(BidRequest.Device.UA)
	if BidRequest.Device.IP == "" {
		query["cc"] = mvutil.RenderQueryMapV("CN")
	}
	query["client_ip"] = mvutil.RenderQueryMapV(BidRequest.Device.IP)
	// CC is needed in case we didn't get IP
	if mvutil.GetPlatformStr(BidRequest.Device.OS) == mvconst.PlatformAndroid {
		query["platform"] = mvutil.RenderQueryMapV("1")
		query["android_id"] = mvutil.RenderQueryMapV(BidRequest.Device.IFA)
	} else if mvutil.GetPlatformStr(BidRequest.Device.OS) == mvconst.PlatformIOS {
		query["platform"] = mvutil.RenderQueryMapV("2")
		query["idfa"] = mvutil.RenderQueryMapV(BidRequest.Device.IFA)
	}
	query["imei"] = mvutil.RenderQueryMapV(BidRequest.Device.IMEI)
	query["brand"] = mvutil.RenderQueryMapV(BidRequest.Device.Make)
	query["model"] = mvutil.RenderQueryMapV(BidRequest.Device.Model)
	query["os_version"] = mvutil.RenderQueryMapV(BidRequest.Device.OSV)
	query["language"] = mvutil.RenderQueryMapV(BidRequest.Device.Language)
	query["package_name"] = mvutil.RenderQueryMapV(BidRequest.App.Bundle)
	query["network_type"] = mvutil.RenderQueryMapV(getCMNetworkType(BidRequest.Device.ConnectionType))
	//避免开发者出现漏传屏幕宽或高的情况
	if BidRequest.Device.W != 0 && BidRequest.Device.H != 0 {
		resolution := strconv.Itoa(BidRequest.Device.W) + "x" + strconv.Itoa(BidRequest.Device.H)
		query["screen_size"] = mvutil.RenderQueryMapV(resolution)
	}

	query["sign"] = mvutil.RenderQueryMapV("NO_CHECK_SIGN")

	cmMap, _ := extractor.GetCM_APPID_AND_UNITID()
	// 处理回传参数。
	r.Param.CMId = BidRequest.Id
	for _, v := range BidRequest.Imp {
		r.Param.CMImpId = v.Id
		cmAppUnit, ok := cmMap[v.TagId]
		if ok {
			query["app_id"] = mvutil.RenderQueryMapV(strconv.FormatInt(cmAppUnit.AppId, 10))
			query["unit_id"] = mvutil.RenderQueryMapV(strconv.FormatInt(cmAppUnit.UnitId, 10))
		}

	}
	// 猎豹仅处理一条广告，故返回一条。
	query["ad_num"] = mvutil.RenderQueryMapV("1")
	// isvast
	query["is_vast"] = mvutil.RenderQueryMapV("true")
	//fmt.Println(query) // Debug
	return query, nil
}

func getCMNetworkType(networkType int) string {
	switch networkType {
	case 2:
		return "9"
	default:
		return "0"
	}
}

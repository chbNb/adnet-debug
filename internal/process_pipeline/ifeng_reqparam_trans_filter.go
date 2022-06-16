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

type IFENGReqparamTransFilter struct {
}

type IFENG struct {
	Id     string      `json:"id"` // return
	Imp    []IFENGImp  `json:"imp"`
	App    IFENGApp    `json:"app"`
	Device IFENGDevice `json:"device"`
}

type IFENGImp struct {
	Id    string       `json:"id"`    // return
	TagId string       `json:"tagid"` // return
	Video []IFENGVideo `json:"video"`
	// Video struct is unnecessary to use, abort
}

type IFENGApp struct {
	Id     string `json:"id"`     // return?
	Name   string `json:"name"`   // app name
	Domain string `json:"domain"` // package name
	Ver    string `json:"ver"`    // app ver
}

type IFENGDevice struct {
	UA             string `json:"ua"`
	IP             string `json:"ip"`
	Model          string `json:"model"`
	OS             string `json:"os"`
	OSV            string `json:"osv"`
	H              int    `json:"h"`
	W              int    `json:"w"`
	ConnectionType int    `json:"connectiontype"`
	IFA            string `json:"ifa"`
	DIDMD5         string `json:"didmd5"`
}

type IFENGVideo struct {
	MinDuration int `json:"minduration"`
	MaxDuration int `json:"maxduration"`
}

func (ifrtf *IFENGReqparamTransFilter) Process(data interface{}) (interface{}, error) {
	in, ok := data.(*http.Request)
	if !ok {
		return nil, errorcode.EXCEPTION_SPECIAL_EMPTY_STRING
	}
	body, _ := ioutil.ReadAll(in.Body)
	in.Body.Close()
	r := mvutil.RequestParams{}
	//var rawQuery mvutil.RequestQueryMap
	rawQuery, err := RenderIFENGReqParam(body, &r)
	if err != nil {
		return nil, errorcode.EXCEPTION_PARAMS_ERROR
	}
	RenderReqParam(in, &r, rawQuery)
	return &r, nil
}

func RenderIFENGReqParam(body []byte, r *mvutil.RequestParams) (map[string][]string, error) {
	var BidRequest IFENG
	err := jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(body, &BidRequest)
	if err != nil {
		return nil, err
	}
	query := make(mvutil.RequestQueryMap)
	query["useragent"] = mvutil.RenderQueryMapV(BidRequest.Device.UA)
	query["client_ip"] = mvutil.RenderQueryMapV(BidRequest.Device.IP)
	if BidRequest.Device.IP == "" {
		query["cc"] = mvutil.RenderQueryMapV("CN")
	}
	if mvutil.GetPlatformStr(BidRequest.Device.OS) == mvconst.PlatformAndroid {
		query["platform"] = mvutil.RenderQueryMapV("1")
		query["imei_md5"] = mvutil.RenderQueryMapV(BidRequest.Device.DIDMD5)
	} else if mvutil.GetPlatformStr(BidRequest.Device.OS) == mvconst.PlatformIOS {
		query["platform"] = mvutil.RenderQueryMapV("2")
		query["idfa"] = mvutil.RenderQueryMapV(BidRequest.Device.IFA)
		query["http_req"] = mvutil.RenderQueryMapV("2")
	}
	query["model"] = mvutil.RenderQueryMapV(BidRequest.Device.Model)
	query["os_version"] = mvutil.RenderQueryMapV(BidRequest.Device.OSV)
	query["package_name"] = mvutil.RenderQueryMapV(BidRequest.App.Domain)
	query["app_version_name"] = mvutil.RenderQueryMapV(BidRequest.App.Ver)
	query["network_type"] = mvutil.RenderQueryMapV(getIFENGNetworkType(BidRequest.Device.ConnectionType))
	//避免开发者出现漏传屏幕宽或高的情况
	if BidRequest.Device.W != 0 && BidRequest.Device.H != 0 {
		resolution := strconv.Itoa(BidRequest.Device.W) + "x" + strconv.Itoa(BidRequest.Device.H)
		query["screen_size"] = mvutil.RenderQueryMapV(resolution)
	}
	//处理回传值
	r.Param.IFENGId = BidRequest.Id
	for _, v := range BidRequest.Imp {
		r.Param.IFENGImpId = v.Id
		r.Param.IFENGTagId = v.TagId
		for _, val := range v.Video {
			r.Param.MaxDuration = val.MaxDuration
			r.Param.MinDuration = val.MinDuration
		}
	}
	// ad num 为1.
	query["ad_num"] = mvutil.RenderQueryMapV("1")
	// using vast
	query["is_vast"] = mvutil.RenderQueryMapV("true")
	// unit_id, sign and app_id missing, ifeng only use ios, still talking, so using migu...
	platformStr := strings.ToLower(BidRequest.Device.OS)
	auConf, ifFind := extractor.GetIFENG_APPID_AND_UNITID()
	if conf, ok := auConf[platformStr]; ok && ifFind {
		query["unit_id"] = mvutil.RenderQueryMapV(strconv.FormatInt(conf.UnitId, 10))
		query["app_id"] = mvutil.RenderQueryMapV(strconv.FormatInt(conf.AppId, 10))
	}
	query["sign"] = mvutil.RenderQueryMapV("NO_CHECK_SIGN")

	return query, nil
}

func getIFENGNetworkType(networkType int) string {
	switch networkType {
	case 1:
		return "9"
	case 2:
		return "9"
	case 3:
		return "2"
	case 4:
		return "2"
	case 5:
		return "3"
	case 6:
		return "4"
	default:
		return "0"
	}
}

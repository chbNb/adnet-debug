package process_pipeline

import (
	"errors"
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

type pptvParam struct {
	ID     string     `json:"id"`
	Device pptvDevice `json:"device"`
	Imp    []pptvImp  `json:"imp"`
	App    pptvApp    `json:"app"`
}

type pptvDevice struct {
	Ua          string `json:"ua"`
	IP          string `json:"ip"`
	Brand       string `json:"make"`
	Model       string `json:"model"`
	Platform    string `json:"os"`
	OsVersion   string `json:"osv"`
	NetworkType int    `json:"connectiontype"`
	DevId       string `json:"ifa"`
	ImeiSha     string `json:"didsha1"`
	ScreenW     int    `json:"w"`
	ScreenH     int    `json:"h"`
	Orientation int    `json:"orientation"`
}

type pptvImp struct {
	HttpReq         int       `json:"https_flag"`
	Video           pptvVideo `json:"video"`
	Interactivetype []int     `json:"interactivetype"`
	TagId           string    `json:"tagid"`
	BidFloor        float64   `json:"bidfloor"`
}

type pptvVideo struct {
	AllowType   []int `json:"allyesadformat"`
	UnitW       int   `json:"w"`
	UnitH       int   `json:"h"`
	MaxDuration int   `json:"maxduration"`
	MinDuration int   `json:"minduration"`
}

type pptvApp struct {
	PackageName    string `json:"bundle"`
	AppVersionCode string `json:"ver"`
}

type PPTVReqParamTransFilter struct {
}

func (pprptf *PPTVReqParamTransFilter) Process(data interface{}) (interface{}, error) {
	in, ok := data.(*http.Request)
	if !ok {
		return nil, errorcode.EXCEPTION_SPECIAL_EMPTY_STRING
	}
	body, _ := ioutil.ReadAll(in.Body)
	in.Body.Close()
	r := mvutil.RequestParams{}
	//var rawQuery mvutil.RequestQueryMap
	rawQuery, err := RenderPPTVReqParam(body, &r)
	if err != nil {
		mvutil.Logger.Runtime.Warnf("pptv render params error. error=[%s]", err)
		return nil, errorcode.EXCEPTION_PARAMS_ERROR
	}
	RenderReqParam(in, &r, rawQuery)
	return &r, nil
}

func RenderPPTVReqParam(body []byte, r *mvutil.RequestParams) (map[string][]string, error) {
	var BidRequest pptvParam
	err := jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(body, &BidRequest)
	if err != nil {
		return nil, err
	}
	query := make(mvutil.RequestQueryMap)
	query["useragent"] = mvutil.RenderQueryMapV(BidRequest.Device.Ua)
	query["client_ip"] = mvutil.RenderQueryMapV(BidRequest.Device.IP)
	query["brand"] = mvutil.RenderQueryMapV(BidRequest.Device.Brand)
	query["model"] = mvutil.RenderQueryMapV(BidRequest.Device.Model)
	query["os_version"] = mvutil.RenderQueryMapV(BidRequest.Device.OsVersion)
	query["imei_md5"] = mvutil.RenderQueryMapV(BidRequest.Device.ImeiSha)
	query["package_name"] = mvutil.RenderQueryMapV(BidRequest.App.PackageName)
	query["app_version_code"] = mvutil.RenderQueryMapV(BidRequest.App.AppVersionCode)
	query["ad_num"] = mvutil.RenderQueryMapV("1")
	query["sign"] = mvutil.RenderQueryMapV("NO_CHECK_SIGN")

	platformStr := strings.ToLower(BidRequest.Device.Platform)
	if platformStr == mvconst.PlatformNameAndroid {
		query["platform"] = mvutil.RenderQueryMapV("1")
		query["android_id"] = mvutil.RenderQueryMapV(BidRequest.Device.DevId)
	} else if platformStr == mvconst.PlatformNameIOS {
		query["platform"] = mvutil.RenderQueryMapV("2")
		query["idfa"] = mvutil.RenderQueryMapV(BidRequest.Device.DevId)
	}

	query["network_type"] = mvutil.RenderQueryMapV(getPPTVNetworkType(BidRequest.Device.NetworkType))
	query["screen_size"] = mvutil.RenderQueryMapV(strconv.Itoa(BidRequest.Device.ScreenW) + "x" + strconv.Itoa(BidRequest.Device.ScreenH))
	if BidRequest.Device.Orientation == 1 {
		query["orientation"] = mvutil.RenderQueryMapV("2")
	} else {
		query["orientation"] = mvutil.RenderQueryMapV("1")
	}
	r.Param.OnlineReqId = BidRequest.ID
	for _, v := range BidRequest.Imp {
		if v.HttpReq == 1 {
			query["http_req"] = mvutil.RenderQueryMapV("2")
		}
		if !mvutil.InArray(6, v.Video.AllowType) {
			return nil, errors.New("no allow format")
		}
		query["unit_size"] = mvutil.RenderQueryMapV(strconv.Itoa(v.Video.UnitW) + "x" + strconv.Itoa(v.Video.UnitH))
		// 素材审核只推15s视频素材
		//r.Param.MaxDuration = v.Video.MaxDuration
		//r.Param.MinDuration = v.Video.MinDuration
		// 与link_type对比
		r.Param.AllowType = v.Interactivetype

		// 价格比较，若比设置值低则过滤
		onlinePriceFloor, ifFind := extractor.GetONLINE_PRICE_FLOOR()
		var resPrice float64
		if ifFind {
			price, ok := onlinePriceFloor["pptv"]
			if ok {
				resPrice = price
			} else {
				resPrice = 10
			}
		}
		if v.BidFloor > resPrice {
			return nil, errors.New("minPrice is too high")
		}
		// 出价价格
		r.Param.MaxPrice = int32(resPrice)
		pptvMap, _ := extractor.GetPPTV_APPID_AND_UNITID()
		pptvMvMap, ok := pptvMap[v.TagId]
		if ok {
			query["app_id"] = mvutil.RenderQueryMapV(strconv.FormatInt(pptvMvMap.AppId, 10))
			query["unit_id"] = mvutil.RenderQueryMapV(strconv.FormatInt(pptvMvMap.UnitId, 10))
		}
		/////////todo
		break
	}
	return query, nil

}

func getPPTVNetworkType(networkType int) string {
	switch networkType {
	case 2:
		return "9"
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

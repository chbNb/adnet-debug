package process_pipeline

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"gitlab.mobvista.com/ADN/adnet/internal/errorcode"
	"gitlab.mobvista.com/ADN/adnet/internal/mvconst"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
)

type MIGUReqparamTransFilter struct {
}

type BidRequest struct {
	Id     string     `json:"id"`
	Imp    []ImpMIGU  `json:"imp"`
	App    AppMIGU    `json:"app"`
	Device DeviceMIGU `json:"device"`
}

type ImpMIGU struct {
	Id     string     `json:"id"`    //Ctrl + C, Ctrl + V给输出，存储至param
	TagId  string     `json:"tagid"` //unitId
	Native NativeMIGU `json:"native"`
	Cur    string     `json:"cur"` //存储至param
}

type NativeMIGU struct {
	Title TitleMIGU `json:"title"`
	Desc  DescMIGU  `json:"desc"`
}

type TitleMIGU struct {
	Len int `json:"len"` //请存储至param
}

type DescMIGU struct {
	Len int `json:"len"` //请存储至param
}

type AppMIGU struct {
	Id     string `json:"id"` //app_id以及sign，注意这个参数同时包含了app_Id以及sign！
	Bundle string `json:"bundle"`
}

type DeviceMIGU struct {
	W              int    `json:"w"`
	H              int    `json:"h"`
	UA             string `json:"ua"`
	IP             string `json:"ip"`
	Did            string `json:"did"`
	Dpid           string `json:"dpid"`
	MAC            string `json:"mac"`
	IFA            string `json:"ifa"`
	Make           string `json:"make"`
	Model          string `json:"model"`
	OS             string `json:"os"`
	OSV            string `json:"osv"`
	Language       string `json:"language"`
	Connectiontype int    `json:"connectiontype"`
}

type Idlist struct {
	AppId string
	Sign  string
}

func (mgrtf *MIGUReqparamTransFilter) Process(data interface{}) (interface{}, error) {
	in, ok := data.(*http.Request)
	if !ok {
		return nil, errorcode.EXCEPTION_SPECIAL_EMPTY_STRING
	}
	body, _ := ioutil.ReadAll(in.Body)
	in.Body.Close()
	r := mvutil.RequestParams{}
	//var rawQuery mvutil.RequestQueryMap
	rawQuery, err := RenderMIGUReqParam(body, &r)
	if err != nil {
		return nil, errorcode.EXCEPTION_PARAMS_ERROR
	}
	RenderReqParam(in, &r, rawQuery)
	return &r, nil
}

func RenderMIGUReqParam(body []byte, r *mvutil.RequestParams) (map[string][]string, error) {
	var BidRequest BidRequest
	var MIGUIdList Idlist
	err := json.Unmarshal(body, &BidRequest)
	if err != nil {
		return nil, err
	}
	//fmt.Println(BidRequest)
	query := make(mvutil.RequestQueryMap)
	query["useragent"] = mvutil.RenderQueryMapV(BidRequest.Device.UA)
	query["client_ip"] = mvutil.RenderQueryMapV(BidRequest.Device.IP)
	if mvutil.GetPlatformStr(BidRequest.Device.OS) == mvconst.PlatformAndroid {
		query["platform"] = mvutil.RenderQueryMapV("1")
		query["android_id"] = mvutil.RenderQueryMapV(BidRequest.Device.Dpid)
		query["gaid"] = mvutil.RenderQueryMapV(BidRequest.Device.IFA)
	} else if mvutil.GetPlatformStr(BidRequest.Device.OS) == mvconst.PlatformIOS {
		query["platform"] = mvutil.RenderQueryMapV("2")
		query["idfa"] = mvutil.RenderQueryMapV(BidRequest.Device.IFA)
	}
	query["imei"] = mvutil.RenderQueryMapV(BidRequest.Device.Did)
	query["brand"] = mvutil.RenderQueryMapV(BidRequest.Device.Make)
	query["model"] = mvutil.RenderQueryMapV(BidRequest.Device.Model)
	query["os_version"] = mvutil.RenderQueryMapV(BidRequest.Device.OSV)
	query["language"] = mvutil.RenderQueryMapV(BidRequest.Device.Language)
	query["mac"] = mvutil.RenderQueryMapV(BidRequest.Device.MAC)
	query["package_name"] = mvutil.RenderQueryMapV(BidRequest.App.Bundle)
	query["network_type"] = mvutil.RenderQueryMapV(getMIGUNetworkType(BidRequest.Device.Connectiontype))
	//避免开发者出现漏传屏幕宽或高的情况
	if BidRequest.Device.W != 0 && BidRequest.Device.H != 0 {
		resolution := strconv.Itoa(BidRequest.Device.W) + "x" + strconv.Itoa(BidRequest.Device.H)
		query["screen_size"] = mvutil.RenderQueryMapV(resolution)
	}
	//特别地，将咪咕需要回传的值传递给param
	r.Param.MIGUId = BidRequest.Id
	//bug fixed,避免无法检索的思路：遍历。。。
	for _, v := range BidRequest.Imp {
		r.Param.MIGUImpId = v.Id
		r.Param.MIGUCur = v.Cur
		r.Param.MIGUDescLen = v.Native.Title.Len
		r.Param.MIGUTitleLen = v.Native.Desc.Len
		query["unit_id"] = mvutil.RenderQueryMapV(v.TagId)
	}
	query["ad_num"] = mvutil.RenderQueryMapV("1") //requeset_param会自动转换，不用担心
	//获取广告unit_Id,sign以及app_Id。
	MIGUIdList = getMIGUAppIdAndSIGN(BidRequest.App.Id)
	query["app_id"] = mvutil.RenderQueryMapV(MIGUIdList.AppId)
	query["sign"] = mvutil.RenderQueryMapV(MIGUIdList.Sign)
	return query, nil
}

func getMIGUNetworkType(networkType int) string {
	switch networkType {
	case 1:
		return "9"
	case 2:
		return "9"
	case 3:
		return "2"
	case 4:
		return "3"
	case 5:
		return "4"
	default:
		return "0"
	}
}

//Need Performance Optimization
func getMIGUAppIdAndSIGN(param string) Idlist {
	pos := strings.Index(param, "ngis")
	var result Idlist
	if pos >= 0 {
		trans := []byte(param)
		result.AppId = string(trans[0:pos])
		result.Sign = string(trans[pos+4:])
	} else {
		result.AppId = "0"
		result.Sign = "0"
	}
	return result
}

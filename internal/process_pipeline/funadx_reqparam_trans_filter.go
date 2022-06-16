package process_pipeline

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	jsoniter "github.com/json-iterator/go"
	"gitlab.mobvista.com/ADN/adnet/internal/errorcode"
	"gitlab.mobvista.com/ADN/adnet/internal/extractor"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
)

type FUNADXReqparamTransFilter struct {
}

type funParam struct {
	ID        string    `json:"id"`
	Version   string    `json:"version"`
	FunImp    []FunImp  `json:"imp"`
	FunDevice FunDevice `json:"device"`
	FunApp    FunApp    `json:"app"`
}

type FunImp struct {
	// Deals Deals
	ID    string `json:"id"`
	Tagid string `json:"tagid"`
}

// type Deals struct {
// 	Id string `json:"id"`
// }

type FunDevice struct {
	Ua  string `json:"ua"`
	IP  string `json:"ip"`
	Did string `json:"did"`
	// DidMd5   string `json:"didmd5"`
	Dpid string `json:"dpid"`
	// DpidMd5  string `json:"dpidmd5"`
	Make     string `json:"make"`
	Model    string `json:"model"`
	Os       string `json:"os"`
	Osv      string `json:"osv"`
	Carrier  string `json:"carrier"`
	Language string `json:"language"`
	Ext      Ext    `json:"ext"`
}

type Ext struct {
	Idfa string `json:"idfa"`
	Mac  string `json:"mac"`
}

type FunApp struct {
	Bundle string `json:"bundle"`
}

func (fartf *FUNADXReqparamTransFilter) Process(data interface{}) (interface{}, error) {
	in, ok := data.(*http.Request)
	if !ok {
		return nil, errorcode.EXCEPTION_SPECIAL_EMPTY_STRING
	}
	body, _ := ioutil.ReadAll(in.Body)
	in.Body.Close()
	r := mvutil.RequestParams{}
	//var rawQuery mvutil.RequestQueryMap
	rawQuery, err := RenderFunReqParam(body, &r)
	if err != nil {
		return nil, errorcode.EXCEPTION_PARAMS_ERROR
	}
	RenderReqParam(in, &r, rawQuery)
	return &r, nil
}

func RenderFunReqParam(body []byte, r *mvutil.RequestParams) (map[string][]string, error) {
	var funParam funParam
	var FunImpIDList []string
	var TagIDList []string
	var TagID string
	err := jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(body, &funParam)
	if err != nil {
		return nil, err
	}
	query := make(mvutil.RequestQueryMap)
	query["useragent"] = mvutil.RenderQueryMapV(funParam.FunDevice.Ua)
	query["client_ip"] = mvutil.RenderQueryMapV(funParam.FunDevice.IP)
	query["imei"] = mvutil.RenderQueryMapV(funParam.FunDevice.Did)
	platformStr := strings.ToLower(funParam.FunDevice.Os)
	if platformStr == "android" {
		query["platform"] = mvutil.RenderQueryMapV("1")
	} else if platformStr == "ios" {
		query["platform"] = mvutil.RenderQueryMapV("2")
		query["idfa"] = mvutil.RenderQueryMapV(funParam.FunDevice.Ext.Idfa)
	}
	query["brand"] = mvutil.RenderQueryMapV(funParam.FunDevice.Make)
	query["model"] = mvutil.RenderQueryMapV(funParam.FunDevice.Model)
	query["os_version"] = mvutil.RenderQueryMapV(funParam.FunDevice.Osv)
	query["mcc"] = mvutil.RenderQueryMapV(funParam.FunDevice.Carrier)
	query["language"] = mvutil.RenderQueryMapV(funParam.FunDevice.Language)
	query["mac"] = mvutil.RenderQueryMapV(funParam.FunDevice.Ext.Mac)
	query["package_name"] = mvutil.RenderQueryMapV(funParam.FunApp.Bundle)
	// 请求id
	r.Param.FunRequestId = funParam.ID
	// 每次请求一条广告
	query["ad_num"] = mvutil.RenderQueryMapV("1")
	query["sign"] = mvutil.RenderQueryMapV("NO_CHECK_SIGN")

	if len(funParam.FunImp) <= 0 {
		return nil, errors.New("param error")
	}
	for _, v := range funParam.FunImp {
		FunImpIDList = append(FunImpIDList, v.ID)
		TagIDList = append(TagIDList, v.Tagid)
	}
	if len(FunImpIDList[0]) > 0 {
		r.Param.FunImpId = FunImpIDList[0]
	}

	if len(TagIDList[0]) > 0 {
		TagID = TagIDList[0]
	}
	// 广告位mapping
	funMap, _ := extractor.GetFUN_MV_MAP()
	funMvMap, ok := funMap[TagID]
	if ok && len(funMvMap["appId"]) > 0 && len(funMvMap["unitId"]) > 0 {
		query["app_id"] = mvutil.RenderQueryMapV(funMvMap["appId"])
		query["unit_id"] = mvutil.RenderQueryMapV(funMvMap["unitId"])
	}
	return query, nil
}

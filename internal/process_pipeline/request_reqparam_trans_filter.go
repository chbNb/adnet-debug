package process_pipeline

import (
	"encoding/base64"
	"net/http"
	"strconv"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
	"gitlab.mobvista.com/ADN/adnet/internal/errorcode"
	"gitlab.mobvista.com/ADN/adnet/internal/mvconst"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
)

type RequestReqparamTransFilter struct {
}

type requestP struct {
	Ip          string `json:"ip"`
	Ua          string `json:"ua"`
	DeviceModel string `json:"dm"`
	OsVersion   string `json:"ov"`
	Platform    string `json:"pf"`
	Orientation int    `json:"ot"`
	Scenario    string `json:"snr"`
	AdType      string `json:"at"`
	Template    string `json:"tpl"`
	Adnum       int    `json:"an"`
	ImageSize   int    `json:"is"`
	ExcludeIds  string `json:"ei"`
	Mnc         string `json:"mnc"`
	Mcc         string `json:"mcc"`
	// timestamp 没用到 jsonp（暂时不用）
	TimeStamp string `json:"t"`
	AppId     string `json:"aid"`
	UnitId    int    `json:"utid"`
	Sign      string `json:"sn"`
	UnitSize  string `json:"us"`
	//RequestType int `json:"rt"`
	PingMode            int    `json:"cp"`
	Category            int    `json:"cat"`
	PackageName         string `json:"pn"`
	SdkVersion          string `json:"sv"`
	AppVersionName      string `json:"vn"`
	AppVersionCode      int    `json:"vc"`
	GooglePlayVersion   string `json:"gpv"`
	Imei                string `json:"im"`
	Mac                 string `json:"mac"`
	DevId               string `json:"did"`
	GoogleAdvertisingId string `json:"adid"`
	ScreenSize          string `json:"ss"`
	NetworkType         int    `json:"nt"` // 发现其实为int
	//NetworkTypeId string `json:"ntid"`
	Language         string `json:"l"`
	ImpressionImage  int    `json:"impression_image"`
	OnlyImpression   int    `json:"only_impression"`
	Network          int    `json:"network"`
	Offset           int    `json:"offset"`
	Timezone         string `json:"tz"`
	OfferPackageName string `json:"pkg"`
}

func (rrtf *RequestReqparamTransFilter) Process(data interface{}) (interface{}, error) {
	in, ok := data.(*http.Request)
	if !ok {
		return nil, errorcode.EXCEPTION_SPECIAL_EMPTY_STRING
		//return mvconst.GetResEmpty(), errors.New(mvconst.EXCEPTION_SPECIAL_RESULT_RETURN_NOENCODE)
	}
	rawQuery := mvutil.RequestQueryMap(in.URL.Query())
	// 判断是否为post请求
	var pStr string
	var formatStr string
	if len(rawQuery) <= 0 {
		pStr = in.PostFormValue("p")
		formatStr = in.PostFormValue("format")
	}
	if len(rawQuery) <= 0 && len(pStr) <= 0 {
		return nil, errorcode.EXCEPTION_PARAMS_ERROR
		//return mvconst.GetResEmpty(), errors.New(mvconst.EXCEPTION_SPECIAL_RESULT_RETURN_NOENCODE)
	}
	// get || post获取p参数
	var p string
	var format string
	if len(rawQuery) > 0 {
		p, _ = rawQuery.GetString("p", true, "")
		format, _ = rawQuery.GetString("p", true, "format")
	} else {
		p = pStr
		format = formatStr
	}
	// 默认为base64
	if format != "json" {
		format = "base64"
	}

	// 去除空格
	p = strings.Replace(p, " ", "+", -1)
	parseP, err := base64.StdEncoding.DecodeString(p)
	if err != nil {
		return nil, errorcode.EXCEPTION_PARAMS_ERROR
	}
	var requestP requestP
	err = jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal([]byte(parseP), &requestP)
	if err != nil {
		return nil, errorcode.EXCEPTION_PARAMS_ERROR
	}

	r := mvutil.RequestParams{}
	if requestP.Ip != "" {
		rawQuery["client_ip"] = mvutil.RenderQueryMapV(requestP.Ip)
	}
	if requestP.Ua != "" {
		rawQuery["useragent"] = mvutil.RenderQueryMapV(requestP.Ua)
	}
	RenderReqParam(in, &r, rawQuery)
	// 处理request接口参数
	RenderRequestParam(&r, requestP)

	return &r, nil
}

func RenderRequestParam(r *mvutil.RequestParams, p requestP) {
	r.Param.Model = p.DeviceModel
	r.Param.OSVersion = p.OsVersion
	r.Param.Platform = mvutil.GetPlatformStr(p.Platform)
	r.Param.Orientation = p.Orientation
	r.Param.Scenario = p.Scenario
	r.Param.AdTypeStr = p.AdType
	if len(r.Param.AdTypeStr) <= 0 && len(r.Param.Scenario) > 0 {
		adTypeArr := mvconst.GetAdtypeFromSnr(r.Param.Scenario)
		r.Param.AdTypeStr = mvutil.RandStringArr(adTypeArr)
	}
	if len(r.Param.Scenario) <= 0 {
		r.Param.Scenario = r.Param.AdTypeStr
	}
	// 获取模板
	if len(p.Template) <= 0 {
		fsTplArr := mvconst.GetFsTpl(r.Param.Orientation)
		p.Template = mvutil.RandStringArr(fsTplArr)
	}
	r.Param.AdNum = int32(p.Adnum)
	r.Param.ImageSizeID = p.ImageSize
	if r.Param.ImageSizeID == 0 {
		r.Param.ImageSizeID = mvutil.RenderImageSizeIdByAdType(r.Param.AdTypeStr, p.Template, r.Param.Orientation)
	}
	r.Param.ExcludeIDS = p.ExcludeIds
	r.Param.MNC = p.Mnc
	r.Param.MCC = p.Mcc
	r.Param.AppID, _ = strconv.ParseInt(p.AppId, 10, 64)
	r.Param.UnitID = int64(p.UnitId)
	r.Param.Sign = p.Sign
	r.Param.UnitSize = p.UnitSize
	r.Param.PingMode = 1
	r.Param.Category = p.Category
	r.Param.PackageName = p.PackageName
	r.Param.SDKVersion = p.SdkVersion
	r.Param.AppVersionName = p.AppVersionName
	r.Param.AppVersionCode = strconv.Itoa(p.AppVersionCode)
	r.Param.GPVersion = p.GooglePlayVersion
	r.Param.IMEI = p.Imei
	r.Param.MAC = p.Mac
	r.Param.AndroidID = p.DevId
	r.Param.GAID = p.GoogleAdvertisingId
	r.Param.ScreenSize = p.ScreenSize
	r.Param.NetworkType = p.NetworkType
	r.Param.Language = p.Language
	r.Param.ImpressionImage = p.ImpressionImage
	r.Param.OnlyImpression = p.OnlyImpression
	r.Param.NetWork = strconv.Itoa(p.Network)
	r.Param.Offset = int32(p.Offset)
	r.Param.TimeZone = p.Timezone
	r.Param.RequestType = mvconst.REQUEST_TYPE_SDK
	r.Param.TimeStamp = p.TimeStamp
	r.Param.OnlyImpression = 1
	r.Param.RequestTime = time.Now().Unix()
}

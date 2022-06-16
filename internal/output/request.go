package output

import (
	"crypto/md5"
	"encoding/base64"
	"net/url"

	"fmt"
	"strconv"

	jsoniter "github.com/json-iterator/go"
	"gitlab.mobvista.com/ADN/adnet/internal/mvconst"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
)

type RequestResult struct {
	Status int        `json:"status"`
	Msg    string     `json:"msg"`
	Data   RequestRep `json:"data"`
}

type RequestRep struct {
	Scenario          string       `json:"scenario"`
	AdType            string       `json:"adType"`
	Orientation       int          `json:"orientation"`
	OnlyImpressionUrl string       `json:"onlyImpressionUrl"`
	Rows              []RequestCam `json:"rows"`
	HtmlUrl           string       `json:"htmlUrl,omitempty"`
}
type RequestCam struct {
	CamId         int64   `json:"campaignId"`
	CrId          int64   `json:"creativeId"`
	AppName       string  `json:"appName"`
	AppDesc       string  `json:"appDesc"`
	AppScore      float32 `json:"appScore"`
	PackageName   string  `json:"packageName"`
	IconUrl       string  `json:"iconUrl"`
	ImageUrl      string  `json:"imageUrl"`
	ImageSize     string  `json:"imageSize"`
	ImpressionUrl string  `json:"impressionUrl"`
	ClickUrl      string  `json:"clickUrl"`
	NoticeUrl     string  `json:"noticeUrl"`
	Fca           int     `json:"fca"`
	Fcb           int     `json:"fcb"`
	ApkUrl        string  `json:"dataUrl"`
	AppSize       string  `json:"appSize"`
}

type ReqPParam struct {
	RequestId   string  `json:"k"`
	AppId       int64   `json:"aid"`
	Scenario    string  `json:"snr"`
	AdTypeStr   string  `json:"at"`
	OsVersion   string  `json:"ov"`
	Imei        string  `json:"im"`
	Mac         string  `json:"mac"`
	DevId       string  `json:"did"`
	Platform    string  `json:"pf"`
	Orientation int     `json:"ot"`
	DeviceModel string  `json:"dm"`
	SdkVersion  string  `json:"sv"`
	RequestTime int64   `json:"t"`
	Extra       string  `json:"ext"`
	CamIds      []int64 `json:"cids"`
	SParam      string  `json:"s"`
}

func RenderRequestRes(mr MobvistaResult, r *mvutil.RequestParams) RequestResult {
	var requestResult RequestResult
	requestResult.Status = 1
	requestResult.Msg = "success"
	var requestRep RequestRep
	requestRep.Scenario = r.Param.Scenario
	requestRep.AdType = r.Param.AdTypeStr
	requestRep.Orientation = r.Param.Orientation
	requestRep.OnlyImpressionUrl = mr.Data.OnlyImpressionURL
	var camIds []int64
	var requestCam RequestCam
	for _, v := range mr.Data.Ads {
		camIds = append(camIds, v.CampaignID)
		requestCam.CamId = v.CampaignID
		if r.Param.VideoCreativeid > 0 {
			requestCam.CrId = r.Param.VideoCreativeid
		} else {
			requestCam.CrId = r.Param.ImageCreativeid
		}
		requestCam.AppName = v.AppName
		requestCam.AppDesc = v.AppDesc
		requestCam.AppScore = v.Rating
		requestCam.PackageName = v.PackageName
		requestCam.IconUrl = v.IconURL
		requestCam.ImageUrl = v.ImageURL
		requestCam.ImageSize = v.ImageSize
		requestCam.ImpressionUrl = v.ImpressionURL
		requestCam.ClickUrl = v.ClickURL
		requestCam.NoticeUrl = v.NoticeURL
		requestCam.Fca = v.FCA
		requestCam.Fcb = v.FCB
		requestCam.ApkUrl = v.ApkUrl
		requestCam.AppSize = v.AppSize
		requestRep.Rows = append(requestRep.Rows, requestCam)
	}
	if r.Param.AdTypeStr == "full_screen" || r.Param.AdTypeStr == "overlay" {
		requestRep.HtmlUrl = renderHtmlUrl(r, camIds)
	}
	requestResult.Data = requestRep
	return requestResult
}

func renderHtmlUrl(r *mvutil.RequestParams, camIds []int64) string {
	var reqPParam ReqPParam
	reqPParam.RequestId = r.Param.RequestID
	reqPParam.AppId = r.Param.AppID
	reqPParam.Scenario = r.Param.Scenario
	reqPParam.AdTypeStr = r.Param.AdTypeStr
	reqPParam.OsVersion = r.Param.OSVersion
	reqPParam.Imei = r.Param.IMEI
	reqPParam.Mac = r.Param.MAC
	reqPParam.DevId = r.Param.AndroidID + r.Param.IDFA
	reqPParam.Platform = mvconst.GetPlatformStr(r.Param.Platform)
	reqPParam.Orientation = r.Param.Orientation
	reqPParam.DeviceModel = r.Param.Model
	reqPParam.SdkVersion = r.Param.SDKVersion
	reqPParam.RequestTime = r.Param.RequestTime
	reqPParam.Extra = r.Param.Extra
	reqPParam.CamIds = camIds
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	jsonEncodeCam, _ := json.Marshal(camIds)

	sParam := strconv.FormatInt(reqPParam.AppId, 10) + "_" + reqPParam.AdTypeStr + "_" +
		strconv.FormatInt(reqPParam.RequestTime, 10) + "_" + mvconst.SHOW_AD_KEY + "_" + string(jsonEncodeCam)
	reqPParam.SParam = fmt.Sprintf("%x", md5.Sum([]byte(sParam)))
	ParamJsonStr, err := json.Marshal(reqPParam)
	if err != nil {
		return ""
	}
	paramStrEncode := base64.StdEncoding.EncodeToString([]byte(ParamJsonStr))
	paramP := url.QueryEscape(paramStrEncode)

	tpl := ""
	if r.Param.AdTypeStr == "full_screen" {
		tpl = mvutil.RandStringArr(mvconst.GetFsTpl(r.Param.Orientation))
	}
	if r.Param.AdTypeStr == "overlay" && (r.Param.Scenario == "quit_overlay" || r.Param.Scenario == "exit") {
		tpl = "quit"
	}
	htmlUrl := "http://" + r.Param.Domain + "/show?p=" + paramP + "&tpl=" + tpl
	if r.Param.SDKVersion == "2.2.0" {
		htmlUrl = htmlUrl + "&direct_load=true"
	}
	return htmlUrl
}

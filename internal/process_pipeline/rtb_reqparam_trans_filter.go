package process_pipeline

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"gitlab.mobvista.com/ADN/adnet/internal/errorcode"
	"gitlab.mobvista.com/ADN/adnet/internal/mvconst"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
)

type RTBReqparamTransFilter struct {
}

type ucParams struct {
	Id     string   `json:"id"`
	Cur    []string `json:"cur"`
	App    App      `json:"app"`
	Device Device   `json:"device"`
	Imp    []Imp    `json:"imp"`
	Tmax   int      `json:"tmax"`
}

type App struct {
	Bundle   string `json:"bundle"`
	Id       string `json:"id"`
	Ver      string `json:"ver"`
	StoreUrl string `json:"storeurl"`
}

type Device struct {
	Ua             string `json:"ua"`
	Ip             string `json:"ip"`
	DeviceType     int    `json:"devicetype"`
	Make           string `json:"make"`
	Model          string `json:"model"`
	Os             string `json:"os"`
	Osv            string `json:"osv"`
	Language       string `json:"language"`
	ConnectionType int    `json:"connectiontype"`
	Ifa            string `json:"ifa"`
}

type Imp struct {
	Id          string  `json:"id"`
	TagId       string  `json:"tagid"`
	Bidfloor    float64 `json:"bidfloor"`
	BidfloorCur string  `json:"bidfloorcur"`
	Native      Native  `json:"native"`
}

type Native struct {
	Request string `json:"request"`
}

type RequestData struct {
	Assets     []Assets   `json:"assets"`
	NativeData NativeData `json:"native"`
}

type NativeData struct {
	Assets []Assets `json:"assets"`
}

type Assets struct {
	Id    int   `json:"id"`
	Data  Data  `json:"data"`
	Title Title `json:"title"`
	Img   Img   `json:"img"`
}

type Img struct {
	Hmin int `json:"hmin"`
	Wmin int `json:"wmin"`
	Type int `json:"type"`
}

type Title struct {
	Len int `json:"len"`
}

type Data struct {
	Type int `json:"type"`
	Len  int `json:"len"`
}

func (rrtf *RTBReqparamTransFilter) Process(data interface{}) (interface{}, error) {
	in, ok := data.(*http.Request)
	if !ok {
		return nil, errorcode.EXCEPTION_SPECIAL_EMPTY_STRING
		//return mvconst.GetResEmpty(), errors.New(mvconst.EXCEPTION_SPECIAL_RESULT_RETURN_NOENCODE)
	}
	body, _ := ioutil.ReadAll(in.Body)
	in.Body.Close()
	r := mvutil.RequestParams{}
	//var rawQuery mvutil.RequestQueryMap
	rawQuery := RenderUcReqParam(body, &r)
	platformId, _ := rawQuery.GetInt("platform", 0)
	if platformId != 1 {
		//return mvconst.EXCEPTION_APP_PLATFORM_ERROR, errors.New("EXCEPTION_APP_PLATFORM_ERROR")
		return nil, errorcode.EXCEPTION_APP_PLATFORM_ERROR
	}
	RenderReqParam(in, &r, rawQuery)
	return &r, nil
}

func RenderUcReqParam(body []byte, r *mvutil.RequestParams) mvutil.RequestQueryMap {
	var ucParams ucParams
	var unitIdList []string
	var requestList []string
	var requestData RequestData
	var assets []Assets
	json.Unmarshal(body, &ucParams)
	ucRequest := ucParams.Imp
	var bidfloor float64
	for _, v := range ucRequest {
		if v.Bidfloor > 0.0 {
			r.Param.NewRTBFlag = true
			bidfloor = v.Bidfloor
		}
		unitIdList = append(unitIdList, v.Id)
		requestList = append(requestList, v.Native.Request)
	}
	WMin := ""
	HMin := ""
	json.Unmarshal([]byte(requestList[0]), &requestData)

	if len(requestData.Assets) > 0 {
		assets = requestData.Assets
	} else {
		assets = requestData.NativeData.Assets
	}
	for _, v := range assets {
		if v.Id == 4 {
			HMin = strconv.Itoa(v.Img.Hmin)
			WMin = strconv.Itoa(v.Img.Wmin)
		}
	}
	query := make(mvutil.RequestQueryMap)
	query["app_id"] = []string{ucParams.App.Id}
	query["package_name"] = []string{ucParams.App.Bundle}
	query["app_version_name"] = []string{ucParams.App.Ver}
	query["useragent"] = []string{ucParams.Device.Ua}
	query["client_ip"] = []string{ucParams.Device.Ip}
	query["brand"] = []string{ucParams.Device.Make}
	query["model"] = []string{ucParams.Device.Model}
	query["platform"] = []string{strconv.Itoa(mvutil.GetPlatformStr(ucParams.Device.Os))}
	query["language"] = []string{ucParams.Device.Language}
	query["network_type"] = []string{getMapNetworkType(strconv.Itoa(ucParams.Device.ConnectionType))}
	if mvutil.GetPlatformStr(ucParams.Device.Os) == mvconst.PlatformAndroid {
		query["gaid"] = []string{ucParams.Device.Ifa}
	} else if mvutil.GetPlatformStr(ucParams.Device.Os) == mvconst.PlatformIOS {
		query["idfa"] = []string{ucParams.Device.Ifa}
	}
	query["unit_id"] = []string{unitIdList[0]}
	imageSize := WMin + "x" + HMin
	query["image_size"] = []string{strconv.Itoa(mvconst.GetImageSizeID(imageSize))}
	query["ad_num"] = []string{"1"}
	query["ecpm_floor"] = []string{strconv.FormatFloat(bidfloor, 'f', 6, 64)}
	r.Param.UCResponseId = ucParams.Id
	return query
}

func getMapNetworkType(networkType string) string {
	switch networkType {
	case "2":
		return "9"
	case "4":
		return "2"
	case "5":
		return "3"
	case "6":
		return "4"
	}
	return "0"
}

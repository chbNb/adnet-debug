package process_pipeline

import (
	"gitlab.mobvista.com/ADN/adnet/internal/mvconst"
	"gitlab.mobvista.com/ADN/adnet/internal/output"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	jsoniter "github.com/json-iterator/go"
	"gitlab.mobvista.com/ADN/adnet/internal/errorcode"
	"gitlab.mobvista.com/ADN/adnet/internal/extractor"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
)

type HUPUReqparamTransFilter struct {
}

type hupuParams struct {
	Id string `json:"id"`
	//App    App
	HupuDevice HupuDevice `json:"device"`
	Imp        []HupuImp  `json:"imp"`
}

type HupuImp struct {
	Id     int `json:"id"`
	TagId  int `json:"tagid"`
	Pmp    Pmp `json:"pmp"`
	Secure int `json:"secure"`
}

type Pmp struct {
	//Id    string  `json:"id"`
	Deals []Deals `json:"deals"`
}

type Deals struct {
	Id string `json:"id"`
}

type HupuDevice struct {
	Ua             string `json:"ua"`
	IP             string `json:"ip"`
	Osv            string `json:"osv"`
	Idfa           string `json:"idfa"`
	Model          string `json:"model"`
	Os             string `json:"os"`
	Imei           string `json:"imei"`
	Androidid      string `json:"androidid"`
	Make           string `json:"make"`
	Oaid           string `json:"oaid"`
	ConnectionType int    `json:"connectiontype"`
	ScreenWidth    int    `json:"screenwidth"`
	ScreenHeight   int    `json:"screenheight"`
}

func (hrtf *HUPUReqparamTransFilter) Process(data interface{}) (interface{}, error) {
	in, ok := data.(*http.Request)
	if !ok {
		return nil, errorcode.EXCEPTION_SPECIAL_EMPTY_STRING
	}
	body, _ := ioutil.ReadAll(in.Body)
	in.Body.Close()
	r := mvutil.RequestParams{}

	rawQuery, err := RenderHupuReqParam(body, &r)
	if err != nil {
		mvutil.Logger.Runtime.Warnf("Hupu Render Params Error. Err=[%s]", err)
		return nil, errorcode.EXCEPTION_PARAMS_ERROR
	}

	RenderReqParam(in, &r, rawQuery)
	return &r, nil
}

func RenderHupuReqParam(body []byte, r *mvutil.RequestParams) (map[string][]string, error) {
	var hupuParams hupuParams
	err := jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(body, &hupuParams)

	if err != nil {
		return nil, err
	}

	impArr := hupuParams.Imp[0]
	pmpDeals := impArr.Pmp.Deals[0]
	pmpId := pmpDeals.Id
	deviceArr := hupuParams.HupuDevice

	query := make(mvutil.RequestQueryMap)

	var platform string
	platformStr := strings.ToLower(hupuParams.HupuDevice.Os)
	if platformStr == "android" {
		platform = "1"
	} else if platformStr == "ios" {
		platform = "2"
	}
	query["platform"] = mvutil.RenderQueryMapV(platform)
	unitId := getUnitId(platform, impArr.TagId, pmpId)
	query["unit_id"] = mvutil.RenderQueryMapV(strconv.FormatInt(unitId, 10))
	appId := getAppId(unitId)
	query["app_id"] = mvutil.RenderQueryMapV(strconv.FormatInt(appId, 10))
	httpReq := getHttpReq(&impArr)
	query["http_req"] = mvutil.RenderQueryMapV(strconv.Itoa(httpReq))

	query["useragent"] = mvutil.RenderQueryMapV(deviceArr.Ua)
	query["client_ip"] = mvutil.RenderQueryMapV(deviceArr.IP)
	query["os_version"] = mvutil.RenderQueryMapV(deviceArr.Osv)
	query["idfa"] = mvutil.RenderQueryMapV(deviceArr.Idfa)
	query["android_id"] = mvutil.RenderQueryMapV(deviceArr.Androidid)
	query["imei"] = mvutil.RenderQueryMapV(deviceArr.Imei)
	query["brand"] = mvutil.RenderQueryMapV(deviceArr.Make)
	query["model"] = mvutil.RenderQueryMapV(hupuParams.HupuDevice.Model)
	networkType := getNetworkType(&deviceArr)
	query["network_type"] = mvutil.RenderQueryMapV(strconv.Itoa(networkType))
	screenSize := getHupuScreenSize(&deviceArr)
	query["screen_size"] = mvutil.RenderQueryMapV(screenSize)
	query["ad_num"] = mvutil.RenderQueryMapV("1")
	query["oaid"] = mvutil.RenderQueryMapV(hupuParams.HupuDevice.Oaid)

	// 处理返回的impid
	if output.IsHupuSplash(pmpId) {
		r.Param.HupuImpId = "1"
		query["image_size"] = mvutil.RenderQueryMapV("5")
	} else {
		r.Param.HupuImpId = "3"
	}
	r.Param.DealId = pmpId
	r.Param.HupuRequestId = hupuParams.Id

	return query, nil
}

// 从配置中取unitid
func getUnitId(platform string, tagid int, pmpId string) int64 {
	tagidStr := strconv.Itoa(tagid)
	//优先使用映射值，若无映射配置，则使用默认广告位
	hupuUnitMapConf := extractor.GetNEW_HUPU_UNITID_MAP()
	if platformConf, ok := hupuUnitMapConf[platform]; ok {
		// 开屏位置。为了改动最小，直接复用NEW_HUPU_UNITID_MAP 配置

		// 改为优先使用pmpid（虎扑口中的dealid）。hupu新版本的流量都使用dealid映射了。
		if mvUnit, ok := platformConf[pmpId]; ok && mvUnit > 0 {
			return mvUnit
		}

		// 信息流位置
		if mvUnit, ok := platformConf[tagidStr]; ok && mvUnit > 0 {
			return mvUnit
		}
	}

	unitConf, ok := extractor.GetHUPU_DEFAULT_UNITID()
	if !ok {
		mvutil.Logger.Runtime.Warnf("HupuReqparamTransFilter get HUPU_DEFAULT_UNITID error")
	}
	unitId := unitConf[platform]
	return unitId
}

func getAppId(unitId int64) int64 {
	//id, err := strconv.Atoi(unitId)
	//if err != nil {
	//	mvutil.Logger.Runtime.Warnf("HupuReqparamTransFilter data type conversion failed")
	//	return 0
	//}

	unitInfo, ok := extractor.GetUnitInfo(unitId)
	if !ok {
		mvutil.Logger.Runtime.Warnf("HupuReqparamTransFilter get unitInfo error")
		return 0
	}
	appId := unitInfo.AppId
	return appId

}

// http&https
func getHttpReq(impArr *HupuImp) int {
	secure := impArr.Secure
	var httpReq int = 1
	if secure == 1 {
		httpReq = 2
	}
	return httpReq
}

// networkType
func getNetworkType(device *HupuDevice) int {
	connectiontype := device.ConnectionType
	// 映射关系
	networkMap := map[int]int{
		0: mvconst.NETWORK_TYPE_UNKNOWN,
		1: mvconst.NETWORK_TYPE_UNKNOWN,
		2: mvconst.NETWORK_TYPE_WIFI,
		3: mvconst.NETWORK_TYPE_UNKNOWN,
		4: mvconst.NETWORK_TYPE_2G,
		5: mvconst.NETWORK_TYPE_3G,
		6: mvconst.NETWORK_TYPE_4G,
		7: mvconst.NETWORK_TYPE_5G,
	}
	networkType, ok := networkMap[connectiontype]
	if !ok {
		return mvconst.NETWORK_TYPE_UNKNOWN
	}
	return networkType
}

func getHupuScreenSize(device *HupuDevice) string {
	scrWidth := device.ScreenWidth
	scrHeight := device.ScreenHeight
	screenSize := strconv.Itoa(scrWidth) + "x" + strconv.Itoa(scrHeight)
	if screenSize != "0x0" {
		return screenSize
	}
	return ""
}

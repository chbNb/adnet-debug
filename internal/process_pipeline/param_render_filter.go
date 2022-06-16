package process_pipeline

import (
	"errors"
	"fmt"
	"hash/crc32"
	"math"
	"math/rand"
	"net/url"
	"strconv"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
	"gitlab.mobvista.com/ADN/adnet/internal/errorcode"
	"gitlab.mobvista.com/ADN/adnet/internal/extractor"
	"gitlab.mobvista.com/ADN/adnet/internal/geo"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/constant"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/filter"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/helpers"
	hbreqctx "gitlab.mobvista.com/ADN/adnet/internal/hb/req_context"
	"gitlab.mobvista.com/ADN/adnet/internal/mkv"
	"gitlab.mobvista.com/ADN/adnet/internal/mvconst"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	"gitlab.mobvista.com/ADN/adnet/internal/output"
	"gitlab.mobvista.com/ADN/adnet/internal/redis"
	"gitlab.mobvista.com/ADN/adnet/internal/req_context"
	"gitlab.mobvista.com/ADN/adnet/internal/utility"
	"gitlab.mobvista.com/ADN/adnet/internal/uuid"
	"gitlab.mobvista.com/ADN/chasm/module/demand"
	supply_mvutil "gitlab.mobvista.com/ADN/chasm/module/supply/mvutil"
	"gitlab.mobvista.com/ADN/exporter/metrics"
	smodel "gitlab.mobvista.com/ADN/structs/model"
	"gitlab.mobvista.com/ADN/treasure_box_sdk/tb_tools"
)

type Ints []int64

type Strs []string

type ParamRenderFilter struct {
}

type DeviceInfo struct {
	Imei         string  `json:"imei"`
	Mac          string  `json:"mac"`
	AndroidID    string  `json:"android_id"`
	Lat          string  `json:"lat"`
	Lng          string  `json:"lng"`
	Gpst         string  `json:"gpst"`
	GpstAccuracy string  `json:"gps_accuracy"`
	OAID         string  `json:"oaid"`
	IMSI         string  `json:"imsi"`
	Dmt          string  `json:"dmt"` // 设备内存总空间
	Dmf          float64 `json:"dmf"` // 设备内存剩余空间
	CpuType      string  `json:"ct"`  // cpu type
	// GpsType      int    `json:"gps_type"`
}

type Tki struct {
	OsVersionUpdateTime string `json:"os_version_up_time"`
	Carrier             string `json:"carrier"`
	Ram                 string `json:"ram"`
	UpdateTime          string `json:"uptime"`
	NewId               string `json:"new_id,omitempty"`
	OldId               string `json:"old_id,omitempty"`
	Abstract            string `json:"abstract,omitempty"`
}

type NewTki struct {
	OsVersionUpdateTime string `json:"4"`
	Carrier             string `json:"1"`
	Ram                 string `json:"2"`
	UpdateTime          string `json:"3"`
	Abstract            string `json:"5"`
}

type TrafficInfo struct {
	StackListStr     string `json:"1"`
	ClassNameListStr string `json:"2"`
	ProtocolListStr  string `json:"3"`
}

func parseDVI(r *mvutil.RequestParams) {
	if len(r.Param.NativeInfo) > 0 {
		nativeInfo := r.Param.NativeInfo
		var nativeInfoList mvutil.NativeInfoEntrys
		err := jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal([]byte(nativeInfo), &nativeInfoList)
		if err != nil {
			nativeInfo = strings.Replace(nativeInfo, "id", "\"id\"", -1)
			nativeInfo = strings.Replace(nativeInfo, "ad_num", "\"ad_num\"", -1)
			err := jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal([]byte(nativeInfo), &nativeInfoList)
			if err == nil {
				r.Param.NativeInfoList = nativeInfoList
			}
		} else {
			r.Param.NativeInfoList = nativeInfoList
		}

	}

	if len(r.Param.DVI) > 0 {
		deviceInfoDecode := decodeByPath(r.Param.RequestPath, r.Param.DVI)
		var deviceInfo DeviceInfo
		err := jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal([]byte(deviceInfoDecode), &deviceInfo)
		if err == nil {
			if len(deviceInfo.Imei) > 0 {
				r.Param.IMEI = deviceInfo.Imei
			}

			if len(deviceInfo.Mac) > 0 {
				r.Param.MAC = deviceInfo.Mac
			}

			if len(deviceInfo.AndroidID) > 0 {
				r.Param.AndroidID = deviceInfo.AndroidID
			}

			if len(deviceInfo.Lat) > 0 {
				r.Param.LAT = deviceInfo.Lat
			}

			if len(deviceInfo.Lng) > 0 {
				r.Param.LNG = deviceInfo.Lng
			}

			if len(deviceInfo.Gpst) > 0 {
				r.Param.GPST = deviceInfo.Gpst
			}

			if len(deviceInfo.GpstAccuracy) > 0 {
				r.Param.GPSAccuracy = deviceInfo.GpstAccuracy
			}

			if len(deviceInfo.OAID) > 0 {
				r.Param.OAID = deviceInfo.OAID
			}

			if len(deviceInfo.IMSI) > 0 {
				r.Param.IMSI = deviceInfo.IMSI
			}

			// 安卓的dmt，dmf，cpu type存储在dvi内，需要区分处理.如在参数获取不到，则从dvi里获取。
			if r.Param.Dmt == 0 {
				dmt, err := strconv.ParseFloat(deviceInfo.Dmt, 64)
				if err == nil {
					r.Param.Dmt = dmt
				}
			}
			if r.Param.Dmf == 0 {
				r.Param.Dmf = deviceInfo.Dmf
			}
			if len(r.Param.Ct) == 0 {
				r.Param.Ct = deviceInfo.CpuType
			}

			// if deviceInfo.GpsType > 0 {
			// 	r.Param.GPSType = deviceInfo.GpsType
			// }
		}
	}

	if len(r.Param.IMEI) == 0 && len(r.Param.D1) > 0 {
		r.Param.IMEI = decodeByPath(r.Param.RequestPath, r.Param.D1)
	}
	if len(r.Param.MAC) == 0 && len(r.Param.D2) > 0 {
		r.Param.MAC = decodeByPath(r.Param.RequestPath, r.Param.D2)
	}
	if len(r.Param.AndroidID) == 0 && len(r.Param.D3) > 0 {
		r.Param.AndroidID = decodeByPath(r.Param.RequestPath, r.Param.D3)
	}

	// trim
	r.Param.IMEI = mvutil.TrimOnlyaA0New(mvutil.TrimBlank(r.Param.IMEI))

	// imei abtest 背景：sdk 在获取不到imei的时候，会用其他信息md5生成一个值，传给adnet。
	// abtest的目的想确认传imei和传错误的imei的效果分别是怎样的
	ImeiAbTest(r)

	r.Param.AndroidID = mvutil.TrimOnlyaA0New(mvutil.TrimBlank(r.Param.AndroidID))
	r.Param.MAC = mvutil.TrimBlank(r.Param.MAC)
	if mvutil.HasInvisibleChar(r.Param.MAC) {
		r.Param.MAC = ""
	}
	r.Param.SDKVersion = mvutil.TrimAa0AndDotNew(r.Param.SDKVersion)
}

func ImeiAbTest(r *mvutil.RequestParams) {
	// 针对sdk 流量处理
	if !mvutil.IsHbOrV3OrV5Request(r.Param.RequestPath) || r.Param.RequestType != mvconst.REQUEST_TYPE_OPENAPI_V3 {
		return
	}
	// 有效的imei 不参与实验
	// 限制切量流量
	if len(r.Param.IMEI) != 32 || r.Param.IMEI == "00000000000000000000000000000000" {
		return
	}

	adnSwitchConfs, _ := extractor.GetADNET_SWITCHS()
	randVal := rand.Intn(100)
	if imeiAbTestRate, ok := adnSwitchConfs["imeiAbTestRate"]; ok {
		// 切量部分把imei清空
		if imeiAbTestRate > randVal {
			r.Param.IMEI = ""
			// 记录标记
			r.Param.ExtDataInit.ImeiABTest = 1
		} else {
			r.Param.ExtDataInit.ImeiABTest = 2
		}
	}
}

func decodeByPath(path string, data string) string {
	var res string
	if path == mvconst.PATHMPNewAD {
		res = mvutil.Base64DecodeMPNewAd(data)
	} else {
		res = mvutil.Base64Decode(data)
	}
	return res
}

// /////////todo
func getScreenSizeNew(screenSize string) (error, string, int, int) {
	screenList := strings.Split(screenSize, "x")
	if len(screenList) != 2 {
		return errors.New("screenSize is invalidate"), "", 0, 0
	}
	width := strings.TrimSpace(screenList[0])
	height := strings.TrimSpace(screenList[1])
	var iwidth, iheight int
	var err error
	if strings.Contains(width, ".") {
		fwidth, err := strconv.ParseFloat(width, 64)
		if err != nil {
			return err, "", 0, 0
		}
		iwidth = int(fwidth)

		fheight, err := strconv.ParseFloat(height, 64)
		if err != nil {
			return err, "", 0, 0
		}
		iheight = int(fheight)
		return nil, fmt.Sprintf("%dx%d", iwidth, iheight), iwidth, iheight
	} else {
		iwidth, err = strconv.Atoi(width)
		if err != nil {
			return err, "", 0, 0
		}

		iheight, err = strconv.Atoi(height)
		if err != nil {
			return err, "", 0, 0
		}
		return nil, screenSize, iwidth, iheight
	}
}

func formatScreenSize(r *mvutil.RequestParams) {
	screeSize := r.Param.ScreenSize
	err, ss, w, h := getScreenSizeNew(screeSize)
	if err != nil {
		return
	}
	r.Param.ScreenSize = ss
	r.Param.ScreenWidth = w
	r.Param.ScreenHeigh = h
}

func handleScenario(r *mvutil.RequestParams) {
	if r.Param.RequestPath == mvconst.PATHREQUEST {
		return
	}
	scenarios, ifFind := extractor.GetOpenapiScenario()
	if !ifFind || len(scenarios) <= 0 {
		r.Param.Scenario = mvconst.SCENARIO_OPENAPI
		return
	}
	if mvutil.InStrArray(r.Param.Scenario, scenarios) {
		return
	}
	r.Param.Scenario = mvconst.SCENARIO_OPENAPI
}

type RankerInfo struct {
	PowerRate      int     `json:"power_rate"`
	Charging       int     `json:"charging"`
	TotalMemory    string  `json:"total_memory"`
	ResidualMemory string  `json:"residual_memory"`
	CID            string  `json:"cid"`
	LAT            string  `json:"lat"`
	LNG            string  `json:"lng"`
	GPST           string  `json:"gpst"`
	GPSAccuracy    string  `json:"gps_accuracy"`
	GPSType        string  `json:"gps_type"`
	Dmt            float64 `json:"dmt,omitempty"`
	Dmf            float64 `json:"dmf,omitempty"`
	CpuType        string  `json:"ct,omitempty"`
	ChannelInfo    string  `json:"topon_info,omitempty"`
	PriceFactor    float64 `json:"pf,omitempty"`
	HBMn           string  `json:"hb_mn,omitempty"`
}

func handleRankerInfo(r *mvutil.RequestParams) {
	adnConf, _ := extractor.GetADNET_SWITCHS()
	var rankerInfo RankerInfo
	rankerInfo.PowerRate = r.Param.PowerRate
	rankerInfo.Charging = r.Param.Charging
	rankerInfo.TotalMemory = r.Param.TotalMemory
	rankerInfo.ResidualMemory = r.Param.ResidualMemory
	rankerInfo.CID = r.Param.CID
	rankerInfo.LAT = r.Param.LAT
	rankerInfo.LNG = r.Param.LNG
	rankerInfo.GPST = r.Param.GPST
	rankerInfo.GPSAccuracy = r.Param.GPSAccuracy
	// todo
	rankerInfo.GPSType = strconv.Itoa(r.Param.GPSType)
	// 设置开关传递dmt，dmf，ct给算法
	if dmSwitch, ok := adnConf["dmSwitch"]; ok && dmSwitch == 1 {
		rankerInfo.Dmt = r.Param.Dmt
		rankerInfo.Dmf = r.Param.Dmf
		rankerInfo.CpuType = r.Param.Ct
	}
	rankerInfo.ChannelInfo = r.Param.ChannelInfo
	fcc := extractor.GetFREQ_CONTROL_CONFIG()
	if fcc != nil && fcc.FreqControlToRs == 1 && fcc.Status == 1 &&
		r.Param.ExtDataInit.PriceFactor > mvconst.PriceFactor_MINValue &&
		r.Param.ExtDataInit.PriceFactor <= mvconst.PriceFactor_MAXValue {
		// 开关开启&符合标准的才传给RS
		rankerInfo.PriceFactor = r.Param.ExtDataInit.PriceFactor
	}
	if len(r.Param.MediationName) > 0 {
		rankerInfo.HBMn = r.Param.MediationName
	}
	jsonvalue, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(rankerInfo)
	r.Param.RankerInfo = string(jsonvalue)
}

func parseNums(r *mvutil.RequestParams) {
	adType := r.Param.AdType
	trueNum := r.Param.TNum
	apiRequestNum := r.Param.ApiRequestNum
	apiCacheNum := r.Param.ApiCacheNum
	if adType == mvconst.ADTypeFeedsVideo || adType == mvconst.ADTypeNative {
		if trueNum > 0 {
			// 配置为开发者实际请求数
			if apiRequestNum == -2 {
				r.Param.AdNum = int32(trueNum)
				// 对于apiRequestNum == -1的情况在参数过滤中已经处理为异常
			} else {
				r.Param.AdNum = apiRequestNum
			}
		} else {
			if apiRequestNum != -2 {
				r.Param.AdNum = apiRequestNum
			}
		}

	} else if adType == mvconst.ADTypeOfferWall || adType == mvconst.ADTypeRewardVideo {
		r.Param.AdNum = apiRequestNum
		r.Param.TNum = apiCacheNum
	} else if adType == mvconst.ADTypeInterstitialSdk {
		r.Param.AdNum = apiRequestNum
		if apiCacheNum > 0 {
			r.Param.TNum = apiCacheNum
		}
	} else if adType == mvconst.ADTypeInterstitialVideo {
		r.Param.AdNum = apiRequestNum
		r.Param.TNum = 1
	} else if mvutil.IsWxAdType(adType) {
		r.Param.TNum = int(apiRequestNum)
		r.Param.AdNum = apiRequestNum
	}

	// give V4 api right AdNum
	if r.Param.RequestPath == mvconst.PATHOpenApiV4 && r.Param.FrameNum >= 0 && len(r.Param.NativeInfo) > 0 {
		// try to solve adnum for V4
		r.Param.RequestType = mvconst.REQUEST_TYPE_OPENAPI_V3 // Attention, V4 for adserver is be equal to V3!
		for _, v := range r.Param.NativeInfoList {
			r.Param.AdNum = int32(r.Param.FrameNum * v.RequireNum) // use adnnet logic
			r.Param.RequireNum = v.RequireNum
			if v.AdTemplate == mvconst.TemplateMultiIcons { // I guess the const should be this...
				break
			}
		}
		r.Param.TNum = int(r.Param.AdNum) // I guess we should use this?
		if r.Param.RequireNum == 0 {
			r.Param.RequireNum = 1
		}
	}
	// adNum验证
	if r.Param.AdNum <= int32(0) {
		r.Param.AdNum = int32(5)
	} else if r.Param.AdNum > int32(50) {
		r.Param.AdNum = int32(50)
	}

	// 对tnum做最后的检查 小于等于0的情况用默认值，大于MaxTrueNum的情况用MaxTrueNum, AdNum不做处理
	if r.Param.TNum <= 0 {
		r.Param.TNum = getDefaultTNum(adType)
	}
	if r.Param.TNum > mvconst.MaxTrueNum {
		r.Param.TNum = mvconst.MaxTrueNum
	}
	// fix hb online api ad_num
	if r.IsHBRequest {
		confMap := extractor.GetOnlinePublisherAdNumConfig()
		if adNum, ok := confMap[strconv.FormatInt(r.Param.PublisherID, 10)]; ok && r.Param.AdNum > 10 {
			r.Param.AdNum = adNum
			r.Param.TNum = int(adNum)
		}
	}
}

func getDefaultTNum(adType int32) int {
	conf := extractor.GetTRUE_NUM_BY_AD_TYPE()
	adTypeStr := strconv.Itoa(int(adType))
	if num, ok := conf[adTypeStr]; ok {
		return num
	}
	if num, ok := conf["default"]; ok {
		return num
	}
	// 没有配置的广告类型，默认为1
	return 1
}

func impCapBlock(r *mvutil.RequestParams) {
	if r.Param.FcaSwitch {
		if mvutil.AppFcaDefault(r.AppInfo.App.FrequencyCap) {
			return
		}
	}

	placementImpCapSwitch := extractor.GetUSE_PLACEMENT_IMP_CAP_SWITCH()
	// 旧逻辑的判断
	oldLogic := !placementImpCapSwitch && ((r.UnitInfo.Unit.ImpressionCap == -1 && r.AppInfo.App.ImpressionCap > 0) ||
		r.UnitInfo.Unit.ImpressionCap > 0)
	// 新逻辑的判断
	newLogic := placementImpCapSwitch && r.PlacementInfo != nil &&
		r.PlacementInfo.ImpressionCap > 0 && r.PlacementInfo.ImpressionCapPeriod > 0
	// 先收集数据，不用开启开关也能收集到对应的数据
	if r.PlacementInfo != nil && r.PlacementInfo.ImpressionCap > 0 && r.PlacementInfo.ImpressionCapPeriod > 0 {
		r.Param.ExtDataInit.ImpressionCap = r.PlacementInfo.ImpressionCap
	}

	// 更新控IsBlockByImpcap
	if (mvutil.IsHbOrV3OrV5Request(r.Param.RequestPath) ||
		r.Param.RequestPath == mvconst.PATHJssdkApi ||
		r.Param.RequestPath == mvconst.PATHOnlineApi ||
		(r.Param.RequestPath == mvconst.PATHMPAD || r.Param.RequestPath == mvconst.PATHMPADV2)) &&
		(oldLogic || newLogic) {
		// watcher.AddWatchValue("impCapQueryRedis", float64(1))
		// 遵循app的配置
		deviceId := mvutil.GetGlobalDeviceTag(&r.Param)
		if deviceId == "" {
			return
		}

		// 去Aerospike获取数据
		data, err := getFreqFromAerospike(r, strings.ToLower(deviceId))
		if err != nil {
			mvutil.Logger.Runtime.Warnf("impCapBlock(),mkv is error!" + err.Error())
			return
		}

		if oldLogic {
			oldImpCapBlock(r, data)
		}
		if newLogic {
			newImpCapBlock(r, data)
		}

		if r.Param.IsBlockByImpCap && extractor.GetONLY_REQUEST_THIRD_DSP_SWITCH() {
			r.Param.OnlyRequestThirdDsp = true
		}
	}
}

func oldImpCapBlock(r *mvutil.RequestParams, data map[string][]byte) {
	// 获取DEV的逻辑
	devCap := 0
	subKey := "if_" + strconv.FormatInt(r.Param.AppID, 10)
	if jsonStr, ok := data[subKey]; ok {
		var json = jsoniter.ConfigCompatibleWithStandardLibrary
		value := new(mvutil.FreqImpCapMkvData)
		err := json.Unmarshal(jsonStr, &value)
		if err != nil {
			mvutil.Logger.Runtime.Warnf("impCapBlock(),mkv is error!" + err.Error())
			return
		}
		if time.Now().Unix()-value.Ts <= 24*60*60 {
			devCap = value.Count
		}
	}

	if r.UnitInfo.Unit.ImpressionCap == -1 && r.AppInfo.App.ImpressionCap <= devCap {
		r.Param.IsBlockByImpCap = true
	}
	if r.UnitInfo.Unit.ImpressionCap > 0 && r.UnitInfo.Unit.ImpressionCap <= devCap {
		r.Param.IsBlockByImpCap = true
	}
}

func newImpCapBlock(r *mvutil.RequestParams, data map[string][]byte) {
	if r.PlacementInfo == nil || r.PlacementInfo.ImpressionCap <= 0 || r.PlacementInfo.ImpressionCapPeriod <= 0 {
		return
	}
	r.Param.ExtDataInit.ImpressionCapTime = 1 // 进入默认值为1
	// 获取DEV的逻辑
	devCap := 0
	subKey := "ic_" + strconv.FormatInt(r.PlacementInfo.PlacementId, 10)
	if jsonStr, ok := data[subKey]; ok {
		var json = jsoniter.ConfigCompatibleWithStandardLibrary
		value := new(mvutil.FreqImpCapMkvData)
		err := json.Unmarshal(jsonStr, &value)
		if err != nil {
			mvutil.Logger.Runtime.Warnf("impCapBlock(),mkv is error!" + err.Error())
			return
		}
		if time.Now().Unix()-value.Ts <= int64(r.PlacementInfo.ImpressionCapPeriod)*3600 {
			devCap = value.Count
		}
		r.Param.ExtDataInit.ImpressionCapTime = value.Ts
	}

	if r.PlacementInfo.ImpressionCap <= devCap {
		r.Param.IsBlockByImpCap = true
	}
}

func renderReduceFillConfig(r *mvutil.RequestParams) {
	reduceFillConfig, found := getReduceFillConfig(r.Param.UnitID, r.Param.AppID, r.Param.Platform, r.Param.CountryCode)
	if found {
		r.Param.ReduceFillConfig = reduceFillConfig
	}
}

func getReduceFillConfig(unitid int64, appid int64, platform int, country string) (cfg *smodel.ConfigAlgorithmFillRate, ifFind bool) {
	cfg, ifFind = extractor.GetFillRateControllConfig(extractor.GetFillRateKey(unitid, appid, platform, country))
	if ifFind {
		return cfg, ifFind
	}

	cfg, ifFind = extractor.GetFillRateControllConfig(extractor.GetFillRateKey(unitid, appid, platform, "ALL"))
	return cfg, ifFind
}

func getFillEcpmFloor(r *mvutil.RequestParams) (float64, error) {
	key := strconv.FormatInt(r.Param.UnitID, 10) + "_" + strings.ToLower(r.Param.CountryCode)
	allkey := strconv.FormatInt(r.Param.UnitID, 10) + "_all"

	// 区分素材2,3期
	if output.IsRvNvNewCreative(r.Param.NewCreativeFlag, r.Param.AdType) {
		key += "_1"
		allkey += "_1"
	} else {
		key += "_0"
		allkey += "_0"
	}
	fillEcpmFloorKey := key
	valueStr, err := redis.LocalRedisAlgoHGet("reduce_fill", key)
	if err != nil {
		fillEcpmFloorKey = allkey
		valueStr, err = redis.LocalRedisAlgoHGet("reduce_fill", allkey)
		if err != nil {
			return 0.0, err
		}
	}
	r.Param.FillEcpmFloorKey = fillEcpmFloorKey
	r.Param.ExtReduceFillValue = valueStr

	ecpmFloorPriceStr := valueStr
	ecpmFloorPriceVer := "0"
	// 兼容带版本号的 value
	// fillrate_yyyyhhmmHHMM
	valList := strings.Split(ecpmFloorPriceStr, "_")
	if len(valList) > 1 {
		ecpmFloorPriceStr = valList[0]
		ecpmFloorPriceVer = valList[1]
	}
	r.Param.FillEcpmFloorVer = ecpmFloorPriceVer

	fillEcpmFloor, err := strconv.ParseFloat(ecpmFloorPriceStr, 64)
	if err != nil {
		mvutil.Logger.Runtime.Warnf("key[%s] reduce fill ecpm_floor[%s] ver[%s] parse err:%v", fillEcpmFloorKey, ecpmFloorPriceStr, ecpmFloorPriceVer, err)
		return 0.0, err
	}

	return fillEcpmFloor, nil
}

func renderEcpmFloor(r *mvutil.RequestParams) float64 {
	// 现有使用新频次控制的实验组的配置
	fcc := extractor.GetFREQ_CONTROL_CONFIG()
	if r.Param.ExtDataInit.PriceFactorTag == mvconst.PriceFactorTag_B &&
		r.Param.ExtDataInit.PriceFactorHit == mvconst.PriceFactorHit_TRUE && // 命中实验组
		r.Param.ReduceFillConfig != nil && r.Param.ReduceFillConfig.ControlMode == mvconst.EcpmFloor &&
		fcc.UseFixedEcpm == 1 && fcc.FixedEcpmFactor > 0 {
		if value, found := extractor.GetSspProfitDistributionRuleByUnitIdAndCountryCode(r.UnitInfo.UnitId, r.Param.CountryCode); found && value.Type == mvconst.SspProfitDistributionRuleFixedEcpm && value.FixedEcpm > 0 {
			return fcc.FixedEcpmFactor * value.FixedEcpm // 返回控制后的floor
		}
	}

	sspEcpmFloor := r.Param.BidFloor
	amEcpmFloor, pubEcpmFloor := 0.0, 0.0
	if len(r.UnitInfo.EcpmFloors) > 0 {
		if efloor, ok := r.UnitInfo.EcpmFloors[r.Param.CountryCode]; ok {
			pubEcpmFloor = efloor
		}
	}

	if r.Param.ReduceFillConfig != nil && r.Param.ReduceFillConfig.ControlMode == mvconst.EcpmFloor {
		amEcpmFloor = r.Param.ReduceFillConfig.EcpmFloor
	}

	if amEcpmFloor > sspEcpmFloor {
		return amEcpmFloor
	}

	if amEcpmFloor <= 0.0 && pubEcpmFloor > sspEcpmFloor {
		return pubEcpmFloor
	}
	return sspEcpmFloor
}

func fillAppInfo(r *mvutil.RequestParams) error {
	appInfo, ifFind := extractor.GetAppInfo(r.Param.AppID)
	if ifFind {
		r.AppInfo = appInfo
	} else {
		mvutil.Logger.Runtime.Warnf("request_id=[%s] has filter by app %d entry is not exists", r.Param.RequestID, r.Param.AppID)
		return errorcode.EXCEPTION_APP_NOT_FOUND
		// return mvconst.EXCEPTION_APP_NOT_FOUND, errors.New("EXCEPTION_APP_NOT_FOUND")
	}
	// todo
	r.Param.AppName = appInfo.App.Name
	r.Param.RealPackageName = appInfo.RealPackageName
	if r.Param.RequestType <= 0 {
		r.Param.RequestType = mvconst.REQUEST_TYPE_OPENAPI_V3
	}

	// 若为dsp流量带来的more_offer请求，则使用请求传来的package_name。因为这种情况下是使用固定的appid请求广告，app配置的包名都是一样的
	if r.Param.DspMof == 1 {
		r.Param.ExtfinalPackageName = r.Param.PackageName
	} else {
		r.Param.ExtfinalPackageName = appInfo.RealPackageName
	}
	// 赋值publisherID, publisherType
	r.Param.PublisherID = appInfo.Publisher.PublisherId
	r.Param.PublisherType = appInfo.Publisher.Type
	return nil
}

func fillPublisherInfo(r *mvutil.RequestParams) error {
	publisherInfo, ifFind := extractor.GetPublisherInfo(r.Param.PublisherID)
	if ifFind {
		r.PublisherInfo = publisherInfo
	} else {
		mvutil.Logger.Runtime.Warnf("request_id=[%s][pid:%d] has filter by publisher entry is not exists", r.Param.RequestID, r.Param.PublisherID)
		return errorcode.EXCEPTION_PUBLISHER_NOT_FOUND
		// return mvconst.EXCEPTION_PUBLISHER_NOT_FOUND, errors.New("EXCEPTION_PUBLISHER_NOT_FOUND")
	}
	return nil
}

func fillPlacementInfo(r *mvutil.RequestParams) error {
	if r.Param.FinalPlacementId == 0 {
		return nil
	}
	placementInfo, ifFind := extractor.GetPlacementInfo(r.Param.FinalPlacementId)
	if ifFind && placementInfo.Status == mvutil.ACTIVE {
		r.PlacementInfo = placementInfo
	} else {
		// mvutil.Logger.Runtime.Warnf("request_id=[%s] has filter by placement entry is not exists", r.Param.RequestID)
		return nil
	}
	return nil
}

func fillUnitInfo(r *mvutil.RequestParams) error {
	unitInfo, ifFind := extractor.GetUnitInfo(r.Param.UnitID)
	if ifFind {
		r.UnitInfo = unitInfo
	} else {
		// 兼容v2接口不传unitid的情况
		if r.Param.RequestPath == mvconst.PATHOpenApiV2 || r.Param.RequestPath == mvconst.PATHREQUEST {
			r.UnitInfo = &smodel.UnitInfo{
				Unit: smodel.Unit{},
			}
			return nil
		}
		mvutil.Logger.Runtime.Warnf("request_id=[%s] has filter by UnitID is not exists", r.Param.RequestID)
		return errorcode.EXCEPTION_UNIT_NOT_FOUND
		// return mvconst.EXCEPTION_UNIT_NOT_FOUND, errors.New("EXCEPTION_UNIT_NOT_FOUND")
	}

	if unitInfo.Setting.ApiRequestNum == -1 {
		mvutil.Logger.Runtime.Warnf("request_id=[%s] has filter by UnitInfo.Setting Exception_Adnnum_Set_NONE", r.Param.RequestID)
		return errorcode.EXCEPTION_ADNUM_SET_NONE
		// return mvconst.EXCEPTION_ADNUM_SET_NONE, errors.New("EXCEPTION_ADNUM_SET_NONE")
	}

	// //有传adType
	// 针对mp请求，不做adtype校验。背景：mp希望使用mp banner广告位来请求mv native广告，但是因adtype被过滤掉了
	// more offer请求时固定传ad_type=3，会造成与unit的adtype校验不一致，导致被过滤。
	// topon native 的adtype为43，44，与unit配置不一致，因此需要排除判断
	if r.Param.AdType > 0 && r.UnitInfo.Unit.AdType != r.Param.AdType && !mvutil.IsMpad(r.Param.RequestPath) && r.Param.Mof != 1 &&
		!(r.Param.RequestPath == mvconst.PATHTOPON && mvutil.InInt32Arr(r.Param.AdType, []int32{mvconst.ADTypeNativeVideo, mvconst.ADTypeNativePic})) {
		mvutil.Logger.Runtime.Warnf("request_id=[%s] has filter by unit adtype is not match", r.Param.RequestID)
		return errorcode.EXCEPTION_UNIT_ADTYPE_ERROR
		// return mvconst.EXCEPTION_UNIT_ADTYPE_ERROR, errors.New("EXCEPTION_UNIT_ADTYPE_ERROR")
	} else {
		// 没有传用unit的AdType
		r.Param.AdType = int32(r.UnitInfo.Unit.AdType)

	}

	// 1 Header Bidding, 2 Traditional
	if r.IsHBRequest && r.UnitInfo.Unit.BiddingType == 2 ||
		(!r.IsHBRequest && r.UnitInfo.Unit.BiddingType == 1) {
		return errorcode.EXCEPTION_UNIT_BIDDING_TYPE_ERROR
	}

	r.Param.ApiRequestNum = unitInfo.Setting.ApiRequestNum
	r.Param.ApiCacheNum = unitInfo.Setting.ApiCacheNum

	//fixed onlineAPI的native， No video
	if r.Param.RequestType == mvconst.REQUEST_TYPE_OPENAPI_AD && r.Param.AdType == mvconst.ADTypeNative {
		r.Param.VideoAdType = mvconst.VideoAdTypeNOVideo
	} else { // 获取unit中的videoAds
		r.Param.VideoAdType = unitInfo.Unit.VideoAds
	}
	return nil
}

func valideIVRecallNet(r *mvutil.RequestParams) (int, error) {
	if r.Param.AdType == mvconst.ADTypeInterstitialVideo {
		if r.Param.Orientation == mvconst.ORIENTATION_BOTH &&
			r.Param.FormatOrientation == mvconst.ORIENTATION_BOTH {
			mvutil.Logger.Runtime.Warnf("request_id=[%s] has filter by AdType InterstitialVideo and orientation is invalidate", r.Param.RequestID)
			return 0, errorcode.EXCEPTION_IV_ORIENTATION_INVALIDATE
			// return mvconst.EXCEPTION_IV_ORIENTATION_INVALIDATE, errors.New("EXCEPTION_IV_ORIENTATION_INVALIDATE")
		}

		var RecallNetArr []int
		for _, RecallNet := range strings.Split(r.UnitInfo.Unit.RecallNet, ";") {
			RecallNet, err := strconv.Atoi(RecallNet)
			if err == nil {
				RecallNetArr = append(RecallNetArr, RecallNet)
			}
		}
		// networktype 在范围之外 并且 （属于2G 3G 4G wifi或者没有勾选other）（other的值为0）。
		// networktype 在范围之内或者networktype不属于2G 3G 4G wifi并且勾选了other，则不过滤。
		if !mvutil.InArray(r.Param.NetworkType, RecallNetArr) && (!mvutil.InArray(mvconst.NETWORK_TYPE_UNKNOWN, RecallNetArr) || NetworkTypeInAllowNormalType(r.Param.NetworkType)) {
			mvutil.Logger.Runtime.Warnf("request_id=[%s] has filter by AdType InterstitialVideo and recallnet is invalidate", r.Param.RequestID)
			return 0, errorcode.EXCEPTION_IV_RECALLNET_INVALIDATE
			// return mvconst.EXCEPTION_IV_RECALLNET_INVALIDATE, errors.New("EXCEPTION_IV_RECALLNET_INVALIDATE")
		}
	}
	return 0, nil
}

func NetworkTypeInAllowNormalType(networkType int) bool {
	return networkType == mvconst.NETWORK_TYPE_2G || networkType == mvconst.NETWORK_TYPE_3G || networkType == mvconst.NETWORK_TYPE_4G ||
		networkType == mvconst.NETWORK_TYPE_5G || networkType == mvconst.NETWORK_TYPE_WIFI
}

func (prf *ParamRenderFilter) isBidRequest(scenario string, adType int32) bool {
	return scenario == mvconst.SCENARIO_OPENAPI
}

func (prf *ParamRenderFilter) filterAndroidLowVersion(unitId, appId, publisherId int64) bool {
	filterConditions, _ := extractor.GetAndroidLowVersionFilterCondition()
	if filterConditions == nil {
		return false
	}
	if len(filterConditions.UnitIds) == 0 && len(filterConditions.AppIds) == 0 && len(filterConditions.PublisherIds) == 0 {
		return false
	}
	if len(filterConditions.UnitIds) > 0 && mvutil.InInt64Arr(unitId, filterConditions.UnitIds) {
		return true
	}

	if len(filterConditions.AppIds) > 0 && mvutil.InInt64Arr(appId, filterConditions.AppIds) {
		return true
	}

	if len(filterConditions.PublisherIds) > 0 && mvutil.InInt64Arr(publisherId, filterConditions.PublisherIds) {
		return true
	}
	return false
}

func (prf *ParamRenderFilter) Process(data interface{}) (interface{}, error) {
	in, ok := data.(*mvutil.RequestParams)
	if !ok {
		return nil, errors.New("ParamRenderFilter input type should be *params.RequestParams")
	}
	// TreasureBox缓存系统
	if tb_tools.Enable() {
		in.Param.ExtDataInit.TreasureBoxAbtestTag = 1
	} else {
		in.Param.ExtDataInit.TreasureBoxAbtestTag = 2
	}

	// app,publisher,unit,都是根据ID去缓存中查出详细信息
	err := fillAppInfo(in)
	if err != nil {
		return nil, err
	}

	err = fillPublisherInfo(in)
	if err != nil {
		return nil, err
	}

	err = fillUnitInfo(in)
	if err != nil {
		return nil, err
	}

	// 获取placementId
	renderPlacementId(in)
	// 根据ip填充placementInfo
	fillPlacementInfo(in)

	// 记录mp 流量的 映射关系   将unit、publisher、placement串成字符串
	in.Param.ExtMpNormalMap = output.RenderExtMpNormalMap(in)

	// handle Scenario
	handleScenario(in)
	// 对加密后的鸿蒙info进行解码
	DecryptHarmonyInfo(in)
	// TODO ? 判断开发者是否在排除头部竞价的app里？
	if in.IsHBRequest {
		pubExcludData, ifFind := extractor.GetHBPubExcludeApp()
		fixApp := false
		if ifFind && len(pubExcludData) > 0 && helpers.InInt64Arr(in.Param.PublisherID, pubExcludData) {
			fixApp = true
		}

		if !fixApp && in.UnitInfo.AppId != in.Param.AppID {
			return in, filter.UnitNotFoundAppError
		}

		if fixApp {
			in.Param.AppID = in.UnitInfo.AppId
		}
		// 判断是否支持的广告类型
		in.Param.AdType = in.UnitInfo.Unit.AdType
		adTypes := []int32{mvconst.ADTypeNative, mvconst.ADTypeRewardVideo, mvconst.ADTypeInterstitialVideo,
			mvconst.ADTypeSdkBanner, mvconst.ADTypeSplash, mvconst.ADTypeNativeH5, mvconst.ADTypeBanner,
			mvconst.ADTypeNative, mvconst.ADTypeOnlineVideo}
		if !helpers.InInt32Arr(in.Param.AdType, adTypes) {
			return in, filter.AdTypeNotSupport
		}
	}
	// int2str
	in.Param.PlatformName = mvconst.GetPlatformStr(in.Param.Platform)
	in.Param.NetworkTypeName = mvconst.GetNetworkName(in.Param.NetworkType)

	// TODO extra都是啥
	in.Param.Extra5 = in.Param.RequestID
	in.Param.Extra4 = in.Param.RequestID
	in.Param.ExtpackageName = in.Param.PackageName
	// 记录apiversion
	in.Param.ExtApiVersion = strconv.FormatFloat(in.Param.ApiVersion, 'f', 1, 64)
	in.Param.ExtcdnType = extractor.GetSYSTEM_AREA()
	// 原生广告
	if in.Param.AdType == mvconst.ADTypeNative {
		nativeInfoStr := in.Param.NativeInfo
		nativeInfo := make([]NativeInfo, 0)
		// nativeInfo 字段为空， 默认为大图
		if len(nativeInfoStr) == 0 {
			nativeInfo = append(nativeInfo, NativeInfo{Id: 2})
		} else {
			_ = jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal([]byte(nativeInfoStr), &nativeInfo)
			if len(nativeInfo) == 0 {
				nativeInfo = append(nativeInfo, NativeInfo{Id: 2})
			}
		}
		in.Param.FormatAdType = mvconst.ADTypeNativePic
		if in.Param.VideoAdType != mvconst.VideoAdTypeNOVideo && nativeInfo[0].Id == 2 && in.Param.VideoVersion != "" {
			in.Param.FormatAdType = mvconst.ADTypeNativeVideo
		}
	} else {
		in.Param.FormatAdType = in.Param.AdType
	}

	// TODO 下面这一坨直接看不懂
	in.Param.Domain = renderTrackDomain(in)
	// 中国tracking专线域名切量实验        一部分流量走中国域名的tracking，这里的域名指的是上报给tracking的url中的域名
	renderCNTrackDomainABTest(in)

	// 处理country_code
	renderCountryCode(in)

	// 根据国家返回对应tracking域名
	// 主要应对埃及屏蔽我们的域名的情况
	renderTrackDomainByCountryCode(in)

	// ad_tracking 里非imp，click的埋点域名切量
	replaceAdTrackingByCdnDomain(in)
	// AB测试，我们可以为同一个优化目标(例如优化购买转化率)制定两个方案(比如两个页面)，让一部分用户使用A方案，另一个部分用户使用B方案，统计并对比不同方案的转化率、点击量、留存率等指标，以判断不同方案的优劣并进行决策，从而提升转化率。
	// imp,click,only_imp域名切量cdn abtest
	renderCdnTrackingDomainAbtest(in)

	// decode dvi
	parseDVI(in)

	renderV5Abtest(in)

	// 获取unit fixed_ecpm
	renderUnitFixedEcpm(in)

	if in.Param.Scenario == mvconst.SCENARIO_OPENAPI &&
		(utility.IsLowFlowUnit(in.Param.UnitID, in.Param.AppID, in.Param.PublisherID, in.Param.CountryCode) ||
			utility.IsLowFlowAdType(in.Param.FormatAdType, in.Param.Platform, in.Param.CountryCode, in.Param.FixedEcpm)) {
		in.Param.IsLowFlowUnitReq = true
	}

	// 获取频次控制整体开关
	in.Param.FcaSwitch, _ = extractor.GetFCA_SWITCH()

	// android low version filter
	if in.Param.Platform == mvconst.PlatformAndroid &&
		in.Param.OSVersionCode < 5000000 &&
		prf.filterAndroidLowVersion(in.UnitInfo.UnitId, in.AppInfo.AppId, in.PublisherInfo.PublisherId) {
		return nil, errorcode.EXCEPTION_OS_VER_LOWER
	}

	// imp cap black
	impCapBlock(in)

	// // playable use adserver
	renderPlayableFlag(in)

	// 素材三期新逻辑切量
	renderNewCreativeFlag(in)

	// more offer切量
	renderNewMoreOfferFlag(in)

	// more offer 批量展示上报逻辑切量
	renderMoreOfferNewImp(in)

	// more offer 提供给h5的abtest标记
	// renderMoreOfferAbFlag(in)

	// 新增more offer adtype,按照unit切量
	changeMoreOfferAdType(in)

	// close button ad unit/整体切量
	renderCloseButtonAdFlag(in)

	// 大模板切量逻辑
	renderBigTemplateFlag(in)

	// polaris切量逻辑
	renderPolarisFlag(in)

	// 处理unit orientation
	handleOrientation(in)

	_, err = valideIVRecallNet(in)
	if err != nil {
		return nil, err
	}

	// 获取sdk banner的unit_size
	handleSdkBannerUnitSize(in)
	// format screenSize
	formatScreenSize(in)
	// 修复screensize
	handleScreensize(in)

	// 新频次控制
	renderPriceFactor(in)

	// handle rankInfo
	handleRankerInfo(in)

	// dco abtest实验切量标记
	renderIfSupDco(in)

	// 通知栏常驻设置。枚举值 0和1。ntbarpt=1表示不常驻 为0或不下发此字段表示常驻。
	renderNtbarpt(in)
	// 设置通知栏是否可针对apk下载执行暂停。枚举值 0和1。ntbarpasbl=1表示可暂停，为0或不下发此字段表示不可暂停
	renderNtbarpasbl(in)
	// 表示控制anpk安装完成后的激活控制逻辑。
	//2或不下发此字段表示安装完成后无额外处理
	//1控制表示需要在检测到用户已经安装完成后，弹窗提示激活。弹窗提示仅在当前广告任务所属广告处于展示阶段下进行，广告内容移除后不进行提示。
	//0控制表示检测到已安装完成后，自动触发激活
	renderAtatType(in)

	// mapping idfa
	renderMappingIdfa(in)

	// mapping idfa给idfa赋值切量
	renderMappingIdfaCoverIdfaABTest(in)

	// 判断是否为低配设备
	isLowModel(in)

	parseNums(in)

	in.Param.ServerIP = req_context.GetInstance().ServerIP()
	if in.IsHBRequest {
		in.Param.ServerIP = hbreqctx.GetInstance().ServerIp
	}
	if in.Param.Platform == mvconst.PlatformIOS {
		in.Param.Brand = "apple"
	}

	// render xxxSize and template
	in.Param.Template = renderTemplate(in)
	renderImageSize(in)

	if in.Param.RequestType == mvconst.REQUEST_TYPE_OPENAPI_AD {
		// offset多样性
		renderOffset(in)
	}

	// sdk iv,rv offset 置为0
	adnConf, _ := extractor.GetADNET_SWITCHS()
	closeIVRVOffsetPoison := 0
	if closeOffsetPoison, ok := adnConf["closeOffsetPoison"]; ok {
		closeIVRVOffsetPoison = closeOffsetPoison
	}
	if closeIVRVOffsetPoison != 1 && in.Param.RequestType == mvconst.REQUEST_TYPE_OPENAPI_V3 &&
		(in.Param.AdType == mvconst.ADTypeRewardVideo || in.Param.AdType == mvconst.ADTypeInterstitialVideo) {
		in.Param.Offset = int32(0)
	}

	// new url
	in.Param.IsNewUrl = false
	if in.Param.ApiVersion >= mvconst.API_VERSION_1_4 {
		in.Param.IsNewUrl = true
	}

	renderSupportTrackingTemplate(in)

	if in.Param.RequestPath == mvconst.PATHOpenApiV2 {
		in.Param.RequestType = mvconst.REQUEST_TYPE_OPENAPI
		if in.Param.AppID == 18349 || in.Param.AppID == 18976 || in.Param.AppID == 19154 {
			in.Param.OnlyImpression = 1
		}
	}
	// 处理自有id
	in.Param.ExtSysId = in.Param.SysId + "," + in.Param.BkupId

	// jssdk独立域名，midway的chet，mtrack，iconUrl
	renderJssdkDomain(in)

	// 是否bid请求
	in.IsBidRequest = prf.isBidRequest(in.Param.Scenario, in.Param.AdType)
	if !in.IsHBRequest && in.IsBidRequest {
		// get reduce fill config
		renderReduceFillConfig(in)
		// 格式化2位小数
		in.Param.BidFloor = mvutil.NumFormat(renderEcpmFloor(in), 2)
	}

	// 拼装MCC MNC
	in.Param.MCCMNC = in.Param.MCC + in.Param.MNC
	// 拼装idfv+openidfa
	if in.Param.Platform == mvconst.PlatformIOS && in.Param.FormatSDKVersion.SDKVersionCode < mvconst.IOSSupportTransferIDFVOpenIDFAInParamC &&
		(len(in.Param.IDFV) > 0 || len(in.Param.OpenIDFA) > 0) {
		in.Param.IDFVOpenIDFA = in.Param.IDFV + "," + in.Param.OpenIDFA
	}

	renderExcludePackageName(in)

	// 根据cc + pkg控制不召回某些包名。
	renderCountryExcludePackage(in)
	// vcn abtest
	renderVcnABTest(in)

	// 处理display_cids
	renderDisplayCamIds(in)

	// 算法 offer retarget 实验
	renderAlgoExperiment(in)

	// dmp abtest框架切量标记分配
	renderDmpTag(in)

	// 整理传给as的package_name
	renderPackageName(in)

	// 针对安卓开屏广告下毒，对于没有传only_impression及ping mode的情况。默认都传1
	androidSplashPoison(in)

	// ios sdk 5.9.0版本及以上，6.1.3以下版本，os_version在大于等于9.0.0,小于10.0.0情况下，sdk有个bug 会导致展示storekit后，关闭sk，用户无法操作。
	renderIosStorekitPoison(in)

	if in.IsHBRequest {
		in.Param.IsBlockByImpCap = false
	}

	parseTki(in)

	// H265 视频素材压缩切量
	renderH265ABTest(in)

	// 处理Skadnetwork 信息
	// tag 决定大小写
	// 根据sdk_version和os_version决定ver
	renderSkadnetwork(in)

	// 对于新版本，adnet 对于没有sysid的情况生成sysid
	NewSysId(in)

	// 解析TrafficInfo
	RenderTrafficInfo(in)

	// more_offer&appwall 迁移pioneer切量
	renderMoreOfferAndAppwallMoveToPioneerABTest(in)

	renderMpToPioneerABTest(in)

	renderHBOfferBidPriceABTest(in)
	// 生成more_offer的requestid
	if in.UnitInfo.Unit.MofUnitId > 0 {
		in.Param.MoreOfferRequestId = mvutil.GetGoTkClickID()
	}

	renderParentId(in)
	// 切量到pioneer 的more_offer提前解析mofdata
	renderMofData(in)

	renderLoadDomainABTest(in)

	isMatch := checkBidDeviceGEOCountry(in)
	if !isMatch {
		// 通过配置控制要过滤哪些聚合和开发者
		if filterBidByDeviceGEOCountry(in) {
			return in, errors.New("device geo country is not match this request")
		}
	}

	// native h5 unit size 兜底
	if in.Param.AdType == mvconst.ADTypeNativeH5 && len(in.Param.UnitSize) == 0 {
		in.Param.UnitSize = "320x250"
	}

	// 获取算法出价调控类型
	renderHBPriceFactorSubsidyType(in)

	return in, nil
}

func renderCdnTrackingDomainAbtest(r *mvutil.RequestParams) {
	// 不支持cn流量
	if r.Param.CountryCode == "CN" {
		return
	}
	// 获取配置
	confs := extractor.GetCdnTrackingDomainABTestConf()
	if len(confs) == 0 {
		return
	}
	var rateMap map[string]int
	for _, conf := range confs {
		if conf == nil {
			continue
		}
		rateMap = GetCdnTrackingDomainRateMap(r, conf)
		if len(rateMap) > 0 {
			break
		}
	}
	if len(rateMap) == 0 {
		return
	}
	// 没有对应需要替换的域名，则不做abtest
	cdnTrackingDomain := extractor.GetCdnTrackingDomain(r.Param.Domain)
	if len(cdnTrackingDomain) == 0 {
		return
	}
	// 设备切量
	res := mvutil.RandByDeviceWithRateMap(rateMap, &r.Param, mvconst.SALT_CDN_TRACKING_DOMAIN)
	if res == mvconst.ABTEST_TEST_GROUP_B {
		r.Param.Domain = cdnTrackingDomain
	}
	// 记录标记
	r.Param.ExtDataInit.CdnTrackingDomainABTestTag = res
}

func GetCdnTrackingDomainRateMap(r *mvutil.RequestParams, conf *mvutil.CdnTrackingDomainABTestConf) (rateMap map[string]int) {
	if len(conf.DomainList) > 0 && !mvutil.InStrArray(r.Param.Domain, conf.DomainList) {
		return
	}
	return conf.TotalRate
}

func replaceAdTrackingByCdnDomain(r *mvutil.RequestParams) {
	// 不支持cn流量
	if r.Param.CountryCode == "CN" {
		return
	}

	// 开关为1，才不切到走cdn域名的逻辑里
	adTrackingdoNotUseCdnDomain := extractor.GetAdnetSwitchConf("adTrackingdoNotUseCdnDomain")
	if adTrackingdoNotUseCdnDomain == 1 {
		return
	}
	// 默认使用走cdn的域名
	// 不在tracking 域名替换列表里就不做替换。
	cdnDomain := extractor.GetCdnTrackingDomain(r.Param.Domain)
	if len(cdnDomain) > 0 {
		r.Param.UseCdnTrackingDomain = 1
		r.Param.TrackingCdnDomain = cdnDomain
	}
}

func renderHBPriceFactorSubsidyType(r *mvutil.RequestParams) {
	// 限制hb流量
	if !r.IsHBRequest {
		return
	}
	HBPriceFactorConf, ifFind := extractor.GetAdxHeaderBiddingPriceFactorConfByUnitIdAndCountryCode(r.UnitInfo.UnitId, r.Param.CountryCode)
	if !ifFind {
		return
	}
	r.Param.ExtDataInit.HBSubsidyType = HBPriceFactorConf.SubsidyType
}

func DecryptHarmonyInfo(r *mvutil.RequestParams) {
	if len(r.Param.EncryptHarmonyInfo) > 0 {
		r.Param.DecryptHarmonyInfo = mvutil.Base64Decode(r.Param.EncryptHarmonyInfo)
	}
}

func renderLoadDomainABTest(r *mvutil.RequestParams) {
	if !r.IsHBRequest || r.Param.RequestType != mvconst.REQUEST_TYPE_OPENAPI_V3 {
		// SDK流量的load才生效
		return
	}

	cfg := extractor.GetLoadDomainABTest()
	if cfg.Rate <= 0 {
		return
	}

	if len(cfg.Region) > 0 && !mvutil.InStrArray(hbreqctx.GetInstance().Region, cfg.Region) {
		return
	}

	if len(cfg.Cloud) > 0 && !mvutil.InStrArray(hbreqctx.GetInstance().Cloud, cfg.Cloud) {
		return
	}

	if len(cfg.Platform) > 0 && !mvutil.InArray(r.Param.Platform, cfg.Platform) {
		return
	}

	if len(cfg.BMediationNames) > 0 && mvutil.InStrArray(r.Param.MediationName, cfg.BMediationNames) {
		return
	}

	if len(cfg.WMediationNames) > 0 && !mvutil.InStrArray(r.Param.MediationName, cfg.WMediationNames) {
		return
	}

	if len(cfg.BCountryCode) > 0 && mvutil.InStrArray(r.Param.CountryCode, cfg.BCountryCode) {
		return
	}

	if len(cfg.WCountryCode) > 0 && !mvutil.InStrArray(r.Param.CountryCode, cfg.WCountryCode) {
		return
	}

	if len(cfg.BPublishIds) > 0 && mvutil.InInt64Arr(r.Param.PublisherID, cfg.BPublishIds) {
		return
	}

	if len(cfg.WPublishIds) > 0 && !mvutil.InInt64Arr(r.Param.PublisherID, cfg.WPublishIds) {
		return
	}

	if mvutil.GetRandByGlobalTagId(&r.Param, mvconst.SALT_WTICK, 100) < cfg.Rate {
		r.Param.ExtDataInit.LoadCDNTag = 1
	} else {
		r.Param.ExtDataInit.LoadCDNTag = 2
	}
}

func renderUnitFixedEcpm(r *mvutil.RequestParams) {
	fixecpmObj, ifFind := extractor.GetSspProfitDistributionRuleByUnitIdAndCountryCode(r.Param.UnitID, r.Param.CountryCode)
	if !ifFind || fixecpmObj == nil {
		return
	}
	if fixecpmObj.Type != mvconst.SspProfitDistributionRuleFixedEcpm {
		return
	}
	r.Param.FixedEcpm = fixecpmObj.FixedEcpm
}

func renderHBOfferBidPriceABTest(r *mvutil.RequestParams) {
	if r.Param.RequestPath != mvconst.PATHBidAds {
		return
	}
	conf := extractor.GetHBOfferBidPriceABTestConf()
	if conf == nil {
		return
	}
	if len(conf.UnitList) > 0 && !mvutil.InInt64Arr(r.Param.UnitID, conf.UnitList) {
		return
	}
	var randVal int
	if conf.RandType == 1 {
		randVal = mvutil.GetRandByGlobalTagId(&r.Param, mvconst.SALT_HB_OFFER_BID_PRICE, 100)
	} else {
		randVal = rand.Intn(100)
	}
	if conf.Rate > randVal {
		r.Param.OnlineApiNeedOfferBidPrice = "1"
	}
	return
}

func renderTrackDomainByCountryCode(r *mvutil.RequestParams) {
	conf := extractor.GetTrackDomainByCountryCodeConf()
	countryConf, ok := conf[r.Param.CountryCode]
	if !ok {
		return
	}
	if countryConf == nil {
		return
	}
	if len(countryConf.CityBlackList) > 0 && mvutil.InInt64Arr(r.Param.CityCode, countryConf.CityBlackList) {
		return
	}
	getTrackDomain(countryConf, r)
}

func getTrackDomain(conf *mvutil.TrackDomainByCountryCodeConf, r *mvutil.RequestParams) {
	domainMap := make(map[string]int)
	confId := make(map[string]int)
	for _, confData := range conf.ConfMap {
		domainMap[confData.Url] = confData.Weight
		confId[confData.Url] = confData.Id
	}
	var trackDomain string
	// randtype为1则表示为设备切量
	if conf.RandType == 1 {
		trackDomain = mvutil.RandByDeviceWithRateMap(domainMap, &r.Param, mvconst.SALT_TRACK_DOMAIN_BY_COUNTRY_CODE)
	} else {
		trackDomain = mvutil.RandByRate3(domainMap)
	}

	// 记录切量标记
	if id, ok := confId[trackDomain]; ok {
		r.Param.ExtDataInit.TrackDomainByCountryCode = id
	}
	if len(trackDomain) > 0 {
		r.Param.Domain = trackDomain
	}
}

func renderParentId(r *mvutil.RequestParams) {
	// more_offer使用传过来的值
	if r.Param.AdType == mvconst.ADTypeMoreOffer {
		return
	}
	// 限制有more_offer unitid 才能传ParentId
	// 有more_offer的unitid 的主unit请求，才能走more_offer cache逻辑
	if r.UnitInfo.Unit.MofUnitId > 0 {
		r.Param.ParentId = r.Param.RequestID
	}
}

func renderCNTrackDomainABTest(r *mvutil.RequestParams) {
	// 走了 cn 集群切量
	if r.Param.ExtDataInit.TKCNABTestTag == 1 {
		return
	}
	// 仅支持CN流量
	if r.Param.CountryCode != "CN" {
		return
	}
	conf := extractor.GetCNTrackingDomainConf()
	if conf == nil {
		return
	}
	// 黑名单过滤
	if mvutil.InInt64Arr(r.Param.UnitID, conf.UnitBList) {
		return
	}
	if mvutil.InInt64Arr(r.Param.AppID, conf.AppBList) {
		return
	}
	if mvutil.InInt64Arr(r.Param.PublisherID, conf.PubBList) {
		return
	}
	if len(conf.Conf) == 0 {
		return
	}
	if mvutil.InInt64Arr(r.Param.UnitID, conf.UnitWList) || conf.Status {
		// 分配tracking域名
		getCNTrackDomain(conf.Conf, r)
	}
}

func getCNTrackDomain(conf []*smodel.CdnSetting, r *mvutil.RequestParams) {
	cdnMap := make(map[string]int)
	cdnId := make(map[string]int)
	for _, confData := range conf {
		cdnMap[confData.Url] = confData.Weight
		cdnId[confData.Url] = confData.Id
	}
	// randtype为1则表示为设备切量
	trackDomain := mvutil.RandByRate3(cdnMap)

	if id, ok := cdnId[trackDomain]; ok {
		r.Param.ExtDataInit.CNTrackDomain = id
	}
	// 对于没有切到中国专线的情况，则不改变原有的tracking域名
	if trackDomain != "oriDomain" && len(trackDomain) > 0 {
		r.Param.Domain = trackDomain
	}
}

func renderMofData(r *mvutil.RequestParams) {
	// 切量到pioneer流量，才需要提前解析
	// 现在都切到pioneer了
	if len(r.Param.MofData) == 0 {
		return
	}
	var mofData mvutil.MofData
	// mof_data进行urldecode处理
	MofDataStr, _ := url.QueryUnescape(r.Param.MofData)

	// sdk 动态view逻辑，会把mof_data做base64处理，因此当解析不出来的时候，尝试base64解密后再unmarshal
	if !strings.HasPrefix(MofDataStr, "{") {
		MofDataStr = mvutil.OriBase64Decode(MofDataStr)
	}

	err := jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal([]byte(MofDataStr), &mofData)
	if err == nil {

		// 即无ifmd5也无vfmd5则判定主offer为三方广告主返回的单子
		if len(mofData.Ifmd5) == 0 && len(mofData.Vfmd5) == 0 {
			r.Param.IsThirdPartyMoreoffer = 1
		}

		// 记录主offer的request_id
		if len(mofData.CrtRid) > 0 {
			// tracking日志也需记录
			r.Param.ExtDataInit.CrtRid = mofData.CrtRid
		}
	}
}

func renderMoreOfferAndAppwallMoveToPioneerABTest(r *mvutil.RequestParams) {

	if !mvutil.IsAppwallOrMoreOffer(r.Param.AdType) {
		return
	}
	conf := extractor.GetMoreOfferAndAppwallMoveToPioneerABTestConf()
	adTypeCodeStr := strconv.Itoa(int(r.Param.AdType))
	adTypeConf, ok := conf[adTypeCodeStr]
	if !ok {
		return
	}
	// unit
	if len(adTypeConf.UnitConf) > 0 {
		unitConf, ok := adTypeConf.UnitConf[strconv.FormatInt(r.Param.UnitID, 10)]
		if ok {
			renderMoreOfferAndAppwallMoveToPioneerTag(r, unitConf, adTypeConf.RandType)
			return
		}
	}

	// app
	if len(adTypeConf.AppConf) > 0 {
		appConf, ok := adTypeConf.AppConf[strconv.FormatInt(r.Param.AppID, 10)]
		if ok {
			renderMoreOfferAndAppwallMoveToPioneerTag(r, appConf, adTypeConf.RandType)
			return
		}
	}

	// pub
	if len(adTypeConf.PublisherConf) > 0 {
		pubConf, ok := adTypeConf.PublisherConf[strconv.FormatInt(r.Param.PublisherID, 10)]
		if ok {
			renderMoreOfferAndAppwallMoveToPioneerTag(r, pubConf, adTypeConf.RandType)
			return
		}
	}

	// 整体切量
	if len(adTypeConf.TotalRate) > 0 {
		renderMoreOfferAndAppwallMoveToPioneerTag(r, adTypeConf.TotalRate, adTypeConf.RandType)
	}
}

func renderMoreOfferAndAppwallMoveToPioneerTag(r *mvutil.RequestParams, rate map[string]int, randType int) {
	// 默认请求切量，1则为设备切量
	res := "a0"
	if randType == 1 {
		res = mvutil.RandByDeviceWithRateMap(rate, &r.Param, mvconst.SALT_MORE_OFFER_MV_TO_PIONEER)
	} else {
		res = mvutil.RandByRate3(rate)
	}
	r.Param.ExtDataInit.MoreofferAndAppwallMvToPioneerTag = res
}

func renderMappingIdfaCoverIdfaABTest(r *mvutil.RequestParams) {
	// 限制sdk流量以及没有idfa部分
	if !mvutil.IsHbOrV3OrV5Request(r.Param.RequestPath) || !demand.IsEmptyIDFA(r.Param.IDFA) || r.Param.PlatformName != mvconst.PlatformNameIOS {
		return
	}
	// 要有mapping idfa
	if len(r.Param.MappingIdfa) == 0 {
		return
	}
	mappingIdfaCoverIdfaABTestConf := extractor.GetMappingIdfaCoverIdfaABTestConf()
	if !isAllowMappingIdfaTest(r, mappingIdfaCoverIdfaABTestConf) {
		return
	}
	var res string
	if mappingIdfaCoverIdfaABTestConf.RandType == 1 {
		res = mvutil.RandByRate3(mappingIdfaCoverIdfaABTestConf.ConfMap)
	} else {
		res = mvutil.RandByMappingIdfa(mappingIdfaCoverIdfaABTestConf.ConfMap, r.Param.MappingIdfa)
	}
	r.Param.ExtDataInit.MappingIdfaCoverIdfaTag = res

	// key为b则表示实验组
	if res == "b" {
		r.Param.IDFA = r.Param.MappingIdfa
	}
}

func renderMappingIdfa(r *mvutil.RequestParams) {
	// 限制sdk流量以及没有idfa部分
	if !mvutil.IsHbOrV3OrV5Request(r.Param.RequestPath) || !demand.IsEmptyIDFA(r.Param.IDFA) || r.Param.PlatformName != mvconst.PlatformNameIOS {
		return
	}
	if len(r.Param.SysId) == 0 && len(r.Param.IDFV) == 0 {
		return
	}
	mappingIdfaConf := extractor.GetMappingIdfaConf()
	if !isAllowMappingIdfaTest(r, mappingIdfaConf) {
		return
	}
	// 表示有命中实验条件，但没查询到idfa
	originTag := "a0"
	metrics.IncCounterWithLabelValues(15)
	// 使用sysid查询
	var mappingIdfa string
	if len(r.Param.SysId) > 0 {
		metrics.IncCounterWithLabelValues(17)
		devKey := "sysid_" + strings.ToLower(r.Param.SysId)
		mappingIdfa = getMappingIdfaFromAerospike(r, devKey)
		if len(mappingIdfa) > 0 {
			metrics.IncCounterWithLabelValues(18)
		}
	}
	// 使用idfv查询
	if len(mappingIdfa) == 0 && len(r.Param.IDFV) > 0 {
		metrics.IncCounterWithLabelValues(19)
		devKey := "idfv_" + strings.ToLower(r.Param.IDFV)
		mappingIdfa = getMappingIdfaFromAerospike(r, devKey)
		if len(mappingIdfa) > 0 {
			metrics.IncCounterWithLabelValues(20)
		}
	}
	r.Param.ExtDataInit.MappingIdfaTag = originTag
	if len(mappingIdfa) == 0 {
		return
	}
	// 有idfa计数
	metrics.IncCounterWithLabelValues(16)
	var res string
	if mappingIdfaConf.RandType == 1 {
		res = mvutil.RandByRate3(mappingIdfaConf.ConfMap)
	} else {
		res = mvutil.RandByMappingIdfa(mappingIdfaConf.ConfMap, mappingIdfa)
	}

	r.Param.ExtDataInit.MappingIdfaTag = res

	// key为b则表示实验组
	if res == "b" {
		r.Param.MappingIdfa = mappingIdfa
	}
}

func getMappingIdfaFromAerospike(r *mvutil.RequestParams, devKey string) string {
	value, err := mkv.GetMKVFieldVal(devKey)
	if err != nil {
		mvutil.Logger.AerospikeLog.Errorf("get aerospike data error.error=[%s],requestid=[%s],key=[%s]", err.Error(), r.Param.RequestID, devKey)
		return ""
	}
	if val, ok := value["idfa"]; ok {
		return mvutil.GetIdfaString(string(val))
	}
	return ""
}

func isAllowMappingIdfaTest(r *mvutil.RequestParams, conf *mvutil.MappingIdfaAbtestConf) bool {
	if conf == nil {
		return false
	}
	if len(conf.AppBlackList) > 0 && mvutil.InInt64Arr(r.Param.AppID, conf.AppBlackList) {
		return false
	}
	if len(conf.PubBlackList) > 0 && mvutil.InInt64Arr(r.Param.PublisherID, conf.PubBlackList) {
		return false
	}
	if len(conf.CountryCodeBlackList) > 0 && mvutil.InStrArray(r.Param.CountryCode, conf.CountryCodeBlackList) {
		return false
	}
	if len(conf.AppList) > 0 && !mvutil.InInt64Arr(r.Param.AppID, conf.AppList) {
		return false
	}
	if len(conf.PubList) > 0 && !mvutil.InInt64Arr(r.Param.PublisherID, conf.PubList) {
		return false
	}
	if len(conf.CountryCodeList) > 0 && !mvutil.InStrArray(r.Param.CountryCode, conf.CountryCodeList) {
		return false
	}
	if len(conf.ConfMap) == 0 {
		return false
	}
	return true
}

func RenderTrafficInfo(r *mvutil.RequestParams) {
	if len(r.Param.TrafficInfo) == 0 {
		return
	}
	trafficInfoStr, err := mvutil.Decrypt(r.Param.TrafficInfo)
	r.Param.DecryptTrafficInfoStr = trafficInfoStr
	if err != nil {
		mvutil.Logger.Runtime.Errorf("str=[%s] decrypt traffic info error. error:%s", r.Param.TrafficInfo, err.Error())
	}
	var trafficInfo TrafficInfo
	err = jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal([]byte(trafficInfoStr), &trafficInfo)
	if err == nil {
		r.Param.StackList = trafficInfo.StackListStr
		r.Param.ClassNameList = trafficInfo.ClassNameListStr
		r.Param.ProtocolList = trafficInfo.ProtocolListStr
		return
	}
}

func NewSysId(r *mvutil.RequestParams) {
	if r.Param.ApiVersion < mvconst.API_VERSION_2_2 {
		return
	}
	if !mvutil.IsHbOrV3OrV5Request(r.Param.RequestPath) {
		return
	}
	if len(r.Param.SysId) > 0 {
		return
	}
	// 需要有设备id才能生成sysid
	devId := mvutil.GetSysIdDevId(&r.Param)
	if len(devId) == 0 {
		return
	}
	r.Param.SysId = uuid.NewV5(uuid.NamespaceDNS, devId).String()
}

func renderDisplayCamIds(r *mvutil.RequestParams) {
	// 优先使用新版本传来的display_info
	if len(r.Param.NewDisplayInfoData) > 0 {
		for _, v := range r.Param.NewDisplayInfoData {
			cid, _ := strconv.ParseInt(v.CampaignId, 10, 64)
			r.Param.DisplayCamIds = append(r.Param.DisplayCamIds, cid)
		}
		return
	}
	if len(r.Param.DisplayInfoData) > 0 {
		for _, v := range r.Param.DisplayInfoData {
			cid, _ := strconv.ParseInt(v.CampaignId, 10, 64)
			r.Param.DisplayCamIds = append(r.Param.DisplayCamIds, cid)
		}
	}
}

func renderSkadnetwork(r *mvutil.RequestParams) {
	if r.Param.Skadnetwork == nil {
		if r.Param.Platform == mvconst.PlatformIOS && r.Param.OSVersionCode >= 11030000 && r.Param.OSVersionCode < 14000000 {
			r.Param.Skadnetwork = &mvutil.Skadnetwork{}
			r.Param.Skadnetwork.Ver = "1.0"
		} else if r.Param.Platform == mvconst.PlatformIOS && r.Param.OSVersionCode >= 14000000 {
			r.Param.Skadnetwork = &mvutil.Skadnetwork{}
			r.Param.Skadnetwork.Ver = "2.0"
		} else {
			return
		}
	}

	// tag 决定大小写
	if r.Param.Skadnetwork.Tag == "1" {
		r.Param.Skadnetwork.Adnetids = []string{
			mvconst.MTG_SK_NETWORK_ID,
		}
	} else if r.Param.Skadnetwork.Tag == "2" {
		r.Param.Skadnetwork.Adnetids = []string{
			mvconst.LOWER_MTG_SK_NETWORK_ID,
		}
	}
	// process s2s sk id list
	var skIds string
	var appVerInSkIds bool
	for _, appSkId := range r.AppInfo.AppSkIds {
		if appSkId.AppVersion == r.Param.AppVersionName {
			skIds = appSkId.SkIds
			appVerInSkIds = true
			break
		}
	}
	if !appVerInSkIds {
		skIds = r.AppInfo.DefaultSkIds
	}
	if len(skIds) > 0 {
		if r.Param.Skadnetwork.Tag == "1" {
			skIds = strings.Replace(skIds, mvconst.LOWER_MTG_SK_NETWORK_ID, mvconst.MTG_SK_NETWORK_ID, -1)
		} else if r.Param.Skadnetwork.Tag == "2" {
			skIds = strings.Replace(skIds, mvconst.MTG_SK_NETWORK_ID, mvconst.LOWER_MTG_SK_NETWORK_ID, -1)
		}
		ids := strings.Split(skIds, ",")
		r.Param.Skadnetwork.Adnetids = mvutil.StringAppendCategory(ids, r.Param.Skadnetwork.Adnetids)
	}

	// skver在2.1版本及以后由服务端控制判断
	skVerCode, err := strconv.ParseFloat(r.Param.Skadnetwork.Ver, 64)
	if err == nil && skVerCode >= 2.1 {
		// os_version为14.6+则为3.0
		if r.Param.OSVersionCode >= 14060000 {
			r.Param.Skadnetwork.Ver = "3.0"
			// 若为大写的sknetworkid，则只能支持2.2
			mtgSkAdNetworkId := demand.GetMTGSKAdnetworkId(r.Param.Skadnetwork.Adnetids)
			if mtgSkAdNetworkId == demand.MTG_SK_NETWORK_ID {
				r.Param.Skadnetwork.Ver = "2.2"
			}
		}
	}

}

func renderAtatType(r *mvutil.RequestParams) {
	paramAbtestConf := extractor.GetADNET_PARAMS_ABTEST_CONFS()
	if atatTypeConf, ok := paramAbtestConf["atatType"]; ok && atatTypeConf != nil {
		conf, randType := output.GetABTestConf(&r.Param, nil, atatTypeConf, 0)
		if len(conf) == 0 {
			// 没命中实验，则使用app维度配置
			r.Param.AtatType = r.AppInfo.App.AtatType
			return
		}
		// 获取abtest结果
		finalVal, randOk := output.GetABTestRes(conf, randType, &r.Param, "atatType")
		if !randOk {
			// 没命中实验，则使用app维度配置
			r.Param.AtatType = r.AppInfo.App.AtatType
			return
		}
		r.Param.AtatType = finalVal
	} else {
		// 没命中实验，则使用app维度配置
		r.Param.AtatType = r.AppInfo.App.AtatType
	}
}

func renderNtbarpasbl(r *mvutil.RequestParams) {
	paramAbtestConf := extractor.GetADNET_PARAMS_ABTEST_CONFS()
	if ntbarpasblConf, ok := paramAbtestConf["ntbarpasbl"]; ok && ntbarpasblConf != nil {
		conf, randType := output.GetABTestConf(&r.Param, nil, ntbarpasblConf, 0)
		if len(conf) == 0 {
			// 没命中实验，则使用app维度配置
			r.Param.Ntbarpasbl = r.AppInfo.App.Ntbarpasbl
			return
		}
		// 获取abtest结果
		finalVal, randOk := output.GetABTestRes(conf, randType, &r.Param, "ntbarpasbl")
		if !randOk {
			// 没命中实验，则使用app维度配置
			r.Param.Ntbarpasbl = r.AppInfo.App.Ntbarpasbl
			return
		}
		r.Param.Ntbarpasbl = finalVal
	} else {
		// 没命中实验，则使用app维度配置
		r.Param.Ntbarpasbl = r.AppInfo.App.Ntbarpasbl
	}
}

func renderNtbarpt(r *mvutil.RequestParams) {
	paramAbtestConf := extractor.GetADNET_PARAMS_ABTEST_CONFS()
	if ntbarptConf, ok := paramAbtestConf["ntbarpt"]; ok && ntbarptConf != nil {
		conf, randType := output.GetABTestConf(&r.Param, nil, ntbarptConf, 0)
		if len(conf) == 0 {
			// 没命中实验，则使用app维度配置
			r.Param.Ntbarpt = r.AppInfo.App.Ntbarpt
			return
		}
		// 获取abtest结果
		finalVal, randOk := output.GetABTestRes(conf, randType, &r.Param, "ntbarpt")
		if !randOk {
			// 没命中实验，则使用app维度配置
			r.Param.Ntbarpt = r.AppInfo.App.Ntbarpt
			return
		}
		r.Param.Ntbarpt = finalVal
	} else {
		// 没命中实验，则使用app维度配置
		r.Param.Ntbarpt = r.AppInfo.App.Ntbarpt
	}
}

func renderH265ABTest(r *mvutil.RequestParams) {
	// 限制sdk 流量
	if !mvutil.IsHbOrV3OrV5Request(r.Param.RequestPath) {
		return
	}
	// 排除topon 流量
	if r.Param.RequestPath == mvconst.PATHTOPON {
		return
	}
	h265VideoABTestConf := extractor.GetH265_VIDEO_ABTEST_CONF()
	if h265VideoABTestConf == nil {
		return
	}
	// 排除不支持H265的流量
	if r.Param.PlatformName == mvconst.PlatformNameAndroid {
		if r.Param.OSVersionCode < h265VideoABTestConf.MinAllowAndroidOsVersionCode {
			return
		}
		if r.Param.FormatSDKVersion.SDKVersionCode < h265VideoABTestConf.MinAllowAndroidSdkVersionCode {
			return
		}
		if mvutil.InInt32Arr(r.Param.OSVersionCode, h265VideoABTestConf.AndroidBlackOsVersionCodeList) {
			return
		}
		if mvutil.InInt32Arr(r.Param.FormatSDKVersion.SDKVersionCode, h265VideoABTestConf.AndroidBlackSdkVersionCodeList) {
			return
		}
	} else {
		if r.Param.OSVersionCode < h265VideoABTestConf.MinAllowIosOsVersionCode {
			return
		}
		if r.Param.FormatSDKVersion.SDKVersionCode < h265VideoABTestConf.MinAllowIosSdkVersionCode {
			return
		}
		if mvutil.InInt32Arr(r.Param.OSVersionCode, h265VideoABTestConf.IosBlackOsVersionCodeList) {
			return
		}
		if mvutil.InInt32Arr(r.Param.FormatSDKVersion.SDKVersionCode, h265VideoABTestConf.IosBlackSdkVersionCodeList) {
			return
		}
	}

	// model 黑名单
	if mvutil.InStrArray(r.Param.Model, h265VideoABTestConf.BlackModelList) {
		return
	}

	// model+os_version 黑名单
	if mvutil.InStrArray(strings.ToLower(r.Param.Model)+"_"+strconv.Itoa(int(r.Param.OSVersionCode)), h265VideoABTestConf.BlackModelOsVersionCodeList) {
		return
	}

	// brand + model 白名单
	if len(h265VideoABTestConf.BrandModelWhiteList) > 0 && !mvutil.InStrArray(strings.ToLower(r.Param.Brand)+"_"+strings.ToLower(r.Param.Model), h265VideoABTestConf.BrandModelWhiteList) {
		return
	}

	// unit 维度切量
	unitConf, ok := h265VideoABTestConf.UnitConf[strconv.FormatInt(r.Param.UnitID, 10)]
	if ok && len(unitConf) > 0 {
		resKey := mvutil.RandByRate3(unitConf)
		ifSupH265Id, _ := strconv.Atoi(resKey)
		r.Param.IfSupH265 = int32(ifSupH265Id)
		return
	}

	// app 维度切量
	appConf, ok := h265VideoABTestConf.AppConf[strconv.FormatInt(r.Param.AppID, 10)]
	if ok && len(appConf) > 0 {
		resKey := mvutil.RandByRate3(appConf)
		ifSupH265Id, _ := strconv.Atoi(resKey)
		r.Param.IfSupH265 = int32(ifSupH265Id)
		return
	}

	// app 维度切量
	adTypeConf, ok := h265VideoABTestConf.AdTypeConf[mvutil.GetAdTypeStr(r.Param.AdType)]
	if ok && len(adTypeConf) > 0 {
		resKey := mvutil.RandByRate3(adTypeConf)
		ifSupH265Id, _ := strconv.Atoi(resKey)
		r.Param.IfSupH265 = int32(ifSupH265Id)
		return
	}

	// 整体 维度切量
	if len(h265VideoABTestConf.TotalRate) > 0 {
		resKey := mvutil.RandByRate3(h265VideoABTestConf.TotalRate)
		ifSupH265Id, _ := strconv.Atoi(resKey)
		r.Param.IfSupH265 = int32(ifSupH265Id)
		return
	}
}

func parseTki(r *mvutil.RequestParams) {
	// 优先取新tki参数
	if len(r.Param.NewTKI) > 0 {
		decryptTki, err := mvutil.Decrypt(r.Param.NewTKI)
		if err != nil {
			mvutil.Logger.Runtime.Errorf("str=[%s] decrypt tki error. error:%s", r.Param.NewTKI, err.Error())
		}
		var newTki NewTki
		err = jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal([]byte(decryptTki), &newTki)
		if err == nil {
			r.Param.OsvUpTime = newTki.OsVersionUpdateTime
			r.Param.Ram = newTki.Ram
			r.Param.UpTime = newTki.UpdateTime
			r.Param.Carrier = newTki.Carrier
			// 新版本中，不会有newid和oldid了
			return
		}
	}

	var tki Tki
	err := jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal([]byte(r.Param.TKI), &tki)
	if err == nil {
		r.Param.OsvUpTime = tki.OsVersionUpdateTime
		r.Param.Ram = tki.Ram
		r.Param.UpTime = tki.UpdateTime
		r.Param.NewId = tki.NewId
		r.Param.OldId = tki.OldId
		r.Param.Carrier = tki.Carrier
	}
}

func renderIfSupDco(r *mvutil.RequestParams) {
	dcoTestConf := extractor.GetDcoTestConf()
	rate := getTestConfRate(r, dcoTestConf)
	randVal := rand.Intn(100)
	if rate > randVal {
		r.Param.IfSupDco = 1
		r.Param.ExtDataInit.IfSupDco = 1
	}
}

func renderIosStorekitPoison(r *mvutil.RequestParams) {
	// 限制ios sdk流量
	if !mvutil.IsHbOrV3OrV5Request(r.Param.RequestPath) {
		return
	}
	if r.Param.Platform != mvconst.PlatformIOS {
		return
	}
	adnConf, _ := extractor.GetADNET_SWITCHS()
	if sdkCloseSkPoison, ok := adnConf["sdkCloseSkPoison"]; ok && sdkCloseSkPoison == 1 {
		return
	}
	skpsDefaultMinOs := 9000000
	skpsdefaultMaxOs := 10000000
	if skpsDefMinOs, ok := adnConf["skpsDefMinOs"]; ok {
		skpsDefaultMinOs = skpsDefMinOs
	}
	if skpsDefMaxOs, ok := adnConf["skpsDefMaxOs"]; ok {
		skpsdefaultMaxOs = skpsDefMaxOs
	}
	// 5.9.0版本及以上，6.1.3以下版本，os_version在大于等于9.0.0,小于10.0.0情况下需要下毒
	if (r.Param.FormatSDKVersion.SDKVersionCode >= 50900 && r.Param.FormatSDKVersion.SDKVersionCode < 60103) &&
		(r.Param.OSVersionCode >= int32(skpsDefaultMinOs) && r.Param.OSVersionCode < int32(skpsdefaultMaxOs)) {
		r.Param.IosStorekitPoisonFlag = true
	}
}

func renderSupportTrackingTemplate(r *mvutil.RequestParams) {
	// 取消掉大模版限制的开关
	adnConf, _ := extractor.GetADNET_SWITCHS()
	if supportTrackingTemplateBigTemplateSwitch, ok := adnConf["supportTrackingTemplateBigTemplateSwitch"]; ok && supportTrackingTemplateBigTemplateSwitch == 1 {
		if r.Param.BigTemplateFlag {
			return
		}
	}

	if !r.Param.IsNewUrl {
		return
	}

	if r.Param.Scenario != mvconst.SCENARIO_OPENAPI || r.Param.RequestType != mvconst.REQUEST_TYPE_OPENAPI_V3 {
		return
	}

	cfg := extractor.GetSupportTrackingTemplateConf()
	if cfg == nil || !cfg.Status {
		return
	}

	if len(cfg.BAppIds) > 0 && mvutil.InInt64Arr(r.Param.AppID, cfg.BAppIds) {
		return
	}

	if len(cfg.BSDKVersion) > 0 && mvutil.InStrArray(r.Param.SDKVersion, cfg.BSDKVersion) {
		return
	}

	if len(cfg.BAPIVersion) > 0 && mvutil.InInt32Arr(r.Param.ApiVersionCode, cfg.BAPIVersion) {
		return
	}
	r.Param.SupportTrackingTemplate = true
}

func androidSplashPoison(r *mvutil.RequestParams) {
	if r.Param.Platform == mvconst.PlatformAndroid && r.Param.AdType == mvconst.ADTypeSplash {
		if r.Param.OnlyImpression == 0 {
			r.Param.OnlyImpression = 1
		}
		if r.Param.PingMode == 0 {
			r.Param.PingMode = 1
		}
	}
}

func renderBigTemplateFlag(r *mvutil.RequestParams) {
	// 限制sdk 流量
	if !mvutil.IsHbOrV3OrV5Request(r.Param.RequestPath) || r.Param.RequestType != mvconst.REQUEST_TYPE_OPENAPI_V3 {
		return
	}
	adnConf, _ := extractor.GetADNET_SWITCHS()
	if filterBigTplApiVersion, ok := adnConf["filterBigTplApiVersion"]; ok && filterBigTplApiVersion == 1 {
		if r.Param.ApiVersion < mvconst.API_VERSION_1_9 {
			return
		}
	}

	// topon 聚合平台的ios的693、694不支持大模版。
	// 背景：ios sdk的693、694版本sdk  在topon 渠道的流量上，出现 大模板（二选一）的广告无法关闭的问题。
	if r.Param.PlatformName == mvconst.PlatformNameIOS && r.Param.Extchannel == constant.ToponMediationIdStr &&
		(r.Param.FormatSDKVersion.SDKVersionCode == mvconst.IosSupportOfferRewardPlusVersion ||
			r.Param.FormatSDKVersion.SDKVersionCode == mvconst.ToponUnSupportBigTempalteVersion) {
		return
	}

	bigTemplateConf := extractor.GetBIG_TEMPLATE_CONF()
	rate := getTestConfRate(r, bigTemplateConf)
	randVal := rand.Intn(100)
	if rate > randVal {
		r.Param.BigTemplateFlag = true
	}
}

func renderPolarisFlag(r *mvutil.RequestParams) {
	// 限制sdk 流量
	if !mvutil.IsHbOrV3OrV5Request(r.Param.RequestPath) || r.Param.RequestType != mvconst.REQUEST_TYPE_OPENAPI_V3 {
		return
	}
	// sdk banner默认支持polaris
	if r.Param.AdType == mvconst.ADTypeSdkBanner || r.Param.AdType == mvconst.ADTypeNativeH5 {
		r.Param.PolarisFlag = true
	}
	polarisConf := extractor.GetPolarisFlagConf()
	rate := getTestConfRate(r, polarisConf)
	randVal := rand.Intn(100)
	if rate > randVal {
		r.Param.PolarisFlag = true
	}
}

func getTestConfRate(r *mvutil.RequestParams, conf *mvutil.TestConf) int {
	rate := 0
	if conf == nil {
		return rate
	}
	// 限制adtype
	if !mvutil.InInt32Arr(r.Param.AdType, conf.SupAdType) {
		return rate
	}
	// 黑名单内不做切量
	if mvutil.InInt64Arr(r.Param.UnitID, conf.UnitBList) {
		return rate
	}
	if mvutil.InInt64Arr(r.Param.AppID, conf.AppBList) {
		return rate
	}
	if mvutil.InInt64Arr(r.Param.PublisherID, conf.PubBList) {
		return rate
	}
	// 获取unit，app，pub维度配置
	unitStr := strconv.FormatInt(r.Param.UnitID, 10)
	if rate, ok := conf.UnitWList[unitStr]; ok {
		return rate
	}
	appStr := strconv.FormatInt(r.Param.AppID, 10)
	if rate, ok := conf.AppWList[appStr]; ok {
		return rate
	}
	pubStr := strconv.FormatInt(r.Param.PublisherID, 10)
	if rate, ok := conf.PubWList[pubStr]; ok {
		return rate
	}
	return conf.TotalRate
}

func renderPlacementId(r *mvutil.RequestParams) {
	// 如果unit 上有配置，则使用unit上配置的placementid。如果unit配置为空值，则使用sdk传入的值。
	if r.UnitInfo.Unit.PlacementId > 0 {
		r.Param.FinalPlacementId = r.UnitInfo.Unit.PlacementId
	} else {
		r.Param.FinalPlacementId = r.Param.PlacementId
	}
	// 记录到日志中
	r.Param.ExtPlacementId = strconv.FormatInt(r.Param.FinalPlacementId, 10)
}

func renderPackageName(r *mvutil.RequestParams) {
	if r.Param.DspMof == 1 {
		r.Param.AppPackageName = r.Param.PackageName
		return
	}
	// 对于ios，原本sdk传递的package_name为bundle id，现需改为pkg
	if r.Param.Platform == mvconst.PlatformIOS {
		r.Param.AppPackageName = r.AppInfo.RealPackageName
		return
	}
	r.Param.AppPackageName = r.Param.PackageName

}

func renderDmpTag(r *mvutil.RequestParams) {
	if len(r.Param.GAID) > 0 {
		r.Param.ExtDataInit.GaidTag = mvutil.GetRandByStrAddSalt(r.Param.GAID, mvconst.RandSum128, mvconst.SALT_DMP_ABTEST)
	}
	if len(r.Param.IDFA) > 0 {
		r.Param.ExtDataInit.IdfaTag = mvutil.GetRandByStrAddSalt(r.Param.IDFA, mvconst.RandSum128, mvconst.SALT_DMP_ABTEST)
	}
	if len(r.Param.IMEI) > 0 {
		r.Param.ExtDataInit.ImeiTag = mvutil.GetRandByStrAddSalt(r.Param.IMEI, mvconst.RandSum128, mvconst.SALT_DMP_ABTEST)
	}
	if len(r.Param.AndroidID) > 0 {
		r.Param.ExtDataInit.AndroidIdTag = mvutil.GetRandByStrAddSalt(r.Param.AndroidID, mvconst.RandSum128, mvconst.SALT_DMP_ABTEST)
	}
	if len(r.Param.ImeiMd5) > 0 {
		r.Param.ExtDataInit.ImeiMd5Tag = mvutil.GetRandByStrAddSalt(r.Param.ImeiMd5, mvconst.RandSum128, mvconst.SALT_DMP_ABTEST)
	}
}

func renderCountryCode(r *mvutil.RequestParams) {
	// 如果传了CC 覆盖countryCode
	if len(r.Param.CC) > 0 {
		if r.Param.CC == "GB" {
			r.Param.CC = "UK"
		}
		adnConf, _ := extractor.GetADNET_SWITCHS()
		var CCButton int
		if ccButtonConf, ok := adnConf["ccButton"]; ok {
			CCButton = ccButtonConf
		}
		// 对于online api，若已传了client_ip，若再传cc，则以client_ip为准
		// 若client_ip解析不出cc，则使用参数cc
		if r.Param.RequestType != mvconst.REQUEST_TYPE_OPENAPI_AD || len(r.Param.ParamCIP) == 0 || r.Param.CountryCode == "**" || CCButton == 1 {
			r.Param.CountryCode = r.Param.CC
		}
	}
}

func renderAlgoExperiment(r *mvutil.RequestParams) {
	r.Param.PubFlowExpectPrice = extractor.GetUnitFixedEcpm(r.Param.UnitID, r.Param.CountryCode)
	if !r.IsHBRequest {
		if fillEcpmFloor, err := getFillEcpmFloor(r); err == nil {
			r.Param.FillEcpmFloor = fillEcpmFloor
		}
	}
}

func renderPlayableFlag(r *mvutil.RequestParams) {
	r.Param.PlayableFlag = false
	if mvutil.IsHbOrV3OrV5Request(r.Param.RequestPath) {
		r.Param.PlayableFlag = true
	}
}

func handleScreensize(r *mvutil.RequestParams) {
	screenSize := r.Param.ScreenSize
	arr := strings.Split(screenSize, "x")
	if len(arr) != 2 {
		return
	}
	width, _ := strconv.ParseInt(arr[0], 10, 64)
	heigh, _ := strconv.ParseInt(arr[1], 10, 64)
	max := ""
	min := ""
	if width >= heigh {
		max = arr[0]
		min = arr[1]
	} else {
		max = arr[1]
		min = arr[0]
	}
	if r.Param.FormatOrientation == mvconst.ORIENTATION_LANDSCAPE {
		screenSize = max + "x" + min
		r.Param.ScreenWidth, _ = strconv.Atoi(max)
		r.Param.ScreenHeigh, _ = strconv.Atoi(min)
	} else if r.Param.FormatOrientation == mvconst.ORIENTATION_PORTRAIT {
		screenSize = min + "x" + max
		r.Param.ScreenWidth, _ = strconv.Atoi(min)
		r.Param.ScreenHeigh, _ = strconv.Atoi(max)
	}
	r.Param.ScreenSize = screenSize
	// 对于native h5的unitsize，使用sdk传来的值取值
	if r.Param.AdType != mvconst.ADTypeNativeH5 {
		// unitSize
		r.Param.UnitSize = r.Param.ScreenSize
	}
}

func renderImageSize(r *mvutil.RequestParams) {
	r.Param.ImageSize = mvconst.GetImageSizeByID(r.Param.ImageSizeID)
	if r.UnitInfo.UnitId <= 0 {
		return
	}
	if r.Param.ImageSizeID > 0 {
		return
	}
	sAdType := strconv.Itoa(int(r.Param.AdType))
	sOrientation := strconv.Itoa(r.Param.FormatOrientation)
	sTemplate := strconv.Itoa(r.Param.Template)
	kStr := sAdType + "_" + sOrientation + "_" + sTemplate
	confs, ifFind := extractor.GetTemplate()
	if !ifFind {
		return
	}
	conf, ok := confs[kStr]
	if !ok {
		return
	}
	r.Param.ImageSize = conf.ImageSize
	r.Param.ImageSizeID = conf.ImageSizeID
	if len(conf.UnitSize) > 0 {
		r.Param.UnitSize = conf.UnitSize
	}
}

func renderTemplate(r *mvutil.RequestParams) int {
	if len(r.UnitInfo.Unit.Templates) <= 0 {
		return mvconst.TemplateDefault
	}
	if r.Param.AdType == mvconst.ADTypeBanner {
		if r.Param.ScreenWidth > 320 && mvutil.InArray(mvconst.TemplateMultiElements, r.UnitInfo.Unit.Templates) {
			return mvconst.TemplateMultiElements
		}
		return mvutil.RandIntArr(r.UnitInfo.Unit.Templates)
	}
	if r.Param.AdType == mvconst.ADTypeNative {
		template := mvconst.TemplateMultiElements
		for _, v := range r.Param.NativeInfoList {
			template = v.AdTemplate
			// r.Param.AdNum = int32(v.RequireNum)
			if v.AdTemplate == mvconst.TemplateMultiIcons {
				break
			}
		}
		return template
	}
	return mvutil.RandIntArr(r.UnitInfo.Unit.Templates)
}

// 若unit没配或者为both，则用流量侧传的orientation
func handleOrientation(r *mvutil.RequestParams) {

	r.Param.FormatOrientation = r.UnitInfo.Unit.Orientation
	if r.Param.FormatOrientation == mvconst.ORIENTATION_BOTH {
		orientation := r.Param.Orientation
		if orientation == mvconst.ORIENTATION_PORTRAIT || orientation == mvconst.ORIENTATION_LANDSCAPE {
			r.Param.FormatOrientation = orientation
		}
	}
	// hb 需要处理 unit 和 request 都为 ORIENTATION_BOTH 的情况
	if mvutil.IsHBS2SUseVideoOrientation(r) && r.Param.FormatOrientation == mvconst.ORIENTATION_BOTH {
		if r.Param.VideoW < r.Param.VideoH {
			r.Param.FormatOrientation = mvconst.ORIENTATION_PORTRAIT
		} else {
			r.Param.FormatOrientation = mvconst.ORIENTATION_LANDSCAPE
		}
	}
	// both 默认横屏
	if r.Param.PlayableFlag && r.Param.FormatOrientation == mvconst.ORIENTATION_BOTH {
		// 开屏广告默认竖屏
		if r.Param.AdType == mvconst.ADTypeSplash {
			r.Param.FormatOrientation = mvconst.ORIENTATION_PORTRAIT
		} else {
			r.Param.FormatOrientation = mvconst.ORIENTATION_LANDSCAPE
		}
	}

	// Orientation 兼容逻辑
	if !mvutil.InArray(r.Param.Orientation, []int{mvconst.ORIENTATION_BOTH, mvconst.ORIENTATION_LANDSCAPE, mvconst.ORIENTATION_PORTRAIT}) {
		// 过滤 Orientation 非法的情况 (非: 0, 1, 2), 统一置为0
		mvutil.Logger.Runtime.Warnf("wrong orientation:[%d] hit with bidId:[%s] appId:[%d] unitId:[%d]", r.Param.Orientation, r.Param.RequestID, r.Param.AppID, r.Param.UnitID)
		r.Param.Orientation = mvconst.ORIENTATION_BOTH
	}
}

// 判断是否为低配设备
func isLowModel(r *mvutil.RequestParams) {
	r.Param.LowDevice = false
	model := r.Param.Model
	if !strings.HasPrefix(model, "iphone") || len(model) <= 6 || strings.Index(model, ",") <= 6 {
		return
	}
	iphoneNum, err := strconv.Atoi(model[6:strings.Index(model, ",")])
	if err != nil {
		return
	}
	if iphoneNum <= 7 {
		r.Param.LowDevice = true
	}
}

func renderNewCreativeFlag(r *mvutil.RequestParams) {
	r.Param.NewCreativeFlag = false
	// 只针对mv sdk 切素材三期
	if mvutil.IsHbOrV3OrV5Request(r.Param.RequestPath) {
		r.Param.NewCreativeFlag = true
	}
	return
}

func renderJssdkDomain(r *mvutil.RequestParams) {
	if mvutil.NeedNewJssdkDomain(r.Param.RequestPath, r.Param.Ndm) {
		jssdkDomainConf, _ := extractor.GetJSSDK_DOMAIN()
		if chetDomain, ok := jssdkDomainConf["chet"]; ok && len(chetDomain) > 0 {
			r.Param.ChetDomain = chetDomain
		}
		if mtrackDomain, ok := jssdkDomainConf["mtrack"]; ok && len(mtrackDomain) > 0 {
			r.Param.MTrackDomain = mtrackDomain
		}
		if cdnDomain, ok := jssdkDomainConf["cdn"]; ok && len(cdnDomain) > 0 {
			r.Param.JssdkCdnDomain = cdnDomain
		}
	}
}

func renderNewMoreOfferFlag(r *mvutil.RequestParams) {
	r.Param.NewMoreOfferFlag = false
	// 只针对 sdk iv，rv流量
	if !mvutil.IsHbOrV3OrV5Request(r.Param.RequestPath) {
		return
	}
	// 删除rv，iv的限制，使banner 也能支持more_offer
	if r.UnitInfo.Unit.MofUnitId == 0 {
		return
	}
	moreOfferConf, _ := extractor.GetMORE_OFFER_CONF()
	if len(moreOfferConf.UnitIds) > 0 {
		if mvutil.InInt64Arr(r.Param.UnitID, moreOfferConf.UnitIds) {
			r.Param.NewMoreOfferFlag = true
			return
		}
	}
	if moreOfferConf.TotalRate == 100 {
		r.Param.NewMoreOfferFlag = true
		return
	}
	rateRand := rand.Intn(100)
	if moreOfferConf.TotalRate > rateRand {
		r.Param.NewMoreOfferFlag = true
	}
}

func renderMoreOfferNewImp(r *mvutil.RequestParams) {
	r.Param.NewMofImpFlag = false
	// 只针对 sdk iv，rv流量
	if !mvutil.IsHbOrV3OrV5Request(r.Param.RequestPath) {
		return
	}
	if r.Param.AdType != mvconst.ADTypeRewardVideo && r.Param.AdType != mvconst.ADTypeInterstitialVideo {
		return
	}
	newMofImpRate, _ := extractor.GetNEW_MOF_IMP_RATE()
	if newMofImpRate == 0 {
		return
	}
	if newMofImpRate == 100 {
		r.Param.NewMofImpFlag = true
		return
	}
	rateRand := rand.Intn(100)
	if newMofImpRate > rateRand {
		r.Param.NewMofImpFlag = true
	}
}

func renderMoreOfferAbFlag(r *mvutil.RequestParams) {
	r.Param.MofAbFlag = false
	mofAbTestRate, _ := extractor.GetMOF_ABTEST_RATE()
	if mofAbTestRate == 0 {
		return
	}
	if mofAbTestRate == 100 {
		r.Param.MofAbFlag = true
		return
	}
	rateRand := rand.Intn(100)
	if mofAbTestRate > rateRand {
		r.Param.MofAbFlag = true
	}
}

func renderExcludePackageName(r *mvutil.RequestParams) {
	if r.Param.ExcludePackageNames == nil {
		r.Param.ExcludePackageNames = make(map[string]bool)
	}

	if r.Param.IsLowFlowUnitReq {
		return
	}

	renderExcludePackageNameByMtgClick(r)
	renderExcludePackageNameByAdnClick(r)
	renderExcludePackageNameByThirdPostback(r)
	renderExcludePackageNameByAnalysisOfflineInstallPackageName(r)
	// 展示过不召回实验
	// TODO
	// 后面频次控制会切到通用频次控制
	// 为帮助通用频次控制效果的数据分析, 展示过不召回通过约定字段单独传给 pioneer
	renderExcludePackageNameByMtgImpression(r)

	renderExcludePackageNameByCityCode(r)
}

func renderExcludePackageNameByCityCode(r *mvutil.RequestParams) {
	if !mvutil.IsHbOrV3OrV5Request(r.Param.RequestPath) && r.Param.RequestType != mvconst.REQUEST_TYPE_OPENAPI_AD &&
		r.Param.RequestType != mvconst.REQUEST_TYPE_SITE {
		return
	}
	excludePackagesByCityCodeConf := extractor.GetEXCLUDE_PACKAGES_BY_CITYCODE_CONF()

	if excludePackagesByCityCodeConf == nil {
		return
	}
	if !excludePackagesByCityCodeConf.Status {
		return
	}
	if !mvutil.InInt64Arr(r.Param.CityCode, excludePackagesByCityCodeConf.CityCodeList) {
		return
	}
	var blockPackagelist []string
	if r.Param.Platform == mvconst.PlatformAndroid {
		blockPackagelist = excludePackagesByCityCodeConf.BlockAndroidPackageList
	} else {
		blockPackagelist = excludePackagesByCityCodeConf.BlockIosPackageList
	}

	for _, pkg := range blockPackagelist {
		r.Param.ExcludePackageNames[pkg] = true
	}
}

func renderExcludePackageNameByMtgImpression(r *mvutil.RequestParams) {
	cfgs := extractor.GetExcludeImpressionPackagesV2()
	// 过滤不符规则的流量
	var cfg *mvutil.ExcludeClickPackages
	found, idx := canExcludePkgTestV2(cfgs, r)
	if !found {
		return
	}
	if idx != -1 {
		cfg = cfgs.Configs[idx]
	}
	// 生成key
	devKey := mvutil.GetGlobalUniqDeviceTag(&r.Param)
	// impression无需加前缀，单独的一条sqs
	devKey = strings.ToLower(devKey)
	// 切量结果，记录标记
	r.Param.ExtDataInit.ImpExcludePkg = renderExcludePkgTest(r, cfg, devKey, true)
}

func OnlineCanExcludePkg(r *mvutil.RequestParams) bool {
	if r.Param.RequestType != mvconst.REQUEST_TYPE_OPENAPI_AD {
		return false
	}
	conf := extractor.GetONLINE_FREQUENCY_CONTROL_CONF()
	rate := getTestConfRate(r, conf)
	var randVal int
	// 按设备切量
	if !mvutil.IsDevidEmpty(&r.Param) {
		randVal = mvutil.GetRandConsiderZero(r.Param.GAID, r.Param.IDFA, mvconst.SALT_ONLINE_FREQUENCY_ABTEST, 100)
	} else if mvutil.HasDevidInAndroidCN(&r.Param) {
		randVal = mvutil.GetRandConsiderZeroWithAndroidIdAndImei(r.Param.AndroidID, r.Param.IMEI, mvconst.SALT_ONLINE_FREQUENCY_ABTEST, 100)
	} else {
		randVal = 0
	}
	if rate > randVal {
		// 记录标记，此标记fluentd会使用，判断是否将设备流发到sqs中。
		r.Param.ExtDataInit.SqsCollect = 1
		return true
	}
	return false
}

func canExcludePkgTestV2(cfgs *mvutil.ExcludeClickPackagesV2, r *mvutil.RequestParams) (bool, int) {
	// mtg 点击过，SDK 流量不召回
	// online api 也按开关控制是否做频次控制
	if (!mvutil.IsHbOrV3OrV5Request(r.Param.RequestPath) || r.Param.RequestType != mvconst.REQUEST_TYPE_OPENAPI_V3) &&
		!OnlineCanExcludePkg(r) {
		return false, -1
	}

	idx := -1
	var filterTag bool
	for i, cfg := range cfgs.Configs {
		if !cfg.Status {
			// return false
			continue
		}

		// filter
		if len(cfg.Platform) > 0 && !mvutil.InArray(r.Param.Platform, cfg.Platform) {
			// return false
			continue
		}

		if len(cfg.PubBList) > 0 && mvutil.InInt64Arr(r.Param.PublisherID, cfg.PubBList) {
			// return false
			continue
		}

		if len(cfg.AppBList) > 0 && mvutil.InInt64Arr(r.Param.AppID, cfg.AppBList) {
			// return false
			continue
		}

		// cfg限制了ad_type，对于online api 不限制ad_type
		if len(cfg.AdTypeBList) > 0 && mvutil.InInt32Arr(r.Param.AdType, cfg.AdTypeBList) && r.Param.ExtDataInit.SqsCollect != 1 {
			// return false
			continue
		}

		if len(cfg.PubWList) > 0 && !mvutil.InInt64Arr(r.Param.PublisherID, cfg.PubWList) {
			// return false
			continue
		}

		if len(cfg.AppWList) > 0 && !mvutil.InInt64Arr(r.Param.AppID, cfg.AppWList) {
			// return false
			continue
		}

		if len(cfg.AdTypeWList) > 0 && !mvutil.InInt32Arr(r.Param.AdType, cfg.AdTypeWList) {
			// return false
			continue
		}

		filterTag = true
		idx = i
		break
	}

	return filterTag, idx
}

func canExcludePkgTest(cfg *mvutil.ExcludeClickPackages, r *mvutil.RequestParams) bool {
	// mtg 点击过，SDK 流量不召回
	// online api 也按开关控制是否做频次控制
	if (!mvutil.IsHbOrV3OrV5Request(r.Param.RequestPath) || r.Param.RequestType != mvconst.REQUEST_TYPE_OPENAPI_V3) &&
		!OnlineCanExcludePkg(r) {
		return false
	}

	if !cfg.Status {
		return false
	}

	// filter
	if len(cfg.Platform) > 0 && !mvutil.InArray(r.Param.Platform, cfg.Platform) {
		return false
	}

	if len(cfg.PubBList) > 0 && mvutil.InInt64Arr(r.Param.PublisherID, cfg.PubBList) {
		return false
	}

	if len(cfg.AppBList) > 0 && mvutil.InInt64Arr(r.Param.AppID, cfg.AppBList) {
		return false
	}

	// cfg限制了ad_type，对于online api 不限制ad_type
	if len(cfg.AdTypeBList) > 0 && mvutil.InInt32Arr(r.Param.AdType, cfg.AdTypeBList) && r.Param.ExtDataInit.SqsCollect != 1 {
		return false
	}

	if len(cfg.PubWList) > 0 && !mvutil.InInt64Arr(r.Param.PublisherID, cfg.PubWList) {
		return false
	}

	if len(cfg.AppWList) > 0 && !mvutil.InInt64Arr(r.Param.AppID, cfg.AppWList) {
		return false
	}

	if len(cfg.AdTypeWList) > 0 && !mvutil.InInt32Arr(r.Param.AdType, cfg.AdTypeWList) {
		return false
	}
	return true
}

func renderExcludePkgTest(r *mvutil.RequestParams, cfg *mvutil.ExcludeClickPackages, devKey string, isImp bool) string {
	rates := cfg.GetRates()
	clickInterval := mvutil.RandByRateInMapV2(rates, func(sum int) int {
		if sum <= 0 {
			return 0
		}
		return mvutil.GetRandByStr(devKey, sum)
	})

	var tag string
	v, ok := cfg.TagRates[clickInterval]
	if !ok {
		return tag
	}
	tag = v.Tag
	if subRate := cfg.GetSubRates(clickInterval); len(subRate) > 0 {
		// 有配置BB实验，则进行BB实验切量，得到BB实验结果标签
		tag = mvutil.RandByRateInMapV2(subRate, func(sum int) int {
			if sum <= 0 {
				return 0
			}
			return mvutil.GetRandByStr(devKey+"_bb", sum)
		})
	}

	value, err := mkv.GetMKVFieldVal(devKey)
	if err != nil {
		mvutil.Logger.AerospikeLog.Errorf("get aerospike data error.error=[%s],requestid=[%s],key=[%s]", err.Error(), r.Param.RequestID, devKey)
		// 这部分认为无点击包名
		return tag + "_1"
	}

	skipTest := false // 默认进行实验
	interval, _ := strconv.Atoi(clickInterval)
	if interval == 0 {
		skipTest = true
		interval = cfg.ControlGroupTime
	}

	// 约定如果tag前缀为bb_开头，那么不进行实质实验，只做切量
	if strings.Index(tag, "bb_") == 0 {
		skipTest = true
	}

	now := time.Now().Unix()
	hasPkg := false
	for pkg, ts := range value {
		if len(cfg.PackageBList) > 0 && mvutil.InStrArray(pkg, cfg.PackageBList) {
			continue
		}

		if len(cfg.PackageWList) > 0 && !mvutil.InStrArray(pkg, cfg.PackageWList) {
			continue
		}

		timestamp, _ := strconv.ParseInt(string(ts), 10, 64)
		if int(now-timestamp) <= interval {
			hasPkg = true
			if skipTest {
				// 如果有包名但是没有实质的实验动作，则退出循环
				break
			}
			r.Param.ExcludePackageNames[pkg] = true
			// 展示过不召回走通用频次控制
			if isImp {
				r.Param.ImpExcPkgNames = append(r.Param.ImpExcPkgNames, pkg)
			}
		}
	}

	if !hasPkg {
		return tag + "_1"
	}
	return tag + "_2"
}

func renderExcludePackageNameByMtgClick(r *mvutil.RequestParams) {
	cfg := extractor.GetExcludeClickPackages()
	// 过滤不符规则的流量
	if !canExcludePkgTest(cfg, r) {
		return
	}

	devKey := mvutil.GetGlobalUniqDeviceTag(&r.Param)
	devKey = "tc_" + strings.ToLower(devKey)
	// 记录标记
	r.Param.ExtDataInit.ExcludePkg = renderExcludePkgTest(r, cfg, devKey, false)
}

func renderExcludePackageNameByAdnClick(r *mvutil.RequestParams) {
	devKey := r.Param.GAID
	if r.Param.Platform == mvconst.PlatformIOS {
		devKey = r.Param.IDFA
	}

	devKey = strings.ToLower(devKey)

	if len(devKey) == 0 || devKey == "0" || devKey == "00000000-0000-0000-0000-000000000000" ||
		devKey == "gaid" || devKey == "idfa" || devKey == "null" || devKey == "-" {
		return
	}
	r.Param.DeviceKey = devKey

	blackConfig, _ := extractor.GetBLACK_FOR_EXCLUDE_PACKAGE_NAME()
	if !blackConfig.Status {
		return
	}

	// 有效时间
	durations := blackConfig.Durations
	if durations <= 0 {
		durations = 86400
	}

	if mvutil.InInt64Arr(r.Param.PublisherID, blackConfig.PublisherIds) {
		return
	}

	if mvutil.InInt64Arr(r.Param.AppID, blackConfig.AppIds) {
		return
	}

	devKey = "pkg:" + devKey // add salt
	value, err := mkv.GetMKVFieldVal(devKey)
	if err != nil {
		mvutil.Logger.AerospikeLog.Errorf("get aerospike data error.error=[%s],requestid=[%s],key=[%s]", err.Error(), r.Param.RequestID, devKey)
		return
	}
	r.Param.DeviceInstalledPackages = value
	now := r.Param.RequestTime
	for pkg, ts := range value {
		if strings.HasPrefix(pkg, "psb_") {
			continue
		}

		if strings.HasPrefix(pkg, "aop_") {
			// 离线上报安装数据
			continue
		}
		timestamp, _ := strconv.ParseInt(string(ts), 10, 64)
		if now-timestamp <= durations {
			r.Param.ExcludePackageNames[pkg] = true
		}
	}
	return
}

func renderExcludePackageNameByAnalysisOfflineInstallPackageName(r *mvutil.RequestParams) {
	// 无设备信息
	if len(r.Param.DeviceKey) == 0 {
		return
	}

	// postback 切量配置
	cfg := extractor.GetAopReduplicateConfig()
	if cfg.Status == 0 {
		return
	}

	if !mvutil.InArray(r.Param.Platform, cfg.Platform) {
		return
	}

	deviceInstalledPackages := r.Param.DeviceInstalledPackages
	devKey := "aop:" + r.Param.DeviceKey

	// 默认为aop_0,当有值时，才进入实验切量，aop_2为对照组, aop1为实验组
	tag := "aop_0"
	pass := false
	if mvutil.GetRandByStr(devKey, 10000) < cfg.Rate {
		pass = true // 实验切量
	}

	wtickPkgs := getWtickPackage()
	for pkg := range deviceInstalledPackages {
		if !strings.HasPrefix(pkg, "aop_") {
			continue
		}

		tag = "aop_2" // 对照组
		if pass {
			// 实验组
			tag = "aop_1"
			if realPkg := strings.TrimPrefix(pkg, "aop_"); !mvutil.InStrArray(realPkg, cfg.BlackPkg) && !mvutil.InStrArray(realPkg, wtickPkgs) {
				r.Param.ExcludePackageNames[realPkg] = true
			}

		}
	}
	r.Param.ExtDataInit.ExcludeAopPkg = tag
	return
}

func getWtickPackage() []string {
	cfg := extractor.GetSupportSmartVBAConfig()
	if !cfg.Status {
		return nil
	}

	pkgs := make([]string, 10)
	for _, item := range cfg.Items {
		if item.WTick != nil {
			if item.WTick.Selector != nil {
				pkgs = append(pkgs, item.WTick.Selector.IncCampaignPackageName...)
			}
		}
	}
	return pkgs
}

// renderExcludePackageNameByThirdPostback third postback 包名去重
func renderExcludePackageNameByThirdPostback(r *mvutil.RequestParams) {
	// 无设备信息
	if len(r.Param.DeviceKey) == 0 {
		return
	}

	// postback 切量配置
	cfg := extractor.GetPsbReduplicateConfig()
	if cfg.Status == 0 {
		return
	}

	if !mvutil.InArray(r.Param.Platform, cfg.Platform) {
		return
	}

	deviceInstalledPackages := r.Param.DeviceInstalledPackages
	devKey := "psb:" + r.Param.DeviceKey

	if !(mvutil.GetRandByStr(devKey, 10000) < cfg.Rate) {
		// 对照组
		r.Param.ExtDataInit.ExcludePsbPkg = "psb_A_1"
		for pkg := range deviceInstalledPackages {
			realPkg := strings.TrimPrefix(pkg, "psb_")
			if strings.HasPrefix(pkg, "psb_") && !mvutil.InStrArray(realPkg, cfg.BlackPkg) {
				r.Param.ExtDataInit.ExcludePsbPkg = "psb_A_2"
				break
			}
		}
	} else {
		// 实验组
		r.Param.ExtDataInit.ExcludePsbPkg = "psb_B_1"
		for pkg := range deviceInstalledPackages {
			realPkg := strings.TrimPrefix(pkg, "psb_")
			if strings.HasPrefix(pkg, "psb_") && !mvutil.InStrArray(realPkg, cfg.BlackPkg) {
				pkg = strings.TrimPrefix(pkg, "psb_")
				r.Param.ExtDataInit.ExcludePsbPkg = "psb_B_2"
				r.Param.ExcludePackageNames[pkg] = true
			}
		}
	}
	return
}

func changeMoreOfferAdType(r *mvutil.RequestParams) {
	if r.Param.Mof == 1 && r.Param.AdType == mvconst.ADTypeAppwall {
		changeAdTypeUnitConf, _ := extractor.GetAPPWALL_TO_MORE_OFFER_UNIT()
		if len(changeAdTypeUnitConf) == 0 {
			return
		}
		unitStr := strconv.FormatInt(r.Param.UnitID, 10)
		if unitRate, ok := changeAdTypeUnitConf[unitStr]; ok {
			unitRateRand := rand.Intn(100)
			if unitRate > int32(unitRateRand) {
				r.Param.AdType = mvconst.ADTypeMoreOffer
			}
		}
	}
}

func renderVcnABTest(r *mvutil.RequestParams) {
	vcnABTestConf, _ := extractor.GetAUTO_LOAD_CACHE_ABTSET()
	if !vcnABTestConf.Status {
		return
	}

	// 针对SDK实验
	if r.Param.Scenario != mvconst.SCENARIO_OPENAPI ||
		r.Param.RequestType != mvconst.REQUEST_TYPE_OPENAPI_V3 {
		return
	}

	// 白名单unit才做实验
	if !mvutil.InInt64Arr(r.Param.UnitID, vcnABTestConf.UnitIds) {
		return
	}

	// 设备id为空值不做实验
	if mvutil.IsDevidEmpty(&r.Param) {
		return
	}

	randVal := mvutil.GetRandConsiderZero(r.Param.GAID, r.Param.IDFA, mvconst.SALT_VCN_ABTEST, 100)
	aabTestVal, randOK := mvutil.RandByRateInMap(vcnABTestConf.Rate, randVal)
	if !randOK {
		return
	}
	r.Param.VcnABTest = aabTestVal
}

func renderCloseButtonAdFlag(r *mvutil.RequestParams) {
	if !mvutil.IsHbOrV3OrV5Request(r.Param.RequestPath) {
		return
	}
	if r.Param.AdType != mvconst.ADTypeRewardVideo && r.Param.AdType != mvconst.ADTypeInterstitialVideo {
		return
	}
	closeButtonAdUnitConf, _ := extractor.GetCLOSE_BUTTON_AD_TEST_UNITS()
	if len(closeButtonAdUnitConf) == 0 {
		return
	}
	unitStr := strconv.FormatInt(r.Param.UnitID, 10)
	if unitRate, ok := closeButtonAdUnitConf[unitStr]; ok {
		unitRateRand := rand.Intn(100)
		if unitRate > int32(unitRateRand) {
			r.Param.ExtDataInit.CloseAdTag = "1"
			return
		}
	}
	// 整体切量逻辑
	adnConf, _ := extractor.GetADNET_SWITCHS()
	if clsAdRate, ok := adnConf["clsAdRate"]; ok {
		rateRand := rand.Intn(100)
		if clsAdRate > rateRand {
			r.Param.ExtDataInit.CloseAdTag = "1"
		}
	}
}

func handleSdkBannerUnitSize(r *mvutil.RequestParams) {
	if !mvutil.IsBannerOrSdkBannerOrSplash(r.Param.AdType) {
		return
	}
	unitSize := r.Param.UnitSize
	arr := strings.Split(unitSize, "x")
	if len(arr) != 2 {
		if r.Param.AdType == mvconst.ADTypeBanner { // onlineAPI banner 的默认值
			r.Param.SdkBannerUnitWidth = 320
			r.Param.SdkBannerUnitHeight = 50
		}
		return
	}
	r.Param.SdkBannerUnitWidth, _ = strconv.ParseInt(arr[0], 10, 64)
	r.Param.SdkBannerUnitHeight, _ = strconv.ParseInt(arr[1], 10, 64)
}

func renderPriceFactor(r *mvutil.RequestParams) {

	// @流量的过滤
	if r.Param.RequestType != mvconst.REQUEST_TYPE_OPENAPI_V3 {
		return
	}

	fcc := extractor.GetFREQ_CONTROL_CONFIG()
	if fcc == nil || fcc.Status != 1 {
		return
	}

	// 放到开关后
	r.Param.ExtDataInit.PriceFactor = mvconst.PriceFactor_DefaultValue
	// 频次控制- 是否发送给RS， 1=发送（Hb不处理价格），2不发送（HB需要处理价格）
	r.Param.ExtDataInit.Send2RS = 1
	// hb 流量才使用 config 里的配置开关
	if r.IsHBRequest {
		if fcc.HbReducedRate <= 0 {
			fcc.HbReducedRate = 1
		}
	}

	rule, ok := getPriceFactorRuleByFilter(r, fcc)
	if !ok {
		mvutil.Logger.Runtime.Debug("renderPriceFactor cannot found rule from config")
		return
	}

	deviceId := mvutil.GetGlobalDeviceTag(&r.Param)
	if deviceId == "" || len(deviceId) <= 0 {
		mvutil.Logger.Runtime.Debug("renderPriceFactor cannot found deviceId")
		return
	}
	deviceId = strings.ToLower(deviceId)

	group, ok := getPriceFactorGroup(deviceId, rule)
	if !ok {
		mvutil.Logger.Runtime.Debug("renderPriceFactor cannot found group from config")
		return
	}

	abTestTag := getAbTestTag(deviceId, group)
	r.Param.ExtDataInit.PriceFactorGroupName = group.GroupName
	r.Param.ExtDataInit.PriceFactorTag = abTestTag

	mkvDataItem, err := getPriceFactorFreqByMkv(r, group, deviceId, r.Param.AdType)
	if err != nil {
		mvutil.Logger.Runtime.Debug(err.Error())
		// return
	}
	// if mkvDataItem != nil {
	t := time.Now()
	var pffInt int = 0
	startTime, err := getFreqControlStartTime(group.TimeWindow, t, r.Param.TimeZone, r.Param.CountryCode)
	if err == nil {
		var pff float64 = 0
		find := false
		if mkvDataItem != nil {
			// 先拿IMP，没有再拿Req
			if mkvDataItem.Imp != nil {
				for _, v := range mkvDataItem.Imp {
					if v.Ts > startTime {
						pff = pff + v.Freq
						find = true
					}
				}
			}
			if !find && mkvDataItem.Req != nil {
				for _, v := range mkvDataItem.Req {
					if v.Ts > startTime {
						pff = pff + v.Freq
					}
				}
			}
			pffInt = int(math.Floor(pff + 0.5))
		}

		r.Param.ExtDataInit.PriceFactorFreq = &pffInt

		rpKeys := getFreqControlReplaceKeys(r)

		if pf, hit := getFreqControlValue(pffInt, group.FreqControl, rpKeys); hit {
			r.Param.ExtDataInit.PriceFactorHit = mvconst.PriceFactorHit_TRUE
			if abTestTag == mvconst.PriceFactorTag_B {
				if !r.IsHBRequest {
					r.Param.ExtDataInit.PriceFactor = pf
				} else {
					setHBPriceFactor(r, pf*fcc.HbReducedRate)
					// 命中实验B才需要给RS传，因为他去的都会是1,传了也没意义
					r.Param.ExtDataInit.Send2RS = fcc.FreqControlToRs
				}
			}
		} else {
			r.Param.ExtDataInit.PriceFactorHit = mvconst.PriceFactorHit_FALSE
		}
		return
	} else {
		mvutil.Logger.Runtime.Warn(err.Error())
	}
	// }
	// default
	r.Param.ExtDataInit.PriceFactorHit = mvconst.PriceFactorHit_FALSE
	r.Param.ExtDataInit.PriceFactorFreq = &pffInt
}

func getFreqControlKeys(r *mvutil.RequestParams) map[string]string {
	isHb := "h0"
	if r.IsHBRequest {
		isHb = "h1"
	}

	rp := map[string]string{
		"{ad_type}":      strconv.FormatInt(int64(r.Param.AdType), 10),
		"{app_id}":       "a" + strconv.FormatInt(r.Param.AppID, 10),
		"{publisher_id}": "pb" + strconv.FormatInt(r.Param.PublisherID, 10),
		"{is_hb}":        isHb,
		"{placement_id}": "pc" + r.Param.ExtPlacementId,
	}

	return rp
}

func getFreqControlReplaceKeys(r *mvutil.RequestParams) map[string]string {
	isHb := "2"
	if r.IsHBRequest {
		isHb = "1"
	}

	rp := map[string]string{
		"{ad_type}":      strconv.FormatInt(int64(r.Param.AdType), 10),
		"{platform}":     r.Param.PlatformName,
		"{country_code}": r.Param.CountryCode,
		"{unit_id}":      strconv.FormatInt(r.Param.UnitID, 10),
		"{app_id}":       strconv.FormatInt(r.Param.AppID, 10),
		"{publisher_id}": strconv.FormatInt(r.Param.PublisherID, 10),
		"{is_hb}":        isHb,
		"{placement_id}": r.Param.ExtPlacementId,
	}

	if mvutil.HasDevid(&r.Param) {
		rp["{has_devid}"] = "1"
	} else {
		rp["{has_devid}"] = "2"
	}

	return rp
}

func getFreqControlValue(pff int, fcg []*mvutil.FreqControlGroupsItemFreqControlItem, rpKeys map[string]string) (float64, bool) {
	keys, hit := getFreqControlKey(pff, fcg) // 能在配置中找到对应的key表示命中规则
	if hit && keys != nil && len(keys) > 0 {
		for _, oriKey := range keys {
			key := oriKey
			for k, v := range rpKeys {
				key = strings.Replace(key, k, v, -1)
			}
			if pf, ifFind := extractor.GetFreqControlPriceFactor(key); ifFind {
				rate := pf.FactorRate
				if rate <= mvconst.PriceFactor_MAXValue && rate > mvconst.PriceFactor_MINValue {
					return rate, true
				} else {
					mvutil.Logger.Runtime.Warnf("GetFreqControlPriceFactor(%s)=%f", key, pf.FactorRate)
				}
			}
		}
	}

	return mvconst.PriceFactor_DefaultValue, hit
}

func setHBPriceFactor(in *mvutil.RequestParams, pf float64) {
	if pf > constant.PriceFactor_MAXValue || pf <= constant.PriceFactor_MINValue {
		in.Param.ExtDataInit.PriceFactor = constant.PriceFactor_DefaultValue
	} else {
		in.Param.ExtDataInit.PriceFactor = pf
	}
}

func getFreqControlKey(pff int, fcg []*mvutil.FreqControlGroupsItemFreqControlItem) ([]string, bool) {

	for _, v := range fcg {
		if pff >= v.Min && pff < v.Max {
			return v.Keys, true
		}
	}

	return nil, false
}

func getFreqFromAerospike(r *mvutil.RequestParams, devicId string) (map[string][]byte, error) {
	if r.FreqDataFromAerospike == nil {
		// 去Aerospike获取数据
		data, err := mkv.GetMKVFieldVal(mvconst.PriceFactor_AerospikePrefix + devicId)
		if err != nil {
			mvutil.Logger.AerospikeLog.Errorf("get aerospike data error.error=[%s],requestid=[%s],key=[%s]", err.Error(), r.Param.RequestID, mvconst.PriceFactor_AerospikePrefix+devicId)
			return nil, errors.New("getFreqFromAerospike(),mkv is error!" + err.Error())
		}
		r.FreqDataFromAerospike = data
	}

	return r.FreqDataFromAerospike, nil
}

func getPriceFactorFreqByMkv(r *mvutil.RequestParams, group *mvutil.FreqControlGroupItem, devicId string, adType int32) (
	*mvutil.FreqControlMkvData, error) {

	if len(devicId) <= 0 {
		return nil, errors.New("getPriceFactorFreqByMkv(),deviceId is errror!" + devicId)
	}

	// 去Aerospike获取数据
	data, err := getFreqFromAerospike(r, devicId)
	if err != nil {
		return nil, errors.New("getPriceFactorFreqByMkv(),mkv is error!" + err.Error())
	}

	var subKey string
	if group.FreqControlKey == "" {
		subKey = "adt_" + strconv.FormatInt(int64(adType), 10)
	} else {
		subKey = group.FreqControlKey
		rpKeys := getFreqControlKeys(r)
		for k, v := range rpKeys {
			subKey = strings.Replace(subKey, k, v, -1)
		}
	}

	if jsonStr, ok := data[subKey]; ok {
		var json = jsoniter.ConfigCompatibleWithStandardLibrary
		value := new(mvutil.FreqControlMkvData)
		err := json.Unmarshal(jsonStr, &value)
		if err != nil {
			return nil, errors.New("getPriceFactorFreqByMkv()," + err.Error())
		}
		return value, nil
		// if value.Imp != nil {
		//	return value.Imp, nil
		// }
		// if value.Req != nil {
		//	return value.Req, nil
		// }
	}
	return nil, nil
}

func getFreqControlStartTime(tw *mvutil.FreqControlGroupsItemTimeWindow, t time.Time, timezone string, countryCode string) (int, error) {
	if tw == nil {
		return 0, errors.New("time_window is empty")
	}

	if tw.Mode == mvconst.FreqControlTimeWindowModeByDate {
		var tzNum int
		var ok bool
		if tzConf := extractor.GetTIMEZONE_CONFIG(); tzConf != nil {
			tzNum, ok = tzConf[timezone]
		}
		if !ok {
			if ccConfig := extractor.GetCOUNTRY_CODE_TIMEZONE_CONFIG(); ccConfig != nil {
				tzNum, ok = ccConfig[countryCode]
			}
		}
		if !ok { // 兜底
			tzNum = 8
		}
		// 获取对应时区的时间
		tt, err := mvutil.GetLocationTimeFromTimezone(tzNum, t)
		if err != nil {
			return 0, err
		}

		newTimeStr := tt.Format("2006-01-02 ") + fmt.Sprintf("%02d", tw.StartHour) + ":00:00"
		newTT, err := time.ParseInLocation("2006-01-02 15:04:05", newTimeStr, tt.Location())
		if err != nil {
			return 0, err
		}

		if tt.Hour() >= tw.StartHour {
			return int(newTT.Unix()), nil
		} else {
			return int(newTT.Unix() - 86400), nil
		}
	} else if tw.Mode == mvconst.FreqControlTimeWindowModeByHour {
		return int(t.Unix()) - tw.WindowSec, nil
	}
	return 0, errors.New("time_window's mode is error")
}

func getAbTestTag(deviceId string, group *mvutil.FreqControlGroupItem) int {

	randInt := mvutil.GetRandConsiderZero(deviceId, "", mvconst.SALT_PRICE_FACTOR_RATE_ABTEST, 100)
	// fmt.Println("RATE:", randInt)
	if group.Rate <= randInt {
		return mvconst.PriceFactorTag_A
	}

	if group.SubRate == 0 { // 表示没有 BB的概率了
		return mvconst.PriceFactorTag_B
	}

	randInt = mvutil.GetRandConsiderZero(deviceId, "", mvconst.SALT_PRICE_FACTOR_SUBRATE_ABTEST, 100)
	// fmt.Println("SUBRATE:", randInt)
	if group.SubRate <= randInt {
		return mvconst.PriceFactorTag_B
	}

	return mvconst.PriceFactorTag_BB
}

// 获取中签group信息
func getPriceFactorGroup(deviceId string, rule *mvutil.FreqControlRule) (*mvutil.FreqControlGroupItem, bool) {
	randInt := mvutil.GetRandConsiderZero(deviceId, "", mvconst.SALT_PRICE_FACTOR_GROUP_ABTEST, 100)
	var max int
	for _, group := range rule.Groups {
		max = max + group.GroupRate
		if randInt < max {
			return group, true
		}
	}
	return nil, false
}

func checkFilterItem(str string, filterItem *mvutil.FreqControlFilterItem) bool {
	if filterItem == nil {
		return true
	}
	ok := false
	if filterItem.Op == "in" {
		ok = mvutil.InStrArray(str, filterItem.Value)
		if ok == false {
			return false
		}
	}
	if filterItem.Op == "not_in" {
		ok = !mvutil.InStrArray(str, filterItem.Value)
	}

	return ok
}

func getPriceFactorRuleByFilter(req *mvutil.RequestParams, fcc *mvutil.FreqControlConfig) (*mvutil.FreqControlRule, bool) {
	if fcc == nil || fcc.Rules == nil || len(fcc.Rules) == 0 {
		return nil, false
	}

	isHb := "2"
	if req.IsHBRequest {
		isHb = "1"
	}

	for _, rule := range fcc.Rules {
		// filter
		if rule.Groups == nil {
			continue
		}
		if rule.Filter != nil {
			if ok := checkFilterItem(req.Param.PlatformName, rule.Filter.Platform); !ok {
				continue
			}
			if ok := checkFilterItem(strconv.FormatInt(int64(req.Param.AdType), 10), rule.Filter.AdType); !ok {
				continue
			}
			if ok := checkFilterItem(req.Param.CountryCode, rule.Filter.CountryCode); !ok {
				continue
			}
			if ok := checkFilterItem(strconv.FormatInt(req.Param.PublisherID, 10), rule.Filter.PublisherId); !ok {
				continue
			}
			if ok := checkFilterItem(strconv.FormatInt(req.Param.AppID, 10), rule.Filter.AppId); !ok {
				continue
			}
			if ok := checkFilterItem(strconv.FormatInt(req.Param.UnitID, 10), rule.Filter.UnitId); !ok {
				continue
			}
			if ok := checkFilterItem(req.Param.ExtPlacementId, rule.Filter.PlacementId); !ok {
				continue
			}
			if ok := checkFilterItem(isHb, rule.Filter.IsHb); !ok {
				continue
			}
			devid := "2"
			if mvutil.HasDevid(&req.Param) {
				devid = "1"
			}
			if ok := checkFilterItem(devid, rule.Filter.HasDevid); !ok {
				continue
			}
		}

		return rule, true
	}

	return nil, false
}

//
//type ExtData2MAS struct {
//	PriceFactor                       float64 `json:"pf,omitempty"`       // 频次控制- 价格系数
//	PriceFactorGroupName              string  `json:"pf_g,omitempty"`     // 频次控制- 实验组名称
//	PriceFactorTag                    int     `json:"pf_t,omitempty"`     // 频次控制- 实验标签，1=A, 2=B, 3=B'
//	PriceFactorFreq                   *int    `json:"pf_f,omitempty"`     // 频次控制- 获取到当前的频次
//	PriceFactorHit                    int     `json:"pf_h,omitempty"`     // 频次控制- 是否能命中概率， 1=命中，2=不命中
//	ImpressionCap                     int     `json:"imp_c,omitempty"`    // placement的impressionCap
//	IfSupDco                          int     `json:"dco,omitempty"`      // 标记dco切量结果
//	V5AbtestTag                       string  `json:"v5_t,omitempty"`     // V5的实验标记， 5_5, 5_3, 或者控
//	SDKOpen                           int     `json:"sdk_open,omitempty"` // SDK是否开源版本
//	BandWidth                         int64   `json:"bw,omitempty"`       // 带宽
//	TKSysTag                          string  `json:"tkst,omitempty"`     // tracking 集群切量
//	UpTime                            string  `json:"uptime,omitempty"`   // 系统开机时间
//	Brt                               string  `json:"brt,omitempty"`      // 屏幕亮度
//	Vol                               string  `json:"vol,omitempty"`      // 音量
//	Lpm                               string  `json:"lpm,omitempty"`      // 是否为低电量模式
//	Font                              string  `json:"font,omitempty"`     // 设备默认字体大小
//	ImeiABTest                        int     `json:"imei_abt,omitempty"` // imei abtest
//	IsReturnWtickTag                  int     `json:"irwt,omitempty"`     // 是否给sdk返回wtick=1
//	Ntbarpt                           int     `json:"ntbarpt,omitempty"`
//	Ntbarpasbl                        int     `json:"ntbarpasbl,omitempty"`
//	AtatType                          int     `json:"atatType,omitempty"`
//	MappingIdfaTag                    string  `json:"mp_idfa,omitempty"`     // mapping idfa abtest标记
//	ExcludeAopPkg                     string  `json:"exc_aop_pkg,omitempty"` // analysis offline package 不召回实验标记
//	TKCNABTestTag                     int     `json:"tkcn"`                  // tracking cn 集群切量标记
//	TKCNABTestAATag                   int     `json:"tkcnaa"`                // tracking cn 集群切量AA标记
//	MappingIdfaCoverIdfaTag           string  `json:"mp_ici,omitempty"`      // mapping idfa 替换idfa的abtest标记
//	MiskSpt                           string  `json:"misk_spt,omitempty"`    // 是否支持小米storekit。1表示支持，0表示不支持。-1表示未安装小米商店
//	MoreofferAndAppwallMvToPioneerTag string  `json:"maamtp,omitempty"`      // more_offer/appwall迁移aabtest标记
//	ParentUnitId                      int64   `json:"parent_id,omitempty"`
//	H5Type                            int     `json:"h5_t,omitempty"`
//	MofType                           int     `json:"mof_type,omitempty"` // 区分是more offer 还是 close button ad
//	CrtRid                            string  `json:"crt_rid,omitempty"`  // 主offer的request_id
//	CNTrackDomain                     int     `json:"cntd,omitempty"`     // 中国专线tracking域名切量标记
//	OnlineApiNeedOfferBidPrice        string  `json:"olnobp,omitempty"`   // 表示算法是否需要针对hb online api请求的每个offer单独出价 1表示需要
//	TrackDomainByCountryCode          int     `json:"tdbcc,omitempty"`    // 根据country code选择的tracking域名
//	Dnt                               string  `json:"dnt,omitempty"`      // dnt值为1时，表示用户退出个性化广告
//}
//
//func handleExtData2MAS(r *mvutil.RequestParams) {
//	var extData2MAS ExtData2MAS
//	needToEncode := true
//	fcc := extractor.GetFREQ_CONTROL_CONFIG()
//	if fcc != nil && fcc.Status == 1 {
//		extData2MAS.PriceFactor = r.Param.ExtDataInit.PriceFactor
//		extData2MAS.PriceFactorGroupName = r.Param.ExtDataInit.PriceFactorGroupName
//		extData2MAS.PriceFactorTag = r.Param.ExtDataInit.PriceFactorTag
//		extData2MAS.PriceFactorFreq = r.Param.ExtDataInit.PriceFactorFreq
//		extData2MAS.PriceFactorHit = r.Param.ExtDataInit.PriceFactorHit
//		needToEncode = true
//	}
//
//	if r.PlacementInfo != nil && r.PlacementInfo.ImpressionCap > 0 {
//		extData2MAS.ImpressionCap = r.PlacementInfo.ImpressionCap
//		needToEncode = true
//	}
//	// v5
//	if len(r.Param.ExtDataInit.V5AbtestTag) > 0 {
//		extData2MAS.V5AbtestTag = r.Param.ExtDataInit.V5AbtestTag
//		needToEncode = true
//	}
//
//	if r.Param.Open == 1 {
//		extData2MAS.SDKOpen = r.Param.Open
//		needToEncode = true
//	}
//
//	if r.Param.BandWidth > 0 {
//		extData2MAS.BandWidth = r.Param.BandWidth
//		needToEncode = true
//	}
//
//	extData2MAS.TKSysTag = r.Param.ExtDataInit.TKSysTag
//
//	// 记录ios 设备信息。
//	extData2MAS.UpTime = r.Param.UpTime
//	extData2MAS.Brt = r.Param.Brt
//	extData2MAS.Vol = r.Param.Vol
//	extData2MAS.Lpm = r.Param.Lpm
//	extData2MAS.Font = r.Param.Font
//	extData2MAS.ImeiABTest = r.Param.ExtDataInit.ImeiABTest
//
//	extData2MAS.IsReturnWtickTag = r.Param.IsReturnWtick
//
//	extData2MAS.Ntbarpt = r.Param.Ntbarpt
//	extData2MAS.Ntbarpasbl = r.Param.Ntbarpasbl
//	extData2MAS.AtatType = r.Param.AtatType
//	extData2MAS.MappingIdfaTag = r.Param.ExtDataInit.MappingIdfaTag
//	extData2MAS.ExcludeAopPkg = r.Param.ExtDataInit.ExcludeAopPkg
//	extData2MAS.MappingIdfaCoverIdfaTag = r.Param.ExtDataInit.MappingIdfaCoverIdfaTag
//
//	extData2MAS.TKCNABTestTag = r.Param.ExtDataInit.TKCNABTestTag
//	extData2MAS.TKCNABTestAATag = r.Param.ExtDataInit.TKCNABTestAATag
//
//	extData2MAS.MiskSpt = r.Param.MiSkSpt
//	extData2MAS.MoreofferAndAppwallMvToPioneerTag = r.Param.ExtDataInit.MoreofferAndAppwallMvToPioneerTag
//
//	if r.Param.UcParentUnitId > 0 {
//		extData2MAS.ParentUnitId = r.Param.UcParentUnitId
//	} else {
//		extData2MAS.ParentUnitId = r.Param.ParentUnitId
//	}
//	extData2MAS.H5Type = r.Param.H5Type
//	extData2MAS.MofType = r.Param.MofType
//	extData2MAS.CrtRid = r.Param.ExtDataInit.CrtRid
//	extData2MAS.CNTrackDomain = r.Param.ExtDataInit.CNTrackDomain
//	extData2MAS.OnlineApiNeedOfferBidPrice = r.Param.OnlineApiNeedOfferBidPrice
//	extData2MAS.TrackDomainByCountryCode = r.Param.ExtDataInit.TrackDomainByCountryCode
//	extData2MAS.Dnt = r.Param.Dnt
//
//	if needToEncode {
//		jsonvalue, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(extData2MAS)
//		r.Param.ExtData2MAS = string(jsonvalue)
//	}
//}

func renderCountryExcludePackage(r *mvutil.RequestParams) {
	// 针对sdk流量生效
	if !mvutil.IsHbOrV3OrV5Request(r.Param.RequestPath) && r.Param.RequestType != mvconst.REQUEST_TYPE_OPENAPI_AD &&
		r.Param.RequestType != mvconst.REQUEST_TYPE_SITE {
		return
	}
	countryBlackPkgListConf := extractor.GetCountryBlackPackageListConf()
	if platformCountryBlackPkgList, ok := countryBlackPkgListConf[r.Param.CountryCode]; ok {
		if countryBlackPkgList, ok := platformCountryBlackPkgList[r.Param.PlatformName]; ok && len(countryBlackPkgList) > 0 {
			r.Param.CountryBlackPackageList = countryBlackPkgList
			// 加入到ExcludePackageName字段里
			if r.Param.ExcludePackageNames == nil {
				r.Param.ExcludePackageNames = make(map[string]bool)
			}
			for _, v := range countryBlackPkgList {
				r.Param.ExcludePackageNames[v] = true
			}
		}
	}
}

func renderV5Abtest(r *mvutil.RequestParams) {
	// 如果传了CC 覆盖countryCode
	if r.Param.RequestPath == mvconst.PATHOpenApiV5 {
		// if r.Method == "POST" { //非POST方法不返回V5结构体
		conf := extractor.GetV5_ABTEST_CONFIG()
		if conf.Switch == 1 { // 全开
			r.Param.ExtDataInit.V5AbtestTag = mvconst.V5_ABTEST_V5_V5
			return
		} else if conf.Switch == 2 { // 使用unit_conf配置
			if randValue, ok := conf.UnitConf[r.Param.UnitID]; ok {
				deviceId := mvutil.GetGlobalDeviceTag(&r.Param)
				if len(deviceId) > 0 {
					randInt := mvutil.GetRandConsiderZero(deviceId, "", "V5_ABTEST_CONFIG", 100)
					if randValue <= randInt { // 命中实验
						r.Param.ExtDataInit.V5AbtestTag = mvconst.V5_ABTEST_V5_V5
						return
					}
				}
			}
		}
		// }
		r.Param.ExtDataInit.V5AbtestTag = mvconst.V5_ABTEST_V5_V3
	} else if mvconst.PATHBid == r.Param.RequestPath || mvconst.PATHMopubBid == r.Param.RequestPath {
		cfg := extractor.GetHBV5ABTestConf() // hb 流量走 v5 切量配置
		for _, c := range cfg.HBV5ABTestConfigs {
			if selectorHBV5ABTestFilter(r, c) {
				r.Param.ExtDataInit.V5AbtestTag = mvconst.V5_ABTEST_V5_V3
				if mvutil.GetRandByGlobalTagId(&r.Param, mvconst.SALT_HB_V5_ABTEST, 100) < c.Rate {
					r.Param.ExtDataInit.V5AbtestTag = mvconst.V5_ABTEST_V5_V5
				}
				break
			}
		}
	}
}

func selectorHBV5ABTestFilter(r *mvutil.RequestParams, c *mvutil.HBV5ABTestConfig) bool {
	if len(c.CountryCode) > 0 && !mvutil.InStrArray(r.Param.CountryCode, c.CountryCode) {
		return false
	}
	if len(c.MediationName) > 0 && !mvutil.InStrArray(r.Param.MediationName, c.MediationName) {
		return false
	}
	if len(c.Platform) > 0 && !mvutil.InArray(r.Param.Platform, c.Platform) {
		return false
	}
	if len(c.SDKVersion) > 0 {
		sdkVersion := supply_mvutil.RenderSDKVersion(c.SDKVersion)
		if r.Param.FormatSDKVersion.SDKType != sdkVersion.SDKType || r.Param.FormatSDKVersion.SDKVersionCode < sdkVersion.SDKVersionCode {
			return false
		}
	}
	if len(c.AdType) > 0 && !mvutil.InInt32Arr(r.Param.AdType, c.AdType) {
		return false
	}
	if len(c.PublisherID) > 0 && !mvutil.InInt64Arr(r.Param.PublisherID, c.PublisherID) {
		return false
	}
	if len(c.AppID) > 0 && !mvutil.InInt64Arr(r.Param.AppID, c.AppID) {
		return false
	}
	if len(c.UnitID) > 0 && !mvutil.InInt64Arr(r.Param.UnitID, c.UnitID) {
		return false
	}

	return true
}

func renderTrackDomain(r *mvutil.RequestParams) string {
	// 对M流量进行切域名abtest
	if extractor.GetSYSTEM() != mvconst.SERVER_SYSTEM_M {
		return extractor.GetDOMAIN_TRACK()
	}

	tracks := extractor.GetDOMAIN_TRACKS()
	if len(tracks) == 0 || r.Param.CountryCode == "CN" {
		return extractor.GetDOMAIN_TRACK()
	}

	key := mvutil.RandStringByRateIntCustonRand(tracks, func() int {
		return int(crc32.ChecksumIEEE([]byte(mvconst.SALT_TRACK_DOMAIN_ABTEST + mvutil.GetGlobalUniqDeviceTag(&r.Param))))
	})

	trackDomain := extractor.GetDOMAIN_TRACK()
	if key != "" {
		// abtest result
		irate := tracks[key]
		if item, ok := irate.(*mvutil.TagRate); ok {
			trackDomain = item.Ext1
			r.Param.ExtDataInit.TKSysTag = item.Tag
		} else {
			trackDomain = extractor.GetDOMAIN_TRACK()
		}
	}

	return trackDomain
}

func renderMpToPioneerABTest(r *mvutil.RequestParams) {
	// 限制mp的流量
	if !mvutil.IsMpad(r.Param.RequestPath) {
		return
	}
	// 限制native 类型
	if r.Param.AdType != mvconst.ADTypeNative {
		return
	}
	conf := extractor.GetMpToPioneerABTestConf()
	if conf == nil {
		return
	}
	// 切量
	renderMpToPioneerTag(r, conf.TotalRate, conf.RandType)
}

func renderMpToPioneerTag(r *mvutil.RequestParams, rate map[string]int, randType int) {
	// 默认请求切量，1则为设备切量
	res := "a0"
	if randType == 1 {
		res = mvutil.RandByDeviceWithRateMap(rate, &r.Param, mvconst.SALT_MP_TO_PIONEER)
	} else {
		res = mvutil.RandByRate3(rate)
	}
	r.Param.ExtDataInit.MpToPioneerTag = res
}

func checkBidDeviceGEOCountry(r *mvutil.RequestParams) (match bool) {
	if len(r.Param.BidDeviceGEOCC) == 0 {
		r.Param.ExtDataInit.DeviceGEOCCMatch = 0
		match = true // 为空可以不用做后面的过滤判断, 所以返回 true
		return
	}
	isoCc := geo.CountryCodeMap[r.Param.CountryCode]
	if isoCc == r.Param.BidDeviceGEOCC {
		r.Param.ExtDataInit.DeviceGEOCCMatch = 1
		match = true
	} else {
		r.Param.ExtDataInit.DeviceGEOCCMatch = 2
		r.Param.ExtDataInit.ThreeLetterCountry = r.Param.BidDeviceGEOCC
	}
	return
}

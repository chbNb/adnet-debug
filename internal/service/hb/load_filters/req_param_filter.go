package load_filters

import (
	"math/rand"
	"strconv"
	"strings"
	"time"

	"gitlab.mobvista.com/ADN/adnet/internal/backend"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"gitlab.mobvista.com/ADN/adnet/internal/corsair_proto"
	"gitlab.mobvista.com/ADN/adnet/internal/extractor"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/constant"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/filter"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/helpers"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/output"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/params"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/req_context"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/storage"
	"gitlab.mobvista.com/ADN/adnet/internal/mvconst"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	"gitlab.mobvista.com/ADN/adnet/internal/process_pipeline"
	"gitlab.mobvista.com/ADN/adnet/internal/watcher"
	"gitlab.mobvista.com/ADN/exporter/metrics"
)

var (
	ReqParamFilterInputError = errors.New("req_param_filter input error")
	MissingRequiredParams    = errors.New("load token missing required params")
	TokenInvalidate          = errors.New("load token is invalidate")
	UnitIDInvalidate         = errors.New("load unit_id is invalidate")
	AdTypeInvalidate         = errors.New("load ad_type is invalidate")
	IllegalParamError        = errors.New("load params is illegal")
	LoadInfoNotFound         = errors.New("load app or unit or publish not found")
)

type ReqParamFilter struct {
}

func (rpf *ReqParamFilter) Process(data interface{}) (interface{}, error) {
	in, ok := data.(*mvutil.RequestParams)
	if !ok {
		return nil, ReqParamFilterInputError
	}

	token, _ := in.QueryMap.GetString("token", true, "")
	if len(token) == 0 {
		return in, TokenInvalidate
	}
	val, err := output.GetBidCache(token)
	if err != nil || (val.ReqParams == nil && val.Result == nil) {
		watcher.AddWatchValue(constant.GetBidCacheMiss, 1)
		metrics.IncCounterWithLabelValues(7, "0")
		if tokenStr := strings.Split(token, "_"); len(tokenStr) == 2 {
			uuidStr := strings.SplitN(tokenStr[0], "-", -1)
			if len(uuidStr) >= 6 {
				nowTS := time.Now()
				i, _ := strconv.ParseInt(uuidStr[5], 10, 64)
				tokenTS := time.Unix(i, 0)
				t := nowTS.Sub(tokenTS)
				if t.Hours() > 1 {
					metrics.IncCounterWithLabelValues(8, "0")
				}
			}
		}
		return in, errors.Wrap(filter.New(filter.QueryBidError), err.Error())
	}
	reqCtx := mvutil.ReqCtx{}
	reqCtx.ReqParams = val.ReqParams
	reqCtx.MaxTimeout = val.MaxTimeout
	reqCtx.FlowTagID = val.FlowTagID
	reqCtx.RandValue = val.RandValue
	reqCtx.AdsTest = val.AdsTest
	reqCtx.Elapsed = val.Elapsed
	reqCtx.Finished = val.Finished
	reqCtx.Backends = val.Backends
	reqCtx.OrderedBackends = val.OrderedBackends
	reqCtx.DebugModeInfo = val.DebugModeInfo
	reqCtx.IsWhiteGdt = val.IsWhiteGdt
	reqCtx.IsNativeVideo = val.IsNativeVideo

	// rerender appInfo from mongoDB
	if val.ReqParams.AppInfo == nil {
		if appInfo, exist := extractor.GetAppInfo(val.ReqParams.Param.AppID); exist {
			reqCtx.ReqParams.AppInfo = appInfo
		} else {
			// reqCtx.ReqParams.AppInfo = &smodel.AppInfo{}
			metrics.IncCounterWithLabelValues(29, "appInfo")
			mvutil.Logger.Runtime.Warnf("Load appInfo not found in Treasure Box, appId:[%d], token: [%s]", val.ReqParams.Param.AppID, token)
			return in, LoadInfoNotFound
		}
	}
	if val.ReqParams.UnitInfo == nil {
		if unitInfo, exist := extractor.GetUnitInfo(val.ReqParams.Param.UnitID); exist {
			reqCtx.ReqParams.UnitInfo = unitInfo
		} else {
			// reqCtx.ReqParams.UnitInfo = &smodel.UnitInfo{}
			metrics.IncCounterWithLabelValues(29, "unitInfo")
			mvutil.Logger.Runtime.Warnf("Load unitInfo not found in Treasure Box, UnitID:[%d], token: [%s]", val.ReqParams.Param.UnitID, token)
			return in, LoadInfoNotFound
		}
	}
	if val.ReqParams.PublisherInfo == nil {
		if publishInfo, exist := extractor.GetPublisherInfo(val.ReqParams.Param.PublisherID); exist {
			reqCtx.ReqParams.PublisherInfo = publishInfo
		} else {
			// reqCtx.ReqParams.PublisherInfo = &smodel.PublisherInfo{}
			metrics.IncCounterWithLabelValues(29, "publisherInfo")
			mvutil.Logger.Runtime.Warnf("Load publishInfo not found in Treasure Box, PublishID:[%d], token: [%s]", val.ReqParams.Param.PublisherID, token)
			return in, LoadInfoNotFound
		}
	}
	res, err := RenderResult(val, &reqCtx, token)
	if err != nil {
		// metric
		metrics.IncCounterWithLabelValues(32, err.Error())
		return in, errors.Wrap(filter.New(filter.BidServerLoadError), err.Error())
	}
	reqCtx.Result = res

	// fix aerospike recard too big
	if dmt, err := in.QueryMap.GetFloat("dmt", 0); err == nil {
		reqCtx.ReqParams.Param.Dmt = dmt
	}
	if dmf, err := in.QueryMap.GetFloat("dmf", 0); err == nil {
		reqCtx.ReqParams.Param.Dmf = dmf
	}
	if powerRate, err := in.QueryMap.GetInt("power_rate", 0); err == nil {
		reqCtx.ReqParams.Param.PowerRate = powerRate
	}
	if charging, err := in.QueryMap.GetInt("charging", 0); err == nil {
		reqCtx.ReqParams.Param.Charging = charging
	}
	reqCtx.ReqParams.Param.TotalMemory, _ = in.QueryMap.GetString("h", true, "")
	if len(reqCtx.ReqParams.Param.TotalMemory) == 0 {
		reqCtx.ReqParams.Param.TotalMemory, _ = in.QueryMap.GetString("cache1", true, "")
	}
	reqCtx.ReqParams.Param.ResidualMemory, _ = in.QueryMap.GetString("i", true, "")
	if len(reqCtx.ReqParams.Param.ResidualMemory) == 0 {
		reqCtx.ReqParams.Param.ResidualMemory, _ = in.QueryMap.GetString("cache2", true, "")
	}
	reqCtx.ReqParams.Param.Ct, _ = in.QueryMap.GetString("ct", true, "")
	if dvi, _ := in.QueryMap.GetString("dvi", true, ""); len(dvi) > 0 {
		renderLoadDeviceInfo(dvi, &reqCtx)
	}
	if chInfo, _ := in.QueryMap.GetString("ch_info", false, ""); len(chInfo) > 0 {
		reqCtx.ReqParams.Param.ChannelInfo = helpers.Base64Decode(chInfo)
	}
	// 需要接收新版本sdk的session id，用于判断是否需要生成加密的sysid，bkupid
	if sessionId, _ := in.QueryMap.GetString("a", true, ""); len(sessionId) > 0 {
		reqCtx.ReqParams.Param.SessionID = sessionId
	}
	if len(reqCtx.ReqParams.Param.SessionID) <= 0 {
		// 记录此次为新session
		reqCtx.ReqParams.Param.IsNewSession = true
		reqCtx.ReqParams.Param.SessionID = mvutil.GetRequestID()
	}

	// 记录app，unit，reward setting abtest 标记
	if appSettingId, _ := in.QueryMap.GetString("a_stid", true, ""); len(appSettingId) > 0 {
		reqCtx.ReqParams.Param.AppSettingId = appSettingId
	}

	if unitSettingId, _ := in.QueryMap.GetString("u_stid", true, ""); len(unitSettingId) > 0 {
		reqCtx.ReqParams.Param.UnitSettingId = unitSettingId
	}

	if rewardSettingId, _ := in.QueryMap.GetString("r_stid", true, ""); len(rewardSettingId) > 0 {
		reqCtx.ReqParams.Param.RewardSettingId = rewardSettingId
	}

	if trafficInfo, _ := in.QueryMap.GetString("j", true, ""); len(trafficInfo) > 0 {
		reqCtx.ReqParams.Param.TrafficInfo = trafficInfo
	}
	// v5 param
	tmpIdsSt, _ := in.QueryMap.GetString("g", false, "")
	if len(tmpIdsSt) == 0 {
		tmpIdsSt, _ = in.QueryMap.GetString("tmp_ids", false, "")
	}
	if len(tmpIdsSt) > 0 {
		reqCtx.ReqParams.Param.TmpIds = make(map[string]bool)
		tmpIdsArr := []string{}
		err := jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal([]byte(tmpIdsSt), &tmpIdsArr)
		if err == nil {
			for _, v := range tmpIdsArr {
				reqCtx.ReqParams.Param.TmpIds[v] = true
			}
		}
	}

	// 解析堆栈信息
	process_pipeline.RenderTrafficInfo(reqCtx.ReqParams)

	// 封装extdata
	reqCtx.ReqParams.Param.ExtData = RenderHBExtData(reqCtx.ReqParams)

	// 根据堆栈信息过滤
	if process_pipeline.FilterRequestByStack(reqCtx.ReqParams) {
		metrics.IncCounterWithLabelValues(26, mvconst.FilterRequestByStack)
		return in, errors.Wrap(filter.New(filter.IllegalParamError), IllegalParamError.Error())
	}

	// load param, don't use bid request param
	assemblyLoadParm(in, reqCtx.ReqParams)

	// 生成sysid，bkupid
	process_pipeline.NewSysId(reqCtx.ReqParams)
	process_pipeline.NewBkupId(reqCtx.ReqParams)

	unitID, _ := in.QueryMap.GetInt64("unit_id", 0)
	reqCtx.ReqParams.Param.BidUnitID = unitID
	adType, _ := in.QueryMap.GetInt("ad_type", 0)
	reqCtx.ReqParams.Param.BidAdType = adType
	hbAPIConfs := extractor.GetHB_API_CONFS()
	expiredOnLoad := hbAPIConfs["cache_expired_on_load"]
	hbLoadFilterConfigs := extractor.GetHBLoadFilterConfigs()
	if hbLoadFilterConfigs.Status && selectorFilter(in, hbLoadFilterConfigs) {
		if unitID == 0 || adType == 0 { // 暂对此情况不做过滤
			reqCtx.ReqParams.LoadRejectCode = filter.MissingRequiredParamsError.Int()
		} else if unitID != reqCtx.ReqParams.Param.UnitID {
			reqCtx.ReqParams.LoadRejectCode = filter.UnitIDInvalidateError.Int()
			metrics.IncCounterDyWithValues("req_param_filter-rejectcode", strconv.Itoa(filter.UnitIDInvalidateError.Int()))
			return in, errors.Wrap(filter.New(filter.UnitIDInvalidateError), UnitIDInvalidate.Error())
		} else if adType != int(reqCtx.ReqParams.Param.AdType) {
			reqCtx.ReqParams.LoadRejectCode = filter.AdTypeInvalidateError.Int()
			metrics.IncCounterDyWithValues("req_param_filter-rejectcode", strconv.Itoa(filter.AdTypeInvalidateError.Int()))
			return in, errors.Wrap(filter.New(filter.AdTypeInvalidateError), AdTypeInvalidate.Error())
		}
	}
	// load 之后将默认一小时过期的缓存设置为1分钟过期
	if expiredOnLoad {
		if req_context.GetInstance().Cfg.ConsulCfg.Aerospike.Enable {
			ratio := extractor.GetUseConsulServicesV2Ratio(mvutil.Cloud(), mvutil.Region(), "hb_aerospike")
			if ratio > 0 && ratio > rand.Float64() {
				key := &storage.ReqCtxKey{Token: token}
				client, e := req_context.GetInstance().GetMkvAerospikeClient()
				if e != nil {
					err = req_context.GetInstance().BidCacheClient.SetExReqCtx(token, val, time.Minute)
				} else {
					err = client.SetExReqCtx(key, val, time.Minute)
				}
			} else {
				err = req_context.GetInstance().BidCacheClient.SetExReqCtx(token, val, time.Minute)
			}
		} else {
			err = req_context.GetInstance().BidCacheClient.SetExReqCtx(token, val, time.Minute)
		}
	}
	return &reqCtx, nil
}

func assemblyLoadParm(req *mvutil.RequestParams, load *mvutil.RequestParams) {
	load.Param.RequestPath = req.Param.RequestPath
	load.Param.RequestURI = req.Param.RequestURI
	load.Param.ClientIP = req.Param.ClientIP
	load.Param.UserAgent = req.Param.UserAgent
	load.Param.Extra9 = req.Param.Extra9
	load.Param.ExtsystemUseragent = req.Param.ExtsystemUseragent
	load.Param.MvLine = req.Param.MvLine
	if len(load.Param.MvLine) > 0 {
		load.Param.Algorithm = load.Param.Algorithm + "-" + load.Param.MvLine
	}
}

func selectorFilter(req *mvutil.RequestParams, selector *mvutil.HBLoadFilterConfigs) bool {
	if len(selector.IncAdTypes) == 0 || (!mvutil.InInt32Arr(req.Param.AdType, selector.IncAdTypes) && !mvutil.InInt32Arr(-1, selector.IncAdTypes)) {
		return false
	}

	if len(selector.ExcAdTypes) > 0 && mvutil.InInt32Arr(req.Param.AdType, selector.ExcAdTypes) {
		return false
	}

	if len(selector.IncPublisherIds) == 0 || (!mvutil.InInt64Arr(req.Param.PublisherID, selector.IncPublisherIds) && !mvutil.InInt64Arr(-1, selector.IncPublisherIds)) {
		return false
	}

	if len(selector.ExcPublisherIds) > 0 && mvutil.InInt64Arr(req.Param.PublisherID, selector.ExcPublisherIds) {
		return false
	}

	if len(selector.IncAppIds) == 0 || (!mvutil.InInt64Arr(req.Param.AppID, selector.IncAppIds) && !mvutil.InInt64Arr(-1, selector.IncAppIds)) {
		return false
	}

	if len(selector.ExcAppIds) > 0 && mvutil.InInt64Arr(req.Param.AppID, selector.ExcAppIds) {
		return false
	}

	if len(selector.IncUnitIds) == 0 || (!mvutil.InInt64Arr(req.Param.UnitID, selector.IncUnitIds) && !mvutil.InInt64Arr(-1, selector.IncUnitIds)) {
		return false
	}

	if len(selector.ExcUnitIds) > 0 && mvutil.InInt64Arr(req.Param.UnitID, selector.ExcUnitIds) {
		return false
	}

	return true
}

// func renderNewMoreOfferFlag(loadReq *params.LoadReqData) {
// 	// if loadReq.AdType != constant.RewardVideo && loadReq.AdType != constant.InterstitialVideo {
// 	// 	return
// 	// }
// 	if loadReq.MofUnitId == 0 {
// 		return
// 	}
// 	moreOfferConf, _ := extractor.GetMoreOfferConf()
// 	if len(moreOfferConf.UnitIds) > 0 {
// 		if helpers.InInt64Arr(loadReq.UnitId, moreOfferConf.UnitIds) {
// 			loadReq.NewMoreOfferFlag = true
// 			return
// 		}
// 	}
// 	if moreOfferConf.TotalRate == 100 {
// 		loadReq.NewMoreOfferFlag = true
// 		return
// 	}
// 	rateRand := rand.Intn(100)
// 	if moreOfferConf.TotalRate > rateRand {
// 		loadReq.NewMoreOfferFlag = true
// 	}
// }

// func renderMoreOfferNewImp(loadReq *params.LoadReqData) {
// 	// if loadReq.AdType != constant.RewardVideo && loadReq.AdType != constant.InterstitialVideo {
// 	// 	return
// 	// }
// 	newMofImpRate, _ := extractor.GetNewMofImpRate()
// 	if newMofImpRate == 0 {
// 		return
// 	}
// 	if newMofImpRate == 100 {
// 		loadReq.NewMofImpFlag = true
// 		return
// 	}
// 	rateRand := rand.Intn(100)
// 	if newMofImpRate > rateRand {
// 		loadReq.NewMofImpFlag = true
// 	}
// }

// func renderMoreOfferAbFlag(loadReq *params.LoadReqData) {
// 	// loadReq.MofAbFlag = false
// 	mofAbTestRate, _ := extractor.GetMofABTestRate()
// 	if mofAbTestRate == 0 {
// 		return
// 	}
// 	if mofAbTestRate == 100 {
// 		loadReq.MofAbFlag = true
// 		return
// 	}
// 	rateRand := rand.Intn(100)
// 	if mofAbTestRate > rateRand {
// 		loadReq.MofAbFlag = true
// 	}
// }

func renderCloseButtonAdFlag(loadReq *params.LoadReqData) {
	// if loadReq.AdType != constant.RewardVideo && loadReq.AdType != constant.InterstitialVideo {
	// 	return
	// }
	closeButtonAdUnitConf, _ := extractor.GetCLOSE_BUTTON_AD_TEST_UNITS()
	// if len(closeButtonAdUnitConf) == 0 {
	// 	return
	// }
	if unitRate, ok := closeButtonAdUnitConf[loadReq.UnitIdStr]; ok {
		unitRateRand := rand.Intn(100)
		if unitRate > int32(unitRateRand) {
			loadReq.CloseAdTag = "1"
			return
		}
	}
	// 整体切量逻辑
	adnConf, _ := extractor.GetADNET_SWITCHS()
	if clsAdRate, ok := adnConf["clsAdRate"]; ok {
		rateRand := rand.Intn(100)
		if clsAdRate > rateRand {
			loadReq.CloseAdTag = "1"
		}
	}
}

func renderLoadDeviceInfo(dvi string, req *mvutil.ReqCtx) {
	if dvi == "" {
		return
	}

	deviceInfoDecode := helpers.Base64Decode(dvi)
	var deviceInfo helpers.DeviceInfo
	if err := jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal([]byte(deviceInfoDecode), &deviceInfo); err != nil {
		return
	}

	if len(deviceInfo.Dmt) > 0 {
		if dmt, err := strconv.ParseFloat(deviceInfo.Dmt, 64); err == nil {
			req.ReqParams.Param.Dmt = dmt
		}
	}
	if deviceInfo.Dmf > 0 {
		req.ReqParams.Param.Dmf = deviceInfo.Dmf
	}
	if len(deviceInfo.Ct) > 0 {
		req.ReqParams.Param.Ct = deviceInfo.Ct
	}
	return
}

func RenderHBExtData(r *mvutil.RequestParams) string {
	var extData mvutil.ExtData
	adnConf, _ := extractor.GetADNET_SWITCHS()
	if len(r.Param.DecryptTrafficInfoStr) > 0 {
		// 抽样记录下来
		randVal := rand.Intn(10000)
		if recordTrafficInfoRate, ok := adnConf["recordTrafficInfoRate"]; ok && recordTrafficInfoRate > randVal {
			extData.DecryptTrafficInfoStr = r.Param.DecryptTrafficInfoStr
		}
	}
	extDataJson, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(extData)
	extDataStr := string(extDataJson)
	return extDataStr
}

// 二阶段 bid 需求
func RenderResult(val *storage.ReqCtxVal, reqCtx *mvutil.ReqCtx, token string) (*corsair_proto.QueryResult_, error) {

	//
	if val == nil {
		return nil, nil
	}

	if val.ReqParams == nil || val.ReqParams.Param.BidServerCtx == nil {
		// 走原有路径
		return val.Result, nil
	}

	// 走到这里代表走到了二阶段load流程, 记录下metric
	metrics.IncCounterWithLabelValues(33, "load")

	//
	backendCtx, err := backend.RenderLoadResult(val, reqCtx)
	if err != nil {
		// 渲染广告失败
		return nil, err
	}
	// 渲染广告后会回填部分 reqCtx 的数据如token 这里需要重新补上
	reqCtx.ReqParams.Token = token

	//
	res := corsair_proto.QueryResult_{
		LogId:           val.Result.LogId,
		RandValue:       val.Result.RandValue,
		FlowTagId:       val.Result.FlowTagId,
		AdBackendConfig: val.Result.AdBackendConfig,
	}
	adCount := int64(0)
	var backendData []string
	process_pipeline.RenderAdsByBackend(reqCtx, mvconst.MAdx, backendCtx, &adCount, backendData, &res)

	return &res, nil
}

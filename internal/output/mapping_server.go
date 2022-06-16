package output

import (
	flatbuffers "github.com/google/flatbuffers/go"
	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
	"gitlab.mobvista.com/ADN/adnet/internal/consuls"
	"gitlab.mobvista.com/ADN/adnet/internal/extractor"
	mapping_server "gitlab.mobvista.com/ADN/adnet/internal/flat_buffers/ml/mapping_service"
	"gitlab.mobvista.com/ADN/adnet/internal/mvconst"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	"gitlab.mobvista.com/ADN/adnet/internal/utility"
	"gitlab.mobvista.com/ADN/exporter/metrics"
	"hash/crc32"
	"strconv"
	"strings"
	"time"
)

type MappingServerResponse struct {
	Message   string            `json:"message,omitempty"`
	Ruid      string            `json:"ruid,omitempty"`
	Code      int               `json:"code,omitempty"`
	DebugInfo map[string]string `json:"debug,omitempty"`
}

func RenderRuid(r *mvutil.RequestParams, fbsFrom string) {
	// 低效流量不走生成逻辑
	if r.Param.IsLowFlowUnitReq {
		return
	}

	// 限制sdk ios 流量
	if r.Param.PlatformName != mvconst.PlatformNameIOS || !mvutil.IsHbOrV3OrV5Request(r.Param.RequestPath) {
		return
	}

	// 切量逻辑
	mappingServerRateConf := extractor.GetMAPPING_SERVER_RATE_CONF()
	if mappingServerRateConf == nil {
		return
	}
	rate := -1
	// 有idfa的
	if idfaRate, ok := mappingServerRateConf["idfa"]; ok && len(mvutil.GetIdfaString(r.Param.IDFA)) > 0 {
		rate = idfaRate
	}
	// 有sys_id的
	if sysIdRate, ok := mappingServerRateConf["sys_id"]; ok && rate < 0 && len(r.Param.SysId) > 0 {
		rate = sysIdRate
	}
	// 整体切量
	if totalRate, ok := mappingServerRateConf["total"]; ok && rate < 0 {
		rate = totalRate
	}
	if rate <= 0 {
		return
	}
	// 设备切量
	randVal := int(crc32.ChecksumIEEE([]byte(mvconst.SALT_RUID+"_"+mvutil.GetGlobalUniqDeviceTag(&r.Param))) % 1000)
	if randVal > rate {
		return
	}
	builder := flatbuffers.NewBuilder(0)
	sysidKey := builder.CreateString("sysid")
	sysidVal := builder.CreateString(r.Param.SysId)

	idfvKey := builder.CreateString("idfv")
	idfvVal := builder.CreateString(r.Param.IDFV)

	bkupidKey := builder.CreateString("bkupid")
	bkupidVal := builder.CreateString(r.Param.BkupId)

	unifiedDevModelKey := builder.CreateString("unifiedDevModel")
	unifiedDevModelVal := builder.CreateString(r.Param.ExtModel)

	osvUpTimeKey := builder.CreateString("osvUpTime")
	osvUpTimeVal := builder.CreateString(r.Param.OsvUpTime)

	platformKey := builder.CreateString("platform")
	platformVal := builder.CreateString(r.Param.PlatformName)

	fromKey := builder.CreateString("from")
	r.Param.MappingServerFrom = "adnet"
	if r.IsHBRequest {
		r.Param.MappingServerFrom = "hb"
	}

	// 保持和setting一样的参数请求mapping server生成加密后的ruid
	fromVal := builder.CreateString(r.Param.MappingServerFrom)
	if len(fbsFrom) > 0 {
		fromVal = builder.CreateString(fbsFrom)
	}

	appidKey := builder.CreateString("appid")
	appidVal := builder.CreateString(strconv.FormatInt(r.Param.AppID, 10))

	ipKey := builder.CreateString("ip")
	ipVal := builder.CreateString(r.Param.ClientIP)

	uaKey := builder.CreateString("ua")
	uaVal := builder.CreateString(r.Param.UserAgent)

	uptimeKey := builder.CreateString("UpTime")
	uptimeVal := builder.CreateString(r.Param.UpTime)

	requestIdKey := builder.CreateString("uuid")
	requestIdVal := builder.CreateString(r.Param.RequestID)

	cityCodeKey := builder.CreateString("citycode")
	cityCodeVal := builder.CreateString(strconv.FormatInt(r.Param.CityCode, 10))

	// sysid
	mapping_server.ItemStart(builder)

	mapping_server.ItemAddKey(builder, sysidKey)
	mapping_server.ItemAddValue(builder, sysidVal)
	sysid := mapping_server.ItemEnd(builder)

	// idfv
	mapping_server.ItemStart(builder)
	mapping_server.ItemAddKey(builder, idfvKey)
	mapping_server.ItemAddValue(builder, idfvVal)
	idfv := mapping_server.ItemEnd(builder)

	// bkupid
	mapping_server.ItemStart(builder)
	mapping_server.ItemAddKey(builder, bkupidKey)
	mapping_server.ItemAddValue(builder, bkupidVal)
	bkupid := mapping_server.ItemEnd(builder)

	// unifiedDevModel
	mapping_server.ItemStart(builder)
	mapping_server.ItemAddKey(builder, unifiedDevModelKey)
	mapping_server.ItemAddValue(builder, unifiedDevModelVal)
	unifiedDevModel := mapping_server.ItemEnd(builder)

	// osvUpTime
	mapping_server.ItemStart(builder)
	mapping_server.ItemAddKey(builder, osvUpTimeKey)
	mapping_server.ItemAddValue(builder, osvUpTimeVal)
	osvUpTime := mapping_server.ItemEnd(builder)

	// platform
	mapping_server.ItemStart(builder)
	mapping_server.ItemAddKey(builder, platformKey)
	mapping_server.ItemAddValue(builder, platformVal)
	platform := mapping_server.ItemEnd(builder)

	// from
	mapping_server.ItemStart(builder)
	mapping_server.ItemAddKey(builder, fromKey)
	mapping_server.ItemAddValue(builder, fromVal)
	from := mapping_server.ItemEnd(builder)

	// appid
	mapping_server.ItemStart(builder)
	mapping_server.ItemAddKey(builder, appidKey)
	mapping_server.ItemAddValue(builder, appidVal)
	appid := mapping_server.ItemEnd(builder)

	// ip
	mapping_server.ItemStart(builder)
	mapping_server.ItemAddKey(builder, ipKey)
	mapping_server.ItemAddValue(builder, ipVal)
	ip := mapping_server.ItemEnd(builder)

	// ua
	mapping_server.ItemStart(builder)
	mapping_server.ItemAddKey(builder, uaKey)
	mapping_server.ItemAddValue(builder, uaVal)
	ua := mapping_server.ItemEnd(builder)

	// UpTime
	mapping_server.ItemStart(builder)
	mapping_server.ItemAddKey(builder, uptimeKey)
	mapping_server.ItemAddValue(builder, uptimeVal)
	UpTime := mapping_server.ItemEnd(builder)

	// uuid adnet传requestid
	mapping_server.ItemStart(builder)
	mapping_server.ItemAddKey(builder, requestIdKey)
	mapping_server.ItemAddValue(builder, requestIdVal)
	uuid := mapping_server.ItemEnd(builder)

	// citycode
	mapping_server.ItemStart(builder)
	mapping_server.ItemAddKey(builder, cityCodeKey)
	mapping_server.ItemAddValue(builder, cityCodeVal)
	cityCode := mapping_server.ItemEnd(builder)

	mapping_server.MessageStartKvpairsVector(builder, 13)
	builder.PrependUOffsetT(sysid)
	builder.PrependUOffsetT(idfv)
	builder.PrependUOffsetT(bkupid)
	builder.PrependUOffsetT(unifiedDevModel)
	builder.PrependUOffsetT(osvUpTime)
	builder.PrependUOffsetT(platform)
	builder.PrependUOffsetT(from)
	builder.PrependUOffsetT(appid)
	builder.PrependUOffsetT(ip)
	builder.PrependUOffsetT(ua)
	builder.PrependUOffsetT(UpTime)
	builder.PrependUOffsetT(uuid)
	builder.PrependUOffsetT(cityCode)
	items := builder.EndVector(13)

	defaultStatus := builder.CreateString("0")
	defaultVersion := builder.CreateString("v1")

	mapping_server.MessageStart(builder)
	mapping_server.MessageAddKvpairs(builder, items)
	mapping_server.MessageAddStatus(builder, defaultStatus)
	mapping_server.MessageAddVersion(builder, defaultVersion)
	message := mapping_server.MessageEnd(builder)

	builder.Finish(message)
	buf := builder.FinishedBytes()

	// 记录日志
	defer mvutil.StatMappingServerLog(&r.Param)

	// 默认超时时间20ms
	timeout := 20
	adnConf, _ := extractor.GetADNET_SWITCHS()
	if mappingServerTimeout, ok := adnConf["mappingServerTimeout"]; ok && mappingServerTimeout > 0 {
		timeout = mappingServerTimeout
	}
	timestamp := time.Now().UTC().Unix()
	timeStart := time.Now()
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.Header.SetMethod("POST")
	host := consuls.GetNode()
	reqUrl := "http://" + host + "/search?from=" + r.Param.MappingServerFrom + "&format=fb&v=1&time=" + strconv.FormatInt(timestamp, 10) + "&debug=" + r.Param.MappingServerDebug
	req.SetRequestURI(reqUrl)
	req.SetBody(buf)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)
	err := utility.HttpClientApp().Do(timeout, req, resp)
	// 统计响应时间
	timeElapsed := (float64)(time.Since(timeStart) / time.Millisecond)
	metrics.SetGaugeWithLabelValues(timeElapsed, 9)
	metrics.AddSummaryWithLabelValues(timeElapsed, 12)
	if err != nil {
		// 记录日志
		errorStr := err.Error()
		if strings.Contains(err.Error(), "write tcp") {
			errorStr = "write tcp i/o timeout"
		}
		metrics.IncCounterWithLabelValues(14, errorStr)
		mvutil.Logger.Runtime.Errorf("mapping_server request error. error:%s", err.Error())
		return
	}
	if resp.StatusCode() != 200 {
		// 记录异常code 日志
		metrics.IncCounterWithLabelValues(13, strconv.Itoa(resp.StatusCode()))
		mvutil.Logger.Runtime.Errorf("mapping_server request error status. statusCode:%v", resp.StatusCode())
	}

	var mappingServerResp MappingServerResponse
	err = jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(resp.Body(), &mappingServerResp)
	if err != nil {
		// 记录日志
		mvutil.Logger.Runtime.Errorf("mapping_server response unmarshal error. error:%s", err.Error())
		return
	}
	// code 监控
	metrics.IncCounterWithLabelValues(10, strconv.Itoa(mappingServerResp.Code))
	// ruid为空的监控
	if len(mappingServerResp.Ruid) == 0 {
		metrics.IncCounterWithLabelValues(11)
	}

	if fbsFrom == "encryption" {
		r.Param.EncryptedRuid = mappingServerResp.Ruid
	} else {
		r.Param.RuId = mappingServerResp.Ruid
	}
	r.Param.MappingServerMessage = mappingServerResp.Message
	r.Param.MappingServerResCode = strconv.Itoa(mappingServerResp.Code)
	debugInfo, err := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(mappingServerResp.DebugInfo)
	if err == nil {
		r.Param.MappingServerDebugInfo = string(debugInfo)
	}
}

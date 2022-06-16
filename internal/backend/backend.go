package backend

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
	"gitlab.mobvista.com/ADN/adnet/internal/clients"
	"gitlab.mobvista.com/ADN/adnet/internal/extractor"
	hbreqctx "gitlab.mobvista.com/ADN/adnet/internal/hb/req_context"
	"gitlab.mobvista.com/ADN/adnet/internal/mvconst"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	"gitlab.mobvista.com/ADN/adnet/internal/utility"
	"gitlab.mobvista.com/ADN/adnet/internal/watcher"
	"gitlab.mobvista.com/ADN/exporter/metrics"
	"gitlab.mobvista.com/adserver/recommend_protocols/go/ad_server"
)

const (
	ERR_OK      = iota + 0
	ERR_Timeout // backend timeout
	ERR_Internal
	ERR_EnqueueFail // enqueue backend ReqQueue fail
	ERR_RpcFail
	ERR_HttpFail
	ERR_ParseRsp
	ERR_NoAds
	ERR_RspErrCode
	ERR_HttpReadRsp
	ERR_Http4xx
	ERR_Http5xx
	ERR_Param
	ERR_BidServerReject
)

const (
	QUEUE_FACTOR  = 8
	ADSERVER_PORT = 9099
)

// 后端抽象接口
// ad_server： 需要实现getCampaigns接口
// 第三方使用http协议的后端：需要实现composeHttpRequest和parseHttpResponse接口
// 所有后端需要实现filterBackend接口
type IBackend interface {
	filterBackend(reqCtx *mvutil.ReqCtx) int
	// getCampaigns(reqCtx *mvutil.ReqCtx, backendCtx *mvutil.BackendCtx) (int, error)
	composeHttpRequest(reqCtx *mvutil.ReqCtx, backendCtx *mvutil.BackendCtx, req *fasthttp.Request) error
	parseHttpResponse(reqCtx *mvutil.ReqCtx, backendCtx *mvutil.BackendCtx) (int, error)
	getRequestNode() string
}

// Backend 定义广告后端
type Backend struct {
	Name    string
	ID      int
	Timeout int
	// CountryCode string
	// AdType      int
	// VideoAdType int
	Workers int
	// reqQueue    chan *mvutil.ReqCtx

	// logical interface
	specific IBackend

	// backend interface status managerment
	ReqCount int32 // per second
	ReqFail  int32 // per second
	Risk     int
	State    bool // true - normal, false - limit reqs to this backend

	// ad_server
	RpcAddr string
	// AsClients *AsList

	// backend use http/https protocol
	HttpURL               string
	HttpsURL              string
	Path                  string
	Method                string
	HttpClient            *http.Client
	AdServerClient        *clients.AdServerClient
	MAdxClient            *clients.MAdxClient
	PioneerClient         *clients.PioneerClient
	GRPCPoolLimit         int
	GRPCPoolFlushInterval int
	GRPCPoolCloseTimeout  int
}

func NewBackend(serviceDetail *mvutil.ServiceDetail) *Backend {
	if serviceDetail == nil {
		return nil
	}
	b := &Backend{
		Name:    serviceDetail.Name,
		ID:      serviceDetail.ID,
		Timeout: serviceDetail.Timeout,
		// CountryCode: serviceDetail.CountryCode,
		// AdType:      serviceDetail.AdType,
		// VideoAdType: serviceDetail.VideoAdType,

		Workers: serviceDetail.Workers,
		// reqQueue: make(chan *mvutil.ReqCtx, serviceDetail.Workers*QUEUE_FACTOR),

		specific: nil,

		ReqCount: 0,
		ReqFail:  0,
		Risk:     0,
		State:    true,

		// RpcAddr: serviceDetail.RpcAddr,

		HttpURL:        serviceDetail.HttpURL,
		HttpsURL:       serviceDetail.HttpsURL,
		Path:           serviceDetail.Path,
		Method:         serviceDetail.Method,
		HttpClient:     nil,
		AdServerClient: nil,
		MAdxClient:     nil,
		PioneerClient:  nil,
		// GRPCPoolLimit:         serviceDetail.Workers,
		// GRPCPoolFlushInterval: serviceDetail.ConsulCfg.Internal,
		// GRPCPoolCloseTimeout:  serviceDetail.Timeout,
	}

	switch b.ID {
	case mvconst.Mobvista:
		b.specific = &AdServer{}
	case mvconst.MAdx:
		b.specific = &MAdxBackend{}
	case mvconst.Pioneer:
		b.specific = &PioneerBackend{}
	default:
		b.specific = nil
	}

	return b
}

// adx请求
func (backend *Backend) dispatchProxyRequest(reqCtx *mvutil.ReqCtx, backendCtx *mvutil.BackendCtx) *mvutil.BackendMetric {
	metricData := &mvutil.BackendMetric{BackendId: backend.ID}
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			mvutil.Logger.Runtime.Errorf("backend ID=[%d] name=[%s] dispatchProxyRequest defer panic: %s",
				backend.ID, backend.Name, err)
			var buf [4096]byte
			n := runtime.Stack(buf[:], false)
			mvutil.Logger.Runtime.Errorf("panic Stack: %s", string(buf[:n]))
			metrics.IncCounterWithLabelValues(25)
		}
	}()
	if reqCtx == nil || backendCtx == nil || reqCtx.ReqParams == nil {
		metricData.FilterCode = mvconst.ParamInvalidate
		return metricData
	}
	backendCtx.ReqPath = backend.specific.getRequestNode()
	if reqCtx.ReqParams.IsHBRequest {
		// use mongodb config
		if adxEndpointConf, ok := extractor.GetHBAdxEndpointV2(mvutil.Config.Cloud, mvutil.Config.Region); ok {
			// mvutil.Logger.Runtime.Debugf("adxEndpoint used mongodb config: %#v", adxEndpointConf)
			backendCtx.ReqPath = adxEndpointConf.Endpoint
			backendCtx.Tmax = adxEndpointConf.Timeout
		}
		if len(hbreqctx.GetInstance().Cfg.AdxCfg.ServiceName) > 0 {
			backendCtx.ReqPath = hbreqctx.GetInstance().Cfg.AdxCfg.ServiceName
		}
		if hbreqctx.GetInstance().Cfg.ConsulCfg.Adx.Enable {
			ratio := extractor.GetUseConsulServicesV2Ratio(mvutil.Cloud(), mvutil.Region(), "hb_adx")
			if ratio > 0 && ratio > rand.Float64() {
				node := hbreqctx.GetInstance().Cfg.AdxConsulBuild.SelectNode()
				if node != nil && len(node.Host) > 0 {
					var addr = node.Host
					if strings.Index(addr, ":") == -1 && node.Port != 0 {
						addr += ":" + strconv.Itoa(node.Port)
					}
					metrics.IncCounterWithLabelValues(22, "hb_adx", mvutil.Zone(), node.Zone)
					backendCtx.ReqPath = "http://" + addr + "/hbrtb"
				}
			}
		}
		if os.Getenv("FORCE_ADX_ENDPIONT_FROM_ENV") == "1" && len(os.Getenv("ADX_SERVICE_NAME")) > 0 {
			backendCtx.ReqPath = os.Getenv("ADX_SERVICE_NAME")
		}
	}
	backendCtx.Method = backend.Method

	// 针对more_offer/appwall请求，调整请求超时时间
	if mvutil.IsAppwallOrMoreOffer(reqCtx.ReqParams.Param.AdType) {
		adnetConf, _ := extractor.GetADNET_SWITCHS()
		if moreOfferRequestPioneerTimeout, ok := adnetConf["moreOfferRequestPioneerTimeout"]; ok && moreOfferRequestPioneerTimeout > 0 {
			backend.Timeout = moreOfferRequestPioneerTimeout
		}
	}
	// 针对mp 单独设置的超时时间
	if mvutil.IsMpad(reqCtx.ReqParams.Param.RequestPath) {
		adnetConf, _ := extractor.GetADNET_SWITCHS()
		if mpRequestPioneerTimeout, ok := adnetConf["mpRequestPioneerTimeout"]; ok && mpRequestPioneerTimeout > 0 {
			backend.Timeout = mpRequestPioneerTimeout
		}
	}
	// 第三方的请求
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	err := backend.composeHttpRequest(reqCtx, backendCtx, req)
	if err != nil {
		mvutil.Logger.Runtime.Errorf("request_id=[%s] backend ID=[%d] name=[%s] error:%s",
			reqCtx.ReqParams.Param.RequestID, backend.ID, backend.Name, err.Error())
		metricData.FilterCode = mvconst.BuildReqError
		return metricData
	}
	// send request
	now := time.Now().UnixNano()
	strBackendID := strconv.Itoa(backend.ID)
	watcher.AddWatchValue("before_req_backend_"+strBackendID, float64(1))
	metricData.IsReqBackend = true
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)
	timeout := backend.Timeout
	if reqCtx.ReqParams.Param.UseDynamicTmax == "1" && extractor.GetTmaxABTestConf().Timeout["backend"] > 0 {
		timeout = int(extractor.GetTmaxABTestConf().Timeout["backend"])
		if reqCtx.ReqParams.Param.HBTmax > 0 && int(reqCtx.ReqParams.Param.HBTmax) < timeout {
			timeout = int(reqCtx.ReqParams.Param.HBTmax)
		}
	}
	err = utility.HttpClientApp().Do(timeout, req, resp)
	// err := backend.MAdxClient.Client.Do(httpReq.WithContext(ctx))
	watcher.AddWatchValue("after_req_backend_"+strBackendID, float64(1))
	backendCtx.Elapsed = int((time.Now().UnixNano() - now) / 1e6)
	watcher.AddAvgWatchValue("AvgD_"+strBackendID, float64(backendCtx.Elapsed))

	metrics.AddSummaryWithLabelValues(float64(backendCtx.Elapsed), 3, strconv.Itoa(backend.ID)) // metrics

	// 单独统计二阶段的响应分位
	if reqCtx.ReqParams.Param.IsHitRequestBidServer == 1 {
		metrics.AddSummaryWithLabelValues(float64(backendCtx.Elapsed), 37, "bid") // metrics
	}

	if err != nil {
		if strings.Contains(err.Error(), "timeout") {
			metrics.IncCounterDyWithValues("backend-tag-ab-error", strconv.Itoa(backend.ID),
				reqCtx.ReqParams.Param.ExtDataInit.UseDynamicTmax, strconv.Itoa(reqCtx.ReqParams.Param.ExtDataInit.TmaxABTestTag), "timeout")
			reqCtx.ReqParams.Param.ExtDataInit.ReqTimeout = 1
			if reqCtx.ReqParams.Param.IsHitRequestBidServer == 1 {
				// 二阶段bid请求timeout
				metrics.IncCounterWithLabelValues(32, BidHttpTimeoutError.Error())
			}
		}
		mvutil.Logger.Runtime.Errorf("request_id=[%s] backend ID=[%d] name=[%s] req_host=[%s] error:%s",
			reqCtx.ReqParams.Param.RequestID, backend.ID, backend.Name, backendCtx.ReqPath, err.Error())
		watcher.AddWatchValue("Do_Err_"+strBackendID, float64(1))
		metrics.IncCounterWithLabelValues(4, strconv.Itoa(backend.ID), strconv.Itoa(0), mvutil.GetAdTypeStr(reqCtx.ReqParams.Param.AdType))                  // metrics
		metrics.IncCounterWithLabelValues(24, strconv.Itoa(backend.ID), strconv.Itoa(resp.StatusCode()), mvutil.GetAdTypeStr(reqCtx.ReqParams.Param.AdType)) // metrics
		// 针对more_offer以及appwall监控
		if mvutil.IsAppwallOrMoreOffer(reqCtx.ReqParams.Param.AdType) {
			metrics.IncCounterWithLabelValues(23, strconv.Itoa(backend.ID), mvutil.GetAdTypeStr(reqCtx.ReqParams.Param.AdType))
		}
		// 针对二阶段bid的监控
		if reqCtx.ReqParams.Param.IsHitRequestBidServer == 1 {
			metrics.IncCounterWithLabelValues(32, BidHttpError.Error())
		}
		metricData.FilterCode = mvconst.HTTPDoError
		return metricData
	}

	metrics.IncCounterWithLabelValues(4, strconv.Itoa(backend.ID), strconv.Itoa(resp.StatusCode()), mvutil.GetAdTypeStr(reqCtx.ReqParams.Param.AdType))  // metrics
	metrics.IncCounterWithLabelValues(24, strconv.Itoa(backend.ID), strconv.Itoa(resp.StatusCode()), mvutil.GetAdTypeStr(reqCtx.ReqParams.Param.AdType)) // metrics
	// defer resp.Body.Close()
	// 对于直接请求pioneer的情况，返回的是9xx的情况。
	if resp.StatusCode() != 200 {
		// 单独记录pioneer返回的错误码
		if backend.ID == mvconst.Pioneer {
			reqCtx.ReqParams.Param.ExtDataInit.PioneerHttpCode = resp.StatusCode()
		}
		// 不为pioneer的情况，或者pioneer返回不为9xx或者2xx，则记录异常code以及retrun
		if backend.ID != mvconst.Pioneer || (!(resp.StatusCode() >= 900 && resp.StatusCode() <= 1000) && resp.StatusCode()/100 != 2) {
			watcher.AddWatchValue("Code_Not_200_"+strBackendID, float64(1))
			metricData.FilterCode = mvconst.HTTPStatusNotOK
			return metricData
		}
	}

	// read response data
	// var data []byte
	// data = resp.Body()
	// resp.Body.Close()
	backendCtx.RespData = resp.Body()
	if len(backendCtx.RespData) == 0 {
		mvutil.Logger.Runtime.Errorf("request_id=[%s] backend ID=[%d] name=[%s] error:%s",
			reqCtx.ReqParams.Param.RequestID, backend.ID, backend.Name, err.Error())
		metricData.FilterCode = mvconst.HTTPReadBodyError
		return metricData
	}

	// 解析第三方广告
	code, err := backend.parseHttpResponse(reqCtx, backendCtx)
	if err != nil {
		mvutil.Logger.Runtime.Warnf("request_id=[%s] ip=[%s] app_id=[%d] flow_tag_id=[%d] backend ID=[%d] name=[%s] error:%s",
			reqCtx.ReqParams.Param.RequestID, reqCtx.ReqParams.Param.ClientIP, reqCtx.ReqParams.Param.AppID, reqCtx.ReqParams.FlowTagID, backend.ID, backend.Name, err.Error())
		if code == ERR_NoAds {
			watcher.AddWatchValue("zero_ads_"+strBackendID, float64(1))
		}
		metricData.FilterCode = 5000 + code
		return metricData
	}

	// fill common part
	backendCtx.Ads.BackendId = int32(backend.ID)
	backendCtx.Ads.RequestKey = backendCtx.AdReqKeyName
	metricData.FilterCode = mvconst.BackendOK
	return metricData
}

// adserver don't filter
func (b Backend) filter(reqCtx *mvutil.ReqCtx, backendCtx *mvutil.BackendCtx) int {
	if reqCtx == nil || backendCtx == nil || reqCtx.ReqParams == nil {
		return mvconst.ParamInvalidate
	}
	// Region过滤广告后端
	if b.ID != mvconst.Mobvista && !mvutil.InStrArray("ALL", backendCtx.Region) &&
		!mvutil.InStrArray(reqCtx.ReqParams.Param.CountryCode, backendCtx.Region) {
		return mvconst.BackendRegionFilter
	}

	// content过滤
	if len(reqCtx.ReqParams.Param.VideoVersion) > 0 && backendCtx.Content != mvconst.VideoAdTypeNOLimit && reqCtx.ReqParams.Param.VideoAdType > 0 &&
		reqCtx.ReqParams.Param.VideoAdType != mvconst.VideoAdTypeNOLimit && backendCtx.Content != int(reqCtx.ReqParams.Param.VideoAdType) {
		return mvconst.BackendContentFilter
	}

	if b.specific != nil {
		return b.specific.filterBackend(reqCtx)
	}

	return mvconst.BackendOK
}

// 将请求参数合并/复合在一起？
func (b *Backend) composeHttpRequest(reqCtx *mvutil.ReqCtx, backendCtx *mvutil.BackendCtx, req *fasthttp.Request) error {
	if reqCtx == nil || backendCtx == nil {
		return errors.New("backend composeHttpRequest Params invalidate")
	}
	if b.specific != nil {
		return b.specific.composeHttpRequest(reqCtx, backendCtx, req)
	}
	return errors.New("backend doesn't implement composeHttpRequest")
}

func (b *Backend) parseHttpResponse(reqCtx *mvutil.ReqCtx, backendCtx *mvutil.BackendCtx) (int, error) {
	if reqCtx == nil || backendCtx == nil {
		return ERR_Param, errors.New("backend parseHttpResponse Params Invalidate")
	}
	if b.specific != nil {
		return b.specific.parseHttpResponse(reqCtx, backendCtx)
	}
	return ERR_Internal, errors.New("backend doesn't implement parseHttpResponse")
}

func IsFilterAdserverRequest(params *mvutil.Params) bool {
	conf := extractor.GetFilterAdserverRequestConf()
	adTypeCodeStr := strconv.Itoa(int(params.AdType))
	if rate, ok := conf[adTypeCodeStr]; ok {
		randVal := rand.Intn(100)
		if rate > randVal {
			return true
		}
	}
	return false
}

func (b *Backend) getCampaigns(reqCtx *mvutil.ReqCtx, backendCtx *mvutil.BackendCtx) *mvutil.BackendMetric {
	metricData := &mvutil.BackendMetric{BackendId: b.ID}
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			mvutil.Logger.Runtime.Errorf("backend ID=[%d] name=[%s] getCampaigns defer panic: %s",
				b.ID, b.Name, err)
			metrics.IncCounterWithLabelValues(25)
		}
	}()
	if reqCtx == nil || backendCtx == nil || reqCtx.ReqParams == nil {
		metricData.FilterCode = mvconst.ParamInvalidate
		return metricData
	}

	// 如果是低效流量，不请求adserver，也不兜底	TODO ?
	if reqCtx.ReqParams.Param.Scenario == mvconst.SCENARIO_OPENAPI && reqCtx.ReqParams.Param.IsLowFlowUnitReq {
		reqCtx.ReqParams.Param.FBFlag = 9
		metricData.FilterCode = mvconst.BackendLowFlowFilter
		return metricData
	}

	// 干掉直接请求as的流量，后续as模块会下掉。分ad_type控制	TODO as是啥？
	if IsFilterAdserverRequest(&reqCtx.ReqParams.Param) {
		metricData.FilterCode = mvconst.BackendLowVersionFilter
		return metricData
	}

	req, err := composeAdServerRequest(reqCtx, b.ID)
	if err != nil {
		mvutil.Logger.Runtime.Errorf("request_id=[%s] backend ID=[%d] name=[%s] composeAdServerRequest error:%s",
			reqCtx.ReqParams.Param.RequestID, b.ID, b.Name, err.Error())
		metricData.FilterCode = mvconst.BuildReqError
		return metricData
	}
	// // 针对素材三期单独设置超时时间
	timeout := b.Timeout
	adnConf, _ := extractor.GetADNET_SWITCHS()
	// debug 超时开关
	closeDebugTimeout, ok := adnConf["closeDebugTimeout"]
	if !ok || closeDebugTimeout != 1 {
		// 在debug及debugmode模式下，设置较长的超时时间
		if reqCtx.ReqParams.Param.Debug > 0 || reqCtx.ReqParams.Param.DebugMode {
			// 默认超时时间是10s
			timeout = 10000
			if reqCtx.ReqParams.Param.DebugModeTimeout > 0 {
				timeout = reqCtx.ReqParams.Param.DebugModeTimeout
			}
		}
	}

	adServerFlag := "_adsev"

	watcher.AddWatchValue("before_req"+adServerFlag, float64(1))
	now := time.Now().UnixNano()

	metricData.IsReqBackend = true
	var r *ad_server.QueryResult_
	// 通过thrift进行rpc调用，获得campaign
	r, err = b.AdServerClient.GetCampaigns(req, timeout)

	backendCtx.Elapsed = int((time.Now().UnixNano() - now) / 1e6)
	watcher.AddAvgWatchValue("avg_cost"+adServerFlag, float64(backendCtx.Elapsed))
	watcher.AddWatchValue("after_req"+adServerFlag, float64(1))

	metrics.AddSummaryWithLabelValues(float64(backendCtx.Elapsed), 3, strconv.Itoa(b.ID)) // metrics

	if err != nil {
		mvutil.Logger.Runtime.Errorf("request_id=[%s] backend ID=[%d] name=[%s] getCampaigns error:%s",
			reqCtx.ReqParams.Param.RequestID, b.ID, b.Name, err.Error())
		metricData.FilterCode = mvconst.BackendUnknownError
		metrics.IncCounterWithLabelValues(4, strconv.Itoa(b.ID), strconv.Itoa(0), mvutil.GetAdTypeStr(reqCtx.ReqParams.Param.AdType))  // metrics
		metrics.IncCounterWithLabelValues(24, strconv.Itoa(b.ID), strconv.Itoa(0), mvutil.GetAdTypeStr(reqCtx.ReqParams.Param.AdType)) // metrics
		// 针对more_offer以及appwall监控
		if mvutil.IsAppwallOrMoreOffer(reqCtx.ReqParams.Param.AdType) {
			metrics.IncCounterWithLabelValues(23, strconv.Itoa(b.ID), mvutil.GetAdTypeStr(reqCtx.ReqParams.Param.AdType))
		}
		return metricData
	}
	metrics.IncCounterWithLabelValues(4, strconv.Itoa(b.ID), strconv.Itoa(200), mvutil.GetAdTypeStr(reqCtx.ReqParams.Param.AdType))  // metrics
	metrics.IncCounterWithLabelValues(24, strconv.Itoa(b.ID), strconv.Itoa(200), mvutil.GetAdTypeStr(reqCtx.ReqParams.Param.AdType)) // metrics

	if (mvutil.Config.CommonConfig.LogConfig.OutputFullReqRes || reqCtx.ReqParams.Param.Debug > 0 || reqCtx.ReqParams.Param.DebugMode) && r != nil {
		rStr, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(r)
		if reqCtx.ReqParams.Param.Debug > 0 {
			reqCtx.ReqParams.DebugInfo += "ad_server_get_campaign_result@@@" + string(rStr) + "<\br>"
		} else if reqCtx.ReqParams.Param.DebugMode {
			var debugInfo []interface{}
			// debugInfo := r.GetDebugInfo()
			for _, info := range r.GetDebugInfo() {
				debugInfo = append(debugInfo, info)
			}
			r.DebugInfo = []string{}
			reqCtx.DebugModeInfo = append(reqCtx.DebugModeInfo, r)
			debugInfo = append(debugInfo, reqCtx.DebugModeInfo...)
			reqCtx.DebugModeInfo = debugInfo
		}
	}

	backendCtx.Ads.BackendId = int32(b.ID)
	backendCtx.Ads.RequestKey = backendCtx.AdReqKeyName
	_, err = doFillAdServerAds(reqCtx, backendCtx.Ads, r, false)
	if err != nil {
		mvutil.Logger.Runtime.Errorf("request_id=[%s] backend ID=[%d] name=[%s] fillAds error:%s",
			reqCtx.ReqParams.Param.RequestID, b.ID, b.Name, err.Error())
		metricData.FilterCode = mvconst.BackendFillError
		return metricData
	}
	if len(backendCtx.Ads.CampaignList) == 0 {
		if metricData.FilterCode == 0 {
			metricData.FilterCode = mvconst.BackendNoAds
		}
		watcher.AddWatchValue("zero_ads_"+strconv.Itoa(b.ID), float64(1))
	}

	if metricData.FilterCode == 0 {
		metricData.FilterCode = mvconst.BackendOK
	}
	return metricData
}

// URL 通过 pathQuery 生成完整的URL
func (b *Backend) URL(query string) string {
	host := b.HttpURL
	if len(query) == 0 {
		return fmt.Sprintf("%s%s", host, b.Path)
	}
	return fmt.Sprintf("%s%s?%s", host, b.Path, query)
}
func (b *Backend) TLSURL(query string) string {
	host := b.HttpsURL
	if len(query) == 0 {
		return fmt.Sprintf("%s%s", host, b.Path)
	}
	return fmt.Sprintf("%s%s?%s", host, b.Path, query)
}

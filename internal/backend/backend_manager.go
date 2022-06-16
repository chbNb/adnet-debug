package backend

import (
	"errors"
	"strconv"
	"sync"
	"time"

	"github.com/mae-pax/logger"
	"gitlab.mobvista.com/ADN/adnet/internal/clients"
	"gitlab.mobvista.com/ADN/adnet/internal/extractor"
	"gitlab.mobvista.com/ADN/adnet/internal/mvconst"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	"gitlab.mobvista.com/ADN/adnet/internal/watcher"
)

type BackendManage struct {
	Backends map[int]*Backend
}

var BackendManager *BackendManage

func init() {
	BackendManager = NewBackendManage()
}

// NewBackendManage 新建一个BackendManage
func NewBackendManage() *BackendManage {
	manager := &BackendManage{}
	manager.Backends = make(map[int]*Backend)
	return manager
}

// AddBackend 添加广告后端
func (manager *BackendManage) AddBackend(serviceDetail *mvutil.ServiceDetail, logger *logger.Log, consulAdxLog *logger.Log, consulWatchLog *logger.Log) error {
	backend := NewBackend(serviceDetail)
	// 不想用map管理
	switch serviceDetail.ID {
	case mvconst.Mobvista:
		adServerClient, err := clients.NewAdServerClient(serviceDetail, logger)
		if err != nil {
			return err
		}
		backend.AdServerClient = adServerClient
		backend.specific = &AdServer{
			Backend{AdServerClient: adServerClient},
		}
	case mvconst.MAdx:
		adxClient, err := clients.NewMAdxClient(serviceDetail, logger, consulAdxLog, consulWatchLog)
		if err != nil {
			return err
		}
		backend.MAdxClient = adxClient
		backend.specific = &MAdxBackend{
			Backend{MAdxClient: adxClient},
		}
	case mvconst.Pioneer:
		pioneerClient, err := clients.NewPioneerClient(serviceDetail, logger)
		if err != nil {
			return err
		}
		backend.PioneerClient = pioneerClient
		backend.specific = &PioneerBackend{
			Backend{PioneerClient: pioneerClient},
		}
	}
	manager.Backends[serviceDetail.ID] = backend
	return nil
}

type AdRequestTask struct {
	ReqCtx  *mvutil.ReqCtx
	Backend *Backend
	Tid     int
	TName   string
	// err     error
}

func (manager *BackendManage) GetAds(reqCtx *mvutil.ReqCtx) (map[int]*mvutil.BackendMetric, error) {
	if reqCtx == nil {
		return nil, errors.New("BackendManage GetAds Params invalidate")
	}
	maxTimeout := 0
	backendData := make(map[int]*mvutil.BackendMetric, len(reqCtx.Backends))
	metricData := make([]*mvutil.BackendMetric, 0, len(reqCtx.Backends))
	var endChan chan struct{} = make(chan struct{}, 1)
	var wg sync.WaitGroup
	// 遍历所有Backend     注意，以前是发给头条等ad_server,现在改成只发给自己的adx
	for id, backendCtx := range reqCtx.Backends {
		backend, ok := manager.Backends[id]
		if !ok {
			mvutil.Logger.Runtime.Warnf("not found backend instance!!!")
			metricData = append(metricData, &mvutil.BackendMetric{FilterCode: mvconst.BackendNotInstance, BackendId: id})
			continue
		}
		// mobvista 不需要过滤
		if id != mvconst.Mobvista && id != mvconst.Pioneer {
			if code := backend.filter(reqCtx, backendCtx); code != mvconst.BackendOK {
				watcher.AddWatchValue("backend_filter_"+strconv.Itoa(id), float64(1))
				metricData = append(metricData, &mvutil.BackendMetric{FilterCode: code, BackendId: id})
				continue
			}
			// cap block filter
			if reqCtx.ReqParams.Param.IsBlockByImpCap {
				if !reqCtx.ReqParams.Param.OnlyRequestThirdDsp || id != mvconst.MAdx {
					watcher.AddWatchValue("imp_cap_block_"+strconv.Itoa(id), float64(1))
					metricData = append(metricData, &mvutil.BackendMetric{FilterCode: mvconst.ImpCapBlock, BackendId: id})
					continue
				}
			}
		}
		if reqCtx.ReqParams.Param.Debug <= 0 && !reqCtx.ReqParams.Param.DebugMode && // 非debug模式
			(id == mvconst.Mobvista || id == mvconst.Pioneer) && reqCtx.ReqParams.Param.IsBlockByImpCap { // 针对AS
			watcher.AddWatchValue("imp_cap_block_"+strconv.Itoa(id), float64(1))
			metricData = append(metricData, &mvutil.BackendMetric{FilterCode: mvconst.ImpCapBlock, BackendId: id})
			continue
		}
		// 判断rv req_type是否为3，为3则做aabtest，对于req_type为3且命中不返回的情况，不请求adserver，也不做兜底
		// 19-11-26改为针对playrix的req_type=3的流量进行下毒，不召回广告。
		if reqTypeFilter(reqCtx) {
			metricData = append(metricData, &mvutil.BackendMetric{FilterCode: mvconst.BackendReqTypeAABTestFilter, BackendId: id})
			continue
		}

		if backend.Timeout > maxTimeout {
			maxTimeout = backend.Timeout
		}
		wg.Add(1)

		go func(reqCtx *mvutil.ReqCtx, backendCtx *mvutil.BackendCtx, backendId int, backend *Backend) {
			switch backendId {
			case mvconst.Mobvista:
				metricData = append(metricData, backend.getCampaigns(reqCtx, backendCtx))
			default:
				metricData = append(metricData, backend.dispatchProxyRequest(reqCtx, backendCtx))
			}
			wg.Done()
		}(reqCtx, backendCtx, id, backend)
	}
	// 做超时的处理
	go func(eChan chan struct{}) {
		wg.Wait()
		eChan <- struct{}{}
	}(endChan)

	var timeoutErr error

	select {
	case <-endChan:
		timeoutErr = nil
	case <-time.After(time.Duration(maxTimeout+5) * time.Millisecond):
		// cancel()
		watcher.AddWatchValue("timeout_backend_req", float64(1))
		timeoutErr = errors.New("AdnTimeout occured")
	}

	// adxFailure := false
	for _, metric := range metricData {
		if metric != nil && metric.BackendId > 0 {
			backendData[metric.BackendId] = metric
		}
	}
	return backendData, timeoutErr
}

// 由原来abtest改为下毒逻辑
func reqTypeFilter(r *mvutil.ReqCtx) bool {
	// mvutil.Logger.Runtime.Info("reqTypeFilter")
	if r.ReqParams.Param.AdType != mvconst.ADTypeRewardVideo || r.ReqParams.Param.ReqType != "3" {
		return false
	}
	reqTypeAABTestConf, _ := extractor.GetREQ_TYPE_AAB_TEST_CONFIG()
	mvutil.Logger.Runtime.Infof("request publisher_id: %d, REQ_TYPE_AAB_TEST_CONFIG: %#v", r.ReqParams.Param.PublisherID, reqTypeAABTestConf)
	if !reqTypeAABTestConf.Status {
		return false
	}

	// 为了屏蔽playrix的type 3的请求。
	if len(reqTypeAABTestConf.PublisherIds) > 0 && !mvutil.InInt64Arr(r.ReqParams.Param.PublisherID, reqTypeAABTestConf.PublisherIds) {
		return false
	}

	if len(reqTypeAABTestConf.AppIds) > 0 && !mvutil.InInt64Arr(r.ReqParams.Param.AppID, reqTypeAABTestConf.AppIds) {
		return false
	}

	if len(reqTypeAABTestConf.UnitIds) > 0 && !mvutil.InInt64Arr(r.ReqParams.Param.UnitID, reqTypeAABTestConf.UnitIds) {
		return false
	}
	r.ReqParams.Param.ReqTypeAABTest = 3
	return true
}

func ExcludeDisplayPackageABTest(r *mvutil.ReqCtx) {
	abTestConf, _ := extractor.GetExcludeDisplayPackageABTest()
	if !abTestConf.Status {
		return
	}

	// 针对SDK 实验
	if r.ReqParams.Param.Scenario != mvconst.SCENARIO_OPENAPI ||
		r.ReqParams.Param.RequestType != mvconst.REQUEST_TYPE_OPENAPI_V3 {
		return
	}

	// 空设备流量不进入实验
	if mvutil.IsDevidEmpty(&r.ReqParams.Param) {
		return
	}

	if mvutil.InInt64Arr(r.ReqParams.Param.PublisherID, abTestConf.BlackPublisherIds) {
		return
	}

	if mvutil.InInt64Arr(r.ReqParams.Param.AppID, abTestConf.BlackAppIds) {
		return
	}

	if mvutil.InInt64Arr(r.ReqParams.Param.UnitID, abTestConf.BlackUnitIds) {
		return
	}

	if len(abTestConf.WhitePublisherIds) > 0 &&
		!mvutil.InInt64Arr(r.ReqParams.Param.PublisherID, abTestConf.WhitePublisherIds) {
		return
	}

	if len(abTestConf.WhiteAppIds) > 0 &&
		!mvutil.InInt64Arr(r.ReqParams.Param.AppID, abTestConf.WhiteAppIds) {
		return
	}

	if len(abTestConf.WhiteUnitIds) > 0 &&
		!mvutil.InInt64Arr(r.ReqParams.Param.UnitID, abTestConf.WhiteUnitIds) {
		return
	}

	rate := abTestConf.Rate
	if len(rate) == 0 {
		return
	}

	randVal := mvutil.GetRandConsiderZero(r.ReqParams.Param.GAID, r.ReqParams.Param.IDFA, mvconst.SALT_DISPLAY_CAMPAIGN_ABTEST, 10000)
	aabTestVal, randOK := mvutil.RandByRateInMap(rate, randVal)
	if !randOK {
		return
	}

	if len(r.ReqParams.Param.DisplayCamIds) > 0 {
		aabTestVal += 3 // 三组测试 base 为 1 2 3，带有display的为4 5 6
	}

	if aabTestVal == 6 {
		if r.ReqParams.Param.ExcludePackageNames == nil {
			r.ReqParams.Param.ExcludePackageNames = make(map[string]bool)
		}

		for _, cid := range r.ReqParams.Param.DisplayCamIds {
			if campaign, find := extractor.GetCampaignInfo(cid); find {
				if campaign.PackageName != "" {
					r.ReqParams.Param.ExcludePackageNames[campaign.PackageName] = true
				}
			}
		}
	}

	r.ReqParams.Param.DisplayCampaignABTest = aabTestVal
}

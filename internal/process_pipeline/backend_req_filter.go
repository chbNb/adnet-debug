package process_pipeline

import (
	"errors"
	"strconv"
	"time"

	jsoniter "github.com/json-iterator/go"
	"gitlab.mobvista.com/ADN/adnet/internal/backend"
	"gitlab.mobvista.com/ADN/adnet/internal/corsair_proto"
	"gitlab.mobvista.com/ADN/adnet/internal/mvconst"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	"gitlab.mobvista.com/ADN/adnet/internal/watcher"
)

type BackendReqFilter struct {
}

func (brf *BackendReqFilter) Process(data interface{}) (interface{}, error) {
	in, ok := data.(*mvutil.ReqCtx)
	if !ok {
		return nil, errors.New("BackendReqFilter input type should be *mvutil.ReqCtx")
	}

	if in.ReqParams.Param.Debug > 0 {
		var json = jsoniter.ConfigCompatibleWithStandardLibrary
		resStr, _ := json.Marshal(in)
		in.ReqParams.DebugInfo += "before_backend_req@@@" + string(resStr) + "<\br>"
	}

	backend.ExcludeDisplayPackageABTest(in)

	// 过了参数过滤，请求广告后端的请求数据
	watcher.AddWatchValue("ad_request", float64(1))
	now := time.Now().UnixNano()
	// 广告请求
	backendMetric, err := backend.BackendManager.GetAds(in)
	if err != nil {
		// 记录下backendid以及错误原因
		// 记录到ext_data中
		for id, _ := range in.Backends {
			in.ReqParams.Param.GetAdsErrbackendList = append(in.ReqParams.Param.GetAdsErrbackendList, id)
		}
		in.ReqParams.Param.GetAdsErr = err.Error()
		// add log
		mvutil.Logger.Runtime.Warnf("request_id=[%s] backend.BackendManager.GetAds() error: %s", in.ReqParams.Param.RequestID, err.Error())
	}

	// fill ads
	adCount := int64(0)
	var backendData []string
	var res corsair_proto.QueryResult_
	res.LogId = in.ReqParams.Param.RequestID
	res.FlowTagId = int32(in.FlowTagID)
	res.RandValue = int32(in.RandValue)
	res.AdBackendConfig = GetBackendConfig(in)
	// append by backend priority order
	for _, id := range in.OrderedBackends {
		if metric, exists := backendMetric[id]; exists {
			indexStr := strconv.Itoa(id)
			// backend改为为mobvista
			if mvutil.IsRequestPioneerDirectly(&in.ReqParams.Param) {
				indexStr = "1"
			}
			if metric != nil && metric.IsReqBackend {
				in.ReqParams.Param.ReqBackends = append(in.ReqParams.Param.ReqBackends, indexStr)
			}
			if metric != nil {
				in.ReqParams.Param.BackendReject = append(in.ReqParams.Param.BackendReject, indexStr+":"+strconv.Itoa(metric.FilterCode))
			}
		}

		backendCtx, ok := in.Backends[id]
		if ok {
			RenderAdsByBackend(in, id, backendCtx, &adCount, backendData, &res)
		}
	}

	if len(res.AdsByBackend) == 0 {
		watcher.AddWatchValue("zero_ads", float64(1))
	}
	in.Elapsed = int((time.Now().UnixNano() - now) / 1e6) // ms
	watcher.AddAvgWatchValue("avg_delay", float64(in.Elapsed))

	if mvutil.Config.CommonConfig.LogConfig.OutputFullReqRes || in.ReqParams.Param.Debug > 0 {
		var json = jsoniter.ConfigCompatibleWithStandardLibrary
		resStr, _ := json.Marshal(&res)
		if in.ReqParams.Param.Debug > 0 {
			in.ReqParams.DebugInfo += "before_output_render_result@@@" + string(resStr) + "<\br>"
		} else {
			mvutil.Logger.Runtime.Debugf("request_id=[%s] response data=%s", in.ReqParams.Param.RequestID, resStr)
		}
	}
	in.Result = &res
	return in, nil
}

func RenderAdsByBackend(in *mvutil.ReqCtx, id int, backendCtx *mvutil.BackendCtx, adCount *int64, backendData []string, res *corsair_proto.QueryResult_) {
	if in.FlowTagID > mvconst.FlowTagDefault {
		// ad filter
		FilterAds(in, id, backendCtx.Ads, adCount, &backendData)
	}

	if len(backendCtx.Ads.CampaignList) > 0 {
		// 对于没有切到大模板的情况下，interstitial只要一个广告
		if in.ReqParams.Param.AdType == mvconst.ADTypeInterstitialVideo {
			if !in.ReqParams.Param.BigTemplateFlag {
				backendCtx.Ads.CampaignList = backendCtx.Ads.CampaignList[:1]
			}
			res.AdsByBackend = append(res.AdsByBackend, backendCtx.Ads)
			return
		}
	}
	res.AdsByBackend = append(res.AdsByBackend, backendCtx.Ads)
	if backendCtx.AsAds != nil && len(backendCtx.AsAds.CampaignList) > 0 {
		res.AdsByBackend = append(res.AdsByBackend, backendCtx.AsAds)
	}
}

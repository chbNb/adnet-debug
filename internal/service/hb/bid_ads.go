package hb

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"gitlab.mobvista.com/ADN/exporter/metrics"

	pf "github.com/easierway/pipefiter_framework/pipefilter"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/constant"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/filter"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/output"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/params"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/req_context"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	pp "gitlab.mobvista.com/ADN/adnet/internal/process_pipeline"
	"gitlab.mobvista.com/ADN/adnet/internal/service/hb/bid_filters"
)

type BidAdsHandler struct {
	Pipeline pf.Filter
}

var bidAdsFilters = []pf.Filter{
	&pp.MVReqparamTransFilter{},
	&pp.RequestParamFilter{},
	&bid_filters.AreaTargetFilter{},
	&pp.UaParserFilter{},
	&pp.OnlineParamFilter{},
	&pp.ParamRenderFilter{},
	&pp.ReplaceBrandFilter{},
	&pp.ReqValidateFilter{},
	&pp.TrafficAllotFilter{},
	&pp.BackendReqFilter{},
	&bid_filters.BidAdsFormatOutputFilter{},
	&pp.AdRenderFilter{},
}

func CreateBidAdsHandler() *BidAdsHandler {
	bidAdsPipeline := &WallTimePipeline{
		Name:        "bid_ads",
		Filters:     &bidAdsFilters,
		TimeElapsed: make([]AtomicInt, len(bidAdsFilters)),
	}
	bidAdsPipeline.Show()
	return &BidAdsHandler{
		Pipeline: bidAdsPipeline,
	}
}

func (bh *BidAdsHandler) ServeHTTP(c http.ResponseWriter, req *http.Request) {
	var bidAdsResp params.MobvistaResult
	bidAdsResp.Status = 204
	bidAdsResp.Msg = "has no bid"
	respByte, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(bidAdsResp)
	// 设置跨域参数
	c.Header().Set("Access-Control-Allow-Origin", "*")
	c.Header().Set("Content-Type", "application/json;charset=UTF-8")
	var isServerBiddingReq bool

	// catch panic
	defer func() {
		if r := recover(); r != nil {
			msg := fmt.Sprintf("BidAdsHandler Process Panic: %v, stack= %s", r, string(debug.Stack()))
			fmt.Println(msg)
			req_context.GetInstance().MLogs.Runtime.Error(msg)
			metrics.IncCounterWithLabelValues(25)
			c.WriteHeader(http.StatusNoContent)
			io.WriteString(c, string(respByte))
		}
	}()

	now := time.Now().UnixNano()
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		//req_context.GetInstance().MLogs.ReqMonitor.Warnf("BidAdsHandler Process ReadAllBody: %s", err.Error())
		WriterData(now, c, []byte(err.Error()), constant.BidAdsCost)
		return
	}
	req.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	// business process
	respData, err := bh.Pipeline.Process(req)
	if err != nil {
		wrapErrMsg, rawErrMsg := filter.FormatErrorMessage(err.Error())
		method := strings.ToUpper(req.Method)
		url := req.URL
		uri := url.RequestURI()
		//var b strings.Builder
		if method == "POST" {
			isServerBiddingReq = true
			//b.Write(body)
		}
		// var errCode int
		// var errMsg string
		switch errors.Cause(err).(type) {
		case *filter.PipeLineHandlerError:
			// filterCode, _ := filter.UnmarshalMessage(rawErrMsg)
			// errCode = filterCode.Int()
			// errMsg = filterCode.String()
			respByte = []byte(rawErrMsg)
		default:
			if filterCode, ok := errors.Cause(err).(filter.FilterCode); ok {
				// errCode = filterCode.Int()
				// errMsg = filterCode.String()
				respByte = []byte(filterCode.Error())
			} else {
				// errCode = filter.BidNoError.Int()
				// errMsg = rawErrMsg
			}
		}
		if bidReq, ok := respData.(*mvutil.ReqCtx); ok {
			req_context.GetInstance().MLogs.Bid.Info(output.FormatBidLog(bidReq, filter.BidNoError.Int(), err.Error()))
		}
		req_context.GetInstance().MLogs.ReqMonitor.Warnf("BidAdsHandler Process Error: %s, Raw Error: %s, uri: %s, method: %s", wrapErrMsg, rawErrMsg, uri, method)
		// WriterDataWithHttpStatusCode(now, c, respByte, constant.BidAdsCost)
		WriterDataWithHttpStatusCode(now, c, http.StatusNoContent, respByte, constant.BidAdsCost, isServerBiddingReq)
		return
	}

	if str, ok := respData.(*[]byte); ok {
		// WriterDataWithHttpStatusCode(now, c, *str, constant.BidAdsCost)
		WriterDataWithHttpStatusCode(now, c, http.StatusOK, *str, constant.BidAdsCost, isServerBiddingReq)
	} else {
		// WriterDataWithHttpStatusCode(now, c, respByte, constant.BidAdsCost)
		if bidReq, ok := respData.(*mvutil.ReqCtx); ok {
			req_context.GetInstance().MLogs.Bid.Info(output.FormatBidLog(bidReq, filter.BidNoError.Int(), filter.BidNoError.String()))
		}
		WriterDataWithHttpStatusCode(now, c, http.StatusNoContent, respByte, constant.BidAdsCost, isServerBiddingReq)
	}
	return
}

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

	pf "github.com/easierway/pipefiter_framework/pipefilter"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"gitlab.mobvista.com/ADN/adnet/internal/errorcode"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/constant"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/filter"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/output"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/params"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/req_context"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	pp "gitlab.mobvista.com/ADN/adnet/internal/process_pipeline"
	"gitlab.mobvista.com/ADN/adnet/internal/service/hb/bid_filters"
	"gitlab.mobvista.com/ADN/adnet/internal/watcher"
	"gitlab.mobvista.com/ADN/exporter/metrics"
)

type BidHandler struct {
	Pipeline pf.Filter
}

// 若是修改这个bidFilters的元素时，请把它下面的labels对应位置也调整下，否则会导致argus数据面板展示混淆
var bidFilters = []pf.Filter{
	&pp.MVReqparamTransFilter{},
	&pp.RequestParamFilter{},
	&bid_filters.AreaTargetFilter{},
	&pp.UaParserFilter{},
	&pp.ParamRenderFilter{},
	&pp.ReplaceBrandFilter{},
	&pp.MappingServerFilter{},
	&pp.ReqValidateFilter{},
	&pp.TrafficAllotFilter{},
	&pp.BackendReqFilter{},
	&bid_filters.FormatOutputFilter{},
	// &pp.AdRenderFilter{},
}

// 量最大(5k+ qps), 可作为pipeline的样本采集，注意与上面面的bidFilters顺序保持一致
var labels = []string{
	"MVReqparamTransFilter",
	"RequestParamFilter",
	"AreaTargetFilter",
	"UaParserFilter",
	"ParamRenderFilter",
	"ReplaceBrandFilter",
	"MappingServerFilter",
	"ReqValidateFilter",
	"TrafficAllotFilter",
	"BackendReqFilter",
	"FormatOutputFilter",
	// "AdRenderFilter",
}

func CreateBidHandler() *BidHandler {
	bidPipeline := &WallTimePipeline{
		Name:        Bid,
		Filters:     &bidFilters,
		TimeElapsed: make([]AtomicInt, len(bidFilters)),
		Labels:      labels,
	}
	bidPipeline.Show()
	return &BidHandler{
		Pipeline: bidPipeline,
	}
}

func (bh *BidHandler) ServeHTTP(c http.ResponseWriter, req *http.Request) {
	var bidResp params.BidResp
	bidResp.Status = 204
	bidResp.Msg = "has no bid"
	respByte, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(bidResp)
	// 设置跨域参数
	c.Header().Set("Access-Control-Allow-Origin", "*")
	c.Header().Set("Content-Type", "application/json;charset=UTF-8")

	var isServerBiddingReq bool
	method := strings.ToUpper(req.Method)
	url := req.URL
	uri := url.RequestURI()
	if method == "POST" {
		isServerBiddingReq = true
	}

	// catch panic
	defer func() {
		if r := recover(); r != nil {
			msg := fmt.Sprintf("BidHandler Process Panic: %v, stack= %s", r, string(debug.Stack()))
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
		//req_context.GetInstance().MLogs.ReqMonitor.Warnf("BidHandler Process ReadAllBody: %s", err.Error())
		WriterDataWithHttpStatusCode(now, c, http.StatusNoContent, []byte(err.Error()), constant.BidCost, isServerBiddingReq)
		return
	}
	req.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	// business process
	respData, err := bh.Pipeline.Process(req)
	if err != nil {
		wrapErrMsg, rawErrMsg := filter.FormatErrorMessage(err.Error())
		var errCode int
		var errMsg string
		switch errors.Cause(err).(type) {
		case *filter.PipeLineHandlerError:
			filterCode, _ := filter.UnmarshalMessage(rawErrMsg)
			errCode = filterCode.Int()
			errMsg = filterCode.String()
			respByte = []byte(rawErrMsg)
		default:
			if filterCode, ok := errors.Cause(err).(filter.FilterCode); ok {
				errCode = filterCode.Int()
				errMsg = filterCode.String()
				respByte = []byte(filterCode.Error())
			} else if filterCode, ok := errors.Cause(err).(errorcode.AdnetCode); ok {
				errCode = int(filterCode)
				errMsg = filterCode.String()
				respByte = []byte(filterCode.Error())
			} else {
				errCode = filter.BidNoError.Int()
				errMsg = err.Error()
			}
		}

		switch respData.(type) {
		case *mvutil.ReqCtx:
			bidReq, _ := respData.(*mvutil.ReqCtx)
			req_context.GetInstance().MLogs.Bid.Info(output.FormatBidLog(bidReq, errCode, errMsg))
		case *mvutil.RequestParams:
			bidReqParams, _ := respData.(*mvutil.RequestParams)
			bidReq := new(mvutil.ReqCtx)
			bidReq.ReqParams = bidReqParams
			req_context.GetInstance().MLogs.Bid.Info(output.FormatBidLog(bidReq, errCode, errMsg))
			//default:
			//req_context.GetInstance().MLogs.ReqMonitor.Warnf("BidHandler Process Error, and unknown respData type for the bid log.")
		}
		req_context.GetInstance().MLogs.ReqMonitor.Warnf("BidHandler Process Error: %s, Raw Error: %s, uri: %s, method: %s", wrapErrMsg, rawErrMsg, uri, method)
		//req_context.GetInstance().MLogs.Runtime.Warnf("BidHandler Process Error: %s, Raw Error: %s, uri: %s, method: %s", wrapErrMsg, rawErrMsg, uri, method) WriterDataWithHttpStatusCode(now, c, http.StatusNoContent, respByte, constant.BidCost, isServerBiddingReq)
		return
	}

	if _, ok := respData.(*[]byte); !ok {
		if bidReq, ok := respData.(*mvutil.ReqCtx); ok {
			req_context.GetInstance().MLogs.Bid.Info(output.FormatBidLog(bidReq, filter.BidNoError.Int(), filter.BidNoError.String()))
		}
		//req_context.GetInstance().MLogs.ReqMonitor.Warnf("Process no error")
		WriterDataWithHttpStatusCode(now, c, http.StatusNoContent, respByte, constant.BidCost, isServerBiddingReq)
		return
	}
	WriterDataWithHttpStatusCode(now, c, http.StatusOK, *respData.(*[]byte), constant.BidCost, isServerBiddingReq)
	return
}

func WriterData(now int64, c http.ResponseWriter, data []byte, watchTag string) {
	// c.Header().Set("Content-length", strconv.Itoa(len(data)))
	c.Write(data)
	watcher.AddAvgWatchValue(watchTag, float64((time.Now().UnixNano()-now)/1e6))
}

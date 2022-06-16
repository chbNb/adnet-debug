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

type MopubBidHandler struct {
	Pipeline pf.Filter
}

var bopubBidHandlers = []pf.Filter{
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

// func init() {
// 	httpExtractFilter := &base_filters.HttpExtractFilter{}
// 	bopubBidHandlers = []pf.Filter{&base_filters.MopubHttpExtractFilter{BaseHttpExtractFilter: httpExtractFilter}, &bid_filters.ReqParamFilter{},
// 		&bid_filters.AreaTargetFilter{}, &bid_filters.UserAgentDataFilter{},
// 		&bid_filters.ReplaceBrandModelFilter{}, &bid_filters.RenderCoreDataFilter{},
// 		&bid_filters.PriceFactorFilter{}, &bid_filters.RankerInfoFilter{},
// 		&bid_filters.BuildAsRequestFilter{}, &bid_filters.BidAdxFilter{}, &bid_filters.BidCacheFilter{},
// 		&bid_filters.MopubFormatOutputFilter{}}
// }

func CreateMopubBidHandler() *MopubBidHandler {
	return &MopubBidHandler{
		Pipeline: &WallTimePipeline{
			Name:        Bid,
			Filters:     &bopubBidHandlers,
			TimeElapsed: make([]AtomicInt, len(bopubBidHandlers)),
		},
	}
}

func (bh *MopubBidHandler) ServeHTTP(c http.ResponseWriter, req *http.Request) {
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
		//var b strings.Builder
		//if isServerBiddingReq {
		//	b.Write(body)
		//}
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
			} else {
				errCode = filter.BidNoError.Int()
				errMsg = err.Error()
			}
		}
		//if bidReq, ok := respData.(*mvutil.ReqCtx); ok {
		//	// if bidReq.Nbr > 0 {
		//	// 	bidReq.RejectData = strconv.Itoa(constant.MAdx) + ":" + strconv.Itoa(bidReq.BidRejectCode)
		//	// 	// req_context.GetInstance().MLogs.MidwayReq.Info(output.FormatMidwayLog(bidReq))
		//	// }
		//	req_context.GetInstance().MLogs.Bid.Info(output.FormatBidLog(bidReq, errCode, errMsg))
		//}
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
			//req_context.GetInstance().MLogs.ReqMonitor.Warnf("Mopub BidHandler Process Error, and unknown respData type for the bid log.")
		}
		req_context.GetInstance().MLogs.ReqMonitor.Warnf("Mopub BidHandler Process Error: %s, Raw Error: %s, uri: %s, method: %s", wrapErrMsg, rawErrMsg, uri, method)
		WriterDataWithHttpStatusCode(now, c, http.StatusNoContent, respByte, constant.BidCost, isServerBiddingReq)
		return
	}

	if _, ok := respData.(*[]byte); !ok {
		if bidReq, ok := respData.(*mvutil.ReqCtx); ok {
			req_context.GetInstance().MLogs.Bid.Info(output.FormatBidLog(bidReq, filter.BidNoError.Int(), filter.BidNoError.String()))
		}
		//req_context.GetInstance().MLogs.ReqMonitor.Warnf("Mopub Process no error")
		WriterDataWithHttpStatusCode(now, c, http.StatusNoContent, respByte, constant.BidCost, isServerBiddingReq)
		return
	}
	WriterDataWithHttpStatusCode(now, c, http.StatusOK, *respData.(*[]byte), constant.BidCost, isServerBiddingReq)
	return
}

func WriterDataWithHttpStatusCode(now int64, c http.ResponseWriter, statusCode int, data []byte, watchTag string, isServerBiddingReq bool) {
	if isServerBiddingReq {
		c.WriteHeader(statusCode)
		if statusCode == http.StatusOK {
			c.Write(data)
		}
	} else {
		c.Write(data)
	}
	watcher.AddAvgWatchValue(watchTag, float64((time.Now().UnixNano()-now)/1e6))
}

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
	"gitlab.mobvista.com/ADN/adnet/internal/service/hb/load_filters"
	"gitlab.mobvista.com/ADN/exporter/metrics"
)

type LoadHandler struct {
	Pipeline pf.Filter
}

var loadFilters = []pf.Filter{
	&pp.MVReqparamTransFilter{},
	&pp.ServiceDegradeFilter{},
	&load_filters.ReqParamFilter{},
	&pp.AdRenderFilter{},
}

// func init() {
// 	loadFilters = []pf.Filter{&base_filters.HttpExtractFilter{}, &load_filters.ReqParamFilter{},
// 		&load_filters.QueryBidFilter{}, &load_filters.RenderAsAdFilter{}, &load_filters.RenderDspAdFilter{},
// 		&load_filters.RenderTrackingLinkFilter{}, &load_filters.FormatOutputFilter{}}
// }

func CreateLoadHandler() *LoadHandler {
	return &LoadHandler{
		Pipeline: &WallTimePipeline{
			Name:        Load,
			Filters:     &loadFilters,
			TimeElapsed: make([]AtomicInt, len(loadFilters)),
		},
	}
}

func (lh *LoadHandler) ServeHTTP(c http.ResponseWriter, req *http.Request) {
	var loadResp params.MobvistaResult
	loadResp.Status = 204
	loadResp.Msg = "load no ad"
	respByte, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(loadResp)

	// 设置跨域参数
	c.Header().Set("Access-Control-Allow-Origin", "*")
	c.Header().Set("Content-Type", "application/json;charset=UTF-8")

	// catch panic
	defer func() {
		if r := recover(); r != nil {
			msg := fmt.Sprintf("LoadHandler Process panic %v, stack=%s", r, debug.Stack())
			fmt.Println(msg)
			req_context.GetInstance().MLogs.Runtime.Error(msg)
			metrics.IncCounterWithLabelValues(25)
			io.WriteString(c, string(respByte))
		}
	}()
	now := time.Now().UnixNano()
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		WriterData(now, c, []byte(err.Error()), constant.LoadCost)
		return
	}
	req.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	// business process
	respData, err := lh.Pipeline.Process(req)
	if err != nil {
		wrapErrMsg, rawErrMsg := filter.FormatErrorMessage(err.Error())
		method := strings.ToUpper(req.Method)
		url := req.URL
		uri := url.RequestURI()
		//var b strings.Builder
		//if method == "POST" {
		//	b.Write(body)
		//}
		var errCode int
		var errMsg string
		switch errors.Cause(err).(type) {
		case *filter.PipeLineHandlerError:
			filterCode, _ := filter.UnmarshalMessage(rawErrMsg)
			errCode = filterCode.Int()
			errMsg = filterCode.String()
			// respByte = []byte(rawErrMsg)
		default:
			if filterCode, ok := errors.Cause(err).(filter.FilterCode); ok {
				errCode = filterCode.Int()
				errMsg = filterCode.String()
				// respByte = []byte(filterCode.Error())
			} else {
				errCode = -1
				errMsg = err.Error()
			}
		}
		//if loadReq, ok := respData.(*mvutil.ReqCtx); ok {
		//	//loadReq.ReqParams.LoadRejectCode = errCode
		//	req_context.GetInstance().MLogs.Load.Info(output.FormatLoadLog(loadReq, errCode, errMsg))
		//	// if loadReq.RejectCode == filter.QueryBidError {
		//	// 	mobvistaResult.Msg = "load token param expired"
		//	// 	respByte, _ = jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(mobvistaResult)
		//	// }
		//} else {
		//	loadRequest := new(mvutil.ReqCtx)
		//	rawQuery := mvutil.RequestQueryMap(req.URL.Query())
		//	if strings.ToUpper(req.Method) == "POST" {
		//		req.ParseForm()
		//		for k, v := range req.PostForm {
		//			rawQuery[k] = v
		//		}
		//	}
		//	loadRequest.ReqParams = &mvutil.RequestParams{}
		//	appID, _ := rawQuery.GetInt64("app_id")
		//	loadRequest.ReqParams.Param.AppID = appID
		//	unitID, _ := rawQuery.GetInt64("unit_id", 0)
		//	loadRequest.ReqParams.Param.UnitID = unitID
		//	token, _ := rawQuery.GetString("token", true, "")
		//	loadRequest.ReqParams.Token = token
		//	loadRequest.ReqParams.LoadRejectCode = errCode
		//	req_context.GetInstance().MLogs.Load.Info(output.FormatLoadLog(loadRequest, errCode, errMsg))
		//	req_context.GetInstance().MLogs.ReqMonitor.Warnf("LoadHandler Process respData isn't conversion to mvutil.ReqCtx")
		//}
		switch respData.(type) {
		case *mvutil.ReqCtx:
			loadReq, _ := respData.(*mvutil.ReqCtx)
			req_context.GetInstance().MLogs.Load.Info(output.FormatLoadLog(loadReq, errCode, errMsg))
		case *mvutil.RequestParams:
			loadReqParams, _ := respData.(*mvutil.RequestParams)
			rawQuery := mvutil.RequestQueryMap(req.URL.Query())
			if strings.ToUpper(req.Method) == "POST" {
				req.ParseForm()
				for k, v := range req.PostForm {
					rawQuery[k] = v
				}
			}
			appID, _ := rawQuery.GetInt64("app_id")
			loadReqParams.Param.AppID = appID
			unitID, _ := rawQuery.GetInt64("unit_id", 0)
			loadReqParams.Param.UnitID = unitID
			token, _ := rawQuery.GetString("token", true, "")
			loadReqParams.Token = token
			loadReq := new(mvutil.ReqCtx)
			loadReq.ReqParams = loadReqParams
			req_context.GetInstance().MLogs.Load.Info(output.FormatLoadLog(loadReq, errCode, errMsg))
			//default:
			//	req_context.GetInstance().MLogs.ReqMonitor.Warnf("LoadHandler Process Error, and unknown respData type for the load log.")
		}
		req_context.GetInstance().MLogs.ReqMonitor.Warnf("LoadHandler Process Error: %s, Raw Error: %s, uri: %s, method: %s", wrapErrMsg, rawErrMsg, uri, method)
		WriterData(now, c, respByte, constant.LoadCost)
		return
	}

	if str, ok := respData.(*[]byte); ok {
		WriterData(now, c, *str, constant.LoadCost)
	} else {
		WriterData(now, c, respByte, constant.LoadCost)
	}
	return
}

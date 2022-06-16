package hb

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	pf "github.com/easierway/pipefiter_framework/pipefilter"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/constant"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/filter"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/req_context"
	"gitlab.mobvista.com/ADN/adnet/internal/service/hb/base_filters"
	"gitlab.mobvista.com/ADN/adnet/internal/service/hb/query_mem_filters"
)

type QueryMemDataHandler struct {
	Pipeline pf.Filter
}

var queryMemDataFilters = []pf.Filter{&base_filters.HttpExtractFilter{}, &query_mem_filters.QueryMemDataFilter{}}

// func init() {
// 	queryMemDataFilters = []pf.Filter{&base_filters.HttpExtractFilter{}, &query_mem_filters.QueryMemDataFilter{}}
// }

func CreateQueryMemDataHandler() *QueryMemDataHandler {
	return &QueryMemDataHandler{
		Pipeline: &WallTimePipeline{
			Name:        QueryMem,
			Filters:     &queryMemDataFilters,
			TimeElapsed: make([]AtomicInt, len(queryMemDataFilters)),
		},
	}
}

func (qmdh *QueryMemDataHandler) ServeHTTP(c http.ResponseWriter, req *http.Request) {
	// catch panic
	defer func() {
		if r := recover(); r != nil {
			msg := fmt.Sprintf("QueryMemDataHandler Process panic %v, stack=%s", r, debug.Stack())
			fmt.Println(msg)
			req_context.GetInstance().MLogs.Runtime.Error(msg)
			c.WriteHeader(http.StatusBadRequest)
		}
	}()
	now := time.Now().UnixNano()
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		WriterData(now, c, []byte(err.Error()), constant.BidCost)
		return
	}
	req.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	// business process
	var bytes []byte
	respData, err := qmdh.Pipeline.Process(req)
	if err != nil {
		wrapErrMsg, rawErrMsg := filter.FormatErrorMessage(err.Error())
		method := strings.ToUpper(req.Method)
		url := req.URL
		uri := url.RequestURI()
		//var b strings.Builder
		//if method == "POST" {
		//	b.Write(body)
		//}
		req_context.GetInstance().MLogs.ReqMonitor.Warnf("QueryMemDataHandler Process Error: %s, Raw Error: %s, uri: %s, method: %s", wrapErrMsg, rawErrMsg, uri, method)
		c.WriteHeader(http.StatusNoContent)
		return
	}
	if res, ok := respData.(string); ok {
		bytes = []byte(res)
		c.WriteHeader(http.StatusOK)
		c.Write(bytes)
	} else {
		c.WriteHeader(http.StatusInternalServerError)
	}
	return
}

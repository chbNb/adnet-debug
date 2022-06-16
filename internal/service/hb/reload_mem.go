package hb

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime/debug"
	"strings"

	pf "github.com/easierway/pipefiter_framework/pipefilter"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/filter"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/req_context"
	pp "gitlab.mobvista.com/ADN/adnet/internal/process_pipeline"
	"gitlab.mobvista.com/ADN/adnet/internal/service/hb/base_filters"
	"gitlab.mobvista.com/ADN/exporter/metrics"
)

type ReloadMemHandler struct {
	Name     string
	Pipeline pf.Filter
}

var ReloadMemAllFilters = []pf.Filter{&base_filters.HttpExtractFilter{}, &pp.ReloadMemAllFilter{}}

func CreateReloadMemAllHandler() *ReloadMemHandler {
	return &ReloadMemHandler{
		Name: ReloadAll,
		Pipeline: &WallTimePipeline{
			Name:        ReloadAll,
			Filters:     &ReloadMemAllFilters,
			TimeElapsed: make([]AtomicInt, len(ReloadMemAllFilters)),
		},
	}
}

var ReloadMemOneFilters = []pf.Filter{&base_filters.HttpExtractFilter{}, &pp.ReloadMemOneFilter{}}

func CreateReloadMemOneHandler() *ReloadMemHandler {
	return &ReloadMemHandler{
		Name: ReloadOne,
		Pipeline: &WallTimePipeline{
			Name:        ReloadOne,
			Filters:     &ReloadMemOneFilters,
			TimeElapsed: make([]AtomicInt, len(ReloadMemOneFilters)),
		},
	}
}

func (handler *ReloadMemHandler) ServeHTTP(c http.ResponseWriter, req *http.Request) {
	// catch panic
	defer func() {
		if r := recover(); r != nil {
			msg := fmt.Sprintf("%s Process panic %v, stack=%s", handler.Name, r, debug.Stack())
			fmt.Println(msg)
			req_context.GetInstance().MLogs.Runtime.Error(msg)
			metrics.IncCounterWithLabelValues(25)
			c.WriteHeader(http.StatusBadRequest)
		}
	}()
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		c.Write([]byte(err.Error()))
		return
	}
	req.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	// business process
	var bytes []byte
	respData, err := handler.Pipeline.Process(req)
	if err != nil {
		wrapErrMsg, rawErrMsg := filter.FormatErrorMessage(err.Error())
		method := strings.ToUpper(req.Method)
		url := req.URL
		uri := url.RequestURI()
		//var b strings.Builder
		//if method == "POST" {
		//	b.Write(body)
		//}
		req_context.GetInstance().MLogs.ReqMonitor.Warnf("%s Process Error: %s, Raw Error: %s, uri: %s, method: %s", handler.Name, wrapErrMsg, rawErrMsg, uri, method)
		c.WriteHeader(http.StatusBadRequest)
		c.Write([]byte(err.Error()))
		return
	}
	if res, ok := respData.(*string); ok {
		bytes = []byte(*res)
		c.WriteHeader(http.StatusOK)
		c.Write(bytes)
	} else {
		c.WriteHeader(http.StatusInternalServerError)
	}
	return
}

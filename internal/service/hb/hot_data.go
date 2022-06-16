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

type HotDataHandler struct {
	Name     string
	Pipeline pf.Filter
}

var HotTableRegistryFilters = []pf.Filter{&base_filters.HttpExtractFilter{}, &pp.HotTableRegistryFilter{}}
var HotDataRangeFilters = []pf.Filter{&base_filters.HttpExtractFilter{}, &pp.HotDataWithRangeFilter{}}
var HotDataAllFilters = []pf.Filter{&base_filters.HttpExtractFilter{}, &pp.HotDataAllFilter{}}

func CreateHotTableRegistryHandler() *HotDataHandler {
	return &HotDataHandler{
		Name: HotTable,
		Pipeline: &WallTimePipeline{
			Name:        HotTable,
			Filters:     &HotTableRegistryFilters,
			TimeElapsed: make([]AtomicInt, len(HotTableRegistryFilters)),
		},
	}
}

func CreateHotDataWithRangeHandler() *HotDataHandler {
	return &HotDataHandler{
		Name: HotDataRange,
		Pipeline: &WallTimePipeline{
			Name:        HotDataRange,
			Filters:     &HotDataRangeFilters,
			TimeElapsed: make([]AtomicInt, len(HotDataRangeFilters)),
		},
	}
}

func CreateHotDataAllHandler() *HotDataHandler {
	return &HotDataHandler{
		Name: HotDataAll,
		Pipeline: &WallTimePipeline{
			Name:        HotDataAll,
			Filters:     &HotDataAllFilters,
			TimeElapsed: make([]AtomicInt, len(HotDataAllFilters)),
		},
	}
}

func (handler *HotDataHandler) ServeHTTP(c http.ResponseWriter, req *http.Request) {
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

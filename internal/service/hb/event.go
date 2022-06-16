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
	pp "gitlab.mobvista.com/ADN/adnet/internal/process_pipeline"
	"gitlab.mobvista.com/ADN/adnet/internal/service/hb/event_filters"
	"gitlab.mobvista.com/ADN/adnet/internal/watcher"
	"gitlab.mobvista.com/ADN/exporter/metrics"
)

type EventHandler struct {
	Pipeline pf.Filter
}

var eventFilters = []pf.Filter{
	&pp.MVReqparamTransFilter{},
	&event_filters.ReqParamFilter{},
}

// func init() {
// 	eventFilters = []pf.Filter{&base_filters.HttpExtractFilter{}, &event_filters.ReqParamFilter{}}
// }

func CreateEventHandler() *EventHandler {
	return &EventHandler{
		Pipeline: &WallTimePipeline{
			Name:        Event,
			Filters:     &eventFilters,
			TimeElapsed: make([]AtomicInt, len(eventFilters)),
		},
	}
}

func (eh *EventHandler) ServeHTTP(c http.ResponseWriter, req *http.Request) {
	// catch panic
	defer func() {
		if r := recover(); r != nil {
			msg := fmt.Sprintf("EventHandler Process panic %v, stack=%s", r, debug.Stack())
			fmt.Println(msg)
			req_context.GetInstance().MLogs.Runtime.Error(msg)
			metrics.IncCounterWithLabelValues(25)
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
	_, err = eh.Pipeline.Process(req)
	if err != nil {
		wrapErrMsg, rawErrMsg := filter.FormatErrorMessage(err.Error())
		method := strings.ToUpper(req.Method)
		url := req.URL
		uri := url.RequestURI()
		//var b strings.Builder
		//if method == "POST" {
		//	b.Write(body)
		//}
		req_context.GetInstance().MLogs.ReqMonitor.Warnf("EventHandler Process Error: %s, Raw Error: %s, uri: %s, method: %s", wrapErrMsg, rawErrMsg, uri, method)
		c.WriteHeader(http.StatusNoContent)
		watcher.AddAvgWatchValue(constant.EventCost, float64((time.Now().UnixNano()-now)/1e6))
		return
	}
	c.WriteHeader(http.StatusOK)
	c.Write(bytes)
	watcher.AddAvgWatchValue(constant.EventCost, float64((time.Now().UnixNano()-now)/1e6))
}

package hb

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"gitlab.mobvista.com/ADN/exporter/metrics"

	pf "github.com/easierway/pipefiter_framework/pipefilter"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/req_context"
)

type HealthHandler struct {
	Pipeline pf.Filter
}

var healthFilters []pf.Filter

func init() {
	healthFilters = []pf.Filter{}
}

func CreateHealthHandler() *HealthHandler {
	return &HealthHandler{
		Pipeline: &WallTimePipeline{
			Name:        Health,
			Filters:     &healthFilters,
			TimeElapsed: make([]AtomicInt, len(healthFilters)),
		},
	}
}

func (eh *HealthHandler) ServeHTTP(c http.ResponseWriter, req *http.Request) {
	// catch panic
	defer func() {
		if r := recover(); r != nil {
			msg := fmt.Sprintf("HealthHandler Process panic %v, stack=%s", r, debug.Stack())
			fmt.Println(msg)
			req_context.GetInstance().MLogs.Runtime.Error(msg)
			metrics.IncCounterWithLabelValues(25)
			c.WriteHeader(http.StatusBadRequest)
		}
	}()
	// business process
	var bytes []byte
	eh.Pipeline.Process(req)
	if _, err := req_context.GetInstance().CreativeCacheClient.Conn.Ping().Result(); err != nil {
		//req_context.GetInstance().MLogs.Runtime.Warnf("HealthHandler Process Error: %s", err.Error())
		c.WriteHeader(http.StatusNoContent)
		return
	}
	c.WriteHeader(http.StatusOK)
	c.Write(bytes)
}

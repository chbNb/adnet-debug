package service

import (
	"fmt"
	"io"
	"net/http"
	"runtime/debug"
	"time"

	"gitlab.mobvista.com/ADN/exporter/metrics"

	pf "github.com/easierway/pipefiter_framework/pipefilter"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	pp "gitlab.mobvista.com/ADN/adnet/internal/process_pipeline"
)

type ToolsHandler struct {
	Pipeline pf.Filter
}

func CreateQueryMemHandler() *ToolsHandler {
	queryMemFilters := []pf.Filter{&pp.QueryMemDataFilter{}}
	queryMemPipeline := &WallTimePipeline{
		Name:        "tools-querymem",
		Filters:     &queryMemFilters,
		TimeElapsed: make([]AtomicInt, len(queryMemFilters)),
	}
	return &ToolsHandler{
		Pipeline: queryMemPipeline,
	}
}

func CreateReloadMemOneHandler() *ToolsHandler {
	filter := []pf.Filter{&pp.ReloadMemOneFilter{}}
	pipeline := &WallTimePipeline{
		Name:        "tools-reload-mem-one",
		Filters:     &filter,
		TimeElapsed: make([]AtomicInt, len(filter)),
	}
	return &ToolsHandler{
		Pipeline: pipeline,
	}
}

func CreateReloadMemAllHandler() *ToolsHandler {
	filter := []pf.Filter{&pp.ReloadMemAllFilter{}}
	pipeline := &WallTimePipeline{
		Name:        "tools-reload-mem-all",
		Filters:     &filter,
		TimeElapsed: make([]AtomicInt, len(filter)),
	}
	return &ToolsHandler{
		Pipeline: pipeline,
	}
}

func CreateHotTableRegistryHandler() *ToolsHandler {
	filter := []pf.Filter{&pp.HotTableRegistryFilter{}}
	pipeline := &WallTimePipeline{
		Name:        "tools-hot-table-registry",
		Filters:     &filter,
		TimeElapsed: make([]AtomicInt, len(filter)),
	}
	return &ToolsHandler{
		Pipeline: pipeline,
	}
}

func CreateHotDataWithRangeHandler() *ToolsHandler {
	filter := []pf.Filter{&pp.HotDataWithRangeFilter{}}
	pipeline := &WallTimePipeline{
		Name:        "tools-hot-data-with-range",
		Filters:     &filter,
		TimeElapsed: make([]AtomicInt, len(filter)),
	}
	return &ToolsHandler{
		Pipeline: pipeline,
	}
}

func CreateHotDataAllHandler() *ToolsHandler {
	filter := []pf.Filter{&pp.HotDataAllFilter{}}
	pipeline := &WallTimePipeline{
		Name:        "tools-hot-data-all",
		Filters:     &filter,
		TimeElapsed: make([]AtomicInt, len(filter)),
	}
	return &ToolsHandler{
		Pipeline: pipeline,
	}
}

func CreateSnapshotHandler() *ToolsHandler {
	snapshotFilters := []pf.Filter{&pp.CaptureAdPackFilter{}}
	snapshotPipeline := &WallTimePipeline{
		Name:        "tools-snapshot",
		Filters:     &snapshotFilters,
		TimeElapsed: make([]AtomicInt, len(snapshotFilters)),
	}
	return &ToolsHandler{
		Pipeline: snapshotPipeline,
	}
}

func CreateVersionHandler() *ToolsHandler {
	snapshotFilters := []pf.Filter{&pp.VersionFilter{}}
	snapshotPipeline := &WallTimePipeline{
		Name:        "tools-vesion",
		Filters:     &snapshotFilters,
		TimeElapsed: make([]AtomicInt, len(snapshotFilters)),
	}
	return &ToolsHandler{
		Pipeline: snapshotPipeline,
	}
}

func (conn *HTTPConnector) InitToolsRouter() {
	conn.router.Handle("/query_mem", CreateQueryMemHandler())
	conn.router.Handle("/snap", CreateSnapshotHandler())
	conn.router.Handle("/version", CreateVersionHandler())
	conn.router.Handle("/reload_mem_one", CreateReloadMemOneHandler())
	conn.router.Handle("/reload_mem_all", CreateReloadMemAllHandler())
	conn.router.Handle("/hot_table_registry", CreateHotTableRegistryHandler())
	conn.router.Handle("/hot_data_range", CreateHotDataWithRangeHandler())
	conn.router.Handle("/hot_data_all", CreateHotDataAllHandler())
}

func (toolsHandler *ToolsHandler) ServeHTTP(c http.ResponseWriter, req *http.Request) {
	//catch panic
	defer func() {
		if r := recover(); r != nil {
			msg := fmt.Sprintf("adnet Panic: %v, stack=[%s]", r, string(debug.Stack()))
			fmt.Println(msg)
			mvutil.Logger.Runtime.Error(msg)
			metrics.IncCounterWithLabelValues(25)
			io.WriteString(c, "tools handler panic")
		}
	}()
	now := time.Now().UnixNano()

	// business process
	ret, err := toolsHandler.Pipeline.Process(req)
	if err != nil {
		mvutil.Logger.Runtime.Warnf("tools handler ServeHTTP get error=[%s]", err.Error())
		pp.WriterData(now, c, err.Error())
		return
	}

	if _, ok := ret.(*string); !ok {
		mvutil.Logger.Runtime.Warnf("tools handler output not string")
		pp.WriterData(now, c, "toos handler output not string")
		return
	}

	pp.WriterData(now, c, *ret.(*string))
}

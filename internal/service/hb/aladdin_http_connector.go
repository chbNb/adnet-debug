package hb

import (
	"context"
	"net/http"
	"net/http/pprof"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/req_context"
	"gitlab.mobvista.com/ADN/adnet/internal/utility"
)

type AladdingHTTPConnector struct {
	httpAddr   string
	httpServer *http.Server
	router     *mux.Router
	reqCtx     *req_context.ReqContext
}

func (conn *AladdingHTTPConnector) Start(ctx context.Context) error {
	return conn.httpServer.ListenAndServe()
}

func (conn *AladdingHTTPConnector) Stop(ctx context.Context) error {
	conn.reqCtx.Close()
	return conn.httpServer.Close()
}

func (conn *AladdingHTTPConnector) InitRouter() {
	conn.router.Use(utility.HttpHandlerInterceptor)    // 上报qps、响应时间等信息
	conn.router.Handle("/metrics", promhttp.Handler()) // prometheus client api
	conn.router.Handle("/bid", CreateBidHandler())     // Pipeline
	conn.router.Handle("/bid_ads", http.AllowQuerySemicolons(CreateBidAdsHandler()))
	conn.router.Handle("/v2/bid", CreateMopubBidHandler())
	conn.router.Handle("/load", CreateLoadHandler())
	conn.router.Handle("/win", CreateEventHandler())
	conn.router.Handle("/loss", CreateEventHandler())
	conn.router.Handle("/billing", CreateEventHandler())
	conn.router.Handle("/health", CreateHealthHandler())
	conn.router.Handle("/query_mem", CreateQueryMemDataHandler())
	conn.router.Handle("/reload_mem_one", CreateReloadMemOneHandler())
	conn.router.Handle("/reload_mem_all", CreateReloadMemAllHandler())
	conn.router.Handle("/hot_table_registry", CreateHotTableRegistryHandler())
	conn.router.Handle("/hot_data_range", CreateHotDataWithRangeHandler())
	conn.router.Handle("/hot_data_all", CreateHotDataAllHandler())
	conn.InitPprof()
	conn.httpServer.Handler = conn.router
}

func (conn *AladdingHTTPConnector) InitPprof() {
	var prefix = "/pprof"
	var names = []string{"allocs", "block", "goroutine", "heap", "mutex", "threadcreate"}
	var handles = map[string]http.HandlerFunc{
		"profile": pprof.Profile,
		"symbol":  pprof.Symbol,
		"trace":   pprof.Trace,
		"cmdline": pprof.Cmdline,
	}
	conn.router.HandleFunc(prefix+"/", pprof.Index)
	for _, name := range names {
		conn.router.Handle(filepath.Join(prefix, name), pprof.Handler(name))
	}
	for path, handle := range handles {
		conn.router.HandleFunc(filepath.Join(prefix, path), handle)
	}
}

func CreateHTTPConnector(httpAddr string) *AladdingHTTPConnector {
	return &AladdingHTTPConnector{
		httpAddr: httpAddr,
		router:   mux.NewRouter(),
		reqCtx:   req_context.GetInstance(),
		httpServer: &http.Server{
			Addr:           httpAddr,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		},
	}
}

package service

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gitlab.mobvista.com/ADN/adnet/internal/req_context"
	"gitlab.mobvista.com/ADN/adnet/internal/utility"
)

type HTTPConnector struct {
	httpAddr   string
	httpServer *http.Server
	router     *mux.Router
	reqCtx     *req_context.ReqContext
}

func (conn *HTTPConnector) Start(ctx context.Context) error {
	return conn.httpServer.ListenAndServe()
}

func (conn *HTTPConnector) Stop(ctx context.Context) error {
	//conn.reqCtx.Close()
	return conn.httpServer.Close()
}

func (conn *HTTPConnector) InitRouter() {
	conn.router.Use(utility.HttpHandlerInterceptor)    // 中间件应该是对每个请求进行监控，发送数据给监控平台
	conn.router.Handle("/metrics", promhttp.Handler()) // prometheus client api
	conn.InitPprof()                                   //  监控"goroutine", "heap", "mutex"等信息
	conn.InitNetworkRouter()                           //  配置Network pipeline，包含很多个url-handler映射，每个url都对应一个pipeline
	conn.InitMobPowerRouter()
	conn.InitToolsRouter()
	conn.httpServer.Handler = conn.router
}

func CreateHTTPConnector(httpAddr string) *HTTPConnector {
	return &HTTPConnector{
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

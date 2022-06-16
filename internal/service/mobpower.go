package service

import (
	"fmt"
	"io"
	"net/http"
	"runtime/debug"
	"time"

	"gitlab.mobvista.com/ADN/exporter/metrics"

	pf "github.com/easierway/pipefiter_framework/pipefilter"
	"gitlab.mobvista.com/ADN/adnet/internal/errorcode"
	"gitlab.mobvista.com/ADN/adnet/internal/mvconst"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	pp "gitlab.mobvista.com/ADN/adnet/internal/process_pipeline"
)

type MobPowerHandler struct {
	Pipeline pf.Filter
}

func CreateMobPowerHandler() *MobPowerHandler {
	filters := []pf.Filter{&pp.MPReqparamTransFilter{}, &pp.RequestParamFilter{}, &pp.IpInfoFilter{},
		&pp.UaParserFilter{}, &pp.MpAdFilter{}, &pp.ParamRenderFilter{}, &pp.ReplaceBrandFilter{},
		&pp.ReqValidateFilter{}, &pp.AdBackendMaker{}, &pp.BackendReqFilter{}, &pp.AdRenderFilter{}}
	pipeline := &WallTimePipeline{
		Name:        "mobpower",
		Filters:     &filters,
		TimeElapsed: make([]AtomicInt, len(filters)),
	}
	return &MobPowerHandler{
		Pipeline: pipeline,
	}
}

func (conn *HTTPConnector) InitMobPowerRouter() {
	conn.router.Handle(mvconst.PATHMPAD, CreateMobPowerHandler())
	conn.router.Handle(mvconst.PATHMPNewAD, CreateMobPowerHandler())
	conn.router.Handle(mvconst.PATHMPADV2, CreateMobPowerHandler())
}

func (mpHandler *MobPowerHandler) ServeHTTP(c http.ResponseWriter, req *http.Request) {
	//catch panic
	defer func() {
		if r := recover(); r != nil {
			msg := fmt.Sprintf("adnet Panic: %v, stack=[%s]", r, string(debug.Stack()))
			fmt.Println(msg)
			mvutil.Logger.Runtime.Error(msg)
			metrics.IncCounterWithLabelValues(25)
			io.WriteString(c, errorcode.EXCEPTION_RETURN_EMPTY.Error())
		}
	}()
	now := time.Now().UnixNano()
	//rateLimit
	if !pp.RateLimitDecorator.TryToGetToken() {
		pp.LimitError(now, c, req)
		return
	}

	//设置跨域参数
	c.Header().Set("Access-Control-Allow-Origin", "*")
	c.Header().Set("Content-Type", "application/json;charset=UTF-8")
	// business process
	ret, err := mpHandler.Pipeline.Process(req)
	if err != nil {
		//catch sepcailCode error
		if sepcailCode, ok := err.(errorcode.SpecialCode); ok {
			pp.WriterData(now, c, sepcailCode.Error())
			return
		}

		if adnetCode, ok := err.(errorcode.AdnetCode); ok {

			if req.URL.Path == mvconst.PATHMPAD {
				pp.WriterData(now, c, mvutil.Base64EncodeMP(adnetCode.Error()))
				return
			}

			// new mp sdk error
			if req.URL.Path == mvconst.PATHMPNewAD {
				pp.WriterData(now, c, mvutil.Base64EncodeMPNew(adnetCode.Error()))
				return
			}

			if req.URL.Path == mvconst.PATHMPADV2 {
				// 无填充则不返回任何值，http code为204
				if adnetCode == mvconst.EXCEPTION_RETURN_EMPTY {
					pp.WriterHeader(now, c, http.StatusNoContent)
					return
				}
				pp.WriterData(now, c, adnetCode.Error())
				return
			}
		}
		mvutil.Logger.Runtime.Errorf("mobpower handler ServeHTTP get error=[%s]", err.Error())
		pp.WriterData(now, c, errorcode.EXCEPTION_RETURN_EMPTY.Error())
		return
	}

	if _, ok := ret.(*string); !ok {
		mvutil.Logger.Runtime.Errorf("mobpower handler output not string")
		pp.WriterData(now, c, errorcode.EXCEPTION_RETURN_EMPTY.Error())
		return
	}

	if req.URL.Path == mvconst.PATHMPAD {
		pp.WriterData(now, c, mvutil.Base64EncodeMP(*ret.(*string)))
		return
	}

	if req.URL.Path == mvconst.PATHMPNewAD {
		pp.WriterData(now, c, mvutil.Base64EncodeMPNew(*ret.(*string)))
		return
	}
	// PATHMPADV2 不加密返回
	pp.WriterData(now, c, *ret.(*string))
}

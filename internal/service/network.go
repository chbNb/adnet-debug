package service

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"gitlab.mobvista.com/ADN/exporter/metrics"

	pf "github.com/easierway/pipefiter_framework/pipefilter"
	"github.com/gogo/protobuf/proto"
	"gitlab.mobvista.com/ADN/adnet/internal/errorcode"
	"gitlab.mobvista.com/ADN/adnet/internal/extractor"
	"gitlab.mobvista.com/ADN/adnet/internal/mvconst"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	"gitlab.mobvista.com/ADN/adnet/internal/output"
	pp "gitlab.mobvista.com/ADN/adnet/internal/process_pipeline"
	"gitlab.mobvista.com/ADN/adnet/internal/protobuf"
)

const APIV3 string = "api_v3"
const MPAD string = "mpad"
const APIONLINE string = "api_online"
const RTB string = "rtb"
const FUNADX string = "funadx"
const APIJSSDK string = "api_jssdk"
const HUPU string = "hupu"
const MIGUADX string = "miguadx"
const MAXADX string = "maxadx"
const CHEETAH string = "cheetah"
const KSYUN string = "ksyun"
const APIV4 string = "api_v4"
const APIV5 string = "api_v5"
const IMAGE string = "image"
const REQUEST string = "request"
const PPTV string = "pptv"
const APIV2 string = "api_v2"
const IFENGADX string = "ifeng"
const QueryMem string = "query_mem"
const Snapshot string = "snap"
const MANGOADX string = "mangoadx"
const TOPON string = "topon"
const MOREOFFER string = "moreoffer"

var networkApis map[string]string

func init() {
	networkApis = map[string]string{
		mvconst.PATHOpenApiV3:        APIV3,
		mvconst.PATHOnlineApi:        APIONLINE,
		mvconst.PATHOpenApiV2:        APIV2,
		mvconst.PATHOpenApiV4:        APIV4,
		mvconst.PATHOpenApiV5:        APIV5,
		mvconst.PATHIMAGE:            IMAGE,
		mvconst.PATHRTB:              RTB,
		mvconst.PATHFUNADX:           FUNADX,
		mvconst.PATHJssdkApi:         APIJSSDK,
		mvconst.PATHHUPUADX:          HUPU,
		mvconst.PATHMIGUADX:          MIGUADX,
		mvconst.PATHMAXADX:           MAXADX,
		mvconst.PATHCHEETAH:          CHEETAH,
		mvconst.PATHKSYUN:            KSYUN,
		mvconst.PATHIFENGADX:         IFENGADX,
		mvconst.PATHREQUEST:          REQUEST,
		mvconst.PATHPPTV:             PPTV,
		mvconst.PATHMANGOADX:         MANGOADX,
		mvconst.PATHTOPON:            TOPON,
		mvconst.PATHOpenApiMoreOffer: MOREOFFER,
	}
}

type NetworkHandler struct {
	Pipeline pf.Filter
}

// 根据不同的name,返回不同的pipeline
func CreatePipeline(name string) ([]pf.Filter, []string, error) {
	switch name {
	case APIV3, APIV4, APIV2, APIV5, MOREOFFER:
		return []pf.Filter{&pp.MVReqparamTransFilter{}, &pp.RequestParamFilter{}, &pp.IpInfoFilter{},
			&pp.UaParserFilter{}, &pp.ParamRenderFilter{}, &pp.ReplaceBrandFilter{}, &pp.MappingServerFilter{}, &pp.ReqValidateFilter{},
			&pp.TrafficAllotFilter{}, &pp.BackendReqFilter{}, &pp.AdRenderFilter{}}, nil, nil
	case APIONLINE: // 量最大(20k+ qps), 可作为pipeline的样本采集
		labels := []string{
			"MVReqparamTransFilter", "RequestParamFilter", "IpInfoFilter", "UaParserFilter",
			"OnlineParamFilter", "ParamRenderFilter", "ReplaceBrandFilter", "ReqValidateFilter",
			"TrafficAllotFilter", "BackendReqFilter", "AdRenderFilter"} // 注意与下面的pipeline顺序保持一致
		return []pf.Filter{&pp.MVReqparamTransFilter{}, &pp.RequestParamFilter{}, &pp.IpInfoFilter{}, &pp.UaParserFilter{},
			&pp.OnlineParamFilter{}, &pp.ParamRenderFilter{}, &pp.ReplaceBrandFilter{}, &pp.ReqValidateFilter{},
			&pp.TrafficAllotFilter{}, &pp.BackendReqFilter{}, &pp.AdRenderFilter{}}, labels, nil
	case IMAGE:
		return []pf.Filter{&pp.MVReqparamTransFilter{}, &pp.RequestParamFilter{}, &pp.ParamRenderFilter{}, &pp.ImageFilter{}}, nil, nil
	case RTB:
		return []pf.Filter{&pp.RTBReqparamTransFilter{}, &pp.RequestParamFilter{}, &pp.IpInfoFilter{},
			&pp.UaParserFilter{}, &pp.OnlineParamFilter{}, &pp.ParamRenderFilter{}, &pp.ReplaceBrandFilter{},
			&pp.ReqValidateFilter{}, &pp.TrafficAllotFilter{}, &pp.BackendReqFilter{}, &pp.RtbAdRenderFilter{}}, nil, nil
	case FUNADX:
		return []pf.Filter{&pp.FUNADXReqparamTransFilter{}, &pp.RequestParamFilter{}, &pp.IpInfoFilter{}, &pp.UaParserFilter{},
			&pp.OnlineParamFilter{}, &pp.ParamRenderFilter{}, &pp.ReplaceBrandFilter{}, &pp.ReqValidateFilter{},
			&pp.TrafficAllotFilter{}, &pp.BackendReqFilter{}, &pp.FunAdRenderFilter{}}, nil, nil
	case APIJSSDK:
		return []pf.Filter{&pp.JssdkReqparamTransFilter{}, &pp.RequestParamFilter{}, &pp.IpInfoFilter{},
			&pp.UaParserFilter{}, &pp.ParamRenderFilter{}, &pp.ReplaceBrandFilter{}, &pp.ReqValidateFilter{},
			&pp.TrafficAllotFilter{}, &pp.BackendReqFilter{}, &pp.JssdkAdRenderFilter{}}, nil, nil
	case HUPU:
		return []pf.Filter{&pp.HUPUReqparamTransFilter{}, &pp.RequestParamFilter{}, &pp.IpInfoFilter{},
			&pp.UaParserFilter{}, &pp.OnlineParamFilter{}, &pp.ParamRenderFilter{}, &pp.ReplaceBrandFilter{},
			&pp.ReqValidateFilter{}, &pp.TrafficAllotFilter{}, &pp.BackendReqFilter{}, &pp.HupuAdRenderFilter{}}, nil, nil
	case MIGUADX:
		return []pf.Filter{&pp.MIGUReqparamTransFilter{}, &pp.RequestParamFilter{}, &pp.IpInfoFilter{},
			&pp.UaParserFilter{}, &pp.OnlineParamFilter{}, &pp.ParamRenderFilter{}, &pp.ReplaceBrandFilter{},
			&pp.ReqValidateFilter{}, &pp.TrafficAllotFilter{}, &pp.BackendReqFilter{}, &pp.MIGUAdRenderFilter{}}, nil, nil
	case MAXADX:
		return []pf.Filter{&pp.MAXReqparamTransFilter{}, &pp.RequestParamFilter{}, &pp.IpInfoFilter{},
			&pp.UaParserFilter{}, &pp.OnlineParamFilter{}, &pp.ParamRenderFilter{}, &pp.ReplaceBrandFilter{},
			&pp.ReqValidateFilter{}, &pp.TrafficAllotFilter{}, &pp.BackendReqFilter{}, &pp.MaxAdRenderFilter{}}, nil, nil
	case CHEETAH:
		return []pf.Filter{&pp.CMReqparamTransFilter{}, &pp.RequestParamFilter{}, &pp.IpInfoFilter{},
			&pp.UaParserFilter{}, &pp.OnlineParamFilter{}, &pp.ParamRenderFilter{}, &pp.ReplaceBrandFilter{},
			&pp.ReqValidateFilter{}, &pp.TrafficAllotFilter{}, &pp.BackendReqFilter{}, &pp.CMAdRenderFilter{}}, nil, nil
	case KSYUN:
		return []pf.Filter{&pp.KSYUNReqParamTransFilter{}, &pp.RequestParamFilter{}, &pp.IpInfoFilter{},
			&pp.UaParserFilter{}, &pp.OnlineParamFilter{}, &pp.ParamRenderFilter{}, &pp.ReplaceBrandFilter{},
			&pp.ReqValidateFilter{}, &pp.TrafficAllotFilter{}, &pp.BackendReqFilter{}, &pp.KSYUNAdRenderFilter{}}, nil, nil
	case IFENGADX:
		return []pf.Filter{&pp.IFENGReqparamTransFilter{}, &pp.RequestParamFilter{}, &pp.IpInfoFilter{},
			&pp.UaParserFilter{}, &pp.OnlineParamFilter{}, &pp.ParamRenderFilter{}, &pp.ReplaceBrandFilter{},
			&pp.ReqValidateFilter{}, &pp.TrafficAllotFilter{}, &pp.BackendReqFilter{}, &pp.IFENGAdRenderFilter{}}, nil, nil
	case REQUEST:
		return []pf.Filter{&pp.RequestReqparamTransFilter{}, &pp.IpInfoFilter{}, &pp.UaParserFilter{},
			&pp.ParamRenderFilter{}, &pp.ReplaceBrandFilter{}, &pp.ReqValidateFilter{},
			&pp.TrafficAllotFilter{}, &pp.BackendReqFilter{}, &pp.RequestAdRenderFilter{}}, nil, nil
	case PPTV:
		return []pf.Filter{&pp.PPTVReqParamTransFilter{}, &pp.RequestParamFilter{}, &pp.IpInfoFilter{},
			&pp.UaParserFilter{}, &pp.OnlineParamFilter{}, &pp.ParamRenderFilter{}, &pp.ReplaceBrandFilter{},
			&pp.ReqValidateFilter{}, &pp.TrafficAllotFilter{}, &pp.BackendReqFilter{}, &pp.PPTVAdRenderFilter{}}, nil, nil
	case MANGOADX:
		return []pf.Filter{&pp.MangoReqparamTransFilter{}, &pp.RequestParamFilter{}, &pp.IpInfoFilter{},
			&pp.UaParserFilter{}, &pp.OnlineParamFilter{}, &pp.ParamRenderFilter{}, &pp.ReplaceBrandFilter{},
			&pp.ReqValidateFilter{}, &pp.TrafficAllotFilter{}, &pp.BackendReqFilter{}, &pp.MangoAdRenderFilter{}}, nil, nil
	case TOPON:
		return []pf.Filter{&pp.ToponReqparamTransFilter{}, &pp.IpInfoFilter{},
			&pp.UaParserFilter{}, &pp.ParamRenderFilter{}, &pp.ReplaceBrandFilter{}, &pp.MappingServerFilter{}, &pp.ReqValidateFilter{},
			&pp.TrafficAllotFilter{}, &pp.BackendReqFilter{}, &pp.ToponAdRenderFilter{}}, nil, nil
	default:
		return nil, nil, pp.NotFoundPipelineError
	}
}

func CreateNetworkHandler(name string) *NetworkHandler {
	filters, labels, _ := CreatePipeline(name)
	pipeline := &WallTimePipeline{
		Name:        name,
		Filters:     &filters,
		TimeElapsed: make([]AtomicInt, len(filters)),
		Labels:      labels,
	}
	pipeline.Show(name)
	return &NetworkHandler{
		Pipeline: pipeline,
	}
}

func (conn *HTTPConnector) InitNetworkRouter() {
	for path, handler := range networkApis {
		conn.router.Handle(path, CreateNetworkHandler(handler))
	}
}

func (nwHandler *NetworkHandler) ServeHTTP(c http.ResponseWriter, req *http.Request) {
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
	// if !pp.RateLimitDecorator.TryToGetToken() {
	// 	pp.LimitError(now, c, req)
	// 	return
	// }
	//设置跨域参数
	c.Header().Set("Access-Control-Allow-Origin", "*")
	c.Header().Set("Content-Type", "application/json;charset=UTF-8")
	// business process
	ret, err := nwHandler.Pipeline.Process(req)
	if err != nil {
		//catch sepcailCode error
		if sepcailCode, ok := err.(errorcode.SpecialCode); ok {
			pp.WriterData(now, c, sepcailCode.Error())
			return
		}

		if adnetCode, ok := err.(errorcode.AdnetCode); ok {
			if req.URL.Path == mvconst.PATHRTB {
				if adnetCode != errorcode.EXCEPTION_RETURN_EMPTY {
					mvutil.Logger.Runtime.Warnf("path=[%s] error code is [%d]", req.URL.Path, adnetCode)
				}
				pp.WriterHeader(now, c, http.StatusNoContent)
				return
			}

			// 风行，咪咕，猎豹和凤凰视频返回httpcode
			if req.URL.Path == mvconst.PATHFUNADX || req.URL.Path == mvconst.PATHMIGUADX || req.URL.Path == mvconst.PATHCHEETAH || req.URL.Path == mvconst.PATHIFENGADX ||
				req.URL.Path == mvconst.PATHPPTV {
				if adnetCode == errorcode.EXCEPTION_RETURN_EMPTY {
					pp.WriterHeader(now, c, http.StatusNoContent)
					return
				} else {
					pp.WriterHeader(now, c, http.StatusBadRequest)
					mvutil.Logger.Runtime.Warnf("path=[%s] error code is [%d]", req.URL.Path, adnetCode)
					return
				}
			}

			if req.URL.Path == mvconst.PATHHUPUADX {
				if adnetCode != errorcode.EXCEPTION_RETURN_EMPTY {
					mvutil.Logger.Runtime.Warnf("Hupu error code is [%d]", adnetCode)
				}
				pp.WriterHeader(now, c, http.StatusNoContent)
				return
			}

			//金山云返回httpcode
			if req.URL.Path == mvconst.PATHKSYUN {
				if adnetCode == errorcode.EXCEPTION_RETURN_EMPTY {
					pp.WriterHeader(now, c, http.StatusNoContent)
					return
				} else {
					mvutil.Logger.Runtime.Warnf("error code is [%d]", adnetCode)
					pp.WriterHeader(now, c, http.StatusBadRequest)
					return
				}
			}
			if req.URL.Path == mvconst.PATHMAXADX {
				mvutil.Logger.Runtime.Warnf("error code is [%d]", adnetCode)
				var response protobuf.BidResponse
				bid := "0"
				response.Bid = &bid
				res, _ := proto.Marshal(&response)
				resultProto := strings.TrimRight(string(res), "\n")
				pp.WriterData(now, c, resultProto)
				return
			}
			if req.URL.Path == mvconst.PATHREQUEST {
				pp.WriterData(now, c, mvutil.EnBase64(adnetCode.Error()))
				return
			}
			if req.URL.Path == mvconst.PATHOnlineApi {
				if adnetCode != errorcode.EXCEPTION_RETURN_EMPTY {
					// 针对韩国开发者
					rawQuery := mvutil.RequestQueryMap(req.URL.Query())
					returnVastApps := extractor.Get_RETRUN_VAST_APP()
					unitId, err := rawQuery.GetInt64("unit_id", 0)
					// 判断是否走vast协议
					isVast, _ := rawQuery.GetBool("is_vast", false)
					appId, _ := rawQuery.GetInt64("app_id", 0)
					if err == nil && isVast && (output.IsAfreecatvUnit(unitId) || mvutil.InInt64Arr(appId, returnVastApps)) {
						rs, _ := output.VastReturnEmpty()
						resultJson := strings.TrimRight(string(rs), "\n")
						c.Header().Set("Content-Type", "application/xml;charset=UTF-8")
						pp.WriterData(now, c, resultJson)
						mvutil.Logger.Runtime.Warnf("error code is [%d] and unit_id is [%d]", adnetCode, unitId)
						return
					}
				}
			}
			if req.URL.Path == mvconst.PATHTOPON {
				if adnetCode == errorcode.EXCEPTION_RETURN_EMPTY {
					pp.WriterHeader(now, c, http.StatusNoContent)
					return
				} else {
					mvutil.Logger.Runtime.Warnf("path=[%s] error code is [%d]", req.URL.Path, adnetCode)
					pp.WriterHeader(now, c, http.StatusBadRequest)
					return
				}
			}
			// 芒果在无广告填充时依旧需要返回json
			if req.URL.Path == mvconst.PATHMANGOADX {
				if adnetCode != errorcode.EXCEPTION_RETURN_EMPTY {
					mvutil.Logger.Runtime.Warnf("mango error code is [%d]", adnetCode)
					pp.WriterHeader(now, c, http.StatusBadRequest)
					var res = []byte(`{"bid" : "", "version" : 3, "err_code": 204 }`)
					resultJson := strings.TrimRight(string(res), "\n")
					pp.WriterData(now, c, resultJson)
					return
				}
			}
			pp.WriterData(now, c, adnetCode.Error())
			return

		}
		mvutil.Logger.Runtime.Warnf("network handler ServeHTTP get error=[%s]", err.Error())
		pp.WriterData(now, c, errorcode.EXCEPTION_RETURN_EMPTY.Error())
		return
	}

	if _, ok := ret.(*string); !ok {
		mvutil.Logger.Runtime.Warnf("network handler output not string")
		pp.WriterData(now, c, errorcode.EXCEPTION_RETURN_EMPTY.Error())
		return
	}

	// request
	if req.URL.Path == mvconst.PATHREQUEST {
		pp.WriterData(now, c, mvutil.EnBase64(*ret.(*string)))
		return
	}

	// vast协议跨域访问
	if strings.Contains(*ret.(*string), xml.Header) && req.URL.Path == mvconst.PATHOnlineApi {
		c.Header().Set("Content-Type", "application/xml;charset=UTF-8")
		c.Header().Set("Access-Control-Allow-Origin", "*")
		c.Header().Set("Access-Control-Allow-Methods", "HEAD, GET, POST")
		c.Header().Set("Access-Control-Allow-Credentials", "true")
	}

	if req.URL.Path == mvconst.PATHMAXADX {
		c.Header().Set("Content-Type", "application/octet-stream;charset=UTF-8")
	}

	if req.URL.Path == mvconst.PATHJssdkApi {
		if strings.Contains(req.RequestURI, "rih=true") {
			c.Header().Set("Content-Type", "text/html;charset=UTF-8")
		}
	}
	pp.WriterData(now, c, *ret.(*string))
}

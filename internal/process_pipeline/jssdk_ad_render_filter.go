package process_pipeline

import (
	"errors"
	"strings"

	"github.com/json-iterator/go"
	"gitlab.mobvista.com/ADN/adnet/internal/errorcode"
	"gitlab.mobvista.com/ADN/adnet/internal/mvconst"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	"gitlab.mobvista.com/ADN/adnet/internal/output"
)

type JssdkAdRenderFilter struct{}

func (jarf *JssdkAdRenderFilter) Process(data interface{}) (interface{}, error) {
	in, ok := data.(*mvutil.ReqCtx)
	if !ok {
		return nil, errors.New("RtbAdRenderFilter input type should be *mvutil.ReqCtx")
	}
	if in.Result != nil {
		result, err := output.RenderOutput(in.ReqParams, in.Result)
		// 无论是否有err 都打request日志
		mvutil.StatRequestLog(&in.ReqParams.Param)
		//降填充日志
		mvutil.StatReduceFillLog(&in.ReqParams.Param)

		//被展示控cap不返回广告
		if in.ReqParams.Param.Debug <= 0 && !in.ReqParams.Param.DebugMode && in.ReqParams.Param.IsBlockByImpCap &&
			!in.ReqParams.Param.RespOnlyRequestThirdDsp { //如果ADX有请求三方则不需要过滤
			return nil, errorcode.EXCEPTION_CAP_BLOCK
		}

		// jssdk支持debug
		if in.ReqParams.Param.Debug > 0 {
			//debug模式输出debug信息
			return &in.ReqParams.DebugInfo, nil
		}
		if err != nil {
			if !in.ReqParams.Param.DebugMode {
				return nil, errorcode.EXCEPTION_RETURN_EMPTY
			} else {
				var mr output.MobvistaResult
				mr.Status = mvconst.EXCEPTION_RETURN_EMPTY
				mr.Msg = "EXCEPTION_RETURN_EMPTY"
				mr.Data.Ads = make([]output.Ad, 0)
				mr.DebugInfo = in.DebugModeInfo
				result = &mr
			}
		}
		// jssdk支持debugmode
		if in.ReqParams.Param.DebugMode {
			result.DebugInfo = in.DebugModeInfo
		}

		if result != nil {
			if len(result.Data.Ads) > 0 {
				result.Data.Setting = output.RenderJssdkRes(in.ReqParams)
			}
			res, err := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(result)
			if err != nil {
				return nil, errorcode.EXCEPTION_RETURN_EMPTY
			}
			resultJson := strings.TrimRight(string(res), "\n")
			if in.ReqParams.Param.ResInHtml {
				resultJson, err = output.RenderJssdkResInHtml(resultJson, in.ReqParams.Param, result)
			}
			if err != nil {
				return nil, errorcode.EXCEPTION_RETURN_EMPTY
			}
			return &resultJson, nil
		}
	}
	return nil, errorcode.EXCEPTION_RETURN_EMPTY
}

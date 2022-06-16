package process_pipeline

import (
	"errors"
	"strconv"
	"strings"

	jsoniter "github.com/json-iterator/go"
	"gitlab.mobvista.com/ADN/adnet/internal/errorcode"
	"gitlab.mobvista.com/ADN/adnet/internal/extractor"
	hboutput "gitlab.mobvista.com/ADN/adnet/internal/hb/output"
	hbreqctx "gitlab.mobvista.com/ADN/adnet/internal/hb/req_context"
	"gitlab.mobvista.com/ADN/adnet/internal/mvconst"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	"gitlab.mobvista.com/ADN/adnet/internal/output"
)

type AdRenderFilter struct {
}

func transLiebaoReqLog(params *mvutil.Params) {
	transMap, _ := extractor.GetMP_MAP_UNIT()
	if len(transMap) <= 0 {
		return
	}
	unitId := params.UnitID
	if unitId <= int64(0) {
		return
	}
	unitIdStr := strconv.FormatInt(unitId, 10)
	key := "mp_" + unitIdStr
	transObj, ok := transMap[key]
	if ok {
		params.UnitID = transObj.UnitID
		params.AppID = transObj.AppID
		params.PublisherID = transObj.PublisherID
	}
}

func (arf *AdRenderFilter) Process(data interface{}) (interface{}, error) {
	in, ok := data.(*mvutil.ReqCtx)
	if !ok {
		return nil, errors.New("AdRenderFilter input type should be *mvutil.ReqCtx")
	}
	if in.Result == nil {
		return mvconst.EXCEPTION_RETURN_EMPTY, errors.New("EXCEPTION_RETURN_EMPTY")
	}

	var (
		resOnline   output.OnlineResult
		resMPAd     output.MPResult
		resUCAd     output.UCResult
		resV2       output.V2Result
		resV4       output.V4Result
		resV5       output.V5Result
		resXMOnline output.XMOnlineResult
		resMPNewAd  output.MPNewAdResult

		xmCanNewReturn bool // 判断小米能否返回删减后的内容
		result         *output.MobvistaResult
		err            error
		resFormat      string
	)
	// 渲染要返回的素材
	result, err = output.RenderOutput(in.ReqParams, in.Result)
	resFormat = "v3"

	if result != nil && err == nil && (in.ReqParams.Param.RequestPath == mvconst.PATHOnlineApi ||
		in.ReqParams.Param.RequestPath == mvconst.PATHBidAds) {
		if len(result.Data.Ads) > 0 {
			xmCanNewReturn = output.XMNewReturn(&in.ReqParams.Param)
			if xmCanNewReturn {
				resXMOnline = output.RenderxmOnlineRes(*result)
			} else {
				resOnline = output.RenderOnlineRes(*result, in.ReqParams)
			}
			resFormat = "onlineapi"
		}
	}

	if err == nil && result != nil && (in.ReqParams.Param.RequestPath == mvconst.PATHMPAD ||
		in.ReqParams.Param.RequestPath == mvconst.PATHMPADV2) {
		if len(result.Data.Ads) > 0 {
			// TODO mpad 要加密输出
			resMPAd = output.RenderMPRes(*result, in.ReqParams)

			resFormat = "mpad"
		}
	}

	if err == nil && result != nil && in.ReqParams.Param.RequestPath == mvconst.PATHMPNewAD {
		if len(result.Data.Ads) > 0 {
			// TODO mpad 要加密输出
			resMPNewAd = output.RenderMPNewAdRes(*result, in.ReqParams)
			resFormat = "mpad"
		}
	}

	if in.ReqParams.Param.RespOnlyRequestThirdDsp { // 如果返回为仅请求三方DSP，则记录为非正常降填充
		in.ReqParams.Param.ExtReduceFillReq = 0
		in.ReqParams.Param.ExtReduceFillResp = 0
	}

	if in.ReqParams.Param.OnlineFilterByBidFloor {
		in.ReqParams.Param.ExtReduceFillReq = 0
		in.ReqParams.Param.ExtReduceFillResp = 0
	}

	// 替换猎豹unit从mp到mv
	logParams := in.ReqParams.Param
	transLiebaoReqLog(&logParams)
	if !in.ReqParams.IsHBRequest {
		// TODO refactor hb and adnet
		// 无论是否有err 都打request日志, hb 在 bid 请求打
		mvutil.StatRequestLog(&logParams)
	}

	// TODO refactor hb and adnet
	// task: remove reduce fill log
	// 降填充日志
	// mvutil.StatReduceFillLog(&logParams)

	// TODO refactor hb and adnet
	if err == nil {
		// capture adpack snapshot request && response ads
		output.CaptureAdPack(&in.ReqParams.Param, result)
	}
	// rtb uc 相关处理
	if err == nil && result != nil && in.ReqParams.Param.RequestPath == mvconst.PATHRTB {
		if len(result.Data.Ads) > 0 {
			resUCAd = output.RenderUCRes(*result, in.ReqParams)
			resFormat = "ucad"
		}
	}

	// v2 接口相关处理, resformat以v2形式在单独处理。 （在第200行+v3接口被封装了。。。）
	if err == nil && result != nil && in.ReqParams.Param.RequestPath == mvconst.PATHOpenApiV2 {
		if len(result.Data.Ads) > 0 {
			resV2 = output.RenderV2Res(*result, in.ReqParams)
			resFormat = "v2"
		}
	}
	// apiv4 相关处理
	if err == nil && result != nil && in.ReqParams.Param.RequestPath == mvconst.PATHOpenApiV4 {
		if len(result.Data.Ads) > 0 {
			resV4 = output.RenderV4Res(*result, in.ReqParams)
			resFormat = "v4"
		}
	}
	// apiv5 相关处理   由 path 改为标记位
	if err == nil && result != nil && in.ReqParams.Param.ExtDataInit.V5AbtestTag == mvconst.V5_ABTEST_V5_V5 {
		if len(result.Data.Ads) > 0 {
			resV5 = output.RenderV5Res(*result, in.ReqParams)
			resFormat = "v5"
		} else {
			result.Version = "v5"
		}
	}
	if in.ReqParams.Param.ExtDataInit.V5AbtestTag == mvconst.V5_ABTEST_V5_V3 {
		result.Version = "v3"
	}

	if in.ReqParams.Param.Debug > 0 {
		// debug模式输出debug信息
		return &in.ReqParams.DebugInfo, nil
	}

	// 被展示控cap不返回广告
	if in.ReqParams.Param.Debug <= 0 && !in.ReqParams.Param.DebugMode && in.ReqParams.Param.IsBlockByImpCap &&
		!in.ReqParams.Param.RespOnlyRequestThirdDsp { // 如果ADX有请求三方则不需要过滤
		// if isTestingAlgorix(in) {
		//	mvutil.Logger.Runtime.Infof("algorix filtered IsBlockByImpCap: request_id=[%s] unitid=[%d]",
		//		in.ReqParams.Param.RequestID, in.ReqParams.Param.UnitID)
		// }
		return nil, errorcode.EXCEPTION_CAP_BLOCK
	}

	if err != nil {
		// if isTestingAlgorix(in) {
		//	mvutil.Logger.Runtime.Infof("algorix filtered has err: request_id=[%s] unitid=[%d], err=[%s]",
		//		in.ReqParams.Param.RequestID, in.ReqParams.Param.UnitID, err.Error())
		// }
		// vast返回为空的情况
		returnVastApps := extractor.Get_RETRUN_VAST_APP()

		if output.Isvast(in.ReqParams) && in.ReqParams.Param.RequestPath == mvconst.PATHOnlineApi &&
			(output.IsAfreecatvUnit(in.ReqParams.Param.UnitID) || mvutil.InInt64Arr(in.ReqParams.Param.AppID, returnVastApps)) {
			rs, err := output.VastReturnEmpty()
			resultJson := strings.TrimRight(string(rs), "\n")
			return &resultJson, err
		}

		if !in.ReqParams.Param.DebugMode {
			return nil, errorcode.EXCEPTION_RETURN_EMPTY
			// return mvconst.EXCEPTION_RETURN_EMPTY, err
		} else {
			var mr output.MobvistaResult
			mr.Status = mvconst.EXCEPTION_RETURN_EMPTY
			mr.Msg = "EXCEPTION_RETURN_EMPTY"
			mr.DebugInfo = in.DebugModeInfo
			mr.Data.Ads = make([]output.Ad, 0)
			result = &mr
		}
	}

	if in.ReqParams.Param.DebugMode {
		result.DebugInfo = in.DebugModeInfo
		result.AsDebugInfo = in.ReqParams.Param.AsDebugParam
		result.MasDebugInfo = in.ReqParams.Param.MasDebugParam
		if in.ReqParams.Param.RequestPath == mvconst.PATHOnlineApi {
			if xmCanNewReturn {
				resXMOnline.DebugInfo = in.DebugModeInfo
			} else {
				resOnline.DebugInfo = in.DebugModeInfo
			}
		}
	}

	var res []byte
	if resFormat == "onlineapi" {
		// 判断是否走vast协议
		if output.Isvast(in.ReqParams) && (in.ReqParams.Param.RequestPath == mvconst.PATHOnlineApi || in.ReqParams.Param.RequestPath == mvconst.PATHBidAds) {
			res, err = output.RanderVastData(in.ReqParams, resOnline)
			// bigo 的特殊vast处理，为了支持deeplink。
			if output.IsVastReturnInJson(&in.ReqParams.Param) {
				res = output.RenderVastReturnInJson(res, resOnline, in.ReqParams)
			}
		} else if xmCanNewReturn {
			res, err = jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(resXMOnline)
		} else {
			res, err = jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(resOnline)
		}
	} else if resFormat == "mpad" {
		res, err = jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(resMPAd)
		if in.ReqParams.Param.RequestPath == mvconst.PATHMPNewAD {
			res, err = jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(resMPNewAd)
		}
	} else if resFormat == "ucad" { // TODO, maybe remove.
		if !resUCAd.NotEmptyAd && in.ReqParams.Param.NewRTBFlag {
			return nil, errorcode.EXCEPTION_UNIT_NO_CONFIG_MANAGE_REVENUE
		}
		res, err = jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(resUCAd)
	} else if resFormat == "v2" {
		res, err = jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(resV2)
	} else if resFormat == "v4" {
		res, err = jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(resV4)
	} else if resFormat == "v5" {
		res, err = jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(resV5)
	} else {
		res, err = jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(result)
	}

	if err != nil {
		mvutil.Logger.Runtime.Warnf("request_id=[%s] result=[%+v] to json error : %s", in.ReqParams.Param.RequestID, *result, err.Error())
		return nil, errorcode.EXCEPTION_RETURN_EMPTY
	}
	if in.ReqParams.IsHBRequest {
		hbreqctx.GetInstance().MLogs.Load.Info(hboutput.FormatLoadLog(in, 0, ""))
		// mvutil.Logger.Runtime.Debugf("=========== response ad data: %s", string(res))
		return &res, nil
	}
	resultJson := strings.TrimRight(string(res), "\n")
	return &resultJson, nil

}

// // isTestingALgorix  TODO 测试完下版本删掉，(目前版本3.5.7
// func isTestingAlgorix(in *mvutil.ReqCtx) bool {
//	return mvutil.InInt64Arr(in.ReqParams.Param.UnitID, []int64{112368, 149783, 110810, 78290, 148140, 143394, 110204})
// }

package process_pipeline

import (
	"errors"
	"strings"

	jsoniter "github.com/json-iterator/go"
	"gitlab.mobvista.com/ADN/adnet/internal/errorcode"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	"gitlab.mobvista.com/ADN/adnet/internal/output"
)

type RtbAdRenderFilter struct{}

func (ratf *RtbAdRenderFilter) Process(data interface{}) (interface{}, error) {
	in, ok := data.(*mvutil.ReqCtx)
	if !ok {
		return nil, errors.New("RtbAdRenderFilter input type should be *mvutil.ReqCtx")
	}
	if in.Result != nil {
		result, err := output.RenderOutput(in.ReqParams, in.Result)

		if in.ReqParams.Param.OnlineFilterByBidFloor {
			in.ReqParams.Param.ExtReduceFillReq = 0
			in.ReqParams.Param.ExtReduceFillResp = 0
		}
		// 无论是否有err 都打request日志
		mvutil.StatRequestLog(&in.ReqParams.Param)

		//降填充日志
		mvutil.StatReduceFillLog(&in.ReqParams.Param)
		if err != nil {
			return nil, errorcode.EXCEPTION_RETURN_EMPTY
		}

		if result != nil && len(result.Data.Ads) > 0 {
			ucRes := output.RenderUCRes(*result, in.ReqParams)
			if !ucRes.NotEmptyAd && in.ReqParams.Param.NewRTBFlag {
				return nil, errorcode.EXCEPTION_UNIT_NO_CONFIG_MANAGE_REVENUE
			}
			res, err := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(ucRes)
			if err != nil {
				return nil, errorcode.EXCEPTION_RETURN_EMPTY
			}
			resultJson := strings.TrimRight(string(res), "\n")
			return &resultJson, nil
		}
	}
	return nil, errorcode.EXCEPTION_RETURN_EMPTY
}

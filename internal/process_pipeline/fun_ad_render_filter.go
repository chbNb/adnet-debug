package process_pipeline

import (
	"errors"
	"strings"

	jsoniter "github.com/json-iterator/go"
	"gitlab.mobvista.com/ADN/adnet/internal/errorcode"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	"gitlab.mobvista.com/ADN/adnet/internal/output"
)

type FunAdRenderFilter struct{}

func (farf *FunAdRenderFilter) Process(data interface{}) (interface{}, error) {
	in, ok := data.(*mvutil.ReqCtx)
	if !ok {
		return nil, errors.New("FunAdRenderFilter input type should be *mvutil.ReqCtx")
	}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	if in.Result != nil {
		result, err := output.RenderOutput(in.ReqParams, in.Result)
		// 无论是否有err 都打request日志
		mvutil.StatRequestLog(&in.ReqParams.Param)
		//降填充日志
		mvutil.StatReduceFillLog(&in.ReqParams.Param)
		if err != nil {
			return nil, errorcode.EXCEPTION_RETURN_EMPTY
		}

		if result != nil && len(result.Data.Ads) > 0 {
			res, err := json.Marshal(output.RenderFunRes(*result, in.ReqParams))
			if err != nil {
				return nil, errorcode.EXCEPTION_RETURN_EMPTY
			}
			resultJson := strings.TrimRight(string(res), "\n")
			return &resultJson, nil
		}
	}
	return nil, errorcode.EXCEPTION_RETURN_EMPTY
}

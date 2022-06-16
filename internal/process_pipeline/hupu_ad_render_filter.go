package process_pipeline

import (
	"errors"
	"strings"

	"github.com/json-iterator/go"
	"gitlab.mobvista.com/ADN/adnet/internal/errorcode"
	"gitlab.mobvista.com/ADN/adnet/internal/extractor"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	"gitlab.mobvista.com/ADN/adnet/internal/output"
)

type HupuAdRenderFilter struct{}

func (harf *HupuAdRenderFilter) Process(data interface{}) (interface{}, error) {
	in, ok := data.(*mvutil.ReqCtx)
	if !ok {
		return nil, errors.New("HupuAdRenderFilter input type should be *mvutil.ReqCtx")
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
			//todo
			creative, _ := extractor.GetCREATIVE_CHECK_HUPU_ADX_CREATIVE_IDS()
			dict, ok := output.RenderHupuRes(*result, in.ReqParams, creative)
			if !ok {
				return nil, errorcode.EXCEPTION_RETURN_EMPTY
			}
			res, err := json.Marshal(dict)
			if err != nil {
				return nil, errorcode.EXCEPTION_RETURN_EMPTY
			}
			resultJson := strings.TrimRight(string(res), "\n")
			return &resultJson, nil
		}
	}
	return nil, errorcode.EXCEPTION_RETURN_EMPTY
}

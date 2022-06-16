package process_pipeline

import (
	"errors"
	"strings"

	jsoniter "github.com/json-iterator/go"
	"gitlab.mobvista.com/ADN/adnet/internal/errorcode"
	"gitlab.mobvista.com/ADN/adnet/internal/extractor"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	"gitlab.mobvista.com/ADN/adnet/internal/output"
)

type MangoAdRenderFilter struct {
}

func (this *MangoAdRenderFilter) Process(data interface{}) (interface{}, error) {
	in, ok := data.(*mvutil.ReqCtx)
	if !ok {
		return nil, errors.New("MangoAdRenderFilter input type should be *mvutil.ReqCtx")
	}
	if in.Result != nil {
		result, err := output.RenderOutput(in.ReqParams, in.Result)
		// 无论是否有err 都打request日志
		mvutil.StatRequestLog(&in.ReqParams.Param)
		if err != nil {
			//fmt.Println(err)
			res, err := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(output.RenderNoRes(in.ReqParams))
			if err != nil {
				return nil, errorcode.EXCEPTION_RETURN_EMPTY
			}
			resultJson := strings.TrimRight(string(res), "\n")
			return &resultJson, nil
		}
		if result != nil && len(result.Data.Ads) > 0 {
			creative := extractor.GetCREATIVE_CHECK_MANGGUO_ADX_CREATIVE_IDS()
			res, err := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(output.RenderMangoRes(*result, in.ReqParams, creative))
			if err != nil {
				return nil, errorcode.EXCEPTION_RETURN_EMPTY
			}
			resultJson := strings.TrimRight(string(res), "\n")
			return &resultJson, nil
		}
	}
	return nil, errorcode.EXCEPTION_RETURN_EMPTY
}

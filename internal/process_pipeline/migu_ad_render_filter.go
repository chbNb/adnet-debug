package process_pipeline

import (
	"encoding/json"
	"errors"
	"strings"

	"gitlab.mobvista.com/ADN/adnet/internal/errorcode"
	"gitlab.mobvista.com/ADN/adnet/internal/extractor"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	"gitlab.mobvista.com/ADN/adnet/internal/output"
)

type MIGUAdRenderFilter struct {
}

func (mgarf *MIGUAdRenderFilter) Process(data interface{}) (interface{}, error) {
	in, ok := data.(*mvutil.ReqCtx)
	if !ok {
		return nil, errors.New("MIGUAdRenderFilter input type should be *mvutil.ReqCtx")
	}
	if in.Result != nil {
		result, err := output.RenderOutput(in.ReqParams, in.Result)
		// 无论是否有err 都打request日志
		mvutil.StatRequestLog(&in.ReqParams.Param)
		//降填充日志
		mvutil.StatReduceFillLog(&in.ReqParams.Param)
		if err != nil {
			//fmt.Println(err)
			return nil, errorcode.EXCEPTION_RETURN_EMPTY
		}
		//fmt.Println(result)
		if result != nil && len(result.Data.Ads) > 0 {
			creative, ok := extractor.GetCREATIVE_CHECK_MIGU_ADX_CREATIVE_IDS()
			if !ok {
				return nil, errorcode.EXCEPTION_RETURN_EMPTY
			}
			dict, ok := output.RenderMIGURes(*result, in.ReqParams, creative)
			if !ok {
				return nil, errorcode.EXCEPTION_RETURN_EMPTY
			}
			res, err := json.Marshal(dict)
			//res, err := json.Marshal(output.RenderMIGURes(*result, in.ReqParams))
			if err != nil {
				//fmt.Println(err)
				return nil, errorcode.EXCEPTION_RETURN_EMPTY
			}
			//fmt.Println(string(res))
			resultJson := strings.TrimRight(string(res), "\n")
			return &resultJson, nil
		}
	}
	return nil, errorcode.EXCEPTION_RETURN_EMPTY
}

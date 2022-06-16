package process_pipeline

import (
	"errors"
	"strings"

	jsoniter "github.com/json-iterator/go"
	"gitlab.mobvista.com/ADN/adnet/internal/errorcode"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	"gitlab.mobvista.com/ADN/adnet/internal/output"
)

type IFENGAdRenderFilter struct {
}

func (ifarf *IFENGAdRenderFilter) Process(data interface{}) (interface{}, error) {
	in, ok := data.(*mvutil.ReqCtx)
	if !ok {
		return nil, errors.New("IFENGAdRenderFilter input type should be *mvutil.ReqCtx")
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
		if result != nil && len(result.Data.Ads) > 0 {
			resOnline := output.RenderOnlineRes(*result, in.ReqParams)
			vastxml, err := output.RanderVastData(in.ReqParams, resOnline)
			if err != nil {
				//fmt.Println(err)
				mvutil.Logger.Runtime.Warnf("IFENGAdRenderFilter get RenderVastData error error code=[%s]", err)
				return nil, errorcode.EXCEPTION_RETURN_EMPTY
			}
			resJson := output.RenderIFENGRes(*result, in.ReqParams, vastxml)
			res, err := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(resJson)
			if err != nil {
				return nil, errorcode.EXCEPTION_RETURN_EMPTY
			}
			resultJson := strings.TrimRight(string(res), "\n")
			return &resultJson, nil
		}
	}
	return nil, errorcode.EXCEPTION_RETURN_EMPTY
}

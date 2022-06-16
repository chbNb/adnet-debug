package process_pipeline

import (
	"errors"
	"strings"

	"github.com/gogo/protobuf/proto"
	"gitlab.mobvista.com/ADN/adnet/internal/errorcode"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	"gitlab.mobvista.com/ADN/adnet/internal/output"
)

type KSYUNAdRenderFilter struct {
}

func (ksarf *KSYUNAdRenderFilter) Process(data interface{}) (interface{}, error) {
	in, ok := data.(*mvutil.ReqCtx)
	if !ok {
		return nil, errors.New("KSYUNAdRenderFilter input type should be *mvutil.ReqCtx")
	}
	if in.Result != nil {
		result, err := output.RenderOutput(in.ReqParams, in.Result)
		// 无论是否有err 都打request日志
		mvutil.StatRequestLog(&in.ReqParams.Param)
		//降填充日志
		mvutil.StatReduceFillLog(&in.ReqParams.Param)
		if err != nil {
			return nil, errorcode.EXCEPTION_RETURN_EMPTY
		}

		if result != nil && len(result.Data.Ads) != 0 {
			KSYUNResult := output.RenderKSYUNRes(*result, in.ReqParams)
			res, err := proto.Marshal(&KSYUNResult)
			if err != nil {
				return nil, errorcode.EXCEPTION_RETURN_EMPTY
			}
			resultProto := strings.TrimRight(string(res), "\n")
			return &resultProto, nil
		}
	}
	return nil, errorcode.EXCEPTION_RETURN_EMPTY
}

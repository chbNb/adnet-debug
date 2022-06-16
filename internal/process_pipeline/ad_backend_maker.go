package process_pipeline

import (
	"errors"

	"gitlab.mobvista.com/ADN/adnet/internal/mvconst"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
)

type AdBackendMaker struct {
}

func (abm *AdBackendMaker) Process(data interface{}) (interface{}, error) {
	// 类型断言
	in, ok := data.(*mvutil.ReqCtx)
	if !ok {
		return nil, errors.New("TrafficAllotFilter input type should be *mvutil.ReqCtx")
	}

	orderedbackends := make([]int, 0)

	if mvutil.IsRequestPioneerDirectly(&in.ReqParams.Param) {
		in.Backends[mvconst.Pioneer] = mvutil.NewMobvistaCtx()
		orderedbackends = append(orderedbackends, mvconst.Pioneer)
	} else {
		in.Backends[mvconst.Mobvista] = mvutil.NewMobvistaCtx()
		orderedbackends = append(orderedbackends, mvconst.Mobvista)
	}

	in.ReqParams.FlowTagID = mvconst.FlowTagDefault
	in.OrderedBackends = orderedbackends
	return in, nil
}

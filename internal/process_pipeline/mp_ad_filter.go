package process_pipeline

import (
	"errors"
	"gitlab.mobvista.com/ADN/adnet/internal/extractor"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
)

type MpAdFilter struct {
}

func (mpnf *MpAdFilter) Process(data interface{}) (interface{}, error) {
	in, ok := data.(*mvutil.RequestParams)
	if !ok {
		return nil, errors.New("MpAdFilter input type should be *params.RequestParams")
	}

	unitInfo, ifFind := extractor.GetUnitInfo(in.Param.UnitID)
	if ifFind && unitInfo.MPToMV != nil {
		mpmv := *unitInfo.MPToMV
		in.Param.AppID = int64(mpmv.AppId)
		in.Param.UnitID = int64(mpmv.UnitId)
		in.Param.Sign = "NO_CHECK_SIGN"
	}
	return in, nil
}

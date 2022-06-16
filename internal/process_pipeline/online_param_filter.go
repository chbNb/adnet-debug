package process_pipeline

import (
	"errors"
	"strconv"

	"gitlab.mobvista.com/ADN/adnet/internal/errorcode"
	"gitlab.mobvista.com/ADN/adnet/internal/extractor"
	"gitlab.mobvista.com/ADN/adnet/internal/mvconst"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	"gitlab.mobvista.com/ADN/adnet/internal/output"
)

type OnlineParamFilter struct {
}

func (opf *OnlineParamFilter) Process(data interface{}) (interface{}, error) {
	in, ok := data.(*mvutil.RequestParams)
	if !ok {
		// errn := errors.New("OnlineParamFilter input type should be *params.RequestParams")
		// errorDetail := errorcode.ErrorDetail{errorcode.EXCEPTION_PARAMS_ERROR, 0, &errn}
		// return nil, errorDetail
		return nil, errors.New("OnlineParamFilter input type should be *params.RequestParams")
	}
	// online api 参数强校验
	err := checkParams(in)
	if err != nil {
		return nil, err
	}
	in.Param.RequestType = mvconst.REQUEST_TYPE_OPENAPI_AD
	in.Param.TNum = int(in.Param.AdNum)
	in.Param.AdNum = int32(40)
	in.Param.OnlyImpression = 1
	// 针对gameloft使传来的price_floor失效
	delPriceFloorUnits, ifFind := extractor.GetDEL_PRICE_FLOOR_UNIT()
	if ifFind {
		if mvutil.InInt64Arr(in.Param.UnitID, delPriceFloorUnits) {
			in.Param.PriceFloor = 0
		}
	}
	// 韩国开发者
	if output.IsAfreecatvUnit(in.Param.UnitID) {
		if in.Param.Platform == mvconst.PlatformAndroid {
			in.Param.GAID = in.Param.AID
		} else {
			in.Param.IDFA = in.Param.AID
		}
	}
	return in, nil
}

func renderOffset(r *mvutil.RequestParams) {
	// online publisher维度offset配置生效
	if r.Param.Offset > int32(0) {
		return
	}
	offsetList := r.PublisherInfo.OffsetList
	if len(offsetList) <= 0 {
		return
	}
	randArr := make(map[int]int, len(offsetList))
	for k, v := range offsetList {
		kInt, err := strconv.Atoi(k)
		if err != nil {
			continue
		}
		vInt := int(v)
		randArr[kInt] = vInt
	}
	if len(randArr) <= 0 {
		return
	}
	randOffset := mvutil.RandByRate(randArr)
	r.Param.Offset = int32(randOffset)
}

func checkParams(r *mvutil.RequestParams) error {
	// 小程序请求ads接口不做参数强校验
	if r.Param.NCP == mvconst.NO_CHECK_PARAMS {
		return nil
	}
	// 配置了开关及不做校验的appid
	noCheckApp, _ := extractor.GetNO_CHECK_PARAM_APP()
	// 开关没开，则不做校验
	if noCheckApp.Status == nil || *noCheckApp.Status != 1 {
		return nil
	}
	// 针对某些app不做校验
	if noCheckApp.AppIds == nil || mvutil.InInt64Arr(r.Param.AppID, *noCheckApp.AppIds) {
		return nil
	}
	if r.Param.Platform == mvconst.PlatformOther {
		// mvutil.Logger.Runtime.Warnf("request_id=[%s] has filter by platform is validate.", r.Param.RequestID)
		return errorcode.EXCEPTION_APP_PLATFORM_ERROR
	}
	if len(r.Param.OSVersion) == 0 {
		// mvutil.Logger.Runtime.Warnf("request_id=[%s] has filter by os_version is validate.", r.Param.RequestID)
		return errorcode.EXCEPTION_SERVICE_REQUEST_OS_VERSION_REQUIRED
	}
	// if len(r.Param.IMEI) == 0 && r.Param.CountryCode == "CN" {
	//	//mvutil.Logger.Runtime.Warnf("request_id=[%s] has filter by imei is validate.", r.Param.RequestID)
	//	return errorcode.EXCEPTION_IMEI_EMPTY
	// }
	// if r.Param.Brand == "0" {
	//	//mvutil.Logger.Runtime.Warnf("request_id=[%s] has filter by brand is validate.", r.Param.RequestID)
	//	return errorcode.EXCEPTION_BRAND_EMPTY
	// }
	// if r.Param.Model == "0" {
	//	//mvutil.Logger.Runtime.Warnf("request_id=[%s] has filter by model is validate.", r.Param.RequestID)
	//	return errorcode.EXCEPTION_MODEL_EMPTY
	// }
	// if r.Param.Platform == mvconst.PlatformAndroid && len(r.Param.AndroidID) == 0 {
	//	//mvutil.Logger.Runtime.Warnf("request_id=[%s] has filter by android_id is validate.", r.Param.RequestID)
	//	return errorcode.EXCEPTION_ANDROIDID_EMPTY
	// }
	// if r.Param.NetworkType == mvconst.NETWORK_TYPE_UNKNOWN {
	//	//mvutil.Logger.Runtime.Warnf("request_id=[%s] has filter by network_type is validate.", r.Param.RequestID)
	//	return errorcode.EXCEPTION_NETWORK_TYPE_EMPTY
	// }
	if r.Param.Platform == mvconst.PlatformIOS && len(r.Param.IDFA) == 0 {
		// mvutil.Logger.Runtime.Warnf("request_id=[%s] has filter by idfa is validate.", r.Param.RequestID)
		return errorcode.EXCEPTION_IDFA_EMPTY
	}
	return nil
}

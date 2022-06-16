package process_pipeline

import (
	"errors"
	"strconv"
	"time"

	"gitlab.mobvista.com/ADN/adnet/internal/errorcode"
	"gitlab.mobvista.com/ADN/adnet/internal/geo"
	"gitlab.mobvista.com/ADN/adnet/internal/mvconst"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	"gitlab.mobvista.com/ADN/exporter/metrics"
)

var RedisBucket chan time.Time

type IpInfoFilter struct {
}

// 调用geo服务，通过IP地址解析出用户所在国家、城市和区域信息
func (iif *IpInfoFilter) Process(data interface{}) (interface{}, error) {
	in, ok := data.(*mvutil.RequestParams)
	if !ok {
		return nil, errors.New("IpInfoFilter input type should be *params.RequestParams")
	}

	key := in.Param.ClientIP
	if len(key) == 0 {
		return mvconst.EXCEPTION_SERVICE_REQUEST_COUNTRY_CODE_REQUIRED, errors.New("EXCEPTION_SERVICE_REQUEST_COUNTRY_CODE_REQUIRED")
	}

	metrics.IncCounterDyWithValues("ip_info_filter-ipcache-rejectcode", "query", "")
	val, hit, err := geo.GetGeo(in.Param.ClientIP)
	if err != nil {
		mvutil.Logger.Runtime.Warnf("query ip :%s error: %s", in.Param.ClientIP, err.Error())
		metrics.IncCounterDyWithValues("ip_info_filter-ipcache-rejectcode", "", strconv.Itoa(mvconst.EXCEPTION_SERVICE_REQUEST_COUNTRY_CODE_REQUIRED))
		return nil, errorcode.EXCEPTION_SERVICE_REQUEST_COUNTRY_CODE_REQUIRED
	}

	if hit {
		metrics.IncCounterDyWithValues("ip_info_filter-ipcache-rejectcode", "hit", "")
	}
	in.Param.CountryCode = val.TwoLetterCountry
	in.Param.CityCode = int64(val.CityCode)
	in.Param.CityString = val.City
	in.Param.RegionString = val.Region
	return in, nil
}

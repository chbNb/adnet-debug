package bid_filters

import (
	"strconv"

	"gitlab.mobvista.com/ADN/adnet/internal/geo"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/filter"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/helpers"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/req_context"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	"gitlab.mobvista.com/ADN/exporter/metrics"
)

type AreaTargetFilter struct {
}

// 根据ip查询详细地址信息
func (atf *AreaTargetFilter) Process(data interface{}) (interface{}, error) {
	in, ok := data.(*mvutil.RequestParams)
	if !ok {
		return nil, filter.AreaTargetFilterInputError
	}
	if len(in.Param.ClientIP) == 0 || !helpers.IsCorrectIp(in.Param.ClientIP) {
		return in, filter.ClientIpInvalidate
	}

	metrics.IncCounterDyWithValues("ip_info_filter-ipcache-rejectcode", "query", "")
	// 根据ip查geo,先查redis缓存，再通过Netacuity查，最后放缓存
	val, hit, err := geo.GetGeo(in.Param.ClientIP)
	if err != nil {
		req_context.GetInstance().MLogs.Runtime.Errorf("query ip %s error: %s", in.Param.ClientIP, err.Error())
		metrics.IncCounterDyWithValues("ip_info_filter-ipcache-rejectcode", "", strconv.Itoa(filter.QueryNetServiceError.Int()))
		return in, filter.QueryNetServiceError
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

package query_mem_filters

import (
	"errors"
	"strconv"
	"strings"

	"github.com/json-iterator/go"
	"gitlab.mobvista.com/ADN/adnet/internal/extractor"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/params"
	"gitlab.mobvista.com/ADN/treasure_box_sdk/tb_tools"
)

var QueryMemDataFilterInputError = errors.New("query_mem_data_filter input error")

// var TokenDataInvalidate = errors.New("token data is invalidate")

type QueryMemDataFilter struct {
}

func (qmdf *QueryMemDataFilter) Process(data interface{}) (interface{}, error) {
	in, ok := data.(*params.HttpReqData)
	if !ok {
		return nil, QueryMemDataFilterInputError
	}

	table := in.QueryData.GetString("table", true)
	key := in.QueryData.GetString("key", true)
	// unitId := in.QueryData.GetString("unit_id", true)
	// countryCode := in.QueryData.GetString("cc", true)
	if len(table) == 0 || len(key) == 0 {
		return nil, errors.New("query data processor input param has no table or key")
	}

	dataTable := strings.ToLower(table)

	var interKey int64
	if inStrArray(dataTable, []string{"configcenter", "advertiser", "app", "publisher", "unit", "campaign", "adx_dsp_config"}) {
		ikey, err := strconv.ParseInt(key, 10, 64)
		if err != nil {
			return nil, err
		}
		interKey = ikey
	}

	var info interface{}
	var ifFind bool
	switch dataTable {
	case "app":
		info, ifFind = extractor.GetAppInfo(interKey)
		if !ifFind {
			return nil, errors.New("app has no key:" + key)
		}
	case "unit":
		info, ifFind = extractor.GetUnitInfo(interKey)
		if !ifFind {
			return nil, errors.New("unit has no key:" + key)
		}
	case "publisher":
		info, ifFind = extractor.GetPublisherInfo(interKey)
		if !ifFind {
			return nil, errors.New("publisher has no key:" + key)
		}
	case "campaign":
		info, ifFind = extractor.GetCampaignInfo(interKey)
		if !ifFind {
			return nil, errors.New("campaign has no key:" + key)
		}
	case "configcenter":
		info, ifFind = extractor.GetConfigcenter(key)
		if !ifFind {
			return nil, errors.New("configcenter has no key:" + key)
		}
	case "advertiser":
		info, ifFind = extractor.GetAdvertiserInfo(interKey)
		if !ifFind {
			return nil, errors.New("advertiser has no key:" + key)
		}
	case "adx_dsp_config":
		info, ifFind = extractor.GetDspConfig(interKey)
		if !ifFind {
			return nil, errors.New("campaign has no key:" + key)
		}
	case "config_algorithm_fillrate":
		info, ifFind = extractor.GetFillRateControllConfig(key)
		if !ifFind {
			return nil, errors.New("campaign has no key:" + key)
		}
	case "getfreq_control_config":
		if info = extractor.GetFREQ_CONTROL_CONFIG(); info == nil {
			ifFind = false
		}
		if !ifFind {
			return nil, errors.New("GetFREQ_CONTROL_CONFIG has no key:" + key)
		}
	case "gettimezone_config":
		if info = extractor.GetTIMEZONE_CONFIG(); info == nil {
			ifFind = false
		}
		if !ifFind {
			return nil, errors.New("GetFREQ_CONTROL_CONFIG has no key:" + key)
		}
	case "getcountry_code_timezone_config":
		if info = extractor.GetCOUNTRY_CODE_TIMEZONE_CONFIG(); info == nil {
			ifFind = false
		}
		if !ifFind {
			return nil, errors.New("GetFREQ_CONTROL_CONFIG has no key:" + key)
		}
	case "getfreqcontrolpricefactor":
		info, ifFind = extractor.GetFreqControlPriceFactor(key)
		if !ifFind {
			return nil, errors.New("GetFreqControlPriceFactor has no key:" + key)
		}
	case "config_hbadxendpoint":
		info, ifFind = extractor.GetHBAdxEndpoint(key)
		if !ifFind {
			return nil, errors.New("GetHBAdxEndpoint has no key:" + key)
		}
	case "config_hbadxendpointv2":
		k := strings.Split(key, ":")
		cloud := k[0]
		region := k[1]
		info, ifFind = extractor.GetHBAdxEndpointV2(cloud, region)
		if !ifFind {
			return nil, errors.New("GetHBAdxEndpointV2 has no key:" + key)
		}
	case "config":
		info, ifFind = extractor.GetMVConfigValue(key)
		if !ifFind {
			return nil, errors.New("GetMVConfigValue has no key:" + key)
		}
	case "tb_tools":
		module := in.QueryData.GetString("module", true)
		var err error
		info, ifFind, err = tb_tools.GetDataByStringFromMemory(module, key)
		if err != nil {
			return nil, err
		}
		if !ifFind {
			return nil, errors.New("tb_tools.GetDataByStringFromMemory(" + module + ", " + key + ") has no key")
		}
	default:
		return nil, errors.New("table is no support. dataTable: " + dataTable)
	}

	return jsoniter.ConfigCompatibleWithStandardLibrary.MarshalToString(info)
}

func inStrArray(val string, array []string) bool {
	for _, v := range array {
		if val == v {
			return true
		}
	}
	return false
}

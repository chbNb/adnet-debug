package process_pipeline

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/json-iterator/go"
	"gitlab.mobvista.com/ADN/adnet/internal/extractor"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	"gitlab.mobvista.com/ADN/treasure_box_sdk/tb_tools"
)

type QueryMemDataFilter struct {
}

func (qmdf *QueryMemDataFilter) Process(data interface{}) (interface{}, error) {
	in, ok := data.(*http.Request)
	if !ok {
		return nil, errors.New("QueryMemDataFilter input type should be *http.Request")
	}

	table := in.URL.Query().Get("table")
	key := in.URL.Query().Get("key")
	if len(table) == 0 || len(key) == 0 {
		return nil, errors.New("query data processor input param has no table or key")
	}

	dataTable := strings.ToLower(table)

	var interKey int64
	if mvutil.InStrArray(dataTable, []string{
		"app", "publisher", "unit", "placement", "campaign",
		"check_campaign", "check_campaign_fake",
	}) {
		ikey, err := strconv.ParseInt(key, 10, 64)
		if err != nil {
			return nil, err
		}
		interKey = ikey
	}

	var jsonData string
	switch dataTable {
	case "app":
		info, ifFind := extractor.GetAppInfo(interKey)
		if !ifFind {
			return nil, errors.New("app has no key:" + key)
		}
		jsonData, _ = jsoniter.ConfigCompatibleWithStandardLibrary.MarshalToString(info)
	case "publisher":
		info, ifFind := extractor.GetPublisherInfo(interKey)
		if !ifFind {
			return nil, errors.New("publisher has no key:" + key)
		}
		jsonData, _ = jsoniter.ConfigCompatibleWithStandardLibrary.MarshalToString(info)
	case "unit":
		info, ifFind := extractor.GetUnitInfo(interKey)
		if !ifFind {
			return nil, errors.New("unit has no key:" + key)
		}
		jsonData, _ = jsoniter.ConfigCompatibleWithStandardLibrary.MarshalToString(info)
	case "placement":
		info, ifFind := extractor.GetPlacementInfo(interKey)
		if !ifFind {
			return nil, errors.New("placement has no key:" + key)
		}
		jsonData, _ = jsoniter.ConfigCompatibleWithStandardLibrary.MarshalToString(info)
	case "campaign":
		info, ifFind := extractor.GetCampaignInfo(interKey)
		if !ifFind {
			return nil, errors.New("campaign has no key:" + key)
		}
		jsonData, _ = jsoniter.ConfigCompatibleWithStandardLibrary.MarshalToString(info)
	case "check_campaign":
		info := extractor.CheckGetCampaignInfo(interKey)
		jsonData, _ = jsoniter.ConfigCompatibleWithStandardLibrary.MarshalToString(info)
	case "freq_control_factor":
		info, ifFind := extractor.GetFreqControlPriceFactor(in.URL.Query().Get("key"))
		if !ifFind {
			return nil, errors.New("freq_control_factor has no key:" + key)
		}
		jsonData, _ = jsoniter.ConfigCompatibleWithStandardLibrary.MarshalToString(info)
	case "getfreq_control_config":
		value := extractor.GetFREQ_CONTROL_CONFIG()
		res, err := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(value)
		if err != nil {
			return nil, err
		}
		jsonData = string(res)
	case "getfillratecontrollconfig":
		info, ifFind := extractor.GetFillRateControllConfig(in.URL.Query().Get("key"))
		if !ifFind {
			return nil, errors.New("getfillratecontrollconfig has no key:" + key)
		}
		jsonData, _ = jsoniter.ConfigCompatibleWithStandardLibrary.MarshalToString(info)

	case "getsspprofitdistributionrule":
		info, ifFind := extractor.GetSspProfitDistributionRule(key)
		if !ifFind {
			return nil, errors.New("getsspprofitdistributionrule has no key:" + key)
		}
		jsonData, _ = jsoniter.ConfigCompatibleWithStandardLibrary.MarshalToString(info)

	case "getadxtrafficmediaconfigbykey":
		//curl "http://127.0.0.1:9099/query_mem?table=getadxtrafficmediaconfigbykey&key=2-21-ALL"
		info, ifFind := extractor.GetAdxTrafficMediaConfigByKey(key)
		if !ifFind {
			return nil, errors.New("getsspprofitdistributionrule has no key:" + key)
		}
		jsonData, _ = jsoniter.ConfigCompatibleWithStandardLibrary.MarshalToString(info)
	case "gettreasureboxidcount":
		info, ifFind := extractor.GetTreasureBoxIDCount()
		if !ifFind {
			return nil, errors.New("GetTreasureBoxIDCount has no key:" + key)
		}
		jsonData, _ = jsoniter.ConfigCompatibleWithStandardLibrary.MarshalToString(info)
	case "tb_tools":
		module := in.URL.Query().Get("module")
		info, ifFind, err := tb_tools.GetDataByStringFromMemory(module, key)
		if err != nil {
			return nil, err
		}
		if !ifFind {
			return nil, errors.New("tb_tools.GetDataByStringFromMemory(" + module + ", " + key + ") has no key")
		}
		jsonData, _ = jsoniter.ConfigCompatibleWithStandardLibrary.MarshalToString(info)
	default:
		return nil, errors.New("table is no support. dataTable: " + dataTable)
	}

	return &jsonData, nil
}

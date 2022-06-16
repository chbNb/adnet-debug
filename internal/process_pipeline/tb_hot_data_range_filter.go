package process_pipeline

import (
	"errors"
	jsoniter "github.com/json-iterator/go"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/params"
	"gitlab.mobvista.com/ADN/treasure_box_sdk/tb_tools"
	"net/http"
	"strconv"
	"time"
)

type HotDataWithRangeFilter struct {
}

func (f *HotDataWithRangeFilter) Process(data interface{}) (interface{}, error) {
	var module, table, start, end string
	switch in := data.(type) {
	case *http.Request:
		module = in.URL.Query().Get("module")
		table = in.URL.Query().Get("table")
		start = in.URL.Query().Get("start")
		end = in.URL.Query().Get("end")
	case *params.HttpReqData:
		module = in.QueryData.GetString("module", true)
		table = in.QueryData.GetString("table", true)
		start = in.QueryData.GetString("start", true)
		end = in.QueryData.GetString("end", true)
	default:
		return nil, errors.New("unknown type of data")
	}

	now := time.Now()
	if len(end) == 0 {
		end = strconv.FormatInt(now.Unix(), 10)
	}

	if len(start) == 0 {
		start = strconv.FormatInt(now.Add(-24*time.Hour).Unix(), 10)
	}

	if _, err := strconv.ParseInt(start, 10, 64); err != nil {
		return nil, err
	}
	if _, err := strconv.ParseInt(end, 10, 64); err != nil {
		return nil, err
	}

	moduleKey := module + ":" + table
	res, err := tb_tools.GetHotDataWithRange(moduleKey, start, end)

	if err != nil {
		return nil, err
	}

	var jsonData string
	jsonData, _ = jsoniter.ConfigCompatibleWithStandardLibrary.MarshalToString(res)
	return &jsonData, nil
}

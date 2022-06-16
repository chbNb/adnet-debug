package process_pipeline

import (
	"errors"
	jsoniter "github.com/json-iterator/go"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/params"
	"gitlab.mobvista.com/ADN/treasure_box_sdk/tb_tools"
	"net/http"
)

type HotDataAllFilter struct {
}

func (f *HotDataAllFilter) Process(data interface{}) (interface{}, error) {
	var module, table string
	switch in := data.(type) {
	case *http.Request:
		module = in.URL.Query().Get("module")
		table = in.URL.Query().Get("table")
	case *params.HttpReqData:
		module = in.QueryData.GetString("module", true)
		table = in.QueryData.GetString("table", true)
	default:
		return nil, errors.New("unknown type of data")
	}

	moduleKey := module + ":" + table
	res, err := tb_tools.GetHotDataUnion(moduleKey)

	if err != nil {
		return nil, err
	}

	var jsonData string
	jsonData, _ = jsoniter.ConfigCompatibleWithStandardLibrary.MarshalToString(res)
	return &jsonData, nil
}

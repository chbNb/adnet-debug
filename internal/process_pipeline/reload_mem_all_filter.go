package process_pipeline

import (
	"gitlab.mobvista.com/ADN/adnet/internal/hb/params"
	"gitlab.mobvista.com/ADN/treasure_box_sdk/tb_tools"
	"net/http"
	"strconv"
)

type ReloadMemAllFilter struct {
}

func (f *ReloadMemAllFilter) Process(data interface{}) (interface{}, error) {
	bs := ""
	switch in := data.(type) {
	case *http.Request:
		bs = in.URL.Query().Get("size")
	case *params.HttpReqData:
		bs = in.QueryData.GetString("size", true)
	}

	batchSize := 1000
	if bs, err := strconv.Atoi(bs); bs > 0 && err == nil {
		batchSize = bs
	}

	var jsonData string
	if err := tb_tools.ReloadAll(batchSize); err != nil {
		return nil, err
	}
	jsonData = "reload ok"

	return &jsonData, nil
}

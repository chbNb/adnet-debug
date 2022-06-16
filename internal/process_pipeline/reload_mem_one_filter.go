package process_pipeline

import (
	"errors"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/params"
	"gitlab.mobvista.com/ADN/treasure_box_sdk/tb_tools"
	"net/http"
	"strconv"
)

type ReloadMemOneFilter struct {
}

func (f *ReloadMemOneFilter) Process(data interface{}) (interface{}, error) {
	moduleKey := ""
	bs := ""
	switch in := data.(type) {
	case *http.Request:
		moduleKey = in.URL.Query().Get("moduleKey")
		bs = in.URL.Query().Get("size")
	case *params.HttpReqData:
		moduleKey = in.QueryData.GetString("moduleKey", true)
		bs = in.QueryData.GetString("size", true)
	}

	if len(moduleKey) == 0 {
		return nil, errors.New("reload mem processor input param has no module key")
	}

	batchSize := 1000
	if bs, err := strconv.Atoi(bs); err == nil && bs > 0 {
		batchSize = bs
	}

	var jsonData string
	if err := tb_tools.ReloadOne(moduleKey, batchSize); err != nil {
		return nil, err
	}
	jsonData = "reload ok"

	return &jsonData, nil
}

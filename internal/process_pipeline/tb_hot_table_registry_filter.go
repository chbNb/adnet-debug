package process_pipeline

import (
	jsoniter "github.com/json-iterator/go"
	"gitlab.mobvista.com/ADN/treasure_box_sdk/tb_tools"
)

type HotTableRegistryFilter struct {
}

func (f *HotTableRegistryFilter) Process(data interface{}) (interface{}, error) {
	var jsonData string
	tableKeys := tb_tools.GetHotDataRegistry()
	jsonData, _ = jsoniter.ConfigCompatibleWithStandardLibrary.MarshalToString(tableKeys)
	return &jsonData, nil
}

package helpers

import (
	jsoniter "github.com/json-iterator/go"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/params"
)

func SerializeMidwayCreative(mwCreative []int) string {
	if len(mwCreative) == 0 {
		return "0"
	}
	creativeMap := make(map[int]int)
	for _, tempId := range mwCreative {
		if _, ok := creativeMap[tempId]; !ok {
			creativeMap[tempId] = 1
			continue
		}
		creativeMap[tempId] += 1
	}
	creative := params.MidwayCreative{TempID: creativeMap}
	creativeTmp, err := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(creative)
	if err != nil {
		return "0"
	}
	return string(creativeTmp)
}

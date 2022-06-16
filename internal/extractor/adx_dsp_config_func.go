package extractor

import (
	"errors"

	"github.com/easierway/concurrent_map"
	smodel "gitlab.mobvista.com/ADN/structs/model"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const DspStatusActive = 1

var dspConfigUpdateInputError = errors.New("dspConfigUpdateFunc failed, query or dbLoaderInfo is nil")

func NewAdxDspConfigExtractor(allFilter, incFilter, selectOptions bson.M) dbLoaderInfo {
	return dbLoaderInfo{
		collection:   "adx_dsp_config",
		querybmDbMap: allFilter,
		queryimDbMap: incFilter,
		querySelect:  selectOptions,
		updateFunc:   dspConfigUpdateFunc,
		dataCur:      concurrent_map.CreateConcurrentMap(concurrentMapPartitionsNum),
	}
}

func dspConfigUpdateFunc(query *mgo.Query, num int, dbLoaderInfo *dbLoaderInfo) (maxUpdate int64, err error) {
	if query == nil || dbLoaderInfo == nil {
		return 0, dspConfigUpdateInputError
	}
	adxDspConfig := &smodel.AdxDspConfig{}
	maxUpdate = 0
	item := query.Iter()
	var incData int
	for idx := 0; idx < num; idx++ {
		if item.Next(adxDspConfig) {
			if int64(adxDspConfig.Updated) > maxUpdate {
				maxUpdate = int64(adxDspConfig.Updated)
			}
			dspIncUpdateFunc(adxDspConfig, dbLoaderInfo)
			incData++
			adxDspConfig = &smodel.AdxDspConfig{}
		} else {
			if item.Timeout() {
				logger.Warn("dspConfigUpdateFunc item.Timeout()")
				break
			}
			// logger.Warnf("dspConfigUpdateFunc item.error=%s", item.Err().Error())
		}
	}

	if err := item.Close(); err != nil {
		logger.Error(err.Error())
	}

	if incData != num {
		logger.Warnf("dspConfigUpdateFunc Size not Equal incData=%d, QueryNum=%d", incData, num)
	}
	return maxUpdate, nil
}

func dspIncUpdateFunc(adxDspConfig *smodel.AdxDspConfig, dbLoaderInfo *dbLoaderInfo) {
	if adxDspConfig.Status != DspStatusActive {
		DbLoaderRegistry[TblAdxDspConfig].dataCur.Del(concurrent_map.I64Key(adxDspConfig.DspID))
		return
	}

	DbLoaderRegistry[TblAdxDspConfig].dataCur.Set(concurrent_map.I64Key(adxDspConfig.DspID), adxDspConfig)
}

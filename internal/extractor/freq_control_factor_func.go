package extractor

import (
	"errors"

	"github.com/easierway/concurrent_map"
	smodel "gitlab.mobvista.com/ADN/structs/model"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const FreqControlFactorActive = 1

var freqControlFactorInputError = errors.New("freqControlFactorFunc failed, query or dbLoaderInfo is nil")

func NewFreqControlFactorExtractor(allFilter, incFilter, selectOptions bson.M) dbLoaderInfo {
	return dbLoaderInfo{
		collection:   "freq_control_factor",
		querybmDbMap: allFilter,
		queryimDbMap: incFilter,
		querySelect:  selectOptions,
		updateFunc:   freqControlFactorFunc,
		dataCur:      concurrent_map.CreateConcurrentMap(concurrentMapPartitionsNum),
	}
}

func freqControlFactorFunc(query *mgo.Query, num int, dbLoaderInfo *dbLoaderInfo) (maxUpdate int64, err error) {
	if query == nil || dbLoaderInfo == nil {
		return 0, freqControlFactorInputError
	}
	freqControlFactorConfig := &smodel.FreqControlFactor{}
	maxUpdate = 0
	item := query.Iter()
	var incData int
	for idx := 0; idx < num; idx++ {
		if item.Next(freqControlFactorConfig) {
			if int64(freqControlFactorConfig.Updated) > maxUpdate {
				maxUpdate = int64(freqControlFactorConfig.Updated)
			}
			freqControlFactorIncUpdateFunc(freqControlFactorConfig, dbLoaderInfo)
			incData++
			freqControlFactorConfig = &smodel.FreqControlFactor{}
		} else {
			if item.Timeout() {
				logger.Warn("freqControlFactorFunc item.Timeout()")
				break
			}
			// logger.Warnf("freqControlFactorFunc item.error=%s", item.Err().Error())
		}
	}

	if err := item.Close(); err != nil {
		logger.Error(err.Error())
	}

	if incData != num {
		logger.Warnf("freqControlFactorFunc Size not Equal incData=%d, QueryNum=%d", incData, num)
	}
	return maxUpdate, nil

}

func freqControlFactorIncUpdateFunc(freqControlFactorConfig *smodel.FreqControlFactor, dbLoaderInfo *dbLoaderInfo) {
	if freqControlFactorConfig.Status != FreqControlFactorActive {
		DbLoaderRegistry[TblFreqControlFactor].dataCur.Del(concurrent_map.StrKey(freqControlFactorConfig.FactorKey))
		return
	}

	DbLoaderRegistry[TblFreqControlFactor].dataCur.Set(concurrent_map.StrKey(freqControlFactorConfig.FactorKey), freqControlFactorConfig)
}

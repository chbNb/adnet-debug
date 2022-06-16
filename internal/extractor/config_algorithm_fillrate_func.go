package extractor

import (
	"errors"

	smodel "gitlab.mobvista.com/ADN/structs/model"

	"github.com/easierway/concurrent_map"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var fillrateUpdateInputError = errors.New("fillrateUpdateFunc failed, query or dbLoaderInfo is nil")

func NewConfigAlgorithmFillRateExtractor(allFilter, incFilter, selectOptions bson.M) dbLoaderInfo {
	return dbLoaderInfo{
		collection:   "config_algorithm_fillrate",
		querybmDbMap: allFilter,
		queryimDbMap: incFilter,
		querySelect:  selectOptions,
		updateFunc:   fillrateUpdateFunc,
		dataCur:      concurrent_map.CreateConcurrentMap(concurrentMapPartitionsNum),
	}
}

func fillrateUpdateFunc(query *mgo.Query, num int, dbLoaderInfo *dbLoaderInfo) (maxUpdate int64, err error) {
	if query == nil || dbLoaderInfo == nil {
		return 0, fillrateUpdateInputError
	}
	ConfigAlgorithmFillRate := &smodel.ConfigAlgorithmFillRate{}
	maxUpdate = 0
	item := query.Iter()
	var incData int
	for idx := 0; idx < num; idx++ {
		if item.Next(ConfigAlgorithmFillRate) {
			if int64(ConfigAlgorithmFillRate.Updated) > maxUpdate {
				maxUpdate = int64(ConfigAlgorithmFillRate.Updated)
			}
			fillrateIncUpdateFunc(ConfigAlgorithmFillRate)
			incData++
			ConfigAlgorithmFillRate = &smodel.ConfigAlgorithmFillRate{}
		} else {
			if item.Timeout() {
				logger.Warn("fillrateUpdateFunc item.Timeout()")
				break
			}
			// logger.Warnf("fillrateUpdateFunc item.error=%s", item.Err().Error())
		}
	}

	if err := item.Close(); err != nil {
		logger.Error(err.Error())
	}

	if incData != num {
		logger.Warnf("fillrateUpdateFunc Size not Equal incData=%d, QueryNum=%d", incData, num)
	}
	return maxUpdate, nil

}

func fillrateIncUpdateFunc(configAlgorithmFillRate *smodel.ConfigAlgorithmFillRate) {
	if !legalFillRate(configAlgorithmFillRate) {
		DbLoaderRegistry[TblConfigAlgorithmFillrate].dataCur.Del(concurrent_map.StrKey(configAlgorithmFillRate.UniqueKey))
		return
	}
	DbLoaderRegistry[TblConfigAlgorithmFillrate].dataCur.Set(concurrent_map.StrKey(configAlgorithmFillRate.UniqueKey),
		configAlgorithmFillRate)
}

func legalFillRate(configAlgorithmFillRate *smodel.ConfigAlgorithmFillRate) bool {
	return configAlgorithmFillRate.Status == mvutil.ACTIVE
}

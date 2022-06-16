package extractor

import (
	"errors"

	smodel "gitlab.mobvista.com/ADN/structs/model"

	"github.com/easierway/concurrent_map"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	Batch    = 4000
	Prefetch = 0.25
)

var unitUpdateInputError = errors.New("unitUpdateFunc failed, query or dbLoaderInfo is nil")

func NewUnitExtractor(allFilter, incFilter, selectOptions bson.M) dbLoaderInfo {
	return dbLoaderInfo{
		collection:   "unit",
		querybmDbMap: allFilter,
		queryimDbMap: incFilter,
		querySelect:  selectOptions,
		updateFunc:   unitUpdateFunc,
		dataCur:      concurrent_map.CreateConcurrentMap(concurrentMapPartitionsNum),
	}
}

func unitUpdateFunc(query *mgo.Query, num int, dbLoaderInfo *dbLoaderInfo) (maxUpdate int64, err error) {
	if query == nil || dbLoaderInfo == nil {
		return 0, unitUpdateInputError
	}
	unitInfo := &smodel.UnitInfo{}
	maxUpdate = 0
	item := query.Iter()
	var incData, updData int
	for idx := 0; idx < num; idx++ {
		if item.Next(unitInfo) {
			if int64(unitInfo.Updated) > maxUpdate {
				maxUpdate = int64(unitInfo.Updated)
			}
			incData++
			if meta.config.UseExpiredMap { // 新方式的时候才需要查
				if _, find := DbLoaderRegistry[TblUnit].dataCur.Get(concurrent_map.I64Key(unitInfo.UnitId)); !find && !dbLoaderInfo.updateByGetAll {
					continue // 找不到  且  是getInc  的就不用更新
				}
			}
			updData++
			unitIncUpdateFunc(unitInfo, dbLoaderInfo)
			unitInfo = &smodel.UnitInfo{}
		} else {
			if item.Timeout() {
				logger.Warn("unitUpdateFunc item.Timeout()")
				break
			}
			// logger.Warnf("unitUpdateFunc item.error=%s", item.Err().Error())
		}
	}

	if err := item.Close(); err != nil {
		logger.Error(err.Error())
	}

	if incData != num {
		logger.Warnf("unitUpdateFunc Size not Equal incData=%d, QueryNum=%d, UpdateNum=%d", incData, num, updData)
	}
	logger.Infof("unitUpdateFunc incData=%d, QueryNum=%d, UpdateNum=%d", incData, num, updData)
	return maxUpdate, nil

}

func unitIncUpdateFunc(unitInfo *smodel.UnitInfo, dbLoaderInfo *dbLoaderInfo) {
	if !legalUnit(unitInfo) {
		if meta.config.UseExpiredMap { // 新的剔除需要保存为 nil
			DbLoaderRegistry[TblUnit].dataCur.Set(concurrent_map.I64Key(unitInfo.UnitId), nil)
		} else {
			DbLoaderRegistry[TblUnit].dataCur.Del(concurrent_map.I64Key(unitInfo.UnitId))
		}
		return
	}

	if len(unitInfo.AdSourceCountry) >= 252 {
		tmpMap := make(map[string]int)
		for k, v := range unitInfo.AdSourceCountry {
			if v == 0 {
				continue
			}
			tmpMap[k] = v
		}
		unitInfo.AdSourceLen = len(unitInfo.AdSourceCountry)
		unitInfo.AdSourceCountry = tmpMap
	}
	DbLoaderRegistry[TblUnit].dataCur.Set(concurrent_map.I64Key(unitInfo.UnitId), unitInfo)
}

func legalUnit(unitInfo *smodel.UnitInfo) bool {
	return unitInfo.Unit.Status == mvutil.ACTIVE
}

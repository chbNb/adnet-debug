package extractor

import (
	"errors"

	smodel "gitlab.mobvista.com/ADN/structs/model"

	"github.com/easierway/concurrent_map"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var placementUpdateInputError = errors.New("placementUpdateFunc failed, query or dbLoaderInfo is nil")

func NewPlacementExtractor(allFilter, incFilter, selectOptions bson.M) dbLoaderInfo {
	return dbLoaderInfo{
		collection:   "placement",
		querybmDbMap: allFilter,
		queryimDbMap: incFilter,
		querySelect:  selectOptions,
		updateFunc:   placementUpdateFunc,
		dataCur:      concurrent_map.CreateConcurrentMap(concurrentMapPartitionsNum),
	}
}

func placementUpdateFunc(query *mgo.Query, num int, dbLoaderInfo *dbLoaderInfo) (maxUpdate int64, err error) {
	if query == nil || dbLoaderInfo == nil {
		return 0, placementUpdateInputError
	}
	placementInfo := &smodel.PlacementInfo{}
	maxUpdate = 0
	item := query.Iter()
	var incData int
	for idx := 0; idx < num; idx++ {
		if item.Next(placementInfo) {
			if int64(placementInfo.Updated) > maxUpdate {
				maxUpdate = int64(placementInfo.Updated)
			}
			placementIncUpdateFunc(placementInfo, dbLoaderInfo)
			placementInfo = &smodel.PlacementInfo{}
			incData++
		} else {
			if item.Timeout() {
				logger.Warn("placementUpdateFunc item.Timeout()")
				break
			}
			// logger.Warnf("placementUpdateFunc item.error=%s", item.Err().Error())
		}
	}

	if err := item.Close(); err != nil {
		logger.Error(err.Error())
	}

	if incData != num {
		logger.Warnf("placementUpdateFunc Size not Equal incData=%d, QueryNum=%d", incData, num)
	}
	return maxUpdate, nil
}

func placementIncUpdateFunc(placementInfo *smodel.PlacementInfo, dbLoaderInfo *dbLoaderInfo) {
	if !legalPlacement(placementInfo) {
		DbLoaderRegistry[TblPlacement].dataCur.Del(concurrent_map.I64Key(placementInfo.PlacementId))
		return
	}
	DbLoaderRegistry[TblPlacement].dataCur.Set(concurrent_map.I64Key(placementInfo.PlacementId), placementInfo)
}

func legalPlacement(placementInfo *smodel.PlacementInfo) bool {
	if placementInfo.Status != mvutil.ACTIVE {
		return false
	}

	return true
}

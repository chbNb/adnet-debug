package extractor

import (
	"errors"

	smodel "gitlab.mobvista.com/ADN/structs/model"

	"github.com/easierway/concurrent_map"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var appPackageMtgIDInputError = errors.New("appPackageMtgIDUpdateFunc failed, query or dbLoaderInfo is nil")

func NewAppPackageMtgIDExtractor(allFilter, incFilter, selectOptions bson.M) dbLoaderInfo {
	return dbLoaderInfo{
		collection:   "app_package_mtg_id",
		querybmDbMap: allFilter,
		queryimDbMap: incFilter,
		querySelect:  selectOptions,
		updateFunc:   appPackageMtgIDUpdateFunc,
		dataCur:      concurrent_map.CreateConcurrentMap(concurrentMapPartitionsNum),
	}
}

func appPackageMtgIDUpdateFunc(query *mgo.Query, num int, dbLoaderInfo *dbLoaderInfo) (maxUpdate int64, err error) {
	if query == nil || dbLoaderInfo == nil {
		return 0, appPackageMtgIDInputError
	}
	appPackageMtgID := &smodel.AppPackageMtgID{}
	maxUpdate = 0
	item := query.Iter()
	var incData int
	for idx := 0; idx < num; idx++ {
		if item.Next(appPackageMtgID) {
			if int64(appPackageMtgID.Updated) > maxUpdate {
				maxUpdate = int64(appPackageMtgID.Updated)
			}
			appPackageMtgIDIncUpdateFunc(appPackageMtgID)
			incData++
			appPackageMtgID = &smodel.AppPackageMtgID{}
		} else {
			if item.Timeout() {
				logger.Warn("appPackageMtgIDUpdateFunc item.Timeout()")
				break
			}
			// logger.Warnf("appPackageMtgIDUpdateFunc item.error=%s", item.Err().Error())
		}
	}

	if err := item.Close(); err != nil {
		logger.Error(err.Error())
	}

	if incData != num {
		logger.Warnf("appPackageMtgIDUpdateFunc Size not Equal incData=%d, QueryNum=%d", incData, num)
	}
	return maxUpdate, nil

}

func appPackageMtgIDIncUpdateFunc(appPackageMtgID *smodel.AppPackageMtgID) {
	DbLoaderRegistry[TblAppPackageMTGId].dataCur.Set(concurrent_map.StrKey(appPackageMtgID.AppPackage), appPackageMtgID)
}

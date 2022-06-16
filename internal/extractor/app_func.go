package extractor

import (
	"errors"

	smodel "gitlab.mobvista.com/ADN/structs/model"

	"github.com/easierway/concurrent_map"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func NewAppExtractor(allFilter, incFilter, selectOptions bson.M) dbLoaderInfo {
	return dbLoaderInfo{
		collection:   "app",
		querybmDbMap: allFilter,
		queryimDbMap: incFilter,
		querySelect:  selectOptions,
		updateFunc:   appUpdateFunc,
		dataCur:      concurrent_map.CreateConcurrentMap(concurrentMapPartitionsNum),
	}
}

var appUpdateInputError = errors.New("appUpdateFunc failed, query or dbLoaderInfo is nil")

func appUpdateFunc(query *mgo.Query, num int, dbLoaderInfo *dbLoaderInfo) (maxUpdate int64, err error) {
	if query == nil || dbLoaderInfo == nil {
		return 0, appUpdateInputError
	}
	appInfo := &smodel.AppInfo{}
	maxUpdate = 0
	item := query.Iter()
	var incData, updData int
	for idx := 0; idx < num; idx++ {
		if item.Next(appInfo) {
			if int64(appInfo.Updated) > maxUpdate {
				maxUpdate = int64(appInfo.Updated)
			}
			incData++
			if meta.config.UseExpiredMap { // 新方式的时候才需要查
				if _, find := DbLoaderRegistry[TblAPP].dataCur.Get(concurrent_map.I64Key(appInfo.AppId)); !find && !dbLoaderInfo.updateByGetAll {
					continue // 找不到  且  是getInc  的就不用更新
				}
			}
			updData++
			appIncUpdateFunc(appInfo, dbLoaderInfo)
			appInfo = &smodel.AppInfo{}
		} else {
			if item.Timeout() {
				logger.Warn("appUpdateFunc item.Timeout()")
				break
			}
			// logger.Warnf("appUpdateFunc item.error=%s", item.Err().Error())
		}
	}

	if err := item.Close(); err != nil {
		logger.Error(err.Error())
	}

	if incData != num {
		logger.Warnf("appUpdateFunc Size not Equal incData=%d, QueryNum=%d, UpdateNum=%d", incData, num, updData)
	}
	logger.Infof("appUpdateFunc incData=%d, QueryNum=%d, UpdateNum=%d", incData, num, updData)
	return maxUpdate, nil
}

func appIncUpdateFunc(appInfo *smodel.AppInfo, dbLoaderInfo *dbLoaderInfo) {
	if !legalApp(appInfo) {
		if meta.config.UseExpiredMap { // 新的剔除需要保存为 nil
			DbLoaderRegistry[TblAPP].dataCur.Set(concurrent_map.I64Key(appInfo.AppId), nil)
		} else {
			DbLoaderRegistry[TblAPP].dataCur.Del(concurrent_map.I64Key(appInfo.AppId))
		}
		return
	}
	DbLoaderRegistry[TblAPP].dataCur.Set(concurrent_map.I64Key(appInfo.AppId), appInfo)
}

func legalApp(appInfo *smodel.AppInfo) bool {
	if appInfo.App.Status != mvutil.ACTIVE && appInfo.App.Status != mvutil.PENDING {
		return false
	}
	return true
}

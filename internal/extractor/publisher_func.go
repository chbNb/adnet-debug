package extractor

import (
	"errors"

	smodel "gitlab.mobvista.com/ADN/structs/model"

	"github.com/easierway/concurrent_map"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var publisherUpdateInputError = errors.New("publisherUpdateFunc failed, query or dbLoaderInfo is nil")

func NewPublisherExtractor(allFilter, incFilter, selectOptions bson.M) dbLoaderInfo {
	return dbLoaderInfo{
		collection:   "publisher",
		querybmDbMap: allFilter,
		queryimDbMap: incFilter,
		querySelect:  selectOptions,
		updateFunc:   publisherUpdateFunc,
		dataCur:      concurrent_map.CreateConcurrentMap(concurrentMapPartitionsNum),
	}
}

func publisherUpdateFunc(query *mgo.Query, num int, dbLoaderInfo *dbLoaderInfo) (maxUpdate int64, err error) {
	if query == nil || dbLoaderInfo == nil {
		return 0, publisherUpdateInputError
	}
	publisherInfo := &smodel.PublisherInfo{}
	maxUpdate = 0
	item := query.Iter()
	var incData, updData int
	for idx := 0; idx < num; idx++ {
		if item.Next(publisherInfo) {
			if int64(publisherInfo.Updated) > maxUpdate {
				maxUpdate = int64(publisherInfo.Updated)
			}
			incData++
			if meta.config.UseExpiredMap { // 新方式的时候才需要查
				if _, find := DbLoaderRegistry[TblPublisher].dataCur.Get(concurrent_map.I64Key(publisherInfo.PublisherId)); !find && !dbLoaderInfo.updateByGetAll {
					continue // 找不到  且  是getInc  的就不用更新
				}
			}
			updData++
			publisherIncUpdateFunc(publisherInfo, dbLoaderInfo)
			publisherInfo = &smodel.PublisherInfo{}
		} else {
			if item.Timeout() {
				logger.Warn("publisherUpdateFunc item.Timeout()")
				break
			}
			// logger.Warnf("publisherUpdateFunc item.error=%s", item.Err().Error())
		}
	}

	if err := item.Close(); err != nil {
		logger.Error(err.Error())
	}

	if incData != num {
		logger.Warnf("publisherUpdateFunc Size not Equal incData=%d, QueryNum=%d", incData, num)
	}
	return maxUpdate, nil
}

func publisherIncUpdateFunc(publisherInfo *smodel.PublisherInfo, dbLoaderInfo *dbLoaderInfo) {
	if !legalPublisher(publisherInfo) {
		if meta.config.UseExpiredMap { // 新的剔除需要保存为 nil
			DbLoaderRegistry[TblPublisher].dataCur.Set(concurrent_map.I64Key(publisherInfo.PublisherId), nil)
		} else {
			DbLoaderRegistry[TblPublisher].dataCur.Del(concurrent_map.I64Key(publisherInfo.PublisherId))
		}
		return
	}
	DbLoaderRegistry[TblPublisher].dataCur.Set(concurrent_map.I64Key(publisherInfo.PublisherId), publisherInfo)
}

func legalPublisher(publisherInfo *smodel.PublisherInfo) bool {
	if publisherInfo.Publisher.Status != mvutil.ACTIVE && publisherInfo.Publisher.Status != mvutil.PAUSED {
		return false
	}
	if len(publisherInfo.Publisher.Apikey) == 0 {
		return false
	}
	return true
}

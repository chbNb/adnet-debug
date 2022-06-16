package extractor

import (
	"errors"

	"github.com/easierway/concurrent_map"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	smodel "gitlab.mobvista.com/ADN/structs/model"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var AdvertiserInputError = errors.New("advertiserUpdateFunc failed, query or dbLoaderInfo is nil")

func NewAdvertiserExtractor(allFilter, incFilter, selectOptions bson.M) dbLoaderInfo {
	return dbLoaderInfo{
		collection:   "advertiser",
		querybmDbMap: allFilter,
		queryimDbMap: incFilter,
		querySelect:  selectOptions,
		updateFunc:   advertiserUpdateFunc,
		dataCur:      concurrent_map.CreateConcurrentMap(concurrentMapPartitionsNum),
	}
}

func advertiserUpdateFunc(query *mgo.Query, num int, dbLoaderInfo *dbLoaderInfo) (maxUpdate int64, err error) {
	if query == nil || dbLoaderInfo == nil {
		return 0, AdvertiserInputError
	}

	advInfo := new(smodel.AdvertiserInfo)
	maxUpdate = 0
	item := query.Iter()
	var incData int
	for idx := 0; idx < num; idx++ {
		if item.Next(advInfo) {
			if int64(advInfo.Updated) > maxUpdate {
				maxUpdate = int64(advInfo.Updated)
			}
			advIncUpdateFunc(advInfo, dbLoaderInfo)
			incData++
			advInfo = new(smodel.AdvertiserInfo)
		} else {
			if item.Timeout() {
				logger.Warn("advertiserUpdateFunc item.Timeout()")
				break
			}
			// logger.Warnf("advertiserUpdateFunc item.error=%s", item.Err().Error())
		}
	}

	if err := item.Close(); err != nil {
		logger.Error(err.Error())
	}

	if incData != num {
		logger.Warnf("advertiserUpdateFunc Size not Equal incData=%d, QueryNum=%d", incData, num)
	}
	return maxUpdate, nil
}

func advIncUpdateFunc(advInfo *smodel.AdvertiserInfo, dbLoaderInfo *dbLoaderInfo) {
	if advInfo.Advertiser.Status != mvutil.ACTIVE {
		DbLoaderRegistry[TblAdvertiser].dataCur.Del(concurrent_map.I64Key(advInfo.AdvertiserId))
		return
	}

	DbLoaderRegistry[TblAdvertiser].dataCur.Set(concurrent_map.I64Key(advInfo.AdvertiserId), advInfo)
	return
}

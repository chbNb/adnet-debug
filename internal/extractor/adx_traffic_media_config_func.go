package extractor

import (
	"errors"
	"strconv"

	smodel "gitlab.mobvista.com/ADN/structs/model"

	"github.com/easierway/concurrent_map"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var adxTrafficMedisConfigInputError = errors.New("adxTrafficMediaConfigFunc failed, query or dbLoaderInfo is nil")

func NewAdxTrafficMediaConfigExtractor(allFilter, incFilter, selectOptions bson.M) dbLoaderInfo {
	return dbLoaderInfo{
		collection:   "adx_traffic_media_config",
		querybmDbMap: allFilter,
		queryimDbMap: incFilter,
		querySelect:  selectOptions,
		updateFunc:   adxTrafficMediaConfigFunc,
		dataCur:      concurrent_map.CreateConcurrentMap(concurrentMapPartitionsNum),
	}
}

func adxTrafficMediaConfigFunc(query *mgo.Query, num int, dbLoaderInfo *dbLoaderInfo) (maxUpdate int64, err error) {
	if query == nil || dbLoaderInfo == nil {
		return 0, adxTrafficMedisConfigInputError
	}
	adxTrafficMediaConfig := &smodel.AdxTrafficMediaConfig{}
	maxUpdate = 0
	item := query.Iter()
	var incData int
	for idx := 0; idx < num; idx++ {
		if item.Next(adxTrafficMediaConfig) {
			if int64(adxTrafficMediaConfig.Updated) > maxUpdate {
				maxUpdate = int64(adxTrafficMediaConfig.Updated)
			}
			adxTrafficIncUpdateFunc(adxTrafficMediaConfig, dbLoaderInfo)
			adxTrafficMediaConfig = &smodel.AdxTrafficMediaConfig{}
			incData++
		} else {
			if item.Timeout() {
				logger.Warn("adxTrafficMediaConfigFunc item.Timeout()")
				break
			}
			// logger.Warnf("adxTrafficMediaConfigFunc item.error=%s", item.Err().Error())
		}
	}

	if err := item.Close(); err != nil {
		logger.Error(err.Error())
	}

	if incData != num {
		logger.Warnf("adxTrafficMediaConfigFunc Size not Equal incData=%d, QueryNum=%d", incData, num)
	}
	return maxUpdate, nil
}

func adxTrafficIncUpdateFunc(adxTrafficMediaConfig *smodel.AdxTrafficMediaConfig, dbLoaderInfo *dbLoaderInfo) {
	keyId := adxTrafficMediaConfig.UnitId
	if adxTrafficMediaConfig.Mode == 2 {
		keyId = adxTrafficMediaConfig.AdType
	}
	key := strconv.Itoa(adxTrafficMediaConfig.Mode) + "-" + strconv.FormatInt(keyId, 10) + "-" + adxTrafficMediaConfig.Area
	if !legalAdxTraffic(adxTrafficMediaConfig) {
		DbLoaderRegistry[TblAdxTrafficMediaConfig].dataCur.Del(concurrent_map.StrKey(key))
		return
	}
	DbLoaderRegistry[TblAdxTrafficMediaConfig].dataCur.Set(concurrent_map.StrKey(key), adxTrafficMediaConfig)
}

func legalAdxTraffic(adxTrafficMediaConfig *smodel.AdxTrafficMediaConfig) bool {
	return adxTrafficMediaConfig.Status == mvutil.ACTIVE
}

package extractor

import (
	"errors"
	"gitlab.mobvista.com/ADN/adnet/internal/mvconst"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	"strconv"

	"github.com/easierway/concurrent_map"
	smodel "gitlab.mobvista.com/ADN/structs/model"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var sspProfitDisUpdateInputError = errors.New("sspProfitDistributionUpdateFunc failed, query or dbLoaderInfo is nil")

func NewSspProfitDistributionRuleExtractor(allFilter, incFilter, selectOptions bson.M) dbLoaderInfo {
	return dbLoaderInfo{
		collection:   "ssp_profit_distribution_rule",
		querybmDbMap: allFilter,
		queryimDbMap: incFilter,
		querySelect:  selectOptions,
		updateFunc:   sspProfitDistributionUpdateFunc,
		dataCur:      concurrent_map.CreateConcurrentMap(concurrentMapPartitionsNum),
	}
}

func sspProfitDistributionUpdateFunc(query *mgo.Query, num int, dbLoaderInfo *dbLoaderInfo) (maxUpdate int64, err error) {
	if query == nil || dbLoaderInfo == nil {
		return 0, sspProfitDisUpdateInputError
	}
	sspProfitDistributionRule := &smodel.SspProfitDistributionRule{}
	maxUpdate = 0
	item := query.Iter()
	var incData int
	for idx := 0; idx < num; idx++ {
		if item.Next(sspProfitDistributionRule) {
			if int64(sspProfitDistributionRule.Updated) > maxUpdate {
				maxUpdate = int64(sspProfitDistributionRule.Updated)
			}
			sspProfitDistributionIncUpdateFunc(sspProfitDistributionRule, dbLoaderInfo)
			incData++
			sspProfitDistributionRule = &smodel.SspProfitDistributionRule{}
		} else {
			if item.Timeout() {
				logger.Warn("sspProfitDistributionUpdateFunc item.Timeout()")
				break
			}
			// logger.Warnf("sspProfitDistributionUpdateFunc item.error=%s", item.Err().Error())
		}
	}

	if err := item.Close(); err != nil {
		logger.Error(err.Error())
	}

	if incData != num {
		logger.Warnf("sspProfitDistributionUpdateFunc Size not Equal incData=%d, QueryNum=%d", incData, num)
	}
	return maxUpdate, nil

}

func sspProfitDistributionIncUpdateFunc(sspProfitDistributionRule *smodel.SspProfitDistributionRule, dbLoaderInfo *dbLoaderInfo) {
	key := strconv.FormatInt(sspProfitDistributionRule.UnitId, 10) + ":" + sspProfitDistributionRule.Area
	if !legalSspProfitDistributionRule(sspProfitDistributionRule) {
		DbLoaderRegistry[TblSSPProfitDistributionRule].dataCur.Del(concurrent_map.StrKey(key))
		return
	}
	DbLoaderRegistry[TblSSPProfitDistributionRule].dataCur.Set(concurrent_map.StrKey(key), sspProfitDistributionRule)
}

func legalSspProfitDistributionRule(sspProfitDistributionRule *smodel.SspProfitDistributionRule) bool {
	if (sspProfitDistributionRule.Type != mvconst.SspProfitDistributionRuleFixedEcpm &&
		sspProfitDistributionRule.Type != mvconst.SspProfitDistributionRuleOnlineApiEcpm) || sspProfitDistributionRule.Status != mvutil.ACTIVE {
		return false
	}
	return true
}

package extractor

import (
	"errors"

	"github.com/easierway/concurrent_map"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	smodel "gitlab.mobvista.com/ADN/structs/model"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type ErrorCampaign struct {
	CampaignId int64 `bson:"campaignId,omitempty" json:"campaignId,omitempty"`
}

var campaignUpdateInputError = errors.New("campaignUpdateFunc failed, query or dbLoaderInfo is nil")

func NewCampaignExtractor(allFilter, incFilter, selectOptions bson.M) dbLoaderInfo {
	return dbLoaderInfo{
		collection:   "campaign",
		querybmDbMap: allFilter,
		queryimDbMap: incFilter,
		querySelect:  selectOptions,
		updateFunc:   campaignUpdateFunc,
		dataCur:      concurrent_map.CreateConcurrentMap(concurrentMapPartitionsNum),
	}
}

func campaignUpdateFunc(query *mgo.Query, num int, dbLoaderInfo *dbLoaderInfo) (maxUpdate int64, err error) {
	if query == nil || dbLoaderInfo == nil {
		return 0, campaignUpdateInputError
	}
	camInfo := &smodel.CampaignInfo{}
	maxUpdate = 0
	item := query.Iter()
	var incData, updData int
	for idx := 0; idx < num; idx++ {
		rawBson := &bson.Raw{}
		if item.Next(rawBson) {
			if err := rawBson.Unmarshal(&camInfo); err != nil {
				errCampaign := ErrorCampaign{}
				err := rawBson.Unmarshal(&errCampaign)
				if err == nil {
					logger.Warnf("campaignUpdateFunc Unmarshal error CampaignID=%d", errCampaign.CampaignId)
				} else {
					logger.Warnf("campaignUpdateFunc Unmarshal error: %s", err.Error())
				}
				continue
			}
			if int64(camInfo.Updated) > maxUpdate {
				maxUpdate = int64(camInfo.Updated)
			}
			incData++
			if meta.config.UseExpiredMap { // 新方式的时候才需要查
				if _, find := DbLoaderRegistry[TblCampaign].dataCur.Get(concurrent_map.I64Key(camInfo.CampaignId)); !find && !dbLoaderInfo.updateByGetAll {
					continue // 找不到  且  是getInc  的就不用更新
				}
			}
			updData++
			camIncUpdateFunc(camInfo, dbLoaderInfo)
			camInfo = &smodel.CampaignInfo{}
		} else {
			if item.Timeout() {
				logger.Warn("campaignUpdateFunc item.Timeout()")
				break
			}
			// logger.Warnf("campaignUpdateFunc item.error=%s", item.Err().Error())
		}
	}

	if err := item.Close(); err != nil {
		logger.Error(err.Error())
	}

	if incData != num {
		logger.Warnf("campaignUpdateFunc Size not Equal incData=%d, QueryNum=%d, UpdateNum=%d", incData, num, updData)
	}
	logger.Infof("campaignUpdateFunc incData=%d, QueryNum=%d, UpdateNum=%d", incData, num, updData)
	return maxUpdate, nil
}

func camIncUpdateFunc(camInfo *smodel.CampaignInfo, dbLoaderInfo *dbLoaderInfo) {
	if !legalCam(camInfo) {
		if meta.config.UseExpiredMap { // 新的剔除需要保存为nil
			DbLoaderRegistry[TblCampaign].dataCur.Set(concurrent_map.I64Key(camInfo.CampaignId), nil)
		} else {
			DbLoaderRegistry[TblCampaign].dataCur.Del(concurrent_map.I64Key(camInfo.CampaignId))
		}
		return
	}
	DbLoaderRegistry[TblCampaign].dataCur.Set(concurrent_map.I64Key(camInfo.CampaignId), camInfo)
}

func legalCam(camInfo *smodel.CampaignInfo) bool {
	return camInfo.Status == mvutil.ACTIVE
}

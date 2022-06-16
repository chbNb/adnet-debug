package extractor

import (
	jsoniter "github.com/json-iterator/go"
	"gitlab.mobvista.com/ADN/adnet/internal/mvconst"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/easierway/concurrent_map"
	"gitlab.mobvista.com/ADN/adnet/internal/expired_map"
	"gitlab.mobvista.com/ADN/adnet/internal/hot_data"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	"gitlab.mobvista.com/ADN/adnet/internal/watcher"
	"gitlab.mobvista.com/ADN/lego/map_key"
	smodel "gitlab.mobvista.com/ADN/structs/model"
	"gitlab.mobvista.com/ADN/treasure_box_sdk/tb_tools"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	TimeoutMongo = 10
)

func getAppInfoFromConcurrentExpiredMap(appId int64) (*smodel.AppInfo, bool) {
	defer mvutil.RecordUseTime("avg_get_app_cem", time.Now().UnixNano(), 1e3)
	appInfoObj, ifFind := getDataInfoFromConcurrentExpiredMap(
		appId,
		TblAPP,
		bson.M{"appId": appId},
		&smodel.AppInfo{},
		func(dataInter interface{}) (interface{}, bool) {
			info, ifFind := dataInter.(*smodel.AppInfo)
			if !ifFind || info == nil || !legalApp(info) {
				dataInter = nil
				ifFind = false
			}
			return dataInter, ifFind
		})
	if !ifFind {
		return nil, false
	}
	appInfo, ifFind := appInfoObj.(*smodel.AppInfo)
	return appInfo, ifFind
}
func getAppInfoFromConcurrentMap(appId int64) (*smodel.AppInfo, bool) {
	defer mvutil.RecordUseTime("avg_get_app_cm", time.Now().UnixNano(), 1e3)
	appInfoObj, ifFind := DbLoaderRegistry[TblAPP].dataCur.Get(concurrent_map.I64Key(appId))
	if !ifFind {
		return nil, false
	}
	appInfo, ifFind := appInfoObj.(*smodel.AppInfo)
	return appInfo, ifFind
}
func getAppInfoFromTreasureBox(appId int64) (*smodel.AppInfo, bool) {
	defer mvutil.RecordUseTime("avg_get_app_tb", time.Now().UnixNano(), 1e3)
	reObj, ifFind := tb_tools.GetData(&smodel.AppInfo{}, map_key.I64Key(appId))
	if ifFind {
		if re, ok := reObj.(*smodel.AppInfo); ok {
			return re, ok
		}
	}
	return nil, false
}
func GetAppInfo(appId int64) (appInfo *smodel.AppInfo, ifFind bool) {
	if tb_tools.IsAvailable(&smodel.AppInfo{}) { //判断 Treasure Box SDK 是否支持
		return getAppInfoFromTreasureBox(appId)
	}
	if meta.config.UseExpiredMap {
		appInfo, ifFind = getAppInfoFromConcurrentExpiredMap(appId)
	} else {
		appInfo, ifFind = getAppInfoFromConcurrentMap(appId)
	}

	if ifFind {
		hot_data.AddToActiveDataCollecter(hot_data.APP, appInfo.AppId)
	}

	return appInfo, ifFind
}

func getCampaignInfoFromConcurrentExpiredMap(camId int64) (camInfo *smodel.CampaignInfo, ifFind bool) {
	defer mvutil.RecordUseTime("avg_get_campaign_cem", time.Now().UnixNano(), 1e3)
	camInfoObj, ifFind := getDataInfoFromConcurrentExpiredMap(
		camId,
		TblCampaign,
		bson.M{"campaignId": camId, "status": 1},
		&smodel.CampaignInfo{},
		func(dataInter interface{}) (interface{}, bool) {
			camInfo, ifFind := dataInter.(*smodel.CampaignInfo)
			if !ifFind || camInfo == nil || !legalCam(camInfo) {
				dataInter = nil
				ifFind = false
			}
			return dataInter, ifFind
		})
	if !ifFind {
		return nil, false
	}
	camInfo, ifFind = camInfoObj.(*smodel.CampaignInfo)
	return camInfo, ifFind
}

func getCampaignInfoFromConcurrentMap(camId int64) (*smodel.CampaignInfo, bool) {
	defer mvutil.RecordUseTime("avg_get_campaign_cm", time.Now().UnixNano(), 1e3)
	camInfoObj, ifFind := DbLoaderRegistry[TblCampaign].dataCur.Get(concurrent_map.I64Key(camId))
	if !ifFind {
		return nil, false
	}
	camInfo, ifFind := camInfoObj.(*smodel.CampaignInfo)
	return camInfo, ifFind
}
func getCampaignInfoFromTreasureBox(camId int64) (*smodel.CampaignInfo, bool) {
	defer mvutil.RecordUseTime("avg_get_campaign_tb", time.Now().UnixNano(), 1e3)
	reObj, ifFind := tb_tools.GetData(&smodel.CampaignInfo{}, map_key.I64Key(camId))
	if ifFind {
		if re, ok := reObj.(*smodel.CampaignInfo); ok {
			return re, ok
		}
	}
	return nil, false
}

func CheckGetCampaignInfo(camId int64) smodel.CheckResult {
	oldData, _ := getCampaignInfoFromConcurrentExpiredMap(camId)
	newData, _ := getCampaignInfoFromTreasureBox(camId)
	result := "ok"
	if !reflect.DeepEqual(oldData, newData) {
		result = "error"
	}
	return smodel.CheckResult{
		Result:  result,
		OldData: oldData,
		NewData: newData,
	}
}
func GetCampaignInfo(camId int64) (camInfo *smodel.CampaignInfo, ifFind bool) {
	if tb_tools.IsAvailable(&smodel.CampaignInfo{}) { //判断 Treasure Box SDK 是否支持
		return getCampaignInfoFromTreasureBox(camId)
	}
	if meta.config.UseExpiredMap {
		camInfo, ifFind = getCampaignInfoFromConcurrentExpiredMap(camId)
	} else {
		camInfo, ifFind = getCampaignInfoFromConcurrentMap(camId)
	}

	if ifFind {
		hot_data.AddToActiveDataCollecter(hot_data.CAMPAIGN, camInfo.CampaignId)
	}

	return camInfo, ifFind
}

func GetAdvertiserInfo(advId int64) (advInfo *smodel.AdvertiserInfo, ifFind bool) {
	if tb_tools.IsAvailable(&smodel.AdvertiserInfo{}) { //判断 Treasure Box SDK 是否支持
		return getAdvertiserInfoFromTreasureBox(advId)
	}
	return getAdvertiserInfoFromConcurrentMap(advId)
}
func getAdvertiserInfoFromTreasureBox(advId int64) (*smodel.AdvertiserInfo, bool) {
	defer mvutil.RecordUseTime("avg_get_advertiser_tb", time.Now().UnixNano(), 1e3)
	reObj, ifFind := tb_tools.GetData(&smodel.AdvertiserInfo{}, map_key.I64Key(advId))
	if ifFind {
		if re, ok := reObj.(*smodel.AdvertiserInfo); ok {
			return re, ok
		}
	}
	return nil, false
}
func getAdvertiserInfoFromConcurrentMap(advId int64) (advInfo *smodel.AdvertiserInfo, ifFind bool) {
	defer mvutil.RecordUseTime("avg_get_advertiser_cm", time.Now().UnixNano(), 1e3)
	advInfoObj, ifFind := DbLoaderRegistry[TblAdvertiser].dataCur.Get(concurrent_map.I64Key(advId))
	if !ifFind {
		return
	}
	advInfo, ifFind = advInfoObj.(*smodel.AdvertiserInfo)
	return advInfo, ifFind
}

func getUnitInfoFromConcurrentExpiredMap(unitId int64) (*smodel.UnitInfo, bool) {
	defer mvutil.RecordUseTime("avg_get_unit_cem", time.Now().UnixNano(), 1e3)
	unitInfoObj, ifFind := getDataInfoFromConcurrentExpiredMap(
		unitId,
		TblUnit,
		bson.M{"unitId": unitId},
		&smodel.UnitInfo{},
		func(dataInter interface{}) (interface{}, bool) {
			info, ifFind := dataInter.(*smodel.UnitInfo)
			if !ifFind || info == nil || !legalUnit(info) {
				dataInter = nil
				ifFind = false
			}
			return dataInter, ifFind
		})
	if !ifFind {
		return nil, false
	}
	unitInfo, ifFind := unitInfoObj.(*smodel.UnitInfo)
	return unitInfo, ifFind
}
func getUnitInfoFromConcurrentMap(unitId int64) (*smodel.UnitInfo, bool) {
	defer mvutil.RecordUseTime("avg_get_unit_cm", time.Now().UnixNano(), 1e3)
	unitInfoObj, ifFind := DbLoaderRegistry[TblUnit].dataCur.Get(concurrent_map.I64Key(unitId))
	if !ifFind {
		return nil, false
	}
	unitInfo, ifFind := unitInfoObj.(*smodel.UnitInfo)
	return unitInfo, ifFind
}
func getUnitInfoFromTreasureBox(unitId int64) (unitInfo *smodel.UnitInfo, ifFind bool) {
	defer mvutil.RecordUseTime("avg_get_unit_tb", time.Now().UnixNano(), 1e3)
	reObj, ifFind := tb_tools.GetData(&smodel.UnitInfo{}, map_key.I64Key(unitId))
	if ifFind {
		if re, ok := reObj.(*smodel.UnitInfo); ok {
			return re, ok
		}
	}
	return nil, false
}
func GetUnitInfo(unitId int64) (unitInfo *smodel.UnitInfo, ifFind bool) {
	if tb_tools.IsAvailable(&smodel.UnitInfo{}) { //判断 Treasure Box SDK 是否支持
		return getUnitInfoFromTreasureBox(unitId)
	}
	if meta.config.UseExpiredMap {
		unitInfo, ifFind = getUnitInfoFromConcurrentExpiredMap(unitId)
	} else {
		unitInfo, ifFind = getUnitInfoFromConcurrentMap(unitId)
	}
	if ifFind {
		hot_data.AddToActiveDataCollecter(hot_data.UNIT, unitInfo.UnitId)
	}
	return unitInfo, ifFind
}

func GetUnitFixedEcpm(unitId int64, cc string) float64 {
	var unitFixedEcpm float64
	unitInfo, found := GetUnitInfo(unitId)
	if found && len(unitInfo.FakeRuleV2) > 0 {
		var fakeRule smodel.FakeRule
		var hasConf bool
		fakeRule, hasConf = unitInfo.FakeRuleV2[cc]
		if !hasConf || fakeRule.Status != 1 || fakeRule.Type != mvutil.FixedEcpm {
			fakeRule, hasConf = unitInfo.FakeRuleV2["ALL"]
		}
		if hasConf && fakeRule.Status == 1 && fakeRule.Type == mvutil.FixedEcpm {
			unitFixedEcpm = fakeRule.Ecppv
		}
	}
	return unitFixedEcpm
}

func getPublisherInfoFromConcurrentExpiredMap(publisherId int64) (*smodel.PublisherInfo, bool) {
	defer mvutil.RecordUseTime("avg_get_publisher_cem", time.Now().UnixNano(), 1e3)
	publisherInfoObj, ifFind := getDataInfoFromConcurrentExpiredMap(
		publisherId,
		TblPublisher,
		bson.M{"publisherId": publisherId},
		&smodel.PublisherInfo{},
		func(dataInter interface{}) (interface{}, bool) {
			info, ifFind := dataInter.(*smodel.PublisherInfo)
			if !ifFind || info == nil || !legalPublisher(info) {
				dataInter = nil
				ifFind = false
			}
			return dataInter, ifFind
		})
	if !ifFind {
		return nil, false
	}
	publisherInfo, ifFind := publisherInfoObj.(*smodel.PublisherInfo)
	return publisherInfo, ifFind
}
func getPublisherInfoFromConcurrentMap(publisherId int64) (*smodel.PublisherInfo, bool) {
	defer mvutil.RecordUseTime("avg_get_publisher_cm", time.Now().UnixNano(), 1e3)
	pubInfoObj, ifFind := DbLoaderRegistry[TblPublisher].dataCur.Get(concurrent_map.I64Key(publisherId))
	if !ifFind {
		return nil, false
	}
	publisherInfo, ifFind := pubInfoObj.(*smodel.PublisherInfo)
	return publisherInfo, ifFind
}
func getPublisherInfoFromTreasureBox(publisherId int64) (*smodel.PublisherInfo, bool) {
	defer mvutil.RecordUseTime("avg_get_publisher_tb", time.Now().UnixNano(), 1e3)
	reObj, ifFind := tb_tools.GetData(&smodel.PublisherInfo{}, map_key.I64Key(publisherId))
	if ifFind {
		if re, ok := reObj.(*smodel.PublisherInfo); ok {
			return re, ok
		}
	}
	return nil, false
}
func GetPublisherInfo(publisherId int64) (publisherInfo *smodel.PublisherInfo, ifFind bool) {
	if tb_tools.IsAvailable(&smodel.PublisherInfo{}) { //判断 Treasure Box SDK 是否支持
		return getPublisherInfoFromTreasureBox(publisherId)
	}
	if meta.config.UseExpiredMap {
		publisherInfo, ifFind = getPublisherInfoFromConcurrentExpiredMap(publisherId)
	} else {
		publisherInfo, ifFind = getPublisherInfoFromConcurrentMap(publisherId)
	}
	if ifFind {
		hot_data.AddToActiveDataCollecter(hot_data.PUBLISHER, publisherInfo.PublisherId)
	}
	return publisherInfo, ifFind
}

func GetPlacementInfo(placementId int64) (*smodel.PlacementInfo, bool) {
	if tb_tools.IsAvailable(&smodel.PlacementInfo{}) { //判断 Treasure Box SDK 是否支持
		return getPlacementInfoTreasureBox(placementId)
	}
	return getPlacementInfoFromConcurrentMap(placementId)
}
func getPlacementInfoTreasureBox(placementId int64) (*smodel.PlacementInfo, bool) {
	defer mvutil.RecordUseTime("avg_get_placement_tb", time.Now().UnixNano(), 1e3)
	reObj, ifFind := tb_tools.GetData(&smodel.PlacementInfo{}, map_key.I64Key(placementId))
	if ifFind {
		if re, ok := reObj.(*smodel.PlacementInfo); ok {
			return re, ok
		}
	}
	return nil, false
}
func getPlacementInfoFromConcurrentMap(placementId int64) (*smodel.PlacementInfo, bool) {
	placementInfoObj, ifFind := DbLoaderRegistry[TblPlacement].dataCur.Get(concurrent_map.I64Key(placementId))
	if !ifFind {
		return nil, false
	}
	placementInfo, ifFind := placementInfoObj.(*smodel.PlacementInfo)
	return placementInfo, ifFind
}

func GetMVConfigValue(key string) (value interface{}, ifFind bool) {
	if !tb_tools.IsAvailable(&smodel.MVConfig{}) { //判断 Treasure Box SDK 是否支持
		return nil, false
	}
	mvConfig, ifFind := getMVConfigValueFromTreasureBox(key)
	if !ifFind {
		return nil, false
	}
	// jsonValue, err := bson.MarshalExtJSON(mvconfig.Value, false, false)
	// logger.Infof("json.Marshal %s", jsonValue)
	return mvConfig.Value, ifFind
}

func getMVConfigValueFromTreasureBox(key string) (*smodel.MVConfig, bool) {
	defer mvutil.RecordUseTime("avg_get_mvconfig_tb", time.Now().UnixNano(), 1e3)
	reObj, ifFind := tb_tools.GetData(&smodel.MVConfig{}, map_key.StrKey(key))
	if ifFind {
		if re, ok := reObj.(*smodel.MVConfig); ok {
			return re, ok
		}
	}
	return nil, false
}
func getMVConfigValueFromConcurrentMap(key string) (*smodel.MVConfig, bool) {
	defer mvutil.RecordUseTime("avg_get_mvconfig_cm", time.Now().UnixNano(), 1e3)
	configObj, ifFind := DbLoaderRegistry[TblConfig].dataCur.Get(concurrent_map.StrKey(key))
	if !ifFind {
		return nil, ifFind
	}
	mvconfig, ifFind := configObj.(*smodel.MVConfig)
	return mvconfig, ifFind
}

func GetConfigcenter(key string) (configcenter *smodel.ConfigCenter, ifFind bool) {
	if tb_tools.IsAvailable(&smodel.ConfigCenter{}) { //判断 Treasure Box SDK 是否支持
		return getConfigcenterFromTreasureBox(key)
	}
	return getConfigcenterFromConcurrentMap(key)
}
func getConfigcenterFromTreasureBox(key string) (*smodel.ConfigCenter, bool) {
	defer mvutil.RecordUseTime("avg_get_configcenter_tb", time.Now().UnixNano(), 1e3)
	reObj, ifFind := tb_tools.GetData(&smodel.ConfigCenter{}, map_key.StrKey(key))
	if ifFind {
		if re, ok := reObj.(*smodel.ConfigCenter); ok {
			return re, ok
		}
	}
	return nil, false
}
func getConfigcenterFromConcurrentMap(key string) (configcenter *smodel.ConfigCenter, ifFind bool) {
	defer mvutil.RecordUseTime("avg_get_configcenter_cm", time.Now().UnixNano(), 1e3)
	configCenterObj, ifFind := DbLoaderRegistry[TblConfigCenter].dataCur.Get(concurrent_map.StrKey(key))
	if !ifFind {
		return
	}
	configcenter, ifFind = configCenterObj.(*smodel.ConfigCenter)
	return configcenter, ifFind
}

func GetAdxTrafficMediaConfigByKey(key string) (adxTrafficMediaConfig *smodel.AdxTrafficMediaConfig, ifFind bool) {
	if tb_tools.IsAvailable(&smodel.AdxTrafficMediaConfig{}) { //判断 Treasure Box SDK 是否支持
		return getAdxTrafficMediaConfigFromTreasureBox(key)
	}
	return getAdxTrafficMediaConfigFromConcurrentMap(key)

}
func getAdxTrafficMediaConfigFromTreasureBox(key string) (*smodel.AdxTrafficMediaConfig, bool) {
	defer mvutil.RecordUseTime("avg_get_AdxTrafficMediaConfig_tb", time.Now().UnixNano(), 1e3)
	reObj, ifFind := tb_tools.GetData(&smodel.AdxTrafficMediaConfig{}, map_key.StrKey(key))
	if ifFind {
		if re, ok := reObj.(*smodel.AdxTrafficMediaConfig); ok {
			return re, ok
		}
	}
	return nil, false
}
func getAdxTrafficMediaConfigFromConcurrentMap(key string) (adxTrafficMediaConfig *smodel.AdxTrafficMediaConfig, ifFind bool) {
	defer mvutil.RecordUseTime("avg_get_AdxTrafficMediaConfig_cm", time.Now().UnixNano(), 1e3)
	cfgObj, ifFind := DbLoaderRegistry[TblAdxTrafficMediaConfig].dataCur.Get(concurrent_map.StrKey(key))
	if !ifFind {
		return adxTrafficMediaConfig, ifFind
	}
	adxTrafficMediaConfig, ifFind = cfgObj.(*smodel.AdxTrafficMediaConfig)
	return adxTrafficMediaConfig, ifFind

}
func GetAdxTrafficMediaConfig(unitId int64, adType int32, country string) (adxTrafficMediaConfig *smodel.AdxTrafficMediaConfig, ifFind bool) {
	// mode-unitId/adtype-cc; mode:1表示 unit 维度，2表示 adtype 维度
	// 查找优先级 unit+cc > unit+ALL > adtype+cc > adtype+ALL
	unitIdStr := strconv.FormatInt(unitId, 10)
	key := "1-" + unitIdStr + "-" + country
	trafficCfg, ok := GetAdxTrafficMediaConfigByKey(key)
	if ok {
		return trafficCfg, ok
	}
	key = "1-" + unitIdStr + "-ALL"
	trafficCfg, ok = GetAdxTrafficMediaConfigByKey(key)
	if ok {
		return trafficCfg, ok
	}
	adTypeStr := strconv.Itoa(mvutil.GetDspAdType(adType))
	key = "2-" + adTypeStr + "-" + country
	trafficCfg, ok = GetAdxTrafficMediaConfigByKey(key)
	if ok {
		return trafficCfg, ok
	}
	key = "2-" + adTypeStr + "-ALL"
	trafficCfg, ok = GetAdxTrafficMediaConfigByKey(key)
	if ok {
		return trafficCfg, ok
	}
	return trafficCfg, ok
}

func GetFillRateKey(unitId int64, appId int64, platform int, country string) string {
	uStr := strconv.FormatInt(unitId, 10)
	aStr := strconv.FormatInt(appId, 10)
	pStr := strconv.Itoa(platform)
	return pStr + "_" + aStr + "_" + uStr + "_" + country
}

func GetEcpmFloor(key string) (ecpmFloor float64, ifFind bool) {
	configObj, ifFind := getConfigAlgorithmFillRate(key)
	if !ifFind {
		return
	}
	ecpmFloor = configObj.EcpmFloor
	return ecpmFloor, ifFind
}

func GetFillRate(key string) (rate int, ifFind bool) {
	rateData, ifFind := getConfigAlgorithmFillRate(key)
	if !ifFind {
		return
	}
	rate = rateData.Rate
	return rate, ifFind
}
func getConfigAlgorithmFillRate(key string) (*smodel.ConfigAlgorithmFillRate, bool) {
	if tb_tools.IsAvailable(&smodel.ConfigAlgorithmFillRate{}) { //判断 Treasure Box SDK 是否支持
		return getConfigAlgorithmFillRateFromTreasureBox(key)
	}
	return getConfigAlgorithmFillRateFromConcurrentMap(key)
}
func getConfigAlgorithmFillRateFromTreasureBox(key string) (*smodel.ConfigAlgorithmFillRate, bool) {
	defer mvutil.RecordUseTime("avg_get_AlgorithmFillRate_tb", time.Now().UnixNano(), 1e3)
	reObj, ifFind := tb_tools.GetData(&smodel.ConfigAlgorithmFillRate{}, map_key.StrKey(key))
	if ifFind {
		if re, ok := reObj.(*smodel.ConfigAlgorithmFillRate); ok {
			return re, ok
		}
	}
	return nil, false
}
func getConfigAlgorithmFillRateFromConcurrentMap(key string) (*smodel.ConfigAlgorithmFillRate, bool) {
	defer mvutil.RecordUseTime("avg_get_AlgorithmFillRate_cm", time.Now().UnixNano(), 1e3)
	rateObj, ifFind := DbLoaderRegistry[TblConfigAlgorithmFillrate].dataCur.Get(concurrent_map.StrKey(key))
	if !ifFind {
		return nil, false
	}
	rateData, ifFind := rateObj.(*smodel.ConfigAlgorithmFillRate)
	return rateData, ifFind
}

///////to-do adnet不应该读取dspconfig的信息
func GetDspConfig(dspId int64) (adxDspConfig *smodel.AdxDspConfig, ifFind bool) {
	if tb_tools.IsAvailable(&smodel.AdxDspConfig{}) { //判断 Treasure Box SDK 是否支持
		return getDspConfigFromTreasureBox(dspId)
	}
	return getDspConfigFromConcurrentMap(dspId)
}

func getDspConfigFromTreasureBox(appId int64) (adxDspConfig *smodel.AdxDspConfig, ifFind bool) {
	defer mvutil.RecordUseTime("avg_get_dsp_config_tb", time.Now().UnixNano(), 1e3)
	reObj, ifFind := tb_tools.GetData(&smodel.AdxDspConfig{}, map_key.I64Key(appId))
	if ifFind {
		if re, ok := reObj.(*smodel.AdxDspConfig); ok {
			return re, ok
		}
	}
	return nil, false
}
func getDspConfigFromConcurrentMap(dspId int64) (adxDspConfig *smodel.AdxDspConfig, ifFind bool) {
	defer mvutil.RecordUseTime("avg_get_dsp_config_cm", time.Now().UnixNano(), 1e3)
	adxDspConfigObj, ifFind := DbLoaderRegistry[TblAdxDspConfig].dataCur.Get(concurrent_map.I64Key(dspId))
	if !ifFind {
		return adxDspConfig, ifFind
	}
	adxDspConfig, ifFind = adxDspConfigObj.(*smodel.AdxDspConfig)
	return adxDspConfig, ifFind
}

func GetFillRateControllConfig(key string) (cfg *smodel.ConfigAlgorithmFillRate, ifFind bool) {
	if tb_tools.IsAvailable(&smodel.ConfigAlgorithmFillRate{}) { //判断 Treasure Box SDK 是否支持
		return getFillRateControllConfigFromTreasureBox(key)
	}
	return getFillRateControllConfigFromConcurrentMap(key)
}
func getFillRateControllConfigFromTreasureBox(key string) (*smodel.ConfigAlgorithmFillRate, bool) {
	defer mvutil.RecordUseTime("avg_get_fillRateControllConfig_tb", time.Now().UnixNano(), 1e3)
	reObj, ifFind := tb_tools.GetData(&smodel.ConfigAlgorithmFillRate{}, map_key.StrKey(key))
	if ifFind {
		if re, ok := reObj.(*smodel.ConfigAlgorithmFillRate); ok {
			return re, ok
		}
	}
	return nil, false
}
func getFillRateControllConfigFromConcurrentMap(key string) (cfg *smodel.ConfigAlgorithmFillRate, ifFind bool) {
	defer mvutil.RecordUseTime("avg_get_fillRateControllConfig_cm", time.Now().UnixNano(), 1e3)
	cfgObj, ifFind := DbLoaderRegistry[TblConfigAlgorithmFillrate].dataCur.Get(concurrent_map.StrKey(key))
	if !ifFind {
		return
	}
	cfg, ifFind = cfgObj.(*smodel.ConfigAlgorithmFillRate)
	return cfg, ifFind
}

func GetFreqControlPriceFactor(key string) (freqControlFactor *smodel.FreqControlFactor, ifFind bool) {
	if tb_tools.IsAvailable(&smodel.FreqControlFactor{}) { //判断 Treasure Box SDK 是否支持
		return getFreqControlPriceFactorFromTreasureBox(key)
	}
	return getFreqControlPriceFactorFromConcurrentMap(key)
}
func getFreqControlPriceFactorFromTreasureBox(key string) (*smodel.FreqControlFactor, bool) {
	defer mvutil.RecordUseTime("avg_get_app_tb", time.Now().UnixNano(), 1e3)
	reObj, ifFind := tb_tools.GetData(&smodel.FreqControlFactor{}, map_key.StrKey(key))
	if ifFind {
		if re, ok := reObj.(*smodel.FreqControlFactor); ok {
			return re, ok
		}
	}
	return nil, false
}
func getFreqControlPriceFactorFromConcurrentMap(key string) (freqControlFactor *smodel.FreqControlFactor, ifFind bool) {
	cfgObj, ifFind := DbLoaderRegistry[TblFreqControlFactor].dataCur.Get(concurrent_map.StrKey(key))
	if !ifFind {
		return freqControlFactor, ifFind
	}
	freqControlFactor, ifFind = cfgObj.(*smodel.FreqControlFactor)
	return freqControlFactor, ifFind

}

func getDataInfoFromConcurrentExpiredMap(dataId int64, tableIndex int, bsonObj bson.M, dataObj interface{},
	assertionFnc func(interface{}) (interface{}, bool)) (dataInfo interface{}, ifFind bool) {
	canTry := true

TryAgain:
	dbLoaderRegistry := DbLoaderRegistry[tableIndex]
	tableName := dbLoaderRegistry.collection
	dataInfoObj, ifFind := dbLoaderRegistry.dataCur.Get(concurrent_map.I64Key(dataId))
	watcher.AddWatchValue("get_new_"+tableName+"_data_count", float64(1))

	needToSaveToCEM := 0
	if !ifFind { // 从Mongo获取
		lock := DbLoaderRegistryLock[tableIndex].TryLock(dataId)
		defer DbLoaderRegistryLock[tableIndex].UnLock(dataId)

		if !lock { // 拿不到锁的直接认为单子不存在
			mvutil.Logger.Runtime.Infof("Get %s Info[id:%d] is locked", tableName, dataId)
			if canTry && meta.config.EMRetryAgrainSleepMicrosecond > 0 {
				canTry = false
				time.Sleep(time.Millisecond * time.Duration(meta.config.EMRetryAgrainSleepMicrosecond))
				watcher.AddWatchValue(""+tableName+"_try_again_count", float64(1))
				goto TryAgain
			}
			watcher.AddWatchValue(""+tableName+"_unlock_count", float64(1))
			return
		}
		// 从mongo获取
		watcher.AddWatchValue(""+tableName+"_from_mongo_count", float64(1))
		if dataInfoMongoObj, err := getDataInfoFromMongo(tableName, bsonObj, dbLoaderRegistry.querySelect, dataObj); err == nil {
			dataInfoObj = dataInfoMongoObj
			needToSaveToCEM = 1
		} else {
			// if err == mongo.ErrNoDocuments {
			if err == mgo.ErrNotFound {
				needToSaveToCEM = 1
			}
			watcher.AddWatchValue(""+tableName+"_not_in_mongo_count", float64(1))
		}
	} else {
		needToSaveToCEM = 2
	}

	dataInfo, ifFind = assertionFnc(dataInfoObj)

	if needToSaveToCEM == 2 { // 从内存中获取到的数据，如果有效时间少于
		if cur, ok := DbLoaderRegistry[tableIndex].dataCur.(*expired_map.ConcurrentExpiredMap); ok {
			cur.ReSetAfterHalfexpire(
				concurrent_map.I64Key(dataId), dataInfo, meta.config.EMExpiredDefaultTime)
		}
	} else if needToSaveToCEM == 1 {
		DbLoaderRegistry[tableIndex].dataCur.Set(concurrent_map.I64Key(dataId), dataInfo)
	}

	return dataInfo, ifFind

}

func getDataInfoFromMongo(tableName string, filter bson.M, querySelect bson.M, dataObj interface{}) (interface{}, error) {
	mongoAddr := meta.config.Mongo
	if meta.UseMongConsul {
		node := meta.Resolver.DiscoverNode()
		if node != nil && len(node.Address) > 0 {
			mongoAddr = node.Address
		}
	}

	session, err := getSession()
	if err != nil {
		logger.Errorf("connect mongo %s failed. err: %s", mongoAddr, err.Error())
		return nil, err
	}
	defer session.Close()
	session.SetSocketTimeout(time.Duration(TimeoutMongo) * time.Second)
	session.SetSyncTimeout(time.Duration(TimeoutMongo) * time.Second)
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("new_adn").C(tableName)

	err = c.Find(filter).Select(querySelect).One(dataObj)
	if err != nil {
		return nil, err
	}

	return dataObj, nil
}

func GetMasAbtest() (value []smodel.MasUnitAbtest, ifFind bool) {
	obj, ifFind := GetMVConfigValue("MAS_ABTEST")
	if !ifFind {
		return nil, false
	}
	value, ifFind = obj.([]smodel.MasUnitAbtest)
	if !ifFind {
		logger.Errorf("GetMasAbtest error: config value fails to type cast to []smodel.MasUnitAbtest")
		return nil, false
	}
	return value, true
}

func getMasAbtest(jsonStr []byte) []smodel.MasUnitAbtest {
	value := []smodel.MasUnitAbtest{}
	err := jsoniter.Unmarshal(jsonStr, &value)
	if err != nil {
		return nil
	}
	return value
}

func GetAppPackageMtgID(appPackage string) int64 {
	// appPackage is number
	if _, err := strconv.Atoi(appPackage); err == nil {
		appPackage = "id" + appPackage
	}
	var (
		appPackageMtgID *smodel.AppPackageMtgID
		ifFind          bool
	)
	if tb_tools.IsAvailable(&smodel.AppPackageMtgID{}) { //判断 Treasure Box SDK 是否支持
		appPackageMtgID, ifFind = getAppPackageMtgIDFromTreasureBox(appPackage)
	} else {
		appPackageMtgID, ifFind = getAppPackageMtgIDFromConcurrentMap(appPackage)
	}
	if !ifFind {
		return 1700785270
	}

	return appPackageMtgID.MtgID
}
func getAppPackageMtgIDFromTreasureBox(appPackage string) (*smodel.AppPackageMtgID, bool) {
	defer mvutil.RecordUseTime("avg_get_app_package_tb", time.Now().UnixNano(), 1e3)
	reObj, ifFind := tb_tools.GetData(&smodel.AppPackageMtgID{}, map_key.StrKey(appPackage))
	if ifFind {
		if re, ok := reObj.(*smodel.AppPackageMtgID); ok {
			return re, ok
		}
	}
	return nil, false
}
func getAppPackageMtgIDFromConcurrentMap(appPackage string) (*smodel.AppPackageMtgID, bool) {
	defer mvutil.RecordUseTime("avg_get_app_package_cm", time.Now().UnixNano(), 1e3)
	value, ifFind := DbLoaderRegistry[TblAppPackageMTGId].dataCur.Get(concurrent_map.StrKey(appPackage))
	if !ifFind {
		return nil, false
	}
	obj, ok := value.(*smodel.AppPackageMtgID)
	return obj, ok
}

func GetSspProfitDistributionRule(key string) (*smodel.SspProfitDistributionRule, bool) {
	if tb_tools.IsAvailable(&smodel.SspProfitDistributionRule{}) { //判断 Treasure Box SDK 是否支持
		return getSspProfitDistributionRuleFromTreasureBox(key)
	}
	return getSspProfitDistributionRuleFromConcurrentMap(key)
}
func getSspProfitDistributionRuleFromTreasureBox(key string) (*smodel.SspProfitDistributionRule, bool) {
	defer mvutil.RecordUseTime("avg_get_ssp_profit_tb", time.Now().UnixNano(), 1e3)
	reObj, ifFind := tb_tools.GetData(&smodel.SspProfitDistributionRule{}, map_key.StrKey(key))
	if ifFind {
		if re, ok := reObj.(*smodel.SspProfitDistributionRule); ok {
			if (re.Type != mvconst.SspProfitDistributionRuleFixedEcpm && re.Type != mvconst.SspProfitDistributionRuleOnlineApiEcpm) || re.Status != mvutil.ACTIVE { //业务代码处理type的逻辑
				return nil, false
			}
			return re, ok
		}
	}
	return nil, false
}
func getSspProfitDistributionRuleFromConcurrentMap(key string) (*smodel.SspProfitDistributionRule, bool) {
	defer mvutil.RecordUseTime("avg_get_ssp_profit_cm", time.Now().UnixNano(), 1e3)
	value, found := DbLoaderRegistry[TblSSPProfitDistributionRule].dataCur.Get(concurrent_map.StrKey(key))
	if !found {
		return nil, false
	}
	if sspProfitDistributionRule, ok := value.(*smodel.SspProfitDistributionRule); ok {
		if (sspProfitDistributionRule.Type != mvconst.SspProfitDistributionRuleFixedEcpm && sspProfitDistributionRule.Type != mvconst.SspProfitDistributionRuleOnlineApiEcpm) || sspProfitDistributionRule.Status != mvutil.ACTIVE {
			return nil, false
		}
		return sspProfitDistributionRule, ok
	} else {
		return nil, false
	}
}

func GetSspProfitDistributionRuleByUnitIdAndCountryCode(unitId int64, countryCode string) (*smodel.SspProfitDistributionRule, bool) {
	countryCode = strings.ToUpper(countryCode)
	key := strconv.FormatInt(unitId, 10) + ":" + countryCode
	if value, found := GetSspProfitDistributionRule(key); found {
		return value, found
	} else {
		if countryCode != "ALL" {
			key = strconv.FormatInt(unitId, 10) + ":ALL"
			if value, found := GetSspProfitDistributionRule(key); found {
				return value, found
			}
		}
	}
	return nil, false
}

func GetAdvOffer(networkCid string) (*smodel.AdvOffer, bool) {
	if tb_tools.IsAvailable(&smodel.AdvOffer{}) { //判断 Treasure Box SDK 是否支持
		return getAdvOfferFromTreasureBox(networkCid)
	}
	return nil, false
}

func getAdvOfferFromTreasureBox(networkCid string) (*smodel.AdvOffer, bool) {
	defer mvutil.RecordUseTime("avg_get_adv_offer_tb", time.Now().UnixNano(), 1e3)
	reObj, ifFind := tb_tools.GetData(&smodel.AdvOffer{}, map_key.StrKey(networkCid))
	if ifFind {
		if re, ok := reObj.(*smodel.AdvOffer); ok {
			return re, ok
		}
	}
	return nil, false
}

func GetCreativePackage(cpId int64) (*smodel.CreativePackage, bool) {
	if tb_tools.IsAvailable(&smodel.CreativePackage{}) {
		return getCreativePackageFromTreasureBox(cpId)
	}
	return nil, false
}

func getCreativePackageFromTreasureBox(cpId int64) (*smodel.CreativePackage, bool) {
	defer mvutil.RecordUseTime("avg_get_creative_package_tb", time.Now().UnixNano(), 1e3)
	reObj, ifFind := tb_tools.GetData(&smodel.CreativePackage{}, map_key.I64Key(cpId))
	if ifFind {
		if re, ok := reObj.(*smodel.CreativePackage); ok {
			return re, ok
		}
	}
	return nil, false
}

func GetAdxHeaderBiddingPriceFactorConfByUnitIdAndCountryCode(unitId int64, countryCode string) (*smodel.AdxHeaderBiddingPriceFactor, bool) {
	countryCode = strings.ToUpper(countryCode)
	key := strconv.FormatInt(unitId, 10) + "_" + countryCode
	if value, found := GetAdxHeaderBiddingPriceFactor(key); found {
		return value, found
	} else {
		if countryCode != "ALL" {
			key = strconv.FormatInt(unitId, 10) + "_ALL"
			if value, found := GetAdxHeaderBiddingPriceFactor(key); found {
				return value, found
			}
		}
	}
	return nil, false
}

func GetAdxHeaderBiddingPriceFactor(key string) (*smodel.AdxHeaderBiddingPriceFactor, bool) {
	if tb_tools.IsAvailable(&smodel.AdxHeaderBiddingPriceFactor{}) {
		return getAdxHeaderBiddingPriceFactorFromTreasureBox(key)
	}
	return nil, false
}

func getAdxHeaderBiddingPriceFactorFromTreasureBox(key string) (*smodel.AdxHeaderBiddingPriceFactor, bool) {
	defer mvutil.RecordUseTime("avg_get_adx_hb_price_factor_tb", time.Now().UnixNano(), 1e3)
	reObj, ifFind := tb_tools.GetData(&smodel.AdxHeaderBiddingPriceFactor{}, map_key.StrKey(key))
	if ifFind {
		if re, ok := reObj.(*smodel.AdxHeaderBiddingPriceFactor); ok {
			return re, ok
		}
	}
	return nil, false
}

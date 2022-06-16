package extractor

import (
	"github.com/easierway/concurrent_map"
	"gitlab.mobvista.com/ADN/adnet/internal/expired_map"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type ConcurrentMapInterf interface {
	Get(key concurrent_map.Partitionable) (interface{}, bool)
	Set(key concurrent_map.Partitionable, v interface{})
	Del(key concurrent_map.Partitionable)
}

type dbLoaderInfo struct {
	collection      string
	querybmDbMap    bson.M                                                                                   // 基准 query 语句
	queryimDbMap    bson.M                                                                                   // 增量 query 语句
	querySelect     bson.M                                                                                   // 需要哪些字段
	updateFunc      func(query *mgo.Query, num int, dbLoaderInfo *dbLoaderInfo) (maxUpdate int64, err error) // 数据更新函数
	dataCur         ConcurrentMapInterf                                                                      // 内存 map
	notNeedToGetAll bool                                                                                     // 不需要全量加载
	updateByGetAll  bool                                                                                     // 是否来自 getAllMongoData
}

const concurrentMapPartitionsNum = 99

var (
	OrigUnitData         *concurrent_map.ConcurrentMap
	DbLoaderRegistry     [TblMax + 1]dbLoaderInfo
	DbLoaderRegistryLock [HotTabNum]*expired_map.ConccurentLock
	DbLoaderRegistryTime [TblMax + 1]int64

	QueryInc       = bson.M{}
	QueryAll       = bson.M{"status": 1}
	ConfigQueryAll = bson.M{"key": bson.M{"$exists": 1}}
	AppStatusArray = [2]int{1, 4}
	AppQueryAll    = bson.M{"app.status": bson.M{"$in": AppStatusArray}}
	AppQuerySelect = bson.M{"appId": 1, "JUMP_TYPE_CONFIG": 1, "realPackageName": 1, "landingPageVersion": 1,
		"updated": 1, "JUMP_TYPE_CONFIG_2": 1, "rewards": 1, "blackCategoryList": 1, "blackPackageList": 1, "directOfferConfig": 1,
		"publisher.publisherId": 1, "publisher.status": 1, "publisher.apiKey": 1, "publisher.type": 1,
		"app.appId": 1, "app.grade": 1, "app.status": 1, "app.platform": 1, "app.isIncent": 1, "app.name": 1,
		"app.devinfoEncrypt": 1, "app.btClass": 1, "app.frequencyCap": 1, "app.storekitLoading": 1, "app.devIdAllowNull": 1, "app.offerPreference": 1,
		"app.impressionCap": 1, "app.domain": 1, "app.domain_verify": 1, "app.plct": 1, "app.coppa": 1,
		"app.iabCategoryV2": 1, "app.storeUrl": 1, "app.m1Num": 1, "app.bundleId": 1, "app.officialName": 1, "app.ccpa": 1, "app.livePlatform": 1,
		"app.ntbarpt": 1, "app.ntbarpasbl": 1, "app.atatType": 1}
	UnitQueryAll    = bson.M{"unit.status": 1}
	UnitQuerySelect = bson.M{"unitId": 1, "appId": 1, "setting": 1, "adSourceCountry": 1, "adSourceData": 1,
		"updated": 1, "virtualReward": 1, "endcard": 1, "JUMP_TYPE_CONFIG": 1, "ecpmFloors": 1,
		"mp2mv": 1, "mv2mp": 1, "JUMP_TYPE_CONFIG_2": 1, "adSourceTime": 1, "CDNSetting": 1,
		"blackCategoryList": 1, "blackPackageList": 1, "directOfferConfig": 1,
		"unit.entraImage": 1, "unit.redPointShow": 1, "unit.redPointShowInterval": 1, "unit.unitId": 1,
		"unit.isIncent": 1, "unit.btClass": 1, "unit.status": 1, "unit.adType": 1, "unit.videoAds": 1,
		"unit.orientation": 1, "unit.templates": 1, "unit.nvTemplate": 1, "unit.videoEndType": 1,
		"unit.recallNet": 1, "unit.devIdAllowNull": 1, "unit.impressionCap": 1, "unit.entranceImg": 1,
		"unit.cookieAchieve": 1, "unit.hang": 1, "unit.endcardTemplate": 1, "unit.isServerCall": 1, "unit.entraTitle": 1,
		"unit.alac": 1, "unit.alecfc": 1, "unit.mof": 1, "unit.mofUnitId": 1, "fakeRuleV2": 1, "unit.plmtId": 1, "unit.biddingType": 1}
	CampaignQuerySelect = bson.M{"campaignId": 1, "advertiserId": 1, "trackingUrl": 1, "trackingUrlHttps": 1, "directUrl": 1,
		"price": 1, "oriPrice": 1, "cityCodeV2": 1, "status": 1, "network": 1, "previewUrl": 1, "packageName": 1,
		"campaignType": 1, "ctype": 1, "appSize": 1, "tag": 1, "adSourceId": 1, "publisherId": 1,
		"frequencyCap": 1, "directPackageName": 1, "sdkPackageName": 1, "advImp": 1, "adUrlList": 1, "jumpType": 1,
		"vbaConnecting": 1, "vbaTrackingLink": 1, "retargetingDevice": 1, "sendDeviceidRate": 1, "endcard": 1,
		"loopback": 1, "belongType": 1, "configVBA": 1, "appPostList": 1, "blackSubidListV2": 1, "btV4": 1,
		"openType": 1, "subCategoryName": 1, "isCampaignCreative": 1, "costType": 1, "source": 1,
		"JUMP_TYPE_CONFIG": 1, "chnId": 1, "thirdParty": 1, "updated": 1, "mChanlPrice": 1,
		"JUMP_TYPE_CONFIG_2": 1, "category": 1, "wxAppId": 1, "wxPath": 1, "bindId": 1, "deepLink": 1, "apkVersion": 1,
		"apkMd5": 1, "apkUrl": 1, "basicCrList": 1, "imageList": 1, "videoList": 1, "readCreative": 1, "createSrc": 1,
		"fakeCreative": 1, "alacRate": 1, "alecfcRate": 1, "mof": 1, "mCountryChanlPrice": 1, "mCountryChanlOriPrice": 1,
		"needToNotice3s": 1, "localPrice": 1, "countryChanlOriPrice": 1, "retargetOffer": 1, "skadnSupport": 1,
		"skadnCampaignId": 1, "skadnNeed": 1}
	AdvertiserQueryAll    = bson.M{"advertiser.status": 1}
	AdvertiserQuerySelect = bson.M{"advertiserId": 1, "advertiser": 1, "updated": 1}
	PublisherStatusArray  = [2]int{1, 2}
	PublisherQueryAll     = bson.M{"publisher.status": bson.M{"$in": PublisherStatusArray}}
	PublisherQuerySelect  = bson.M{"publisherId": 1, "updated": 1, "JUMP_TYPE_CONFIG": 1, "JUMP_TYPE_CONFIG_2": 1,
		"offsetList": 1, "publisher.publisherId": 1, "publisher.status": 1, "publisher.apiKey": 1, "publisher.type": 1, "publisher.openSource": 1}
	ConfigQuerySelect                  = bson.M{"key": 1, "value": 1, "updated": 1}
	ConfigCenterQuerySelect            = bson.M{"key": 1, "area": 1, "value": 1, "updated": 1}
	AppPackageMtgIDQuerySelect         = bson.M{"mtgId": 1, "appPackage": 1, "updated": 1}
	AdxDspConfigQuerySelect            = bson.M{"dspId": 1, "status": 1, "updated": 1, "target": 1}
	ConfigAlgorithmFillRateQuerySelect = bson.M{"uniqueKey": 1, "rate": 1, "updated": 1, "status": 1, "ecpmFloor": 1, "controlMode": 1}
	FreqControlFactorQuerySelect       = bson.M{"factorKey": 1, "factorRate": 1, "status": 1, "updated": 1}
	PlacementQuerySelect               = bson.M{"placementId": 1, "impressionCap": 1, "impressionCapPeriod": 1, "status": 1, "updated": 1}
	AdxTrafficMediaConfigQuerySelect   = bson.M{"unitId": 1, "trafficType": 1, "mode": 1, "updated": 1, "status": 1, "adType": 1, "deviceId": 1, "area": 1, "dspWhiteList": 1}

	SspProfitDistributionRuleQueryAll    = bson.M{"status": 1}
	SspProfitDistributionRuleQuerySelect = bson.M{"type": 1, "plusMax": 1, "fixedEcpm": 1, "unitId": 1, "area": 1, "isDailyDeficitCap": 1, "updated": 1, "updatedDate": 1, "status": 1}
)

func GetDbLoaderByCollection(col string) (int, *dbLoaderInfo, int64) {
	for i := 0; i < len(DbLoaderRegistry); i++ {
		if col == DbLoaderRegistry[i].collection {
			return i, &DbLoaderRegistry[i], DbLoaderRegistryTime[i]
		}
	}
	return 0, nil, 0
}

//func InitMConfigInstanceOnlyForUnitTest() {
//	meta.stopChan = make(chan bool)
//	ccObj = new(ConfigCenterObj)
//mvConfigObj = new(MVConfigObj)
//NewMDbLoaderRegistry()
//}

func NewMDbLoaderRegistry() {
	// DbLoaderRegistry = make(map[string]dbLoaderInfo)
	// DbLoaderRegistryLock = make(map[string]*expired_map.ConccurentLock)
	// DbLoaderRegistryTime = make(map[string]int)
	OrigUnitData = concurrent_map.CreateConcurrentMap(concurrentMapPartitionsNum)

	DbLoaderRegistry[TblAPP] = NewAppExtractor(AppQueryAll, QueryInc, AppQuerySelect)
	DbLoaderRegistryLock[TblAPP] = expired_map.CreateCoccurrentLock()

	DbLoaderRegistry[TblUnit] = NewUnitExtractor(UnitQueryAll, QueryInc, UnitQuerySelect)
	DbLoaderRegistryLock[TblUnit] = expired_map.CreateCoccurrentLock()

	DbLoaderRegistry[TblCampaign] = NewCampaignExtractor(QueryAll, QueryInc, CampaignQuerySelect)
	DbLoaderRegistryLock[TblCampaign] = expired_map.CreateCoccurrentLock()

	DbLoaderRegistry[TblAdvertiser] = NewAdvertiserExtractor(AdvertiserQueryAll, QueryInc, AdvertiserQuerySelect)

	DbLoaderRegistry[TblPublisher] = NewPublisherExtractor(PublisherQueryAll, QueryInc, PublisherQuerySelect)
	DbLoaderRegistryLock[TblPublisher] = expired_map.CreateCoccurrentLock()

	DbLoaderRegistry[TblConfig] = NewConfigExtractor(ConfigQueryAll, QueryInc, ConfigQuerySelect)

	DbLoaderRegistry[TblPlacement] = NewPlacementExtractor(QueryAll, QueryInc, PlacementQuerySelect)

	DbLoaderRegistry[TblConfigCenter] = NewConfigCenterExtractor(ConfigQueryAll, QueryInc, ConfigCenterQuerySelect)

	DbLoaderRegistry[TblAdxTrafficMediaConfig] = NewAdxTrafficMediaConfigExtractor(QueryAll, QueryInc, AdxTrafficMediaConfigQuerySelect)

	DbLoaderRegistry[TblConfigAlgorithmFillrate] = NewConfigAlgorithmFillRateExtractor(QueryAll, QueryInc, ConfigAlgorithmFillRateQuerySelect)

	DbLoaderRegistry[TblAdxDspConfig] = NewAdxDspConfigExtractor(QueryAll, QueryInc, AdxDspConfigQuerySelect)

	DbLoaderRegistry[TblAppPackageMTGId] = NewAppPackageMtgIDExtractor(QueryInc, QueryInc, AppPackageMtgIDQuerySelect)

	DbLoaderRegistry[TblSSPProfitDistributionRule] = NewSspProfitDistributionRuleExtractor(SspProfitDistributionRuleQueryAll, QueryInc, SspProfitDistributionRuleQuerySelect)

	DbLoaderRegistry[TblFreqControlFactor] = NewFreqControlFactorExtractor(QueryAll, QueryInc, FreqControlFactorQuerySelect)
}

func NewHBDbLoaderRegistry() {
	DbLoaderRegistry[TblAPP] = NewAppExtractor(AppQueryAll, QueryInc, AppQuerySelect)

	unitQueryAll := bson.M{"unit.status": 1, "unit.adType": bson.M{"$in": []int{42, 94, 287, 296, 297, 298}}}
	unitQueryInc := bson.M{"unit.adType": bson.M{"$in": []int{42, 94, 287, 296, 297, 298}}}
	DbLoaderRegistry[TblUnit] = NewUnitExtractor(unitQueryAll, unitQueryInc, UnitQuerySelect)

	DbLoaderRegistry[TblPublisher] = NewPublisherExtractor(PublisherQueryAll, QueryInc, PublisherQuerySelect)

	campaignQueryAll := bson.M{"status": 1, "advertiserId": 903}
	campaignQueryInc := bson.M{"advertiserId": 903}
	DbLoaderRegistry[TblCampaign] = NewCampaignExtractor(campaignQueryAll, campaignQueryInc, CampaignQuerySelect)

	DbLoaderRegistry[TblConfig] = NewConfigExtractor(ConfigQueryAll, QueryInc, ConfigQuerySelect)

	DbLoaderRegistry[TblConfigCenter] = NewConfigCenterExtractor(ConfigQueryAll, QueryInc, ConfigCenterQuerySelect)

	DbLoaderRegistry[TblAdvertiser] = NewAdvertiserExtractor(AdvertiserQueryAll, QueryInc, AdvertiserQuerySelect)

	DbLoaderRegistry[TblAppPackageMTGId] = NewAppPackageMtgIDExtractor(QueryInc, QueryInc, AppPackageMtgIDQuerySelect)

	DbLoaderRegistry[TblAdxDspConfig] = NewAdxDspConfigExtractor(QueryAll, QueryInc, AdxDspConfigQuerySelect)

	DbLoaderRegistry[TblConfigAlgorithmFillrate] = NewConfigAlgorithmFillRateExtractor(QueryAll, QueryInc, ConfigAlgorithmFillRateQuerySelect)

	DbLoaderRegistry[TblSSPProfitDistributionRule] = NewSspProfitDistributionRuleExtractor(SspProfitDistributionRuleQueryAll, QueryInc, SspProfitDistributionRuleQuerySelect)

	DbLoaderRegistry[TblFreqControlFactor] = NewFreqControlFactorExtractor(QueryAll, QueryInc, FreqControlFactorQuerySelect)

	DbLoaderRegistry[TblPlacement] = NewPlacementExtractor(QueryAll, QueryInc, PlacementQuerySelect)

	DbLoaderRegistry[TblAdxTrafficMediaConfig] = NewAdxTrafficMediaConfigExtractor(QueryAll, QueryInc, AdxTrafficMediaConfigQuerySelect)
}

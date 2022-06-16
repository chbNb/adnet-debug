package extractor

import (
	"errors"
	"fmt"
	"runtime/debug"
	"strings"
	"time"

	"github.com/easierway/go-kit/balancer"
	mlogger "github.com/mae-pax/logger"
	"gitlab.mobvista.com/ADN/adnet/internal/expired_map"
	"gitlab.mobvista.com/ADN/adnet/internal/hot_data"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	"gitlab.mobvista.com/ADN/adnet/internal/req_context"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type extrMeta struct {
	stopChan      chan bool
	mgoSession    *mgo.Session
	mgoCollection *mgo.Collection
	config        *mvutil.ExtraConfig
	ConsulConfig  *mvutil.ConsulConfig
	Resolver      *balancer.ConsulResolver
	UseMongConsul bool
}

var meta extrMeta

var logger *mlogger.Log

func setLog(log *mlogger.Log) {
	logger = log
}

func InitConfig(log *mlogger.Log) {
	setLog(log)
	//initMVConfigObj()
	//initCCObj()
	//decodeConfig()
}

func InitExtractor(config *mvutil.ExtraConfig, consulConfig *mvutil.ConsulConfig, log *mlogger.Log, useMongoConsul bool) error {
	//if ccObj == nil {
	//	ccObj = new(ConfigCenterObj)
	//}
	//if mvConfigObj == nil {
	//	mvConfigObj = new(MVConfigObj)
	//}
	meta.stopChan = make(chan bool)
	meta.config = config

	logger = log
	if useMongoConsul {
		meta.ConsulConfig = consulConfig
		Resolver, err := balancer.NewConsulResolver(consulConfig.Address, consulConfig.Service, consulConfig.MyService,
			time.Duration(consulConfig.Internal)*time.Millisecond, consulConfig.ServiceRatio, consulConfig.CpuThreshold)
		if err != nil {
			logger.Errorf("NewConsulResolver error: %s will use Mongo ELB.", err.Error())
		} else {
			logger.Info("NewConsulResolver ok will use Mongo Consul.")
			meta.UseMongConsul = true
			Resolver.SetLogger(logger)
			meta.Resolver = Resolver
		}
	}
	for i, table := range meta.config.DbConfig {
		col, _, _ := GetDbLoaderByCollection(table.Collection)
		meta.config.DbConfig[i].Index = col
	}
	if config.UseExpiredMap {
		logger.Infof("UseExpiredMap, EMBatchDeleteTime: %d; EMExpiredDefaultTime: %d; EMRetryAgrainSleepMicrosecondM: %d",
			meta.config.EMBatchDeleteTime, meta.config.EMExpiredDefaultTime, meta.config.EMRetryAgrainSleepMicrosecond)
		useExpirdMapByKey(TblCampaign, "campaignId", bson.M{"status": 1})
		useExpirdMapByKey(TblPublisher, "publisherId", bson.M{})
		useExpirdMapByKey(TblAPP, "appId", bson.M{})
		useExpirdMapByKey(TblUnit, "unitId", bson.M{})
	}
	logger.Infof("mongoDataExtr get config: %+v", config)
	if err := getAllMongoData(); err != nil {
		logger.Errorf("getAllMongoData failed, err: %s", err.Error())
		return err
	}
	logger.Info("mongoDataExtr getAllMongoData success, init success.")
	return nil
}

func useExpirdMapByKey(tableIndex int, keyName string, extBsonM bson.M) {
	camp := DbLoaderRegistry[tableIndex] // 拿不到让他直接报错好了，一定有问题的
	tableName := camp.collection
	camp.dataCur = expired_map.CreateConcurrentExpiredMap(
		concurrentMapPartitionsNum,
		meta.config.EMBatchDeleteTime,
		meta.config.EMExpiredDefaultTime)
	camp.notNeedToGetAll = true
	ids, err := hot_data.GetActiveDatas(tableName)
	if err != nil {
		logger.Infof("Get %s data from Redis Error: %s", tableName, err.Error())
	} else {
		logger.Infof("Get %s data from Redis is: %d", tableName, len(ids))
		if len(ids) > 0 {
			camp.querybmDbMap = bson.M{keyName: bson.M{"$in": ids}}
			for k, v := range extBsonM {
				camp.querybmDbMap[k] = v
			}
			camp.notNeedToGetAll = false
		}
	}
	DbLoaderRegistry[tableIndex] = camp
	DbLoaderRegistryTime[tableIndex] = time.Now().Unix() - 3 // 让自增的从前3s前开始获取数据
}

//func DecodeConfigWorking() {
//	logger.Info("DecodeConfigWorking")
//	go func() {
//		inc := time.After(time.Duration(meta.config.ModifyInterval*meta.config.IntervalFactor) * time.Second)
//		for {
//			select {
//			case <-inc:
//				// logger.Infof("DecodeConfigWorking %d", meta.config.ModifyInterval*meta.config.IntervalFactor)
//				DecodeConfig()
//				inc = time.After(time.Duration(meta.config.ModifyInterval*meta.config.IntervalFactor) * time.Second)
//			case <-meta.stopChan:
//				logger.Info("mongoExtr get stop signal, will stop")
//				return
//			}
//		}
//	}()
//}

func decodeConfig() {
	//ccObj.DomainTrack = ""
	//ccObj.DomainTracks = make(map[string]mvutil.IRate)
	//ccObj.Domain = ""
	//ccObj.System = ""
	//ccObj.SYSTEM_AREA = ""
	//ccObj.MP_DOMAIN_CONF = nil
	//ccObj.JSSDK_DOMAIN_TRACK = ""
	//ccObj.CHET_DOMAIN = ""
	//ccObj.CLOUD_NAME = ""
	//ccObj.RateLimit = nil

	//mvConfigObj.TreasureBoxIDCount = getTreasureBoxIDCount(nil)
	//
	//mvConfigObj.OpenapiScenaria, _ = getOpenapiScenario(nil)
	//mvConfigObj.LowGradeAppId, _ = getConfigLowGradeApp(nil)
	//mvConfigObj.ReplenishApp, _ = getReplenishApp(nil)
	//mvConfigObj.AdvBlackSubIdListV2, _ = getAdvBlackSubIdList(nil)
	//mvConfigObj.AppPostListV2, _ = getAppPostList(nil)
	//mvConfigObj.OfferwallGuidelines, _ = getOfferwallGuidelines(nil)
	//mvConfigObj.EndCardConfig, _ = getEndcard(nil)
	//mvConfigObj.DefRvTemp, _ = getDefRVTemplate(nil)
	//mvConfigObj.VersionCompare, _ = getVersionCompare(nil)
	//mvConfigObj.PlayableTest, _ = getPlayableTest(nil)
	//mvConfigObj.TrackUrlConfigNew, _ = getTRACK_URL_CONFIG_NEW(nil)
	//mvConfigObj.S3TrackingDomain, _ = get3S_CHINA_DOMAIN(nil)
	//mvConfigObj.JumpTypeConfig, _ = getJUMP_TYPE_CONFIG(nil)
	//mvConfigObj.JumpTypeIOS, _ = getJUMP_TYPE_CONFIG_IOS(nil)
	//mvConfigObj.JumpTypeSDKVersion, _ = getJUMPTYPE_SDKVERSION(nil)
	//mvConfigObj.AdStacking, _ = getADSTACKING(nil)
	//mvConfigObj.SettingConfig, _ = getSETTING_CONFIG(nil)
	//mvConfigObj.Template, _ = getTemplate(nil)
	//mvConfigObj.OfferWallUrls, _ = getofferwall_urls(nil)
	//mvConfigObj.NoCheckApps, _ = getSIGN_NO_CHECK_APPS(nil)
	//mvConfigObj.AdServerTestConfig, _ = getAdServerTestConfig(nil)
	//mvConfigObj.CToi, _ = getCToi(nil)
	//mvConfigObj.MP_MAP_UNIT, _ = getMP_MAP_UNIT(nil)
	//mvConfigObj.DEL_PRICE_FLOOR_UNITS, _ = getDEL_PRICE_FLOOR_UNIT(nil)
	//mvConfigObj.DecodeSwith, _ = getDecoodeSwith(nil)
	//mvConfigObj.IA_APPWALL, _ = getIA_APPWALL(nil)
	//mvConfigObj.CAN_CLICK_MODE_SIX_PUBLISHER, _ = getCAN_CLICK_MODE_SIX_PUBLISHER(nil)
	//mvConfigObj.TOUTIAO_ALGO, _ = getTOUTIAO_ALGO(nil)
	//mvConfigObj.TO_NEW_CDN_APPS, _ = getTO_NEW_CDN_APPS(nil)
	//mvConfigObj.UA_AB_TEST_CONFIG, _ = getUA_AB_TEST_CONFIG(nil)
	//mvConfigObj.GO_TRACK, _ = getGO_TRACK(nil)
	//mvConfigObj.JUMP_TYPE_CONFIG_ADV, _ = getJUMP_TYPE_CONFIG_ADV(nil)
	//mvConfigObj.FUN_MV_MAP, _ = getFUN_MV_MAP(nil)
	//mvConfigObj.NO_CHECK_PARAM_APP, _ = getNO_CHECK_PARAM_APP(nil)
	//mvConfigObj.HAS_EXTENSIONS_UNIT, _ = getHAS_EXTENSIONS_UNIT(nil)
	//mvConfigObj.GAMELOFT_CREATIVE_URLS, _ = getGAMELOFT_CREATIVE_URLS(nil)
	//mvConfigObj.HAS_EXTENSIONS_UNIT, _ = getHAS_EXTENSIONS_UNIT(nil)
	//mvConfigObj.HUPU_DEFAULT_UNITID, _ = getHUPU_DEFAULT_UNITID(nil)
	//mvConfigObj.HUPU_DEFAULT_PRICE, _ = getHUPU_DEFAULT_PRICE(nil)
	//mvConfigObj.UA_AB_TEST_SDK_OS_CONFIG, _ = getUA_AB_TEST_SDK_OS_CONFIG(nil)
	//mvConfigObj.UA_AB_TEST_THIRD_PARTY_CONFIG, _ = getUA_AB_TEST_THIRD_PARTY_CONFIG(nil)
	//mvConfigObj.UA_AB_TEST_CAMPAIGN_CONFIG, _ = getUA_AB_TEST_CAMPAIGN_CONFIG(nil)
	//mvConfigObj.ONLINE_PRICE_FLOOD_APPID, _ = getONLINE_PRICE_FLOOR_APPID(nil)
	//mvConfigObj.CREATIVE_CHECK_MIGU_ADX_CREATIVE_IDS, _ = getCREATIVE_CHECK_MIGU_ADX_CREATIVE_IDS(nil)
	//mvConfigObj.JUMP_TYPE_CONFIG_THIRD_PARTY, _ = getJUMP_TYPE_CONFIG_THIRD_PARTY(nil)
	//mvConfigObj.ONLINE_PRICE_FLOOR, _ = getONLINE_PRICE_FLOOR(nil)
	//mvConfigObj.MAX_APPID_AND_UNITID, _ = getMAX_APPID_AND_UNITID(nil)
	//mvConfigObj.LINKTYPE_UNITID, _ = getLINKTYPE_UNITID(nil)
	////mvConfigObj.CHEETAH_CONFIG, _ = getCHEETAH_CONFIG(nil)
	//mvConfigObj.PLAYABLE_TEST_UNITS, _ = getPLAYABLE_TEST_UNITS(nil)
	//mvConfigObj.ABTEST_GAIDIDFA, _ = getABTEST_GAIDIDFA(nil)
	//mvConfigObj.REPLACE_BRAND_MODEL, _ = getREPLACE_BRAND_MODEL(nil)
	////mvConfigObj.REN_CLE_RES_SDK_OS_CONFIG, _ = getREN_CLE_RES_SDK_OS_CONFIG(nil)
	//mvConfigObj.CREATIVE_CHECK_HUPU_ADX_CREATIVE_IDS, _ = getCREATIVE_CHECK_HUPU_ADX_CREATIVE_IDS(nil)
	//mvConfigObj.NEW_PLAYABLE_SWITCH, _ = getNEW_PLAYABLE_SWITCH(nil)
	//mvConfigObj.PLAYABLE_ABTEST_RATE, _ = getPLAYABLE_ABTEST_RATE(nil)
	//mvConfigObj.PUB_CC_CDN, _ = getPUB_CC_CDN(nil)
	//mvConfigObj.TP_WITHOUT_RV, _ = getTP_WITHOUT_RV(nil)
	//mvConfigObj.PPTV_APPID_AND_UNITID, _ = getPPTV_APPID_AND_UNITID(nil)
	//mvConfigObj.REQUEST_BLACKLIST, _ = getREQUEST_BLACKLIST(nil)
	//mvConfigObj.IFENG_APPID_AND_UNITID, _ = getIFENG_APPID_AND_UNITID(nil)
	//mvConfigObj.PPTV_DEAL_ID, _ = getPPTV_DEAL_ID(nil)
	//mvConfigObj.CREATIVE_CHECK_PPTV_ADX_CREATIVE_IDS, _ = getCREATIVE_CHECK_PPTV_ADX_CREATIVE_IDS(nil)
	//mvConfigObj.CM_APPID_AND_UNITID, _ = getCM_APPID_AND_UNITID(nil)
	//mvConfigObj.FCA_SWITCH, _ = getFCA_SWITCH(nil)
	//mvConfigObj.FCA_CAMIDS, _ = getFCA_CAMIDS(nil)
	//mvConfigObj.NEW_DEFAULT_FCA, _ = getNEW_DEFAULT_FCA(nil)
	//mvConfigObj.CM_MP_APPIDS, _ = getCM_MP_APPIDS(nil)
	//mvConfigObj.LOW_FLOW_UNITS, _ = getLOW_FLOW_UNITS(nil)
	//mvConfigObj.NEW_CREATIVE_ABTEST_RATE, _ = getNEW_CREATIVE_ABTEST_RATE(nil)
	//mvConfigObj.NEW_CREATIVE_TEST_UNITS, _ = getNEW_CREATIVE_TEST_UNITS(nil)
	//mvConfigObj.TEMPLATE_MAP, _ = getTEMPLATE_MAP(nil)
	//mvConfigObj.CREATIVE_ABTEST, _ = getCREATIVE_ABTEST(nil)
	//mvConfigObj.BACKEND_BT, _ = getBackendBT(nil)
	//mvConfigObj.SS_ABTEST_CAMPAIGN, _ = getSS_ABTEST_CAMPAIGN(nil)
	//mvConfigObj.NEW_CREATIVE_TIMEOUT, _ = getNEW_CREATIVE_TIMEOUT(nil)
	//mvConfigObj.XM_NEW_RETURN_UNITS, _ = getXM_NEW_RETURN_UNITS(nil)
	//mvConfigObj.WEBVIEW_PRISON_PUBIDS_APPIDS, _ = getWEBVIEW_PRISON_PUBIDS_APPIDS(nil)
	//mvConfigObj.CONFIG_OFFER_PLCT, _ = getCONFIG_OFFER_PLCT(nil)
	//mvConfigObj.CCT_ABTEST_CONF, _ = getCCT_ABTEST_CONF(nil)
	//mvConfigObj.JSSDK_DOMAIN, _ = getJSSDK_DOMAIN(nil)
	//mvConfigObj.ADSOURCE_PACKAGE_BLACKLIST, _ = getADSOURCE_PACKAGE_BLACKLIST(nil)
	//mvConfigObj.CHET_URL_UNIT, _ = getCHET_URL_UNIT(nil)
	//mvConfigObj.MD5_FILE_PRISON_PUB, _ = getMD5_FILE_PRISON_PUB(nil)
	//mvConfigObj.ANR_PRISION_BY_PUB_AND_OSV, _ = getANR_PRISION_BY_PUB_AND_OSV(nil)
	//mvConfigObj.UNIT_WITHOUT_VIDEO_START, _ = getUNIT_WITHOUT_VIDEO_START(nil)
	//mvConfigObj.REDUCE_FILL_SWITCH, _ = getREDUCE_FILL_SWITCH(nil)
	//mvConfigObj.MORE_OFFER_CONF, _ = getMORE_OFFER_CONF(nil)
	//mvConfigObj.NEW_MOF_IMP_RATE, _ = getNEW_MOF_IMP_RATE(nil)
	//mvConfigObj.MOF_ABTEST_RATE, _ = getMOF_ABTEST_RATE(nil)
	//mvConfigObj.GDT_KEY_WHITELIST, _ = getGdtKeyWhiteList(nil)
	//mvConfigObj.ADX_FAKE_DSP_PRICE_DEFAULT_FACTOR, _ = getADX_FAKE_DSP_PRICE_DEFAULT_FACTOR(nil)
	//mvConfigObj.ADNET_SWITCHS, _ = getADNET_SWITCHS(nil)
	//mvConfigObj.BLACK_FOR_EXCLUDE_PACKAGE_NAME, _ = getBLACK_FOR_EXCLUDE_PACKAGE_NAME(nil)
	//mvConfigObj.APPWALL_TO_MORE_OFFER_UNIT, _ = getAPPWALL_TO_MORE_OFFER_UNIT(nil)
	//mvConfigObj.REQ_TYPE_AAB_TEST_CONFIG, _ = getREQ_TYPE_AAB_TEST_CONFIG(nil)
	//mvConfigObj.CleanEmptyDeviceABTest, _ = getCleanEmptyDeviceABTest(nil)
	//mvConfigObj.REDUCE_FILL_FLINK_SWITCH, _ = getREDUCE_FILL_FLINK_SWITCH(nil)
	//mvConfigObj.OnlineEmptyDeviceNoServerJump, _ = getOnlineEmptyDeviceNoServerJump(nil)
	//mvConfigObj.ExcludeDisplayPackageABTest, _ = getExcludeDisplayPackageABTest(nil)
	//mvConfigObj.AUTO_LOAD_CACHE_ABTSET_CONFIG, _ = getAUTO_LOAD_CACHE_ABTSET(nil)
	//mvConfigObj.OnlineEmptyDeviceIPUAABTest, _ = getOnlineEmptyDeviceIPUAABTest(nil)
	//mvConfigObj.CTN_SIZE_ABTEST, _ = getCTN_SIZE_ABTEST(nil)
	//mvConfigObj.BLOCK_MORE_OFFER_CONFIG, _ = getBLOCK_MORE_OFFER_CONFIG(nil)
	//mvConfigObj.CLOSE_BUTTON_AD_TEST_UNITS, _ = getCLOSE_BUTTON_AD_TEST_UNITS(nil)
	//mvConfigObj.MIGRATION_ABTEST, _ = getMIGRATION_ABTEST(nil)
	//mvConfigObj.RETURN_PARAM_K_UNIT, _ = getRETURN_PARAM_K_UNIT(nil)
	//mvConfigObj.ADD_UNIT_ID_ON_TRACKING_URL_UNIT, _ = getADD_UNIT_ID_ON_TRACKING_URL_UNIT(nil)
	//mvConfigObj.CLICK_IN_ADTRACKING_THIRD_PARTY_CONFIG, _ = getCLICK_IN_ADTRACKING_THIRD_PARTY_CONFIG(nil)
	//mvConfigObj.PsbReduplicateConfig = getPsbReduplicateConfig(nil)
	//mvConfigObj.ALAC_PRISON_CONFIG, _ = getALAC_PRISON_CONFIG(nil)
	//mvConfigObj.ABTEST_FIELDS, _ = getABTEST_FIELDS(nil)
	//mvConfigObj.ABTEST_CONFS, _ = getABTEST_CONFS(nil)
	//mvConfigObj.ChetLinkConfig, _ = getChetLinkConfigs(nil)
	//mvConfigObj.AndroidLowVersionFilterCondition, _ = getAndroidLowVersionFilterCondition(nil)
	//mvConfigObj.BANNER_HTML_STR, _ = getBANNER_HTML_STR(nil)
	//mvConfigObj.ThirdPartyWhiteList = getThirdPartyWhiteList(nil)
	//mvConfigObj.STOREKIT_TIME_PRISON_CONF = getSTOREKIT_TIME_PRISON_CONF(nil)
	//mvConfigObj.IPUA_WHITE_LIST_RATE_CONF = getIPUA_WHITE_LIST_RATE_CONF(nil)
	//mvConfigObj.IS_RETURN_PAUSEMOD = getIS_RETURN_PAUSEMOD(nil)
	//mvConfigObj.CREATIVE_CHECK_MANGGUO_ADX_CREATIVE_IDS = getCREATIVE_CHECK_MANGGUO_ADX_CREATIVE_IDS(nil)
	//mvConfigObj.MANGO_APPID_AND_UNITID = getMANGO_APPID_AND_UNITID(nil)
	//mvConfigObj.CLICK_IN_SERVER_CONF_NEW = getCLICK_IN_SERVER_CONF_NEW(nil)
	//mvConfigObj.ExcludeClickPackages = getExcludeClickPackages(nil)
	//mvConfigObj.ADNET_CONF_LIST = getADNET_CONF_LIST(nil)
	//mvConfigObj.AdChoiceConfig = getAdChoiceConfigData(nil)
	//mvConfigObj.FillRateEcpmFloorSwitch = getFillRateEcpmFloorSwitch(nil)
	//mvConfigObj.MoattagConfig = getMoattagConfig(nil)
	//mvConfigObj.BAD_REQUEST_FILTER_CONF = getBAD_REQUEST_FILTER_CONF(nil)
	//mvConfigObj.BIG_TEMPLATE_CONF = getBIG_TEMPLATE_CONF(nil)
	//mvConfigObj.CN_TRACKING_DOMAIN_CONF = getCNTrackingDomainConf(nil)
	//mvConfigObj.NoticeClickUniq = getNoticeClickUniq(nil)
	////mvConfigObj.AdjustPostbackABTest = getAdjustPostbackABTest(nil)
	//mvConfigObj.ExcludeImpressionPackages = getExcludeImpressionPackages(nil)
	//mvConfigObj.PassthroughData = getPassthroughData(nil)
	//mvConfigObj.PolarisFlagConf = getPolarisFlagConf(nil)
	//mvConfigObj.AdnetLangCreativeABTestConf = getAdnetLangCreativeABTestConf(nil)
	//mvConfigObj.FreqControlConf = getFREQ_CONTROL_CONFIG(nil)
	//mvConfigObj.TimezonConfig = getTIMEZONE_CONFIG(nil)
	//mvConfigObj.CountryCodeTimezoneConfig = getCOUNTRY_CODE_TIMEZONE_CONFIG(nil)
	//mvConfigObj.DspTplAbtestConf = getDSP_TPL_ABTEST_CONF(nil)
	//mvConfigObj.BannerToAsABTestConf = getBANNER_TO_AS_ABTEST_CONF(nil)
	//mvConfigObj.ONLY_REQUEST_THIRD_DSP_SWITCH = getONLY_REQUEST_THIRD_DSP_SWITCH(nil)
	//mvConfigObj.DcoTestConf = getDcoTestConf(nil)
	//mvConfigObj.USE_PLACEMENT_IMP_CAP_SWITCH = getUSE_PLACEMENT_IMP_CAP_SWITCH(nil)
	//mvConfigObj.LOW_FLOW_ADTYPE, _ = getLOW_FLOW_ADTYPE(nil)
	//mvConfigObj.TaoBaoOfferID = getTaoBaoOfferID(nil)
	//mvConfigObj.ADNET_DEFAULT_VALUE = getADNET_DEFAULT_VALUE(nil)
	//mvConfigObj.CREATIVE_COMPRESS_ABTEST_CONF_V3 = getADNET_CREATIVE_COMPRESS_ABTEST_CONF_V3(nil)
	//mvConfigObj.CREATIVE_VIDEO_COMPRESS_ABTEST = getCREATIVE_VIDEO_COMPRESS_ABTEST(nil)
	//mvConfigObj.CREATIVE_IMG_COMPRESS_ABTEST = getCREATIVE_IMG_COMPRESS_ABTEST(nil)
	//mvConfigObj.CREATIVE_ICON_COMPRESS_ABTEST = getCREATIVE_ICON_COMPRESS_ABTEST(nil)
	//mvConfigObj.HBBidSdkVersionConfig = getHBBidSdkVersionConfig(nil)
	//mvConfigObj.HBExchangeRate = getHBExchangeRate(nil)
	//mvConfigObj.HBPublisherCurrency = getHBPublisherCurrency(nil)
	//mvConfigObj.HBBlacklist = getHBBlacklist(nil)
	//mvConfigObj.HBPubExcludeAppCheck = getHBPubExcludeAppCheck(nil)
	//mvConfigObj.HBBidFloorPubWhiteList = getHBBidFloorPubWhiteList(nil)
	//mvConfigObj.HBAdxEndpoint = getHBAdxEndpoint(nil)
	//mvConfigObj.HBAdxEndpointV2 = getHBAdxEndpointV2(nil)
	//mvConfigObj.HBCDNDomainABTestPubs = getHBCDNDomainABTestPubs(nil)
	//mvConfigObj.AerospikeUseConsul = getAerospikeUseConsul(nil)
	//mvConfigObj.UseConsulServices = getUseConsulServices(nil)
	//mvConfigObj.UseConsulServicesV2 = getUseConsulServicesV2(nil)
	//mvConfigObj.CountryBlackPackageListConf = getCountryBlackPackageListConf(nil)
	//mvConfigObj.OnlineApiPubBidPriceConf = getOnlineApiPubBidPriceConf(nil)
	//mvConfigObj.SupportTrackingTemplateConf = getSupportTrackingTemplateConf(nil)
	//mvConfigObj.YLHClickModeTestConfig = getYLHClickModeTestConfig(nil)
	//mvConfigObj.RETRUN_VAST_APP = getRETRUN_VAST_APP(nil)
	//mvConfigObj.V5AbtestConf = getV5_ABTEST_CONFIG(nil)
	//mvConfigObj.REPLACE_TRACKING_DOMAIN_CONF = getREPLACE_TRACKING_DOMAIN_CONF(nil)
	//mvConfigObj.NEW_AFREECATV_UNIT = getNEW_AFREECATV_UNIT(nil)
	//mvConfigObj.SHAREIT_AF_WHITE_LIST_CONF = getSHAREIT_AF_WHITE_LIST_CONF(nil)
	//mvConfigObj.ONLINE_API_USE_ADX = getONLINE_API_USE_ADX(nil)
	//mvConfigObj.ONLINE_API_USE_ADX_MAS = getONLINE_API_USE_ADX_MAS(nil)
	//mvConfigObj.ONLINE_API_SUPPORT_DEEPLINK_V2 = getONLINE_API_SUPPORT_DEEPLINK_V2(nil)
	//mvConfigObj.ONLINE_API_MAX_BID_PRICE = getONLINE_API_MAX_BID_PRICE(nil)
	//mvConfigObj.NEW_HUPU_UNITID_MAP = getNEW_HUPU_UNITID_MAP(nil)
	//mvConfigObj.NEW_CDN_TEST = getNEW_CDN_TEST(nil)
	//mvConfigObj.ONLINE_FREQUENCY_CONTROL_CONF = getONLINE_FREQUENCY_CONTROL_CONF(nil)
	//mvConfigObj.CHANGE_ONLINE_DEEPLINK_WAY_PACKAGE_LIST = getCHANGE_ONLINE_DEEPLINK_WAY_PACKAGE_LIST(nil)
	//mvConfigObj.DEBUG_BID_FLOOR_AND_BID_PRICE_CONF = getDEBUG_BID_FLOOR_AND_BID_PRICE_CONF(nil)
	//mvConfigObj.H265_VIDEO_ABTEST_CONF = getH265_VIDEO_ABTEST_CONF(nil)
	//mvConfigObj.REPLACE_TEMPLATE_URL_CONF = getREPLACE_TEMPLATE_URL_CONF(nil)
	//mvConfigObj.TPL_CREATIVE_DOMAIN_CONF = getTPL_CREATIVE_DOMAIN_CONF(nil)
	//mvConfigObj.RETURN_WTICK_CONF = getRETURN_WTICK_CONF(nil)
	//mvConfigObj.EXCLUDE_PACKAGES_BY_CITYCODE_CONF = getEXCLUDE_PACKAGES_BY_CITYCODE_CONF(nil)
	//mvConfigObj.ADNET_PARAMS_ABTEST_CONFS = getADNET_PARAMS_ABTEST_CONFS(nil)
	logger.Info("DecodeConfigWorking runing end")
}

func Working() {
	logger.Info("mongoExtr will working")
	go func() {
		inc := time.After(time.Duration(meta.config.ModifyInterval) * time.Second)
		for {
			select {
			case <-inc:
				getIncMongoData()
				// getOtherConfig()
				inc = time.After(time.Duration(meta.config.ModifyInterval) * time.Second)
			case <-meta.stopChan:
				logger.Info("mongoExtr get stop signal, will stop")
				return
			}
		}
	}()
}

func Stop() {
	meta.stopChan <- true
}

// getAll  获取全量
func getAllMongoData() error {
	for _, dbconf := range meta.config.DbConfig {
		startTime := time.Now().UnixNano() / 1e6
		dbLoaderInfo := DbLoaderRegistry[dbconf.Index]
		if err := connect(meta.config.Db, dbconf.Collection); err != nil {
			logger.Errorf("mongo connect failed, db: %s, collection: %s, err: %s", meta.config.Db, dbconf.Collection, err.Error())
			return err
		}
		if dbLoaderInfo.notNeedToGetAll {
			logger.Infof("%s set notNeedToGetAll is true", dbconf.Collection)
			continue
		}
		// query
		res := meta.mgoCollection.Find(dbLoaderInfo.querybmDbMap).Select(dbLoaderInfo.querySelect)
		num, err := res.Count()
		if err != nil {
			logger.Errorf("getAllMongoData query error, db: %s, collection: %s, error: %s", meta.config.Db, dbconf.Collection, err.Error())
			return err
		}
		logger.Infof("getAllMongoData, Sql: %+v, db: %s, collection: %s, result num: %d", dbLoaderInfo.querybmDbMap, meta.config.Db, dbconf.Collection, num)
		dbLoaderInfo.updateByGetAll = true // 用于标记从哪进入到updateFunc的
		maxUpdate, err := dbLoaderInfo.updateFunc(res, num, &dbLoaderInfo)
		if err != nil {
			logger.Errorf(" getAllMongoData all update failed, db: %s, collection: %s, err: %s", meta.config.Db, dbconf.Collection, err.Error())
			return err
		}

		lastUpdate := time.Now().Unix()
		if num > 0 && maxUpdate > 0 {
			lastUpdate = maxUpdate
		}
		DbLoaderRegistryTime[dbconf.Index] = lastUpdate
		logger.Infof("extractor getAllMongoData db: %s, collection: %s, num: %d, mongo time: %d, use_time_ms: %d", meta.config.Db, dbconf.Collection, num, DbLoaderRegistryTime[dbconf.Index], time.Now().UnixNano()/1e6-startTime)
	}
	return nil
}

// getIncMongoData  获取增量
func getIncMongoData() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			logger.Warnf("%#v", err)
			debug.PrintStack()
		}
	}()
	logger.Info("process serverip info")
	req_context.GetInstance().UpdateServerIpInfo()
	logger.Info("getIncMongoData begin")
	for _, dbconf := range meta.config.DbConfig {
		startTime := time.Now().UnixNano() / 1e6
		dbLoaderInfo := DbLoaderRegistry[dbconf.Index]
		lastUpdate := DbLoaderRegistryTime[dbconf.Index]
		if err := connect(meta.config.Db, dbconf.Collection); err != nil {
			logger.Errorf("mongo connect failed, db: %s, collection: %s, err: %s", meta.config.Db, dbconf.Collection, err.Error())
			continue
		}
		if lastUpdate == 0 && !(dbconf.Collection == "config" || dbconf.Collection == "configcenter") {
			lastUpdate = time.Now().Unix() - int64(meta.config.ModifyInterval)
		}

		queryimDbMap := dbLoaderInfo.queryimDbMap
		queryimDbMap["updated"] = bson.M{"$gte": lastUpdate - int64(meta.config.UpdateOffset), "$lte": int(time.Now().Unix())}
		// query
		res := meta.mgoCollection.Find(queryimDbMap).Select(dbLoaderInfo.querySelect)
		num, err := res.Count()
		if err != nil {
			logger.Errorf("getIncMongoData query error, db: %s, collection: %s, error: %s", meta.config.Db, dbconf.Collection, err.Error())
			continue
		}
		logger.Infof("getIncMongoData, Sql: %+v, db: %s, collection: %s, result num: %d", queryimDbMap, meta.config.Db, dbconf.Collection, num)
		dbLoaderInfo.updateByGetAll = false // 用于标记从哪进入到updateFunc的
		maxUpdate, err := dbLoaderInfo.updateFunc(res, num, &dbLoaderInfo)
		if err != nil {
			logger.Errorf(" getIncMongoData inc update failed, db: %s, collection: %s, err: %s", meta.config.Db, dbconf.Collection, err.Error())
			continue
		}

		if num > 0 && maxUpdate > 0 {
			DbLoaderRegistryTime[dbconf.Index] = maxUpdate
		}
		logger.Infof("extractor getIncMongoData db: %s, collection: %s, num: %d, mongo time: %d, use_time_ms: %d", meta.config.Db, dbconf.Collection, num, DbLoaderRegistryTime[dbconf.Index], time.Now().UnixNano()/1e6-startTime)
	}
	logger.Info("getIncMongoData end")
}

func getSession() (*mgo.Session, error) {
	mongoAddr := meta.config.Mongo

	if meta.UseMongConsul {
		node := meta.Resolver.DiscoverNode()
		if node != nil && len(node.Address) > 0 {
			mongoAddr = node.Address
		}
	}

	dialInfo := &mgo.DialInfo{
		Addrs:     strings.Split(mongoAddr, ","),
		Direct:    true,
		Timeout:   time.Duration(meta.config.TimeOut) * time.Second,
		PoolLimit: 10,
	}
	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		logger.Errorf("connect mongo %s failed. err: %s", mongoAddr, err.Error())
		return nil, errors.New("connect mongo failed")
	}
	session.SetSocketTimeout(time.Duration(meta.config.TimeOut) * time.Second)
	session.SetSyncTimeout(time.Duration(meta.config.TimeOut) * time.Second)
	session.SetMode(mgo.Monotonic, true)
	logger.Infof("mongo addrs: %s", mongoAddr)
	return session, nil
}

func connect(curDb, curCollection string) (err error) {
	session, err := getSession()
	if err != nil {
		return err
	}

	if meta.mgoSession != nil {
		meta.mgoSession.Close()
	}
	meta.mgoSession = session
	db := session.DB(curDb)
	meta.mgoCollection = db.C(curCollection)
	if meta.mgoSession == nil || meta.mgoCollection == nil {
		return errors.New("mongo connection is nil, " + curCollection)
	}
	return nil
}

func getNewAdnDatabase() (*mgo.Database, error) {
	session, err := getSession()
	if err != nil {
		return nil, err
	}

	return session.DB("new_adn"), nil
}

// func init() {
// 	meta.stopChan = make(chan bool)
// 	NewMDbLoaderRegistry()
// }

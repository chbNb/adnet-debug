package extractor

import (
	"errors"

	"gitlab.mobvista.com/ADN/chasm/module/demand"

	jsoniter "github.com/json-iterator/go"
	smodel "gitlab.mobvista.com/ADN/structs/model"
	"gitlab.mobvista.com/ADN/treasure_box_sdk/tb_tools"
)

//func initMVConfigObj() {
//	mvConfigObj = new(MVConfigObj)
//}

func ConfigPreProc(i interface{}) error {
	ptr, ok := i.(*tb_tools.ExtractorInterface)
	if !ok {
		return errors.New("ConfigPreProc failed: type cast to *tb_tools.ExtractorInterface")
	}
	if mvConfig, ok := (*ptr).(*smodel.MVConfig); ok {
		value, err := configUpdateFunc(mvConfig)
		if err != nil {
			logger.Errorf("ConfigPreProc error: %v, key: %v", err.Error(), mvConfig.Key)
		}
		if value == nil {
			*ptr = nil
		} else {
			mvConfig.Value = value
			*ptr = mvConfig
		}
		return err
	} else {
		return errors.New("ConfigPreProc failed: type cast to *smodel.MVConfig failed")
	}
}

func configUpdateFunc(mvConfig *smodel.MVConfig) (interface{}, error) {
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	jsonStr, _ := json.Marshal(mvConfig.Value)
	switch mvConfig.Key {
	case "MAS_ABTEST":
		value := getMasAbtest(jsonStr)
		return value, nil
	case "SUPPORT_VIDEO_BANNER":
		value := getSupportVideoBanner(jsonStr)
		return value, nil
	case "ONLINE_API_SUPPORT_DEEPLINK_V2":
		value := getONLINE_API_SUPPORT_DEEPLINK_V2(jsonStr)
		return value, nil
	//case "GET_CAMPAIGN_USE_ZEUS_REDIS_DECODE":
	//	value, _ := getDecoodeSwith(jsonStr)
	//	return value, nil
	case "OPENAPI_SCENARIO":
		value, _ := getOpenapiScenario(jsonStr)
		return value, nil
	case "LOW_GRADE_APP_ID":
		value, _ := getConfigLowGradeApp(jsonStr)
		return value, nil
	case "REPLENISH_APP":
		value, _ := getReplenishApp(jsonStr)
		return value, nil
	case "ADV_BLACK_SUBID_LIST_V2":
		value, _ := getAdvBlackSubIdList(jsonStr)
		return value, nil
	case "APP_POST_LIST_V2":
		value, _ := getAppPostList(jsonStr)
		return value, nil
	case "offerwall_guidelines":
		value, _ := getOfferwallGuidelines(jsonStr)
		return value, nil
	case "ENDCARD_CONFIG":
		value, _ := getEndcard(jsonStr)
		return value, nil
	case "ALGO_TEST_CONFIG":
		value, _ := getAdServerTestConfig(jsonStr)
		return value, nil
	case "DEF_RV_TEMPLATE":
		value, _ := getDefRVTemplate(jsonStr)
		return value, nil
	case "VERSION_COMPARE":
		value, _ := getVersionCompare(jsonStr)
		return value, nil
	case "PLAYABLE_TEST":
		value, _ := getPlayableTest(jsonStr)
		return value, nil
	case "TRACK_URL_CONFIG_NEW":
		value, _ := getTRACK_URL_CONFIG_NEW(jsonStr)
		return value, nil
	case "3S_CHINA_DOMAIN":
		value, _ := get3S_CHINA_DOMAIN(jsonStr)
		return value, nil
	case "JUMP_TYPE_CONFIG":
		value, _ := getJUMP_TYPE_CONFIG(jsonStr)
		return value, nil
	case "JUMP_TYPE_CONFIG_IOS":
		value, _ := getJUMP_TYPE_CONFIG_IOS(jsonStr)
		return value, nil
	case "JUMPTYPE_SDKVERSION":
		value, _ := getJUMPTYPE_SDKVERSION(jsonStr)
		return value, nil
	case "ADSTACKING":
		value, _ := getADSTACKING(jsonStr)
		return value, nil
	case "SETTING_CONFIG":
		value, _ := getSETTING_CONFIG(jsonStr)
		return value, nil
	case "template":
		value, _ := getTemplate(jsonStr)
		return value, nil
	case "offerwall_urls":
		value, _ := getofferwall_urls(jsonStr)
		return value, nil
	case "SIGN_NO_CHECK_APPS":
		value, _ := getSIGN_NO_CHECK_APPS(jsonStr)
		return value, nil
	case "C_TOI":
		value, _ := getCToi(jsonStr)
		return value, nil
	case "MP_MAP_UNIT_V2":
		value, _ := getMP_MAP_UNIT(jsonStr)
		return value, nil
	case "CAN_CLICK_MODE_SIX_PUBLISHER":
		value, _ := getCAN_CLICK_MODE_SIX_PUBLISHER(jsonStr)
		return value, nil
	case "DEL_PRICE_FLOOR_UNIT":
		value, _ := getDEL_PRICE_FLOOR_UNIT(jsonStr)
		return value, nil
	case "ONLINE_PRICE_FLOOR_APPID":
		value, _ := getONLINE_PRICE_FLOOR_APPID(jsonStr)
		return value, nil
	case "CREATIVE_CHECK_MIGU_ADX_CREATIVE_IDS":
		value, _ := getCREATIVE_CHECK_MIGU_ADX_CREATIVE_IDS(jsonStr)
		return value, nil
	case "CREATIVE_CHECK_HUPU_ADX_CREATIVE_IDS":
		value, _ := getCREATIVE_CHECK_HUPU_ADX_CREATIVE_IDS(jsonStr)
		return value, nil
	case "IA_APPWALL":
		value, _ := getIA_APPWALL(jsonStr)
		return value, nil
	//case "GDT_KEY_WHITE_LIST":
	//	value, _ := getGdtKeyWhiteList(jsonStr)
	//	return value, nil
	//case "TOUTIAO_ALGO":
	//	value, _ := getTOUTIAO_ALGO(jsonStr)
	//	return value, nil
	case "TO_NEW_CDN_APPS":
		value, _ := getTO_NEW_CDN_APPS(jsonStr)
		return value, nil
	case "UA_AB_TEST_CONFIG":
		value, _ := getUA_AB_TEST_CONFIG(jsonStr)
		return value, nil
	//case "GO_TRACK":
	//	value, _ := getGO_TRACK(jsonStr)
	//	return value, nil
	case "JUMP_TYPE_CONFIG_ADV":
		value, _ := getJUMP_TYPE_CONFIG_ADV(jsonStr)
		return value, nil
	case "FUN_MV_MAP":
		value, _ := getFUN_MV_MAP(jsonStr)
		return value, nil
	case "NO_CHECK_PARAM_APP":
		value, _ := getNO_CHECK_PARAM_APP(jsonStr)
		return value, nil
	case "HAS_EXTENSIONS_UNIT":
		value, _ := getHAS_EXTENSIONS_UNIT(jsonStr)
		return value, nil
	case "RETRUN_VAST_APP":
		value := getRETRUN_VAST_APP(jsonStr)
		return value, nil
	case "UA_AB_TEST_SDK_OS_CONFIG":
		value, _ := getUA_AB_TEST_SDK_OS_CONFIG(jsonStr)
		return value, nil
	case "GAMELOFT_CREATIVE_URLS":
		value, _ := getGAMELOFT_CREATIVE_URLS(jsonStr)
		return value, nil
	case "UA_AB_TEST_THIRD_PARTY_CONFIG":
		value, _ := getUA_AB_TEST_THIRD_PARTY_CONFIG(jsonStr)
		return value, nil
	case "HUPU_DEFAULT_UNITID":
		value, _ := getHUPU_DEFAULT_UNITID(jsonStr)
		return value, nil
	case "UA_AB_TEST_CAMPAIGN_CONFIG":
		value, _ := getUA_AB_TEST_CAMPAIGN_CONFIG(jsonStr)
		return value, nil
	case "HUPU_DEFAULT_PRICE":
		value, _ := getHUPU_DEFAULT_PRICE(jsonStr)
		return value, nil
	case "JUMP_TYPE_CONFIG_THIRD_PARTY":
		value, _ := getJUMP_TYPE_CONFIG_THIRD_PARTY(jsonStr)
		return value, nil
	case "ONLINE_PRICE_FLOOR":
		value, _ := getONLINE_PRICE_FLOOR(jsonStr)
		return value, nil
	case "LINKTYPE_UNITID":
		value, _ := getLINKTYPE_UNITID(jsonStr)
		return value, nil
	//case "PLAYABLE_TEST_UNITS":
	//	value, _ := getPLAYABLE_TEST_UNITS(jsonStr)
	//	return value, nil
	case "MAX_APPID_AND_UNITID":
		value, _ := getMAX_APPID_AND_UNITID(jsonStr)
		return value, nil
	//case "ABTEST_GAIDIDFA":
	//	value, _ := getABTEST_GAIDIDFA(jsonStr)
	//	return value, nil
	case "REPLACE_BRAND_MODEL":
		value, _ := getREPLACE_BRAND_MODEL(jsonStr)
		return value, nil
	//case "NEW_PLAYABLE_SWITCH":
	//	value, _ := getNEW_PLAYABLE_SWITCH(jsonStr)
	//	return value, nil
	//case "PLAYABLE_ABTEST_RATE":
	//	value, _ := getPLAYABLE_ABTEST_RATE(jsonStr)
	//	return value, nil
	case "PUB_CC_CDN":
		value, _ := getPUB_CC_CDN(jsonStr)
		return value, nil
	case "TP_WITHOUT_RV":
		value, _ := getTP_WITHOUT_RV(jsonStr)
		return value, nil
	case "PPTV_APPID_AND_UNITID":
		value, _ := getPPTV_APPID_AND_UNITID(jsonStr)
		return value, nil
	case "IFENG_APPID_AND_UNITID":
		value, _ := getIFENG_APPID_AND_UNITID(jsonStr)
		return value, nil
	case "REQUEST_BLACKLIST":
		value, _ := getREQUEST_BLACKLIST(jsonStr)
		return value, nil
	case "PPTV_DEAL_ID":
		value, _ := getPPTV_DEAL_ID(jsonStr)
		return value, nil
	case "CREATIVE_CHECK_PPTV_ADX_CREATIVE_IDS":
		value, _ := getCREATIVE_CHECK_PPTV_ADX_CREATIVE_IDS(jsonStr)
		return value, nil
	case "CM_APPID_AND_UNITID":
		value, _ := getCM_APPID_AND_UNITID(jsonStr)
		return value, nil
	case "FCA_SWITCH":
		value, _ := getFCA_SWITCH(jsonStr)
		return value, nil
	case "FCA_CAMIDS":
		value, _ := getFCA_CAMIDS(jsonStr)
		return value, nil
	case "NEW_DEFAULT_FCA":
		value, _ := getNEW_DEFAULT_FCA(jsonStr)
		return value, nil
	//case "CM_MP_APPIDS":
	//	value, _ := getCM_MP_APPIDS(jsonStr)
	//	return value, nil
	case "LOW_FLOW_UNITS":
		value, _ := getLOW_FLOW_UNITS(jsonStr)
		return value, nil
	case "SS_ABTEST_CAMPAIGN":
		value, _ := getSS_ABTEST_CAMPAIGN(jsonStr)
		return value, nil
	//case "NEW_CREATIVE_ABTEST_RATE":
	//	value, _ := getNEW_CREATIVE_ABTEST_RATE(jsonStr)
	//	return value, nil
	//case "NEW_CREATIVE_TEST_UNITS":
	//	value, _ := getNEW_CREATIVE_TEST_UNITS(jsonStr)
	//	return value, nil
	//case "CREATIVE_ABTEST":
	//	value, _ := getCREATIVE_ABTEST(jsonStr)
	//	return value, nil
	case "TEMPLATE_MAP":
		value, _ := getTEMPLATE_MAP(jsonStr)
		return value, nil
	case "NEW_CREATIVE_TIMEOUT":
		value, _ := getNEW_CREATIVE_TIMEOUT(jsonStr)
		return value, nil
	case "XM_NEW_RETURN_UNITS":
		value, _ := getXM_NEW_RETURN_UNITS(jsonStr)
		return value, nil
	case "WEBVIEW_PRISON_PUBIDS_APPIDS":
		value, _ := getWEBVIEW_PRISON_PUBIDS_APPIDS(jsonStr)
		return value, nil
	case "CONFIG_OFFER_PLCT":
		value, _ := getCONFIG_OFFER_PLCT(jsonStr)
		return value, nil
	case "CCT_ABTEST_CONF":
		value, _ := getCCT_ABTEST_CONF(jsonStr)
		return value, nil
	case "JSSDK_DOMAIN":
		value, _ := getJSSDK_DOMAIN(jsonStr)
		return value, nil
	case "CHET_URL_UNIT":
		value, _ := getCHET_URL_UNIT(jsonStr)
		return value, nil
	//case "MD5_FILE_PRISON_PUB":
	//	value, _ := getMD5_FILE_PRISON_PUB(jsonStr)
	//	return value, nil
	case "ANR_PRISION_BY_PUB_AND_OSV":
		value, _ := getANR_PRISION_BY_PUB_AND_OSV(jsonStr)
		return value, nil
	case "UNIT_WITHOUT_VIDEO_START":
		value, _ := getUNIT_WITHOUT_VIDEO_START(jsonStr)
		return value, nil
	case "REDUCE_FILL_SWITCH":
		value, _ := getREDUCE_FILL_SWITCH(jsonStr)
		return value, nil
	case "MORE_OFFER_CONF":
		value, _ := getMORE_OFFER_CONF(jsonStr)
		return value, nil
	case "NEW_MOF_IMP_RATE":
		value, _ := getNEW_MOF_IMP_RATE(jsonStr)
		return value, nil
	case "MOF_ABTEST_RATE":
		value, _ := getMOF_ABTEST_RATE(jsonStr)
		return value, nil
	case "adx_fake_dsp_price_default_factor":
		value, _ := getADX_FAKE_DSP_PRICE_DEFAULT_FACTOR(jsonStr)
		return value, nil
	case "ADNET_SWITCHS":
		value, _ := getADNET_SWITCHS(jsonStr)
		return value, nil
	case "BLACK_FOR_EXCLUDE_PACKAGE_NAME":
		value, _ := getBLACK_FOR_EXCLUDE_PACKAGE_NAME(jsonStr)
		return value, nil
	case "APPWALL_TO_MORE_OFFER_UNIT":
		value, _ := getAPPWALL_TO_MORE_OFFER_UNIT(jsonStr)
		return value, nil
	case "CLEAN_EMPTY_DEVICE_ABTEST":
		value, _ := getCleanEmptyDeviceABTest(jsonStr)
		return value, nil
	case "REQ_TYPE_AAB_TEST_CONFIG_NEW":
		value, _ := getREQ_TYPE_AAB_TEST_CONFIG(jsonStr)
		return value, nil
	case "ExcludeDisplayPackageABTest":
		value, _ := getExcludeDisplayPackageABTest(jsonStr)
		return value, nil
	case "REDUCE_FILL_FLINK_SWITCH":
		value, _ := getREDUCE_FILL_FLINK_SWITCH(jsonStr)
		return value, nil
	case "ONLINE_EMPTY_DEVICE_IPUA_ABTEST":
		value, _ := getOnlineEmptyDeviceIPUAABTest(jsonStr)
		return value, nil
	case "AUTO_LOAD_CACHE_ABTSET":
		value, _ := getAUTO_LOAD_CACHE_ABTSET(jsonStr)
		return value, nil
	case "CTN_SIZE_ABTEST":
		value, _ := getCTN_SIZE_ABTEST(jsonStr)
		return value, nil
	case "BLOCK_MORE_OFFER_CONFIG":
		value, _ := getBLOCK_MORE_OFFER_CONFIG(jsonStr)
		return value, nil
	case "CLOSE_BUTTON_AD_TEST_UNITS":
		value, _ := getCLOSE_BUTTON_AD_TEST_UNITS(jsonStr)
		return value, nil
	//case "MIGRATION_ABTEST":
	//	value, _ := getMIGRATION_ABTEST(jsonStr)
	//	return value, nil
	case "RETURN_PARAM_K_UNIT":
		value, _ := getRETURN_PARAM_K_UNIT(jsonStr)
		return value, nil
	case "ADD_UNIT_ID_ON_TRACKING_URL_UNIT":
		value, _ := getADD_UNIT_ID_ON_TRACKING_URL_UNIT(jsonStr)
		return value, nil
	case "CLICK_IN_ADTRACKING_THIRD_PARTY_CONFIG":
		value, _ := getCLICK_IN_ADTRACKING_THIRD_PARTY_CONFIG(jsonStr)
		return value, nil
	case "PsbReduplicateConfig":
		value := getPsbReduplicateConfig(jsonStr)
		return value, nil
	case "AopReduplicateConfig":
		value := getAopReduplicateConfig(jsonStr)
		return value, nil
	case "ALAC_PRISON_CONFIG":
		value, _ := getALAC_PRISON_CONFIG(jsonStr)
		return value, nil
	case "ADNET_COM_ABTEST_FIELDS":
		value, _ := getABTEST_FIELDS(jsonStr)
		return value, nil
	case "ADNET_COM_ABTEST_CONFS":
		value, _ := getABTEST_CONFS(jsonStr)
		return value, nil
	case "Adnet_Chet_Links":
		value, _ := getChetLinkConfigs(jsonStr)
		return value, nil
	case "Android_Low_Version_Config":
		value, _ := getAndroidLowVersionFilterCondition(jsonStr)
		return value, nil
	case "BANNER_HTML_STR":
		value, _ := getBANNER_HTML_STR(jsonStr)
		return value, nil
	case "STOREKIT_TIME_PRISON_CONF":
		value := getSTOREKIT_TIME_PRISON_CONF(jsonStr)
		return value, nil
	//case "IPUA_WHITE_LIST_RATE_CONF":
	//	value := getIPUA_WHITE_LIST_RATE_CONF(jsonStr)
	//	return value, nil
	case "IS_RETURN_PAUSEMOD":
		value := getIS_RETURN_PAUSEMOD(jsonStr)
		return value, nil
	case "CREATIVE_CHECK_MANGGUO_ADX_CREATIVE_IDS":
		value := getCREATIVE_CHECK_MANGGUO_ADX_CREATIVE_IDS(jsonStr)
		return value, nil
	case "MANGO_APPID_AND_UNITID":
		value := getMANGO_APPID_AND_UNITID(jsonStr)
		return value, nil
	case "ExcludeClickPackages":
		value := getExcludeClickPackages(jsonStr)
		return value, nil
	case "ADNET_CONF_LIST":
		value := getADNET_CONF_LIST(jsonStr)
		return value, nil
	case "ADCHOICE_CONFIG":
		value := getAdChoiceConfigData(jsonStr)
		return value, nil
	case "FILLRATE_ECPM_FLOOR_SWITCH":
		value := getFillRateEcpmFloorSwitch(jsonStr)
		return value, nil
	case "MOATTAG_CONFIG":
		value := getMoattagConfig(jsonStr)
		return value, nil
	case "BAD_REQUEST_FILTER_CONF":
		value := getBAD_REQUEST_FILTER_CONF(jsonStr)
		return value, nil
	case "BIG_TEMPLATE_CONF":
		value := getBIG_TEMPLATE_CONF(jsonStr)
		return value, nil
	case "CN_TRACKING_DOMAIN_CONF":
		value := getCNTrackingDomainConf(jsonStr)
		return value, nil
	case "NoticeClickUniq":
		value := getNoticeClickUniq(jsonStr)
		return value, nil
	case "ExcludeImpressionPackages":
		value := getExcludeImpressionPackages(jsonStr)
		return value, nil
	case "ExcludeImpressionPackagesV2":
		value := getExcludeImpressionPackagesV2(jsonStr)
		return value, nil
	case "PassthroughData":
		value := getPassthroughData(jsonStr)
		return value, nil
	case "PolarisFlagConf":
		value := getPolarisFlagConf(jsonStr)
		return value, nil
	//case "ADNET_LANG_CREATIVE_ABTEST_CONF":
	//	value := getAdnetLangCreativeABTestConf(jsonStr)
	//	return value, nil
	case "FREQ_CONTROL_CONFIG":
		value := getFREQ_CONTROL_CONFIG(jsonStr)
		return value, nil
	case "TIMEZONE_CONFIG":
		value := getTIMEZONE_CONFIG(jsonStr)
		return value, nil
	case "COUNTRY_CODE_TIMEZONE_CONFIG":
		value := getCOUNTRY_CODE_TIMEZONE_CONFIG(jsonStr)
		return value, nil
	case "ADNET_DSP_TPL_ABTEST_CONF":
		value := getDSP_TPL_ABTEST_CONF(jsonStr)
		return value, nil
	case "BANNER_TO_AS_ABTEST_CONF":
		value := getBANNER_TO_AS_ABTEST_CONF(jsonStr)
		return value, nil
	case "ONLY_REQUEST_THIRD_DSP_SWITCH":
		value := getONLY_REQUEST_THIRD_DSP_SWITCH(jsonStr)
		return value, nil
	case "ADNET_DCO_TEST_CONF":
		value := getDcoTestConf(jsonStr)
		return value, nil
	case "USE_PLACEMENT_IMP_CAP_SWITCH":
		value := getUSE_PLACEMENT_IMP_CAP_SWITCH(jsonStr)
		return value, nil
	case "LOW_FLOW_ADTYPE":
		value, _ := getLOW_FLOW_ADTYPE(jsonStr)
		return value, nil
	case "TaoBaoOfferID":
		value := getTaoBaoOfferID(jsonStr)
		return value, nil
	case "ADN_DEFAULT_VALUE":
		value := getADNET_DEFAULT_VALUE(jsonStr)
		return value, nil
	case "ADNET_CREATIVE_COMPRESS_ABTEST_CONF_V3":
		value := getADNET_CREATIVE_COMPRESS_ABTEST_CONF_V3(jsonStr)
		return value, nil
	case "CREATIVE_VIDEO_COMPRESS_ABTEST":
		value := getCREATIVE_VIDEO_COMPRESS_ABTEST(jsonStr)
		return value, nil
	case "CREATIVE_IMG_COMPRESS_ABTEST":
		value := getCREATIVE_IMG_COMPRESS_ABTEST(jsonStr)
		return value, nil
	case "CREATIVE_ICON_COMPRESS_ABTEST":
		value := getCREATIVE_ICON_COMPRESS_ABTEST(jsonStr)
		return value, nil
	case "BID_SDK_VERSION_CONFIG":
		value := getHBBidSdkVersionConfig(jsonStr)
		return value, nil
	case "EXCHANGE_RATE":
		value := getHBExchangeRate(jsonStr)
		return value, nil
	case "PUBLISHER_CURRENCY":
		value := getHBPublisherCurrency(jsonStr)
		return value, nil
	//case "HEADER_BIDDING_BLACKLIST":
	//	value := getHBBlacklist(jsonStr)
	//	return value, nil
	case "HBPubExcludeAppCheck":
		value := getHBPubExcludeAppCheck(jsonStr)
		return value, nil
	//case "HBBidFloorPubWhiteList":
	//	value := getHBBidFloorPubWhiteList(jsonStr)
	//	return value, nil
	case "HBAdxEndpoint":
		value := getHBAdxEndpoint(jsonStr)
		return value, nil
	case "HBAdxEndpointV2":
		value := getHBAdxEndpointV2(jsonStr)
		return value, nil
	case "HBCDNDomainABTestPubs":
		value := getHBCDNDomainABTestPubs(jsonStr)
		return value, nil
	//case "AerospikeUseConsul":
	//	value := getAerospikeUseConsul(jsonStr)
	//	return value, nil
	//case "UseConsulServices":
	//	value := getUseConsulServices(jsonStr)
	//	return value, nil
	case "UseConsulServicesV2Ratio":
		value := getUseConsulServicesV2Ratio(jsonStr)
		return value, nil
	case "HBAerospikeStorageConf":
		value := getHBAerospikeStorageConfArr(jsonStr)
		return value, nil
	case "HBRequestBidServerConf":
		value := getHBRequestBidServerConfArr(jsonStr)
		return value, nil
	case "UseConsulServicesV2":
		value := getUseConsulServicesV2(jsonStr)
		return value, nil
	case "TreasureBoxIDCount":
		value := getTreasureBoxIDCount(jsonStr)
		return value, nil
	case "CountryBlackPackageListConf":
		value := getCountryBlackPackageListConf(jsonStr)
		return value, nil
	case "OnlineApiPubBidPriceConf":
		value := getOnlineApiPubBidPriceConf(jsonStr)
		return value, nil
	case "SupportTrackingTemplateConf":
		value := getSupportTrackingTemplateConf(jsonStr)
		return value, nil
	case "YLHClickModeTestConfig":
		value := getYLHClickModeTestConfig(jsonStr)
		return value, nil
	case "V5_ABTEST_CONFIG":
		value := getV5_ABTEST_CONFIG(jsonStr)
		return value, nil
	case "REPLACE_TRACKING_DOMAIN_CONF":
		value := getREPLACE_TRACKING_DOMAIN_CONF(jsonStr)
		return value, nil
	case "NEW_AFREECATV_UNIT":
		value := getNEW_AFREECATV_UNIT(jsonStr)
		return value, nil
	case "SHAREIT_AF_WHITE_LIST_CONF":
		value := getSHAREIT_AF_WHITE_LIST_CONF(jsonStr)
		return value, nil
	case "ONLINE_API_USE_ADX":
		value := getONLINE_API_USE_ADX(jsonStr)
		return value, nil
	case "ONLINE_API_USE_ADX_MAS":
		value := getONLINE_API_USE_ADX_MAS(jsonStr)
		return value, nil
	case "ONLINE_API_MAX_BID_PRICE":
		value := getONLINE_API_MAX_BID_PRICE(jsonStr)
		return value, nil
	case "NEW_HUPU_UNITID_MAP":
		value := getNEW_HUPU_UNITID_MAP(jsonStr)
		return value, nil
	case "NEW_CDN_TEST":
		value := getNEW_CDN_TEST(jsonStr)
		return value, nil
	case "ONLINE_FREQUENCY_CONTROL_CONF":
		value := getONLINE_FREQUENCY_CONTROL_CONF(jsonStr)
		return value, nil
	case "CHANGE_ONLINE_DEEPLINK_WAY_PACKAGE_LIST":
		value := getCHANGE_ONLINE_DEEPLINK_WAY_PACKAGE_LIST(jsonStr)
		return value, nil
	case "DEBUG_BID_FLOOR_AND_BID_PRICE_CONF":
		value := getDEBUG_BID_FLOOR_AND_BID_PRICE_CONF(jsonStr)
		return value, nil
	case "REPLACE_TEMPLATE_URL_CONF":
		value := getREPLACE_TEMPLATE_URL_CONF(jsonStr)
		return value, nil
	case "TPL_CREATIVE_DOMAIN_CONF":
		value := getTPL_CREATIVE_DOMAIN_CONF(jsonStr)
		return value, nil
	case "RETURN_WTICK_CONF":
		value := getRETURN_WTICK_CONF(jsonStr)
		return value, nil
	case "EXCLUDE_PACKAGES_BY_CITYCODE_CONF":
		value := getEXCLUDE_PACKAGES_BY_CITYCODE_CONF(jsonStr)
		return value, nil
	case "H265_VIDEO_ABTEST_CONF":
		value := getH265_VIDEO_ABTEST_CONF(jsonStr)
		return value, nil
	case "ADNET_PARAMS_ABTEST_CONFS":
		value := getADNET_PARAMS_ABTEST_CONFS(jsonStr)
		return value, nil
	case "MAPPING_SERVER_RATE_CONF":
		value := getMAPPING_SERVER_RATE_CONF(jsonStr)
		return value, nil
	case "HB_API_CONFS":
		value := getHB_API_CONFS(jsonStr)
		return value, nil
	case "MEDIATION_CHANNEL_ID":
		value := getMEDIATION_CHANNEL_ID(jsonStr)
		return value, nil
	case "NEW_JUMP_TYPE_CONFIG_THIRD_PARTY":
		value := getNEW_JUMP_TYPE_CONFIG_THIRD_PARTY(jsonStr)
		return value, nil
	case "AD_PACKAGE_NAME_REPLACE_CONF":
		value := getAdPackageNameReplace(jsonStr)
		return value, nil
	case "FILTER_BY_STACK_CONF":
		value := getFilterByStackConf(jsonStr)
		return value, nil
	case "HB_REQUEST_SNAPSHOT":
		value := getHBRequestSnapshot(jsonStr)
		return value, nil
	case "MAPPING_IDFA_CONF":
		value := getMappingIdfaConf(jsonStr)
		return value, nil
	case "ORIENTATION_POISON":
		value := getOrientationPoisonConf(jsonStr)
		return value, nil
	case "TRACKING_CN_ABTEST_CONF":
		value := getTrackingCNABTestConf(jsonStr)
		return value, nil
	case "MAPPING_IDFA_COVER_IDFA_ABTEST_CONF":
		value := getMappingIdfaCoverIdfaABTestConf(jsonStr)
		return value, nil
	case "SupportSmartVBAConfig":
		value := getSupportSmartVBAConfig(jsonStr)
		return value, nil
	case "MOREOFFER_AND_APPWALL_MOVE_TO_PIONEER_ABTEST_CONF":
		value := getMoreOfferAndAppwallMoveToPioneerABTestConf(jsonStr)
		return value, nil
	case "MORE_OFFER_REQUEST_DOMAIN":
		value := getMORE_OFFER_REQUEST_DOMAIN(jsonStr)
		return value, nil
	case "ONLINE_PUBLISHER_ADNUM_CONF":
		value := getOnlinePublisherAdNumConfig(jsonStr)
		return value, nil
	case "TRUE_NUM_BY_AD_TYPE":
		value := getTRUE_NUM_BY_AD_TYPE(jsonStr)
		return value, nil
	case "MEDIATION_NOTICE_URL_MACRO_CONF":
		value := getMediationNoticeURLMacroConfig(jsonStr)
		return value, nil
	case "HB_CACHED_ADNUM_CONF":
		value := getHBCachedAdNumConfig(jsonStr)
		return value, nil
	case "HB_OFFER_BID_PRICE_ABTEST_CONF":
		value := getHBOfferBidPriceABTestConf(jsonStr)
		return value, nil
	case "TRACK_DOMAIN_BY_COUNTRY_CODE_CONF":
		value := getTrackDomainByCountryCodeConf(jsonStr)
		return value, nil
	case "ServiceDegradeRate":
		value := getServiceDegradeRate(jsonStr)
		return value, nil
	case "HB_LOAD_DOMAIN_BY_COUNTRY_CODE_CONF":
		value := getHBLoadDomainByCountryCodeConf(jsonStr)
		return value, nil
	case "MORE_OFFER_REQUEST_DOMAIN_BY_COUNTRY_CODE_CONF":
		value := getMoreOfferRequestDomainByCountryCodeConf(jsonStr)
		return value, nil
	case "LOAD_DOMAIN_ABTEST":
		value := getLoadDomainABTest(jsonStr)
		return value, nil
	case "HB_AEROSPIKE_CONF":
		value := getHBAerospikeConf(jsonStr)
		return value, nil
	case "VAST_BANNER_DSP":
		return getVastBannerDsp(jsonStr)
	case "RESIDUAL_MEMORY_FILTER_CONF":
		value := getResidualMemoryFilterConf(jsonStr)
		return value, nil
	case "TMAX_ABTEST_CONF":
		value := getTmaxABTestConf(jsonStr)
		return value, nil
	case "MP_TO_PIONEER_ABTEST_CONF":
		value := getMpToPioneerABTestConf(jsonStr)
		return value, nil
	case "HB_EVENT_HTTP_PROTOCOL":
		value := getHBEeventHTTPProtocolConf(jsonStr)
		return value, nil
	case "HB_V5_ABTEST_CONF":
		value := getHBV5ABTestConf(jsonStr)
		return value, nil
	case "FILTER_ADSERVER_REQUEST_CONF":
		value := getFilterAdserverRequestConf(jsonStr)
		return value, nil
	case "HB_CHECK_DEVICE_GEO_CC":
		value := getHBCheckDeviceGEOCCConf(jsonStr)
		return value, nil
	case "TEMPLATE_MAP_V2":
		value := getTEMPLATEMAPV2(jsonStr)
		return value, nil
	case "HBLoadFilterConfigs":
		value := getHBLoadFilterConfigs(jsonStr)
		return value, nil
	case "FILTER_AUTO_CLICK_CONF":
		value := getFilterAutoClickConf(jsonStr)
		return value, nil
	case "TRACKING_CDN_DOMAIN_MAP":
		value := getTrackingCdnDomainMap(jsonStr)
		return value, nil
	case "CDN_TRACKING_DOMAIN_ABTEST_CONF":
		value := getCdnTrackingDomainABTestConf(jsonStr)
		return value, nil
	default:
		if parserFunc, err := demand.GetGlobalConfigParserByBytes(mvConfig.Key); err == nil {
			value, _ := parserFunc(jsonStr)
			return value, nil
		}
	}
	return nil, nil
}

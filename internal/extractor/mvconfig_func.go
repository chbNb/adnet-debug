package extractor

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/easierway/concurrent_map"
	jsoniter "github.com/json-iterator/go"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	smodel "gitlab.mobvista.com/ADN/structs/model"
	"gitlab.mobvista.com/ADN/treasure_box_sdk/tb_tools"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//type MVConfigObj struct {
//	OpenapiScenaria                         []string
//	LowGradeAppId                           int64
//	ReplenishApp                            []int64
//	AdvBlackSubIdListV2                     map[string]map[string]map[string]string
//	AppPostListV2                           map[int32]*smodel.AppPostList
//	OfferwallGuidelines                     map[int32]string
//	EndCardConfig                           smodel.EndCard
//	DefRvTemp                               *smodel.VideoTemplateUrlItem
//	VersionCompare                          map[string]mvutil.VersionCompare
//	PlayableTest                            map[int64]mvutil.PlatableTest
//	TrackUrlConfigNew                       map[int32]mvutil.TRACK_URL_CONFIG_NEW
//	S3TrackingDomain                        mvutil.CONFIG_3S_CHINA_DOMAIN
//	JumpTypeConfig                          map[string]int32
//	JumpTypeIOS                             map[string]int32
//	JumpTypeSDKVersion                      map[string]map[string]string
//	AdStacking                              mvutil.ADSTACKING
//	SettingConfig                           mvutil.SETTING_CONFIG
//	Template                                map[string]mvutil.Template
//	OfferWallUrls                           mvutil.COfferwallUrls
//	NoCheckApps                             []int64
//	AdServerTestConfig                      mvutil.AdServerTestConfig
//	CToi                                    int32
//	MP_MAP_UNIT                             map[string]mvutil.MP_MAP_UNIT_
//	DEL_PRICE_FLOOR_UNITS                   []int64
//	CAN_CLICK_MODE_SIX_PUBLISHER            []int64
//	IA_APPWALL                              string
//	TOUTIAO_ALGO                            map[int64][]int64
//	TO_NEW_CDN_APPS                         []int64
//	UA_AB_TEST_CONFIG                       map[string]int
//	DecodeSwith                             map[string]int
//	GO_TRACK                                map[string]map[string]int
//	JUMP_TYPE_CONFIG_ADV                    map[string]map[string]map[string]int32
//	FUN_MV_MAP                              map[string]map[string]string
//	NO_CHECK_PARAM_APP                      mvutil.NO_CHECK_PARAM_APP
//	HAS_EXTENSIONS_UNIT                     []int64
//	GAMELOFT_CREATIVE_URLS                  map[string]map[string]string
//	HUPU_DEFAULT_UNITID                     map[string]int64
//	HUPU_DEFAULT_PRICE                      int
//	UA_AB_TEST_SDK_OS_CONFIG                map[string]string
//	UA_AB_TEST_THIRD_PARTY_CONFIG           []string
//	UA_AB_TEST_CAMPAIGN_CONFIG              []int64
//	ONLINE_PRICE_FLOOD_APPID                map[string]float64
//	CREATIVE_CHECK_MIGU_ADX_CREATIVE_IDS    map[string]string
//	ONLINE_PRICE_FLOOR                      map[string]float64
//	MAX_APPID_AND_UNITID                    map[string]mvutil.ONLINE_APPID_UNITID
//	JUMP_TYPE_CONFIG_THIRD_PARTY            map[string]map[string]map[string]int32
//	LINKTYPE_UNITID                         map[string]int
//	PLAYABLE_TEST_UNITS                     []int64
//	ABTEST_GAIDIDFA                         map[string][]string
//	REPLACE_BRAND_MODEL                     map[string]map[string]string
//	CREATIVE_CHECK_HUPU_ADX_CREATIVE_IDS    []int64
//	NEW_PLAYABLE_SWITCH                     bool
//	PLAYABLE_ABTEST_RATE                    int
//	PUB_CC_CDN                              map[string]map[string]string
//	TP_WITHOUT_RV                           []int64
//	PPTV_APPID_AND_UNITID                   map[string]mvutil.ONLINE_APPID_UNITID
//	REQUEST_BLACKLIST                       mvutil.REQUEST_BLACKLIST
//	IFENG_APPID_AND_UNITID                  map[string]mvutil.ONLINE_APPID_UNITID
//	PPTV_DEAL_ID                            string
//	CREATIVE_CHECK_PPTV_ADX_CREATIVE_IDS    map[string]string
//	CM_APPID_AND_UNITID                     map[string]mvutil.ONLINE_APPID_UNITID
//	FCA_SWITCH                              bool
//	FCA_CAMIDS                              []int64
//	NEW_DEFAULT_FCA                         int
//	CM_MP_APPIDS                            []int64
//	LOW_FLOW_UNITS                          map[string]int
//	NEW_CREATIVE_ABTEST_RATE                int
//	NEW_CREATIVE_TEST_UNITS                 map[string]int32
//	CREATIVE_ABTEST                         map[string]int
//	TEMPLATE_MAP                            map[string]map[string]string
//	BACKEND_BT                              map[int64]map[int64][]int64
//	SS_ABTEST_CAMPAIGN                      map[string]*mvutil.SSTestTracfficRule
//	NEW_CREATIVE_TIMEOUT                    int
//	XM_NEW_RETURN_UNITS                     []int64
//	WEBVIEW_PRISON_PUBIDS_APPIDS            map[string][]int64
//	CONFIG_OFFER_PLCT                       map[string]map[string]mvutil.OFFER_PLCT
//	CCT_ABTEST_CONF                         map[string]map[string]int
//	JSSDK_DOMAIN                            map[string]string
//	ADSOURCE_PACKAGE_BLACKLIST              map[string]map[string][]string //adsource包黑名单
//	CHET_URL_UNIT                           []int64
//	MD5_FILE_PRISON_PUB                     []int64
//	ANR_PRISION_BY_PUB_AND_OSV              mvutil.PUB_ADN_OSV
//	UNIT_WITHOUT_VIDEO_START                []int64
//	REDUCE_FILL_SWITCH                      mvutil.REDUCE_FILL_SWITCH // 降填充开关
//	MORE_OFFER_CONF                         mvutil.MORE_OFFER_CONFIG
//	NEW_MOF_IMP_RATE                        int
//	MOF_ABTEST_RATE                         int
//	GDT_KEY_WHITELIST                       []string
//	ADX_FAKE_DSP_PRICE_DEFAULT_FACTOR       map[string]float64
//	ADNET_SWITCHS                           map[string]int
//	BLACK_FOR_EXCLUDE_PACKAGE_NAME          mvutil.BLACK_FOR_EXCLUDE_PACKAGE_NAME
//	APPWALL_TO_MORE_OFFER_UNIT              map[string]int32
//	REQ_TYPE_AAB_TEST_CONFIG                mvutil.REQ_TYPE_AAB_TEST_CONFIG
//	CleanEmptyDeviceABTest                  mvutil.CleanEmptyDeviceABTest
//	REDUCE_FILL_FLINK_SWITCH                mvutil.REDUCE_FILL_FLINK_SWITCH
//	OnlineEmptyDeviceNoServerJump           mvutil.OnlineEmptyDeviceNoServerJump
//	ExcludeDisplayPackageABTest             mvutil.ExcludeDisplayPackageABTest
//	AUTO_LOAD_CACHE_ABTSET_CONFIG           mvutil.AUTO_LOAD_CACHE_ABTSET_CONFIG
//	OnlineEmptyDeviceIPUAABTest             mvutil.OnlineEmptyDeviceIPUA
//	CTN_SIZE_ABTEST                         map[string]map[string]int32
//	BLOCK_MORE_OFFER_CONFIG                 mvutil.ENDCARD_BLOCK_MORE_OFFER_CONFIG
//	CLOSE_BUTTON_AD_TEST_UNITS              map[string]int32
//	CLICK_IN_ADTRACKING_THIRD_PARTY_CONFIG  map[string]int
//	MIGRATION_ABTEST                        map[string]int32
//	RETURN_PARAM_K_UNIT                     []int64
//	ADD_UNIT_ID_ON_TRACKING_URL_UNIT        []int64
//	PsbReduplicateConfig                    mvutil.PsbReduplicateConfig
//	ALAC_PRISON_CONFIG                      map[string][]int
//	ABTEST_FIELDS                           mvutil.ABTEST_FIELDS
//	ABTEST_CONFS                            map[string][]mvutil.ABTEST_CONF
//	ChetLinkConfig                          []*mvutil.ChetLinkConfigItem
//	AndroidLowVersionFilterCondition        *mvutil.AndroidLowVersionFilterCondition
//	BANNER_HTML_STR                         string
//	ThirdPartyWhiteList                     mvutil.ThirdPartyWhiteList
//	STOREKIT_TIME_PRISON_CONF               mvutil.StorekitTimePrisonConf
//	IPUA_WHITE_LIST_RATE_CONF               map[string]map[string]int
//	CREATIVE_CHECK_MANGGUO_ADX_CREATIVE_IDS map[string]string
//	MANGO_APPID_AND_UNITID                  map[string]map[string]mvutil.ONLINE_APPID_UNITID
//	IS_RETURN_PAUSEMOD                      mvutil.IsReturnPauseModConf
//	ExcludeClickPackages                    *mvutil.ExcludeClickPackages
//	ADNET_CONF_LIST                         map[string][]int64
//	AdChoiceConfig                          *mvutil.AdChoiceConfigData
//	FillRateEcpmFloorSwitch                 bool
//	CLICK_IN_SERVER_CONF_NEW                mvutil.ClickInServerConf
//	MoattagConfig                           *mvutil.MoattagConfigData
//	BAD_REQUEST_FILTER_CONF                 []*mvutil.BadRequestFilterConf
//	BIG_TEMPLATE_CONF                       *mvutil.TestConf
//	CN_TRACKING_DOMAIN_CONF                 *mvutil.CNTrackingDomainTestConf
//	NoticeClickUniq                         *mvutil.NoticeClickUniq
//	//AdjustPostbackABTest                    *mvutil.AdjustPostbackABTest
//	ExcludeImpressionPackages               *mvutil.ExcludeClickPackages
//	PolarisFlagConf                         *mvutil.TestConf
//	TemplateMapV2                           map[string]map[string][]mvutil.TemplateMap
//	BigTemplateMap                          map[string][]mvutil.TemplateMap
//	PassthroughData                         []string
//	AdnetLangCreativeABTestConf             map[string][]map[string]int64
//	FreqControlConf                         *mvutil.FreqControlConfig
//	TimezonConfig                           map[string]int
//	CountryCodeTimezoneConfig               map[string]int
//	DspTplAbtestConf                        map[string][]*mvutil.DspTplAbtest
//	BannerToAsABTestConf                    map[string]int
//	ONLY_REQUEST_THIRD_DSP_SWITCH           bool
//	DcoTestConf                             *mvutil.TestConf
//	USE_PLACEMENT_IMP_CAP_SWITCH            bool
//	LOW_FLOW_ADTYPE                         map[string]float64
//	TaoBaoOfferID                           []int64
//	ADNET_DEFAULT_VALUE                     map[string]string
//	CREATIVE_COMPRESS_ABTEST_CONF_V3        map[string]*mvutil.CreativeCompressABTestV3Data
//	CREATIVE_VIDEO_COMPRESS_ABTEST          map[string]int32
//	CREATIVE_IMG_COMPRESS_ABTEST            map[string]int32
//	CREATIVE_ICON_COMPRESS_ABTEST           map[string]int32
//	HBBidSdkVersionConfig                   map[string]map[string]string
//	HBExchangeRate                          map[string]float64
//	HBPublisherCurrency                     map[string]string
//	HBBlacklist                             map[string][]string
//	HBPubExcludeAppCheck                    []int64
//	HBBidFloorPubWhiteList                  []int64
//	HBAdxEndpoint                           map[string]mvutil.HBAdxEndpointMkvData
//	HBAdxEndpointV2                         map[string]map[string]mvutil.HBAdxEndpointMkvData
//	HBCDNDomainABTestPubs                   map[string]string
//	TreasureBoxIDCount                      map[string]bool
//	CountryBlackPackageListConf             map[string]map[string][]string
//	OnlineApiPubBidPriceConf                *mvutil.OnlineApiPubBidPriceConf
//	SupportTrackingTemplateConf             *mvutil.SupportTrackingTemplateConf
//	YLHClickModeTestConfig                  map[string]int
//	AerospikeUseConsul                      bool
//	UseConsulServices                       map[string]bool
//	UseConsulServicesV2                     map[string]map[string]map[string]bool
//	RETRUN_VAST_APP                         []int64
//	V5AbtestConf                            *mvutil.V5AbtestConf
//	REPLACE_TRACKING_DOMAIN_CONF            *mvutil.ReplaceTrackingDomainConf
//	NEW_AFREECATV_UNIT                      []int64
//	SHAREIT_AF_WHITE_LIST_CONF              *mvutil.ShareitAfWhiteListConf
//	ONLINE_API_USE_ADX                      map[string]map[int64]int
//	ONLINE_API_USE_ADX_MAS                  map[string]map[int64]int
//	ONLINE_API_SUPPORT_DEEPLINK_V2          map[string]map[int64]int32
//	ONLINE_API_MAX_BID_PRICE                *mvutil.OnlineApiMaxBidPrice
//	NEW_HUPU_UNITID_MAP                     map[string]map[string]int64
//	NEW_CDN_TEST                            map[string]map[string][]*smodel.CdnSetting
//	ONLINE_FREQUENCY_CONTROL_CONF           *mvutil.TestConf
//	CHANGE_ONLINE_DEEPLINK_WAY_PACKAGE_LIST []string
//	DEBUG_BID_FLOOR_AND_BID_PRICE_CONF      map[string]*mvutil.DebugBidFloorAndBidPriceConf
//	H265_VIDEO_ABTEST_CONF                  *mvutil.H265VideoABTestConf
//	REPLACE_TEMPLATE_URL_CONF               *mvutil.ReplaceTemplateUrlConf
//	TPL_CREATIVE_DOMAIN_CONF                map[string][]*mvutil.TemplateCreativeDomainMap
//	RETURN_WTICK_CONF                       *mvutil.ReturnWtickConf
//	EXCLUDE_PACKAGES_BY_CITYCODE_CONF       *mvutil.ExcludePackagesByCityCodeConf
//	ADNET_PARAMS_ABTEST_CONFS               map[string][]mvutil.ABTEST_CONF
//}

//var mvConfigObj *MVConfigObj

var mvconfigUpdateInputError = errors.New("mvconfigUpdateFunc failed, query or dbLoaderInfo is nil")

func NewConfigExtractor(allFilter, incFilter, selectOptions bson.M) dbLoaderInfo {
	return dbLoaderInfo{
		collection:   "config",
		querybmDbMap: allFilter,
		queryimDbMap: incFilter,
		querySelect:  selectOptions,
		updateFunc:   mvconfigUpdateFunc,
		dataCur:      concurrent_map.CreateConcurrentMap(concurrentMapPartitionsNum),
	}
}

func mvconfigUpdateFunc(query *mgo.Query, num int, dbLoaderInfo *dbLoaderInfo) (maxUpdate int64, err error) {
	if query == nil || dbLoaderInfo == nil {
		return 0, mvconfigUpdateInputError
	}
	mvConfig := &smodel.MVConfig{}
	maxUpdate = 0
	item := query.Iter()
	for item.Next(mvConfig) {
		mvconfigIncUpdateFunc(mvConfig, dbLoaderInfo)
		mvConfig = &smodel.MVConfig{}
	}
	if item.Err() != nil {
		logger.Warnf("mvconfigUpdateFunc err: %s", err.Error())
	}
	if err := item.Close(); err != nil {
		logger.Error(err.Error())
	}
	return maxUpdate, nil
}

func mvconfigIncUpdateFunc(mvconfig *smodel.MVConfig, dbLoaderInfo *dbLoaderInfo) {
	value, err := configUpdateFunc(mvconfig)
	if err != nil {
		logger.Warnf("configUpdateFunc err: %s", err.Error())
	}
	if value != nil {
		mvconfig.Value = value
		DbLoaderRegistry[TblConfig].dataCur.Set(concurrent_map.StrKey(mvconfig.Key), mvconfig)
	}
}

//func GetDecoodeSwith() (map[string]int, bool) {
//	return mvConfigObj.DecodeSwith, true
//}
//
//func getDecoodeSwith(jsonStr []byte) (map[string]int, bool) {
//
//	var values map[string]int
//
//	var json = jsoniter.ConfigCompatibleWithStandardLibrary
//	err := json.Unmarshal(jsonStr, &values)
//	if err != nil {
//		return values, false
//	}
//	return values, true
//}

func GetOpenapiScenario() ([]string, bool) {
	value, found := GetMVConfigValue("OPENAPI_SCENARIO")
	if !found {
		return nil, true
	}
	v, ok := value.([]string)
	if !ok {
		return nil, true
	}
	return v, true
}

func getOpenapiScenario(jsonStr []byte) ([]string, bool) {

	var values []string

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &values)
	if err != nil {
		return values, false
	}
	return values, true
}

func GetConfigLowGradeApp() (int64, bool) {
	value, found := GetMVConfigValue("LOW_GRADE_APP_ID")
	if !found {
		return 0, true
	}
	v, ok := value.(int64)
	if !ok {
		return 0, true
	}
	return v, true
}

func getConfigLowGradeApp(jsonStr []byte) (int64, bool) {
	var value int64
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return int64(0), false
	}
	return value, true
}

func GetReplenishApp() ([]int64, bool) {
	value, found := GetMVConfigValue("REPLENISH_APP")
	if !found {
		return []int64{}, true
	}
	v, ok := value.([]int64)
	if !ok {
		return []int64{}, true
	}
	return v, true
}
func getReplenishApp(jsonStr []byte) ([]int64, bool) {

	var value []int64
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return []int64{}, false
	}
	return value, true
}

func GetAdvBlackSubIdList() (map[string]map[string]map[string]string, bool) {
	value, found := GetMVConfigValue("ADV_BLACK_SUBID_LIST_V2")
	if !found {
		return make(map[string]map[string]map[string]string), true
	}
	v, ok := value.(map[string]map[string]map[string]string)
	if !ok {
		return make(map[string]map[string]map[string]string), true
	}
	return v, true
}

func getAdvBlackSubIdList(jsonStr []byte) (map[string]map[string]map[string]string, bool) {

	value := make(map[string]map[string]map[string]string)

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value, false
	}
	return value, true
}

func GetAppPostList() (map[int32]*smodel.AppPostList, bool) {
	value, found := GetMVConfigValue("APP_POST_LIST_V2")
	if !found {
		return make(map[int32]*smodel.AppPostList), true
	}
	v, ok := value.(map[int32]*smodel.AppPostList)
	if !ok {
		return make(map[int32]*smodel.AppPostList), true
	}
	return v, true
}

func getAppPostList(jsonStr []byte) (map[int32]*smodel.AppPostList, bool) {

	value := make(map[int32]*smodel.AppPostList)

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value, false
	}
	return value, true
}

func GetOfferwallGuidelines() (map[int32]string, bool) {
	value, found := GetMVConfigValue("offerwall_guidelines")
	if !found {
		return make(map[int32]string), true
	}
	v, ok := value.(map[int32]string)
	if !ok {
		return make(map[int32]string), true
	}
	return v, true
}

func getOfferwallGuidelines(jsonStr []byte) (map[int32]string, bool) {

	value := make(map[int32]string)

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value, false
	}
	return value, true
}

func GetEndcard() (smodel.EndCard, bool) {
	value, found := GetMVConfigValue("ENDCARD_CONFIG")
	if !found {
		return smodel.EndCard{}, true
	}
	v, ok := value.(smodel.EndCard)
	if !ok {
		return smodel.EndCard{}, true
	}
	return v, true
}

func getEndcard(jsonStr []byte) (smodel.EndCard, bool) {

	var value smodel.EndCard

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value, false
	}
	return value, true
}

func GetAdServerTestConfig() mvutil.AdServerTestConfig {
	value, found := GetMVConfigValue("ALGO_TEST_CONFIG")
	if !found {
		return mvutil.AdServerTestConfig{}
	}
	v, ok := value.(mvutil.AdServerTestConfig)
	if !ok {
		return mvutil.AdServerTestConfig{}
	}
	return v
}

func getAdServerTestConfig(jsonStr []byte) (mvutil.AdServerTestConfig, bool) {

	var value mvutil.AdServerTestConfig

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value, false
	}
	return value, true
}

//func GetDecRvTemplate() {
//
//}

func GetDefRVTemplate() (*smodel.VideoTemplateUrlItem, bool) {
	value, found := GetMVConfigValue("DEF_RV_TEMPLATE")
	if !found {
		return &smodel.VideoTemplateUrlItem{}, true
	}
	v, ok := value.(*smodel.VideoTemplateUrlItem)
	if !ok {
		return &smodel.VideoTemplateUrlItem{}, true
	}
	return v, true
}

func getDefRVTemplate(jsonStr []byte) (*smodel.VideoTemplateUrlItem, bool) {

	value := &smodel.VideoTemplateUrlItem{}

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value, false
	}
	return value, true
}

func GetVersionCompare() (map[string]mvutil.VersionCompare, bool) {
	value, found := GetMVConfigValue("VERSION_COMPARE")
	if !found {
		return make(map[string]mvutil.VersionCompare), true
	}
	v, ok := value.(map[string]mvutil.VersionCompare)
	if !ok {
		return make(map[string]mvutil.VersionCompare), true
	}
	return v, true
}

func getVersionCompare(jsonStr []byte) (map[string]mvutil.VersionCompare, bool) {

	value := make(map[string]mvutil.VersionCompare)

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value, false
	}
	return value, true
}

func GetPlayableTest() (map[int64]mvutil.PlatableTest, bool) {
	value, found := GetMVConfigValue("PLAYABLE_TEST")
	if !found {
		return make(map[int64]mvutil.PlatableTest), true
	}
	v, ok := value.(map[int64]mvutil.PlatableTest)
	if !ok {
		return make(map[int64]mvutil.PlatableTest), true
	}
	return v, true
}

func getPlayableTest(jsonStr []byte) (map[int64]mvutil.PlatableTest, bool) {

	value := make(map[int64]mvutil.PlatableTest)

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value, false
	}
	return value, true
}

func GetTRACK_URL_CONFIG_NEW() (map[int32]mvutil.TRACK_URL_CONFIG_NEW, bool) {
	value, found := GetMVConfigValue("TRACK_URL_CONFIG_NEW")
	if !found {
		return make(map[int32]mvutil.TRACK_URL_CONFIG_NEW), true
	}
	v, ok := value.(map[int32]mvutil.TRACK_URL_CONFIG_NEW)
	if !ok {
		return make(map[int32]mvutil.TRACK_URL_CONFIG_NEW), true
	}
	return v, true
}

func getTRACK_URL_CONFIG_NEW(jsonStr []byte) (map[int32]mvutil.TRACK_URL_CONFIG_NEW, bool) {

	value := make(map[int32]mvutil.TRACK_URL_CONFIG_NEW)

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value, false
	}
	return value, true
}

func Get3S_CHINA_DOMAIN() (mvutil.CONFIG_3S_CHINA_DOMAIN, bool) {
	value, found := GetMVConfigValue("3S_CHINA_DOMAIN")
	if !found {
		return mvutil.CONFIG_3S_CHINA_DOMAIN{}, true
	}
	v, ok := value.(mvutil.CONFIG_3S_CHINA_DOMAIN)
	if !ok {
		return mvutil.CONFIG_3S_CHINA_DOMAIN{}, true
	}
	return v, true
}

func get3S_CHINA_DOMAIN(jsonStr []byte) (mvutil.CONFIG_3S_CHINA_DOMAIN, bool) {

	var value mvutil.CONFIG_3S_CHINA_DOMAIN

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value, false
	}
	return value, true
}

func GetJUMP_TYPE_CONFIG() (map[string]int32, bool) {
	value, found := GetMVConfigValue("JUMP_TYPE_CONFIG")
	if !found {
		return make(map[string]int32), true
	}
	v, ok := value.(map[string]int32)
	if !ok {
		return make(map[string]int32), true
	}
	return v, true
}

func getJUMP_TYPE_CONFIG(jsonStr []byte) (map[string]int32, bool) {

	value := make(map[string]int32)

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value, false
	}
	return value, true
}

func GetJUMP_TYPE_CONFIG_IOS() (map[string]int32, bool) {
	value, found := GetMVConfigValue("JUMP_TYPE_CONFIG_IOS")
	if !found {
		return make(map[string]int32), true
	}
	v, ok := value.(map[string]int32)
	if !ok {
		return make(map[string]int32), true
	}
	return v, true
}

func getJUMP_TYPE_CONFIG_IOS(jsonStr []byte) (map[string]int32, bool) {

	value := make(map[string]int32)

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value, false
	}
	return value, true
}

func GetJUMPTYPE_SDKVERSION() (map[string]map[string]string, bool) {
	value, found := GetMVConfigValue("JUMPTYPE_SDKVERSION")
	if !found {
		return make(map[string]map[string]string), true
	}
	v, ok := value.(map[string]map[string]string)
	if !ok {
		return make(map[string]map[string]string), true
	}
	return v, true
}

func getJUMPTYPE_SDKVERSION(jsonStr []byte) (map[string]map[string]string, bool) {

	value := make(map[string]map[string]string)

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value, false
	}
	return value, true
}

func GetADSTACKING() (mvutil.ADSTACKING, bool) {
	value, found := GetMVConfigValue("ADSTACKING")
	if !found {
		return mvutil.ADSTACKING{}, true
	}
	v, ok := value.(mvutil.ADSTACKING)
	if !ok {
		return mvutil.ADSTACKING{}, true
	}
	return v, true
}

func getADSTACKING(jsonStr []byte) (mvutil.ADSTACKING, bool) {

	var value mvutil.ADSTACKING

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value, false
	}
	return value, true
}

func GetSETTING_CONFIG() (mvutil.SETTING_CONFIG, bool) {
	value, found := GetMVConfigValue("SETTING_CONFIG")
	if !found {
		return mvutil.SETTING_CONFIG{}, true
	}
	v, ok := value.(mvutil.SETTING_CONFIG)
	if !ok {
		return mvutil.SETTING_CONFIG{}, true
	}
	return v, true
}

func getSETTING_CONFIG(jsonStr []byte) (mvutil.SETTING_CONFIG, bool) {

	var value mvutil.SETTING_CONFIG

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value, false
	}
	return value, true
}

func GetTemplate() (map[string]mvutil.Template, bool) {
	value, found := GetMVConfigValue("template")
	if !found {
		return make(map[string]mvutil.Template), true
	}
	v, ok := value.(map[string]mvutil.Template)
	if !ok {
		return make(map[string]mvutil.Template), true
	}
	return v, true
}

func getTemplate(jsonStr []byte) (map[string]mvutil.Template, bool) {

	value := make(map[string]mvutil.Template)

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value, false
	}
	return value, true
}

func Getofferwall_urls() (mvutil.COfferwallUrls, bool) {
	value, found := GetMVConfigValue("offerwall_urls")
	if !found {
		return mvutil.COfferwallUrls{}, true
	}
	v, ok := value.(mvutil.COfferwallUrls)
	if !ok {
		return mvutil.COfferwallUrls{}, true
	}
	return v, true
}

func getofferwall_urls(jsonStr []byte) (mvutil.COfferwallUrls, bool) {

	var value mvutil.COfferwallUrls

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value, false
	}
	return value, true
}

func GetSIGN_NO_CHECK_APPS() ([]int64, bool) {
	value, found := GetMVConfigValue("SIGN_NO_CHECK_APPS")
	if !found {
		return nil, true
	}
	v, ok := value.([]int64)
	if !ok {
		return nil, true
	}
	return v, true
}

func getSIGN_NO_CHECK_APPS(jsonStr []byte) ([]int64, bool) {

	var value []int64

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value, false
	}
	return value, true
}

func GetC_TOI() (int32, bool) {
	value, found := GetMVConfigValue("C_TOI")
	if !found {
		return 0, true
	}
	v, ok := value.(int32)
	if !ok {
		return 0, true
	}
	return v, true
}

func getCToi(jsonStr []byte) (int32, bool) {

	var value int32
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value, false
	}
	return value, true
}

func GetMP_MAP_UNIT() (map[string]mvutil.MP_MAP_UNIT_, bool) {
	value, found := GetMVConfigValue("MP_MAP_UNIT_V2")
	if !found {
		return make(map[string]mvutil.MP_MAP_UNIT_), true
	}
	v, ok := value.(map[string]mvutil.MP_MAP_UNIT_)
	if !ok {
		return make(map[string]mvutil.MP_MAP_UNIT_), true
	}
	return v, true
}

func getMP_MAP_UNIT(jsonStr []byte) (map[string]mvutil.MP_MAP_UNIT_, bool) {

	value := make(map[string]mvutil.MP_MAP_UNIT_)

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value, false
	}
	return value, true
}

func GetCAN_CLICK_MODE_SIX_PUBLISHER() ([]int64, bool) {
	value, found := GetMVConfigValue("CAN_CLICK_MODE_SIX_PUBLISHER")
	if !found {
		return nil, true
	}
	v, ok := value.([]int64)
	if !ok {
		return nil, true
	}
	return v, true
}

func getCAN_CLICK_MODE_SIX_PUBLISHER(jsonStr []byte) ([]int64, bool) {

	var value []int64

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value, false
	}
	return value, true
}

func GetDEL_PRICE_FLOOR_UNIT() ([]int64, bool) {
	value, found := GetMVConfigValue("DEL_PRICE_FLOOR_UNIT")
	if !found {
		return nil, true
	}
	v, ok := value.([]int64)
	if !ok {
		return nil, true
	}
	return v, true
}

func getDEL_PRICE_FLOOR_UNIT(jsonStr []byte) ([]int64, bool) {

	var value []int64

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value, false
	}
	return value, true
}

// 第三方API用底价查询
func GetONLINE_PRICE_FLOOR_APPID() (map[string]float64, bool) {
	value, found := GetMVConfigValue("ONLINE_PRICE_FLOOR_APPID")
	if !found {
		return map[string]float64{}, true
	}
	v, ok := value.(map[string]float64)
	if !ok {
		return map[string]float64{}, true
	}
	return v, true
}

func getONLINE_PRICE_FLOOR_APPID(jsonStr []byte) (map[string]float64, bool) {

	var value map[string]float64

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return nil, false
	}
	return value, true
}

func GetCREATIVE_CHECK_MIGU_ADX_CREATIVE_IDS() (map[string]string, bool) {
	value, found := GetMVConfigValue("CREATIVE_CHECK_MIGU_ADX_CREATIVE_IDS")
	if !found {
		return map[string]string{}, true
	}
	v, ok := value.(map[string]string)
	if !ok {
		return map[string]string{}, true
	}
	return v, true
}

// 咪咕专用creative_id映射
func getCREATIVE_CHECK_MIGU_ADX_CREATIVE_IDS(jsonStr []byte) (map[string]string, bool) {

	var value map[string]string

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return nil, false
	}
	return value, true
}

func GetCREATIVE_CHECK_HUPU_ADX_CREATIVE_IDS() ([]int64, bool) {
	value, found := GetMVConfigValue("CREATIVE_CHECK_HUPU_ADX_CREATIVE_IDS")
	if !found {
		return nil, true
	}
	v, ok := value.([]int64)
	if !ok {
		return nil, true
	}
	return v, true
}

// 虎扑防未审核单子投放机制
func getCREATIVE_CHECK_HUPU_ADX_CREATIVE_IDS(jsonStr []byte) ([]int64, bool) {

	var value []int64

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return nil, false
	}
	return value, true
}

func GetIA_APPWALL() (string, bool) {
	value, found := GetMVConfigValue("IA_APPWALL")
	if !found {
		return "", true
	}
	v, ok := value.(string)
	if !ok {
		return "", true
	}
	return v, true
}

func getIA_APPWALL(jsonStr []byte) (string, bool) {

	var value string
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return "", false
	}
	return value, true
}

//func GetGdtKeyWhiteList() ([]string, bool) {
//	value, found := GetMVConfigValue("GDT_KEY_WHITE_LIST")
//	if !found {
//		return nil, true
//	}
//	v, ok := value.
//	([]string)
//	if !ok {
//		return nil, true
//	}
//	return v, true
//}

//func getGdtKeyWhiteList(jsonStr []byte) ([]string, bool) {
//
//	var value []string
//
//	var json = jsoniter.ConfigCompatibleWithStandardLibrary
//	err := json.Unmarshal(jsonStr, &value)
//	if err != nil {
//		return value, false
//	}
//	return value, true
//}

//func GetTOUTIAO_ALGO() (map[int64][]int64, bool) {
//	value, found := GetMVConfigValue("TOUTIAO_ALGO")
//	if !found {
//		return make(map[int64][]int64), true
//	}
//	v, ok := value.
//	(map[int64][]int64)
//	if !ok {
//		return make(map[int64][]int64), true
//	}
//	return v, true
//}

//func getTOUTIAO_ALGO(jsonStr []byte) (map[int64][]int64, bool) {
//
//	value := make(map[int64][]int64)
//
//	var json = jsoniter.ConfigCompatibleWithStandardLibrary
//	err := json.Unmarshal(jsonStr, &value)
//	if err != nil {
//		return value, false
//	}
//	return value, true
//}

func GetTO_NEW_CDN_APPS() ([]int64, bool) {
	value, found := GetMVConfigValue("TO_NEW_CDN_APPS")
	if !found {
		return nil, true
	}
	v, ok := value.([]int64)
	if !ok {
		return nil, true
	}
	return v, true
}

func getTO_NEW_CDN_APPS(jsonStr []byte) ([]int64, bool) {

	var value []int64

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value, false
	}
	return value, true
}

func GetUA_AB_TEST_CONFIG() (map[string]int, bool) {
	value, found := GetMVConfigValue("UA_AB_TEST_CONFIG")
	if !found {
		return make(map[string]int), true
	}
	v, ok := value.(map[string]int)
	if !ok {
		return make(map[string]int), true
	}
	return v, true
}

func getUA_AB_TEST_CONFIG(jsonStr []byte) (map[string]int, bool) {

	value := make(map[string]int)

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value, false
	}
	return value, true
}

//func GetGO_TRACK() (map[string]map[string]int, bool) {
//	return mvConfigObj.GO_TRACK, true
//}
//
//func getGO_TRACK(jsonStr []byte) (map[string]map[string]int, bool) {
//
//	value := make(map[string]map[string]int)
//
//	var json = jsoniter.ConfigCompatibleWithStandardLibrary
//	err := json.Unmarshal(jsonStr, &value)
//	if err != nil {
//		return value, false
//	}
//	return value, true
//}

// platformid->advid->conf
func GetJUMP_TYPE_CONFIG_ADV() (map[string]map[string]map[string]int32, bool) {
	value, found := GetMVConfigValue("JUMP_TYPE_CONFIG_ADV")
	if !found {
		return make(map[string]map[string]map[string]int32), true
	}
	v, ok := value.(map[string]map[string]map[string]int32)
	if !ok {
		return make(map[string]map[string]map[string]int32), true
	}
	return v, true
}

func getJUMP_TYPE_CONFIG_ADV(jsonStr []byte) (map[string]map[string]map[string]int32, bool) {

	value := make(map[string]map[string]map[string]int32)

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value, false
	}
	return value, true
}

// 风行mtg广告位map
func GetFUN_MV_MAP() (map[string]map[string]string, bool) {
	value, found := GetMVConfigValue("FUN_MV_MAP")
	if !found {
		return make(map[string]map[string]string), true
	}
	v, ok := value.(map[string]map[string]string)
	if !ok {
		return make(map[string]map[string]string), true
	}
	return v, true
}

func getFUN_MV_MAP(jsonStr []byte) (map[string]map[string]string, bool) {

	value := make(map[string]map[string]string)

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value, false
	}
	return value, true
}

// 不做强校验的appids
func GetNO_CHECK_PARAM_APP() (mvutil.NO_CHECK_PARAM_APP, bool) {
	value, found := GetMVConfigValue("NO_CHECK_PARAM_APP")
	if !found {
		return mvutil.NO_CHECK_PARAM_APP{}, true
	}
	v, ok := value.(mvutil.NO_CHECK_PARAM_APP)
	if !ok {
		return mvutil.NO_CHECK_PARAM_APP{}, true
	}
	return v, true
}

func getNO_CHECK_PARAM_APP(jsonStr []byte) (mvutil.NO_CHECK_PARAM_APP, bool) {

	var value mvutil.NO_CHECK_PARAM_APP

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value, false
	}
	return value, true
}

func GetHAS_EXTENSIONS_UNIT() ([]int64, bool) {
	value, found := GetMVConfigValue("HAS_EXTENSIONS_UNIT")
	if !found {
		return nil, true
	}
	v, ok := value.([]int64)
	if !ok {
		return nil, true
	}
	return v, true
}

func getHAS_EXTENSIONS_UNIT(jsonStr []byte) ([]int64, bool) {

	var value []int64

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value, false
	}
	return value, true
}

func Get_RETRUN_VAST_APP() []int64 {
	value, found := GetMVConfigValue("RETRUN_VAST_APP")
	if !found {
		return nil
	}
	v, ok := value.([]int64)
	if !ok {
		return nil
	}
	return v
}

func getRETRUN_VAST_APP(jsonStr []byte) []int64 {

	var value []int64

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value
	}
	return value
}

func GetUA_AB_TEST_SDK_OS_CONFIG() (map[string]string, bool) {
	value, found := GetMVConfigValue("UA_AB_TEST_SDK_OS_CONFIG")
	if !found {
		return make(map[string]string), true
	}
	v, ok := value.(map[string]string)
	if !ok {
		return make(map[string]string), true
	}
	return v, true
}

func getUA_AB_TEST_SDK_OS_CONFIG(jsonStr []byte) (map[string]string, bool) {

	value := make(map[string]string)

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value, false
	}
	return value, true
}

func GetGAMELOFT_CREATIVE_URLS() (map[string]map[string]string, bool) {
	value, found := GetMVConfigValue("GAMELOFT_CREATIVE_URLS")
	if !found {
		return make(map[string]map[string]string), true
	}
	v, ok := value.(map[string]map[string]string)
	if !ok {
		return make(map[string]map[string]string), true
	}
	return v, true
}

func getGAMELOFT_CREATIVE_URLS(jsonStr []byte) (map[string]map[string]string, bool) {

	value := make(map[string]map[string]string)

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value, false
	}
	return value, true
}

func GetUA_AB_TEST_THIRD_PARTY_CONFIG() ([]string, bool) {
	value, found := GetMVConfigValue("UA_AB_TEST_THIRD_PARTY_CONFIG")
	if !found {
		return nil, true
	}
	v, ok := value.([]string)
	if !ok {
		return nil, true
	}
	return v, true
}

func getUA_AB_TEST_THIRD_PARTY_CONFIG(jsonStr []byte) ([]string, bool) {

	var value []string

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value, false
	}
	return value, true
}

func GetHUPU_DEFAULT_UNITID() (map[string]int64, bool) {
	value, found := GetMVConfigValue("HUPU_DEFAULT_UNITID")
	if !found {
		return make(map[string]int64), true
	}
	v, ok := value.(map[string]int64)
	if !ok {
		return make(map[string]int64), true
	}
	return v, true
}

func getHUPU_DEFAULT_UNITID(jsonStr []byte) (map[string]int64, bool) {

	value := make(map[string]int64)

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value, false
	}
	return value, true
}

func GetUA_AB_TEST_CAMPAIGN_CONFIG() ([]int64, bool) {
	value, found := GetMVConfigValue("UA_AB_TEST_CAMPAIGN_CONFIG")
	if !found {
		return nil, true
	}
	v, ok := value.([]int64)
	if !ok {
		return nil, true
	}
	return v, true
}

func getUA_AB_TEST_CAMPAIGN_CONFIG(jsonStr []byte) ([]int64, bool) {

	var value []int64

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value, false
	}
	return value, true
}

func GetHUPU_DEFAULT_PRICE() (int, bool) {
	value, found := GetMVConfigValue("HUPU_DEFAULT_PRICE")
	if !found {
		return 0, true
	}
	v, ok := value.(int)
	if !ok {
		return 0, true
	}
	return v, true
}

func getHUPU_DEFAULT_PRICE(jsonStr []byte) (int, bool) {

	var value int
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return int(0), false
	}
	return value, true
}

func GetJUMP_TYPE_CONFIG_THIRD_PARTY() (map[string]map[string]map[string]int32, bool) {
	value, found := GetMVConfigValue("JUMP_TYPE_CONFIG_THIRD_PARTY")
	if !found {
		return make(map[string]map[string]map[string]int32), true
	}
	v, ok := value.(map[string]map[string]map[string]int32)
	if !ok {
		return make(map[string]map[string]map[string]int32), true
	}
	return v, true
}
func getJUMP_TYPE_CONFIG_THIRD_PARTY(jsonStr []byte) (map[string]map[string]map[string]int32, bool) {

	value := make(map[string]map[string]map[string]int32)

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value, false
	}
	return value, true
}

// 第三方API底价
func GetONLINE_PRICE_FLOOR() (map[string]float64, bool) {
	value, found := GetMVConfigValue("ONLINE_PRICE_FLOOR")
	if !found {
		return map[string]float64{}, true
	}
	v, ok := value.(map[string]float64)
	if !ok {
		return map[string]float64{}, true
	}
	return v, true
}

func getONLINE_PRICE_FLOOR(jsonStr []byte) (map[string]float64, bool) {

	var value map[string]float64

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value, false
	}
	return value, true
}

func GetLINKTYPE_UNITID() (map[string]int, bool) {
	value, found := GetMVConfigValue("LINKTYPE_UNITID")
	if !found {
		return make(map[string]int), true
	}
	v, ok := value.(map[string]int)
	if !ok {
		return make(map[string]int), true
	}
	return v, true
}
func getLINKTYPE_UNITID(jsonStr []byte) (map[string]int, bool) {

	value := make(map[string]int)

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value, false
	}
	return value, true
}

//func GetPLAYABLE_TEST_UNITS() ([]int64, bool) {
//	return mvConfigObj.PLAYABLE_TEST_UNITS, true
//}
//
//func getPLAYABLE_TEST_UNITS(jsonStr []byte) ([]int64, bool) {
//
//	var value []int64
//
//	var json = jsoniter.ConfigCompatibleWithStandardLibrary
//	err := json.Unmarshal(jsonStr, &value)
//	if err != nil {
//		return value, false
//	}
//	return value, true
//}

// 360max的appid，unitid
func GetMAX_APPID_AND_UNITID() (map[string]mvutil.ONLINE_APPID_UNITID, bool) {
	value, found := GetMVConfigValue("MAX_APPID_AND_UNITID")
	if !found {
		return map[string]mvutil.ONLINE_APPID_UNITID{}, true
	}
	v, ok := value.(map[string]mvutil.ONLINE_APPID_UNITID)
	if !ok {
		return map[string]mvutil.ONLINE_APPID_UNITID{}, true
	}
	return v, true
}

func getMAX_APPID_AND_UNITID(jsonStr []byte) (map[string]mvutil.ONLINE_APPID_UNITID, bool) {

	var value map[string]mvutil.ONLINE_APPID_UNITID

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value, false
	}
	return value, true
}

// func GetCHEETAH_CONFIG() (map[string]string, bool) {
//	return mvConfigObj.CHEETAH_CONFIG, true
// }
//
// // Config mapping for Cheetah
// func getCHEETAH_CONFIG() (map[string]string, bool) {
//	jsonStr, ifFind := GetMVConfigValue("CHEETAH_CONFIG")
//	var value map[string]string
//	if ifFind == false {
//		return value, false
//	}
//	var json = jsoniter.ConfigCompatibleWithStandardLibrary
//	err := json.Unmarshal(jsonStr, &value)
//	if err != nil {
//		return value, false
//	}
//	return value, true
// }

//func GetABTEST_GAIDIDFA() (map[string][]string, bool) {
//	return mvConfigObj.ABTEST_GAIDIDFA, true
//}
//
//func getABTEST_GAIDIDFA(jsonStr []byte) (map[string][]string, bool) {
//
//	value := make(map[string][]string)
//
//	var json = jsoniter.ConfigCompatibleWithStandardLibrary
//	err := json.Unmarshal(jsonStr, &value)
//	if err != nil {
//		return value, false
//	}
//	return value, true
//}

func GetREPLACE_BRAND_MODEL() (map[string]map[string]string, bool) {
	value, found := GetMVConfigValue("REPLACE_BRAND_MODEL")
	if !found {
		return make(map[string]map[string]string), true
	}
	v, ok := value.(map[string]map[string]string)
	if !ok {
		return make(map[string]map[string]string), true
	}
	return v, true
}
func getREPLACE_BRAND_MODEL(jsonStr []byte) (map[string]map[string]string, bool) {

	value := make(map[string]map[string]string)

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value, false
	}
	return value, true
}

//func GetNEW_PLAYABLE_SWITCH() (bool, bool) {
//	return mvConfigObj.NEW_PLAYABLE_SWITCH, true
//}
//
//func getNEW_PLAYABLE_SWITCH(jsonStr []byte) (bool, bool) {
//
//	var value bool
//
//	var json = jsoniter.ConfigCompatibleWithStandardLibrary
//	err := json.Unmarshal(jsonStr, &value)
//	if err != nil {
//		return value, false
//	}
//	return value, true
//}

//func GetPLAYABLE_ABTEST_RATE() (int, bool) {
//	return mvConfigObj.PLAYABLE_ABTEST_RATE, true
//}

//func getPLAYABLE_ABTEST_RATE(jsonStr []byte) (int, bool) {
//
//	var value int
//	var json = jsoniter.ConfigCompatibleWithStandardLibrary
//	err := json.Unmarshal(jsonStr, &value)
//	if err != nil {
//		return int(0), false
//	}
//	return value, true
//}

func GetPUB_CC_CDN() (map[string]map[string]string, bool) {
	value, found := GetMVConfigValue("PUB_CC_CDN")
	if !found {
		return map[string]map[string]string{}, true
	}
	v, ok := value.(map[string]map[string]string)
	if !ok {
		return map[string]map[string]string{}, true
	}
	return v, true
}

// Config mapping for Cheetah
func getPUB_CC_CDN(jsonStr []byte) (map[string]map[string]string, bool) {

	var value map[string]map[string]string

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value, false
	}
	return value, true
}

func GetTP_WITHOUT_RV() ([]int64, bool) {
	value, found := GetMVConfigValue("TP_WITHOUT_RV")
	if !found {
		return []int64{}, true
	}
	v, ok := value.([]int64)
	if !ok {
		return []int64{}, true
	}
	return v, true
}
func getTP_WITHOUT_RV(jsonStr []byte) ([]int64, bool) {

	var value []int64
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return []int64{}, false
	}
	return value, true
}

// pptv的appid，unitid
func GetPPTV_APPID_AND_UNITID() (map[string]mvutil.ONLINE_APPID_UNITID, bool) {
	value, found := GetMVConfigValue("PPTV_APPID_AND_UNITID")
	if !found {
		return map[string]mvutil.ONLINE_APPID_UNITID{}, true
	}
	v, ok := value.(map[string]mvutil.ONLINE_APPID_UNITID)
	if !ok {
		return map[string]mvutil.ONLINE_APPID_UNITID{}, true
	}
	return v, true
}

func getPPTV_APPID_AND_UNITID(jsonStr []byte) (map[string]mvutil.ONLINE_APPID_UNITID, bool) {

	var value map[string]mvutil.ONLINE_APPID_UNITID

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value, false
	}
	return value, true
}

// 凤凰视频的appid，unitid
func GetIFENG_APPID_AND_UNITID() (map[string]mvutil.ONLINE_APPID_UNITID, bool) {
	value, found := GetMVConfigValue("IFENG_APPID_AND_UNITID")
	if !found {
		return map[string]mvutil.ONLINE_APPID_UNITID{}, true
	}
	v, ok := value.(map[string]mvutil.ONLINE_APPID_UNITID)
	if !ok {
		return map[string]mvutil.ONLINE_APPID_UNITID{}, true
	}
	return v, true
}

func getIFENG_APPID_AND_UNITID(jsonStr []byte) (map[string]mvutil.ONLINE_APPID_UNITID, bool) {

	var value map[string]mvutil.ONLINE_APPID_UNITID

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value, false
	}
	return value, true
}

func GetREQUEST_BLACKLIST() (mvutil.REQUEST_BLACKLIST, bool) {
	value, found := GetMVConfigValue("REQUEST_BLACKLIST")
	if !found {
		return mvutil.REQUEST_BLACKLIST{}, true
	}
	v, ok := value.(mvutil.REQUEST_BLACKLIST)
	if !ok {
		return mvutil.REQUEST_BLACKLIST{}, true
	}
	return v, true
}

func getREQUEST_BLACKLIST(jsonStr []byte) (mvutil.REQUEST_BLACKLIST, bool) {

	var value mvutil.REQUEST_BLACKLIST

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value, false
	}
	return value, true
}

func GetPPTV_DEAL_ID() (string, bool) {
	value, found := GetMVConfigValue("PPTV_DEAL_ID")
	if !found {
		return "", true
	}
	v, ok := value.(string)
	if !ok {
		return "", true
	}
	return v, true
}

func getPPTV_DEAL_ID(jsonStr []byte) (string, bool) {

	var value string
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return "", false
	}
	return value, true
}

func GetCREATIVE_CHECK_PPTV_ADX_CREATIVE_IDS() (map[string]string, bool) {
	value, found := GetMVConfigValue("CREATIVE_CHECK_PPTV_ADX_CREATIVE_IDS")
	if !found {
		return map[string]string{}, true
	}
	v, ok := value.(map[string]string)
	if !ok {
		return map[string]string{}, true
	}
	return v, true
}

// 咪咕专用creative_id映射
func getCREATIVE_CHECK_PPTV_ADX_CREATIVE_IDS(jsonStr []byte) (map[string]string, bool) {

	var value map[string]string

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return nil, false
	}
	return value, true
}

// 猎豹的appid，unitid
func GetCM_APPID_AND_UNITID() (map[string]mvutil.ONLINE_APPID_UNITID, bool) {
	value, found := GetMVConfigValue("CM_APPID_AND_UNITID")
	if !found {
		return map[string]mvutil.ONLINE_APPID_UNITID{}, true
	}
	v, ok := value.(map[string]mvutil.ONLINE_APPID_UNITID)
	if !ok {
		return map[string]mvutil.ONLINE_APPID_UNITID{}, true
	}
	return v, true
}

func getCM_APPID_AND_UNITID(jsonStr []byte) (map[string]mvutil.ONLINE_APPID_UNITID, bool) {

	var value map[string]mvutil.ONLINE_APPID_UNITID

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value, false
	}
	return value, true
}

// 展示频次控制整体开关，true表示关闭频次控制
func GetFCA_SWITCH() (bool, bool) {
	value, found := GetMVConfigValue("FCA_SWITCH")
	if !found {
		return false, true
	}
	v, ok := value.(bool)
	if !ok {
		return false, true
	}
	return v, true
}

func getFCA_SWITCH(jsonStr []byte) (bool, bool) {

	var value bool

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value, false
	}
	return value, true
}

func GetFCA_CAMIDS() ([]int64, bool) {
	value, found := GetMVConfigValue("FCA_CAMIDS")
	if !found {
		return []int64{}, true
	}
	v, ok := value.([]int64)
	if !ok {
		return []int64{}, true
	}
	return v, true
}
func getFCA_CAMIDS(jsonStr []byte) ([]int64, bool) {

	var value []int64
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return []int64{}, false
	}
	return value, true
}

func GetNEW_DEFAULT_FCA() (int, bool) {
	value, found := GetMVConfigValue("NEW_DEFAULT_FCA")
	if !found {
		return 0, true
	}
	v, ok := value.(int)
	if !ok {
		return 0, true
	}
	return v, true
}

func getNEW_DEFAULT_FCA(jsonStr []byte) (int, bool) {

	var value int
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return int(0), false
	}
	return value, true
}

// 猎豹下毒app ids
//func GetCM_MP_APPIDS() ([]int64, bool) {
//	value, found := GetMVConfigValue("CM_MP_APPIDS")
//	if !found {
//		return []int64{}, true
//	}
//	v, ok := value.
//	([]int64)
//	if !ok {
//		return []int64{}, true
//	}
//	return v, true
//}
//func getCM_MP_APPIDS(jsonStr []byte) ([]int64, bool) {
//
//	var value []int64
//	var json = jsoniter.ConfigCompatibleWithStandardLibrary
//	err := json.Unmarshal(jsonStr, &value)
//	if err != nil {
//		return []int64{}, false
//	}
//	return value, true
//}

// 低效流量unit ids
func GetLOW_FLOW_UNITS() (map[string]int, bool) {
	value, found := GetMVConfigValue("LOW_FLOW_UNITS")
	if !found {
		return map[string]int{}, true
	}
	v, ok := value.(map[string]int)
	if !ok {
		return map[string]int{}, true
	}
	return v, true
}
func getLOW_FLOW_UNITS(jsonStr []byte) (map[string]int, bool) {

	var value map[string]int

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value, false
	}
	return value, true
}

func GetSS_ABTEST_CAMPAIGN() (map[string]*mvutil.SSTestTracfficRule, bool) {
	value, found := GetMVConfigValue("SS_ABTEST_CAMPAIGN")
	if !found {
		return map[string]*mvutil.SSTestTracfficRule{}, true
	}
	v, ok := value.(map[string]*mvutil.SSTestTracfficRule)
	if !ok {
		return map[string]*mvutil.SSTestTracfficRule{}, true
	}
	return v, true
}

func getSS_ABTEST_CAMPAIGN(jsonStr []byte) (map[string]*mvutil.SSTestTracfficRule, bool) {

	var value map[string]*mvutil.SSTestTracfficRule

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value, false
	}
	return value, true
}

//func GetNEW_CREATIVE_ABTEST_RATE() (int, bool) {
//	return mvConfigObj.NEW_CREATIVE_ABTEST_RATE, true
//}
//
//func getNEW_CREATIVE_ABTEST_RATE(jsonStr []byte) (int, bool) {
//
//	var value int
//	var json = jsoniter.ConfigCompatibleWithStandardLibrary
//	err := json.Unmarshal(jsonStr, &value)
//	if err != nil {
//		return int(0), false
//	}
//	return value, true
//}

//func GetNEW_CREATIVE_TEST_UNITS() (map[string]int32, bool) {
//	return mvConfigObj.NEW_CREATIVE_TEST_UNITS, true
//}
//
//func getNEW_CREATIVE_TEST_UNITS(jsonStr []byte) (map[string]int32, bool) {
//
//	var value map[string]int32
//
//	var json = jsoniter.ConfigCompatibleWithStandardLibrary
//	err := json.Unmarshal(jsonStr, &value)
//	if err != nil {
//		return value, false
//	}
//	return value, true
//}

// func GetREN_CLE_RES_SDK_OS_CONFIG() (map[string]string, bool) {
//	return mvConfigObj.REN_CLE_RES_SDK_OS_CONFIG, true
// }

func GetCREATIVE_ABTEST() (map[string]int, bool) {
	value, found := GetMVConfigValue("CREATIVE_ABTEST")
	if !found {
		return make(map[string]int), true
	}
	v, ok := value.(map[string]int)
	if !ok {
		return make(map[string]int), true
	}
	return v, true
}

//func GetCREATIVE_ABTEST() (map[string]int, bool) {
//	return mvConfigObj.CREATIVE_ABTEST, true
//}

//func getCREATIVE_ABTEST(jsonStr []byte) (map[string]int, bool) {
//
//	value := make(map[string]int)
//
//	var json = jsoniter.ConfigCompatibleWithStandardLibrary
//	err := json.Unmarshal(jsonStr, &value)
//	if err != nil {
//		return value, false
//	}
//	return value, true
//}

func GetTEMPLATE_MAP() (mvutil.GlobalTemplateMap, bool) {
	value, found := GetMVConfigValue("TEMPLATE_MAP")
	if !found {
		return mvutil.GlobalTemplateMap{}, true
	}
	v, ok := value.(mvutil.GlobalTemplateMap)
	if !ok {
		return mvutil.GlobalTemplateMap{}, true
	}
	return v, true
}

func getTEMPLATE_MAP(jsonStr []byte) (mvutil.GlobalTemplateMap, bool) {

	value := mvutil.GlobalTemplateMap{}

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value, false
	}
	return value, true
}

func GetNEW_CREATIVE_TIMEOUT() (int, bool) {
	value, found := GetMVConfigValue("NEW_CREATIVE_TIMEOUT")
	if !found {
		return 0, true
	}
	v, ok := value.(int)
	if !ok {
		return 0, true
	}
	return v, true
}

func getNEW_CREATIVE_TIMEOUT(jsonStr []byte) (int, bool) {

	var value int
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return int(0), false
	}
	return value, true
}

func GetXM_NEW_RETURN_UNITS() ([]int64, bool) {
	value, found := GetMVConfigValue("XM_NEW_RETURN_UNITS")
	if !found {
		return nil, true
	}
	v, ok := value.([]int64)
	if !ok {
		return nil, true
	}
	return v, true
}
func getXM_NEW_RETURN_UNITS(jsonStr []byte) ([]int64, bool) {

	var value []int64
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return []int64{}, false
	}
	return value, true
}

func GetWEBVIEW_PRISON_PUBIDS_APPIDS() (map[string][]int64, bool) {
	value, found := GetMVConfigValue("WEBVIEW_PRISON_PUBIDS_APPIDS")
	if !found {
		return map[string][]int64{}, true
	}
	v, ok := value.(map[string][]int64)
	if !ok {
		return map[string][]int64{}, true
	}
	return v, true
}
func getWEBVIEW_PRISON_PUBIDS_APPIDS(jsonStr []byte) (map[string][]int64, bool) {

	var value map[string][]int64
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return map[string][]int64{}, false
	}
	return value, true
}

func GetCONFIG_OFFER_PLCT() (map[string]map[string]mvutil.OFFER_PLCT, bool) {
	value, found := GetMVConfigValue("CONFIG_OFFER_PLCT")
	if !found {
		return map[string]map[string]mvutil.OFFER_PLCT{}, true
	}
	v, ok := value.(map[string]map[string]mvutil.OFFER_PLCT)
	if !ok {
		return map[string]map[string]mvutil.OFFER_PLCT{}, true
	}
	return v, true
}
func getCONFIG_OFFER_PLCT(jsonStr []byte) (map[string]map[string]mvutil.OFFER_PLCT, bool) {

	var value map[string]map[string]mvutil.OFFER_PLCT
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return map[string]map[string]mvutil.OFFER_PLCT{}, false
	}
	return value, true
}

func GetCCT_ABTEST_CONF() (map[string]map[string]int, bool) {
	value, found := GetMVConfigValue("CCT_ABTEST_CONF")
	if !found {
		return map[string]map[string]int{}, true
	}
	v, ok := value.(map[string]map[string]int)
	if !ok {
		return map[string]map[string]int{}, true
	}
	return v, true
}

func getCCT_ABTEST_CONF(jsonStr []byte) (map[string]map[string]int, bool) {

	var value map[string]map[string]int

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value, false
	}
	return value, true
}

func GetJSSDK_DOMAIN() (map[string]string, bool) {
	value, found := GetMVConfigValue("JSSDK_DOMAIN")
	if !found {
		return map[string]string{}, true
	}
	v, ok := value.(map[string]string)
	if !ok {
		return map[string]string{}, true
	}
	return v, true
}

func getJSSDK_DOMAIN(jsonStr []byte) (map[string]string, bool) {

	var value map[string]string
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return map[string]string{}, false
	}
	return value, true
}

func GetCHET_URL_UNIT() ([]int64, bool) {
	value, found := GetMVConfigValue("CHET_URL_UNIT")
	if !found {
		return nil, true
	}
	v, ok := value.([]int64)
	if !ok {
		return nil, true
	}
	return v, true
}

func getCHET_URL_UNIT(jsonStr []byte) ([]int64, bool) {

	var value []int64

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value, false
	}
	return value, true
}

//func GetMD5_FILE_PRISON_PUB() ([]int64, bool) {
//	return mvConfigObj.MD5_FILE_PRISON_PUB, true
//}
//
//func getMD5_FILE_PRISON_PUB(jsonStr []byte) ([]int64, bool) {
//
//	var value []int64
//
//	var json = jsoniter.ConfigCompatibleWithStandardLibrary
//	err := json.Unmarshal(jsonStr, &value)
//	if err != nil {
//		return value, false
//	}
//	return value, true
//}

func GetANR_PRISION_BY_PUB_AND_OSV() (mvutil.PUB_ADN_OSV, bool) {
	value, found := GetMVConfigValue("ANR_PRISION_BY_PUB_AND_OSV")
	if !found {
		return mvutil.PUB_ADN_OSV{}, true
	}
	v, ok := value.(mvutil.PUB_ADN_OSV)
	if !ok {
		return mvutil.PUB_ADN_OSV{}, true
	}
	return v, true
}

func getANR_PRISION_BY_PUB_AND_OSV(jsonStr []byte) (mvutil.PUB_ADN_OSV, bool) {

	var value mvutil.PUB_ADN_OSV

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value, false
	}
	return value, true
}

func GetUNIT_WITHOUT_VIDEO_START() ([]int64, bool) {
	value, found := GetMVConfigValue("UNIT_WITHOUT_VIDEO_START")
	if !found {
		return nil, true
	}
	v, ok := value.([]int64)
	if !ok {
		return nil, true
	}
	return v, true
}

func getUNIT_WITHOUT_VIDEO_START(jsonStr []byte) ([]int64, bool) {

	var value []int64

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value, false
	}
	return value, true
}

func GetREDUCE_FILL_SWITCH() (mvutil.REDUCE_FILL_SWITCH, bool) {
	value, found := GetMVConfigValue("REDUCE_FILL_SWITCH")
	if !found {
		return mvutil.REDUCE_FILL_SWITCH{}, true
	}
	v, ok := value.(mvutil.REDUCE_FILL_SWITCH)
	if !ok {
		return mvutil.REDUCE_FILL_SWITCH{}, true
	}
	return v, true
}

func getREDUCE_FILL_SWITCH(jsonStr []byte) (mvutil.REDUCE_FILL_SWITCH, bool) {

	var value mvutil.REDUCE_FILL_SWITCH

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value, false
	}
	return value, true
}

func GetMORE_OFFER_CONF() (mvutil.MORE_OFFER_CONFIG, bool) {
	value, found := GetMVConfigValue("MORE_OFFER_CONF")
	if !found {
		return mvutil.MORE_OFFER_CONFIG{}, true
	}
	v, ok := value.(mvutil.MORE_OFFER_CONFIG)
	if !ok {
		return mvutil.MORE_OFFER_CONFIG{}, true
	}
	return v, true
}

func getMORE_OFFER_CONF(jsonStr []byte) (mvutil.MORE_OFFER_CONFIG, bool) {

	var value mvutil.MORE_OFFER_CONFIG

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value, false
	}
	return value, true
}

func GetNEW_MOF_IMP_RATE() (int, bool) {
	value, found := GetMVConfigValue("NEW_MOF_IMP_RATE")
	if !found {
		return 0, true
	}
	v, ok := value.(int)
	if !ok {
		return 0, true
	}
	return v, true
}

func getNEW_MOF_IMP_RATE(jsonStr []byte) (int, bool) {

	var value int
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return int(0), false
	}
	return value, true
}

func GetMOF_ABTEST_RATE() (int, bool) {
	value, found := GetMVConfigValue("MOF_ABTEST_RATE")
	if !found {
		return 0, true
	}
	v, ok := value.(int)
	if !ok {
		return 0, true
	}
	return v, true
}

func getMOF_ABTEST_RATE(jsonStr []byte) (int, bool) {

	var value int
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return int(0), false
	}
	return value, true
}

func GetADX_FAKE_DSP_PRICE_DEFAULT_FACTOR() (map[string]float64, bool) {
	value, found := GetMVConfigValue("adx_fake_dsp_price_default_factor")
	if !found {
		return map[string]float64{}, true
	}
	v, ok := value.(map[string]float64)
	if !ok {
		return map[string]float64{}, true
	}
	return v, true
}

func getADX_FAKE_DSP_PRICE_DEFAULT_FACTOR(jsonStr []byte) (map[string]float64, bool) {

	var value map[string]float64

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return nil, false
	}
	return value, true
}

func GetADNET_SWITCHS() (map[string]int, bool) {
	value, found := GetMVConfigValue("ADNET_SWITCHS")
	if !found {
		return map[string]int{}, true
	}
	v, ok := value.(map[string]int)
	if !ok {
		return map[string]int{}, true
	}
	return v, true
}

func getADNET_SWITCHS(jsonStr []byte) (map[string]int, bool) {

	var value map[string]int

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value, false
	}
	return value, true
}

func GetBLACK_FOR_EXCLUDE_PACKAGE_NAME() (mvutil.BLACK_FOR_EXCLUDE_PACKAGE_NAME, bool) {
	value, found := GetMVConfigValue("BLACK_FOR_EXCLUDE_PACKAGE_NAME")
	if !found {
		return mvutil.BLACK_FOR_EXCLUDE_PACKAGE_NAME{}, true
	}
	v, ok := value.(mvutil.BLACK_FOR_EXCLUDE_PACKAGE_NAME)
	if !ok {
		return mvutil.BLACK_FOR_EXCLUDE_PACKAGE_NAME{}, true
	}
	return v, true
}

func getBLACK_FOR_EXCLUDE_PACKAGE_NAME(jsonStr []byte) (mvutil.BLACK_FOR_EXCLUDE_PACKAGE_NAME, bool) {

	var value mvutil.BLACK_FOR_EXCLUDE_PACKAGE_NAME

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value, false
	}
	return value, true
}

func GetAPPWALL_TO_MORE_OFFER_UNIT() (map[string]int32, bool) {
	value, found := GetMVConfigValue("APPWALL_TO_MORE_OFFER_UNIT")
	if !found {
		return map[string]int32{}, true
	}
	v, ok := value.(map[string]int32)
	if !ok {
		return map[string]int32{}, true
	}
	return v, true
}

func getAPPWALL_TO_MORE_OFFER_UNIT(jsonStr []byte) (map[string]int32, bool) {

	var value map[string]int32

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value, false
	}
	return value, true
}

func GetCleanEmptyDeviceABTest() (mvutil.CleanEmptyDeviceABTest, bool) {
	value, found := GetMVConfigValue("CLEAN_EMPTY_DEVICE_ABTEST")
	if !found {
		return mvutil.CleanEmptyDeviceABTest{}, true
	}
	v, ok := value.(mvutil.CleanEmptyDeviceABTest)
	if !ok {
		return mvutil.CleanEmptyDeviceABTest{}, true
	}
	return v, true
}

func getCleanEmptyDeviceABTest(jsonStr []byte) (mvutil.CleanEmptyDeviceABTest, bool) {

	var value mvutil.CleanEmptyDeviceABTest

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value, false
	}
	return value, true
}

func GetREQ_TYPE_AAB_TEST_CONFIG() (mvutil.REQ_TYPE_AAB_TEST_CONFIG, bool) {
	value, found := GetMVConfigValue("REQ_TYPE_AAB_TEST_CONFIG_NEW")
	if !found {
		return mvutil.REQ_TYPE_AAB_TEST_CONFIG{}, true
	}
	v, ok := value.(mvutil.REQ_TYPE_AAB_TEST_CONFIG)
	if !ok {
		return mvutil.REQ_TYPE_AAB_TEST_CONFIG{}, true
	}
	return v, true
}

func getREQ_TYPE_AAB_TEST_CONFIG(jsonStr []byte) (mvutil.REQ_TYPE_AAB_TEST_CONFIG, bool) {
	value := mvutil.REQ_TYPE_AAB_TEST_CONFIG{}
	if jsonStr == nil {
		logger.Infof("REQ_TYPE_AAB_TEST_CONFIG_NEW: Nil")
		return value, false
	}

	logger.Infof("REQ_TYPE_AAB_TEST_CONFIG_NEW: %s", string(jsonStr))

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value, false
	}
	return value, true
}

func GetExcludeDisplayPackageABTest() (mvutil.ExcludeDisplayPackageABTest, bool) {
	value, found := GetMVConfigValue("ExcludeDisplayPackageABTest")
	if !found {
		return mvutil.ExcludeDisplayPackageABTest{}, true
	}
	v, ok := value.(mvutil.ExcludeDisplayPackageABTest)
	if !ok {
		return mvutil.ExcludeDisplayPackageABTest{}, true
	}
	return v, true
}

func getExcludeDisplayPackageABTest(jsonStr []byte) (mvutil.ExcludeDisplayPackageABTest, bool) {

	value := mvutil.ExcludeDisplayPackageABTest{}

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value, false
	}
	return value, true
}

func GetREDUCE_FILL_FLINK_SWITCH() (mvutil.REDUCE_FILL_FLINK_SWITCH, bool) {
	value, found := GetMVConfigValue("REDUCE_FILL_FLINK_SWITCH")
	if !found {
		return mvutil.REDUCE_FILL_FLINK_SWITCH{}, true
	}
	v, ok := value.(mvutil.REDUCE_FILL_FLINK_SWITCH)
	if !ok {
		return mvutil.REDUCE_FILL_FLINK_SWITCH{}, true
	}
	return v, true
}

func getREDUCE_FILL_FLINK_SWITCH(jsonStr []byte) (mvutil.REDUCE_FILL_FLINK_SWITCH, bool) {

	var value mvutil.REDUCE_FILL_FLINK_SWITCH

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value, false
	}
	return value, true
}

func GetOnlineEmptyDeviceIPUAABTest() (mvutil.OnlineEmptyDeviceIPUA, bool) {
	value, found := GetMVConfigValue("ONLINE_EMPTY_DEVICE_IPUA_ABTEST")
	if !found {
		return mvutil.OnlineEmptyDeviceIPUA{}, true
	}
	v, ok := value.(mvutil.OnlineEmptyDeviceIPUA)
	if !ok {
		return mvutil.OnlineEmptyDeviceIPUA{}, true
	}
	return v, true
}

func getOnlineEmptyDeviceIPUAABTest(jsonStr []byte) (mvutil.OnlineEmptyDeviceIPUA, bool) {

	var value mvutil.OnlineEmptyDeviceIPUA

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value, false
	}
	return value, true
}

func GetAUTO_LOAD_CACHE_ABTSET() (mvutil.AUTO_LOAD_CACHE_ABTSET_CONFIG, bool) {
	value, found := GetMVConfigValue("AUTO_LOAD_CACHE_ABTSET")
	if !found {
		return mvutil.AUTO_LOAD_CACHE_ABTSET_CONFIG{}, true
	}
	v, ok := value.(mvutil.AUTO_LOAD_CACHE_ABTSET_CONFIG)
	if !ok {
		return mvutil.AUTO_LOAD_CACHE_ABTSET_CONFIG{}, true
	}
	return v, true
}

func getAUTO_LOAD_CACHE_ABTSET(jsonStr []byte) (mvutil.AUTO_LOAD_CACHE_ABTSET_CONFIG, bool) {

	value := mvutil.AUTO_LOAD_CACHE_ABTSET_CONFIG{}

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value, false
	}
	return value, true
}

func GetCTN_SIZE_ABTEST() (map[string]map[string]int32, bool) {
	value, found := GetMVConfigValue("CTN_SIZE_ABTEST")
	if !found {
		return map[string]map[string]int32{}, true
	}
	v, ok := value.(map[string]map[string]int32)
	if !ok {
		return map[string]map[string]int32{}, true
	}
	return v, true
}

func getCTN_SIZE_ABTEST(jsonStr []byte) (map[string]map[string]int32, bool) {

	var value map[string]map[string]int32

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value, false
	}
	return value, true
}

func GetBLOCK_MORE_OFFER_CONFIG() (mvutil.ENDCARD_BLOCK_MORE_OFFER_CONFIG, bool) {
	value, found := GetMVConfigValue("BLOCK_MORE_OFFER_CONFIG")
	if !found {
		return mvutil.ENDCARD_BLOCK_MORE_OFFER_CONFIG{}, true
	}
	v, ok := value.(mvutil.ENDCARD_BLOCK_MORE_OFFER_CONFIG)
	if !ok {
		return mvutil.ENDCARD_BLOCK_MORE_OFFER_CONFIG{}, true
	}
	return v, true
}

func getBLOCK_MORE_OFFER_CONFIG(jsonStr []byte) (mvutil.ENDCARD_BLOCK_MORE_OFFER_CONFIG, bool) {

	value := mvutil.ENDCARD_BLOCK_MORE_OFFER_CONFIG{}

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value, false
	}
	return value, true
}

func GetCLOSE_BUTTON_AD_TEST_UNITS() (map[string]int32, bool) {
	value, found := GetMVConfigValue("CLOSE_BUTTON_AD_TEST_UNITS")
	if !found {
		return map[string]int32{}, true
	}
	v, ok := value.(map[string]int32)
	if !ok {
		return map[string]int32{}, true
	}
	return v, true
}

func getCLOSE_BUTTON_AD_TEST_UNITS(jsonStr []byte) (map[string]int32, bool) {

	var value map[string]int32

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value, false
	}
	return value, true
}

//func GetMIGRATION_ABTEST() (map[string]int32, bool) {
//	return mvConfigObj.MIGRATION_ABTEST, true
//}
//
//func getMIGRATION_ABTEST(jsonStr []byte) (map[string]int32, bool) {
//
//	var value map[string]int32
//
//	var json = jsoniter.ConfigCompatibleWithStandardLibrary
//	err := json.Unmarshal(jsonStr, &value)
//	if err != nil {
//		return value, false
//	}
//	return value, true
//}

func GetRETURN_PARAM_K_UNIT() ([]int64, bool) {
	value, found := GetMVConfigValue("RETURN_PARAM_K_UNIT")
	if !found {
		return nil, true
	}
	v, ok := value.([]int64)
	if !ok {
		return nil, true
	}
	return v, true
}

func getRETURN_PARAM_K_UNIT(jsonStr []byte) ([]int64, bool) {

	var value []int64

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value, false
	}
	return value, true
}

func GetADD_UNIT_ID_ON_TRACKING_URL_UNIT() ([]int64, bool) {
	value, found := GetMVConfigValue("ADD_UNIT_ID_ON_TRACKING_URL_UNIT")
	if !found {
		return nil, true
	}
	v, ok := value.([]int64)
	if !ok {
		return nil, true
	}
	return v, true
}

func getADD_UNIT_ID_ON_TRACKING_URL_UNIT(jsonStr []byte) ([]int64, bool) {

	var value []int64

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value, false
	}
	return value, true
}

func GetCLICK_IN_ADTRACKING_THIRD_PARTY_CONFIG() (map[string]int, bool) {
	value, found := GetMVConfigValue("CLICK_IN_ADTRACKING_THIRD_PARTY_CONFIG")
	if !found {
		return map[string]int{}, true
	}
	v, ok := value.(map[string]int)
	if !ok {
		return map[string]int{}, true
	}
	return v, true
}

func getCLICK_IN_ADTRACKING_THIRD_PARTY_CONFIG(jsonStr []byte) (map[string]int, bool) {

	var value map[string]int

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value, false
	}
	return value, true
}

func GetPsbReduplicateConfig() mvutil.PsbReduplicateConfig {
	value, found := GetMVConfigValue("PsbReduplicateConfig")
	if !found {
		return mvutil.PsbReduplicateConfig{}
	}
	v, ok := value.(mvutil.PsbReduplicateConfig)
	if !ok {
		return mvutil.PsbReduplicateConfig{}
	}
	return v
}

func getPsbReduplicateConfig(jsonStr []byte) mvutil.PsbReduplicateConfig {

	value := mvutil.PsbReduplicateConfig{}

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	json.Unmarshal(jsonStr, &value)
	return value
}

func GetAopReduplicateConfig() mvutil.PsbReduplicateConfig {
	value, found := GetMVConfigValue("AopReduplicateConfig")
	if !found {
		return mvutil.PsbReduplicateConfig{}
	}
	v, ok := value.(mvutil.PsbReduplicateConfig)
	if !ok {
		return mvutil.PsbReduplicateConfig{}
	}
	return v
}

func getAopReduplicateConfig(jsonStr []byte) mvutil.PsbReduplicateConfig {

	value := mvutil.PsbReduplicateConfig{}

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	json.Unmarshal(jsonStr, &value)
	return value
}

func GetALAC_PRISON_CONFIG() (map[string][]int, bool) {
	value, found := GetMVConfigValue("ALAC_PRISON_CONFIG")
	if !found {
		return map[string][]int{}, true
	}
	v, ok := value.(map[string][]int)
	if !ok {
		return map[string][]int{}, true
	}
	return v, true
}

func getALAC_PRISON_CONFIG(jsonStr []byte) (map[string][]int, bool) {

	value := map[string][]int{}

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value, false
	}
	return value, true
}

func GetABTEST_FIELDS() (mvutil.ABTEST_FIELDS, bool) {
	value, found := GetMVConfigValue("ADNET_COM_ABTEST_FIELDS")
	if !found {
		return mvutil.ABTEST_FIELDS{}, true
	}
	v, ok := value.(mvutil.ABTEST_FIELDS)
	if !ok {
		return mvutil.ABTEST_FIELDS{}, true
	}
	return v, true
}

func getABTEST_FIELDS(jsonStr []byte) (mvutil.ABTEST_FIELDS, bool) {

	value := mvutil.ABTEST_FIELDS{}

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value, false
	}
	return value, true
}

func GetABTEST_CONFS() (map[string][]mvutil.ABTEST_CONF, bool) {
	value, found := GetMVConfigValue("ADNET_COM_ABTEST_CONFS")
	if !found {
		return map[string][]mvutil.ABTEST_CONF{}, true
	}
	v, ok := value.(map[string][]mvutil.ABTEST_CONF)
	if !ok {
		return map[string][]mvutil.ABTEST_CONF{}, true
	}
	return v, true
}

func getABTEST_CONFS(jsonStr []byte) (map[string][]mvutil.ABTEST_CONF, bool) {

	value := map[string][]mvutil.ABTEST_CONF{}

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value, false
	}
	return value, true
}

func GetChetLinkConfigs() ([]*mvutil.ChetLinkConfigItem, bool) {
	value, found := GetMVConfigValue("Adnet_Chet_Links")
	if !found {
		return nil, false
	}
	v, ok := value.([]*mvutil.ChetLinkConfigItem)
	if !ok {
		return nil, false
	}
	if len(v) == 0 {
		return nil, false
	}
	return v, true
}

func getChetLinkConfigs(jsonStr []byte) ([]*mvutil.ChetLinkConfigItem, bool) {

	var value []*mvutil.ChetLinkConfigItem

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value, false
	}
	return value, true
}

func GetAndroidLowVersionFilterCondition() (*mvutil.AndroidLowVersionFilterCondition, bool) {
	value, found := GetMVConfigValue("Android_Low_Version_Config")
	if !found {
		return nil, false
	}
	v, ok := value.(*mvutil.AndroidLowVersionFilterCondition)
	if !ok {
		return nil, false
	}
	if v == nil {
		return nil, false
	}
	return v, true
}

func getAndroidLowVersionFilterCondition(jsonStr []byte) (*mvutil.AndroidLowVersionFilterCondition, bool) {

	var value *mvutil.AndroidLowVersionFilterCondition

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value, false
	}
	return value, true
}

func GetBANNER_HTML_STR() (string, bool) {
	value, found := GetMVConfigValue("BANNER_HTML_STR")
	if !found {
		return "", true
	}
	v, ok := value.(string)
	if !ok {
		return "", true
	}
	return v, true
}

func getBANNER_HTML_STR(jsonStr []byte) (string, bool) {

	var value string
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return "", false
	}
	return value, true
}

func GetThirdPartyWhiteList() mvutil.ThirdPartyWhiteList {
	value, found := GetMVConfigValue("ThirdPartyWhiteList")
	if !found {
		return mvutil.ThirdPartyWhiteList{}
	}
	v, ok := value.(mvutil.ThirdPartyWhiteList)
	if !ok {
		return mvutil.ThirdPartyWhiteList{}
	}
	return v
}

func getThirdPartyWhiteList(jsonStr []byte) mvutil.ThirdPartyWhiteList {

	var value mvutil.ThirdPartyWhiteList

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value
	}
	return value
}

func GetSTOREKIT_TIME_PRISON_CONF() mvutil.StorekitTimePrisonConf {
	value, found := GetMVConfigValue("STOREKIT_TIME_PRISON_CONF")
	if !found {
		return mvutil.StorekitTimePrisonConf{}
	}
	v, ok := value.(mvutil.StorekitTimePrisonConf)
	if !ok {
		return mvutil.StorekitTimePrisonConf{}
	}
	return v
}

func getSTOREKIT_TIME_PRISON_CONF(jsonStr []byte) mvutil.StorekitTimePrisonConf {

	var value mvutil.StorekitTimePrisonConf

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value
	}
	return value
}

//func GetIPUA_WHITE_LIST_RATE_CONF() map[string]map[string]int {
//	return mvConfigObj.IPUA_WHITE_LIST_RATE_CONF
//}
//
//func getIPUA_WHITE_LIST_RATE_CONF(jsonStr []byte) map[string]map[string]int {
//
//	var value map[string]map[string]int
//
//	var json = jsoniter.ConfigCompatibleWithStandardLibrary
//	err := json.Unmarshal(jsonStr, &value)
//	if err != nil {
//		return value
//	}
//	return value
//}

func GetIS_RETURN_PAUSEMOD() mvutil.IsReturnPauseModConf {
	value, found := GetMVConfigValue("IS_RETURN_PAUSEMOD")
	if !found {
		return mvutil.IsReturnPauseModConf{}
	}
	v, ok := value.(mvutil.IsReturnPauseModConf)
	if !ok {
		return mvutil.IsReturnPauseModConf{}
	}
	return v
}

func getIS_RETURN_PAUSEMOD(jsonStr []byte) mvutil.IsReturnPauseModConf {

	var value mvutil.IsReturnPauseModConf

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value
	}
	return value
}

func GetCREATIVE_CHECK_MANGGUO_ADX_CREATIVE_IDS() map[string]string {
	value, found := GetMVConfigValue("CREATIVE_CHECK_MANGGUO_ADX_CREATIVE_IDS")
	if !found {
		return map[string]string{}
	}
	v, ok := value.(map[string]string)
	if !ok {
		return map[string]string{}
	}
	return v
}

// 芒果专用cteative_id映射，目前不确定是否需要使用。
func getCREATIVE_CHECK_MANGGUO_ADX_CREATIVE_IDS(jsonStr []byte) map[string]string {

	var value map[string]string
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value
	}
	return value
}

// 芒果tv的appid，unitid
func GetMANGO_APPID_AND_UNITID() map[string]map[string]mvutil.ONLINE_APPID_UNITID {
	value, found := GetMVConfigValue("MANGO_APPID_AND_UNITID")
	if !found {
		return map[string]map[string]mvutil.ONLINE_APPID_UNITID{}
	}
	v, ok := value.(map[string]map[string]mvutil.ONLINE_APPID_UNITID)
	if !ok {
		return map[string]map[string]mvutil.ONLINE_APPID_UNITID{}
	}
	return v
}

func getMANGO_APPID_AND_UNITID(jsonStr []byte) map[string]map[string]mvutil.ONLINE_APPID_UNITID {

	var value map[string]map[string]mvutil.ONLINE_APPID_UNITID

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value
	}
	return value
}

// func GetCLICK_IN_SERVER_CONF() mvutil.ClickInServerConf {
//	return mvConfigObj.CLICK_IN_SERVER_CONF
// }
//
// func getCLICK_IN_SERVER_CONF() mvutil.ClickInServerConf {
//	jsonStr, ifFind := GetMVConfigValue("CLICK_IN_SERVER_CONF")
//	var value mvutil.ClickInServerConf
//	if !ifFind {
//		return value
//	}
//	var json = jsoniter.ConfigCompatibleWithStandardLibrary
//	err := json.Unmarshal(jsonStr, &value)
//	if err != nil {
//		return value
//	}
//	return value
// }

func GetExcludeClickPackages() *mvutil.ExcludeClickPackages {
	value, found := GetMVConfigValue("ExcludeClickPackages")
	if !found {
		return new(mvutil.ExcludeClickPackages)
	}
	v, ok := value.(*mvutil.ExcludeClickPackages)
	if !ok {
		return new(mvutil.ExcludeClickPackages)
	}
	return v
}

func getExcludeClickPackages(jsonStr []byte) *mvutil.ExcludeClickPackages {

	value := new(mvutil.ExcludeClickPackages)

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, value)
	if err != nil {
		return value
	}
	return value
}

func GetADNET_CONF_LIST() map[string][]int64 {
	value, found := GetMVConfigValue("ADNET_CONF_LIST")
	if !found {
		return map[string][]int64{}
	}
	v, ok := value.(map[string][]int64)
	if !ok {
		return map[string][]int64{}
	}
	return v
}

func getADNET_CONF_LIST(jsonStr []byte) map[string][]int64 {

	var value map[string][]int64

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value
	}
	return value
}

func GetAdChoiceConfigData() *mvutil.AdChoiceConfigData {
	value, found := GetMVConfigValue("ADCHOICE_CONFIG")
	if !found {
		return new(mvutil.AdChoiceConfigData)
	}
	v, ok := value.(*mvutil.AdChoiceConfigData)
	if !ok {
		return new(mvutil.AdChoiceConfigData)
	}
	return v
}

func getAdChoiceConfigData(jsonStr []byte) *mvutil.AdChoiceConfigData {

	value := new(mvutil.AdChoiceConfigData)

	configData := mvutil.AdChoiceConfig{}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &configData)
	if err != nil {
		return value
	}
	value.AppIds = configData.AppIds
	value.UnitIds = configData.UnitIds
	var IOSSDKVerion, AndroidSDKVerion mvutil.DataSDKVersion
	for _, item := range configData.SdkVersions.Include {
		if strings.HasPrefix(strings.ToLower(item.Max), "mi_") {
			IOSSDKVerion.Include = append(IOSSDKVerion.Include, converionCode(item, 3))
		} else {
			AndroidSDKVerion.Include = append(AndroidSDKVerion.Include, converionCode(item, 4))
		}
	}

	for _, item := range configData.SdkVersions.Exclude {
		if strings.HasPrefix(strings.ToLower(item.Max), "mi_") {
			IOSSDKVerion.Exclude = append(IOSSDKVerion.Exclude, converionCode(item, 3))
		} else {
			AndroidSDKVerion.Exclude = append(AndroidSDKVerion.Exclude, converionCode(item, 4))
		}
	}

	if len(IOSSDKVerion.Exclude) > 0 || len(IOSSDKVerion.Include) > 0 {
		value.IOSSDKVersions = &IOSSDKVerion
	}

	if len(AndroidSDKVerion.Exclude) > 0 || len(AndroidSDKVerion.Include) > 0 {
		value.AndroidSDKVersions = &AndroidSDKVerion
	}
	return value
}

func GetFillRateEcpmFloorSwitch() bool {
	value, found := GetMVConfigValue("FILLRATE_ECPM_FLOOR_SWITCH")
	if !found {
		return false
	}
	v, ok := value.(bool)
	if !ok {
		return false
	}
	return v
}

func converionCode(item *mvutil.InfoSDKVersionItem, index int) *mvutil.DataSDKVersionItem {
	var dataItem mvutil.DataSDKVersionItem
	maxCode := mvutil.GetVersionCode(item.Max[index:])
	dataItem.Max = maxCode
	minCode := mvutil.GetVersionCode(item.Min[index:])
	dataItem.Min = minCode
	return &dataItem
}

// {"dspIds":[10,12]}
func getFillRateEcpmFloorSwitch(jsonStr []byte) bool {

	var value bool

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value
	}
	return value
}

func GetCLICK_IN_SERVER_CONF_NEW() mvutil.ClickInServerConf {
	value, found := GetMVConfigValue("CLICK_IN_SERVER_CONF_NEW")
	if !found {
		return mvutil.ClickInServerConf{}
	}
	v, ok := value.(mvutil.ClickInServerConf)
	if !ok {
		return mvutil.ClickInServerConf{}
	}
	return v
}

func getCLICK_IN_SERVER_CONF_NEW(jsonStr []byte) mvutil.ClickInServerConf {

	var value mvutil.ClickInServerConf

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value
	}
	return value
}

func GetMoattagConfig() *mvutil.MoattagConfigData {
	value, found := GetMVConfigValue("MOATTAG_CONFIG")
	if !found {
		return new(mvutil.MoattagConfigData)
	}
	v, ok := value.(*mvutil.MoattagConfigData)
	if !ok {
		return new(mvutil.MoattagConfigData)
	}
	return v
}

func getMoattagConfig(jsonStr []byte) *mvutil.MoattagConfigData {

	value := new(mvutil.MoattagConfigData)

	configData := mvutil.MoattagConfig{}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &configData)
	if err != nil {
		return value
	}
	var IOSSDKVerion, AndroidSDKVerion mvutil.DataSDKVersion
	for _, item := range configData.SdkVersions.Include {
		if strings.HasPrefix(strings.ToLower(item.Max), "mi_") {
			IOSSDKVerion.Include = append(IOSSDKVerion.Include, converionCode(item, 3))
		} else {
			AndroidSDKVerion.Include = append(AndroidSDKVerion.Include, converionCode(item, 4))
		}
	}

	for _, item := range configData.SdkVersions.Exclude {
		if strings.HasPrefix(strings.ToLower(item.Max), "mi_") {
			IOSSDKVerion.Exclude = append(IOSSDKVerion.Exclude, converionCode(item, 3))
		} else {
			AndroidSDKVerion.Exclude = append(AndroidSDKVerion.Exclude, converionCode(item, 4))
		}
	}

	if len(IOSSDKVerion.Exclude) > 0 || len(IOSSDKVerion.Include) > 0 {
		value.IOSSDKVersions = &IOSSDKVerion
	}

	if len(AndroidSDKVerion.Exclude) > 0 || len(AndroidSDKVerion.Include) > 0 {
		value.AndroidSDKVersions = &AndroidSDKVerion
	}
	return value
}

func GetBAD_REQUEST_FILTER_CONF() []*mvutil.BadRequestFilterConf {
	value, found := GetMVConfigValue("BAD_REQUEST_FILTER_CONF")
	if !found {
		return nil
	}
	v, ok := value.([]*mvutil.BadRequestFilterConf)
	if !ok {
		return nil
	}
	return v
}

func getBAD_REQUEST_FILTER_CONF(jsonStr []byte) []*mvutil.BadRequestFilterConf {

	var value []*mvutil.BadRequestFilterConf

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value
	}
	return value
}

func GetBIG_TEMPLATE_CONF() *mvutil.TestConf {
	value, found := GetMVConfigValue("BIG_TEMPLATE_CONF")
	if !found {
		return &mvutil.TestConf{}
	}
	v, ok := value.(*mvutil.TestConf)
	if !ok {
		return &mvutil.TestConf{}
	}
	return v
}

func getBIG_TEMPLATE_CONF(jsonStr []byte) *mvutil.TestConf {

	var value *mvutil.TestConf

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value
	}
	return value
}

func GetCNTrackingDomainConf() *mvutil.CNTrackingDomainTestConf {
	value, found := GetMVConfigValue("CN_TRACKING_DOMAIN_CONF")
	if !found {
		return &mvutil.CNTrackingDomainTestConf{}
	}
	v, ok := value.(*mvutil.CNTrackingDomainTestConf)
	if !ok {
		return &mvutil.CNTrackingDomainTestConf{}
	}
	return v
}

func getCNTrackingDomainConf(jsonStr []byte) *mvutil.CNTrackingDomainTestConf {

	var value *mvutil.CNTrackingDomainTestConf

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value
	}
	return value
}

func GetNoticeClickUniq() *mvutil.NoticeClickUniq {
	value, found := GetMVConfigValue("NoticeClickUniq")
	if !found {
		return new(mvutil.NoticeClickUniq)
	}
	v, ok := value.(*mvutil.NoticeClickUniq)
	if !ok {
		return new(mvutil.NoticeClickUniq)
	}
	return v
}

func getNoticeClickUniq(jsonStr []byte) *mvutil.NoticeClickUniq {

	value := new(mvutil.NoticeClickUniq)

	err := json.Unmarshal(jsonStr, value)
	if err != nil {
		return value
	}
	return value
}

// func GetAdjustPostbackABTest() *mvutil.AdjustPostbackABTest {
// 	return mvConfigObj.AdjustPostbackABTest
// }

// func getAdjustPostbackABTest() *mvutil.AdjustPostbackABTest {
// 	jsonStr, ifFind := GetMVConfigValue("AdjustPostbackABTest")
// 	value := new(mvutil.AdjustPostbackABTest)
// 	if !ifFind {
// 		return value
// 	}
// 	var json = jsoniter.ConfigCompatibleWithStandardLibrary
// 	err := json.Unmarshal(jsonStr, value)
// 	if err != nil {
// 		return value
// 	}
// 	return value
// }

// GetSupportVideoBanner 开关： 是否支持在video请求中并列加一个banner
// 1：开， 2：关
func GetSupportVideoBanner() bool {
	value, found := GetMVConfigValue("SUPPORT_VIDEO_BANNER")
	if !found {
		return false
	}
	v, ok := value.(int8)
	if !ok {
		logger.Errorf("GetSupportVideoBanner error: config value fails to type cast to int8")
		return false
	}
	return v == 1
}

func getSupportVideoBanner(jsonStr []byte) int8 {
	var value int8
	err := jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(jsonStr, &value)
	if err != nil {
		return int8(0)
	}
	return value
}

func GetExcludeImpressionPackagesV2() *mvutil.ExcludeClickPackagesV2 {
	value, found := GetMVConfigValue("ExcludeImpressionPackagesV2")
	if !found {
		return new(mvutil.ExcludeClickPackagesV2)
	}
	v, ok := value.(*mvutil.ExcludeClickPackagesV2)
	if !ok {
		return new(mvutil.ExcludeClickPackagesV2)
	}
	return v
}

func getExcludeImpressionPackagesV2(jsonStr []byte) *mvutil.ExcludeClickPackagesV2 {
	value := new(mvutil.ExcludeClickPackagesV2)
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, value)
	if err != nil {
		return value
	}
	return value
}

func GetExcludeImpressionPackages() *mvutil.ExcludeClickPackages {
	value, found := GetMVConfigValue("ExcludeImpressionPackages")
	if !found {
		return new(mvutil.ExcludeClickPackages)
	}
	v, ok := value.(*mvutil.ExcludeClickPackages)
	if !ok {
		return new(mvutil.ExcludeClickPackages)
	}
	return v
}

func getExcludeImpressionPackages(jsonStr []byte) *mvutil.ExcludeClickPackages {

	value := new(mvutil.ExcludeClickPackages)

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, value)
	if err != nil {
		return value
	}
	return value
}

func GetPassthroughData() []string {
	// 兜底配置
	defaultValue := []string{
		"OAID",
		"IDFV",
		"OsvUpTime",
	}
	value, found := GetMVConfigValue("PassthroughData")
	if !found {
		return defaultValue
	}
	v, ok := value.([]string)
	if !ok {
		return defaultValue
	}
	return v
}

func getPassthroughData(jsonStr []byte) []string {

	// 兜底配置
	value := []string{
		"OAID",
		"IDFV",
		"OsvUpTime",
	}

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value
	}
	return value
}

func GetPolarisFlagConf() *mvutil.TestConf {
	value, found := GetMVConfigValue("PolarisFlagConf")
	if !found {
		return &mvutil.TestConf{}
	}
	v, ok := value.(*mvutil.TestConf)
	if !ok {
		return &mvutil.TestConf{}
	}
	return v
}

func getPolarisFlagConf(jsonStr []byte) *mvutil.TestConf {

	var value *mvutil.TestConf

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value
	}
	return value
}

//func GetAdnetLangCreativeABTestConf() map[string][]map[string]int64 {
//	value, found := GetMVConfigValue("ADNET_LANG_CREATIVE_ABTEST_CONF")
//	if !found {
//		return nil
//	}
//	v, ok := value.(map[string][]map[string]int64)
//	if !ok {
//		return nil
//	}
//	return v
//}
//
//func getAdnetLangCreativeABTestConf(jsonStr []byte) map[string][]map[string]int64 {
//
//	var value map[string][]map[string]int64
//
//	var json = jsoniter.ConfigCompatibleWithStandardLibrary
//	err := json.Unmarshal(jsonStr, &value)
//	if err != nil {
//		return value
//	}
//	return value
//}

//func GetTemplateMapV2() map[string]map[string][]mvutil.TemplateMap {
//	return mvConfigObj.TemplateMapV2
//}
//
//func getTemplateMapV2() map[string]map[string][]mvutil.TemplateMap {
//	jsonStr, ifFind := GetMVConfigValue("TEMPLATE_MAP_V2")
//	var value map[string]map[string][]mvutil.TemplateMap
//	if !ifFind {
//		return value
//	}
//	var json = jsoniter.ConfigCompatibleWithStandardLibrary
//	err := json.Unmarshal(jsonStr, &value)
//	if err != nil {
//		return value
//	}
//	return value
//}

func GetFREQ_CONTROL_CONFIG() *mvutil.FreqControlConfig {
	value, found := GetMVConfigValue("FREQ_CONTROL_CONFIG")
	if !found {
		return new(mvutil.FreqControlConfig)
	}
	v, ok := value.(*mvutil.FreqControlConfig)
	if !ok {
		return new(mvutil.FreqControlConfig)
	}
	return v
}

func getFREQ_CONTROL_CONFIG(jsonStr []byte) *mvutil.FreqControlConfig {

	value := new(mvutil.FreqControlConfig)

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, value)
	if err != nil {
		return value
	}

	return value
}

func GetTIMEZONE_CONFIG() map[string]int {
	value, found := GetMVConfigValue("TIMEZONE_CONFIG")
	if !found {
		return map[string]int{}
	}
	v, ok := value.(map[string]int)
	if !ok {
		return map[string]int{}
	}
	return v
}

func getTIMEZONE_CONFIG(jsonStr []byte) map[string]int {

	var value map[string]int

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value
	}
	return value
}

//func GetBigTemplateMap() map[string][]mvutil.TemplateMap {
//	return mvConfigObj.BigTemplateMap
//}
//
//func getBigTemplateMap() map[string][]mvutil.TemplateMap {
//	jsonStr, ifFind := GetMVConfigValue("BIG_TEMPLATE_MAP")
//	var value map[string][]mvutil.TemplateMap
//	if !ifFind {
//		return value
//	}
//	var json = jsoniter.ConfigCompatibleWithStandardLibrary
//	err := json.Unmarshal(jsonStr, &value)
//	if err != nil {
//		return value
//	}
//	return value
//}
func GetCOUNTRY_CODE_TIMEZONE_CONFIG() map[string]int {
	value, found := GetMVConfigValue("COUNTRY_CODE_TIMEZONE_CONFIG")
	if !found {
		return map[string]int{}
	}
	v, ok := value.(map[string]int)
	if !ok {
		return map[string]int{}
	}
	return v
}

func getCOUNTRY_CODE_TIMEZONE_CONFIG(jsonStr []byte) map[string]int {

	var value map[string]int

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value
	}
	return value
}

func GetDSP_TPL_ABTEST_CONF() map[string][]*mvutil.DspTplAbtest {
	value, found := GetMVConfigValue("ADNET_DSP_TPL_ABTEST_CONF")
	if !found {
		return map[string][]*mvutil.DspTplAbtest{}
	}
	v, ok := value.(map[string][]*mvutil.DspTplAbtest)
	if !ok {
		return map[string][]*mvutil.DspTplAbtest{}
	}
	return v
}

func getDSP_TPL_ABTEST_CONF(jsonStr []byte) map[string][]*mvutil.DspTplAbtest {

	var value map[string][]*mvutil.DspTplAbtest

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value
	}
	return value
}

func GetBANNER_TO_AS_ABTEST_CONF() map[string]int {
	value, found := GetMVConfigValue("BANNER_TO_AS_ABTEST_CONF")
	if !found {
		return map[string]int{}
	}
	v, ok := value.(map[string]int)
	if !ok {
		return map[string]int{}
	}
	return v
}

func getBANNER_TO_AS_ABTEST_CONF(jsonStr []byte) map[string]int {

	var value map[string]int

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value
	}
	return value
}

func GetONLY_REQUEST_THIRD_DSP_SWITCH() bool {
	value, found := GetMVConfigValue("ONLY_REQUEST_THIRD_DSP_SWITCH")
	if !found {
		return false
	}
	v, ok := value.(bool)
	if !ok {
		return false
	}
	return v
}

func GetUSE_PLACEMENT_IMP_CAP_SWITCH() bool {
	value, found := GetMVConfigValue("USE_PLACEMENT_IMP_CAP_SWITCH")
	if !found {
		return false
	}
	v, ok := value.(bool)
	if !ok {
		return false
	}
	return v
}

func getONLY_REQUEST_THIRD_DSP_SWITCH(jsonStr []byte) bool {

	var value int
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return false
	}
	return value == 1
}

func GetDcoTestConf() *mvutil.TestConf {
	value, found := GetMVConfigValue("ADNET_DCO_TEST_CONF")
	if !found {
		return &mvutil.TestConf{}
	}
	v, ok := value.(*mvutil.TestConf)
	if !ok {
		return &mvutil.TestConf{}
	}
	return v
}

func getDcoTestConf(jsonStr []byte) *mvutil.TestConf {

	var value *mvutil.TestConf

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value
	}
	return value
}

func getUSE_PLACEMENT_IMP_CAP_SWITCH(jsonStr []byte) bool {

	var value int
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return false
	}
	return value == 1
}

// 低效流量unit ids
func GetLOW_FLOW_ADTYPE() (map[string]float64, bool) {
	value, found := GetMVConfigValue("LOW_FLOW_ADTYPE")
	if !found {
		return map[string]float64{}, true
	}
	v, ok := value.(map[string]float64)
	if !ok {
		return map[string]float64{}, true
	}
	return v, true
}

func getLOW_FLOW_ADTYPE(jsonStr []byte) (map[string]float64, bool) {

	var value map[string]float64

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value, false
	}
	return value, true
}

func GetTaoBaoOfferID() []int64 {
	value, found := GetMVConfigValue("TaoBaoOfferID")
	if !found {
		return []int64{}
	}
	v, ok := value.([]int64)
	if !ok {
		return []int64{}
	}
	return v
}

func getTaoBaoOfferID(jsonStr []byte) []int64 {

	value := []int64{}

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value
	}
	return value
}

func GetADNET_DEFAULT_VALUE() map[string]string {
	value, found := GetMVConfigValue("ADN_DEFAULT_VALUE")
	if !found {
		return map[string]string{}
	}
	v, ok := value.(map[string]string)
	if !ok {
		return map[string]string{}
	}
	return v
}

func getADNET_DEFAULT_VALUE(jsonStr []byte) map[string]string {

	var value map[string]string

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value
	}
	return value
}

func GetADNET_CREATIVE_COMPRESS_ABTEST_CONF_V3() map[string]*mvutil.CreativeCompressABTestV3Data {
	value, found := GetMVConfigValue("ADNET_CREATIVE_COMPRESS_ABTEST_CONF_V3")
	if !found {
		return map[string]*mvutil.CreativeCompressABTestV3Data{}
	}
	v, ok := value.(map[string]*mvutil.CreativeCompressABTestV3Data)
	if !ok {
		return map[string]*mvutil.CreativeCompressABTestV3Data{}
	}
	return v
}

func getADNET_CREATIVE_COMPRESS_ABTEST_CONF_V3(jsonStr []byte) map[string]*mvutil.CreativeCompressABTestV3Data {

	var value map[string]*mvutil.CreativeCompressABTestV3Data

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value
	}
	return value
}

func GetCREATIVE_VIDEO_COMPRESS_ABTEST() map[string]int32 {
	value, found := GetMVConfigValue("CREATIVE_VIDEO_COMPRESS_ABTEST")
	if !found {
		return map[string]int32{}
	}
	v, ok := value.(map[string]int32)
	if !ok {
		return map[string]int32{}
	}
	return v
}

func getCREATIVE_VIDEO_COMPRESS_ABTEST(jsonStr []byte) map[string]int32 {

	var value map[string]int32

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value
	}
	return value
}

func GetCREATIVE_IMG_COMPRESS_ABTEST() map[string]int32 {
	value, found := GetMVConfigValue("CREATIVE_IMG_COMPRESS_ABTEST")
	if !found {
		return map[string]int32{}
	}
	v, ok := value.(map[string]int32)
	if !ok {
		return map[string]int32{}
	}
	return v
}

func getCREATIVE_IMG_COMPRESS_ABTEST(jsonStr []byte) map[string]int32 {

	var value map[string]int32

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value
	}
	return value
}

func GetCREATIVE_ICON_COMPRESS_ABTEST() map[string]int32 {
	value, found := GetMVConfigValue("CREATIVE_ICON_COMPRESS_ABTEST")
	if !found {
		return map[string]int32{}
	}
	v, ok := value.(map[string]int32)
	if !ok {
		return map[string]int32{}
	}
	return v
}

func getCREATIVE_ICON_COMPRESS_ABTEST(jsonStr []byte) map[string]int32 {

	var value map[string]int32

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value
	}
	return value
}

func GetHBBidSDKVersionData(appId, unitId, os string) (string, bool) {
	// app_id+unit_id => app_id+all => all+os
	value, found := GetMVConfigValue("BID_SDK_VERSION_CONFIG")
	if !found {
		return "", false
	}
	v, ok := value.(map[string]map[string]string)
	if !ok {
		return "", false
	}
	bidSdkVerConf, ok := v[appId][unitId]
	if !ok {
		bidSdkVerConf, ok = v[appId]["all"]
		if !ok {
			bidSdkVerConf, ok = v["all"][os]
		}
	}
	if !ok || len(bidSdkVerConf) <= 0 {
		return "", false
	}
	return bidSdkVerConf, true
}

func getHBBidSdkVersionConfig(jsonStr []byte) map[string]map[string]string {

	var value map[string]map[string]string

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value
	}
	return value
}

func GetHBCurrencyRate(currency string) (float64, bool) {
	value, found := GetMVConfigValue("EXCHANGE_RATE")
	if !found {
		return 0.0, false
	}
	v, ok := value.(map[string]float64)
	if !ok {
		return 0.0, false
	}
	if currencyRate, ok := v[currency]; ok {
		return currencyRate, true
	}
	return 0.0, false
}

func getHBExchangeRate(jsonStr []byte) map[string]float64 {

	var value map[string]float64

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value
	}
	return value
}

func GetHBPublisherCurrency(publisherId string) (string, bool) {
	value, found := GetMVConfigValue("PUBLISHER_CURRENCY")
	if !found {
		return "", false
	}
	v, ok := value.(map[string]string)
	if !ok {
		return "", false
	}
	if publisherCurrency, ok := v[publisherId]; ok {
		return publisherCurrency, true
	}
	return "", false
}

func getHBPublisherCurrency(jsonStr []byte) map[string]string {

	var value map[string]string

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value
	}
	return value
}

//func GetHBBlackList(unitId string) ([]string, bool) {
//	if blackList, ok := mvConfigObj.HBBlacklist[unitId]; ok {
//		return blackList, true
//	}
//	return nil, false
//}
//
//func getHBBlacklist(jsonStr []byte) map[string][]string {
//
//	var value map[string][]string
//
//	var json = jsoniter.ConfigCompatibleWithStandardLibrary
//	err := json.Unmarshal(jsonStr, &value)
//	if err != nil {
//		return value
//	}
//	return value
//}

func GetHBPubExcludeApp() ([]int64, bool) {
	value, found := GetMVConfigValue("HBPubExcludeAppCheck")
	if !found {
		return nil, false
	}
	v, ok := value.([]int64)
	if !ok {
		return nil, false
	}
	return v, true
}

func getHBPubExcludeAppCheck(jsonStr []byte) []int64 {

	var value []int64

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value
	}
	return value
}

//func GetHBBidFloorPubWhiteList() ([]int64, bool) {
//	value, found := GetMVConfigValue("HBBidFloorPubWhiteList")
//	if !found {
//		return ([]int64, false
//		return nil, true
//	}
//	v, ok := value.
//	([]int64)
//	if !ok {
//		return ([]int64, false
//		return nil, true
//	}
//	return v, true
//	if len(mvConfigObj.HBBidFloorPubWhiteList) == 0 {
//		return nil, false
//	}
//	return mvConfigObj.HBBidFloorPubWhiteList, true
//}

//func getHBBidFloorPubWhiteList(jsonStr []byte) []int64 {
//
//	var value []int64
//
//	var json = jsoniter.ConfigCompatibleWithStandardLibrary
//	err := json.Unmarshal(jsonStr, &value)
//	if err != nil {
//		return value
//	}
//	return value
//}

func GetHBAdxEndpoint(region string) (*mvutil.HBAdxEndpointMkvData, bool) {
	value, found := GetMVConfigValue("HBAdxEndpoint")
	if !found {
		return nil, false
	}
	v, ok := value.(map[string]mvutil.HBAdxEndpointMkvData)
	if !ok {
		return nil, false
	}
	if endpoint, ok := v[region]; ok {
		return &endpoint, true
	}
	return nil, false
}

func getHBAdxEndpoint(jsonStr []byte) map[string]mvutil.HBAdxEndpointMkvData {

	var value map[string]mvutil.HBAdxEndpointMkvData

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value
	}
	return value
}

func GetHBAdxEndpointV2(cloud, region string) (*mvutil.HBAdxEndpointMkvData, bool) {
	value, found := GetMVConfigValue("HBAdxEndpointV2")
	if !found {
		return nil, false
	}
	v, ok := value.(map[string]map[string]mvutil.HBAdxEndpointMkvData)
	if !ok {
		return nil, false
	}
	if data, ok := v[cloud]; ok {
		if endpoint, ok := data[region]; ok {
			return &endpoint, true
		}
	}
	return nil, false
}

func getHBAdxEndpointV2(jsonStr []byte) map[string]map[string]mvutil.HBAdxEndpointMkvData {

	var value map[string]map[string]mvutil.HBAdxEndpointMkvData

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value
	}
	return value
}

func GetHBCDNDomainABTestPubs(publisherId string) (string, bool) {
	value, found := GetMVConfigValue("HBCDNDomainABTestPubs")
	if !found {
		return "", false
	}
	v, ok := value.(map[string]string)
	if !ok {
		return "", false
	}
	if domain, ok := v[publisherId]; ok {
		return domain, true
	}
	return "", false
}

func getHBCDNDomainABTestPubs(jsonStr []byte) map[string]string {

	var value map[string]string

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value
	}
	return value
}

//func GetAerospikeUseConsul() bool {
//	return mvConfigObj.AerospikeUseConsul
//}
//
//func getAerospikeUseConsul(jsonStr []byte) bool {
//
//	var value int
//	var json = jsoniter.ConfigCompatibleWithStandardLibrary
//	err := json.Unmarshal(jsonStr, &value)
//	if err != nil {
//		return false
//	}
//	return value == 1
//}

//func GetUseConsulServices(service string) bool {
//	useConsul, ok := mvConfigObj.UseConsulServices[service]
//	if !ok {
//		return false
//	}
//	return useConsul
//}

//func getUseConsulServices(jsonStr []byte) map[string]bool {
//
//	var value map[string]bool
//
//	var json = jsoniter.ConfigCompatibleWithStandardLibrary
//	err := json.Unmarshal(jsonStr, &value)
//	if err != nil {
//		return value
//	}
//	return value
//}

//
func GetHBAerospikeGzipRate(cloud, region string) float64 {
	value, found := GetMVConfigValue("HBAerospikeStorageConf")
	if !found {
		return 0
	}

	confArr, ok := value.(*mvutil.HBAerospikeStorageConfArr)
	if !ok {
		return 0
	}

	for _, conf := range confArr.Configs {
		if conf != nil && conf.Cloud == cloud && conf.Region == region {
			return conf.GzipRate
		}
	}
	//
	return 0
}

func GetRequestBidServerRate() *mvutil.HBRequestBidServerConfArr {
	value, found := GetMVConfigValue("HBRequestBidServerConf")
	if !found {
		return nil
	}
	confArr, ok := value.(*mvutil.HBRequestBidServerConfArr)
	if !ok {
		return nil
	}
	//
	return confArr
}

//
func GetHBAerospikeRemoveRedundancyRate(cloud, region string) float64 {
	value, found := GetMVConfigValue("HBAerospikeStorageConf")
	if !found {
		return 0
	}
	confArr, ok := value.(*mvutil.HBAerospikeStorageConfArr)
	if !ok {
		return 0
	}
	for _, conf := range confArr.Configs {
		if conf != nil && conf.Cloud == cloud && conf.Region == region {
			return conf.RemoveRedundancyRate
		}
	}
	return 0
}

//
func getHBAerospikeStorageConfArr(jsonStr []byte) *mvutil.HBAerospikeStorageConfArr {
	v := make([]*mvutil.HBAerospikeStorageConf, 0)
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &v)
	if err != nil {
		return &mvutil.HBAerospikeStorageConfArr{}
	}
	return &mvutil.HBAerospikeStorageConfArr{Configs: v}
}

//
func getHBRequestBidServerConfArr(jsonStr []byte) *mvutil.HBRequestBidServerConfArr {
	v := make([]*mvutil.HBRequestBidServerConf, 0)
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &v)
	if err != nil {
		return &mvutil.HBRequestBidServerConfArr{}
	}
	return &mvutil.HBRequestBidServerConfArr{Configs: v}
}

func GetUseConsulServicesV2Ratio(cloud, region, service string) float64 {
	value, found := GetMVConfigValue("UseConsulServicesV2Ratio")
	if !found {
		return 0
	}
	v, ok := value.(map[string]map[string]map[string]float64)
	if !ok {
		return 0
	}
	data, ok := v[cloud]
	if !ok {
		return 0
	}
	useConsulMap, ok := data[region]
	if !ok {
		return 0
	}
	r, ok := useConsulMap[service]
	if !ok {
		return 0
	}

	return r
}

func GetUseConsulServicesV2(cloud, region, service string) bool {
	value, found := GetMVConfigValue("UseConsulServicesV2")
	if !found {
		return false
	}
	v, ok := value.(map[string]map[string]map[string]bool)
	if !ok {
		return false
	}
	data, ok := v[cloud]
	if !ok {
		return false
	}
	useConsulMap, ok := data[region]
	if !ok {
		return false
	}
	useConsul := useConsulMap[service]
	return useConsul
}

func getUseConsulServicesV2Ratio(jsonStr []byte) map[string]map[string]map[string]float64 {

	var value map[string]map[string]map[string]float64

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value
	}
	return value
}
func getUseConsulServicesV2(jsonStr []byte) map[string]map[string]map[string]bool {

	var value map[string]map[string]map[string]bool

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value
	}
	return value
}

func GetTreasureBoxIDCount() (map[string]bool, bool) {
	value, found := GetMVConfigValue("TreasureBoxIDCount")
	if !found {
		return map[string]bool{}, true
	}
	v, ok := value.(map[string]bool)
	if !ok {
		return map[string]bool{}, true
	}
	return v, true
}

func getTreasureBoxIDCount(jsonStr []byte) map[string]bool {

	var value map[string]bool

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)

	// 直接设置配置--不管是否启用都无所谓，每分钟执行一次
	tb_tools.SetCollectCannotFoundIds(value)
	if err != nil {
		return value
	}
	return value
}

func GetCountryBlackPackageListConf() map[string]map[string][]string {
	value, found := GetMVConfigValue("CountryBlackPackageListConf")
	if !found {
		return map[string]map[string][]string{}
	}
	v, ok := value.(map[string]map[string][]string)
	if !ok {
		return map[string]map[string][]string{}
	}
	return v
}

func getCountryBlackPackageListConf(jsonStr []byte) map[string]map[string][]string {

	var value map[string]map[string][]string

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value
	}
	return value
}

func GetOnlineApiPubBidPriceConf() *mvutil.OnlineApiPubBidPriceConf {
	value, found := GetMVConfigValue("OnlineApiPubBidPriceConf")
	if !found {
		return &mvutil.OnlineApiPubBidPriceConf{}
	}
	v, ok := value.(*mvutil.OnlineApiPubBidPriceConf)
	if !ok {
		return &mvutil.OnlineApiPubBidPriceConf{}
	}
	return v
}

func getOnlineApiPubBidPriceConf(jsonStr []byte) *mvutil.OnlineApiPubBidPriceConf {

	var value *mvutil.OnlineApiPubBidPriceConf

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value
	}
	return value
}

func GetSupportTrackingTemplateConf() *mvutil.SupportTrackingTemplateConf {
	value, found := GetMVConfigValue("SupportTrackingTemplateConf")
	if !found {
		return new(mvutil.SupportTrackingTemplateConf)
	}
	v, ok := value.(*mvutil.SupportTrackingTemplateConf)
	if !ok {
		return new(mvutil.SupportTrackingTemplateConf)
	}
	return v
}

func getSupportTrackingTemplateConf(jsonStr []byte) *mvutil.SupportTrackingTemplateConf {

	value := new(mvutil.SupportTrackingTemplateConf)

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, value)
	if err != nil {
		return value
	}
	return value
}

// YLHClickModeTestConfig
func GetYLHClickModeTestConfig() map[string]int {
	value, found := GetMVConfigValue("YLHClickModeTestConfig")
	if !found {
		return map[string]int{}
	}
	v, ok := value.(map[string]int)
	if !ok {
		return map[string]int{}
	}
	return v
}
func getYLHClickModeTestConfig(jsonStr []byte) map[string]int {

	var value map[string]int

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value
	}
	return value
}

func GetV5_ABTEST_CONFIG() *mvutil.V5AbtestConf {
	value, found := GetMVConfigValue("V5_ABTEST_CONFIG")
	if !found {
		return new(mvutil.V5AbtestConf)
	}
	v, ok := value.(*mvutil.V5AbtestConf)
	if !ok {
		return new(mvutil.V5AbtestConf)
	}
	return v
}

func getV5_ABTEST_CONFIG(jsonStr []byte) *mvutil.V5AbtestConf {

	value := new(mvutil.V5AbtestConf)

	err := jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(jsonStr, value)
	if err != nil {
		return value
	}

	return value
}

func GetREPLACE_TRACKING_DOMAIN_CONF() *mvutil.ReplaceTrackingDomainConf {
	value, found := GetMVConfigValue("REPLACE_TRACKING_DOMAIN_CONF")
	if !found {
		return &mvutil.ReplaceTrackingDomainConf{}
	}
	v, ok := value.(*mvutil.ReplaceTrackingDomainConf)
	if !ok {
		return &mvutil.ReplaceTrackingDomainConf{}
	}
	return v
}

func getREPLACE_TRACKING_DOMAIN_CONF(jsonStr []byte) *mvutil.ReplaceTrackingDomainConf {

	var value *mvutil.ReplaceTrackingDomainConf

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value
	}
	return value
}

func GetNEW_AFREECATV_UNIT() []int64 {
	value, found := GetMVConfigValue("NEW_AFREECATV_UNIT")
	if !found {
		return nil
	}
	v, ok := value.([]int64)
	if !ok {
		return nil
	}
	return v
}

func getNEW_AFREECATV_UNIT(jsonStr []byte) []int64 {

	var value []int64

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value
	}
	return value
}

func GetSHAREIT_AF_WHITE_LIST_CONF() *mvutil.ShareitAfWhiteListConf {
	value, found := GetMVConfigValue("SHAREIT_AF_WHITE_LIST_CONF")
	if !found {
		return new(mvutil.ShareitAfWhiteListConf)
	}
	v, ok := value.(*mvutil.ShareitAfWhiteListConf)
	if !ok {
		return new(mvutil.ShareitAfWhiteListConf)
	}
	return v
}

func getSHAREIT_AF_WHITE_LIST_CONF(jsonStr []byte) *mvutil.ShareitAfWhiteListConf {

	value := new(mvutil.ShareitAfWhiteListConf)

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, value)
	if err != nil {
		return value
	}
	return value
}

func GetONLINE_API_USE_ADX() map[string]map[int64]int {
	value, found := GetMVConfigValue("ONLINE_API_USE_ADX")
	if !found {
		return map[string]map[int64]int{}
	}
	v, ok := value.(map[string]map[int64]int)
	if !ok {
		return map[string]map[int64]int{}
	}
	return v
}

func getONLINE_API_USE_ADX(jsonStr []byte) map[string]map[int64]int {

	var value map[string]map[int64]int

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value
	}
	return value
}

func GetONLINE_API_USE_ADX_MAS() map[string]map[int64]int {
	value, found := GetMVConfigValue("ONLINE_API_USE_ADX_MAS")
	if !found {
		return map[string]map[int64]int{}
	}
	v, ok := value.(map[string]map[int64]int)
	if !ok {
		return map[string]map[int64]int{}
	}
	return v
}

func getONLINE_API_USE_ADX_MAS(jsonStr []byte) map[string]map[int64]int {

	var value map[string]map[int64]int

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value
	}
	return value
}

func GetONLINE_API_SUPPORT_DEEPLINK_V2() map[string]map[int64]int32 {
	value, found := GetMVConfigValue("ONLINE_API_SUPPORT_DEEPLINK_V2")
	if !found {
		return map[string]map[int64]int32{}
	}
	v, ok := value.(map[string]map[int64]int32)
	if !ok {
		return map[string]map[int64]int32{}
	}
	return v
}

func getONLINE_API_SUPPORT_DEEPLINK_V2(jsonStr []byte) map[string]map[int64]int32 {
	var value map[string]map[int64]int32

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value
	}
	return value
}

func GetONLINE_API_MAX_BID_PRICE() *mvutil.OnlineApiMaxBidPrice {
	value, found := GetMVConfigValue("ONLINE_API_MAX_BID_PRICE")
	if !found {
		return new(mvutil.OnlineApiMaxBidPrice)
	}
	v, ok := value.(*mvutil.OnlineApiMaxBidPrice)
	if !ok {
		return new(mvutil.OnlineApiMaxBidPrice)
	}
	return v
}

func getONLINE_API_MAX_BID_PRICE(jsonStr []byte) *mvutil.OnlineApiMaxBidPrice {

	value := new(mvutil.OnlineApiMaxBidPrice)

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, value)
	if err != nil {
		return value
	}
	return value
}

func GetNEW_HUPU_UNITID_MAP() map[string]map[string]int64 {
	value, found := GetMVConfigValue("NEW_HUPU_UNITID_MAP")
	if !found {
		return make(map[string]map[string]int64)
	}
	v, ok := value.(map[string]map[string]int64)
	if !ok {
		return make(map[string]map[string]int64)
	}
	return v
}

func getNEW_HUPU_UNITID_MAP(jsonStr []byte) map[string]map[string]int64 {

	value := make(map[string]map[string]int64)

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value
	}
	return value
}

func GetNEW_CDN_TEST() map[string]map[string][]*smodel.CdnSetting {
	value, found := GetMVConfigValue("NEW_CDN_TEST")
	if !found {
		return map[string]map[string][]*smodel.CdnSetting{}
	}
	v, ok := value.(map[string]map[string][]*smodel.CdnSetting)
	if !ok {
		return map[string]map[string][]*smodel.CdnSetting{}
	}
	return v
}

func getNEW_CDN_TEST(jsonStr []byte) map[string]map[string][]*smodel.CdnSetting {

	var value map[string]map[string][]*smodel.CdnSetting

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value
	}
	return value
}

func GetONLINE_FREQUENCY_CONTROL_CONF() *mvutil.TestConf {
	value, found := GetMVConfigValue("ONLINE_FREQUENCY_CONTROL_CONF")
	if !found {
		return &mvutil.TestConf{}
	}
	v, ok := value.(*mvutil.TestConf)
	if !ok {
		return &mvutil.TestConf{}
	}
	return v
}

func getONLINE_FREQUENCY_CONTROL_CONF(jsonStr []byte) *mvutil.TestConf {

	var value *mvutil.TestConf

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value
	}
	return value
}

func GetCHANGE_ONLINE_DEEPLINK_WAY_PACKAGE_LIST() []string {
	value, found := GetMVConfigValue("CHANGE_ONLINE_DEEPLINK_WAY_PACKAGE_LIST")
	if !found {
		return nil
	}
	v, ok := value.([]string)
	if !ok {
		return nil
	}
	return v
}

func getCHANGE_ONLINE_DEEPLINK_WAY_PACKAGE_LIST(jsonStr []byte) []string {

	var value []string

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value
	}
	return value
}

func GetDEBUG_BID_FLOOR_AND_BID_PRICE_CONF() map[string]*mvutil.DebugBidFloorAndBidPriceConf {
	value, found := GetMVConfigValue("DEBUG_BID_FLOOR_AND_BID_PRICE_CONF")
	if !found {
		return map[string]*mvutil.DebugBidFloorAndBidPriceConf{}
	}
	v, ok := value.(map[string]*mvutil.DebugBidFloorAndBidPriceConf)
	if !ok {
		return map[string]*mvutil.DebugBidFloorAndBidPriceConf{}
	}
	return v
}

func getDEBUG_BID_FLOOR_AND_BID_PRICE_CONF(jsonStr []byte) map[string]*mvutil.DebugBidFloorAndBidPriceConf {

	var value map[string]*mvutil.DebugBidFloorAndBidPriceConf

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value
	}
	return value
}

func GetH265_VIDEO_ABTEST_CONF() *mvutil.H265VideoABTestConf {
	value, found := GetMVConfigValue("H265_VIDEO_ABTEST_CONF")
	if !found {
		return &mvutil.H265VideoABTestConf{}
	}
	v, ok := value.(*mvutil.H265VideoABTestConf)
	if !ok {
		return &mvutil.H265VideoABTestConf{}
	}
	return v
}

func getH265_VIDEO_ABTEST_CONF(jsonStr []byte) *mvutil.H265VideoABTestConf {
	var value *mvutil.H265VideoABTestConf
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value
	}
	return value
}

func GetREPLACE_TEMPLATE_URL_CONF() *mvutil.ReplaceTemplateUrlConf {
	value, found := GetMVConfigValue("REPLACE_TEMPLATE_URL_CONF")
	if !found {
		return &mvutil.ReplaceTemplateUrlConf{}
	}
	v, ok := value.(*mvutil.ReplaceTemplateUrlConf)
	if !ok {
		return &mvutil.ReplaceTemplateUrlConf{}
	}
	return v
}

func getREPLACE_TEMPLATE_URL_CONF(jsonStr []byte) *mvutil.ReplaceTemplateUrlConf {

	var value *mvutil.ReplaceTemplateUrlConf

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value
	}
	return value
}

func GetTPL_CREATIVE_DOMAIN_CONF() map[string][]*mvutil.TemplateCreativeDomainMap {
	value, found := GetMVConfigValue("TPL_CREATIVE_DOMAIN_CONF")
	if !found {
		return map[string][]*mvutil.TemplateCreativeDomainMap{}
	}
	v, ok := value.(map[string][]*mvutil.TemplateCreativeDomainMap)
	if !ok {
		return map[string][]*mvutil.TemplateCreativeDomainMap{}
	}
	return v
}

func getTPL_CREATIVE_DOMAIN_CONF(jsonStr []byte) map[string][]*mvutil.TemplateCreativeDomainMap {

	var value map[string][]*mvutil.TemplateCreativeDomainMap

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value
	}
	return value
}

func GetRETURN_WTICK_CONF() *mvutil.ReturnWtickConf {
	value, found := GetMVConfigValue("RETURN_WTICK_CONF")
	if !found {
		return &mvutil.ReturnWtickConf{}
	}
	v, ok := value.(*mvutil.ReturnWtickConf)
	if !ok {
		return &mvutil.ReturnWtickConf{}
	}
	return v
}

func getRETURN_WTICK_CONF(jsonStr []byte) *mvutil.ReturnWtickConf {

	var value *mvutil.ReturnWtickConf

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value
	}
	return value
}

func GetEXCLUDE_PACKAGES_BY_CITYCODE_CONF() *mvutil.ExcludePackagesByCityCodeConf {
	value, found := GetMVConfigValue("EXCLUDE_PACKAGES_BY_CITYCODE_CONF")
	if !found {
		return new(mvutil.ExcludePackagesByCityCodeConf)
	}
	v, ok := value.(*mvutil.ExcludePackagesByCityCodeConf)
	if !ok {
		return new(mvutil.ExcludePackagesByCityCodeConf)
	}
	return v
}

func getEXCLUDE_PACKAGES_BY_CITYCODE_CONF(jsonStr []byte) *mvutil.ExcludePackagesByCityCodeConf {
	value := new(mvutil.ExcludePackagesByCityCodeConf)
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, value)
	if err != nil {
		return value
	}
	return value
}

func GetADNET_PARAMS_ABTEST_CONFS() map[string][]mvutil.ABTEST_CONF {
	value, found := GetMVConfigValue("ADNET_PARAMS_ABTEST_CONFS")
	if !found {
		return map[string][]mvutil.ABTEST_CONF{}
	}
	v, ok := value.(map[string][]mvutil.ABTEST_CONF)
	if !ok {
		return map[string][]mvutil.ABTEST_CONF{}
	}
	return v
}

func getADNET_PARAMS_ABTEST_CONFS(jsonStr []byte) map[string][]mvutil.ABTEST_CONF {
	value := map[string][]mvutil.ABTEST_CONF{}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value
	}
	return value
}

func GetHB_API_CONFS() map[string]bool {
	value, found := GetMVConfigValue("HB_API_CONFS")
	if !found {
		return map[string]bool{}
	}
	v, ok := value.(map[string]bool)
	if !ok {
		return map[string]bool{}
	}
	return v
}

func getHB_API_CONFS(jsonStr []byte) map[string]bool {
	var value map[string]bool
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value
	}
	return value
}

func GetMAPPING_SERVER_RATE_CONF() map[string]int {
	value, found := GetMVConfigValue("MAPPING_SERVER_RATE_CONF")
	if !found {
		return nil
	}
	v, ok := value.(map[string]int)
	if !ok {
		return nil
	}
	return v
}

func GetMEDIATION_CHANNEL_ID() map[string]string {
	value, found := GetMVConfigValue("MEDIATION_CHANNEL_ID")
	if !found {
		return map[string]string{}
	}
	v, ok := value.(map[string]string)
	if !ok {
		return map[string]string{}
	}
	return v
}

func getMAPPING_SERVER_RATE_CONF(jsonStr []byte) map[string]int {
	var value map[string]int
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value
	}
	return value
}
func getMEDIATION_CHANNEL_ID(jsonStr []byte) map[string]string {
	var value map[string]string
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value
	}
	return value
}

func GetNEW_JUMP_TYPE_CONFIG_THIRD_PARTY() map[string]map[string]map[string]int32 {
	value, found := GetMVConfigValue("NEW_JUMP_TYPE_CONFIG_THIRD_PARTY")
	if !found {
		return make(map[string]map[string]map[string]int32)
	}
	v, ok := value.(map[string]map[string]map[string]int32)
	if !ok {
		return make(map[string]map[string]map[string]int32)
	}
	return v
}
func getNEW_JUMP_TYPE_CONFIG_THIRD_PARTY(jsonStr []byte) map[string]map[string]map[string]int32 {
	value := make(map[string]map[string]map[string]int32)
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value
	}
	return value
}

func GetAdPackageNameReplace() *mvutil.AdPackageNameReplaceConf {
	value, found := GetMVConfigValue("AD_PACKAGE_NAME_REPLACE_CONF")
	if !found {
		return new(mvutil.AdPackageNameReplaceConf)
	}
	v, ok := value.(*mvutil.AdPackageNameReplaceConf)
	if !ok {
		return new(mvutil.AdPackageNameReplaceConf)
	}
	return v
}

func getAdPackageNameReplace(jsonStr []byte) *mvutil.AdPackageNameReplaceConf {
	value := new(mvutil.AdPackageNameReplaceConf)
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, value)
	if err != nil {
		return value
	}
	return value
}

func GetFilterByStackConf() []*mvutil.FilterByStackConf {
	value, found := GetMVConfigValue("FILTER_BY_STACK_CONF")
	if !found {
		return nil
	}
	v, ok := value.([]*mvutil.FilterByStackConf)
	if !ok {
		return nil
	}
	return v
}

func GetHBRequestSnapshot() *mvutil.HBRequestSnapshot {
	value, found := GetMVConfigValue("HB_REQUEST_SNAPSHOT")
	if !found {
		return new(mvutil.HBRequestSnapshot)
	}
	v, ok := value.(*mvutil.HBRequestSnapshot)
	if !ok {
		return new(mvutil.HBRequestSnapshot)
	}
	return v
}

func getFilterByStackConf(jsonStr []byte) []*mvutil.FilterByStackConf {
	value := []*mvutil.FilterByStackConf{}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value
	}
	return value
}

func getHBRequestSnapshot(jsonStr []byte) *mvutil.HBRequestSnapshot {
	value := new(mvutil.HBRequestSnapshot)
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, value)
	if err != nil {
		return value
	}
	return value
}

func getMappingIdfaConf(jsonStr []byte) *mvutil.MappingIdfaAbtestConf {
	value := new(mvutil.MappingIdfaAbtestConf)
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, value)
	if err != nil {
		return value
	}
	return value
}

func GetMappingIdfaConf() *mvutil.MappingIdfaAbtestConf {
	value, found := GetMVConfigValue("MAPPING_IDFA_CONF")
	if !found {
		return new(mvutil.MappingIdfaAbtestConf)
	}
	v, ok := value.(*mvutil.MappingIdfaAbtestConf)
	if !ok {
		return new(mvutil.MappingIdfaAbtestConf)
	}
	return v
}

func GetOrientationPoisonConf() *mvutil.OrientationPoisonConf {
	value, found := GetMVConfigValue("ORIENTATION_POISON")
	if !found {
		return new(mvutil.OrientationPoisonConf)
	}
	v, ok := value.(*mvutil.OrientationPoisonConf)
	if !ok {
		return new(mvutil.OrientationPoisonConf)
	}
	return v
}

func GetTrackingCNABTestConf() *mvutil.TrackingCNABTestConf {
	value, found := GetMVConfigValue("TRACKING_CN_ABTEST_CONF")
	if !found {
		return new(mvutil.TrackingCNABTestConf)
	}
	v, ok := value.(*mvutil.TrackingCNABTestConf)
	if !ok {
		return new(mvutil.TrackingCNABTestConf)
	}
	return v
}

func getOrientationPoisonConf(jsonStr []byte) *mvutil.OrientationPoisonConf {
	value := new(mvutil.OrientationPoisonConf)
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, value)
	if err != nil {
		return value
	}
	return value
}

func getTrackingCNABTestConf(jsonStr []byte) *mvutil.TrackingCNABTestConf {
	value := new(mvutil.TrackingCNABTestConf)
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, value)
	if err != nil {
		return value
	}
	return value
}

func getMappingIdfaCoverIdfaABTestConf(jsonStr []byte) *mvutil.MappingIdfaAbtestConf {
	value := new(mvutil.MappingIdfaAbtestConf)
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, value)
	if err != nil {
		return value
	}
	return value
}

func GetMappingIdfaCoverIdfaABTestConf() *mvutil.MappingIdfaAbtestConf {
	value, found := GetMVConfigValue("MAPPING_IDFA_COVER_IDFA_ABTEST_CONF")
	if !found {
		return new(mvutil.MappingIdfaAbtestConf)
	}
	v, ok := value.(*mvutil.MappingIdfaAbtestConf)
	if !ok {
		return new(mvutil.MappingIdfaAbtestConf)
	}
	return v
}

func getSupportSmartVBAConfig(jsonStr []byte) *mvutil.SupportSmartVBAConfig {
	value := new(mvutil.SupportSmartVBAConfig)
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, value)
	if err != nil {
		return value
	}
	return value
}
func GetSupportSmartVBAConfig() *mvutil.SupportSmartVBAConfig {
	value, found := GetMVConfigValue("SupportSmartVBAConfig")
	if !found {
		return new(mvutil.SupportSmartVBAConfig)
	}
	v, ok := value.(*mvutil.SupportSmartVBAConfig)
	if !ok {
		return new(mvutil.SupportSmartVBAConfig)
	}
	return v
}

func getMoreOfferAndAppwallMoveToPioneerABTestConf(jsonStr []byte) map[string]*mvutil.MoreOfferAndAppwallMoveToPioneerABTestConf {
	var value map[string]*mvutil.MoreOfferAndAppwallMoveToPioneerABTestConf
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value
	}
	return value
}
func GetMoreOfferAndAppwallMoveToPioneerABTestConf() map[string]*mvutil.MoreOfferAndAppwallMoveToPioneerABTestConf {
	value, found := GetMVConfigValue("MOREOFFER_AND_APPWALL_MOVE_TO_PIONEER_ABTEST_CONF")
	if !found {
		return map[string]*mvutil.MoreOfferAndAppwallMoveToPioneerABTestConf{}
	}
	v, ok := value.(map[string]*mvutil.MoreOfferAndAppwallMoveToPioneerABTestConf)
	if !ok {
		return map[string]*mvutil.MoreOfferAndAppwallMoveToPioneerABTestConf{}
	}
	return v
}

func GetMORE_OFFER_REQUEST_DOMAIN() map[string]string {
	value, found := GetMVConfigValue("MORE_OFFER_REQUEST_DOMAIN")
	if !found {
		return map[string]string{}
	}
	v, ok := value.(map[string]string)
	if !ok {
		return map[string]string{}
	}
	return v
}

func GetOnlinePublisherAdNumConfig() map[string]int32 {
	value, found := GetMVConfigValue("ONLINE_PUBLISHER_ADNUM_CONF")
	if !found {
		return make(map[string]int32)
	}
	v, ok := value.(map[string]int32)
	if !ok {
		return make(map[string]int32)
	}
	return v
}

func getMORE_OFFER_REQUEST_DOMAIN(jsonStr []byte) map[string]string {
	var value map[string]string
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return nil
	}
	return value
}

func getOnlinePublisherAdNumConfig(jsonStr []byte) map[string]int32 {
	value := make(map[string]int32)
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value
	}
	return value
}

func GetTRUE_NUM_BY_AD_TYPE() map[string]int {
	value, found := GetMVConfigValue("TRUE_NUM_BY_AD_TYPE")
	if !found {
		return make(map[string]int)
	}
	v, ok := value.(map[string]int)
	if !ok {
		return make(map[string]int)
	}
	return v
}

func getTRUE_NUM_BY_AD_TYPE(jsonStr []byte) map[string]int {
	value := make(map[string]int)
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value
	}
	return value
}

func GetHBCachedAdNumConfig() map[string]int32 {
	value, found := GetMVConfigValue("HB_CACHED_ADNUM_CONF")
	if !found {
		return make(map[string]int32)
	}
	v, ok := value.(map[string]int32)
	if !ok {
		return make(map[string]int32)
	}
	return v
}

func getHBCachedAdNumConfig(jsonStr []byte) map[string]int32 {
	value := make(map[string]int32)
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value
	}
	return value
}

func GetMediationNoticeURLMacroConfig() *mvutil.MediationNoticeURLMacroConfValue {
	value, found := GetMVConfigValue("MEDIATION_NOTICE_URL_MACRO_CONF")
	if !found {
		return new(mvutil.MediationNoticeURLMacroConfValue)
	}
	v, ok := value.(*mvutil.MediationNoticeURLMacroConfValue)
	if !ok {
		return new(mvutil.MediationNoticeURLMacroConfValue)
	}
	return v
}

func getMediationNoticeURLMacroConfig(jsonStr []byte) *mvutil.MediationNoticeURLMacroConfValue {
	value := new(mvutil.MediationNoticeURLMacroConfValue)
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value
	}
	return value
}

func getHBOfferBidPriceABTestConf(jsonStr []byte) *mvutil.HBOfferBidPriceConf {
	value := new(mvutil.HBOfferBidPriceConf)
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, value)
	if err != nil {
		return value
	}
	return value
}

func GetTrackDomainByCountryCodeConf() map[string]*mvutil.TrackDomainByCountryCodeConf {
	value, found := GetMVConfigValue("TRACK_DOMAIN_BY_COUNTRY_CODE_CONF")
	if !found {
		return map[string]*mvutil.TrackDomainByCountryCodeConf{}
	}
	v, ok := value.(map[string]*mvutil.TrackDomainByCountryCodeConf)
	if !ok {
		return map[string]*mvutil.TrackDomainByCountryCodeConf{}
	}
	return v
}

func getTrackDomainByCountryCodeConf(jsonStr []byte) map[string]*mvutil.TrackDomainByCountryCodeConf {
	var value map[string]*mvutil.TrackDomainByCountryCodeConf
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value
	}
	return value
}

func GetHBOfferBidPriceABTestConf() *mvutil.HBOfferBidPriceConf {
	value, found := GetMVConfigValue("HB_OFFER_BID_PRICE_ABTEST_CONF")
	if !found {
		return new(mvutil.HBOfferBidPriceConf)
	}
	v, ok := value.(*mvutil.HBOfferBidPriceConf)
	if !ok {
		return new(mvutil.HBOfferBidPriceConf)
	}
	return v
}

func GetServiceDegradeRate(cloud, region, requestPath string) float64 {
	value, found := GetMVConfigValue("ServiceDegradeRate")
	if !found {
		return 0
	}
	v, ok := value.(map[string]map[string]map[string]float64)
	if !ok {
		return 0
	}
	data, ok := v[cloud]
	if !ok {
		return 0
	}
	dataDeep, ok := data[region]
	if !ok {
		return 0
	}
	r, ok := dataDeep[requestPath]
	if !ok {
		return 0
	}
	return r
}

func getServiceDegradeRate(jsonStr []byte) map[string]map[string]map[string]float64 {
	var value map[string]map[string]map[string]float64
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value
	}
	return value
}

func GetHBLoadDomainByCountryCodeConf() map[string]string {
	value, found := GetMVConfigValue("HB_LOAD_DOMAIN_BY_COUNTRY_CODE_CONF")
	if !found {
		return make(map[string]string)
	}
	v, ok := value.(map[string]string)
	if !ok {
		return make(map[string]string)
	}
	return v
}

func getHBLoadDomainByCountryCodeConf(jsonStr []byte) map[string]string {
	var value map[string]string
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value
	}
	return value
}

func GetMoreOfferRequestDomainByCountryCodeConf() map[string]map[string]string {
	value, found := GetMVConfigValue("MORE_OFFER_REQUEST_DOMAIN_BY_COUNTRY_CODE_CONF")
	if !found {
		return map[string]map[string]string{}
	}
	v, ok := value.(map[string]map[string]string)
	if !ok {
		return map[string]map[string]string{}
	}
	return v
}

func getMoreOfferRequestDomainByCountryCodeConf(jsonStr []byte) map[string]map[string]string {
	var value map[string]map[string]string
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return nil
	}
	return value
}

func GetLoadDomainABTest() *mvutil.LoadDomainABTest {
	value, found := GetMVConfigValue("LOAD_DOMAIN_ABTEST")
	if !found {
		return new(mvutil.LoadDomainABTest)
	}
	v, ok := value.(*mvutil.LoadDomainABTest)
	if !ok {
		return new(mvutil.LoadDomainABTest)
	}
	return v
}

func getLoadDomainABTest(jsonStr []byte) *mvutil.LoadDomainABTest {
	value := new(mvutil.LoadDomainABTest)
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, value)
	if err != nil {
		return nil
	}
	return value
}

func GetHBAerospikeConf() *mvutil.HBAerospikeConfMap {
	value, found := GetMVConfigValue("HB_AEROSPIKE_CONF")
	if !found {
		return new(mvutil.HBAerospikeConfMap)
	}
	v, ok := value.(*mvutil.HBAerospikeConfMap)
	if !ok {
		return new(mvutil.HBAerospikeConfMap)
	}
	return v
}

func getHBAerospikeConf(jsonStr []byte) *mvutil.HBAerospikeConfMap {
	value := new(mvutil.HBAerospikeConfMap)
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, value)
	if err != nil {
		return nil
	}
	return value
}

func getVastBannerDsp(jsonStr []byte) ([]int64, error) {
	var value []int64
	err := jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(jsonStr, &value)
	return value, err
}

// 对于banner请求 返回一个vast. vast里包括图片+模板url
func IsVastBannerDsp(dspId int64) bool {
	value, found := GetMVConfigValue("VAST_BANNER_DSP")
	if !found {
		return false
	}
	v, ok := value.([]int64)
	if !ok {
		return false
	}
	for _, id := range v {
		if id == dspId {
			return true
		}
	}
	return false
}
func GetResidualMemoryFilterConf() *mvutil.ResidualMemoryFilterConf {
	value, found := GetMVConfigValue("RESIDUAL_MEMORY_FILTER_CONF")
	if !found {
		return new(mvutil.ResidualMemoryFilterConf)
	}
	v, ok := value.(*mvutil.ResidualMemoryFilterConf)
	if !ok {
		return new(mvutil.ResidualMemoryFilterConf)
	}
	return v
}

func getResidualMemoryFilterConf(jsonStr []byte) *mvutil.ResidualMemoryFilterConf {
	value := new(mvutil.ResidualMemoryFilterConf)
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, value)
	if err != nil {
		return value
	}
	return value
}

func GetTmaxABTestConf() *mvutil.TmaxABTestConfMap {
	value, found := GetMVConfigValue("TMAX_ABTEST_CONF")
	if !found {
		return new(mvutil.TmaxABTestConfMap)
	}
	v, ok := value.(*mvutil.TmaxABTestConfMap)
	if !ok {
		return new(mvutil.TmaxABTestConfMap)
	}
	return v
}

func getTmaxABTestConf(jsonStr []byte) *mvutil.TmaxABTestConfMap {
	value := new(mvutil.TmaxABTestConfMap)
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, value)
	if err != nil {
		return nil
	}
	return value
}

func GetHBEeventHTTPProtocolConf() map[string]int32 {
	value, found := GetMVConfigValue("HB_EVENT_HTTP_PROTOCOL")
	if !found {
		return make(map[string]int32)
	}
	v, ok := value.(map[string]int32)
	if !ok {
		return make(map[string]int32)
	}
	return v
}

func getHBEeventHTTPProtocolConf(jsonStr []byte) map[string]int32 {
	value := make(map[string]int32)
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value
	}
	return value
}

func GetHBV5ABTestConf() *mvutil.HBV5ABTestConf {
	value, found := GetMVConfigValue("HB_V5_ABTEST_CONF")
	if !found {
		return new(mvutil.HBV5ABTestConf)
	}
	v, ok := value.(*mvutil.HBV5ABTestConf)
	if !ok {
		return new(mvutil.HBV5ABTestConf)
	}
	return v
}

func getHBV5ABTestConf(jsonStr []byte) *mvutil.HBV5ABTestConf {
	value := new(mvutil.HBV5ABTestConf)
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, value)
	if err != nil {
		return value
	}
	return value
}

//
//func GetAerospikeRemoveRedundancyRate(cloud, region string) float64 {
//	value, found := GetMVConfigValue("AerospikeRemoveRedundancyRate")
//	if !found {
//		return 0
//	}
//	v, ok := value.(map[string]map[string]map[string]float64)
//	if !ok {
//		return 0
//	}
//	data, ok := v[cloud]
//	if !ok {
//		return 0
//	}
//	useConsulMap, ok := data[region]
//	if !ok {
//		return 0
//	}
//	r, ok := useConsulMap["rate"]
//	if !ok {
//		return 0
//	}
//	return r
//}

//
//func getAerospikeRemoveRedundancyRate(jsonStr []byte) map[string]map[string]map[string]float64 {
//	var value map[string]map[string]map[string]float64
//	var json = jsoniter.ConfigCompatibleWithStandardLibrary
//	err := json.Unmarshal(jsonStr, &value)
//	if err != nil {
//		return value
//	}
//	return value
//}

func getMpToPioneerABTestConf(jsonStr []byte) *mvutil.MpToPioneerABTestConf {
	var value *mvutil.MpToPioneerABTestConf
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value
	}
	return value
}
func GetMpToPioneerABTestConf() *mvutil.MpToPioneerABTestConf {
	value, found := GetMVConfigValue("MP_TO_PIONEER_ABTEST_CONF")
	if !found {
		return new(mvutil.MpToPioneerABTestConf)
	}
	v, ok := value.(*mvutil.MpToPioneerABTestConf)
	if !ok {
		return new(mvutil.MpToPioneerABTestConf)
	}
	return v
}

func GetFilterAdserverRequestConf() map[string]int {
	value, found := GetMVConfigValue("FILTER_ADSERVER_REQUEST_CONF")
	if !found {
		return make(map[string]int)
	}
	v, ok := value.(map[string]int)
	if !ok {
		return make(map[string]int)
	}
	return v
}

func getFilterAdserverRequestConf(jsonStr []byte) map[string]int {
	value := make(map[string]int)
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value
	}
	return value
}

func GetHBCheckDeviceGEOCCConf() *mvutil.HBFilterDeviceGEOConf {
	value, found := GetMVConfigValue("HB_CHECK_DEVICE_GEO_CC")
	if !found {
		return new(mvutil.HBFilterDeviceGEOConf)
	}
	v, ok := value.(*mvutil.HBFilterDeviceGEOConf)
	if !ok {
		return new(mvutil.HBFilterDeviceGEOConf)
	}
	return v
}

func getHBCheckDeviceGEOCCConf(jsonStr []byte) *mvutil.HBFilterDeviceGEOConf {
	value := new(mvutil.HBFilterDeviceGEOConf)
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value
	}
	return value
}

func GetTEMPLATEMAPV2() map[string]map[string][]*mvutil.TemplateWeightMap {
	value, found := GetMVConfigValue("TEMPLATE_MAP_V2")
	if !found {
		return map[string]map[string][]*mvutil.TemplateWeightMap{}
	}
	v, ok := value.(map[string]map[string][]*mvutil.TemplateWeightMap)
	if !ok {
		return map[string]map[string][]*mvutil.TemplateWeightMap{}
	}
	return v
}

func getTEMPLATEMAPV2(jsonStr []byte) map[string]map[string][]*mvutil.TemplateWeightMap {
	value := map[string]map[string][]*mvutil.TemplateWeightMap{}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value
	}
	return value
}

func GetHBLoadFilterConfigs() *mvutil.HBLoadFilterConfigs {
	value, found := GetMVConfigValue("HBLoadFilterConfigs")
	if !found {
		return new(mvutil.HBLoadFilterConfigs)
	}
	v, ok := value.(*mvutil.HBLoadFilterConfigs)
	if !ok {
		return new(mvutil.HBLoadFilterConfigs)
	}
	return v
}

func getHBLoadFilterConfigs(jsonStr []byte) *mvutil.HBLoadFilterConfigs {
	value := new(mvutil.HBLoadFilterConfigs)
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value
	}
	return value
}

func GetAdnetSwitchConf(key string) int {
	conf, ok := GetADNET_SWITCHS()
	if !ok {
		return 0
	}
	val, ok := conf[key]
	if !ok {
		return 0
	}
	return val
}

func GetFilterAutoClickConf() *mvutil.FilterAutoClickConf {
	value, found := GetMVConfigValue("FILTER_AUTO_CLICK_CONF")
	if !found {
		return new(mvutil.FilterAutoClickConf)
	}
	v, ok := value.(*mvutil.FilterAutoClickConf)
	if !ok {
		return new(mvutil.FilterAutoClickConf)
	}
	return v
}

func getFilterAutoClickConf(jsonStr []byte) *mvutil.FilterAutoClickConf {
	value := new(mvutil.FilterAutoClickConf)
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, value)
	if err != nil {
		return value
	}
	return value
}

func GetTrackingCdnDomainMap() map[string]string {
	value, found := GetMVConfigValue("TRACKING_CDN_DOMAIN_MAP")
	if !found {
		return map[string]string{}
	}
	v, ok := value.(map[string]string)
	if !ok {
		return map[string]string{}
	}
	return v
}

func getTrackingCdnDomainMap(jsonStr []byte) map[string]string {
	value := make(map[string]string)
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value
	}
	return value
}

func GetCdnTrackingDomainABTestConf() []*mvutil.CdnTrackingDomainABTestConf {
	value, found := GetMVConfigValue("CDN_TRACKING_DOMAIN_ABTEST_CONF")
	if !found {
		return []*mvutil.CdnTrackingDomainABTestConf{}
	}
	v, ok := value.([]*mvutil.CdnTrackingDomainABTestConf)
	if !ok {
		return []*mvutil.CdnTrackingDomainABTestConf{}
	}
	return v
}

func getCdnTrackingDomainABTestConf(jsonStr []byte) []*mvutil.CdnTrackingDomainABTestConf {
	value := []*mvutil.CdnTrackingDomainABTestConf{}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal(jsonStr, &value)
	if err != nil {
		return value
	}
	return value
}

func GetCdnTrackingDomain(oriDomain string) (finalDomain string) {
	trackingCdnDomainMap := GetTrackingCdnDomainMap()
	originDomainKey := strings.ReplaceAll(oriDomain, ".", "_")
	cdnDomain, ok := trackingCdnDomainMap[originDomainKey]
	if !ok {
		return
	}
	return cdnDomain
}

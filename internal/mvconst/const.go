package mvconst

const (
	OffTypeMtg = 0
)

const (
	ContentAll   = 1
	ContentImg   = 2
	ContentVideo = 3
)

const (
	VideoAdTypeNOLimit = iota + 1
	VideoAdTypeNOVideo
	VideoAdTypeOnlyVideo
)

const (
	CampaignSourceOfferSync     = 0
	CampaignSourceAdnPortal     = 1
	CampaignSourceSSPlatform    = 2
	CampaignSourceSSAdvPlatform = 3
	CampaignSourceSSRMABTest    = 4
)

const (
	NormalAgency = 903
)

// campaignID 前 30亿段留给MV
const MaxMVOfferID = 3000000000
const BACKENDLEN = 1000000000

const MaxTrueNum = 100

const (
	FlowTagDefault = iota
	FlowTagDefaultM
	FlowTagSdkBanner = 101
	FlowTagMedia     = 102
)

const (
	DeviceUnknown = 11
	DevicePhone   = 4
	DeviceTablet  = 5
)

const (
	FixPriority = iota + 1
	Guarnteed
)

// 广告请求const杂项

// app grade
const (
	GRADE_A = iota + 1
	GRADE_B
	GRADE_C
	GRADE_D
)

// ctype

var CTALang = map[string]string{
	"ar":      "تثبيت",
	"en":      "Install",
	"fr":      "Installer",
	"de":      "Installieren",
	"id":      "Memasang",
	"it":      "Installare",
	"ja":      "インストール",
	"ko":      "설치",
	"nb":      "Installere",
	"pt":      "Instalar",
	"ru":      "Установить",
	"zh":      "安装",
	"es":      "Instalar",
	"sv":      "Installera",
	"th":      "ติดตั้ง",
	"zh-Hant": "安裝",
	"zh-Hans": "安装",
	"tr":      "İNDİR",
	"vi":      "cài đặt, dựng lên",
}

var CTAViewLang = map[string]string{
	"ar":      "للعرض",
	"en":      "View",
	"fr":      "Voir plus",
	"de":      "mehr sehen",
	"id":      "lihat lebih",
	"it":      "vedi dettagli",
	"ja":      "もっと",
	"ko":      "더 많은",
	"nb":      "se mer ",
	"pt":      "Veja mais ",
	"ru":      "Больше",
	"zh":      "查看",
	"es":      "ver más ",
	"sv":      " visa mer ",
	"th":      " ดูเพิ่มเติม",
	"zh-Hant": "查看",
	"zh-Hans": "查看",
	"tr":      "görünüm",
	"vi":      "lượt xem",
}

var SplashInstallLang = map[string]string{
	"zh":      "下载第三方应用",
	"zh-Hans": "下载第三方应用",
}

var SplashViewLang = map[string]string{
	"zh":      "浏览第三方页面",
	"zh-Hans": "浏览第三方页面",
}

const (
	Install = "Install"
	View    = "View"
	Open    = "Open"
)

const DspPublisherID = int64(6028)
const TOPONADXPublisherID = int64(22942)

const (
	DEVINFO_ENCRYPT_DONT = iota + 1
	DEVINFO_ENCRYPT_DO
	DEVINFO_ENCRYPT_NULL
)

// rand salt
const (
	SALT_RV_TEMPLATE                  = "salt_rv_template"
	SALT_ENDCARD                      = "salt_endcard"
	SALT_PLAYABLE                     = "salt_playable"
	SALT_CHECK_DEVID                  = "salt_check_devid"
	SALT_TTALGO                       = "salt_ttalgo"
	SALT_GOTRACK                      = "salt_gotrack"
	SALT_3SLINE                       = "salt_3sline"
	SALT_CCTABTEST                    = "salt_cct_abtest"
	SALT_REQTYPE_AABTEST              = "salt_reqtype_aabtest"
	SALT_DISPLAY_CAMPAIGN_ABTEST      = "salt_display_campaign_abtest"
	SALT_VCN_ABTEST                   = "salt_vcn_aabtest"
	SALT_CDN_ABTEST                   = "salt_cdn_abtest"
	SALT_DMP_ABTEST                   = "salt_dmp_abtest"
	SALT_PRICE_FACTOR_GROUP_ABTEST    = "salt_price_factor_group_abtest"
	SALT_PRICE_FACTOR_RATE_ABTEST     = "salt_price_factor_rate_abtest"
	SALT_PRICE_FACTOR_SUBRATE_ABTEST  = "salt_price_factor_subrate_abtest"
	SALT_SHAREIT_AF_WHITE_LIST_ABTEST = "salt_shareit_af_white_list_abtest"
	SALT_CLICK_IN_SERVER__ABTEST      = "salt_click_in_server"
	SALT_ONLINE_FREQUENCY_ABTEST      = "salt_online_frequency"
	SALT_TRACK_DOMAIN_ABTEST          = "salt_track_domain"
	SALT_TRACK_CN_ABTEST              = "salt_track_cn_abtest"
	SALT_RETURN_WTICK                 = "salt_return_wtick"
	SALT_RUID                         = "salt_ruid"
	SALT_APPSFLYER_UA_ABTEST          = "salt_appsflyer_ua_abtest"
	SALT_PACKAGE_NAME_REPLACE_ABTEST  = "salt_package_name_replace_abtest"
	SALT_MAPPING_IDFA                 = "salt_mapping_idfa"
	SALT_WTICK                        = "salt_wtick"
	SALT_REPLACE_PACKCAGE             = "salt_replace_package"
	SALT_MORE_OFFER_MV_TO_PIONEER     = "salt_more_offer_mv_to_pioneer"
	SALT_CN_TRACK_DOMAIN              = "salt_cn_track_domain"
	SALT_HB_OFFER_BID_PRICE           = "salt_hb_offer_bid_price"
	SALT_TRACK_DOMAIN_BY_COUNTRY_CODE = "salt_track_domain_by_country_code"
	SALT_TMAX_ABTEST                  = "salt_tmax_abtest"
	SALT_TMAX_GROUPS                  = "salt_tmax_groups"
	SALT_MP_TO_PIONEER                = "salt_mp_to_pioneer"
	SALT_HB_V5_ABTEST                 = "salt_hb_v5_abtest"
	SALT_HB_REQUEST_BID_SERVER        = "salt_hb_request_bid_server"
	SALT_THIRDPARTY_DSP_TPL_ABTEST    = "salt_thirdparty_dsp_tpl_abtest"
	SALT_CDN_TRACKING_DOMAIN          = "salt_cdn_tracking_domain"
)

// api version
const (
	API_VERSION_1_0 = float64(1)
	API_VERSION_1_1 = float64(1.1)
	API_VERSION_1_2 = float64(1.2)
	API_VERSION_1_3 = float64(1.3)
	API_VERSION_1_4 = float64(1.4) // 支持url新模板
	API_VERSION_1_5 = float64(1.5) // 支持新pv统计逻辑
	API_VERSION_1_9 = float64(1.9)
	API_VERSION_2_0 = float64(2)
	API_VERSION_2_1 = float64(2.1)
	API_VERSION_2_2 = float64(2.2)
	API_VERSION_2_3 = float64(2.3) // 新插屏
)

// orientation BOTH 竖屏 横屏
const (
	ORIENTATION_BOTH = iota
	ORIENTATION_PORTRAIT
	ORIENTATION_LANDSCAPE
)

const (
	INSTALL_FROM_CLICK = iota + 1
	INSTALL_FROM_IMPRESSION
)

func GetPlatformStr(platform int) string {
	plMap := map[int]string{
		PlatformAndroid: PlatformNameAndroid,
		PlatformIOS:     PlatformNameIOS,
	}
	name, ok := plMap[platform]
	if ok {
		return name
	}
	return PlatformNameOther
}

const (
	SERVER_SYSTEM_M        = "M"
	SERVER_SYSTEM_MOBPOWER = "MP"
	SERVER_SYSTEM_SA       = "SA"
	SERVER_SYSTEM_TEST     = "TEST"
)

const (
	CIRCUIT_HTTP   = "http_request"
	CIRCUIT_ADSERV = "adserver_request"
)

const (
	NO_CHECK_SIGN = "NO_CHECK_SIGN"
)

const (
	DELETE_DIVICEID_TRUE = iota + 1
	DELETE_DIVICEID_FALSE
	DELETE_DIVICEID_BUT_NOT_IMPRESSION
)

var AndroidVersionMap = map[int]string{
	28: "9",
	27: "8.1.0",
	26: "8.0",
	25: "7.1.1",
	24: "7.0",
	23: "6.0",
	22: "5.1",
	21: "5",
	19: "4.4",
	18: "4.3",
	17: "4.2",
	16: "4.1",
	15: "4.0.3",
	14: "4",
	13: "3.2",
	12: "3.1",
	11: "3",
	10: "2.3.3",
	9:  "2.3",
	8:  "2.2",
	7:  "2.1",
	6:  "2.0.1",
	5:  "2",
	4:  "1.6",
	3:  "1.5",
	2:  "1.1",
	1:  "1",
}

func GetResEmpty() *string {
	str := ""
	return &str
}

func GetResZero() *string {
	str := "0"
	return &str
}

// 3s开发者
const (
	THIRD_PARTY_APPSFLYER = "appsflyer"
	THIRD_PARTY_S2S       = "s2s"
	THIRD_PARTY_ADJUST    = "adjust"
)

const (
	IARST_PLAYABLE = iota + 1
	IARST_APPWALL
)

const (
	LINEAR_ADS           = 1
	SKIPPABLE_LINEAR_ADS = 2
)

const (
	PER_START          = 0
	PER_FIRST_QUARTILE = 25
	PER_MIDPOINT       = 50
	PER_THIRD_QUARTILE = 75
	PER_COMPLETE       = 100
)

func GetVastPercentage(perRate int) string {
	perArr := map[int]string{
		PER_START:          "start",
		PER_FIRST_QUARTILE: "firstQuartile",
		PER_MIDPOINT:       "midpoint",
		PER_THIRD_QUARTILE: "thirdQuartile",
		PER_COMPLETE:       "complete",
	}
	val, ok := perArr[perRate]
	if ok {
		return val
	} else {
		return ""
	}
}

const (
	APP_STROE_ADDR     = "itunes.apple.com"
	NEW_APP_STROE_ADDR = "apps.apple.com"
)

// 小程序单子
const (
	NETWORK_SMALL_ROUTINE = 141
)

// 小程序不做校验逻辑md5
const (
	NO_CHECK_PARAMS = "a90a6bffb70d95562baa43cf7fd6bf18"
)

const (
	GAME        = 1
	APPLICATION = 2
)

// abtest gaid or idfa
const (
	GAIDIDFA_ABTEST_PLAYABLE     = "playable3"
	GAIDIDFA_ABTEST_NEW_CREATIVE = "newCreative"
)

func GetAdtypeFromSnr(snr string) []string {
	switch snr {
	case "banner":
		return []string{"banner"}
	case "interstitial":
		return []string{"full_screen", "appwall", "overlay"}
	case "splash":
		return []string{"full_screen"}
	case "exit":
		return []string{"overlay"}
	case "return":
		return []string{"full_screen"}
	case "default":
		return []string{"full_screen"}
	case "none":
		return []string{"none"}
	}
	return []string{"full_screen"}
}

func GetFsTpl(ori int) []string {
	// 2 => horizontal 1=> vertical
	if ori == 2 {
		return []string{"h_sunlight", "h_crystal", "h_balloon_winter"}
	}
	return []string{"rotation", "dive", "santa", "crystal", "sunlight", "lock"}
}

const (
	SHOW_AD_KEY = "uiKFULDJFFJ84NBKFFDq23"
)

const (
	ABTEST_CREATIVE_MONGO = "2"
	ABTEST_CREATIVE_REDIS = "1"
)

const (
	PUB_XIAOMI = 12361
	PUB_BIGO   = 23357
)

const (
	REDUCE_FILL_ECPMFLOOR_MAX = 1.0e5
	REDUCE_FILL_ECPMFLOOR_MIN = 1.0e-11
)

const (
	CLOUD_NAME_ALI = "ali"
	CLOUD_NAME_AWS = "aws"
	CLOUD_NAME_HW  = "hw"
)

const (
	MOF_TYPE_MORE_OFFER                 = 1
	MOF_TYPE_CLOSE_BUTTON_AD            = 2
	MOF_TYPE_CLOSE_BUTTON_AD_MORE_OFFER = 3 // 关闭广告场景下，以more offer的形式展现
)

const (
	RANDOM_BY_DEVICE = 1
	RANDOM_BY_REQ    = 2
)

const (
	ABTEST_TRUE  = 1
	ABTEST_FALSE = 2
)

const (
	APPSFLYER_TOKEN_PARAMS = "&af_token=[af_token]"
)

const (
	RandSum128 = 128
)

const (
	UnsupportBigTpl     = 0   // 没有切量大模板
	SupportButNotBigTpl = 101 // 切量大模板，但是算法选择走旧逻辑
)

const (
	AdnetStartModeTag = "netg"
	AdxStartModeTag   = "adxg"
)

const (
	StorekitLoad    = 1
	StorekitNotLoad = 2
)

const (
	Material      = "material"
	OfferMaterial = "offer_material"
)

const (
	CreativeCompress = "1"
)

const (
	BeforeCompressPath = "cdn-adn/"
	AfterCompressPath  = "cdn-adn/abtestv2/"
)

const (
	BigTemplateUrlFake = "fakeTemplate"
)

// 记录到请求日志中的abtest tag
const (
	ABTestTagVideoCompress = "v_cp"
	ABTestTagImageCompress = "i_cp"
	ABTestTagIconCompress  = "ic_cp"
)

const (
	V5_ABTEST_V5_V5 = "5_5" //请求的V5接口，返回v5的结构体
	V5_ABTEST_V5_V3 = "5_3" //请求的V5接口，返回v3的结构体
)

const (
	SspProfitDistributionRuleFixedEcpm     = 1
	SspProfitDistributionRuleOnlineApiEcpm = 6
)

const (
	BannerIosSDKVersion     = "6.3.0"
	BannerAndriodSDKVersion = "14.0.0"
)

const (
	MTG_SK_NETWORK_ID       = "KBD757YWX3.skadnetwork"
	LOWER_MTG_SK_NETWORK_ID = "kbd757ywx3.skadnetwork"
)

const (
	REQUEST_PIONEER = "b"
)

const AdSpaceTypeFullScreen = 1
const AdSpaceTypeHalfScreen = 2

const MaterialTypeVideoEndcard = 0
const MaterialTypeEndcard = 1
const MaterialTypeVideo = 2

const (
	TEMPLATE_TYPE_TYPE_VIDEO   = "1"
	TEMPLATE_TYPE_TYPE_ENDCARD = "2"
	TEMPLATE_TYPE_TYPE_CAMTPL  = "3"
	TEMPLATE_TYPE_TYPE_PAUSE   = "4"
)

const (
	ABTEST_TEST_GROUP_B = "b"
)

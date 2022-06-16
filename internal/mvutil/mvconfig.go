package mvutil

import (
	smodel "gitlab.mobvista.com/ADN/structs/model"
)

type MVConfig struct {
	Key     string      `bson:"key,omitempty" json:"key"`
	Value   interface{} `bson:"value,omitempty" json:"value"`
	Updated int64       `bson:"updated,omitempty" json:"updated"`
}

type VersionCompare struct {
	Android VersionCompareItem `bson:"android,omitempty" json:"android"`
	IOS     VersionCompareItem `bson:"ios,omitempty" json:"ios"`
}

type VersionCompareItem struct {
	Version        int32   `bson:"version,omitempty" json:"version"`
	ExcludeVersion []int32 `bson:"exclude_version,omitempty" json:"exclude_version"`
}

type PlatableTest struct {
	Url  string      `bson:"url,omitempty" json:"url"`
	Rate map[int]int `bson:"rate,omitempty" json:"rate"`
}

type TRACK_URL_CONFIG_NEW struct {
	Android string `bson:"android,omitempty" json:"android"`
	IOS     string `bson:"ios,omitempty" json:"ios"`
}

type CONFIG_3S_CHINA_DOMAIN struct {
	Domains  []string `bson:"domains,omitempty" json:"domains"`
	CNDomain string   `bson:"cndomain,omitempty" json:"cndomain"`
	Countrys []string `bson:"countrys,omitempty" json:"countrys"`
	Units    []string `bson:"units,omitempty" json:"units"`
	Rate     int      `bson:"rate,omitempty" json:"rate"`
	CNLineDo string   `bson:"cnlinedo,omitempty" json:"cnlinedo"`
}

type ADSTACKING struct {
	Android string `bson:"platform_android,omitempty" json:"platform_android"`
	IOS     string `bson:"platform_ios,omitempty" json:"platform_ios"`
}

type Template struct {
	UnitSize    string `bson:"unitSize,omitempty" json:"unitSize"`
	ImageSize   string `bson:"imageSize,omitempty" json:"imageSize"`
	ImageSizeID int    `bson:"imageSizeId,omitempty" json:"imageSizeId"`
}

type COfferwallUrls struct {
	HTTP  COfferwallUrls_ `bson:"http,omitempty" json:"http"`
	HTTPS COfferwallUrls_ `bson:"https,omitempty" json:"https"`
}

type COfferwallUrls_ struct {
	Offerwall              string `bson:"offerwall,omitempty" json:"offerwall"`
	Interstitial_sdk       string `bson:"interstitial_sdk,omitempty" json:"interstitial_sdk"`
	End_screen             string `bson:"end_screen,omitempty" json:"end_screen"`
	Rewardvideo_end_screen string `bson:"rewardvideo_end_screen,omitempty" json:"rewardvideo_end_screen"`
}

type MP_MAP_UNIT_ struct {
	UnitID      int64 `bson:"unid,omitempty" json:"unid"`
	AppID       int64 `bson:"aid,omitempty" json:"aid"`
	PublisherID int64 `bson:"pid,omitempty" json:"pid"`
}

type SETTING_CONFIG struct {
	UPAL   int64  `bson:"upal,omitempty" json:"upal"`
	CFC    int    `bson:"cfc,omitempty" json:"cfc"`
	GETPF  int64  `bson:"getpf,omitempty" json:"getpf"`
	UPLC   int    `bson:"uplc,omitempty" json:"uplc"`
	FCB    bool   `bson:"cfb,omitempty" json:"cfb"`
	CCT    int    `bson:"cct,omitempty" json:"cct"`
	PCS    bool   `bson:"pcs,omitempty" json:"pcs"`
	PCT    bool   `bson:"pct,omitempty" json:"pct"`
	MO     int    `bson:"mo,omitempty" json:"mo"`
	TTCT   int    `bson:"ttct,omitempty" json:"ttct"`
	AWTTCT int    `bson:"awttct,omitempty" json:"awttct"`
	DMCR   int64  `bson:"dmcr,omitempty" json:"dmcr"`
	NCDN   string `bson:"ncdn,omitempty" json:"ncdn"`
	JABTS  bool   `bson:"jabts,omitempty" json:"jabts"`
	JT3R   int    `bson:"jt3r,omitempty" json:"jt3r"`
	ACCT   int    `bson:"acct,omitempty" json:"acct"` // 安卓的cct值（点击缓存时间）
}

type NO_CHECK_PARAM_APP struct {
	Status *int     `bson:"status,omitempty" json:"status"`
	AppIds *[]int64 `bson:"appId,omitempty" json:"appId"`
}

type ONLINE_APPID_UNITID struct {
	AppId  int64 `bson:"appId,omitempty" json:"appId"`
	UnitId int64 `bson:"unitId,omitempty" json:"unitId"`
}

type REQUEST_BLACKLIST struct {
	AppIds       *[]int64  `bson:"apps,omitempty" json:"apps"`
	DeviceModels *[]string `bson:"device_models,omitempty" json:"device_models"`
	Countries    *[]string `bson:"countries,omitempty" json:"countries"`
}

type SSTestTracfficRule struct {
	CampaignId int `bson:"campaignId"` // 切量后的campaign ID
	ABRatio    int `bson:"abRatio"`    // 分母为10000, 切量10%即设置值为1000
	AARatio    int `bson:"aaRatio"`    // 分母为10000, 切量10%即设置值为1000
}

type OFFER_PLCT struct {
	AppIds []int64 `bson:"appIds"` // specific app
	Plct   int     `bson:"plct"`   // 有效缓存时间，单位s
	Plctb  int     `bson:"plctb"`  // 备用缓存时间，单位s
}

type PUB_ADN_OSV struct {
	PubIds []int64  `bson:"pubIds"`
	Osvs   []string `bson:"osv" json:"osv"`
}

// 降填充开关
type REDUCE_FILL_SWITCH struct {
	Total     int            `bson:"total"`     // 总开关 0.关 1.开
	WhiteList map[string]int `bson:"whitelist"` // 白名单,只要key存在就在白名单内
}

type REDUCE_FILL_FLINK_SWITCH struct {
	Total     int            `bson:"total"`     // 总开关 0.关 1.开
	WhiteList map[string]int `bson:"whitelist"` // 白名单,只要key存在就在白名单内
}

type MORE_OFFER_CONFIG struct {
	TotalRate int     `bson:"totalRate"`
	UnitIds   []int64 `bson:"unitIds"`
}

type BLACK_FOR_EXCLUDE_PACKAGE_NAME struct {
	Status       bool    `bson:"status"`
	PublisherIds []int64 `bson:"publisherIds"`
	AppIds       []int64 `bson:"appIds"`
	Durations    int64   `bson:"durations,omitempty"`
}

type CleanEmptyDeviceABTest struct {
	Status            bool                        `bson:"status"`
	Rate              int                         `bson:"rate"`            // 清空设备abtest 比例
	BlackThirdParty   []string                    `bson:"blackThirdParty"` // 三方黑名单
	WhiteThirdParty   []string                    `bson:"whiteThirdParty"` // 三方白名单
	ThirdPartyInjects map[string]ThirdPartyInject `bson:"thirdPartyInjects"`
}

type ThirdPartyInject struct {
	InjectParams string         `bson:"injectParams"`
	CleanDevId   bool           `bson:"cleanDevId"`
	SubRate      map[string]int `bson:"subRate"` // action -> rate
}

type ExcludeDisplayPackageABTest struct {
	Status            bool           `bson:"status"`
	BlackPublisherIds []int64        `bson:"blackPublisherIds"` // 开发者黑名单
	BlackAppIds       []int64        `bson:"blackAppIds"`       // APP黑名单
	BlackUnitIds      []int64        `bson:"blackUnitIds"`      // Unit黑名单
	WhitePublisherIds []int64        `bson:"whitePublisherIds"` // 开发者白名单
	WhiteAppIds       []int64        `bson:"whiteAppIds"`       // APP白名单
	WhiteUnitIds      []int64        `bson:"whiteUnitIds"`      // Unit白名单
	Rate              map[string]int `bson:"rate"`              // 1 2 为对照组、3 为实验组
}

type REQ_TYPE_AAB_TEST_CONFIG struct {
	Status       bool    `bson:"status"`
	PublisherIds []int64 `bson:"publisherIds"`
	AppIds       []int64 `bson:"appIds"`
	UnitIds      []int64 `bson:"unitIds"`
	Rate         int     `bson:"rate"`
}

type OnlineEmptyDeviceNoServerJump struct {
	Status          bool     `bson:"status,omitempty" json:"status,omitempty"`
	ThirdParty      []string `bson:"thirdParty,omitempty" json:"thirdParty,omitempty"` // 针对三方生效
	BlackThirdParty []string `bson:"blackThirdParty,omitempty" json:"blackThirdParty,omitempty"`
}

type OnlineEmptyDeviceIPUA struct {
	Status bool `bson:"status,omitempty" json:"status,omitempty"`
	// ThirdParty []string `bson:"thirdParty,omitempty" json:"thirdParty,omitempty"` // 针对三方生效
}

type AUTO_LOAD_CACHE_ABTSET_CONFIG struct {
	Status  bool           `bson:"status"`
	UnitIds []int64        `bson:"unitIds"`
	Rate    map[string]int `bson:"adnetRate,omitempty" json:"adnetRate,omitempty"`
}

type ENDCARD_BLOCK_MORE_OFFER_CONFIG struct {
	Status  bool    `bson:"status"`
	UnitIds []int64 `bson:"unitIds"`
}

type IOS_TC_DELAY_CONFIG struct {
	BlackAppList []int64        `bson:"blackAppList,omitempty" json:"blackAppList,omitempty"`
	BlackPkgList []string       `bson:"blackPkgList,omitempty" json:"blackPkgList,omitempty"`
	WhiteAppList []int64        `bson:"whiteAppList,omitempty" json:"whiteAppList,omitempty"`
	WhitePkgList []string       `bson:"whitePkgList,omitempty" json:"whitePkgList,omitempty"`
	OldVerRate   int            `bson:"oldVerRate,omitempty" json:"oldVerRate,omitempty"`
	NewVerRate   int            `bson:"newVerRate,omitempty" json:"newVerRate,omitempty"`
	OldVerConf   map[string]int `bson:"oldVerConf,omitempty" json:"oldVerConf,omitempty"`
	NewVerConf   map[string]int `bson:"newVerConf,omitempty" json:"newVerConf,omitempty"`
}

type TCReduplicateConfig struct {
	Status         bool  `bson:"status,omitempty" json:"status,omitempty"`
	Platform       []int `bson:"platform,omitempty" json:"platform,omitempty"`
	Rate           int   `bson:"rate,omitempty" json:"rate,omitempty"`
	Durations      int64 `bson:"durations,omitempty" json:"durations,omitempty"`
	Type4Durations int64 `bson:"type4Durations,omitempty" json:"type4Durations,omitempty"`
}

type PsbReduplicateConfig struct {
	Status   int      `bson:"status,omitempty" json:"status,omitempty"`
	Platform []int    `bson:"platform,omitempty" json:"platform,omitempty"`
	Rate     int      `bson:"rate,omitempty" json:"rate,omitempty"`
	BlackPkg []string `bson:"blackPkg,omitempty" json:"blackPkg,omitempty"` // 黑名单，加入黑名单之后，注意确认是否需要清理老数据
}

type CreativeBTConfig struct {
	Status        bool    `bson:"status,omitempty" json:"status,omitempty"`
	AdvertiserIds []int64 `bson:"advertiserIds,omitempty" json:"advertiserIds,omitempty"`
	CampaignIds   []int64 `bson:"campaignIds,omitempty" json:"campaignIds,omitempty"`
	Rate          int     `bson:"rate,omitempty" json:"rate,omitempty"`
}

type ABTEST_FIELDS struct {
	AdParams          []string `bson:"adParams,omitempty" json:"adParams,omitempty"`
	EndcardUrlParams  []string `bson:"ecUrlParams,omitempty" json:"ecUrlParams,omitempty"`
	BigTemplateParams []string `bson:"bigTemplateParams,omitempty" json:"bigTemplateParams,omitempty"`
}

type ABTEST_CONF struct {
	CIdWList        []int64        `bson:"cidWList,omitempty" json:"cidWList,omitempty"`
	CIdBList        []int64        `bson:"cidBList,omitempty" json:"cidBList,omitempty"`
	UnitWList       []int64        `bson:"unitWList,omitempty" json:"unitWList,omitempty"`
	UnitBList       []int64        `bson:"unitBList,omitempty" json:"unitBList,omitempty"`
	AppWList        []int64        `bson:"appWList,omitempty" json:"appWList,omitempty"`
	AppBList        []int64        `bson:"appBList,omitempty" json:"appBList,omitempty"`
	PubWList        []int64        `bson:"pubWList,omitempty" json:"pubWList,omitempty"`
	PubBList        []int64        `bson:"pubBList,omitempty" json:"pubBList,omitempty"`
	MaxSdkVer       int32          `bson:"maxSdkVer,omitempty" json:"maxSdkVer,omitempty"`
	MinSdkVer       int32          `bson:"minSdkVer,omitempty" json:"minSdkVer,omitempty"`
	SdkVerBList     []int32        `bson:"sdkVerBList,omitempty" json:"sdkVerBList,omitempty"`
	CountryWList    []string       `bson:"countryWList,omitempty" json:"countryWList,omitempty"`
	CountryBList    []string       `bson:"countryBList,omitempty" json:"countryBList,omitempty"`
	AdTypeWList     []int32        `bson:"adTypeWList,omitempty" json:"adTypeWList,omitempty"`
	AdTypeBList     []int32        `bson:"adTypeBList,omitempty" json:"adTypeBList,omitempty"`
	PlatformWList   []int          `bson:"platformWList,omitempty" json:"platformWList,omitempty"`
	PlatformBList   []int          `bson:"platformBList,omitempty" json:"platformBList,omitempty"`
	ApiVerWList     []int32        `bson:"apiVerWList,omitempty" json:"apiVerWList,omitempty"`
	ApiVerBList     []int32        `bson:"apiVerBList,omitempty" json:"apiVerBList,omitempty"`
	ThirdPartyWList []string       `bson:"thirdPartyWList,omitempty" json:"thirdPartyWList,omitempty"`
	ThirdPartyBList []string       `bson:"thirdPartyBList,omitempty" json:"thirdPartyBList,omitempty"`
	MaxOsVer        int32          `bson:"maxOsVer,omitempty" json:"maxOsVer,omitempty"`
	MinOsVer        int32          `bson:"minOsVer,omitempty" json:"minOsVer,omitempty"`
	OsVerBList      []int32        `bson:"osVerBList,omitempty" json:"osVerBList,omitempty"`
	DspIdWList      []int64        `bson:"dspIdWList,omitempty" json:"dspIdWList,omitempty"`
	DspIdBList      []int64        `bson:"dspIdBList,omitempty" json:"dspIdBList,omitempty"`
	TotalRate       int            `bson:"totalRate,omitempty" json:"totalRate,omitempty"`
	RandType        int32          `bson:"randType,omitempty" json:"randType,omitempty"`
	RateMap         map[string]int `bson:"rateMap,omitempty" json:"rateMap,omitempty"`
}

type ChetLinkConfigItem struct {
	Condition *ChetLinkCondition `bson:"condition,omitempty" json:"condition,omitempty"`
	Link      string             `bson:"link,omitempty" json:"link,omitempty"`
}

type ChetLinkCondition struct {
	IOSSDKVersions     []int32  `bson:"iOSSdkVersions,omitempty" json:"iOSSdkVersions,omitempty"`
	AndroidSDKVersions []int32  `bson:"androidSdkVersions,omitempty" json:"androidSdkVersions,omitempty"`
	CountryCodes       []string `bson:"countryCodes,omitempty" json:"countryCodes,omitempty"`
	AppIds             []int64  `bson:"appIds,omitempty" json:"appIds,omitempty"`
	UnitIds            []int64  `bson:"unitIds,omitempty" json:"unitIds,omitempty"`
}

type AndroidLowVersionFilterCondition struct {
	PublisherIds []int64 `bson:"publisherIds,omitempty" json:"publisherIds,omitempty"`
	AppIds       []int64 `bson:"appIds,omitempty" json:"appIds,omitempty"`
	UnitIds      []int64 `bson:"unitIds,omitempty" json:"unitIds,omitempty"`
}

type PathThirdPartyWhiteList struct {
	CampaignIds []int64 `bson:"campaignIds,omitempty" json:"campaignIds,omitempty"`
	Rate        int     `bson:"rate,omitempty" json:"rate,omitempty"`
}

type ThirdPartyWhiteList struct {
	Status       bool                                `bson:"status,omitempty" json:"status,omitempty"`
	PublisherIds []int64                             `bson:"publisherIds,omitempty" json:"publisherIds,omitempty"`
	AppIds       []int64                             `bson:"appIds,omitempty" json:"appIds,omitempty"`
	UnitIds      []int64                             `bson:"unitIds,omitempty" json:"unitIds,omitempty"`
	Platform     []string                            `bson:"platform,omitempty" json:"platform,omitempty"`
	CountryCode  []string                            `bson:"countryCode,omitempty" json:"countryCode,omitempty"`
	JSSDKRate    map[string]*PathThirdPartyWhiteList `bson:"jssdkRate,omitempty" json:"jssdkRate,omitempty"`
}

type StorekitTimePrisonConf struct {
	MinOsVer        int32 `bson:"minOsVer,omitempty" json:"minOsVer,omitempty"`
	MaxOsVer        int32 `bson:"maxOsVer,omitempty" json:"maxOsVer,omitempty"`
	MinSdkVer       int32 `bson:"minSdkVer,omitempty" json:"minSdkVer,omitempty"`
	MaxSdkVer       int32 `bson:"maxSdkVer,omitempty" json:"maxSdkVer,omitempty"`
	StorekitTimeVal int32 `bson:"storekitTime,omitempty" json:"storekitTime,omitempty"`
}

type IsReturnPauseModConf struct {
	TotalRate int     `bson:"totalRate,omitempty" json:"totalRate,omitempty"`
	UnitIds   []int64 `bson:"unitIds,omitempty" json:"unitIds,omitempty"`
}

// type ClickInServerConf struct {
//	PubBList       []int64                   `bson:"pubIdBList,omitempty" json:"pubIdBList,omitempty"`
//	AppBList       []int64                   `bson:"appIdBList,omitempty" json:"appIdBList,omitempty"`
//	UnitBList      []int64                   `bson:"unitIdBList,omitempty" json:"unitIdBList,omitempty"`
//	OfferBList     []int64                   `bson:"offerIdBList,omitempty" json:"offerIdBList,omitempty"`
//	ClickModeBList []string                  `bson:"clickModeBList,omitempty" json:"clickModeBList,omitempty"`
//	HasDevIdConf   map[string]map[string]int `bson:"hasDevIdConf,omitempty" json:"hasDevIdConf,omitempty"`
//	NoDevIdConf    map[string]map[string]int `bson:"noDevIdConf,omitempty" json:"noDevIdConf,omitempty"`
// }

type TagRate struct {
	Tag  string `bson:"tag,omitempty" json:"tag,omitempty"`
	Rate int    `bson:"rate,omitempty" json:"rate,omitempty"`
	Ext1 string `bson:"ext1,omitempty" json:"ext1,omitempty"`
}

func (t TagRate) GetRate() int {
	return t.Rate
}

type ExcludeClickPackagesV2 struct {
	Configs []*ExcludeClickPackages `bson:"configs,omitempty" json:"configs,omitempty"`
}

type ExcludeClickPackages struct {
	Status           bool                  `bson:"status,omitempty" json:"status,omitempty"`
	Platform         []int                 `bson:"platform,omitempty" json:"platform,omitempty"`
	PubBList         []int64               `bson:"pubBList,omitempty" json:"pubBList,omitempty"`
	AppBList         []int64               `bson:"appBList,omitempty" json:"appBList,omitempty"`
	AdTypeBList      []int32               `bson:"adTypeBList,omitempty" json:"adTypeBList,omitempty"`
	PackageBList     []string              `bson:"packageBList,omitempty" json:"packageBList,omitempty"`
	PubWList         []int64               `bson:"pubWList,omitempty" json:"pubWList,omitempty"`
	AppWList         []int64               `bson:"appWList,omitempty" json:"appWList,omitempty"`
	AdTypeWList      []int32               `bson:"adTypeWList,omitempty" json:"adTypeWList,omitempty"`
	PackageWList     []string              `bson:"packageWList,omitempty" json:"packageWList,omitempty"`
	ControlGroupTime int                   `bson:"controlGroupTime,omitempty" json:"controlGroupTime,omitempty"` // 对照组 命中的时间限制
	TagRates         map[string]*TagRate   `bson:"rate,omitempty" json:"rate,omitempty"`                         // 切量，key为时间，单位为s，value 为切量比例。key为0 即不做过滤
	SubTagRates      map[string][]*TagRate `bson:"bbRate,omitempty" json:"bbRate,omitempty"`                     // 切量实验组可以通过该配置，配置比例不做实验，起着BB对照.
	rate             map[string]IRate
	subRate          map[string]IRate
}

func (e ExcludeClickPackages) GetRates() map[string]IRate {
	if e.rate != nil {
		return e.rate
	}

	rate := make(map[string]IRate)
	if len(e.TagRates) != 0 {
		for key, val := range e.TagRates {
			rate[key] = val
		}
	}

	e.rate = rate
	return e.rate
}

func (e ExcludeClickPackages) GetSubRates(key string) map[string]IRate {
	var res map[string]IRate
	if key == "" || e.SubTagRates == nil {
		return res
	}

	tags, ok := e.SubTagRates[key]
	if !ok {
		return res
	}

	res = make(map[string]IRate)
	for _, val := range tags {
		res[val.Tag] = val
	}
	return res
}

type AdChoiceConfig struct {
	AppIds      []int64        `bson:"appIds,omitempty" json:"appIds,omitempty"`
	UnitIds     []int64        `bson:"unitIds,omitempty" json:"unitIds,omitempty"`
	SdkVersions InfoSDKVersion `bson:"sdkVersions,omitempty" json:"sdkVersions,omitempty"`
}

type AdChoiceConfigData struct {
	AppIds             []int64
	UnitIds            []int64
	IOSSDKVersions     *DataSDKVersion
	AndroidSDKVersions *DataSDKVersion
}

type ClickInServerConf struct {
	RandDev        bool                    `bson:"randNum,omitempty" json:"randNum,omitempty"`
	PubBList       []int64                 `bson:"pubIdBList,omitempty" json:"pubIdBList,omitempty"`
	AppBList       []int64                 `bson:"appIdBList,omitempty" json:"appIdBList,omitempty"`
	UnitBList      []int64                 `bson:"unitIdBList,omitempty" json:"unitIdBList,omitempty"`
	OfferBList     []int64                 `bson:"offerIdBList,omitempty" json:"offerIdBList,omitempty"`
	ClickModeBList []string                `bson:"clickModeBList,omitempty" json:"clickModeBList,omitempty"`
	OpenTypeBList  []int                   `bson:"openTypeBList,omitempty" json:"openTypeBList,omitempty"`
	HasDevIdConf   []ClickInServerRateConf `bson:"hasDevIdConf,omitempty" json:"hasDevIdConf,omitempty"`
	NoDevIdConf    []ClickInServerRateConf `bson:"noDevIdConf,omitempty" json:"noDevIdConf,omitempty"`
}

type ClickInServerRateConf struct {
	Region          []string `bson:"reWList,omitempty" json:"reWList,omitempty"`
	BlackRegion     []string `bson:"reBList,omitempty" json:"reBList,omitempty"`
	Country         []string `bson:"ccWList,omitempty" json:"ccWList,omitempty"`
	BlackCountry    []string `bson:"ccBList,omitempty" json:"ccBList,omitempty"`
	ThirdParty      []string `bson:"tpWList,omitempty" json:"tpWList,omitempty"`
	BlackThirdParty []string `bson:"tpBList,omitempty" json:"tpBList,omitempty"`
	Platform        []int    `bson:"plWList,omitempty" json:"plWList,omitempty"`
	Rate            int      `bson:"rate,omitempty" json:"rate,omitempty"`
}

type MoattagConfig struct {
	SdkVersions InfoSDKVersion `bson:"sdkVersions,omitempty" json:"sdkVersions,omitempty"`
}

type MoattagConfigData struct {
	IOSSDKVersions     *DataSDKVersion
	AndroidSDKVersions *DataSDKVersion
}

type BadRequestFilterConf struct {
	UnitWList           []int64  `bson:"unitWList,omitempty" json:"unitWList,omitempty"`
	AppWList            []int64  `bson:"appWList,omitempty" json:"appWList,omitempty"`
	PubWList            []int64  `bson:"pubWList,omitempty" json:"pubWList,omitempty"`
	MaxSdkVer           int32    `bson:"maxSdkVer,omitempty" json:"maxSdkVer,omitempty"`
	MinSdkVer           int32    `bson:"minSdkVer,omitempty" json:"minSdkVer,omitempty"`
	CountryWList        []string `bson:"countryWList,omitempty" json:"countryWList,omitempty"`
	AdTypeWList         []int32  `bson:"adTypeWList,omitempty" json:"adTypeWList,omitempty"`
	PlatformWList       []int    `bson:"platformWList,omitempty" json:"platformWList,omitempty"`
	ApiVerWList         []int32  `bson:"apiVerWList,omitempty" json:"apiVerWList,omitempty"`
	MaxOsVer            int32    `bson:"maxOsVer,omitempty" json:"maxOsVer,omitempty"`
	MinOsVer            int32    `bson:"minOsVer,omitempty" json:"minOsVer,omitempty"`
	DeviceTypeWList     []int    `bson:"deviceTypeWList,omitempty" json:"deviceTypeWList,omitempty"`
	ModelWList          []string `bson:"modelWList,omitempty" json:"modelWList,omitempty"`
	AppVersionCodeList  []string `bson:"appVersionCodeList,omitempty" json:"appVersionCodeList,omitempty"`
	AppVersionNameList  []string `bson:"appVersionNameList,omitempty" json:"appVersionNameList,omitempty"`
	HModelWList         []string `bson:"hModelWList,omitempty" json:"hModelWList,omitempty"`
	HModelMatchingWList []string `bson:"hModelMatchingWList,omitempty" json:"hModelMatchingWList,omitempty"`
}

type BigTemplateConf struct {
	UnitWList map[string]int `bson:"unitWList,omitempty" json:"unitWList,omitempty"`
	AppWList  map[string]int `bson:"appWList,omitempty" json:"appWList,omitempty"`
	PubWList  map[string]int `bson:"pubWList,omitempty" json:"pubWList,omitempty"`
	UnitBList []int64        `bson:"unitBList,omitempty" json:"unitBList,omitempty"`
	AppBList  []int64        `bson:"appBList,omitempty" json:"appBList,omitempty"`
	PubBList  []int64        `bson:"pubBList,omitempty" json:"pubBList,omitempty"`
	SupAdType []int32        `bson:"supAdType,omitempty" json:"supAdType,omitempty"`
	TotalRate int            `bson:"totalRate,omitempty" json:"totalRate,omitempty"`
}
type CNTrackingDomainTestConf struct {
	UnitBList []int64              `bson:"unitBList,omitempty" json:"unitBList,omitempty"`
	AppBList  []int64              `bson:"appBList,omitempty" json:"appBList,omitempty"`
	PubBList  []int64              `bson:"pubBList,omitempty" json:"pubBList,omitempty"`
	UnitWList []int64              `bson:"unitWList,omitempty" json:"unitWList,omitempty"`
	RandType  int                  `bson:"randType,omitempty" json:"randType,omitempty"`
	Conf      []*smodel.CdnSetting `bson:"conf,omitempty" json:"conf,omitempty"`
	Status    bool                 `bson:"status,omitempty" json:"status,omitempty"` // 整体切量开关
}

type NoticeClickUniq struct {
	Status              bool           `bson:"status,omitempty" json:"status,omitempty"`
	ServerClickUniq     bool           `bson:"serverClickUniq,omitempty" json:"serverClickUniq,omitempty"`
	ServerClickUniqTime map[string]int `bson:"scut,omitempty" json:"scut,omitempty"`
}

// type AdjustPostbackABTest struct {
// 	Status        bool     `bson:"status,omitempty" json:"status,omitempty"`
// 	Rate          int      `bson:"rate,omitempty" json:"rate,omitempty"`         // 切量比例
// 	Scenario      []string `bson:"scenario,omitempty" json:"scenario,omitempty"` // 限制场景
// 	BListCampaign []int64  `bson:"bListCampaign,omitempty" json:"bListCampaign,omitempty"`
// 	WListCampaign []int64  `bson:"wListCampaign,omitempty" json:"wListCampaign,omitempty"`
// }

// 通用切量配置
type TestConf struct {
	UnitWList map[string]int `bson:"unitWList,omitempty" json:"unitWList,omitempty"`
	AppWList  map[string]int `bson:"appWList,omitempty" json:"appWList,omitempty"`
	PubWList  map[string]int `bson:"pubWList,omitempty" json:"pubWList,omitempty"`
	UnitBList []int64        `bson:"unitBList,omitempty" json:"unitBList,omitempty"`
	AppBList  []int64        `bson:"appBList,omitempty" json:"appBList,omitempty"`
	PubBList  []int64        `bson:"pubBList,omitempty" json:"pubBList,omitempty"`
	SupAdType []int32        `bson:"supAdType,omitempty" json:"supAdType,omitempty"`
	TotalRate int            `bson:"totalRate,omitempty" json:"totalRate,omitempty"`
}

// offer模版&大模版
type TemplateMap struct {
	Id     int32  `bson:"id,omitempty" json:"id,omitempty"`
	Url    string `bson:"url,omitempty" json:"url,omitempty"`
	Weight int32  `bson:"weight,omitempty" json:"weight,omitempty"`
}

type FreqControlConfig struct {
	FreqControlToRs int                `bson:"freq_control_to_rs,omitempty" json:"freq_control_to_rs,omitempty"`
	HbReducedRate   float64            `bson:"hb_reduced_rate,omitempty" json:"hb_reduced_rate,omitempty"`
	Status          int                `bson:"status,omitempty" json:"status,omitempty"`
	UseFixedEcpm    int                `bson:"use_fixed_ecpm,omitempty" json:"use_fixed_ecpm,omitempty"`
	FixedEcpmFactor float64            `bson:"fixed_ecpm_factor,omitempty" json:"fixed_ecpm_factor,omitempty"`
	Rules           []*FreqControlRule `bson:"rules,omitempty" json:"rules,omitempty"`
}

type FreqControlRule struct {
	Filter *FreqControlFilter      `bson:"filter,omitempty" json:"filter,omitempty"`
	Groups []*FreqControlGroupItem `bson:"groups,omitempty" json:"groups,omitempty"`
}

type FreqControlFilter struct {
	Platform    *FreqControlFilterItem `bson:"platform,omitempty" json:"platform,omitempty"`
	AdType      *FreqControlFilterItem `bson:"ad_type,omitempty" json:"ad_type,omitempty"`
	CountryCode *FreqControlFilterItem `bson:"country_code,omitempty" json:"country_code,omitempty"`
	PublisherId *FreqControlFilterItem `bson:"publisher_id,omitempty" json:"publisher_id,omitempty"`
	AppId       *FreqControlFilterItem `bson:"app_id,omitempty" json:"app_id,omitempty"`
	UnitId      *FreqControlFilterItem `bson:"unit_id,omitempty" json:"unit_id,omitempty"`
	PlacementId *FreqControlFilterItem `bson:"placement_id,omitempty" json:"placement_id,omitempty"`
	HasDevid    *FreqControlFilterItem `bson:"has_devid,omitempty" json:"has_devid,omitempty"`
	IsHb        *FreqControlFilterItem `bson:"is_hb,omitempty" json:"is_hb,omitempty"`
}

type FreqControlGroupItem struct {
	GroupName      string                                  `bson:"group_name,omitempty" json:"group_name,omitempty"`
	FreqControlKey string                                  `bson:"freq_control_key,omitempty" json:"freq_control_key,omitempty"`
	GroupRate      int                                     `bson:"group_rate,omitempty" json:"group_rate,omitempty"`
	Rate           int                                     `bson:"rate,omitempty" json:"rate,omitempty"`
	SubRate        int                                     `bson:"sub_rate,omitempty" json:"sub_rate,omitempty"`
	TimeWindow     *FreqControlGroupsItemTimeWindow        `bson:"time_window,omitempty" json:"time_window,omitempty"`
	FreqControl    []*FreqControlGroupsItemFreqControlItem `bson:"freq_control,omitempty" json:"freq_control,omitempty"`
}

type FreqControlFilterItem struct {
	Op    string   `bson:"op,omitempty" json:"op,omitempty"`
	Value []string `bson:"value,omitempty" json:"value,omitempty"`
}

type FreqControlGroupsItemTimeWindow struct {
	Mode      int `bson:"mode,omitempty" json:"mode,omitempty"`
	StartHour int `bson:"start_hour,omitempty" json:"start_hour,omitempty"`
	WindowSec int `bson:"window_sec,omitempty" json:"window_sec,omitempty"`
}

type FreqControlGroupsItemFreqControlItem struct {
	Min  int      `bson:"min,omitempty" json:"min,omitempty"`
	Max  int      `bson:"max,omitempty" json:"max,omitempty"`
	Keys []string `bson:"keys,omitempty" json:"keys,omitempty"`
}

type FreqControlMkvDataItem struct {
	Ts   int     `bson:"ts,omitempty" json:"ts,omitempty"`
	Freq float64 `bson:"freq,omitempty" json:"freq,omitempty"`
}

type FreqControlMkvData struct {
	Imp []*FreqControlMkvDataItem `bson:"imp,omitempty" json:"imp,omitempty"`
	Req []*FreqControlMkvDataItem `bson:"req,omitempty" json:"req,omitempty"`
}

type DspTplAbtest struct {
	Id     int `bson:"id,omitempty" json:"id,omitempty"`
	Weight int `bson:"weight,omitempty" json:"weight,omitempty"`
}

type FreqImpCapMkvData struct {
	Ts    int64 `bson:"ts,omitempty" json:"ts,omitempty"`
	Count int   `bson:"count,omitempty" json:"count,omitempty"`
}

type CreativeCompressABTestV3Data struct {
	RateMap       map[string]int `bson:"rateMap,omitempty" json:"rateMap,omitempty"`
	CampaignBList []int64        `bson:"campaignBList,omitempty" json:"campaignBList,omitempty"`
	AppBList      []int64        `bson:"appBList,omitempty" json:"appBList,omitempty"`
	UnitBList     []int64        `bson:"unitBList,omitempty" json:"unitBList,omitempty"`
	PubBList      []int64        `bson:"pubBList,omitempty" json:"pubBList,omitempty"`
}

type TcFilterByCityConf struct {
	CityCodeBlackList    []int64  `bson:"cityCodeBList,omitempty" json:"cityCodeBList,omitempty"`
	CityCodeWlackList    []int64  `bson:"cityCodeWList,omitempty" json:"cityCodeWList,omitempty"`
	CountryCodeBlackList []string `bson:"countryCodeBList,omitempty" json:"countryCodeBList,omitempty"`
}

type HBAdxEndpointMkvData struct {
	Endpoint string `bson:"endpoint,omitempty" json:"endpoint,omitempty"`
	Timeout  int    `bson:"timeout,omitempty" json:"timeout,omitempty"`
}

type OnlineApiPubBidPriceConf struct {
	PubConf  map[string]float64 `bson:"publisher,omitempty" json:"publisher,omitempty"`
	AppConf  map[string]float64 `bson:"app,omitempty" json:"app,omitempty"`
	UnitConf map[string]float64 `bson:"unit,omitempty" json:"unit,omitempty"`
}

type SupportTrackingTemplateConf struct {
	Status      bool     `bson:"status,omitempty" json:"status,omitempty"`
	BAppIds     []int64  `bson:"bappIds,omitempty" json:"bappIds,omitempty"`
	BSDKVersion []string `bson:"bsdkVersion,omitempty" json:"bsdkVersion,omitempty"`
	BAPIVersion []int32  `bson:"bapiVersion,omitempty" json:"bapiVersion,omitempty"`
}

type V5AbtestConf struct {
	Switch   int           `json:"switch,omitempty"` //0=关闭，1=全部开启，2=实验，符合unit_conf配置才返回v5结构
	UnitConf map[int64]int `json:"unit_conf,omitempty"`
}

type ReplaceTrackingDomainConf struct {
	PubConfs  map[string]map[string]*TrackingDomainActionConf `bson:"pubConfs,omitempty" json:"pubConfs,omitempty"`
	AppConfs  map[string]map[string]*TrackingDomainActionConf `bson:"appConfs,omitempty" json:"appConfs,omitempty"`
	UnitConfs map[string]map[string]*TrackingDomainActionConf `bson:"unitConfs,omitempty" json:"unitConfs,omitempty"`
	TotalConf map[string]*TrackingDomainActionConf            `bson:"totalConf,omitempty" json:"totalConf,omitempty"`
}

type TrackingDomainActionConf struct {
	ImpressionTrackConf []*TrackingDomainWeightConf `bson:"imp,omitempty" json:"imp,omitempty"`
	ClickTrackConf      []*TrackingDomainWeightConf `bson:"click,omitempty" json:"click,omitempty"`
}

type TrackingDomainWeightConf struct {
	Id     int    `bson:"id,omitempty" json:"id,omitempty"`
	Domain string `bson:"domain,omitempty" json:"domain,omitempty"`
	Weight int    `bson:"weight,omitempty" json:"weight,omitempty"`
}

type ShareitAfWhiteListConf struct {
	AllowPackage                     []string          `bson:"allowPackages,omitempty" json:"allowPackages,omitempty"`
	UnitConfs                        map[string]int    `bson:"unitConfs,omitempty" json:"unitConfs,omitempty"`
	CloseShareitAfWhiteListQfLossbid bool              `bson:"closeShareitAfWhiteListQfLossbid,omitempty" json:"closeShareitAfWhiteListQfLossbid,omitempty"`
	IsMtgPidValue                    string            `bson:"isMtgPidVal,omitempty" json:"isMtgPidVal,omitempty"`       // 是否mtg
	TotalStatus                      bool              `bson:"totalStatus,omitempty" json:"totalStatus,omitempty"`       // 整体开关
	PackageNameMap                   map[string]string `bson:"packageNameMap,omitempty" json:"packageNameMap,omitempty"` // qcc接口流量测包名
}

type OnlineApiMaxBidPrice struct {
	TotalBidPrice float64            `bson:"totalBidPrice,omitempty" json:"totalBidPrice,omitempty"`
	PubBidPrice   map[string]float64 `bson:"pubBidPrice,omitempty" json:"pubBidPrice,omitempty"`
	AppBidPrice   map[string]float64 `bson:"appBidPrice,omitempty" json:"appBidPrice,omitempty"`
	UnitBidPrice  map[string]float64 `bson:"unitBidPrice,omitempty" json:"unitBidPrice,omitempty"`
}

type DebugBidFloorAndBidPriceConf struct {
	DebugBidFloor float64 `bson:"debugBidFloor,omitempty" json:"debugBidFloor,omitempty"`
	DebugBidPrice float64 `bson:"debugBidPrice,omitempty" json:"debugBidPrice,omitempty"`
}

type H265VideoABTestConf struct {
	UnitConf                       map[string]map[string]int `bson:"unitConf,omitempty" json:"unitConf,omitempty"`
	AppConf                        map[string]map[string]int `bson:"appConf,omitempty" json:"appConf,omitempty"`
	AdTypeConf                     map[string]map[string]int `bson:"adTypeConf,omitempty" json:"adTypeConf,omitempty"`
	TotalRate                      map[string]int            `bson:"totalRate,omitempty" json:"totalRate,omitempty"`
	IosBlackSdkVersionCodeList     []int32                   `bson:"iosBlackSdkVersionCodeList,omitempty" json:"iosBlackSdkVersionCodeList,omitempty"`
	AndroidBlackSdkVersionCodeList []int32                   `bson:"androidBlackSdkVersionCodeList,omitempty" json:"androidBlackSdkVersionCodeList,omitempty"`
	IosBlackOsVersionCodeList      []int32                   `bson:"iosBlackOsVersionCodeList,omitempty" json:"iosBlackOsVersionCodeList,omitempty"`
	AndroidBlackOsVersionCodeList  []int32                   `bson:"androidBlackOsVersionCodeList,omitempty" json:"androidBlackOsVersionCodeList,omitempty"`
	MinAllowIosOsVersionCode       int32                     `bson:"minAllowIosOsVersionCode,omitempty" json:"minAllowIosOsVersionCode,omitempty"`
	MinAllowAndroidOsVersionCode   int32                     `bson:"minAllowAndroidOsVersionCode,omitempty" json:"minAllowAndroidOsVersionCode,omitempty"`
	MinAllowIosSdkVersionCode      int32                     `bson:"minAllowIosSdkVersionCode,omitempty" json:"minAllowIosSdkVersionCode,omitempty"`
	MinAllowAndroidSdkVersionCode  int32                     `bson:"minAllowAndroidSdkVersionCode,omitempty" json:"minAllowAndroidSdkVersionCode,omitempty"`
	BlackModelList                 []string                  `bson:"blackModelList,omitempty" json:"blackModelList,omitempty"`
	BlackModelOsVersionCodeList    []string                  `bson:"blackModelOsVersionCodeList,omitempty" json:"blackModelOsVersionCodeList,omitempty"`
	BrandModelWhiteList            []string                  `bson:"brandModelWhiteList,omitempty" json:"brandModelWhiteList,omitempty"`
}

type TemplateWeightMap struct {
	Id     int    `bson:"id,omitempty" json:"id,omitempty"`
	Url    string `bson:"url,omitempty" json:"url,omitempty"`
	Weight int    `bson:"weight,omitempty" json:"weight,omitempty"`
}

type ReplaceTemplateUrlConf struct {
	Endcard    map[string][]*TemplateWeightMap `bson:"endcard,omitempty" json:"endcard,omitempty"`
	RvTemplate map[string][]*TemplateWeightMap `bson:"rv_template,omitempty" json:"rv_template,omitempty"`
}

type TemplateCreativeDomainMap struct {
	Id     int    `bson:"id,omitempty" json:"id,omitempty"`
	Domain string `bson:"domain,omitempty" json:"domain,omitempty"`
	Weight int    `bson:"weight,omitempty" json:"weight,omitempty"`
}

type ReturnWtickConf struct {
	UnitConf  map[string]int `bson:"unitConf,omitempty" json:"unitConf,omitempty"`
	TotalConf int            `bson:"totalConf,omitempty" json:"totalConf,omitempty"`
}

type AdPackageNameReplaceConf struct {
	PackageNames map[string]*AdPackageNameReplace
}

type AdPackageNameReplace struct {
	ReplacePackageName string `bson:"replacePackageName,omitempty" json:"replacePackageName,omitempty"`
	Rate               int    `bson:"rate,omitempty" json:"rate,omitempty"`
}

type ExcludePackagesByCityCodeConf struct {
	CityCodeList            []int64  `bson:"cityCodeList,omitempty" json:"cityCodeList,omitempty"`
	BlockIosPackageList     []string `bson:"blockIosPackageList,omitempty" json:"blockIosPackageList,omitempty"`
	BlockAndroidPackageList []string `bson:"blockAndroidPackageList,omitempty" json:"blockAndroidPackageList,omitempty"`
	Status                  bool     `bson:"status,omitempty" json:"status,omitempty"`
}

type FilterByStackConf struct {
	StackName string         `bson:"stackName,omitempty" json:"stackName,omitempty"`
	AppConf   map[string]int `bson:"appConf,omitempty" json:"appConf,omitempty"`
	PubConf   map[string]int `bson:"pubConf,omitempty" json:"pubConf,omitempty"`
	TotalConf int            `bson:"totalConf,omitempty" json:"totalConf,omitempty"`
}

type HBRequestSnapshot struct {
	Enable bool `bson:"enable,omitempty" json:"enable,omitempty"`
	Rate   int  `bson:"rate,omitempty" json:"rate,omitempty"`
}

type TrackingCNABTestConf struct {
	Enable   bool            `bson:"enable,omitempty" json:"enable,omitempty"`
	Unit     TagListRateConf `bson:"unit,omitempty" json:"unit,omitempty"`
	App      TagListRateConf `bson:"app,omitempty" json:"app,omitempty"`
	Pub      TagListRateConf `bson:"pub,omitempty" json:"pub,omitempty"`
	AdType   TagListRateConf `bson:"ad_type,omitempty" json:"ad_type,omitempty"`
	Platform TagListRateConf `bson:"platform,omitempty" json:"platform,omitempty"`
	Rate     int             `bson:"rate,omitempty" json:"rate,omitempty"` //全量比例控制
	Domain   string          `bson:"domain,omitempty" json:"domain,omitempty"`
	JsDomain string          `bson:"js_domain,omitempty" json:"js_domain,omitempty"`
}

type TagListRateConf struct {
	Rate    int     `bson:"rate,omitempty" json:"rate,omitempty"`
	TagList []int64 `bson:"list,omitempty" json:"list,omitempty"`
}

type MappingIdfaAbtestConf struct {
	AppList              []int64        `bson:"appList,omitempty" json:"appList,omitempty"`
	AppBlackList         []int64        `bson:"appBlackList,omitempty" json:"appBlackList,omitempty"`
	PubList              []int64        `bson:"pubList,omitempty" json:"pubList,omitempty"`
	PubBlackList         []int64        `bson:"pubBlackList,omitempty" json:"pubBlackList,omitempty"`
	CountryCodeList      []string       `bson:"countryCodeList,omitempty" json:"countryCodeList,omitempty"`
	CountryCodeBlackList []string       `bson:"countryCodeBlackList,omitempty" json:"countryCodeBlackList,omitempty"`
	ConfMap              map[string]int `bson:"confMap,omitempty" json:"confMap,omitempty"`
	RandType             int            `bson:"randType,omitempty" json:"randType,omitempty"`
}

type OrientationPoisonConf struct {
	PubList       []int64 `bson:"pubList,omitempty" json:"pubList,omitempty"`
	AppList       []int64 `bson:"appList,omitempty" json:"appList,omitempty"`
	UnitList      []int64 `bson:"unitList,omitempty" json:"unitList,omitempty"`
	AppBlackList  []int64 `bson:"appBlackList,omitempty" json:"appBlackList,omitempty"`
	UnitBlackList []int64 `bson:"unitBlackList,omitempty" json:"unitBlackList,omitempty"`
}

type SupportSmartVBAConfig struct {
	Status          bool               `bson:"status,omitempty" json:"status,omitempty"`
	SmartVBASetting *SmartVBASetting   `bson:"smart_vba_setting,omitempty" json:"smart_vba_setting,omitempty"`
	Items           []*SupportSmartVBA `bson:"items,omitempty" json:"items,omitempty"`
}

type SupportSmartVBA struct {
	WTick              *SmartVBAWTick               `bson:"wtick,omitempty" json:"wtick,omitempty"`
	ReplacePackageName *SmartVBAWReplacePackageName `bson:"replace_package_name,omitempty" json:"replace_package_name,omitempty"`
	SmartVBA           []*SmartVBA                  `bson:"smart_vba,omitempty" json:"smart_vba,omitempty"`
}

type SmartVBAWTick struct {
	Rate     int               `bson:"rate,omitempty" json:"rate,omitempty"`
	Selector *SmartVBASelector `bson:"selector,omitempty" json:"selector,omitempty"`
}

type SmartVBAWReplacePackageName struct {
	Rate               int               `bson:"rate,omitempty" json:"rate,omitempty"`
	Selector           *SmartVBASelector `bson:"selector,omitempty" json:"selector,omitempty"`
	ReplacePackageName string            `bson:"replace_package_name,omitempty" json:"replace_package_name,omitempty"`
}

type SmartVBA struct {
	Rate            int               `bson:"rate,omitempty" json:"rate,omitempty"`
	Tag             string            `bson:"tag,omitempty" json:"tag,omitempty"`
	Selector        *SmartVBASelector `bson:"selector,omitempty" json:"selector,omitempty"`
	Action          string            `bson:"action,omitempty" json:"action,omitempty"`
	Key             string            `bson:"key,omitempty" json:"key,omitempty"`
	PlayRate        int               `bson:"play_rate,omitempty" json:"play_rate,omitempty"`
	DontFakeImp     bool              `bson:"dont_fake_imp,omitempty" json:"dont_fake_imp,omitempty"`
	ClientClick     bool              `bson:"client_click,omitempty" json:"client_click,omitempty"`
	MinDuration     int               `bson:"min_duration,omitempty" json:"min_duration,omitempty"`
	MaxDuration     int               `bson:"max_duration,omitempty" json:"max_duration,omitempty"`
	ClickDelayAlpha float64           `bson:"click_delay_alpha,omitempty" json:"click_delay_alpha,omitempty"`
	ClickDelayBeta  float64           `bson:"click_delay_beta,omitempty" json:"click_delay_beta,omitempty"`
	ImpDelayAlpha   float64           `bson:"imp_delay_alpha,omitempty" json:"imp_delay_alpha,omitempty"`
	ImpDelayBeta    float64           `bson:"imp_delay_beta,omitempty" json:"imp_delay_beta,omitempty"`
}

type SmartVBASetting struct {
	MinDuration     int     `bson:"min_duration,omitempty" json:"min_duration,omitempty"`
	MaxDuration     int     `bson:"max_duration,omitempty" json:"max_duration,omitempty"`
	ClickDelayAlpha float64 `bson:"click_delay_alpha,omitempty" json:"click_delay_alpha,omitempty"`
	ClickDelayBeta  float64 `bson:"click_delay_beta,omitempty" json:"click_delay_beta,omitempty"`
	ImpDelayAlpha   float64 `bson:"imp_delay_alpha,omitempty" json:"imp_delay_alpha,omitempty"`
	ImpDelayBeta    float64 `bson:"imp_delay_beta,omitempty" json:"imp_delay_beta,omitempty"`
}

type SmartVBASelector struct {
	IncReplacePackage      []int    `bson:"inc_replace_package,omitempty" json:"inc_replace_package,omitempty"`
	ExcReplacePackage      []int    `bson:"exc_replace_package,omitempty" json:"exc_replace_package,omitempty"`
	IncWtick               []int    `bson:"inc_wtick,omitempty" json:"inc_wtick,omitempty"`
	ExcWtick               []int    `bson:"exc_wtick,omitempty" json:"exc_wtick,omitempty"`
	IncCampaignIds         []int64  `bson:"inc_campaign_ids,omitempty" json:"inc_campaign_ids,omitempty"`
	ExcCampaignIds         []int64  `bson:"exc_campaign_ids,omitempty" json:"exc_campaign_ids,omitempty"`
	IncCampaignPackageName []string `bson:"inc_campaign_package_name,omitempty" json:"inc_campaign_package_name,omitempty"`
	ExcCampaignPackageName []string `bson:"exc_campaign_package_name,omitempty" json:"exc_campaign_package_name,omitempty"`
	IncRequestTypes        []int    `bson:"inc_request_types,omitempty" json:"inc_request_types,omitempty"`
	ExcRequestTypes        []int    `bson:"exc_request_types,omitempty" json:"exc_request_types,omitempty"`
	IncPackageNames        []string `bson:"inc_package_names,omitempty" json:"inc_package_names,omitempty"`
	ExcPackageNames        []string `bson:"exc_package_names,omitempty" json:"exc_package_names,omitempty"`
	IncAdTypes             []int    `bson:"inc_ad_types,omitempty" json:"inc_ad_types,omitempty"`
	ExcAdTypes             []int    `bson:"exc_ad_types,omitempty" json:"exc_ad_types,omitempty"`
	IncPublisherIds        []int64  `bson:"inc_publisher_ids,omitempty" json:"inc_publisher_ids,omitempty"`
	ExcPublisherIds        []int64  `bson:"exc_publisher_ids,omitempty" json:"exc_publisher_ids,omitempty"`
	IncAppIds              []int64  `bson:"inc_app_ids,omitempty" json:"inc_app_ids,omitempty"`
	ExcAppIds              []int64  `bson:"exc_app_ids,omitempty" json:"exc_app_ids,omitempty"`
	IncUnitIds             []int64  `bson:"inc_unit_ids,omitempty" json:"inc_unit_ids,omitempty"`
	ExcUnitIds             []int64  `bson:"exc_unit_ids,omitempty" json:"exc_unit_ids,omitempty"`
}

type MoreOfferAndAppwallMoveToPioneerABTestConf struct {
	PublisherConf map[string]map[string]int `bson:"pubConf,omitempty" json:"pubConf,omitempty"`
	AppConf       map[string]map[string]int `bson:"appConf,omitempty" json:"appConf,omitempty"`
	UnitConf      map[string]map[string]int `bson:"unitConf,omitempty" json:"unitConf,omitempty"`
	TotalRate     map[string]int            `bson:"totalRate,omitempty" json:"totalRate,omitempty"`
	RandType      int                       `bson:"randType,omitempty" json:"randType,omitempty"`
}

type HBOfferBidPriceConf struct {
	UnitList []int64 `bson:"unitList,omitempty" json:"unitList,omitempty"`
	RandType int     `bson:"randType,omitempty" json:"randType,omitempty"`
	Rate     int     `bson:"rate,omitempty" json:"rate,omitempty"`
}

type MediationNoticeURLMacroConfValue struct {
	Configs []*MediationNoticeURLMacroConf `bson:"configs,omitempty" json:"configs,omitempty"`
}

type MediationNoticeURLMacroConf struct {
	ChannelID    string  `bson:"channel_id,omitempty" json:"channel_id,omitempty"`
	UnitIDs      []int64 `bson:"unit_ids,omitempty" json:"unit_ids,omitempty"`
	AppIDs       []int64 `bson:"app_ids,omitempty" json:"app_ids,omitempty"`
	PublisherIDs []int64 `bson:"publisher_ids,omitempty" json:"publisher_ids,omitempty"`
	LURLMacro    string  `bson:"lurl_marco,omitempty" json:"lurl_marco,omitempty"`
}

type TrackDomainByCountryCodeConf struct {
	CityBlackList []int64              `bson:"cityBlackList,omitempty" json:"cityBlackList,omitempty"`
	RandType      int                  `bson:"randType,omitempty" json:"randType,omitempty"`
	ConfMap       []*smodel.CdnSetting `bson:"confMap,omitempty" json:"confMap,omitempty"`
}

type LoadDomainABTest struct {
	Region          []string `bson:"region,omitempty" json:"region,omitempty"`
	Cloud           []string `bson:"cloud,omitempty" json:"cloud,omitempty"`
	Platform        []int    `bson:"platform,omitempty" json:"platform,omitempty"`
	WMediationNames []string `bson:"wMediationNames,omitempty" json:"wMediationNames,omitempty"`
	BMediationNames []string `bson:"bMediationNames,omitempty" json:"bMediationNames,omitempty"`
	WCountryCode    []string `bson:"wCountryCode,omitempty" json:"wCountryCode,omitempty"`
	BCountryCode    []string `bson:"bCountryCode,omitempty" json:"bCountryCode,omitempty"`
	WPublishIds     []int64  `bson:"wPublishIds,omitempty" json:"wPublishIds,omitempty"`
	BPublishIds     []int64  `bson:"bPublishIds,omitempty" json:"bPublishIds,omitempty"`
	Rate            int      `bson:"rate,omitempty" json:"rate,omitempty"`
	CDNPrefix       string   `bson:"cdnPrefix,omitempty" json:"cdnPrefix,omitempty"`
}

type HBAerospikeConfMap struct {
	// key is region
	ConfMap map[string]*HBAerospikeConf `bson:"configs,omitempty" json:"configs,omitempty"`
}

type HBAerospikeConf struct {
	Endpoint        string `bson:"endpoint,omitempty" json:"endpoint,omitempty"`
	MigrateEnable   bool   `bson:"migrate_enable,omitempty" json:"migrate_enable,omitempty"`
	MigrateEndpoint string `bson:"migrate_endpoint,omitempty" json:"migrate_endpoint,omitempty"`
	MigrateRate     int    `bson:"migrate_rate,omitempty" json:"migrate_rate,omitempty"`
	ReadTimeout     int64  `bson:"read_timeout,omitempty" json:"read_timeout,omitempty"`
	ReadRetry       int    `bson:"read_retry,omitempty" json:"read_retry,omitempty"`
	WriteTimeout    int64  `bson:"write_timeout,omitempty" json:"write_timeout,omitempty"`
	WriteRetry      int    `bson:"write_retry,omitempty" json:"write_retry,omitempty"`
	Expiration      int64  `bson:"expiration,omitempty" json:"expiration,omitempty"`
	ConnectionSize  int    `bson:"connection_size,omitempty" json:"connection_size,omitempty"`
	Namespace       string `bson:"namespace,omitempty" json:"namespace,omitempty"`
	SetName         string `bson:"set_name,omitempty" json:"set_name,omitempty"`
	Updated         int64  `bson:"updated,omitempty" json:"updated,omitempty"`
}

type ResidualMemoryFilterConf struct {
	AdTypeList     []int32 `bson:"adTypeList,omitempty" json:"adTypeList,omitempty"`
	ResidualMemory float64 `bson:"residualMemory,omitempty" json:"residualMemory,omitempty"`
	PlatformList   []int   `bson:"platformList,omitempty" json:"platformList,omitempty"`
	Rate           int     `bson:"rate,omitempty" json:"rate,omitempty"`
}

type HBV5ABTestConf struct {
	HBV5ABTestConfigs []*HBV5ABTestConfig `bson:"hb_v5_abtest_configs,omitempty" json:"hb_v5_abtest_configs,omitempty"`
}

type HBV5ABTestConfig struct {
	CountryCode   []string `bson:"cc,omitempty" json:"cc,omitempty"`
	MediationName []string `bson:"mediation_name,omitempty" json:"mediation_name,omitempty"`
	Platform      []int    `bson:"platform,omitempty" json:"platform,omitempty"`
	AdType        []int32  `bson:"ad_type,omitempty" json:"ad_type,omitempty"`
	PublisherID   []int64  `bson:"publisher_id,omitempty" json:"publisher_id,omitempty"`
	AppID         []int64  `bson:"app_id,omitempty" json:"app_id,omitempty"`
	UnitID        []int64  `bson:"unit_id,omitempty" json:"unit_id,omitempty"`
	SDKVersion    string   `bson:"sdk_version,omitempty" json:"sdk_version,omitempty"`
	Rate          int      `bson:"rate,omitempty" json:"rate,omitempty"`
}

type TmaxABTestConfMap struct {
	UseBidRequestTmax bool              `bson:"use_bid_tmax,omitempty" json:"use_bid_tmax,omitempty"`
	Timeout           map[string]int32  `bson:"timeout,omitempty" json:"timeout,omitempty"`
	NetworkCost       map[string]int32  `bson:"network_cost,omitempty" json:"network_cost,omitempty"`
	TmaxABTestConfigs []*TmaxABTestConf `bson:"tmax_abtest_configs,omitempty" json:"tmax_abtest_configs,omitempty"`
}

type TmaxABTestConf struct {
	Cloud         []string `bson:"cloud,omitempty" json:"cloud,omitempty"`
	Region        []string `bson:"region,omitempty" json:"region,omitempty"`
	CountryCode   []string `bson:"cc,omitempty" json:"cc,omitempty"`
	IsHB          []int    `bson:"is_hb,omitempty" json:"is_hb,omitempty"`
	MediationName []string `bson:"mediation_name,omitempty" json:"mediation_name,omitempty"`
	Platform      []int    `bson:"platform,omitempty" json:"platform,omitempty"`
	RequestType   []int    `bson:"request_type,omitempty" json:"request_type,omitempty"`
	Scenario      []string `bson:"scenario,omitempty" json:"scenario,omitempty"`
	AdType        []int32  `bson:"ad_type,omitempty" json:"ad_type,omitempty"`
	PublisherID   []int64  `bson:"publisher_id,omitempty" json:"publisher_id,omitempty"`
	AppID         []int64  `bson:"app_id,omitempty" json:"app_id,omitempty"`
	UnitID        []int64  `bson:"unit_id,omitempty" json:"unit_id,omitempty"`
	Rate          int      `bson:"rate,omitempty" json:"rate,omitempty"`
	Tmax          []int32  `bson:"tmax,omitempty" json:"tmax,omitempty"`
}

// Aerospike 的存储配置
type HBAerospikeStorageConf struct {
	Cloud  string `bson:"cloud,omitempty" json:"cloud,omitempty"`
	Region string `bson:"region,omitempty" json:"region,omitempty"`
	// 切量比例 [0, 1] - AerospikeGzip压缩比例
	GzipRate float64 `bson:"gzip_rate,omitempty" json:"gzip_rate,omitempty"`
	// 切量比例 [0, 1] - Aerospike移除AppInfo, UnitInfo, PublisherInfo三个冗余对象
	RemoveRedundancyRate float64 `bson:"remove_redundancy_rate,omitempty" json:"remove_redundancy_rate,omitempty"`
}
type HBAerospikeStorageConfArr struct {
	Configs []*HBAerospikeStorageConf `bson:"configs,omitempty" json:"configs,omitempty"`
}

type MpToPioneerABTestConf struct {
	TotalRate map[string]int `bson:"totalRate,omitempty" json:"totalRate,omitempty"`
	RandType  int            `bson:"randType,omitempty" json:"randType,omitempty"`
}

type GlobalTemplateMap struct {
	Video            map[string]string          `bson:"video,omitempty" json:"video"`
	EndScreen        map[string]string          `bson:"endscreen,omitempty" json:"endscreen"`
	MiniCard         map[string]string          `bson:"minicard,omitempty" json:"minicard"`
	BigTempalte      map[string]string          `bson:"bigTemplate,omitempty" json:"bigTemplate"`
	DiverseVideo     map[string]DiverseTemplate `bson:"diverseVideo,omitempty" json:"diverseVideo"`
	DiverseEndScreen map[string]DiverseTemplate `bson:"diverseEndscreen,omitempty" json:"diverseEndscreen"`
}

type DiverseTemplate struct {
	FullScreen string `bson:"full-screen,omitempty" json:"full-screen"`
	HalfScreen string `bson:"half-screen,omitempty" json:"half-screen"`
}

type HBFilterDeviceGEOConf struct {
	HBFilterDeviceGEOConfigs []*HBFilterDeviceGEOConfig `bson:"configs,omitempty" json:"configs,omitempty"`
}

type HBFilterDeviceGEOConfig struct {
	CountryCode   []string `bson:"cc,omitempty" json:"cc,omitempty"`
	MediationName []string `bson:"mediation_name,omitempty" json:"mediation_name,omitempty"`
	PublisherID   []int64  `bson:"publisher_id,omitempty" json:"publisher_id,omitempty"`
	AppID         []int64  `bson:"app_id,omitempty" json:"app_id,omitempty"`
}

// bid server ab test 配置
type HBRequestBidServerConf struct {
	//Cloud  string `bson:"cloud,omitempty" json:"cloud,omitempty"`
	//Region string `bson:"region,omitempty" json:"region,omitempty"`
	//Rate   int    `bson:"rate,omitempty" json:"rate,omitempty"` // rate 为 (0, 10000], 如 100 表示 100/10000 = 1% 的切量比例
	Cloud         []string `bson:"cloud,omitempty" json:"cloud,omitempty"`
	Region        []string `bson:"region,omitempty" json:"region,omitempty"`
	CountryCode   []string `bson:"cc,omitempty" json:"cc,omitempty"`
	MediationName []string `bson:"mediation_name,omitempty" json:"mediation_name,omitempty"`
	Platform      []int    `bson:"platform,omitempty" json:"platform,omitempty"`
	RequestType   []int    `bson:"request_type,omitempty" json:"request_type,omitempty"`
	Scenario      []string `bson:"scenario,omitempty" json:"scenario,omitempty"`
	AdType        []int32  `bson:"ad_type,omitempty" json:"ad_type,omitempty"`
	PublisherID   []int64  `bson:"publisher_id,omitempty" json:"publisher_id,omitempty"`
	AppID         []int64  `bson:"app_id,omitempty" json:"app_id,omitempty"`
	UnitID        []int64  `bson:"unit_id,omitempty" json:"unit_id,omitempty"`
	Rate          int      `bson:"rate,omitempty" json:"rate,omitempty"` // rate 为 (0, 10000], 如 100 表示 100/10000 = 1% 的切量比例
}

type HBRequestBidServerConfArr struct {
	Configs []*HBRequestBidServerConf `bson:"configs,omitempty" json:"configs,omitempty"`
}

type HBLoadFilterConfigs struct {
	Status          bool    `bson:"status,omitempty" json:"status,omitempty"`
	IncAdTypes      []int32 `bson:"inc_ad_types,omitempty" json:"inc_ad_types,omitempty"`
	ExcAdTypes      []int32 `bson:"exc_ad_types,omitempty" json:"exc_ad_types,omitempty"`
	IncPublisherIds []int64 `bson:"inc_publisher_ids,omitempty" json:"inc_publisher_ids,omitempty"`
	ExcPublisherIds []int64 `bson:"exc_publisher_ids,omitempty" json:"exc_publisher_ids,omitempty"`
	IncAppIds       []int64 `bson:"inc_app_ids,omitempty" json:"inc_app_ids,omitempty"`
	ExcAppIds       []int64 `bson:"exc_app_ids,omitempty" json:"exc_app_ids,omitempty"`
	IncUnitIds      []int64 `bson:"inc_unit_ids,omitempty" json:"inc_unit_ids,omitempty"`
	ExcUnitIds      []int64 `bson:"exc_unit_ids,omitempty" json:"exc_unit_ids,omitempty"`
}

type FilterAutoClickConf struct {
	TotalStatus  bool    `bson:"totalStatus,omitempty" json:"totalStatus,omitempty"`
	DspList      []int64 `bson:"dspList,omitempty" json:"dspList,omitempty"`
	PubList      []int64 `bson:"pubList,omitempty" json:"pubList,omitempty"`
	DspBlacklist []int64 `bson:"dspBlacklist,omitempty" json:"dspBlacklist,omitempty"`
	PubBlacklist []int64 `bson:"pubBlacklist,omitempty" json:"pubBlacklist,omitempty"`
}

type CdnTrackingDomainABTestConf struct {
	TotalRate  map[string]int `bson:"totalRate,omitempty" json:"totalRate,omitempty"`
	DomainList []string       `bson:"domainList,omitempty" json:"domainList,omitempty"`
}

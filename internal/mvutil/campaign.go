package mvutil

import (
	"gitlab.mobvista.com/ADN/adnet/internal/mvconst"
	smodel "gitlab.mobvista.com/ADN/structs/model"
)

type CampaignResult struct {
	RequestID    string         `json:"requestID"`
	CampaignList []CampaignInfo `json:"campaign"`
	RetCode      int            `json:"retCode"`
	RetMsg       string         `json:"retMsg"`
}

type CampaignInfo struct {
	CampaignId         int64                        `bson:"campaignId,omitempty" json:"campaignId,omitempty"`
	AdvertiserId       int32                        `bson:"advertiserId,omitempty" json:"advertiserId,omitempty"`
	TrackingUrl        string                       `bson:"trackingUrl,omitempty" json:"trackingUrl,omitempty"`
	TrackingUrlHttps   string                       `bson:"trackingUrlHttps,omitempty" json:"trackingUrlHttps,omitempty"`
	DirectUrl          string                       `bson:"directUrl,omitempty" json:"directUrl,omitempty"`
	Price              float64                      `bson:"price,omitempty" json:"price,omitempty"`
	OriPrice           float64                      `bson:"oriPrice,omitempty" json:"oriPrice,omitempty"`
	CityCodeV2         map[string][]int32           `bson:"cityCodeV2,omitempty" json:"cityCodeV2,omitempty"`
	Status             int32                        `bson:"status,omitempty" json:"status,omitempty"`
	Network            int32                        `bson:"network,omitempty" json:"network,omitempty"`
	PreviewUrl         string                       `bson:"previewUrl,omitempty" json:"previewUrl,omitempty"`
	PackageName        string                       `bson:"packageName,omitempty" json:"packageName,omitempty"`
	CampaignType       int32                        `bson:"campaignType,omitempty" json:"campaignType,omitempty"`
	Ctype              int32                        `bson:"ctype,omitempty" json:"ctype,omitempty"`
	AppSize            string                       `bson:"appSize,omitempty" json:"appSize,omitempty"`
	Tag                int32                        `bson:"tag,omitempty" json:"tag,omitempty"`
	AdSourceId         int32                        `bson:"adSourceId,omitempty" json:"adSourceId,omitempty"`
	PublisherId        int64                        `bson:"publisherId,omitempty" json:"publisherId,omitempty"`
	FrequencyCap       int32                        `bson:"frequencyCap,omitempty" json:"frequencyCap,omitempty"`
	DirectPackageName  string                       `bson:"directPackageName,omitempty" json:"directPackageName,omitempty"`
	SdkPackageName     string                       `bson:"sdkPackageName,omitempty" json:"sdkPackageName,omitempty"`
	AdvImp             []*AdvImp                    `bson:"advImp,omitempty" json:"advImp,omitempty"`
	AdUrlList          []string                     `bson:"adUrlList,omitempty" json:"adUrlList,omitempty"`
	JumpType           int32                        `bson:"jumpType,omitempty" json:"jumpType,omitempty"`
	VbaConnecting      int32                        `bson:"vbaConnecting,omitempty" json:"vbaConnecting,omitempty"`
	VbaTrackingLink    string                       `bson:"vbaTrackingLink,omitempty" json:"vbaTrackingLink,omitempty"`
	RetargetingDevice  int32                        `bson:"retargetingDevice,omitempty" json:"retargetingDevice,omitempty"`
	SendDeviceidRate   int32                        `bson:"sendDeviceidRate,omitempty" json:"sendDeviceidRate,omitempty"`
	Endcard            map[string]*EndCard          `bson:"endcard,omitempty" json:"endcard,omitempty"`
	Loopback           *LoopBack                    `bson:"loopback,omitempty" json:"loopback,omitempty"`
	BelongType         int32                        `bson:"belongType,omitempty" json:"belongType,omitempty"`
	ConfigVBA          *ConfigVBA                   `bson:"configVBA,omitempty" json:"configVBA,omitempty"`
	AppPostList        *AppPostList                 `bson:"appPostList,omitempty" json:"appPostList,omitempty"`
	BlackSubidListV2   map[string]map[string]string `bson:"blackSubidListV2,omitempty" json:"blackSubidListV2,omitempty"`
	BtV4               *BtV4                        `bson:"btV4,omitempty" json:"btV4,omitempty"`
	OpenType           int32                        `bson:"openType,omitempty" json:"openType,omitempty"`
	SubCategoryName    []string                     `bson:"subCategoryName,omitempty" json:"subCategoryName,omitempty"`
	IsCampaignCreative int32                        `bson:"isCampaignCreative,omitempty" json:"isCampaignCreative,omitempty"`
	CostType           int32                        `bson:"costType,omitempty" json:"costType,omitempty"`
	Source             int32                        `bson:"source,omitempty" json:"source,omitempty"`
	ChnID              int                          `bson:"chnId,omitempty" json:"chnId,omitempty"`
	ThirdParty         string                       `bson:"thirdParty,omitempty" json:"thirdParty,omitempty"`

	JumpTypeConfig  map[string]int32 `bson:"JUMP_TYPE_CONFIG,omitempty" json:"JUMP_TYPE_CONFIG,omitempty"`
	JumpTypeConfig2 map[string]int32 `bson:"JUMP_TYPE_CONFIG_2,omitempty" json:"JUMP_TYPE_CONFIG_2,omitempty"`

	//JUMPTYPECONFIG     map[string]int32            `bson:"JUMP_TYPE_CONFIG,omitempty" json:"JUMP_TYPE_CONFIG,omitempty"`
	//JUMPTYPECONFIGV2 map[string]int32 `bson:"JUMP_TYPE_CONFIG_2,omitempty" json:"JUMP_TYPE_CONFIG_2,omitempty"`
	Updated     int             `bson:"updated,omitempty" json:"updated,omitempty"`
	MChanlPrice []*ChannelPrice `bson:"mChanlPrice,omitempty" json:"mChanlPrice,omitempty"`
	Category    int32           `bson:"category,omitempty" json:"category,omitempty"`
	WxAppId     string          `bson:"wxAppId,omitempty" json:"wxAppId,omitempty"`
	WxPath      string          `bson:"wxPath,omitempty" json:"wxPath,omitempty"`
	BindId      string          `bson:"bindId,omitempty" json:"bindId,omitempty"`
	DeepLink    string          `bson:"deepLink,omitempty" json:"deepLink,omitempty"`
	ApkVersion  string          `bson:"apkVersion,omitempty" json:"apkVersion,omitempty"`
	ApkMd5      string          `bson:"apkMd5,omitempty" json:"apkMd5,omitempty"`
	ApkUrl      string          `bson:"apkUrl,omitempty" json:"apkUrl,omitempty"`
	// Creative 相关的字段
	BasicCrList *BasicCrList `bson:"basicCrList,omitempty" json:"basicCrList,omitempty"`
	//ImageList    map[string][]map[string]interface{} `bson:"imageList,omitempty" json:"imageList,omitempty"`
	//VideoList    map[string][]map[string]interface{} `bson:"videoList,omitempty" json:"videoList,omitempty"`
	ReadCreative          int                                 `bson:"readCreative,omitempty" json:"readCreative,omitempty"`
	CreateSrc             int                                 `bson:"createSrc,omitempty" json:"createSrc,omitempty"`
	FakeCreative          map[string]map[string]*FakeCreative `bson:"fakeCreative,omitempty" json:"fakeCreative,omitempty"`
	AlacRate              int                                 `bson:"alacRate,omitempty" json:"alacRate,omitempty"`
	AlecfcRate            int                                 `bson:"alecfcRate,omitempty" json:"alecfcRate,omitempty"`
	Mof                   int                                 `bson:"mof,omitempty" json:"mof,omitempty"`
	MCountryChanlPrice    map[string]float64                  `bson:"mCountryChanlPrice,omitempty" json:"mCountryChanlPrice,omitempty"`       // M  按 国家 + 渠道的  given price
	MCountryChanlOriPrice map[string]float64                  `bson:"mCountryChanlOriPrice,omitempty" json:"mCountryChanlOriPrice,omitempty"` // M  按 国家 + 渠道的  receive price
	NeedToNotice3s        int                                 `bson:"needToNotice3s,omitempty" json:"needToNotice3s,omitempty"`               // 判断是否需要通知3s
}

// IsSSPlatform 判断是否为SS平台单子
func (c *CampaignInfo) IsSSPlatform() bool {
	return c.CreateSrc == mvconst.CampaignSourceSSPlatform ||
		c.CreateSrc == mvconst.CampaignSourceSSAdvPlatform
}

type FakeCreative struct {
	Id     string `bson:"id,omitempty" json:"id,omitempty"`
	Rate   int    `bson:"rate,omitempty" json:"rate,omitempty"`
	Name   string `bson:"name,omitempty" json:"name,omitempty"`
	AdType string `bson:"adType,omitempty" json:"adType,omitempty"`
}

type NewSubIdInfo struct {
	Rate        int                `bson:"rate,omitempty" json:"rate,omitempty"`
	SubId       int64              `bson:"subId,omitempty" json:"subId,omitempty"`
	AdType      int                `bson:"adType,omitempty" json:"adType,omitempty"`
	PackageName string             `bson:"packageName,omitempty" json:"packageName,omitempty"`
	Creative    []*NewFakeCreative `bson:"creative,omitempty" json:"creative,omitempty"`
	DspSubIds   []*NewSubIdInfo    `bson:"dspSubIds,omitempty" json:"dspSubIds,omitempty"`
}

func (ns *NewSubIdInfo) GetRate() int {
	return ns.Rate
}

type NewFakeCreative struct {
	Id            int    `bson:"id,omitempty" json:"id,omitempty"`
	AdnCreativeId int    `bson:"adnCrId,omitempty" json:"adnCrId,omitempty"`
	Rate          int    `bson:"rate,omitempty" json:"rate,omitempty"`
	Name          string `bson:"name,omitempty" json:"name,omitempty"`
	AdType        int    `bson:"adType,omitempty" json:"adType,omitempty"`
	AdvEndCardId  int    `bson:"advEcId,omitempty" json:"advEcId,omitempty"`
	AdnEndCardId  int    `bson:"adnEcId,omitempty" json:"adnEcId,omitempty"`
	EndCardType   int    `bson:"ecType,omitempty" json:"ecType,omitempty"`
	CreativeType  int    `bson:"crType,omitempty" json:"crType,omitempty"`
}

func (nfc *NewFakeCreative) GetRate() int {
	return nfc.Rate
}

type AppPostList struct {
	Include []string `bson:"include,omitempty" json:"include"`
	Exclude []string `bson:"exclude,omitempty" json:"exclude"`
}

type JumpParam struct {
	B2t              int32 `bson:"b2t,omitempty" json:"b2t"`
	B2tStatus        int32 `bson:"b2tStatus,omitempty" json:"b2tStatus"`
	NoDeviceId       int32 `bson:"noDeviceId,omitempty" json:"noDeviceId"`
	NoDeviceIdStatus int32 `bson:"noDeviceIdStatus,omitempty" json:"noDeviceIdStatus"`
}

type ReduceRuleItem struct {
	Priority int32 `bson:"priority,omitempty" json:"priority"`
	Install  int32 `bson:"install,omitempty" json:"install"`
	Status   int32 `bson:"status,omitempty" json:"status"`
	Start    int64 `bson:"start,omitempty" json:"start"`
}

type LoopBack struct {
	Domain string `bson:"domain,omitempty" json:"domain,omitempty"`
	Key    string `bson:"key,omitempty" json:"key,omitempty"`
	Value  string `bson:"value,omitempty" json:"value,omitempty"`
	Rate   int32  `bson:"rate,omitempty" json:"rate,omitempty"`
}

type EndCard struct {
	Urls             []*EndCardUrls          `bson:"urls,omitempty" json:"urls,omitempty"`
	Status           int32                   `bson:"status,omitempty" json:"status,omitempty"`
	Orientation      int32                   `bson:"orientation,omitempty" json:"orientation,omitempty"`
	VideoTemplateUrl []*VideoTemplateUrlItem `bson:"videoTemplateUrl,omitempty" json:"videoTemplateUrl,omitempty"`
	EndcardProtocal  int                     `bson:"endcardProtocol,omitempty" json:"endcardProtocol,omitempty"`
	EndcardRate      map[string]int          `bson:"endcardRate,omitempty" json:"endcardRate,omitempty"`
	EndcardType      int32                   `bson:"endcardType,omitempty" json:"endcardType,omitempty"`
}

type EndcardItem struct {
	Url             string `bson:"url,omitempty" json:"url,omitempty"`
	UrlV2           string `bson:"url_v2,omitempty" json:"url_v2,omitempty"`
	Orientation     int32  `bson:"orientation,omitempty" json:"orientation,omitempty"`
	ID              int32  `bson:"id,omitempty" json:"id,omitempty"`
	EndcardProtocal int
	EndcardRate     map[string]int
}

type VideoTemplateUrlItem struct {
	ID           int32  `bson:"id,omitempty" json:"id,omitempty"`
	URL          string `bson:"url,omitempty" json:"url,omitempty"`
	URLZip       string `bson:"url_zip,omitempty" json:"url_zip,omitempty"`
	Weight       int32  `bson:"weight,omitempty" json:"weight,omitempty"`
	PausedURL    string `bson:"paused_url,omitempty" json:"paused_url,omitempty"`
	PausedURLZip string `bson:"paused_url_zip,omitempty" json:"paused_url_zip,omitempty"`
}

type EndCardUrls struct {
	Id     int32  `bson:"id,omitempty" json:"id,omitempty"`
	Url    string `bson:"url,omitempty" json:"url,omitempty"`
	Weight int32  `bson:"weight,omitempty" json:"weight,omitempty"`
	UrlV2  string `bson:"url_v2,omitempty" json:"url_v2,omitempty"`
}

type AdvImp struct {
	Sec int32  `bson:"sec,omitempty" json:"sec,omitempty"`
	Url string `bson:"url,omitempty" json:"url,omitempty"`
}

type Creative struct {
	CampaignId     int64       `bson:"campaignId,omitempty" json:"campaignId,omitempty"`
	CreativeId     int64       `bson:"creativeId,omitempty" json:"creativeId,omitempty"`
	Lang           int32       `bson:"lang,omitempty" json:"lang,omitempty"`
	Type           int32       `bson:"type,omitempty" json:"type,omitempty"`
	Width          int32       `bson:"width,omitempty" json:"width,omitempty"`
	Height         int32       `bson:"height,omitempty" json:"height,omitempty"`
	ImageSize      string      `bson:"imageSize,omitempty" json:"imageSize,omitempty"`
	ImageSizeId    int32       `bson:"imageSizeId,omitempty" json:"imageSizeId,omitempty"`
	Name           string      `bson:"name,omitempty" json:"name,omitempty"`
	TextVideo      interface{} `bson:"textVideo,omitempty" json:"textVideo,omitempty"`
	VideoUrlEncode string      `bson:"videoUrlEncode,omitempty" json:"videoUrlEncode,omitempty"`
	VideoUrl       string      `bson:"videoUrl,omitempty" json:"videoUrl,omitempty"`
	ImageUrl       string      `bson:"imageUrl,omitempty" json:"imageUrl,omitempty"`
	Comment        string      `bson:"comment,omitempty" json:"comment,omitempty"`
	CreativeCta    string      `bson:"creativeCta,omitempty" json:"creativeCta,omitempty"`
	Status         int32       `bson:"status,omitempty" json:"status,omitempty"`
	Tag            int32       `bson:"tag,omitempty" json:"tag,omitempty"`
	Created        int64       `bson:"created,omitempty" json:"created,omitempty"`
	ResourceType   int32       `bson:"resourceType,omitempty" json:"resourceType,omitempty"`
	Mime           []string    `bson:"mime,omitempty" json:"mime,omitempty"`
	Attribute      []int32     `bson:"attribute,omitempty" json:"attribute,omitempty"`
	TemplateType   interface{} `bson:"templateType,omitempty" json:"templateType,omitempty"`
	TagCode        string      `bson:"tagCode,omitempty" json:"tagCode,omitempty"`
	ShowType       int32       `bson:"showType,omitempty" json:"showType,omitempty"`
}

type BtV2 struct {
	SubIds  interface{}         `bson:"subIds,omitempty" json:"subIds"`
	BtClass map[string]*BtClass `bson:"btClass,omitempty" json:"btClass"`
}

type BtClass struct {
	Percent   float64 `bson:"percent,omitempty" json:"percent"`
	CapMargin int32   `bson:"capMargin,omitempty" json:"capMargin"`
	Status    int32   `bson:"status,omitempty" json:"status"`
}

type BtV4 struct {
	SubIds  map[string]*SubInfo `bson:"subIds,omitempty" json:"subIds"`
	BtClass map[string]*BtClass `bson:"btClass,omitempty" json:"btClass"`
}

// type BtV3 struct {
// 	//SubIds  map[int64]SubInfo `bson:"subIds,omitempty" json:"subIds"`
// 	SubIds  interface{}     `bson:"subIds,omitempty" json:"subIds"`
// 	BtClass map[int]BtClass `bson:"btClass,omitempty" json:"btClass"`
// }

type SubInfo struct {
	Rate        int                    `bson:"rate,omitempty" json:"rate"`
	PackageName string                 `bson:"packageName,omitempty" json:"packageName"`
	DspSubIds   map[string]*DspSubInfo `bson:"dspSubIds,omitempty" json:"dspSubIds"`
}

type SubInfoe struct {
	Rate        int     `bson:"rate,omitempty" json:"rate"`
	PackageName string  `bson:"packageName,omitempty" json:"packageName"`
	DspSubIds   []int32 `bson:"dspSubIds,omitempty" json:"dspSubIds"`
}

type DspSubInfo struct {
	Rate        int    `bson:"rate,omitempty" json:"rate"`
	PackageName string `bson:"packageName,omitempty" json:"packageName"`
}

type ConfigVBA struct {
	UseVBA       int `bson:"useVBA,omitempty" json:"useVBA"`
	FrequencyCap int `bson:"frequencyCap,omitempty" json:"frequencyCap"`
	Status       int `bson:"status,omitempty" json:"status"`
}

type ChannelPrice struct {
	Chanl string  `bson:"chanl,omitempty" json:"chanl"`
	Price float64 `bson:"price,omitempty" json:"price"`
}

type ExtCreativeNew struct {
	PlayWithoutVideo int    `json:"pwv,omitempty"` // playable_ads_without_video
	VideoEndType     int    `json:"vet,omitempty"` // VideoEndType
	TemplateGroupId  *int   `json:"t_group,omitempty"`
	EndScreenId      string `json:"es_id,omitempty"`
	IsCreativeNew    bool   `json:"is_new,omitempty"`
}

type BasicCrList struct {
	AppName   string  `bson:"401,omitempty" json:"401"`
	AppDesc   string  `bson:"402,omitempty" json:"402"`
	AppRate   float64 `bson:"403,omitempty" json:"403"`
	CtaButton string  `bson:"404,omitempty" json:"404"`
	AppIcon   string  `bson:"405,omitempty" json:"405"`
	NumRating int     `bson:"406,omitempty" json:"406"`
}

//type ImageInfo struct {
//	AdvCreativeId int    `bson:"adv_creative_id,omitempty", json:"adv_creative_id"`
//	Attribute     string `bson:"attribute,omitempty", json:"attribute"`
//	Cname         string `bson:"cname,omitempty", json:"cname"`
//	Ext           string `bson:"ext,omitempty", json:"ext"`
//	Fmd5          string `bson:"fmd5,omitempty", json:"fmd5"`
//	Mime          string `bson:"mime,omitempty", json:"mime"`
//	PkgNoRecall   int    `bson:"pkg_no_recall,omitempty", json:"pkg_no_recall"`
//	Resolution    string `bson:"resolution,omitempty", json:"resolution"`
//	Url           string `bson:"url,omitempty", json:"url"`
//	AdnCreativeId string `bson:"adnCreativeId,omitempty", json:"adnCreativeId"`
//}
//
//type VideoInfo struct {
//	Url             string `bson:"url,omitempty", json:"url"`
//	VideoLength     int    `bson:"video_length,omitempty", json:"video_length"`
//	VideoSize       int    `bson:"video_size,omitempty", json:"video_size"`
//	VideoResolution string `bson:"video_resolution,omitempty", json:"video_resolution"`
//	Width           int    `bson:"width,omitempty", json:"width"`
//	Height          int    `bson:"height,omitempty", json:"height"`
//	VideoTruncation int    `bson:"video_truncation,omitempty", json:"video_truncation"`
//	WatchMile       int    `bson:"watch_mile,omitempty", json:"watch_mile"`
//	BitRate         int    `bson:"bit_rate,omitempty", json:"bit_rate"`
//	ScreenShot      string `bson:"screen_show,omitempty", json:"screen_show"`
//	Mime            string `bson:"mime,omitempty", json:"mime"`
//	Attribute       string `bson:"attribute,omitempty", json:"attribute"`
//	Fmd5            string `bson:"f_md5,omitempty", json:"f_md5"`
//	AdvCreativeId   int    `bson:"adv_creative_id,omitempty", json:"adv_creative_id"`
//	SourceMark      string `bson:"source_mark,omitempty", json:"source_mark"`
//	Orientation     string `bson:"orientation,omitempty", json:"orientation"`
//	Clarity         int    `bson:"clarity,omitempty", json:"clarity"`
//	PkgNoCall       int    `bson:"pkg_no_call,omitempty", json:"pkg_no_call"`
//	Cname           string `bson:"cname,omitempty", json:"cname"`
//	AdnCreativeId   string `bson:"adnCreativeId,omitempty", json:"adnCreativeId"`
//}

const (
	ATTR_BRAND_OFFER = 1
	ATTR_VTA_OFFER   = 2
	ATTR_CITY_OFFER  = 4
)

// func GetBtV3(btV3 interface{}) map[int64]SubInfo {
// 	subInfo, ok := btV3.(map[int64]SubInfo)
// 	if ok {
// 		return subInfo
// 	}
// 	subInfoe, ok := btV3.(map[int64]SubInfoe)
// 	result := make(map[int64]SubInfo, len(subInfoe))
// 	if ok {
// 		for k, val := range subInfoe {
// 			var sub SubInfo
// 			sub.Rate = val.Rate
// 			sub.PackageName = val.PackageName
// 			sub.DspSubIds = make(map[int64]DspSubInfo)
// 			result[k] = sub
// 		}
// 		return result
// 	}
// 	return nil
// }

func GetCityCode(cityCode interface{}) map[string]int64 {
	res, ok := cityCode.(map[string]int64)
	if ok {
		return res
	}
	return nil
}

func GetFCA(r *RequestParams, campaign *smodel.CampaignInfo, confCamIds []int64, newDefaultFca int) int {
	// 频次控制整体开关开启
	if r.Param.FcaSwitch {
		// 默认为1000
		fcaVal := 1000
		// 若有配置，则取配置的值
		if newDefaultFca != 0 {
			fcaVal = newDefaultFca
		}
		if AppFcaDefault(r.AppInfo.App.FrequencyCap) &&
			!InInt64Arr(campaign.CampaignId, confCamIds) {
			return fcaVal
		}
	}
	// 如果是VBA，则获取config里面的fca
	if IsVBA(campaign) && campaign.ConfigVBA != nil {
		configVBA := *(campaign.ConfigVBA)
		fca := configVBA.FrequencyCap
		if fca <= 0 {
			return 1
		} else {
			return fca
		}
	}
	if campaign.FrequencyCap > 0 {
		return int(campaign.FrequencyCap)
	}
	if r.AppInfo.App.FrequencyCap > 0 {
		return r.AppInfo.App.FrequencyCap
	}
	return 5
}

func GetFCB(source int32) int {
	if source == int32(1) {
		return 1
	}
	return 2
}

func IsVBA(campaign *smodel.CampaignInfo) bool {
	// 不能是VTA单子
	if IsVTA(campaign) {
		return false
	}
	// 配置了VBA
	if campaign.ConfigVBA == nil {
		return false
	}
	if campaign.ConfigVBA.Status != 1 || campaign.ConfigVBA.UseVBA != 1 {
		return false
	}
	return true
}

func IsVTA(campaign *smodel.CampaignInfo) bool {
	if campaign.VbaConnecting != 1 {
		return false
	}
	if campaign.VbaTrackingLink == "" {
		return false
	}
	return true
}

func IsBrandOffer(campaign *smodel.CampaignInfo) bool {
	if campaign.Tag == 0 {
		return false
	}
	return campaign.Tag == 4
}

func IsCityOffer(campaign *smodel.CampaignInfo) bool {
	return len(campaign.CityCodeV2) > 0
}

func GetVTALink(campaign *smodel.CampaignInfo) string {
	if campaign.VbaTrackingLink == "" {
		return ""
	}
	return campaign.VbaTrackingLink
}

func GetVTATag(campaign *smodel.CampaignInfo) int {
	if IsVTA(campaign) {
		return 1
	}
	if IsVBA(campaign) {
		return 2
	}
	vtaLink := GetVTALink(campaign)
	if len(vtaLink) > 0 {
		return 3
	}
	return 0
}

func IsJsVideo(adType int32) bool {
	if adType == mvconst.ADTypeJSBannerVideo || adType == mvconst.ADTypeJSNativeVideo {
		return true
	} else {
		return false
	}
}

func IsWxAdType(adType int32) bool {
	return adType == mvconst.ADTypeWXBanner || adType == mvconst.ADTypeWXAppwall || adType == mvconst.ADTypeWXNative || adType == mvconst.ADTypeWXRewardImg
}

// 是否为ss单子流量，且clickmode为6。包括现有的ss abtest
func IsSSOfferAndCM6(c *smodel.CampaignInfo, params *Params) bool {
	return params.Extra10 == mvconst.JUMP_TYPE_CLIENT_DO_ALL &&
		(c.IsSSPlatform() || c.CreateSrc == mvconst.CampaignSourceSSRMABTest)
}

func GetCloudExtra(cloud string) string {
	if len(cloud) > 0 {
		return cloud
	}
	return mvconst.CLOUD_NAME_AWS
}

package corsair_proto

import (
	"fmt"

	"gitlab.mobvista.com/adserver/recommend_protocols/go/ad_server"
)

type Campaign struct {
	CampaignId         string
	AdSource           ad_server.ADSource
	AdTemplate         *ad_server.ADTemplate
	ImageSizeId        ad_server.ImageSizeEnum
	OfferType          *int32
	BtType             ad_server.BtType
	CreativeId         *string
	CreativeTypeIdMap  map[ad_server.CreativeType]int64
	AdElementTemplate  *ad_server.AdElementTemplate
	Downloadtest       *int32
	CreativeId2        *string
	CreativeTypeIdMap2 map[ad_server.CreativeType]int64
	DynamicCreative    map[ad_server.CreativeType]string
	Playable           bool
	EndcardUrl         *string
	ExtPlayable        *int32
	VideoEndTypeAs     *int32
	Orientation        *ad_server.Orientation
	UsageVideo         *bool
	TemplateGroup      *ad_server.TemplateGroup
	VideoTemplateId    *ad_server.VideoTemplateId
	EndCardTemplateId  *ad_server.EndCardTemplateId
	MiniCardTemplateId *ad_server.MiniCardTemplateId
	// unused fields # 25 to 30
	AppName                 *string
	AppDesc                 *string
	PackageName             *string
	AppSize                 *string
	IconURL                 *string
	ImageURL                *string
	ImageSize               *string
	VideoURL                *string
	VideoLength             *int32
	VideoSize               *int32
	BitRate                 *int32
	VideoResolution         *string
	CType                   *int32
	AdvImpList              []*AdvImp
	AdURLList               []string
	AdvertiserId            *int32
	ClickURL                *string
	Price                   *float64
	OfferName               *string
	EndcardURL              *string
	InstallToken            *string
	ClickMode               *int32
	Rating                  *float64
	LandingType             *int32
	CtaText                 *string
	LinkType                *int32
	GuideLines              *string
	AdTracking              *AdTracking
	Rv                      *RewardVideo
	LoopBack                map[string]string
	VideoEndType            *int32
	PlayableAdsWithoutVideo *int32
	WatchMile               *int32
	TImp                    *int32
	FCA                     *int32
	FCB                     *int32
	ClickCacheTime          *int32
	RewardAmount            *int32
	RewardName              *string
	RetargetOffer           *int32
	StatsURL                *string
	NoticeURL               *string
	ExtStats2               *string
	NumberRating            *int32
	RawCampaignId           *string
	SSABTest                *bool
	SSAATest                *bool
	BidPrice                *float64
	UseAlgoPrice            *bool
	AlgoPriceIn             *float64
	AlgoPriceOut            *float64
	OriCampaignId           *string
	CreativeDataList        []*CreativeData
	SlotId                  *int32
	UrlTemplate             *string
	HtmlTemplate            *string
	DeepLink                *string
	CpdIds                  []string
	AsABTestResTag          *string
	AKS                     map[string]string
	Skad                    *Skad
	AdHtml                  *string

	ImageResolution   *string
	ImageMime         *string
	GifWidgets        *string
	WTick             *int32
	Floatball         *int32
	FloatballSkipTime *int32
	ImageWidth        *int
	ImageHeight       *int
	//IconResolution  *string
	//IconMime        *string
}

type Skad struct {
	Version         string
	Network         string
	AppleCampaignId string
	Targetid        string
	Nonce           string
	Sourceid        string
	Timestamp       string
	Sign            string
	Need            int32
}

func NewCampaign() *Campaign {
	return &Campaign{}
}

type CreativeData struct {
	DocId              string
	CreativeTypeIdList []*CreativeTypeId
}

func NewCreativeData() *CreativeData {
	return &CreativeData{}
}

type AdvImp struct {
	Second int32
	URL    string
}

func NewAdvImp() *AdvImp {
	return &AdvImp{}
}

type RewardVideo struct {
	Template    *int32
	Orientation *int32
	TemplateURL *string
	PausedURL   *string
	Images      map[string][]string
}

func NewRewardVideo() *RewardVideo {
	return &RewardVideo{}
}

type CreativeTypeId struct {
	Type       ad_server.CreativeType
	CreativeId int64
	MId        *int64
	CpdId      *int64
	OmId       *int64
}

func NewCreativeTypeId() *CreativeTypeId {
	return &CreativeTypeId{}
}

type AdTracking struct {
	Start            []string
	FirstQuart       []string
	Mid              []string
	ThirdQuart       []string
	Complete         []string
	Mute             []string
	Unmute           []string
	Impression       []string
	Click            []string
	EndcardShow      []string
	Close            []string
	PlayPerct        []*PlayPercent
	Pause            []string
	ApkDownloadStart []string
	ApkDownloadEnd   []string
	ApkInstall       []string
	PubImp           []string
	VideoClick       []string
	ImpressionT2     []string
}

func NewAdTracking() *AdTracking {
	return &AdTracking{}
}

type PlayPercent struct {
	Rate int32
	URL  string
}

func NewPlayPercent() *PlayPercent {
	return &PlayPercent{}
}

type Frame struct {
	CampaignList []*Campaign
	AdTemplate   ad_server.ADTemplate
}

func NewFrame() *Frame {
	return &Frame{}
}

type BigTemplate struct {
	BigTemplateId          ad_server.BigTemplateId
	SlotIndexCampaignIdMap map[int32]int64
}

func NewBigTemplate() *BigTemplate {
	return &BigTemplate{}
}

type RunTimeVariable struct {
	NumRecalled                 *int32
	GetCampaignPositionInfoTime *int32
	GetCampaignInfoTime         *int32
	RankTime                    *int32
	GetCampaignIdTime           *int32
	GetCampaignDetailTime       *int32
	PublisherId                 *int32
	RecallSize                  *int32
	RecallSizeAll               *int32
}

func NewRunTimeVariable() *RunTimeVariable {
	return &RunTimeVariable{}
}

type FilterReason struct {
	CampaignId int64
	Reason     string
}

func NewFilterReason() *FilterReason {
	return &FilterReason{}
}

type NoneResultReason int64

const (
	NoneResultReason_UNKNOWN      NoneResultReason = 0
	NoneResultReason_NOOFFER      NoneResultReason = 1
	NoneResultReason_STRICTFILTER NoneResultReason = 2
	NoneResultReason_BEYONDOFFSET NoneResultReason = 3
)

func (p NoneResultReason) String() string {
	switch p {
	case NoneResultReason_UNKNOWN:
		return "UNKNOWN"
	case NoneResultReason_NOOFFER:
		return "NOOFFER"
	case NoneResultReason_STRICTFILTER:
		return "STRICTFILTER"
	case NoneResultReason_BEYONDOFFSET:
		return "BEYONDOFFSET"
	}
	return "<UNSET>"
}

func NoneResultReasonFromString(s string) (NoneResultReason, error) {
	switch s {
	case "UNKNOWN":
		return NoneResultReason_UNKNOWN, nil
	case "NOOFFER":
		return NoneResultReason_NOOFFER, nil
	case "STRICTFILTER":
		return NoneResultReason_STRICTFILTER, nil
	case "BEYONDOFFSET":
		return NoneResultReason_BEYONDOFFSET, nil
	}
	return NoneResultReason(0), fmt.Errorf("not a valid NoneResultReason string")
}

type QueryResult_ struct {
	LogId           string
	RandValue       int32
	FlowTagId       int32
	AdBackendConfig string
	AdsByBackend    []*BackendAds
}

func NewQueryResult_() *QueryResult_ {
	return &QueryResult_{}
}

var Campaign_BidPrice_DEFAULT float64

func (p *Campaign) GetBidPrice() float64 {
	if p.BidPrice == nil {
		return Campaign_BidPrice_DEFAULT
	}
	return *p.BidPrice
}

var Campaign_AdTracking_DEFAULT *AdTracking

func (p *Campaign) GetAdTracking() *AdTracking {
	if p.AdTracking == nil {
		return Campaign_AdTracking_DEFAULT
	}
	return p.AdTracking
}

var Campaign_SSABTest_DEFAULT bool

func (p *Campaign) GetSSABTest() bool {
	if p.SSABTest == nil {
		return Campaign_SSABTest_DEFAULT
	}
	return *p.SSABTest
}

var Campaign_RawCampaignId_DEFAULT string

func (p *Campaign) GetRawCampaignId() string {
	if p.RawCampaignId == nil {
		return Campaign_RawCampaignId_DEFAULT
	}
	return *p.RawCampaignId
}

var Campaign_SSAATest_DEFAULT bool

func (p *Campaign) GetSSAATest() bool {
	if p.SSAATest == nil {
		return Campaign_SSAATest_DEFAULT
	}
	return *p.SSAATest
}

var Campaign_TemplateGroup_DEFAULT ad_server.TemplateGroup

func (p *Campaign) GetTemplateGroup() ad_server.TemplateGroup {
	if p.TemplateGroup == nil {
		return Campaign_TemplateGroup_DEFAULT
	}
	return *p.TemplateGroup
}

var Campaign_EndCardTemplateId_DEFAULT ad_server.EndCardTemplateId

func (p *Campaign) GetEndCardTemplateId() ad_server.EndCardTemplateId {
	if p.EndCardTemplateId == nil {
		return Campaign_EndCardTemplateId_DEFAULT
	}
	return *p.EndCardTemplateId
}

var Campaign_VideoTemplateId_DEFAULT ad_server.VideoTemplateId

func (p *Campaign) GetVideoTemplateId() ad_server.VideoTemplateId {
	if p.VideoTemplateId == nil {
		return Campaign_VideoTemplateId_DEFAULT
	}
	return *p.VideoTemplateId
}

var Campaign_PackageName_DEFAULT string

func (p *Campaign) GetPackageName() string {
	if p.PackageName == nil {
		return Campaign_PackageName_DEFAULT
	}
	return *p.PackageName
}

var Campaign_WatchMile_DEFAULT int32

func (p *Campaign) GetWatchMile() int32 {
	if p.WatchMile == nil {
		return Campaign_WatchMile_DEFAULT
	}
	return *p.WatchMile
}

var Campaign_AdvertiserId_DEFAULT int32

func (p *Campaign) GetAdvertiserId() int32 {
	if p.AdvertiserId == nil {
		return Campaign_AdvertiserId_DEFAULT
	}
	return *p.AdvertiserId
}

var Campaign_AppSize_DEFAULT string

func (p *Campaign) GetAppSize() string {
	if p.AppSize == nil {
		return Campaign_AppSize_DEFAULT
	}
	return *p.AppSize
}

var Campaign_Rating_DEFAULT float64

func (p *Campaign) GetRating() float64 {
	if p.Rating == nil {
		return Campaign_Rating_DEFAULT
	}
	return *p.Rating
}

var Campaign_NumberRating_DEFAULT int32

func (p *Campaign) GetNumberRating() int32 {
	if p.NumberRating == nil {
		return Campaign_NumberRating_DEFAULT
	}
	return *p.NumberRating
}

var Campaign_CtaText_DEFAULT string

func (p *Campaign) GetCtaText() string {
	if p.CtaText == nil {
		return Campaign_CtaText_DEFAULT
	}
	return *p.CtaText
}

var RewardVideo_Template_DEFAULT int32

func (p *RewardVideo) GetTemplate() int32 {
	if p.Template == nil {
		return RewardVideo_Template_DEFAULT
	}
	return *p.Template
}

var RewardVideo_Orientation_DEFAULT int32

func (p *RewardVideo) GetOrientation() int32 {
	if p.Orientation == nil {
		return RewardVideo_Orientation_DEFAULT
	}
	return *p.Orientation
}

var RewardVideo_TemplateURL_DEFAULT string

func (p *RewardVideo) GetTemplateURL() string {
	if p.TemplateURL == nil {
		return RewardVideo_TemplateURL_DEFAULT
	}
	return *p.TemplateURL
}

var RewardVideo_PausedURL_DEFAULT string

func (p *RewardVideo) GetPausedURL() string {
	if p.PausedURL == nil {
		return RewardVideo_PausedURL_DEFAULT
	}
	return *p.PausedURL
}

var Campaign_PlayableAdsWithoutVideo_DEFAULT int32

func (p *Campaign) GetPlayableAdsWithoutVideo() int32 {
	if p.PlayableAdsWithoutVideo == nil {
		return Campaign_PlayableAdsWithoutVideo_DEFAULT
	}
	return *p.PlayableAdsWithoutVideo
}

var Campaign_VideoURL_DEFAULT string

func (p *Campaign) GetVideoURL() string {
	if p.VideoURL == nil {
		return Campaign_VideoURL_DEFAULT
	}
	return *p.VideoURL
}

func (p *Campaign) GetCampaignId() string {
	return p.CampaignId
}

var Campaign_AppName_DEFAULT string

func (p *Campaign) GetAppName() string {
	if p.AppName == nil {
		return Campaign_AppName_DEFAULT
	}
	return *p.AppName
}

var Campaign_ImageURL_DEFAULT string

func (p *Campaign) GetImageURL() string {
	if p.ImageURL == nil {
		return Campaign_ImageURL_DEFAULT
	}
	return *p.ImageURL
}

var Campaign_ClickURL_DEFAULT string

func (p *Campaign) GetClickURL() string {
	if p.ClickURL == nil {
		return Campaign_ClickURL_DEFAULT
	}
	return *p.ClickURL
}

var Campaign_LinkType_DEFAULT int32

func (p *Campaign) GetLinkType() int32 {
	if p.LinkType == nil {
		return Campaign_LinkType_DEFAULT
	}
	return *p.LinkType
}

var Campaign_VideoLength_DEFAULT int32

func (p *Campaign) GetVideoLength() int32 {
	if p.VideoLength == nil {
		return Campaign_VideoLength_DEFAULT
	}
	return *p.VideoLength
}

func (p *AdTracking) GetClick() []string {
	return p.Click
}

func (p *QueryResult_) GetRandValue() int32 {
	return p.RandValue
}

func (p *QueryResult_) GetFlowTagId() int32 {
	return p.FlowTagId
}

func (p *QueryResult_) GetAdBackendConfig() string {
	return p.AdBackendConfig
}

var Campaign_OfferType_DEFAULT int32

func (p *Campaign) GetOfferType() int32 {
	if p.OfferType == nil {
		return Campaign_OfferType_DEFAULT
	}
	return *p.OfferType
}

func (p *Campaign) GetUrlTemplate() string {
	if p.UrlTemplate == nil {
		return ""
	}
	return *p.UrlTemplate
}

func (p *Campaign) GetHtmlTemplate() string {
	if p.HtmlTemplate == nil {
		return ""
	}
	return *p.HtmlTemplate
}

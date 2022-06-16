package params

type Ad struct {
	CampaignId              int64             `json:"id"`
	OfferID                 int               `json:"-"`
	AppName                 string            `json:"title"`
	AppDesc                 string            `json:"desc"`
	PackageName             string            `json:"package_name"`
	IconURL                 string            `json:"icon_url"`
	ImageURL                string            `json:"image_url"`
	ImageSize               string            `json:"image_size"`
	ImpressionURL           string            `json:"impression_url"`
	VideoURL                string            `json:"video_url"`
	VideoLength             int               `json:"video_length"`
	VideoSize               int               `json:"video_size"`
	VideoResolution         string            `json:"video_resolution"`
	VideoEndType            int32             `json:"video_end_type,omitempty"`
	PlayableAdsWithoutVideo int               `json:"playable_ads_without_video,omitempty"`
	EndcardUrl              string            `json:"endcard_url,omitempty"`
	WatchMile               int               `json:"watch_mile"`
	CType                   int               `json:"ctype"`
	AdvImp                  []CAdvImp         `json:"adv_imp,omitempty"`
	AdURLList               []string          `json:"ad_url_list,omitempty"`
	TImp                    int               `json:"t_imp"`
	AdvID                   int               `json:"adv_id"`
	ClickURL                string            `json:"click_url"`
	NoticeURL               string            `json:"notice_url"`
	Price                   float64           `json:"-"`
	OfferName               string            `json:"offer_name,omitempty"`
	InstallToken            string            `json:"install_token,omitempty"`
	FCA                     int               `json:"fca"`
	FCB                     int32             `json:"fcb"`
	Template                int               `json:"template"`
	BigTemlate              string            `json:"big_template"`
	AdSourceID              int               `json:"ad_source_id"`
	AppSize                 string            `json:"app_size"`
	ClickMode               int               `json:"click_mode"`
	Rating                  float64           `json:"rating"`
	LandingType             int               `json:"landing_type"`
	CtaText                 string            `json:"ctatext"`
	ClickCacheTime          int               `json:"c_ct"`
	CampaignType            int32             `json:"link_type"`
	Guidelines              string            `json:"guidelines"`
	RewardAmount            int               `json:"reward_amount"`
	RewardName              string            `json:"reward_name"`
	OfferType               int               `json:"offer_type,omitempty"`
	RetargetOffer           int               `json:"retarget_offer,omitempty"`
	StatsURL                string            `json:"stats_url,omitempty"`
	AdTracking              CAdTracking       `json:"-"`
	AdTrackingPoint         *CAdTracking      `json:"ad_tracking,omitempty"`
	LoopBack                map[string]string `json:"loopback,omitempty"`
	Rv                      RV                `json:"-"`
	RvPoint                 *RV               `json:"rv,omitempty"`
	Storekit                *int              `json:"storekit,omitempty"`
	Md5File                 *string           `json:"md5_file,omitempty"`
	GifURL                  string            `json:"gif_url,omitempty"`
	NumberRating            int               `json:"number_rating,omitempty"`
	IconMime                string            `json:"icon_mime,omitempty"`
	IconResolution          string            `json:"icon_resolution,omitempty"`
	ImageMime               string            `json:"image_mime,omitempty"`
	ImageResolution         string            `json:"image_resolution,omitempty"`
	VideoWidth              int32             `json:"video_width,omitempty"`
	VideoHeight             int32             `json:"video_height,omitempty"`
	Bitrate                 int32             `json:"bitrate,omitempty"`
	VideoMime               string            `json:"video_mime,omitempty"`
	NVT2                    int32             `json:"nv_t2,omitempty"`
	SubCategoryName         []string          `json:"sub_category_name,omitempty"`
	StoreKitTime            int32             `json:"storekit_time,omitempty"`
	EndcardClickResult      int32             `json:"endcard_click_result,omitempty"`
	CToi                    int32             `json:"c_toi,omitempty"`
	ImpUa                   int               `json:"imp_ua,omitempty"`
	ClickUa                 int               `json:"c_ua,omitempty"`
	PreviewUrl              string            `json:"-"`
	JMPD                    *int              `json:"jm_pd,omitempty"`
	AKS                     *AdsKey           `json:"aks,omitempty"`
	Category                int32             `json:"-"`
	WxAppId                 string            `json:"gh_id,omitempty"`
	WxPath                  string            `json:"gh_path,omitempty"`
	BindId                  string            `json:"bind_id,omitempty"`
	DeepLink                string            `json:"deep_link,omitempty"`
	ApkVersion              string            `json:"-"`
	ApkMd5                  string            `json:"-"`
	ApkUrl                  string            `json:"-"`
	OfferCType              int               `json:"oc_type,omitempty"`
	OfferCTime              int               `json:"oc_time,omitempty"`
	TokenList               []COfferToken     `json:"t_list,omitempty"`
	AdChoice                *AdChoice         `json:"adchoice,omitempty"`
	Plct                    int               `json:"plct,omitempty"`
	Plctb                   int               `json:"plctb,omitempty"`
	ExtImg                  string            `json:"ext_img,omitempty"`
	CreativeId              int64             `json:"creative_id,omitempty"`
	WithOutInstallCheck     int               `json:"wtick,omitempty"`
	FakeExt                 string            `json:"-"`
	IsDownload              bool              `json:"-"`
	ReadyRate               int               `json:"ready_rate,omitempty"`
	OfferExtData            *OfferExtData     `json:"ext_data,omitempty"`
	ParamQ                  string            `json:"-"`
	ParamR                  string            `json:"-"`
	ParamAL                 string            `json:"-"`
	CamHtml                 string            `json:"cam_html,omitempty"`
	CamTplUrl               string            `json:"cam_tpl_url,omitempty"`
}

type OfferExtData struct {
	SlotId int32 `json:"slot_id"`
}

type AdsKey struct {
	K      *string `json:"k,omitempty"`
	Q      *string `json:"q,omitempty"`
	R      *string `json:"r,omitempty"`
	AL     *string `json:"al,omitempty"`
	CSP    *string `json:"csp,omitempty"`
	MP     *string `json:"mp,omitempty"`
	AdType *string `json:"t,omitempty"`
}

type RV struct {
	VideoTemplate int    `json:"video_template"`
	TemplateUrl   string `json:"template_url"`
	Orientation   int    `json:"orientation"`
	PausedUrl     string `json:"paused_url"`
}

type CAdTracking struct {
	Start            []string        `json:"start,omitempty"`
	First_quartile   []string        `json:"first_quartile,omitempty"`
	Midpoint         []string        `json:"midpoint,omitempty"`
	Third_quartile   []string        `json:"third_quartile,omitempty"`
	Complete         []string        `json:"complete,omitempty"`
	Mute             []string        `json:"mute,omitempty"`
	UnMute           []string        `json:"unmute,omitempty"`
	Impression       []string        `json:"impression,omitempty"`
	Click            []string        `json:"click,omitempty"`
	EndCardShow      []string        `json:"endcard_show,omitempty"`
	Close            []string        `json:"close,omitempty"`
	PlayPercentage   []CPlayTracking `json:"play_percentage,omitempty"`
	Pause            []string        `json:"pause,omitempty"`
	VideoClick       []string        `json:"video_click,omitempty"`
	ImpressionT2     []string        `json:"impression_t2,omitempty"`
	ApkDownloadStart []string        `json:"apk_download_start,omitempty"`
	ApkDownloadEnd   []string        `json:"apk_download_end,omitempty"`
	ApkInstall       []string        `json:"apk_install,omitempty"`
	Fcb              int             `json:"fcb,omitempty"`
	Dropout          []string        `json:"dropout_track,omitempty"`
	Plycmpt          []string        `json:"plycmpt_track,omitempty"`
	PubImp           []string        `json:"pub_imp,omitempty"`
	ExaClick         []string        `json:"exa_click,omitempty"`
	ExaImp           []string        `json:"exa_imp,omitempty"`
	// EndCardShow      []string        `json:"-"`
}

type CPlayTracking struct {
	Rate int    `json:"rate"`
	Url  string `json:"url"`
}

type CRewardVideo struct {
	VideoTemplate int                 `json:"video_template"`
	Orientation   int                 `json:"orientation"`
	TemplateUrl   string              `json:"template_url"`
	PausedUrl     string              `json:"paused_url"`
	Images        map[string][]string `json:"image,omitempty"`
}

type CAdvImp struct {
	Sec int    `json:"sec"`
	Url string `json:"url"`
}

type AdChoice struct {
	AdLogolink   string `json:"ad_logo_link,omitempty"`
	AdchoiceIcon string `json:"adchoice_icon,omitempty"`
	AdchoiceLink string `json:"adchoice_link,omitempty"`
	AdchoiceSize string `json:"adchoice_size,omitempty"`
	AdvLogo      string `json:"adv_logo,omitempty"`
	AdvName      string `json:"adv_name,omitempty"`
	PlatformLogo string `json:"platform_logo,omitempty"`
	PlatformName string `json:"platform_name,omitempty"`
}

type COfferToken struct {
	CacheOfferToken string `bson:"token,omitempty" json:"token,omitempty"`
	CacheOfferTitle string `bson:"titleMd5,omitempty" json:"title,omitempty"`
}

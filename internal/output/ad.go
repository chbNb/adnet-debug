package output

import (
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	smodel "gitlab.mobvista.com/ADN/structs/model"
)

// Ad 定义返回的广告数据
type Ad struct {
	CampaignID              int64             `json:"id"`
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
	VideoEndType            int               `json:"video_end_type,omitempty"`
	PlayableAdsWithoutVideo int               `json:"playable_ads_without_video,omitempty"`
	EndcardUrl              string            `json:"endcard_url,omitempty"`
	WatchMile               int               `json:"watch_mile"`
	CType                   int               `json:"ctype"`
	AdvImp                  []CAdvImp         `json:"adv_imp"`
	AdURLList               []string          `json:"ad_url_list"`
	TImp                    int               `json:"t_imp"`
	AdvID                   int               `json:"adv_id"`
	ClickURL                string            `json:"click_url"`
	NoticeURL               string            `json:"notice_url"`
	Price                   float32           `json:"-"`
	OfferName               string            `json:"offer_name,omitempty"`
	InstallToken            string            `json:"install_token,omitempty"`
	FCA                     int               `json:"fca"`
	FCB                     int               `json:"fcb"`
	Template                int               `json:"template"`
	AdSourceID              int               `json:"ad_source_id"`
	AppSize                 string            `json:"app_size"`
	ClickMode               int               `json:"click_mode"`
	Rating                  float32           `json:"rating"`
	LandingType             int               `json:"landing_type"`
	CtaText                 string            `json:"ctatext"`
	ClickCacheTime          int               `json:"c_ct"`
	CampaignType            int               `json:"link_type"`
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
	Storekit                int               `json:"storekit,omitempty"`
	//Md5File                 string               `json:"md5_file,omitempty"`
	GifURL              string            `json:"gif_url,omitempty"`
	NumberRating        int               `json:"number_rating,omitempty"`
	IconMime            string            `json:"icon_mime,omitempty"`
	IconResolution      string            `json:"icon_resolution,omitempty"`
	ImageMime           string            `json:"image_mime,omitempty"`
	ImageResolution     string            `json:"image_resolution,omitempty"`
	VideoWidth          int32             `json:"video_width,omitempty"`
	VideoHeight         int32             `json:"video_height,omitempty"`
	Bitrate             int32             `json:"bitrate,omitempty"`
	VideoMime           string            `json:"video_mime,omitempty"`
	NVT2                int32             `json:"nv_t2,omitempty"`
	SubCategoryName     []string          `json:"sub_category_name,omitempty"`
	StoreKitTime        int32             `json:"storekit_time,omitempty"`
	EndcardClickResult  int32             `json:"endcard_click_result,omitempty"`
	CToi                int32             `json:"c_toi,omitempty"`
	ImpUa               int               `json:"imp_ua,omitempty"`
	CUA                 int               `json:"c_ua,omitempty"`
	PreviewUrl          string            `json:"-"`
	JMPD                *int              `json:"jm_pd,omitempty"`
	AKS                 map[string]string `json:"aks,omitempty"`
	Category            int32             `json:"-"`
	WxAppId             string            `json:"gh_id,omitempty"`
	WxPath              string            `json:"gh_path,omitempty"`
	BindId              string            `json:"bind_id,omitempty"`
	DeepLink            string            `json:"deep_link,omitempty"`
	ApkVersion          string            `json:"-"`
	ApkMd5              string            `json:"-"`
	ApkUrl              string            `json:"-"`
	AdChoice            *AdChoice         `json:"adchoice,omitempty"`
	Plct                int               `json:"plct,omitempty"`
	Plctb               int               `json:"plctb,omitempty"`
	ExtImg              string            `json:"ext_img,omitempty"`
	CreativeId          int64             `json:"creative_id,omitempty"`
	ParamK              string            `json:"-"`
	OmSDK               []mvutil.OmSDK    `json:"omid,omitempty"`
	Mraid               string            `json:"mraid,omitempty"`
	ReadyRate           int               `json:"ready_rate,omitempty"`
	OfferExtData        *OfferExtData     `json:"ext_data,omitempty"`
	CamHtml             string            `json:"cam_html,omitempty"`
	CamTplUrl           string            `json:"cam_tpl_url,omitempty"`
	ApkAlt              int               `json:"apk_alt,omitempty"`
	Skadnetwork         *Skadnetwork      `json:"skad,omitempty"`
	OnlineApiBidPrice   float64           `json:"-"`
	UserActivation      bool              `json:"user_activation,omitempty"`
	AdHtml              string            `json:"ad_html,omitempty"`
	WithOutInstallCheck int               `json:"wtick,omitempty"`
	SkImp               *SkImp            `json:"skimp,omitempty"`
	Pcm                 *Pcm              `json:"pcm,omitempty"`
	ApkFMd5             string            `json:"akdlui,omitempty"` // apk download unique id
	Ntbarpt             int               `json:"ntbarpt,omitempty"`
	Ntbarpasbl          int               `json:"ntbarpasbl,omitempty"`
	AtatType            int               `json:"atat_type,omitempty"`
	RewardPlus          *RewardPlus       `json:"rw_pl,omitempty"`
	ApkInfo             *ApkInfo          `json:"apk_info,omitempty"`
	Maitve              int32             `json:"maitve,omitempty"`
	MaitveSrc           string            `json:"maitve_src,omitempty"`
	Floatball           int32             `json:"flb,omitempty"`
	FloatballSkipTime   int32             `json:"flb_skiptime,omitempty"`
	HBOnlineBidPrice    float64           `json:"bid_price,omitempty"`
	VideoCtnType        int32             `json:"vctn_t,omitempty"`
	VideoCheckType      int32             `json:"vck_t,omitempty"`
	RsIgnoreCheckRule   []int32           `json:"rs_ignc_r,omitempty"`
	ViewCompletedTime   int               `json:"view_com_time,omitempty"`
	AdspaceType         int32             `json:"adspace_t,omitempty"`
	CloseButtonDelay    *int32            `json:"cbd,omitempty"`
	VideoSkipTime       *int32            `json:"vst,omitempty"`
	ThirdPartyOffer     int               `json:"tp_offer,omitempty"`
	FilterAutoClick     int               `json:"fac,omitempty"`
}

type ApkInfo struct {
	AppName                   string   `json:"app_name,omitempty"`
	SensitivePermission       []string `json:"perm_desc,omitempty"`
	OriginSensitivePermission []string `json:"ori_perm_desc,omitempty"`
	PrivacyUrl                string   `json:"pri_url,omitempty"`
	AppVersionUpdateTime      string   `json:"upd_time,omitempty"`
	AppVersion                string   `json:"app_ver,omitempty"`
	DeveloperName             string   `json:"dev_name,omitempty"`
}

type RewardPlus struct {
	CurrencyId         int    `json:"currency_id"`
	CurrencyDesc       string `json:"virtual_currency"`
	CurrencyName       string `json:"name"`
	CurrencyReward     int    `json:"amount"`
	CurrencyRewardPlus int    `json:"amount_max"`
	CurrencyCbType     int    `json:"callback_rule"`
	CurrencyIcon       string `json:"icon"`
}

type OfferExtData struct {
	SlotId     int32  `json:"slot_id"`
	GifWidgets string `json:"gif_wgs,omitempty"`
}

type AdsKey struct {
	K   *string `json:"k,omitempty"`
	Q   *string `json:"q,omitempty"`
	R   *string `json:"r,omitempty"`
	AL  *string `json:"al,omitempty"`
	CSP *string `json:"csp,omitempty"`
	MP  *string `json:"mp,omitempty"`
}

type RV struct {
	VideoTemplate int    `json:"video_template"`
	TemplateUrl   string `json:"template_url"`
	Orientation   int    `json:"orientation"`
	PausedUrl     string `json:"paused_url"`
	Image         *Image `json:"image,omitempty"`
}

type Image struct {
	IdcdImg []string `json:"idcd_img,omitempty"`
}

type CAdTracking struct {
	Start            []string        `json:"start,omitempty"`
	First_quartile   []string        `json:"first_quartile,omitempty"`
	Midpoint         []string        `json:"midpoint,omitempty"`
	Third_quartile   []string        `json:"third_quartile,omitempty"`
	Complete         []string        `json:"complete,omitempty"`
	Mute             []string        `json:"mute,omitempty"`
	Unmute           []string        `json:"unmute,omitempty"`
	Impression       []string        `json:"impression,omitempty"`
	Click            []string        `json:"click,omitempty"`
	Endcard_show     []string        `json:"endcard_show,omitempty"`
	Close            []string        `json:"close,omitempty"`
	Play_percentage  []CPlayTracking `json:"play_percentage,omitempty"`
	Pause            []string        `json:"pause,omitempty"`
	Video_Click      []string        `json:"video_click,omitempty"`
	Impression_t2    []string        `json:"impression_t2,omitempty"`
	ApkDownloadStart []string        `json:"apk_download_start,omitempty"`
	ApkDownloadEnd   []string        `json:"apk_download_end,omitempty"`
	ApkInstall       []string        `json:"apk_install,omitempty"`
	Dropout          []string        `json:"dropout_track,omitempty"`
	Plycmpt          []string        `json:"plycmpt_track,omitempty"`
	PubImp           []string        `json:"pub_imp,omitempty"`
	ExaClick         []string        `json:"exa_click,omitempty"`
	ExaImp           []string        `json:"exa_imp,omitempty"`
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

//type CLoopBack struct {
//	Domain string `json:"domain,omitempty"`
//	Key    string `json:"key,omitempty"`
//	Value  string `json:"value,omitempty"`
//}

type MobvistaData struct {
	SessionID           string            `json:"session_id,omitempty"`
	NewVersionSessionID string            `json:"a,omitempty"`
	EncryptedSysId      string            `json:"b,omitempty"`
	EncryptedBkupId     string            `json:"c,omitempty"`
	ParentSessionID     string            `json:"parent_session_id"`
	AdType              int               `json:"ad_type"`
	Template            int               `json:"template"`
	UnitSize            string            `json:"unit_size"`
	Ads                 []Ad              `json:"ads"`
	HTMLURL             string            `json:"html_url"`
	EndScreenURL        string            `json:"end_screen_url"`
	OnlyImpressionURL   string            `json:"only_impression_url,omitempty"`
	IAIcon              *string           `json:"ia_icon,omitempty"`
	IARst               *int              `json:"ia_rst,omitempty"`
	IAUrl               *string           `json:"ia_url,omitempty"`
	IAOri               *int              `json:"ia_ori,omitempty"`
	RKS                 map[string]string `json:"rks,omitempty"`
	Setting             *Setting          `json:"setting,omitempty"`
	BannerUrl           string            `json:"banner_url,omitempty"`
	BannerHtml          string            `json:"banner_html,omitempty"`
	Nscpt               int               `json:"nscpt,omitempty"`
	BigTemplateUrl      string            `json:"mof_template_url,omitempty"`
	BigTemplateId       int64             `json:"mof_tplid,omitempty"`
	PvUrls              []string          `json:"pv_urls,omitempty"`
	AdTplUrl            string            `json:"ad_tpl_url,omitempty"`
	AdHtml              string            `json:"ad_html,omitempty"`
	RequestExtData      *RequestExtData   `json:"req_ext_data,omitempty"`
	Vcn                 int32             `json:"vcn,omitempty"`
	TokenRule           int32             `json:"token_r,omitempty"`
	EncryptPrice        string            `json:"encrypt_p,omitempty"`
}

type Setting struct {
	AdSourceTime         *map[string]int64 `json:"ad_source_time,omitempty"`
	CookieAchieve        *int              `json:"cookie_achieve,omitempty"`
	Offset               *int              `json:"offset,omitempty"`
	Autoplay             *int              `json:"autoplay,omitempty"`
	Clicktype            *int              `json:"click_type,omitempty"`
	DLNet                *int              `json:"dlnet,omitempty"`
	ShowImage            *int              `json:"show_image,omitempty"`
	Reward               *[]smodel.Reward  `json:"reward,omitempty"`
	RecallNet            *string           `json:"recall_net,omitempty"`
	Hang                 *int              `json:"hang,omitempty"`
	EndcardTemplate      *string           `json:"endcard_template,omitempty"`
	IsIncent             *int              `json:"is_incent,omitempty"`
	IsServerCall         *int              `json:"is_server_call,omitempty"`
	VideoSkipTime        *int              `json:"video_skip_time,omitempty"`
	DailyPlayCap         *int              `json:"daily_play_cap,omitempty"`
	Orientation          *int              `json:"orientation,omitempty"`
	CloseButtonDelay     *int              `json:"close_button_delay,omitempty"`
	OffsetMax            *int              `json:"offset_max,omitempty"`
	Plct                 *int              `json:"plct,omitempty"`
	VideoInteractiveType *int              `json:"video_interactive_type,omitempty"`
	MuteMode             *int              `json:"mute_mode,omitempty"`
	IsReady              *int              `json:"is_ready,omitempty"`
	ApiCacheNum          *int              `json:"vcn,omitempty"`
	IconImg              *string           `json:"icon_img,omitempty"`
	IconTitle            *string           `json:"icon_t,omitempty"`
	GifCtd               *int              `json:"gif_ctd,omitempty"`
	RefreshFq            *int              `json:"refresh_fq,omitempty"`
}
type MobvistaResult struct {
	Status       int           `json:"status"`
	Msg          string        `json:"msg"`
	Data         MobvistaData  `json:"data"`
	DebugInfo    []interface{} `json:"debuginfo,omitempty"`
	AsDebugInfo  interface{}   `json:"asdebuginfo,omitempty"`
	MasDebugInfo interface{}   `json:"masdebuginfo,omitempty"`
	Version      string        `json:"version,omitempty"` //接口返回的结构体版本（v3/v5)【注：只有v5接口才返回】
}

type ErrorInfo struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
}

type RequestKey struct {
	SH *string `json:"sh,omitempty"`
	DO *string `json:"do,omitempty"`
	Z  *string `json:"z,omitempty"`
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

type Skadnetwork struct {
	Version         string `json:"ver,omitempty"`
	Network         string `json:"nid,omitempty"`
	AppleCampaignId string `json:"cid,omitempty"`
	Targetid        string `json:"targetid,omitempty"`
	Nonce           string `json:"nonce,omitempty"`
	Sourceid        string `json:"sourceid,omitempty"`
	Timestamp       string `json:"tmp,omitempty"`
	Sign            string `json:"sign,omitempty"`
	Need            int    `json:"need,omitempty"`
}

type SkImp struct {
	ViewSign      string `json:"view_sign,omitempty"`
	AdType        string `json:"adtype,omitempty"`
	AdDesc        string `json:"ad_desc,omitempty"`
	PurchaserName string `json:"purchaser_n,omitempty"`
}

type Pcm struct {
	SourceId          int    `json:"source_id,omitempty"`
	DestinationUrl    string `json:"dest_url,omitempty"`
	SourceDescription string `json:"source_desc,omitempty"`
	Purchaser         string `json:"purchaser,omitempty"`
}

type RequestExtData struct {
	ParentId         string `json:"parent_id,omitempty"`
	MofRequestDomain string `json:"mof_domain,omitempty"`
}

package params

type AdData struct {
	SessionID         string      `json:"session_id"`
	ParentSessionID   string      `json:"parent_session_id"`
	AdType            int         `json:"ad_type"`
	Template          int         `json:"template"`
	UnitSize          string      `json:"unit_size"`
	Ads               []Ad        `json:"ads"`
	HTMLURL           string      `json:"html_url"`
	EndScreenURL      string      `json:"end_screen_url"`
	OnlyImpressionURL string      `json:"only_impression_url,omitempty"`
	IAIcon            *string     `json:"ia_icon,omitempty"`
	IARst             *int        `json:"ia_rst,omitempty"`
	IAUrl             *string     `json:"ia_url,omitempty"`
	IAOri             *int        `json:"ia_ori,omitempty"`
	RKS               *RequestKey `json:"rks,omitempty"`
	Setting           *Setting    `json:"setting,omitempty"`
	BannerUrl         string      `json:"banner_url,omitempty"`
	BannerHtml        string      `json:"banner_html,omitempty"`
	SplashAdUrl       string      `json:"ad_tpl_url,omitempty"`
	SplashAdHtml      string      `json:"ad_html,omitempty"`
	Nscpt             int         `json:"nscpt,omitempty"`
	BigTemplateUrl    string      `json:"mof_template_url,omitempty"`
	BigTemplateId     int64       `json:"mof_tplid,omitempty"`
	PvUrls            []string    `json:"pv_urls,omitempty"`
}

type Setting struct {
	AdSourceTime         *map[string]int64 `json:"ad_source_time,omitempty"`
	CookieAchieve        *int              `json:"cookie_achieve,omitempty"`
	Offset               *int              `json:"offset,omitempty"`
	Autoplay             *int              `json:"autoplay,omitempty"`
	Clicktype            *int              `json:"click_type,omitempty"`
	DLNet                *int              `json:"dlnet,omitempty"`
	ShowImage            *int              `json:"show_image,omitempty"`
	Reward               *[]Reward         `json:"reward,omitempty"`
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
}

type MobvistaResult struct {
	Status    int           `json:"status"`
	Msg       string        `json:"msg"`
	Data      *AdData       `json:"data,omitempty"`
	DebugInfo []interface{} `json:"debuginfo,omitempty"`
}

type Reward struct {
	ID     int64  `bson:"id,omitempty" json:"id,omitempty"`
	Name   string `bson:"name,omitempty" json:"name,omitempty"`
	Amount int64  `bson:"amount,omitempty" json:"amount,omitempty"`
}

type RequestKey struct {
	SH *string `json:"sh,omitempty"`
	DO *string `json:"do,omitempty"`
	Z  *string `json:"z,omitempty"`
}

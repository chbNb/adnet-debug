package output

import (
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
)

type V5Result struct {
	Status       int           `json:"status"`
	Msg          string        `json:"msg"`
	Data         V5Data        `json:"data"`
	DebugInfo    []interface{} `json:"debuginfo,omitempty"`
	AsDebugInfo  interface{}   `json:"asdebuginfo,omitempty"`
	MasDebugInfo interface{}   `json:"masdebuginfo,omitempty"`
	Version      string        `json:"version"` //接口返回的结构体版本（v3/v5)
}

type V5Data struct {
	SessionID           string                 `json:"session_id,omitempty"`
	NewVersionSessionID string                 `json:"a,omitempty"`
	EncryptedSysId      string                 `json:"b,omitempty"`
	EncryptedBkupId     string                 `json:"c,omitempty"`
	AdType              int                    `json:"ad_type"`
	Template            int                    `json:"template"`
	UnitSize            string                 `json:"unit_size"`
	Ads                 []*V5Ads               `json:"ads"`
	ReplaceTmp          map[string]interface{} `json:"replace_tmp"`
	HTMLURL             string                 `json:"html_url"`
	EndScreenURL        string                 `json:"end_screen_url"`
	OnlyImpressionURL   string                 `json:"only_impression_url,omitempty"`
	IAIcon              *string                `json:"ia_icon,omitempty"`
	IARst               *int                   `json:"ia_rst,omitempty"`
	IAUrl               *string                `json:"ia_url,omitempty"`
	IAOri               *int                   `json:"ia_ori,omitempty"`
	RKS                 map[string]string      `json:"rks,omitempty"`
	Setting             *Setting               `json:"setting,omitempty"`
	BannerUrl           string                 `json:"banner_url,omitempty"`
	BannerHtml          string                 `json:"banner_html,omitempty"`
	Nscpt               int                    `json:"nscpt,omitempty"`
	BigTemplateUrl      string                 `json:"mof_template_url,omitempty"`
	BigTemplateId       int64                  `json:"mof_tplid,omitempty"`
	PvUrls              []string               `json:"pv_urls,omitempty"`
	AdTplUrl            string                 `json:"ad_tpl_url,omitempty"`
	AdHtml              string                 `json:"ad_html,omitempty"`
	RequestExtData      *RequestExtData        `json:"req_ext_data,omitempty"`
	Vcn                 int32                  `json:"vcn,omitempty"`
	TokenRule           int32                  `json:"token_r,omitempty"`
	EncryptPrice        string                 `json:"encrypt_p,omitempty"`
}

// =========replace obj 定义，json的key不能有重复的===================
type BaseObj struct {
	AppName     string            `json:"title"`
	AppDesc     string            `json:"desc"`
	PackageName string            `json:"package_name"`
	AppSize     string            `json:"app_size"`
	Rating      float32           `json:"rating"`
	CType       int               `json:"ctype"`
	LoopBack    map[string]string `json:"loopback,omitempty"`
	CtaText     string            `json:"ctatext"`
	ApkAlt      int               `json:"apk_alt,omitempty"`
	ApkInfo     *ApkInfo          `json:"apk_info,omitempty"`
	Maitve      int32             `json:"maitve,omitempty"`
	MaitveSrc   string            `json:"maitve_src,omitempty"`
}

type ImageObj struct {
	IconURL   string `json:"icon_url"`
	ImageURL  string `json:"image_url"`
	ImageSize string `json:"image_size"`
	ExtImg    string `json:"ext_img,omitempty"`
}

type UnitSettingObj struct {
	RewardPlus *RewardPlus `json:"rw_pl,omitempty"`
}

type VideoObj struct {
	VideoURL        string `json:"video_url"`
	VideoLength     int    `json:"video_length"`
	VideoSize       int    `json:"video_size"`
	VideoResolution string `json:"video_resolution"`
	Md5File         string `json:"md5_file,omitempty"`
}

type TkBaseObj struct {
	NoticeURL     string `json:"notice_url"`
	ImpressionURL string `json:"impression_url"`
}

type CSTObj struct {
	FCA                     int     `json:"fca"`
	FCB                     int     `json:"fcb"`
	ClickCacheTime          int     `json:"c_ct"`
	VideoEndType            int     `json:"video_end_type,omitempty"`
	PlayableAdsWithoutVideo int     `json:"playable_ads_without_video,omitempty"`
	RetargetOffer           int     `json:"retarget_offer,omitempty"`
	Storekit                int     `json:"storekit,omitempty"`
	ImpUa                   int     `json:"imp_ua,omitempty"`
	CUA                     int     `json:"c_ua,omitempty"`
	StoreKitTime            int32   `json:"storekit_time,omitempty"`
	EndcardClickResult      int32   `json:"endcard_click_result,omitempty"`
	CToi                    int32   `json:"c_toi,omitempty"`
	NVT2                    int32   `json:"nv_t2,omitempty"`
	Plct                    int     `json:"plct,omitempty"`
	Plctb                   int     `json:"plctb,omitempty"`
	JMPD                    *int    `json:"jm_pd,omitempty"`
	ReadyRate               int     `json:"ready_rate,omitempty"`
	WithOutInstallCheck     int     `json:"wtick,omitempty"`
	ApkFMd5                 string  `json:"akdlui,omitempty"` // apk download unique id
	Ntbarpt                 int     `json:"ntbarpt,omitempty"`
	Ntbarpasbl              int     `json:"ntbarpasbl,omitempty"`
	AtatType                int     `json:"atat_type,omitempty"`
	Floatball               int32   `json:"flb,omitempty"`
	FloatballSkipTime       int32   `json:"flb_skiptime,omitempty"`
	VideoCtnType            int32   `json:"vctn_t,omitempty"`
	VideoCheckType          int32   `json:"vck_t,omitempty"`
	RsIgnoreCheckRule       []int32 `json:"rs_ignc_r,omitempty"`
	ViewCompletedTime       int     `json:"view_com_time,omitempty"`
	AdspaceType             int32   `json:"adspace_t,omitempty"`
	CloseButtonDelay        *int32  `json:"cbd,omitempty"`
	VideoSkipTime           *int32  `json:"vst,omitempty"`
	ThirdPartyOffer         int     `json:"tp_offer,omitempty"`
	FilterAutoClick         int     `json:"fac,omitempty"`
}

type TrackingObj struct {
	AdTrackingPoint *CAdTracking `json:"ad_tracking,omitempty"`
}

type OwSpecObj struct {
	Guidelines   string `json:"guidelines"`
	RewardAmount int    `json:"reward_amount"`
	RewardName   string `json:"reward_name"`
}

type RvPoitObj struct {
	RvPoint *RV `json:"rv,omitempty"`
}

type WxObj struct {
	WxAppId string `json:"gh_id,omitempty"`
	WxPath  string `json:"gh_path,omitempty"`
	BindId  string `json:"bind_id,omitempty"`
}

type AdChoiceObj struct {
	AdChoice *AdChoice `json:"adchoice,omitempty"`
}

type OmSDKObj struct {
	OmSDK []mvutil.OmSDK `json:"omid,omitempty"`
}

type SkadnetworkObj struct {
	Skadnetwork    *Skadnetwork `json:"skad,omitempty"`
	ImpSkadnetwork *SkImp       `json:"skimp,omitempty"`
}

type PcmObj struct {
	Pcm *Pcm `json:"pcm,omitempty"`
}

type V5Ads struct {
	//V5模板ID
	TmpIds []string `json:"tmp_ids"`

	//待确认
	//WatchMile int       `json:"watch_mile"`
	AdvImp    []CAdvImp `json:"adv_imp"` //如果没有不返回
	AdURLList []string  `json:"ad_url_list"`
	Template  int       `json:"template"`

	//确认无用删除
	//InstallToken    string               `json:"install_token,omitempty"`
	//StatsURL        string               `json:"stats_url,omitempty"`
	//IconMime        string               `json:"icon_mime,omitempty"`
	//IconResolution  string               `json:"icon_resolution,omitempty"`
	//ImageMime       string               `json:"image_mime,omitempty"`
	//ImageResolution string               `json:"image_resolution,omitempty"`
	//VideoWidth      int32                `json:"video_width,omitempty"`
	//VideoHeight     int32                `json:"video_height,omitempty"`
	//VideoMime       string               `json:"video_mime,omitempty"`
	//Bitrate         int32                `json:"bitrate,omitempty"`
	//SubCategoryName []string             `json:"sub_category_name,omitempty"`
	//OfferCType      int                  `json:"oc_type,omitempty"`
	//OfferCTime      int                  `json:"oc_time,omitempty"`
	//TokenList       []mvutil.COfferToken `json:"t_list,omitempty"`

	//每次必返回区域
	CampaignID     int64             `json:"id"`
	EndcardUrl     string            `json:"endcard_url,omitempty"`
	ClickURL       string            `json:"click_url"`
	OfferName      string            `json:"offer_name,omitempty"`
	ClickMode      int               `json:"click_mode"`
	LandingType    int               `json:"landing_type"`
	CampaignType   int               `json:"link_type"`
	GifURL         string            `json:"gif_url,omitempty"`
	NumberRating   int               `json:"number_rating,omitempty"`
	DeepLink       string            `json:"deep_link,omitempty"`
	CreativeId     int64             `json:"creative_id,omitempty"`
	Mraid          string            `json:"mraid,omitempty"`
	OfferExtData   *OfferExtData     `json:"ext_data,omitempty"`
	CamHtml        string            `json:"cam_html,omitempty"`
	CamTplUrl      string            `json:"cam_tpl_url,omitempty"`
	AKS            map[string]string `json:"aks,omitempty"`
	UserActivation bool              `json:"user_activation,omitempty"`
}

func RenderV5Res(mr MobvistaResult, r *mvutil.RequestParams) V5Result {
	var result V5Result

	result.Status = 1
	result.Msg = "success"
	result.DebugInfo = mr.DebugInfo
	result.AsDebugInfo = mr.AsDebugInfo
	result.MasDebugInfo = mr.MasDebugInfo
	result.Version = "v5"
	//Data
	result.Data.SessionID = mr.Data.SessionID
	result.Data.NewVersionSessionID = mr.Data.NewVersionSessionID
	result.Data.EncryptedSysId = mr.Data.EncryptedSysId
	result.Data.EncryptedBkupId = mr.Data.EncryptedBkupId
	result.Data.AdType = mr.Data.AdType
	result.Data.Template = mr.Data.Template
	result.Data.UnitSize = mr.Data.UnitSize
	result.Data.HTMLURL = mr.Data.HTMLURL
	result.Data.EndScreenURL = mr.Data.EndScreenURL
	result.Data.OnlyImpressionURL = mr.Data.OnlyImpressionURL
	result.Data.IAIcon = mr.Data.IAIcon
	result.Data.IARst = mr.Data.IARst
	result.Data.IAUrl = mr.Data.IAUrl
	result.Data.IAOri = mr.Data.IAOri
	result.Data.RKS = mr.Data.RKS
	result.Data.Setting = mr.Data.Setting
	result.Data.BannerUrl = mr.Data.BannerUrl
	result.Data.BannerHtml = mr.Data.BannerHtml
	result.Data.Nscpt = mr.Data.Nscpt
	result.Data.BigTemplateUrl = mr.Data.BigTemplateUrl
	result.Data.BigTemplateId = mr.Data.BigTemplateId
	result.Data.PvUrls = mr.Data.PvUrls
	result.Data.AdTplUrl = mr.Data.AdTplUrl
	result.Data.AdHtml = mr.Data.AdHtml
	result.Data.RequestExtData = mr.Data.RequestExtData
	result.Data.Vcn = mr.Data.Vcn
	result.Data.TokenRule = mr.Data.TokenRule
	result.Data.EncryptPrice = mr.Data.EncryptPrice

	result.Data.ReplaceTmp = make(map[string]interface{})

	result.Data.Ads = make([]*V5Ads, 0, len(mr.Data.Ads))

	for _, ad := range mr.Data.Ads {
		v5ads, tmpMap := renderObj(ad, r.Param.TmpIds)
		result.Data.Ads = append(result.Data.Ads, v5ads)
		for k, v := range tmpMap {
			if _, ok := result.Data.ReplaceTmp[k]; !ok { //不存在的就要加入
				result.Data.ReplaceTmp[k] = v
			}
		}
	}

	return result
}

func addTmpMap(tmpMap map[string]interface{}, key string, value interface{}, ignoreMap map[string]bool) {
	if _, ok := ignoreMap[key]; !ok {
		tmpMap[key] = value
	}
}

func renderObj(ad Ad, ignoreMap map[string]bool) (*V5Ads, map[string]interface{}) {
	v5ads := &V5Ads{}
	tmpMap := make(map[string]interface{})

	v5ads.AdvImp = ad.AdvImp
	v5ads.AdURLList = ad.AdURLList
	v5ads.Template = ad.Template
	v5ads.CampaignID = ad.CampaignID
	v5ads.EndcardUrl = ad.EndcardUrl
	v5ads.ClickURL = ad.ClickURL
	v5ads.OfferName = ad.OfferName
	v5ads.ClickMode = ad.ClickMode
	v5ads.LandingType = ad.LandingType
	v5ads.CampaignType = ad.CampaignType
	v5ads.GifURL = ad.GifURL
	v5ads.NumberRating = ad.NumberRating
	v5ads.DeepLink = ad.DeepLink
	v5ads.CreativeId = ad.CreativeId
	v5ads.Mraid = ad.Mraid
	v5ads.OfferExtData = ad.OfferExtData
	v5ads.CamHtml = ad.CamHtml
	v5ads.CamTplUrl = ad.CamTplUrl
	v5ads.AKS = ad.AKS
	v5ads.UserActivation = ad.UserActivation

	bObj := BaseObj{}
	bObj.AppName = ad.AppName
	bObj.AppDesc = ad.AppDesc
	bObj.PackageName = ad.PackageName
	bObj.AppSize = ad.AppSize
	bObj.Rating = ad.Rating
	bObj.CType = ad.CType
	bObj.LoopBack = ad.LoopBack
	bObj.CtaText = ad.CtaText
	bObj.ApkAlt = ad.ApkAlt
	bObj.ApkInfo = ad.ApkInfo
	bObj.Maitve = ad.Maitve
	bObj.MaitveSrc = ad.MaitveSrc
	if key, err := mvutil.APHashByObj(bObj); err == nil {
		v5ads.TmpIds = append(v5ads.TmpIds, key)
		addTmpMap(tmpMap, key, bObj, ignoreMap)
	}

	imgObj := ImageObj{}
	imgObj.IconURL = ad.IconURL
	imgObj.ImageURL = ad.ImageURL
	imgObj.ImageSize = ad.ImageSize
	imgObj.ExtImg = ad.ExtImg
	if key, err := mvutil.APHashByObj(imgObj); err == nil {
		v5ads.TmpIds = append(v5ads.TmpIds, key)
		addTmpMap(tmpMap, key, imgObj, ignoreMap)
	}

	videoObj := VideoObj{}
	videoObj.VideoURL = ad.VideoURL
	videoObj.VideoLength = ad.VideoLength
	videoObj.VideoSize = ad.VideoSize
	videoObj.VideoResolution = ad.VideoResolution
	//videoObj.Md5File = ad.Md5File
	if key, err := mvutil.APHashByObj(videoObj); err == nil {
		v5ads.TmpIds = append(v5ads.TmpIds, key)
		addTmpMap(tmpMap, key, videoObj, ignoreMap)
	}

	tkBaseObj := TkBaseObj{}
	tkBaseObj.NoticeURL = ad.NoticeURL
	tkBaseObj.ImpressionURL = ad.ImpressionURL
	if key, err := mvutil.APHashByObj(tkBaseObj); err == nil {
		v5ads.TmpIds = append(v5ads.TmpIds, key)
		addTmpMap(tmpMap, key, tkBaseObj, ignoreMap)
	}

	cstObj := CSTObj{}
	cstObj.FCA = ad.FCA
	cstObj.FCB = ad.FCB
	cstObj.ClickCacheTime = ad.ClickCacheTime
	cstObj.VideoEndType = ad.VideoEndType
	cstObj.PlayableAdsWithoutVideo = ad.PlayableAdsWithoutVideo
	cstObj.RetargetOffer = ad.RetargetOffer
	cstObj.Storekit = ad.Storekit
	cstObj.ImpUa = ad.ImpUa
	cstObj.CUA = ad.CUA
	cstObj.StoreKitTime = ad.StoreKitTime
	cstObj.EndcardClickResult = ad.EndcardClickResult
	cstObj.CToi = ad.CToi
	cstObj.NVT2 = ad.NVT2
	cstObj.Plct = ad.Plct
	cstObj.Plctb = ad.Plctb
	cstObj.JMPD = ad.JMPD
	cstObj.ReadyRate = ad.ReadyRate
	cstObj.WithOutInstallCheck = ad.WithOutInstallCheck
	cstObj.ApkFMd5 = ad.ApkFMd5
	cstObj.Ntbarpt = ad.Ntbarpt
	cstObj.Ntbarpasbl = ad.Ntbarpasbl
	cstObj.AtatType = ad.AtatType
	cstObj.Floatball = ad.Floatball
	cstObj.FloatballSkipTime = ad.FloatballSkipTime
	cstObj.VideoCtnType = ad.VideoCtnType
	cstObj.VideoCheckType = ad.VideoCheckType
	cstObj.RsIgnoreCheckRule = ad.RsIgnoreCheckRule
	cstObj.ViewCompletedTime = ad.ViewCompletedTime
	cstObj.AdspaceType = ad.AdspaceType
	cstObj.CloseButtonDelay = ad.CloseButtonDelay
	cstObj.VideoSkipTime = ad.VideoSkipTime
	cstObj.ThirdPartyOffer = ad.ThirdPartyOffer
	cstObj.FilterAutoClick = ad.FilterAutoClick
	if key, err := mvutil.APHashByObj(cstObj); err == nil {
		v5ads.TmpIds = append(v5ads.TmpIds, key)
		addTmpMap(tmpMap, key, cstObj, ignoreMap)
	}

	if ad.AdTrackingPoint != nil {
		trackingObj := TrackingObj{}
		trackingObj.AdTrackingPoint = ad.AdTrackingPoint
		if key, err := mvutil.APHashByObj(trackingObj); err == nil {
			v5ads.TmpIds = append(v5ads.TmpIds, key)
			addTmpMap(tmpMap, key, trackingObj, ignoreMap)
		}
	}

	owSpecObj := OwSpecObj{}
	owSpecObj.Guidelines = ad.Guidelines
	owSpecObj.RewardAmount = ad.RewardAmount
	owSpecObj.RewardName = ad.RewardName
	if key, err := mvutil.APHashByObj(owSpecObj); err == nil {
		v5ads.TmpIds = append(v5ads.TmpIds, key)
		addTmpMap(tmpMap, key, owSpecObj, ignoreMap)
	}

	if ad.RvPoint != nil {
		rvPoitObj := RvPoitObj{}
		rvPoitObj.RvPoint = ad.RvPoint
		if key, err := mvutil.APHashByObj(rvPoitObj); err == nil {
			v5ads.TmpIds = append(v5ads.TmpIds, key)
			addTmpMap(tmpMap, key, rvPoitObj, ignoreMap)
		}
	}

	if ad.WxAppId != "" || ad.WxPath != "" || ad.BindId != "" {
		wxObj := WxObj{}
		wxObj.WxAppId = ad.WxAppId
		wxObj.WxPath = ad.WxPath
		wxObj.BindId = ad.BindId
		if key, err := mvutil.APHashByObj(wxObj); err == nil {
			v5ads.TmpIds = append(v5ads.TmpIds, key)
			addTmpMap(tmpMap, key, wxObj, ignoreMap)
		}
	}

	if ad.AdChoice != nil {
		adChoiceObj := AdChoiceObj{}
		adChoiceObj.AdChoice = ad.AdChoice
		if key, err := mvutil.APHashByObj(adChoiceObj); err == nil {
			v5ads.TmpIds = append(v5ads.TmpIds, key)
			addTmpMap(tmpMap, key, adChoiceObj, ignoreMap)
		}
	}

	if len(ad.OmSDK) > 0 {
		omSDKObj := OmSDKObj{}
		omSDKObj.OmSDK = ad.OmSDK
		if key, err := mvutil.APHashByObj(omSDKObj); err == nil {
			v5ads.TmpIds = append(v5ads.TmpIds, key)
			addTmpMap(tmpMap, key, omSDKObj, ignoreMap)
		}
	}

	if ad.Pcm != nil {
		pcmObj := PcmObj{}
		pcmObj.Pcm = ad.Pcm
		if key, err := mvutil.APHashByObj(pcmObj); err == nil {
			v5ads.TmpIds = append(v5ads.TmpIds, key)
			addTmpMap(tmpMap, key, pcmObj, ignoreMap)
		}
	}

	var skaObj SkadnetworkObj
	if ad.Skadnetwork != nil {
		skaObj.Skadnetwork = ad.Skadnetwork
	}
	if ad.SkImp != nil {
		skaObj.ImpSkadnetwork = ad.SkImp
	}
	if skaObj.Skadnetwork != nil || skaObj.ImpSkadnetwork != nil {
		if key, err := mvutil.APHashByObj(skaObj); err == nil {
			v5ads.TmpIds = append(v5ads.TmpIds, key)
			addTmpMap(tmpMap, key, skaObj, ignoreMap)
		}
	}

	if ad.RewardPlus != nil {
		unitSettingObj := UnitSettingObj{}
		unitSettingObj.RewardPlus = ad.RewardPlus
		if key, err := mvutil.APHashByObj(unitSettingObj); err == nil {
			v5ads.TmpIds = append(v5ads.TmpIds, key)
			addTmpMap(tmpMap, key, unitSettingObj, ignoreMap)
		}
	}

	return v5ads, tmpMap
}

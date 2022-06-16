package output

import (
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
)

type MPResult struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
	Data   MPData `json:"data"`
}

type MPNewAdResult struct {
	Status int         `json:"status"`
	Msg    string      `json:"msg"`
	Data   MPNewAdData `json:"data"`
}

type MPData struct {
	AdType            int    `json:"fomate"`
	UnitSize          string `json:"p_size"`
	Ads               []MPAd `json:"ad_info"`
	HTMLURL           string `json:"h5url"`
	OnlyImpressionURL string `json:"only_im"`
}

type MPNewAdData struct {
	AdType            int       `json:"adtype"`
	UnitSize          string    `json:"sltsize"`
	Ads               []MPNewAd `json:"adsdetails"`
	HTMLURL           string    `json:"htmlurl"`
	OnlyImpressionURL string    `json:"oly_imp"`
}

type MPAd struct {
	CampaignID      int64         `json:"id"`
	AppName         string        `json:"title,omitempty"`
	AppDesc         string        `json:"body,omitempty"`
	PackageName     string        `json:"package_name"`
	IconURL         string        `json:"icon_a,omitempty"`
	ImageURL        string        `json:"img_a,omitempty"`
	ImageSize       string        `json:"image_size,omitempty"`
	ImpressionURL   string        `json:"impression_url,omitempty"`
	ClickURL        string        `json:"click_url"`
	NoticeURL       string        `json:"n_url"`
	AdSourceID      int           `json:"ads_id"`
	FCA             int           `json:"ica,omitempty"`
	ClickMode       int           `json:"cltype"`
	Rating          float32       `json:"star,omitempty"`
	CampaignType    int           `json:"ltype,omitempty"`
	CtaText         string        `json:"cta,omitempty"`
	ClickCacheTime  int           `json:"uct,omitempty"`
	CType           int           `json:"costmode,omitempty"`
	OfferType       int           `json:"offer_type"`
	AdURLList       []string      `json:"urlgroup,omitempty"`
	AdvImp          []CAdvImp     `json:"impression_adv,omitempty"`
	VideoURL        string        `json:"video_a,omitempty"`
	VideoLength     int           `json:"video_l,omitempty"`
	VideoSize       int           `json:"video_s,omitempty"`
	VideoResolution string        `json:"video_r,omitempty"`
	AdTrackingPoint *MPAdTracking `json:"video_tracking,omitempty"`
	ImpUa           int           `json:"imua,omitempty"`
	CUA             int           `json:"clua,omitempty"`
	JMPD            *int          `json:"ic_dlc,omitempty"`
	RIP             string        `json:"rip"`
	RUA             string        `json:"rua"`
	Storekit        int           `json:"stkit,omitempty"`
}

type MPNewAd struct {
	CampaignID      int64            `json:"id"`
	AppName         string           `json:"title,omitempty"`
	AppDesc         string           `json:"descbody,omitempty"`
	PackageName     string           `json:"pkagename"`
	IconURL         string           `json:"icn_url,omitempty"`
	ImageURL        string           `json:"imga_url,omitempty"`
	ImageSize       string           `json:"image_size,omitempty"`
	ImpressionURL   string           `json:"impressionurl,omitempty"`
	ClickURL        string           `json:"clickurl"`
	NoticeURL       string           `json:"ntic_url"`
	AdSourceID      int              `json:"adsscrid"`
	FCA             int              `json:"showcap,omitempty"`
	ClickMode       int              `json:"clikmode"`
	Rating          float32          `json:"staring,omitempty"`
	CampaignType    int              `json:"lkmode,omitempty"`
	CtaText         string           `json:"ctastring,omitempty"`
	ClickCacheTime  int              `json:"mcct,omitempty"`
	CType           int              `json:"costmode,omitempty"`
	OfferType       int              `json:"offtype"`
	AdURLList       []string         `json:"adurlglist,omitempty"`
	AdvImp          []CAdvImp        `json:"imp_adv_url,omitempty"`
	VideoURL        string           `json:"vid_url,omitempty"`
	VideoLength     int              `json:"vid_l,omitempty"`
	VideoSize       int              `json:"vid_s,omitempty"`
	VideoResolution string           `json:"vid_r,omitempty"`
	AdTrackingPoint *MPNewAdTracking `json:"vid_tk,omitempty"`
	ImpUa           int              `json:"impuag,omitempty"`
	CUA             int              `json:"clkuag,omitempty"`
	JMPD            *int             `json:"ic_dlc,omitempty"`
	RIP             string           `json:"vpkir"`
	RUA             string           `json:"vakur"`
	Storekit        int              `json:"tksit,omitempty"`
}

type MPNewAdTracking struct {
	MPPlayComplete   []MPPlayTracking `json:"plycplt,omitempty"`
	ApkDownloadStart []string         `json:"appdstk,omitempty"`
	ApkDownloadEnd   []string         `json:"appdetk,omitempty"`
	ApkInstall       []string         `json:"appictk,omitempty"`
	ExaClick         []string         `json:"clkexdsn,omitempty"`
	ExaImp           []string         `json:"impexds,omitempty"`
}

type MPAdTracking struct {
	MPPlayComplete   []MPPlayTracking `json:"play_complete,omitempty"`
	ApkDownloadStart []string         `json:"apk_download_start,omitempty"`
	ApkDownloadEnd   []string         `json:"apk_download_end,omitempty"`
	ApkInstall       []string         `json:"apk_install,omitempty"`
	ExaClick         []string         `json:"exa_click,omitempty"`
	ExaImp           []string         `json:"exa_imp,omitempty"`
}

type MPPlayTracking struct {
	Percentage int    `json:"percentage"`
	Url        string `json:"url"`
}

func RenderMPRes(mr MobvistaResult, r *mvutil.RequestParams) MPResult {
	var mpr MPResult
	mpr.Status = 1
	mpr.Msg = "success"
	mpr.Data.AdType = mr.Data.AdType
	mpr.Data.UnitSize = mr.Data.UnitSize
	mpr.Data.HTMLURL = mr.Data.HTMLURL
	mpr.Data.OnlyImpressionURL = mr.Data.OnlyImpressionURL

	var mpadList []MPAd
	for _, v := range mr.Data.Ads {
		mpad := renderMPAd(v, *r)
		mpadList = append(mpadList, mpad)
	}
	mpr.Data.Ads = mpadList

	return mpr
}

func renderMPAd(ad Ad, r mvutil.RequestParams) MPAd {
	var mpad MPAd

	mpad.CampaignID = ad.CampaignID
	mpad.PackageName = ad.PackageName
	mpad.ClickURL = ad.ClickURL
	mpad.NoticeURL = ad.NoticeURL
	mpad.AdSourceID = ad.AdSourceID
	mpad.ClickMode = ad.ClickMode
	mpad.OfferType = ad.OfferType
	mpad.CUA = ad.CUA
	mpad.RIP = ""
	mpad.RUA = ""
	mpad.AppName = ad.AppName
	mpad.AppDesc = ad.AppDesc
	mpad.IconURL = ad.IconURL
	mpad.ImageURL = ad.ImageURL
	mpad.ImageSize = ad.ImageSize
	mpad.ImpressionURL = ad.ImpressionURL
	mpad.FCA = ad.FCA
	mpad.Rating = ad.Rating
	mpad.CampaignType = ad.CampaignType
	mpad.CtaText = ad.CtaText
	mpad.ClickCacheTime = ad.ClickCacheTime
	mpad.CType = ad.CType
	mpad.AdURLList = ad.AdURLList
	mpad.AdvImp = ad.AdvImp
	mpad.VideoURL = ad.VideoURL
	mpad.VideoLength = ad.VideoLength
	mpad.VideoSize = ad.VideoSize
	mpad.VideoResolution = ad.VideoResolution
	mpad.ImpUa = ad.ImpUa
	mpad.JMPD = ad.JMPD
	mpad.Storekit = ad.Storekit

	renderMPAdTracking(&mpad, ad)
	return mpad
}

func renderMPAdTracking(mpad *MPAd, ad Ad) {
	if ad.AdTrackingPoint == nil {
		return
	}
	var percentage MPPlayTracking
	percentage.Percentage = 100
	adTracking := *(ad.AdTrackingPoint)
	for _, v := range adTracking.Play_percentage {
		if v.Rate == 100 {
			percentage.Url = v.Url
		}
	}
	// 兼容返回adsense点击上报，exa_click也为空才不返回tracking
	if len(percentage.Url) == 0 {
		// 对于adsense点击上报，处理没有percentage上报的情况
		if len(ad.AdTracking.ExaClick) > 0 && len(ad.AdTracking.ExaImp) > 0 {
			var mpAdTracking MPAdTracking
			mpAdTracking.ExaClick = ad.AdTracking.ExaClick
			mpAdTracking.ExaImp = ad.AdTracking.ExaImp
			mpad.AdTrackingPoint = &(mpAdTracking)
		}
		return
	}
	percentageList := make([]MPPlayTracking, 1)
	percentageList[0] = percentage
	var mpAdTracking MPAdTracking
	mpAdTracking.MPPlayComplete = percentageList
	mpAdTracking.ApkDownloadStart = ad.AdTracking.ApkDownloadStart
	mpAdTracking.ApkDownloadEnd = ad.AdTracking.ApkDownloadEnd
	mpAdTracking.ApkInstall = ad.AdTracking.ApkInstall
	mpAdTracking.ExaClick = ad.AdTracking.ExaClick
	mpAdTracking.ExaImp = ad.AdTracking.ExaImp
	mpad.AdTrackingPoint = &(mpAdTracking)
}

func RenderMPInfo(ad *Ad) {
	ad.ClickCacheTime = 60
}

func RenderMPNewAdRes(mr MobvistaResult, r *mvutil.RequestParams) MPNewAdResult {
	var mpr MPNewAdResult
	mpr.Status = 1
	mpr.Msg = "success"
	mpr.Data.AdType = mr.Data.AdType
	mpr.Data.UnitSize = mr.Data.UnitSize
	mpr.Data.HTMLURL = mr.Data.HTMLURL
	mpr.Data.OnlyImpressionURL = mr.Data.OnlyImpressionURL

	var mpadList []MPNewAd
	for _, v := range mr.Data.Ads {
		mpad := renderMPNewNormalAd(v, *r)
		mpadList = append(mpadList, mpad)
	}
	mpr.Data.Ads = mpadList

	return mpr
}

func renderMPNewNormalAd(ad Ad, r mvutil.RequestParams) MPNewAd {
	var mpad MPNewAd

	mpad.CampaignID = ad.CampaignID
	mpad.PackageName = ad.PackageName
	mpad.ClickURL = ad.ClickURL
	mpad.NoticeURL = ad.NoticeURL
	mpad.AdSourceID = ad.AdSourceID
	mpad.ClickMode = ad.ClickMode
	mpad.OfferType = ad.OfferType
	mpad.CUA = ad.CUA
	mpad.RIP = ""
	mpad.RUA = ""
	mpad.AppName = ad.AppName
	mpad.AppDesc = ad.AppDesc
	mpad.IconURL = ad.IconURL
	mpad.ImageURL = ad.ImageURL
	mpad.ImageSize = ad.ImageSize
	mpad.ImpressionURL = ad.ImpressionURL
	mpad.FCA = ad.FCA
	mpad.Rating = ad.Rating
	mpad.CampaignType = ad.CampaignType
	mpad.CtaText = ad.CtaText
	mpad.ClickCacheTime = ad.ClickCacheTime
	mpad.CType = ad.CType
	mpad.AdURLList = ad.AdURLList
	mpad.AdvImp = ad.AdvImp
	mpad.VideoURL = ad.VideoURL
	mpad.VideoLength = ad.VideoLength
	mpad.VideoSize = ad.VideoSize
	mpad.VideoResolution = ad.VideoResolution
	mpad.ImpUa = ad.ImpUa
	mpad.JMPD = ad.JMPD
	mpad.Storekit = ad.Storekit

	renderMPNewNormalAdTracking(&mpad, ad)
	return mpad
}

func renderMPNewNormalAdTracking(mpad *MPNewAd, ad Ad) {
	if ad.AdTrackingPoint == nil {
		return
	}
	var percentage MPPlayTracking
	percentage.Percentage = 100
	adTracking := *(ad.AdTrackingPoint)
	for _, v := range adTracking.Play_percentage {
		if v.Rate == 100 {
			percentage.Url = v.Url
		}
	}
	// 兼容返回adsense点击上报，exa_click也为空才不返回tracking
	if len(percentage.Url) == 0 {
		// 对于adsense点击上报，处理没有percentage上报的情况
		if len(ad.AdTracking.ExaClick) > 0 && len(ad.AdTracking.ExaImp) > 0 {
			var mpAdTracking MPNewAdTracking
			mpAdTracking.ExaClick = ad.AdTracking.ExaClick
			mpAdTracking.ExaImp = ad.AdTracking.ExaImp
			mpad.AdTrackingPoint = &(mpAdTracking)
		}
		return
	}
	percentageList := make([]MPPlayTracking, 1)
	percentageList[0] = percentage
	var mpAdTracking MPNewAdTracking
	mpAdTracking.MPPlayComplete = percentageList
	mpAdTracking.ApkDownloadStart = ad.AdTracking.ApkDownloadStart
	mpAdTracking.ApkDownloadEnd = ad.AdTracking.ApkDownloadEnd
	mpAdTracking.ApkInstall = ad.AdTracking.ApkInstall
	mpAdTracking.ExaClick = ad.AdTracking.ExaClick
	mpAdTracking.ExaImp = ad.AdTracking.ExaImp
	mpad.AdTrackingPoint = &(mpAdTracking)
}

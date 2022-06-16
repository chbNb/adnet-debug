package output

import (
	"gitlab.mobvista.com/ADN/adnet/internal/extractor"
	"gitlab.mobvista.com/ADN/adnet/internal/mvconst"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
)

type XMOnlineAd struct {
	CampaignID      int64        `json:"id"`
	AppName         string       `json:"title"`
	AppDesc         string       `json:"desc"`
	PackageName     string       `json:"package_name"`
	IconURL         string       `json:"icon_url"`
	ImageURL        string       `json:"image_url"`
	ImpressionURL   string       `json:"impression_url"`
	ClickURL        string       `json:"click_url"`
	Rating          float32      `json:"rating"`
	CtaText         string       `json:"ctatext"`
	CampaignType    int          `json:"link_type"`
	PreviewUrl      string       `json:"preview_link"`
	ClickMode       int          `json:"click_mode"`
	VideoURL        string       `json:"video_url,omitempty"`
	VideoLength     int          `json:"video_length,omitempty"`
	VideoSize       int          `json:"video_size,omitempty"`
	VideoResolution string       `json:"video_resolution,omitempty"`
	AdTrackingPoint *CAdTracking `json:"ad_tracking,omitempty"`
	DeepLink        string       `json:"deep_link,omitempty"`
}

type XMOnlineResult struct {
	Status    int           `json:"status"`
	Msg       string        `json:"msg"`
	Data      XMOnlineData  `json:"data"`
	DebugInfo []interface{} `json:"debuginfo,omitempty"`
}

type XMOnlineData struct {
	XMAds []XMOnlineAd `json:"ads"`
}

func RenderxmOnlineRes(mr MobvistaResult) XMOnlineResult {
	var or XMOnlineResult
	or.Status = 1
	or.Msg = "success"

	var oadList []XMOnlineAd
	for _, v := range mr.Data.Ads {
		oad := renderXmOnlineAd(v)
		oadList = append(oadList, oad)
	}
	or.Data.XMAds = oadList
	return or
}

func renderXmOnlineAd(ad Ad) XMOnlineAd {
	var xmOad XMOnlineAd
	xmOad.CampaignID = ad.CampaignID
	xmOad.AppName = ad.AppName
	xmOad.AppDesc = ad.AppDesc
	xmOad.PackageName = ad.PackageName
	xmOad.IconURL = ad.IconURL
	xmOad.ImageURL = ad.ImageURL
	xmOad.ImpressionURL = ad.ImpressionURL
	xmOad.ClickURL = ad.ClickURL
	xmOad.Rating = ad.Rating
	xmOad.CtaText = ad.CtaText
	xmOad.CampaignType = ad.CampaignType
	xmOad.PreviewUrl = ad.PreviewUrl
	xmOad.ClickMode = ad.ClickMode
	xmOad.VideoURL = mvutil.Base64Decode(ad.VideoURL)
	xmOad.VideoLength = ad.VideoLength
	xmOad.VideoSize = ad.VideoSize
	xmOad.VideoResolution = ad.VideoResolution
	if ad.AdTrackingPoint != nil {
		xmOad.AdTrackingPoint = ad.AdTrackingPoint
	}
	if len(ad.DeepLink) > 0 {
		xmOad.DeepLink = ad.DeepLink
	}
	return xmOad
}

func XMNewReturn(param *mvutil.Params) bool {
	if param.PublisherID != mvconst.PUB_XIAOMI || param.IsVast {
		return false
	}
	xmUnitsConf, _ := extractor.GetXM_NEW_RETURN_UNITS()
	// 若不配置，则默认都切新返回。若配置，则表示处于灰度阶段。
	if len(xmUnitsConf) <= 0 {
		return true
	}
	if len(xmUnitsConf) > 0 && mvutil.InInt64Arr(param.UnitID, xmUnitsConf) {
		return true
	}
	return false
}

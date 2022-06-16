package output

import (
	"hash/crc32"
	"strconv"

	"gitlab.mobvista.com/ADN/adnet/internal/extractor"
	"gitlab.mobvista.com/ADN/adnet/internal/mvconst"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
)

type OnlineAd struct {
	CampaignID      int64        `json:"id"`
	AppName         string       `json:"title"`
	AppDesc         string       `json:"desc"`
	PackageName     string       `json:"package_name"`
	IconURL         string       `json:"icon_url"`
	ImageURL        string       `json:"image_url"`
	ImpressionURL   string       `json:"impression_url"`
	VideoURL        string       `json:"video_url"`
	VideoLength     int          `json:"video_length"`
	VideoSize       int          `json:"video_size"`
	VideoResolution string       `json:"video_resolution"`
	CType           int          `json:"ctype"`
	AdvImp          []CAdvImp    `json:"adv_imp"`
	AdURLList       []string     `json:"ad_url_list"`
	ClickURL        string       `json:"click_url"`
	NoticeURL       string       `json:"notice_url"`
	WinNoticeURL    string       `json:"win_url"`
	Rating          float32      `json:"rating"`
	CtaText         string       `json:"ctatext"`
	CampaignType    int          `json:"link_type"`
	AdTrackingPoint *CAdTracking `json:"ad_tracking,omitempty"`
	NumberRating    int          `json:"number_rating,omitempty"`
	SubCategoryName []string     `json:"sub_category_name,omitempty"`
	PricePoint      *float32     `json:"payout,omitempty"`
	PreviewUrl      string       `json:"preview_link"`
	ClickMode       int          `json:"click_mode"`
	AdHtml          string       `json:"ad_html,omitempty"`

	// VideoCreativeid string `json:"video_creativeid"`
	VideoHeight    int32  `json:"video_height,omitempty"`
	VideoWidth     int32  `json:"Video_width,omitempty"`
	VideoMime      string `json:"video_mime,omitempty"`
	IconResolution string `json:"icon_resolution,omitempty"`
	IconMime       string `json:"icon_mime,omitempty"`
	// ImageCreativeid string `json:"image_creativeid"`
	ImageResolution string  `json:"image_resolution,omitempty"`
	ImageMime       string  `json:"image_mime,omitempty"`
	Bitrate         int32   `json:"bitrate,omitempty"`
	ParamK          string  `json:"param_k,omitempty"`
	BidPrice        float64 `json:"bid_price,omitempty"`
	DeepLink        string  `json:"deep_link,omitempty"`
	Vast            string  `json:"vast,omitempty"`
	CreativeId      int64   `json:"creative_id,omitempty"`
}

type OnlineResult struct {
	Status    int           `json:"status"`
	Msg       string        `json:"msg"`
	Data      OnlineData    `json:"data"`
	DebugInfo []interface{} `json:"debuginfo,omitempty"`
}

type OnlineData struct {
	Ads []OnlineAd `json:"ads"`
}

func RenderOnlineRes(mr MobvistaResult, r *mvutil.RequestParams) OnlineResult {
	var or OnlineResult
	or.Status = 1
	or.Msg = "success"

	var oadList []OnlineAd
	for _, v := range mr.Data.Ads {
		oad := renderOnlineAd(v, r)
		oadList = append(oadList, oad)
	}
	or.Data.Ads = oadList

	return or
}

func renderOnlineAd(ad Ad, r *mvutil.RequestParams) OnlineAd {
	var oad OnlineAd
	oad.CampaignID = ad.CampaignID
	oad.AppName = ad.AppName
	oad.AppDesc = ad.AppDesc
	oad.PackageName = ad.PackageName
	oad.IconURL = ad.IconURL
	oad.ImageURL = ad.ImageURL
	oad.ImpressionURL = ad.ImpressionURL
	oad.VideoURL = mvutil.Base64Decode(ad.VideoURL)
	oad.VideoLength = ad.VideoLength
	oad.VideoSize = ad.VideoSize
	oad.VideoResolution = ad.VideoResolution
	oad.CType = ad.CType
	oad.AdvImp = ad.AdvImp
	oad.AdURLList = ad.AdURLList
	oad.ClickURL = ad.ClickURL
	oad.NoticeURL = ad.NoticeURL
	oad.Rating = ad.Rating
	oad.CtaText = ad.CtaText
	oad.CampaignType = ad.CampaignType
	if ad.AdTrackingPoint != nil {
		oad.AdTrackingPoint = ad.AdTrackingPoint
	}
	oad.NumberRating = ad.NumberRating
	oad.SubCategoryName = ad.SubCategoryName
	oad.PreviewUrl = ad.PreviewUrl
	oad.ClickMode = ad.ClickMode
	oad.DeepLink = ad.DeepLink

	if Isvast(r) || r.Param.RequestPath == mvconst.PATHCHEETAH {
		// 兼容第三方无下发videoHeight和videoWidth的情况
		compatibleVideoHW(&ad)
		oad.VideoHeight = ad.VideoHeight
		oad.VideoWidth = ad.VideoWidth
		oad.VideoMime = ad.VideoMime
		oad.IconResolution = ad.IconResolution
		oad.IconMime = ad.IconMime
		oad.ImageResolution = ad.ImageResolution
		oad.ImageMime = ad.ImageMime
		oad.Bitrate = ad.Bitrate
	}

	// 输出价格 todo-fix
	if r.PublisherInfo.Publisher.Type == mvconst.PublisherTypeM {
		if r.UnitInfo.Unit.IsIncent == 1 {
			oad.PricePoint = &(ad.Price)
		}
	} else {
		if r.AppInfo.App.IsIncent == 1 {
			oad.PricePoint = &(ad.Price)
		}
	}
	// 返回k参数，用于开发者点击对数
	returnKUnits, _ := extractor.GetRETURN_PARAM_K_UNIT()
	if len(returnKUnits) > 0 && mvutil.InInt64Arr(r.Param.UnitID, returnKUnits) {
		oad.ParamK = ad.ParamK
	}
	// 获取出价，返回给开发者。vivo需要支持返回价格。
	// renderBidPrice(&oad, r)
	// 统一改为用portal的配置
	oad.BidPrice = ad.OnlineApiBidPrice
	if r.Param.RequestPath == mvconst.PATHBidAds {
		oad.BidPrice = r.BidPrice
		oad.WinNoticeURL = r.BidWinUrl
	}
	// 没有配置，则使用原有config表的配置兜底
	if oad.BidPrice == 0 {
		renderBidPrice(&oad, r)
	}

	if mvutil.IsBigoRequest(r.Param.PublisherID, r.Param.RequestType) {
		renderBigoCreativeId(r, &ad)
		oad.CreativeId = ad.CreativeId
	}

	// 对于走as dsp或as的情况，bigo 的vast需要使用preview兜底
	if IsVastReturnInJson(&r.Param) {
		renderBigoVastData(&oad, &ad, r)
	}

	oad.AdHtml = ad.AdHtml

	return oad
}

func renderBigoCreativeId(r *mvutil.RequestParams, ad *Ad) {
	// vast表示视频广告位，因此选视频素材返回。
	if r.Param.IsVast {
		// 若为adnet渲染素材，则使用内部的creativeid
		if r.Param.VideoCreativeid > 0 {
			ad.CreativeId = r.Param.VideoCreativeid
		} else {
			ad.CreativeId = int64(crc32.ChecksumIEEE([]byte(ad.VideoURL)))
		}
	} else {
		// 不走vast，一定是只使用大图素材。
		// 若为adnet渲染素材，则使用内部的creativeid
		if r.Param.ImageCreativeid > 0 {
			ad.CreativeId = r.Param.ImageCreativeid
		} else {
			ad.CreativeId = int64(crc32.ChecksumIEEE([]byte(ad.ImageURL)))
		}
	}
}

func renderBigoVastData(oad *OnlineAd, ad *Ad, r *mvutil.RequestParams) {
	// deeplink 单子才这样处理
	if len(ad.DeepLink) == 0 {
		return
	}
	if len(ad.PreviewUrl) == 0 {
		return
	}
	// as dsp 才会处理
	if r.DspExtData != nil && r.DspExtData.DspId != mvconst.FakeAdserverDsp {
		return
	}
	// 把click_url放到ad_tracking中
	if len(ad.ClickURL) > 0 {
		if oad.AdTrackingPoint != nil {
			oad.AdTrackingPoint.Click = append(oad.AdTrackingPoint.Click, ad.ClickURL)
		} else {
			var adtracking CAdTracking
			adtracking.Click = append(adtracking.Click, ad.ClickURL)
			oad.AdTrackingPoint = &adtracking
		}
	}
	// preview url 放到click_url中，作为兜底。
	oad.ClickURL = ad.PreviewUrl
}

func compatibleVideoHW(ad *Ad) {
	// 兼容第三方无下发videoHeight和videoWidth的情况
	if ad.VideoHeight == 0 && ad.VideoWidth == 0 {
		videoWidthStr, videoHeightStr := getWidthAndHeight(ad.VideoResolution)
		videoWidth, _ := strconv.Atoi(videoWidthStr)
		videoHeight, _ := strconv.Atoi(videoHeightStr)
		ad.VideoWidth = int32(videoWidth)
		ad.VideoHeight = int32(videoHeight)
	}
}

func renderBidPrice(oad *OnlineAd, r *mvutil.RequestParams) {
	onlineApiPubBidPriceConf := extractor.GetOnlineApiPubBidPriceConf()
	if onlineApiPubBidPriceConf == nil {
		return
	}

	if bidPrice, ok := onlineApiPubBidPriceConf.UnitConf[strconv.FormatInt(r.Param.UnitID, 10)]; ok && bidPrice > 0 {
		oad.BidPrice = bidPrice
		return
	}

	if bidPrice, ok := onlineApiPubBidPriceConf.AppConf[strconv.FormatInt(r.Param.AppID, 10)]; ok && bidPrice > 0 {
		oad.BidPrice = bidPrice
		return
	}

	if bidPrice, ok := onlineApiPubBidPriceConf.PubConf[strconv.FormatInt(r.Param.PublisherID, 10)]; ok && bidPrice > 0 {
		oad.BidPrice = bidPrice
	}
}

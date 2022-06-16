package output

import (
	"strconv"

	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
)

func RenderFunRes(mr MobvistaResult, r *mvutil.RequestParams) UCResult {
	var fun UCResult
	var bid Bid
	var seatbid Seatbid
	// 风行请求id
	fun.UCResponseId = r.Param.FunRequestId
	fun.Bidid = r.Param.RequestID
	for _, v := range mr.Data.Ads {
		bid = renderFunAd(v, *r)
		seatbid.Bid = append(seatbid.Bid, bid)
	}
	fun.Seatbid = append(fun.Seatbid, seatbid)
	return fun
}

func renderFunAd(ad Ad, r mvutil.RequestParams) Bid {
	var bid Bid
	// 先定位requestid
	bid.Id = r.Param.RequestID
	bid.Price = 0.0
	// 这里指的是曝光ID
	bid.UnitId = r.Param.FunImpId
	bid.Adm = mvutil.Base64Decode(ad.VideoURL)
	bid.FunCampaignId = strconv.FormatInt(ad.CampaignID, 10)
	bid.Ext = renderExt(ad)
	return bid
}

func renderExt(ad Ad) Ext {
	var ext Ext
	var pm Pm
	ext.PreviewUrl = ad.PreviewUrl
	// 在point=0时加入impression_url
	pm.Point = 0
	pm.TrackingUrl = ad.ImpressionURL
	ext.Pm = append(ext.Pm, pm)
	for _, v := range ad.AdTracking.Play_percentage {
		if v.Rate == 0 {
			// 视频播放开始url
			pm.Point = 0
			pm.TrackingUrl = v.Url
			ext.Pm = append(ext.Pm, pm)
		} else if v.Rate == 100 {
			// 视频播放结束url
			pm.Point = -1
			pm.TrackingUrl = v.Url
			ext.Pm = append(ext.Pm, pm)
		}
	}
	ext.Cm = append(ext.Cm, ad.ClickURL)
	return ext
}

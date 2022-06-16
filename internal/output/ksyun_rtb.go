package output

import (
	"strconv"

	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	"gitlab.mobvista.com/ADN/adnet/internal/protobuf"
)

func RenderKSYUNRes(mr MobvistaResult, r *mvutil.RequestParams) protobuf.STExchangeRsp {
	var res protobuf.STExchangeRsp
	res.RequestId = r.Param.KSYUNRequestID
	for _, v := range mr.Data.Ads {
		//////////todo
		ad, ok := renderKSYUNAD(v, r)
		if !ok {
			res.ErrorCode = 204
			return res
		}
		res.Ads = append(res.Ads, ad)
	}
	return res
}

func renderKSYUNAD(ad Ad, r *mvutil.RequestParams) (*protobuf.Ad, bool) {
	var adlist protobuf.Ad
	adlist.AdslotId = r.Param.KSYUNAdSlotID
	// adlist.MaxCpm = 2500                                // 分/CPM，写配置！
	adlist.MaxCpm = r.Param.KSYUNMaxCPM
	adlist.ExpirationTime = 3600                        // 自定3600秒，具体时间待确定。
	adlist.AdKey = strconv.FormatInt(ad.CampaignID, 10) //基数指代进制。
	adlist.Vid = r.Param.RequestID
	adlist.SearchKey = r.Param.KSYUNSKey
	// starting rendering monitor url staff...
	adlist.AdTracking = append(adlist.AdTracking,
		renderKSYUNTracking(ad, protobuf.TrackingEvent_AD_EXPOSURE),
		renderKSYUNTracking(ad, protobuf.TrackingEvent_AD_CLOSE),
		renderKSYUNTracking(ad, protobuf.TrackingEvent_VIDEO_AD_START),
		renderKSYUNTracking(ad, protobuf.TrackingEvent_VIDEO_AD_END))
	if len(ad.NoticeURL) > 0 {
		adlist.AdTracking = append(adlist.AdTracking,
			renderKSYUNTracking(ad, protobuf.TrackingEvent_AD_CLICK),
			renderKSYUNTracking(ad, protobuf.TrackingEvent_APP_AD_DOWNLOAD),
			renderKSYUNTracking(ad, protobuf.TrackingEvent_APP_AD_INSTALL),
			renderKSYUNTracking(ad, protobuf.TrackingEvent_APP_AD_START_DOWNLOAD),
		)
	}
	// then, give ad url staff...
	//////////////todo
	adtracking, ok := renderKSYUNMeta(ad)
	if !ok {
		return nil, false
	}
	adlist.MetaGroup = append(adlist.MetaGroup, adtracking)
	return &adlist, true
}

func renderKSYUNTracking(ad Ad, event protobuf.TrackingEvent) *protobuf.Tracking {
	var tracking protobuf.Tracking
	tracking.TrackingEvent = event
	switch event {
	case protobuf.TrackingEvent_AD_CLICK:
		if len(ad.NoticeURL) > 0 {
			tracking.TrackingUrl = append(tracking.TrackingUrl, ad.NoticeURL)
		}
		// 兼容头条点击tracking
		if len(ad.AdTracking.Click) > 0 {
			tracking.TrackingUrl = append(tracking.TrackingUrl, ad.AdTracking.Click...)
			// for _, v := range ad.AdTracking.Click {
			// 	Tracking.TrackingUrl = append(Tracking.TrackingUrl, v)
			// }
		}

	case protobuf.TrackingEvent_AD_EXPOSURE:
		tracking.TrackingUrl = append(tracking.TrackingUrl, ad.ImpressionURL)
		// 兼容头条展示tracking
		if len(ad.AdTracking.Impression) > 0 {
			tracking.TrackingUrl = append(tracking.TrackingUrl, ad.AdTracking.Impression...)
			// for _, v := range ad.AdTracking.Impression {
			// 	tracking.TrackingUrl = append(tracking.TrackingUrl, v)
			// }
		}
	case protobuf.TrackingEvent_AD_CLOSE:
		// 关闭上报
		tracking.TrackingUrl = ad.AdTracking.Close
	case protobuf.TrackingEvent_VIDEO_AD_START:
		for _, v := range ad.AdTracking.Play_percentage {
			if v.Rate == 0 {
				tracking.TrackingUrl = append(tracking.TrackingUrl, v.Url)
			}
		}
	case protobuf.TrackingEvent_VIDEO_AD_END:
		for _, v := range ad.AdTracking.Play_percentage {
			if v.Rate == 100 {
				tracking.TrackingUrl = append(tracking.TrackingUrl, v.Url)
			}
		}
	case protobuf.TrackingEvent_APP_AD_DOWNLOAD:
		// 中途关闭以及回看等无法提供
		tracking.TrackingUrl = ad.AdTracking.ApkDownloadEnd
	case protobuf.TrackingEvent_APP_AD_INSTALL:
		tracking.TrackingUrl = ad.AdTracking.ApkInstall
	case protobuf.TrackingEvent_APP_AD_START_DOWNLOAD:
		tracking.TrackingUrl = ad.AdTracking.ApkDownloadStart
	}
	return &tracking
}

func renderKSYUNMeta(ad Ad) (*protobuf.MaterialMeta, bool) {
	var Meta protobuf.MaterialMeta
	Meta.CreativeType = 9
	Meta.InteractionType = 2
	// 若不为apk单子，则使用浏览器打开网页
	if ad.CampaignType != 3 {
		Meta.InteractionType = 1
	}
	Meta.ClickUrl = ad.ClickURL // 下载链接
	Meta.Title = []byte(ad.OfferName)
	// convert desc to [][]byte
	desc := make([][]byte, 0)
	desc = append(desc, []byte(ad.AppDesc))
	Meta.Description = desc
	Meta.IconSrc = append(Meta.IconSrc, ad.IconURL)
	Meta.ImageSrc = append(Meta.ImageSrc, ad.ImageURL)
	Meta.AppPackage = ad.PackageName
	// 头条单子没有传此值
	if len(ad.AppSize) == 0 {
		ad.AppSize = "0"
	}
	size, err := strconv.ParseInt(ad.AppSize, 10, 64)
	if err != nil {
		return nil, false
	}
	Meta.AppSize = uint32(size)
	Meta.VideoUrl = mvutil.Base64Decode(ad.VideoURL)
	Meta.VideoDuration = uint32(ad.VideoLength)
	Meta.MaterialHeight = uint32(ad.VideoHeight)
	Meta.MaterialWidth = uint32(ad.VideoWidth)
	Meta.BrandName = ad.AppName
	Meta.AppName = Meta.BrandName
	return &Meta, true
}

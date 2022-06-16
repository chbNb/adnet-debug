package output

import (
	"net/url"
	"strconv"

	"gitlab.mobvista.com/ADN/adnet/internal/mvconst"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	"gitlab.mobvista.com/ADN/adnet/internal/protobuf"
)

type MaxRes struct {
}

func RenderMaxRes(mr MobvistaResult, r *mvutil.RequestParams) protobuf.BidResponse {
	var res protobuf.BidResponse
	maxPrice := uint32(r.Param.MaxPrice)
	iconVal := int32(128)
	var ads protobuf.BidResponse_Ads
	subTitle := "Mintegral"
	for i := range mr.Data.Ads {
		var creative protobuf.NativeAd_Creative
		var native protobuf.NativeAd
		var contentImg protobuf.NativeAd_Creative_Image
		var logo protobuf.NativeAd_Creative_Image
		var link protobuf.NativeAd_Creative_Link
		var video protobuf.NativeAd_Creative_Video
		v := mr.Data.Ads[i]
		creative.Title = &(v.AppName)
		creative.SubTitle = &subTitle
		creative.Description = &(v.AppDesc)
		creative.ButtonName = &(v.CtaText)
		contentImg.ImageUrl = &(v.ImageURL)
		contentImg.ImageWidth = &(r.Param.MaxImgW)
		contentImg.ImageHeight = &(r.Param.MaxImgH)
		creative.ContentImage = &contentImg

		logo.ImageUrl = &(v.IconURL)
		logo.ImageWidth = &iconVal
		logo.ImageHeight = &iconVal
		creative.Logo = &logo
		// 处理跳转链接
		clickUrl := "%%CLICK_URL_UNESC%%" + url.QueryEscape(v.ClickURL)
		link.ClickUrl = &(clickUrl)
		var landingType int32
		if v.CampaignType == 3 {
			landingType = 1
		} else {
			landingType = 0
		}
		link.LandingType = &landingType
		// 兼容头条点击tracking
		var clickTracks []string
		if len(v.AdTracking.Click) > 0 {
			for _, clickThroughVal := range v.AdTracking.Click {
				ttClickUrl := "%%CLICK_URL_UNESC%%" + url.QueryEscape(clickThroughVal)
				link.ClickTracks = append(clickTracks, ttClickUrl)
			}
		}
		creative.Link = &link
		videoUrl := mvutil.Base64Decode(v.VideoURL)
		video.VideoUrl = &videoUrl
		videoLength := int32(v.VideoLength)
		video.Duration = &videoLength
		// 添加视频监测链接
		for _, playperVal := range v.AdTracking.Play_percentage {
			if playperVal.Rate == 0 {
				video.VideoStartTracks = append([]string{}, playperVal.Url)
			}
			if playperVal.Rate == 100 {
				video.VideoComplateTracks = append([]string{}, playperVal.Url)
			}

		}
		creative.Video = &video
		var creatives []*protobuf.NativeAd_Creative
		native.Creatives = append(creatives, &creative)
		var impressionTracks []string
		native.ImpressionTracks = append(impressionTracks, v.ImpressionURL)
		// 兼容头条impression上报
		if len(v.AdTracking.Impression) > 0 {
			for _, impTrackVal := range v.AdTracking.Impression {
				native.ImpressionTracks = append(impressionTracks, impTrackVal)
			}
		}
		if v.Category == mvconst.GAME {
			native.Category = append([]int32{}, 80)
		} else {
			native.Category = append([]int32{}, 234)
		}
		native.DestinationUrl = append([]string{}, v.PreviewUrl)
		advId := strconv.Itoa(v.AdvID)
		native.AdvertiserId = &advId
		campaignId := strconv.FormatInt(v.CampaignID, 10)
		native.CreativeId = &campaignId
		native.DealId = &(r.Param.MaxDealId)
		native.TemplateId = &(r.Param.TemplateId)
		native.MaxCpmPrice = &maxPrice
		// AppAttr
		var appAttr protobuf.NativeAd_AppAttr
		appAttr.AppName = &(v.AppName)
		appAttr.AppPkg = &(v.PackageName)
		appAttr.AppMd5 = &(v.ApkMd5)
		apkVersionInt, _ := strconv.ParseInt(v.ApkVersion, 10, 32)
		avInt32 := int32(apkVersionInt)
		appAttr.AppVc = &(avInt32)
		appSizeInt, _ := strconv.ParseInt(v.AppSize, 10, 32)
		asInt32 := int32(appSizeInt)
		appAttr.AppSize = &(asInt32)
		native.AppAttr = &appAttr
		ads.NativeAd = append(ads.NativeAd, &native)
	}
	ads.AdslotId = &(r.Param.AdslotId)
	ads.MaxCpmPrice = &maxPrice
	res.Bid = &(r.Param.MaxBid)
	res.Ads = append(res.Ads, &ads)
	return res
}

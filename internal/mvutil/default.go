package mvutil

import (
	"gitlab.mobvista.com/ADN/adnet/internal/corsair_proto"
	"gitlab.mobvista.com/ADN/adnet/internal/mvconst"
	"gitlab.mobvista.com/ADN/mtg_openrtb/pkg/mtgrtb"
	"gitlab.mobvista.com/adserver/recommend_protocols/go/ad_server"
)

// SetCampaignDefaultFields - set campaign default values
func SetCampaignDefaultFields(adTmp *corsair_proto.Campaign) {
	appName := ""
	adTmp.AppName = &appName
	desc := ""
	adTmp.AppDesc = &desc
	packageName := ""
	adTmp.PackageName = &packageName
	iconURL := ""
	adTmp.IconURL = &iconURL
	imgURL := ""
	adTmp.ImageURL = &imgURL
	htmlTemp := ""
	adTmp.HtmlTemplate = &htmlTemp
	imgSize := "1200x627"
	adTmp.ImageSize = &imgSize
	clickURL := ""
	adTmp.ClickURL = &clickURL
	noticeURL := ""
	adTmp.NoticeURL = &noticeURL

	videoURL := ""
	adTmp.VideoURL = &videoURL
	videoLength := int32(0)
	adTmp.VideoLength = &videoLength
	videlSize := int32(0)
	adTmp.VideoSize = &videlSize
	bitRate := int32(0)
	adTmp.BitRate = &bitRate
	videoResolution := ""
	adTmp.VideoResolution = &videoResolution
	videoEndType := int32(2)
	adTmp.VideoEndType = &videoEndType
	playableAdsWithoutVideo := int32(1)
	adTmp.PlayableAdsWithoutVideo = &playableAdsWithoutVideo
	endcardURL := ""
	adTmp.EndcardURL = &endcardURL
	watchMile := int32(0)
	adTmp.WatchMile = &watchMile
	ctype := int32(0)
	adTmp.CType = &ctype
	timp := int32(0)
	adTmp.TImp = &timp
	advid := int32(0)
	adTmp.AdvertiserId = &advid
	price := 0.0
	adTmp.Price = &price
	offerName := ""
	adTmp.OfferName = &offerName
	installToken := ""
	adTmp.InstallToken = &installToken
	fca := int32(0)
	adTmp.FCA = &fca
	fcb := int32(0)
	adTmp.FCB = &fcb

	//to-do change to 大图的template
	//var adTemplate ad_server.ADTemplate
	adTemplate := ad_server.ADTemplate(mvconst.TemplateSinglePic)
	adTmp.AdTemplate = &adTemplate

	//to-do change ad_source_id
	adTmp.AdSource = mvconst.ADSourceAPIOffer

	// 确认Mtg的offer_type
	offerType := int32(mvconst.OffTypeMtg)
	adTmp.OfferType = &offerType

	appSize := ""
	adTmp.AppSize = &appSize
	clickMode := int32(0)
	adTmp.ClickMode = &clickMode
	rating := GetThirdPartRating()
	adTmp.Rating = &rating
	landingType := int32(0)
	adTmp.LandingType = &landingType
	ctaText := ""
	adTmp.CtaText = &ctaText
	clickCacheTime := int32(0)
	adTmp.ClickCacheTime = &clickCacheTime
	linkType := int32(mtgrtb.LinkType_DEFAULTBROWSER)
	adTmp.LinkType = &linkType
	guideLines := ""
	adTmp.GuideLines = &guideLines
	rewardAmount := int32(0)
	adTmp.RewardAmount = &rewardAmount
	rewardName := ""
	adTmp.RewardName = &rewardName
	retargetOffer := int32(2)
	adTmp.RetargetOffer = &retargetOffer
	statsURL := ""
	adTmp.StatsURL = &statsURL

	imgResolution := ""
	adTmp.ImageResolution = &imgResolution
	imageMime := ""
	adTmp.ImageMime = &imageMime
	imageHeight := 0
	imageWidth := 0
	adTmp.ImageHeight = &imageHeight
	adTmp.ImageWidth = &imageWidth
	//iconResolution := ""
	//adTmp.IconResolution = &iconResolution
	//iconMime := ""
	//adTmp.IconMime = &iconMime
}

func GetDefaultVideoResolution(orientation int) string {
	result := "1080x720"
	if orientation == mvconst.ORIENTATION_PORTRAIT {
		result = "720x1080"
	}
	return result
}

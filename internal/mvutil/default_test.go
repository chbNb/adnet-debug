package mvutil

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"gitlab.mobvista.com/ADN/adnet/internal/corsair_proto"
)

func TestSetCampaignDefaultFields(t *testing.T) {
	Convey("SetCampaignDefaultFields", t, func() {
		appName := "test_app_name"
		appDesc := "test_app_desc"
		packageName := "test_package_name"
		iconURL := "http://test.com"
		imageSize := "300x300"
		price := 42.0
		adTmp := corsair_proto.Campaign{
			AppName:     &appName,
			AppDesc:     &appDesc,
			PackageName: &packageName,
			IconURL:     &iconURL,
			ImageSize:   &imageSize,
			Price:       &price,
		}
		SetCampaignDefaultFields(&adTmp)

		So(*adTmp.AppName, ShouldEqual, "")
		So(*adTmp.AppDesc, ShouldEqual, "")
		So(*adTmp.PackageName, ShouldEqual, "")
		So(*adTmp.IconURL, ShouldEqual, "")
		So(*adTmp.ImageSize, ShouldEqual, "1200x627")
		So(*adTmp.Price, ShouldEqual, 0.0)
		So(*adTmp.VideoLength, ShouldEqual, 0)
		So(*adTmp.VideoSize, ShouldEqual, 0)
		So(*adTmp.PlayableAdsWithoutVideo, ShouldEqual, 1)
		So(*adTmp.VideoEndType, ShouldEqual, 2)
		So(*adTmp.WatchMile, ShouldEqual, 0)
		So(*adTmp.CType, ShouldEqual, 0)
		So(*adTmp.TImp, ShouldEqual, 0)
		So(*adTmp.AdvertiserId, ShouldEqual, 0)
		So(*adTmp.OfferName, ShouldEqual, "")
		So(*adTmp.InstallToken, ShouldEqual, "")
		So(*adTmp.FCA, ShouldEqual, 0)
		So(*adTmp.FCB, ShouldEqual, 0)
		So(*adTmp.AdTemplate, ShouldEqual, 1)
		So(adTmp.AdSource, ShouldEqual, 1)
		So(*adTmp.AppSize, ShouldEqual, "")
		So(*adTmp.OfferType, ShouldEqual, 0)
		So(*adTmp.ClickMode, ShouldEqual, 0)
		//So(*adTmp.Rating, ShouldEqual, 0)
		So(*adTmp.LandingType, ShouldEqual, 0)
		So(*adTmp.CtaText, ShouldEqual, "")
		So(*adTmp.ClickCacheTime, ShouldEqual, 0)
		So(*adTmp.LinkType, ShouldEqual, 9)
		So(*adTmp.GuideLines, ShouldEqual, "")
		So(*adTmp.RewardAmount, ShouldEqual, 0)
		So(*adTmp.RewardName, ShouldEqual, "")
		So(*adTmp.RetargetOffer, ShouldEqual, 2)
		So(*adTmp.StatsURL, ShouldEqual, "")
	})
}

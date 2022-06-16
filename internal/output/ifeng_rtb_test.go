package output

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
)

func TestRenderIFENGRes(t *testing.T) {
	convey.Convey("RenderIFENGRes ok", t, func() {
		var mr MobvistaResult
		mr.Status = 1
		mr.Msg = "success"
		mr.Data.Ads = make([]Ad, 0)
		var ad Ad
		ad.CampaignID = int64(123)
		ad.AppName = "appName"
		ad.AppDesc = "appDesc"
		ad.PackageName = "packageName"
		ad.IconURL = "iconUrl"
		ad.ImageURL = "imageUrl"
		ad.ImpressionURL = "impressionUrl"
		ad.ClickURL = "clickUrl"
		ad.Rating = 3.14
		ad.CtaText = "ctaText"
		ad.CampaignType = 3
		ad.PreviewUrl = "previewUrl"
		ad.ClickMode = 5
		ad.VideoURL = "J75AJcKAJdzuYrh="
		ad.VideoResolution = "300x400"
		ad.AdTracking.Play_percentage = []CPlayTracking{CPlayTracking{Rate: 0, Url: "v0-url"}, CPlayTracking{Rate: 100, Url: "v100-url"}}
		mr.Data.Ads = append(mr.Data.Ads, ad)
		r := &mvutil.RequestParams{
			Param: mvutil.Params{
				IFENGId:    "IFENGId",
				RequestID:  "requestid",
				IFENGTagId: "IFENGTagId",
				IFENGImpId: "IFENGImpId",
				IFENGAdId:  int64(12345),
			},
		}
		vastxml := "vastxml"
		ret := RenderIFENGRes(mr, r, []byte(vastxml))
		convey.So(ret, convey.ShouldNotBeNil)
		convey.So(len(ret.Ad), convey.ShouldEqual, 1)
	})
}

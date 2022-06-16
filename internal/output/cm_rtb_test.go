package output

import (
	"testing"

	"github.com/bouk/monkey"
	"github.com/smartystreets/goconvey/convey"
	"gitlab.mobvista.com/ADN/adnet/internal/extractor"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
)

func TestRenderCMRes(t *testing.T) {
	convey.Convey("RenderCMRes ok", t, func() {
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
		ad.VideoResolution = "300x400"
		mr.Data.Ads = append(mr.Data.Ads, ad)
		r := &mvutil.RequestParams{
			Param: mvutil.Params{
				CMId:      "cmid",
				RequestID: "requestid",
				CMImpId:   "cmpImpId",
				AppID:     234,
			},
		}
		guard := monkey.Patch(extractor.GetONLINE_PRICE_FLOOR_APPID, func() (map[string]float64, bool) {
			return map[string]float64{"234": 3.14}, true
		})
		defer guard.Unpatch()
		vastxml := "vastxml"
		ret, _ := RenderCMRes(mr, r, []byte(vastxml))
		convey.So(ret, convey.ShouldNotBeNil)
		convey.So(len(ret.SeatBid), convey.ShouldEqual, 1)
	})
}

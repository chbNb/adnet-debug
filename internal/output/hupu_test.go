package output

import (
	"testing"

	"github.com/bouk/monkey"
	"github.com/smartystreets/goconvey/convey"
	"gitlab.mobvista.com/ADN/adnet/internal/extractor"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
)

func TestRenderHupuRes(t *testing.T) {
	convey.Convey("RenderHupuRes ok", t, func() {
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
				HupuRequestId: "hupuRequestId",
				RequestID:     "requestid",
				HupuImpId:     "12345",
				AppID:         234,
				DealId:        "DealId",
			},
		}
		guard := monkey.Patch(extractor.GetHUPU_DEFAULT_PRICE, func() (int, bool) {
			return 234, true
		})
		defer guard.Unpatch()
		guard = monkey.Patch(extractor.GetADNET_CONF_LIST, func() map[string][]int64 {
			return map[string][]int64{}
		})
		defer guard.Unpatch()
		creative := []int64{int64(123)}
		ret, _ := RenderHupuRes(mr, r, creative)
		convey.So(ret, convey.ShouldNotBeNil)
		convey.So(len(ret.Seatbid), convey.ShouldEqual, 1)
	})
}

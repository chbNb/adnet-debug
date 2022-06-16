package output

import (
	"testing"

	"github.com/bouk/monkey"
	"github.com/smartystreets/goconvey/convey"
	"gitlab.mobvista.com/ADN/adnet/internal/extractor"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
)

func TestRenderMIGURes(t *testing.T) {
	convey.Convey("RenderMIGURes ok", t, func() {
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
		ad.OfferName = "offerName"
		ad.CampaignType = 3
		ad.PreviewUrl = "previewUrl"
		ad.ClickMode = 5
		ad.NoticeURL = "NoticeURL"
		ad.VideoURL = "J75AJcKAJdzuYrh="
		ad.VideoLength = 15
		ad.VideoWidth = int32(300)
		ad.VideoHeight = int32(400)
		ad.VideoResolution = "300x400"
		ad.ImageResolution = "1200x768"
		ad.AdTracking.Click = []string{"click1", "click2"}
		ad.AdTracking.Impression = []string{"impression1", "impression2"}
		ad.AdTracking.Close = []string{"close", "close1"}
		ad.AdTracking.ApkDownloadStart = []string{"apkDownloadStart"}
		ad.AdTracking.ApkInstall = []string{"ApkInstall"}
		ad.AdTracking.ApkDownloadStart = []string{"ApkDownloadStart"}
		ad.AdTracking.Play_percentage = []CPlayTracking{CPlayTracking{Rate: 0, Url: "v0-url"}, CPlayTracking{Rate: 100, Url: "v100-url"}}
		mr.Data.Ads = append(mr.Data.Ads, ad)
		r := &mvutil.RequestParams{
			Param: mvutil.Params{
				MIGUId:    "MIGUId",
				RequestID: "RequestId",
				MIGUImpId: "MIGUImpId",
				AppID:     456,
			},
		}
		guard := monkey.Patch(extractor.GetONLINE_PRICE_FLOOR_APPID, func() (map[string]float64, bool) {
			return map[string]float64{"456": 3.14}, true
		})
		defer guard.Unpatch()
		creative := map[string]string{"123": "creativeid"}
		ret, _ := RenderMIGURes(mr, r, creative)
		convey.So(ret, convey.ShouldNotBeNil)
		convey.So(len(ret.Seatbid), convey.ShouldEqual, 1)
	})
}

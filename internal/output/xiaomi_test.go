package output

import (
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	"testing"

	"github.com/bouk/monkey"
	"github.com/smartystreets/goconvey/convey"
	"gitlab.mobvista.com/ADN/adnet/internal/extractor"
	"gitlab.mobvista.com/ADN/adnet/internal/mvconst"
)

func TestXMNewReturn(t *testing.T) {
	convey.Convey("XMNewReturn false", t, func() {
		var params mvutil.Params
		params.PublisherID = int64(123)
		params.UnitID = int64(0)
		params.IsVast = false
		ret := XMNewReturn(&params)
		convey.So(ret, convey.ShouldBeFalse)
	})

	convey.Convey("XMNewReturn true", t, func() {
		guard := monkey.Patch(extractor.GetXM_NEW_RETURN_UNITS, func() ([]int64, bool) {
			return []int64{678}, true
		})
		defer guard.Unpatch()
		var params mvutil.Params
		params.PublisherID = int64(mvconst.PUB_XIAOMI)
		params.UnitID = int64(678)
		params.IsVast = false
		ret := XMNewReturn(&params)
		convey.So(ret, convey.ShouldBeTrue)
	})
}

func TestRenderxmOnlineRes(t *testing.T) {
	convey.Convey("RenderxmOnlineRes ok", t, func() {
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
		mr.Data.Ads = append(mr.Data.Ads, ad)
		ret := RenderxmOnlineRes(mr)
		convey.So(ret, convey.ShouldNotBeNil)
		convey.So(len(ret.Data.XMAds), convey.ShouldEqual, 1)
	})
}

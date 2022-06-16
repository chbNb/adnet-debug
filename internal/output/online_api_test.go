package output

import (
	smodel "gitlab.mobvista.com/ADN/structs/model"
	"testing"

	. "github.com/bouk/monkey"
	. "github.com/smartystreets/goconvey/convey"
	"gitlab.mobvista.com/ADN/adnet/internal/extractor"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
)

func TestRenderOnlineRes(t *testing.T) {
	Convey("test RenderOnlineRes", t, func() {
		var mr MobvistaResult
		var r mvutil.RequestParams
		var res OnlineResult

		guard := Patch(renderOnlineAd, func(ad Ad, r *mvutil.RequestParams) OnlineAd {
			return OnlineAd{
				CampaignID:  int64(233),
				AppName:     "233_app",
				AppDesc:     "233_desc",
				PackageName: "233_pack",
			}
		})
		defer guard.Unpatch()

		guard = Patch(extractor.GetRETURN_PARAM_K_UNIT, func() ([]int64, bool) {
			return []int64{}, true
		})
		defer guard.Unpatch()

		guard = Patch(extractor.GetOnlineApiPubBidPriceConf, func() *mvutil.OnlineApiPubBidPriceConf {
			return &mvutil.OnlineApiPubBidPriceConf{}
		})
		defer guard.Unpatch()

		Convey("data", func() {
			r = mvutil.RequestParams{
				UnitInfo:      &smodel.UnitInfo{},
				AppInfo:       &smodel.AppInfo{},
				PublisherInfo: &smodel.PublisherInfo{},
				DebugInfo:     "debugInfo",
			}
			mr.Data = MobvistaData{
				Ads: []Ad{
					Ad{
						CampaignID:  int64(100001),
						OfferID:     int(100002),
						AppName:     "test_app_name",
						AppDesc:     "test_app_desc",
						PackageName: "test_package",
						IconURL:     "test_icon",
						ImageURL:    "test_image_url",
						ImageSize:   "test_image_size",
					},
					Ad{
						CampaignID:  int64(200001),
						OfferID:     int(200002),
						AppName:     "test_app_name_2",
						AppDesc:     "test_app_desc_2",
						PackageName: "test_package_2",
						IconURL:     "test_icon_2",
						ImageURL:    "test_image_url_2",
						ImageSize:   "test_image_size_2",
					},
				},
			}

			res = RenderOnlineRes(mr, &r)
			So(res.Status, ShouldEqual, 1)
			So(res.Msg, ShouldEqual, "success")
			So(res.Data, ShouldResemble, OnlineData{
				Ads: []OnlineAd{
					OnlineAd{
						CampaignID:  int64(233),
						AppName:     "233_app",
						AppDesc:     "233_desc",
						PackageName: "233_pack",
					},
					OnlineAd{
						CampaignID:  int64(233),
						AppName:     "233_app",
						AppDesc:     "233_desc",
						PackageName: "233_pack",
					},
				},
			})
		})
	})
}

func TestRenderOnlineAd(t *testing.T) {
	Convey("test renderOnlineAd", t, func() {
		var ad Ad
		var r mvutil.RequestParams
		var res OnlineAd
		guardConf := Patch(extractor.GetRETURN_PARAM_K_UNIT, func() ([]int64, bool) {
			return []int64{12345}, true
		})
		defer guardConf.Unpatch()

		guard := Patch(extractor.GetOnlineApiPubBidPriceConf, func() *mvutil.OnlineApiPubBidPriceConf {
			return &mvutil.OnlineApiPubBidPriceConf{}
		})
		defer guard.Unpatch()
		Convey("data", func() {
			r = mvutil.RequestParams{
				UnitInfo:      &smodel.UnitInfo{},
				AppInfo:       &smodel.AppInfo{},
				PublisherInfo: &smodel.PublisherInfo{},
				DebugInfo:     "debugInfo",
			}

			ad = Ad{
				CampaignID:  int64(100003),
				OfferID:     int(100003),
				AppName:     "test_app_name_3",
				AppDesc:     "test_app_desc_3",
				PackageName: "test_package_3",
				IconURL:     "test_icon_3",
				ImageURL:    "test_image_url_3",
				ImageSize:   "test_image_size_3",
			}

			res = renderOnlineAd(ad, &r)
			So(res.CampaignID, ShouldEqual, int64(100003))
			So(res.AppName, ShouldEqual, "test_app_name_3")
			So(res.AppDesc, ShouldEqual, "test_app_desc_3")
			So(res.PackageName, ShouldEqual, "test_package_3")
			So(res.IconURL, ShouldEqual, "test_icon_3")
			So(res.ImageURL, ShouldEqual, "test_image_url_3")
		})
	})
}

package output

import (
	"testing"

	"github.com/bouk/monkey"
	"github.com/smartystreets/goconvey/convey"
	"gitlab.mobvista.com/ADN/adnet/internal/extractor"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
)

func TestRanderVastData(t *testing.T) {

	convey.Convey("test RanderVastData", t, func() {
		params := mvutil.RequestParams{}
		guard := monkey.Patch(extractor.GetHAS_EXTENSIONS_UNIT, func() ([]int64, bool) {
			return []int64{23, 45, 678}, true
		})
		defer guard.Unpatch()

		guard = monkey.Patch(extractor.GetGAMELOFT_CREATIVE_URLS, func() (map[string]map[string]string, bool) {
			return map[string]map[string]string{
				"1": {
					"systemIcon": "http://cdn-adn.rayjump.com/cdn-adn/default-icon/mintergral_icon.png",
					"fullStar":   "http://cdn-adn.rayjump.com/cdn-adn/default-icon/fullStar.png",
					"halfStar":   "http://cdn-adn.rayjump.com/cdn-adn/default-icon/halfStar.png",
					"emptyStar":  "http://cdn-adn.rayjump.com/cdn-adn/default-icon/emptyStar.png",
					"ctaButton":  "http://cdn-adn.rayjump.com/cdn-adn/default-icon/button.png",
				},
				"2": {
					"systemIcon": "https://cdn-adn-https.rayjump.com/cdn-adn/default-icon/mintergral_icon.png",
					"fullStar":   "https://cdn-adn-https.rayjump.com/cdn-adn/default-icon/fullStar.png",
					"halfStar":   "https://cdn-adn-https.rayjump.com/cdn-adn/default-icon/halfStar.png",
					"emptyStar":  "https://cdn-adn-https.rayjump.com/cdn-adn/default-icon/emptyStar.png",
					"ctaButton":  "https://cdn-adn-https.rayjump.com/cdn-adn/default-icon/button.png",
				},
			}, true
		})
		defer guard.Unpatch()

		guard = monkey.Patch(extractor.GetUNIT_WITHOUT_VIDEO_START, func() ([]int64, bool) {
			return []int64{}, true
		})
		defer guard.Unpatch()

		guard = monkey.Patch(extractor.GetNEW_AFREECATV_UNIT, func() []int64 {
			return []int64{}
		})
		defer guard.Unpatch()

		convey.Convey("xml encoding success", func() {
			resOnline := OnlineResult{
				Data: OnlineData{
					Ads: []OnlineAd{
						OnlineAd{
							CampaignID:  int64(233),
							AppName:     "233_app",
							AppDesc:     "233_desc",
							PackageName: "233_pack",
							AdTrackingPoint: &CAdTracking{
								Mute: []string{"233", "253"},
								Play_percentage: []CPlayTracking{
									CPlayTracking{0, "daskgja"},
									CPlayTracking{25, "sdagaag"},
									CPlayTracking{75, "sdagaag"},
								},
							},
						},
						OnlineAd{
							CampaignID:  int64(233),
							AppName:     "233_app",
							AppDesc:     "233_desc",
							PackageName: "233_pack",
							AdTrackingPoint: &CAdTracking{
								Mute: []string{"233", "283"},
								Play_percentage: []CPlayTracking{
									CPlayTracking{0, "daskgja"},
									CPlayTracking{25, "sdagaag"},
									CPlayTracking{75, "sdagaag"},
								},
							},
						},
					},
				},
			}
			res, err := RanderVastData(&params, resOnline)
			convey.So(err, convey.ShouldBeNil)
			convey.So(res, convey.ShouldNotBeNil)
		})

	})

}

func TestIsVastReturnInJson(t *testing.T) {
	convey.Convey("test IsVastReturnInJson return false", t, func() {
		params := mvutil.Params{}
		params.IsVast = false
		params.PublisherID = int64(123)
		res := IsVastReturnInJson(&params)
		convey.So(res, convey.ShouldBeFalse)
	})
	convey.Convey("test IsVastReturnInJson return true", t, func() {
		params := mvutil.Params{}
		params.IsVast = true
		params.PublisherID = int64(123)
		guard := monkey.Patch(extractor.GetADNET_CONF_LIST, func() map[string][]int64 {
			return map[string][]int64{
				"vastReturnInJsonPub": []int64{
					int64(123),
				},
			}
		})
		defer guard.Unpatch()
		res := IsVastReturnInJson(&params)
		convey.So(res, convey.ShouldBeTrue)
	})
}

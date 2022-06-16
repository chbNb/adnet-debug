package process_pipeline

import (
	"testing"

	"gitlab.mobvista.com/ADN/adnet/internal/corsair_proto"
	//"gitlab.mobvista.com/ADN/adnet/internal/extractor"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"

	//. "github.com/bouk/monkey"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGenAdBackendData(t *testing.T) {
	Convey("Test genAdBackendData", t, func() {
		res := genAdBackendData(42, "test_c_id", 1)
		So(res, ShouldEqual, "42:test_c_id:1;")
	})
}

func TestGetBackendConfig(t *testing.T) {
	Convey("Test GetBackendConfig", t, func() {
		var reqCtx mvutil.ReqCtx
		var res string

		reqCtx = mvutil.ReqCtx{
			OrderedBackends: []int{1, 2, 3},
			Backends: map[int]*mvutil.BackendCtx{
				1: &mvutil.BackendCtx{},
				2: &mvutil.BackendCtx{
					AdReqKeyName: "t_ad_req_name_2",
				},
			},
		}

		res = GetBackendConfig(&reqCtx)
		So(res, ShouldEqual, "1:0:0;2:t_ad_req_name_2:0")
	})
}

func TestFilterAds(t *testing.T) {
	Convey("Test FilterAds", t, func() {
		var reqCtx mvutil.ReqCtx
		var id int
		var ads corsair_proto.BackendAds
		var adCount int64
		var backendData []string

		Convey("When ads is empty", func() {
			FilterAds(&reqCtx, id, &ads, &adCount, &backendData)
			So(ads.CampaignList, ShouldBeNil)
		})

		Convey("When ads is not empty, and adtype is native(42)", func() {
			id = 2
			//tNum := int64(10)
			vVersion := "1"
			adType := int32(42)
			reqCtx = mvutil.ReqCtx{
				OrderedBackends: []int{1, 2, 3},
				Backends: map[int]*mvutil.BackendCtx{
					1: &mvutil.BackendCtx{},
					2: &mvutil.BackendCtx{
						AdReqKeyName: "t_ad_req_name_2",
					},
				},
				ReqParams: &mvutil.RequestParams{
					Param: mvutil.Params{
						TNum:         10,
						VideoVersion: vVersion,
						AdType:       adType,
					},
				},
			}

			oType1 := int32(0)
			oType2 := int32(1)
			vURL1 := "v_url_1"
			vURL2 := "v_url_2"
			ads = corsair_proto.BackendAds{
				BackendId: int32(101),
				CampaignList: []*corsair_proto.Campaign{
					&corsair_proto.Campaign{
						CampaignId: "10001",
						OfferType:  &oType1,
						VideoURL:   &vURL1,
					},
					&corsair_proto.Campaign{
						CampaignId: "10002",
						OfferType:  &oType2,
						VideoURL:   &vURL2,
					},
				},
			}

			FilterAds(&reqCtx, id, &ads, &adCount, &backendData)
			So(ads.CampaignList[0].CampaignId, ShouldEqual, "10001")
			So(*(ads.CampaignList[0].OfferType), ShouldEqual, 0)
		})

		Convey("When ads is not empty, and adtype is rv(94)", func() {
			id = 2
			//tNum := int64(10)
			vVersion := ""
			adType := int32(94)
			reqCtx = mvutil.ReqCtx{
				OrderedBackends: []int{1, 2, 3},
				Backends: map[int]*mvutil.BackendCtx{
					1: &mvutil.BackendCtx{},
					2: &mvutil.BackendCtx{
						AdReqKeyName: "t_ad_req_name_2",
					},
				},
				ReqParams: &mvutil.RequestParams{
					Param: mvutil.Params{
						TNum:         10,
						VideoVersion: vVersion,
						AdType:       adType,
					},
				},
			}

			oType1 := int32(99)
			oType2 := int32(1)
			vURL1 := "v_url_1"
			vURL2 := "v_url_2"
			ads = corsair_proto.BackendAds{
				BackendId: int32(101),
				CampaignList: []*corsair_proto.Campaign{
					&corsair_proto.Campaign{
						CampaignId: "20001",
						OfferType:  &oType1,
						VideoURL:   &vURL1,
					},
					&corsair_proto.Campaign{
						CampaignId: "20002",
						OfferType:  &oType2,
						VideoURL:   &vURL2,
					},
				},
			}

			// guard := Patch(extractor.GetAdsourcePkgBlackList, func(platform int, dspId int64, backendId int64) []string {
			// 	return []string{"test_pkg"}
			// })
			// defer guard.Unpatch()

			FilterAds(&reqCtx, id, &ads, &adCount, &backendData)
			So(ads.CampaignList[0].CampaignId, ShouldEqual, "20001")
			So(*(ads.CampaignList[0].OfferType), ShouldEqual, 99)

			So(ads.CampaignList[1].CampaignId, ShouldEqual, "20002")
			So(*(ads.CampaignList[1].OfferType), ShouldEqual, 1)
		})
	})
}

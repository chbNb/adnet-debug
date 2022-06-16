package output

import (
	"fmt"
	"math/rand"
	"testing"

	supply_mvutil "gitlab.mobvista.com/ADN/chasm/module/supply/mvutil"

	. "bou.ke/monkey"
	. "github.com/smartystreets/goconvey/convey"
	"gitlab.mobvista.com/ADN/adnet/internal/corsair_proto"
	"gitlab.mobvista.com/ADN/adnet/internal/extractor"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	smodel "gitlab.mobvista.com/ADN/structs/model"
)

func TestRenderThirdPartyCampaign(t *testing.T) {
	Convey("test RenderThirdPartyCampaign", t, func() {
		r := mvutil.RequestParams{AppInfo: &smodel.AppInfo{}, UnitInfo: &smodel.UnitInfo{}, PublisherInfo: &smodel.PublisherInfo{}}
		var corsairCampaign corsair_proto.Campaign
		var res Ad

		guard := Patch(RenderThirdPartyUrls, func(params *mvutil.Params, ad *Ad, r *mvutil.RequestParams) {
		})
		defer guard.Unpatch()

		guard = Patch(extractor.GetDefRVTemplate, func() (*smodel.VideoTemplateUrlItem, bool) {
			return &smodel.VideoTemplateUrlItem{}, true
		})
		defer guard.Unpatch()

		guard = Patch(GetUrlScheme, func(params *mvutil.Params) string {
			return "url_scheme"
		})
		defer guard.Unpatch()

		guard = Patch(RenderThirdAdtracking, func(ad *Ad, corsairCampaign corsair_proto.Campaign, params *mvutil.Params,
			backendID int32, r *mvutil.RequestParams) {
		})
		defer guard.Unpatch()

		guard = Patch(HandleNumberRating, func(numberRating int) int {
			return 100
		})
		defer guard.Unpatch()

		guard = Patch(extractor.GetC_TOI, func() (int32, bool) {
			return int32(42), true
		})
		defer guard.Unpatch()

		guard = Patch(extractor.GetTP_WITHOUT_RV, func() ([]int64, bool) {
			return []int64{}, true
		})
		defer guard.Unpatch()

		guard = Patch(extractor.GetTEMPLATE_MAP, func() (mvutil.GlobalTemplateMap, bool) {
			return mvutil.GlobalTemplateMap{
				EndScreen: map[string]string{
					"401": "hybird.rayjump.com/offerwall/tpl/mintegral/endscreen.v4.html",
				},
			}, true
		})

		guard = Patch(renderPlct, func(ad *Ad, r *mvutil.Params, backendId int, dspId int64) {
		})
		defer guard.Unpatch()

		guard = Patch(extractor.GetADNET_SWITCHS, func() (map[string]int, bool) {
			return map[string]int{}, true
		})
		defer guard.Unpatch()

		Convey("test", func() {
			res = RenderThirdPartyCampaign(&r, corsairCampaign, int32(1), 0)
			So(res.NumberRating, ShouldEqual, 100)
		})
	})
}

func TestRenderThirdAdtracking(t *testing.T) {
	Convey("test RenderThirdAdtracking", t, func() {
		var ad Ad
		var corsairCampaign corsair_proto.Campaign
		var params mvutil.Params
		var r mvutil.RequestParams
		guard := Patch(hitChetLinkConfig, func(params *mvutil.Params) (string, bool) {
			return "", false
		})
		defer guard.Unpatch()

		Convey("corsairCampaign.AdTracking == nil", func() {
			corsairCampaign.AdTracking = nil
			RenderThirdAdtracking(&ad, corsairCampaign, &params, int32(1), &r)
		})

		Convey("corsairCampaign.AdTracking != nil", func() {
			adTra := corsair_proto.AdTracking{
				Start:      []string{"start_1", "start_2"},
				FirstQuart: []string{"first_1", "first_2"},
				Mid:        []string{"mid_1", "mid_2"},
				ThirdQuart: []string{"third_q_1", "third_q_2"},
				PlayPerct: []*corsair_proto.PlayPercent{
					&corsair_proto.PlayPercent{
						Rate: int32(40),
						URL:  "url_1",
					},
					&corsair_proto.PlayPercent{
						Rate: int32(60),
						URL:  "url_2",
					},
				},
			}
			corsairCampaign.AdTracking = &adTra
			RenderThirdAdtracking(&ad, corsairCampaign, &params, int32(0), &r)
			So(ad.AdTracking.Start, ShouldResemble, []string{"start_1", "start_2"})
			So(ad.AdTracking.Midpoint, ShouldResemble, []string{"mid_1", "mid_2"})
			So(ad.AdTracking.Play_percentage, ShouldResemble, []CPlayTracking{
				CPlayTracking{
					Rate: int(40),
					Url:  "url_1",
				},
				CPlayTracking{
					Rate: int(60),
					Url:  "url_2",
				},
			})
		})
	})
}

func TestRenderThirdPartyUrls(t *testing.T) {
	Convey("test RenderThirdPartyUrls", t, func() {
		var params mvutil.Params
		var ad Ad

		Convey("test", func() {
			params = mvutil.Params{
				Platform:  1,
				GAID:      "test_gaid",
				BackendID: 6,
				RequestID: "test_r_id",
				Domain:    "test_domain",
			}
			r := mvutil.RequestParams{
				DspExt: "{dspid: 1}",
			}

			RenderThirdPartyUrls(&params, &ad, &r)
			So(ad.ImpressionURL, ShouldEqual, "http://test_domain/impression?k=test_r_id&mp=J75AJcKB%2BFQ36aSIidMM6aHIfUEMGUcI6aSIi%2BMM6deI6deI6aSI6aSI6deI6deI6dxQhbx6HFcuHdeIidMM6aSIidMFGUEB6deI6deORgxQY%2BSzHoR14BRMRUEe6%2B2I6aSI6deI6deIideIideI")
			So(ad.NoticeURL, ShouldEqual, "http://test_domain/click?k=test_r_id&mp=J75AJcKB%2BFQ36aSIidMM6aHIfUEMGUcI6aSIi%2BMM6deI6deI6aSI6aSI6deI6deI6dxQhbx6HFcuHdeIidMM6aSIidMFGUEB6deI6deORgxQY%2BSzHoR14BRMRUEe6%2B2I6aSI6deI6deIideIideI")
		})
	})
}

func TestRenderMWParams(t *testing.T) {
	Convey("test RenderMWParams", t, func() {
		var params mvutil.Params
		var ad Ad
		var isPv bool

		Convey("isPv is true", func() {
			isPv = true
			paramsTest := mvutil.Params{
				Platform:    1,
				GAID:        "test_gaid",
				BackendID:   6,
				RequestID:   "test_r_id",
				Domain:      "test_domain",
				BackendList: []int32{5, 4, 3, 2, 1},
			}
			RenderMWParams(&paramsTest, &ad, isPv)
			So(paramsTest.MWadBackend, ShouldEqual, "5,4,3,2,1")
			So(paramsTest.MWadBackendData, ShouldEqual, "0")
			So(paramsTest.MWbackendConfig, ShouldEqual, "")
		})

		Convey("isPv is true, params.BackendList is empty", func() {
			isPv = true
			params = mvutil.Params{
				Platform:  1,
				GAID:      "test_gaid",
				BackendID: 6,
				RequestID: "test_r_id",
				Domain:    "test_domain",
			}
			RenderMWParams(&params, &ad, isPv)
			So(params.MWadBackend, ShouldEqual, "0")
			So(params.MWadBackendData, ShouldEqual, "0")
			So(params.MWbackendConfig, ShouldEqual, "")
		})

		Convey("isPv is false, params.BackendList is empty", func() {
			isPv = false
			params = mvutil.Params{
				Platform:  1,
				GAID:      "test_gaid",
				BackendID: 6,
				RequestID: "test_r_id",
				Domain:    "test_domain",
			}
			ad.VideoURL = "ad_url"
			RenderMWParams(&params, &ad, isPv)
			So(params.MWadBackend, ShouldEqual, "6")
			So(params.MWadBackendData, ShouldEqual, "6:0:1")
			So(params.MWbackendConfig, ShouldEqual, "6::3")
		})
	})
}

func TestRenderClickmode(t *testing.T) {
	Convey("test renderClickmode", t, func() {
		var params mvutil.Params
		var ad Ad
		var corsairCampaign corsair_proto.Campaign
		var r mvutil.RequestParams
		Convey("OriCampaignId is nil", func() {
			renderClickmode(&ad, &params, corsairCampaign, &r)
			So(ad.ClickMode, ShouldEqual, 0)
			So(params.ThirdPartyABTestRes.ClickmodeRes, ShouldEqual, 0)
		})
		Convey("OriCampaignId is not number", func() {
			str := "abc"
			corsairCampaign.OriCampaignId = &str
			renderClickmode(&ad, &params, corsairCampaign, &r)
			So(ad.ClickMode, ShouldEqual, 0)
			So(params.ThirdPartyABTestRes.ClickmodeRes, ShouldEqual, 0)
		})
		Convey("campaignInfo empty", func() {
			str := "123"
			corsairCampaign.OriCampaignId = &str
			guard := Patch(extractor.GetCampaignInfo, func(camId int64) (camInfo *smodel.CampaignInfo, ifFind bool) {
				return &smodel.CampaignInfo{}, false
			})
			defer guard.Unpatch()
			renderClickmode(&ad, &params, corsairCampaign, &r)
			So(ad.ClickMode, ShouldEqual, 0)
			So(params.ThirdPartyABTestRes.ClickmodeRes, ShouldEqual, 0)
		})
		Convey("clickmode config without 6&12", func() {
			str := "123"
			corsairCampaign.OriCampaignId = &str
			guard := Patch(extractor.GetCampaignInfo, func(camId int64) (camInfo *smodel.CampaignInfo, ifFind bool) {
				return &smodel.CampaignInfo{}, true
			})
			defer guard.Unpatch()
			guard = Patch(getJumpTypeArr, func(r *mvutil.RequestParams, camInfo *smodel.CampaignInfo, params *mvutil.Params) map[string]int32 {
				return map[string]int32{
					"5": 100,
				}
			})
			defer guard.Unpatch()
			renderClickmode(&ad, &params, corsairCampaign, &r)
			So(ad.ClickMode, ShouldEqual, 0)
			So(params.ThirdPartyABTestRes.ClickmodeRes, ShouldEqual, 0)
		})
		Convey("clickmode config without 6&12 but not match", func() {
			str := "123"
			corsairCampaign.OriCampaignId = &str
			guard := Patch(extractor.GetCampaignInfo, func(camId int64) (camInfo *smodel.CampaignInfo, ifFind bool) {
				return &smodel.CampaignInfo{}, true
			})
			defer guard.Unpatch()
			guard = Patch(getJumpTypeArr, func(r *mvutil.RequestParams, camInfo *smodel.CampaignInfo, params *mvutil.Params) map[string]int32 {
				return map[string]int32{
					"6": 100,
				}
			})
			defer guard.Unpatch()
			guard = Patch(extractor.GetADNET_SWITCHS, func() (map[string]int, bool) {
				return map[string]int{
					"cmTestRate": 0,
				}, true
			})
			defer guard.Unpatch()
			guard = Patch(rand.Intn, func(n int) int {
				return 100
			})
			defer guard.Unpatch()
			renderClickmode(&ad, &params, corsairCampaign, &r)
			So(ad.ClickMode, ShouldEqual, 0)
			So(params.ThirdPartyABTestRes.ClickmodeRes, ShouldEqual, 2)
		})
		Convey("clickmode config without 6&12 and match", func() {
			str := "123"
			corsairCampaign.OriCampaignId = &str
			guard := Patch(extractor.GetCampaignInfo, func(camId int64) (camInfo *smodel.CampaignInfo, ifFind bool) {
				return &smodel.CampaignInfo{}, true
			})
			defer guard.Unpatch()
			guard = Patch(getJumpTypeArr, func(r *mvutil.RequestParams, camInfo *smodel.CampaignInfo, params *mvutil.Params) map[string]int32 {
				return map[string]int32{
					"6": 100,
				}
			})
			defer guard.Unpatch()
			guard = Patch(extractor.GetADNET_SWITCHS, func() (map[string]int, bool) {
				return map[string]int{
					"cmTestRate": 100,
				}, true
			})
			defer guard.Unpatch()
			guard = Patch(rand.Intn, func(n int) int {
				return 0
			})
			defer guard.Unpatch()
			renderClickmode(&ad, &params, corsairCampaign, &r)
			So(ad.ClickMode, ShouldEqual, 6)
			So(params.ThirdPartyABTestRes.ClickmodeRes, ShouldEqual, 1)
		})
	})
}

func TestMatchSdkVersion(t *testing.T) {
	// isMatch, err := MatchSdkVersion("mi_2.0.1", "mi_1.0.0", "mi_1.1.2")
	Convey("start test", t, func() {
		trans := func(version string) supply_mvutil.SDKVersionItem {
			return supply_mvutil.RenderSDKVersion(version)
		}
		isMatch1, _ := MatchSdkVersion(trans("mi_2.0.1"), trans("mi_1.0.0"), trans("mi_1.1.2"))
		isMatch2, _ := MatchSdkVersion(trans("mi_2.0.1"), trans("mi_1.0.0"), trans(""))
		isMatch3, _ := MatchSdkVersion(trans("mal_2.0.1"), trans("mi_1.0.0"), trans("mi_1.1.2"))
		isMatch4, _ := MatchSdkVersion(trans("mi_1.1.0"), trans(""), trans("mi_1.1.2"))
		isMatch5, _ := MatchSdkVersion(trans("mi_2.0.1"), trans("mi_1.0.0"), trans("mi_3.1.2"))
		isMatch6, _ := MatchSdkVersion(trans("mal_2.0.1"), trans("mi_1.0.0"), trans("mi_3.1.2"))
		So(isMatch1, ShouldEqual, false)
		So(isMatch2, ShouldEqual, true)
		So(isMatch3, ShouldEqual, false)
		So(isMatch4, ShouldEqual, true)
		So(isMatch5, ShouldEqual, true)
		So(isMatch6, ShouldEqual, false)
	})
}

//
func TestIsConfigV9Template(t *testing.T) {
	Convey("start test", t, func() {

		trans := func(version string) supply_mvutil.SDKVersionItem {
			return supply_mvutil.RenderSDKVersion(version)
		}

		Convey("nil", func() {
			res, err := IsSupportSDKV9Template(&mvutil.RequestParams{
				UnitInfo: &smodel.UnitInfo{
					TemplateConf: nil,
				},

				Param: mvutil.Params{
					FormatOrientation: 1,
					OSVersionCode:     11080000,
					//OSVersionCode: 14080000,
					//OSVersionCode: 14080000,
					//OSVersionCode: 14080000,
					//SDKVersion: "mi_7.0.5",
					//SDKVersion:    "mal_7.0.3",
					//SDKVersion:    "mi_7.0.3.2",
					//SDKVersion:    "mi_7",
					FormatSDKVersion: trans("mi_7.0.5"),
				},
			})
			fmt.Println(err)
			So(res, ShouldEqual, true)
		})
		Convey("match orientation", func() {
			res, err := IsSupportSDKV9Template(&mvutil.RequestParams{
				UnitInfo: &smodel.UnitInfo{
					TemplateConf: map[string]smodel.TemplateConf{
						"10001": {
							Orientation: 1,
							SdkVersion: map[string][]smodel.SdkVersionRule{
								"include": {
									{
										Max: "mi_7.0.9",
										Min: "mi_7.0.4",
									},
								},
							},
							OSMin: 0,
							OSMax: 0,
							VideoTemplate: []smodel.VideoTemplate{
								{
									Type:   8,
									Weight: 0,
								},
							},
						},
					},
				},
				Param: mvutil.Params{
					FormatOrientation: 0,
					OSVersionCode:     14080000,
					//SDKVersion:    "mi_7.0.3",
					FormatSDKVersion: trans("mi_7.0.3"),
				},
			})
			fmt.Println(err)
			So(res, ShouldEqual, true)
		})
		Convey("match sdk version", func() {
			res, err := IsSupportSDKV9Template(&mvutil.RequestParams{
				UnitInfo: &smodel.UnitInfo{
					TemplateConf: map[string]smodel.TemplateConf{
						"10001": {
							Orientation: 1,
							SdkVersion: map[string][]smodel.SdkVersionRule{
								"include": {
									{
										Max: "mi_7.0.9",
										Min: "mi_7.0.1",
									},
								},
								"exclude": {
									{
										Max: "mi_7.0.4",
										Min: "mi_7.0.4",
									},
								},
							},
							OSMin: 0,
							OSMax: 0,
							VideoTemplate: []smodel.VideoTemplate{
								{
									Type:   8,
									Weight: 0,
								},
								{
									Type:   9,
									Weight: 0,
								},
							},
						},
						"10002": {
							Orientation: 1,
							SdkVersion: map[string][]smodel.SdkVersionRule{
								"include": {
									{
										Max: "mal_7.0.9",
										Min: "mal_7.0.1",
									},
								},
								"exclude": {
									{
										Max: "mal_7.0.4",
										Min: "mal_7.0.4",
									},
								},
							},
							OSMin: 0,
							OSMax: 0,
							VideoTemplate: []smodel.VideoTemplate{
								{
									Type:   8,
									Weight: 0,
								},
								{
									Type:   9,
									Weight: 0,
								},
							},
						},
					},
				},
				Param: mvutil.Params{
					FormatOrientation: 1,
					OSVersionCode:     14080000,
					//SDKVersion:    "mi_7.0.8",
					//SDKVersion:    "mal_7.0.3",
					//SDKVersion:    "mi_7.0.3.2",
					//SDKVersion:    "mi_7",
					FormatSDKVersion: trans("mal_7.0.3"),
				},
			})

			fmt.Println(err)
			So(res, ShouldEqual, true)
		})
		Convey("match os version", func() {
			res, err := IsSupportSDKV9Template(&mvutil.RequestParams{
				UnitInfo: &smodel.UnitInfo{
					TemplateConf: map[string]smodel.TemplateConf{
						"10001": {
							Orientation: 1,
							SdkVersion: map[string][]smodel.SdkVersionRule{
								"include": {
									{
										Max: "mi_7.0.9",
										Min: "mi_7.0.1",
									},
								},
								"exclude": {
									{
										Max: "mi_7.0.4",
										Min: "mi_7.0.4",
									},
								},
							},
							OSMin: 11080000,
							OSMax: 12080000,
							VideoTemplate: []smodel.VideoTemplate{
								{
									Type:   9,
									Weight: 0,
								},
								{
									Type:   8,
									Weight: 0,
								},
								{
									Type:   2,
									Weight: 0,
								},
							},
						},
					},
				},
				Param: mvutil.Params{
					FormatOrientation: 1,
					OSVersionCode:     12000000,
					//OSVersionCode: 14080000,
					//OSVersionCode: 14080000,
					//OSVersionCode: 14080000,
					//SDKVersion: "mi_7.0.5",
					//SDKVersion:    "mal_7.0.3",
					//SDKVersion:    "mi_7.0.3.2",
					//SDKVersion:    "mi_7",
					FormatSDKVersion: trans("mi_7.0.5"),
				},
			})
			fmt.Println(err)
			So(res, ShouldEqual, true)
		})
	})
}

package output

import (
	"math/rand"
	"testing"

	"gitlab.mobvista.com/ADN/chasm/module/demand"
	chasm_extractor "gitlab.mobvista.com/ADN/chasm/module/extractor"

	. "github.com/bouk/monkey"
	. "github.com/smartystreets/goconvey/convey"
	"gitlab.mobvista.com/ADN/adnet/internal/corsair_proto"
	"gitlab.mobvista.com/ADN/adnet/internal/extractor"
	"gitlab.mobvista.com/ADN/adnet/internal/mvconst"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	"gitlab.mobvista.com/ADN/adnet/internal/protobuf"
	smodel "gitlab.mobvista.com/ADN/structs/model"
	"gitlab.mobvista.com/adserver/recommend_protocols/go/ad_server"
)

func TestRenderParams(t *testing.T) {
	var ad Ad
	var params mvutil.Params
	campaign := &smodel.CampaignInfo{}
	var corsairCampaign corsair_proto.Campaign

	Convey("test renderParams", t, func() {
		// Mock 方法
		guard := Patch(mvutil.GetGoTkClickID, func() string {
			return "test_go_click_id"
		})
		defer guard.Unpatch()

		guard = Patch(renderExtAttr, func(campaign *smodel.CampaignInfo) int {
			return 42
		})
		defer guard.Unpatch()

		guard = Patch(mvutil.GetVTALink, func(campaign *smodel.CampaignInfo) string {
			return "test_vta_link"
		})
		defer guard.Unpatch()

		guard = Patch(renderNativeVideoFlag, func(ad Ad, params *mvutil.Params) {
			params.ExtnativeVideo = 1
		})
		defer guard.Unpatch()

		guard = Patch(renderNativex, func(params *mvutil.Params, campaign *smodel.CampaignInfo) {
			params.Extnativex = 1
		})
		defer guard.Unpatch()

		guard = Patch(RenderExtBp, func(params *mvutil.Params, campaign *smodel.CampaignInfo) string {
			return "test_ext_bp"
		})
		defer guard.Unpatch()

		guard = Patch(renderExtra20, func(params *mvutil.Params, campaign *smodel.CampaignInfo) []string {
			return []string{"test_extr20_1", "test_extr20_2", "test_extr20_3"}
		})
		defer guard.Unpatch()

		guard2 := Patch(renderExtPlayable, func(p *mvutil.Params, c *smodel.CampaignInfo) string {
			return "test_ext_playable"
		})
		defer guard2.Unpatch()

		guard3 := Patch(RenderMWParams, func(params *mvutil.Params, ad *Ad, isPv bool) {
			params.MWadBackend = "1"
			params.MWadBackendData = "2"
			params.MWbackendConfig = "test_mw_backend_config"
		})
		defer guard3.Unpatch()

		Convey("campaign.PublisherID 为 0，返回 -1", func() {
			pID := int64(0)
			campaign.PublisherId = pID
			renderParams(ad, &params, campaign, corsairCampaign)
			So(params.Extra8, ShouldEqual, -1)
		})

		Convey("campaign.PublisherID 不为0 && AdvertiverId 不为 nil，返回 advID", func() {
			campaign.PublisherId = 1
			advID := int32(42)
			campaign.AdvertiserId = advID
			renderParams(ad, &params, campaign, corsairCampaign)
			So(params.Extra8, ShouldEqual, 42)
		})

		Convey("campaign.PublisherID && AdvertiverId 为 0，返回 -1", func() {
			campaign.PublisherId = 0
			campaign.AdvertiserId = 0
			renderParams(ad, &params, campaign, corsairCampaign)
			So(params.Extra8, ShouldEqual, -1)
		})

		Convey("campaign.PublisherID 为 正常值，返回该值 0", func() {
			pID := int64(42)
			campaign.PublisherId = pID
			advID := int32(21)
			campaign.AdvertiserId = advID
			renderParams(ad, &params, campaign, corsairCampaign)
			So(params.Extra8, ShouldEqual, 21)
		})

		Convey("corsairCampaign.AdTemplate 为空，赋值 0", func() {
			corsairCampaign.AdTemplate = nil
			renderParams(ad, &params, campaign, corsairCampaign)
			So(params.Extra16, ShouldEqual, 0)
		})

		Convey("corsairCampaign.AdTemplate 不为空，赋值 3", func() {
			adTemplate := ad_server.ADTemplate(3)
			corsairCampaign.AdTemplate = &adTemplate
			renderParams(ad, &params, campaign, corsairCampaign)
			So(params.Extra16, ShouldEqual, 3)
		})

		Convey("测试 Extalgo", func() {
			params.AlgoMap = map[int64]string{1: "test_algo_1"}
			campaign.CampaignId = 5
			renderParams(ad, &params, campaign, corsairCampaign)
			So(params.Extalgo, ShouldEqual, "")

			campaign.CampaignId = 1
			renderParams(ad, &params, campaign, corsairCampaign)
			So(params.Extalgo, ShouldEqual, "test_algo_1")
		})

		Convey("测试 AdvertiserID", func() {
			advID := int32(5)
			campaign.AdvertiserId = advID
			adsourceID := int32(6)
			campaign.AdSourceId = adsourceID
			cType := int32(7)
			campaign.Ctype = cType
			source := int32(8)
			campaign.Source = source
			cOfferType := int32(9)
			corsairCampaign.OfferType = &cOfferType
			params.FlowTagID = 1

			renderParams(ad, &params, campaign, corsairCampaign)

			So(params.AdvertiserID, ShouldEqual, 5)
			So(params.Extra13, ShouldEqual, 6)
			So(params.Extctype, ShouldEqual, 7)
			// So(params.RequestID, ShouldEqual, "test_go_click_id")
			So(params.Extattr, ShouldEqual, 42)
			So(params.ExtinstallFrom, ShouldEqual, 1)
			//So(params.Extbp, ShouldEqual, "test_ext_bp")
			So(params.Extsource, ShouldEqual, 8)
			So(params.OfferType, ShouldEqual, 9)
			So(params.Extra20List, ShouldResemble, []string{"test_extr20_1", "test_extr20_2", "test_extr20_3"})
			So(params.ExtPlayableList, ShouldEqual, "test_ext_playable")
			So(params.MWadBackend, ShouldEqual, "1")
			So(params.MWadBackendData, ShouldEqual, "2")
			So(params.MWbackendConfig, ShouldEqual, "test_mw_backend_config")
		})
	})
}

func TestRenderDeleteDevID(t *testing.T) {
	var ad Ad
	var params mvutil.Params
	campaign := &smodel.CampaignInfo{}

	Convey("campaign.SendDeviceidRate 为空，不做任何处理", t, func() {
		renderDeleteDevID(&ad, &params, campaign)
		So(ad.ClickMode, ShouldEqual, 0)
	})

	Convey("campaign.SendDeviceidRate 不在 [0, 100)，不作任何处理", t, func() {
		rate := int32(100)
		campaign = &smodel.CampaignInfo{SendDeviceidRate: rate}
		renderDeleteDevID(&ad, &params, campaign)
		So(ad.ClickMode, ShouldEqual, 0)
	})

	Convey("campaign.SendDeviceidRate 为 [0, 100) 时，随机", t, func() {
		rate := int32(0)
		campaign = &smodel.CampaignInfo{SendDeviceidRate: rate}
		renderDeleteDevID(&ad, &params, campaign)
		So(ad.ClickMode, ShouldBeIn, []int{0, 5})
		So(params.ExtdeleteDevid, ShouldBeIn, []int{1, 2})
		So(params.Extra10, ShouldBeIn, []string{"0", "5"})
	})
}

func TestRenderExtra20(t *testing.T) {
	Convey("默认值，0 转换为空字符串", t, func() {
		params := &mvutil.Params{}
		campaign := &smodel.CampaignInfo{}
		res := renderExtra20(params, campaign)
		So(res, ShouldResemble, []string{"", "", "", "", "", "", ""})
	})

	Convey("默认值，0 转换为空字符串", t, func() {
		params := &mvutil.Params{
			CampaignID:   100,
			AdvertiserID: 101,
			OfferType:    0,
			CUA:          103,
		}
		configVBA := smodel.ConfigVBA{Status: 1, UseVBA: 1}
		campaign := &smodel.CampaignInfo{
			ConfigVBA: &configVBA,
		}
		res := renderExtra20(params, campaign)
		So(res, ShouldResemble, []string{"100", "101", "", "", "2", "103", ""})
	})
}

func TestRenderExtBp(t *testing.T) {
	Convey("campaign 为空，使用默认值填充", t, func() {
		campaign := &smodel.CampaignInfo{}
		params := &mvutil.Params{}
		res := RenderExtBp(params, campaign)
		So(res, ShouldResemble, "[\"0\",\"0\",\"0\",\"0\",\"0\",\"0\"]")
	})

	Convey("campaign 为空，使用默认值填充", t, func() {
		oriPrice := float64(1.23)
		price := float64(6.66)
		ctypeDef := int32(42)
		costType := int32(21)
		params := &mvutil.Params{
			PriceIn:             oriPrice,
			PriceOut:            price,
			LocalCurrency:       1,
			LocalChannelPriceIn: 2.3,
		}

		campaign := &smodel.CampaignInfo{
			OriPrice: oriPrice,
			Price:    price,
			Ctype:    ctypeDef,
			CostType: costType,
		}

		res := RenderExtBp(params, campaign)
		So(res, ShouldResemble, "[\"1.23\",\"6.66\",\"42\",\"21\",\"1\",\"2.3\"]")
	})
}

func TestRenderNativex(t *testing.T) {
	Convey("belong type 为空，nx 标示 2", t, func() {
		params := mvutil.Params{}
		campaign := &smodel.CampaignInfo{}
		renderNativex(&params, campaign)

		So(params.Extnativex, ShouldEqual, 2)
	})

	Convey("belong type != 1，nx 标示 2", t, func() {
		params := mvutil.Params{}
		belong := int32(2)
		campaign := &smodel.CampaignInfo{BelongType: belong}
		renderNativex(&params, campaign)

		So(params.Extnativex, ShouldEqual, 2)
	})

	Convey("belong type == 1，nx 标示 1", t, func() {
		params := mvutil.Params{}
		belong := int32(1)
		campaign := &smodel.CampaignInfo{BelongType: belong}
		renderNativex(&params, campaign)

		So(params.Extnativex, ShouldEqual, 1)
	})
}

func TestRenderNativeVideoFlag(t *testing.T) {
	Convey("params 空，返回 0", t, func() {
		ad := Ad{}
		params := mvutil.Params{}
		renderNativeVideoFlag(ad, &params)
		So(params.ExtnativeVideo, ShouldEqual, 0)
	})

	Convey("ApiVersion 大于 1.0，返回 0", t, func() {
		ad := Ad{}
		params := mvutil.Params{
			ApiVersion: float64(2.0),
		}
		renderNativeVideoFlag(ad, &params)
		So(params.ExtnativeVideo, ShouldEqual, 0)
	})

	Convey("否则 native video，返回 1", t, func() {
		ad := Ad{
			VideoURL: "http://test.video.com/1.mp4",
		}
		params := mvutil.Params{
			AdType:       42,
			VideoVersion: "test_version",
		}
		renderNativeVideoFlag(ad, &params)
		So(params.ExtnativeVideo, ShouldEqual, 1)
	})
}

func TestRenderVideoInfo(t *testing.T) {
	var ad Ad
	var params mvutil.Params
	var r mvutil.RequestParams
	guard := Patch(extractor.GetADNET_SWITCHS, func() (map[string]int, bool) {
		return map[string]int{}, true
	})
	defer guard.Unpatch()

	Convey("测试整理 videoEndType", t, func() {
		ad = Ad{
			VideoEndType:            0,
			PlayableAdsWithoutVideo: 0,
		}

		Convey("如果 apiVersion >= 1.2，不处理", func() {
			params = mvutil.Params{
				ApiVersion: float64(2.0),
				Platform:   1,
			}
			r = mvutil.RequestParams{
				UnitInfo: &smodel.UnitInfo{
					Unit: smodel.Unit{VideoEndType: 10},
				},
			}
			renderVideoInfo(&ad, &params, &r)

			So(ad.VideoEndType, ShouldEqual, 10)
		})

		Convey("如果 apiVersion < 1.2", func() {
			params = mvutil.Params{
				ApiVersion: float64(1.0),
				Platform:   1,
			}

			Convey("如果 unit 维度 videoEndType <= 0，返回 2", func() {
				r = mvutil.RequestParams{
					UnitInfo: &smodel.UnitInfo{
						Unit: smodel.Unit{VideoEndType: 0},
					},
				}
				renderVideoInfo(&ad, &params, &r)

				So(ad.VideoEndType, ShouldEqual, 2)
			})

			Convey("如果 unit 维度 videoEndType > 0", func() {
				Convey("如果 unit 维度 videoEndType > 5，返回 2", func() {
					r = mvutil.RequestParams{
						UnitInfo: &smodel.UnitInfo{
							Unit: smodel.Unit{VideoEndType: 10},
						},
					}
					renderVideoInfo(&ad, &params, &r)

					So(ad.VideoEndType, ShouldEqual, 2)
				})

				Convey("如果 unit 维度 videoEndType <= 5，返回 videoEndType", func() {
					r = mvutil.RequestParams{
						UnitInfo: &smodel.UnitInfo{
							Unit: smodel.Unit{VideoEndType: 3},
						},
					}
					renderVideoInfo(&ad, &params, &r)

					So(ad.VideoEndType, ShouldEqual, 3)
				})
			})
		})
	})

	Convey("测试整理 storekit", t, func() {
		ad = Ad{
			VideoEndType:            10,
			PlayableAdsWithoutVideo: 0,
		}
		Convey("如果 platform 为 iOS", func() {
			params.Platform = 2

			Convey("如果 apiVersion >= 1.3 && adtype 为 native && Extnvt2 = 5，storekit 为 1", func() {
				params.ApiVersion = float64(1.5)
				params.AdType = 42
				params.Extnvt2 = int32(5)

				renderVideoInfo(&ad, &params, &r)
				So(ad.Storekit, ShouldEqual, 1)
			})

			Convey("ad.Rv.VideoTemplate == 4，storekit 为 1", func() {
				ad.Rv.VideoTemplate = 4

				renderVideoInfo(&ad, &params, &r)
				So(ad.Storekit, ShouldEqual, 1)
			})

			Convey("ad.VideoEndType == 6，storekit 为 1", func() {
				ad.VideoEndType = 6

				renderVideoInfo(&ad, &params, &r)
				So(ad.Storekit, ShouldEqual, 1)
			})

			Convey("其他情况，storekit 为 0", func() {
				params.ApiVersion = float64(0.1)
				ad.Rv.VideoTemplate = 1
				ad.VideoEndType = 1

				renderVideoInfo(&ad, &params, &r)
				So(ad.Storekit, ShouldEqual, 0)
			})
		})
	})
	//素材三期迁移至renderPlayableAdsWithoutVideo处理
	//Convey("测试整理 PlayableAdsWithoutVideo", t, func() {
	//	Convey("如果 PlayableAdsWithoutVideo 不为 0，返回原值", func() {
	//		ad.PlayableAdsWithoutVideo = 3
	//		renderVideoInfo(&ad, &params, &r)
	//		So(ad.PlayableAdsWithoutVideo, ShouldEqual, 3)
	//	})
	//
	//	Convey("如果 PlayableAdsWithoutVideo 为 0，返回 1", func() {
	//		ad.PlayableAdsWithoutVideo = 0
	//		renderVideoInfo(&ad, &params, &r)
	//		So(ad.PlayableAdsWithoutVideo, ShouldEqual, 1)
	//	})
	//})
}

func TestRenderLoopback(t *testing.T) {
	var ad *Ad
	campaign := &smodel.CampaignInfo{}

	Convey("campaign.loopback 为空，不做处理", t, func() {
		renderLoopback(ad, campaign)
		So(ad, ShouldBeNil)
	})

	Convey("loopback.rate 为空，不做处理", t, func() {
		back := smodel.LoopBack{}
		campaign.Loopback = &back
		renderLoopback(ad, campaign)
		So(ad, ShouldBeNil)
	})

	Convey("loopback.rate > 100，修改 loopback 内容", t, func() {
		adVal := Ad{}
		ad = &adVal
		rate := int32(200)
		domain := "test_domain"
		key := "test_key"
		val := "test_val"
		newLoopBack := smodel.LoopBack{
			Rate:   rate,
			Domain: domain,
			Key:    key,
			Value:  val,
		}
		campaign.Loopback = &newLoopBack
		renderLoopback(ad, campaign)
		So(ad.LoopBack, ShouldResemble, map[string]string{"domain": domain, "key": key, "value": val})
	})
}

//func TestRenderJumpTypeInfo(t *testing.T) {
//	Convey("test renderJumpTypeInfo", t, func() {
//		ad := Ad{}
//		r := &mvutil.RequestParams{}
//		campaign := &smodel.CampaignInfo{}
//		params := &mvutil.Params{}
//
//		Convey("空参数", func() {
//			renderJumpTypeInfo(&ad, r, campaign, params)
//			So(params.Extra10, ShouldEqual, "0")
//			So(ad.PackageName, ShouldEqual, "")
//			So(params.JumpType, ShouldEqual, "0")
//			So(ad.ClickMode, ShouldEqual, 0)
//		})
//
//		Convey("jumptype = JUMP_TYPE_CLIENT_DO_ALL(6) ", func() {
//			guard := Patch(RandJumpType, func(r *mvutil.RequestParams, campaign *smodel.CampaignInfo, params *mvutil.Params) string {
//				return "6"
//			})
//			defer guard.Unpatch()
//
//			sdkPkName := "test_sdk_package_name"
//			campaign.SdkPackageName = sdkPkName
//			renderJumpTypeInfo(&ad, r, campaign, params)
//			So(params.Extra10, ShouldEqual, "6")
//			So(ad.PackageName, ShouldEqual, "test_sdk_package_name")
//			So(params.JumpType, ShouldEqual, "6")
//			So(ad.ClickMode, ShouldEqual, 6)
//		})
//
//		Convey("jumptype = JUMP_TYPE_CLIENT_DO_ALL_PING_SERVER(12) ", func() {
//			guard := Patch(RandJumpType, func(r *mvutil.RequestParams, campaign *smodel.CampaignInfo, params *mvutil.Params) string {
//				return "12"
//			})
//			defer guard.Unpatch()
//
//			sdkPkName := "test_sdk_package_name_12"
//			campaign.SdkPackageName = sdkPkName
//			renderJumpTypeInfo(&ad, r, campaign, params)
//			So(params.Extra10, ShouldEqual, "12")
//			So(ad.PackageName, ShouldEqual, "test_sdk_package_name_12")
//			So(params.JumpType, ShouldEqual, "6")
//			So(ad.ClickMode, ShouldEqual, 6)
//		})
//
//		Convey("jumptype = JUMP_TYPE_CLIENT_SEND_DEVID_PING_SERVER(11) ", func() {
//			guard := Patch(RandJumpType, func(r *mvutil.RequestParams, campaign *smodel.CampaignInfo, params *mvutil.Params) string {
//				return "11"
//			})
//			defer guard.Unpatch()
//
//			sdkPkName := "test_sdk_package_name_11"
//			campaign.SdkPackageName = sdkPkName
//			renderJumpTypeInfo(&ad, r, campaign, params)
//			So(params.Extra10, ShouldEqual, "11")
//			So(ad.PackageName, ShouldEqual, "")
//			So(params.JumpType, ShouldEqual, "5")
//			So(ad.ClickMode, ShouldEqual, 5)
//		})
//
//		Convey("jumptype = 其他 ", func() {
//			guard := Patch(RandJumpType, func(r *mvutil.RequestParams, campaign *smodel.CampaignInfo, params *mvutil.Params) string {
//				return "1"
//			})
//			defer guard.Unpatch()
//
//			sdkPkName := "test_sdk_package_name_1"
//			campaign.SdkPackageName = sdkPkName
//			renderJumpTypeInfo(&ad, r, campaign, params)
//			So(params.Extra10, ShouldEqual, "1")
//			So(ad.PackageName, ShouldEqual, "")
//			So(params.JumpType, ShouldEqual, "1")
//			So(ad.ClickMode, ShouldEqual, 1)
//		})
//	})
//}

// func TestRenderPlayableNew(t *testing.T) {
// 	Convey("TestRenderPlayableNew", t, func() {
// 		var ad *Ad
// 		campaign := &smodel.CampaignInfo{}
// 		//var params *mvutil.Params
// 		r := mvutil.RequestParams{}

// 		Convey("非 rewarded video 或非 interstitial video，不处理", func() {
// 			adData := Ad{
// 				EndcardUrl:              "test_endcard_url",
// 				PlayableAdsWithoutVideo: 10,
// 				VideoURL:                "test_video_url",
// 			}
// 			ad = &adData
// 			params := &mvutil.Params{
// 				AdType:      1,
// 				Extplayable: 42,
// 			}

// 			renderPlayableNew(ad, campaign, params, &r)
// 			So(params.Extplayable, ShouldEqual, 42)
// 			So(ad.EndcardUrl, ShouldEqual, "test_endcard_url")
// 			So(ad.PlayableAdsWithoutVideo, ShouldEqual, 10)
// 			So(ad.VideoURL, ShouldEqual, "test_video_url")
// 		})

// 		Convey("campaign.Endcard 为空，不处理", func() {
// 			adData := Ad{
// 				EndcardUrl:              "test_endcard_url",
// 				PlayableAdsWithoutVideo: 10,
// 				VideoURL:                "test_video_url",
// 			}
// 			ad = &adData
// 			params := &mvutil.Params{
// 				AdType:      94,
// 				Extplayable: 42,
// 			}
// 			renderPlayableNew(ad, campaign, params, &r)
// 			So(params.Extplayable, ShouldEqual, 42)
// 			So(ad.EndcardUrl, ShouldEqual, "test_endcard_url")
// 			So(ad.PlayableAdsWithoutVideo, ShouldEqual, 10)
// 			So(ad.VideoURL, ShouldEqual, "test_video_url")
// 		})

// 		Convey("!Compare，不处理", func() {
// 			adData := Ad{
// 				EndcardUrl:              "test_endcard_url",
// 				PlayableAdsWithoutVideo: 10,
// 				VideoURL:                "test_video_url",
// 			}
// 			ad = &adData
// 			params := &mvutil.Params{
// 				AdType:      94,
// 				Extplayable: 42,
// 			}
// 			endcardData := map[string]*smodel.EndCard{
// 				"endcard_1": &smodel.EndCard{},
// 			}
// 			campaign.Endcard = endcardData
// 			r.Param.SDKVersion = "test_0.0.0"
// 			renderPlayableNew(ad, campaign, params, &r)
// 			So(params.Extplayable, ShouldEqual, 42)
// 			So(ad.EndcardUrl, ShouldEqual, "test_endcard_url")
// 			So(ad.PlayableAdsWithoutVideo, ShouldEqual, 10)
// 			So(ad.VideoURL, ShouldEqual, "test_video_url")
// 		})

// 		Convey("gaid 为空，不处理", func() {
// 			// todo
// 		})
// 	})
// }

// func TestRenderPlayable(t *testing.T) {
// 	Convey("test RenderPlayable", t, func() {
// 		var ad Ad
// 		campaign := &smodel.CampaignInfo{}
// 		var params mvutil.Params
// 		r := mvutil.RequestParams{}

// 		guard := Patch(Compare, func(r *mvutil.RequestParams, compareType string) bool {
// 			return true
// 		})
// 		defer guard.Unpatch()

// 		params.AdType = int32(94)

// 		Convey("params.AdType != 94，不做处理", func() {
// 			params.AdType = int32(1)
// 			renderPlayable(&ad, campaign, &params, &r)
// 		})

// 		Convey("Compare 返回 false，不做处理", func() {
// 			guard = Patch(Compare, func(r *mvutil.RequestParams, compareType string) bool {
// 				return false
// 			})
// 			defer guard.Unpatch()
// 			renderPlayable(&ad, campaign, &params, &r)
// 		})

// 		Convey("GetPlayableTest 返回 false，不做处理", func() {
// 			guard := Patch(extractor.GetPlayableTest, func() (map[int64]mvutil.PlatableTest, bool) {
// 				playtable := map[int64]mvutil.PlatableTest{}
// 				return playtable, false
// 			})
// 			defer guard.Unpatch()

// 			renderPlayable(&ad, campaign, &params, &r)
// 		})

// 		Convey("GetPlayableTest 返回 true，但没有 campaignId 的配置，不做处理", func() {
// 			guard := Patch(extractor.GetPlayableTest, func() (map[int64]mvutil.PlatableTest, bool) {
// 				playtable := map[int64]mvutil.PlatableTest{}
// 				return playtable, true
// 			})
// 			defer guard.Unpatch()
// 			campaign.CampaignId = int64(2)

// 			renderPlayable(&ad, campaign, &params, &r)
// 		})

// 		Convey("GetPlayableTest 返回 true，有 campaignId 的配置，GetRandConsiderZero 返回 -1，不做处理", func() {
// 			guard := Patch(extractor.GetPlayableTest, func() (map[int64]mvutil.PlatableTest, bool) {
// 				playtable := map[int64]mvutil.PlatableTest{int64(2): mvutil.PlatableTest{}}
// 				return playtable, true
// 			})
// 			defer guard.Unpatch()

// 			guard = Patch(mvutil.GetRandConsiderZero, func(gaid string, idfa string, salt string, randSum int) int {
// 				return -1
// 			})
// 			defer guard.Unpatch()

// 			campaign.CampaignId = int64(2)

// 			renderPlayable(&ad, campaign, &params, &r)
// 		})

// 		Convey("GetPlayableTest 返回 true，有 campaignId 的配置，GetRandConsiderZero 不返回 -1，case1", func() {
// 			guard := Patch(extractor.GetPlayableTest, func() (map[int64]mvutil.PlatableTest, bool) {
// 				playtable := map[int64]mvutil.PlatableTest{int64(2): mvutil.PlatableTest{
// 					Url:  "test_url",
// 					Rate: map[int]int{1: 88},
// 				}}
// 				return playtable, true
// 			})
// 			defer guard.Unpatch()

// 			guard = Patch(mvutil.GetRandConsiderZero, func(gaid string, idfa string, salt string, randSum int) int {
// 				return 1
// 			})
// 			defer guard.Unpatch()

// 			campaign.CampaignId = int64(2)

// 			renderPlayable(&ad, campaign, &params, &r)
// 			So(params.Extplayable, ShouldEqual, 1)
// 			So(ad.EndcardUrl, ShouldEqual, "")
// 			So(ad.PlayableAdsWithoutVideo, ShouldEqual, 0)
// 			So(ad.VideoURL, ShouldEqual, "")
// 		})

// 		Convey("GetPlayableTest 返回 true，有 campaignId 的配置，GetRandConsiderZero 不返回 -1，case2", func() {
// 			guard := Patch(extractor.GetPlayableTest, func() (map[int64]mvutil.PlatableTest, bool) {
// 				playtable := map[int64]mvutil.PlatableTest{int64(2): mvutil.PlatableTest{
// 					Url:  "test_url",
// 					Rate: map[int]int{2: 77},
// 				}}
// 				return playtable, true
// 			})
// 			defer guard.Unpatch()

// 			guard = Patch(mvutil.GetRandConsiderZero, func(gaid string, idfa string, salt string, randSum int) int {
// 				return 1
// 			})
// 			defer guard.Unpatch()

// 			campaign.CampaignId = int64(2)

// 			renderPlayable(&ad, campaign, &params, &r)
// 			So(params.Extplayable, ShouldEqual, 2)
// 			So(ad.EndcardUrl, ShouldEqual, "http://test_url")
// 			So(ad.PlayableAdsWithoutVideo, ShouldEqual, 0)
// 			So(ad.VideoURL, ShouldEqual, "")
// 		})

// 		Convey("GetPlayableTest 返回 true，有 campaignId 的配置，GetRandConsiderZero 不返回 -1，case3", func() {
// 			guard := Patch(extractor.GetPlayableTest, func() (map[int64]mvutil.PlatableTest, bool) {
// 				playtable := map[int64]mvutil.PlatableTest{int64(2): mvutil.PlatableTest{
// 					Url:  "test_url",
// 					Rate: map[int]int{3: 66},
// 				}}
// 				return playtable, true
// 			})
// 			defer guard.Unpatch()

// 			guard = Patch(mvutil.GetRandConsiderZero, func(gaid string, idfa string, salt string, randSum int) int {
// 				return 1
// 			})
// 			defer guard.Unpatch()

// 			campaign.CampaignId = int64(2)

// 			renderPlayable(&ad, campaign, &params, &r)
// 			So(params.Extplayable, ShouldEqual, 3)
// 			So(ad.EndcardUrl, ShouldEqual, "http://test_url")
// 			So(ad.PlayableAdsWithoutVideo, ShouldEqual, 2)
// 			So(ad.VideoURL, ShouldEqual, "")
// 		})
// 	})
// }

func TestRenderLinkType(t *testing.T) {
	var ad Ad
	campaign := &smodel.CampaignInfo{}
	var params mvutil.Params
	params.UnitID = 12345

	Convey("linkType 测试", t, func() {
		guard := Patch(extractor.GetLINKTYPE_UNITID, func() (map[string]int, bool) {
			lt := map[string]int{
				"54297": 9,
				"54301": 9,
			}
			return lt, true
		})
		defer guard.Unpatch()

		Convey("openType 存在，api version < 1.3，openType 为 8、9 时，campaignType 为 4，LinkType 为 4", func() {
			cOpenType := int32(8)
			campaign.OpenType = cOpenType
			params.ApiVersion = float64(0.1)
			params.RequestPath = "/openapi/ad/v3"
			renderLinkType(&ad, campaign, &params)
			So(ad.CampaignType, ShouldEqual, 4)
			So(params.LinkType, ShouldEqual, 4)
		})

		Convey("openType 存在，为mp的，api version < 1.3，openType 为 8、9 时，campaignType 为 8，LinkType 为 8,case1", func() {
			cOpenType := int32(8)
			campaign.OpenType = cOpenType
			params.ApiVersion = float64(0.1)
			params.RequestPath = "/mpapi/ad"
			renderLinkType(&ad, campaign, &params)
			So(ad.CampaignType, ShouldEqual, 8)
			So(params.LinkType, ShouldEqual, 8)
		})

		Convey("openType 存在，为mp的，api version < 1.3，openType 为 8、9 时，campaignType 为 8，LinkType 为 8,case2", func() {
			cOpenType := int32(8)
			campaign.OpenType = cOpenType
			params.SDKVersion = "mp_5.1.1"
			params.ApiVersion = float64(0.1)
			params.RequestPath = "/mpapi/ad"
			renderLinkType(&ad, campaign, &params)
			So(ad.CampaignType, ShouldEqual, 8)
			So(params.LinkType, ShouldEqual, 8)
		})

		Convey("openType 存在，mpsdkversion <4.xx，openType 为 8、9 时，campaignType 为 4，LinkType 为 4", func() {
			cOpenType := int32(8)
			campaign.OpenType = cOpenType
			params.SDKVersion = "mp_3.1.1"
			params.RequestPath = "/mpapi/ad"
			renderLinkType(&ad, campaign, &params)
			So(ad.CampaignType, ShouldEqual, 4)
			So(params.LinkType, ShouldEqual, 4)
		})

		Convey("OpenType 不存在，CampaignType 不为空", func() {

			Convey("CampaignType in [5,6,7]，CampaignType & LinkType 设置为 4", func() {
				cType := int32(5)
				campaign.CampaignType = cType
				renderLinkType(&ad, campaign, &params)
				So(ad.CampaignType, ShouldEqual, 4)
				So(params.LinkType, ShouldEqual, 4)

				cType = int32(6)
				campaign.CampaignType = cType
				renderLinkType(&ad, campaign, &params)
				So(ad.CampaignType, ShouldEqual, 4)
				So(params.LinkType, ShouldEqual, 4)

				cType = int32(7)
				campaign.CampaignType = cType
				renderLinkType(&ad, campaign, &params)
				So(ad.CampaignType, ShouldEqual, 4)
				So(params.LinkType, ShouldEqual, 4)
			})

			Convey("CampaignType not in [5,6,7]，CampaignType & LinkType 设置为 campaign.CampaignType", func() {
				cType := int32(1)
				campaign.CampaignType = cType
				renderLinkType(&ad, campaign, &params)
				So(ad.CampaignType, ShouldEqual, 4)
				So(params.LinkType, ShouldEqual, 4)
			})

			Convey("CampaignType not in [5,6,7]，in 8||9 CampaignType & LinkType 设置为 campaign.CampaignType", func() {
				cType := int32(8)
				campaign.CampaignType = cType
				params.SDKVersion = "mp_3.1.1"
				renderLinkType(&ad, campaign, &params)
				So(ad.CampaignType, ShouldEqual, 4)
				So(params.LinkType, ShouldEqual, 4)
			})

		})

		//Convey("OpenType 不存在，CampaignType 为空，CampaignType & LinkType 设置为 0", func() {
		//	campaign.OpenType = nil
		//	campaign.CampaignType = nil
		//	renderLinkType(&ad, campaign, &params)
		//	So(ad.CampaignType, ShouldEqual, 0)
		//	So(params.LinkType, ShouldEqual, 0)
		//})
	})
}

func TestRenderOfferwall(t *testing.T) {
	Convey("test renderOfferwall", t, func() {
		var ad Ad
		r := mvutil.RequestParams{UnitInfo: &smodel.UnitInfo{}}
		var offerType int32
		var price float64

		Convey("offerType = 99，不做处理", func() {
			offerType = int32(99)
			renderOfferwall(&ad, &r, offerType, price)
		})

		Convey("r.Param.AdType != 278，不做处理", func() {
			r.Param.AdType = int32(277)
			renderOfferwall(&ad, &r, offerType, price)
		})

		Convey("处理逻辑", func() {
			r.Param.AdType = int32(278)
			r.UnitInfo.VirtualReward = smodel.VirtualReward{
				Name:         "test_name",
				ExchangeRate: 2,
				StaticReward: 2,
			}

			guard := Patch(extractor.GetOfferwallGuidelines, func() (map[int32]string, bool) {
				return map[int32]string{int32(1): "test_val_1"}, true
			})
			defer guard.Unpatch()

			Convey("offer type = 1", func() {
				offerType = int32(1)
				renderOfferwall(&ad, &r, offerType, price)

				So(ad.Guidelines, ShouldEqual, "test_val_1")
				So(ad.RewardName, ShouldEqual, "test_name")
				So(ad.RewardAmount, ShouldEqual, 2)
			})

			Convey("offer type = 2", func() {
				offerType = int32(2)
				price = float64(10)
				renderOfferwall(&ad, &r, offerType, price)

				So(ad.Guidelines, ShouldEqual, "")
				So(ad.RewardName, ShouldEqual, "test_name")
				So(ad.RewardAmount, ShouldEqual, 20)
			})
		})
	})
}

func TestRenderImpUrl(t *testing.T) {
	var ad Ad
	campaign := &smodel.CampaignInfo{}
	r := &mvutil.RequestParams{
		AppInfo: &smodel.AppInfo{
			App: smodel.App{
				DevinfoEncrypt: 1,
			},
		},
	}
	params := &mvutil.Params{
		GAID:                "test_gaid",
		IDFA:                "test_idfa",
		AndroidID:           "test_andorid_id",
		IMEI:                "test_imei",
		MAC:                 "test_mac",
		ExtfinalPackageName: "test_final_packagename",
		Extra14:             int64(2),
	}

	Convey("整理 AdvImp", t, func() {
		campaign.AdUrlList = nil
		sec1 := int32(41)
		sec2 := int32(42)
		url1 := "http://www.test.com?gaid={gaid}&imei={imei}"
		url2 := "url2?mac={mac}&package_name={package_name}&sub={subId}"
		url3 := "url3?mac={idfa}&sub={subId}&dev={devId}"

		advImp := []*smodel.AdvImp{
			&smodel.AdvImp{Sec: sec1, Url: url1},
			&smodel.AdvImp{Url: url2},
			&smodel.AdvImp{Sec: sec2, Url: url3},
		}
		campaign.AdvImp = advImp
		renderImpUrl(&ad, r, campaign, params)
		So(len(ad.AdvImp), ShouldEqual, 3)
		So(ad.AdvImp[0].Sec, ShouldEqual, 41)
		So(ad.AdvImp[1].Sec, ShouldEqual, 0)
		So(ad.AdvImp[2].Sec, ShouldEqual, 42)
		So(ad.AdvImp[0].Url, ShouldEqual, "http://www.test.com?gaid=test_gaid&imei=test_imei")
		So(ad.AdvImp[1].Url, ShouldEqual, "url2?mac=test_mac&package_name=test_final_packagename&sub=2")
		So(ad.AdvImp[2].Url, ShouldEqual, "url3?mac=test_idfa&sub=2&dev=test_andorid_id")
	})

	Convey("整理 AdURLList", t, func() {
		campaign.AdvImp = nil

		url1 := "http://www.test.com?gaid={gaid}&imei={imei}"
		url2 := "url2?mac={mac}&package_name={package_name}&sub={subId}"
		url3 := "url3?mac={idfa}&sub={subId}&dev={devId}"

		adURLs := []string{
			url1,
			url2,
			url3,
		}
		campaign.AdUrlList = adURLs
		renderImpUrl(&ad, r, campaign, params)
		So(len(ad.AdURLList), ShouldEqual, 3)
		So(ad.AdURLList[0], ShouldEqual, "http://www.test.com?gaid=test_gaid&imei=test_imei")
		So(ad.AdURLList[1], ShouldEqual, "url2?mac=test_mac&package_name=test_final_packagename&sub=2")
		So(ad.AdURLList[2], ShouldEqual, "url3?mac=test_idfa&sub=2&dev=test_andorid_id")
	})
}

func TestHandleCTAButton(t *testing.T) {
	var campaignType int32
	params := &mvutil.Params{}

	Convey("空内容，返回 install", t, func() {
		res := handleCTAButton(campaignType, params)
		So(res, ShouldEqual, "install")
	})

	Convey("campaignType in [1,2,3]", t, func() {
		Convey("campaignType = 1", func() {
			Convey("CTALang 存在，返回对应 CTA", func() {
				campaignType := int32(1)
				params.Language = "zh-TW"
				res := handleCTAButton(campaignType, params)
				So(res, ShouldEqual, "安裝")
			})
			Convey("CTALang 不存在，返回 install", func() {
				campaignType := int32(1)
				params.Language = "test-none-exsit-lang"
				res := handleCTAButton(campaignType, params)
				So(res, ShouldEqual, "install")
			})
		})

		Convey("campaignType = 2", func() {
			Convey("CTALang 存在，返回对应 CTA", func() {
				campaignType := int32(1)
				params.Language = "ko"
				res := handleCTAButton(campaignType, params)
				So(res, ShouldEqual, "설치")
			})
		})

		Convey("campaignType = 3", func() {
			Convey("CTALang 存在，返回对应 CTA", func() {
				campaignType := int32(1)
				params.Language = "ar-"
				res := handleCTAButton(campaignType, params)
				So(res, ShouldEqual, "تثبيت")
			})
		})
	})

	Convey("campaignType in [6,7]", t, func() {
		Convey("campaignType = 6", func() {
			Convey("CTAViewLang 存在，返回对应 CTA", func() {
				campaignType := int32(6)
				params.Language = "ja-test"
				res := handleCTAButton(campaignType, params)
				So(res, ShouldEqual, "もっと")
			})
			Convey("CTAViewLang 不存在，返回 install", func() {
				campaignType := int32(6)
				params.Language = "test-none-exsit-lang"
				res := handleCTAButton(campaignType, params)
				So(res, ShouldEqual, "install")
			})
		})

		Convey("campaignType = 7", func() {
			Convey("CTAViewLang 存在，返回对应 CTA", func() {
				campaignType := int32(7)
				params.Language = "zh-Hant"
				res := handleCTAButton(campaignType, params)
				So(res, ShouldEqual, "查看")
			})
		})
	})
}

func TestHandleLanguage(t *testing.T) {
	var language string

	Convey("不包含 -，直接返回", t, func() {
		language = "test"
		res := handleLanguage(language)
		So(res, ShouldEqual, "test")
	})

	Convey("包含 zh-Hant，返回 zh-Hant", t, func() {
		language = "zh-Hant"
		res := handleLanguage(language)
		So(res, ShouldEqual, "zh-Hant")

		language = "test-zh-Hant-a"
		res = handleLanguage(language)
		So(res, ShouldEqual, "zh-Hant")
	})

	Convey("包含 zh-Hans，返回 zh-Hans", t, func() {
		language = "zh-Hans"
		res := handleLanguage(language)
		So(res, ShouldEqual, "zh-Hans")

		language = "test-zh-Hans-a"
		res = handleLanguage(language)
		So(res, ShouldEqual, "zh-Hans")
	})

	Convey("包含 zh-TW，返回 zh-Hant", t, func() {
		language = "zh-TW"
		res := handleLanguage(language)
		So(res, ShouldEqual, "zh-Hant")

		language = "test-zh-TW-a"
		res = handleLanguage(language)
		So(res, ShouldEqual, "zh-Hant")
	})

	Convey("其他包含 -，返回 - 前部分", t, func() {
		language = "zh-tw"
		res := handleLanguage(language)
		So(res, ShouldEqual, "zh")

		language = "a-b-c-d"
		res = handleLanguage(language)
		So(res, ShouldEqual, "a")
	})
}

func TestHandleNumberRating(t *testing.T) {
	var num int
	var res int
	Convey("在 (0, 1000) 之间，返回该值", t, func() {
		num = 1
		res = HandleNumberRating(num)
		So(res, ShouldBeGreaterThan, 10000)

		num = 12345
		res = HandleNumberRating(num)
		So(res, ShouldEqual, 12345)
	})

	Convey("不在 (0, 1000) 之间，返回 [10000, 50001] 间的随机值", t, func() {
		num = 0
		res = HandleNumberRating(num)
		So(res, ShouldBeGreaterThan, 9999)
		So(res, ShouldBeLessThan, 50002)

		num = 999999
		So(res, ShouldBeGreaterThan, 9999)
		So(res, ShouldBeLessThan, 50002)
	})
}

func TestRenderGInfo(t *testing.T) {
	Convey("将字段组合起来", t, func() {
		crType := ad_server.CreativeType(42)
		cr := protobuf.Creative{}
		cr.CreativeId = 1000
		cr.AdvCreativeId = "2000"
		cr.VideoResolution = "test"
		// cr := mvutil.Content{
		// 	CreativeId:      1000,
		// 	AdvCreativeId:   "2000",
		// 	VideoResolution: "test",
		// }

		res := renderGInfo(crType, &cr)
		So(res, ShouldEqual, "42,1000,2000,test,,0")
	})
}

func TestRenderCreative(t *testing.T) {
	Convey("test renderCreative", t, func() {
		var ad Ad
		var params mvutil.Params
		var corsairCampaign corsair_proto.Campaign
		var content map[ad_server.CreativeType]*protobuf.Creative
		var campaign smodel.CampaignInfo

		guard := Patch(extractor.GetADNET_SWITCHS, func() (map[string]int, bool) {
			return map[string]int{}, false
		})
		defer guard.Unpatch()

		Convey("空参数", func() {
			renderCreativePb(&ad, &params, corsairCampaign, content, &campaign)
		})

		Convey("real", func() {
			adEleTemplate := ad_server.AdElementTemplate(1990)
			ad = Ad{}
			params = mvutil.Params{}
			corsairCampaign = corsair_proto.Campaign{
				CampaignId:        "123",
				AdSource:          ad_server.ADSource(42),
				AdElementTemplate: &adEleTemplate,
				CreativeTypeIdMap: map[ad_server.CreativeType]int64{
					ad_server.CreativeType(1):   int64(100001),
					ad_server.CreativeType(2):   int64(100002),
					ad_server.CreativeType(301): int64(100301),
					ad_server.CreativeType(42):  int64(100042),
					ad_server.CreativeType(401): int64(1000401),
					ad_server.CreativeType(402): int64(1000402),
					ad_server.CreativeType(403): int64(1000403),
					ad_server.CreativeType(404): int64(1000404),
					ad_server.CreativeType(405): int64(1000405),
					ad_server.CreativeType(101): int64(1000101),
					ad_server.CreativeType(201): int64(1000201),
					ad_server.CreativeType(501): int64(1000501),
					ad_server.CreativeType(601): int64(1000601),
					ad_server.CreativeType(701): int64(1000701),
				},
				ImageSizeId: ad_server.ImageSizeEnum(404),
			}
			content = make(map[ad_server.CreativeType]*protobuf.Creative)
			isCamCre := int32(202)
			oriPrice := float64(101.1)
			campaign = smodel.CampaignInfo{
				IsCampaignCreative: isCamCre,
				OriPrice:           oriPrice,
				Price:              oriPrice,
			}

			renderCreativePb(&ad, &params, corsairCampaign, content, &campaign)
		})
	})
}

func TestGetMainAdvCreativeMap(t *testing.T) {
	bannerAdvCrIDMap := map[string]*mvutil.AdvCreativeMap{
		"image": &mvutil.AdvCreativeMap{
			AdvCreativeId: "image_val",
		},
		"icon": &mvutil.AdvCreativeMap{
			AdvCreativeId: "icon_val",
		},
	}

	Convey("banner 返回 image", t, func() {
		res := getMainAdvCreativeMap(bannerAdvCrIDMap)
		So(res.AdvCreativeId, ShouldEqual, "image_val")
	})

	interstitialAdvCrIDMap := map[string]*mvutil.AdvCreativeMap{
		"image": &mvutil.AdvCreativeMap{
			AdvCreativeId: "image_val",
		},
		"icon": &mvutil.AdvCreativeMap{
			AdvCreativeId: "icon_val",
		},
	}
	Convey("INTERSTITIAL 返回 image", t, func() {
		res := getMainAdvCreativeMap(interstitialAdvCrIDMap)
		So(res.AdvCreativeId, ShouldEqual, "image_val")
	})

	nativeAdvCrIDMap := map[string]*mvutil.AdvCreativeMap{
		"image": &mvutil.AdvCreativeMap{
			AdvCreativeId: "image_val",
		},
		"icon": &mvutil.AdvCreativeMap{
			AdvCreativeId: "icon_val",
		},
	}
	Convey("NATIVE 返回 image", t, func() {
		res := getMainAdvCreativeMap(nativeAdvCrIDMap)
		So(res.AdvCreativeId, ShouldEqual, "image_val")
	})

	appwallAdvCrIDMap := map[string]*mvutil.AdvCreativeMap{
		"icon": &mvutil.AdvCreativeMap{
			AdvCreativeId: "icon_val",
		},
	}
	Convey("APPWALL 返回 icon", t, func() {
		res := getMainAdvCreativeMap(appwallAdvCrIDMap)
		So(res.AdvCreativeId, ShouldEqual, "icon_val")
	})

	offerwallAdvCrIDMap := map[string]*mvutil.AdvCreativeMap{
		"icon": &mvutil.AdvCreativeMap{
			AdvCreativeId: "icon_val",
		},
	}
	Convey("OFFERWALL 返回 icon", t, func() {
		res := getMainAdvCreativeMap(offerwallAdvCrIDMap)
		So(res.AdvCreativeId, ShouldEqual, "icon_val")
	})

	nativeVideoAdvCrIDMap := map[string]*mvutil.AdvCreativeMap{
		"video": &mvutil.AdvCreativeMap{
			AdvCreativeId: "video_val",
		},
		"icon": &mvutil.AdvCreativeMap{
			AdvCreativeId: "icon_val",
		},
	}
	Convey("NATIVE_VIDEO 返回 video", t, func() {
		res := getMainAdvCreativeMap(nativeVideoAdvCrIDMap)
		So(res.AdvCreativeId, ShouldEqual, "video_val")
	})

	rewardedVideoAdvCrIDMap := map[string]*mvutil.AdvCreativeMap{
		"video": &mvutil.AdvCreativeMap{
			AdvCreativeId: "video_val",
		},
		"icon": &mvutil.AdvCreativeMap{
			AdvCreativeId: "icon_val",
		},
	}
	Convey("REWARDED_VIDEO 返回 video", t, func() {
		res := getMainAdvCreativeMap(rewardedVideoAdvCrIDMap)
		So(res.AdvCreativeId, ShouldEqual, "video_val")
	})
}

func TestGetQueryRAdType(t *testing.T) {
	var params mvutil.Params
	var res int
	ad := Ad{
		VideoURL: "test_video_url",
	}

	Convey("ADTypeBanner 时返回 1", t, func() {
		params.AdType = 2
		res = getQueryRAdType(params, ad)
		So(res, ShouldEqual, 1)
	})

	Convey("ADTypeInterstitial 时返回 1", t, func() {
		params.AdType = 29
		res = getQueryRAdType(params, ad)
		So(res, ShouldEqual, 2)
	})

	Convey("ADTypeNative 时", t, func() {
		params.AdType = 42
		Convey("videoVersion > 0 & videoUrl > 0 时，返回 7", func() {
			params.VideoVersion = "1.0"
			res = getQueryRAdType(params, ad)
			So(res, ShouldEqual, 7)
		})

		Convey("否则，返回 3", func() {
			params.VideoVersion = ""
			res = getQueryRAdType(params, ad)
			So(res, ShouldEqual, 3)
		})
	})

	Convey("ADTypeAppwall 时返回 1", t, func() {
		params.AdType = 3
		res = getQueryRAdType(params, ad)
		So(res, ShouldEqual, 4)
	})

	Convey("ADTypeOfferWall 时返回 1", t, func() {
		params.AdType = 278
		res = getQueryRAdType(params, ad)
		So(res, ShouldEqual, 5)
	})

	Convey("ADTypeRewardVideo 时返回 1", t, func() {
		params.AdType = 94
		res = getQueryRAdType(params, ad)
		So(res, ShouldEqual, 8)
	})

	Convey("其他时返回 0", t, func() {
		params.AdType = 1000
		res = getQueryRAdType(params, ad)
		So(res, ShouldEqual, 0)
	})
}

func TestGetCreativeGroupID(t *testing.T) {
	var res string

	Convey("空 list 返回空字符串", t, func() {
		res = getCreativeGroupID([]int{})
		So(res, ShouldEqual, "")
	})

	Convey("", t, func() {
		gidList := []int{99, 4, 42, 1, 9, 5, 6}
		res = getCreativeGroupID(gidList)
		So(res, ShouldEqual, "a2f8c46eb8dc36e4e53b68dcb4506830")
	})
}

// func TestRenderRewardVideoTemplate(t *testing.T) {
// 	Convey("test renderRewardVideoTemplate", t, func() {
// 		var ad Ad
// 		var params mvutil.Params
// 		r := mvutil.RequestParams{}
// 		var campaign smodel.CampaignInfo

// 		Convey("空请求", func() {
// 			renderRewardVideoTemplate(&ad, &params, &r, &campaign)
// 		})

// 		Convey("正常数据", func() {
// 			// guard := Patch(getRVTemplate, func(r mvutil.RequestParams, campaign *smodel.CampaignInfo) *smodel.VideoTemplateUrlItem {
// 			// 	id := int32(101)
// 			// 	url := "test_url"
// 			// 	return &smodel.VideoTemplateUrlItem{
// 			// 		ID:  id,
// 			// 		URL: url,
// 			// 	}
// 			// })
// 			// defer guard.Unpatch()

// 			guard := Patch(handleRVTemplate, func(ad *Ad, params *mvutil.Params, r *mvutil.RequestParams, template *smodel.VideoTemplateUrlItem) {
// 			})
// 			defer guard.Unpatch()

// 			r.Param.AdType = 94
// 			ad.CampaignType = 3
// 			ad.VideoURL = "test_ad_url"
// 			ad.Rv.TemplateUrl = "test_template_url"
// 			params.HTTPReq = int32(2)
// 			renderRewardVideoTemplate(&ad, &params, &r, &campaign)
// 		})
// 	})
// }

// func TestHandleRVTemplate(t *testing.T) {
// 	Convey("test handleRVTemplate", t, func() {
// 		var ad Ad
// 		var params mvutil.Params
// 		r := mvutil.RequestParams{}
// 		var template smodel.VideoTemplateUrlItem

// 		Convey("r.Param.ApiVersion >= mvconst.API_VERSION_1_2", func() {
// 			guard := Patch(extractor.GetDefRVTemplate, func() (*smodel.VideoTemplateUrlItem, bool) {
// 				id := int32(101)
// 				url := "test_url"
// 				urlZip := "test_zip"
// 				weight := int32(2)
// 				pURL := "test_p_url"
// 				pZip := "test_z_url"
// 				return &smodel.VideoTemplateUrlItem{
// 						ID:           id,
// 						URL:          url,
// 						URLZip:       urlZip,
// 						Weight:       weight,
// 						PausedURL:    pURL,
// 						PausedURLZip: pZip,
// 					},
// 					true
// 			})
// 			defer guard.Unpatch()

// 			ad.Rv.VideoTemplate = 3
// 			r.Param.ApiVersion = float64(1.3)
// 			tID := int32(2)
// 			tURL := "t_URL"
// 			template.ID = tID
// 			template.URL = tURL
// 			handleRVTemplate(&ad, &params, &r, &template)
// 		})

// 		Convey("ad.Rv.VideoTemplate == 4", func() {
// 			guard := Patch(extractor.GetDefRVTemplate, func() (*smodel.VideoTemplateUrlItem, bool) {
// 				id := int32(101)
// 				url := "test_url"
// 				urlZip := "test_zip"
// 				weight := int32(2)
// 				pURL := "test_p_url"
// 				pZip := "test_z_url"
// 				return &smodel.VideoTemplateUrlItem{
// 						ID:           id,
// 						URL:          url,
// 						URLZip:       urlZip,
// 						Weight:       weight,
// 						PausedURL:    pURL,
// 						PausedURLZip: pZip,
// 					},
// 					true
// 			})
// 			defer guard.Unpatch()

// 			ad.Rv.VideoTemplate = 4
// 			r.Param.ApiVersion = float64(1.3)
// 			tID := int32(2)
// 			tURL := "t_URL"
// 			template.ID = tID
// 			template.URL = tURL
// 			handleRVTemplate(&ad, &params, &r, &template)
// 		})
// 	})
// }

// func TestCheckResolution(t *testing.T) {
// 	var res bool
// 	Convey("符合尺寸", t, func() {
// 		res = checkResolution("32x18")
// 		So(res, ShouldBeTrue)
// 	})

// 	Convey("不符合尺寸", t, func() {
// 		res = checkResolution("180x18")
// 		So(res, ShouldBeFalse)
// 	})
// }

// func TestIsOriPortrait(t *testing.T) {
// 	var res bool
// 	Convey("0、1 返回 true", t, func() {
// 		res = isOriPortrait(0)
// 		So(res, ShouldBeTrue)

// 		res = isOriPortrait(1)
// 		So(res, ShouldBeTrue)
// 	})

// 	Convey("其他返回 false", t, func() {
// 		res = isOriPortrait(2)
// 		So(res, ShouldBeFalse)
// 	})
// }

// func TestIsUnitIOS(t *testing.T) {
// 	Convey("ios 返回 true", t, func() {
// 		platform := 2
// 		res := isUnitIOS(platform)
// 		So(res, ShouldBeTrue)
// 	})

// 	Convey("其他返回 false", t, func() {
// 		platform := 1
// 		res := isUnitIOS(platform)
// 		So(res, ShouldBeFalse)
// 	})
// }

// func TestHandleOrientation(t *testing.T) {
// 	var res int
// 	Convey("3、4 返回 1", t, func() {
// 		res = handleOrientation(3)
// 		So(res, ShouldEqual, 1)

// 		res = handleOrientation(4)
// 		So(res, ShouldEqual, 1)
// 	})

// 	Convey("其他返回 0", t, func() {
// 		res = handleOrientation(1)
// 		So(res, ShouldEqual, 0)
// 	})
// }

// func TestGetRVTemplate(t *testing.T) {
// 	Convey("test getRVTemplate", t, func() {
// 		r := mvutil.RequestParams{UnitInfo: &mvutil.UnitInfo{}}
// 		var campaign smodel.CampaignInfo
// 		var res *smodel.VideoTemplateUrlItem

// 		guard := Patch(extractor.GetEndcard, func() (smodel.EndCard, bool) {
// 			endCURLs := []*smodel.EndCardUrls{}
// 			vTURL := []*smodel.VideoTemplateUrlItem{}
// 			status := int32(1)
// 			endcardPro := int(3)
// 			rate := map[string]int{"1": 1, "2": 2}
// 			return smodel.EndCard{
// 				Urls:             endCURLs,
// 				Status:           status,
// 				Orientation:      status,
// 				VideoTemplateUrl: vTURL,
// 				EndcardProtocal:  endcardPro,
// 				EndcardRate:      rate,
// 			}, true
// 		})
// 		defer guard.Unpatch()

// 		Convey("offer + unit 维度", func() {
// 			guard = Patch(randRVTemplate, func(endcard *smodel.EndCard, r mvutil.RequestParams) *smodel.VideoTemplateUrlItem {
// 				id := int32(101)
// 				url := "test_url"
// 				urlZip := "test_zip"
// 				weight := int32(2)
// 				pURL := "test_p_url"
// 				pZip := "test_z_url"
// 				return &smodel.VideoTemplateUrlItem{
// 					ID:           id,
// 					URL:          url,
// 					URLZip:       urlZip,
// 					Weight:       weight,
// 					PausedURL:    pURL,
// 					PausedURLZip: pZip,
// 				}
// 			})
// 			defer guard.Unpatch()
// 			endC := map[string]*smodel.EndCard{
// 				"test_1": &smodel.EndCard{},
// 			}
// 			campaign.Endcard = endC
// 			r.Param.UnitID = int64(64100)

// 			res = getRVTemplate(r, &campaign)
// 			So(res.ID, ShouldResemble, int32(101))
// 			So(res.URL, ShouldResemble, "test_url")
// 			So(res.URLZip, ShouldEqual, "test_zip")
// 			So(res.Weight, ShouldEqual, int32(2))
// 			So(res.PausedURL, ShouldEqual, "test_p_url")
// 			So(res.PausedURLZip, ShouldEqual, "test_z_url")
// 		})

// 		Convey("offer 维度", func() {
// 			guard = Patch(randRVTemplate, func(endcard *smodel.EndCard, r mvutil.RequestParams) *smodel.VideoTemplateUrlItem {
// 				if endcard.Status == int32(111) {
// 					id := int32(101)
// 					url := "test_url"
// 					urlZip := "test_zip"
// 					weight := int32(2)
// 					pURL := "test_p_url"
// 					pZip := "test_z_url"
// 					return &smodel.VideoTemplateUrlItem{
// 						ID:           id,
// 						URL:          url,
// 						URLZip:       urlZip,
// 						Weight:       weight,
// 						PausedURL:    pURL,
// 						PausedURLZip: pZip,
// 					}
// 				}
// 				return &smodel.VideoTemplateUrlItem{}
// 			})
// 			defer guard.Unpatch()

// 			status := int32(111)
// 			statusTwo := int32(2)
// 			endC := map[string]*smodel.EndCard{
// 				"64100": &smodel.EndCard{
// 					Status: statusTwo,
// 				},
// 				"ALL": &smodel.EndCard{
// 					Status: status,
// 				},
// 			}
// 			campaign.Endcard = endC
// 			r.Param.UnitID = int64(64100)

// 			res = getRVTemplate(r, &campaign)
// 			So(res.ID, ShouldResemble, int32(101))
// 			So(res.URL, ShouldResemble, "test_url")
// 			So(res.URLZip, ShouldEqual, "test_zip")
// 			So(res.Weight, ShouldEqual, int32(2))
// 			So(res.PausedURL, ShouldEqual, "test_p_url")
// 			So(res.PausedURLZip, ShouldEqual, "test_z_url")
// 		})

// 		Convey("unit 维度", func() {
// 			guard = Patch(randRVTemplate, func(endcard *smodel.EndCard, r mvutil.RequestParams) *smodel.VideoTemplateUrlItem {
// 				id := int32(1011)
// 				url := "test_url_unit"
// 				urlZip := "test_zip_unit"
// 				weight := int32(22)
// 				pURL := "test_p_url_unit"
// 				pZip := "test_z_url_unit"
// 				return &smodel.VideoTemplateUrlItem{
// 					ID:           id,
// 					URL:          url,
// 					URLZip:       urlZip,
// 					Weight:       weight,
// 					PausedURL:    pURL,
// 					PausedURLZip: pZip,
// 				}
// 			})
// 			defer guard.Unpatch()

// 			status := int32(111)
// 			endC := smodel.EndCard{
// 				Status: status,
// 			}
// 			campaign.Endcard = nil
// 			r.UnitInfo.Endcard = &endC
// 			r.Param.UnitID = int64(64100)

// 			res = getRVTemplate(r, &campaign)
// 			So(res.ID, ShouldResemble, int32(1011))
// 			So(res.URL, ShouldResemble, "test_url_unit")
// 			So(res.URLZip, ShouldEqual, "test_zip_unit")
// 			So(res.Weight, ShouldEqual, int32(22))
// 			So(res.PausedURL, ShouldEqual, "test_p_url_unit")
// 			So(res.PausedURLZip, ShouldEqual, "test_z_url_unit")
// 		})

// 		Convey("config 维度", func() {
// 			guard = Patch(randRVTemplate, func(endcard *smodel.EndCard, r mvutil.RequestParams) *smodel.VideoTemplateUrlItem {
// 				id := int32(10111)
// 				url := "test_url_config"
// 				urlZip := "test_zip_config"
// 				weight := int32(222)
// 				pURL := "test_p_url_config"
// 				pZip := "test_z_url_config"
// 				return &smodel.VideoTemplateUrlItem{
// 					ID:           id,
// 					URL:          url,
// 					URLZip:       urlZip,
// 					Weight:       weight,
// 					PausedURL:    pURL,
// 					PausedURLZip: pZip,
// 				}
// 			})
// 			defer guard.Unpatch()

// 			campaign.Endcard = nil
// 			r.UnitInfo.Endcard = nil

// 			res = getRVTemplate(r, &campaign)
// 			So(res.ID, ShouldResemble, int32(10111))
// 			So(res.URL, ShouldResemble, "test_url_config")
// 			So(res.URLZip, ShouldEqual, "test_zip_config")
// 			So(res.Weight, ShouldEqual, int32(222))
// 			So(res.PausedURL, ShouldEqual, "test_p_url_config")
// 			So(res.PausedURLZip, ShouldEqual, "test_z_url_config")
// 		})
// 	})
// }

// func TestRandRVTemplate(t *testing.T) {
// 	var res *smodel.VideoTemplateUrlItem
// 	var endcard *smodel.EndCard
// 	r := mvutil.RequestParams{}
// 	exp := &smodel.VideoTemplateUrlItem{}

// 	status := int32(1)
// 	templateURL := []*smodel.VideoTemplateUrlItem{}
// 	endcard = &smodel.EndCard{
// 		Status:           status,
// 		VideoTemplateUrl: templateURL,
// 	}

// 	Convey("status 非 1，直接返回 endcard", t, func() {
// 		endcard.Status = 0
// 		res = randRVTemplate(endcard, r)
// 		So(res, ShouldResemble, exp)
// 	})

// 	Convey("VideoTemplateUrl 为空，直接返回 endcard", t, func() {
// 		endcard.Status = status
// 		endcard.VideoTemplateUrl = nil
// 		res = randRVTemplate(endcard, r)
// 		So(res, ShouldResemble, exp)
// 	})

// 	Convey("VideoTemplateUrl 内容为空，直接返回 endcard", t, func() {
// 		endcard.VideoTemplateUrl = templateURL
// 		res = randRVTemplate(endcard, r)
// 		So(res, ShouldResemble, exp)
// 	})

// 	Convey("VideoTemplateUrl 非空，直接返回 endcard", t, func() {
// 		weight1 := int32(5)
// 		templateURL = []*smodel.VideoTemplateUrlItem{
// 			&smodel.VideoTemplateUrlItem{Weight: weight1},
// 		}
// 		endcard.VideoTemplateUrl = templateURL
// 		res = randRVTemplate(endcard, r)
// 		exp := &smodel.VideoTemplateUrlItem{}
// 		exp.Weight = weight1
// 		So(res, ShouldResemble, exp)
// 	})
// }

// func TestRenderEndcard(t *testing.T) {
// 	Convey("test renderEndcard", t, func() {
// 		var ad Ad
// 		var params mvutil.Params
// 		var r mvutil.RequestParams
// 		var campaign smodel.CampaignInfo

// 		Convey("素材二期，优先使用 adserver 返回的 endcard", func() {
// 			params.EndcardCreativeID = int64(999)
// 			ad.EndcardUrl = "ad_endcard_url"
// 			guard := Patch(RenderEndcardUrl, func(params *mvutil.Params, endcardUrl string) string {
// 				return "test_render_endcard_url"
// 			})
// 			defer guard.Unpatch()

// 			renderEndcard(&ad, &params, &r, &campaign)
// 			So(params.Extendcard, ShouldEqual, "999")
// 			So(ad.EndcardUrl, ShouldEqual, "test_render_endcard_url")
// 		})

// 		Convey("IsReturnEndcard 返回 false，不做处理", func() {
// 			params.EndcardCreativeID = int64(0)
// 			guard := Patch(IsReturnEndcard, func(r *mvutil.RequestParams) bool {
// 				return false
// 			})
// 			defer guard.Unpatch()

// 			renderEndcard(&ad, &params, &r, &campaign)
// 		})

// 		Convey("campaign.Endcard == nil，不做处理", func() {
// 			params.EndcardCreativeID = int64(0)
// 			guard := Patch(IsReturnEndcard, func(r *mvutil.RequestParams) bool {
// 				return true
// 			})
// 			defer guard.Unpatch()

// 			campaign.Endcard = nil
// 			renderEndcard(&ad, &params, &r, &campaign)
// 		})

// 		// todo: 剩余逻辑...
// 	})
// }

// func TestGetEndcard(t *testing.T) {
// 	var res mvutil.EndcardItem
// 	var confs map[string]*smodel.EndCard
// 	r := mvutil.RequestParams{}

// 	ori := int32(1)
// 	endcardProto := int(8)
// 	endcardRate := map[string]int{"401": 411, "402": 412}
// 	ecID1 := int32(1)
// 	URL1 := "test_url_1"
// 	urls := []*smodel.EndCardUrls{
// 		&smodel.EndCardUrls{Id: ecID1, Url: URL1},
// 	}
// 	status := int32(1)
// 	conf := &smodel.EndCard{
// 		Status:          status,
// 		Urls:            urls,
// 		Orientation:     ori,
// 		EndcardProtocal: endcardProto,
// 		EndcardRate:     endcardRate,
// 	}
// 	confs = map[string]*smodel.EndCard{
// 		"101": conf,
// 	}

// 	Convey("endcar.Url 存在，返回 endcard", t, func() {
// 		r.Param.UnitID = int64(101)
// 		res = getEndcard(confs, &r)
// 		exp := mvutil.EndcardItem{
// 			Url:             "test_url_1",
// 			UrlV2:           "",
// 			Orientation:     1,
// 			ID:              1,
// 			EndcardProtocal: 8,
// 			EndcardRate:     map[string]int{"401": 411, "402": 412},
// 		}
// 		So(res, ShouldResemble, exp)
// 	})

// 	Convey("endcar.Url 不存在，返回 endcard", t, func() {
// 		r.Param.UnitID = int64(102)
// 		res = getEndcard(confs, &r)
// 		var nowEndRate map[string]int
// 		exp := mvutil.EndcardItem{
// 			Url:             "",
// 			UrlV2:           "",
// 			Orientation:     0,
// 			ID:              0,
// 			EndcardProtocal: 0,
// 			EndcardRate:     nowEndRate,
// 		}
// 		So(res, ShouldResemble, exp)
// 	})
// }

func TestRandEndcard(t *testing.T) {
	var res mvutil.EndcardItem
	var conf smodel.EndCard
	r := mvutil.RequestParams{}

	Convey("RandEndcard", t, func() {
		ori := int32(1)
		endcardProto := int(8)
		endcardRate := map[string]int{"401": 411, "402": 412}
		ecID2 := int32(2)
		URL2 := "test_url_2"
		urls := []*smodel.EndCardUrls{
			&smodel.EndCardUrls{Id: ecID2, Url: URL2},
		}

		status := int32(1)
		conf = smodel.EndCard{
			Status:          status,
			Urls:            urls,
			Orientation:     ori,
			EndcardProtocal: endcardProto,
			EndcardRate:     endcardRate,
		}
		res = RandEndcard(&conf, &r)
		exp := mvutil.EndcardItem{
			Url:             "test_url_2",
			UrlV2:           "",
			Orientation:     1,
			ID:              2,
			EndcardProtocal: 8,
			EndcardRate:     map[string]int{"401": 411, "402": 412},
		}
		So(res, ShouldResemble, exp)
	})
}

func TestCheckOrientation(t *testing.T) {
	var res bool

	Convey("任一参数为 0，返回 true", t, func() {
		res = checkOrientation(0, 10)
		So(res, ShouldBeTrue)

		res = checkOrientation(10, 0)
		So(res, ShouldBeTrue)
	})

	Convey("否则，判断两个数是否相等", t, func() {
		res = checkOrientation(1, 1)
		So(res, ShouldBeTrue)

		res = checkOrientation(1, 12)
		So(res, ShouldBeFalse)
	})
}

func TestRenderExtAttr(t *testing.T) {
	var res int

	Convey("renderExtAttr", t, func() {
		tag := int32(4)
		vba := int32(1)
		link := "test_link"
		campaign := smodel.CampaignInfo{
			Tag:             tag,
			VbaConnecting:   vba,
			VbaTrackingLink: link,
			CityCodeV2:      make(map[string][]int32),
		}
		res = renderExtAttr(&campaign)
		So(res, ShouldEqual, 3)
	})
}

func TestRenderCampaignInfo(t *testing.T) {
	Convey("test renderCampaignInfo", t, func() {
		var ad Ad
		r := mvutil.RequestParams{AppInfo: &smodel.AppInfo{}, UnitInfo: &smodel.UnitInfo{}}
		var campaign smodel.CampaignInfo
		var params mvutil.Params
		var corsairCampaign corsair_proto.Campaign

		guard := Patch(renderImpUrl, func(ad *Ad, r *mvutil.RequestParams, campaign *smodel.CampaignInfo, params *mvutil.Params) {
		})
		defer guard.Unpatch()

		guard = Patch(extractor.GetFCA_CAMIDS, func() ([]int64, bool) {
			return []int64{}, false
		})
		defer guard.Unpatch()

		guard = Patch(extractor.GetNEW_DEFAULT_FCA, func() (int, bool) {
			return 0, false
		})
		defer guard.Unpatch()

		guard = Patch(renderOfferwall, func(ad *Ad, r *mvutil.RequestParams, offerType int32, price float64) {
		})
		defer guard.Unpatch()

		guard = Patch(renderLinkType, func(ad *Ad, campaign *smodel.CampaignInfo, params *mvutil.Params) {
		})
		defer guard.Unpatch()

		// guard = Patch(renderRewardVideoTemplate, func(ad *Ad, params *mvutil.Params, r *mvutil.RequestParams, campaign *smodel.CampaignInfo) {
		// })
		// defer guard.Unpatch()

		// guard = Patch(renderEndcard, func(ad *Ad, params *mvutil.Params, r *mvutil.RequestParams, campaign *smodel.CampaignInfo) {
		// })
		// defer guard.Unpatch()

		// guard = Patch(renderPlayableNew, func(ad *Ad, campaign *smodel.CampaignInfo, params *mvutil.Params, r *mvutil.RequestParams) {
		// })
		// defer guard.Unpatch()

		//guard = Patch(renderJumpTypeInfo, func(ad *Ad, r *mvutil.RequestParams, campaign *smodel.CampaignInfo, params *mvutil.Params) {
		//})
		//defer guard.Unpatch()

		guard = Patch(renderDeleteDevID, func(ad *Ad, params *mvutil.Params, campaign *smodel.CampaignInfo) {
		})
		defer guard.Unpatch()

		guard = Patch(renderSettingInfo, func(ad *Ad, params *mvutil.Params, campaign *smodel.CampaignInfo) {
		})
		defer guard.Unpatch()

		guard = Patch(renderUa, func(ad *Ad, c *smodel.CampaignInfo, p *mvutil.Params) {
		})
		defer guard.Unpatch()

		guard = Patch(renderLoopback, func(ad *Ad, campaign *smodel.CampaignInfo) {
		})
		defer guard.Unpatch()

		guard = Patch(extractor.GetC_TOI, func() (int32, bool) {
			return int32(42), true
		})
		defer guard.Unpatch()

		guard = Patch(toNewCDN, func(ad *Ad, p *mvutil.Params) {
		})
		defer guard.Unpatch()

		guard = Patch(RenderCreativeUrls, func(ad *Ad, r *mvutil.RequestParams, params *mvutil.Params) {
		})
		defer guard.Unpatch()

		guard = Patch(renderPlct, func(ad *Ad, r *mvutil.Params, backendId int, dspId int64) {
		})
		defer guard.Unpatch()

		guard = Patch(extractor.GetADNET_SWITCHS, func() (map[string]int, bool) {
			return map[string]int{}, true
		})
		defer guard.Unpatch()

		guard = Patch(extractor.GetCampaignInfo, func(camId int64) (camInfo *smodel.CampaignInfo, ifFind bool) {
			return &smodel.CampaignInfo{}, true
		})
		defer guard.Unpatch()

		guard = Patch(NewDemandContext, func(params *mvutil.Params, r *mvutil.RequestParams) (context *demand.Context) {
			return new(demand.Context)
		})
		defer guard.Unpatch()

		guard = Patch(demand.RenderClickMode, func(ctx *demand.Context, dao chasm_extractor.IDAO) (string, string) {
			return "", ""
		})
		defer guard.Unpatch()

		Convey("空参数", func() {
			renderCampaignInfo(&ad, &r, &campaign, &params, corsairCampaign)
		})
	})
}

func TestCanUAAbtest(t *testing.T) {
	c := smodel.CampaignInfo{
		CampaignId: 1234567,
	}
	Convey("安卓可做的情况", t, func() {
		var params mvutil.Params
		camGuard := Patch(extractor.GetUA_AB_TEST_CAMPAIGN_CONFIG, func() ([]int64, bool) {
			return []int64{1234567}, true
		})
		defer camGuard.Unpatch()
		thirdGuard := Patch(extractor.GetUA_AB_TEST_THIRD_PARTY_CONFIG, func() ([]string, bool) {
			return []string{}, true
		})
		defer thirdGuard.Unpatch()
		guard := Patch(extractor.GetUA_AB_TEST_SDK_OS_CONFIG, func() (map[string]string, bool) {
			return map[string]string{"1": "8.7.4"}, true
		})
		defer guard.Unpatch()
		params = mvutil.Params{
			Platform:   1,
			SDKVersion: "mal_8.8.8",
			GAID:       "aaddddddcccccccbbbb",
		}
		res := canUAAbTest(&params, &c)
		So(res, ShouldBeTrue)
	})

	Convey("安卓不可做的情况", t, func() {
		var params mvutil.Params
		camGuard := Patch(extractor.GetUA_AB_TEST_CAMPAIGN_CONFIG, func() ([]int64, bool) {
			return []int64{1234567}, true
		})
		defer camGuard.Unpatch()
		thirdGuard := Patch(extractor.GetUA_AB_TEST_THIRD_PARTY_CONFIG, func() ([]string, bool) {
			return []string{}, true
		})
		defer thirdGuard.Unpatch()
		guard := Patch(extractor.GetUA_AB_TEST_SDK_OS_CONFIG, func() (map[string]string, bool) {
			return map[string]string{"1": "8.7.4"}, true
		})
		defer guard.Unpatch()
		params = mvutil.Params{
			Platform:   1,
			SDKVersion: "mal_6.8.8",
			GAID:       "aaddddddcccccccbbbb",
		}
		res := canUAAbTest(&params, &c)
		So(res, ShouldBeFalse)
	})

	Convey("ios可做的情况", t, func() {
		var params mvutil.Params
		camGuard := Patch(extractor.GetUA_AB_TEST_CAMPAIGN_CONFIG, func() ([]int64, bool) {
			return []int64{1234567}, true
		})
		defer camGuard.Unpatch()
		thirdGuard := Patch(extractor.GetUA_AB_TEST_THIRD_PARTY_CONFIG, func() ([]string, bool) {
			return []string{}, true
		})
		defer thirdGuard.Unpatch()
		guard := Patch(extractor.GetUA_AB_TEST_SDK_OS_CONFIG, func() (map[string]string, bool) {
			return map[string]string{"2": "3.3.5"}, true
		})
		defer guard.Unpatch()
		params = mvutil.Params{
			Platform:   2,
			SDKVersion: "mi_6.8.8",
			IDFA:       "aabbccddeeffgg",
		}
		res := canUAAbTest(&params, &c)
		So(res, ShouldBeTrue)
	})

	Convey("ios不可做的情况", t, func() {
		var params mvutil.Params
		camGuard := Patch(extractor.GetUA_AB_TEST_CAMPAIGN_CONFIG, func() ([]int64, bool) {
			return []int64{1234567}, true
		})
		defer camGuard.Unpatch()
		thirdGuard := Patch(extractor.GetUA_AB_TEST_THIRD_PARTY_CONFIG, func() ([]string, bool) {
			return []string{}, true
		})
		defer thirdGuard.Unpatch()
		guard := Patch(extractor.GetUA_AB_TEST_SDK_OS_CONFIG, func() (map[string]string, bool) {
			return map[string]string{"2": "3.3.5"}, true
		})
		defer guard.Unpatch()
		params = mvutil.Params{
			Platform:   2,
			SDKVersion: "mi_2.8.8",
			IDFA:       "aabbccddeeffgg",
		}
		res := canUAAbTest(&params, &c)
		So(res, ShouldBeFalse)
	})
}

// func TestRenderPlayable3(t *testing.T) {
// 	Convey("TestRenderPlayable3", t, func() {
// 		var ad *Ad
// 		//campaign := &smodel.CampaignInfo{}
// 		//var params *mvutil.Params
// 		r := mvutil.RequestParams{}
// 		var corsairCam corsair_proto.Campaign
// 		var param mvutil.Params

// 		Convey("corsairCampaign.Playable ==true, ExtPlayable 为1的情况", func() {
// 			playable := int32(1)
// 			videoEndType := int32(10)
// 			endcardUrl := "test_endcard_url"
// 			corsairCam = corsair_proto.Campaign{
// 				ExtPlayable:    &playable,
// 				VideoEndTypeAs: &videoEndType,
// 				EndcardUrl:     &endcardUrl,
// 				Playable:       true,
// 			}
// 			ad = &Ad{
// 				PlayableAdsWithoutVideo: 10,
// 				VideoURL:                "test_video_url",
// 			}

// 			renderPlayable3(ad, corsairCam, &param, &r)
// 			So(ad.EndcardUrl, ShouldEqual, "")
// 			So(param.Extplayable, ShouldEqual, 10110)
// 			So(ad.PlayableAdsWithoutVideo, ShouldEqual, 10)
// 			So(ad.VideoURL, ShouldEqual, "test_video_url")
// 		})

// 		Convey("corsairCampaign.Playable ==true, ExtPlayable 为2的情况", func() {
// 			playable := int32(2)
// 			videoEndType := int32(10)
// 			endcardUrl := "test_endcard_url"
// 			corsairCam = corsair_proto.Campaign{
// 				ExtPlayable:    &playable,
// 				VideoEndTypeAs: &videoEndType,
// 				EndcardUrl:     &endcardUrl,
// 				Playable:       true,
// 			}
// 			ad = &Ad{
// 				PlayableAdsWithoutVideo: 10,
// 				VideoURL:                "test_video_url",
// 			}

// 			renderPlayable3(ad, corsairCam, &param, &r)
// 			So(ad.EndcardUrl, ShouldEqual, "http://test_endcard_url")
// 			So(param.Extplayable, ShouldEqual, 10210)
// 			So(ad.PlayableAdsWithoutVideo, ShouldEqual, 10)
// 			So(ad.VideoURL, ShouldEqual, "test_video_url")
// 		})

// 		Convey("corsairCampaign.Playable ==true, ExtPlayable 为3的情况", func() {
// 			playable := int32(3)
// 			videoEndType := int32(10)
// 			endcardUrl := "test_endcard_url"
// 			corsairCam = corsair_proto.Campaign{
// 				ExtPlayable:    &playable,
// 				VideoEndTypeAs: &videoEndType,
// 				EndcardUrl:     &endcardUrl,
// 				Playable:       true,
// 			}
// 			ad = &Ad{
// 				PlayableAdsWithoutVideo: 10,
// 				VideoURL:                "test_video_url",
// 			}

// 			renderPlayable3(ad, corsairCam, &param, &r)
// 			So(ad.EndcardUrl, ShouldEqual, "http://test_endcard_url")
// 			So(param.Extplayable, ShouldEqual, 10310)
// 			So(ad.PlayableAdsWithoutVideo, ShouldEqual, 2)
// 			So(ad.VideoURL, ShouldEqual, "")
// 		})

// 		Convey("corsairCampaign.Playable ==false, ExtPlayable 为1的情况", func() {
// 			playable := int32(1)
// 			videoEndType := int32(10)
// 			endcardUrl := "test_endcard_url"
// 			corsairCam = corsair_proto.Campaign{
// 				ExtPlayable:    &playable,
// 				VideoEndTypeAs: &videoEndType,
// 				EndcardUrl:     &endcardUrl,
// 				Playable:       false,
// 			}
// 			ad = &Ad{
// 				PlayableAdsWithoutVideo: 10,
// 				VideoURL:                "test_video_url",
// 			}

// 			renderPlayable3(ad, corsairCam, &param, &r)
// 			So(ad.EndcardUrl, ShouldEqual, "")
// 			So(param.Extplayable, ShouldEqual, 20110)
// 			So(ad.PlayableAdsWithoutVideo, ShouldEqual, 10)
// 			So(ad.VideoURL, ShouldEqual, "test_video_url")

// 		})

// 		Convey("corsairCampaign.Playable ==false, ExtPlayable 为0的情况", func() {
// 			playable := int32(0)
// 			videoEndType := int32(10)
// 			endcardUrl := "test_endcard_url"
// 			corsairCam = corsair_proto.Campaign{
// 				ExtPlayable:    &playable,
// 				VideoEndTypeAs: &videoEndType,
// 				EndcardUrl:     &endcardUrl,
// 				Playable:       false,
// 			}
// 			ad = &Ad{
// 				PlayableAdsWithoutVideo: 10,
// 				VideoURL:                "test_video_url",
// 			}

// 			renderPlayable3(ad, corsairCam, &param, &r)
// 			So(ad.EndcardUrl, ShouldEqual, "")
// 			So(param.Extplayable, ShouldEqual, 0)
// 			So(ad.PlayableAdsWithoutVideo, ShouldEqual, 10)
// 			So(ad.VideoURL, ShouldEqual, "test_video_url")

// 		})
// 	})

// }

func TestRenderEndcardProperty(t *testing.T) {
	Convey("TestRenderEndcardProperty", t, func() {
		var ad *Ad
		campaign := &smodel.CampaignInfo{}
		a := 1
		r := &mvutil.RequestParams{
			UnitInfo: &smodel.UnitInfo{
				Unit: smodel.Unit{
					Alac:   &a,
					Alecfc: &a,
					Mof:    &a,
				},
			},
		}
		guard := Patch(extractor.GetALAC_PRISON_CONFIG, func() (map[string][]int, bool) {
			return map[string][]int{}, true
		})
		defer guard.Unpatch()
		guard = Patch(extractor.GetABTEST_FIELDS, func() (mvutil.ABTEST_FIELDS, bool) {
			return mvutil.ABTEST_FIELDS{}, true
		})
		defer guard.Unpatch()
		guard = Patch(extractor.GetABTEST_CONFS, func() (map[string][]mvutil.ABTEST_CONF, bool) {
			return map[string][]mvutil.ABTEST_CONF{}, true
		})
		defer guard.Unpatch()
		guard = Patch(extractor.GetADNET_SWITCHS, func() (map[string]int, bool) {
			return map[string]int{}, true
		})
		defer guard.Unpatch()
		guard = Patch(extractor.GetADNET_CONF_LIST, func() map[string][]int64 {
			return map[string][]int64{}
		})
		defer guard.Unpatch()

		guard = Patch(extractor.GetTPL_CREATIVE_DOMAIN_CONF, func() map[string][]*mvutil.TemplateCreativeDomainMap {
			return map[string][]*mvutil.TemplateCreativeDomainMap{}
		})
		defer guard.Unpatch()

		Convey("endcard_url和endscreen_url均为空值则不处理", func() {
			ad = &Ad{
				EndcardUrl: "",
			}
			params := &mvutil.Params{
				IAPlayableUrl: "",
			}
			renderEndcardProperty(ad, campaign, params, r)
			So(ad.EndcardUrl, ShouldEqual, "")
			So(params.IAPlayableUrl, ShouldEqual, "")
		})
		Convey("ad_type不为iv，rv，interactive", func() {
			ad = &Ad{
				EndcardUrl: "test_endcard_url",
			}
			params := &mvutil.Params{
				IAPlayableUrl: "test_endcard_url",
				AdType:        42,
			}
			renderEndcardProperty(ad, campaign, params, r)
			So(ad.EndcardUrl, ShouldEqual, "test_endcard_url")
			So(params.IAPlayableUrl, ShouldEqual, "test_endcard_url")
		})
		Convey("rv 情况下，均不下发alac，alecfc，mof", func() {
			ad = &Ad{
				EndcardUrl: "test_endcard_url",
			}
			params := &mvutil.Params{
				IAPlayableUrl: "test_endcard_url",
				AdType:        94,
			}
			mof := 2
			campaign.Mof = mof
			renderEndcardProperty(ad, campaign, params, r)
			So(ad.EndcardUrl, ShouldEqual, "test_endcard_url&clsdly=3")
			So(params.IAPlayableUrl, ShouldEqual, "test_endcard_url")
		})
		Convey("rv 情况下，下发alac", func() {
			ad = &Ad{
				EndcardUrl: "test_endcard_url",
			}
			params := &mvutil.Params{
				IAPlayableUrl: "test_endcard_url",
				AdType:        94,
			}
			mof := 2
			campaign.Mof = mof
			alacRate := 100
			campaign.AlacRate = alacRate
			renderEndcardProperty(ad, campaign, params, r)
			So(ad.EndcardUrl, ShouldEqual, "test_endcard_url&alac=1")
			So(params.IAPlayableUrl, ShouldEqual, "test_endcard_url")
			So(params.ExtDataInit.Alac, ShouldEqual, 1)
		})
		Convey("rv 情况下，下发alecfc,", func() {
			ad = &Ad{
				EndcardUrl: "test_endcard_url",
			}
			params := &mvutil.Params{
				IAPlayableUrl: "test_endcard_url",
				AdType:        94,
			}
			mof := 2
			campaign.Mof = mof
			alecfcRate := 100
			campaign.AlecfcRate = alecfcRate
			renderEndcardProperty(ad, campaign, params, r)
			So(ad.EndcardUrl, ShouldEqual, "test_endcard_url&clsdly=3&alecfc=1")
			So(params.IAPlayableUrl, ShouldEqual, "test_endcard_url")
			So(params.ExtDataInit.Alecfc, ShouldEqual, 1)
		})
		Convey("rv 情况下，下发mof,", func() {
			ad = &Ad{
				EndcardUrl: "test_endcard_url",
			}
			params := &mvutil.Params{
				IAPlayableUrl: "test_endcard_url",
				AdType:        94,
			}
			mof := 1
			campaign.Mof = mof
			renderEndcardProperty(ad, campaign, params, r)
			So(ad.EndcardUrl, ShouldEqual, "test_endcard_url&clsdly=3&mof=1&ec_id=&rv_tid=0&tplgp=0&v_fmd5=&i_fmd5=&mcc=&mnc=")
			So(params.IAPlayableUrl, ShouldEqual, "test_endcard_url")
			So(params.ExtDataInit.Mof, ShouldEqual, 1)
		})
		Convey("interactive 情况下，均不下发alac，alecfc，mof", func() {
			ad = &Ad{
				EndcardUrl: "test_endcard_url",
			}
			params := &mvutil.Params{
				IAPlayableUrl: "test_endcard_url",
				AdType:        288,
			}
			mof := 2
			campaign.Mof = mof
			renderEndcardProperty(ad, campaign, params, r)
			So(ad.EndcardUrl, ShouldEqual, "test_endcard_url")
			So(params.IAPlayableUrl, ShouldEqual, "test_endcard_url&clsdly=3")
		})
		Convey("interactive 情况下，下发alac", func() {
			ad = &Ad{
				EndcardUrl: "test_endcard_url",
			}
			params := &mvutil.Params{
				IAPlayableUrl: "test_endcard_url",
				AdType:        288,
			}
			mof := 2
			campaign.Mof = mof
			alacRate := 100
			campaign.AlacRate = alacRate
			renderEndcardProperty(ad, campaign, params, r)
			So(ad.EndcardUrl, ShouldEqual, "test_endcard_url")
			So(params.IAPlayableUrl, ShouldEqual, "test_endcard_url&alac=1")
			So(params.ExtDataInit.Alac, ShouldEqual, 1)
		})
		Convey("interactive 情况下，下发alecfc,", func() {
			ad = &Ad{
				EndcardUrl: "test_endcard_url",
			}
			params := &mvutil.Params{
				IAPlayableUrl: "test_endcard_url",
				AdType:        288,
			}
			mof := 2
			campaign.Mof = mof
			alecfcRate := 100
			campaign.AlecfcRate = alecfcRate
			renderEndcardProperty(ad, campaign, params, r)
			So(ad.EndcardUrl, ShouldEqual, "test_endcard_url")
			So(params.IAPlayableUrl, ShouldEqual, "test_endcard_url&clsdly=3&alecfc=1")
			So(params.ExtDataInit.Alecfc, ShouldEqual, 1)
		})
		Convey("interactive 情况下，下发mof,", func() {
			ad = &Ad{
				EndcardUrl: "test_endcard_url",
			}
			params := &mvutil.Params{
				IAPlayableUrl: "test_endcard_url",
				AdType:        288,
			}
			mof := 1
			campaign.Mof = mof
			renderEndcardProperty(ad, campaign, params, r)
			So(ad.EndcardUrl, ShouldEqual, "test_endcard_url")
			So(params.IAPlayableUrl, ShouldEqual, "test_endcard_url&clsdly=3&mof=1&ec_id=&rv_tid=0&tplgp=0&v_fmd5=&i_fmd5=&mcc=&mnc=")
			So(params.ExtDataInit.Mof, ShouldEqual, 1)
		})
	})
}

func TestJudgeByUnitAndOffer(t *testing.T) {
	Convey("TestRenderEndcardProperty", t, func() {
		Convey("默认配置", func() {
			var unitSwitch int
			var rate int
			res := judgeByUnitAndOffer(&unitSwitch, rate)
			So(res, ShouldEqual, false)
		})
		Convey("unit关闭", func() {
			unitSwitch := 2
			var rate int
			res := judgeByUnitAndOffer(&unitSwitch, rate)
			So(res, ShouldEqual, false)
		})
		Convey("unit开启，offer维度没配", func() {
			unitSwitch := 1
			var rate int
			res := judgeByUnitAndOffer(&unitSwitch, rate)
			So(res, ShouldEqual, false)
		})
		Convey("unit开启，offer维度80%", func() {
			unitSwitch := 1
			rate := 80
			guardRand := Patch(rand.Intn, func(n int) int {
				return 1
			})
			defer guardRand.Unpatch()
			res := judgeByUnitAndOffer(&unitSwitch, rate)
			So(res, ShouldEqual, true)
		})
		Convey("unit关闭，offer维度100%", func() {
			unitSwitch := 2
			rate := 100
			res := judgeByUnitAndOffer(&unitSwitch, rate)
			So(res, ShouldEqual, false)
		})
		Convey("unit开启，offer维度0%", func() {
			unitSwitch := 1
			rate := 0
			res := judgeByUnitAndOffer(&unitSwitch, rate)
			So(res, ShouldEqual, false)
		})
	})
}

func Test_fixAndroidJumpType(t *testing.T) {
	type args struct {
		jumpType string
		r        *mvutil.RequestParams
		campaign *smodel.CampaignInfo
		params   *mvutil.Params
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "非android",
			args: args{
				r: &mvutil.RequestParams{
					Param: mvutil.Params{
						Platform: mvconst.PlatformIOS,
					},
				},
				jumpType: mvconst.JUMP_TYPE_CLIENT_DO_ALL_PING_SERVER,
			},
			want: mvconst.JUMP_TYPE_CLIENT_DO_ALL_PING_SERVER,
		},
		{
			name: "非sdk",
			args: args{
				r: &mvutil.RequestParams{
					Param: mvutil.Params{
						Platform: mvconst.PlatformAndroid,
					},
				},
				jumpType: mvconst.JUMP_TYPE_CLIENT_DO_ALL_PING_SERVER,
			},
			want: mvconst.JUMP_TYPE_CLIENT_DO_ALL_PING_SERVER,
		},
		{
			name: "android sdk 非 6/12",
			args: args{
				r: &mvutil.RequestParams{
					Param: mvutil.Params{
						Platform:    mvconst.PlatformAndroid,
						RequestType: mvconst.REQUEST_TYPE_OPENAPI_V3,
					},
				},
				jumpType: mvconst.JUMP_TYPE_CLIENT_SEND_DEVID,
			},
			want: mvconst.JUMP_TYPE_CLIENT_SEND_DEVID,
		},
		{
			name: "android sdk 6/12 device 为空",
			args: args{
				r: &mvutil.RequestParams{
					Param: mvutil.Params{
						Platform:    mvconst.PlatformAndroid,
						RequestType: mvconst.REQUEST_TYPE_OPENAPI_V3,
					},
				},
				params:   &mvutil.Params{},
				jumpType: mvconst.JUMP_TYPE_CLIENT_DO_ALL_PING_SERVER,
			},
			want: mvconst.JUMP_TYPE_NORMAL,
		},
		{
			name: "android sdk 6 -> 5",
			args: args{
				r: &mvutil.RequestParams{
					Param: mvutil.Params{
						Platform:    mvconst.PlatformAndroid,
						RequestType: mvconst.REQUEST_TYPE_OPENAPI_V3,
					},
				},
				params: &mvutil.Params{
					GAID: "xx-xx-xx-xx",
				},
				campaign: &smodel.CampaignInfo{},
				jumpType: mvconst.JUMP_TYPE_CLIENT_DO_ALL,
			},
			want: mvconst.JUMP_TYPE_CLIENT_SEND_DEVID,
		},
		{
			name: "android sdk 12 -> 5",
			args: args{
				r: &mvutil.RequestParams{
					Param: mvutil.Params{
						Platform:    mvconst.PlatformAndroid,
						RequestType: mvconst.REQUEST_TYPE_OPENAPI_V3,
					},
				},
				params: &mvutil.Params{
					GAID: "xx-xx-xx-xx",
				},
				campaign: &smodel.CampaignInfo{},
				jumpType: mvconst.JUMP_TYPE_CLIENT_DO_ALL,
			},
			want: mvconst.JUMP_TYPE_CLIENT_SEND_DEVID,
		},
		// TODO CASE 没完善
	}
	for _, tt := range tests {
		guard := Patch(extractor.GetADNET_SWITCHS, func() (map[string]int, bool) {
			return map[string]int{"icGaid": 1}, false
		})
		defer guard.Unpatch()
		t.Run(tt.name, func(t *testing.T) {
			if got := fixAndroidJumpType(tt.args.jumpType, tt.args.r, tt.args.campaign, tt.args.params); got != tt.want {
				t.Errorf("fixAndroidJumpType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCanUseAfWhiteList(t *testing.T) {
	Convey("TestCanUseAfWhiteList", t, func() {
		var ad Ad
		var params mvutil.Params
		campaign := &smodel.CampaignInfo{}

		Convey("非903单子", func() {
			advId := int32(902)
			campaign.AdvertiserId = advId
			canUseAfWhiteList(&params, campaign, &ad)
			So(params.ExtDataInit.ClickInServer, ShouldEqual, 0)
			So(ad.CUA, ShouldEqual, 0)
			So(params.Extra10, ShouldEqual, "")
		})
		params.RequestPath = mvconst.PATHOpenApiV3
		params.Scenario = mvconst.SCENARIO_OPENAPI
		Convey("有deviceid的情况", func() {
			advId := int32(903)
			campaign.AdvertiserId = advId
			params.Platform = mvconst.PlatformIOS
			params.IDFA = "2F5BBAA9-07CE-40F7-8747-FD343974F62E"
			Convey("被offer黑名单过滤", func() {
				campaign.CampaignId = 1234567
				guard := Patch(extractor.GetCLICK_IN_SERVER_CONF_NEW, func() mvutil.ClickInServerConf {
					return mvutil.ClickInServerConf{
						OfferBList: []int64{
							1234567,
						},
					}
				})
				defer guard.Unpatch()
				canUseAfWhiteList(&params, campaign, &ad)
				So(params.ExtDataInit.ClickInServer, ShouldEqual, 0)
				So(ad.CUA, ShouldEqual, 0)
				So(params.Extra10, ShouldEqual, "")
			})

			Convey("无黑名单过滤，无此三方配置", func() {
				campaign.ThirdParty = "tengjin"
				guard := Patch(extractor.GetCLICK_IN_SERVER_CONF_NEW, func() mvutil.ClickInServerConf {
					return mvutil.ClickInServerConf{}
				})
				defer guard.Unpatch()
				guard = Patch(extractor.GetSYSTEM_AREA, func() string {
					return "sg"
				})
				defer guard.Unpatch()
				canUseAfWhiteList(&params, campaign, &ad)
				So(params.ExtDataInit.ClickInServer, ShouldEqual, 0)
				So(ad.CUA, ShouldEqual, 0)
				So(params.Extra10, ShouldEqual, "")
			})

			Convey("无黑名单过滤，有三方配置，命中切量", func() {
				campaign.ThirdParty = "Appsflyer"
				guard := Patch(extractor.GetCLICK_IN_SERVER_CONF_NEW, func() mvutil.ClickInServerConf {
					return mvutil.ClickInServerConf{}
				})
				defer guard.Unpatch()
				guardRand := Patch(rand.Intn, func(n int) int {
					return 1
				})
				defer guardRand.Unpatch()

				guard = Patch(getClickInServerRate, func(params *mvutil.Params, campaign *smodel.CampaignInfo, isDeviceEmpty bool) int {
					return 50
				})
				defer guard.Unpatch()
				guard = Patch(extractor.GetSYSTEM_AREA, func() string {
					return "sg"
				})
				defer guard.Unpatch()

				canUseAfWhiteList(&params, campaign, &ad)
				So(params.ExtDataInit.ClickInServer, ShouldEqual, 3)
				So(ad.CUA, ShouldEqual, 1)
				So(params.Extra10, ShouldEqual, "13")
			})

			Convey("无黑名单过滤，有三方配置，无命中切量", func() {
				campaign.ThirdParty = "Appsflyer"
				guard := Patch(extractor.GetCLICK_IN_SERVER_CONF_NEW, func() mvutil.ClickInServerConf {
					return mvutil.ClickInServerConf{}
				})
				defer guard.Unpatch()
				guardRand := Patch(rand.Intn, func(n int) int {
					return 100
				})
				defer guardRand.Unpatch()
				guard = Patch(extractor.GetSYSTEM_AREA, func() string {
					return "sg"
				})
				defer guard.Unpatch()

				guard = Patch(getClickInServerRate, func(params *mvutil.Params, campaign *smodel.CampaignInfo, isDeviceEmpty bool) int {
					return 50
				})
				defer guard.Unpatch()

				canUseAfWhiteList(&params, campaign, &ad)
				So(params.ExtDataInit.ClickInServer, ShouldEqual, 4)
				So(ad.CUA, ShouldEqual, 0)
				So(params.Extra10, ShouldEqual, "")
			})
		})

		Convey("无deviceid的情况", func() {
			advId := int32(903)
			campaign.AdvertiserId = advId
			params.Platform = mvconst.PlatformIOS
			params.IDFA = ""
			Convey("被clickmode黑名单过滤", func() {
				params.Extra10 = "11"
				guard := Patch(extractor.GetCLICK_IN_SERVER_CONF_NEW, func() mvutil.ClickInServerConf {
					return mvutil.ClickInServerConf{
						ClickModeBList: []string{
							"11",
						},
					}
				})
				defer guard.Unpatch()
				canUseAfWhiteList(&params, campaign, &ad)
				So(params.ExtDataInit.ClickInServer, ShouldEqual, 0)
				So(ad.CUA, ShouldEqual, 0)
				So(params.Extra10, ShouldEqual, "11")
			})

			Convey("无黑名单过滤，无此三方配置", func() {
				campaign.ThirdParty = "tengjin"
				guard := Patch(extractor.GetCLICK_IN_SERVER_CONF_NEW, func() mvutil.ClickInServerConf {
					return mvutil.ClickInServerConf{}
				})
				defer guard.Unpatch()
				guard = Patch(extractor.GetSYSTEM_AREA, func() string {
					return "sg"
				})
				defer guard.Unpatch()
				canUseAfWhiteList(&params, campaign, &ad)
				So(params.ExtDataInit.ClickInServer, ShouldEqual, 0)
				So(ad.CUA, ShouldEqual, 0)
				So(params.Extra10, ShouldEqual, "")
			})

			Convey("无黑名单过滤，有三方配置，命中切量", func() {
				campaign.ThirdParty = "Appsflyer"
				guard := Patch(extractor.GetCLICK_IN_SERVER_CONF_NEW, func() mvutil.ClickInServerConf {
					return mvutil.ClickInServerConf{}
				})
				defer guard.Unpatch()
				guardRand := Patch(rand.Intn, func(n int) int {
					return 1
				})
				defer guardRand.Unpatch()

				guard = Patch(getClickInServerRate, func(params *mvutil.Params, campaign *smodel.CampaignInfo, isDeviceEmpty bool) int {
					return 50
				})
				defer guard.Unpatch()
				guard = Patch(extractor.GetSYSTEM_AREA, func() string {
					return "sg"
				})
				defer guard.Unpatch()

				canUseAfWhiteList(&params, campaign, &ad)
				So(params.ExtDataInit.ClickInServer, ShouldEqual, 1)
				So(ad.CUA, ShouldEqual, 1)
				So(params.Extra10, ShouldEqual, "13")
			})

			Convey("无黑名单过滤，有三方配置，无命中切量", func() {
				campaign.ThirdParty = "Appsflyer"
				guard := Patch(extractor.GetCLICK_IN_SERVER_CONF_NEW, func() mvutil.ClickInServerConf {
					return mvutil.ClickInServerConf{}
				})
				defer guard.Unpatch()
				guardRand := Patch(rand.Intn, func(n int) int {
					return 100
				})
				defer guardRand.Unpatch()
				guard = Patch(extractor.GetSYSTEM_AREA, func() string {
					return "sg"
				})
				defer guard.Unpatch()

				guard = Patch(getClickInServerRate, func(params *mvutil.Params, campaign *smodel.CampaignInfo, isDeviceEmpty bool) int {
					return 50
				})
				defer guard.Unpatch()

				canUseAfWhiteList(&params, campaign, &ad)
				So(params.ExtDataInit.ClickInServer, ShouldEqual, 2)
				So(ad.CUA, ShouldEqual, 0)
				So(params.Extra10, ShouldEqual, "")
			})
		})
	})
}

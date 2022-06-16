package output

import (
	"testing"

	"gitlab.mobvista.com/ADN/adnet/internal/mvconst"

	. "github.com/bouk/monkey"
	mlogger "github.com/mae-pax/logger"
	. "github.com/smartystreets/goconvey/convey"
	"gitlab.mobvista.com/ADN/adnet/internal/corsair_proto"
	"gitlab.mobvista.com/ADN/adnet/internal/extractor"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	"gitlab.mobvista.com/ADN/adnet/internal/protobuf"
	"gitlab.mobvista.com/ADN/adnet/internal/redis"
	"gitlab.mobvista.com/ADN/adnet/internal/watcher"
	smodel "gitlab.mobvista.com/ADN/structs/model"
	"gitlab.mobvista.com/adserver/recommend_protocols/go/ad_server"
)

func setupTest() {
	extractor.NewMDbLoaderRegistry()
}

func TestRenderOutput(t *testing.T) {
	setupTest()
	Convey("test RenderOutput", t, func() {
		r := mvutil.RequestParams{AppInfo: &smodel.AppInfo{}, UnitInfo: &smodel.UnitInfo{}, Param: mvutil.Params{ReduceFillConfig: &smodel.ConfigAlgorithmFillRate{}}}
		var res corsair_proto.QueryResult_

		guard := Patch(RenderUnitEndcard, func(r *mvutil.RequestParams) {
		})
		defer guard.Unpatch()

		guard = Patch(extractor.GetDOMAIN_TRACK, func() string {
			return "test_get_domain_track"
		})
		defer guard.Unpatch()

		guard = Patch(extractor.GetDOMAIN_TRACKS, func() map[string]mvutil.IRate {
			return nil
		})
		defer guard.Unpatch()

		guard = Patch(extractor.GetSYSTEM, func() string {
			return mvconst.SERVER_SYSTEM_M
		})
		defer guard.Unpatch()

		// guard = Patch(extractor.GetSYSTEM_AREA, func() string {
		// 	return "test_get_sys_area"
		// })
		// defer guard.Unpatch()

		guard = Patch(getCloudExtra, func() string {
			return "aws"
		})
		defer guard.Unpatch()

		guard = Patch(renderMobvistaCampaigns, func(r *mvutil.RequestParams, ads *corsair_proto.BackendAds) []Ad {
			return []Ad{
				Ad{
					CampaignID: int64(111),
					OfferID:    int(1111),
					AppName:    "test_app_name",
					AppDesc:    "test_app_desc",
				},
			}
		})
		defer guard.Unpatch()

		guard = Patch(renderThirdPartyCampaigns, func(r *mvutil.RequestParams, ads *corsair_proto.BackendAds) []Ad {
			return []Ad{
				Ad{
					CampaignID: int64(112),
					OfferID:    int(1112),
					AppName:    "test_app_name_2",
					AppDesc:    "test_app_desc_2",
				},
			}
		})
		defer guard.Unpatch()

		guard = Patch(CreateOnlyImpressionUrl, func(params mvutil.Params, r *mvutil.RequestParams) string {
			return "test_impression_url"
		})
		defer guard.Unpatch()

		// guard = Patch(IsGoTrack, func(r *mvutil.RequestParams) {
		// 	r.Param.Extabtest1 = 0
		// })

		guard = Patch(Is3SCNLine, func(r *mvutil.RequestParams) {
			r.Param.CNDomainTest = 0
		})

		guard = Patch(mvutil.GetGoTkClickID, func() string {
			return "test_go_click_id"
		})

		guard = Patch(extractor.GetADNET_SWITCHS, func() (map[string]int, bool) {
			return map[string]int{}, true
		})
		defer guard.Unpatch()

		guard = Patch(extractor.GetONLY_REQUEST_THIRD_DSP_SWITCH, func() bool {
			return false
		})
		defer guard.Unpatch()

		guard = Patch(extractor.GetREPLACE_TRACKING_DOMAIN_CONF, func() *mvutil.ReplaceTrackingDomainConf {
			return &mvutil.ReplaceTrackingDomainConf{}
		})
		defer guard.Unpatch()

		Convey("正常数据", func() {
			c := mlogger.NewFromYaml("./testdata/watch_log.yaml")
			logger := c.InitLogger("time", "", true, true)
			watcher.Init(logger)
			res.AdsByBackend = []*corsair_proto.BackendAds{
				&corsair_proto.BackendAds{
					BackendId:  int32(101),
					RequestKey: "test_r_key_1",
				},
				&corsair_proto.BackendAds{
					BackendId:  int32(102),
					RequestKey: "test_r_key_2",
				},
			}
			mvutil.Config = &mvutil.AdnetConfig{
				AreaConfig: &mvutil.AreaConfig{
					HttpConfig: mvutil.HttpConfig{
						RegionName: "singapore",
					},
				},
			}
			result, errs := RenderOutput(&r, &res)
			So((*result).Status, ShouldEqual, 1)
			So((*result).Msg, ShouldEqual, "success")
			So((*result).Data.Ads[0].CampaignID, ShouldEqual, int64(112))
			So((*result).Data.Ads[0].OfferID, ShouldEqual, int64(1112))
			So((*result).Data.Ads[0].AppDesc, ShouldEqual, "test_app_desc_2")
			So((*result).Data.Ads[0].AppName, ShouldEqual, "test_app_name_2")
			//So(r.Param.ExtcdnType, ShouldEqual, "test_get_sys_area")
			So(errs, ShouldBeNil)
		})
	})
}

func TestRenderThirdPartyCampaigns(t *testing.T) {
	Convey("test renderThirdPartyCampaigns", t, func() {
		r := mvutil.RequestParams{}
		var ads corsair_proto.BackendAds
		var res []Ad

		guard := Patch(RenderThirdPartyCampaign, func(r *mvutil.RequestParams, corsairCampaign corsair_proto.Campaign, backnedID int32, i int) Ad {
			return Ad{
				CampaignID:  int64(100001),
				OfferID:     int(100002),
				AppName:     *corsairCampaign.AppName,
				AppDesc:     "test_app_desc",
				PackageName: "test_package",
				IconURL:     "test_icon_url",
			}
		})
		defer guard.Unpatch()

		guard = Patch(redis.LocalRedisAlgoHExists, func(key, field string) (bool, error) {
			return false, nil
		})
		defer guard.Unpatch()

		Convey("ads.CampaignList 为空", func() {
			res = renderThirdPartyCampaigns(&r, &ads)
			So(res, ShouldBeNil)
		})

		Convey("ads.CampaignList 不为空", func() {
			appName1 := "app_name_1"
			appName2 := "app_name_2"
			ads.CampaignList = []*corsair_proto.Campaign{
				&corsair_proto.Campaign{
					CampaignId: "12345",
					AdSource:   ad_server.ADSource(1),
					AppName:    &appName1,
				},
				&corsair_proto.Campaign{
					CampaignId: "12346",
					AdSource:   ad_server.ADSource(2),
					AppName:    &appName2,
				},
				nil,
			}
			res = renderThirdPartyCampaigns(&r, &ads)
			So(r.Param.ThirdCidList, ShouldResemble, []string{"12345", "12346"})
			So(res, ShouldResemble, []Ad{
				Ad{
					CampaignID:  int64(100001),
					OfferID:     int(100002),
					AppName:     "app_name_1",
					AppDesc:     "test_app_desc",
					PackageName: "test_package",
					IconURL:     "test_icon_url",
				},
				Ad{
					CampaignID:  int64(100001),
					OfferID:     int(100002),
					AppName:     "app_name_2",
					AppDesc:     "test_app_desc",
					PackageName: "test_package",
					IconURL:     "test_icon_url",
				},
			})
		})
	})
}

func TestRenderMobvistaCampaigns(t *testing.T) {
	Convey("test renderMobvistaCampaigns", t, func() {
		r := mvutil.RequestParams{}
		var ads corsair_proto.BackendAds
		var res []Ad

		guard := Patch(renderAllParams, func(r *mvutil.RequestParams, ads corsair_proto.BackendAds) {
		})
		defer guard.Unpatch()

		guard = Patch(GetCampaignsInfo, func(campaignIds []int64) (map[int64]*smodel.CampaignInfo, error) {
			//apkURL := "apk_URL"
			advID := int32(42)
			return map[int64]*smodel.CampaignInfo{
				int64(10001): &smodel.CampaignInfo{
					CampaignId: int64(10001),
					//ApkUrl:       &apkURL,
					AdvertiserId: advID,
				},
			}, nil
		})
		defer guard.Unpatch()

		guard = Patch(renderAdchoice4Ad, func(ad *Ad, campaign *smodel.CampaignInfo) {
			// not render
			return
		})
		defer guard.Unpatch()

		guard = Patch(RenderCampaignWithCreative, func(r *mvutil.RequestParams, corsairCampaign corsair_proto.Campaign, campaign *smodel.CampaignInfo, cContent map[ad_server.CreativeType]*protobuf.Creative) (Ad, error) {
			return Ad{
				CampaignID: int64(201),
				OfferID:    int(202),
				AppName:    "test_app_name",
				AppDesc:    "app_desc",
			}, nil
		})
		defer guard.Unpatch()

		guard = Patch(watcher.AddWatchValue, func(key string, val float64) {})
		defer guard.Unpatch()

		guard = Patch(watcher.AddAvgWatchValue, func(key string, val float64) {})
		defer guard.Unpatch()

		guard = Patch(extractor.GetSS_ABTEST_CAMPAIGN, func() (map[string]*mvutil.SSTestTracfficRule, bool) {
			return make(map[string]*mvutil.SSTestTracfficRule), true
		})
		defer guard.Unpatch()

		guard = Patch(redis.LocalRedisAlgoHExists, func(key, field string) (bool, error) {
			return false, nil
		})
		defer guard.Unpatch()

		Convey("return []Ad{}", func() {
			res = renderMobvistaCampaigns(&r, &ads)
			So(res, ShouldResemble, []Ad{})
		})

		Convey("正常数据", func() {
			appName1 := "app_name_1"
			appName2 := "app_name_2"
			ads.CampaignList = []*corsair_proto.Campaign{
				&corsair_proto.Campaign{
					CampaignId: "10001",
					AdSource:   ad_server.ADSource(1),
					AppName:    &appName1,
				},
				&corsair_proto.Campaign{
					CampaignId: "10002",
					AdSource:   ad_server.ADSource(2),
					AppName:    &appName2,
				},
			}
			r.Param.Extra20List = []string{"extra20_list_1", "extra20_list_2"}
			r.Param.ExtPlayableList = "ex_playbale_list"
			res = renderMobvistaCampaigns(&r, &ads)
			So(res, ShouldNotResemble, []Ad{})
		})
	})
}

func TestRenderAllParams(t *testing.T) {
	Convey("test renderAllParams", t, func() {
		var r mvutil.RequestParams
		var ads corsair_proto.BackendAds

		guard := Patch(renderAlgo, func(algo string) map[int64]string {
			return map[int64]string{
				int64(101): "algo_101",
				int64(102): "algo_102",
			}
		})
		defer guard.Unpatch()

		Convey("ads.IsAdServerTest != nil && int(*(ads.IsAdServerTest)) == 1", func() {
			isAds := int32(1)
			ifLowerImp := int32(4)
			adTemp := ad_server.ADTemplate(42)
			noneRes := corsair_proto.NoneResultReason(32)
			algoFeat := "test_algo"

			ads = corsair_proto.BackendAds{
				BackendId:        int32(13579),
				RequestKey:       "test_r_key",
				CampaignList:     []*corsair_proto.Campaign{},
				Strategy:         "test_strategy",
				RunTimeVariables: &corsair_proto.RunTimeVariable{},
				AdTemplate:       &adTemp,
				NoneResultReason: &noneRes,
				FilterReason:     []*corsair_proto.FilterReason{},
				FrameList:        []*corsair_proto.Frame{},
				AlgoFeatInfo:     &algoFeat,
				IfLowerImp:       &ifLowerImp,
				IsAdServerTest:   &isAds,
			}

			r.Param = mvutil.Params{
				Extra:         "old_extra",
				AlgoMap:       map[int64]string{int64(42): "old_algo"},
				Extalgo:       "old_extalgo",
				Extra5:        "old_extra5",
				ExtflowTagId:  int(998),
				Extra3:        "old_extra3",
				ExtifLowerImp: int32(997),
				BackendID:     int32(996),
				RequestKey:    "old_extra3",
				FlowTagID:     123,
			}

			renderAllParams(&r, ads)
			So(r.Param.Extra, ShouldEqual, "old_extra_adserver")
			So(r.Param.AlgoMap, ShouldResemble, map[int64]string{
				int64(101): "algo_101",
				int64(102): "algo_102",
			})
			So(r.Param.Extalgo, ShouldEqual, "test_algo")
			So(r.Param.Extra5, ShouldEqual, "")
			So(r.Param.ExtflowTagId, ShouldEqual, 123)
			So(r.Param.Extra3, ShouldEqual, "test_strategy")
			So(r.Param.ExtifLowerImp, ShouldEqual, 4)
			So(r.Param.BackendID, ShouldEqual, 13579)
			So(r.Param.RequestKey, ShouldEqual, "test_r_key")
		})

		Convey("else", func() {
			isAds := int32(2)
			ifLowerImp := int32(4)
			adTemp := ad_server.ADTemplate(42)
			noneRes := corsair_proto.NoneResultReason(32)
			algoFeat := "test_algo"

			ads = corsair_proto.BackendAds{
				BackendId:        int32(13579),
				RequestKey:       "test_r_key",
				CampaignList:     []*corsair_proto.Campaign{},
				Strategy:         "test_strategy",
				RunTimeVariables: &corsair_proto.RunTimeVariable{},
				AdTemplate:       &adTemp,
				NoneResultReason: &noneRes,
				FilterReason:     []*corsair_proto.FilterReason{},
				FrameList:        []*corsair_proto.Frame{},
				AlgoFeatInfo:     &algoFeat,
				IfLowerImp:       &ifLowerImp,
				IsAdServerTest:   &isAds,
			}

			r.Param = mvutil.Params{
				Extra:         "old_extra",
				AlgoMap:       map[int64]string{int64(42): "old_algo"},
				Extalgo:       "old_extalgo",
				Extra5:        "old_extra5",
				ExtflowTagId:  int(998),
				Extra3:        "old_extra3",
				ExtifLowerImp: int32(997),
				BackendID:     int32(996),
				RequestKey:    "old_extra3",
				FlowTagID:     123,
			}

			renderAllParams(&r, ads)
			So(r.Param.Extra, ShouldEqual, "old_extra_adserver")
		})
	})
}

func TestRenderAlgo(t *testing.T) {
	Convey("test renderAlgo", t, func() {
		var algo string
		var res map[int64]string

		Convey("len(algo) <= 0", func() {
			algo = ""
			res = renderAlgo(algo)
			So(res, ShouldResemble, map[int64]string{})
		})

		Convey("len(algoList) <= 0", func() {
			algo = "test"
			res = renderAlgo(algo)
			So(res, ShouldResemble, map[int64]string{})
		})
	})
}

func TestRenderBigTplUrl(t *testing.T) {
	Convey("test renderBigTplUrl", t, func() {
		tpl201 := "201_big_tpl_url"
		tplFake := "fakeTemplate_big_tpl_url"
		guard := Patch(extractor.GetTEMPLATE_MAP, func() (mvutil.GlobalTemplateMap, bool) {
			return mvutil.GlobalTemplateMap{
				BigTempalte: map[string]string{
					"201":          tpl201,
					"fakeTemplate": tplFake,
				},
			}, true
		})
		defer guard.Unpatch()
		Convey("pioneer 返回结果", func() {
			bigTplUrl := "pioneer_big_tpl_url"
			r := &mvutil.RequestParams{
				Param: mvutil.Params{
					BigTemplateUrl: bigTplUrl,
				},
			}
			thirdDspId := int64(mvconst.MAS)
			backend := &corsair_proto.BackendAds{}
			res := renderBigTplUrl(backend, r, thirdDspId)
			So(res, ShouldEqual, bigTplUrl)
		})
		Convey("as dsp 走201大模版，返回非v11模版", func() {
			r := &mvutil.RequestParams{
				Param: mvutil.Params{
					BigTemplateUrl: "",
					BigTemplateId:  int64(201),
					ExtBigTemId:    "201",
				},
			}
			thirdDspId := int64(mvconst.FakeAdserverDsp)
			backend := &corsair_proto.BackendAds{
				CampaignList: []*corsair_proto.Campaign{},
			}
			res := renderBigTplUrl(backend, r, thirdDspId)
			So(res, ShouldEqual, "http://"+tpl201)
		})
		Convey("as dsp 走101大模版，返回非v11模版", func() {
			r := &mvutil.RequestParams{
				Param: mvutil.Params{
					BigTemplateUrl: "",
					BigTemplateId:  int64(101),
					ExtBigTemId:    "101",
				},
			}
			thirdDspId := int64(mvconst.FakeAdserverDsp)
			backend := &corsair_proto.BackendAds{}
			res := renderBigTplUrl(backend, r, thirdDspId)
			So(res, ShouldEqual, "")
		})
		Convey("as dsp 走101大模版，返回v11模版", func() {
			r := &mvutil.RequestParams{
				Param: mvutil.Params{
					BigTemplateUrl: "",
					BigTemplateId:  int64(101),
					ExtBigTemId:    "101",
				},
			}
			thirdDspId := int64(mvconst.FakeAdserverDsp)
			tplId := ad_server.VideoTemplateId_V11_ZIP
			backend := &corsair_proto.BackendAds{
				CampaignList: []*corsair_proto.Campaign{
					0: &corsair_proto.Campaign{
						VideoTemplateId: &tplId,
					},
				},
			}
			res := renderBigTplUrl(backend, r, thirdDspId)
			So(res, ShouldEqual, "http://"+tplFake)
		})
		Convey("as dsp 不切大模版，返回v11模版", func() {
			r := &mvutil.RequestParams{
				Param: mvutil.Params{
					BigTemplateUrl: "",
				},
			}
			thirdDspId := int64(mvconst.FakeAdserverDsp)
			tplId := ad_server.VideoTemplateId_V11_ZIP
			backend := &corsair_proto.BackendAds{
				CampaignList: []*corsair_proto.Campaign{
					0: &corsair_proto.Campaign{
						VideoTemplateId: &tplId,
					},
				},
			}
			res := renderBigTplUrl(backend, r, thirdDspId)
			So(res, ShouldEqual, "http://"+tplFake)
		})
	})
}

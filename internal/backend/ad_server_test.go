package backend

import (
	"testing"

	. "github.com/bouk/monkey"
	. "github.com/smartystreets/goconvey/convey"
	"gitlab.mobvista.com/ADN/adnet/internal/corsair_proto"
	"gitlab.mobvista.com/ADN/adnet/internal/extractor"
	"gitlab.mobvista.com/ADN/adnet/internal/mvconst"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	smodel "gitlab.mobvista.com/ADN/structs/model"
	"gitlab.mobvista.com/adserver/recommend_protocols/go/ad_server"
)

func TestFilterBackend(t *testing.T) {
	Convey("filterBackend 返回 false", t, func() {
		backend := AdServer{
			Backend{},
		}
		reqCtx := mvutil.ReqCtx{}

		res := backend.filterBackend(&reqCtx)
		So(res, ShouldEqual, mvconst.BackendOK)
	})
}

func TestComposeAdServerRequest(t *testing.T) {
	Convey("omposeAdServerRequest params is invalidate", t, func() {
		reqCtx := mvutil.ReqCtx{}
		bid := 1
		res, err := composeAdServerRequest(&reqCtx, bid)
		So(err, ShouldNotBeNil)
		So(res, ShouldBeNil)
	})

	Convey("omposeAdServerRequest params is validate", t, func() {
		bid := 1
		appID := int64(1000)
		unitID := int64(2000)
		scenario := "online"
		adType := int32(42)
		realAppID := int64(1001)

		param := mvutil.Params{
			AppID:         appID,
			UnitID:        unitID,
			Scenario:      scenario,
			AdType:        adType,
			DisplayCamIds: []int64{1, 2, 3},
			RealAppID:     realAppID,
			Debug:         1,
		}
		unitInfo := smodel.UnitInfo{Unit: smodel.Unit{PlacementId: int64(123)}}
		//unitInfo := mvutil.UnitInfo{Unit: mvutil.Unit{PlacementId: int64(123)}}

		reqParams := mvutil.RequestParams{Param: param, UnitInfo: &unitInfo}
		reqCtx := mvutil.ReqCtx{ReqParams: &reqParams}

		var commonConfig mvutil.CommonConfig
		commonConfig.LogConfig.OutputFullReqRes = false
		mvutil.Config.CommonConfig = &commonConfig

		guard := Patch(extractor.GetPassthroughData, func() []string {
			return []string{}
		})
		defer guard.Unpatch()

		res, err := composeAdServerRequest(&reqCtx, bid)
		So(err, ShouldBeNil)

		//So((*res).Timestamp, ShouldEqual, 10000)
		So(*(*res).AppId, ShouldEqual, 1000)
		So(*(*res).UnitId, ShouldEqual, 2000)
		So(*(*res).Scenario, ShouldEqual, "online")
		So(*(*res).AdTypeStr, ShouldEqual, "native")

		//imp_cap
		reqCtx.ReqParams.Param.IsBlockByImpCap = true
		res, err = composeAdServerRequest(&reqCtx, bid)
		So(err, ShouldBeNil)
		So(*(*res).TrueNum, ShouldEqual, int64(0))
	})
}

func TestFillCampaign(t *testing.T) {
	Convey("fillAdServerAds params invalidate", t, func() {
		coCampaign := corsair_proto.Campaign{}
		adCampaign := ad_server.Campaign{}
		fillCampaign(&coCampaign, &adCampaign)
		So(coCampaign.CampaignId, ShouldEqual, "0")

		cAdTemplate := ad_server.ADTemplate(98)
		cOfferType := int32(95)
		cCreativeID := "corsair_creative_id"
		cAdEleTempl := ad_server.AdElementTemplate(93)
		cDownload := int32(92)
		coCampaign = corsair_proto.Campaign{
			CampaignId:        "corsair_c_id",
			AdSource:          ad_server.ADSource(99),
			AdTemplate:        &cAdTemplate,
			ImageSizeId:       ad_server.ImageSizeEnum(97),
			OfferType:         &cOfferType,
			BtType:            ad_server.BtType(94),
			CreativeId:        &cCreativeID,
			CreativeTypeIdMap: map[ad_server.CreativeType]int64{},
			AdElementTemplate: &cAdEleTempl,
			Downloadtest:      &cDownload,
		}
		adAdTemplate := ad_server.ADTemplate(32)
		adOfferType := int32(35)
		adAdEleTemp := ad_server.AdElementTemplate(37)
		adDownload := int32(38)
		adCreativeID := "adserver_creative_id"
		adCampaign = ad_server.Campaign{
			CampaignId:         int64(10000001),
			AdSource:           ad_server.ADSource(31),
			AdTemplate:         &adAdTemplate,
			ImageSizeId:        ad_server.ImageSizeEnum(33),
			OfferType:          &adOfferType,
			BtType:             ad_server.BtType(36),
			CreativeId:         &adCreativeID,
			CreativeTypeIdMap:  map[ad_server.CreativeType]int64{ad_server.CreativeType(301): int64(3001)},
			AdElementTemplate:  &adAdEleTemp,
			Downloadtest:       &adDownload,
			CreativeId2:        &adCreativeID,
			CreativeTypeIdMap2: map[ad_server.CreativeType]int64{},
		}
		fillCampaign(&coCampaign, &adCampaign)
		So(coCampaign.CampaignId, ShouldEqual, "10000001")
		So(coCampaign.AdSource, ShouldEqual, 31)
		So(coCampaign.AdTemplate, ShouldEqual, &adAdTemplate)
		So(coCampaign.ImageSizeId, ShouldEqual, 33)
		So(coCampaign.OfferType, ShouldEqual, &adOfferType)
		So(coCampaign.BtType, ShouldEqual, 36)
		So(coCampaign.CreativeId, ShouldEqual, &adCreativeID)
		So(coCampaign.CreativeTypeIdMap, ShouldResemble, map[ad_server.CreativeType]int64{ad_server.CreativeType(301): int64(3001)})
		So(coCampaign.CreativeId2, ShouldEqual, &adCreativeID)
		So(coCampaign.CreativeTypeIdMap2, ShouldResemble, map[ad_server.CreativeType]int64{})
		So(coCampaign.AdElementTemplate, ShouldEqual, &adAdEleTemp)
	})
}

func TestRenderPassthroughData(t *testing.T) {
	param := new(mvutil.Params)
	param.OAID = "oaid"
	param.IDFV = "idfa"
	param.OsvUpTime = "osv_uptime"
	param.UpTime = "uptime"
	param.CachedCampaignIds = "cached_campaign_ids"
	param.OnlineApiNeedOfferBidPrice = "online_api_need_offer_bid_price"
	param.UseDynamicTmax = "use_dynamic_tmax"
	param.ImpExcPkgNames = []string{"aa", "bb", "cc"}
	Convey("Test RenderPassthroughData", t, func() {
		guard := Patch(extractor.GetPassthroughData, func() []string {
			return []string{"OAID", "IDFV", "OsvUpTime", "UpTime", "CachedCampaignIds", "OnlineApiNeedOfferBidPrice", "UseDynamicTmax", "ImpExcPkgNames"}
		})
		defer guard.Unpatch()

		res := RenderPassthroughData(param)
		So(res["OAID"], ShouldEqual, "oaid")
		So(res["IDFV"], ShouldEqual, "idfa")
		So(res["OsvUpTime"], ShouldEqual, "osv_uptime")
		So(res["UpTime"], ShouldEqual, "uptime")
		So(res["CachedCampaignIds"], ShouldEqual, "cached_campaign_ids")
		So(res["OnlineApiNeedOfferBidPrice"], ShouldEqual, "online_api_need_offer_bid_price")
		So(res["UseDynamicTmax"], ShouldEqual, "use_dynamic_tmax")
		So(res["ImpExcPkgNames"], ShouldEqual, "aa,bb,cc")
	})
}

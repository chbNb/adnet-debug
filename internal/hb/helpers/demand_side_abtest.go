package helpers

import (
	"gitlab.mobvista.com/ADN/adnet/internal/hb/constant"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/params"
	"gitlab.mobvista.com/ADN/structs/model"
	smodel "gitlab.mobvista.com/ADN/structs/model"
)

type abTestFunc func(ad *params.Ad, campaign *smodel.CampaignInfo, loadReq *params.LoadReqData) (key string, value interface{})

var abTestFuncs = []abTestFunc{
	//adJustPostbackABTest,
	impExcludePkgABTest,
	// 新频次控制新增字段
	priceFactorABTest,
	PriceFactorGroupNameABTest,
	PriceFactorTagABTest,
	PriceFactorFreqABTest,
	PriceFactorHitABTest,
	DcoAbTest,
}

func RenderDemandSideABTest(ad *params.Ad, campaign *smodel.CampaignInfo, loadReq *params.LoadReqData) map[string]interface{} {
	labelsMap := make(map[string]interface{})

	for _, funcv := range abTestFuncs {
		key, val := funcv(ad, campaign, loadReq)
		if key == "" || val == nil {
			continue
		}

		labelsMap[key] = val
	}

	return labelsMap
}

func DcoAbTest(ad *params.Ad, campaign *model.CampaignInfo, loadReq *params.LoadReqData) (key string, value interface{}) {
	if loadReq.DcoTestTag == 0 {
		return
	}
	return constant.EXT_DATA_KEY_DCO_TEST_FLAG, loadReq.DcoTestTag
}

func impExcludePkgABTest(ad *params.Ad, campaign *smodel.CampaignInfo, loadReq *params.LoadReqData) (key string, value interface{}) {
	if len(loadReq.ImpExcludePkg) == 0 {
		return
	}
	return constant.EXT_DATA_KEY_IMP_EXCLUDE_PKG, loadReq.ImpExcludePkg
}

// func adJustPostbackABTest(ad *params.Ad, campaign *smodel.CampaignInfo, loadReq *params.LoadReqData) (key string, value interface{}) {
// 	cfg := extractor.GetAdjustPostbackABTest()
// 	if !cfg.Status {
// 		return
// 	}

// 	if !campaign.IsSSPlatform() {
// 		return
// 	}

// 	if strings.ToLower(campaign.ThirdParty) != constant.THIRD_PARTY_ADJUST {
// 		return
// 	}

// 	if len(cfg.Scenario) > 0 && !utils.StringInSlice(loadReq.Scenario, cfg.Scenario) {
// 		return
// 	}

// 	if len(cfg.BListCampaign) > 0 && utils.Int64InSlice(campaign.CampaignId, cfg.BListCampaign) {
// 		return
// 	}

// 	if len(cfg.WListCampaign) > 0 && !utils.Int64InSlice(campaign.CampaignId, cfg.WListCampaign) {
// 		return
// 	}

// 	if rand.Intn(10000) < cfg.Rate {
// 		loadReq.BidId = loadReq.BidId[:len(loadReq.BidId)-1] + "v"
// 		return constant.EXT_DATA_KEY_ADJUST_POSTBACK, 2
// 	}

// 	return constant.EXT_DATA_KEY_ADJUST_POSTBACK, 1
// }

func priceFactorABTest(ad *params.Ad, campaign *smodel.CampaignInfo, loadReq *params.LoadReqData) (key string, value interface{}) {
	if loadReq.Add2ExtData2.PriceFactor == 0 {
		return
	}
	return constant.EXT_DATA_KEY_PRICE_FACTOR, loadReq.Add2ExtData2.PriceFactor
}
func PriceFactorGroupNameABTest(ad *params.Ad, campaign *smodel.CampaignInfo, loadReq *params.LoadReqData) (key string, value interface{}) {
	if len(loadReq.Add2ExtData2.PriceFactorGroupName) == 0 {
		return
	}
	return constant.EXT_DATA_KEY_PRICE_FACTOR_GROUP_NAME, loadReq.Add2ExtData2.PriceFactorGroupName
}
func PriceFactorTagABTest(ad *params.Ad, campaign *smodel.CampaignInfo, loadReq *params.LoadReqData) (key string, value interface{}) {
	if loadReq.Add2ExtData2.PriceFactorTag == 0 {
		return
	}
	return constant.EXT_DATA_KEY_PRICE_FACTOR_TAG, loadReq.Add2ExtData2.PriceFactorTag
}
func PriceFactorFreqABTest(ad *params.Ad, campaign *smodel.CampaignInfo, loadReq *params.LoadReqData) (key string, value interface{}) {
	if loadReq.Add2ExtData2.PriceFactorFreq == nil {
		return
	}
	return constant.EXT_DATA_KEY_PRICE_FACTOR_FREQ, loadReq.Add2ExtData2.PriceFactorFreq
}
func PriceFactorHitABTest(ad *params.Ad, campaign *smodel.CampaignInfo, loadReq *params.LoadReqData) (key string, value interface{}) {
	if loadReq.Add2ExtData2.PriceFactorHit == 0 {
		return
	}
	return constant.EXT_DATA_KEY_PRICE_FACTOR_HIT, loadReq.Add2ExtData2.PriceFactorHit

}

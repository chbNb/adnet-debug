package c2s

import (
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"gitlab.mobvista.com/ADN/adnet/internal/extractor"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/helpers"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/output"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/params"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/req_context"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
)

var RenderInputError = errors.New("c2s render input error")

func RenderOutput(bidData *mvutil.ReqCtx) (*params.BidResp, error) {
	if bidData == nil || bidData.Result == nil {
		return nil, RenderInputError
	}

	var bidResp params.BidResp
	bidResp.Status = 200
	bidResp.Msg = "ok"
	cur := bidData.ReqParams.BidCur
	// 美分转美元
	bidPrice := bidData.ReqParams.Price
	if price, err := output.PriceConversion(bidData.ReqParams.Price, 100.0, "div"); err == nil {
		bidPrice = price
	}
	// publisher bid response currency proccess
	publisherIdStr := strconv.FormatInt(bidData.ReqParams.Param.PublisherID, 10)
	if publisherCurrency, ok := extractor.GetHBPublisherCurrency(publisherIdStr); ok {
		if currencyRate, ok := extractor.GetHBCurrencyRate(publisherCurrency); ok && currencyRate > 0 {
			// 人民币元
			if toCurrencyPrice, err := output.PriceConversion(bidPrice, currencyRate, "mul"); err == nil {
				bidPrice = toCurrencyPrice
				currency := strings.Split(publisherCurrency, "_")
				if len(currency) > 1 {
					cur = currency[1]
				} else {
					cur = currency[0]
				}
				bidData.ReqParams.CurrencyPrice = bidPrice * 100
				bidData.ReqParams.BidIsNotUSDCur = 1
			}
		}
	}
	bidData.ReqParams.Currency = cur

	// 为节省出口流量费用，取消下发td参数
	// c2s量不大，直接切
	adnetConf, _ := extractor.GetADNET_SWITCHS()
	c2sDelTd, ok := adnetConf["c2sDelTd"]
	var lossUrl, winUrl string
	if ok && c2sDelTd == 1 {
		lossUrl = GetLossWithMacrosWithoutTd()
		winUrl = GetWinWithMacorsWithoutTd()
	} else {
		lossUrl = GetLossWithMacros()
		winUrl = GetWinWithMacors()
	}

	bidResp.Data = &params.BidRespData{
		Bid:     bidData.ReqParams.BidRespID,
		Token:   bidData.ReqParams.Token,
		Price:   bidPrice,
		Cur:     cur,
		LossUrl: lossUrl,
		WinUrl:  winUrl,
		Macors:  make(map[string]string, 2),
	}

	bidResp.Data.Macors["sd"] = helpers.GenFullUrl(int(bidData.ReqParams.Param.HTTPReq),
		output.GetDomain(bidData.ReqParams.Param.CountryCode, req_context.GetInstance().Cfg.ServerCfg.AerospikeMultiZone))
	if ok && c2sDelTd == 1 {
		bidResp.Data.Macors["token"] = bidData.ReqParams.Token
	} else {
		bidResp.Data.Macors["td"] = output.GenTokenData(bidData)
	}
	return &bidResp, nil
}

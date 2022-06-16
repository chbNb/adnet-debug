package bid_filters

import (
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"gitlab.mobvista.com/ADN/adnet/internal/extractor"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/constant"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/filter"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/helpers"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/output"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/req_context"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	adn_output "gitlab.mobvista.com/ADN/adnet/internal/output"
)

type BidAdsFormatOutputFilter struct {
}

func (bafof *BidAdsFormatOutputFilter) Process(data interface{}) (interface{}, error) {
	in, ok := data.(*mvutil.ReqCtx)
	if !ok {
		return nil, filter.FormatOutputFilterInputError
	}

	// 记录需要记录到hb request日志的字段
	adn_output.RenderHBRequestLogParam(in.ReqParams)

	// 记录了 loss request 就不记录 request
	if !in.ReqParams.Param.LossReqFlag {
		mvutil.StatHBRequestLog(in)
	}

	if in.Result == nil || in.ReqParams.Nbr != constant.OK {
		return in, errors.New(strings.Join(in.ReqParams.Param.BackendReject, ";"))
	}

	// add ab flag
	in.ReqParams.Param.HBExtPfData = in.ReqParams.Param.ExtData

	cur := in.ReqParams.BidCur
	// 美分转美元
	bidPrice := in.ReqParams.Price
	if price, err := output.PriceConversion(in.ReqParams.Price, 100.0, "div"); err == nil {
		bidPrice = price
	}
	// publisher bid response currency proccess
	publisherIdStr := strconv.FormatInt(in.ReqParams.Param.PublisherID, 10)
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
				in.ReqParams.CurrencyPrice = bidPrice * 100
				in.ReqParams.BidIsNotUSDCur = 1
			}
		}
	}

	in.ReqParams.Currency = cur
	in.ReqParams.BidPrice = bidPrice
	prefix := helpers.GenFullUrl(int(in.ReqParams.Param.HTTPReq),
		output.GetDomain(in.ReqParams.Param.CountryCode, req_context.GetInstance().Cfg.ServerCfg.AerospikeMultiZone))

	td := output.GenTokenData(in)
	in.ReqParams.BidWinUrl = prefix + "/win?td=" + td
	// bid log
	req_context.GetInstance().MLogs.Bid.Info(output.FormatBidLog(in, 0, ""))

	return in, nil
}

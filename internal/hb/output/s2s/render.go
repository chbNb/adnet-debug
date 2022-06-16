package s2s

import (
	"math/rand"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"gitlab.mobvista.com/ADN/adnet/internal/extractor"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/constant"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/helpers"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/output"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/req_context"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	"gitlab.mobvista.com/ADN/mtg_openrtb/pkg/mtg_hb_rtb"
)

const (
	LossEventURI    = "/loss?td="
	WinEventURI     = "/win?td="
	BillingURI      = "/billing?td="
	MaxLossEventURI = "/loss?token="
	MaxWinEventURI  = "/win?token="
)

var RenderInputError = errors.New("s2s render input error")

func RenderOutput(bidData *mvutil.ReqCtx) (*mtg_hb_rtb.BidResponse, error) {
	if bidData == nil || bidData.Result == nil {
		return nil, RenderInputError
	}

	var bidResp mtg_hb_rtb.BidResponse
	bidResp.Id = &bidData.ReqParams.Param.HBS2SBidID
	cur := bidData.ReqParams.BidCur
	if bidData.ReqParams.Nbr != constant.OK {
		nbr := mtg_hb_rtb.BidResponse_NoBidReason(int32(bidData.ReqParams.Nbr))
		bidResp.Nbr = &nbr
	} else {
		bidResp.Bidid = &bidData.ReqParams.Token
		bid := mtg_hb_rtb.BidResponse_SeatBid_Bid{}
		bid.Impid = &bidData.ReqParams.Param.ImpID
		bidPrice := bidData.ReqParams.Price
		// 美分转美元
		if price, err := output.PriceConversion(bidPrice, 100.0, "div"); err == nil {
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
					bidData.ReqParams.CurrencyPrice = toCurrencyPrice * 100
					bidData.ReqParams.BidIsNotUSDCur = 1
				}
			}
		}
		bidData.ReqParams.Currency = cur
		bid.Price = &bidPrice
		bid.Id = &bidData.ReqParams.Token
		bid.Adm = &bidData.ReqParams.Token
		var lURL, nURL, bURL string
		httpReqProc := int(bidData.ReqParams.Param.HTTPReq)
		cfg := extractor.GetHBEeventHTTPProtocolConf()
		if v, ok := cfg[strings.ToUpper(bidData.ReqParams.Param.MediationName)]; ok {
			httpReqProc = int(v)
		}
		prefix := helpers.GenFullUrl(httpReqProc, output.GetDomain(bidData.ReqParams.Param.CountryCode, req_context.GetInstance().Cfg.ServerCfg.AerospikeMultiZone))

		//  不返回td，改为返回token，节省流量成本
		adnetConf, _ := extractor.GetADNET_SWITCHS()
		s2sDelTdRate, ok := adnetConf["s2sDelTdRate"]
		// 控制比例，避免areospike扛不住
		randVal := rand.Intn(100)
		if (ok && (s2sDelTdRate > randVal)) || constant.GetMediationChnId(strings.ToUpper(bidData.ReqParams.Param.MediationName)) == constant.MaxMediation {
			lURL = prefix + MaxLossEventURI + bidData.ReqParams.Token
			nURL = prefix + MaxWinEventURI + bidData.ReqParams.Token
		} else {
			tokenData := output.GenTokenData(bidData)
			lURL = prefix + LossEventURI + tokenData
			nURL = prefix + WinEventURI + tokenData
		}
		val := extractor.GetMediationNoticeURLMacroConfig()
		for _, conf := range val.Configs {
			if selectorMacroConfigFilter(bidData.ReqParams.Param, conf) {
				lURL = lURL + conf.LURLMacro
			}
		}
		bid.Lurl = &lURL
		bid.Nurl = &nURL
		bURL = prefix + BillingURI + output.GenShortTokenData(bidData)
		bid.Burl = &bURL

		// 不下发Skadn了，sdk和聚合平台都没有用到
		confs, _ := extractor.GetADNET_SWITCHS()
		if bidSkadnSwitch, ok := confs["bidSkadnSwitch"]; ok && bidSkadnSwitch == 1 {
			if bidData.ReqParams.BidSkAdNetwork != nil {
				if bid.Ext == nil {
					bid.Ext = new(mtg_hb_rtb.BidResponse_SeatBid_Bid_Ext)
				}
				bid.Ext.Skadn = bidData.ReqParams.BidSkAdNetwork
			}
		}
		bidResp.Seatbid = []*mtg_hb_rtb.BidResponse_SeatBid{{Bid: []*mtg_hb_rtb.BidResponse_SeatBid_Bid{&bid}}}
	}
	bidResp.Cur = &cur
	return &bidResp, nil
}

func selectorMacroConfigFilter(req mvutil.Params, conf *mvutil.MediationNoticeURLMacroConf) bool {
	if req.Extchannel != conf.ChannelID {
		return false
	}
	if !mvutil.InInt64Arr(-1, conf.UnitIDs) && !mvutil.InInt64Arr(req.UnitID, conf.UnitIDs) {
		return false
	}
	if !mvutil.InInt64Arr(-1, conf.AppIDs) && !mvutil.InInt64Arr(req.AppID, conf.AppIDs) {
		return false
	}
	if !mvutil.InInt64Arr(-1, conf.PublisherIDs) && !mvutil.InInt64Arr(req.PublisherID, conf.PublisherIDs) {
		return false
	}
	return true
}

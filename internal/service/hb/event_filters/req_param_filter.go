package event_filters

import (
	"errors"
	"net/url"
	"strconv"
	"strings"
	"time"

	"gitlab.mobvista.com/ADN/adnet/internal/hb/constant"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/helpers"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/output"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/req_context"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	"gitlab.mobvista.com/ADN/adnet/internal/watcher"
	"gitlab.mobvista.com/ADN/exporter/metrics"
)

var ReqParamFilterInputError = errors.New("req_param_filter input error")
var TokenDataInvalidate = errors.New("token data is invalidate")
var TokenInvalidate = errors.New("token is invalidate")

type ReqParamFilter struct {
}

func (rpf *ReqParamFilter) Process(data interface{}) (interface{}, error) {
	in, ok := data.(*mvutil.RequestParams)
	if !ok {
		return nil, ReqParamFilterInputError
	}

	var rawTd string
	var eventType int
	switch in.Param.RequestPath {
	case "/win":
		eventType = 1
	case "/loss":
		eventType = 2
	case "/billing":
		eventType = 3
	}

	td, _ := in.QueryMap.GetString("td", true, "")
	td = strings.Replace(td, " ", "+", -1)
	winPrice, _ := in.QueryMap.GetString("win_price", true, "")
	auctionCurrency, _ := in.QueryMap.GetString("auction_currency", true, "")
	// TODO remove td
	if len(td) > 0 {
		rawTd = helpers.Base64Decode(td)
		if strings.Count(rawTd, constant.SplitterLine) < 19 {
			return in, TokenDataInvalidate
		}
	} else {
		token, _ := in.QueryMap.GetString("token", true, "")
		if len(token) == 0 {
			return in, TokenInvalidate
		}
		val, err := output.GetBidCache(token)
		// val, err = req_context.GetInstance().BidCacheClient.GetReqCtx(token)
		if err != nil || val.ReqParams == nil {
			watcher.AddWatchValue(constant.GetBidCacheMiss, 1)
			metrics.IncCounterWithLabelValues(7, strconv.Itoa(eventType))
			if tokenStr := strings.Split(token, "_"); len(tokenStr) == 2 {
				uuidStr := strings.SplitN(tokenStr[0], "-", -1)
				if len(uuidStr) >= 6 {
					nowTS := time.Now()
					i, _ := strconv.ParseInt(uuidStr[5], 10, 64)
					tokenTS := time.Unix(i, 0)
					t := nowTS.Sub(tokenTS)
					if t.Hours() > 1 {
						metrics.IncCounterWithLabelValues(8, strconv.Itoa(eventType))
					}
				}
			}
			return in, errors.New("get cache miss from aerospike of the event request")
		}
		reqCtx := mvutil.ReqCtx{}
		reqCtx.ReqParams = val.ReqParams
		var eventData []string
		publisherIdStr := strconv.FormatInt(reqCtx.ReqParams.Param.PublisherID, 10)
		eventData = append(eventData, publisherIdStr)
		appIdStr := strconv.FormatInt(reqCtx.ReqParams.Param.AppID, 10)
		eventData = append(eventData, appIdStr)
		unitIdStr := strconv.FormatInt(reqCtx.ReqParams.Param.UnitID, 10)
		eventData = append(eventData, unitIdStr)
		algorithm := reqCtx.ReqParams.Param.Extra
		eventData = append(eventData, algorithm)
		scenario := reqCtx.ReqParams.Param.Scenario
		eventData = append(eventData, scenario)
		adTypeStr := mvutil.GetAdTypeStr(reqCtx.ReqParams.Param.AdType)
		eventData = append(eventData, adTypeStr)
		os := helpers.GetOs(reqCtx.ReqParams.Param.Platform)
		eventData = append(eventData, os)
		sdkVersion := strings.ToLower(reqCtx.ReqParams.Param.SDKVersion)
		eventData = append(eventData, sdkVersion)
		appVersionName := reqCtx.ReqParams.Param.AppVersionName
		eventData = append(eventData, appVersionName)
		countryCode := reqCtx.ReqParams.Param.CountryCode
		eventData = append(eventData, countryCode)
		cityCode := strconv.FormatInt(reqCtx.ReqParams.Param.CityCode, 10)
		eventData = append(eventData, cityCode)
		apiOffer := strconv.Itoa(constant.APIOffer)
		eventData = append(eventData, apiOffer)
		channel := reqCtx.ReqParams.Param.Extchannel
		eventData = append(eventData, channel)
		dspID := strconv.FormatInt(reqCtx.ReqParams.DspExtData.DspId, 10)
		eventData = append(eventData, dspID, token)
		mediationName := reqCtx.ReqParams.Param.MediationName
		eventData = append(eventData, mediationName)
		price := reqCtx.ReqParams.PriceBigDecimal
		eventData = append(eventData, price)
		extra3 := url.QueryEscape(reqCtx.ReqParams.Param.Extra3)
		eventData = append(eventData, extra3)
		extAlgo := url.QueryEscape(reqCtx.ReqParams.Param.Extalgo)
		eventData = append(eventData, extAlgo)
		extAdxAlgo := url.QueryEscape(reqCtx.ReqParams.Param.ExtAdxAlgo)
		eventData = append(eventData, extAdxAlgo)
		clientIp := reqCtx.ReqParams.Param.ClientIP
		eventData = append(eventData, clientIp)
		userAgent := mvutil.RawUrlEncode(reqCtx.ReqParams.Param.UserAgent)
		eventData = append(eventData, userAgent)
		osVersion := reqCtx.ReqParams.Param.OSVersion
		eventData = append(eventData, osVersion)
		idfa := reqCtx.ReqParams.Param.IDFA
		eventData = append(eventData, idfa)
		idfv := reqCtx.ReqParams.Param.IDFV
		eventData = append(eventData, idfv)
		imei := reqCtx.ReqParams.Param.IMEI
		eventData = append(eventData, imei)
		androidId := reqCtx.ReqParams.Param.AndroidID
		eventData = append(eventData, androidId)
		gaid := reqCtx.ReqParams.Param.GAID
		eventData = append(eventData, gaid)
		requestId := reqCtx.ReqParams.Param.RequestID
		eventData = append(eventData, requestId)
		rawTd = strings.Join(eventData, "|")
	}
	req_context.GetInstance().MLogs.Event.Info(output.FormatEventLog(eventType, rawTd, winPrice, auctionCurrency))
	return in, nil
}

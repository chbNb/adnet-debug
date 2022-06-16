package output

import (
	"bytes"
	"gitlab.mobvista.com/ADN/adnet/internal/extractor"
	"gitlab.mobvista.com/ADN/adnet/internal/mvconst"
	"math/rand"
	"net/url"
	"strconv"
	"strings"

	"gitlab.mobvista.com/ADN/adnet/internal/hb/constant"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/helpers"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
)

func GenTokenData(bidData *mvutil.ReqCtx) string {
	var dspID string
	dspExt, err := bidData.ReqParams.GetDspExt()
	if err == nil {
		dspID = strconv.FormatInt(dspExt.DspId, 10)
	}

	extAlgo := bidData.ReqParams.Param.Extalgo
	// 对于online api的流量，不会下发下发token再查as，因此缩减ext_algo达到节省流量目的
	if bidData.ReqParams.Param.RequestType == mvconst.REQUEST_TYPE_OPENAPI_AD {
		adnetConf, _ := extractor.GetADNET_SWITCHS()
		onlineApiDelTdRate, ok := adnetConf["onlineApiDelTdRate"]
		randVal := rand.Intn(100)
		if ok && (onlineApiDelTdRate > randVal) {
			extAlgo = ""
		}
	}

	// publisher_id|app_id|unit_id|algorithm|scenario|ad_type|platform|sdk_version|app_version|country_code|city_code|
	// ad_source_id|channel|dsp_id|token|mediation_name|price|strategy|ext_algo|ext_adx_algo|client_ip|user_agent|
	// os_version|idfa|idfv|imei|android_id|gaid|requestid
	var buf bytes.Buffer
	buf.WriteString(strconv.FormatInt(bidData.ReqParams.Param.PublisherID, 10))
	buf.WriteString(constant.SplitterLine)
	buf.WriteString(strconv.FormatInt(bidData.ReqParams.Param.AppID, 10))
	buf.WriteString(constant.SplitterLine)
	buf.WriteString(strconv.FormatInt(bidData.ReqParams.Param.UnitID, 10))
	buf.WriteString(constant.SplitterLine)
	buf.WriteString(bidData.ReqParams.Param.Extra)
	buf.WriteString(constant.SplitterLine)
	buf.WriteString(bidData.ReqParams.Param.Scenario)
	buf.WriteString(constant.SplitterLine)
	buf.WriteString(mvutil.GetAdTypeStr(bidData.ReqParams.Param.AdType))
	buf.WriteString(constant.SplitterLine)
	buf.WriteString(helpers.GetOs(bidData.ReqParams.Param.Platform))
	buf.WriteString(constant.SplitterLine)
	buf.WriteString(strings.ToLower(bidData.ReqParams.Param.SDKVersion))
	buf.WriteString(constant.SplitterLine)
	buf.WriteString(bidData.ReqParams.Param.AppVersionName)
	buf.WriteString(constant.SplitterLine)
	buf.WriteString(bidData.ReqParams.Param.CountryCode)
	buf.WriteString(constant.SplitterLine)
	buf.WriteString(strconv.FormatInt(bidData.ReqParams.Param.CityCode, 10))
	buf.WriteString(constant.SplitterLine)
	buf.WriteString(strconv.Itoa(constant.APIOffer))
	buf.WriteString(constant.SplitterLine)
	buf.WriteString(bidData.ReqParams.Param.Extchannel)
	buf.WriteString(constant.SplitterLine)
	buf.WriteString(dspID)
	buf.WriteString(constant.SplitterLine)
	buf.WriteString(bidData.ReqParams.Token)
	buf.WriteString(constant.SplitterLine)
	buf.WriteString(bidData.ReqParams.Param.MediationName)
	buf.WriteString(constant.SplitterLine)
	buf.WriteString(bidData.ReqParams.PriceBigDecimal)
	buf.WriteString(constant.SplitterLine)
	buf.WriteString(url.QueryEscape(bidData.ReqParams.Param.Extra3))
	buf.WriteString(constant.SplitterLine)
	buf.WriteString(url.QueryEscape(extAlgo))
	buf.WriteString(constant.SplitterLine)
	buf.WriteString(url.QueryEscape(bidData.ReqParams.Param.ExtAdxAlgo))
	buf.WriteString(constant.SplitterLine)
	buf.WriteString(bidData.ReqParams.Param.ClientIP)
	buf.WriteString(constant.SplitterLine)
	buf.WriteString(mvutil.RawUrlEncode(bidData.ReqParams.Param.UserAgent))
	buf.WriteString(constant.SplitterLine)
	buf.WriteString(bidData.ReqParams.Param.OSVersion)
	buf.WriteString(constant.SplitterLine)
	buf.WriteString(bidData.ReqParams.Param.IDFA)
	buf.WriteString(constant.SplitterLine)
	buf.WriteString(bidData.ReqParams.Param.IDFV)
	buf.WriteString(constant.SplitterLine)
	buf.WriteString(bidData.ReqParams.Param.IMEI)
	buf.WriteString(constant.SplitterLine)
	buf.WriteString(bidData.ReqParams.Param.AndroidID)
	buf.WriteString(constant.SplitterLine)
	buf.WriteString(bidData.ReqParams.Param.GAID)
	buf.WriteString(constant.SplitterLine)
	buf.WriteString(bidData.ReqParams.Param.RequestID)
	return helpers.Base64Encode(buf.String())
}

func GenShortTokenData(bidData *mvutil.ReqCtx) string {
	var dspID string
	dspExt, err := bidData.ReqParams.GetDspExt()
	if err == nil {
		dspID = strconv.FormatInt(dspExt.DspId, 10)
	}
	// publisher_id|app_id|unit_id|algorithm|scenario|ad_type|platform|sdk_version|app_version|country_code|city_code|
	// ad_source_id|channel|dsp_id|token|mediation_name|price|strategy|ext_algo|ext_adx_algo|client_ip|user_agent|
	// os_version|idfa|idfv|imei|android_id|gaid|requestid
	var buf bytes.Buffer
	buf.WriteString(strconv.FormatInt(bidData.ReqParams.Param.PublisherID, 10))
	buf.WriteString(constant.SplitterLine)
	buf.WriteString(strconv.FormatInt(bidData.ReqParams.Param.AppID, 10))
	buf.WriteString(constant.SplitterLine)
	buf.WriteString(strconv.FormatInt(bidData.ReqParams.Param.UnitID, 10))
	buf.WriteString(constant.SplitterLine)
	buf.WriteString("")
	buf.WriteString(constant.SplitterLine)
	buf.WriteString("")
	buf.WriteString(constant.SplitterLine)
	buf.WriteString(mvutil.GetAdTypeStr(bidData.ReqParams.Param.AdType))
	buf.WriteString(constant.SplitterLine)
	buf.WriteString(helpers.GetOs(bidData.ReqParams.Param.Platform))
	buf.WriteString(constant.SplitterLine)
	buf.WriteString(strings.ToLower(bidData.ReqParams.Param.SDKVersion))
	buf.WriteString(constant.SplitterLine)
	buf.WriteString("")
	buf.WriteString(constant.SplitterLine)
	buf.WriteString(bidData.ReqParams.Param.CountryCode)
	buf.WriteString(constant.SplitterLine)
	buf.WriteString("")
	buf.WriteString(constant.SplitterLine)
	buf.WriteString(strconv.Itoa(constant.APIOffer))
	buf.WriteString(constant.SplitterLine)
	buf.WriteString(bidData.ReqParams.Param.Extchannel)
	buf.WriteString(constant.SplitterLine)
	buf.WriteString(dspID)
	buf.WriteString(constant.SplitterLine)
	buf.WriteString(bidData.ReqParams.Token)
	buf.WriteString(constant.SplitterLine)
	buf.WriteString("")
	buf.WriteString(constant.SplitterLine)
	buf.WriteString("")
	buf.WriteString(constant.SplitterLine)
	buf.WriteString("")
	buf.WriteString(constant.SplitterLine)
	buf.WriteString("")
	buf.WriteString(constant.SplitterLine)
	buf.WriteString("")
	buf.WriteString(constant.SplitterLine)
	buf.WriteString("")
	buf.WriteString(constant.SplitterLine)
	buf.WriteString("")
	buf.WriteString(constant.SplitterLine)
	buf.WriteString("")
	buf.WriteString(constant.SplitterLine)
	buf.WriteString("")
	buf.WriteString(constant.SplitterLine)
	buf.WriteString("")
	buf.WriteString(constant.SplitterLine)
	buf.WriteString("")
	buf.WriteString(constant.SplitterLine)
	buf.WriteString("")
	buf.WriteString(constant.SplitterLine)
	buf.WriteString("")
	buf.WriteString(constant.SplitterLine)
	buf.WriteString(bidData.ReqParams.Param.RequestID)
	return helpers.Base64Encode(buf.String())
}

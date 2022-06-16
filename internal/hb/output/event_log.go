package output

import (
	"bytes"
	"net/url"
	"strconv"
	"strings"
	"time"

	"gitlab.mobvista.com/ADN/adnet/internal/hb/constant"
)

func FormatEventLog(eventType int, tokenData, winPrice, currency string) string {
	// tokenData format
	// publisher_id|app_id|unit_id|algorithm|scenario|ad_type|platform|sdk_version|app_version|country_code|city_code|
	// ad_source_id|channel|dsp_id|token|mediation_name|price|strategy|ext_algo|ext_adx_algo|client_ip|user_agent|
	// os_version|idfa|idfv|imei|android_id|gaid|requestid
	tokenArr := strings.SplitN(tokenData, "|", 29)
	var price, strategy, extAlgo, extAdxAlgo, clientIp, userAgent, osVersion, idfa, idfv, imei, androidId, gaid, requestId string
	price, _ = url.QueryUnescape(tokenArr[16])
	strategy, _ = url.QueryUnescape(tokenArr[17])
	extAlgo, _ = url.QueryUnescape(tokenArr[18])
	extAdxAlgo, _ = url.QueryUnescape(tokenArr[19])
	clientIp = tokenArr[20]
	userAgent = tokenArr[21]
	if len(tokenArr) > 22 {
		osVersion = tokenArr[22]
		idfa = tokenArr[23]
		idfv = tokenArr[24]
		imei = tokenArr[25]
		androidId = tokenArr[26]
		gaid = tokenArr[27]
		requestId = tokenArr[28]
	}
	var buf bytes.Buffer
	buf.WriteString(strconv.FormatInt(time.Now().Unix(), 10))
	buf.WriteString(constant.SplitterTab)
	buf.WriteString(tokenArr[0]) // publisher_id
	buf.WriteString(constant.SplitterTab)
	buf.WriteString(tokenArr[1]) // app_id
	buf.WriteString(constant.SplitterTab)
	buf.WriteString(tokenArr[2]) // unit_id
	buf.WriteString(constant.SplitterTab)
	buf.WriteString(tokenArr[3]) // algorithm
	buf.WriteString(constant.SplitterTab)
	buf.WriteString(tokenArr[4]) // scenario
	buf.WriteString(constant.SplitterTab)
	buf.WriteString(tokenArr[5]) // ad_type
	buf.WriteString(constant.SplitterTab)
	buf.WriteString(tokenArr[6]) // platform
	buf.WriteString(constant.SplitterTab)
	buf.WriteString(tokenArr[7]) // sdk_version
	buf.WriteString(constant.SplitterTab)
	buf.WriteString(tokenArr[8]) // app_version
	buf.WriteString(constant.SplitterTab)
	buf.WriteString(tokenArr[9]) // country_code
	buf.WriteString(constant.SplitterTab)
	buf.WriteString(tokenArr[10]) // city_code
	buf.WriteString(constant.SplitterTab)
	buf.WriteString(tokenArr[11]) // ad_source_id
	buf.WriteString(constant.SplitterTab)
	buf.WriteString(tokenArr[12]) // channel
	buf.WriteString(constant.SplitterTab)
	buf.WriteString(tokenArr[13]) // dsp_id
	buf.WriteString(constant.SplitterTab)
	buf.WriteString(tokenArr[14]) // token
	buf.WriteString(constant.SplitterTab)
	buf.WriteString(requestId)
	buf.WriteString(constant.SplitterTab)
	buf.WriteString(price)
	buf.WriteString(constant.SplitterTab)
	buf.WriteString("0")
	buf.WriteString(constant.SplitterTab)
	buf.WriteString("0")
	buf.WriteString(constant.SplitterTab)
	buf.WriteString("0")
	buf.WriteString(constant.SplitterTab)
	buf.WriteString(strconv.Itoa(eventType))
	buf.WriteString(constant.SplitterTab)
	buf.WriteString(clientIp)
	buf.WriteString(constant.SplitterTab)
	buf.WriteString(tokenArr[15]) // mediation_name
	buf.WriteString(constant.SplitterTab)
	buf.WriteString(strategy)
	buf.WriteString(constant.SplitterTab)
	buf.WriteString(extAlgo)
	buf.WriteString(constant.SplitterTab)
	buf.WriteString(extAdxAlgo)
	buf.WriteString(constant.SplitterTab)
	buf.WriteString(userAgent)
	buf.WriteString(constant.SplitterTab)
	buf.WriteString(osVersion)
	buf.WriteString(constant.SplitterTab)
	buf.WriteString(idfa)
	buf.WriteString(constant.SplitterTab)
	buf.WriteString(idfv)
	buf.WriteString(constant.SplitterTab)
	buf.WriteString(imei)
	buf.WriteString(constant.SplitterTab)
	buf.WriteString(androidId)
	buf.WriteString(constant.SplitterTab)
	buf.WriteString(gaid)
	buf.WriteString(constant.SplitterTab)
	buf.WriteString(winPrice)
	buf.WriteString(constant.SplitterTab)
	buf.WriteString(currency)
	return buf.String()
}

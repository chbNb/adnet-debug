package errorcode

import jsoniter "github.com/json-iterator/go"

type AdnetCode int

const (
	MESSAGE_SUCCESS                                 AdnetCode = 0
	MESSAGE_SUCCESS_OLD                             AdnetCode = 200
	EXCEPTION_RETURN_EMPTY                          AdnetCode = -1
	EXCEPTION_RETURN_EMPTY_OLD                      AdnetCode = 300
	EXCEPTION_PARAMS_EMPTY                          AdnetCode = -2
	EXCEPTION_PARAMS_ERROR                          AdnetCode = -3
	EXCEPTION_TIMEOUT                               AdnetCode = -9
	EXCEPTION_SIGN_ERROR                            AdnetCode = -10
	EXCEPTION_COUNTRY_NOT_ALLOW                     AdnetCode = -11
	EXCEPTION_TOKEN_FORMAT_ERROR                    AdnetCode = -12
	EXCEPTION_DOMAIN_ERROR                          AdnetCode = -13
	EXCEPTION_NETWORK_ERROR                         AdnetCode = -14
	EXCEPTION_IP_NOT_ALLOW                          AdnetCode = -2107
	EXCEPTION_UNIT_ID_EMPTY                         AdnetCode = -1202
	EXCEPTION_UNIT_NOT_FOUND                        AdnetCode = -1201
	EXCEPTION_UNIT_NOT_FOUND_IN_APP                 AdnetCode = -1203
	EXCEPTION_UNIT_NOT_APPWALL                      AdnetCode = -1204
	EXCEPTION_UNIT_ADTYPE_ERROR                     AdnetCode = -1205
	EXCEPTION_UNIT_NOT_ACTIVE                       AdnetCode = -1206
	EXCEPTION_UNIT_NO_CONFIG_MANAGE_REVENUE         AdnetCode = -1207
	EXCEPTION_UNIT_BIDDING_TYPE_ERROR               AdnetCode = -1208
	EXCEPTION_APP_ID_EMPTY                          AdnetCode = -1301
	EXCEPTION_APP_NOT_FOUND                         AdnetCode = -1302
	EXCEPTION_APP_BANNED                            AdnetCode = -1303
	EXCEPTION_PUBLISHER_BANNED                      AdnetCode = -1304
	EXCEPTION_APP_PLATFORM_ERROR                    AdnetCode = -1305
	EXCEPTION_PUBLISHER_ID_EMPTY                    AdnetCode = -1306
	EXCEPTION_PUBLISHER_NOT_FOUND                   AdnetCode = -1307
	EXCEPTION_CAMPAIGN_NOT_FOUND                    AdnetCode = -1401
	EXCEPTION_CAMPAIGN_ID_EMPTY                     AdnetCode = -1402
	EXCEPTION_CAMPAIGN_NOT_ACTIVE                   AdnetCode = -1403
	EXCEPTION_REQUEST_ID_EMPTY                      AdnetCode = -1501
	EXCEPTION_SERVICE_REQUEST_OS_VERSION_REQUIRED   AdnetCode = -2102
	EXCEPTION_SERVICE_REQUEST_COUNTRY_CODE_REQUIRED AdnetCode = -2103
	EXCEPTION_SERVICE_REQUEST_AD_SOURCE_CLOSED      AdnetCode = -2104
	EXCEPTION_GAID_EMPTY                            AdnetCode = -1801
	EXCEPTION_IDFA_EMPTY                            AdnetCode = -1802
	EXCEPTION_IMEI_EMPTY                            AdnetCode = -1803
	EXCEPTION_BRAND_EMPTY                           AdnetCode = -1804
	EXCEPTION_MODEL_EMPTY                           AdnetCode = -1805
	EXCEPTION_ANDROIDID_EMPTY                       AdnetCode = -1806
	EXCEPTION_NETWORK_TYPE_EMPTY                    AdnetCode = -1807
	EXCEPTION_BANNER_UNIT_SIZE_EMPTY                AdnetCode = -1808
	EXCEPTION_ADNUM_SET_NONE                        AdnetCode = -1901
	EXCEPTION_CATEGORY_ERROR                        AdnetCode = -1902
	EXCEPTION_IV_ORIENTATION_INVALIDATE             AdnetCode = -1903
	EXCEPTION_IV_RECALLNET_INVALIDATE               AdnetCode = -1904
	EXCEPTION_CAP_BLOCK                             AdnetCode = -1905
	EXCEPTION_TC_FILTER_PLAND                       AdnetCode = -1906
	EXCEPTION_TC_FILTER_REDUPLI                     AdnetCode = -1907
	EXCEPTION_RATELIMIT                             AdnetCode = 503
	EXCEPTION_404_NOT_FOUND                         AdnetCode = 404
	EXCEPTION_OS_VER_LOWER                          AdnetCode = -1908
	EXCEPTION_TC_FILTER_BLACK_PKG                   AdnetCode = -1909
	EXCEPTION_TC_FILTER_BY_TYPE                     AdnetCode = -1910
	EXCEPTION_TC_FILTER_BY_ISTC                     AdnetCode = -1911
	EXCEPTION_TC_FILTER_BY_BAD_REQUEST              AdnetCode = -1912
	EXCEPTION_TC_FILTER_BY_QCC_BID                  AdnetCode = -1913
	EXCEPTION_FILTER_BY_ILLEGAL_GP_SDK_VERSION      AdnetCode = -1914
	EXCEPTION_FILTER_BY_PLACEMENTID_INCONSISTENT    AdnetCode = -1915
	EXCEPTION_TC_FILTER_BY_CITY_OR_COUNTRY          AdnetCode = -1916
	EXCEPTION_TC_FILTER_BY_APP_RATE                 AdnetCode = -1917
	EXCEPTION_OP_FORBIDDEN                          AdnetCode = -1918
	ExCEPTION_OP_SIGN_CHECK_ERROR                   AdnetCode = -1919
	EXCEPTION_TC_FILTER_BY_QCC_QF_NO_BID            AdnetCode = -1920
	EXCEPTION_TC_FILTER_BY_TC_TYPE                  AdnetCode = -1921
)

type Msg struct {
	Code int    `json:"status"`
	Msg  string `json:"msg"`
}

func (code AdnetCode) String() string {
	switch code {
	case MESSAGE_SUCCESS:
		return "ok"
	case MESSAGE_SUCCESS_OLD:
		return "ok"
	case EXCEPTION_RETURN_EMPTY:
		return "EXCEPTION_RETURN_EMPTY"
	case EXCEPTION_RETURN_EMPTY_OLD:
		return "EXCEPTION_RETURN_EMPTY"
	case EXCEPTION_PARAMS_EMPTY:
		return "EXCEPTION_PARAMS_EMPTY"
	case EXCEPTION_PARAMS_ERROR:
		return "EXCEPTION_PARAMS_ERROR"
	case EXCEPTION_TIMEOUT:
		return "EXCEPTION_TIMEOUT"
	case EXCEPTION_SIGN_ERROR:
		return "EXCEPTION_SIGN_ERROR"
	case EXCEPTION_COUNTRY_NOT_ALLOW:
		return "EXCEPTION_COUNTRY_NOT_ALLOW"
	case EXCEPTION_TOKEN_FORMAT_ERROR:
		return "EXCEPTION_TOKEN_FORMAT_ERROR"
	case EXCEPTION_DOMAIN_ERROR:
		return "EXCEPTION_DOMAIN_ERROR"
	case EXCEPTION_NETWORK_ERROR:
		return "EXCEPTION_NETWORK_ERROR"
	case EXCEPTION_IP_NOT_ALLOW:
		return "EXCEPTION_IP_NOT_ALLOW"
	case EXCEPTION_UNIT_ID_EMPTY:
		return "EXCEPTION_UNIT_ID_EMPTY"
	case EXCEPTION_UNIT_NOT_FOUND:
		return "EXCEPTION_UNIT_NOT_FOUND"
	case EXCEPTION_UNIT_NOT_FOUND_IN_APP:
		return "EXCEPTION_UNIT_NOT_FOUND_IN_APP"
	case EXCEPTION_UNIT_NOT_APPWALL:
		return "EXCEPTION_UNIT_NOT_APPWALL"
	case EXCEPTION_UNIT_ADTYPE_ERROR:
		return "EXCEPTION_UNIT_ADTYPE_ERROR"
	case EXCEPTION_UNIT_NOT_ACTIVE:
		return "EXCEPTION_UNIT_NOT_ACTIVE"
	case EXCEPTION_UNIT_NO_CONFIG_MANAGE_REVENUE:
		return "EXCEPTION_UNIT_NO_CONFIG_MANAGE_REVENUE"
	case EXCEPTION_UNIT_BIDDING_TYPE_ERROR:
		return "EXCEPTION_UNIT_BIDDING_TYPE_ERROR"
	case EXCEPTION_APP_ID_EMPTY:
		return "EXCEPTION_APP_ID_EMPTY"
	case EXCEPTION_APP_NOT_FOUND:
		return "EXCEPTION_APP_NOT_FOUND"
	case EXCEPTION_APP_BANNED:
		return "EXCEPTION_APP_BANNED"
	case EXCEPTION_PUBLISHER_BANNED:
		return "EXCEPTION_PUBLISHER_BANNED"
	case EXCEPTION_APP_PLATFORM_ERROR:
		return "EXCEPTION_APP_PLATFORM_ERROR"
	case EXCEPTION_PUBLISHER_ID_EMPTY:
		return "EXCEPTION_PUBLISHER_ID_EMPTY"
	case EXCEPTION_PUBLISHER_NOT_FOUND:
		return "EXCEPTION_PUBLISHER_NOT_FOUND"
	case EXCEPTION_CAMPAIGN_NOT_FOUND:
		return "EXCEPTION_CAMPAIGN_NOT_FOUND"
	case EXCEPTION_CAMPAIGN_ID_EMPTY:
		return "EXCEPTION_CAMPAIGN_ID_EMPTY"
	case EXCEPTION_CAMPAIGN_NOT_ACTIVE:
		return "EXCEPTION_CAMPAIGN_NOT_ACTIVE"
	case EXCEPTION_REQUEST_ID_EMPTY:
		return "EXCEPTION_REQUEST_ID_EMPTY"
	case EXCEPTION_SERVICE_REQUEST_OS_VERSION_REQUIRED:
		return "EXCEPTION_SERVICE_REQUEST_OS_VERSION_REQUIRED"
	case EXCEPTION_SERVICE_REQUEST_COUNTRY_CODE_REQUIRED:
		return "EXCEPTION_SERVICE_REQUEST_COUNTRY_CODE_REQUIRED"
	case EXCEPTION_SERVICE_REQUEST_AD_SOURCE_CLOSED:
		return "EXCEPTION_SERVICE_REQUEST_AD_SOURCE_CLOSED"
	case EXCEPTION_GAID_EMPTY:
		return "EXCEPTION_GAID_EMPTY"
	case EXCEPTION_IDFA_EMPTY:
		return "EXCEPTION_IDFA_EMPTY"
	case EXCEPTION_IMEI_EMPTY:
		return "EXCEPTION_IMEI_EMPTY"
	case EXCEPTION_BRAND_EMPTY:
		return "EXCEPTION_BRAND_EMPTY"
	case EXCEPTION_MODEL_EMPTY:
		return "EXCEPTION_MODEL_EMPTY"
	case EXCEPTION_ANDROIDID_EMPTY:
		return "EXCEPTION_ANDROIDID_EMPTY"
	case EXCEPTION_NETWORK_TYPE_EMPTY:
		return "EXCEPTION_NETWORK_TYPE_EMPTY"
	case EXCEPTION_BANNER_UNIT_SIZE_EMPTY:
		return "EXCEPTION_BANNER_UNIT_SIZE_EMPTY"
	case EXCEPTION_ADNUM_SET_NONE:
		return "EXCEPTION_ADNUM_SET_NONE"
	case EXCEPTION_CATEGORY_ERROR:
		return "EXCEPTION_CATEGORY_ERROR"
	case EXCEPTION_IV_ORIENTATION_INVALIDATE:
		return "EXCEPTION_IV_ORIENTATION_INVALIDATE"
	case EXCEPTION_IV_RECALLNET_INVALIDATE:
		return "EXCEPTION_IV_RECALLNET_INVALIDATE"
	case EXCEPTION_CAP_BLOCK:
		return "EXCEPTION_CAP_BLOCK"
	case EXCEPTION_RATELIMIT:
		return "EXCEPTION_RATELIMIT"
	case EXCEPTION_404_NOT_FOUND:
		return "EXCEPTION_404_NOT_FOUND"
	case EXCEPTION_TC_FILTER_PLAND:
		return "EXCEPTION_TC_FILTER_PLAND"
	case EXCEPTION_TC_FILTER_REDUPLI:
		return "EXCEPTION_TC_FILTER_REDUPLI"
	case EXCEPTION_OS_VER_LOWER:
		return "EXCEPTION_OS_VER_LOWER"
	case EXCEPTION_TC_FILTER_BLACK_PKG:
		return "EXCEPTION_TC_FILTER_BLACK_PKG"
	case EXCEPTION_TC_FILTER_BY_TYPE:
		return "EXCEPTION_TC_FILTER_BY_TYPE"
	case EXCEPTION_TC_FILTER_BY_ISTC:
		return "EXCEPTION_TC_FILTER_BY_ISTC"
	case EXCEPTION_TC_FILTER_BY_BAD_REQUEST:
		return "EXCEPTION_TC_FILTER_BY_BAD_REQUEST"
	case EXCEPTION_TC_FILTER_BY_QCC_BID:
		return "EXCEPTION_TC_FILTER_BY_QCC_BID"
	case EXCEPTION_FILTER_BY_ILLEGAL_GP_SDK_VERSION:
		return "EXCEPTION_FILTER_BY_ILLEGAL_GP_SDK_VERSION"
	case EXCEPTION_FILTER_BY_PLACEMENTID_INCONSISTENT:
		return "EXCEPTION_FILTER_BY_PLACEMENTID_INCONSISTENT"
	case EXCEPTION_TC_FILTER_BY_CITY_OR_COUNTRY:
		return "EXCEPTION_TC_FILTER_BY_CITY_OR_COUNTRY"
	case EXCEPTION_TC_FILTER_BY_APP_RATE:
		return "EXCEPTION_TC_FILTER_BY_APP_RATE"
	case EXCEPTION_OP_FORBIDDEN:
		return "EXCEPTION_LOG_FORBIDDEN"
	case ExCEPTION_OP_SIGN_CHECK_ERROR:
		return "ExCEPTION_LOG_ERROR"
	case EXCEPTION_TC_FILTER_BY_QCC_QF_NO_BID:
		return "EXCEPTION_TC_FILTER_BY_QCC_QF_NO_BID"
	case EXCEPTION_TC_FILTER_BY_TC_TYPE:
		return "EXCEPTION_TC_FILTER_BY_TC_TYPE"
	default:
		return "unset"
	}
}

func (a AdnetCode) Message() string {
	var msg Msg
	msg.Code = int(a)
	msg.Msg = a.String()
	rJson, err := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(&msg)
	if err != nil {
		return "json format error"
	}
	return string(rJson)
}

func (a AdnetCode) Error() string {
	return a.Message()
}

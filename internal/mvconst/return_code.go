package mvconst

const (
	MESSAGE_SUCCESS                                 = 1
	MESSAGE_SUCCESS_OLD                             = 200
	EXCEPTION_RETURN_EMPTY                          = -1
	EXCEPTION_RETURN_EMPTY_OLD                      = 300
	EXCEPTION_PARAMS_EMPTY                          = -2
	EXCEPTION_PARAMS_ERROR                          = -3
	EXCEPTION_TIMEOUT                               = -9
	EXCEPTION_SIGN_ERROR                            = -10
	EXCEPTION_COUNTRY_NOT_ALLOW                     = -11
	EXCEPTION_TOKEN_FORMAT_ERROR                    = -12
	EXCEPTION_IP_NOT_ALLOW                          = -2107
	EXCEPTION_UNIT_ID_EMPTY                         = -1202
	EXCEPTION_UNIT_NOT_FOUND                        = -1201
	EXCEPTION_UNIT_NOT_FOUND_IN_APP                 = -1203
	EXCEPTION_UNIT_NOT_APPWALL                      = -1204
	EXCEPTION_UNIT_ADTYPE_ERROR                     = -1205
	EXCEPTION_UNIT_NOT_ACTIVE                       = -1206
	EXCEPTION_UNIT_NO_CONFIG_MANAGE_REVENUE         = -1207
	EXCEPTION_APP_ID_EMPTY                          = -1301
	EXCEPTION_APP_NOT_FOUND                         = -1302
	EXCEPTION_APP_BANNED                            = -1303
	EXCEPTION_PUBLISHER_BANNED                      = -1304
	EXCEPTION_APP_PLATFORM_ERROR                    = -1305
	EXCEPTION_PUBLISHER_ID_EMPTY                    = -1306
	EXCEPTION_PUBLISHER_NOT_FOUND                   = -1307
	EXCEPTION_CAMPAIGN_NOT_FOUND                    = -1401
	EXCEPTION_CAMPAIGN_ID_EMPTY                     = -1402
	EXCEPTION_CAMPAIGN_NOT_ACTIVE                   = -1403
	EXCEPTION_REQUEST_ID_EMPTY                      = -1501
	EXCEPTION_SERVICE_REQUEST_OS_VERSION_REQUIRED   = -2102
	EXCEPTION_SERVICE_REQUEST_COUNTRY_CODE_REQUIRED = -2103
	EXCEPTION_SERVICE_REQUEST_AD_SOURCE_CLOSED      = -2104
	EXCEPTION_GAID_EMPTY                            = -1801
	EXCEPTION_IDFA_EMPTY                            = -1802
	EXCEPTION_ADNUM_SET_NONE                        = -1901
	EXCEPTION_CATEGORY_ERROR                        = -1902
	EXCEPTION_IV_ORIENTATION_INVALIDATE             = -1903
	EXCEPTION_IV_RECALLNET_INVALIDATE               = -1904
	EXCEPTION_CAP_BLOCK                             = -1905
	EXCEPTION_RATELIMIT                             = 503
	EXCEPTION_404_NOT_FOUND                         = 404

	// exception 特殊result需要返回
	EXCEPTION_SPECIAL_RESULT_RETURN          = "SPECIAL_RESULT_RETURN"
	EXCEPTION_SPECIAL_RESULT_RETURN_NOENCODE = "SPECIAL_RESULT_RETURN_NOENCODE"
)

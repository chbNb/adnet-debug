package mvconst

const (
	BackendOK                    = 2000
	ParamInvalidate              = 4001
	ConstructReqDataError        = 4002
	BuildReqError                = 4003
	HTTPDoError                  = 4004
	HTTPStatusNotOK              = 4005
	HTTPReadBodyError            = 4006
	ThriftTSocketError           = 4007
	ThriftTransportOpenError     = 4008
	BackendLowFlowFilter         = 4009
	BackendFallback              = 4010
	ThriftClientError            = 4011
	BackendDoError               = 4012
	BackendReqTypeAABTestFilter  = 4013
	BackendLowVersionFilter      = 4014 // 不能切量pioneer的流量
	BackendNoAds                 = 5007
	BackendNotInstance           = 6001
	BackendRegionFilter          = 6002
	BackendContentFilter         = 6003
	ImpCapBlock                  = 6004
	BackendQueryFailureException = 6005
	BackendReadTimeout           = 6006
	BackendNetTemporary          = 6007
	BackendUnknownError          = 6008
	BackendFillError             = 6009
	BackendDeviceFilter          = 6010
	BackendTrackingFilter        = 6011
	BackendPackageFilter         = 6012
	BackendAdTypeFilter          = 6013
	BackendPlatformFilter        = 6014
	BackendParserError           = 6015
	BackendSDKVersionFilter      = 6016
	BackendFillRateFilter        = 6017
)

const (
	FilterRequestByImpBlock       = "imp_block"
	FilterRequestByLowFlow        = "low_flow"
	FilterRequestByStack          = "stack_block"
	FilterRequestByResidualMemory = "memory_block"
	FilterRequestByIpBlacklist    = "ip_block"
)

package mvconst

const (
	AdTracking_Click = iota + 1
	AdTracking_Imp
	AdTracking_Both_Click_Imp
)

const (
	AdTrackingAndroidRV                    = 80300
	AdTrackingAndroidNativeVideo           = 80201
	AdTrackingAndroidNativePicImp          = 80600
	AdTrackingAndroidNativePicClick        = 80400
	AdTrackingAndroidIV                    = 81000
	AdTrackingIOSRV                        = 20800
	ADTRACKINGIOSADX                       = 20600
	AdTrackingIOSNativeImp                 = 30100
	AdTrackingIOSNativeClick               = 20800
	AdTrackingIOSIV                        = 30600
	AdTrackingIOSStorekit                  = 30200
	AdTrackingIOSExpV1                     = 30900
	AdTrackingIOSExpV2                     = 30901
	MoreOfferBlock                         = 90300
	IOSSupportVideoUrlWithParams           = 40100
	AndroidSupDeliverClickUrl              = 80127
	IOSSupportVideoSizeZero                = 40900 // 低于这个版本的sdk强制需要videosize
	IOSSupportVideoUrlWithParamsExpV1      = 50805
	IOSSupportVideoUrlWithParamsExpV2      = 50806
	IOSSupportVideoUrlWithParamsExpV3      = 50807
	IOSSupportVideoUrlWithParamsExpV4      = 50808
	FilterIllegalGPSdkVersion              = 100109
	FilterIllegalGPSdkVersionAndroidXMin   = 120000
	FilterIllegalGPSdkVersionAndroidXMax   = 120101
	IOSSupportTransferIDFVOpenIDFAInParamC = 60307
	AndroidSupportApkInfoVersion           = 150530
	IosSupportOfferRewardPlusVersion       = 60903
	ToponUnSupportBigTempalteVersion       = 60904
	AndroidCorectNetworkType               = 150601
)

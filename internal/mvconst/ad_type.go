package mvconst

const (
	ADTypeUnknown = iota
	ADTypeText
	ADTypeBanner
	ADTypeAppwall
	ADTypeOverlay
	ADTypeFullScreen             //全屏
	ADTypeInterstitial      = 29 //插屏
	ADTypeNative            = 42 //原生
	ADTypeNativeVideo       = 43 //原生视频
	ADTypeNativePic         = 44 //原生图片
	ADTypeRewardVideo       = 94 //激进视频
	ADTypeFeedsVideo        = 95
	ADTypeOfferWall         = 278
	ADTypeInterstitialSdk   = 279 // 与MADX交互时的插屏是这个， 不是29
	ADTypeOnlineVideo       = 284
	ADTypeJSNativeVideo     = 285
	ADTypeJSBannerVideo     = 286
	ADTypeInterstitialVideo = 287
	ADTypeInteractive       = 288 // IA
	ADTypeJMIcon            = 289
	ADTypeWXNative          = 291
	ADTypeWXAppwall         = 292
	ADTypeWXBanner          = 293
	ADTypeWXRewardImg       = 294
	ADTypeMoreOffer         = 295
	ADTypeSdkBanner         = 296
	ADTypeSplash            = 297
	ADTypeNativeH5          = 298
)

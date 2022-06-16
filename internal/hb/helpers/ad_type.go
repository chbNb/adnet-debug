package helpers

import (
	"gitlab.mobvista.com/ADN/adnet/internal/hb/constant"
)

func AdTypeStr(adType int) string {
	result := "unknown"
	switch adType {
	case constant.Native:
		result = "native"
	case constant.RewardVideo:
		result = "rewarded_video"
	case constant.InterstitialVideo:
		result = "interstitial_video"
	case constant.Banner:
		result = "banner"
	case constant.SplashAd:
		result = "splash"
	}
	return result
}

func GenAdTypeStr(adType int, videoUrl string) string {
	paramType := ""
	if adType == constant.Native && len(videoUrl) > 0 {
		paramType = constant.NativeVideoStr
	} else if adType == constant.RewardVideo {
		paramType = constant.RewardVideoStr
	} else if adType == constant.InterstitialVideo {
		paramType = constant.InterstitialVideoStr
	} else if adType == constant.Banner {
		paramType = constant.BannerStr
	} else if adType == constant.SplashAd {
		paramType = constant.SplashAdStr
	}
	return paramType
}

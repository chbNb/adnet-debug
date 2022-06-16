package helpers

import (
	"strings"

	"gitlab.mobvista.com/ADN/adnet/internal/hb/constant"
)

func GetPlatform(os string) int {
	switch strings.ToLower(os) {
	case constant.PlatformIOS:
		return constant.IOS
	case constant.PlatformAndroid:
		return constant.Android
	case constant.PlatformFireOs:
		return constant.Android
	default:
		return constant.Other
	}
}

func GetOs(platform int) string {
	switch platform {
	case constant.IOS:
		return constant.PlatformIOS
	case constant.Android:
		return constant.PlatformAndroid
	default:
		return constant.PlatformOther
	}
}

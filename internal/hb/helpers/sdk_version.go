package helpers

import (
	"strings"

	"gitlab.mobvista.com/ADN/adnet/internal/hb/constant"
	model "gitlab.mobvista.com/ADN/adnet/internal/mvutil"
)

func RenderSdkVersion(sdkversion string) *model.SDKVersionItem {
	var item model.SDKVersionItem
	if !strings.Contains(sdkversion, "_") {
		item.SDKNumber = sdkversion
	} else {
		arr := strings.Split(sdkversion, "_")
		if len(arr[0]) > 0 {
			item.SDKType = arr[0]
		}
		if len(arr[1]) > 0 {
			item.SDKNumber = arr[1]
		}
	}

	if InStrArray(item.SDKType, constant.SdkVersionPrefix) {
	}

	code, err := VersionCode(item.SDKNumber)
	if err == nil {
		item.SDKVersionCode = code
	}
	return &item
}

package output

import (
	supply_mvutil "gitlab.mobvista.com/ADN/chasm/module/supply/mvutil"
	"strings"

	"gitlab.mobvista.com/ADN/adnet/internal/extractor"
	"gitlab.mobvista.com/ADN/adnet/internal/mvconst"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
)

func IsReturnEndcard(r *mvutil.RequestParams) bool {
	// os_version 不为空
	if len(r.Param.OSVersion) <= 0 {
		return false
	}
	// 不能是ipad
	if strings.Contains(strings.ToLower(r.Param.Model), "ipad") {
		return false
	}
	// 只能Android或者IOS
	if r.Param.Platform != mvconst.PlatformAndroid && r.Param.Platform != mvconst.PlatformIOS {
		return false
	}
	// 获取confs
	versionCompare, ifFind := extractor.GetVersionCompare()
	if !ifFind {
		return false
	}
	confs := versionCompare["OS_VERSION_ENDCARD"]
	var conf mvutil.VersionCompareItem
	if r.Param.Platform == mvconst.PlatformAndroid {
		conf = confs.Android
	} else {
		conf = confs.IOS
	}
	if conf.Version <= int32(0) {
		return false
	}
	if r.Param.OSVersionCode > conf.Version && !mvutil.InInt32Arr(r.Param.OSVersionCode, conf.ExcludeVersion) {
		return true
	}
	return false
}

func Compare(r *mvutil.RequestParams, compareType string) bool {
	if len(r.Param.SDKVersion) <= 0 {
		return false
	}
	versionItem := supply_mvutil.RenderSDKVersion(r.Param.SDKVersion)
	allowType := []string{"mi", "mal"}
	if !mvutil.InStrArray(versionItem.SDKType, allowType) {
		return false
	}
	// 只能Android或者IOS
	if r.Param.Platform != mvconst.PlatformAndroid && r.Param.Platform != mvconst.PlatformIOS {
		return false
	}
	// 获取confs
	versionCompare, ifFind := extractor.GetVersionCompare()
	if !ifFind {
		return false
	}
	confs, ok := versionCompare[compareType]
	if !ok {
		return false
	}
	var conf mvutil.VersionCompareItem
	if r.Param.Platform == mvconst.PlatformAndroid {
		conf = confs.Android
	} else {
		conf = confs.IOS
	}
	if conf.Version <= int32(0) {
		return false
	}
	versionCode := mvutil.GetVersionCode(versionItem.SDKNumber)
	if versionCode > conf.Version && !mvutil.InInt32Arr(versionCode, conf.ExcludeVersion) {
		return true
	}
	return false
}

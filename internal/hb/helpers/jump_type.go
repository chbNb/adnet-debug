package helpers

func JumpTypeVal(key string, jumpTypeData map[string]int32) int32 {
	if jumpVal, ok := jumpTypeData[key]; ok {
		return jumpVal
	}
	return int32(0)
}

func CompareSDKVersionForJumpType(jumpType string, sdkVersion string, jumpTypeData map[string]map[string]string) bool {
	jumpTypeCfg, ok := jumpTypeData[jumpType]
	if !ok {
		return false
	}

	sdkversionItem := RenderSdkVersion(sdkVersion)
	if len(sdkversionItem.SDKType) <= 0 {
		return false
	}
	confVersion, ok := jumpTypeCfg[sdkversionItem.SDKType]
	if !ok {
		return false
	}
	confVersionCode, err := VersionCode(confVersion)
	if err != nil {
		return false
	}
	return sdkversionItem.SDKVersionCode >= confVersionCode
}

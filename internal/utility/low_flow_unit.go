package utility

import (
	"math/rand"
	"strconv"

	"gitlab.mobvista.com/ADN/adnet/internal/extractor"
)

// 低效流量过滤
func IsLowFlowUnit(unitId, appId, publisherId int64, countryCode string) bool {
	lowUnits, _ := extractor.GetLOW_FLOW_UNITS()
	if len(lowUnits) == 0 {
		return false
	}

	// 按照概率是否在adnet降填充
	unitStr := strconv.FormatInt(unitId, 10)
	unitCCKey := unitStr + "_" + countryCode
	rate, ok := lowUnits[unitCCKey]
	if ok {
		randInt := rand.Intn(100)
		return randInt >= rate
	}

	rate, ok = lowUnits[unitStr]
	if ok {
		randInt := rand.Intn(100)
		return randInt >= rate
	}

	appStr := "app_" + strconv.FormatInt(appId, 10)
	appCCKey := appStr + "_" + countryCode

	rate, ok = lowUnits[appCCKey]
	if ok {
		randInt := rand.Intn(100)
		return randInt >= rate
	}

	rate, ok = lowUnits[appStr]
	if ok {
		randInt := rand.Intn(100)
		return randInt >= rate
	}

	pubStr := "pub_" + strconv.FormatInt(publisherId, 10)
	pubCCKey := pubStr + "_" + countryCode

	rate, ok = lowUnits[pubCCKey]
	if ok {
		randInt := rand.Intn(100)
		return randInt >= rate
	}

	rate, ok = lowUnits[pubStr]
	if ok {
		randInt := rand.Intn(100)
		return randInt >= rate
	}

	return false
}

func IsLowFlowAdType(adType int32, platform int, countryCode string, fixedEcpm float64) bool {
	lowAdTypes, ifFind := extractor.GetLOW_FLOW_ADTYPE()
	if !ifFind || len(lowAdTypes) == 0 {
		return false
	}

	// key: adtype_platform_cc
	adTypeKey := strconv.FormatInt(int64(adType), 10) + "_" + strconv.Itoa(platform) + "_" + countryCode
	configFloor, ok := lowAdTypes[adTypeKey]
	if !ok {
		return false
	}

	if fixedEcpm > configFloor {
		return true
	}
	return false
}

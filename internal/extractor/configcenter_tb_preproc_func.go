package extractor

import (
	"errors"
	smodel "gitlab.mobvista.com/ADN/structs/model"
	"gitlab.mobvista.com/ADN/treasure_box_sdk/tb_tools"
)

//func initCCObj() {
//	ccObj = new(ConfigCenterObj)
//}

func ConfigCenterPreProc(i interface{}) error {
	ptr, ok := i.(*tb_tools.ExtractorInterface)
	if !ok {
		return errors.New("ConfigCenterPreProc failed: type cast to *tb_tools.ExtractorInterface")
	}
	if configCenter, ok := (*ptr).(*smodel.ConfigCenter); ok {
		cc, err := configCenterUpdateFunc(configCenter)
		if err != nil {
			logger.Errorf("ConfigCenterPreProc error: %v, key: %v", err.Error(), configCenter.Key)
		}
		if cc == nil {
			*ptr = nil
		} else {
			*ptr = cc
		}
		return err
	} else {
		return errors.New("ConfigCenterPreProc failed: type cast to *smodel.ConfigCenter failed")
	}
}

func configCenterUpdateFunc(cc *smodel.ConfigCenter) (*smodel.ConfigCenter, error) {
	if cc == nil {
		return nil, errors.New("configCenterUpdateFunc: nil ConfigCenter")
	}
	for k, v := range cc.Value {
		switch k {
		case "DOMAIN_TRACK":
			cc.Value[k] = getDOMAIN_TRACK(v)
		case "DOMAIN_TRACKS":
			cc.Value[k] = getDOMAIN_TRACKS(v)
		case "DOMAIN":
			cc.Value[k] = getDOMAIN(v)
		case "SYSTEM":
			cc.Value[k] = getSYSTEM(v)
		case "SYSTEM_AREA":
			cc.Value[k] = getSYSTEM_AREA(v)
		case "MP_DOMAIN_CONF":
			cc.Value[k] = getMP_DOMAIN_CONF(v)
		case "JSSDK_DOMAIN_TRACK":
			cc.Value[k] = getJSSDK_DOMAIN_TRACK(v)
		case "CHET_DOMAIN":
			cc.Value[k] = getCHET_DOMAIN(v)
		case "CLOUD_NAME":
			cc.Value[k] = getCLOUD_NAME(v)
		case "RateLimit":
			cc.Value[k] = getRateLimit(v)
		}
	}
	return cc, nil
}

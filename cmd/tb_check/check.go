package main

import (
	"fmt"
	"github.com/pkg/errors"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	"gitlab.mobvista.com/ADN/lego/schema"
	smodel "gitlab.mobvista.com/ADN/structs/model"
	"reflect"
	"strings"
)

func main() {
	var (
		newMap, oldMap map[string]schema.Field
	)

	newMap = getField(smodel.PublisherInfo{})
	oldMap = getField(mvutil.PublisherInfo{})
	if key, err := check(oldMap, newMap, "PublisherInfo", []string{}); err != nil {
		fmt.Println(key, err.Error())
	} else {
		fmt.Println("PublisherInfo is ok!")
	}
	//反向检测 PublisherInfo
	if key, err := check(newMap, oldMap, "PublisherInfo[2]", []string{
		//@todo
		"PublisherInfo[2].Id",
	}); err != nil {
		fmt.Println(key, err.Error())
	} else {
		fmt.Println("【反】PublisherInfo is ok!")
	}

	//appInfo
	newMap = getField(smodel.AppInfo{})
	oldMap = getField(mvutil.AppInfo{})
	if key, err := check(oldMap, newMap, "AppInfo", []string{}); err != nil {
		fmt.Println(key, err.Error())
	} else {
		fmt.Println("AppInfo is ok!")
	}
	//反向检测 AppInfo
	if key, err := check(newMap, oldMap, "AppInfo[2]", []string{
		//@todo
		"AppInfo[2].PostbackM",
		"AppInfo[2].RuduceRule",
		"AppInfo[2].App.Desc",
		"AppInfo[2].SspCampaignIds",
		"AppInfo[2].Id",
		"AppInfo[2].App.Postback",
		"AppInfo[2].ReduceRuleV2",
		"AppInfo[2].PublisherId",
		"AppInfo[2].App.PublisherID",
		"AppInfo[2].BlendTraffic",
		"AppInfo[2].App.EnableHb",
		"AppInfo[2].BtCampaignIds",
		"AppInfo[2].App.PostbackURL",
	}); err != nil {
		fmt.Println(key, err.Error())
	} else {
		fmt.Println("【反】AppInfo is ok!")
	}

	// unitInfo
	newMap = getField(smodel.UnitInfo{})
	oldMap = getField(mvutil.UnitInfo{})
	if key, err := check(oldMap, newMap, "UnitInfo", []string{}); err != nil {
		fmt.Println(key, err.Error())
	} else {
		fmt.Println("UnitInfo is ok!")
	}
	//反向检测 UnitInfo
	if key, err := check(newMap, oldMap, "UnitInfo[2]", []string{
		"UnitInfo[2].Unit.PublisherId",
		"UnitInfo[2].DeductConfig",
		"UnitInfo[2].SubsidyRule",
		"UnitInfo[2].SecurityKey",
		"UnitInfo[2].Id", //  这字段需要确认是否有模块使用，没有需要删掉
		"UnitInfo[2].PublisherId",
		"UnitInfo[2].FakeRule",
		"UnitInfo[2].CallbackURL",
		"UnitInfo[2].Unit.AppId",
		"UnitInfo[2].FakePriceRuleV2",
	}); err != nil {
		fmt.Println(key, err.Error())
	} else {
		fmt.Println("【反】UnitInfo is ok!")
	}

	//AdvertiserInfo
	newMap = getField(smodel.AdvertiserInfo{})
	oldMap = getField(mvutil.AdvertiserInfo{})
	if key, err := check(oldMap, newMap, "AdvertiserInfo", []string{}); err != nil {
		fmt.Println(key, err.Error())
	} else {
		fmt.Println("AdvertiserInfo is ok!")
	}
	//反向检测 UnitInfo
	if key, err := check(newMap, oldMap, "AdvertiserInfo[2]", []string{}); err != nil {
		fmt.Println(key, err.Error())
	} else {
		fmt.Println("【反】AdvertiserInfo is ok!")
	}

	//AdxTrafficMediaConfig
	newMap = getField(smodel.AdxTrafficMediaConfig{})
	oldMap = getField(mvutil.AdxTrafficMediaConfig{})
	if key, err := check(oldMap, newMap, "AdxTrafficMediaConfig", []string{}); err != nil {
		fmt.Println(key, err.Error())
	} else {
		fmt.Println("AdxTrafficMediaConfig is ok!")
	}
	//反向检测 UnitInfo
	if key, err := check(newMap, oldMap, "AdxTrafficMediaConfig[2]", []string{}); err != nil {
		fmt.Println(key, err.Error())
	} else {
		fmt.Println("【反】AdxTrafficMediaConfig is ok!")
	}

	//AdxDspConfig
	newMap = getField(smodel.AdxDspConfig{})
	oldMap = getField(mvutil.AdxDspConfig{})
	if key, err := check(oldMap, newMap, "AdxDspConfig", []string{}); err != nil {
		fmt.Println(key, err.Error())
	} else {
		fmt.Println("AdxDspConfig is ok!")
	}
	//反向检测 UnitInfo
	if key, err := check(newMap, oldMap, "AdxDspConfig[2]", []string{}); err != nil {
		fmt.Println(key, err.Error())
	} else {
		fmt.Println("【反】AdxDspConfig is ok!")
	}

	//AppPackageMtgID
	newMap = getField(smodel.AppPackageMtgID{})
	oldMap = getField(mvutil.AppPackageMtgID{})
	if key, err := check(oldMap, newMap, "AppPackageMtgID", []string{}); err != nil {
		fmt.Println(key, err.Error())
	} else {
		fmt.Println("AppPackageMtgID is ok!")
	}
	//反向检测 AppPackageMtgID
	if key, err := check(newMap, oldMap, "AppPackageMtgID[2]", []string{}); err != nil {
		fmt.Println(key, err.Error())
	} else {
		fmt.Println("【反】AppPackageMtgID is ok!")
	}

	//ConfigCenter
	newMap = getField(smodel.ConfigCenter{})
	oldMap = getField(mvutil.ConfigCenter{})
	if key, err := check(oldMap, newMap, "ConfigCenter", []string{}); err != nil {
		fmt.Println(key, err.Error())
	} else {
		fmt.Println("ConfigCenter is ok!")
	}
	//反向检测 ConfigCenter
	if key, err := check(newMap, oldMap, "ConfigCenter[2]", []string{}); err != nil {
		fmt.Println(key, err.Error())
	} else {
		fmt.Println("【反】ConfigCenter is ok!")
	}

	//ConfigAlgorithmFillRate
	newMap = getField(smodel.ConfigAlgorithmFillRate{})
	oldMap = getField(mvutil.ConfigAlgorithmFillRate{})
	if key, err := check(oldMap, newMap, "ConfigAlgorithmFillRate", []string{}); err != nil {
		fmt.Println(key, err.Error())
	} else {
		fmt.Println("ConfigAlgorithmFillRate is ok!")
	}
	//反向检测 ConfigAlgorithmFillRate
	if key, err := check(newMap, oldMap, "ConfigAlgorithmFillRate[2]", []string{}); err != nil {
		fmt.Println(key, err.Error())
	} else {
		fmt.Println("【反】ConfigAlgorithmFillRate is ok!")
	}

	//FreqControlFactor
	newMap = getField(smodel.FreqControlFactor{})
	oldMap = getField(mvutil.FreqControlFactor{})
	if key, err := check(oldMap, newMap, "FreqControlFactor", []string{}); err != nil {
		fmt.Println(key, err.Error())
	} else {
		fmt.Println("FreqControlFactor is ok!")
	}
	//反向检测 FreqControlFactor
	if key, err := check(newMap, oldMap, "FreqControlFactor[2]", []string{}); err != nil {
		fmt.Println(key, err.Error())
	} else {
		fmt.Println("【反】FreqControlFactor is ok!")
	}

	//MVConfig
	newMap = getField(smodel.MVConfig{})
	oldMap = getField(mvutil.MVConfig{})
	if key, err := check(oldMap, newMap, "MVConfig", []string{}); err != nil {
		fmt.Println(key, err.Error())
	} else {
		fmt.Println("MVConfig is ok!")
	}
	//反向检测 MVConfig
	if key, err := check(newMap, oldMap, "MVConfig[2]", []string{}); err != nil {
		fmt.Println(key, err.Error())
	} else {
		fmt.Println("【反】MVConfig is ok!")
	}

	//CheckResult
	newMap = getField(smodel.CheckResult{})
	oldMap = getField(mvutil.CheckResult{})
	if key, err := check(oldMap, newMap, "CheckResult", []string{}); err != nil {
		fmt.Println(key, err.Error())
	} else {
		fmt.Println("CheckResult is ok!")
	}
	//反向检测 CheckResult
	if key, err := check(newMap, oldMap, "CheckResult[2]", []string{}); err != nil {
		fmt.Println(key, err.Error())
	} else {
		fmt.Println("【反】CheckResult is ok!")
	}

	//PlacementInfo
	newMap = getField(smodel.PlacementInfo{})
	oldMap = getField(mvutil.PlacementInfo{})
	if key, err := check(oldMap, newMap, "PlacementInfo", []string{}); err != nil {
		fmt.Println(key, err.Error())
	} else {
		fmt.Println("PlacementInfo is ok!")
	}
	//反向检测 PlacementInfo
	if key, err := check(newMap, oldMap, "PlacementInfo[2]", []string{}); err != nil {
		fmt.Println(key, err.Error())
	} else {
		fmt.Println("【反】PlacementInfo is ok!")
	}

	//SspProfitDistributionRule
	newMap = getField(smodel.SspProfitDistributionRule{})
	oldMap = getField(mvutil.SspProfitDistributionRule{})
	if key, err := check(oldMap, newMap, "SspProfitDistributionRule", []string{}); err != nil {
		fmt.Println(key, err.Error())
	} else {
		fmt.Println("SspProfitDistributionRule is ok!")
	}
	//反向检测 SspProfitDistributionRule
	if key, err := check(newMap, oldMap, "SspProfitDistributionRule[2]", []string{}); err != nil {
		fmt.Println(key, err.Error())
	} else {
		fmt.Println("【反】SspProfitDistributionRule is ok!")
	}
}

func check(oldMap, newMap map[string]schema.Field, preKey string, ignoreKeys []string) (string, error) {
	var oldType, newType string
	for k, v := range oldMap {
		fullKey := preKey + "." + k
		if mvutil.InStrArray(fullKey, ignoreKeys) {
			continue
		}
		newV, ok := newMap[k]
		if !ok {
			return fullKey, errors.New("Not in new Map")
		} else {
			if v.Kind != newV.Kind {
				return fullKey, errors.New("Kind is err!(old: " + v.Kind + ", new: " + newV.Kind + ")")
			}
			if strings.Contains(v.Type, ".") {
				oldType = strings.Split(v.Type, ".")[1]
			} else {
				oldType = v.Type
			}
			if strings.Contains(newV.Type, ".") {
				newType = strings.Split(newV.Type, ".")[1]
			} else {
				newType = newV.Type
			}

			if oldType != newType {
				return fullKey, errors.New("Type is err!(old: " + oldType + ", new: " + newType + ")")
			}
			//fmt.Println(v.Name, "  ", " TYPE:", oldType, " ", newType, " TAG:", "`", v.Tag, "`  `", newV.Tag, "`")
			if v.Tag != newV.Tag {
				return fullKey, errors.New("Tag is err!(old: " + v.Tag + ", new: " + newV.Tag + ")")
			}
		}
		if len(v.Schema) > 0 {
			if key, err := check(getMap(&v.Schema), getMap(&newV.Schema), fullKey, ignoreKeys); err != nil {
				return key, err
			}
		}
	}
	return "", nil
}

func getField(obj interface{}) (resultMap map[string]schema.Field) {
	//if sc, err := schema.NewSchemaByTag(obj, "bson"); err != nil {
	if sc, err := schema.NewSchema(obj); err != nil {
		fmt.Println(err.Error(), obj)
	} else {
		return getMap(sc)
	}
	return
}

func getMap(sc *schema.Schema) map[string]schema.Field {
	resultMap := make(map[string]schema.Field)
	fields := (*[]schema.Field)(sc)
	for _, v := range *fields {
		resultMap[v.Name] = v
	}
	return resultMap
}

func print(obj interface{}, ignoreArr map[string]int) error {

	v := reflect.ValueOf(obj)
	// 修改值必须是指针类型否则不可行
	if v.Kind() != reflect.Ptr {
		return errors.New("不是指针类型，没法进行修改操作")
	}

	// 获取指针所指向的元素
	v = v.Elem()

	rtype := reflect.TypeOf(obj).Elem()
	for i := 0; i < rtype.NumField(); i++ {
		field := rtype.Field(i)
		if _, ok := ignoreArr[field.Name]; ok {
			continue
		}
		fmt.Printf("\"%s\":1,\n", field.Name)
	}
	return nil
}

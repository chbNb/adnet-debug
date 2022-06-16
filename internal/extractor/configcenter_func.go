package extractor

import (
	"errors"

	"github.com/easierway/concurrent_map"
	jsoniter "github.com/json-iterator/go"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	smodel "gitlab.mobvista.com/ADN/structs/model"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//type ConfigCenterObj struct {
//	TrackingDB         mvutil.TRACKING_DB
//	DomainTrack        string
//	DomainTracks       map[string]mvutil.IRate
//	Domain             string
//	System             string
//	SYSTEM_AREA        string
//	DOMAIN_GO_GRAY     string
//	DOMAIN_SS_PLATFORM string
//	MP_DOMAIN_CONF     []mvutil.MP_DOMAIN_CONF
//	JSSDK_DOMAIN_TRACK string
//	CHET_DOMAIN        string
//	CLOUD_NAME         string
//	RateLimit          *mvutil.RateLimit
//	TC_DELAY_SQSURL    string
//}
//
//var ccObj *ConfigCenterObj

var configcenterUpdateInputError = errors.New("configcenterUpdateFunc failed, query or dbLoaderInfo is nil")

func NewConfigCenterExtractor(allFilter, incFilter, selectOptions bson.M) dbLoaderInfo {
	return dbLoaderInfo{
		collection:   "configcenter",
		querybmDbMap: allFilter,
		queryimDbMap: incFilter,
		querySelect:  selectOptions,
		updateFunc:   configcenterUpdateFunc,
		dataCur:      concurrent_map.CreateConcurrentMap(concurrentMapPartitionsNum),
	}
}

func configcenterUpdateFunc(query *mgo.Query, num int, dbLoaderInfo *dbLoaderInfo) (maxUpdate int64, err error) {
	if query == nil || dbLoaderInfo == nil {
		return 0, configcenterUpdateInputError
	}
	configCenter := &smodel.ConfigCenter{}
	maxUpdate = 0
	item := query.Iter()
	for item.Next(configCenter) {
		configcenterIncUpdateFunc(configCenter, dbLoaderInfo)
		configCenter = &smodel.ConfigCenter{}
	}
	if item.Err() != nil {
		logger.Warnf("configcenterUpdateFunc err: %s", err.Error())
	}
	if err := item.Close(); err != nil {
		logger.Error(err.Error())
	}
	return maxUpdate, nil
}

func configcenterIncUpdateFunc(configCenter *smodel.ConfigCenter, dbLoaderInfo *dbLoaderInfo) {
	//err, _ := configCenterUpdateFunc(configCenter)
	//if err != nil {
	//	logger.Warnf("configcenterIncUpdateFunc err: %s", err.Error())
	//}
	//DbLoaderRegistry[TblConfigCenter].dataCur.Set(concurrent_map.StrKey(configCenter.Key), configCenter)
}

func getConfigCenterValue(key string) (interface{}, error) {
	ccKey := mvutil.Config.AreaConfig.HttpConfig.ConfigCenterKey
	var cc *smodel.ConfigCenter
	var ifFind bool
	cc, ifFind = GetConfigcenter(ccKey)
	if !ifFind {
		return nil, errors.New("key not found")
	}
	value := cc.Value
	if v, ok := value[key]; ok {
		return v, nil
	}
	ccKey = cc.Area
	cc, ifFind = GetConfigcenter(ccKey)
	if !ifFind {
		return nil, errors.New("area not found")
	}
	value = cc.Value
	if v, ok := value[key]; ok {
		return v, nil
	}
	return nil, errors.New("cckey not found")
}

func GetDOMAIN_TRACKS() map[string]mvutil.IRate {
	value, err := getConfigCenterValue("DOMAIN_TRACKS")
	res := make(map[string]mvutil.IRate)
	if err != nil {
		return res
	}
	v, ok := value.(map[string]mvutil.IRate)
	if !ok {
		return res
	}
	return v
}

func getDOMAIN_TRACKS(v interface{}) map[string]mvutil.IRate {
	res := make(map[string]mvutil.IRate)

	tres := make(map[string]*mvutil.TagRate)
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	jsonStr, _ := json.Marshal(v)
	_ = json.Unmarshal(jsonStr, &tres)
	for k, v := range tres {
		res[k] = v
	}
	return res
}

func GetDOMAIN_TRACK() string {
	value, err := getConfigCenterValue("DOMAIN_TRACK")
	res := ""
	if err != nil {
		return res
	}
	v, ok := value.(string)
	if !ok {
		return res
	}
	return v
}

func getDOMAIN_TRACK(v interface{}) string {
	res := ""
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	jsonStr, _ := json.Marshal(v)
	_ = json.Unmarshal(jsonStr, &res)
	return res
}

func GetDOMAIN() string {
	value, err := getConfigCenterValue("DOMAIN")
	res := ""
	if err != nil {
		return res
	}
	v, ok := value.(string)
	if !ok {
		return res
	}
	return v
}

func getDOMAIN(v interface{}) string {
	res := ""
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	jsonStr, _ := json.Marshal(v)
	_ = json.Unmarshal(jsonStr, &res)
	return res
}

func GetSYSTEM() string {
	value, err := getConfigCenterValue("SYSTEM")
	res := ""
	if err != nil {
		return res
	}
	v, ok := value.(string)
	if !ok {
		return res
	}
	return v
}

func getSYSTEM(v interface{}) string {
	res := ""
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	jsonStr, _ := json.Marshal(v)
	_ = json.Unmarshal(jsonStr, &res)
	return res
}

func GetSYSTEM_AREA() string {
	value, err := getConfigCenterValue("SYSTEM_AREA")
	res := ""
	if err != nil {
		return res
	}
	v, ok := value.(string)
	if !ok {
		return res
	}
	return v
}

func getSYSTEM_AREA(v interface{}) string {
	res := ""
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	jsonStr, _ := json.Marshal(v)
	_ = json.Unmarshal(jsonStr, &res)
	return res
}

func GetMP_DOMAIN_CONF() []mvutil.MP_DOMAIN_CONF {
	value, err := getConfigCenterValue("MP_DOMAIN_CONF")
	var res []mvutil.MP_DOMAIN_CONF
	if err != nil {
		return res
	}
	v, ok := value.([]mvutil.MP_DOMAIN_CONF)
	if !ok {
		return res
	}
	return v
}

func getMP_DOMAIN_CONF(v interface{}) []mvutil.MP_DOMAIN_CONF {
	var res []mvutil.MP_DOMAIN_CONF
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	jsonStr, _ := json.Marshal(v)
	_ = json.Unmarshal(jsonStr, &res)
	return res
}

func GetJSSDK_DOMAIN_TRACK() string {
	value, err := getConfigCenterValue("JSSDK_DOMAIN_TRACK")
	var res = ""
	if err != nil {
		return res
	}
	v, ok := value.(string)
	if !ok {
		return res
	}
	return v
}

func getJSSDK_DOMAIN_TRACK(v interface{}) string {
	res := ""
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	jsonStr, _ := json.Marshal(v)
	_ = json.Unmarshal(jsonStr, &res)
	return res
}

func GetCHET_DOMAIN() string {
	value, err := getConfigCenterValue("CHET_DOMAIN")
	var res = ""
	if err != nil {
		return res
	}
	v, ok := value.(string)
	if !ok {
		return res
	}
	return v
}

func getCHET_DOMAIN(v interface{}) string {
	res := ""
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	jsonStr, _ := json.Marshal(v)
	_ = json.Unmarshal(jsonStr, &res)
	return res
}

func GetCLOUD_NAME() string {
	value, err := getConfigCenterValue("CLOUD_NAME")
	res := ""
	if err != nil {
		return res
	}
	v, ok := value.(string)
	if !ok {
		return res
	}
	return v
}
func getCLOUD_NAME(v interface{}) string {
	res := ""
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	jsonStr, _ := json.Marshal(v)
	_ = json.Unmarshal(jsonStr, &res)
	return res
}

func getRateLimit(v interface{}) *mvutil.RateLimit {
	var res mvutil.RateLimit
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	jsonStr, _ := json.Marshal(v)
	err := json.Unmarshal(jsonStr, &res)
	if err != nil {
		return nil
	}
	return &res
}

func GetRateLimit() *mvutil.RateLimit {
	value, err := getConfigCenterValue("RateLimit")
	var res = &mvutil.RateLimit{}
	if err != nil {
		return res
	}
	v, ok := value.(*mvutil.RateLimit)
	if !ok {
		return res
	}
	return v
}

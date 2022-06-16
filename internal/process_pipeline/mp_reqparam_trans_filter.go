package process_pipeline

import (
	"net/http"
	"strconv"
	"strings"

	"gitlab.mobvista.com/ADN/adnet/internal/errorcode"
	"gitlab.mobvista.com/ADN/adnet/internal/extractor"
	"gitlab.mobvista.com/ADN/adnet/internal/mvconst"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
)

type MPReqparamTransFilter struct {
}

var mpQueryMap = map[string]string{
	"system":       "platform",
	"os_v":         "os_version",
	"app_pname":    "package_name",
	"app_vn":       "app_version_name",
	"app_vc":       "app_version_code",
	"direction":    "orientation",
	"brand":        "brand",
	"model":        "model",
	"adid":         "gaid",
	"mnc":          "mnc",
	"mcc":          "mcc",
	"network":      "network_type",
	"language":     "language",
	"timezone":     "timezone",
	"ua":           "useragent",
	"sdkversion":   "sdk_version",
	"screen_size":  "screen_size",
	"ma":           "d1",
	"mb":           "d2",
	"mc":           "d3",
	"appid":        "app_id",
	"placement_id": "unit_id",
	"code":         "sign",
	"ad_cat":       "category",
	"r_num":        "ad_num",
	"n_num":        "tnum",
	"noticetype":   "ping_mode",
	"p_size":       "unit_size",
	"removeads":    "exclude_ids",
	"ins_ads":      "install_ids",
	"offset":       "offset",
	"only_im":      "only_impression",
	"ads_id":       "ad_source_id",
	"fomat":        "ad_type",
	"pre_pid":      "pre_pid",
	"ima":          "dvi",
	"idfa":         "idfa",
}

var mpNewAdQueryMap = map[string]string{
	"plattem":     "platform",
	"oson":        "os_version",
	"pkna":        "package_name",
	"appsionna":   "app_version_name",
	"apver":       "app_version_code",
	"orition":     "orientation",
	"bran":        "brand",
	"mdel":        "model",
	"adrdeid":     "android_id",
	"phemi":       "imei",
	"pheca":       "mac",
	"adverid":     "gaid",
	"iidiosa":     "idfa",
	"bnc":         "mnc",
	"bmc":         "mcc",
	"nwork":       "network_type",
	"plguage":     "language",
	"timone":      "timezone",
	"userent":     "useragent",
	"sdkersi":     "sdk_version",
	"gleverm":     "gp_version",
	"ssize":       "screen_size",
	"b1":          "d1",
	"b2":          "d2",
	"b3":          "d3",
	"imaextra":    "dvi",
	"osvnm":       "os_vv",
	"adsnorl":     "ad_s",
	"apmid":       "app_id",
	"unplid":      "unit_id",
	"sgncod":      "sign",
	"ctgor":       "category",
	"t3num":       "ad_num",
	"n3num":       "tnum",
	"b3num":       "b_num",
	"nttymode":    "ping_mode",
	"sltsize":     "unit_size",
	"reoffers":    "exclude_ids",
	"instloffers": "install_ids",
	"offerset3":   "offset",
	"oly3_imp":    "only_impression",
	"sourid":      "ad_source_id",
	"adtyp":       "ad_type",
	"adscenar":    "scenario",
}

func (mprtf *MPReqparamTransFilter) Process(data interface{}) (interface{}, error) {
	in, ok := data.(*http.Request)
	if !ok {
		return nil, errorcode.EXCEPTION_SPECIAL_EMPTY_STRING
		//return mvconst.GetResEmpty(), errors.New(mvconst.EXCEPTION_SPECIAL_RESULT_RETURN_NOENCODE)
	}

	rawQuery := mvutil.RequestQueryMap(in.URL.Query())
	if len(rawQuery) <= 0 {
		return nil, errorcode.EXCEPTION_SPECIAL_EMPTY_STRING
		//return mvconst.GetResEmpty(), errors.New(mvconst.EXCEPTION_SPECIAL_RESULT_RETURN_NOENCODE)
	}

	query := make(mvutil.RequestQueryMap, len(rawQuery))
	if in.URL.Path == mvconst.PATHMPNewAD {
		query = translateMpMap(rawQuery, mpNewAdQueryMap)
	} else {
		query = translateMpMap(rawQuery, mpQueryMap)
	}

	r := mvutil.RequestParams{}
	r.Param.MpReqDomain = in.Host
	RenderReqParam(in, &r, query)
	RenderMPParam(&r)
	transLieBaoUnit(&r)
	return &r, nil
}

func translateMpMap(rawQuery mvutil.RequestQueryMap, mpMap map[string]string) mvutil.RequestQueryMap {
	query := make(mvutil.RequestQueryMap, len(rawQuery))
	for k, v := range rawQuery {
		mKey, ok := mpMap[k]
		if ok {
			query[mKey] = v
		} else {
			query[k] = v
		}
	}
	return query
}

func transLieBaoUnit(r *mvutil.RequestParams) {
	transMap, _ := extractor.GetMP_MAP_UNIT()
	if len(transMap) <= 0 {
		return
	}
	unitId, _ := r.QueryMap.GetInt("unit_id", 0)
	if unitId <= 0 {
		return
	}
	unitIdStr := strconv.Itoa(unitId)
	key1 := "mv_" + unitIdStr
	transObj, ok := transMap[key1]
	if ok {
		r.QueryMap["app_id"] = mvutil.RenderQueryMapV(strconv.FormatInt(transObj.AppID, 10))
		r.QueryMap["unit_id"] = mvutil.RenderQueryMapV(strconv.FormatInt(transObj.UnitID, 10))
		r.QueryMap["ad_type"] = mvutil.RenderQueryMapV("0")
	} else {
		key2 := "power_" + unitIdStr
		transObj2, ok2 := transMap[key2]
		if ok2 {
			r.QueryMap["app_id"] = mvutil.RenderQueryMapV(strconv.FormatInt(transObj2.AppID, 10))
		}
	}
}

func RenderMPParam(r *mvutil.RequestParams) {
	r.QueryMap["version_flag"] = mvutil.RenderQueryMapV("1")
	osVersion := renderMPOSVersion(*r)
	r.QueryMap["os_version"] = mvutil.RenderQueryMapV(osVersion)
}

func renderMPOSVersion(r mvutil.RequestParams) string {
	// 若有传os_vv，则使用os_vv的值作为os_version
	osVersionNew, _ := r.QueryMap.GetString("os_vv", true, "")
	if len(osVersionNew) > 0 {
		return osVersionNew
	}
	// 优先从ua中获取osversion
	userAgent, _ := r.QueryMap.GetString("useragent", true, "")
	osVersion, _ := r.QueryMap.GetString("os_version", true, "")
	res := ""
	if len(userAgent) > 0 {
		os := mvutil.UaParser.ParseOs(userAgent)
		res = strings.ToLower(os.ToVersionString())
		if len(res) > 0 {
			return res
		}
	}
	osVersionInt, err := strconv.Atoi(osVersion)
	if err == nil {
		res = mvutil.GetAndroidOSVersion(osVersionInt)
	}
	return res
}

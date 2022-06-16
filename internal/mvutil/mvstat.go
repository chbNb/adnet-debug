package mvutil

import (
	"bytes"
	"reflect"
	"strconv"
)

// impression click track等中的z，去除掉设备信息等的p。p中部分设备信息SDK带不过来，也放在此字段
var FIELDS_Z = []string{
	"ImageSize",
	"Extra3",
	"Extra4",
	"TNum",
	"ExtBrand",
	"ExtModel",
	"-",
	"ExtApiVersion",
	"MCCMNC",
	"OAID",
	"ExtBigTemId",
	"ExtPlacementId",
	"RespFillEcpmFloor",
	"StartModeTagsStr",
	"IDFVOpenIDFA",
	"JunoCommonLogInfoJson",
	"-",
	"ExtDeviceId",
}

// onlyImpression 的z，相当于以前的p
var FIELDS_OZ = []string{
	"PublisherID",
	"AppID",
	"UnitID",
	"AdvertiserID",
	"CampaignID",
	"-",
	"Scenario",
	"-",
	"ImageSize",
	"RequestType",
	"-",
	"-",
	"-",
	"-",
	"-",
	"-",
	"CountryCode",
	"-",
	"-",
	"MCCMNC",
	"Extra",
	"-",
	"Extra3",
	"Extra4",
	"Extra5",
	"-",
	"Extra7",
	"Extra8",
	"Extra9",
	"Extra10",
	"RequestID",
	"-",
	"-",
	"-",
	"-",
	"ServerIP",
	"-",
	"-",
	"-",
	"-",
	"-",
	"AppVersionName",
	"-",
	"RemoteIP",
	"-",
	"-",
	"-",
	"CityCode",
	"Extra13",
	"Extra14",
	"Extra15",
	"Extra16",
	"IDFVOpenIDFA",
	"-",
	"-",
	"Extra20",
	"-",
	"Extfinalsubid",
	"-",
	"-",
	"-",
	"ExtpackageName",
	"-",
	"ExtflowTagId",
	"-",
	"Extendcard",
	"ExtrushNoPre",
	"-",
	"ExtfinalPackageName",
	"Extnativex",
	"-",
	"-",
	"-",
	"Extctype",
	"Extrvtemplate",
	"-",
	"-",
	"-",
	"-",
	"-",
	"-",
	"Extb2t",
	"Extchannel",
	"-",
	"-",
	"-",
	"Extbp",
	"Extsource",
	"-",
	"Extalgo",
	"ExtthirdCid",
	"ExtifLowerImp",
	"-",
	"-",
	"ExtsystemUseragent",
	"-",
	"-",
	"-",
	"-",
	"-",
	"ExtMpNormalMap",
	"_",
	"_",
	"ExtBrand",
	"ExtModel",
	"ExtData",
	"_",
	"_",
	"_",
	"_",
	"_",
	"_",
	"_",
	"_",
	"_",
	"_",
	"_",
	"_",
	"_",
	"_",
	"_",
	"_",
	"_",
	"_",
	"_",
	"_",
	"_",
	"_",
	"_",
	"_",
	"_",
	"OAID",
	"_",
	"_",
	"_",
	"ExtBigTplOfferData",
	"ExtBigTemId",
	"_",
	"_",
	"ExtPlacementId",
	"_",
	"_",
	"_",
	"_",
	"_",
	"_",
	"_",
	"_",
	"_",
	"_",
	"_",
	"_",
	"_",
	"_",
	"ExtDeviceId",
}

func SerializeZ(params *Params) string {
	res := Serialize(params, FIELDS_Z, "|", "")
	res = Base64([]byte(res))
	return res
}

func SerializeOZ(params *Params) string {
	res := Serialize(params, FIELDS_OZ, "|", "")
	res = Base64([]byte(res))
	return res
}

func Serialize(params *Params, fields []string, split string, def string) string {
	var buf bytes.Buffer
	v := reflect.ValueOf(*params)
	for i, field := range fields {
		if field == "-" {
			if len(def) > 0 {
				buf.WriteString(def)
			}
			if i < len(fields)-1 {
				buf.WriteString(split)
			}
			continue
		}
		vfield := v.FieldByName(field)
		switch vfield.Kind() {
		case reflect.Bool:
			if vfield.Bool() {
				buf.WriteString("1")
			} else {
				buf.WriteString("0")
			}

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			buf.WriteString(strconv.FormatInt(vfield.Int(), 10))

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			buf.WriteString(strconv.FormatUint(vfield.Uint(), 10))

		case reflect.String:
			buf.WriteString(vfield.String())

		case reflect.Float32, reflect.Float64:
			buf.WriteString(strconv.FormatFloat(vfield.Float(), 'f', 6, 64))
		}

		if i < len(fields)-1 {
			buf.WriteString(split)
		}
	}
	return buf.String()
}

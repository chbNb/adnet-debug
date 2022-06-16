package helpers

import (
	"net/http"
	"strings"
	"testing"

	jsoniter "github.com/json-iterator/go"
	"github.com/smartystreets/goconvey/convey"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/params"
)

func TestBase64Decode(t *testing.T) {
	convey.Convey("TestTD", t, func() {
		// s := "foo|bar|foo|bar|foo|bar|foo|bar|foo|bar|foo|bar|foo|bar|foo"
		// d1 := Base64Encode(s)
		// fmt.Println(d1)
		// d2 := Base64Decode(d1)
		// fmt.Println(d2)
		// d3 := url.QueryEscape(d1) // note: it's must where the client request param
		// fmt.Println(d3)

		raws := "A0635584-FCB1-4106-B924-A80C29150E4D|||XDFGG1245gfsys_id|XDFGG1245gfbkup_id|apple|iPhone10,3|9|mi_5.4.0|2436x1125|Mozilla/5.0 (iPhone; CPU iPhone OS 11_1_2 like Mac OS X) AppleWebKit/604.3.5 (KHTML, like Gecko) Mobile/15B202||||1.9"
		td := Base64Encode(raws)
		t.Logf("%s", td)
		td = strings.Replace(td, " ", "+", -1)
		rawTd := Base64Decode(td)
		convey.So(strings.Count(rawTd, "|"), convey.ShouldEqual, 14)

		req, err := http.NewRequest("GET", "http://hb.com/load?app_id=124751&unit_id=204142&sign=db0efd5d688b0a59dab12268901997c7&req_type=2&ad_num=1&tnum=1&only_impression=1&ping_mode=1&display_info=%5B%5D&ad_source_id=1&ad_type=287&offset=0&token=023fe3e8-2e3c-4a45-9638-8117e4175301_hk&channel=&platform=1&os_version=8.1.0&package_name=com.GMA.Ball.Sort.Puzzle&app_version_name=1.2&app_version_code=36&orientation=1&model=DUB-LX1&brand=HUAWEI&gaid=&gaid2=inQQfUjbiFNwHkftGoTTinh9WnttirjwfkieDnVMDAfrGkfQ&mnc=&mcc=&network_type=9&network_str=&language=en&timezone=&useragent=Mozilla%2F5.0+%28Linux%3B+Android+8.1.0%3B+DUB-LX1+Build%2FHUAWEIDUB-LX1%3B+wv%29+AppleWebKit%2F537.36+%28KHTML%2C+like+Gecko%29+Version%2F4.0+Chrome%2F80.0.3987.119+Mobile+Safari%2F537.36&sdk_version=MAL_12.1.51&gp_version=18.9.11-all+%5B0%5D+%5BPR%5D+295870256&screen_size=720x1520&version_flag=1&cache1=57581&cache2=34424&power_rate=4&charging=0&has_wx=false&pkg_source=com.android.vending&http_req=2&dvi=4BzXDkQ3RUE0fFNbiFxQH7DwHrDbf02QDri2WkHrGkRwf7zrHr5QDrD2fAc0R0M0DFQ3RUE0ioRsRrxwJoR1RURPfnD0Woz3YkD0GUjMfA3sRrfTRUE0kFcBYnDTW%2BD9DZMlD%2BzwHkc0LZ2FfFjsR7cBYk5tDrQJRgT%3D&unknown_source=1&sys_id=736f0b62-ca92-5f22-b686-6a18ed380e46&api_version=1.8", nil)
		convey.So(err, convey.ShouldBeNil)

		httpReqData := &params.HttpReqData{}
		httpReqData.Path = req.URL.Path
		httpReqData.Host = req.Host
		httpReqData.QueryData = params.HttpQueryMap(req.URL.Query())

		bDvi := httpReqData.QueryData.GetString("dvi", true, "")
		dvi := Base64Decode(bDvi)
		t.Logf("%s", dvi)
		deviceInfo := DeviceInfo{}
		err = jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal([]byte(dvi), &deviceInfo)
		convey.So(err, convey.ShouldBeNil)
		t.Logf("%+v", deviceInfo)
		convey.So(deviceInfo.Dmt, convey.ShouldEqual, "2956")
		convey.So(deviceInfo.Dmf, convey.ShouldEqual, 1079)
		convey.So(deviceInfo.Ct, convey.ShouldEqual, "[arm64-v8a, armeabi-v7a, armeabi]")

		req, err = http.NewRequest("GET", "http://hb.com/bid?app_id=124716&unit_id=203755&bid_floor=0&exclude_ids=&install_ids=%5B340942772%5D&display_info=%5B%5D&channel=&platform=1&os_version=8.0.0&package_name=com.leshu.wordPK&app_version_name=1.6&app_version_code=7&orientation=1&model=XT1924-9&brand=motorola&gaid=&gaid2=iAQ0GaV2i7NwDAjMHoTTH75rWnttHUlwHnveiAz0iFjBfaxU&d1=iAV2fnVTia3eiaDAfAjF&d3=f7VFiAV9G7DBGkHriUNeGv%3D%3D&mnc=01&mcc=460&network_type=9&network_str=&language=zh&timezone=&useragent=Mozilla%2F5.0%20%28Linux%3B%20Android%208.0.0%3B%20XT1924-9%20Build%2FOCP27.91-38-31-33%3B%20wv%29%20AppleWebKit%2F537.36%20%28KHTML%2C%20like%20Gecko%29%20Version%2F4.0%20Chrome%2F64.0.3282.137%20Mobile%20Safari%2F537.36&sdk_version=MAL_10.1.31&gp_version=14.3.30-all%20%5B2%5D%20%5BPR%5D%20242795620&screen_size=720x1440&cache1=56495&cache2=38750&power_rate=59&charging=1&has_wx=true&pkg_source=com.android.packageinstaller&http_req=2&dvi=4BzuYk5uRUE0iAV2fnVTia3eiaDAfAjFR0M0YkcURUE0HaxUGnx0fUhAfAVPR0M0Lk2ALZR1RUNFiavefaRPGn3eiA3TGoRsRrKtLkN0G0R0WoztYrxBYFQ3%2BFQ3RUE0f7VFiAV9G7DBGkHriUNeGoRsRrfuHoR1RUv06N%3D%3D&unknown_source=1&sys_id=5a611aba-43c0-5d35-b1da-0ccfab302aa5&api_version=1.8", nil)
		convey.So(err, convey.ShouldBeNil)
		httpReqData.QueryData = params.HttpQueryMap(req.URL.Query())
		bDvi = httpReqData.QueryData.GetString("dvi", true, "")
		dvi = Base64Decode(bDvi)
		t.Logf("%s", dvi)
		deviceInfo = DeviceInfo{}
		err = jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal([]byte(dvi), &deviceInfo)
		convey.So(err, convey.ShouldBeNil)
		t.Logf("%+v", deviceInfo)
		convey.So(deviceInfo.Dmt, convey.ShouldEqual, "")
		convey.So(deviceInfo.Dmf, convey.ShouldEqual, 0)
		convey.So(deviceInfo.Ct, convey.ShouldEqual, "")

		req, err = http.NewRequest("GET", "http://13.251.178.19/bid?debug=0&app_id=92762&unit_id=21310&sign=4c3814fdba4acd6f3b092b3df43f8d9a&req_type=2&ad_num=1&tnum=1&only_impression=1&ping_mode=1&exclude_ids=%5B20436864413%2C20975844226%2C20347290220%5D&ad_source_id=1&session_id=5be52853c6c1e21284109e47&ad_type=94&offset=0&channel=&platform=1&os_version=8.0.0&package_name=com.mintegral.sdk.demo&app_version_name=1.0&app_version_code=1&orientation=1&model=SM-G9550&brand=samsung&gaid=&gaid2=faiFGkiMi7jwGaiAioTTDkfUWnt0DkiwiU5tiUxUHU5tikVT&mnc=07&mcc=460&network_type=9&network_str=&language=zh&timezone=&useragent=Mozilla%2F5.0%20%28Linux%3B%20Android%208.0.0%3B%20SM-G9550%20Build%2FR16NW%3B%20wv%29%20AppleWebKit%2F537.36%20%28KHTML%2C%20like%20Gecko%29%20Version%2F4.0%20Chrome%2F70.0.3538.80%20Mobile%20Safari%2F537.36&sdk_version=MAL_99.2.0&gp_version=11.4.16-all%20%5B0%5D%20%5BPR%5D%20209796717&screen_size=1080x2220&version_flag=1&cache1=57928&cache2=49827&power_rate=100&charging=1&has_wx=true&http_req=2&unknown_source=1&api_version=1.4&imei=352097093386308&dvi=4BzuYk5uRUE0iAVBia3bia3AiAlFiAv9R0M0YkcURUE0faNPinDMfAx0Dn5UR0M0DkP3hrKuHcKuHoR1RUxUfrH0inH3HnRBDkieHaV0WozULkN0G0ReiUlbioRsRoz3Y%2bN0G0ReinjeR0MlRrxwH0R1iURBi0MlRrfTRUE0kbt9WdQP%2bZzK", nil)
		convey.So(err, convey.ShouldBeNil)
		httpReqData.QueryData = params.HttpQueryMap(req.URL.Query())
		bDvi = httpReqData.QueryData.GetString("dvi", true, "")
		dvi = Base64Decode(bDvi)
		t.Logf("%s", dvi)
		deviceInfo = DeviceInfo{}
		err = jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal([]byte(dvi), &deviceInfo)
		convey.So(err, convey.ShouldBeNil)
		t.Logf("%+v", deviceInfo)
		convey.So(deviceInfo.Dmt, convey.ShouldEqual, "1111")
		convey.So(deviceInfo.Dmf, convey.ShouldEqual, 2222)
		convey.So(deviceInfo.Ct, convey.ShouldEqual, "[xx,yy]")

		req, err = http.NewRequest("GET", "https://hb.rayjump.com/bid?api_version=1.9&app_id=118692&app_version_name=6.1.0&ats=4BzGVTcsY7KbhTcBDrQThrcB4VeXDkxARUuThg5Q6N%3D%3D&bid_floor=0&cache1=63937.039062&cache2=20672.832031&ch_info=hgDlHF5TR7zuHv%3D%3D&charging=1&ct=%5B16777228%5D&dmf=81.75&dmt=2820.046875&http_req=2&idfa=3AB54860-09DF-4B0E-A3F4-CBC2106EE43E&idfv=4AE40087-D094-4A2A-B523-5D97E7B066CC&language=zh-Hans-CN&mcc=460&mnc=02&model=iPhone10%2C3&network_str=&network_type=9&openidfa=&orientation=1&os_version=13.3.1&package_name=com.mobvista.ui.test&platform=2&power_rate=100&screen_size=1125.000000x2436.000000&sdk_version=MI_6.1.0&sub_ip=192.168.4.22&sys_id=fc0cf406-2ec1-59ab-a580-3ad25e322122&ui_orientation=30&unit_id=146892&useragent=Mozilla/5.0%20%28iPhone%3B%20CPU%20iPhone%20OS%2013_3_1%20like%20Mac%20OS%20X%29%20AppleWebKit/605.1.15%20%28KHTML%2C%20like%20Gecko%29%20Mobile/15E148", nil)
		convey.So(err, convey.ShouldBeNil)
		httpReqData.QueryData = params.HttpQueryMap(req.URL.Query())
		chInfo := httpReqData.QueryData.GetString("ch_info", false, "")
		convey.So(chInfo, convey.ShouldEqual, "hgDlHF5TR7zuHv==")
		chInfo = Base64Decode(chInfo)
		convey.So(chInfo, convey.ShouldEqual, "rv get bid")
	})

	convey.Convey("test td decode for event log", t, func() {
		var mediationName string
		rawTd := "14228|111219|203066|hb|openapi|banner|ios|mi_5.8.7|146|JP|35690|1|0|2|8f4c7e52-b2d9-4faa-bad7-38c47c4213bb_vg"
		tokenArr := strings.SplitN(rawTd, "|", 16)
		convey.So(len(tokenArr), convey.ShouldEqual, 15)
		if len(tokenArr) > 15 {
			mediationName = tokenArr[15]
		}
		convey.So(mediationName, convey.ShouldEqual, "")
		rawTd = "14228|111219|203066|hb|openapi|banner|ios|mi_5.8.7|146|JP|35690|1|0|2|8f4c7e52-b2d9-4faa-bad7-38c47c4213bb_vg|"
		tokenArr = strings.SplitN(rawTd, "|", 16)
		convey.So(len(tokenArr), convey.ShouldEqual, 16)
		if len(tokenArr) > 15 {
			mediationName = tokenArr[15]
		}
		convey.So(mediationName, convey.ShouldEqual, "")
		rawTd = "14228|111219|203066|hb|openapi|banner|ios|mi_5.8.7|146|JP|35690|1|0|2|8f4c7e52-b2d9-4faa-bad7-38c47c4213bb_vg|applovin_mediation"
		tokenArr = strings.SplitN(rawTd, "|", 16)
		convey.So(len(tokenArr), convey.ShouldEqual, 16)
		if len(tokenArr) > 15 {
			mediationName = tokenArr[15]
		}
		convey.So(mediationName, convey.ShouldEqual, "applovin_mediation")
	})

	convey.Convey("test buyeruid has \\ char", t, func() {
		uid := `fUQcNVfofaiwiAJafBTTGVDFWVjFGnRwiA3eiAz7N3VTGViA6ahMfnlFxaSaWnR9iVVwfaNTNZ2Sf3ReWnNbfaiAxnjTfUj9xdeIHnc0iFHQi7Dwfnz0DBT2fAvPWn3efalwDU3BHU3FHrVBDk536deIL5SEYFPQinjsGdMP6j2z+AD\/iB9b6alBGo9MiavMiaS9inhPi09MiavMiaSInkK1LkesDZI2WUvlp7QNL7K\/HnslN2S5R7QNL7K\/HZSyVBvei2IFR7euLFVlnkcURjKnRcluRjcMh7eQ5F50ZFQTWADMfZ9eWUj2RotWZcxfnoMlY7Q8HZSdHkf8YB3lnkK0LkeQWAj2xnjTGdeIk22IinVPfAvTfaRPGaV2GdMBWUSIiUlAiZ9AiUleiU5IinVPWUD9fA5IkAjFfAhbiUR9++MPGdMM6ajBfA39GZ9TfUlbfnSIinvAGnV\/iAleGaiF`
		uid = strings.Replace(uid, " ", "+", -1)
		t.Log(uid)
		d1 := Base64Decode(uid)
		t.Log(d1)
		uid = strings.Replace(uid, "\\", "", -1)
		t.Log(uid)
		d2 := Base64Decode(uid)
		t.Log(d2)
	})
}

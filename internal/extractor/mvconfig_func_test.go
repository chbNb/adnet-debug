package extractor

import (
	. "github.com/bouk/monkey"
	. "github.com/smartystreets/goconvey/convey"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	smodel "gitlab.mobvista.com/ADN/structs/model"
	"testing"
)

func TestGetOpenapiScenario(t *testing.T) {
	Convey("test GetOpenapiScenario", t, func() {
		guard := Patch(GetMVConfigValue, func(key string) (interface{}, bool) {
			return []string{"s1", "s2"}, true
		})
		defer guard.Unpatch()
		res, bo := GetOpenapiScenario()
		So(res, ShouldResemble, []string{"s1", "s2"})
		So(bo, ShouldBeTrue)
	})
}

func TestPrivategetOpenapiScenario(t *testing.T) {
	Convey("test getOpenapiScenario", t, func() {
		//guard := Patch(GetMVConfigValue, func(key string) (value []byte, ifFind bool) {
		//	return []byte(key), true
		//})
		//defer guard.Unpatch()

		res, b := getOpenapiScenario([]byte(""))
		var exp []string
		So(res, ShouldResemble, exp)
		So(b, ShouldBeFalse)
	})
}

func TestGetConfigLowGradeApp(t *testing.T) {
	Convey("test GetConfigLowGradeApp", t, func() {
		guard := Patch(GetMVConfigValue, func(key string) (interface{}, bool) {
			return int64(42), true
		})
		defer guard.Unpatch()
		res, bo := GetConfigLowGradeApp()
		So(res, ShouldResemble, int64(42))
		So(bo, ShouldBeTrue)
	})
}

func TestPrivategetConfigLowGradeApp(t *testing.T) {
	Convey("test getConfigLowGradeApp", t, func() {
		//guard := Patch(GetMVConfigValue, func(key string) (value []byte, ifFind bool) {
		//	return []byte(key), true
		//})
		//defer guard.Unpatch()

		res, b := getConfigLowGradeApp([]byte(""))
		So(res, ShouldEqual, int64(0))
		So(b, ShouldBeFalse)
	})
}

func TestGetReplenishApp(t *testing.T) {
	Convey("test GetReplenishApp", t, func() {
		guard := Patch(GetMVConfigValue, func(key string) (interface{}, bool) {
			return []int64{1, 3, 5}, true
		})
		defer guard.Unpatch()
		res, bo := GetReplenishApp()
		So(res, ShouldResemble, []int64{1, 3, 5})
		So(bo, ShouldBeTrue)
	})
}

func TestPrivategetReplenishApp(t *testing.T) {
	Convey("test getReplenishApp", t, func() {
		//guard := Patch(GetMVConfigValue, func(key string) (value []byte, ifFind bool) {
		//	return []byte(key), true
		//})
		//defer guard.Unpatch()

		res, b := getReplenishApp([]byte(""))
		So(res, ShouldResemble, []int64{})
		So(b, ShouldBeFalse)
	})
}

func TestGetAdvBlackSubIdList(t *testing.T) {
	Convey("test GetAdvBlackSubIdList", t, func() {
		guard := Patch(GetMVConfigValue, func(key string) (interface{}, bool) {
			return map[string]map[string]map[string]string{}, true
		})
		defer guard.Unpatch()
		res, bo := GetAdvBlackSubIdList()
		So(res, ShouldResemble, map[string]map[string]map[string]string{})
		So(bo, ShouldBeTrue)
	})
}

func TestPrivategetAdvBlackSubIdList(t *testing.T) {
	Convey("test getAdvBlackSubIdList", t, func() {
		//guard := Patch(GetMVConfigValue, func(key string) (value []byte, ifFind bool) {
		//	return []byte(key), true
		//})
		//defer guard.Unpatch()

		res, b := getAdvBlackSubIdList([]byte(""))
		So(res, ShouldResemble, map[string]map[string]map[string]string{})
		So(b, ShouldBeFalse)
	})
}

func TestGetAppPostList(t *testing.T) {
	Convey("test GetAppPostList", t, func() {
		guard := Patch(GetMVConfigValue, func(key string) (interface{}, bool) {
			return map[int32]*smodel.AppPostList{}, true
		})
		defer guard.Unpatch()
		res, bo := GetAppPostList()
		So(res, ShouldResemble, map[int32]*smodel.AppPostList{})
		So(bo, ShouldBeTrue)
	})
}

func TestPrivategetAppPostList(t *testing.T) {
	Convey("test getAppPostList", t, func() {
		//guard := Patch(GetMVConfigValue, func(key string) (value []byte, ifFind bool) {
		//	return []byte(key), true
		//})
		//defer guard.Unpatch()

		res, b := getAppPostList([]byte(""))
		So(res, ShouldResemble, map[int32]*smodel.AppPostList{})
		So(b, ShouldBeFalse)
	})
}

func TestGetOfferwallGuidelines(t *testing.T) {
	Convey("test GetOfferwallGuidelines", t, func() {
		guard := Patch(GetMVConfigValue, func(key string) (interface{}, bool) {
			return map[int32]string{}, true
		})
		defer guard.Unpatch()
		res, bo := GetOfferwallGuidelines()
		So(res, ShouldResemble, map[int32]string{})
		So(bo, ShouldBeTrue)
	})
}

func TestPrivategetOfferwallGuidelines(t *testing.T) {
	Convey("test getOfferwallGuidelines", t, func() {
		//guard := Patch(GetMVConfigValue, func(key string) (value []byte, ifFind bool) {
		//	return []byte(key), true
		//})
		//defer guard.Unpatch()

		res, b := getOfferwallGuidelines([]byte(""))
		So(res, ShouldResemble, map[int32]string{})
		So(b, ShouldBeFalse)
	})
}

func TestPrivategetEndcard(t *testing.T) {
	Convey("test getEndcard", t, func() {
		//guard := Patch(GetMVConfigValue, func(key string) (value []byte, ifFind bool) {
		//	return []byte(key), true
		//})
		//defer guard.Unpatch()

		res, b := getEndcard([]byte(""))
		var exp smodel.EndCard
		So(res, ShouldResemble, exp)
		So(b, ShouldBeFalse)
	})
}

func TestGetAdServerTestConfig(t *testing.T) {
	Convey("test GetAdServerTestConfig", t, func() {
		guard := Patch(GetMVConfigValue, func(key string) (interface{}, bool) {
			return mvutil.AdServerTestConfig{}, true
		})
		defer guard.Unpatch()
		res := GetAdServerTestConfig()
		So(res, ShouldResemble, mvutil.AdServerTestConfig{})
	})
}

func TestPrivategetAdServerTestConfig(t *testing.T) {
	Convey("test getAdServerTestConfig", t, func() {
		//guard := Patch(GetMVConfigValue, func(key string) (value []byte, ifFind bool) {
		//	return []byte(key), true
		//})
		//defer guard.Unpatch()

		res, b := getAdServerTestConfig([]byte(""))
		var exp mvutil.AdServerTestConfig
		So(res, ShouldResemble, exp)
		So(b, ShouldBeFalse)
	})
}

func TestGetDefRVTemplate(t *testing.T) {
	Convey("test GetDefRVTemplate", t, func() {
		guard := Patch(GetMVConfigValue, func(key string) (interface{}, bool) {
			return &smodel.VideoTemplateUrlItem{}, true
		})
		defer guard.Unpatch()
		res, bo := GetDefRVTemplate()
		So(res, ShouldResemble, &smodel.VideoTemplateUrlItem{})
		So(bo, ShouldBeTrue)
	})
}

func TestPrivategetDefRVTemplate(t *testing.T) {
	Convey("test getDefRVTemplate", t, func() {
		//guard := Patch(GetMVConfigValue, func(key string) (value []byte, ifFind bool) {
		//	return []byte(key), true
		//})
		//defer guard.Unpatch()

		res, b := getDefRVTemplate([]byte(""))
		exp := &smodel.VideoTemplateUrlItem{}
		So(res, ShouldResemble, exp)
		So(b, ShouldBeFalse)
	})
}

func TestGetVersionCompare(t *testing.T) {
	Convey("test GetVersionCompare", t, func() {
		guard := Patch(GetMVConfigValue, func(key string) (interface{}, bool) {
			return map[string]mvutil.VersionCompare{}, true
		})
		defer guard.Unpatch()
		res, bo := GetVersionCompare()
		So(res, ShouldResemble, map[string]mvutil.VersionCompare{})
		So(bo, ShouldBeTrue)
	})
}

func TestPrivategetVersionCompare(t *testing.T) {
	Convey("test getVersionCompare", t, func() {
		//guard := Patch(GetMVConfigValue, func(key string) (value []byte, ifFind bool) {
		//	return []byte(key), true
		//})
		//defer guard.Unpatch()

		res, b := getVersionCompare([]byte(""))
		exp := map[string]mvutil.VersionCompare{}
		So(res, ShouldResemble, exp)
		So(b, ShouldBeFalse)
	})
}

func TestGetPlayableTest(t *testing.T) {
	Convey("test GetPlayableTest", t, func() {
		guard := Patch(GetMVConfigValue, func(key string) (interface{}, bool) {
			return map[int64]mvutil.PlatableTest{}, true
		})
		defer guard.Unpatch()
		res, bo := GetPlayableTest()
		So(res, ShouldResemble, map[int64]mvutil.PlatableTest{})
		So(bo, ShouldBeTrue)
	})
}

func TestPrivategetPlayableTest(t *testing.T) {
	Convey("test getPlayableTest", t, func() {
		//guard := Patch(GetMVConfigValue, func(key string) (value []byte, ifFind bool) {
		//	return []byte(key), true
		//})
		//defer guard.Unpatch()

		res, b := getPlayableTest([]byte(""))
		exp := map[int64]mvutil.PlatableTest{}
		So(res, ShouldResemble, exp)
		So(b, ShouldBeFalse)
	})
}

func TestGetTRACK_URL_CONFIG_NEW(t *testing.T) {
	Convey("test GetTRACK_URL_CONFIG_NEW", t, func() {
		guard := Patch(GetMVConfigValue, func(key string) (interface{}, bool) {
			return map[int32]mvutil.TRACK_URL_CONFIG_NEW{}, true
		})
		defer guard.Unpatch()
		res, bo := GetTRACK_URL_CONFIG_NEW()
		So(res, ShouldResemble, map[int32]mvutil.TRACK_URL_CONFIG_NEW{})
		So(bo, ShouldBeTrue)
	})
}

func TestPrivategetTRACK_URL_CONFIG_NEW(t *testing.T) {
	Convey("test getTRACK_URL_CONFIG_NEW", t, func() {
		//guard := Patch(GetMVConfigValue, func(key string) (value []byte, ifFind bool) {
		//	return []byte(key), true
		//})
		//defer guard.Unpatch()

		res, b := getTRACK_URL_CONFIG_NEW([]byte(""))
		exp := map[int32]mvutil.TRACK_URL_CONFIG_NEW{}
		So(res, ShouldResemble, exp)
		So(b, ShouldBeFalse)
	})
}

func TestGet3S_CHINA_DOMAIN(t *testing.T) {
	Convey("test Get3S_CHINA_DOMAIN", t, func() {
		guard := Patch(GetMVConfigValue, func(key string) (interface{}, bool) {
			return mvutil.CONFIG_3S_CHINA_DOMAIN{}, true
		})
		defer guard.Unpatch()
		res, bo := Get3S_CHINA_DOMAIN()
		So(res, ShouldResemble, mvutil.CONFIG_3S_CHINA_DOMAIN{})
		So(bo, ShouldBeTrue)
	})
}

func TestPrivateget3S_CHINA_DOMAIN(t *testing.T) {
	Convey("test get3S_CHINA_DOMAIN", t, func() {
		//guard := Patch(GetMVConfigValue, func(key string) (value []byte, ifFind bool) {
		//	return []byte(key), true
		//})
		//defer guard.Unpatch()

		res, b := get3S_CHINA_DOMAIN([]byte(""))
		var exp mvutil.CONFIG_3S_CHINA_DOMAIN
		So(res, ShouldResemble, exp)
		So(b, ShouldBeFalse)
	})
}

func TestGetJUMP_TYPE_CONFIG(t *testing.T) {
	Convey("test GetJUMP_TYPE_CONFIG", t, func() {
		guard := Patch(GetMVConfigValue, func(key string) (interface{}, bool) {
			return map[string]int32{}, true
		})
		defer guard.Unpatch()
		res, bo := GetJUMP_TYPE_CONFIG()
		So(res, ShouldResemble, map[string]int32{})
		So(bo, ShouldBeTrue)
	})
}

func TestPrivategetJUMP_TYPE_CONFIG(t *testing.T) {
	Convey("test getJUMP_TYPE_CONFIG", t, func() {
		//guard := Patch(GetMVConfigValue, func(key string) (value []byte, ifFind bool) {
		//	return []byte(key), true
		//})
		//defer guard.Unpatch()

		res, b := getJUMP_TYPE_CONFIG([]byte(""))
		exp := map[string]int32{}
		So(res, ShouldResemble, exp)
		So(b, ShouldBeFalse)
	})
}

func TestGetJUMP_TYPE_CONFIG_IOS(t *testing.T) {
	Convey("test GetJUMP_TYPE_CONFIG_IOS", t, func() {
		guard := Patch(GetMVConfigValue, func(key string) (interface{}, bool) {
			return map[string]int32{}, true
		})
		defer guard.Unpatch()
		res, bo := GetJUMP_TYPE_CONFIG_IOS()
		So(res, ShouldResemble, map[string]int32{})
		So(bo, ShouldBeTrue)
	})
}

func TestPrivategetJUMP_TYPE_CONFIG_IOS(t *testing.T) {
	Convey("test getJUMP_TYPE_CONFIG_IOS", t, func() {
		//guard := Patch(GetMVConfigValue, func(key string) (value []byte, ifFind bool) {
		//	return []byte(key), true
		//})
		//defer guard.Unpatch()

		res, b := getJUMP_TYPE_CONFIG_IOS([]byte(""))
		exp := map[string]int32{}
		So(res, ShouldResemble, exp)
		So(b, ShouldBeFalse)
	})
}

func TestGetJUMPTYPE_SDKVERSION(t *testing.T) {
	Convey("test GetJUMPTYPE_SDKVERSION", t, func() {
		guard := Patch(GetMVConfigValue, func(key string) (interface{}, bool) {
			return map[string]map[string]string{}, true
		})
		defer guard.Unpatch()
		res, bo := GetJUMPTYPE_SDKVERSION()
		So(res, ShouldResemble, map[string]map[string]string{})
		So(bo, ShouldBeTrue)
	})
}

func TestPrivategetJUMPTYPE_SDKVERSION(t *testing.T) {
	Convey("test getJUMPTYPE_SDKVERSION", t, func() {
		//guard := Patch(GetMVConfigValue, func(key string) (value []byte, ifFind bool) {
		//	return []byte(key), true
		//})
		//defer guard.Unpatch()

		res, b := getJUMPTYPE_SDKVERSION([]byte(""))
		exp := map[string]map[string]string{}
		So(res, ShouldResemble, exp)
		So(b, ShouldBeFalse)
	})
}

func TestGetADSTACKING(t *testing.T) {
	Convey("test GetADSTACKING", t, func() {
		guard := Patch(GetMVConfigValue, func(key string) (interface{}, bool) {
			return mvutil.ADSTACKING{}, true
		})
		defer guard.Unpatch()
		res, bo := GetADSTACKING()
		So(res, ShouldResemble, mvutil.ADSTACKING{})
		So(bo, ShouldBeTrue)
	})
}

func TestPrivategetADSTACKING(t *testing.T) {
	Convey("test getADSTACKING", t, func() {
		//guard := Patch(GetMVConfigValue, func(key string) (value []byte, ifFind bool) {
		//	return []byte(key), true
		//})
		//defer guard.Unpatch()

		res, b := getADSTACKING([]byte(""))
		exp := mvutil.ADSTACKING{}
		So(res, ShouldResemble, exp)
		So(b, ShouldBeFalse)
	})
}

func TestGetSETTING_CONFIG(t *testing.T) {
	Convey("test GetSETTING_CONFIG", t, func() {
		guard := Patch(GetMVConfigValue, func(key string) (interface{}, bool) {
			return mvutil.SETTING_CONFIG{}, true
		})
		defer guard.Unpatch()
		res, bo := GetSETTING_CONFIG()
		So(res, ShouldResemble, mvutil.SETTING_CONFIG{})
		So(bo, ShouldBeTrue)
	})
}

func TestPrivategetSETTING_CONFIG(t *testing.T) {
	Convey("test getSETTING_CONFIG", t, func() {
		//guard := Patch(GetMVConfigValue, func(key string) (value []byte, ifFind bool) {
		//	return []byte(key), true
		//})
		//defer guard.Unpatch()

		res, b := getSETTING_CONFIG([]byte(""))
		exp := mvutil.SETTING_CONFIG{}
		So(res, ShouldResemble, exp)
		So(b, ShouldBeFalse)
	})
}

func TestGetTemplate(t *testing.T) {
	Convey("test GetTemplate", t, func() {
		guard := Patch(GetMVConfigValue, func(key string) (interface{}, bool) {
			return map[string]mvutil.Template{}, true
		})
		defer guard.Unpatch()
		res, bo := GetTemplate()
		So(res, ShouldResemble, map[string]mvutil.Template{})
		So(bo, ShouldBeTrue)
	})
}

func TestPrivategetTemplate(t *testing.T) {
	Convey("test getTemplate", t, func() {
		//guard := Patch(GetMVConfigValue, func(key string) (value []byte, ifFind bool) {
		//	return []byte(key), true
		//})
		//defer guard.Unpatch()

		res, b := getTemplate([]byte(""))
		exp := map[string]mvutil.Template{}
		So(res, ShouldResemble, exp)
		So(b, ShouldBeFalse)
	})
}

func TestGetofferwall_urls(t *testing.T) {
	Convey("test Getofferwall_urls", t, func() {
		guard := Patch(GetMVConfigValue, func(key string) (interface{}, bool) {
			return mvutil.COfferwallUrls{}, true
		})
		defer guard.Unpatch()
		res, bo := Getofferwall_urls()
		So(res, ShouldResemble, mvutil.COfferwallUrls{})
		So(bo, ShouldBeTrue)
	})
}

func TestPrivategetofferwall_urls(t *testing.T) {
	Convey("test getofferwall_urls", t, func() {
		//guard := Patch(GetMVConfigValue, func(key string) (value []byte, ifFind bool) {
		//	return []byte(key), true
		//})
		//defer guard.Unpatch()

		res, b := getofferwall_urls([]byte(""))
		var exp mvutil.COfferwallUrls
		So(res, ShouldResemble, exp)
		So(b, ShouldBeFalse)
	})
}

func TestGetSIGN_NO_CHECK_APPS(t *testing.T) {
	Convey("test GetSIGN_NO_CHECK_APPS", t, func() {
		guard := Patch(GetMVConfigValue, func(key string) (interface{}, bool) {
			return []int64{}, true
		})
		defer guard.Unpatch()
		res, bo := GetSIGN_NO_CHECK_APPS()
		So(res, ShouldResemble, []int64{})
		So(bo, ShouldBeTrue)
	})
}

func TestPrivategetSIGN_NO_CHECK_APPS(t *testing.T) {
	Convey("test getSIGN_NO_CHECK_APPS", t, func() {
		//guard := Patch(GetMVConfigValue, func(key string) (value []byte, ifFind bool) {
		//	return []byte(key), true
		//})
		//defer guard.Unpatch()

		res, b := getSIGN_NO_CHECK_APPS([]byte(""))
		var exp []int64
		So(res, ShouldResemble, exp)
		So(b, ShouldBeFalse)
	})
}

func TestGetC_TOI(t *testing.T) {
	Convey("test GetC_TOI", t, func() {
		guard := Patch(GetMVConfigValue, func(key string) (interface{}, bool) {
			return int32(99), true
		})
		defer guard.Unpatch()
		res, bo := GetC_TOI()
		So(res, ShouldResemble, int32(99))
		So(bo, ShouldBeTrue)
	})
}

func TestPrivategetCToi(t *testing.T) {
	Convey("test getCToi", t, func() {
		//guard := Patch(GetMVConfigValue, func(key string) (value []byte, ifFind bool) {
		//	return []byte(key), true
		//})
		//defer guard.Unpatch()

		res, b := getCToi([]byte(""))
		var exp int32
		So(res, ShouldResemble, exp)
		So(b, ShouldBeFalse)
	})
}

func TestGetMP_MAP_UNIT(t *testing.T) {
	Convey("test GetMP_MAP_UNIT", t, func() {
		guard := Patch(GetMVConfigValue, func(key string) (interface{}, bool) {
			return map[string]mvutil.MP_MAP_UNIT_{}, true
		})
		defer guard.Unpatch()
		res, bo := GetMP_MAP_UNIT()
		So(res, ShouldResemble, map[string]mvutil.MP_MAP_UNIT_{})
		So(bo, ShouldBeTrue)
	})
}

func TestPrivategetMP_MAP_UNIT(t *testing.T) {
	Convey("test getMP_MAP_UNIT", t, func() {
		//guard := Patch(GetMVConfigValue, func(key string) (value []byte, ifFind bool) {
		//	return []byte(key), true
		//})
		//defer guard.Unpatch()

		res, b := getMP_MAP_UNIT([]byte(""))
		exp := map[string]mvutil.MP_MAP_UNIT_{}
		So(res, ShouldResemble, exp)
		So(b, ShouldBeFalse)
	})
}

func TestGetCREATIVE_ABTEST(t *testing.T) {
	Convey("test GetCREATIVE_ABTEST", t, func() {
		guard := Patch(GetMVConfigValue, func(key string) (interface{}, bool) {
			return map[string]int{"1": 50}, true
		})
		defer guard.Unpatch()
		res, bo := GetCREATIVE_ABTEST()
		So(res, ShouldResemble, map[string]int{"1": 50})
		So(bo, ShouldBeTrue)
	})
}

func TestGetTEMPLATE_MAP(t *testing.T) {
	Convey("test GetTEMPLATE_MAP", t, func() {
		guard := Patch(GetMVConfigValue, func(key string) (interface{}, bool) {
			return mvutil.GlobalTemplateMap{
				EndScreen: map[string]string{
					"401": "hybird.rayjump.com/rv/endv4.html",
				}}, true
		})
		defer guard.Unpatch()
		res, bo := GetTEMPLATE_MAP()
		So(bo, ShouldBeTrue)
		So(res.EndScreen, ShouldContainKey, "401")
	})
}

func TestGetHBAerospikeGzipRate(t *testing.T) {
	Convey("test GetHBAerospikeGzipRate", t, func() {
		Convey("Nil Pointer", func() {
			guard := Patch(GetMVConfigValue, func(key string) (interface{}, bool) {
				return nil, true
			})
			defer guard.Unpatch()
			//
			res := GetHBAerospikeGzipRate("aws", "fk")
			So(res, ShouldEqual, 0)
		})
		Convey("Empty", func() {
			guard := Patch(GetMVConfigValue, func(key string) (interface{}, bool) {
				return &mvutil.HBAerospikeStorageConfArr{}, true
			})
			defer guard.Unpatch()
			//
			res := GetHBAerospikeGzipRate("aws", "fk")
			res2 := GetHBAerospikeRemoveRedundancyRate("aws", "fk")
			So(res, ShouldEqual, 0)
			So(res2, ShouldEqual, 0)
		})
		Convey("Find", func() {
			guard := Patch(GetMVConfigValue, func(key string) (interface{}, bool) {
				return &mvutil.HBAerospikeStorageConfArr{
					Configs: []*mvutil.HBAerospikeStorageConf{
						{
							Cloud:    "aws",
							Region:   "fk",
							GzipRate: 0.1,
						},
						{
							Cloud:    "aws",
							Region:   "se",
							GzipRate: 0.2,
						},
						{
							Cloud:    "aws",
							Region:   "vg",
							GzipRate: 0.3,
						},
						{
							Cloud:    "aws",
							Region:   "hk",
							GzipRate: 0.4,
						},
					},
				}, true
			})
			defer guard.Unpatch()
			//
			res1 := GetHBAerospikeGzipRate("aws", "fk")
			res2 := GetHBAerospikeGzipRate("aws", "se")
			res3 := GetHBAerospikeGzipRate("aws", "vg")
			res4 := GetHBAerospikeGzipRate("aws", "hk")
			res5 := GetHBAerospikeRemoveRedundancyRate("aws", "vg")
			So(res1, ShouldEqual, 0.1)
			So(res2, ShouldEqual, 0.2)
			So(res3, ShouldEqual, 0.3)
			So(res4, ShouldEqual, 0.4)
			So(res5, ShouldEqual, 0)
		})
		Convey("Member null pointer", func() {
			guard := Patch(GetMVConfigValue, func(key string) (interface{}, bool) {
				return &mvutil.HBAerospikeStorageConfArr{
					Configs: []*mvutil.HBAerospikeStorageConf{
						nil,
						nil,
					},
				}, true
			})
			defer guard.Unpatch()
			//
			res1 := GetHBAerospikeGzipRate("aws", "fk")
			res2 := GetHBAerospikeGzipRate("aws", "se")
			res3 := GetHBAerospikeGzipRate("aws", "vg")
			res4 := GetHBAerospikeGzipRate("aws", "hk")
			So(res1, ShouldEqual, 0)
			So(res2, ShouldEqual, 0)
			So(res3, ShouldEqual, 0)
			So(res4, ShouldEqual, 0)
		})
	})
}
func TestGetHBAerospikeGzipRate2(t *testing.T) {
	Convey("test getHBAerospikeGzipRate", t, func() {
		Convey("Nil Pointer", func() {
			res := getHBAerospikeStorageConfArr(nil)
			So(res, ShouldResemble, &mvutil.HBAerospikeStorageConfArr{})
		})
		Convey("Right", func() {
			jsonStr := `[{"region": "fk", "cloud": "aws", "gzip_rate": 0.1, "remove_redundancy_rate": 0.2}]`
			res := getHBAerospikeStorageConfArr([]byte(jsonStr))
			So(res, ShouldResemble, &mvutil.HBAerospikeStorageConfArr{Configs: []*mvutil.HBAerospikeStorageConf{
				{
					Cloud:                "aws",
					Region:               "fk",
					GzipRate:             0.1,
					RemoveRedundancyRate: 0.2,
				},
			}})
		})

		Convey("Empty remove rate", func() {
			jsonStr := `[{"region": "fk", "cloud": "aws", "gzip_rate": 0.1}, {"region": "vg", "cloud": "aws", "gzip_rate": 0.2, "remove_redundancy_rate": 0.3}]`
			res := getHBAerospikeStorageConfArr([]byte(jsonStr))
			So(res, ShouldResemble, &mvutil.HBAerospikeStorageConfArr{Configs: []*mvutil.HBAerospikeStorageConf{
				{
					Cloud:                "aws",
					Region:               "fk",
					GzipRate:             0.1,
					RemoveRedundancyRate: 0,
				},
				{
					Cloud:                "aws",
					Region:               "vg",
					GzipRate:             0.2,
					RemoveRedundancyRate: 0.3,
				},
			}})
		})
	})
}

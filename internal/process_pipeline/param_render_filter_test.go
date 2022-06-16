package process_pipeline

import (
	"errors"
	"math/rand"
	"strconv"
	"strings"
	"testing"
	"time"

	. "github.com/bouk/monkey"
	pf "github.com/easierway/pipefiter_framework/pipefilter"
	. "github.com/smartystreets/goconvey/convey"
	"gitlab.mobvista.com/ADN/adnet/internal/extractor"
	"gitlab.mobvista.com/ADN/adnet/internal/mvconst"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	"gitlab.mobvista.com/ADN/adnet/internal/redis"
	"gitlab.mobvista.com/ADN/adnet/internal/utility"
	smodel "gitlab.mobvista.com/ADN/structs/model"
)

func TestParseDVI(t *testing.T) {
	r := &mvutil.RequestParams{
		Param: mvutil.Params{
			DVI: "4BztYrxBYFQ3+FQ3RUE0irjeHkjTiAiPfaHtinlTioRsRrfuHoR1RUjMiahbiADAfoRsRozuYk5uRUE0GaDAinlAiaiTiaRbiAlAR0MlRr2tDBR1RUvBiavMiavMiavMioRsRretJoR1RUjBiB9MioRsRre/HBR1RURAfo9MioRsRozghdfTRUE0HbSAJdt94dl0WozghdfV4+SQRUEeWozdVcfSDFf2hrcU4ZR1RgDeiUi/iUi/fozK",
			D1:  "GaDAinlAiaiTiaRbiAlA",
			D2:  "irjeHkjTiAiPfaHtinlTiv==",
			D3:  "irjeHkjTiAiPfaHtinlTiv==",
		},
	}
	var mconfig mvutil.CommonConfig
	mconfig.DVIConfig.DVIKeys = []string{"imei", "mac", "devId", "lat", "lng", "gpst", "gpsAccuracy", "gpsType"}
	mvutil.Config = &mvutil.AdnetConfig{CommonConfig: &mconfig}
	Convey("Test parseDVI", t, func() {
		parseDVI(r)
		// parseDVIOld(r2)
		So(r.Param.AndroidID, ShouldNotBeEmpty)
		// So(r2.Param.AndroidID, ShouldEqual, r.Param.AndroidID)
	})
}

func BenchmarkParseDVI(b *testing.B) {
	r := &mvutil.RequestParams{
		Param: mvutil.Params{
			DVI: "4BztYrxBYFQ3+FQ3RUE0irjeHkjTiAiPfaHtinlTioRsRrfuHoR1RUjMiahbiADAfoRsRozuYk5uRUE0GaDAinlAiaiTiaRbiAlAR0MlRr2tDBR1RUvBiavMiavMiavMioRsRretJoR1RUjBiB9MioRsRre/HBR1RURAfo9MioRsRozghdfTRUE0HbSAJdt94dl0WozghdfV4+SQRUEeWozdVcfSDFf2hrcU4ZR1RgDeiUi/iUi/fozK",
			D1:  "GaDAinlAiaiTiaRbiAlA",
			D2:  "irjeHkjTiAiPfaHtinlTiv==",
			D3:  "irjeHkjTiAiPfaHtinlTiv==",
		},
	}
	var mconfig mvutil.CommonConfig
	mconfig.DVIConfig.DVIKeys = []string{"imei", "mac", "devId", "lat", "lng", "gpst", "gpsAccuracy", "gpsType"}
	mvutil.Config = &mvutil.AdnetConfig{CommonConfig: &mconfig}
	b.Run("ParseDVI", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			parseDVI(r)
		}
	})
	// b.Run("ParseDVIOld", func(b *testing.B) {
	// 	for i := 0; i < b.N; i++ {
	// 		//parseDVIOld(r)
	// 	}
	// })
}

func TestFormatScreenSizeV2(t *testing.T) {
	r := &mvutil.RequestParams{
		Param: mvutil.Params{
			ScreenSize: "1200.00x827.00",
		},
	}
	// r1 := &mvutil.RequestParams{
	// 	Param: mvutil.Params{
	// 		ScreenSize: "1200.00x827.00",
	// 	},
	// }
	Convey("Test formatScreenSize", t, func() {
		formatScreenSize(r)
		//formatScreenSizeOld(r1)
		So(r.Param.ScreenSize, ShouldEqual, "1200x827")
		So(r.Param.ScreenHeigh, ShouldEqual, int64(827))
		So(r.Param.ScreenWidth, ShouldEqual, int64(1200))
		// So(r1.Param.ScreenSize, ShouldEqual, r.Param.ScreenSize)
		// So(r1.Param.ScreenWidth, ShouldEqual, r.Param.ScreenWidth)
		// So(r1.Param.ScreenHeigh, ShouldEqual, r.Param.ScreenHeigh)
	})
}

func BenchmarkIntFormat(b *testing.B) {
	str := "123345"
	b.Run("ParseInt-32", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			strconv.ParseInt(str, 10, 32)
		}
	})
	b.Run("ParseInt-64", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			strconv.ParseInt(str, 10, 64)
		}
	})
	b.Run("ParseInt", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			strconv.Atoi(str)
		}
	})
}

func BenchmarkStringFind(b *testing.B) {
	str := "1200.00x827.00"
	b.Run("Contains", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			strings.Contains(str, ".")
		}
	})
	b.Run("Index", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			strings.Index(str, ".")
		}
	})
}

func BenchmarkFormatScreenSize(b *testing.B) {
	r := &mvutil.RequestParams{
		Param: mvutil.Params{
			ScreenSize: "1200x827",
		},
	}
	b.Run("FormatScreenSize", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			formatScreenSize(r)
		}
	})
	// b.Run("FormatScreenSizeNew", func(b *testing.B) {
	// 	for i := 0; i < b.N; i++ {
	// 		formatScreenSizeOld(r)
	// 	}
	// })
}

func TestFormatScreenSize(t *testing.T) {
	Convey("Test formatScreenSize", t, func() {
		r := mvutil.RequestParams{
			Param: mvutil.Params{
				ScreenSize: "1200x827",
			},
			QueryMap:      mvutil.RequestQueryMap{},
			UnitInfo:      &smodel.UnitInfo{},
			AppInfo:       &smodel.AppInfo{},
			PublisherInfo: &smodel.PublisherInfo{},
			DebugInfo:     "debug info",
		}
		formatScreenSize(&r)
		So(r.Param.ScreenSize, ShouldEqual, "1200x827")
		So(r.Param.ScreenHeigh, ShouldEqual, int64(827))
		So(r.Param.ScreenWidth, ShouldEqual, int64(1200))
	})
}

func TestHandleScenario(t *testing.T) {
	Convey("Test handleScenario", t, func() {
		var r mvutil.RequestParams

		Convey("When get scenario return empty", func() {
			guard := Patch(extractor.GetOpenapiScenario, func() ([]string, bool) {
				return []string{}, true
			})
			defer guard.Unpatch()

			r = mvutil.RequestParams{
				Param: mvutil.Params{
					ScreenSize: "1200x827",
				},
				QueryMap:      mvutil.RequestQueryMap{},
				UnitInfo:      &smodel.UnitInfo{},
				AppInfo:       &smodel.AppInfo{},
				PublisherInfo: &smodel.PublisherInfo{},
				DebugInfo:     "debug info",
			}
			handleScenario(&r)
			So(r.Param.Scenario, ShouldEqual, "openapi")
		})

		Convey("When scenario in scenarios", func() {
			guard := Patch(extractor.GetOpenapiScenario, func() ([]string, bool) {
				return []string{"param_scenario"}, true
			})
			defer guard.Unpatch()

			r = mvutil.RequestParams{
				Param: mvutil.Params{
					ScreenSize: "1200x827",
					Scenario:   "param_scenario",
				},
				QueryMap:      mvutil.RequestQueryMap{},
				UnitInfo:      &smodel.UnitInfo{},
				AppInfo:       &smodel.AppInfo{},
				PublisherInfo: &smodel.PublisherInfo{},
				DebugInfo:     "debug info",
			}
			handleScenario(&r)
			So(r.Param.Scenario, ShouldEqual, "param_scenario")
		})

		Convey("When scenario not in scenarios", func() {
			guard := Patch(extractor.GetOpenapiScenario, func() ([]string, bool) {
				return []string{"param_scenario"}, true
			})
			defer guard.Unpatch()

			r = mvutil.RequestParams{
				Param: mvutil.Params{
					ScreenSize: "1200x827",
					Scenario:   "not_in_sce",
				},
				QueryMap:      mvutil.RequestQueryMap{},
				UnitInfo:      &smodel.UnitInfo{},
				AppInfo:       &smodel.AppInfo{},
				PublisherInfo: &smodel.PublisherInfo{},
				DebugInfo:     "debug info",
			}
			handleScenario(&r)
			So(r.Param.Scenario, ShouldEqual, "openapi")
		})
	})
}

func TestHandleRankerInfo(t *testing.T) {
	Convey("Test handleRankerInfo", t, func() {
		r := mvutil.RequestParams{
			Param: mvutil.Params{
				ScreenSize:  "1200x827",
				Scenario:    "not_in_sce",
				PowerRate:   123,
				Charging:    9,
				TotalMemory: "total_mem",
				CID:         "test_cid",
			},
			QueryMap:      mvutil.RequestQueryMap{},
			UnitInfo:      &smodel.UnitInfo{},
			AppInfo:       &smodel.AppInfo{},
			PublisherInfo: &smodel.PublisherInfo{},
			DebugInfo:     "debug info",
		}
		guard := Patch(extractor.GetADNET_SWITCHS, func() (map[string]int, bool) {
			return map[string]int{}, false
		})
		defer guard.Unpatch()

		guard2 := Patch(extractor.GetFREQ_CONTROL_CONFIG, func() *mvutil.FreqControlConfig {
			return &mvutil.FreqControlConfig{}
		})
		defer guard2.Unpatch()

		handleRankerInfo(&r)
		So(r.Param.RankerInfo, ShouldEqual, "{\"power_rate\":123,\"charging\":9,\"total_memory\":\"total_mem\",\"residual_memory\":\"\",\"cid\":\"test_cid\",\"lat\":\"\",\"lng\":\"\",\"gpst\":\"\",\"gps_accuracy\":\"\",\"gps_type\":\"0\"}")
		// So(r.Param.RankerInfo, ShouldBeEmpty)
	})
}

func TestParseNums(t *testing.T) {
	Convey("Test parseNums", t, func() {
		Convey("When adtype is native 42, and tnum > 0", func() {
			r := mvutil.RequestParams{
				Param: mvutil.Params{
					ScreenSize:    "1200x827",
					Scenario:      "not_in_sce",
					PowerRate:     123,
					Charging:      9,
					TotalMemory:   "total_mem",
					CID:           "test_cid",
					AdType:        42,
					TNum:          10,
					ApiRequestNum: -2,
					AdNum:         0,
				},
				QueryMap:      mvutil.RequestQueryMap{},
				UnitInfo:      &smodel.UnitInfo{},
				AppInfo:       &smodel.AppInfo{},
				PublisherInfo: &smodel.PublisherInfo{},
				DebugInfo:     "debug info",
			}

			parseNums(&r)
			So(r.Param.TNum, ShouldEqual, 10)
			So(r.Param.AdNum, ShouldEqual, 10)
		})

		Convey("When adtype is native 42, and tnum > 0, and apiRequestNum != -2", func() {
			r := mvutil.RequestParams{
				Param: mvutil.Params{
					ScreenSize:    "1200x827",
					Scenario:      "not_in_sce",
					PowerRate:     123,
					Charging:      9,
					TotalMemory:   "total_mem",
					CID:           "test_cid",
					AdType:        42,
					TNum:          10,
					ApiRequestNum: 1,
					AdNum:         0,
				},
				QueryMap:      mvutil.RequestQueryMap{},
				UnitInfo:      &smodel.UnitInfo{},
				AppInfo:       &smodel.AppInfo{},
				PublisherInfo: &smodel.PublisherInfo{},
				DebugInfo:     "debug info",
			}

			parseNums(&r)
			So(r.Param.TNum, ShouldEqual, 10)
			So(r.Param.AdNum, ShouldEqual, 1)
		})

		Convey("When adtype is native 42, and tnum <= 0", func() {
			r := mvutil.RequestParams{
				Param: mvutil.Params{
					ScreenSize:    "1200x827",
					Scenario:      "not_in_sce",
					PowerRate:     123,
					Charging:      9,
					TotalMemory:   "total_mem",
					CID:           "test_cid",
					AdType:        42,
					TNum:          0,
					ApiRequestNum: 1,
					AdNum:         0,
				},
				QueryMap:      mvutil.RequestQueryMap{},
				UnitInfo:      &smodel.UnitInfo{},
				AppInfo:       &smodel.AppInfo{},
				PublisherInfo: &smodel.PublisherInfo{},
				DebugInfo:     "debug info",
			}

			parseNums(&r)
			So(r.Param.TNum, ShouldEqual, 1)
			So(r.Param.AdNum, ShouldEqual, 1)
		})

		Convey("When adtype is feedsvideo 95", func() {
			r := mvutil.RequestParams{
				Param: mvutil.Params{
					ScreenSize:    "1200x827",
					Scenario:      "not_in_sce",
					PowerRate:     123,
					Charging:      9,
					TotalMemory:   "total_mem",
					CID:           "test_cid",
					AdType:        95,
					TNum:          0,
					ApiRequestNum: -2,
					AdNum:         0,
				},
				QueryMap:      mvutil.RequestQueryMap{},
				UnitInfo:      &smodel.UnitInfo{},
				AppInfo:       &smodel.AppInfo{},
				PublisherInfo: &smodel.PublisherInfo{},
				DebugInfo:     "debug info",
			}

			parseNums(&r)
			So(r.Param.TNum, ShouldEqual, 1)
			So(r.Param.AdNum, ShouldEqual, 5)
		})

		Convey("When adtype is offerwall 278", func() {
			r := mvutil.RequestParams{
				Param: mvutil.Params{
					ScreenSize:    "1200x827",
					Scenario:      "not_in_sce",
					PowerRate:     123,
					Charging:      9,
					TotalMemory:   "total_mem",
					CID:           "test_cid",
					AdType:        278,
					TNum:          0,
					ApiRequestNum: 111,
					AdNum:         0,
					ApiCacheNum:   23,
				},
				QueryMap:      mvutil.RequestQueryMap{},
				UnitInfo:      &smodel.UnitInfo{},
				AppInfo:       &smodel.AppInfo{},
				PublisherInfo: &smodel.PublisherInfo{},
				DebugInfo:     "debug info",
			}

			parseNums(&r)
			So(r.Param.TNum, ShouldEqual, 23)
			So(r.Param.AdNum, ShouldEqual, 50)
		})

		Convey("When adtype is 279, and apicache > 0", func() {
			r := mvutil.RequestParams{
				Param: mvutil.Params{
					ScreenSize:    "1200x827",
					Scenario:      "not_in_sce",
					PowerRate:     123,
					Charging:      9,
					TotalMemory:   "total_mem",
					CID:           "test_cid",
					AdType:        279,
					TNum:          0,
					ApiRequestNum: 111,
					AdNum:         0,
					ApiCacheNum:   23,
				},
				QueryMap:      mvutil.RequestQueryMap{},
				UnitInfo:      &smodel.UnitInfo{},
				AppInfo:       &smodel.AppInfo{},
				PublisherInfo: &smodel.PublisherInfo{},
				DebugInfo:     "debug info",
			}

			parseNums(&r)
			So(r.Param.TNum, ShouldEqual, 23)
			So(r.Param.AdNum, ShouldEqual, 50)
		})

		Convey("When adtype is 279, and apicache <= 0", func() {
			r := mvutil.RequestParams{
				Param: mvutil.Params{
					ScreenSize:    "1200x827",
					Scenario:      "not_in_sce",
					PowerRate:     123,
					Charging:      9,
					TotalMemory:   "total_mem",
					CID:           "test_cid",
					AdType:        279,
					TNum:          0,
					ApiRequestNum: 111,
					AdNum:         0,
					ApiCacheNum:   0,
				},
				QueryMap:      mvutil.RequestQueryMap{},
				UnitInfo:      &smodel.UnitInfo{},
				AppInfo:       &smodel.AppInfo{},
				PublisherInfo: &smodel.PublisherInfo{},
				DebugInfo:     "debug info",
			}

			parseNums(&r)
			So(r.Param.TNum, ShouldEqual, 1)
			So(r.Param.AdNum, ShouldEqual, 50)
		})

		Convey("When adtype is ADTypeInterstitialVideo 287, and apicache <= 0", func() {
			r := mvutil.RequestParams{
				Param: mvutil.Params{
					ScreenSize:    "1200x827",
					Scenario:      "not_in_sce",
					PowerRate:     123,
					Charging:      9,
					TotalMemory:   "total_mem",
					CID:           "test_cid",
					AdType:        287,
					TNum:          0,
					ApiRequestNum: 222,
					AdNum:         0,
					ApiCacheNum:   0,
				},
				QueryMap:      mvutil.RequestQueryMap{},
				UnitInfo:      &smodel.UnitInfo{},
				AppInfo:       &smodel.AppInfo{},
				PublisherInfo: &smodel.PublisherInfo{},
				DebugInfo:     "debug info",
			}

			parseNums(&r)
			So(r.Param.TNum, ShouldEqual, 1)
			So(r.Param.AdNum, ShouldEqual, 50)
		})

		Convey("When adtype is other 111111, and apicache <= 0", func() {
			r := mvutil.RequestParams{
				Param: mvutil.Params{
					ScreenSize:    "1200x827",
					Scenario:      "not_in_sce",
					PowerRate:     123,
					Charging:      9,
					TotalMemory:   "total_mem",
					CID:           "test_cid",
					AdType:        111111,
					TNum:          0,
					ApiRequestNum: 222,
					AdNum:         0,
					ApiCacheNum:   0,
				},
				QueryMap:      mvutil.RequestQueryMap{},
				UnitInfo:      &smodel.UnitInfo{},
				AppInfo:       &smodel.AppInfo{},
				PublisherInfo: &smodel.PublisherInfo{},
				DebugInfo:     "debug info",
			}

			parseNums(&r)
			So(r.Param.TNum, ShouldEqual, 1)
			So(r.Param.AdNum, ShouldEqual, 5)
		})
	})
}

//func setupTest() {
//	extractor.InitMConfigInstanceOnlyForUnitTest()
//}

func TestParamRenderFilterProcess(t *testing.T) {
	//setupTest()
	Convey("test Process", t, func() {
		guard := Patch(redis.LocalRedisAlgoHGet, func(key, field string) (string, error) {
			return "", errors.New("key not found")
		})
		defer guard.Unpatch()

		guard = Patch(extractor.GetAppInfo, func(appid int64) (appInfo *smodel.AppInfo, ifFind bool) {
			return &smodel.AppInfo{}, true
		})
		defer guard.Unpatch()

		guard = Patch(extractor.GetPublisherInfo, func(publisherid int64) (publisherInfo *smodel.PublisherInfo, ifFind bool) {
			return &smodel.PublisherInfo{}, true
		})
		defer guard.Unpatch()

		guard = Patch(extractor.GetUnitInfo, func(unitid int64) (unitInfo *smodel.UnitInfo, ifFind bool) {
			return &smodel.UnitInfo{
				Unit: smodel.Unit{
					Orientation: 2,
					RecallNet:   "1;2;3",
				},
			}, true
		})
		defer guard.Unpatch()

		guard = Patch(extractor.GetFCA_SWITCH, func() (bool, bool) {
			return false, false
		})
		defer guard.Unpatch()

		guard = Patch(extractor.GetADNET_SWITCHS, func() (map[string]int, bool) {
			return map[string]int{}, false
		})
		defer guard.Unpatch()

		guard = Patch(parseDVI, func(r *mvutil.RequestParams) {
		})
		defer guard.Unpatch()

		// guard = Patch(renderPlayableFlag, func(r *mvutil.RequestParams) {
		// })
		// defer guard.Unpatch()

		guard = Patch(renderNewCreativeFlag, func(r *mvutil.RequestParams) {
			return
		})
		defer guard.Unpatch()

		guard = Patch(renderNewMoreOfferFlag, func(r *mvutil.RequestParams) {
			return
		})
		defer guard.Unpatch()

		guard = Patch(renderMoreOfferNewImp, func(r *mvutil.RequestParams) {
			return
		})
		defer guard.Unpatch()

		guard = Patch(renderMoreOfferAbFlag, func(r *mvutil.RequestParams) {
			return
		})
		defer guard.Unpatch()

		guard = Patch(changeMoreOfferAdType, func(r *mvutil.RequestParams) {
			return
		})
		defer guard.Unpatch()

		guard = Patch(handleOrientation, func(r *mvutil.RequestParams) {
		})
		defer guard.Unpatch()

		guard = Patch(formatScreenSize, func(r *mvutil.RequestParams) {
		})
		defer guard.Unpatch()

		guard = Patch(handleScreensize, func(r *mvutil.RequestParams) {
		})
		defer guard.Unpatch()

		guard = Patch(handleScenario, func(r *mvutil.RequestParams) {
		})
		defer guard.Unpatch()

		guard = Patch(handleRankerInfo, func(r *mvutil.RequestParams) {
		})
		defer guard.Unpatch()

		guard = Patch(parseNums, func(r *mvutil.RequestParams) {
		})
		defer guard.Unpatch()

		guard = Patch(renderTemplate, func(r *mvutil.RequestParams) int {
			return 9083
		})
		defer guard.Unpatch()

		guard = Patch(renderVcnABTest, func(r *mvutil.RequestParams) {
			return
		})
		defer guard.Unpatch()

		guard = Patch(extractor.GetExcludeClickPackages, func() *mvutil.ExcludeClickPackages {
			return &mvutil.ExcludeClickPackages{}
		})
		defer guard.Unpatch()

		guard = Patch(extractor.GetExcludeImpressionPackages, func() *mvutil.ExcludeClickPackages {
			return &mvutil.ExcludeClickPackages{}
		})
		defer guard.Unpatch()

		guard = Patch(extractor.GetFREQ_CONTROL_CONFIG, func() *mvutil.FreqControlConfig {
			return &mvutil.FreqControlConfig{}
		})
		defer guard.Unpatch()
		guard = Patch(extractor.GetUSE_PLACEMENT_IMP_CAP_SWITCH, func() bool {
			return true
		})
		defer guard.Unpatch()

		guard = Patch(extractor.GetSYSTEM_AREA, func() string {
			return "test_get_sys_area"
		})
		defer guard.Unpatch()

		guard = Patch(extractor.GetSYSTEM, func() string {
			return ""
		})
		defer guard.Unpatch()

		guard = Patch(extractor.GetDOMAIN_TRACK, func() string {
			return ""
		})
		defer guard.Unpatch()

		guard1 := Patch(utility.IsLowFlowUnit, func(int64, int64, int64, string) bool {
			return false
		})
		defer guard1.Unpatch()

		guard = Patch(extractor.GetDcoTestConf, func() *mvutil.TestConf {
			return &mvutil.TestConf{}
		})
		defer guard.Unpatch()

		guard = Patch(extractor.GetSspProfitDistributionRuleByUnitIdAndCountryCode, func(int64, string) (*smodel.SspProfitDistributionRule, bool) {
			return nil, false
		})
		defer guard.Unpatch()

		Convey("Given a pipeline", func() {
			filters := []pf.Filter{&ParamRenderFilter{}}

			std := pf.StraightPipeline{
				Name:    "Standard",
				Filters: &filters,
			}
			in := mvutil.RequestParams{
				Param: mvutil.Params{
					ScreenSize:  "1200x827",
					Scenario:    "not_in_sce",
					PowerRate:   123,
					Charging:    9,
					TotalMemory: "total_mem",
					CID:         "test_cid",
					RequestType: -1,
					Orientation: 1,
					NetworkType: 1,
				},
				QueryMap: mvutil.RequestQueryMap{},
				UnitInfo: &smodel.UnitInfo{
					Unit: smodel.Unit{
						Orientation: 2,
						RecallNet:   "1",
					},
				},
				AppInfo:       &smodel.AppInfo{},
				PublisherInfo: &smodel.PublisherInfo{},
				DebugInfo:     "debug info",
			}
			ret, err := std.Process(&in)

			Convey("Then RequestType = 7", func() {
				So(err, ShouldBeNil)
				So(ret.(*mvutil.RequestParams).Param.RequestType, ShouldEqual, 7)
			})
			Convey("Then Template = ", func() {
				So(err, ShouldBeNil)
				So(ret.(*mvutil.RequestParams).Param.Template, ShouldEqual, 9083)
			})
			Convey("Then Brand = ", func() {
				So(err, ShouldBeNil)
				So(ret.(*mvutil.RequestParams).Param.Brand, ShouldEqual, "")
			})
			Convey("Then ServerIP = ", func() {
				So(err, ShouldBeNil)
				So(ret.(*mvutil.RequestParams).Param.ServerIP, ShouldEqual, "")
			})
			Convey("Then CC = ", func() {
				So(err, ShouldBeNil)
				So(ret.(*mvutil.RequestParams).Param.CC, ShouldEqual, "")
			})
			Convey("Then NetworkTypeName = ", func() {
				So(err, ShouldBeNil)
				So(ret.(*mvutil.RequestParams).Param.NetworkTypeName, ShouldEqual, "unknown")
			})
			Convey("Then VideoAdType = ", func() {
				So(err, ShouldBeNil)
				So(ret.(*mvutil.RequestParams).Param.VideoAdType, ShouldEqual, 0)
			})
			Convey("Then ApiCacheNum = ", func() {
				So(err, ShouldBeNil)
				So(ret.(*mvutil.RequestParams).Param.ApiCacheNum, ShouldEqual, 0)
			})
			Convey("Then ApiRequestNum = ", func() {
				So(err, ShouldBeNil)
				So(ret.(*mvutil.RequestParams).Param.ApiRequestNum, ShouldEqual, 0)
			})
			Convey("Then UnitSize = ", func() {
				So(err, ShouldBeNil)
				So(ret.(*mvutil.RequestParams).Param.UnitSize, ShouldEqual, "")
			})
			Convey("Then ExtcdnType = ", func() {
				So(err, ShouldBeNil)
				So(ret.(*mvutil.RequestParams).Param.ExtcdnType, ShouldEqual, "test_get_sys_area")
			})
		})
	})
}

func TestHandleScreensize(t *testing.T) {
	Convey("Test handleScreensize", t, func() {
		r := mvutil.RequestParams{UnitInfo: &smodel.UnitInfo{}}

		Convey("When screen size is not right", func() {
			handleScreensize(&r)
		})

		Convey("When width >= height", func() {
			r.Param.ScreenSize = "1200x827"
			r.Param.FormatOrientation = 1
			handleScreensize(&r)
			So(r.Param.ScreenSize, ShouldEqual, "827x1200")
		})

		Convey("When width < height", func() {
			r.Param.ScreenSize = "120x827"
			r.Param.FormatOrientation = 2
			handleScreensize(&r)
			So(r.Param.ScreenSize, ShouldEqual, "827x120")
		})
	})
}

func TestRenderImageSize(t *testing.T) {
	Convey("Test renderImageSize", t, func() {
		r := mvutil.RequestParams{UnitInfo: &smodel.UnitInfo{}}

		guard := Patch(extractor.GetTemplate, func() (map[string]mvutil.Template, bool) {
			return map[string]mvutil.Template{
				"1_2_3": mvutil.Template{
					ImageSize:   "120x345",
					ImageSizeID: 123321,
					UnitSize:    "222x333",
				},
			}, true
		})
		defer guard.Unpatch()

		Convey("When unit id <= 0", func() {
			r.UnitInfo.UnitId = 0
			renderImageSize(&r)
		})

		Convey("When image size id > 0", func() {
			r = mvutil.RequestParams{
				Param: mvutil.Params{
					RequestPath:       "/mpapi/ad",
					Scenario:          "openapi_download_other_test",
					AdType:            1,
					Template:          3,
					ImageSizeID:       0,
					FormatOrientation: 2,
				},
				UnitInfo: &smodel.UnitInfo{
					UnitId: 998,
					Unit: smodel.Unit{
						Orientation: 2,
					},
				},
			}
			renderImageSize(&r)
			So(r.Param.ImageSize, ShouldEqual, "120x345")
			So(r.Param.ImageSizeID, ShouldEqual, 123321)
			So(r.Param.UnitSize, ShouldEqual, "222x333")
		})
	})
}

func TestRenderTemplate(t *testing.T) {
	Convey("Test renderTemplate", t, func() {
		r := mvutil.RequestParams{UnitInfo: &smodel.UnitInfo{}}

		Convey("When unit templates is empty", func() {
			renderTemplate(&r)
		})

		Convey("When Adtype is banner, and screen width > 320", func() {
			r = mvutil.RequestParams{
				Param: mvutil.Params{
					RequestPath: "/mpapi/ad",
					Scenario:    "openapi_download_other_test",
					AdType:      2,
					Template:    3,
					ImageSizeID: 0,
					ScreenWidth: 321,
				},
				UnitInfo: &smodel.UnitInfo{
					UnitId: 998,
					Unit: smodel.Unit{
						Orientation: 2,
						Templates:   []int{1, 2, 3, 4},
					},
				},
			}

			res := renderTemplate(&r)
			So(res, ShouldEqual, 2)
		})

		Convey("When Adtype is banner, and screen width <= 320", func() {
			r = mvutil.RequestParams{
				Param: mvutil.Params{
					RequestPath: "/mpapi/ad",
					Scenario:    "openapi_download_other_test",
					AdType:      2,
					Template:    3,
					ImageSizeID: 0,
					ScreenWidth: 320,
				},
				UnitInfo: &smodel.UnitInfo{
					UnitId: 998,
					Unit: smodel.Unit{
						Orientation: 2,
						Templates:   []int{4},
					},
				},
			}

			res := renderTemplate(&r)
			So(res, ShouldEqual, 4)
		})

		Convey("When Adtype is native", func() {
			r = mvutil.RequestParams{
				Param: mvutil.Params{
					RequestPath: "/mpapi/ad",
					Scenario:    "openapi_download_other_test",
					AdType:      42,
					Template:    3,
					ImageSizeID: 0,
					ScreenWidth: 320,
					NativeInfoList: []mvutil.NativeInfoEntry{
						mvutil.NativeInfoEntry{
							AdTemplate: 9,
						},
					},
				},
				UnitInfo: &smodel.UnitInfo{
					UnitId: 998,
					Unit: smodel.Unit{
						Orientation: 2,
						Templates:   []int{4},
					},
				},
			}

			res := renderTemplate(&r)
			So(res, ShouldEqual, 9)
		})
	})
}

func TestHandleOrientation(t *testing.T) {
	Convey("Test handleOrientation", t, func() {
		r := mvutil.RequestParams{
			Param: mvutil.Params{
				RequestPath: "/mpapi/ad",
				Scenario:    "openapi_download_other_test",
				AdType:      42,
				Template:    3,
				ImageSizeID: 0,
				ScreenWidth: 320,
				NativeInfoList: []mvutil.NativeInfoEntry{
					mvutil.NativeInfoEntry{
						AdTemplate: 9,
					},
				},
				Orientation: 1,
			},
			UnitInfo: &smodel.UnitInfo{
				UnitId: 998,
				Unit: smodel.Unit{
					Orientation: 0,
					Templates:   []int{4},
				},
			},
		}

		handleOrientation(&r)
		So(r.Param.FormatOrientation, ShouldEqual, 1)
	})
}

func TestChangeMoreOfferAdType(t *testing.T) {
	Convey("Test changeMoreOfferAdType", t, func() {
		Convey("mof条件不符合", func() {
			r := mvutil.RequestParams{
				Param: mvutil.Params{
					Mof:    0,
					AdType: 3,
					UnitID: 12345,
				},
			}
			guard := Patch(extractor.GetAPPWALL_TO_MORE_OFFER_UNIT, func() (map[string]int32, bool) {
				return map[string]int32{
					"12345": 100,
				}, true
			})
			defer guard.Unpatch()
			guardRand := Patch(rand.Intn, func(n int) int {
				return 0
			})
			defer guardRand.Unpatch()
			changeMoreOfferAdType(&r)
			So(r.Param.AdType, ShouldEqual, 3)
		})

		Convey("adtype条件不符合", func() {
			r := mvutil.RequestParams{
				Param: mvutil.Params{
					Mof:    1,
					AdType: 4,
					UnitID: 12345,
				},
			}
			guard := Patch(extractor.GetAPPWALL_TO_MORE_OFFER_UNIT, func() (map[string]int32, bool) {
				return map[string]int32{
					"12345": 100,
				}, true
			})
			defer guard.Unpatch()
			guardRand := Patch(rand.Intn, func(n int) int {
				return 0
			})
			defer guardRand.Unpatch()
			changeMoreOfferAdType(&r)
			So(r.Param.AdType, ShouldEqual, 4)
		})

		Convey("配置条件不符合", func() {
			r := mvutil.RequestParams{
				Param: mvutil.Params{
					Mof:    1,
					AdType: 3,
					UnitID: 12345,
				},
			}
			guard := Patch(extractor.GetAPPWALL_TO_MORE_OFFER_UNIT, func() (map[string]int32, bool) {
				return map[string]int32{
					"123456": 100,
				}, true
			})
			defer guard.Unpatch()
			guardRand := Patch(rand.Intn, func(n int) int {
				return 0
			})
			defer guardRand.Unpatch()
			changeMoreOfferAdType(&r)
			So(r.Param.AdType, ShouldEqual, 3)
		})

		Convey("条件均符合", func() {
			r := mvutil.RequestParams{
				Param: mvutil.Params{
					Mof:    1,
					AdType: 3,
					UnitID: 12345,
				},
			}
			guard := Patch(extractor.GetAPPWALL_TO_MORE_OFFER_UNIT, func() (map[string]int32, bool) {
				return map[string]int32{
					"12345": 100,
				}, true
			})
			defer guard.Unpatch()
			guardRand := Patch(rand.Intn, func(n int) int {
				return 0
			})
			defer guardRand.Unpatch()
			changeMoreOfferAdType(&r)
			So(r.Param.AdType, ShouldEqual, 295)
		})

	})
}

func TestRenderPriceFactor(t *testing.T) {
	Convey("Test renderPriceFactor", t, func() {
		Convey("request type 条件不符合", func() {
			r := mvutil.RequestParams{
				Param: mvutil.Params{
					RequestType: mvconst.REQUEST_TYPE_OPENAPI,
				},
			}
			renderPriceFactor(&r)
			So(r.Param.ExtDataInit.PriceFactor, ShouldEqual, 0)
			So(r.Param.ExtDataInit.PriceFactorHit, ShouldEqual, 0)
			So(r.Param.ExtDataInit.PriceFactorTag, ShouldEqual, 0)
			So(r.Param.ExtDataInit.PriceFactorGroupName, ShouldEqual, "")
			So(r.Param.ExtDataInit.PriceFactorFreq, ShouldEqual, nil)
		})

		Convey("配置 Status 不是1 ", func() {
			guard := Patch(extractor.GetFREQ_CONTROL_CONFIG, func() *mvutil.FreqControlConfig {
				return &mvutil.FreqControlConfig{
					FreqControlToRs: 0,
					Status:          0,
					Rules:           nil,
				}
			})
			defer guard.Unpatch()
			r := mvutil.RequestParams{
				Param: mvutil.Params{
					RequestType: mvconst.REQUEST_TYPE_OPENAPI_V3,
				},
			}
			renderPriceFactor(&r)
			So(r.Param.ExtDataInit.PriceFactor, ShouldEqual, 0)
			So(r.Param.ExtDataInit.PriceFactorHit, ShouldEqual, 0)
			So(r.Param.ExtDataInit.PriceFactorTag, ShouldEqual, 0)
			So(r.Param.ExtDataInit.PriceFactorGroupName, ShouldEqual, "")
			So(r.Param.ExtDataInit.PriceFactorFreq, ShouldEqual, nil)
		})

		Convey("正确值", func() {

			guard := Patch(extractor.GetFREQ_CONTROL_CONFIG, func() *mvutil.FreqControlConfig {
				return &mvutil.FreqControlConfig{
					FreqControlToRs: 1,
					Status:          1,
					Rules: []*mvutil.FreqControlRule{
						{
							Filter: &mvutil.FreqControlFilter{
								Platform: &mvutil.FreqControlFilterItem{
									Op:    "in",
									Value: []string{"ios", "android"},
								},
							},
							Groups: []*mvutil.FreqControlGroupItem{
								{
									GroupName: "AAA",
									GroupRate: 100,
									Rate:      100,
									SubRate:   0,

									TimeWindow: &mvutil.FreqControlGroupsItemTimeWindow{
										Mode:      1,
										WindowSec: 86400,
									},
									FreqControl: []*mvutil.FreqControlGroupsItemFreqControlItem{
										{
											Min:  0,
											Max:  1,
											Keys: []string{"AERASPIKE_1"},
										},
										{
											Min:  1,
											Max:  10000,
											Keys: []string{"AERASPIKE_2"},
										},
									},
								},
							},
						},
					},
				}
			})
			defer guard.Unpatch()
			guard2 := Patch(getPriceFactorFreqByMkv, func(r *mvutil.RequestParams, group *mvutil.FreqControlGroupItem, devicId string, adType int32) (*mvutil.FreqControlMkvData, error) {
				return &mvutil.FreqControlMkvData{
					Imp: []*mvutil.FreqControlMkvDataItem{
						{
							Ts:   int(time.Now().Unix()),
							Freq: 10.333,
						},
					},
				}, nil
			})
			defer guard2.Unpatch()
			guard3 := Patch(extractor.GetFreqControlPriceFactor, func(key string) (freqControlFactor *smodel.FreqControlFactor, ifFind bool) {
				return &smodel.FreqControlFactor{
					FactorRate: 1.3,
				}, true
			})
			defer guard3.Unpatch()
			r := &mvutil.RequestParams{
				Param: mvutil.Params{
					RequestType:  mvconst.REQUEST_TYPE_OPENAPI_V3,
					Platform:     mvconst.PlatformAndroid,
					PlatformName: mvconst.PlatformNameAndroid,
					GAID:         "12345678-2211-ABCD-0000-000000000000",
				},
			}

			renderPriceFactor(r)
			So(r.Param.ExtDataInit.PriceFactor, ShouldEqual, 1.3)
			So(r.Param.ExtDataInit.PriceFactorHit, ShouldEqual, 1)
			So(r.Param.ExtDataInit.PriceFactorTag, ShouldEqual, 2)
			So(r.Param.ExtDataInit.PriceFactorGroupName, ShouldEqual, "AAA")
			So(*r.Param.ExtDataInit.PriceFactorFreq, ShouldEqual, 10)
		})

	})
}

func TestGetFreqControlReplaceKeys(t *testing.T) {
	Convey("Test getFreqControlReplaceKeys，配置时间大于当前时间", t, func() {
		r := &mvutil.RequestParams{
			Param: mvutil.Params{
				AdType: 22,
			},
		}
		m := getFreqControlReplaceKeys(r)

		So(m["{ad_type}"], ShouldEqual, "22")
		So(m["{is_hb}"], ShouldEqual, "2")
		So(m["{has_devid}"], ShouldEqual, "2")
	})
}

func TestGetFreqControlValue(t *testing.T) {

	guard := Patch(extractor.GetFreqControlPriceFactor,
		func(key string) (freqControlFactor *smodel.FreqControlFactor, ifFind bool) {
			if key == "A2" {
				return &smodel.FreqControlFactor{
					FactorRate: 1.5,
				}, true
			}
			return nil, false
		})
	defer guard.Unpatch()
	fcg := []*mvutil.FreqControlGroupsItemFreqControlItem{
		&mvutil.FreqControlGroupsItemFreqControlItem{
			Min:  0,
			Max:  1,
			Keys: []string{"A1", "A2"},
		},
		&mvutil.FreqControlGroupsItemFreqControlItem{
			Min:  1,
			Max:  5,
			Keys: []string{"B1", "B2"},
		},
		&mvutil.FreqControlGroupsItemFreqControlItem{
			Min:  5,
			Max:  10,
			Keys: []string{"C1", "C2"},
		},
	}
	rpKeys := map[string]string{}

	Convey("Test getFreqControlValue，不命中", t, func() {
		pff := 10
		m, hit := getFreqControlValue(pff, fcg, rpKeys)
		So(m, ShouldEqual, 1)
		So(hit, ShouldEqual, false)
	})
	Convey("Test getFreqControlValue，命中, 单找不到对应的rate", t, func() {
		pff := 4
		m, hit := getFreqControlValue(pff, fcg, rpKeys)
		So(m, ShouldEqual, 1)
		So(hit, ShouldEqual, true)
	})
	Convey("Test getFreqControlValue，命中且找到对应的rate", t, func() {
		pff := 0
		m, hit := getFreqControlValue(pff, fcg, rpKeys)
		So(m, ShouldEqual, 1.5)
		So(hit, ShouldEqual, true)
	})
}

func TestGetAbTestTag(t *testing.T) {

	Convey("Test getAbTestTag，不命中,结果A", t, func() {
		group := &mvutil.FreqControlGroupItem{
			Rate:    0,
			SubRate: 50,
		}
		m := getAbTestTag("AAAA", group)
		So(m, ShouldEqual, 1)
	})
	Convey("Test getAbTestTag，不命中,结果B", t, func() {
		group := &mvutil.FreqControlGroupItem{
			Rate:    100,
			SubRate: 0,
		}
		m := getAbTestTag("AAAA", group)
		So(m, ShouldEqual, 2)
	})
	Convey("Test getAbTestTag，不命中,结果B", t, func() {
		group := &mvutil.FreqControlGroupItem{
			Rate:    100,
			SubRate: 100,
		}
		m := getAbTestTag("AAAA", group)
		So(m, ShouldEqual, 3)
	})
	Convey("Test getAbTestTag，不命中,结果B", t, func() {
		group := &mvutil.FreqControlGroupItem{
			Rate:    50,
			SubRate: 50,
		}

		m := getAbTestTag("AAA", group)
		So(m, ShouldEqual, 3)
		m = getAbTestTag("BBB", group)
		So(m, ShouldEqual, 2)
		m = getAbTestTag("CCC", group)
		So(m, ShouldEqual, 1)
	})
}

func TestGetFreqControlStartTime(t *testing.T) {

	timeFormat := "2006-01-02 15:04:05"
	beijing, _ := time.LoadLocation("Asia/Shanghai")

	guard3 := Patch(extractor.GetTIMEZONE_CONFIG, func() map[string]int {
		return map[string]int{
			"GMT+01:00": 1,
			"GMT+02:00": 2,
			"GMT+08:00": 2,
		}
	})
	defer guard3.Unpatch()
	guard4 := Patch(extractor.GetCOUNTRY_CODE_TIMEZONE_CONFIG, func() map[string]int {
		return map[string]int{

			"CN": 8,
		}
	})
	defer guard4.Unpatch()

	Convey("Test getFreqControlStartTime，配置时间大于当前时间", t, func() {
		tw := &mvutil.FreqControlGroupsItemTimeWindow{
			Mode:      2,
			StartHour: 10,
		}

		t, _ := time.ParseInLocation(timeFormat, "2020-03-23 23:15:17", beijing)

		startTime, _ := getFreqControlStartTime(tw, t, "GMT+01:00", "CN")
		tt, _ := time.ParseInLocation(timeFormat, "2020-03-23 17:00:00", beijing)
		resultTime := tt.Unix()

		So(startTime, ShouldEqual, resultTime)
	})
	Convey("Test getFreqControlStartTime，配置时间大于当前时间", t, func() {
		tw := &mvutil.FreqControlGroupsItemTimeWindow{
			Mode:      2,
			StartHour: 5,
		}
		t, _ := time.Parse(timeFormat, "2020-03-23 12:15:17")
		startTime, _ := getFreqControlStartTime(tw, t, "GMT+01:00", "CN")
		// 结果
		tt, _ := time.ParseInLocation(timeFormat, "2020-03-23 12:00:00", beijing)
		resultTime := tt.Unix()

		So(startTime, ShouldEqual, resultTime)
	})
	Convey("Test getFreqControlStartTime，配置时间小于当前时间", t, func() {
		tw := &mvutil.FreqControlGroupsItemTimeWindow{
			Mode:      2,
			StartHour: 6,
		}
		t, _ := time.ParseInLocation(timeFormat, "2020-03-23 12:15:17", beijing)
		startTime, _ := getFreqControlStartTime(tw, t, "GMT+01:00", "CN")
		// 结果
		tt, _ := time.ParseInLocation(timeFormat, "2020-03-22 13:00:00", beijing)
		resultTime := tt.Unix()

		So(startTime, ShouldEqual, resultTime)
	})

	Convey("Test getFreqControlStartTime，没有timezone，使用国家维度", t, func() {
		tw := &mvutil.FreqControlGroupsItemTimeWindow{
			Mode:      2,
			StartHour: 6,
		}
		t, _ := time.ParseInLocation(timeFormat, "2020-03-23 6:15:17", beijing)
		startTime, _ := getFreqControlStartTime(tw, t, "", "CN")
		// 结果
		tt, _ := time.ParseInLocation(timeFormat, "2020-03-23 6:00:00", beijing)
		resultTime := tt.Unix()

		So(startTime, ShouldEqual, resultTime)
	})
}

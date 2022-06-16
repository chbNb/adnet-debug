package output

import (
	supply_mvutil "gitlab.mobvista.com/ADN/chasm/module/supply/mvutil"
	"testing"

	. "github.com/bouk/monkey"
	. "github.com/smartystreets/goconvey/convey"
	"gitlab.mobvista.com/ADN/adnet/internal/extractor"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
)

func TestIsReturnEndcard(t *testing.T) {
	Convey("test IsReturnEndcard", t, func() {
		r := mvutil.RequestParams{}
		var res bool

		guard := Patch(extractor.GetVersionCompare, func() (map[string]mvutil.VersionCompare, bool) {
			return map[string]mvutil.VersionCompare{
				"OS_VERSION_ENDCARD": mvutil.VersionCompare{
					Android: mvutil.VersionCompareItem{
						Version:        int32(111),
						ExcludeVersion: []int32{2, 4},
					},
					IOS: mvutil.VersionCompareItem{
						Version:        int32(222),
						ExcludeVersion: []int32{5, 7},
					},
				},
			}, true
		})
		defer guard.Unpatch()

		Convey("len(r.Param.OSVersion) <= 0", func() {
			r.Param.OSVersion = ""
			res = IsReturnEndcard(&r)
			So(res, ShouldBeFalse)
		})

		Convey("is iPad", func() {
			r.Param.Model = "ipad"
			res = IsReturnEndcard(&r)
			So(res, ShouldBeFalse)
		})

		Convey("非 ios Android", func() {
			r.Param.Platform = 100
			res = IsReturnEndcard(&r)
			So(res, ShouldBeFalse)
		})

		Convey("非 ios Android 数据", func() {
			r.Param.Platform = 1
			r.Param.Model = "iPhone"
			res = IsReturnEndcard(&r)
			So(res, ShouldBeFalse)
		})
	})
}

func TestCompare(t *testing.T) {
	Convey("test Compare", t, func() {
		Convey("sdk type not allow", func() {
			guard := Patch(supply_mvutil.RenderSDKVersion, func(sdkversion string) supply_mvutil.SDKVersionItem {
				return supply_mvutil.SDKVersionItem{
					SDKType:        "test",
					SDKNumber:      "sdk_num",
					SDKVersionCode: int32(123),
				}
			})
			defer guard.Unpatch()

			r := mvutil.RequestParams{
				Param: mvutil.Params{
					Platform:   1110,
					SDKVersion: "sdk_version",
				},
			}
			compareType := "test_1"
			res := Compare(&r, compareType)
			So(res, ShouldBeFalse)
		})

		Convey("is not android && is not ios", func() {
			guard := Patch(supply_mvutil.RenderSDKVersion, func(sdkversion string) supply_mvutil.SDKVersionItem {
				return supply_mvutil.SDKVersionItem{
					SDKType:        "mi",
					SDKNumber:      "sdk_num",
					SDKVersionCode: int32(123),
				}
			})
			defer guard.Unpatch()

			guard = Patch(extractor.GetVersionCompare, func() (map[string]mvutil.VersionCompare, bool) {
				return map[string]mvutil.VersionCompare{
					"test_1": mvutil.VersionCompare{
						Android: mvutil.VersionCompareItem{},
						IOS:     mvutil.VersionCompareItem{},
					},
				}, true
			})
			defer guard.Unpatch()

			r := mvutil.RequestParams{
				Param: mvutil.Params{
					Platform:   1111,
					SDKVersion: "sdk_version",
				},
			}
			compareType := "test_1"
			res := Compare(&r, compareType)
			So(res, ShouldBeFalse)
		})

		Convey("platform is android, version <= 0", func() {
			guard := Patch(supply_mvutil.RenderSDKVersion, func(sdkversion string) supply_mvutil.SDKVersionItem {
				return supply_mvutil.SDKVersionItem{
					SDKType:        "mi",
					SDKNumber:      "sdk_num",
					SDKVersionCode: int32(123),
				}
			})
			defer guard.Unpatch()

			guard = Patch(extractor.GetVersionCompare, func() (map[string]mvutil.VersionCompare, bool) {
				return map[string]mvutil.VersionCompare{
					"test_1": mvutil.VersionCompare{
						Android: mvutil.VersionCompareItem{},
						IOS:     mvutil.VersionCompareItem{},
					},
				}, true
			})
			defer guard.Unpatch()

			r := mvutil.RequestParams{
				Param: mvutil.Params{
					Platform:   1,
					SDKVersion: "sdk_version",
				},
			}
			compareType := "test_1"
			res := Compare(&r, compareType)
			So(res, ShouldBeFalse)
		})

		Convey("platform is android, version > 0", func() {
			guard := Patch(supply_mvutil.RenderSDKVersion, func(sdkversion string) supply_mvutil.SDKVersionItem {
				return supply_mvutil.SDKVersionItem{
					SDKType:        "mi",
					SDKNumber:      "sdk_num",
					SDKVersionCode: int32(123),
				}
			})
			defer guard.Unpatch()

			guard = Patch(extractor.GetVersionCompare, func() (map[string]mvutil.VersionCompare, bool) {
				return map[string]mvutil.VersionCompare{
					"test_1": mvutil.VersionCompare{
						Android: mvutil.VersionCompareItem{
							Version:        123,
							ExcludeVersion: []int32{1, 2},
						},
						IOS: mvutil.VersionCompareItem{},
					},
				}, true
			})
			defer guard.Unpatch()

			guard = Patch(mvutil.GetVersionCode, func(version string) int32 {
				return int32(11111)
			})
			defer guard.Unpatch()

			r := mvutil.RequestParams{
				Param: mvutil.Params{
					Platform:   1,
					SDKVersion: "sdk_version",
				},
			}
			compareType := "test_1"
			res := Compare(&r, compareType)
			So(res, ShouldBeTrue)
		})
	})
}

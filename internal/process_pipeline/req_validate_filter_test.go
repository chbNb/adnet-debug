package process_pipeline

import (
	smodel "gitlab.mobvista.com/ADN/structs/model"
	"reflect"
	"testing"

	. "github.com/bouk/monkey"
	pf "github.com/easierway/pipefiter_framework/pipefilter"
	. "github.com/smartystreets/goconvey/convey"
	"gitlab.mobvista.com/ADN/adnet/internal/extractor"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
)

func TestReqValidateFilterProcess(t *testing.T) {
	Convey("test Process", t, func() {
		guard := Patch(mvutil.GetRandConsiderZero, func(gaid string, idfa string, salt string, randSum int) int {
			return 100
		})
		defer guard.Unpatch()

		var u *smodel.UnitInfo
		guardIns := PatchInstanceMethod(reflect.TypeOf(u), "GetAdSourceID", func(_ *smodel.UnitInfo, _ string, _ int) int {
			return 10086
		})
		defer guardIns.Unpatch()

		Convey("Given a pipeline", func() {
			filters := []pf.Filter{&ReqValidateFilter{}}

			std := pf.StraightPipeline{
				Name:    "Standard",
				Filters: &filters,
			}

			Convey("When error data", func() {
				in := "error_data"
				ret, err := std.Process(&in)
				So(err, ShouldBeError)
				So(ret, ShouldBeNil)
			})

			Convey("When data is mvutil.RequestParams", func() {
				in := mvutil.RequestParams{
					Param: mvutil.Params{
						RequestPath: "/mpapi/ad",
						Scenario:    "openapi_download_other_test",
						AdType:      42,
						Template:    3,
						ImageSizeID: 0,
						AppID:       101,
						CountryCode: "CN",
						AdSourceID:  1,
						Platform:    1,
						Category:    250,
						SDKVersion:  "",
						Sign:        "NO_CHECK_SIGN",
						Debug:       1,
					},
					UnitInfo: &smodel.UnitInfo{
						AppId:  101,
						UnitId: 998,
						Unit: smodel.Unit{
							Orientation: 2,
							Status:      1,
							AdType:      42,
						},
						AdSourceCountry: map[string]int{"CN": 1, "US": 2},
					},
					AppInfo: &smodel.AppInfo{
						App: smodel.App{
							Platform: 1,
						},
					},
					PublisherInfo: &smodel.PublisherInfo{
						Publisher: smodel.Publisher{
							Type: 0,
						},
					},
				}
				guard = Patch(extractor.GetBAD_REQUEST_FILTER_CONF, func() []*mvutil.BadRequestFilterConf {
					return []*mvutil.BadRequestFilterConf{}
				})
				defer guard.Unpatch()
				ret, err := std.Process(&in)
				So(err, ShouldBeNil)
				So(ret.(*mvutil.ReqCtx).ReqParams.Param.RequestPath, ShouldEqual, "/mpapi/ad")
				So(ret.(*mvutil.ReqCtx).ReqParams.Param.AdType, ShouldEqual, 42)
				So(ret.(*mvutil.ReqCtx).ReqParams.Param.DeviceType, ShouldEqual, 4)
			})
		})
	})
}

func TestNeedCheckSign(t *testing.T) {
	Convey("Test needCheckSign", t, func() {
		r := mvutil.RequestParams{Param: mvutil.Params{}}

		guard := Patch(extractor.GetSIGN_NO_CHECK_APPS, func() ([]int64, bool) {
			return []int64{int64(123)}, true
		})
		defer guard.Unpatch()

		Convey("When no need check sign", func() {
			r.Param.Sign = "NO_CHECK_SIGN"
			res := needCheckSign(&r)
			So(res, ShouldBeFalse)
		})

		Convey("When app is in no check list", func() {
			r.Param.Sign = ""
			r.Param.AppID = int64(123)
			res := needCheckSign(&r)
			So(res, ShouldBeFalse)
		})

		Convey("When app is not in no check list", func() {
			r.Param.Sign = ""
			r.Param.AppID = int64(2)
			res := needCheckSign(&r)
			So(res, ShouldBeTrue)
		})

		Convey("When no check list", func() {
			guard = Patch(extractor.GetSIGN_NO_CHECK_APPS, func() ([]int64, bool) {
				return []int64{}, false
			})
			defer guard.Unpatch()

			r.Param.Sign = ""
			r.Param.AppID = int64(2)
			res := needCheckSign(&r)
			So(res, ShouldBeTrue)
		})
	})
}

func TestCheckSystemErr(t *testing.T) {
	Convey("Test checkSystemErr", t, func() {
		r := mvutil.RequestParams{PublisherInfo: &smodel.PublisherInfo{}}

		Convey("When pSystem = 0", func() {
			r.PublisherInfo.Publisher.Type = 0

			res := checkSystemErr(&r)
			So(res, ShouldBeFalse)
		})

		Convey("When get system is test", func() {
			guard := Patch(extractor.GetSYSTEM, func() string {
				return "TEST"
			})
			defer guard.Unpatch()

			r.PublisherInfo.Publisher.Type = 1

			res := checkSystemErr(&r)
			So(res, ShouldBeFalse)
		})

		Convey("When get system is M", func() {
			guard := Patch(extractor.GetSYSTEM, func() string {
				return "M"
			})
			defer guard.Unpatch()

			r.PublisherInfo.Publisher.Type = 1

			res := checkSystemErr(&r)
			So(res, ShouldBeFalse)
		})

		Convey("When get system is SA", func() {
			guard := Patch(extractor.GetSYSTEM, func() string {
				return "M"
			})
			defer guard.Unpatch()

			r.PublisherInfo.Publisher.Type = 4

			res := checkSystemErr(&r)
			So(res, ShouldBeTrue)
		})

		Convey("When get system is MP", func() {
			guard := Patch(extractor.GetSYSTEM, func() string {
				return "M"
			})
			defer guard.Unpatch()

			r.PublisherInfo.Publisher.Type = 6

			res := checkSystemErr(&r)
			So(res, ShouldBeTrue)
		})
	})
}

func TestCanRequestWithoutDevId(t *testing.T) {
	Convey("Test canRequestWithoutDevId", t, func() {
		r := mvutil.RequestParams{UnitInfo: &smodel.UnitInfo{}, AppInfo: &smodel.AppInfo{}}

		Convey("When iscan = 0", func() {
			r.UnitInfo.Unit.DevIdAllowNull = 0
			r.AppInfo.App.DevIdAllowNull = 2
			res := canRequestWithoutDevId(&r)
			So(res, ShouldBeFalse)
		})

		Convey("When iscan = 2", func() {
			r.UnitInfo.Unit.DevIdAllowNull = 0
			r.AppInfo.App.DevIdAllowNull = 1
			res := canRequestWithoutDevId(&r)
			So(res, ShouldBeTrue)
		})

		Convey("When iscan = 10", func() {
			r.UnitInfo.Unit.DevIdAllowNull = 0
			r.AppInfo.App.DevIdAllowNull = 1
			res := canRequestWithoutDevId(&r)
			So(res, ShouldBeTrue)
		})
	})
}

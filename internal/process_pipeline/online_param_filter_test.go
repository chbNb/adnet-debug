package process_pipeline

import (
	smodel "gitlab.mobvista.com/ADN/structs/model"
	"testing"

	"gitlab.mobvista.com/ADN/adnet/internal/mvconst"

	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"

	"gitlab.mobvista.com/ADN/adnet/internal/extractor"

	. "github.com/bouk/monkey"
	pf "github.com/easierway/pipefiter_framework/pipefilter"
	. "github.com/smartystreets/goconvey/convey"
)

func TestOnlineParamFilterProcess(t *testing.T) {
	Convey("test Process", t, func() {

		Convey("Given a pipeline", func() {
			filters := []pf.Filter{&OnlineParamFilter{}}

			std := pf.StraightPipeline{
				Name:    "Standard",
				Filters: &filters,
			}

			Convey("data illegal", func() {
				in := "error_data"
				ret, err := std.Process(&in)

				Convey("Then get the excepted result", func() {
					So(err, ShouldBeError)
					So(ret, ShouldBeNil)
				})
			})

			Convey("When param is mvutil.RequestParams", func() {
				in := mvutil.RequestParams{
					Param: mvutil.Params{
						AdNum: int32(123),
					},
				}
				guard1 := Patch(extractor.GetHAS_EXTENSIONS_UNIT, func() ([]int64, bool) {
					return []int64{23, 45, 678}, true
				})
				defer guard1.Unpatch()

				guard := Patch(extractor.GetDEL_PRICE_FLOOR_UNIT, func() ([]int64, bool) {
					return []int64{123}, true
				})
				defer guard.Unpatch()
				noCheckAppGuard := Patch(extractor.GetNO_CHECK_PARAM_APP, func() (mvutil.NO_CHECK_PARAM_APP, bool) {
					status := 0
					appIds := []int64{12345}
					return mvutil.NO_CHECK_PARAM_APP{
						Status: &status,
						AppIds: &appIds,
					}, true
				})
				defer noCheckAppGuard.Unpatch()

				guard = Patch(extractor.GetNEW_AFREECATV_UNIT, func() []int64 {
					return []int64{}
				})
				defer guard.Unpatch()
				ret, err := std.Process(&in)

				Convey("Then get the excepted result", func() {
					So(err, ShouldBeNil)
					So((*(ret.(*mvutil.RequestParams))).Param.RequestType, ShouldEqual, 10)
					So((*(ret.(*mvutil.RequestParams))).Param.TNum, ShouldEqual, 123)
					So((*(ret.(*mvutil.RequestParams))).Param.AdNum, ShouldEqual, 40)
					So((*(ret.(*mvutil.RequestParams))).Param.OnlyImpression, ShouldEqual, 1)
				})
			})
		})
	})
}

func TestRenderOffset(t *testing.T) {
	Convey("test renderOffset", t, func() {
		r := mvutil.RequestParams{PublisherInfo: &smodel.PublisherInfo{}}

		Convey("When param.Offset > 0", func() {
			r.Param.Offset = 10
			renderOffset(&r)
			So(r.Param.Offset, ShouldEqual, 10)
		})

		Convey("When param.Offset <= 0 & len(r.PublisherInfo.OffsetList) <= 0", func() {
			r.Param.Offset = 0
			r.PublisherInfo.OffsetList = map[string]int32{}
			renderOffset(&r)
			So(r.Param.Offset, ShouldEqual, 0)
		})

		Convey("When param.Offset <= 0 & len(r.PublisherInfo.OffsetList) > 0", func() {
			guard := Patch(mvutil.RandByRate, func(rateMap map[int]int) int {
				return 110
			})
			defer guard.Unpatch()

			r.Param.Scenario = mvconst.SCENARIO_OPENAPI
			r.Param.Offset = int32(0)
			r.Param.Scenario = mvconst.SCENARIO_OPENAPI
			r.PublisherInfo.OffsetList = map[string]int32{
				"2": int32(9),
				"3": int32(8),
			}
			renderOffset(&r)
			So(r.Param.Offset, ShouldEqual, 110)
		})
	})
}

func TestCheckParams(t *testing.T) {
	Convey("测试online api设备参数强校验", t, func() {
		var params mvutil.RequestParams

		Convey("开关不开的情况", func() {
			noCheckAppGuard := Patch(extractor.GetNO_CHECK_PARAM_APP, func() (mvutil.NO_CHECK_PARAM_APP, bool) {
				status := 0
				appIds := []int64{12345}
				return mvutil.NO_CHECK_PARAM_APP{
					Status: &status,
					AppIds: &appIds,
				}, true
			})
			defer noCheckAppGuard.Unpatch()
			res := checkParams(&params)
			So(res, ShouldBeNil)
		})
		Convey("开关开启，没有配置appid的情况", func() {
			noCheckAppGuard := Patch(extractor.GetNO_CHECK_PARAM_APP, func() (mvutil.NO_CHECK_PARAM_APP, bool) {
				status := 1
				return mvutil.NO_CHECK_PARAM_APP{
					Status: &status,
				}, true
			})
			defer noCheckAppGuard.Unpatch()
			res := checkParams(&params)
			So(res, ShouldBeNil)
		})
		Convey("开关开启，配置appid的情况，且命中的情况", func() {
			noCheckAppGuard := Patch(extractor.GetNO_CHECK_PARAM_APP, func() (mvutil.NO_CHECK_PARAM_APP, bool) {
				status := 1
				appIds := []int64{12345}
				return mvutil.NO_CHECK_PARAM_APP{
					Status: &status,
					AppIds: &appIds,
				}, true
			})
			defer noCheckAppGuard.Unpatch()
			params.Param.AppID = 12345
			res := checkParams(&params)
			So(res, ShouldBeNil)
		})
	})
}

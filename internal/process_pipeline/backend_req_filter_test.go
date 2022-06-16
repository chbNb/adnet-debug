package process_pipeline

import (
	"reflect"
	"testing"

	"gitlab.mobvista.com/ADN/adnet/internal/watcher"

	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"

	"gitlab.mobvista.com/ADN/adnet/internal/backend"

	. "github.com/bouk/monkey"
	pf "github.com/easierway/pipefiter_framework/pipefilter"
	. "github.com/smartystreets/goconvey/convey"
)

func TestBackendReqFilterProcess(t *testing.T) {
	Convey("test Process", t, func() {
		var manager *backend.BackendManage
		guardIns := PatchInstanceMethod(reflect.TypeOf(manager), "GetAds", func(_ *backend.BackendManage, _ *mvutil.ReqCtx) (map[int]*mvutil.BackendMetric, error) {
			return nil, nil
		})
		defer guardIns.Unpatch()

		guard := Patch(watcher.AddWatchValue, func(key string, value float64) {
		})
		defer guard.Unpatch()

		guard = Patch(watcher.AddAvgWatchValue, func(key string, value float64) {
		})
		defer guard.Unpatch()

		guard = Patch(GetBackendConfig, func(reqCtx *mvutil.ReqCtx) string {
			return "test_config"
		})
		defer guard.Unpatch()

		guard = Patch(backend.ExcludeDisplayPackageABTest, func(reqCtx *mvutil.ReqCtx) {
		})
		defer guard.Unpatch()

		// guard = Patch(WriteServerLog, func(reqCtx *mvutil.ReqCtx, res *corsair_proto.QueryResult_) {
		// 	return
		// })
		//defer guard.Unpatch()

		Convey("Given a pipeline", func() {
			filters := []pf.Filter{&BackendReqFilter{}}

			std := pf.StraightPipeline{
				Name:    "Standard",
				Filters: &filters,
			}

			Convey("data illlegal", func() {
				in := "error_data"
				ret, err := std.Process(&in)

				Convey("Then get the excepted result", func() {
					So(err, ShouldBeError)
					So(ret, ShouldBeNil)
				})
			})
		})
	})
}

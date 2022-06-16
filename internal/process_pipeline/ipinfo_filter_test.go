package process_pipeline

import (
	"testing"

	pf "github.com/easierway/pipefiter_framework/pipefilter"
	. "github.com/smartystreets/goconvey/convey"
)

func TestIpinfoFilterProcess(t *testing.T) {
	Convey("test Process", t, func() {
		Convey("Given a pipeline", func() {
			filters := []pf.Filter{&IpInfoFilter{}}

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

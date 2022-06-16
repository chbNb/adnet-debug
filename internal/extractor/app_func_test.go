package extractor

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	smodel "gitlab.mobvista.com/ADN/structs/model"
)

func TestLegalApp(t *testing.T) {
	Convey("test legalApp", t, func() {
		var res bool

		Convey("return false", func() {
			appInfo := smodel.AppInfo{
				App: smodel.App{
					Status: 10,
				},
			}
			res = legalApp(&appInfo)
			So(res, ShouldBeFalse)
		})

		Convey("return true", func() {
			appInfo := smodel.AppInfo{
				App: smodel.App{
					Status: 1,
				},
			}
			res = legalApp(&appInfo)
			So(res, ShouldBeTrue)
		})
	})
}

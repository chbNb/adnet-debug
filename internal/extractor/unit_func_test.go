package extractor

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	smodel "gitlab.mobvista.com/ADN/structs/model"
)

func TestLegalUnit(t *testing.T) {
	Convey("test legalUnit", t, func() {
		Convey("not active", func() {
			unit := smodel.UnitInfo{
				Unit: smodel.Unit{
					Status: 4,
				},
			}
			res := legalUnit(&unit)
			So(res, ShouldBeFalse)
		})

		Convey("active", func() {
			unit := smodel.UnitInfo{
				Unit: smodel.Unit{
					Status: 1,
				},
			}
			res := legalUnit(&unit)
			So(res, ShouldBeTrue)
		})
	})
}

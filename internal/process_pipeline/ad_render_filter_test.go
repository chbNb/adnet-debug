package process_pipeline

import (
	"testing"

	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"

	"gitlab.mobvista.com/ADN/adnet/internal/extractor"

	. "github.com/bouk/monkey"
	. "github.com/smartystreets/goconvey/convey"
)

func TestTransLiebaoReqLog(t *testing.T) {
	Convey("test transLiebaoReqLog", t, func() {
		var params mvutil.Params

		Convey("len transMap <= 0", func() {
			guard := Patch(extractor.GetMP_MAP_UNIT, func() (map[string]mvutil.MP_MAP_UNIT_, bool) {
				mpUnit := map[string]mvutil.MP_MAP_UNIT_{}
				return mpUnit, true
			})
			defer guard.Unpatch()

			transLiebaoReqLog(&params)
		})

		Convey("unitId <= 0", func() {
			params = mvutil.Params{
				UnitID: 0,
			}
			gurad := Patch(extractor.GetMP_MAP_UNIT, func() (map[string]mvutil.MP_MAP_UNIT_, bool) {
				mpUnit := map[string]mvutil.MP_MAP_UNIT_{
					"123": mvutil.MP_MAP_UNIT_{
						UnitID: int64(123),
					},
				}
				return mpUnit, true
			})
			defer gurad.Unpatch()

			transLiebaoReqLog(&params)
		})

		Convey("unitId > 0", func() {
			params = mvutil.Params{
				UnitID: 123,
			}
			gurad := Patch(extractor.GetMP_MAP_UNIT, func() (map[string]mvutil.MP_MAP_UNIT_, bool) {
				mpUnit := map[string]mvutil.MP_MAP_UNIT_{
					"mp_123": mvutil.MP_MAP_UNIT_{
						UnitID: int64(223),
					},
				}
				return mpUnit, true
			})
			defer gurad.Unpatch()

			transLiebaoReqLog(&params)
			So(params.UnitID, ShouldEqual, 223)
		})
	})
}

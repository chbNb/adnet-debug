package mvutil

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"gitlab.mobvista.com/ADN/adnet/internal/mvconst"
)

func TestIsThirdDsp(t *testing.T) {
	Convey("IsThirdDsp", t, func() {
		So(IsThirdDsp(mvconst.FakeAdserverDsp), ShouldBeFalse)
		So(IsThirdDsp(mvconst.MVDSP), ShouldBeFalse)
		So(IsThirdDsp(mvconst.MAS), ShouldBeFalse)
		So(IsThirdDsp(mvconst.FakeGDT), ShouldBeTrue)
		So(IsThirdDsp(mvconst.FakeToutiao), ShouldBeTrue)
		So(IsThirdDsp(mvconst.PokktDsp), ShouldBeTrue)
	})
}

// func TestGetBackendId(t *testing.T) {
// 	Convey("GetBackendId" ,t, func() {
// 		So(GetBackendId(mvconst.FakeToutiao), ShouldEqual, mvconst.TouTiao)
// 	})
// }

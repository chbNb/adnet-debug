package errorcode

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestString(t *testing.T) {
	Convey("test String", t, func() {
		Convey("返回 err", func() {
			zero := EXCEPTION_SPECIAL_ZERO
			So(zero.String(), ShouldEqual, "0")
		})
		Convey("返回 unset", func() {
			unset := SpecialCode(3)
			So(unset.String(), ShouldEqual, "unset")
		})
		Convey("返回 empty", func() {
			unset := EXCEPTION_SPECIAL_EMPTY_STRING
			So(unset.String(), ShouldEqual, "")
		})
	})
}

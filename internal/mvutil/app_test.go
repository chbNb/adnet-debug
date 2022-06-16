package mvutil

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestIsDevinfoEncrypt(t *testing.T) {
	Convey("Given APP.DevinfoEncrypt not 1 means deviceinfo is encrypted", t, func() {
		res := IsDevinfoEncrypt(&AppInfo{})
		So(res, ShouldBeTrue)
	})

	Convey("Given APP.DevinfoEncrypt =1 means deviceinfo is not encrypted", t, func() {
		res := IsDevinfoEncrypt(&AppInfo{
			App: App{DevinfoEncrypt: 1},
		})
		So(res, ShouldBeFalse)
	})

}

func TestAppFcaDefault(t *testing.T) {
	Convey("AppFcaDefault", t, func() {
		So(AppFcaDefault(0), ShouldBeTrue)
		So(AppFcaDefault(2), ShouldBeTrue)
		So(AppFcaDefault(1), ShouldBeFalse)
	})

}

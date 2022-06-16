package helpers

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRenderSdkVersion(t *testing.T) {
	Convey("mal ok", t, func() {
		// 541, 542 ,550, 551
		// 5040100, 5040200, 5050000, 5050100
		version := RenderSdkVersion("MAL_9.11.0")
		So(version.SDKType, ShouldEqual, "MAL")
		So(version.SDKNumber, ShouldEqual, "9.11.0")
		So(version.SDKVersionCode, ShouldEqual, 9110000)
		version = RenderSdkVersion("MI_5.4.1")
		So(version.SDKType, ShouldEqual, "MI")
		So(version.SDKNumber, ShouldEqual, "5.4.1")
		So(version.SDKVersionCode, ShouldEqual, 5040100)
		version = RenderSdkVersion("MI_5.4.2")
		So(version.SDKType, ShouldEqual, "MI")
		So(version.SDKNumber, ShouldEqual, "5.4.2")
		So(version.SDKVersionCode, ShouldEqual, 5040200)
		version = RenderSdkVersion("MI_5.5.0")
		So(version.SDKType, ShouldEqual, "MI")
		So(version.SDKNumber, ShouldEqual, "5.5.0")
		So(version.SDKVersionCode, ShouldEqual, 5050000)
		version = RenderSdkVersion("MI_5.5.1")
		So(version.SDKType, ShouldEqual, "MI")
		So(version.SDKNumber, ShouldEqual, "5.5.1")
		So(version.SDKVersionCode, ShouldEqual, 5050100)
	})
}

package mvutil

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRenderSDKVersion(t *testing.T) {
	Convey("RenderSDKVersion 空字符串返回 SDKVersionItem 对象内容空", t, func() {
		res := RenderSDKVersion("")
		So(res, ShouldResemble, SDKVersionItem{
			SDKNumber:      "",
			SDKVersionCode: 0,
		})
	})

	Convey("RenderSDKVersion 不包含 _", t, func() {
		res := RenderSDKVersion("1.0.1")
		So(res, ShouldResemble, SDKVersionItem{
			SDKType:        "",
			SDKNumber:      "1.0.1",
			SDKVersionCode: 10001,
		})
	})

	Convey("RenderSDKVersion 包含 _，前缀不合法", t, func() {
		res := RenderSDKVersion("test_1.2.1")
		So(res, ShouldResemble, SDKVersionItem{
			SDKType:        "test",
			SDKNumber:      "test_1.2.1",
			SDKVersionCode: 201,
		})
	})

	Convey("RenderSDKVersion 包含 _", t, func() {
		res := RenderSDKVersion("mi_1.2.3")
		So(res, ShouldResemble, SDKVersionItem{
			SDKType:        "mi",
			SDKNumber:      "1.2.3",
			SDKVersionCode: 10203,
		})
	})
}

package process_pipeline

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func Test_replaceBrand(t *testing.T) {
	Convey("replaceBrand", t, func() {
		model := "huaweihuawei huawei_huawei-huawei huaweihuawei huawei-lh008 huawei"
		ret, err := replaceBrand("huawei", "huawei", model)
		So(err, ShouldBeNil)
		So(ret, ShouldEqual, "lh008 huawei")

	})
}

func Test_replaceModel(t *testing.T) {
	Convey("replaceBrand", t, func() {
		brand := "0"
		model := "huawei huawei huawei huawei lh008 huawei"
		newBrand, newModel, ok := trimModel(brand, model, map[string]string{"huawei": "1"})
		So(ok, ShouldBeTrue)
		So(newBrand, ShouldEqual, "huawei")
		So(newModel, ShouldEqual, "lh008 huawei")
	})
}

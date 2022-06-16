package mvconst

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetImageSizeID(t *testing.T) {
	Convey("GetImageSizeID 图片 320x50 返回 2", t, func() {
		k := GetImageSizeID("320x50")
		So(k, ShouldEqual, 2)
	})

	Convey("GetImageSizeID 图片 300x250 返回 3", t, func() {
		k := GetImageSizeID("300x250")
		So(k, ShouldEqual, 3)
	})

	Convey("GetImageSizeID 图片 480x320 返回 4", t, func() {
		k := GetImageSizeID("480x320")
		So(k, ShouldEqual, 4)
	})

	Convey("GetImageSizeID 图片 320x480 返回 5", t, func() {
		k := GetImageSizeID("320x480")
		So(k, ShouldEqual, 5)
	})

	Convey("GetImageSizeID 图片 300x300 返回 6", t, func() {
		k := GetImageSizeID("300x300")
		So(k, ShouldEqual, 6)
	})

	Convey("GetImageSizeID 图片 1200x627 返回 7", t, func() {
		k := GetImageSizeID("1200x627")
		So(k, ShouldEqual, 7)
	})

	Convey("GetImageSizeID VIDEO 返回 8", t, func() {
		k := GetImageSizeID("VIDEO")
		So(k, ShouldEqual, 8)
	})

	Convey("GetImageSizeID 其他字符返回 0", t, func() {
		k := GetImageSizeID("other")
		So(k, ShouldEqual, 0)
	})
}

func TestGetImageSizeByID(t *testing.T) {
	Convey("GetImageSizeByID ID 不存在 ImageSizeMap 中时返回 128x128", t, func() {
		imageSize := GetImageSizeByID(0)
		So(imageSize, ShouldEqual, "128x128")
	})
	Convey("GetImageSizeByID ID = 2 时返回 320x50", t, func() {
		imageSize := GetImageSizeByID(2)
		So(imageSize, ShouldEqual, "320x50")
	})
	Convey("GetImageSizeByID ID = 3 时返回 300x250", t, func() {
		imageSize := GetImageSizeByID(3)
		So(imageSize, ShouldEqual, "300x250")
	})
	Convey("GetImageSizeByID ID = 4 时返回 480x320", t, func() {
		imageSize := GetImageSizeByID(4)
		So(imageSize, ShouldEqual, "480x320")
	})
	Convey("GetImageSizeByID ID = 5 时返回 320x480", t, func() {
		imageSize := GetImageSizeByID(5)
		So(imageSize, ShouldEqual, "320x480")
	})
	Convey("GetImageSizeByID ID = 6 时返回 300x300", t, func() {
		imageSize := GetImageSizeByID(6)
		So(imageSize, ShouldEqual, "300x300")
	})
	Convey("GetImageSizeByID ID = 7 时返回 1200x627", t, func() {
		imageSize := GetImageSizeByID(7)
		So(imageSize, ShouldEqual, "1200x627")
	})
	Convey("GetImageSizeByID ID = 8 时返回 VIDEO", t, func() {
		imageSize := GetImageSizeByID(8)
		So(imageSize, ShouldEqual, "VIDEO")
	})
}

func TestGetPlatformStr(t *testing.T) {
	Convey("GetPlatformStr 1 返回 android", t, func() {
		res := GetPlatformStr(1)
		So(res, ShouldEqual, "android")
	})

	Convey("GetPlatformStr 2 返回 ios", t, func() {
		res := GetPlatformStr(2)
		So(res, ShouldEqual, "ios")
	})

	Convey("GetPlatformStr 其他数字返回 other", t, func() {
		res := GetPlatformStr(3)
		So(res, ShouldEqual, "other")
	})
}

func TestGetNetworkName(t *testing.T) {
	Convey("GetNetworkName 0 返回 unknown", t, func() {
		res := GetNetworkName(0)
		So(res, ShouldEqual, "unknown")
	})

	Convey("GetNetworkName 2 返回 2g", t, func() {
		res := GetNetworkName(2)
		So(res, ShouldEqual, "2g")
	})

	Convey("GetNetworkName 3 返回 3g", t, func() {
		res := GetNetworkName(3)
		So(res, ShouldEqual, "3g")
	})

	Convey("GetNetworkName 4 返回 4g", t, func() {
		res := GetNetworkName(4)
		So(res, ShouldEqual, "4g")
	})

	Convey("GetNetworkName 9 返回 wifi", t, func() {
		res := GetNetworkName(9)
		So(res, ShouldEqual, "wifi")
	})

	Convey("GetNetworkName 其他数字返回 unknown", t, func() {
		res := GetNetworkName(100)
		So(res, ShouldEqual, "unknown")
	})
}

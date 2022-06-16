package extractor

import (
	"testing"

	smodel "gitlab.mobvista.com/ADN/structs/model"

	. "github.com/bouk/monkey"
	. "github.com/smartystreets/goconvey/convey"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
)

func TestGetConfigCenterValue(t *testing.T) {
	Convey("test getConfigCenterValue", t, func() {
		var areaConfig mvutil.AreaConfig
		areaConfig.HttpConfig.ConfigCenterKey = "test_key"
		mvutil.Config.AreaConfig = &areaConfig

		Convey("return error", func() {
			guard := Patch(GetConfigcenter, func(key string) (configcenter *smodel.ConfigCenter, ifFind bool) {
				return &smodel.ConfigCenter{}, false
			})
			defer guard.Unpatch()

			res, err := getConfigCenterValue("test")
			So(res, ShouldBeNil)
			So(err, ShouldBeError)
		})

		Convey("get data, return area not found", func() {
			guard := Patch(GetConfigcenter, func(key string) (configcenter *smodel.ConfigCenter, ifFind bool) {
				return &smodel.ConfigCenter{
					Key:   key,
					Value: map[string]interface{}{key: key + "_val"},
					Area:  "test_area",
				}, true
			})
			defer guard.Unpatch()

			res, err := getConfigCenterValue("hello")
			So(res, ShouldBeNil)
			So(err, ShouldBeError)
		})

		Convey("get data, area is found", func() {
			guard := Patch(GetConfigcenter, func(key string) (configcenter *smodel.ConfigCenter, ifFind bool) {
				return &smodel.ConfigCenter{
					Key:   key,
					Value: map[string]interface{}{key: key + "_val"},
					Area:  "test_key",
				}, true
			})
			defer guard.Unpatch()

			res, err := getConfigCenterValue("test_key")
			So(res, ShouldEqual, "test_key_val")
			So(err, ShouldBeNil)
		})

		Convey("get data, area is found, key not found", func() {
			guard := Patch(GetConfigcenter, func(key string) (configcenter *smodel.ConfigCenter, ifFind bool) {
				return &smodel.ConfigCenter{
					Key:   "fixed_key",
					Value: map[string]interface{}{"fixed_key": key + "_val"},
					Area:  "fixed_key_2",
				}, true
			})
			defer guard.Unpatch()

			res, err := getConfigCenterValue("test_key")
			So(res, ShouldBeNil)
			So(err, ShouldBeError)
		})
	})
}

func TestGetDOMAIN_TRACK(t *testing.T) {
	Convey("test GetDOMAIN_TRACK", t, func() {
		guard := Patch(getConfigCenterValue, func(key string) (interface{}, error) {
			return "domain_track", nil
		})
		defer guard.Unpatch()

		res := GetDOMAIN_TRACK()
		So(res, ShouldResemble, "domain_track")
	})
}

func TestPrivateGetDOMAIN_TRACK(t *testing.T) {
	Convey("test getDOMAIN_TRACK", t, func() {
		Convey("get error", func() {
			res := getDOMAIN_TRACK(nil)
			So(res, ShouldResemble, "")
		})
	})
}

func TestGetDOMAIN(t *testing.T) {
	Convey("test GetDOMAIN", t, func() {
		guard := Patch(getConfigCenterValue, func(key string) (interface{}, error) {
			return "domain", nil
		})
		defer guard.Unpatch()
		res := GetDOMAIN()
		So(res, ShouldResemble, "domain")
	})
}

func TestPrivateGetDOMAIN(t *testing.T) {
	Convey("test getDOgetDOMAINMAIN_TRACK", t, func() {
		Convey("get error", func() {
			res := getDOMAIN(nil)
			So(res, ShouldResemble, "")
		})

		Convey("no error", func() {
			res := getDOMAIN(mvutil.TRACKING_DB{Write: []string{"ww2"}})
			So(res, ShouldResemble, "")
		})
	})
}

func TestGetSYSTEM(t *testing.T) {
	Convey("test GetSYSTEM", t, func() {
		guard := Patch(getConfigCenterValue, func(key string) (interface{}, error) {
			return "sys", nil
		})
		defer guard.Unpatch()
		res := GetSYSTEM()
		So(res, ShouldResemble, "sys")
	})
}

func TestPrivateGetSYSTEM(t *testing.T) {
	Convey("test getSYSTEM", t, func() {
		Convey("get error", func() {
			res := getSYSTEM(nil)
			So(res, ShouldResemble, "")
		})

		Convey("no error", func() {
			res := getSYSTEM("test_system")
			So(res, ShouldResemble, "test_system")
		})
	})
}

func TestGetSYSTEM_AREA(t *testing.T) {
	Convey("test GetSYSTEM_AREA", t, func() {
		guard := Patch(getConfigCenterValue, func(key string) (interface{}, error) {
			return "sys_area", nil
		})
		defer guard.Unpatch()
		res := GetSYSTEM_AREA()
		So(res, ShouldResemble, "sys_area")
	})
}

func TestPrivategetSYSTEM_AREA(t *testing.T) {
	Convey("test getSYSTEM_AREA", t, func() {
		Convey("get error", func() {
			res := getSYSTEM_AREA(nil)
			So(res, ShouldResemble, "")
		})

		Convey("no error", func() {
			res := getSYSTEM_AREA("getSYSTEM_AREA")
			So(res, ShouldResemble, "getSYSTEM_AREA")
		})
	})
}

package mvutil

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	testConf           = "testdata/adnet_server_test.conf"
	testConfSdkVersion = "testdata/adnet_server_test_sdkversion.conf"
)

func TestLoadConfigFile(t *testing.T) {
	config := Config
	//load success
	Convey("Given load a adnet_server config file", t, func() {
		res := config.LoadConfigFile(testConf)
		Convey("Then load the config file success", func() {
			So(res, ShouldNotBeNil)
		})
	})

	//load failed
	Convey("Given load a wrong config file", t, func() {
		res := config.LoadConfigFile("testdata/xxx.conf")
		Convey("Then load the config file failed", func() {
			So(res, ShouldNotBeNil)
			So(res.Error(), ShouldEqual, "load config file:testdata/xxx.conf config file not found")
		})
	})
}

func TestParseSDKVersion(t *testing.T) {
	config := Config

	Convey("Given load the config file success", t, func() {
		config.LoadConfigFile(testConf)

		Convey("Then parse SDK version success", func() {
			res := ParseSDKVersion([]string{"7.6.6", "a.b.c"})
			So(res, ShouldBeFalse)
			//So(SDKVersions, ShouldResemble, []string{"MAL_7.6.6", "1.5.0"})

		})
	})

	Convey("Given load the config file with wrong SDkVersions", t, func() {
		config.LoadConfigFile(testConfSdkVersion)

		Convey("Then parse SDK version failed", func() {
			res := ParseSDKVersion([]string{"7.6.6", "a.b.c"})
			So(res, ShouldBeFalse)
			//So(config.HttpConfig.SDkVersions, ShouldResemble, []string{"7.6.6", "a.b.c"})
		})
	})
}

func TestInitUaParser(t *testing.T) {
	// TODO:
	//func InitUaParser()写死了配置文件的地址
}

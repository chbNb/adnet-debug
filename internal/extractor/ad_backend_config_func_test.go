package extractor

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAdBackendConfigUpdateFunc(t *testing.T) {
	Convey("test adBackendConfigUpdateFunc", t, func() {
		// var query mgo.Query
		// var num int
		// var dbLoader dbLoaderInfo

		// guard := Patch(adBackendConfigIncUpdateFunc, func(adBackendConfigInfo *mvutil.AdBackendConfigInfo, dbLoaderInfo *dbLoaderInfo) {
		// 	return
		// })
		// guard.Unpatch()

		// Convey("query.All(&elems)", func() {
		// 	num = 10
		// 	res, err := adBackendConfigUpdateFunc(&query, num, &dbLoaderInfo)
		// 	So(res, ShouldEqual, 0)
		// 	So(err, ShouldBeError)
		// })
	})
}

func TestLegalAdBackendConfig(t *testing.T) {
	// var res bool
	// var adBackendConfigInfo mvutil.AdBackendConfigInfo
	// Convey("test legalAdBackendConfig", t, func() {
	//	Convey("空参数返回 false", func() {
	//		res = legalAdBackendConfig(&adBackendConfigInfo)
	//		So(res, ShouldBeFalse)
	//	})
	//
	//	Convey("status 不为 1 返回 false", func() {
	//		adBackendConfigInfo.Status = 2
	//		res = legalAdBackendConfig(&adBackendConfigInfo)
	//		So(res, ShouldBeFalse)
	//	})
	//
	//	Convey("status 为 1 返回 true", func() {
	//		adBackendConfigInfo.Status = 1
	//		res = legalAdBackendConfig(&adBackendConfigInfo)
	//		So(res, ShouldBeTrue)
	//	})
	// })
}

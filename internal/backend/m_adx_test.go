package backend

import (
	"bou.ke/monkey"
	supply_mvutil "gitlab.mobvista.com/ADN/chasm/module/supply/mvutil"
	"testing"

	. "github.com/bouk/monkey"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/valyala/fasthttp"
	"gitlab.mobvista.com/ADN/adnet/internal/corsair_proto"
	"gitlab.mobvista.com/ADN/adnet/internal/extractor"
	"gitlab.mobvista.com/ADN/adnet/internal/mvconst"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	"gitlab.mobvista.com/ADN/adnet/internal/watcher"
	rtb "gitlab.mobvista.com/ADN/mtg_openrtb/pkg/mtgrtb"
	smodel "gitlab.mobvista.com/ADN/structs/model"
)

func setupTest() {
	extractor.NewMDbLoaderRegistry()
}

func TestMAdxFilterBackend(t *testing.T) {
	setupTest()
	Convey("test filterBackend", t, func() {
		backend := MAdxBackend{}

		var reqCtx mvutil.ReqCtx
		Convey("return true", func() {
			vFlag := int32(11)
			reqCtx = mvutil.ReqCtx{
				ReqParams: &mvutil.RequestParams{
					Param: mvutil.Params{
						VersionFlag: vFlag,
					},
				},
			}
			res := backend.filterBackend(&reqCtx)
			So(res, ShouldEqual, mvconst.BackendOK)
		})

		Convey("return true-2", func() {
			idfa := "test_idfa"
			devID := "test_devid"
			// adType := int32(42)
			vFlag := int32(1)
			reqCtx = mvutil.ReqCtx{
				ReqParams: &mvutil.RequestParams{
					Param: mvutil.Params{
						Platform:    mvconst.PlatformIOS,
						IDFA:        idfa,
						AndroidID:   devID,
						AdType:      mvconst.ADTypeRewardVideo,
						VersionFlag: vFlag,
						FormatSDKVersion: supply_mvutil.SDKVersionItem{
							SDKVersionCode: 20500,
						},
					},
				},
			}

			res := backend.filterBackend(&reqCtx)
			So(res, ShouldEqual, mvconst.BackendOK)
		})
	})
}

func TestMAdxComposeHttpRequest(t *testing.T) {
	setupTest()
	guard := Patch(watcher.AddWatchValue, func(key string, val float64) {
	})
	defer guard.Unpatch()
	Convey("test composeHttpRequest", t, func() {
		var b MAdxBackend
		b.ID = 8
		Convey("reqCtx is nil", func() {
			err := b.composeHttpRequest(nil, nil, nil)
			So(err, ShouldBeError)
		})

		Convey("reqCtx is not nil, ios", func() {
			switchGuard := Patch(extractor.GetFillRateEcpmFloorSwitch, func() bool {
				return true
			})
			defer switchGuard.Unpatch()

			guard := monkey.Patch(extractor.GetDEBUG_BID_FLOOR_AND_BID_PRICE_CONF, func() map[string]*mvutil.DebugBidFloorAndBidPriceConf {
				return map[string]*mvutil.DebugBidFloorAndBidPriceConf{}
			})
			defer guard.Unpatch()

			data := `{"ads":[{"ad_id":"1599802266710029","creative":{"app_name":"剑与家园_iOS首发","button_text":"点击下载","click_url":["https://lf.snssdk.com/api/ad/union/event/?req_id=916176be5b18b910c6c1e23d81f8c8eeu9410\u0026extra=Izew6KCpNmbGQNJfcMQPuhqDP8BkX8OPbgLJjaBzIW6xeW5dxUbOEjq2LEmAu4uSe%2F5wS752xqGM6jxgaaSRoACVmNgdl0jxDoehyWG%2FWCLKxDKuLIriS2wvBxmaDsQtiqx0C%2Fp0W6faHSKcW0ssU22%2BT0uIAFmbufoWXuWbfYJK5LFRmyTgNefRkDrIEmSC6nSyrLkGBE277HkhMdUE8%2FO13OxjPoAPifMDqDC3GCyDGv4EI2R3b2LNDhoN2FNB8CMr2PD%2BLwjbEdoDV6LfQL4lSOVz9R%2BwX4j6pybSDIwsD8QSV%2FfZ57pABeLsh0uiBMVbqekmqi0jfXqJ%2BEQJzLW6hi%2F0XmJ7DGwUUrlFkQvKVtkmJPTWpwvijyXkJbE2uwNL31aw7zfWdlhzxgFmAQwxlbSCL7Gh%2FhC13Soi8d3%2B4CCNQJLLKEMUmTv0RYbWxogXg9mBtgsKdIU2FyNhJBPcoZ0CvEA1BtNKAXiZ5BgGBI5Vws9rutiLcMcz0nXjsn4kbXHrT3KfHRK2Z28H%2F0FKnf5EUkYqjExQhm5XWNsNFWyU3QOdzZEmrrRnccfew36i14FUNh5kep8%2FMQDfu0VMUKy18up%2FYmGebVCSBEEXsvFyikHwcEz3pbmIFuKWMykGc52t8u%2FresQhrT7cw2m3BxEG4G10qinnOSbkPSEbl13AwIPLPOr09SQxRi4moMhvSPWBPxq4Nlp7RH4GlAvOuOxiOUyT4x7MH4M4VXnA7cKD3fp%2BL1f7cllERfCTLxLIJltsDgeXeqRu2yhZjr3IegOMcAXERY%2BH8CEdPEQUBmLyFEDbIbGVFoaXQrAtFPHJBWyqYApGQWPducXcZvhszyHewiuXbJa7uxoI0ie9L4JddAR9k%2FZXmQzib7NW6I6I%2FCq7yMsMoPaJakAsLxl%2Fm%2FGT1KqfkX7gfLL1UoRnY3EMZYGe3CNhgb%2BuS5Ag%2BeUKreH01pKOTn8lzoVhpWiQ9e4RI%2FxqPfJ4tVF0qTJSmJ8OTn885xbuVml2LiQEGFLJJIfGEIFqRudC5rHQOKf5BZihJ7ID6DIs7rSpPbZ5VflB2mFOqi811MFCzwBAfP%2FzChh6xlqkPhoaz5%2BJOIL684isk4Q7Z9etYQS3Uy03zqnhu4cJtzJtPnA3SvOdJchZq2qJQhFM4%2Flmzku3FA9UAbsLxUoE60GmNQ9YPPwyep3Ii%2FkMkgmtn2THtD%2FbAIHzuItNX6o2SGKXYN9QtmLlSzXkDepK9iRCO7%2FGOSVzpk9URdRee0dTmOWaX4kTpeACg2fnOSien7qwtJ9mZWafV%2FjcG%2FID2AgSFwMkc6M%3D\u0026source_type=1\u0026pack_time=1528346897.0"],"creative_type":6,"description":"太嗨了，玩这游戏每天充电五次，无限体力一玩就上瘾","download_url":"https://lf.snssdk.com/api/ad/union/redirect/?req_id=916176be5b18b910c6c1e23d81f8c8eeu9410\u0026use_pb=1\u0026rit=901287679\u0026call_back=03M4TqoDRiBzTyomsTlYNlVoqNUL542DrOoJrKy3Fenl5JlR%2F4nw0UULWppSpuuRdYF45S7qzvBxqdIu5v4HSI0yHBtA7gzt9ZV8PGlGlh1mn1f43BvyA9gIEhcDJHOj\u0026extra=Sw9qWTzJg8YAJYRq%2Bp9T3z87wcXIfskDeygZG8IzWcWiq%2FRDk3cTND7E05UojdMer6twm8vEXy9VxjRPc%2B18FkkaiquqG%2Fljev0M4DYiasZbokwry2zu5BbbHzEkhdTJkin3sVQlSyOAfujgXF2xz3Nz3bD27kiY8tC05ysJeaHxBjCPE%2BTbOeMCEsslvCz1q%2FabPa0hzeUcHeYVrubFqdYWCtDb5NPczi%2FdRZILtic4RSemysyekvJdalWi%2FY%2BNtustFon3rRGIjBQxwBZf5iXIWatqiUIRTOP5Zs5LtxQPVAG7C8VKBOtBpjUPWDz8zqTF0%2BLdmaRx8v%2FsnNUA%2B1G9mBgeS3RF%2FuZDC9l4bL57%2FnBLvnbGoYzqPGBppJGgAJWY2B2XSPEOh6HJYb9YIsrEMq4siuJLbC8HGZoOxC2KrHQL%2BnRbp9odIpxbSyxT36d%2FjPAjjnynPD9qQK1zz%2Bv9PDdQKmFlWge2v2IRce7jUynYK6epFyRiF4%2FcZAExwAELcL5kByViwjTPXfdvKwuEMmWZfD4t%2BH0irvALGvM%3D\u0026source_type=1\u0026pack_time=1528346897.0\u0026active_extra=KX4P%2Fs1PqRaSxDTC%2Fp2Kxq1AxYOQnXaqfnUdL9rPio4TsD2LzvIDK0Tq2Xw1Guj85TxGUI3JCcP%2B2Me9277jZA%3D%3D","ext":"Izew6KCpNmbGQNJfcMQPuhqDP8BkX8OPbgLJjaBzIW6xeW5dxUbOEjq2LEmAu4uSe%2F5wS752xqGM6jxgaaSRoACVmNgdl0jxDoehyWG%2FWCLKxDKuLIriS2wvBxmaDsQtiqx0C%2Fp0W6faHSKcW0ssU22%2BT0uIAFmbufoWXuWbfYJK5LFRmyTgNefRkDrIEmSC6nSyrLkGBE277HkhMdUE8%2FO13OxjPoAPifMDqDC3GCyDGv4EI2R3b2LNDhoN2FNB8CMr2PD%2BLwjbEdoDV6LfQFDHvevgv4TgEgfq8An8pl0sD8QSV%2FfZ57pABeLsh0uiBMVbqekmqi0jfXqJ%2BEQJzLW6hi%2F0XmJ7DGwUUrlFkQvKVtkmJPTWpwvijyXkJbE2uwNL31aw7zfWdlhzxgFmAQwxlbSCL7Gh%2FhC13Soi8d3%2B4CCNQJLLKEMUmTv0RYbWxogXg9mBtgsKdIU2FyNhJBPcoZ0CvEA1BtNKAXiZ5BgGBI5Vws9rutiLcMcz0nXjsn4kbXHrT3KfHRK2Z28H%2F0FKnf5EUkYqjExQhm5XWNsNFWyU3QOdzZEmrrRnccfew36i14FUNh5kep8%2FMQDfu0VMUKy18up%2FYmGebVCSBEEXsvFyikHwcEz3pbmIFuKWMykGc52t8u%2FresQhrT7cw2m3BxEG4G10qinnOSbkPSEbl13AwIPLPOr09SQxRi4moMhvSPWBPxq4Nlp7RH4GlAvOuOxiOUyT4x7MH4M4VXnA7cKD3fp%2BL1f7cllERfCTLxLIJltsDgeXeqRu2yhZjr3IegOMcAXERY%2BH8CEdPEQUBmLyFEDbIbGVFoaXQrAtFPHJBWyqYApGQWPducXcZvhszyHewiuXbJa7uxoI0ie9L4JddAR9k%2FZXmQzib7NW6I6I%2FCq7yMsMoPaJakAsLxl%2Fm%2FGT1KqfkX7gfLL1UoRnY3EMZYGe3CNhgb%2BuS5Ag%2BeUKreH01pKOTn8lzoVhpWiQ9e4RI%2FxqPfJ4tVF0qTJSmJ8OTn885xbuVml2LiQEGFLJJIfGEIFqRudC5rHQOKf5BZihJ7ID6DIs7rSpPbZ5VflB2mFOqi811MFCzwBAfP%2FzChh6xlqkPhoaz5%2BJOIL684isk4Q7Z9etYQS3Uy03zqnhu4cJtzJtPnA3SvOdJchZq2qJQhFM4%2Flmzku3FA9UAbsLxUoE60GmNQ9YPPwyep3Ii%2FkMkgmtn2THtD%2FbAIHzuItNX6o2SGKXYN9QtmLlSzXkDepK9iRCO7%2FGOSVzpk9URdRee0dTmOWaX4kTpeACg2fnOSien7qwtJ9mZWafV%2FjcG%2FID2AgSFwMkc6M%3D","icon":"http://sf1-ttcdn-tos.pstatp.com/img/ad.union.api/cc9ac6fbcb37cf12d602bf07efc37141~c1_0x0_q100.jpeg","image":{"height":720,"url":"http://sf3-ttcdn-tos.pstatp.com/obj/web.business.image/201805075d0dc058170cb2654bf0b81f","width":1280},"image_list":[{"height":720,"url":"http://sf3-ttcdn-tos.pstatp.com/obj/web.business.image/201805075d0dc058170cb2654bf0b81f","width":1280}],"image_mode":5,"interaction_type":4,"package_name":"com.lilithgame.sgame","phone_num":"","show_url":["https://lf.snssdk.com/api/ad/union/show_event/?req_id=916176be5b18b910c6c1e23d81f8c8eeu9410\u0026extra=Izew6KCpNmbGQNJfcMQPuhqDP8BkX8OPbgLJjaBzIW6xeW5dxUbOEjq2LEmAu4uSe%2F5wS752xqGM6jxgaaSRoACVmNgdl0jxDoehyWG%2FWCLKxDKuLIriS2wvBxmaDsQtiqx0C%2Fp0W6faHSKcW0ssU22%2BT0uIAFmbufoWXuWbfYJK5LFRmyTgNefRkDrIEmSC6nSyrLkGBE277HkhMdUE8%2FO13OxjPoAPifMDqDC3GCyDGv4EI2R3b2LNDhoN2FNB8CMr2PD%2BLwjbEdoDV6LfQL4lSOVz9R%2BwX4j6pybSDIwsD8QSV%2FfZ57pABeLsh0uiBMVbqekmqi0jfXqJ%2BEQJzLW6hi%2F0XmJ7DGwUUrlFkQvKVtkmJPTWpwvijyXkJbE2uwNL31aw7zfWdlhzxgFmAQwxlbSCL7Gh%2FhC13Soi8d3%2B4CCNQJLLKEMUmTv0RYbWxogXg9mBtgsKdIU2FyNhJBPcoZ0CvEA1BtNKAXiZ5BgGBI5Vws9rutiLcMcz0nXjsn4kbXHrT3KfHRK2Z28H%2F0FKnf5EUkYqjExQhm5XWNsNFWyU3QOdzZEmrrRnccfew36i14FUNh5kep8%2FMQDfu0VMUKy18up%2FYmGebVCSBEEXsvFyikHwcEz3pbmIFuKWMykGc52t8u%2FresQhrT7cw2m3BxEG4G10qinnOSbkPSEbl13AwIPLPOr09SQxRi4moMhvSPWBPxq4Nlp7RH4GlAvOuOxiOUyT4x7MH4M4VXnA7cKD3fp%2BL1f7cllERfCTLxLIJltsDgeXeqRu2yhZjr3IegOMcAXERY%2BH8CEdPEQUBmLyFEDbIbGVFoaXQrAtFPHJBWyqYApGQWPducXcZvhszyHewiuXbJa7uxoI0ie9L4JddAR9k%2FZXmQzib7NW6I6I%2FCq7yMsMoPaJakAsLxl%2Fm%2FGT1KqfkX7gfLL1UoRnY3EMZYGe3CNhgb%2BuS5Ag%2BeUKreH01pKOTn8lzoVhpWiQ9e4RI%2FxqPfJ4tVF0qTJSmJ8OTn885xbuVml2LiQEGFLJJIfGEIFqRudC5rHQOKf5BZihJ7ID6DIs7rSpPbZ5VflB2mFOqi811MFCzwBAfP%2FzChh6xlqkPhoaz5%2BJOIL684isk4Q7Z9etYQS3Uy03zqnhu4cJtzJtPnA3SvOdJchZq2qJQhFM4%2Flmzku3FA9UAbsLxUoE60GmNQ9YPPwyep3Ii%2FkMkgmtn2THtD%2FbAIHzuItNX6o2SGKXYN9QtmLlSzXkDepK9iRCO7%2FGOSVzpk9URdRee0dTmOWaX4kTpeACg2fnOSien7qwtJ9mZWafV%2FjcG%2FID2AgSFwMkc6M%3D\u0026source_type=1\u0026pack_time=1528346897.0"],"target_url":"https://lf.snssdk.com/api/ad/union/redirect/?req_id=916176be5b18b910c6c1e23d81f8c8eeu9410\u0026use_pb=1\u0026rit=901287679\u0026call_back=03M4TqoDRiBzTyomsTlYNlVoqNUL542DrOoJrKy3Fenl5JlR%2F4nw0UULWppSpuuRdYF45S7qzvBxqdIu5v4HSI0yHBtA7gzt9ZV8PGlGlh1mn1f43BvyA9gIEhcDJHOj\u0026extra=Sw9qWTzJg8YAJYRq%2Bp9T3z87wcXIfskDeygZG8IzWcWiq%2FRDk3cTND7E05UojdMer6twm8vEXy9VxjRPc%2B18FkkaiquqG%2Fljev0M4DYiasZbokwry2zu5BbbHzEkhdTJkin3sVQlSyOAfujgXF2xz3Nz3bD27kiY8tC05ysJeaHxBjCPE%2BTbOeMCEsslvCz1q%2FabPa0hzeUcHeYVrubFqdYWCtDb5NPczi%2FdRZILtic4RSemysyekvJdalWi%2FY%2BNtustFon3rRGIjBQxwBZf5iXIWatqiUIRTOP5Zs5LtxQPVAG7C8VKBOtBpjUPWDz8zqTF0%2BLdmaRx8v%2FsnNUA%2B1G9mBgeS3RF%2FuZDC9l4bL57%2FnBLvnbGoYzqPGBppJGgAJWY2B2XSPEOh6HJYb9YIsrEMq4siuJLbC8HGZoOxC2KrHQL%2BnRbp9odIpxbSyxT36d%2FjPAjjnynPD9qQK1zz%2Bv9PDdQKmFlWge2v2IRce7jUynYK6epFyRiF4%2FcZAExwAELcL5kByViwjTPXfdvKwuEMmWZfD4t%2BH0irvALGvM%3D\u0026source_type=1\u0026pack_time=1528346897.0\u0026active_extra=KX4P%2Fs1PqRaSxDTC%2Fp2Kxq1AxYOQnXaqfnUdL9rPio4TsD2LzvIDK0Tq2Xw1Guj85TxGUI3JCcP%2B2Me9277jZA%3D%3D","temp_extra_info":"{\"img_gen_type\": 0, \"img_md5\": \"\", \"template_id\": 0}","template_id":"","title":"太嗨了，玩这游戏每天充电五次，无限体力一玩就上瘾","video":{"cover_height":720,"cover_url":"https://sf3-ttcdn-tos.pstatp.com/img/web.business.image/201805075d0dc058170cb2654bf0b81f~noop.jpg","cover_width":1280,"endcard":"https://www.toutiaopage.com/union/endcard/1599802266710029/?style_id=3","resolution":"1280x720","size":8645919,"video_duration":29.44,"video_url":"https://v3-ad.ixigua.com/529f7c8324545c09bf78a6982b9e5280/5b18c159/video/m/220c69dc22ef804487591c27189601aa99d1156af4200004ccfc71e3fef/"}},"filter_words":[{"id":"4:2","is_selected":false,"name":"看过了"}]}],"did":9655564071995315,"processing_time_ms":86,"request_id":"916176be5b18b910c6c1e23d81f8c8ee","status_code":20000}`
			backendCtx := mvutil.NewBackendCtx("MAdx_1", "MAdx_1", "package_name", 1.00, 2, []string{"CN"})
			backendCtx.RespData = []byte(data)
			reqCtx := mvutil.NewReqCtx()
			var param mvutil.Params
			param.RequestPath = mvconst.PATHOnlineApi
			param.NetworkType = mvconst.NETWORK_TYPE_4G
			param.Platform = mvconst.PlatformIOS
			param.RequestID = "5b167141c6c1e24c7cf21b9b"
			param.IDFA = "35951D6A-DB66-47A1-8CD0-DB1E9E0D33CD"
			param.ClientIP = "192.168.1.1"
			param.UserAgent = "Mozilla%252f5.0%2B%28iPhone%25253B%2BCPU%2BiPhone%2BOS%2B9_3_5%2Blike%2BMac%2BOS%2BX%29%2BAppleWebKit%252f601.1.46%2B%28KHTML%25252C%2Blike%2BGecko%29%2BMobile%252f13G36"
			param.TNum = 1
			param.OSVersion = "9.3.5"
			param.Model = "iphone6%2C2"
			param.UnitID = 128
			param.AppName = "iOSTest"
			param.PackageName = "com.tianye1.mvsdk"
			param.FormatAdType = mvconst.ADTypeRewardVideo
			param.AdType = mvconst.ADTypeNative
			param.VideoVersion = "1.1"
			param.VideoH = 300
			param.VideoW = 400
			param.FormatOrientation = mvconst.ORIENTATION_LANDSCAPE
			reqParam := &mvutil.RequestParams{Param: param, UnitInfo: &smodel.UnitInfo{}, AppInfo: &smodel.AppInfo{}}
			//reqParam.UnitInfo.Unit.Orientation = mvconst.ORIENTATION_LANDSCAPE
			reqParam.UnitInfo.Unit.VideoAds = mvconst.VideoAdTypeNOLimit
			reqCtx.ReqParams = reqParam
			reqCtx.ReqParams.PublisherInfo = &smodel.PublisherInfo{}

			err := b.composeHttpRequest(reqCtx, backendCtx, &fasthttp.Request{})
			// So(httpReq, ShouldNotBeEmpty)
			// So(res2, ShouldEqual, "{\"request_id\":\"\",\"api_version\":\"api123\",\"uid\":\"\",\"source_type\":\"app\",\"ua\":\"test_ua\",\"ip\":\"\",\"app\":{\"appid\":\"123\"},\"device\":{\"did\":\"test_idfa\",\"imei\":\"\",\"type\":1,\"os\":2,\"os_version\":\"\",\"vendor\":\"\",\"model\":\"\",\"language\":\"\",\"conn_type\":1,\"mac\":\"\",\"screen_width\":0,\"screen_height\":0,\"orientation\":2},\"adslots\":[{\"id\":\"lotid1\",\"adtype\":7,\"pos\":5,\"accepted_size\":[{\"width\":1200,\"height\":627}],\"ad_count\":1}]}")
			So(err, ShouldBeNil)
		})
	})
}

func TestMAdxParseHttpResponse(t *testing.T) {
	Convey("test parseHttpResponse", t, func() {
		backend := &MAdxBackend{}
		backend.ID = 8
		var reqCtx mvutil.ReqCtx
		var backendCtx mvutil.BackendCtx

		guard := Patch(watcher.AddWatchValue, func(key string, val float64) {
		})
		defer guard.Unpatch()

		Convey("nil", func() {
			res, err := backend.parseHttpResponse(nil, nil)
			So(res, ShouldEqual, 12)
			So(err, ShouldBeError)
		})

		Convey("not nil", func() {
			// guard1 := Patch(addBackendLog, func(adTmp *corsair_proto.Campaign, demandExt *DemandExt, param *mvutil.Params) {

			// })
			// defer guard1.Unpatch()
			idfa := "test_idfa"
			devID := "test_devid"
			adType := int32(42)
			vFlag := int32(1)
			reqID := "req_id"
			uaStr := "test_ua"
			// ntType := enum.NetworkType(9)
			reqCtx = mvutil.ReqCtx{
				ReqParams: &mvutil.RequestParams{
					Param: mvutil.Params{
						IDFA:        idfa,
						AndroidID:   devID,
						AdType:      adType,
						VersionFlag: vFlag,
						RequestID:   reqID,
						UserAgent:   uaStr,
						NetworkType: 9,
					},
					UnitInfo: &smodel.UnitInfo{
						Unit: smodel.Unit{
							Orientation: 0,
						},
					},
				},
			}

			// file, err := os.Open("/Users/ledu/Documents/文件/res.txt") // For read access.
			// if err != nil {
			//	fmt.Printf("error: %v", err)
			//	return
			// }
			// defer file.Close()
			// data, err := ioutil.ReadAll(file)

			backendCtx = mvutil.BackendCtx{
				RespData: []byte("{\"id\":\"29e8574f-1add-423a-a083-e47c5623841f\",\"seatbid\":[{\"bid\":[{\"id\":\"29e8574f-1add-423a-a083-e47c5623841f:1\",\"impid\":\"1\",\"price\":1.833,\"adid\":\"255\",\"nurl\":\"http://tracking.url:8766/?bidRequestId=${AUCTION_ID}&winPrice=${AUCTION_PRICE}\",\"adm\":\"<?xml version=\"1.0\" encoding=\"UTF-8\"?><VAST version=\"2.0\"><Ad id=\"255\"><InLine><AdSystem>Exchange Inc.</AdSystem><AdTitle>Test Video Ad</AdTitle><Description>test video ad</Description><Survey/><Error/><Impression><![CDATA[http://tracking.url:8767/?type=IMPRESSION&bidRequestId=${AUCTION_ID}&winPrice=${AUCTION_PRICE}]]></Impression><Creatives><Creative AdID=\"banner1\"><Linear><Duration>00:00:16</Duration><TrackingEvents><Tracking event=\"start\"><![CDATA[http://tracking.url:8767/?type=START&bidRequestId=${AUCTION_ID}]]></Tracking><Tracking event=\"firstQuartile\"><![CDATA[http://tracking.url:8767/?type=FIRST_QUARTILE&bidRequestId=${AUCTION_ID}]]></Tracking><Tracking event=\"midpoint\"><![CDATA[http://tracking.url:8767/?type=MIDPOINT&bidRequestId=${AUCTION_ID}]]></Tracking><Tracking event=\"thirdQuartile\"><![CDATA[http://tracking.url:8767/?type=THIRD_QUARTILE&bidRequestId=${AUCTION_ID}]]></Tracking><Tracking event=\"complete\"><![CDATA[http://tracking.url:8767/?type=COMPLETE&bidRequestId=${AUCTION_ID}]]></Tracking></TrackingEvents><VideoClicks><ClickThrough><![CDATA[http: //www.amadoad.ru]]></ClickThrough><ClickTracking><![CDATA[http: //tracking.url:8767/?type=CLICK&bidRequestId=${AUCTION_ID}]]></ClickTracking></VideoClicks><MediaFiles><MediaFile delivery=\"progressive\" type=\"video/mp4\" bitrate=\"251\" width=\"480\" height=\"320\" scalable=\"true\" maintainAspectRatio=\"true\">https://dashboard.amadoad.com/_banners/df/45/df45bd40f590b8d3ea16a38731cc59178847af4b.mp4</MediaFile></MediaFiles></Linear></Creative><Creative>\n                    <CompanionAds>\n                        <Companion width=\"300\" height=\"550\">\n                            <HTMLResource><![CDATA[<img src='https://node208.fractionalmedia.com/track_img_bin?bid_id=9ff8f27c-fc9d-4d81-9e59-637e45e3a5d8_171019-08__208' width='1' height='1' border='0' /><script type='text/javascript' src='mraid.js'></script><script charset='utf-8' type='text/javascript'>var adInfo = {clickUrl: 'https://node208.fractionalmedia.com/click2_bin?bid_id=9ff8f27c-fc9d-4d81-9e59-637e45e3a5d8_171019-08__208_83557&idfa=499764ae-ddf5-4c9f-a748-61c1fdf62b2c',bidId: '9ff8f27c-fc9d-4d81-9e59-637e45e3a5d8_171019-08__208',campaignId: '83557',creativeId: '27156',playableId: '70114',deviceModel: 'Apple_iPod Touch',idfa: '499764ae-ddf5-4c9f-a748-61c1fdf62b2c',appName: 'Fleet+Battle+-+Apple',node: '208',os: 'iOS',osVersion: 'iOS_10.0',ssp: 'appodeal',adDimensions: '320_480',adSdkVersion: 'appodeal_2.1.4',gameOptions:{'INSTALL_BANNER_TEXT_ID':'1','tower1UpgradeTime':'20','ENDCARD_HEADERTEXT_VICTORY_TOGGLE':'on','ENDCARD_MESSAGE_TIMELIMIT_ID':'1','upgrade_building_message':'on','buttonInstallNowShadowcolor':'#fdff2c','ENDCARD_MESSAGE_TIMELIMIT_TOGGLE':'on','bannerToggle_landscape':'off','didInteractTimeLimitEnabled':'off','purchase_farm_message':'on','inputMessageShowOnPreloader':'off','wall3BaseStrength':'200','inputMessageShowOnGame':'off','INSTALL_BANNER_BOX_FILL_COLOR':'rgba(0,0,0,1)','BOSS_HEALTH':'10','whiteoutEndCardColorAlpha':'1','preloaderScreenCrackLocation':'0','GameID':'FF_Siege','scoreBackgroundAlpha':'0.75','Property2':'off','osID':'iOS','Property3':'2','Property4':'off','wall2Price':'50','Property5':'off','scoreBackgroundHeaderColor':'#FFFFFF','Property6':'off','closeButtonTimer':'5','FileID':'SGF_1.1.0_Venice.js','Property7':'off','Property8':'off','upgradeTowerScore':'10','upgrade_wall_message':'on','ENDCARD_HEADERTEXT_VICTORY_ID':'1','INSTALL_BANNER_TEXT_FILL_COLOR':'rgba(255,255,255,1)','Property1':'off','TrafficID':'bid_playable','showPreloaderOverlay':'on','preloaderStartCountdown':'off','redirectOnGameStart':'0','TROOP_HEALTH':'2','STARTING_GOLD':'100','scoreBackgroundHeaderFontSize':'17','ENEMY_DAMAGE':'0.1','competitiveNameColor':'#e42f43','showScoreBackgroundHeader':'on','TROOP_ATTACK_INTERVAL':'1','endCardType':'0','WALL_ATTACK_DISTANCE':'0.9','endCardLogo':'on','disclaimerToggle':'on','tower2UpgradeTime':'20','banner_clickable_on_show':'off','BOSS_INTERVAL':'4','competitionName':'gamerChick24','hideCloseButtonTime':'5','custom4':'game_action_wall_lvl_3_upgraded','custom5':'game_action_barracks_used','custom2':'milestone_state_game_over_victory','custom3':'game_action_wall_lvl_2_upgraded','killMechScore':'30','closeButtonBackground_landscape':'rgba(0,0,0,1)','custom1':'milestone_state_game_over_defeated','competition':'off','ENDCARD_HEADERTEXT_TIMELIMIT_TOGGLE':'on','INSTALL_BANNER_FONT_SIZE':'40','TROOP_DAMAGE':'1','competitiveAlwaysWin':'off','farm2GoldIncomeBoost':'1.5','disclaimerFontColor':'rgba(128,128,128,1)','tutorial1_pointer':'on','preloaderScreenCrackDuration':'.25','buildTowerScore':'50','ENEMY_MIN_SPAWN_RATE':'2','tower3RateOfFire':'0.5','towerPrice':'40','endCardFullScreen':'off','TestID':'appodeal_TierUS','inputMessageTimeOut':'5','tower2Damage':'0.5','endCardTransitionTime':'0.25','disclaimerMessageID':'1','overrideJSON':'{\"Default_SGF\":{\"milestone\":\"30\"},\"Steroids_SGF\":{\"BASE_GOLD_INCOME\": \"3\",\"BOSS_DAMAGE\": \"1\",\"BOSS_HEALTH\":\"2\",\"BOSS_INTERVAL\":\"0.5\",\"ENEMY_DAMAGE\":\"0.13\",\"ENEMY_HEALTH\": \"0.5\",\"ENEMY_MIN_SPAWN_RATE\": \"0.25\",\"ENEMY_SPAWN_INTERVAL\": \"1\",\"ENEMY_SPAWN_RATE_INCREASE\": \"0.05\",\"farm1GoldIncomeBoost\": \"3.3\",\"farm1UpGradeTime\": \"7\",\"farm2GoldIncomeBoost\": \"4.5\",\"farm3GoldIncomeBoost\": \"6\",\"tower1RateOfFire\": \"0.11\",\"tower2RateOfFire\": \"0.08\",\"tower3Damage\": \"0.75\",\"tower3RateOfFire\": \"0.05\",\"TROOP_ATTACK_INTERVAL\":\"0.13\",\"TROOP_HEALTH\":\"4\"},\"SteroidsHarder_SGF\": {\"ENEMY_GAP\": \"0.5\",\"BASE_GOLD_INCOME\": \"3\",\"BOSS_DAMAGE\": \"1\",\"BOSS_HEALTH\": \"4\",\"BOSS_INTERVAL\": \"0.5\",\"ENEMY_DAMAGE\": \"0.13\",\"ENEMY_HEALTH\": \"1.5\",\"ENEMY_MIN_SPAWN_RATE\": \"0.25\",\"ENEMY_SPAWN_INTERVAL\": \"1\",\"ENEMY_SPAWN_RATE_INCREASE\": \"0.05\",\"farm1GoldIncomeBoost\": \"3.3\",\"farm1UpGradeTime\": \"7\",\"farm2GoldIncomeBoost\": \"4.5\",\"farm3GoldIncomeBoost\": \"6\",\"tower1RateOfFire\": \"0.11\",\"tower2RateOfFire\": \"0.08\",\"tower3Damage\": \"0.75\",\"tower3RateOfFire\": \"0.05\",\"TROOP_ATTACK_INTERVAL\":\"0.13\",\"TROOP_HEALTH\": \"4\"}}','killSoldierScore':'10','tutorial2_pointer':'on','whiteoutEndCard':'off','didInteractTimeLimit':'15','countDownCloseButton':'off','tower1UpgradeCost':'25','toolTipType':'1','farm1UpGradeCost':'10','farm3GoldIncomeBoost':'2','purchase_barrack_message':'on','ENDCARD_HEADERTEXT_DEFEATED_TOGGLE':'on','EndcardTriggersOption':'3','barrackPrice':'30','closeButtonBackground_portrait':'rgba(0,0,0,1)','overrideKey':'Steroids_SGF','farm1GoldIncomeBoost':'1','ENEMY_HEALTH':'2','inputMessageTextID':'0','ENDCARD_MESSAGE_VICTORY_ID':'1','showScoreBackground':'on','tutorial1':'on','tutorial2':'on','farmPrice':'20','MAX_PLAY_TIME':'150','tower2RateOfFire':'0.75','whiteoutEndCardDuration':'1','farm2UpGradeTime':'15','milestone':'30','country':'','whiteoutEndCardColor':'#ffffff','redirectOnGameEnd':'0','tower1Damage':'0.5','tower1RateOfFire':'1','tower2UpgradeCost':'50','gamePlayLogo':'on','ENDCARD_MESSAGE_VICTORY_TOGGLE':'on','disclaimerFontFamily':'arial','wall1BaseStrength':'50','farm1UpGradeTime':'10','ENDCARD_HEADERTEXT_DEFEATED_ID':'1','ENEMY_SPAWN_RATE_INCREASE':'0.5','ENEMY_GAP':'1','preloaderScreenCrackEffect':'off','SizeID':'320x480','EndcardTriggersDelay':'6','EndcardTriggers':'off','BASE_GOLD_INCOME':'1','preloader_type':'3','showCompetitionScore':'on','INSTALL_BANNER_FONT_FAMILY':'Arial','wall3Price':'100','playerNameColor':'#67b252','playerName':'YOU','purchase_tower_message':'on','preloader':'on','ENDCARD_HEADERTEXT_TIMELIMIT_ID':'1','wall2BaseStrength':'100','characterEndCard':'0','bannerToggle_portrait':'off','tower3Damage':'0.5','tutorial_first_action_message':'on','characterToolTip':'0','farm2UpGradeCost':'20','BOSS_DAMAGE':'2','ENDCARD_MESSAGE_DEFEATED_ID':'1','buttonInstallNowScalling':'on','ENEMY_SPAWN_INTERVAL':'8','disclaimerFontSize':'9','tutorial':'on','ENDCARD_MESSAGE_DEFEATED_TOGGLE':'on'},telemetryUrl: 'https://node208.fractionalmedia.com/telim_direct',country:'US'};window.adInfo = adInfo;</script><div id='gamebox'><div id='preloader'>Connecting to Game Server...</div><canvas id='canvas'>Connecting to Game Server...</canvas></div><script type='text/javascript' src='https://cdns3.fractionalmedia.com/fm.creatives/27156/SGF_1.1.0_Venice.js' crossorigin='anonymous'></script>]]></HTMLResource>\n                            <TrackingEvents>\n                                <Tracking event=\"creativeView\">\n                                    <![CDATA[http://myTrackingURL/firstCompanionCreativeView]]>\n                                </Tracking>\n                            </TrackingEvents>\n                            <CompanionClickThrough>http://www.tremormedia.com</CompanionClickThrough>\n                        </Companion>\n                    </CompanionAds>\n               </Creative> \n</Creatives></InLine></Ad></VAST>\",\"adomain\":[\"www.amadoad.ru\"],\"iurl\":\"https://i.url/image.jpg\",\"cid\":\"campaign178-ads263\",\"crid\":\"banner1\",\"cat\":[\"IAB1-4\",\"IAB1-5\"],\"attr\":[6],\"h\":360,\"w\":480}],\"seat\":\"0\"}]}"),
				// RespData: data,
				// Ads: &corsair_proto.BackendAds{
				//	BackendId: int32(13),
				// },
				// AdReqKeyValue: "ad_req_val1=val1&ad_req_val2=val2&appid=123&adslotid=lotid1&apiver=api123",
			}
			res, _ := backend.parseHttpResponse(&reqCtx, &backendCtx)
			So(res, ShouldEqual, 6)
			// So(err, ShouldBeNil)
		})
	})
}

func TestMAdxFillAd(t *testing.T) {
	Convey("test fillAd", t, func() {

		var rdata rtb.BidResponse
		var reqCtx mvutil.ReqCtx
		var backendCtx mvutil.BackendCtx

		// guard := Patch(redis.RedisGet, func(key string) (val string, err error) {
		//	return "redis_img_url", nil
		// })
		// defer guard.Unpatch()

		guard := Patch(watcher.AddWatchValue, func(key string, val float64) {
		})
		defer guard.Unpatch()

		// guard1 := Patch(addBackendLog, func(adTmp *corsair_proto.Campaign, demandExt *DemandExt, param *mvutil.Params) {

		// })
		// defer guard1.Unpatch()

		id := "29e8574f-1add-423a-a083-e47c5623841f"
		bidId := "29e8574f-1add-423a-a083-e47c5623841f:1"
		impId := "1"
		price := 1.833
		adid := "255"
		nurl := "http://tracking.url:8766/?"
		adomain := []string{"www.amadoad.ru"}
		iurl := "campaign178-ads263"
		crid := "banner1"
		cat := []string{"IAB1-4", "IAB1-5"}
		attr := []rtb.CreativeAttribute{rtb.CreativeAttribute_VIDEO_IN_BANNER_AUTO_PLAY}
		// h := int32(360)
		// w := int32(480)

		Convey("nil", func() {
			fillAd(&rdata, nil, nil)
		})

		Convey("not nil, OpenApiV3, RewardVideo", func() {
			dp := "deeplink"
			adm := `<?xml version=\"1.0\" encoding=\"UTF-8\"?><VAST version=\"2.0\"><Ad id=\"255\"><InLine><AdSystem>Exchange Inc.</AdSystem><AdTitle>Test Video Ad</AdTitle><Description>test video ad</Description><Survey/><Error/><Impression><![CDATA[http://tracking.url:8767/?type=IMPRESSION&bidRequestId=${AUCTION_ID}&winPrice=${AUCTION_PRICE}]]></Impression><Creatives><Creative AdID=\"banner1\"><Linear><Duration>00:00:16</Duration><TrackingEvents><Tracking event=\"start\"><![CDATA[http://tracking.url:8767/?type=START&bidRequestId=${AUCTION_ID}]]></Tracking><Tracking event=\"firstQuartile\"><![CDATA[http://tracking.url:8767/?type=FIRST_QUARTILE&bidRequestId=${AUCTION_ID}]]></Tracking><Tracking event=\"midpoint\"><![CDATA[http://tracking.url:8767/?type=MIDPOINT&bidRequestId=${AUCTION_ID}]]></Tracking><Tracking event=\"thirdQuartile\"><![CDATA[http://tracking.url:8767/?type=THIRD_QUARTILE&bidRequestId=${AUCTION_ID}]]></Tracking><Tracking event=\"complete\"><![CDATA[http://tracking.url:8767/?type=COMPLETE&bidRequestId=${AUCTION_ID}]]></Tracking></TrackingEvents><VideoClicks><ClickThrough><![CDATA[http: //www.amadoad.ru]]></ClickThrough><ClickTracking><![CDATA[http: //tracking.url:8767/?type=CLICK&bidRequestId=${AUCTION_ID}]]></ClickTracking></VideoClicks><MediaFiles><MediaFile delivery=\"progressive\" type=\"video/mp4\" bitrate=\"251\" width=\"480\" height=\"320\" scalable=\"true\" maintainAspectRatio=\"true\">https://dashboard.amadoad.com/_banners/df/45/df45bd40f590b8d3ea16a38731cc59178847af4b.mp4</MediaFile></MediaFiles></Linear></Creative><Creative><CompanionAds><Companion width=\"300\" height=\"550\"><HTMLResource><![CDATA[<img src='https://node208.fractionalmedia.com/track_img_bin?bid_id=9ff8f27c-fc9d-4d81-9e59-637e45e3a5d8_171019-08__208' width='1' height='1' border='0' /><script type='text/javascript' src='mraid.js'></script><script charset='utf-8' type='text/javascript'>var adInfo = {clickUrl: 'https://node208.fractionalmedia.com/click2_bin?bid_id=9ff8f27c-fc9d-4d81-9e59-637e45e3a5d8_171019-08__208_83557&idfa=499764ae-ddf5-4c9f-a748-61c1fdf62b2c',bidId: '9ff8f27c-fc9d-4d81-9e59-637e45e3a5d8_171019-08__208',campaignId: '83557',creativeId: '27156',playableId: '70114',deviceModel: 'Apple_iPod Touch',idfa: '499764ae-ddf5-4c9f-a748-61c1fdf62b2c',appName: 'Fleet+Battle+-+Apple',node: '208',os: 'iOS',osVersion: 'iOS_10.0',ssp: 'appodeal',adDimensions: '320_480',adSdkVersion: 'appodeal_2.1.4',gameOptions:{'INSTALL_BANNER_TEXT_ID':'1','tower1UpgradeTime':'20','ENDCARD_HEADERTEXT_VICTORY_TOGGLE':'on','ENDCARD_MESSAGE_TIMELIMIT_ID':'1','upgrade_building_message':'on','buttonInstallNowShadowcolor':'#fdff2c','ENDCARD_MESSAGE_TIMELIMIT_TOGGLE':'on','bannerToggle_landscape':'off','didInteractTimeLimitEnabled':'off','purchase_farm_message':'on','inputMessageShowOnPreloader':'off','wall3BaseStrength':'200','inputMessageShowOnGame':'off','INSTALL_BANNER_BOX_FILL_COLOR':'rgba(0,0,0,1)','BOSS_HEALTH':'10','whiteoutEndCardColorAlpha':'1','preloaderScreenCrackLocation':'0','GameID':'FF_Siege','scoreBackgroundAlpha':'0.75','Property2':'off','osID':'iOS','Property3':'2','Property4':'off','wall2Price':'50','Property5':'off','scoreBackgroundHeaderColor':'#FFFFFF','Property6':'off','closeButtonTimer':'5','FileID':'SGF_1.1.0_Venice.js','Property7':'off','Property8':'off','upgradeTowerScore':'10','upgrade_wall_message':'on','ENDCARD_HEADERTEXT_VICTORY_ID':'1','INSTALL_BANNER_TEXT_FILL_COLOR':'rgba(255,255,255,1)','Property1':'off','TrafficID':'bid_playable','showPreloaderOverlay':'on','preloaderStartCountdown':'off','redirectOnGameStart':'0','TROOP_HEALTH':'2','STARTING_GOLD':'100','scoreBackgroundHeaderFontSize':'17','ENEMY_DAMAGE':'0.1','competitiveNameColor':'#e42f43','showScoreBackgroundHeader':'on','TROOP_ATTACK_INTERVAL':'1','endCardType':'0','WALL_ATTACK_DISTANCE':'0.9','endCardLogo':'on','disclaimerToggle':'on','tower2UpgradeTime':'20','banner_clickable_on_show':'off','BOSS_INTERVAL':'4','competitionName':'gamerChick24','hideCloseButtonTime':'5','custom4':'game_action_wall_lvl_3_upgraded','custom5':'game_action_barracks_used','custom2':'milestone_state_game_over_victory','custom3':'game_action_wall_lvl_2_upgraded','killMechScore':'30','closeButtonBackground_landscape':'rgba(0,0,0,1)','custom1':'milestone_state_game_over_defeated','competition':'off','ENDCARD_HEADERTEXT_TIMELIMIT_TOGGLE':'on','INSTALL_BANNER_FONT_SIZE':'40','TROOP_DAMAGE':'1','competitiveAlwaysWin':'off','farm2GoldIncomeBoost':'1.5','disclaimerFontColor':'rgba(128,128,128,1)','tutorial1_pointer':'on','preloaderScreenCrackDuration':'.25','buildTowerScore':'50','ENEMY_MIN_SPAWN_RATE':'2','tower3RateOfFire':'0.5','towerPrice':'40','endCardFullScreen':'off','TestID':'appodeal_TierUS','inputMessageTimeOut':'5','tower2Damage':'0.5','endCardTransitionTime':'0.25','disclaimerMessageID':'1','overrideJSON':'{\"Default_SGF\":{\"milestone\":\"30\"},\"Steroids_SGF\":{\"BASE_GOLD_INCOME\": \"3\",\"BOSS_DAMAGE\": \"1\",\"BOSS_HEALTH\":\"2\",\"BOSS_INTERVAL\":\"0.5\",\"ENEMY_DAMAGE\":\"0.13\",\"ENEMY_HEALTH\": \"0.5\",\"ENEMY_MIN_SPAWN_RATE\": \"0.25\",\"ENEMY_SPAWN_INTERVAL\": \"1\",\"ENEMY_SPAWN_RATE_INCREASE\": \"0.05\",\"farm1GoldIncomeBoost\": \"3.3\",\"farm1UpGradeTime\": \"7\",\"farm2GoldIncomeBoost\": \"4.5\",\"farm3GoldIncomeBoost\": \"6\",\"tower1RateOfFire\": \"0.11\",\"tower2RateOfFire\": \"0.08\",\"tower3Damage\": \"0.75\",\"tower3RateOfFire\": \"0.05\",\"TROOP_ATTACK_INTERVAL\":\"0.13\",\"TROOP_HEALTH\":\"4\"},\"SteroidsHarder_SGF\": {\"ENEMY_GAP\": \"0.5\",\"BASE_GOLD_INCOME\": \"3\",\"BOSS_DAMAGE\": \"1\",\"BOSS_HEALTH\": \"4\",\"BOSS_INTERVAL\": \"0.5\",\"ENEMY_DAMAGE\": \"0.13\",\"ENEMY_HEALTH\": \"1.5\",\"ENEMY_MIN_SPAWN_RATE\": \"0.25\",\"ENEMY_SPAWN_INTERVAL\": \"1\",\"ENEMY_SPAWN_RATE_INCREASE\": \"0.05\",\"farm1GoldIncomeBoost\": \"3.3\",\"farm1UpGradeTime\": \"7\",\"farm2GoldIncomeBoost\": \"4.5\",\"farm3GoldIncomeBoost\": \"6\",\"tower1RateOfFire\": \"0.11\",\"tower2RateOfFire\": \"0.08\",\"tower3Damage\": \"0.75\",\"tower3RateOfFire\": \"0.05\",\"TROOP_ATTACK_INTERVAL\":\"0.13\",\"TROOP_HEALTH\": \"4\"}}','killSoldierScore':'10','tutorial2_pointer':'on','whiteoutEndCard':'off','didInteractTimeLimit':'15','countDownCloseButton':'off','tower1UpgradeCost':'25','toolTipType':'1','farm1UpGradeCost':'10','farm3GoldIncomeBoost':'2','purchase_barrack_message':'on','ENDCARD_HEADERTEXT_DEFEATED_TOGGLE':'on','EndcardTriggersOption':'3','barrackPrice':'30','closeButtonBackground_portrait':'rgba(0,0,0,1)','overrideKey':'Steroids_SGF','farm1GoldIncomeBoost':'1','ENEMY_HEALTH':'2','inputMessageTextID':'0','ENDCARD_MESSAGE_VICTORY_ID':'1','showScoreBackground':'on','tutorial1':'on','tutorial2':'on','farmPrice':'20','MAX_PLAY_TIME':'150','tower2RateOfFire':'0.75','whiteoutEndCardDuration':'1','farm2UpGradeTime':'15','milestone':'30','country':'','whiteoutEndCardColor':'#ffffff','redirectOnGameEnd':'0','tower1Damage':'0.5','tower1RateOfFire':'1','tower2UpgradeCost':'50','gamePlayLogo':'on','ENDCARD_MESSAGE_VICTORY_TOGGLE':'on','disclaimerFontFamily':'arial','wall1BaseStrength':'50','farm1UpGradeTime':'10','ENDCARD_HEADERTEXT_DEFEATED_ID':'1','ENEMY_SPAWN_RATE_INCREASE':'0.5','ENEMY_GAP':'1','preloaderScreenCrackEffect':'off','SizeID':'320x480','EndcardTriggersDelay':'6','EndcardTriggers':'off','BASE_GOLD_INCOME':'1','preloader_type':'3','showCompetitionScore':'on','INSTALL_BANNER_FONT_FAMILY':'Arial','wall3Price':'100','playerNameColor':'#67b252','playerName':'YOU','purchase_tower_message':'on','preloader':'on','ENDCARD_HEADERTEXT_TIMELIMIT_ID':'1','wall2BaseStrength':'100','characterEndCard':'0','bannerToggle_portrait':'off','tower3Damage':'0.5','tutorial_first_action_message':'on','characterToolTip':'0','farm2UpGradeCost':'20','BOSS_DAMAGE':'2','ENDCARD_MESSAGE_DEFEATED_ID':'1','buttonInstallNowScalling':'on','ENEMY_SPAWN_INTERVAL':'8','disclaimerFontSize':'9','tutorial':'on','ENDCARD_MESSAGE_DEFEATED_TOGGLE':'on'},telemetryUrl: 'https://node208.fractionalmedia.com/telim_direct',country:'US'};window.adInfo = adInfo;</script><div id='gamebox'><div id='preloader'>Connecting to Game Server...</div><canvas id='canvas'>Connecting to Game Server...</canvas></div><script type='text/javascript' src='https://cdns3.fractionalmedia.com/fm.creatives/27156/SGF_1.1.0_Venice.js' crossorigin='anonymous'></script>]]></HTMLResource><TrackingEvents><Tracking event=\"creativeView\"><![CDATA[http://myTrackingURL/firstCompanionCreativeView]]></Tracking></TrackingEvents><CompanionClickThrough>http://www.tremormedia.com</CompanionClickThrough></Companion></CompanionAds></Creative></Creatives></InLine></Ad></VAST>`
			rdata = rtb.BidResponse{
				Id: &id,
				Seatbid: []*rtb.BidResponse_SeatBid{
					{
						Bid: []*rtb.BidResponse_SeatBid_Bid{
							{
								Id:      &bidId,
								Impid:   &impId,
								Price:   &price,
								Adid:    &adid,
								Nurl:    &nurl,
								Adm:     &adm,
								Adomain: adomain,
								Iurl:    &iurl,
								Crid:    &crid,
								Cat:     cat,
								Attr:    attr,
								Ext: &rtb.BidResponse_SeatBid_Bid_Ext{
									Deeplink: &dp,
								},
								// H:       &h,
								// W:       &w,
							},
						},
					},
				},
			}

			idfa := "test_idfa"
			devID := "test_devid"
			// adType := int32(42)
			vFlag := int32(1)
			reqID := "req_id"
			uaStr := "test_ua"
			// ntType := enum.NetworkType(9)
			reqCtx = mvutil.ReqCtx{
				ReqParams: &mvutil.RequestParams{
					Param: mvutil.Params{
						RequestPath: mvconst.PATHOpenApiV3,
						IDFA:        idfa,
						AndroidID:   devID,
						AdType:      mvconst.ADTypeRewardVideo,
						VersionFlag: vFlag,
						RequestID:   reqID,
						UserAgent:   uaStr,
						NetworkType: 9,
					},
					UnitInfo: &smodel.UnitInfo{
						Unit: smodel.Unit{
							Orientation: 0,
						},
					},
				},
			}

			backendCtx = mvutil.BackendCtx{
				RespData: []byte("test"),
				Ads: &corsair_proto.BackendAds{
					BackendId: int32(13),
				},
				AdReqKeyValue: "ad_req_val1=val1&ad_req_val2=val2&appid=123&adslotid=lotid1&apiver=api123",
			}

			err := fillAd(&rdata, &reqCtx, &backendCtx)
			So(err, ShouldNotBeNil)
			// So(*(backendCtx.Ads.CampaignList[0]).Price, ShouldEqual, 1.833)
		})

		Convey("not nil, OpenApiV3, NativeVideo", func() {
			adm := `{\"native\":{\"link \": {\"url \": \"http: //i.am.a/URL\"},\"assets\": [{\"id\": 4,\"video\": {\"vasttag\":\"<VAST version='2.0'></VAST>\"}},{\"id\":123,\"required\": 1,\"title\": {\"text\": \"Watch this awesome thing\"}},{\"id\": 1,\"required\": 1,\"img\": {\"url\":\"http://www.myads.com/thumbnail1.png\"}},{\"id\": 8,\"required\": 1,\"img\": {\"url\":\"http://www.myads.com/largethumb1.png\"}},{\"id\": 126,\"required\": 1,\"data\": {\"value\": \"My Brand\"}},{\"id\": 127,\"required\": 1,\"data\": {\"value\": \"Watch all about this awesome story of someone using my product.\"}}]}}`
			rdata = rtb.BidResponse{
				Id: &id,
				Seatbid: []*rtb.BidResponse_SeatBid{
					{
						Bid: []*rtb.BidResponse_SeatBid_Bid{
							{
								Id:      &bidId,
								Impid:   &impId,
								Price:   &price,
								Adid:    &adid,
								Nurl:    &nurl,
								Adm:     &adm,
								Adomain: adomain,
								Iurl:    &iurl,
								Crid:    &crid,
								Cat:     cat,
								Attr:    attr,
								// H:       &h,
								// W:       &w,
							},
						},
					},
				},
			}

			idfa := "test_idfa"
			devID := "test_devid"
			// adType := int32(42)
			vFlag := int32(1)
			reqID := "req_id"
			uaStr := "test_ua"
			// ntType := enum.NetworkType(9)
			reqCtx = mvutil.ReqCtx{
				ReqParams: &mvutil.RequestParams{
					Param: mvutil.Params{
						RequestPath: mvconst.PATHOpenApiV3,
						IDFA:        idfa,
						AndroidID:   devID,
						AdType:      mvconst.ADTypeNative,
						VersionFlag: vFlag,
						RequestID:   reqID,
						UserAgent:   uaStr,
						NetworkType: 9,
					},
					UnitInfo: &smodel.UnitInfo{
						Unit: smodel.Unit{
							Orientation: 0,
						},
					},
				},
			}

			backendCtx = mvutil.BackendCtx{
				RespData: []byte("test"),
				Ads: &corsair_proto.BackendAds{
					BackendId: int32(13),
				},
				AdReqKeyValue: "ad_req_val1=val1&ad_req_val2=val2&appid=123&adslotid=lotid1&apiver=api123",
			}

			err := fillAd(&rdata, &reqCtx, &backendCtx)
			So(err, ShouldNotBeNil)
			// /So(*(backendCtx.Ads.CampaignList[0]).ImageURL, ShouldEqual, "http://www.myads.com/largethumb1.png")
		})
	})
}

// func Test_getVastVideoSize(t *testing.T) {
// 	uri := "uri"
// 	md5Uri := mvutil.Md5(uri)
// 	patch := Patch(redis.LocalRedisGet, func(key string) (val string, err error) {
// 		if key == md5Uri {
// 			return "{\"size\": 123456}", nil
// 		}
// 		if key == mvutil.Md5("testempty") {
// 			return "", nil
// 		}
// 		return "", errors.New("test err")
// 	})
// 	defer patch.Unpatch()
// 	Convey("getVastVideoSize", t, func() {
// 		mediaFile := vast.MediaFile{
// 			URI: uri,
// 		}
// 		size, key, err := getVastVideoSize(mediaFile)
// 		So(err, ShouldBeNil)
// 		So(key, ShouldEqual, md5Uri)
// 		So(size, ShouldEqual, 123456)

// 		mediaFile.URI = "testempty"
// 		size, _, err = getVastVideoSize(mediaFile)
// 		So(err, ShouldNotBeNil)
// 		So(size, ShouldEqual, 0)

// 		mediaFile.URI = "xxx"
// 		_, _, err = getVastVideoSize(mediaFile)
// 		So(err, ShouldNotBeNil)
// 	})
// }

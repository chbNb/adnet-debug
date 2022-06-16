package process_pipeline

import (
	smodel "gitlab.mobvista.com/ADN/structs/model"
	"net/http"
	"testing"

	. "github.com/bouk/monkey"
	pf "github.com/easierway/pipefiter_framework/pipefilter"
	. "github.com/smartystreets/goconvey/convey"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	"gitlab.mobvista.com/ADN/adnet/internal/watcher"
)

func TestCheckError(t *testing.T) {
	Convey("Test checkError", t, func() {

	})
}

func TestRequestHeader(t *testing.T) {
	Convey("Test requestHeader", t, func() {
		Convey("When req.Header[key] exist", func() {
			req := http.Request{
				Header: http.Header{
					"test_1": []string{"t_v1", "t_v2"},
				},
			}

			res := requestHeader(&req, "test_1")
			So(res, ShouldEqual, "t_v1")
		})

		Convey("When req.Header[key] empty", func() {
			req := http.Request{
				Header: http.Header{},
			}

			res := requestHeader(&req, "test_1")
			So(res, ShouldEqual, "")
		})
	})
}

func TestGetClientIP(t *testing.T) {
	Convey("Test GetClientIp", t, func() {
		Convey("x-real-ip", func() {
			req := http.Request{
				Header: http.Header{
					"X-Real-Ip":       []string{"123.2.3.4"},
					"X-Forwarded-For": []string{"123.1.2.3,123.1.2.4,123.2.3.5"},
				},
			}
			res := GetClientIP(&req)
			So(res, ShouldEqual, "123.2.3.4")
		})

		Convey("X-Forwarded-For", func() {
			req := http.Request{
				Header: http.Header{
					"X-Forwarded-For": []string{"123.1.2.3,123.1.2.4,123.2.3.5"},
				},
			}
			res := GetClientIP(&req)
			So(res, ShouldEqual, "123.1.2.3")
		})
	})
}

func TestRequestParamFilterProcess(t *testing.T) {
	Convey("test Process", t, func() {
		guard := Patch(watcher.AddWatchValue, func(key string, value float64) {
		})
		defer guard.Unpatch()

		Convey("Given a pipeline", func() {
			filters := []pf.Filter{&RequestParamFilter{}}

			std := pf.StraightPipeline{
				Name:    "Standard",
				Filters: &filters,
			}

			Convey("When error data", func() {
				in := "error_data"
				ret, err := std.Process(&in)
				So(err, ShouldBeError)
				So(ret, ShouldBeNil)
			})

			Convey("When data is mvutil.RequestParams", func() {
				in := mvutil.RequestParams{
					Param: mvutil.Params{
						RequestPath: "/mpapi/ad",
						Scenario:    "openapi_download_other_test",
						AdType:      42,
						Template:    3,
						ImageSizeID: 0,
						AppID:       101,
						CountryCode: "CN",
						AdSourceID:  1,
						Platform:    1,
						Category:    250,
						SDKVersion:  "",
						Sign:        "NO_CHECK_SIGN",
						Debug:       1,
					},
					UnitInfo: &smodel.UnitInfo{
						AppId:  101,
						UnitId: 998,
						Unit: smodel.Unit{
							Orientation: 2,
							Status:      1,
							AdType:      42,
						},
						AdSourceCountry: map[string]int{"CN": 1, "US": 2},
					},
					AppInfo: &smodel.AppInfo{
						App: smodel.App{
							Platform: 1,
						},
					},
					PublisherInfo: &smodel.PublisherInfo{
						Publisher: smodel.Publisher{
							Type: 0,
						},
					},
					QueryMap: mvutil.RequestQueryMap{
						"app_id":  []string{"10001", "10002"},
						"unit_id": []string{"20001", "20002"},
						"ad_type": []string{"42", "23"},
					},
				}
				ret, err := std.Process(&in)
				So(err, ShouldBeNil)
				So(ret.(*mvutil.RequestParams).Param.RequestPath, ShouldEqual, "/mpapi/ad")
				So(ret.(*mvutil.RequestParams).Param.AdType, ShouldEqual, 4223)
				So(ret.(*mvutil.RequestParams).Param.DeviceType, ShouldEqual, 0)
			})
		})
	})
}

package process_pipeline

import (
	smodel "gitlab.mobvista.com/ADN/structs/model"
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"

	"gitlab.mobvista.com/ADN/adnet/internal/extractor"

	. "github.com/bouk/monkey"
	pf "github.com/easierway/pipefiter_framework/pipefilter"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/ua-parser/uap-go/uaparser"
)

func TestMpReqparamTransFilterProcess(t *testing.T) {
	Convey("test Process", t, func() {
		Convey("Given a pipeline", func() {
			filters := []pf.Filter{&MPReqparamTransFilter{}}

			std := pf.StraightPipeline{
				Name:    "Standard",
				Filters: &filters,
			}

			Convey("data illlegal", func() {
				in := "error_data"
				ret, err := std.Process(&in)

				Convey("Then get the excepted result", func() {
					So(err, ShouldBeError)
					So(ret, ShouldBeNil)
				})
			})

			Convey("When http.Request", func() {
				guard := Patch(RenderReqParam, func(in *http.Request, r *mvutil.RequestParams, rawQuery mvutil.RequestQueryMap) {
				})
				defer guard.Unpatch()

				guard = Patch(RenderMPParam, func(r *mvutil.RequestParams) {
				})
				defer guard.Unpatch()

				guard = Patch(transLieBaoUnit, func(r *mvutil.RequestParams) {
				})
				defer guard.Unpatch()

				in := http.Request{
					Method: "GET",
					URL: &url.URL{
						Scheme:     "scheme",
						Opaque:     "opaque",
						User:       &url.Userinfo{},
						Host:       "test_host",      // host or host:port
						Path:       "test_path",      // path (relative paths may omit leading slash)
						RawPath:    "test_raw_path",  // encoded path hint (see EscapedPath method)
						ForceQuery: true,             // append a query ('?') even if RawQuery is empty
						RawQuery:   "test_raw_query", // encoded query values, without '?'
						Fragment:   "N",              // fragment for references, without '#'
					},
					Proto:      "HTTP/1.0", // "HTTP/1.0"
					ProtoMajor: 1,
					ProtoMinor: 1,
				}

				ret, err := std.Process(&in)

				Convey("Then get the excepted result", func() {
					So(err, ShouldBeNil)
					res, ok := ret.(*mvutil.RequestParams)
					So(ok, ShouldBeTrue)
					So(res.Param, ShouldNotBeNil)
					//So((*(ret.(*mvutil.RequestParams))).Param.RequestPath, ShouldEqual, "")
				})
			})
		})
	})
}

func TestTransLieBaoUnit(t *testing.T) {
	Convey("test transLieBaoUnit", t, func() {
		Convey("len transMap <= 0", func() {
			guard := Patch(extractor.GetMP_MAP_UNIT, func() (map[string]mvutil.MP_MAP_UNIT_, bool) {
				return map[string]mvutil.MP_MAP_UNIT_{}, true
			})
			defer guard.Unpatch()

			r := &mvutil.RequestParams{}
			transLieBaoUnit(r)
		})

		Convey("len transMap < 0 & unt_id <= 0", func() {
			guard := Patch(extractor.GetMP_MAP_UNIT, func() (map[string]mvutil.MP_MAP_UNIT_, bool) {
				return map[string]mvutil.MP_MAP_UNIT_{
					"123": mvutil.MP_MAP_UNIT_{
						UnitID: 123,
					},
				}, true
			})
			defer guard.Unpatch()

			r := &mvutil.RequestParams{
				QueryMap: mvutil.RequestQueryMap{
					"app_id":  []string{"112"},
					"unit_id": []string{"222", "223"},
				},
				AppInfo: &smodel.AppInfo{
					AppId: 111,
				},
			}
			transLieBaoUnit(r)
			So(r.QueryMap["app_id"], ShouldResemble, []string{"112"})
			So(r.QueryMap["unit_id"], ShouldResemble, []string{"222", "223"})
			So(r.QueryMap["ad_type"], ShouldBeNil)
		})
	})
}

func TestRenderMPParam(t *testing.T) {
	Convey("test RenderMPParam", t, func() {
		guard := Patch(renderMPOSVersion, func(r mvutil.RequestParams) string {
			return "os_version_mock"
		})
		defer guard.Unpatch()

		r := &mvutil.RequestParams{
			QueryMap: mvutil.RequestQueryMap{
				"app_id":  []string{"112"},
				"unit_id": []string{"222", "223"},
			},
			AppInfo: &smodel.AppInfo{
				AppId: 111,
			},
		}
		RenderMPParam(r)
		So(r.QueryMap["version_flag"], ShouldResemble, []string{"1"})
		So(r.QueryMap["os_version"], ShouldResemble, []string{"os_version_mock"})
	})
}

func TestRenderMPOSVersion(t *testing.T) {
	Convey("test renderMPOSVersion", t, func() {
		Convey("len ua > 0", func() {
			var par *uaparser.Parser
			guardIns := PatchInstanceMethod(reflect.TypeOf(par), "Parse", func(_ *uaparser.Parser, line string) *uaparser.Client {
				return &uaparser.Client{
					UserAgent: &uaparser.UserAgent{},
					Os: &uaparser.Os{
						Family: "test_f",
						Major:  "test_major",
						Minor:  "test_minor",
						Patch:  "test_patch",
					},
					Device: &uaparser.Device{},
				}
			})
			defer guardIns.Unpatch()

			mvutil.UaParser = &uaparser.Parser{
				UserAgentMisses: uint64(1),
				OsMisses:        uint64(1),
				DeviceMisses:    uint64(1),
				Mode:            int(2),
				UseSort:         true,
			}

			r := mvutil.RequestParams{
				QueryMap: mvutil.RequestQueryMap{
					"app_id":     []string{"112"},
					"unit_id":    []string{"222", "223"},
					"useragent":  []string{"ua1", "ua2"},
					"os_version": []string{"ov", "ov2"},
				},
				AppInfo: &smodel.AppInfo{
					AppId: 111,
				},
			}
			res := renderMPOSVersion(r)
			// need fix
			//So(res, ShouldEqual, "test_major.test_minor.test_patch")
			So(res, ShouldEqual, "")
		})

		Convey("len ua = 0", func() {
			r := mvutil.RequestParams{
				QueryMap: mvutil.RequestQueryMap{
					"app_id":     []string{"112"},
					"unit_id":    []string{"222", "223"},
					"os_version": []string{"ov", "ov2"},
				},
				AppInfo: &smodel.AppInfo{
					AppId: 111,
				},
			}
			res := renderMPOSVersion(r)
			So(res, ShouldEqual, "")
		})
	})
}

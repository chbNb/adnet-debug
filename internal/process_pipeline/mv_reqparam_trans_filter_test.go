package process_pipeline

import (
	"net/http"
	"net/url"
	"testing"

	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"

	. "github.com/bouk/monkey"
	pf "github.com/easierway/pipefiter_framework/pipefilter"
	. "github.com/smartystreets/goconvey/convey"
)

func TestMvReqparamTransFilterProcess(t *testing.T) {
	Convey("test Process", t, func() {
		Convey("Given a pipeline", func() {
			filters := []pf.Filter{&MVReqparamTransFilter{}}

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

			Convey("When param is http.Request", func() {
				guard := Patch(RenderReqParam, func(in *http.Request, r *mvutil.RequestParams, rawQuery mvutil.RequestQueryMap) {
					r.DebugInfo = "mock info"
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
					So(ret.(*mvutil.RequestParams).DebugInfo, ShouldEqual, "mock info")
				})
			})
		})
	})
}

func TestRenderReqParam(t *testing.T) {
	Convey("test RenderReqParam", t, func() {
		var in http.Request
		r := mvutil.RequestParams{}
		var rawQuery mvutil.RequestQueryMap

		guard := Patch(GetClientIP, func(req *http.Request) string {
			return "test_client_ip"
		})
		defer guard.Unpatch()

		in = http.Request{
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
			Header: http.Header{
				"User-Agent": []string{"in_ua"},
			},
		}

		Convey("When no client_ip and no useragent", func() {
			rawQuery = mvutil.RequestQueryMap{}

			RenderReqParam(&in, &r, rawQuery)
			So(r.Param.ClientIP, ShouldEqual, "test_client_ip")
			So(r.Param.UserAgent, ShouldEqual, "in_ua")
		})

		Convey("When has illegal client_ip and no useragent", func() {
			rawQuery = mvutil.RequestQueryMap{
				"client_ip": []string{"raw_client_ip"},
			}

			RenderReqParam(&in, &r, rawQuery)
			So(r.Param.ClientIP, ShouldEqual, "test_client_ip")
			So(r.Param.UserAgent, ShouldEqual, "in_ua")
		})

		Convey("When has right client_ip and right useragent", func() {
			rawQuery = mvutil.RequestQueryMap{
				"client_ip": []string{"100.110.120.130"},
				"useragent": []string{"raw_ua"},
			}

			RenderReqParam(&in, &r, rawQuery)
			So(r.Param.ClientIP, ShouldEqual, "100.110.120.130")
			So(r.Param.UserAgent, ShouldEqual, "raw_ua")
		})
	})
}

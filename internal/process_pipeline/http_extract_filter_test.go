package process_pipeline

import (
	"net/http"
	"net/url"
	"testing"

	pf "github.com/easierway/pipefiter_framework/pipefilter"
	. "github.com/smartystreets/goconvey/convey"
)

func TestHttpExtractFilterProcess(t *testing.T) {
	Convey("test Process", t, func() {
		Convey("Given a pipeline", func() {
			filters := []pf.Filter{&HttpExtractFilter{}}

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
					So((ret.(*http.Request)).Method, ShouldEqual, "GET")
				})
			})
		})
	})
}

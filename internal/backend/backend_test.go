package backend

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"bou.ke/monkey"
	"github.com/golang/protobuf/proto"
	mlogger "github.com/mae-pax/logger"
	"github.com/prometheus/client_golang/prometheus"
	. "github.com/smartystreets/goconvey/convey"
	"gitlab.mobvista.com/ADN/adnet/internal/clients"
	"gitlab.mobvista.com/ADN/adnet/internal/extractor"
	"gitlab.mobvista.com/ADN/adnet/internal/mvconst"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	"gitlab.mobvista.com/ADN/adnet/internal/watcher"
	"gitlab.mobvista.com/ADN/exporter/metrics"
	"gitlab.mobvista.com/ADN/mtg_openrtb/pkg/mtgrtb"
	smodel "gitlab.mobvista.com/ADN/structs/model"
)

func Init() error {
	if len(metrics.CollectorrRegistry) != 0 {
		for _, cc := range metrics.CollectorrRegistry {
			prometheus.Unregister(cc)
		}
		metrics.CollectorrRegistry = metrics.CollectorrRegistry[:0] // 清零
	}

	if err := metrics.InitMetrics("ali-hk", "0.0.1", "./testdata/metrics.yaml", ""); err != nil {
		return err
	}
	return nil
}

func TestInit(t *testing.T) {
	Convey("TestInit", t, func() {
		Convey("init is nil", func() {
			err := Init()
			So(err, ShouldBeNil)
			//So(err, ShouldNotBeNil)
		})
	})
}

func TestNewBackend(t *testing.T) {
	Convey("test NewBackend", t, func() {
		var b *Backend
		var serviceDetail mvutil.ServiceDetail

		Convey("serviceDetail is nil", func() {
			b = NewBackend(nil)
			So(b, ShouldBeNil)
		})

		Convey("b.ID = mvconst.AdInall(3)", func() {
			serviceDetail = mvutil.ServiceDetail{
				Name: "test_adInall",
				ID:   3,
			}
			b = NewBackend(&serviceDetail)
		})

		Convey("b.ID = mvconst.Inmobi(2)", func() {
			serviceDetail = mvutil.ServiceDetail{
				Name: "Inmobi",
				ID:   2,
			}
			b = NewBackend(&serviceDetail)
		})

		Convey("b.ID = mvconst.Mobvista(1)", func() {
			serviceDetail = mvutil.ServiceDetail{
				Name: "test_Mobvista",
				ID:   1,
			}
			b = NewBackend(&serviceDetail)
		})

		Convey("b.ID = 4", func() {
			serviceDetail = mvutil.ServiceDetail{
				Name: "test_Mobvista",
				ID:   4,
			}
			b = NewBackend(&serviceDetail)
		})

		Convey("b.ID = 5", func() {
			serviceDetail = mvutil.ServiceDetail{
				Name: "test_Mobvista",
				ID:   5,
			}
			b = NewBackend(&serviceDetail)
		})

		Convey("b.ID = 6", func() {
			serviceDetail = mvutil.ServiceDetail{
				Name: "test_Mobvista",
				ID:   6,
			}
			b = NewBackend(&serviceDetail)
		})

		Convey("b.ID = 7", func() {
			serviceDetail = mvutil.ServiceDetail{
				Name: "test_Mobvista",
				ID:   7,
			}
			b = NewBackend(&serviceDetail)
		})

		Convey("b.ID = 8", func() {
			serviceDetail = mvutil.ServiceDetail{
				Name: "test_Mobvista",
				ID:   8,
			}
			b = NewBackend(&serviceDetail)
		})
	})
}

func TestBackendComposeHttpRequest(t *testing.T) {
	Convey("test composeHttpRequest", t, func() {
		b := &Backend{}

		Convey("nil", func() {
			err := b.composeHttpRequest(nil, nil, nil)
			So(err, ShouldBeError)
		})

		Convey("b.specific is nil", func() {
			b = &Backend{}
			err := b.composeHttpRequest(nil, nil, nil)
			So(err, ShouldBeError)
		})
	})
}

func TestBackendURL(t *testing.T) {
	Convey("test URL", t, func() {
		b := &Backend{}
		Convey("query len is 0", func() {
			query := ""
			res := b.URL(query)
			So(res, ShouldEqual, "")
		})

		Convey("query len is not 0", func() {
			query := "test_query"
			res := b.URL(query)
			So(res, ShouldEqual, "?test_query")
		})
	})
}

func TestBackendTLSURL(t *testing.T) {
	Convey("test URL", t, func() {
		b := &Backend{}
		Convey("query len is 0", func() {
			query := ""
			res := b.TLSURL(query)
			So(res, ShouldEqual, "")
		})

		Convey("query len is not 0", func() {
			query := "test_query"
			res := b.TLSURL(query)
			So(res, ShouldEqual, "?test_query")
		})
	})
}

func TestFilter(t *testing.T) {
	Convey("test Filter", t, func() {
		b := &Backend{}
		Convey("return nil", func() {
			res := b.filter(nil, nil)
			So(res, ShouldEqual, mvconst.ParamInvalidate)
		})

		Convey("filter region", func() {
			b := &Backend{ID: 2}
			backendCtx := &mvutil.BackendCtx{Region: []string{"HK"}}
			reqCtx := &mvutil.ReqCtx{
				ReqParams: &mvutil.RequestParams{
					Param: mvutil.Params{
						CountryCode: "JP",
					},
				},
			}
			res := b.filter(reqCtx, backendCtx)
			So(res, ShouldEqual, mvconst.BackendRegionFilter)
		})
		Convey("filter context", func() {
			b := &Backend{ID: 2}
			backendCtx := &mvutil.BackendCtx{Region: []string{"HK"}, Content: 2}
			reqCtx := &mvutil.ReqCtx{
				ReqParams: &mvutil.RequestParams{
					Param: mvutil.Params{
						CountryCode:  "HK",
						VideoVersion: "hello",
						VideoAdType:  3,
					},
				},
			}
			res := b.filter(reqCtx, backendCtx)
			So(res, ShouldEqual, mvconst.BackendContentFilter)
		})
		Convey("filter ok", func() {
			b := &Backend{ID: 2}
			backendCtx := &mvutil.BackendCtx{Region: []string{"HK"}, Content: 2}
			reqCtx := &mvutil.ReqCtx{
				ReqParams: &mvutil.RequestParams{
					Param: mvutil.Params{
						CountryCode:  "HK",
						VideoVersion: "hello",
						VideoAdType:  2,
					},
				},
			}
			res := b.filter(reqCtx, backendCtx)
			So(res, ShouldEqual, mvconst.BackendOK)
		})
	})
}

func TestDispatchProxyRequest(t *testing.T) {
	Convey("test dispatchProxyRequest", t, func() {
		//err := Init()
		//t.Error(err)

		serviceDetail := &mvutil.ServiceDetail{
			Name:      "MAdx",
			ID:        17,
			Timeout:   2000,
			Workers:   4,
			HttpURL:   "http://adx-sg.rayjump.com/open_rtb",
			HttpsURL:  "http://adx-sg.rayjump.com/open_rtb",
			Path:      "/open_rtb",
			Method:    "POST",
			UseConsul: false,
		}
		b := NewBackend(serviceDetail)
		madxClient, _ := clients.NewMAdxClient(serviceDetail, nil, nil, nil)
		b.specific = &MAdxBackend{
			Backend{MAdxClient: madxClient},
		}
		Convey("return nil", func() {
			res := b.dispatchProxyRequest(nil, nil)
			So(res, ShouldNotBeNil)
			So(res.FilterCode, ShouldEqual, mvconst.ParamInvalidate)
		})

		Convey("return not ok", func() {
			backendCtx := &mvutil.BackendCtx{Region: []string{"HK"}, Content: 2}
			reqCtx := &mvutil.ReqCtx{
				ReqParams: &mvutil.RequestParams{
					Param: mvutil.Params{
						CountryCode:  "HK",
						VideoVersion: "hello",
						VideoAdType:  3,
						Debug:        2,
					},
					//UnitInfo: &mvutil.UnitInfo{},
					//AppInfo:  &mvutil.AppInfo{},
					UnitInfo:      &smodel.UnitInfo{},
					AppInfo:       &smodel.AppInfo{},
					PublisherInfo: &smodel.PublisherInfo{},
				},
			}
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				id := "id"
				seatbid := []*mtgrtb.BidResponse_SeatBid{}
				w.WriteHeader(200)
				resp := mtgrtb.BidResponse{
					Id:      &id,
					Seatbid: seatbid,
				}
				b, _ := proto.Marshal(&resp)
				w.Write(b)
			}))
			serviceDetail.HttpURL = srv.URL
			serviceDetail.HttpsURL = srv.URL
			defer srv.Close()
			c := mlogger.NewFromYaml("./testdata/watch_log.yaml")
			logger := c.InitLogger("time", "", true, true)
			watcher.Init(logger)

			c = mlogger.NewFromYaml("./testdata/run_log.yaml")
			runLogger := c.InitLogger("time", "level", false, true)
			mvutil.Logger = &mvutil.MidwayLog{Runtime: runLogger}
			watcher.Init(logger)
			var areaConfig mvutil.AreaConfig
			// areaConfig.HttpConfig.UseAdxConsul = false
			mvutil.Config.AreaConfig = &areaConfig
			// extractor.GetFillRateEcpmFloorSwitch()
			guard1 := monkey.Patch(extractor.GetFillRateEcpmFloorSwitch, func() bool {
				return true
			})
			defer guard1.Unpatch()

			guard1 = monkey.Patch(extractor.GetDEBUG_BID_FLOOR_AND_BID_PRICE_CONF, func() map[string]*mvutil.DebugBidFloorAndBidPriceConf {
				return map[string]*mvutil.DebugBidFloorAndBidPriceConf{}
			})
			defer guard1.Unpatch()
			madxClient, _ := clients.NewMAdxClient(serviceDetail, runLogger, nil, nil)
			b.MAdxClient = madxClient
			res := b.dispatchProxyRequest(reqCtx, backendCtx)
			So(res, ShouldNotBeNil)
			// So(res.FilterCode, ShouldEqual, mvconst.HTTPStatusNotOK)
		})
	})
}

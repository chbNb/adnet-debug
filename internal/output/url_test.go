package output

import (
	"gitlab.mobvista.com/ADN/adnet/internal/mvconst"
	"math/rand"
	"reflect"
	"testing"

	. "github.com/bouk/monkey"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/ua-parser/uap-go/uaparser"
	"gitlab.mobvista.com/ADN/adnet/internal/extractor"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	"gitlab.mobvista.com/ADN/adx_common/model"
	smodel "gitlab.mobvista.com/ADN/structs/model"
)

func TestURLRenderImpUrl(t *testing.T) {
	var res string
	Convey("整理 url", t, func() {
		r := &mvutil.RequestParams{
			AppInfo: &smodel.AppInfo{
				App: smodel.App{
					DevinfoEncrypt: 1,
				},
			},
		}
		params := &mvutil.Params{
			GAID:                "test_gaid",
			IDFA:                "test_idfa",
			AndroidID:           "test_andorid_id",
			IMEI:                "test_imei",
			MAC:                 "test_mac",
			ExtfinalPackageName: "test_final_packagename",
			Extra14:             int64(2),
		}
		url := "http://www.test.com?gaid={gaid}&imei={imei}&mac={mac}&package_name={package_name}&sub={subId}&dev={devId}"
		res = RenderImpUrl(url, params, r)
		So(res, ShouldEqual, "http://www.test.com?gaid=test_gaid&imei=test_imei&mac=test_mac&package_name=test_final_packagename&sub=2&dev=test_andorid_id")
	})
}

func TestGetDevId(t *testing.T) {
	var res string
	Convey("整理 DevId", t, func() {
		r := mvutil.RequestParams{
			AppInfo: &smodel.AppInfo{
				App: smodel.App{
					DevinfoEncrypt: 1,
				},
			},
		}
		params := mvutil.Params{
			GAID:                "test_gaid",
			IDFA:                "test_idfa",
			AndroidID:           "test_andorid_id",
			IMEI:                "test_imei",
			MAC:                 "test_mac",
			ExtfinalPackageName: "test_final_packagename",
			Extra14:             int64(2),
		}
		res = getDevId(&r, &params)
		So(res, ShouldEqual, "test_andorid_id")
	})

	Convey("整理 DevId", t, func() {
		r := mvutil.RequestParams{
			AppInfo: &smodel.AppInfo{
				App: smodel.App{},
			},
		}
		params := mvutil.Params{
			GAID:                "test_gaid",
			IDFA:                "test_idfa",
			AndroidID:           "test_andorid_id",
			IMEI:                "test_imei",
			MAC:                 "test_mac",
			ExtfinalPackageName: "test_final_packagename",
			Extra14:             int64(2),
		}
		res = getDevId(&r, &params)
		So(res, ShouldEqual, "")
	})
}

func TestGetSubId(t *testing.T) {
	var res int64
	Convey("整理 Subid，存在 Extfinalsubid", t, func() {
		params := mvutil.Params{
			GAID:                "test_gaid",
			IDFA:                "test_idfa",
			AndroidID:           "test_andorid_id",
			IMEI:                "test_imei",
			MAC:                 "test_mac",
			ExtfinalPackageName: "test_final_packagename",
			Extfinalsubid:       int64(100),
			Extra14:             int64(2),
			AppID:               int64(42),
		}
		res = getSubId(&params)
		So(res, ShouldEqual, 100)
	})

	Convey("整理 Subid，存在 Extra14", t, func() {
		params := mvutil.Params{
			GAID:                "test_gaid",
			IDFA:                "test_idfa",
			AndroidID:           "test_andorid_id",
			IMEI:                "test_imei",
			MAC:                 "test_mac",
			ExtfinalPackageName: "test_final_packagename",
			Extra14:             int64(2),
			AppID:               int64(42),
		}
		res = getSubId(&params)
		So(res, ShouldEqual, 2)
	})

	Convey("整理 Subid，不存在 Extfinalsubid、Extra14", t, func() {
		params := mvutil.Params{
			GAID:                "test_gaid",
			IDFA:                "test_idfa",
			AndroidID:           "test_andorid_id",
			IMEI:                "test_imei",
			MAC:                 "test_mac",
			ExtfinalPackageName: "test_final_packagename",
			AppID:               int64(42),
		}
		res = getSubId(&params)
		So(res, ShouldEqual, 42)
	})
}

func TestGetUrlScheme(t *testing.T) {
	var res string
	var params mvutil.Params

	Convey("httpReq	 == 2，返回 https", t, func() {
		params.HTTPReq = 2
		res = GetUrlScheme(&params)
		So(res, ShouldEqual, "https")
	})

	Convey("否则，返回 http", t, func() {
		params.HTTPReq = 1
		res = GetUrlScheme(&params)
		So(res, ShouldEqual, "http")
	})
}

func TestNeedSchemeHttps(t *testing.T) {
	var res bool
	Convey("返回 true", t, func() {
		res = NeedSchemeHttps(int32(2))
		So(res, ShouldBeTrue)
	})

	Convey("返回 false", t, func() {
		res = NeedSchemeHttps(int32(1))
		So(res, ShouldBeFalse)
	})
}

func TestRenderEndcardUrl(t *testing.T) {
	var res string
	var params mvutil.Params
	var URL string

	Convey("url 长度 <= 0，直接返回", t, func() {
		res = RenderEndcardUrl(&params, URL)
		So(res, ShouldEqual, "")
	})

	Convey("url 长度 <= 0，直接返回", t, func() {
		params.HTTPReq = 2
		URL = "mobvista.com"
		res = RenderEndcardUrl(&params, URL)
		So(res, ShouldEqual, "https://mobvista.com")
	})
}

func TestRenderUrls(t *testing.T) {
	Convey("test RenderUrls", t, func() {
		var r mvutil.RequestParams
		var ad Ad
		var params mvutil.Params
		var campaign smodel.CampaignInfo

		guard := Patch(mvutil.SerializeQ, func(params *mvutil.Params, campaign *smodel.CampaignInfo) string {
			return "serialize_q"
		})
		defer guard.Unpatch()

		guard = Patch(GetJumpUrl, func(r *mvutil.RequestParams, params *mvutil.Params, campaign *smodel.CampaignInfo, queryQ string, ad *Ad) string {
			return "serialize_q"
		})
		defer guard.Unpatch()

		guard = Patch(createClickUrl, func(params *mvutil.Params, isNotice bool) string {
			return "c_click_url"
		})
		defer guard.Unpatch()

		guard = Patch(createImpressionUrl, func(params *mvutil.Params) string {
			return "c_imp_url"
		})
		defer guard.Unpatch()

		guard = Patch(renderAdTracking, func(r *mvutil.RequestParams, ad *Ad, params *mvutil.Params) {
		})
		defer guard.Unpatch()

		Convey("PingMode = 1", func() {
			params.PingMode = 1
			guard = Patch(extractor.GetCHET_URL_UNIT, func() ([]int64, bool) {
				return []int64{}, false
			})
			defer guard.Unpatch()
			guard = Patch(extractor.GetTaoBaoOfferID, func() []int64 {
				return []int64{}
			})
			defer guard.Unpatch()
			RenderUrls(&r, &ad, &params, &campaign)
			So(ad.ClickURL, ShouldEqual, "serialize_q")
			So(ad.NoticeURL, ShouldEqual, "c_click_url")
			So(ad.ImpressionURL, ShouldEqual, "c_imp_url")
		})

		Convey("PingMode != 1", func() {
			params.PingMode = 0
			guard = Patch(extractor.GetCHET_URL_UNIT, func() ([]int64, bool) {
				return []int64{}, false
			})
			defer guard.Unpatch()
			guard = Patch(extractor.GetTaoBaoOfferID, func() []int64 {
				return []int64{}
			})
			defer guard.Unpatch()
			r.Param.RequestType = mvconst.REQUEST_TYPE_OPENAPI_V3
			RenderUrls(&r, &ad, &params, &campaign)
			So(ad.ClickURL, ShouldEqual, "c_click_url")
			So(ad.NoticeURL, ShouldEqual, "c_click_url&useless_notice=1")
			So(ad.ImpressionURL, ShouldEqual, "c_imp_url")
		})
	})
}

func TestRenderAdTracking(t *testing.T) {
	Convey("空", t, func() {
		r := mvutil.RequestParams{}
		ad := Ad{}
		params := mvutil.Params{}
		// campaign := smodel.CampaignInfo{}
		// queryQ := "test_query"
		renderAdTracking(&r, &ad, &params)

		So(ad.AdTrackingPoint, ShouldBeNil)
	})

	Convey("非空", t, func() {
		r := mvutil.RequestParams{}
		ad := Ad{
			AdTracking: CAdTracking{
				Start: []string{"test_start"},
			},
			VideoURL: "v_url",
		}
		params := mvutil.Params{
			AdType:       42,
			VideoVersion: "v_version",
		}
		// campaign := smodel.CampaignInfo{}
		// queryQ := "test_query"
		renderAdTracking(&r, &ad, &params)

		So(*ad.AdTrackingPoint, ShouldResemble, ad.AdTracking)
	})
}

func TestRenderRewardVideoAdTracking(t *testing.T) {
	Convey("test renderRewardVideoAdTracking", t, func() {
		guard := Patch(getAdTrackingUrlList, func(key string, subUrl string, isPercentage bool, rate int) []string {
			return []string{"t_url_1", "t_url_2", "t_url_3"}
		})
		defer guard.Unpatch()

		guard = Patch(renderApkUrls, func(ad *Ad, subUrl string) {
		})
		defer guard.Unpatch()

		guard = Patch(renderPlayPercentage, func(ad *Ad, subUrl string) {
		})
		defer guard.Unpatch()

		var ad Ad
		subURL := "test_sub_url"
		renderRewardVideoAdTracking(&ad, subURL)
		So(ad.AdTracking.Mute, ShouldResemble, []string{"t_url_1", "t_url_2", "t_url_3"})
		So(ad.AdTracking.Endcard_show, ShouldResemble, []string{"t_url_1", "t_url_2", "t_url_3"})
	})
}

func TestRenderApkUrls(t *testing.T) {
	ad := Ad{
		CampaignType: 3,
		AdTracking: CAdTracking{
			Start:            []string{"test_start"},
			ApkDownloadStart: []string{"download_start"},
			ApkDownloadEnd:   []string{"download_end"},
			ApkInstall:       []string{"apk_install"},
		},
	}

	Convey("render APK urls", t, func() {
		subURL := "http://sub_url"
		renderApkUrls(&ad, subURL)
		So(ad.AdTracking.ApkDownloadStart, ShouldResemble, []string{"http://sub_url&key=apk_download_start"})
		So(ad.AdTracking.ApkDownloadEnd, ShouldResemble, []string{"http://sub_url&key=apk_download_end"})
		So(ad.AdTracking.ApkInstall, ShouldResemble, []string{"http://sub_url&key=apk_install"})

	})
}

func TestRenderPlayPercentage(t *testing.T) {
	ad := Ad{
		CampaignType: 3,
		AdTracking: CAdTracking{
			Start:            []string{"test_start"},
			ApkDownloadStart: []string{"download_start"},
			ApkDownloadEnd:   []string{"download_end"},
			ApkInstall:       []string{"apk_install"},
		},
	}

	Convey("renderPlayPercentage", t, func() {
		subURL := "http://sub_url"
		renderPlayPercentage(&ad, subURL)
		exp := []CPlayTracking{
			CPlayTracking{Rate: 0, Url: "http://sub_url&key=play_percentage&rate=0"},
			CPlayTracking{Rate: 25, Url: "http://sub_url&key=play_percentage&rate=25"},
			CPlayTracking{Rate: 50, Url: "http://sub_url&key=play_percentage&rate=50"},
			CPlayTracking{Rate: 75, Url: "http://sub_url&key=play_percentage&rate=75"},
			CPlayTracking{Rate: 100, Url: "http://sub_url&key=play_percentage&rate=100"},
		}
		So(ad.AdTracking.Play_percentage, ShouldResemble, exp)

	})
}

func TestRenderNativeVideoAdTracking(t *testing.T) {
	ad := Ad{
		CampaignType: 3,
		AdTracking: CAdTracking{
			Start:            []string{"test_start"},
			ApkDownloadStart: []string{"download_start"},
			ApkDownloadEnd:   []string{"download_end"},
			ApkInstall:       []string{"apk_install"},
		},
	}
	params := mvutil.Params{
		GAID:                "test_gaid",
		IDFA:                "test_idfa",
		AndroidID:           "test_andorid_id",
		IMEI:                "test_imei",
		MAC:                 "test_mac",
		ExtfinalPackageName: "test_final_packagename",
		AppID:               int64(42),
		ApiVersion:          float64(2.0),
		Extnvt2:             int32(3),
	}

	Convey("renderNativeVideoAdTracking", t, func() {
		subURL := "test.com/suburl"
		renderNativeVideoAdTracking(&ad, subURL, &params)
		So(ad.AdTracking.Endcard_show, ShouldResemble, []string{"test.com/suburl&key=endcard_show"})
		So(ad.AdTracking.Video_Click, ShouldResemble, []string{"test.com/suburl&key=video_click"})
	})

	Convey("renderNativeVideoAdTracking", t, func() {
		params.ApiVersion = float64(1.1)
		params.Extnvt2 = int32(1)
		subURL := "test.com/suburl/2"
		renderNativeVideoAdTracking(&ad, subURL, &params)
		So(ad.AdTracking.Endcard_show, ShouldResemble, []string{"test.com/suburl&key=endcard_show"})
		So(ad.AdTracking.Video_Click, ShouldResemble, []string{"test.com/suburl/2&key=video_click"})
	})
}

func TestRenderNVDefAdTracking(t *testing.T) {
	ad := Ad{
		CampaignType: 3,
		AdTracking: CAdTracking{
			Start:            []string{"test_start"},
			ApkDownloadStart: []string{"download_start"},
			ApkDownloadEnd:   []string{"download_end"},
			ApkInstall:       []string{"apk_install"},
			Play_percentage: []CPlayTracking{
				CPlayTracking{Rate: 100, Url: "url_c_tracking"},
			},
		},
	}

	Convey("renderNativeVideoAdTracking", t, func() {
		subURL := "test.com/suburl_for_renderNVDefAdTracking"
		renderNVDefAdTracking(&ad, subURL)
		So(ad.AdTracking.Start, ShouldResemble, []string{"test.com/suburl_for_renderNVDefAdTracking&key=play_percentage&rate=0"})
		So(ad.AdTracking.First_quartile, ShouldResemble, []string{"test.com/suburl_for_renderNVDefAdTracking&key=play_percentage&rate=25"})
		So(ad.AdTracking.Midpoint, ShouldResemble, []string{"test.com/suburl_for_renderNVDefAdTracking&key=play_percentage&rate=50"})
		So(ad.AdTracking.Third_quartile, ShouldResemble, []string{"test.com/suburl_for_renderNVDefAdTracking&key=play_percentage&rate=75"})
		So(ad.AdTracking.Complete, ShouldResemble, []string{"test.com/suburl_for_renderNVDefAdTracking&key=play_percentage&rate=100"})
	})
}

func TestRenderImpressionT2(t *testing.T) {
	ad := Ad{
		CampaignType: 3,
		AdTracking: CAdTracking{
			Start:            []string{"test_start"},
			ApkDownloadStart: []string{"download_start"},
			ApkDownloadEnd:   []string{"download_end"},
			ApkInstall:       []string{"apk_install"},
			Play_percentage: []CPlayTracking{
				CPlayTracking{Rate: 100, Url: "url_c_tracking"},
			},
		},
	}

	Convey("renderImpressionT2", t, func() {
		subURL := "test.com/renderImpressionT2"
		nvt2 := int32(42)
		renderImpressionT2(&ad, subURL, nvt2)
		So(ad.AdTracking.Impression_t2, ShouldResemble, []string{"test.com/renderImpressionT2&key=impression_t2&nv_t2=42"})
	})
}

func TestGetAdTrackingUrlList(t *testing.T) {
	var res []string

	Convey("is percent", t, func() {
		key := "test_key"
		subURL := "test_sub_url"

		res = getAdTrackingUrlList(key, subURL, true, 1)
		So(res, ShouldResemble, []string{"test_sub_url&key=test_key&rate=1"})
	})

	Convey("is not percent", t, func() {
		key := "test_key"
		subURL := "test_sub_url"

		res = getAdTrackingUrlList(key, subURL, false, 42)
		So(res, ShouldResemble, []string{"test_sub_url&key=test_key"})
	})
}

func TestGetAdTrackingUrl(t *testing.T) {
	var res string

	Convey("is percent", t, func() {
		key := "test_key_1"
		subURL := "test_sub_url_1"

		res = getAdTrackingUrl(key, subURL, true, 1)
		So(res, ShouldResemble, "test_sub_url_1&key=test_key_1&rate=1")
	})

	Convey("is not percent", t, func() {
		key := "test_key_2"
		subURL := "test_sub_url_2"

		res = getAdTrackingUrl(key, subURL, false, 42)
		So(res, ShouldResemble, "test_sub_url_2&key=test_key_2")
	})
}

func TestGetJumpUrl(t *testing.T) {
	var res string
	params := mvutil.Params{
		GAID:                "test_gaid",
		IDFA:                "test_idfa",
		AndroidID:           "test_andorid_id",
		IMEI:                "test_imei",
		MAC:                 "test_mac",
		ExtfinalPackageName: "test_final_packagename",
		AppID:               int64(42),
		ApiVersion:          float64(2.0),
		Extnvt2:             int32(3),
	}
	r := mvutil.RequestParams{
		AppInfo: &smodel.AppInfo{
			App: smodel.App{
				DevinfoEncrypt: 1,
			},
		},
	}
	tag := int32(4)
	vba := int32(1)
	link := "test_link"
	dURL := "test_direct_url?idfa={idfa}&mac={mac}&subId={subId}&dev={devId}&pack={package_name}"
	tURL := "test_tracking_url?idfa={idfa}&mac={mac}&subId={subId}&dev={devId}&pack={package_name}"
	campaign := smodel.CampaignInfo{
		Tag:             tag,
		VbaConnecting:   vba,
		VbaTrackingLink: link,
		CityCodeV2:      map[string][]int32{},
		DirectUrl:       dURL,
		TrackingUrl:     tURL,
	}
	queryQ := "test_query_q"
	ad := &Ad{}

	Convey("extra10 = JUMP_TYPE_SDK_TO_MARKET", t, func() {
		// params.Extra10 = "4"
		// res = GetJumpUrl(r, params, campaign, queryQ)
		// So(res, ShouldEqual, "test_sub_url_1&key=test_key_1&rate=1")
	})
	Convey("extra10 = JUMP_TYPE_CLIENT_SEND_DEVID_PING_SERVER", t, func() {
		params.Extra10 = "11"
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
		guard := Patch(extractor.GetSETTING_CONFIG, func() (mvutil.SETTING_CONFIG, bool) {
			return mvutil.SETTING_CONFIG{}, true
		})
		defer guard.Unpatch()

		guard = Patch(extractor.GetADNET_SWITCHS, func() (map[string]int, bool) {
			return make(map[string]int), true
		})
		defer guard.Unpatch()
		res = GetJumpUrl(&r, &params, &campaign, queryQ, ad)
		So(res, ShouldEqual, "test_direct_url?idfa=test_idfa&mac=test_mac&subId={subId}&dev=test_andorid_id&pack=test_final_packagename")
	})
	Convey("extra10 = JUMP_TYPE_CLIENT_DO_ALL_PING_SERVER", t, func() {
		params.Extra10 = "12"
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
		guard := Patch(extractor.GetSETTING_CONFIG, func() (mvutil.SETTING_CONFIG, bool) {
			return mvutil.SETTING_CONFIG{}, true
		})
		defer guard.Unpatch()

		guard = Patch(extractor.GetADNET_SWITCHS, func() (map[string]int, bool) {
			return make(map[string]int), true
		})
		defer guard.Unpatch()
		res = GetJumpUrl(&r, &params, &campaign, queryQ, ad)
		So(res, ShouldEqual, "test_direct_url?idfa=test_idfa&mac=test_mac&subId={subId}&dev=test_andorid_id&pack=test_final_packagename")
	})
	Convey("extra10 = 其他", t, func() {
		// params.Extra10 = "101010"
		// campaign.Network = nil
		// res = GetJumpUrl(r, params, campaign, queryQ)
		// So(res, ShouldEqual, "test_sub_url_&rate=1")
	})
}

func TestGetMd5(t *testing.T) {
	Convey("空字符串", t, func() {
		res := getMd5("")
		So(res, ShouldEqual, "")
	})

	Convey("非空字符串", t, func() {
		res := getMd5("testMd5_123.ABC")
		So(res, ShouldEqual, "81D31B7C27120BE1EBD98B3F80A877FE")
	})
}

func TestGetMd5Lower(t *testing.T) {
	Convey("空字符串", t, func() {
		res := getMd5Lower("")
		So(res, ShouldEqual, "")
	})

	Convey("非空字符串", t, func() {
		res := getMd5Lower("testMd5_123.ABC")
		So(res, ShouldEqual, "5fb83ced574e70deec59d05063c6a840")
	})
}

func TestGetSha1(t *testing.T) {
	Convey("空字符串", t, func() {
		res := getSha1("")
		So(res, ShouldEqual, "")
	})

	Convey("非空字符串", t, func() {
		res := getSha1("testMd5_123.ABC")
		So(res, ShouldEqual, "823b75ebdbab7e83e00849b25f944232beba27c2")
	})
}

func TestHandleUaParam(t *testing.T) {
	Convey("空字符串", t, func() {
		res := handleUaParam("")
		So(res, ShouldEqual, "")
	})

	Convey("空字符串", t, func() {
		res := handleUaParam("test_hander_ua_[]param & a=hello & b =world")
		So(res, ShouldEqual, "test_hander_ua_%5B%5Dparam+%26+a%3Dhello+%26+b+%3Dworld")
	})
}

func TestGetAdvSubid(t *testing.T) {
	Convey("参数 0", t, func() {
		res := getAdvSubid(int64(0), 0)
		So(res, ShouldEqual, "5e20663dadd1e483")
	})

	Convey("参数其他", t, func() {
		res := getAdvSubid(int64(1000), 42)
		So(res, ShouldEqual, "2f725cce4de9f95e")
	})
}

func TestCreateDirectUrl(t *testing.T) {
	var res string
	var campaign smodel.CampaignInfo

	params := mvutil.Params{
		GAID:                "test_gaid",
		IDFA:                "test_idfa",
		AndroidID:           "test_andorid_id",
		IMEI:                "test_imei",
		MAC:                 "test_mac",
		ExtfinalPackageName: "test_final_packagename",
		AppID:               int64(42),
		ApiVersion:          float64(2.0),
		Extnvt2:             int32(3),
	}
	r := mvutil.RequestParams{
		AppInfo: &smodel.AppInfo{
			App: smodel.App{
				DevinfoEncrypt: 1,
			},
		},
	}

	Convey("direct url 不存在", t, func() {
		res = createDirectUrl(&r, &params, &campaign)
		So(res, ShouldEqual, "")
	})

	Convey("direct url 存在", t, func() {
		link := "test_link"
		dURL := "test_direct_url?idfa={idfa}&mac={mac}&subId={subId}&dev={devId}&pack={package_name}"
		tURL := "test_tracking_url?idfa={idfa}&mac={mac}&subId={subId}&dev={devId}&pack={package_name}"
		campaign := smodel.CampaignInfo{
			VbaTrackingLink: link,
			CityCodeV2:      map[string][]int32{},
			DirectUrl:       dURL,
			TrackingUrl:     tURL,
		}
		campaign.DirectUrl = dURL
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

		guard := Patch(extractor.GetADNET_SWITCHS, func() (map[string]int, bool) {
			return make(map[string]int), true
		})
		defer guard.Unpatch()
		res = createDirectUrl(&r, &params, &campaign)
		So(res, ShouldEqual, "test_direct_url?idfa=test_idfa&mac=test_mac&subId={subId}&dev=test_andorid_id&pack=test_final_packagename")
	})
}

func TestGetMarketUrl(t *testing.T) {
	Convey("test getMarketUrl", t, func() {
		var params mvutil.Params
		var campaign smodel.CampaignInfo
		var res string

		guard := Patch(extractor.GetADSTACKING, func() (mvutil.ADSTACKING, bool) {
			return mvutil.ADSTACKING{
				Android: "test_and_{package_name}",
				IOS:     "test_ios",
			}, true
		})
		defer guard.Unpatch()

		Convey("params.Platform == mvconst.PlatformAndroid", func() {
			params = mvutil.Params{
				Platform:  1,
				GAID:      "test_gaid",
				BackendID: 6,
				RequestID: "test_r_id",
				Domain:    "test_domain",
			}
			pk := "offer_package"
			campaign.PackageName = pk
			res = getMarketUrl(&params, &campaign)
			So(res, ShouldEqual, "test_and_offer_package")
		})

		Convey("params.Platform == mvconst.PlatformIOS", func() {
			params = mvutil.Params{
				Platform:  2,
				GAID:      "test_gaid",
				BackendID: 6,
				RequestID: "test_r_id",
				Domain:    "test_domain",
			}
			pk := "offer_package"
			campaign.PackageName = pk
			res = getMarketUrl(&params, &campaign)
			So(res, ShouldEqual, "test_ios")
		})

		Convey("params.Platform == mvconst.PlatformIOS & packagename is empty", func() {
			params = mvutil.Params{
				Platform:  2,
				GAID:      "test_gaid",
				BackendID: 6,
				RequestID: "test_r_id",
				Domain:    "test_domain",
			}
			campaign.PackageName = ""
			res = getMarketUrl(&params, &campaign)
			So(res, ShouldEqual, "")
		})
	})
}

func TestCreateTrackUrl(t *testing.T) {
	Convey("test CreateTrackUrl", t, func() {
		var r mvutil.RequestParams
		var params mvutil.Params
		var campaign smodel.CampaignInfo
		var trackURL string
		var queryQ string
		var res string

		guard := Patch(extractor.GetTRACK_URL_CONFIG_NEW, func() (map[int32]mvutil.TRACK_URL_CONFIG_NEW, bool) {
			return map[int32]mvutil.TRACK_URL_CONFIG_NEW{
				int32(4): mvutil.TRACK_URL_CONFIG_NEW{
					Android: "and_track_url",
					IOS:     "ios_track_url",
				},
			}, true
		})
		defer guard.Unpatch()

		guard = Patch(getDevId, func(r *mvutil.RequestParams, params *mvutil.Params) string {
			return "test_dev_id"
		})
		defer guard.Unpatch()

		guard = Patch(getSubId, func(params *mvutil.Params) int64 {
			return int64(42)
		})
		defer guard.Unpatch()

		guard = Patch(extractor.GetSYSTEM, func() string {
			return "SA"
		})
		defer guard.Unpatch()

		guard = Patch(extractor.GetDOMAIN, func() string {
			return "test_domain"
		})
		defer guard.Unpatch()

		guard = Patch(extractor.GetDOMAIN_TRACK, func() string {
			return "test_domain_track"
		})
		defer guard.Unpatch()

		guard = Patch(change3SCNDomain, func(trackUrl string, params *mvutil.Params) string {
			return "3s_cdn_" + trackUrl
		})
		defer guard.Unpatch()

		Convey("len(trackUrl) = 0 return ''", func() {
			trackURL = ""
			campaign.TrackingUrl = ""
			res = CreateTrackUrl(&r, &params, &campaign, trackURL, queryQ)
			So(res, ShouldEqual, "")
		})

		Convey("len(trackUrl)> 0 & delteDevid = 0", func() {
			trackURL = ""
			cTURL := "original_tracking_url_{gaid}_{idfa}_{devId}_{imei}_{mac}_{package_name}+{subId}&{clickId}|{creativeId}&&{adType}"
			campaign.TrackingUrl = cTURL
			cNetWork := int32(111)
			campaign.Network = cNetWork
			params.ExtdeleteDevid = 0
			params.GAID = "test_gaid_23"
			params.IDFA = "idfa_13"
			params.HTTPReq = 2
			res = CreateTrackUrl(&r, &params, &campaign, trackURL, queryQ)
			So(res, ShouldEqual, "3s_cdn_original_tracking_url_test_gaid_23_idfa_13_test_dev_id___+42&a_i09M6deI6aSIidMM6aSIidMM6aSI6aSI6aSI6aSIidMM6deIideIW%2BMM6aSIidMM6aSIidMM6deI6aSI6aSIidMM6dMM6aSI6aSIideIideIidMM6aSI6aSI6dMMWUvM6dMM6deI6deI6deIiv%3D%3D|0&&0")
		})

		Convey("len(trackUrl)> 0 & deleteDevid = 1", func() {
			trackURL = ""
			cTURL := "original_tracking_url_{gaid}_{idfa}_{devId}_{imei}_{mac}_{package_name}+{subId}&{clickId}|{creativeId}&&{adType}"
			campaign.TrackingUrl = cTURL
			cNetWork := int32(111)
			campaign.Network = cNetWork
			params.ExtdeleteDevid = 1
			params.GAID = "test_gaid_23"
			params.IDFA = "idfa_13"
			params.HTTPReq = 2
			res = CreateTrackUrl(&r, &params, &campaign, trackURL, queryQ)
			So(res, ShouldEqual, "3s_cdn_original_tracking_url___test_dev_id___+42&a_i09M6deI6aSIidMM6aSIidMM6aSI6aSI6aSI6aSIidMM6deIideIW%2BMM6aSIidMe6aSIidMM6deI6aSI6aSIidMM6dMM6aSI6aSIideIideIidMM6aSI6aSI6dMMWUvM6dMM6deI6deI6deIiv%3D%3D|0&&0")
		})
	})
}

func TestChange3SCNDomain(t *testing.T) {
	Convey("test change3SCNDomain", t, func() {
		var trackURL string
		var params mvutil.Params
		var res string
		guard := Patch(extractor.Get3S_CHINA_DOMAIN, func() (mvutil.CONFIG_3S_CHINA_DOMAIN, bool) {
			return mvutil.CONFIG_3S_CHINA_DOMAIN{
				Domains:  []string{"d1", "d2"},
				CNDomain: "3s_domain",
				Countrys: []string{"c1", "c2"},
			}, true
		})
		defer guard.Unpatch()

		Convey("country not in countryies", func() {
			params.CountryCode = "c2"
			trackURL = "d1.hello.com/test"
			res = change3SCNDomain(trackURL, &params)
			So(res, ShouldEqual, "3s_domain.hello.com/test&mb_trackingcn=1")
		})
	})
}

func TestCreateClickUrl(t *testing.T) {
	Convey("test createClickUrl", t, func() {
		var params mvutil.Params
		// var campaign smodel.CampaignInfo
		// var queryQ string
		var res string

		// campaign = smodel.CampaignInfo{
		// 	CampaignId: int64(1024),
		// 	Status:     int32(1),
		// }
		params = mvutil.Params{
			QueryR: "test_q",
		}
		// queryQ = "param_q"

		guard := Patch(CreateCsp, func(params *mvutil.Params, dspExt string) string {
			return "test_csp"
		})
		defer guard.Unpatch()

		res = createImpressionUrl(&params)
		So(res, ShouldEqual, "http:///impression?k=&p=&q=&x=0&r=test_q&al=&csp=")
	})
}

func TestCreateCsp(t *testing.T) {
	Convey("test CreateCsp", t, func() {
		Convey("return ''", func() {
			params := mvutil.Params{
				QueryR:      "test_q",
				MWFlowTagID: 0,
			}
			res := CreateCsp(&params, "")
			So(res, ShouldEqual, "")
		})

		Convey("creative csp", func() {
			params := mvutil.Params{
				QueryR:      "test_q",
				MWFlowTagID: 1,
			}
			res := CreateCsp(&params, "")
			So(res, ShouldEqual, "6dMe6aSI6aSIidM=")
		})
	})
}

func TestRenderHttpsUrls(t *testing.T) {
	Convey("test RenderHttpsUrls", t, func() {
		var ad Ad
		var params mvutil.Params

		Convey("!NeedSchemeHttps", func() {
			params = mvutil.Params{
				QueryR:      "test_q",
				MWFlowTagID: 1,
				HTTPReq:     int32(1),
			}

			RenderHttpsUrls(&ad, params)
		})

		Convey("NeedSchemeHttps", func() {
			ad = Ad{
				IconURL:  "test_icon",
				ImageURL: "http://res.rayjump.com/image_url",
				VideoURL: "http://res.rayjump.com/video_url",
				ExtImg:   "http://res.rayjump.com/ext_img_url",
			}
			params = mvutil.Params{
				QueryR:      "test_q",
				MWFlowTagID: 1,
				HTTPReq:     int32(2),
			}

			RenderHttpsUrls(&ad, params)
			So(ad.IconURL, ShouldEqual, "test_icon")
			So(ad.ImageURL, ShouldEqual, "https://res-https.rayjump.com/image_url")
			So(ad.VideoURL, ShouldEqual, "http")
			So(ad.ExtImg, ShouldEqual, "https://res-https.rayjump.com/ext_img_url")
		})
	})
}

func TestRenderCDNUrl2Https(t *testing.T) {
	Convey("test renderCDNUrl2Https", t, func() {
		res := renderCDNUrl2Https("http://cdn-adn.rayjump.com&http://d11kdtiohse1a9.cloudfront.net/test-url&http://res.rayjump.com/hello")
		So(res, ShouldEqual, "https://cdn-adn-https.rayjump.com&https://res-https.rayjump.com/test-url&https://res-https.rayjump.com/hello")
	})
}
func TestCreateOnlyImpressionUrl(t *testing.T) {
	Convey("test CreateOnlyImpressionUrl", t, func() {
		var params mvutil.Params
		var res string

		Convey("onlyImpression != 1", func() {
			params.OnlyImpression = 2
			r := mvutil.RequestParams{
				DspExtData: &model.DspExt{},
			}
			res = CreateOnlyImpressionUrl(params, &r)
			So(res, ShouldEqual, "")
		})
	})
}

func TestToNewCDN(t *testing.T) {
	Convey("test toNewCDN", t, func() {
		var ad Ad
		var params mvutil.Params
		Convey("!NeedSchemeHttps", func() {
			params = mvutil.Params{
				HTTPReq: int32(1),
			}
			ad = Ad{
				IconURL:  "test_icon",
				ImageURL: "https://cdn-adn-https.rayjump.com/image_url",
				VideoURL: "https://cdn-adn-https.rayjump.com/video_url",
				ExtImg:   "https://cdn-adn-https.rayjump.com/ext_img_url",
			}
			toNewCDN(&ad, &params)
			So(ad.IconURL, ShouldEqual, "test_icon")
			So(ad.ImageURL, ShouldEqual, "https://cdn-adn-https.rayjump.com/image_url")
			So(ad.VideoURL, ShouldEqual, "https://cdn-adn-https.rayjump.com/video_url")
			So(ad.ExtImg, ShouldEqual, "https://cdn-adn-https.rayjump.com/ext_img_url")
		})

		Convey("NeedSchemeHttps", func() {
			ad = Ad{
				IconURL:  "test_icon",
				ImageURL: "https://cdn-adn-https.rayjump.com/image_url",
				VideoURL: mvutil.Base64Encode("https://cdn-adn-https.rayjump.com/video_url"),
				ExtImg:   "https://cdn-adn-https.rayjump.com/ext_img_url",
			}
			params = mvutil.Params{
				HTTPReq: int32(2),
				AppID:   int64(12345),
			}
			guard := Patch(extractor.GetTO_NEW_CDN_APPS, func() ([]int64, bool) {
				return []int64{12345}, true
			})
			defer guard.Unpatch()

			guardRand := Patch(rand.Intn, func(n int) int {
				return 1
			})
			defer guardRand.Unpatch()

			toNewCDN(&ad, &params)
			So(ad.IconURL, ShouldEqual, "test_icon")
			So(ad.ImageURL, ShouldEqual, "https://cdn-adn-https-new.rayjump.com/image_url")
			So(ad.VideoURL, ShouldEqual, mvutil.Base64Encode("https://cdn-adn-https-new.rayjump.com/video_url"))
			So(ad.ImageURL, ShouldEqual, "https://cdn-adn-https-new.rayjump.com/image_url")
		})

	})
}

func TestRenderNewCDNUrl(t *testing.T) {
	Convey("renderNewCDNUrl", t, func() {
		res := renderNewCDNUrl("https://cdn-adn-https.rayjump.com")
		So(res, ShouldEqual, "https://cdn-adn-https-new.rayjump.com")
	})
}

func TestRenderCreativeUrls(t *testing.T) {
	//guard := Patch(mvutil.RandByRate2, func(rateMap map[string]int) string {
	//	return "confluence.mobvista.com"
	//})
	//defer guard.Unpatch()

	Convey("renderCreativeUrls", t, func() {
		guardCdn := Patch(extractor.GetPUB_CC_CDN, func() (map[string]map[string]string, bool) {
			return map[string]map[string]string{}, false
		})
		defer guardCdn.Unpatch()
		guard := Patch(extractor.GetNEW_CDN_TEST, func() map[string]map[string][]*smodel.CdnSetting {
			return map[string]map[string][]*smodel.CdnSetting{}
		})
		defer guard.Unpatch()

		guard = Patch(extractor.GetADNET_SWITCHS, func() (map[string]int, bool) {
			return map[string]int{}, false
		})
		defer guard.Unpatch()
		ad := &Ad{
			IconURL:  "test_icon",
			ImageURL: "http://res.rayjump.com/image_url",
			VideoURL: mvutil.Base64Encode("http://res.rayjump.com/video_url"),
		}
		r := &mvutil.RequestParams{
			UnitInfo: &smodel.UnitInfo{
				CdnSetting: []*smodel.CdnSetting{
					{
						1,
						"confluence.mobvista.com",
						2,
					},
				},
			},
		}
		params := &mvutil.Params{
			IDFA: "test_idfa",
		}

		RenderCreativeUrls(ad, r, params)
		So(ad.ImageURL, ShouldEqual, "http://confluence.mobvista.com/image_url")
		So(ad.VideoURL, ShouldEqual, mvutil.Base64Encode("http://confluence.mobvista.com/video_url"))
	})
}

func Test_renderEncodeURIComponent(t *testing.T) {
	type args struct {
		durl string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "check render encode uri component",
			args: args{
				durl: "https//tracking.mintegral.com/?creative={encode_uri_component_{中文}}",
			},
			want: "https//tracking.mintegral.com/?creative=%E4%B8%AD%E6%96%87",
		},
		{
			name: "check render encode uri component 2",
			args: args{
				durl: "https//tracking.mintegral.com/?creative={encode_uri_component_{中文}}&creativeId={encode_uri_component_{aa bc}}",
			},
			want: "https//tracking.mintegral.com/?creative=%E4%B8%AD%E6%96%87&creativeId=aa+bc",
		},
		{
			name: "check render encode uri component 3",
			args: args{
				durl: "https//tracking.mintegral.com/?creativeId={encode_uri_component_{aa bc 中文}}&creative={encode_uri_component_{中文}}",
			},
			want: "https//tracking.mintegral.com/?creativeId=aa+bc+%E4%B8%AD%E6%96%87&creative=%E4%B8%AD%E6%96%87",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := renderEncodeURIComponent(tt.args.durl); got != tt.want {
				t.Errorf("renderEncodeURIComponent() = %v, want %v", got, tt.want)
			}
		})
	}
}

package base_filters_test

import (
	"net/http"
	"testing"

	"github.com/smartystreets/goconvey/convey"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/params"
)

func TestHttpQueryMap(t *testing.T) {
	convey.Convey("Test HttpQueryMap", t, func() {
		req, err := http.NewRequest("GET", "http://hb.com/load?ad_num=1&ad_source_id=1&ad_type=287&api_version=1.8&app_id=118827&app_version_name=1.0.72&ats=4MElRozGVTcsY7KbhTcBDrQThrcB4VeXDkxAR0v1RdxBJkVp6N%3D%3D&cache1=15989.485352&cache2=1587.494873&category=0&charging=0&country_code=KH&ct=%5B16777228%5D&display_info=%5B%0A%20%20%7B%0A%20%20%20%20%22cid%22%20%3A%20%22313596486%22%2C%0A%20%20%20%20%22rid%22%20%3A%20%225e5a7c96cd214c7a31278bdy%22%0A%20%20%7D%0A%5D&dmf=129.87109375&dmt=969&exclude_ids=%5B%5D&http_req=2&idfa=4940CDE1-18AC-49D3-9939-146160AC3B2F&idfv=FD3C28C5-3F95-47DA-9D2B-4D0A59149218&keyword=&language=en-KH&mcc=456&mnc=01&model=iPhone7%2C1&network_str=&network_type=9&offset=0&only_impression=1&openidfa=&orientation=1&os_version=12.4.5&package_name=com.azurgames.stackball&ping_mode=1&platform=2&power_rate=24&req_type=2&screen_size=1242.000000x2208.000000&sdk_version=MI_5.8.8&sign=6846f6fc3571a45ba3e38751a6a0bcd0&sub_ip=192.168.1.204&sys_id=7f993d86-f434-532e-8e45-597948978870&timezone=GMT%2B07%3A00&tnum=1&token=831fc07f-3a4a-4a5b-bf96-a0d969e5c16f_vg&ui_orientation=2&unit_id=198608&useragent=Mozilla/5.0%20%28iPhone%3B%20CPU%20iPhone%20OS%2012_4_5%20like%20Mac%20OS%20X%29%20AppleWebKit/605.1.15%20%28KHTML%2C%20like%20Gecko%29%20Mobile/15E148&version_flag=1", nil)
		convey.So(err, convey.ShouldBeNil)

		httpReqData := &params.HttpReqData{}
		httpReqData.Path = req.URL.Path
		httpReqData.Host = req.Host
		httpReqData.QueryData = params.HttpQueryMap(req.URL.Query())
		convey.So(req.URL.Scheme, convey.ShouldEqual, "http")
		convey.So(httpReqData.Host, convey.ShouldEqual, "hb.com")
		convey.So(httpReqData.Path, convey.ShouldEqual, "/load")

		dmt := httpReqData.QueryData.GetString("dmt", true)
		dmf := httpReqData.QueryData.GetString("dmf", true)
		ct := httpReqData.QueryData.GetString("ct", true)
		convey.So(dmt, convey.ShouldEqual, "969")
		convey.So(dmf, convey.ShouldEqual, "129.87109375")
		convey.So(ct, convey.ShouldEqual, "[16777228]")
	})
}

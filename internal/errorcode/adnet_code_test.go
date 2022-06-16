package errorcode

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestAdnCodeString(t *testing.T) {
	Convey("test adn code String", t, func() {
		Convey("返回 success", func() {
			zero := MESSAGE_SUCCESS
			So(zero.String(), ShouldEqual, "ok")
		})
		Convey("返回 ok", func() {
			ok := MESSAGE_SUCCESS_OLD
			So(ok.String(), ShouldEqual, "ok")
		})
		Convey("返回 empty", func() {
			empty := EXCEPTION_RETURN_EMPTY
			So(empty.String(), ShouldEqual, "EXCEPTION_RETURN_EMPTY")
		})
		Convey("返回 empty2", func() {
			empty := EXCEPTION_RETURN_EMPTY_OLD
			So(empty.String(), ShouldEqual, "EXCEPTION_RETURN_EMPTY")
		})
		Convey("返回 empty3", func() {
			empty := EXCEPTION_PARAMS_EMPTY
			So(empty.String(), ShouldEqual, "EXCEPTION_PARAMS_EMPTY")
		})
		Convey("返回 empty4", func() {
			empty := EXCEPTION_PARAMS_ERROR
			So(empty.String(), ShouldEqual, "EXCEPTION_PARAMS_ERROR")
		})
		Convey("返回 empty5", func() {
			empty := EXCEPTION_TIMEOUT
			So(empty.String(), ShouldEqual, "EXCEPTION_TIMEOUT")
		})
		Convey("返回 empty6", func() {
			empty := EXCEPTION_SIGN_ERROR
			So(empty.String(), ShouldEqual, "EXCEPTION_SIGN_ERROR")
		})
		Convey("返回 empty7", func() {
			empty := EXCEPTION_COUNTRY_NOT_ALLOW
			So(empty.String(), ShouldEqual, "EXCEPTION_COUNTRY_NOT_ALLOW")
		})
		Convey("返回 empty8", func() {
			empty := EXCEPTION_TOKEN_FORMAT_ERROR
			So(empty.String(), ShouldEqual, "EXCEPTION_TOKEN_FORMAT_ERROR")
		})
		Convey("返回 empty9", func() {
			empty := EXCEPTION_DOMAIN_ERROR
			So(empty.String(), ShouldEqual, "EXCEPTION_DOMAIN_ERROR")
		})
		Convey("返回 empty10", func() {
			empty := EXCEPTION_NETWORK_ERROR
			So(empty.String(), ShouldEqual, "EXCEPTION_NETWORK_ERROR")
		})
		Convey("返回 empty11", func() {
			empty := EXCEPTION_IP_NOT_ALLOW
			So(empty.String(), ShouldEqual, "EXCEPTION_IP_NOT_ALLOW")
		})
		Convey("返回 empty12", func() {
			empty := EXCEPTION_UNIT_ID_EMPTY
			So(empty.String(), ShouldEqual, "EXCEPTION_UNIT_ID_EMPTY")
		})
		Convey("返回 empty13", func() {
			empty := EXCEPTION_UNIT_NOT_FOUND
			So(empty.String(), ShouldEqual, "EXCEPTION_UNIT_NOT_FOUND")
		})
		Convey("返回 empty14", func() {
			empty := EXCEPTION_UNIT_NOT_FOUND_IN_APP
			So(empty.String(), ShouldEqual, "EXCEPTION_UNIT_NOT_FOUND_IN_APP")
		})
		Convey("返回 empty15", func() {
			empty := EXCEPTION_UNIT_NOT_APPWALL
			So(empty.String(), ShouldEqual, "EXCEPTION_UNIT_NOT_APPWALL")
		})
		Convey("返回 empty16", func() {
			empty := EXCEPTION_UNIT_ADTYPE_ERROR
			So(empty.String(), ShouldEqual, "EXCEPTION_UNIT_ADTYPE_ERROR")
		})
		Convey("返回 empty17", func() {
			empty := EXCEPTION_UNIT_NOT_ACTIVE
			So(empty.String(), ShouldEqual, "EXCEPTION_UNIT_NOT_ACTIVE")
		})
		Convey("返回 empty18", func() {
			empty := EXCEPTION_UNIT_NO_CONFIG_MANAGE_REVENUE
			So(empty.String(), ShouldEqual, "EXCEPTION_UNIT_NO_CONFIG_MANAGE_REVENUE")
		})
		Convey("返回 empty19", func() {
			empty := EXCEPTION_APP_ID_EMPTY
			So(empty.String(), ShouldEqual, "EXCEPTION_APP_ID_EMPTY")
		})
		Convey("返回 empty20", func() {
			empty := EXCEPTION_APP_NOT_FOUND
			So(empty.String(), ShouldEqual, "EXCEPTION_APP_NOT_FOUND")
		})
		Convey("返回 empty21", func() {
			empty := EXCEPTION_APP_BANNED
			So(empty.String(), ShouldEqual, "EXCEPTION_APP_BANNED")
		})
		Convey("返回 empty22", func() {
			empty := EXCEPTION_PUBLISHER_BANNED
			So(empty.String(), ShouldEqual, "EXCEPTION_PUBLISHER_BANNED")
		})
		Convey("返回 empty23", func() {
			empty := EXCEPTION_APP_PLATFORM_ERROR
			So(empty.String(), ShouldEqual, "EXCEPTION_APP_PLATFORM_ERROR")
		})
		Convey("返回 empty24", func() {
			empty := EXCEPTION_PUBLISHER_ID_EMPTY
			So(empty.String(), ShouldEqual, "EXCEPTION_PUBLISHER_ID_EMPTY")
		})
		Convey("返回 empty25", func() {
			empty := EXCEPTION_PUBLISHER_NOT_FOUND
			So(empty.String(), ShouldEqual, "EXCEPTION_PUBLISHER_NOT_FOUND")
		})
		Convey("返回 empty26", func() {
			empty := EXCEPTION_CAMPAIGN_NOT_FOUND
			So(empty.String(), ShouldEqual, "EXCEPTION_CAMPAIGN_NOT_FOUND")
		})
		Convey("返回 empty27", func() {
			empty := EXCEPTION_CAMPAIGN_ID_EMPTY
			So(empty.String(), ShouldEqual, "EXCEPTION_CAMPAIGN_ID_EMPTY")
		})
		Convey("返回 empty28", func() {
			empty := EXCEPTION_CAMPAIGN_NOT_ACTIVE
			So(empty.String(), ShouldEqual, "EXCEPTION_CAMPAIGN_NOT_ACTIVE")
		})
		Convey("返回 empty29", func() {
			empty := EXCEPTION_REQUEST_ID_EMPTY
			So(empty.String(), ShouldEqual, "EXCEPTION_REQUEST_ID_EMPTY")
		})
		Convey("返回 empty30", func() {
			empty := EXCEPTION_SERVICE_REQUEST_OS_VERSION_REQUIRED
			So(empty.String(), ShouldEqual, "EXCEPTION_SERVICE_REQUEST_OS_VERSION_REQUIRED")
		})
		Convey("返回 empty31", func() {
			empty := EXCEPTION_SERVICE_REQUEST_COUNTRY_CODE_REQUIRED
			So(empty.String(), ShouldEqual, "EXCEPTION_SERVICE_REQUEST_COUNTRY_CODE_REQUIRED")
		})
		Convey("返回 empty32", func() {
			empty := EXCEPTION_SERVICE_REQUEST_AD_SOURCE_CLOSED
			So(empty.String(), ShouldEqual, "EXCEPTION_SERVICE_REQUEST_AD_SOURCE_CLOSED")
		})
		Convey("返回 empty33", func() {
			empty := EXCEPTION_GAID_EMPTY
			So(empty.String(), ShouldEqual, "EXCEPTION_GAID_EMPTY")
		})
		Convey("返回 empty34", func() {
			empty := EXCEPTION_IDFA_EMPTY
			So(empty.String(), ShouldEqual, "EXCEPTION_IDFA_EMPTY")
		})
		Convey("返回 empty35", func() {
			empty := EXCEPTION_IMEI_EMPTY
			So(empty.String(), ShouldEqual, "EXCEPTION_IMEI_EMPTY")
		})
		Convey("返回 empty36", func() {
			empty := EXCEPTION_BRAND_EMPTY
			So(empty.String(), ShouldEqual, "EXCEPTION_BRAND_EMPTY")
		})
		Convey("返回 empty37", func() {
			empty := EXCEPTION_MODEL_EMPTY
			So(empty.String(), ShouldEqual, "EXCEPTION_MODEL_EMPTY")
		})
		Convey("返回 empty38", func() {
			empty := EXCEPTION_ANDROIDID_EMPTY
			So(empty.String(), ShouldEqual, "EXCEPTION_ANDROIDID_EMPTY")
		})
		Convey("返回 empty39", func() {
			empty := EXCEPTION_NETWORK_TYPE_EMPTY
			So(empty.String(), ShouldEqual, "EXCEPTION_NETWORK_TYPE_EMPTY")
		})
		Convey("返回 empty40", func() {
			empty := EXCEPTION_ADNUM_SET_NONE
			So(empty.String(), ShouldEqual, "EXCEPTION_ADNUM_SET_NONE")
		})
		Convey("返回 empty41", func() {
			empty := EXCEPTION_CATEGORY_ERROR
			So(empty.String(), ShouldEqual, "EXCEPTION_CATEGORY_ERROR")
		})
	})
}

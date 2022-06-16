package mvutil

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"gitlab.mobvista.com/ADN/adnet/internal/corsair_proto"
)

func TestNewReqCtx(t *testing.T) {
	Convey("New ReqCtx", t, func() {
		res := NewReqCtx()
		valOfRes := *res
		So(valOfRes.MaxTimeout, ShouldEqual, 0)
		So(valOfRes.FlowTagID, ShouldEqual, 0)
		So(valOfRes.RandValue, ShouldEqual, 0)
		So(valOfRes.Elapsed, ShouldEqual, 0)
	})
}

func TestNewBackendCtx(t *testing.T) {
	Convey("New BackendCtx", t, func() {
		res := NewBackendCtx("test_key_name", "test_req_key", "test_package_name", 1.00, 42, []string{"CN"})
		valOfRes := *res
		So(valOfRes.AdReqKeyName, ShouldEqual, "test_key_name")
		So(valOfRes.AdReqKeyValue, ShouldEqual, "test_req_key")
		So(valOfRes.Content, ShouldEqual, 42)
		So(valOfRes.Elapsed, ShouldEqual, -1)
		So(valOfRes.Region, ShouldResemble, []string{"CN"})
		So(valOfRes.RespData, ShouldResemble, make([]byte, 0))
		valOfAds := *valOfRes.Ads
		So(valOfAds, ShouldHaveSameTypeAs, corsair_proto.BackendAds{})
	})
}

func TestNewMobvistaCtx(t *testing.T) {
	Convey("New MobvistaCtx", t, func() {
		res := NewMobvistaCtx()
		valOfRes := *res
		So(valOfRes.AdReqKeyName, ShouldEqual, "")
		So(valOfRes.AdReqKeyValue, ShouldEqual, "")
		So(valOfRes.Content, ShouldEqual, 0)
		So(valOfRes.Elapsed, ShouldEqual, -1)
		So(valOfRes.Region, ShouldResemble, []string{})
		So(valOfRes.RespData, ShouldResemble, []byte{})
		valOfAds := *valOfRes.Ads
		So(valOfAds, ShouldHaveSameTypeAs, corsair_proto.BackendAds{})
	})
}

func TestRequestParams_GetDspExt(t *testing.T) {
	// TODO
}

package mvutil

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"gitlab.mobvista.com/ADN/adnet/internal/mvconst"
	"gitlab.mobvista.com/ADN/adx_common/constant"
)

func TestGetDspAdType(t *testing.T) {
	Convey("GetDspAdType", t, func() {
		So(GetDspAdType(mvconst.ADTypeRewardVideo), ShouldEqual, constant.AD_TYPE_RV)
		So(GetDspAdType(mvconst.ADTypeOnlineVideo), ShouldEqual, constant.AD_TYPE_NV)
		So(GetDspAdType(mvconst.ADTypeInterstitialVideo), ShouldEqual, constant.AD_TYPE_IV)
		So(GetDspAdType(mvconst.ADTypeNativeVideo), ShouldEqual, constant.AD_TYPE_NV)
		So(GetDspAdType(mvconst.ADTypeNativePic), ShouldEqual, constant.AD_TYPE_NI)
		So(GetDspAdType(mvconst.ADTypeSdkBanner), ShouldEqual, constant.AD_TYPE_BANNER)
		So(GetDspAdType(mvconst.ADTypeInteractive), ShouldEqual, constant.AD_TYPE_IA)
		So(GetDspAdType(mvconst.ADTypeSplash), ShouldEqual, constant.AD_TYPE_SPLASH)
		So(GetDspAdType(1), ShouldEqual, 0)
	})
}

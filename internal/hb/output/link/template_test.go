package link

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/constant"
)

func TestNoticeTemplate(t *testing.T) {
	Convey("NoticeTemplate ok", t, func() {
		noticeUrl := NoticeTemplate(constant.LogMidway)
		So(noticeUrl, ShouldEqual, "{sh}://{do}/click?k={k}&mp={mp}")
	})
}

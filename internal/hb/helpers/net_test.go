package helpers_test

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/helpers"
)

func TestIsCorrectIp(t *testing.T) {
	convey.Convey("Test IsCorrectIp", t, func() {
		t := helpers.IsCorrectIp("58.248.11.4")
		convey.So(t, convey.ShouldBeTrue)
		t = helpers.IsCorrectIp("2404:c800:c206:a1:bf22:ae98:9:8")
		convey.So(t, convey.ShouldBeTrue)
		t = helpers.IsCorrectIp("127.0.0.1")
		convey.So(t, convey.ShouldBeTrue)
		t = helpers.IsCorrectIp("255.255.255.0.0.1")
		convey.So(t, convey.ShouldBeFalse)
		t = helpers.IsCorrectIp("ip")
		convey.So(t, convey.ShouldBeFalse)
		t = helpers.IsCorrectIp("ip:ip")
		convey.So(t, convey.ShouldBeFalse)
		t = helpers.IsCorrectIp(":")
		convey.So(t, convey.ShouldBeFalse)
	})
}

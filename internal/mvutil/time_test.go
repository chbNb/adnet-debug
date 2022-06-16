package mvutil

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetLocationStrByUtc(t *testing.T) {
	Convey("GetLocationStrByUtc(0) 返回 Etc/GMT+0", t, func() {
		res, _ := GetLocationStrByUtc(0)
		So(res, ShouldEqual, "Etc/GMT+0")
	})
	Convey("GetLocationStrByUtc(-12) 返回 Etc/GMT+12", t, func() {
		res, _ := GetLocationStrByUtc(-12)
		So(res, ShouldEqual, "Etc/GMT+12")
	})
	Convey("GetLocationStrByUtc(14) 返回 Etc/GMT-14", t, func() {
		res, _ := GetLocationStrByUtc(14)
		So(res, ShouldEqual, "Etc/GMT-14")
	})
}

func TestGetLocationTimeFromTimezone(t *testing.T) {
	time1 := time.Now()
	timeFormat := "2006-01-02 15:04:05"
	Convey("GetLocationTimeFromTimezone(0, time.Now()) 返回time.Now().UTC()", t, func() {
		res, _ := GetLocationTimeFromTimezone(0, time1)
		So(res.Format(timeFormat), ShouldEqual, time1.UTC().Format(timeFormat))
	})

	Convey("GetLocationTimeFromTimezone(-12, time.Now()) ", t, func() {
		res, _ := GetLocationTimeFromTimezone(-12, time1)
		time2 := time1
		local, _ := time.LoadLocation("Etc/GMT+12")
		So(res.Format(timeFormat), ShouldEqual, time2.In(local).Format(timeFormat))
	})

	Convey("GetLocationTimeFromTimezone(12, time.Now()) ", t, func() {
		res, _ := GetLocationTimeFromTimezone(12, time1)
		time2 := time1
		local, _ := time.LoadLocation("Etc/GMT-12")
		So(res.Format(timeFormat), ShouldEqual, time2.In(local).Format(timeFormat))
	})
}

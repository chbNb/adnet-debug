package hot_data

import (
	"strconv"
	"testing"
	"time"

	mlogger "github.com/mae-pax/logger"
	. "github.com/smartystreets/goconvey/convey"
)

var loggerT *mlogger.Log

func beforeTest(t *testing.T) {
	c := mlogger.New()
	loggerT = c.InitLogger("time", "level", false, true)
}

func TestSplitStr(t *testing.T) {
	beforeTest(t)
	Convey("test TestSplitStr", t, func() {
		str := []string{"ddd", "ff", "d", "sssss"}
		// So(len(arr), ShouldEqual, 11)
		InitActiveDataCollecter(str, loggerT)
		for k := range activeDataCollecter {
			So(k, ShouldBeIn, []string{"ddd", "ff", "d", "sssss"})
		}
		So(len(activeDataCollecter), ShouldEqual, 4)
	})
}

func TestInitActiveData(t *testing.T) {
	beforeTest(t)
	Convey("test TestInitActiveData", t, func() {
		str := []string{"AAA", "BBB", "CCC"}
		InitActiveDataCollecter(str, loggerT)
		for k := range activeDataCollecter {
			So(k, ShouldBeIn, []string{"AAA", "BBB", "CCC"})
		}
		So(len(activeDataCollecter), ShouldEqual, 3)
	})
}

func TestCanUse(t *testing.T) {
	beforeTest(t)
	Convey("test TestCanUse", t, func() {
		str := []string{"AAA", "BBB", "CCC"}
		InitActiveDataCollecter(str, loggerT)
		So(canUse("AAA"), ShouldBeTrue)
		So(canUse("BBB"), ShouldBeTrue)
		So(canUse("CCC"), ShouldBeTrue)
		So(canUse("ddd"), ShouldBeFalse)
	})
}

func TestAddToActiveDataCollecter(t *testing.T) {
	beforeTest(t)
	Convey("test AddToActiveDataCollecter", t, func() {
		str := []string{"AAA", "BBB", "CCC"}
		InitActiveDataCollecter(str, loggerT)

		AddToActiveDataCollecter("ddd", 123)
		So(activeDataCollecter["ddd"], ShouldBeNil)

		AddToActiveDataCollecter("AAA", 123)
		So(len(activeDataCollecter["AAA"].ids), ShouldEqual, 1)
		AddToActiveDataCollecter("AAA", 123)
		So(len(activeDataCollecter["AAA"].ids), ShouldEqual, 1)
		AddToActiveDataCollecter("AAA", 4567)
		So(len(activeDataCollecter["AAA"].ids), ShouldEqual, 2)

		AddToActiveDataCollecter("BBB", 4567)
		AddToActiveDataCollecter("BBB", 123)
		AddToActiveDataCollecter("BBB", 4567)
		AddToActiveDataCollecter("BBB", 123)
		AddToActiveDataCollecter("BBB", 123)
		AddToActiveDataCollecter("BBB", 1234)
		So(len(activeDataCollecter["BBB"].ids), ShouldEqual, 3)

		So(len(activeDataCollecter["AAA"].ids), ShouldEqual, 2)

	})
}

func TestGetRedisKeyArrForGetData(t *testing.T) {
	Convey("test getRedisKeyArrForGetData", t, func() {
		// InitActiveDataCollecter("AAA,BBB,CCC")
		var suffix string
		if time.Now().Hour() < 2 { // 凌晨2点前使用前一天的数据
			suffix = time.Now().Add(-24 * time.Hour).Format("2006-01-02")
		} else {
			suffix = time.Now().Format("2006-01-02")
		}
		keyArr := getRedisKeyArrForGetData("AAA")
		// fmt.Println(keyArr)
		So(keyArr[6], ShouldEqual, "active_data_AAA_"+suffix)
		now := time.Now()
		for i := 0; i < 6; i++ {
			t := now.Add(time.Duration(-1*(i*10+1)) * time.Minute)
			hourStr := strconv.Itoa(t.Hour())
			minuteStr := strconv.Itoa(t.Minute() / 10)
			suffix = t.Format("2006-01-02") + "-" + hourStr + "-" + minuteStr
			So(keyArr[i], ShouldEqual, "active_data_AAA_"+suffix)
		}

	})
}
func TestGetRedisKey(t *testing.T) {
	Convey("test getRedisKey", t, func() {
		// InitActiveDataCollecter("AAA,BBB,CCC")
		now := time.Now()
		hourStr := strconv.Itoa(now.Hour())
		minuteStr := strconv.Itoa(now.Minute() / 10)
		suffix := now.Format("2006-01-02") + "-" + hourStr + "-" + minuteStr
		key := getRedisKey("AAA")
		So(key, ShouldEqual, "active_data_AAA_"+suffix)
	})
}

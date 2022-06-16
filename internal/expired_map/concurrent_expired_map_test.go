package expired_map

import (
	"github.com/easierway/concurrent_map"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestCreateConcurrentExpiredMap(t *testing.T) {
	Convey("test CreateConcurrentExpiredMap", t, func() {
		cem := CreateConcurrentExpiredMap(99, 2, 1)

		cem.Set(concurrent_map.I64Key(12345), "aaaaaa111")
		v, ok := cem.Get(concurrent_map.I64Key(12345))
		So(v, ShouldEqual, "aaaaaa111")
		So(ok, ShouldBeTrue)
		//set nil value
		key := concurrent_map.I64Key(12345)
		cem.Set(key, nil)
		v, ok = cem.Get(key)
		So(v, ShouldBeNil)
		So(ok, ShouldBeTrue)
		//delete
		cem.Del(concurrent_map.I64Key(12345))
		v, ok = cem.Get(key)
		So(v, ShouldBeNil)
		So(ok, ShouldBeFalse)

		testArr := []struct {
			k    int64
			v    interface{}
			res1 interface{}
			res2 bool
		}{
			{1, "101", "101", true},
			{1, nil, nil, true},
			{1234567890123456789, nil, nil, true},       //19位ID
			{1234567890123456789, "test", "test", true}, //19位ID
			{1234567890123456789, nil, nil, true},       //19位ID
		}

		for _, test := range testArr {
			cem.Set(concurrent_map.I64Key(test.k), test.v)
			v, ok = cem.Get(concurrent_map.I64Key(test.k))
			So(v, ShouldEqual, test.res1)
			So(ok, ShouldEqual, test.res2)
		}

		type campaignInfo struct {
			Id           int64
			CampaignName string
		}

		CampaignObj := campaignInfo{
			1234567890123456789,
			"testname",
		}

		cem.Set(concurrent_map.I64Key(CampaignObj.Id), &CampaignObj)
		cObj, find := cem.Get(concurrent_map.I64Key(CampaignObj.Id))
		campInfo, _ := cObj.(*campaignInfo)
		So(campInfo.Id, ShouldEqual, CampaignObj.Id)
		So(campInfo.CampaignName, ShouldEqual, CampaignObj.CampaignName)
		So(find, ShouldBeTrue)

		v, ok = cem.Get(concurrent_map.I64Key(404))
		So(v, ShouldEqual, nil)
		So(ok, ShouldEqual, false)

		//过期时间   会比实际多一秒
		time.Sleep(time.Second * 2)

		for _, test := range testArr {
			v, ok = cem.Get(concurrent_map.I64Key(test.k))
			So(v, ShouldBeNil)
			So(ok, ShouldBeFalse)
		}

		cObj, find = cem.Get(concurrent_map.I64Key(CampaignObj.Id))
		So(cObj, ShouldBeNil)
		So(ok, ShouldBeFalse)

		time.Sleep(time.Second)

		//暂时不会用到的方法
		cem.SetWithTime(concurrent_map.I64Key(333), "7777", 2)

		cem.Close()
	})
}

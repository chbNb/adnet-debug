package mvutil

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetAdSourceID(t *testing.T) {
	unitInfo := UnitInfo{
		AdSourceCountry: map[string]int{
			"CN": 1,
			"US": 2,
			"JP": 3,
		},
		AdSourceData: []AdSourceDataInfo{
			[]AdSourceDataEntry{
				AdSourceDataEntry{
					AdSourceId: 1001,
					Status:     1,
					Priority:   10,
				},
				AdSourceDataEntry{
					AdSourceId: 1002,
					Status:     1,
					Priority:   9,
				},
			},
			[]AdSourceDataEntry{
				AdSourceDataEntry{
					AdSourceId: 2001,
					Status:     1,
					Priority:   10,
				},
				AdSourceDataEntry{
					AdSourceId: 2002,
					Status:     1,
					Priority:   9,
				},
			},
		},
	}

	Convey("返回 AdSourceID", t, func() {
		res := unitInfo.GetAdSourceID("CN", 2001)
		So(res, ShouldEqual, 2001)
	})

	Convey("不存在的 AdSourceID", t, func() {
		res := unitInfo.GetAdSourceID("CN", 2004)
		So(res, ShouldEqual, 0)
	})

	Convey("不存在的 AdSourceIdList", t, func() {
		res := unitInfo.GetAdSourceID("NONE", 2001)
		So(res, ShouldEqual, 0)
	})

	Convey("返回 ADSourceAPIOffer", t, func() {
		res := unitInfo.GetAdSourceID("CN", 0)
		So(res, ShouldEqual, 1)
	})
}

func TestGetAdSourceIDList(t *testing.T) {
	Convey("不存在的 country", t, func() {
		unitInfo := UnitInfo{
			AdSourceCountry: map[string]int{
				"CN": 1,
				"US": 2,
				"JP": 3,
			},
		}

		res := unitInfo.GetAdSourceIDList("TW")
		So(res, ShouldBeNil)
	})

	Convey("超出 AdSourceData 范围", t, func() {
		unitInfo := UnitInfo{
			AdSourceCountry: map[string]int{
				"CN": 1,
				"US": 2,
				"JP": 3,
			},
			AdSourceData: []AdSourceDataInfo{
				[]AdSourceDataEntry{
					AdSourceDataEntry{
						AdSourceId: 1001,
						Status:     1,
						Priority:   10,
					},
					AdSourceDataEntry{
						AdSourceId: 1002,
						Status:     1,
						Priority:   9,
					},
				},
			},
		}

		res := unitInfo.GetAdSourceIDList("CN")
		So(res, ShouldBeNil)
	})

	Convey("返回正常 AdSourceID", t, func() {
		unitInfo := UnitInfo{
			AdSourceCountry: map[string]int{
				"CN": 1,
				"US": 2,
				"JP": 3,
			},
			AdSourceData: []AdSourceDataInfo{
				[]AdSourceDataEntry{
					AdSourceDataEntry{
						AdSourceId: 1001,
						Status:     1,
						Priority:   10,
					},
					AdSourceDataEntry{
						AdSourceId: 1002,
						Status:     1,
						Priority:   9,
					},
				},
				[]AdSourceDataEntry{
					AdSourceDataEntry{
						AdSourceId: 2001,
						Status:     1,
						Priority:   10,
					},
					AdSourceDataEntry{
						AdSourceId: 2002,
						Status:     1,
						Priority:   9,
					},
				},
			},
		}

		res := unitInfo.GetAdSourceIDList("CN")
		So(res, ShouldResemble, []int{2001, 2002})
	})
}

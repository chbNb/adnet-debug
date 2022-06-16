package helpers

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetRequestID(t *testing.T) {
	Convey("request ok", t, func() {
		reqId := GetRequestID()
		So(reqId, ShouldNotBeNil)
		So(len(reqId), ShouldEqual, 24)
	})

}

func TestGetBidID(t *testing.T) {
	Convey("bid id", t, func() {
		bidIds := make([]string, 0)
		for i := 0; i < 100000; i++ {
			bid := GetBidId(GetRequestID())
			// log.Printf("bid: %s\n", bid)
			bidIds = append(bidIds, bid)
		}
		So(len(getDuplicateElement(bidIds)), ShouldEqual, 0)
	})
}

func getDuplicateElement(bidIds []string) []string {
	result := make([]string, 0, len(bidIds))
	temp := map[string]struct{}{}
	for _, item := range bidIds {
		_, ok := temp[item]
		if ok {
			result = append(result, item)
		} else {
			temp[item] = struct{}{}
		}
	}
	return result
}

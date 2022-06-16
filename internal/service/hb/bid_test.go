package hb

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestCreateBidHandler(t *testing.T) {
	Convey("createbidHandler ok", t, func() {
		handler := CreateBidHandler()
		So(handler, ShouldNotBeNil)
	})
}

// func TestBidServeHTTP(t *testing.T) {
// 	Convey("serveHTTP ok", t, func() {
// 		request := httptest.NewRequest("Get", "/bid", nil)
// 		writer := httptest.NewRecorder()
// 		handler := CreateBidHandler()
// 		So(handler, ShouldNotBeNil)
// 		handler.ServeHTTP(writer, request)
// 	})
// }

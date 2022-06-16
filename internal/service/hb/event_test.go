package hb

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCreateEventHandler(t *testing.T) {
	Convey("createEventHandler ok", t, func() {
		handler := CreateEventHandler()
		So(handler, ShouldNotBeNil)
	})
}

// func TestEventServeHTTP(t *testing.T) {
// 	Convey("serveHTTP ok", t, func() {
// 		request := httptest.NewRequest("Get", "/win", nil)
// 		writer := httptest.NewRecorder()
// 		handler := CreateEventHandler()
// 		So(handler, ShouldNotBeNil)
// 		handler.ServeHTTP(writer, request)
// 	})
// }

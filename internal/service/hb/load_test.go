package hb

import (
	//"net/http/httptest"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCreateLoadHandler(t *testing.T) {
	Convey("createLoadHandler ok", t, func() {
		handler := CreateLoadHandler()
		So(handler, ShouldNotBeNil)
	})
}

//func TestLoadServeHTTP(t *testing.T) {
//	Convey("serveHTTP ok", t, func() {
//		request := httptest.NewRequest("Get", "/load", nil)
//		writer := httptest.NewRecorder()
//		handler := CreateLoadHandler()
//		So(handler, ShouldNotBeNil)
//		handler.ServeHTTP(writer, request)
//	})
//}

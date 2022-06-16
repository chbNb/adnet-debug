package hb

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCreateHTTPConnector(t *testing.T) {
	Convey("create http connector ok", t, func() {
		conn := CreateHTTPConnector(":9102")
		So(conn, ShouldNotBeNil)
	})
}

func TestInitRouter(t *testing.T) {
	Convey("init router", t, func() {
		conn := CreateHTTPConnector(":9102")
		So(conn, ShouldNotBeNil)
		conn.InitRouter()
		// ctx, cancel := context.WithCancel(context.Background())
		// defer cancel()
		// conn.Start(ctx)
		// conn.Stop(ctx)
	})
}

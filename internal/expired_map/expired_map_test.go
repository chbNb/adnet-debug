package expired_map

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestCreateExpiredMap(t *testing.T) {

	type testStruct struct {
		a string
	}

	Convey("test TestCreateExpiredMap", t, func() {
		em := CreateExpiredMap(2)
		em2 := CreateExpiredMap(5)

		trickTime := em.GetTickerTime()
		So(trickTime, ShouldEqual, 2)

		key1 := 12345
		key2 := 78911
		key3 := "3333"

		em.Set(key1, "aaaaaa", 1)
		v, ok := em.Get(key1)
		So(v, ShouldEqual, "aaaaaa")
		So(ok, ShouldBeTrue)

		em.Set(key1, nil, 1)
		v, ok = em.Get(key1)
		So(v, ShouldBeNil)
		So(ok, ShouldBeTrue)

		//em.Delete(key1)
		em.Remove(key1)
		v, ok = em.Get(key1)
		So(v, ShouldBeNil)
		So(ok, ShouldBeFalse)

		em.Set(key2, nil, 1)
		em2.Set(key2, nil, 1)

		em.Set(key1, "aaaaaa", 3)
		ttl := em.TTL(key1)
		So(ttl, ShouldEqual, 3)
		time.Sleep(time.Second) ////////////////SLEEP//////////////
		ttl = em.TTL(key1)
		So(ttl, ShouldEqual, 2)
		em.Set(key1, "aaaaaa", 5)
		ttl = em.TTL(key1)
		So(ttl, ShouldEqual, 5)
		em.Set(key1, "aaaaaa", 1)
		ttl = em.TTL(key1)
		So(ttl, ShouldEqual, 1)

		leng := em.Length()
		So(leng, ShouldEqual, 2)

		time.Sleep(time.Second * 2)
		//leng = em.Length()
		//So(leng, ShouldEqual, 1)

		v, ok = em.Get(key2)
		So(v, ShouldBeNil)
		So(ok, ShouldBeFalse)

		ttl = em.TTL(key2)
		So(ttl, ShouldEqual, -1)
		ttl = em2.TTL(key2)
		So(ttl, ShouldEqual, -1)

		em2.Set(404, "AAAAA", 0) //cannot set

		// new test
		var test *testStruct
		arr := make(map[string]*testStruct)
		if test, ok = arr[key3]; !ok {
			test = nil
		}
		em.Set(key3, test, 10)
		vv, ok := em.Get(key3)
		So(vv, ShouldBeNil)
		So(ok, ShouldBeTrue)

		em.Clear()
		em.Close()
		em2.Stop()
	})
}

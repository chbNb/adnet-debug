package expired_map

import (
	. "github.com/smartystreets/goconvey/convey"
	"math/rand"
	"testing"
	"time"
)

func TestTryLock(t *testing.T) {
	Convey("test TryLock", t, func() {
		cl := CreateCoccurrentLock()
		key1 := 12345
		key2 := 789

		lckOk := cl.TryLock(key1)
		So(lckOk, ShouldBeTrue)

		lckOk = cl.TryLock(key1)
		So(lckOk, ShouldBeFalse)

		lckOk = cl.TryLock(key1)
		So(lckOk, ShouldBeFalse)

		lckOk = cl.TryLock(key2)
		So(lckOk, ShouldBeTrue)

		cl.UnLock(key1) //解锁

		lckOk = cl.TryLock(key1)
		So(lckOk, ShouldBeTrue)

		lckOk = cl.TryLock(key2)
		So(lckOk, ShouldBeFalse)
	})
}

//goos: darwin
//goarch: amd64
//pkg: gitlab.mobvista.com/ADN/adnet/internal/expired_map
//	3000000	       799 ns/op
//  ten repeat
//  1000000	      1805 ns/op
//	PASS
func BenchmarkConccurentLock_TryLock(b *testing.B) {
	cl := CreateCoccurrentLock()
	//var key int64 = 89988812321117283
	r := rand.New(rand.NewSource(time.Now().Unix()))
	key := r.Float64()

	for i := 0; i < b.N; i++ {
		func() {
			cl.TryLock(key)
			//cl.TryLock(key)
			//cl.TryLock(key)
			//cl.TryLock(key)
			//cl.TryLock(key)
			//cl.TryLock(key)
			//cl.TryLock(key)
			//cl.TryLock(key)
			//cl.TryLock(key)
			//cl.TryLock(key)
			defer cl.UnLock(key)
		}()
	}
}

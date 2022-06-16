package backend

import (
	"testing"

	. "github.com/bouk/monkey"
	mlogger "github.com/mae-pax/logger"
	. "github.com/smartystreets/goconvey/convey"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
)

func TestNewBackendManage(t *testing.T) {
	Convey("test NewBackendManage", t, func() {
		res := NewBackendManage()
		exp := BackendManage{
			Backends: map[int]*Backend{},
		}
		So(*res, ShouldResemble, exp)
	})
}

func TestAddBackend(t *testing.T) {
	Convey("test AddBackend", t, func() {
		guard := Patch(NewBackend, func(serviceDetail *mvutil.ServiceDetail) *Backend {
			return &Backend{}
		})
		defer guard.Unpatch()

		manager := &BackendManage{
			Backends: map[int]*Backend{},
		}
		serviceDetail := mvutil.ServiceDetail{
			ID:   2,
			Name: "test_1",
		}
		c := mlogger.NewFromYaml("./testdata/run_log.yaml")
		runLogger := c.InitLogger("time", "level", false, true)
		err := manager.AddBackend(&serviceDetail, runLogger, nil, nil)
		So(err, ShouldBeNil)
	})
}

// func TestGetBackendName(t *testing.T) {
// 	Convey("test GetBackendName", t, func() {
// 		var res string

// 		manager := &BackendManage{
// 			Backends: map[int]*Backend{
// 				1: &Backend{
// 					Name: "test_name_1",
// 				},
// 			},
// 		}
// 		id := int(1)

// 		res, _ = manager.GetBackendName(id)
// 		So(res, ShouldEqual, "test_name_1")
// 	})
// }

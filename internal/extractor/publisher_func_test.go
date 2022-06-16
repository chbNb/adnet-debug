package extractor

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	smodel "gitlab.mobvista.com/ADN/structs/model"
)

func TestLegalPublisher(t *testing.T) {
	Convey("test legalPublisher", t, func() {
		Convey("return false status != 1", func() {
			publisherInfo := smodel.PublisherInfo{
				Publisher: smodel.Publisher{
					Status: 21,
				},
			}
			res := legalPublisher(&publisherInfo)
			So(res, ShouldBeFalse)
		})

		Convey("return false apikey == 0", func() {
			publisherInfo := smodel.PublisherInfo{
				Publisher: smodel.Publisher{
					Status: 1,
					Apikey: "",
				},
			}
			res := legalPublisher(&publisherInfo)
			So(res, ShouldBeFalse)
		})

		Convey("return true", func() {
			publisherInfo := smodel.PublisherInfo{
				Publisher: smodel.Publisher{
					Status: 1,
					Apikey: "test",
				},
			}
			res := legalPublisher(&publisherInfo)
			So(res, ShouldBeTrue)
		})
	})
}

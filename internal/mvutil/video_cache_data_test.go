package mvutil

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestVideoMetaDataWithoutUrl(t *testing.T) {
	Convey("VideoMetaDataWithoutUrl ok", t, func() {
		cacheData := "255_00:00:15_300x400"
		data, err := VideoMetaDataWithoutUrl(cacheData)
		So(err, ShouldBeNil)
		So(data, ShouldNotBeNil)
	})
	Convey("VideoMetaDataWithUrl ok", t, func() {
		cacheData := "255_00:00:15_300x400_xxxxxx"
		data, err := VideoMetaDataWithUrl(cacheData)
		So(err, ShouldBeNil)
		So(data, ShouldNotBeNil)
	})
}

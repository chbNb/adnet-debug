package helpers

import (
	"errors"
	"net/url"
	"strconv"
	"strings"
)

func DecodeVideoInfosWithCdn(data string) (int32, int32, string, string, error) {
	size := int32(0)
	length := int32(0)
	resolution := ""
	cdnUrl := ""
	vInfos := strings.Split(data, "_")
	if len(vInfos) != 4 {
		return size, length, resolution, cdnUrl, errors.New("input params not in the correct format data=" + data)
	}

	vsize, err := strconv.ParseInt(vInfos[0], 10, 32)
	if err != nil {
		return size, length, resolution, cdnUrl, errors.New("Convert size data error:" + err.Error())
	}
	size = int32(vsize)
	length = int32(GetVideoLength(vInfos[1]))
	resolution = strings.TrimSpace(vInfos[2])
	videoUrl, err := url.QueryUnescape(strings.TrimSpace(vInfos[3]))
	if err != nil {
		return size, length, resolution, cdnUrl, errors.New("Convert cdnUrl data error:" + err.Error())
	}
	cdnUrl = videoUrl
	return size, length, resolution, cdnUrl, nil
}

func DecodeVideoInfos(data string) (int32, int32, string, error) {
	size := int32(0)
	length := int32(0)
	resolution := ""
	vInfos := strings.Split(data, "_")
	if len(vInfos) != 3 {
		return size, length, resolution, errors.New("input params not in the correct format data=" + data)
	}

	vsize, err := strconv.ParseInt(vInfos[0], 10, 32)
	if err != nil {
		return size, length, resolution, errors.New("Convert size data error:" + err.Error())
	}
	size = int32(vsize)
	length = int32(GetVideoLength(vInfos[1]))
	resolution = strings.TrimSpace(vInfos[2])
	return size, length, resolution, nil
}

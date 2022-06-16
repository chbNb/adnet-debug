package mvutil

import (
	"errors"
	"net/url"
	"strconv"
	"strings"
)

type VideoCacheData struct {
	VideoSize       int32
	VideoLen        int32
	VideoResolution string
	VideoMTGCdnURL  string
}

func GetVideoLength(lenth string) (int, error) {
	result := 0
	lenArr := strings.Split(lenth, ":")
	if len(lenArr) != 3 {
		return result, errors.New("input params not in the correct format length=" + lenth)
	}
	hour, err := strconv.Atoi(lenArr[0])
	if err != nil {
		return result, errors.New("convert hour data error:" + err.Error())
	}
	result += hour * 3600
	min, err := strconv.Atoi(lenArr[1])
	if err != nil {
		return result, errors.New("convert min data error:" + err.Error())
	}
	result += min * 60
	secArr := strings.Split(lenArr[2], ".")
	if len(secArr) == 0 {
		return result, errors.New("convert sec data error not in correct format data=" + lenArr[2])
	}
	sec, err := strconv.Atoi(secArr[0])
	if err != nil {
		return result, errors.New("convert sec data error:" + err.Error())
	}
	result += sec
	// if len(secArr) > 1 {
	// 	result += 1
	// }
	return result, nil
}

func VideoMetaDataWithoutUrl(cacheData string) (*VideoCacheData, error) {
	return videoMetaData(cacheData, false)
}

func VideoMetaDataWithUrl(cacheData string) (*VideoCacheData, error) {
	return videoMetaData(cacheData, true)
}

func videoMetaData(cacheData string, decodeCdnUrl bool) (*VideoCacheData, error) {
	videoData := &VideoCacheData{}
	vInfos := strings.Split(cacheData, "_")
	if len(vInfos) < 3 {
		return videoData, errors.New("input params not in the correct format data=" + cacheData)
	}

	vsize, err := strconv.ParseInt(vInfos[0], 10, 32)
	if err != nil {
		return videoData, errors.New("ConvertReq size data error:" + err.Error())
	}
	videoData.VideoSize = int32(vsize)
	vlen, err := GetVideoLength(vInfos[1])
	if err != nil {
		return videoData, errors.New("ConvertReq length data error:" + err.Error())
	}
	videoData.VideoLen = int32(vlen)
	videoData.VideoResolution = strings.TrimSpace(vInfos[2])
	if decodeCdnUrl && len(vInfos) > 3 {
		videoUrl, err := url.QueryUnescape(strings.TrimSpace(vInfos[3]))
		if err != nil {
			return videoData, errors.New("ConvertReq cdnUrl data error:" + err.Error())
		}
		videoData.VideoMTGCdnURL = videoUrl
	}

	return videoData, nil
}

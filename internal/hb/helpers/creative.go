package helpers

import (
	"sort"
	"strconv"
	"strings"

	"gitlab.mobvista.com/ADN/adnet/internal/hb/constant"
	"gitlab.mobvista.com/ADN/adnet/internal/protobuf"
	"gitlab.mobvista.com/adserver/recommend_protocols/go/ad_server"
)

func RenderGInfo(crType ad_server.CreativeType, cr *protobuf.Creative) string {
	gInfoList := []string{
		strconv.FormatInt(int64(crType), 10),
		strconv.FormatInt(cr.CreativeId, 10),
		cr.AdvCreativeId,
		cr.VideoResolution,
	}
	return strings.Join(gInfoList, ",")
}

func GetQueryRAdvID(cadtype int, advCrIDMap map[string]string) string {
	// 按素材优先级上报，不区分adtype
	// 素材优先级为video->playable->image->icon
	if val, ok := advCrIDMap["video"]; ok && val != "" && val != "0" {
		return val
	}

	if val, ok := advCrIDMap["playable"]; ok && val != "" && val != "0" {
		return val
	}

	if val, ok := advCrIDMap["image"]; ok && val != "" && val != "0" {
		return val
	}

	if val, ok := advCrIDMap["icon"]; ok && val != "" && val != "0" {
		return val
	}

	return "0"
}

func GetQueryRAdType(adType int, videoUrl string) int {
	switch adType {
	case constant.Native:
		if len(videoUrl) > 0 {
			return constant.CREATIVE_AD_TYPE_NATIVE_VIDEO
		}
		return constant.CREATIVE_AD_TYPE_NATIVE
	case constant.RewardVideo:
		return constant.CREATIVE_AD_TYPE_REWARDED_VIDEO
	case constant.InterstitialVideo:
		return constant.CREATIVE_AD_TYPE_INTERSTITIAL_VIDEO
	}
	return 0
}

func GetCreativeGroupID(gidList []int) string {
	if len(gidList) <= 0 {
		return ""
	}
	// sort
	sort.Ints(gidList)
	// join
	var gidStrList []string
	for _, v := range gidList {
		gidStrList = append(gidStrList, strconv.Itoa(v))
	}
	gidStr := strings.Join(gidStrList, ",")
	// md5
	return Md5(gidStr)
}

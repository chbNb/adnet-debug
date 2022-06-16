package mvutil

/*
#cgo LDFLAGS: -L./cgo/lib/  -ldecrypter
#cgo CFLAGS: -I./cgo/include/
#include "decrypter.h"
*/

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"hash/crc32"
	"io"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	jsoniter "github.com/json-iterator/go"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/constant"
	"gitlab.mobvista.com/mae/go-kit/aes"

	"gitlab.mobvista.com/ADN/adnet/internal/mvconst"
	"gitlab.mobvista.com/ADN/adnet/internal/watcher"
	smodel "gitlab.mobvista.com/ADN/structs/model"
	"gopkg.in/mgo.v2/bson"
)

const (
	ACTIVE  = 1
	PAUSED  = 2
	DELETED = 3
	PENDING = 4
)

const (
	ANDROIDPLATFORM = 1
	IOSPLATFORM     = 2
)

func SliceIndex(limit int, predicate func(i int) bool) int {
	for i := 0; i < limit; i++ {
		if predicate(i) {
			return i
		}
	}
	return -1
}

type MwCreative struct {
	TempID map[int32]int `json:"tempId"`
}

// 对字符串进行SHA1哈希
func Sha1(data string) string {
	t := sha1.New()
	io.WriteString(t, data)
	return fmt.Sprintf("%x", t.Sum(nil))
}

func Md5(data string) string {
	t := md5.New()
	io.WriteString(t, data)
	return fmt.Sprintf("%x", t.Sum(nil))
}

func Base64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func Base64WithURLEncoding(data []byte) string {
	return base64.URLEncoding.EncodeToString(data)
}

func NumFormat(v float64, n int) float64 {
	pow10_n := math.Pow10(n)
	return math.Trunc((v)*pow10_n) / pow10_n
}

func GetThirdPartRating() float64 {
	v := 4 + rand.Float64()
	return NumFormat(v, 1)
}

var encodeCharMapping []byte
var decodeCharMapping []byte

var urlsafeEncodeCharMapping []byte
var urlsafeDecodeCharMapping []byte

var mvEncodeCharMapping []byte
var mvDecodeCharMapping []byte

var mpEncodeCharMapping []byte
var mpDecodeCharMapping []byte

var aA0Reg *regexp.Regexp
var aA0AndDot *regexp.Regexp
var Digit *regexp.Regexp

// 有感知sdk 码表
var mpNewEncodeCharMapping []byte
var mpNewDecodeCharMapping []byte

// 有感知mobpower AD解析码表
var mpNewAdEncodeCharMapping []byte
var mpNewAdDecodeCharMapping []byte

// 自有id key
var sysIdPrivateKey = "a453e81eaa7698e70926ca6f67474976"[:16]
var sysIdKey = "4af59e374dd9662dcc1a1a1e44679dce"
var tkiKey = "ebmclXzZOhtU2sRlZxGL8A"

func init() {
	chars := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	charsto := "vSoajc7dRzpWifGyNxZnV5k+DHLYhJ46lt0U3QrgEuq8sw/XMeBAT2Fb9P1OIKmC"
	encodeCharMapping = make([]byte, 256)
	decodeCharMapping = make([]byte, 256)
	for i := 0; i < 256; i++ {
		encodeCharMapping[i] = byte(i)
		decodeCharMapping[i] = byte(i)
	}
	for i := range chars {
		encodeCharMapping[int(chars[i])] = charsto[i]
		decodeCharMapping[int(charsto[i])] = chars[i]
	}

	urlsafeChars := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	urlsafeCharsto := "vSoajc7dRzpWifGyNxZnV5k-DHLYhJ46lt0U3QrgEuq8sw_XMeBAT2Fb9P1OIKmC"
	urlsafeEncodeCharMapping = make([]byte, 256)
	urlsafeDecodeCharMapping = make([]byte, 256)
	for i := 0; i < 256; i++ {
		urlsafeEncodeCharMapping[i] = byte(i)
		urlsafeDecodeCharMapping[i] = byte(i)
	}
	for i := range chars {
		urlsafeEncodeCharMapping[int(urlsafeChars[i])] = urlsafeCharsto[i]
		urlsafeDecodeCharMapping[int(urlsafeCharsto[i])] = urlsafeChars[i]
	}

	mvChars := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	mvCharsto := "uVUoXc3pCnDFvb8lNJj9ZHEia7QYrfSmROkGKA0ehIdtzB64Mq2gP5syTL1wWx+/"
	mvEncodeCharMapping = make([]byte, 256)
	mvDecodeCharMapping = make([]byte, 256)
	for i := 0; i < 256; i++ {
		mvEncodeCharMapping[i] = byte(i)
		mvDecodeCharMapping[i] = byte(i)
	}
	for i := range mvChars {
		mvEncodeCharMapping[int(mvChars[i])] = mvCharsto[i]
		mvDecodeCharMapping[int(mvCharsto[i])] = mvChars[i]
	}

	mpChars := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	mpCharsto := "r0J+IvLNO92/5fjSqGT7R8x3BFkEumlpZVciPHAstC4UXa6QDw1gozYdMWnhbeyK"
	mpEncodeCharMapping = make([]byte, 256)
	mpDecodeCharMapping = make([]byte, 256)
	for i := 0; i < 256; i++ {
		mpEncodeCharMapping[i] = byte(i)
		mpDecodeCharMapping[i] = byte(i)
	}
	for i := range mpChars {
		mpEncodeCharMapping[int(mpChars[i])] = mpCharsto[i]
		mpDecodeCharMapping[int(mpCharsto[i])] = mpChars[i]
	}

	mpNewChars := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	mpNewCharsto := "BFkEr0JIvT7R8x3LNstC4UXa6QO92zYdMWnhbeyK5fjSqGumlpZVciPHADw1go+/"
	mpNewEncodeCharMapping = make([]byte, 256)
	mpNewDecodeCharMapping = make([]byte, 256)
	for i := 0; i < 256; i++ {
		mpNewEncodeCharMapping[i] = byte(i)
		mpNewDecodeCharMapping[i] = byte(i)
	}
	for i := range mpNewChars {
		mpNewEncodeCharMapping[int(mpNewChars[i])] = mpNewCharsto[i]
		mpNewDecodeCharMapping[int(mpNewCharsto[i])] = mpNewChars[i]
	}

	mpNewPowerChars := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	mpNewPowerCharsto := "6QO9BFkEr0JIvYdMT7R8x3LNstC4UXa2zWnhbeyK5fjSw1goqGumlpZVciPHAD+/"
	mpNewAdEncodeCharMapping = make([]byte, 256)
	mpNewAdDecodeCharMapping = make([]byte, 256)
	for i := 0; i < 256; i++ {
		mpNewAdEncodeCharMapping[i] = byte(i)
		mpNewAdDecodeCharMapping[i] = byte(i)
	}
	for i := range mpNewPowerChars {
		mpNewAdEncodeCharMapping[int(mpNewPowerChars[i])] = mpNewPowerCharsto[i]
		mpNewAdDecodeCharMapping[int(mpNewPowerCharsto[i])] = mpNewPowerChars[i]
	}

	aA0Reg = regexp.MustCompile(`\W`)
	aA0AndDot = regexp.MustCompile(`[^\w\.-]`)
	Digit, _ = regexp.Compile("^\\d+$")
}

func Base64Encode(data string) string {
	base64Data := Base64([]byte(data))
	result := make([]byte, len(base64Data))
	for i := range base64Data {
		result[i] = encodeCharMapping[int(base64Data[i])]
	}
	return string(result)
}

func Base64EncodeWithURLEncoding(data string) string {
	base64Data := Base64WithURLEncoding([]byte(data))
	result := make([]byte, len(base64Data))
	for i := range base64Data {
		result[i] = urlsafeEncodeCharMapping[int(base64Data[i])]
	}
	return string(result)
}

func Base64DecodeWithURLEncoding(data string) string {
	decodeData := []byte(data)
	result := make([]byte, len(decodeData))
	for i := range decodeData {
		result[i] = urlsafeDecodeCharMapping[int(decodeData[i])]
	}
	ret, err := base64.URLEncoding.DecodeString(string(result))
	if err != nil {
		return string(ret)
	}
	return string(ret)
}

func Base64Decode(data string) string {
	decodeData := []byte(data)
	result := make([]byte, len(decodeData))
	for i := range decodeData {
		result[i] = decodeCharMapping[int(decodeData[i])]
	}
	ret, err := base64.StdEncoding.DecodeString(string(result))
	if err != nil {
		return string(ret)
	}
	return string(ret)
}

func Base64EncodeMV(data string) string {
	base64Data := Base64([]byte(data))
	result := make([]byte, len(base64Data))
	for i := range base64Data {
		result[i] = mvEncodeCharMapping[int(base64Data[i])]
	}
	return string(result)
}

func Base64DecodeMV(data string) string {
	decodeData := []byte(data)
	result := make([]byte, len(decodeData))
	for i := range decodeData {
		result[i] = mvDecodeCharMapping[int(decodeData[i])]
	}
	ret, err := base64.StdEncoding.DecodeString(string(result))
	if err != nil {
		return string(ret)
	}
	return string(ret)
}

func Base64EncodeMP(data string) string {
	base64Data := Base64([]byte(data))
	result := make([]byte, len(base64Data))
	for i := range base64Data {
		result[i] = mpEncodeCharMapping[int(base64Data[i])]
	}
	return string(result)
}

func Base64DecodeMP(data string) string {
	decodeData := []byte(data)
	result := make([]byte, len(decodeData))
	for i := range decodeData {
		result[i] = mpDecodeCharMapping[int(decodeData[i])]
	}
	ret, err := base64.StdEncoding.DecodeString(string(result))
	if err != nil {
		return string(ret)
	}
	return string(ret)
}

func Base64EncodeMPNew(data string) string {
	base64Data := Base64([]byte(data))
	result := make([]byte, len(base64Data))
	for i := range base64Data {
		result[i] = mpNewEncodeCharMapping[int(base64Data[i])]
	}
	return string(result)
}

func Base64DecodeMPNew(data string) string {
	decodeData := []byte(data)
	result := make([]byte, len(decodeData))
	for i := range decodeData {
		result[i] = mpNewDecodeCharMapping[int(decodeData[i])]
	}
	ret, err := base64.StdEncoding.DecodeString(string(result))
	if err != nil {
		return string(ret)
	}
	return string(ret)
}

func Base64DecodeMPNewAd(data string) string {
	decodeData := []byte(data)
	result := make([]byte, len(decodeData))
	for i := range decodeData {
		result[i] = mpNewAdDecodeCharMapping[int(decodeData[i])]
	}
	ret, err := base64.StdEncoding.DecodeString(string(result))
	if err != nil {
		return string(ret)
	}
	return string(ret)
}

func OriBase64Decode(data string) string {
	res, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return data
	}
	return string(res)
}

func DeBase64(src []byte) ([]byte, error) {
	return base64.StdEncoding.DecodeString(string(src))
}

func EnBase64(src string) string {
	return base64.StdEncoding.EncodeToString([]byte(src))
}

// func base64url_encode(data []byte) string {
// 	ret := base64.StdEncoding.EncodeToString(data)
// 	return strings.Map(func(r rune) rune {
// 		switch r {
// 		case '+':
// 			return '-'
// 		case '/':
// 			return '_'
// 		}

// 		return r
// 	}, ret)
// }

// todo-fix
// func base64url_decode(s string) ([]byte, error) {
// 	base64Str := strings.Map(func(r rune) rune {
// 		switch r {
// 		case '-':
// 			return '+'
// 		case '_':
// 			return '/'
// 		}

// 		return r
// 	}, s)

// 	if pad := len(base64Str) % 4; pad > 0 {
// 		base64Str += strings.Repeat("=", 4-pad)
// 	}

// 	return base64.StdEncoding.DecodeString(base64Str)
// }

// 判断是否是字符串还是数字
// func IsNum(str string) bool {
// 	isNum, _ := regexp.MatchString("(^[0-9]+$)", str)
// 	return isNum
// }

// func GetStrFromNum(num int64) string {
// 	return strconv.FormatInt(num, 10)
// }

func RequestHeader(req *http.Request, key string) string {
	if values, ok := req.Header[key]; ok && len(values) > 0 {
		return values[0]
	}
	return ""
}

func GetRequestID() string {
	return bson.NewObjectId().Hex()
	// return primitive.NewObjectID().Hex()
}

func GetGoTkClickIDNew(reqId string) string {
	return reqId[:len(reqId)-1] + "x"
}

// go tracking clickid
func GetGoTkClickID() string {
	requestID := GetRequestID()
	return SubString(requestID, 0, 18) + SubString(requestID, 19, 5) + "x"
}

func GetSSPlatformClickID() string {
	requestID := GetRequestID()
	return SubString(requestID, 0, 18) + SubString(requestID, 19, 5) + "y"
}

func GetGoTkCNClickID() string {
	requestID := GetRequestID()
	return SubString(requestID, 0, 18) + SubString(requestID, 19, 5) + "v"
}

// func GetSSAdjustPostbackClickID() string {
// 	requestID := GetRequestID()
// 	return SubString(requestID, 0, 18) + SubString(requestID, 19, 5) + "v"
// }

// Gen32BitsIDFromRequestID generate a 32 length uuid from requestid(mvutil.GetRequestID)
// a stupid implementation to avoid saving session id
// @ requestId: generated by mvutil.GetRequestID
// return a id of length 32
// func Gen32BitsIDFromRequestID(requestID string) string {
// 	if len(requestID) != 24 {
// 		return ""
// 	}
// 	return fmt.Sprintf("%s%s", "916176be", requestID)
// }

// func ClientIP(req *http.Request) string {
// 	forwardedByClientIP := true
// 	if forwardedByClientIP {
// 		clientIP := strings.TrimSpace(RequestHeader(req, "X-Real-Ip"))
// 		if len(clientIP) > 0 {
// 			return clientIP
// 		}
// 		clientIP = RequestHeader(req, "X-Forwarded-For")
// 		if index := strings.IndexByte(clientIP, ','); index >= 0 {
// 			clientIP = clientIP[0:index]
// 		}
// 		clientIP = strings.TrimSpace(clientIP)
// 		if len(clientIP) > 0 {
// 			return clientIP
// 		}
// 	}
// 	if ip, _, err := net.SplitHostPort(strings.TrimSpace(req.RemoteAddr)); err == nil {
// 		return ip
// 	}
// 	return ""
// }

// func GetServerIP(url string) (string, error) {
// 	resp, err := http.Get(url)
// 	if err != nil {
// 		return "", err
// 	}
// 	defer resp.Body.Close()
// 	content, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		return "", err
// 	}
// 	return string(content), nil
// }

// VerCampare 版本比较
func VerCampare(ver1 string, ver2 string) (int, error) {
	IntVer1, err := IntVer(ver1)
	if err != nil {
		return -2, err
	}
	IntVer2, err := IntVer(ver2)
	if err != nil {
		return -2, err
	}
	if IntVer1 == IntVer2 {
		return 0, nil
	}
	if IntVer1 > IntVer2 {
		return 1, nil
	} else {
		return -1, nil
	}
}

// IntVer把字符串版本转换成一个数字版本
func IntVer(osVersion string) (int32, error) {
	nums := strings.Split(osVersion, ".")
	osValue := 0
	for i := 0; i < 4; i++ {
		num := 0
		if len(nums) > i {
			var err error
			num, err = strconv.Atoi(nums[i])
			if err != nil {
				return 0, err
			}
		}
		osValue = osValue*100 + num
	}
	return int32(osValue), nil
}

// IntVer 把字符串版本转化成一个数字版本
// func IntVerOld(v string) (int32, error) {
// 	sections := strings.Split(v, ".")
// 	intVerSection := func(v string, n int) string {
// 		if n < len(sections) {
// 			return fmt.Sprintf("%02s", sections[n])
// 		} else {
// 			return "00"
// 		}
// 	}
// 	s := ""
// 	for i := 0; i < 4; i++ {
// 		s += intVerSection(v, i)
// 	}

// 	i64, err := strconv.ParseInt(s, 10, 32)
// 	if err != nil {
// 		return 0, err
// 	}
// 	return int32(i64), nil
// }

func InStrArray(val string, array []string) bool {
	for _, v := range array {
		if val == v {
			return true
		}
	}
	return false
}

func InArray(val int, array []int) bool {
	for _, v := range array {
		if val == v {
			return true
		}
	}
	return false
}

func InInt64Arr(val int64, array []int64) bool {
	for _, v := range array {
		if val == v {
			return true
		}
	}
	return false
}

func InInt32Arr(val int32, array []int32) bool {
	for _, v := range array {
		if val == v {
			return true
		}
	}
	return false
}

func InInt8Array(val int8, array []int8) bool {
	for _, v := range array {
		if val == v {
			return true
		}
	}
	return false
}

// func GetInternalIP() (string, error) {
// 	addrs, err := net.InterfaceAddrs()
// 	if err != nil {
// 		return "", err
// 	}
// 	for _, a := range addrs {
// 		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
// 			if ipnet.IP.To4() != nil {
// 				return ipnet.IP.String(), nil
// 			}
// 		}
// 	}
// 	return "", errors.New("have get on ip info")
// }

// todo-fix
// func GetEndCardUrl(unitid int, sdk_version string) (string, error) {
// 	if len(Rewardvideo_end_screen.Value.Http.Rewardvideo_end_screen) == 0 {
// 		return "", errors.New("offerwall_urls.http.rewardvideo_end_screen is empty")
// 	}

// 	if strings.Contains(Rewardvideo_end_screen.Value.Http.Rewardvideo_end_screen, "?") {
// 		return fmt.Sprintf("%s&unit_id=%d&sdk_version=%s",
// 			Rewardvideo_end_screen.Value.Http.Rewardvideo_end_screen,
// 			unitid, sdk_version), nil
// 	} else {
// 		return fmt.Sprintf("%s?unit_id=%d&sdk_version=%s",
// 			Rewardvideo_end_screen.Value.Http.Rewardvideo_end_screen,
// 			unitid, sdk_version), nil
// 	}
// }

// func GetRvImagesKey(height, width string) string {
// 	if (height == "250") && (width == "300") {
// 		return "play_all_300x250"
// 	} else if (height == "500") && (width == "500") {
// 		return "click_all_500x500"
// 	} else if (height == "800") && (width == "600") {
// 		return "endcard_portrait_600x800"
// 	} else if (height == "600") && (width == "800") {
// 		return "endcard_landscape_800x600"
// 	}
// 	return ""
// }

// return screen width and height
// func FormatScreenSize(requestId, screenSize string) (int, int) {
// 	if strings.Contains(screenSize, "x") {
// 		sizeList := strings.Split(screenSize, "x")
// 		var result []int64
// 		for _, data := range sizeList {
// 			f64, err := strconv.ParseFloat(strings.TrimSpace(data), 64)
// 			if err != nil {
// 				i64, err := strconv.ParseInt(strings.TrimSpace(data), 10, 64)
// 				if err != nil {
// 					// Logger.Runtime.Debugf("request_id=[%s] formatScreenSize error:%s", requestId, err.Error())
// 					break
// 				}
// 				result = append(result, i64)
// 				continue
// 			}
// 			i64 := int64(f64)
// 			result = append(result, i64)
// 		}
// 		if len(result) == 2 {
// 			// screenSize = fmt.Sprintf("%dx%d", result[0], result[1])
// 			return int(result[0]), int(result[1])
// 		}
// 	}
// 	return 0, 0
// }

func GetUrlScheme(reqType int) string {
	if reqType == 2 {
		return "https://"
	}
	return "http://"
}

func GetAdTypeStr(adType int32) string {
	switch adType {
	case mvconst.ADTypeUnknown:
		return "unknown"
	case mvconst.ADTypeText:
		return "text"
	case mvconst.ADTypeBanner:
		return "banner"
	case mvconst.ADTypeAppwall:
		return "appwall"
	case mvconst.ADTypeOverlay:
		return "overlay"
	case mvconst.ADTypeFullScreen:
		return "full_screen"
	case mvconst.ADTypeInterstitial:
		return "interstitial"
	case mvconst.ADTypeNative:
		return "native"
	case mvconst.ADTypeRewardVideo:
		return "rewarded_video"
	case mvconst.ADTypeFeedsVideo:
		return "feeds_video"
	case mvconst.ADTypeOfferWall:
		return "offerwall"
	case mvconst.ADTypeInterstitialSdk:
		return "interstitial_sdk"
	case mvconst.ADTypeOnlineVideo:
		return "online_video"
	case mvconst.ADTypeJSNativeVideo:
		return "js_native_video"
	case mvconst.ADTypeJSBannerVideo:
		return "js_banner_video"
	case mvconst.ADTypeInterstitialVideo:
		return "interstitial_video"
	case mvconst.ADTypeInteractive:
		return "interactive"
	case mvconst.ADTypeJMIcon:
		return "jm_icon"
	case mvconst.ADTypeWXNative:
		return "wx_native"
	case mvconst.ADTypeWXAppwall:
		return "wx_appwall"
	case mvconst.ADTypeWXBanner:
		return "wx_banner"
	case mvconst.ADTypeWXRewardImg:
		return "wx_reward_image"
	case mvconst.ADTypeMoreOffer:
		return "more_offer"
	case mvconst.ADTypeSdkBanner:
		return "sdk_banner"
	case mvconst.ADTypeSplash:
		return "splash"
	case mvconst.ADTypeNativeH5:
		return "native_h5"
	default:
		return ""
	}
}

func SerializeMP(params *Params, mwCreative string) string {
	return Base64Encode(strings.Join([]string{
		params.ExtDataInit.ExpIds,
		params.RequestID,
		strconv.FormatInt(params.PublisherID, 10),
		strconv.FormatInt(params.AppID, 10),
		strconv.FormatInt(params.UnitID, 10),
		params.MWadBackend,
		params.MWadBackendData,
		params.CountryCode,
		strconv.FormatInt(params.CityCode, 10),
		strconv.Itoa(params.Platform),
		strconv.FormatInt(int64(params.AdType), 10),
		params.OSVersion,
		params.SDKVersion,
		params.AppVersionName,
		params.Brand,
		params.Model,
		params.ScreenSize,
		strconv.Itoa(params.Orientation),
		params.Language,
		strconv.Itoa(params.NetworkType),
		"",
		params.ClientIP,
		"",
		params.ServerIP,
		params.IMEI,
		params.MAC,
		params.AndroidID,
		params.GAID,
		params.IDFA,
		strconv.Itoa(params.MWFlowTagID),
		strconv.FormatInt(int64(params.AdNum), 10),
		strconv.Itoa(params.TNum),
		strconv.Itoa(params.MWRandValue),
		params.MWbackendConfig,
		params.MWplayInfo,
		params.Scenario,
		"",
		params.Extra,
		mwCreative,
		params.DspExt,
		strconv.FormatInt(params.DPrice, 10),
		"", // req_backend
		"", // reject_code
		"", // imp_timeout
		params.PriceFactor,
		params.Extchannel,
		params.ThirdPartyABTestStr,
		strconv.Itoa(params.RequestType),
		params.ExtPlacementId,
		strconv.Itoa(params.Open),
		params.OnlineApiBidPrice,
		params.StartModeTagsStr,
	}, "|"))
}

// func SerializeMPOld(params *Params, mwCreative string) string {
// 	var buf bytes.Buffer
// 	buf.WriteString(params.RequestID + "|")
// 	buf.WriteString(strconv.FormatInt(params.PublisherID, 10) + "|")
// 	buf.WriteString(strconv.FormatInt(params.AppID, 10) + "|")
// 	buf.WriteString(strconv.FormatInt(params.UnitID, 10) + "|")
// 	buf.WriteString(params.MWadBackend + "|")
// 	buf.WriteString(params.MWadBackendData + "|")
// 	buf.WriteString(params.CountryCode + "|")
// 	buf.WriteString(strconv.FormatInt(params.CityCode, 10) + "|")
// 	buf.WriteString(strconv.Itoa(params.Platform) + "|")
// 	buf.WriteString(strconv.FormatInt(int64(params.AdType), 10) + "|")
// 	buf.WriteString(params.OSVersion + "|")
// 	buf.WriteString(params.SDKVersion + "|")
// 	buf.WriteString(params.AppVersionName + "|")
// 	buf.WriteString(params.Brand + "|")
// 	buf.WriteString(params.Model + "|")
// 	buf.WriteString(params.ScreenSize + "|")
// 	buf.WriteString(strconv.Itoa(params.Orientation) + "|")
// 	buf.WriteString(params.Language + "|")
// 	buf.WriteString(strconv.Itoa(params.NetworkType) + "|")
// 	buf.WriteString("|")
// 	buf.WriteString(params.ClientIP + "|")
// 	buf.WriteString("|")
// 	buf.WriteString(params.ServerIP + "|")
// 	buf.WriteString(params.IMEI + "|")
// 	buf.WriteString(params.MAC + "|")
// 	buf.WriteString(params.AndroidID + "|")
// 	buf.WriteString(params.GAID + "|")
// 	buf.WriteString(params.IDFA + "|")
// 	buf.WriteString(strconv.Itoa(params.MWFlowTagID) + "|")
// 	buf.WriteString(strconv.FormatInt(int64(params.AdNum), 10) + "|")
// 	buf.WriteString(strconv.Itoa(params.TNum) + "|")
// 	buf.WriteString(strconv.Itoa(params.MWRandValue) + "|")
// 	buf.WriteString(params.MWbackendConfig + "|")
// 	buf.WriteString(params.MWplayInfo + "|")
// 	buf.WriteString(params.Scenario + "||")
// 	buf.WriteString(params.Extra + "|")
// 	buf.WriteString(mwCreative)
// 	return Base64Encode(buf.String())
// }

func SerializeMwCreative(mwCreative []int32, requestId string) string {
	if len(mwCreative) == 0 {
		return "0"
	}
	creativeMap := make(map[int32]int)
	for _, tempId := range mwCreative {
		if _, ok := creativeMap[tempId]; !ok {
			creativeMap[tempId] = 1
			continue
		}
		creativeMap[tempId] += 1
	}
	creative := MwCreative{TempID: creativeMap}
	creativeTmp, err := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(creative)
	if err != nil {
		Logger.Runtime.Warnf("request_id=[%s] marshal creativeMap in write midway request error:%s", requestId, err.Error())
		return "0"
	}
	return string(creativeTmp)
}

// 字段梳理
func SerializeMPPart(reqParam *RequestParams, adBackend, adBackendData, backendConfig, dspExt, fakePrice, priceFactor string) string {
	return strings.Join([]string{
		reqParam.Param.RequestID,
		strconv.FormatInt(reqParam.Param.PublisherID, 10),
		strconv.FormatInt(reqParam.Param.AppID, 10),
		strconv.FormatInt(reqParam.Param.UnitID, 10),
		adBackend,
		adBackendData,
		reqParam.Param.CountryCode,
		strconv.FormatInt(reqParam.Param.CityCode, 10),
		strconv.Itoa(reqParam.Param.Platform),
		strconv.FormatInt(int64(reqParam.Param.AdType), 10),
		reqParam.Param.OSVersion,
		reqParam.Param.SDKVersion,
		reqParam.Param.AppVersionName,
		reqParam.Param.Brand,
		reqParam.Param.Model,
		reqParam.Param.ScreenSize,
		strconv.Itoa(reqParam.Param.Orientation),
		reqParam.Param.Language,
		strconv.Itoa(reqParam.Param.NetworkType),
		reqParam.Param.MCC + reqParam.Param.MNC,
		"0|0|0",
		reqParam.Param.IMEI,
		reqParam.Param.MAC,
		reqParam.Param.AndroidID,
		reqParam.Param.GAID,
		reqParam.Param.IDFA,
		strconv.Itoa(reqParam.FlowTagID),
		strconv.FormatInt(int64(reqParam.Param.AdNum), 10),
		strconv.Itoa(reqParam.Param.TNum),
		strconv.Itoa(reqParam.RandValue),
		backendConfig,
		reqParam.Param.MWplayInfo,
		reqParam.Param.Scenario,
		"||0",
		dspExt,
		fakePrice,
		"",
		"",
		"",
		priceFactor,
		reqParam.Param.Extchannel, "", strconv.Itoa(reqParam.Param.RequestType),
	}, "|")
}

// func SerializeMPPartOld(reqParam *RequestParams, adBackend, adBackendData, backendConfig string) string {
// 	var buf bytes.Buffer
// 	buf.WriteString(reqParam.Param.RequestID + "|")
// 	buf.WriteString(strconv.FormatInt(reqParam.Param.PublisherID, 10) + "|")
// 	buf.WriteString(strconv.FormatInt(reqParam.Param.AppID, 10) + "|")
// 	buf.WriteString(strconv.FormatInt(reqParam.Param.UnitID, 10) + "|")
// 	buf.WriteString(adBackend + "|")
// 	buf.WriteString(adBackendData + "|")
// 	buf.WriteString(reqParam.Param.CountryCode + "|")
// 	buf.WriteString(strconv.FormatInt(reqParam.Param.CityCode, 10) + "|")
// 	buf.WriteString(strconv.Itoa(reqParam.Param.Platform) + "|")
// 	buf.WriteString(strconv.FormatInt(int64(reqParam.Param.AdType), 10) + "|")
// 	buf.WriteString(reqParam.Param.OSVersion + "|")
// 	buf.WriteString(reqParam.Param.SDKVersion + "|")
// 	buf.WriteString(reqParam.Param.AppVersionName + "|")
// 	buf.WriteString(reqParam.Param.Brand + "|")
// 	buf.WriteString(reqParam.Param.Model + "|")
// 	buf.WriteString(reqParam.Param.ScreenSize + "|")
// 	buf.WriteString(strconv.Itoa(reqParam.Param.Orientation) + "|")
// 	buf.WriteString(reqParam.Param.Language + "|")
// 	buf.WriteString(strconv.Itoa(reqParam.Param.NetworkType) + "|")
// 	buf.WriteString(reqParam.Param.MCC + reqParam.Param.MNC + "|")
// 	buf.WriteString("0|0|0|")
// 	buf.WriteString(reqParam.Param.IMEI + "|")
// 	buf.WriteString(reqParam.Param.MAC + "|")
// 	buf.WriteString(reqParam.Param.GAID + "|")
// 	buf.WriteString(reqParam.Param.AndroidID + "|")
// 	buf.WriteString(reqParam.Param.IDFA + "|")
// 	buf.WriteString(strconv.Itoa(reqParam.FlowTagID) + "|")
// 	buf.WriteString(strconv.FormatInt(int64(reqParam.Param.AdNum), 10) + "|")
// 	buf.WriteString(strconv.Itoa(reqParam.Param.TNum) + "|")
// 	buf.WriteString(strconv.Itoa(reqParam.RandValue) + "|")
// 	buf.WriteString(backendConfig + "|")
// 	buf.WriteString(reqParam.Param.MWplayInfo + "|")
// 	buf.WriteString(reqParam.Param.Scenario + "|||0") //增加第三方模版统计数据
// 	return buf.String()
// }

func SerializeP(params *Params) string {
	return Base64([]byte(strings.Join([]string{
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		GetAdTypeStr(params.AdType),
		params.ImageSize,
		"",
		params.PlatformName,
		params.OSVersion,
		params.SDKVersion,
		params.Model,
		params.ScreenSize,
		strconv.Itoa(params.Orientation),
		"",
		params.Language,
		params.NetworkTypeName,
		params.MCC + params.MNC,
		"",
		"",
		params.Extra3,
		params.Extra4,
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		params.ClientIP,
		params.IMEI,
		params.MAC,
		params.AndroidID,
		"",
		"",
		"",
		"",
		params.GAID,
		params.IDFA,
		"",
		params.Brand,
		params.RemoteIP,
		params.SessionID,
		params.ParentSessionID,
		"",
		"",
		"",
		"",
		strconv.Itoa(params.TNum),
		"",
		params.IDFV + "," + params.OpenIDFA,
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		params.Extstats,
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		params.Extstats2,
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		params.ReplaceBrand,
		params.ReplaceModel,
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		params.ExtSysId,
		strconv.FormatFloat(params.ApiVersion, 'f', 1, 64),
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		params.OAID,
		"",
		"",
		"",
		"",
		params.ExtBigTemId,
		"",
		"",
		params.ExtPlacementId,
		"",
		strconv.FormatFloat(params.RespFillEcpmFloor, 'f', 6, 64),
		"",
		"",
		"",
		params.StartModeTagsStr,
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		params.ExtDeviceId,
		"",
		"",
		"",
		"",
		"",
		"",
		params.JunoCommonLogInfoJson,
	}, "|")))
}

// func SerializePOld(params *Params, campaign *CampaignInfo) string {
// 	var buf bytes.Buffer
// 	buf.WriteString("|||||||" + GetAdTypeStr(params.AdType) + "|")
// 	buf.WriteString(params.ImageSize + "||")
// 	buf.WriteString(params.PlatformName + "|")
// 	buf.WriteString(params.OSVersion + "|")
// 	buf.WriteString(params.SDKVersion + "|")
// 	buf.WriteString(params.Model + "|")
// 	buf.WriteString(params.ScreenSize + "|")
// 	buf.WriteString(strconv.Itoa(params.Orientation) + "||")
// 	buf.WriteString(params.Language + "|")
// 	buf.WriteString(params.NetworkTypeName + "|")
// 	buf.WriteString(params.MCC + params.MNC + "|||")
// 	buf.WriteString(params.Extra3 + "|")
// 	buf.WriteString(params.Extra4 + "||||||||")
// 	buf.WriteString(params.ClientIP + "|")
// 	buf.WriteString(params.IMEI + "|")
// 	buf.WriteString(params.MAC + "|")
// 	buf.WriteString(params.AndroidID + "|||||")
// 	buf.WriteString(params.GAID + "|")
// 	buf.WriteString(params.IDFA + "||")
// 	buf.WriteString(params.Brand + "|")
// 	buf.WriteString(params.RemoteIP + "|")
// 	buf.WriteString(params.SessionID + "|")
// 	buf.WriteString(params.ParentSessionID + "|||||")
// 	buf.WriteString(strconv.Itoa(params.TNum) + "||")
// 	buf.WriteString(params.IDFV + "," + params.OpenIDFA + "||||||||")
// 	buf.WriteString(params.Extstats + "||||||||||||||||||||||||||||||||||||")
// 	return Base64(buf.Bytes())
// }

func SerializeQ(params *Params, campaign *smodel.CampaignInfo) string {
	oriPrice := params.PriceIn
	if oriPrice == 0 && campaign.OriPrice != 0 {
		oriPrice = campaign.OriPrice
	}

	price := params.PriceOut
	if price == 0 && campaign.Price != 0 {
		price = campaign.Price
	}

	if params.UseAlgoPrice {
		if params.AlgoPriceIn > 0 {
			oriPrice = params.AlgoPriceIn
		}

		if params.AlgoPriceOut > 0 {
			price = params.AlgoPriceOut
		}
	}

	return "a_" + Base64Encode(strings.Join([]string{
		"2.0",
		params.Domain,
		params.RequestID,
		params.Extra5,
		strconv.FormatInt(params.PublisherID, 10),
		strconv.FormatInt(params.AppID, 10),
		strconv.FormatInt(params.UnitID, 10),
		strconv.FormatInt(int64(params.AdvertiserID), 10),
		strconv.FormatInt(params.CampaignID, 10),
		FormatFloat64(oriPrice),
		FormatFloat64(price),
		params.Extra,
		strconv.Itoa(params.RequestType),
		params.CountryCode,
		strconv.Itoa(params.Extra8),
		params.Extra10,
		strconv.FormatInt(int64(params.Extra13), 10),
		strconv.FormatInt(params.Extra14, 10),
		strconv.FormatInt(params.CityCode, 10),
		params.AppVersionName,
		params.Scenario,
		strconv.Itoa(params.Extra16), "|-",
		strconv.Itoa(params.Extra7),
		strconv.Itoa(params.Extbtclass),
		strconv.FormatInt(params.Extfinalsubid, 10),
		strconv.Itoa(params.ExtdeleteDevid),
		strconv.Itoa(params.ExtinstallFrom),
		strconv.Itoa(params.ExtnativeVideo),
		strconv.Itoa(params.ExtflowTagId), "",
		params.Extendcard, "",
		strconv.FormatInt(params.ExtdspRealAppid, 10),
		params.ExtfinalPackageName,
		strconv.Itoa(params.Extnativex),
		strconv.FormatInt(params.CreativeId, 10),
		strconv.Itoa(params.Extattr), "",
		strconv.FormatInt(int64(params.Extctype), 10),
		strconv.Itoa(params.Extrvtemplate), "",
		strconv.Itoa(params.Extplayable),
		strconv.Itoa(params.Extb2t),
		params.Extchannel,
		strconv.Itoa(params.Extabtest2),
		params.Extbp,
		strconv.FormatInt(int64(params.Extsource), 10),
		strconv.Itoa(params.ExtappearUa),
		strconv.Itoa(params.ExtCDNAbTest),
		params.ExtMpNormalMap,
		strconv.Itoa(params.Extplayable2),
		params.ReqType,
		params.ExtData2,
		strconv.FormatFloat(params.BidPrice, 'f', 2, 64),
		params.ABTestTagStr,
		strconv.FormatInt(params.ExtMtgId, 10),
		params.ExtSlotId,
		params.ExtCreativeNew,
		"",
		"",
		"",
		params.AsABTestResTag,
		params.OnlineApiBidPrice,
		params.SkadnetworkDataStr,
		strconv.Itoa(params.LinkType),
	}, "|"))
}

func SerializeCSP(params *Params, dspExt string) string {
	return Base64Encode(strings.Join([]string{
		params.MWadBackend,
		params.MWadBackendData,
		strconv.Itoa(params.MWFlowTagID),
		strconv.Itoa(params.MWRandValue),
		params.MWbackendConfig,
		strconv.FormatInt(int64(params.AdNum), 10),
		strconv.Itoa(params.TNum),
		dspExt,
	}, "|"))
}

func SerializeC(params *Params) string {
	return RawUrlEncode(strings.Join([]string{
		strconv.Itoa(int(params.FormatAdType)),
		strconv.Itoa(params.Platform),
		params.OSVersion,
		params.SDKVersion,
		params.Model,
		params.ScreenSize,
		strconv.Itoa(params.Orientation),
		params.Language,
		strconv.Itoa(params.NetworkType),
		params.MCCMNC,
		params.IMEI,
		params.MAC,
		params.AndroidID,
		params.GAID,
		params.IDFA,
		params.Brand,
		params.IDFVOpenIDFA,
	}, "|"))

}

// func SerializeCSPOld(params *Params) string {
// 	var buf bytes.Buffer
// 	buf.WriteString(params.MWadBackend + "|")
// 	buf.WriteString(params.MWadBackendData + "|")
// 	buf.WriteString(strconv.Itoa(params.MWFlowTagID) + "|")
// 	buf.WriteString(strconv.Itoa(params.MWRandValue) + "|")
// 	buf.WriteString(params.MWbackendConfig + "|")
// 	buf.WriteString(strconv.FormatInt(int64(params.AdNum), 10) + "|")
// 	buf.WriteString(strconv.Itoa(params.TNum))
// 	return Base64Encode(buf.String())
// }

// func SerializeCSPOldV2(params *Params) string {
// 	return Base64Encode(fmt.Sprintf("%s|%s|%d|%d|%s|%d|%d",
// 		params.MWadBackend, params.MWadBackendData, params.MWFlowTagID,
// 		params.MWRandValue, params.MWbackendConfig, params.AdNum, params.TNum))
// }

func CheckParam(param string) string {
	return strings.Replace(param, "|", "", -1)
}

func HasInvisibleChar(str string) bool {
	f := func(c rune) bool {
		return c < 32 || c > 126
	}

	return strings.IndexFunc(str, f) != -1
}

func TrimBlank(data string) string {
	data = strings.TrimSpace(data)
	data = strings.Replace(data, "\t", "", -1)
	data = strings.Replace(data, "\n", "", -1)
	return data
}

func TrimOnlyaA0New(str string) string {
	return aA0Reg.ReplaceAllString(str, "")
}

func TrimOnlyaA0(str string) string {
	reg := regexp.MustCompile(`\W`)
	return reg.ReplaceAllString(str, "")
}

func TrimAa0AndDotNew(str string) string {
	return aA0AndDot.ReplaceAllString(str, "")
}

func TrimAa0AndDot(str string) string {
	reg := regexp.MustCompile(`[^\w\.-]`)
	return reg.ReplaceAllString(str, "")
}

func CleanH5Data(data string) string {
	data = strings.TrimSpace(data)
	data = strings.Replace(data, "\t", "", -1)
	data = strings.Replace(data, "\n", "", -1)
	data = strings.Replace(data, "|", "", -1)
	return data
}

func RandByRate(rateMap map[int]int) int {
	sum := 0
	for _, v := range rateMap {
		sum = sum + v
	}
	if sum <= 0 {
		return 0
	}
	rand := rand.Intn(sum)
	i := 0
	for k, v := range rateMap {
		i = i + v
		if rand < i {
			return k
		}
	}
	return 0
}

// IRate IRate
type IRate interface {
	GetRate() int
}

// RandStringByRateInt RandStringByRateInt
func RandStringByRateInt(rateMap map[string]IRate) string {
	var sum int
	ks := make([]string, 0, len(rateMap))
	for k, v := range rateMap {
		sum = sum + v.GetRate()
		ks = append(ks, k)
	}
	if sum <= 0 {
		return ""
	}

	sort.Sort(sort.StringSlice(ks))
	rand := int(rand.Intn(sum))
	var i int
	for _, k := range ks {
		v := rateMap[k]
		i = i + v.GetRate()
		if rand < i {
			return k
		}
	}
	return ""
}

// RandStringByRateIntCustonRand RandStringByRateIntCustonRand
func RandStringByRateIntCustonRand(rateMap map[string]IRate, randf func() int) string {
	var sum int
	ks := make([]string, 0, len(rateMap))
	for k, v := range rateMap {
		sum = sum + v.GetRate()
		ks = append(ks, k)
	}
	if sum <= 0 {
		return ""
	}

	sort.Sort(sort.StringSlice(ks))
	rand := randf() % sum
	var i int
	for _, k := range ks {
		v := rateMap[k]
		i = i + v.GetRate()
		if rand < i {
			return k
		}
	}
	return ""
}

func RandByRate2(rateMap map[string]int) string {
	sum := 0
	for _, v := range rateMap {
		sum = sum + v
	}
	if sum <= 0 {
		return ""
	}
	rand := rand.Intn(sum)
	i := 0
	for k, v := range rateMap {
		i = i + v
		if rand < i {
			return k
		}
	}
	return ""
}

func RandByDeviceRateWithIntMap(rateMap map[int]int, params *Params) int {
	sum := 0
	for _, v := range rateMap {
		sum = sum + v
	}
	if sum <= 0 {
		return 0
	}
	rand := 0
	if !IsDevidEmpty(params) {
		rand = GetRandConsiderZero(params.GAID, params.IDFA, mvconst.SALT_CDN_ABTEST, sum)
	} else if HasDevidInAndroidCN(params) {
		rand = GetRandConsiderZeroWithAndroidIdAndImei(params.AndroidID, params.IMEI, mvconst.SALT_CDN_ABTEST, sum)
	} else {
		// 空设备的情况
		return 0
	}
	val, ok := RandByRateInMapInInt(rateMap, rand)
	if ok {
		return val
	}
	return 0
}

func RandByDeviceRate(rateMap map[string]int, params *Params) string {
	sum := 0
	for _, v := range rateMap {
		sum = sum + v
	}
	if sum <= 0 {
		return ""
	}
	rand := 0
	if !IsDevidEmpty(params) {
		rand = GetRandConsiderZero(params.GAID, params.IDFA, mvconst.SALT_CDN_ABTEST, sum)
	} else if HasDevidInAndroidCN(params) {
		rand = GetRandConsiderZeroWithAndroidIdAndImei(params.AndroidID, params.IMEI, mvconst.SALT_CDN_ABTEST, sum)
	} else {
		// 空设备的情况
		return ""
	}
	val, ok := RandByRateInMapInString(rateMap, rand)
	if ok {
		return val
	}
	return ""
}

func RandByDeviceWithRateMap(rateMap map[string]int, params *Params, salt string) string {
	sum := 0
	for _, v := range rateMap {
		sum = sum + v
	}
	if sum <= 0 {
		return ""
	}
	randVal := GetRandByGlobalTagId(params, salt, sum)
	val, ok := RandByRateInMapInString(rateMap, randVal)
	if ok {
		return val
	}
	return ""
}

// 使用Int排序
func RandByRate3WithIntMap(rateMap map[int]int) int {
	sum := 0
	for _, v := range rateMap {
		sum = sum + v
	}
	if sum <= 0 {
		return 0
	}
	rand := rand.Intn(sum)
	val, ok := RandByRateInMapInInt(rateMap, rand)
	if ok {
		return val
	}
	return 0
}

// 使用string排序
func RandByRate3(rateMap map[string]int) string {
	sum := 0
	for _, v := range rateMap {
		sum = sum + v
	}
	if sum <= 0 {
		return ""
	}
	rand := rand.Intn(sum)
	val, ok := RandByRateInMapInString(rateMap, rand)
	if ok {
		return val
	}
	return ""
}

func RandByMappingIdfa(rateMap map[string]int, idfa string) string {
	sum := 0
	for _, v := range rateMap {
		sum = sum + v
	}
	if sum <= 0 {
		return ""
	}
	randVal := int(crc32.ChecksumIEEE([]byte(mvconst.SALT_MAPPING_IDFA+"_"+idfa))) % sum
	val, ok := RandByRateInMapInString(rateMap, randVal)
	if ok {
		return val
	}
	return ""
}

func GetRandConsiderZeroWithAndroidIdAndImei(androidId string, imei string, salt string, randSum int) int {
	str := androidId + imei
	if len(str) <= 0 {
		return -1
	}
	str = str + salt
	return int(crc32.ChecksumIEEE([]byte(str))) % randSum
}

func HasDevidInAndroidCN(params *Params) bool {
	if params.Platform == mvconst.PlatformAndroid && params.CountryCode == "CN" && (len(params.AndroidID) > 0 || len(params.IMEI) > 0) {
		return true
	}
	return false
}

func RandIntArr(arr []int) int {
	llen := len(arr)
	if llen <= 0 {
		return 0
	}
	rand := rand.Intn(llen)
	return arr[rand]
}

func RandInt64Arr(arr []int64) int64 {
	llen := len(arr)
	if llen <= 0 {
		return int64(0)
	}
	rand := rand.Intn(llen)
	return arr[rand]
}

func RandStringArr(arr []string) string {
	llen := len(arr)
	if llen <= 0 {
		return ""
	}
	rand := rand.Intn(llen)
	return arr[rand]
}

func GetRandByStr(str string, randSum int) int {
	return int(crc32.ChecksumIEEE([]byte(str))) % randSum
}

func GetRandByGlobalTagId(params *Params, salt string, randSum int) int {
	return int(crc32.ChecksumIEEE([]byte(salt+"_"+GetGlobalUniqDeviceTag(params)))) % randSum
}

func GetRandConsiderZero(gaid string, idfa string, salt string, randSum int) int {
	str := getDeviceString(gaid, idfa)
	if len(str) <= 0 {
		return -1
	}
	str = str + salt
	return int(crc32.ChecksumIEEE([]byte(str))) % randSum
}

func GetRandByStrAddSalt(str string, randSum int, salt string) int {
	str = str + salt
	return int(crc32.ChecksumIEEE([]byte(str))) % randSum
}

// GetPureRand 返回一个 0<=n< randSum 的随机数
func GetPureRand(randSum int) int {
	if randSum == 0 {
		return 0
	}
	return rand.Intn(randSum) // 使用完全随机代替原来的按设备随机
}

func getDeviceString(gaid string, idfa string) string {
	if len(idfa) <= 0 || idfa == "00000000-0000-0000-0000-000000000000" || idfa == "idfa" || idfa == "-" {
		idfa = ""
	}
	if len(gaid) <= 0 || gaid == "gaid" || gaid == "00000000-0000-0000-0000-000000000000" || gaid == "-" {
		gaid = ""
	}
	return idfa + gaid
}

func GetIdfaString(idfa string) string {
	if len(idfa) <= 0 || idfa == "00000000-0000-0000-0000-000000000000" || idfa == "idfa" || idfa == "-" {
		idfa = ""
	}
	return idfa
}

func RandArr(randArr map[int]int, gaid string, idfa string, salt string) int {
	sum := 0
	for _, v := range randArr {
		sum = sum + v
	}
	randInt := GetRandConsiderZero(gaid, idfa, salt, sum)
	i := 0
	for k, v := range randArr {
		i = i + v
		if randInt < i {
			return k
		}
	}
	return 0
}

func IsDevidEmpty(param *Params) bool {
	str := getDeviceString(param.GAID, param.IDFA)
	// if len(str) > 0 {
	// 	return false
	// }
	// return true
	return len(str) <= 0
}

func Max(x int, y int) int {
	if x > y {
		return x
	}
	return y
}

func SubString(str string, begin, length int) string {
	rs := []rune(str)
	llen := len(rs)
	if begin < 0 {
		begin = 0
	}
	if begin >= llen {
		begin = llen
	}
	end := begin + length
	if end > llen {
		end = llen
	}
	return string(rs[begin:end])
}

func GetVersionCode(version string) int32 {
	if len(version) <= 0 {
		return int32(0)
	}
	arr := strings.Split(version, ".")
	if len(arr) < 3 {
		return int32(0)
	}
	code := int32(0)
	for i := 0; i < 3; i++ {
		tmpInt32, _ := strconv.ParseInt(arr[i], 10, 32)
		code = code*100 + int32(tmpInt32)
	}
	return code
}

// net/url 会对query进行按key排序
func HttpBuildQuery(queryMap map[string]string) string {
	urlV := url.Values{}
	i := 0
	for k, v := range queryMap {
		if i == 0 {
			urlV.Set(k, v)
		} else {
			urlV.Add(k, v)
		}
		i++
	}
	return urlV.Encode()
}

func UrlEncode(value string) string {
	return url.QueryEscape(value)
}

func RawUrlEncode(value string) string {
	value = UrlEncode(value)
	return strings.Replace(value, "+", "%20", -1)
}

func RawUrlDecode(value string) string {
	value = strings.Replace(value, "%20", "+", -1)
	return UrlDecode(value)
}

func UrlDecode(value string) string {
	res, err := url.QueryUnescape(value)
	if err != nil {
		return value
	}
	return res
}

func SubUtf8Str(str string, n int) string {
	if utf8.RuneCountInString(str) <= n {
		return str
	}
	runeStr := []rune(str)
	str = string(runeStr[0:n]) + "......"
	return str
}

func mappingServerLog(params *Params) string {
	ctime := params.RequestTime
	return strings.Join([]string{
		time.Unix(ctime, 0).Format("20060102"),
		time.Unix(ctime, 0).Format("150405"),
		strconv.FormatInt(ctime, 10),
		params.SysId,
		params.IDFV,
		params.BkupId,
		params.ExtModel,
		params.OsvUpTime,
		params.PlatformName,
		params.MappingServerFrom,
		strconv.FormatInt(params.AppID, 10),
		params.ClientIP,
		params.UserAgent,
		params.UpTime,
		params.RequestID,
		params.MappingServerMessage,
		params.RuId,
		params.MappingServerResCode,
		params.MappingServerDebugInfo,
		strconv.FormatInt(params.CityCode, 10),
		params.EncryptedRuid,
	}, "\t")
}

func statLog(params *Params) string {
	filterTag := ""
	if params.IsBlockByImpCap {
		filterTag = mvconst.FilterRequestByImpBlock
	} else if params.IsLowFlowUnitReq {
		filterTag = mvconst.FilterRequestByLowFlow
	} else {
		filterTag = params.FilterRequestReason
	}
	ctime := params.RequestTime
	return strings.Join([]string{
		time.Unix(ctime, 0).Format("20060102"),
		time.Unix(ctime, 0).Format("150405"),
		strconv.FormatInt(ctime, 10),
		strconv.FormatInt(params.PublisherID, 10),
		strconv.FormatInt(params.AppID, 10),
		strconv.FormatInt(params.UnitID, 10),
		"0",
		"0",
		"0",
		params.Scenario,
		GetAdTypeStr(params.AdType),
		params.ImageSize,
		strconv.Itoa(params.RequestType),
		params.PlatformName,
		params.OSVersion,
		params.SDKVersion,
		params.Model,
		params.ScreenSize,
		strconv.Itoa(params.Orientation),
		params.CountryCode,
		params.Language,
		params.NetworkTypeName,
		params.MCC + params.MNC,
		params.Extra,
		params.Extra2,
		params.Extra3,
		params.Extra4,
		params.Extra5,
		strconv.Itoa(params.Extra6),
		strconv.Itoa(params.Extra7),
		strconv.Itoa(params.Extra8),
		params.Extra9,
		params.Extra10,
		params.RequestID,
		params.ClientIP,
		params.IMEI,
		params.MAC,
		params.AndroidID,
		params.ServerIP,
		"0",
		"0",
		"0",
		params.GAID,
		params.IDFA,
		params.AppVersionName,
		params.Brand,
		params.RemoteIP,
		params.SessionID,
		params.ParentSessionID,
		params.Extra11,
		strconv.FormatInt(params.CityCode, 10),
		strconv.FormatInt(int64(params.Extra13), 10),
		"0",
		strconv.Itoa(params.Extra15),
		strconv.Itoa(params.Extra16),
		params.IDFV + "," + params.OpenIDFA,
		"0",
		"0",
		params.Extra20,
		"0",
		"0",
		"0",
		"0",
		"0",
		params.ExtpackageName,
		"0",
		strconv.Itoa(params.ExtflowTagId),
		params.ExtcdnType,
		params.Extendcard,
		params.ChannelInfo,
		"0",
		params.ExtfinalPackageName,
		"0",
		"0",
		"0",
		"0",
		"0",
		strconv.Itoa(params.Extrvtemplate),
		"0",
		"0",
		"0",
		params.ExtPlayableArr,
		"0",
		"0",
		strconv.Itoa(params.Extb2t),
		params.Extchannel,
		"0",
		"0",
		"0",
		"0",
		"0",
		params.Extreject,
		params.Extalgo,
		params.ExtthirdCid,
		strconv.FormatInt(int64(params.ExtifLowerImp), 10),
		"0",
		"0",
		params.ExtsystemUseragent,
		"0",
		"0",
		"0",
		"0",
		"0",
		params.ExtMpNormalMap,
		"0",
		"0",
		params.ReplaceBrand,
		params.ReplaceModel,
		params.ExtData,
		"0",
		params.ExtCreativeNew,
		params.ExtTagArrStr,
		strconv.Itoa(params.Extplayable2),
		params.ReqType,
		"0",
		"0",
		"0",
		params.ExtSysId,
		params.ExtApiVersion,
		"0",
		"0",
		"0",
		"0",
		strconv.FormatFloat(params.BidFloor, 'f', 2, 64),
		params.BidsPrice,
		params.DspExt,
		strconv.Itoa(params.ExtReduceFillReq),
		strconv.Itoa(params.ExtReduceFillResp),
		"0",
		params.ExtReduceFillValue,
		"0",
		"0",
		params.MofData,
		params.ABTestTagStr,
		params.OAID,
		params.IMSI,
		strconv.FormatBool(params.WaterfallFallback),
		"0",
		params.ExtBigTplOfferData,
		params.ExtBigTemId,
		"0",
		"0",
		params.ExtPlacementId,
		params.ExtAdxAlgo,
		strconv.FormatFloat(params.RespFillEcpmFloor, 'f', 6, 64),
		strings.Join(params.ReqBackends, ","),
		strings.Join(params.BackendReject, ","),
		params.ReduceFill,
		params.StartModeTagsStr,
		params.PolarisTplData,
		params.PolarisCreativeData,
		params.TplGrayTag,
		filterTag,
		"0",
		"0",
		"0",
		"0",
		params.ExtDeviceId,
		params.AsABTestResTag,
		strconv.Itoa(params.Open),
		"0",
		"0",
		"0",
		"0",
		params.JunoCommonLogInfoJson,
		"0",
		"0",
		params.ParentId,
		params.PioneerExtdataInfo,
		params.PioneerOfferExtdataInfo,
	}, "\t")
}

func statHBLog(in *ReqCtx) string {
	params := in.ReqParams.Param
	filterTag := ""
	if params.IsBlockByImpCap {
		filterTag = mvconst.FilterRequestByImpBlock
	} else if params.IsLowFlowUnitReq {
		filterTag = mvconst.FilterRequestByLowFlow
	} else {
		filterTag = params.FilterRequestReason
	}
	ctime := params.RequestTime
	return strings.Join([]string{
		time.Unix(ctime, 0).Format("20060102"),
		time.Unix(ctime, 0).Format("150405"),
		strconv.FormatInt(ctime, 10),
		strconv.FormatInt(params.PublisherID, 10),
		strconv.FormatInt(params.AppID, 10),
		strconv.FormatInt(params.UnitID, 10),
		"0",
		"0",
		"0",
		params.Scenario,
		GetAdTypeStr(params.AdType),
		params.ImageSize,
		strconv.Itoa(params.RequestType),
		params.PlatformName,
		params.OSVersion,
		params.SDKVersion,
		params.Model,
		params.ScreenSize,
		strconv.Itoa(params.Orientation),
		params.CountryCode,
		params.Language,
		params.NetworkTypeName,
		params.MCC + params.MNC,
		params.Extra,
		params.Extra2,
		params.Extra3,
		params.Extra4,
		params.Extra5,
		strconv.Itoa(params.Extra6),
		strconv.Itoa(params.Extra7),
		strconv.Itoa(params.Extra8),
		RawUrlEncode(params.UserAgent),
		params.Extra10,
		params.RequestID,
		params.ClientIP,
		params.IMEI,
		params.MAC,
		params.AndroidID,
		params.ServerIP,
		"0",
		"0",
		"0",
		params.GAID,
		params.IDFA,
		params.AppVersionName,
		params.Brand,
		params.RemoteIP,
		params.SessionID,
		params.ParentSessionID,
		params.Extra11,
		strconv.FormatInt(params.CityCode, 10),
		strconv.FormatInt(int64(params.Extra13), 10),
		"0",
		strconv.Itoa(params.Extra15),
		strconv.Itoa(params.Extra16),
		params.IDFV + "," + params.OpenIDFA,
		"0",
		"0",
		params.Extra20,
		"0",
		"0",
		"0",
		"0",
		"0",
		params.ExtpackageName,
		"0",
		strconv.Itoa(in.FlowTagID),
		params.ExtcdnType,
		params.Extendcard,
		params.ChannelInfo,
		"0",
		params.ExtfinalPackageName,
		"0",
		"0",
		"0",
		"0",
		"0",
		strconv.Itoa(params.Extrvtemplate),
		"0",
		"0",
		"0",
		params.ExtPlayableArr,
		"0",
		"0",
		strconv.Itoa(params.Extb2t),
		params.Extchannel,
		"0",
		"0",
		"0",
		"0",
		"0",
		params.Extreject,
		params.Extalgo,
		params.ExtthirdCid,
		strconv.FormatInt(int64(params.ExtifLowerImp), 10),
		"0",
		"0",
		params.ExtsystemUseragent,
		"0",
		"0",
		"0",
		"0",
		"0",
		params.ExtMpNormalMap,
		"0",
		"0",
		params.ReplaceBrand,
		params.ReplaceModel,
		params.ExtData,
		"0",
		params.ExtCreativeNew,
		params.ExtTagArrStr,
		strconv.Itoa(params.Extplayable2),
		params.ReqType,
		"0",
		"0",
		"0",
		params.ExtSysId,
		params.ExtApiVersion,
		"0",
		"0",
		"0",
		"0",
		strconv.FormatFloat(params.BidFloor, 'f', 2, 64),
		params.BidsPrice,
		params.DspExt,
		strconv.Itoa(params.ExtReduceFillReq),
		strconv.Itoa(params.ExtReduceFillResp),
		"0",
		params.ExtReduceFillValue,
		"0",
		"0",
		params.MofData,
		params.ABTestTagStr,
		params.OAID,
		params.IMSI,
		strconv.FormatBool(params.WaterfallFallback),
		"0",
		params.ExtBigTplOfferData,
		params.ExtBigTemId,
		"0",
		"0",
		params.ExtPlacementId,
		params.ExtAdxAlgo,
		strconv.FormatFloat(params.RespFillEcpmFloor, 'f', 6, 64),
		strings.Join(params.ReqBackends, ","),
		strings.Join(params.BackendReject, ","),
		params.ReduceFill,
		params.StartModeTagsStr,
		params.PolarisTplData,
		params.PolarisCreativeData,
		params.TplGrayTag,
		filterTag,
		"0",
		"0",
		"0",
		"0",
		params.ExtDeviceId,
		params.AsABTestResTag,
		strconv.Itoa(params.Open),
		"0",
		"0",
		"0",
		"0",
		params.JunoCommonLogInfoJson,
		"0",
		"0",
		params.ParentId,
		params.PioneerExtdataInfo,
		params.PioneerOfferExtdataInfo,
	}, "\t")
}

func StatMappingServerLog(params *Params) {
	Logger.MappingServerLog.Info(mappingServerLog(params))
}

func StatLossRequestLog(params *Params) {
	params.LossReqFlag = true
	Logger.LossRequest.Info(statLog(params))
}

func StatRequestLog(params *Params) {
	Logger.Request.Info(statLog(params))
}

func StatHBRequestLog(in *ReqCtx) {
	Logger.Request.Info(statHBLog(in))
}

// p参数中的主要字段
func SerializeOImpP(params *Params) string {
	return Base64([]byte(strings.Join([]string{
		strconv.FormatInt(params.PublisherID, 10),
		strconv.FormatInt(params.AppID, 10),
		strconv.FormatInt(params.UnitID, 10),
		strconv.FormatInt(int64(params.AdvertiserID), 10),
		strconv.FormatInt(params.CampaignID, 10),
		"",
		params.Scenario,
		GetAdTypeStr(params.AdType),
		params.ImageSize,
		strconv.Itoa(params.RequestType),
		params.PlatformName,
		params.OSVersion,
		params.SDKVersion,
		params.Model,
		params.ScreenSize,
		strconv.Itoa(params.Orientation),
		params.CountryCode,
		params.Language,
		params.NetworkTypeName,
		params.MCC + params.MNC,
		params.Extra,
		"",
		params.Extra3,
		params.Extra4,
		params.Extra5,
		"",
		strconv.Itoa(params.Extra7),
		strconv.Itoa(params.Extra8),
		params.Extra9,
		params.Extra10,
		params.RequestID,
		params.ClientIP,
		params.IMEI,
		params.MAC,
		params.AndroidID,
		params.ServerIP,
		"",
		"",
		"",
		params.GAID,
		params.IDFA,
		params.AppVersionName,
		params.Brand,
		params.RemoteIP,
		params.SessionID,
		params.ParentSessionID,
		"",
		strconv.FormatInt(params.CityCode, 10),
		strconv.FormatInt(int64(params.Extra13), 10),
		strconv.FormatInt(params.Extra14, 10),
		strconv.Itoa(params.Extra15),
		strconv.Itoa(params.Extra16),
		params.IDFV + "," + params.OpenIDFA,
		"",
		"",
		params.Extra20,
		"",
		strconv.FormatInt(params.Extfinalsubid, 10),
		"",
		"",
		"",
		params.ExtpackageName,
		"",
		strconv.Itoa(params.ExtflowTagId),
		"",
		params.Extendcard,
		strconv.Itoa(params.ExtrushNoPre),
		"",
		params.ExtfinalPackageName,
		strconv.Itoa(params.Extnativex),
		"",
		"",
		"",
		strconv.FormatInt(int64(params.Extctype), 10),
		strconv.Itoa(params.Extrvtemplate),
		strconv.Itoa(params.Extabtest1),
		"",
		"",
		"",
		"",
		"",
		strconv.Itoa(params.Extb2t),
		params.Extchannel,
		"",
		"",
		"",
		params.Extbp,
		strconv.FormatInt(int64(params.Extsource), 10),
		"",
		params.Extalgo,
		params.ExtthirdCid,
		strconv.FormatInt(int64(params.ExtifLowerImp), 10),
		"",
		"",
		params.ExtsystemUseragent,
		"",
		"",
		"",
		"",
		"",
		params.ExtMpNormalMap,
		"",
		"",
		params.ExtBrand,
		params.ExtModel,
		params.ExtData,
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		params.ExtSysId,
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		params.OAID,
		"",
		"",
		"",
		params.ExtBigTplOfferData,
		params.ExtBigTemId,
		"",
		"",
		params.ExtPlacementId,
		"",
		"",
		"",
		"",
		"",
		"",
		params.PolarisTplData,
		params.PolarisCreativeData,
		params.TplGrayTag,
		"",
		"",
		"",
		"",
		"",
		params.ExtDeviceId,
	}, "|")))
}

func GetAndroidOSVersion(versionCode int) string {
	version, ok := mvconst.AndroidVersionMap[versionCode]
	if ok {
		return version
	}
	return ""
}

func GetPlatformStr(platform string) int {
	platform = strings.ToLower(platform)
	switch platform {
	case mvconst.PlatformNameAndroid:
		return 1
	case mvconst.PlatformNameIOS:
		return 2
	case mvconst.PlatformNameOther:
		return 0
	}
	return 0
}

func GetMPSdkVersionCompare(sdkVersion string) bool {
	sdkVersionArr := strings.Split(sdkVersion, "_")
	if sdkVersionArr[0] == "mp" {
		sdkVersionCode := strings.Split(sdkVersionArr[1], ".")
		firstCodeInt, err := strconv.Atoi(sdkVersionCode[0])
		if err != nil {
			return false
		}
		if firstCodeInt < 4 {
			return true
		} else {
			return false
		}
	}
	return false
}

func UrlDecodeReplace(data string) string {
	data = strings.Replace(data, "%3D", "=", -1)
	data = strings.Replace(data, "%2B", "+", -1)
	data = strings.Replace(data, "%2F", "/", -1)
	return data
}

func IntGetData(num int, defValue int) *int {
	if &num == nil {
		return &defValue
	}
	return &num
}

func RenderImageSizeIdByAdType(adType string, fsTpl string, ori int) int {
	switch adType {
	case "full_screen":
		if InStrArray(fsTpl, []string{"stamp", "flow"}) {
			return mvconst.IMAGE_SIZE_ID_300X250
		} else if InStrArray(fsTpl, []string{"flyin", "webview"}) {
			return mvconst.IMAGE_SIZE_ID_480X320
		} else {
			return mvconst.IMAGE_SIZE_ID_UNKOWN
		}
	case "overlay":
		return mvconst.IMAGE_SIZE_ID_300X250
	case "appwall":
		if ori == 2 {
			return mvconst.IMAGE_SIZE_ID_480X320
		} else {
			return mvconst.IMAGE_SIZE_ID_320X480
		}
	case "banner":
		return mvconst.IMAGE_SIZE_ID_320X50
	}
	return mvconst.IMAGE_SIZE_ID_UNKOWN
}

// GetDeviceTag GetDeviceTag
func GetDeviceTag(params *Params) string {
	if params.Platform == mvconst.PlatformIOS {
		idfa := params.IDFA
		if !(len(idfa) <= 0 || idfa == "00000000-0000-0000-0000-000000000000" || idfa == "idfa" || idfa == "-") {
			return idfa
		}
	}

	if params.Platform == mvconst.PlatformAndroid {
		gaid := params.GAID
		if !(len(gaid) <= 0 || gaid == "gaid" || gaid == "00000000-0000-0000-0000-000000000000" || gaid == "-") {
			return gaid
		}

		if len(params.AndroidID) > 0 && params.AndroidID != "0" {
			return params.AndroidID
		}

		if len(params.IMEI) > 0 && params.IMEI != "0" {
			return params.IMEI
		}
	}

	if len(params.SysId) > 0 {
		return params.SysId
	}

	if len(params.BkupId) > 0 {
		return params.BkupId
	}

	return Md5(params.ClientIP + params.Model)
}

func FormatFloat64(f float64) string {
	str := strconv.FormatFloat(f, 'f', 6, 64)
	seg := len(str)
	for index := len(str) - 1; index > 0; index-- {
		if str[index] == '0' {
			seg = index
			continue
		} else if str[index] == '.' {
			seg = index
		}
		break
	}
	return str[:seg]
}

func IsIvOrRv(adtype int32) bool {
	if adtype == mvconst.ADTypeRewardVideo || adtype == mvconst.ADTypeInterstitialVideo {
		return true
	}
	return false
}

func GetGaidOrAndroid(gaid string, android string) string {
	if len(gaid) == 0 {
		return android
	}
	return gaid
}

func RandByRateInMapV2(valInMap map[string]IRate, getRate func(sum int) int) string {
	if len(valInMap) == 0 {
		return ""
	}

	var keyList []string
	var sum int
	for k, v := range valInMap {
		keyList = append(keyList, k)
		sum += v.GetRate()
	}

	rate := getRate(sum)
	sort.Strings(keyList)
	i := 0
	for _, key := range keyList {
		val := valInMap[key]
		i += val.GetRate()
		if rate < i {
			return key
		}
	}
	return ""
}

func RandByRateInMap(valInMap map[string]int, randVal int) (int, bool) {
	valArr := make(map[int]int, len(valInMap))
	var keyList []int
	for k, v := range valInMap {
		cctVal, err := strconv.Atoi(k)
		if err != nil {
			continue
		}
		valArr[cctVal] = v
		keyList = append(keyList, cctVal)
	}
	if len(keyList) == 0 {
		return 0, false
	}
	sort.Ints(keyList)
	i := 0
	for _, Val := range keyList {
		i = i + valArr[Val]
		if randVal < i {
			return Val, true
		}
	}
	return 0, false
}

func RandByRateInMapInString(valInMap map[string]int, randVal int) (string, bool) {
	valArr := make(map[string]int, len(valInMap))
	var keyList []string
	for k, v := range valInMap {
		valArr[k] = v
		keyList = append(keyList, k)
	}
	if len(keyList) == 0 {
		return "", false
	}
	sort.Strings(keyList)
	i := 0
	for _, Val := range keyList {
		i = i + valArr[Val]
		if randVal < i {
			return Val, true
		}
	}
	return "", false
}

func RandByRateInMapInInt(valInMap map[int]int, randVal int) (int, bool) {
	valArr := make(map[int]int, len(valInMap))
	var keyList []int
	for k, v := range valInMap {
		valArr[k] = v
		keyList = append(keyList, k)
	}
	if len(keyList) == 0 {
		return 0, false
	}
	sort.Ints(keyList)
	i := 0
	for _, Val := range keyList {
		i = i + valArr[Val]
		if randVal < i {
			return Val, true
		}
	}
	return 0, false
}

func NeedNewJssdkDomain(path string, flag int) bool {
	if path == mvconst.PATHJssdkApi && flag == 1 {
		return true
	}
	return false
}

// func GetIconUrl(param *Params, defaultUrl string) string {
// 	defaultIconUrl := fmt.Sprintf("%s%s", GetUrlScheme(int(param.HTTPReq)),
// 		defaultUrl)
// 	// jssdk new domain
// 	if NeedNewJssdkDomain(param.RequestPath, param.Ndm) && len(param.JssdkCdnDomain) > 0 {
// 		u, err := url.Parse(defaultIconUrl)
// 		if err == nil {
// 			u.Host = param.JssdkCdnDomain
// 			return u.String()
// 		}
// 	}
// 	return defaultIconUrl
// }

func IsMpad(path string) bool {
	return path == mvconst.PATHMPAD ||
		path == mvconst.PATHMPNewAD || path == mvconst.PATHMPADV2
}

func Int32Join(list []int32, delim string) string {
	var str string
	var i int = 0
	for _, val := range list {
		if i != 0 {
			str += delim
		}
		str += strconv.FormatInt(int64(val), 10)
		i += 1
	}

	return str
}

func StatReduceFillLog(params *Params) {
	// ctime := time.Now().Unix()
	ctime := params.RequestTime
	Logger.ReduceFill.Info(strings.Join([]string{
		time.Unix(ctime, 0).Format("2006-01-02 15:04:05"),
		params.RequestID,
		strconv.FormatInt(params.UnitID, 10),
		params.CountryCode,
		strings.Join(params.ReqBackends, "|"),
		Int32Join(params.FillBackendId, "|"),
		strings.Join(params.ReduceFillList, "|"),
	}, "\t"))
}

func IsAppwallOrMoreOffer(adType int32) bool {
	return adType == mvconst.ADTypeAppwall || adType == mvconst.ADTypeMoreOffer
}

func IsMoreOfferAndAppwallRequestPioneer(adType int32, tag string) bool {
	return IsAppwallOrMoreOffer(adType) && tag == mvconst.REQUEST_PIONEER
}

func IsRequestPioneerDirectly(params *Params) bool {
	return (IsAppwallOrMoreOffer(params.AdType) && params.ExtDataInit.MoreofferAndAppwallMvToPioneerTag == mvconst.REQUEST_PIONEER) ||
		params.ExtDataInit.MpToPioneerTag == mvconst.REQUEST_PIONEER
}

func RenderStringToIntList(listStr string) []int64 {
	var intList []int64
	var eList []int64
	err := jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal([]byte(listStr), &eList)
	if err == nil && len(eList) > 0 {
		for _, v := range eList {
			intList = append(intList, v)
		}
		return intList
	}
	var eStrList []string
	err = jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal([]byte(listStr), &eStrList)
	if err == nil && len(eStrList) > 0 {
		for _, v := range eStrList {
			vInt, err2 := strconv.ParseInt(v, 10, 64)
			if err2 == nil {
				intList = append(intList, vInt)
			}
		}
		return intList
	}
	return intList
}

func SetFieldValueForABTest(vfield reflect.Value, value int) bool {
	if !vfield.CanSet() {
		return false
	}

	switch vfield.Kind() {
	case reflect.Bool:
		vfield.SetBool(value == mvconst.ABTEST_TRUE)
		return true

	case reflect.String:
		if value == mvconst.ABTEST_FALSE {
			vfield.SetString("")
			return true
		}

	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int, reflect.Int64:
		vfield.SetInt(int64(value))
		return true
	}
	return false
}

func getParamStr(key string, queryMap RequestQueryMap) string {
	res, err := queryMap.GetString(key, true, "")
	if err != nil {
		return ""
	}
	return res
}

func CisIsDevidEmpty(params *Params) bool {
	if IsDevidEmpty(params) {
		// 对于中国流量，需判断Android和imei
		if params.Platform == mvconst.PlatformAndroid && params.CountryCode == "CN" && (len(params.AndroidID) > 0 || len(params.IMEI) > 0) {
			return false
		}
		return true
	}
	return false
}

// 是否有deviceid（ios有idfa，Android有gaid）
func HasDevid(params *Params) bool {
	if params.Platform == mvconst.PlatformIOS {
		idfa := params.IDFA
		if !(len(idfa) <= 0 || idfa == "00000000-0000-0000-0000-000000000000" || idfa == "idfa" || idfa == "-") {
			return true
		}
	}
	if params.Platform == mvconst.PlatformAndroid {
		gaid := params.GAID
		if !(len(gaid) <= 0 || gaid == "gaid" || gaid == "00000000-0000-0000-0000-000000000000" || gaid == "-") {
			return true
		}
	}
	return false
}

func GetGlobalDeviceTag(params *Params) string {
	// iOS：idfa>idfv>自有id；android：gaid>imei>oaid>android>自有id
	if params.Platform == mvconst.PlatformIOS {
		idfa := params.IDFA
		if !(len(idfa) <= 0 || idfa == "00000000-0000-0000-0000-000000000000" || idfa == "idfa" || idfa == "-") {
			return idfa
		}
		if len(params.IDFV) > 0 && params.IDFV != "0" {
			return params.IDFV
		}
	}

	if params.Platform == mvconst.PlatformAndroid {
		gaid := params.GAID
		if !(len(gaid) <= 0 || gaid == "gaid" || gaid == "00000000-0000-0000-0000-000000000000" || gaid == "-") {
			return gaid
		}

		if len(params.IMEI) > 0 && params.IMEI != "0" {
			return params.IMEI
		}

		if len(params.OAID) > 0 && params.OAID != "0" {
			return params.OAID
		}

		if len(params.AndroidID) > 0 && params.AndroidID != "0" {
			return params.AndroidID
		}

	}
	if len(params.SysId) > 0 {
		return params.SysId
	}

	if len(params.BkupId) > 0 {
		return params.BkupId
	}

	return ""
}

func GetGlobalUniqDeviceTag(params *Params) string {
	// iOS：idfa>idfv>自有id>ipua；android：gaid>imei>oaid>android>自有id>ipua
	tag := GetGlobalDeviceTag(params)
	if tag != "" {
		return tag
	}

	return Md5(params.ClientIP + params.Model)
}

func GetParamValInString(v reflect.Value) string {
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.String:
		return v.String()
	case reflect.Bool:
		return strconv.FormatBool(v.Bool())
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(v.Float(), 'f', -1, 64)
	case reflect.Slice:
		var val []string
		for i := 0; i < v.Len(); i++ {
			val = append(val, v.Index(i).String())
		}
		return strings.Join(val, ",")
	default:
		return ""
	}
}

func IsBannerOrSplashOrDI(adType int32) bool {
	return adType == mvconst.ADTypeSdkBanner || adType == mvconst.ADTypeSplash || adType == mvconst.ADTypeInterstitialSdk
}

func IsBannerOrSdkBannerOrSplash(adType int32) bool {
	return IsBannerOrSplashOrDI(adType) || adType == mvconst.ADTypeBanner
}

func IsBannerOrSplashOrNativeH5(adType int32) bool {
	return IsBannerOrSplashOrDI(adType) || adType == mvconst.ADTypeNativeH5
}

func IsBannerOrSdkBannerOrSplashOrNativeH5(adType int32) bool {
	return IsBannerOrSplashOrNativeH5(adType) || adType == mvconst.ADTypeBanner
}

func IsNative(adType int32) bool {
	return adType == mvconst.ADTypeNativePic || adType == mvconst.ADTypeNativeVideo
}

func RecordUseTime(key string, startNanoTime, divisor int64) {
	watcher.AddAvgWatchValue(key, float64((time.Now().UnixNano()-startNanoTime)/divisor))
}

type CheckResult struct {
	Result  string      `json:"result"`
	OldData interface{} `json:"old_data"`
	NewData interface{} `json:"new_data"`
}

func IsHbOrV3OrV5Request(path string) bool {
	return path == mvconst.PATHOpenApiV3 ||
		path == mvconst.PATHOpenApiV5 ||
		path == mvconst.PATHBid ||
		path == mvconst.PATHMopubBid ||
		path == mvconst.PATHLoad ||
		path == mvconst.PATHTOPON ||
		path == mvconst.PATHOpenApiMoreOffer
}

func IsHbOrV3OrV5OrOnlineAPIRequest(path string, requestType int) bool {
	return IsHbOrV3OrV5Request(path) || requestType == mvconst.REQUEST_TYPE_OPENAPI_AD
}

func RemoveExtractorCollections(dbc []DbConfig, collections map[string]bool) []DbConfig {
	re := []DbConfig{}
	for _, v := range dbc {
		if !collections[v.Collection] {
			re = append(re, v)
		}

	}
	return re
}

func IsBigoRequest(publisherId int64, requestType int) bool {
	return publisherId == mvconst.PUB_BIGO && requestType == mvconst.REQUEST_TYPE_OPENAPI_AD
}

func GetDebugBidFloorAndBidPriceKey(params *Params) string {
	devId := GetDeviceTag(params)
	appIdStr := strconv.FormatInt(params.AppID, 10)
	return appIdStr + "_" + devId
}

// 原来setting 获取设备id用于生成sys_id的逻辑
func GetSysIdDevId(params *Params) string {
	if params.PlatformName == constant.PlatformIOS {
		if params.IDFA != "00000000-0000-0000-0000-000000000000" && params.IDFA != "" {
			return params.IDFA
		}
	} else {
		// 这里估计还需要排出掉00000000-0000-0000-0000-000000000000，setting 之前是没有排出掉这种情况的
		if params.GAID != "00000000-0000-0000-0000-000000000000" && params.GAID != "" {
			return params.GAID
		}
		if params.AndroidID != "" {
			return params.AndroidID
		}
		if params.IMEI != "" {
			return params.IMEI
		}
	}
	return ""
}

func GetOriginIdFromEncryptSysId(str string) (string, int64) {
	// 获取解密结果
	decryptStr, err := DevIdCbcDecrypt(str)
	// 解密错误，可能是sdk 升级后传来的明文的sysid，bkupid。也可能是其他原因。
	// 解密错误则返回原来的值
	if err != nil {
		Logger.Runtime.Errorf("str=[%s] decrypt id error. error:%s", str, err.Error())
		return str, 0
	}
	// 切割得到原始的id值
	decryptData := strings.SplitN(decryptStr, "_", 2)
	if len(decryptData) == 2 {
		timestamp, _ := strconv.ParseInt(decryptData[0], 10, 64)
		return decryptData[1], timestamp
	}
	return str, 0
}

func DevIdCbcDecrypt(str string) (val string, err error) {
	devIdEncrypt := aes.NewAESCBCEncrypt([]byte(sysIdKey), []byte(sysIdPrivateKey))
	decrypt, err := devIdEncrypt.CbcDecrypt(str)
	if err != nil {
		return "", err
	}
	return decrypt, nil
}

func DevIdCbcEncrypt(str string) (val string, err error) {
	devIdEncrypt := aes.NewAESCBCEncrypt([]byte(sysIdKey), []byte(sysIdPrivateKey))
	encrypt, err := devIdEncrypt.CbcEncrypt(str)
	if err != nil {
		return str, err
	}
	return encrypt, nil
}

func Decrypt(str string) (val string, err error) {
	keyByte := sha512.Sum384([]byte(tkiKey))
	keyStr := keyByte[:32]
	iv := keyByte[32:48]
	tkiEncrypt := aes.NewAESCBCEncrypt(keyStr, iv)
	// 从url中获取时会进行urldecode
	// urlencode后再获取
	str = url.QueryEscape(str)
	decrypt, err := tkiEncrypt.CbcPkCS7Decrypt(str)
	if err != nil {
		return str, err
	}
	return decrypt, nil
}

func StringAppendCategory(a []string, b []string) []string {
	check := make(map[string]int)
	d := append(a, b...)
	res := make([]string, 0)
	for _, val := range d {
		check[val] = 1
	}

	for letter, _ := range check {
		res = append(res, letter)
	}

	return res
}

func RenderIntslice(sliceStr string) []int32 {
	var eList []int32
	err := jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal([]byte(sliceStr), &eList)
	if err != nil {
		return []int32{}
	}
	return eList
}

func DelSubStr(str string, subStr string) string {
	strArr := strings.Split(str, subStr)
	return strings.Join(strArr, "")
}

func IsHBS2SUseVideoOrientation(r *RequestParams) bool {
	return r.IsHBRequest && len(r.Param.HBS2SBidID) > 0 &&
		(r.Param.AdType == mvconst.ADTypeRewardVideo || r.Param.AdType == mvconst.ADTypeInterstitialVideo)
}

func IsInterfaceNil(i interface{}) bool {
	vi := reflect.ValueOf(i)
	if vi.Kind() == reflect.Ptr {
		return vi.IsNil()
	}
	return false
}

func IsNewIv(params *Params) bool {
	return params.ApiVersion >= mvconst.API_VERSION_2_3 && params.AdType == mvconst.ADTypeInterstitialVideo
}

// 判断宽高尺寸是否与给定的比例相等
func IsEqualProportion(width, height, numerator, denominator int) bool {
	if width == 0 || height == 0 || numerator == 0 || denominator == 0 {
		return false
	}

	// 判断w/h是否与n/d相等
	width = width * denominator
	numerator = numerator * height
	if width == numerator {
		return true
	}

	return false
}

package helpers

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	"math"
	"unicode/utf8"
)

var encodeCharMapping []byte
var decodeCharMapping []byte

var chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
var charsto = "vSoajc7dRzpWifGyNxZnV5k+DHLYhJ46lt0U3QrgEuq8sw/XMeBAT2Fb9P1OIKmC"

type DeviceInfo struct {
	Imei      string  `json:"imei"`
	Mac       string  `json:"mac"`
	AndroidID string  `json:"android_id"`
	OAID      string  `json:"oaid"`
	Dmt       string  `json:"dmt"`
	Dmf       float64 `json:"dmf"`
	Ct        string  `json:"ct"`
	// Lat       string `json:"lat"` // aladdin暂未使用
	// Lng       string `json:"lng"` // aladdin暂未使用
	// Gpst      string `json:"gpst"` // aladdin暂未使用
	// GpstAccuracy string `json:"gps_accuracy"` // aladdin暂未使用
	// IMSI         string `json:"imsi"` // aladdin暂未使用
	// GpsType      int    `json:"gps_type"`
}

type RankerInfo struct {
	PowerRate      int     `json:"power_rate,omitempty"`
	Charging       int     `json:"charging,omitempty"`
	TotalMemory    string  `json:"total_memory,omitempty"`
	ResidualMemory string  `json:"residual_memory,omitempty"`
	CID            string  `json:"cid,omitempty"`
	LAT            string  `json:"lat,omitempty"`
	LNG            string  `json:"lng,omitempty"`
	GPST           string  `json:"gpst,omitempty"`
	GPSAccuracy    string  `json:"gps_accuracy,omitempty"`
	GPSType        string  `json:"gps_type,omitempty"`
	Dmt            float64 `json:"dmt,omitempty"`
	Dmf            float64 `json:"dmf,omitempty"`
	CpuType        string  `json:"ct,omitempty"`
	ToponChInfo    string  `json:"topon_info,omitempty"`
	PriceFactor    float64 `json:"pf,omitempty"`    // 频次控制- 价格系数
	HBMn           string  `json:"hb_mn,omitempty"` // hb 聚合平台名字
}

func init() {
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
}

func Base64Encode(data string) string {
	base64Data := base64.StdEncoding.EncodeToString([]byte(data))
	result := make([]byte, len(base64Data))
	for i := range base64Data {
		result[i] = encodeCharMapping[int(base64Data[i])]
	}
	return string(result)
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

func SubUtf8Str(str string, n int) string {
	if utf8.RuneCountInString(str) <= n {
		return str
	}
	runeStr := []rune(str)
	str = string(runeStr[0:n]) + "......"
	return str
}

func NumFormat(v float64, n int) float64 {
	pow10_n := math.Pow10(n)
	return math.Trunc((v)*pow10_n) / pow10_n
}

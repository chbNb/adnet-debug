package helpers

import (
	"hash/crc32"
	"math/rand"
)

func InArray(val int, array []int) bool {
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

func AnyExists(data map[string]bool, keys ...string) bool {
	if len(data) == 0 {
		return false
	}
	for _, k := range keys {
		if len(k) == 0 {
			continue
		}
		if _, ok := data[k]; ok {
			return true
		}
	}
	return false
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

func RandInt64Arr(arr []int64) int64 {
	llen := len(arr)
	if llen <= 0 {
		return int64(0)
	}
	rand := rand.Intn(llen)
	return arr[rand]
}

func InInt64Arr(val int64, array []int64) bool {
	for _, v := range array {
		if val == v {
			return true
		}
	}
	return false
}

func getRandString(gaid string, idfa string) string {
	if len(idfa) <= 0 || idfa == "00000000-0000-0000-0000-000000000000" || idfa == "idfa" || idfa == "-" {
		idfa = ""
	}
	if len(gaid) <= 0 || gaid == "gaid" || gaid == "00000000-0000-0000-0000-000000000000" || gaid == "-" {
		gaid = ""
	}
	return idfa + gaid
}

func GetRandConsiderZero(gaid string, idfa string, salt string, randSum int) int {
	str := getRandString(gaid, idfa)
	if len(str) <= 0 {
		return -1
	}
	str = str + salt
	return int(crc32.ChecksumIEEE([]byte(str))) % randSum
}
func GetRandByStr(str string, randSum int) int {
	return int(crc32.ChecksumIEEE([]byte(str))) % randSum
}

func GetRandByStrAddSalt(str string, randSum int, salt string) int {
	str = str + salt
	return int(crc32.ChecksumIEEE([]byte(str))) % randSum
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

func PickItem(data []string) string {
	dataLen := len(data)
	if dataLen == 1 {
		return data[0]
	} else if dataLen > 1 {
		index := rand.Intn(dataLen)
		return data[index]
	} else {
		return ""
	}
}

func InStrArray(val string, array []string) bool {
	for _, v := range array {
		if val == v {
			return true
		}
	}
	return false
}

package helpers

import (
	"strconv"
	"strings"
)

func GetVideoLength(lenth string) int {
	result := 0
	lenArr := strings.Split(lenth, ":")
	if len(lenArr) != 3 {
		return result
	}
	hour, err := strconv.Atoi(lenArr[0])
	if err != nil {
		return result
	}
	result += hour * 3600
	min, err := strconv.Atoi(lenArr[1])
	if err != nil {
		return result
	}
	result += min * 60
	secArr := strings.Split(lenArr[2], ".")
	if len(secArr) == 0 {
		return result
	}
	sec, err := strconv.Atoi(secArr[0])
	if err != nil {
		return result
	}
	result += sec

	return result
}

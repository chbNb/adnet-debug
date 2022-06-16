package helpers

import (
	"strconv"
	"strings"
)

func VersionCode(osVersion string) (int32, error) {
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

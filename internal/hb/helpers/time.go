package helpers

import (
	"errors"
	"strconv"
	"time"
)

//设置UTC的时区
func GetLocationStrByUtc(timeInt int) (string, error) {

	if timeInt > 14 || timeInt < -12 {
		return "", errors.New(strconv.Itoa(timeInt) + " timezone is error!")
	}

	timeInt = 0 - timeInt
	timeStr := strconv.Itoa(timeInt)

	if timeInt >= 0 {
		return "Etc/GMT+" + timeStr, nil
	} else {
		return "Etc/GMT" + timeStr, nil
	}
}

func GetLocationTimeFromTimezone(timeInt int, t time.Time) (*time.Time, error) {
	tz, err := GetLocationStrByUtc(timeInt)
	if err != nil {
		return nil, err
	}
	local, err := time.LoadLocation(tz)
	if err != nil {
		return nil, err
	}
	t = t.In(local)

	return &t, nil
}

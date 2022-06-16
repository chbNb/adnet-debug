package params

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type HttpQueryMap map[string][]string

type HttpReqData struct {
	Path      string
	Host      string
	PostData  []byte
	UserAgent string
	ClientIp  string
	QueryData HttpQueryMap
}

var HttpQueryMapStringKeyNotFound = errors.New("key is not exists and have no defaultValue")

func (v HttpQueryMap) GetString(key string, filterSpace bool, def ...string) string {
	if val, ok := v[key]; ok {
		rawString := strings.Join(val, "")
		rawString = strings.TrimSpace(rawString)
		if filterSpace {
			rawString = strings.Replace(rawString, "\t", "", -1)
			rawString = strings.Replace(rawString, "\n", "", -1)
		}
		return rawString
	}
	if len(def) > 0 {
		return def[0]
	}
	return ""
}

func (v HttpQueryMap) GetInt(key string, def ...int) (int, error) {
	if val, ok := v[key]; ok {
		return strconv.Atoi(strings.TrimSpace(strings.Join(val, "")))
	}
	if len(def) > 0 {
		return def[0], nil
	}
	return 0, HttpQueryMapStringKeyNotFound
}

func (v HttpQueryMap) GetInt64(key string, def ...int64) (int64, error) {
	if val, ok := v[key]; ok {
		return strconv.ParseInt(strings.TrimSpace(strings.Join(val, "")), 10, 64)
	}
	if len(def) > 0 {
		return def[0], nil
	}
	return -1, HttpQueryMapStringKeyNotFound
}

func (v HttpQueryMap) GetFloat(key string, def ...float64) (float64, error) {
	if val, ok := v[key]; ok {
		return strconv.ParseFloat(strings.TrimSpace(strings.Join(val, "")), 64)
	}
	if len(def) > 0 {
		return def[0], nil
	}
	return 0.0, fmt.Errorf("parse key=[%s] is not exists and have no defaultValue", key)
}

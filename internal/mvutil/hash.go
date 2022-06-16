package mvutil

import (
	"github.com/json-iterator/go"
	"strconv"
)

func APHash(byteStr []byte) string {
	var hash uint32 = 0
	for i := 0; i < len(byteStr); i++ {
		if (i & 1) == 0 {
			hash ^= ((hash << 7) ^ uint32(byteStr[i]) ^ (hash >> 3))
		} else {
			hash ^= (^((hash << 11) ^ uint32(byteStr[i]) ^ (hash >> 5)) + 1)
		}
	}

	return strconv.FormatUint(uint64(hash&0x7FFFFFFF), 36) + "-" + strconv.FormatUint(uint64(len(string(byteStr))), 36)
}

func APHashByObj(obj interface{}) (string, error) {

	byteStr, err := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(obj)
	if err != nil {
		return "", err
	}

	var hash uint32 = 0
	for i := 0; i < len(byteStr); i++ {
		if (i & 1) == 0 {
			hash ^= ((hash << 7) ^ uint32(byteStr[i]) ^ (hash >> 3))
		} else {
			hash ^= (^((hash << 11) ^ uint32(byteStr[i]) ^ (hash >> 5)) + 1)
		}
	}

	return strconv.FormatUint(uint64(hash&0x7FFFFFFF), 36) + "-" + strconv.FormatUint(uint64(len(string(byteStr))), 36), nil
}

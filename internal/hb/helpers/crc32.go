package helpers

import (
	"fmt"
	"hash/crc32"
)

func GetStrCrc32(s string) string {
	return fmt.Sprintf("%08x", crc32.ChecksumIEEE([]byte(s)))
}

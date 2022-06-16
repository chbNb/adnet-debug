package helpers

import "gitlab.mobvista.com/ADN/adnet/internal/hb/constant"

func ImageSizeStr(id int) string {
	imageSizeStr, ok := constant.ImageSizeMap[id]
	if !ok {
		return constant.IMAGE_SIZE_STR_128x128
	}
	return imageSizeStr
}

package helpers

import "gitlab.mobvista.com/ADN/adnet/internal/hb/constant"

func TemplateName(templateId int32) string {
	switch templateId {
	case constant.MIDDLE_BLACK_SCREEN:
		return "middle_black_screen"
	case constant.MIDDLE_FUR_SCREEN:
		return "middle_fur_screen"
	case constant.ABOVE_VIDEO:
		return "above_video"
	case constant.STOREKIT_VIODE:
		return "storekit_video"
	case constant.IMAGE_VIDOE:
		return "image_video"
	case constant.STRETCH_SCREEN:
		return "stretch_screen"
	default:
		return ""
	}
}

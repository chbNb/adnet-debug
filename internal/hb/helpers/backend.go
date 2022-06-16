package helpers

import "gitlab.mobvista.com/ADN/adnet/internal/hb/constant"

func GetBackendId(dspId int64) int {
	switch dspId {
	case constant.FakeToutiao:
		return 8
	case constant.FakeGdt:
		return 10
	default:
		return 0
	}
}

func GetDspId(backendId int) int {
	switch backendId {
	case constant.TouTiao:
		return constant.FakeToutiao
	case constant.Gdt:
		return constant.FakeGdt
	default:
		return 0
	}
}

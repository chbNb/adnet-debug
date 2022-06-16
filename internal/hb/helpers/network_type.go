package helpers

import (
	"gitlab.mobvista.com/ADN/adnet/internal/hb/constant"
)

func NetworkTypeStr(networkType int) string {
	switch networkType {
	case constant.T2G:
		return constant.N2G
	case constant.T3G:
		return constant.N3G
	case constant.T4G:
		return constant.N4G
	case constant.TWIFI:
		return constant.NWIFI
	}
	return constant.NUNKNOWN
}

package mvconst

const (
	NETWORK_TYPE_UNKNOWN   = 0
	NETWORK_TYPE_2G        = 2
	NETWORK_TYPE_3G        = 3
	NETWORK_TYPE_4G        = 4
	NETWORK_TYPE_5G        = 5
	NETWORK_TYPE_UNKNOWN_1 = 6 // sdk读取网络状态之前默认是1，但如果调用系统API读取了，但获取的状态无法识别，则用6表示 unknown，666之前的版本这里为1,v666开始设置为6(下一代网络)
	NETWORK_TYPE_WIFI      = 9
)

const (
	NETWORK_TYPE_NAME_UNKNOWN   = "unknown"
	NETWORK_TYPE_NAME_2G        = "2g"
	NETWORK_TYPE_NAME_3G        = "3g"
	NETWORK_TYPE_NAME_4G        = "4g"
	NETWORK_TYPE_NAME_5G        = "5g"
	NETWORK_TYPE_NAME_UNKNOWN_1 = "unknown"
	NETWORK_TYPE_NAME_WIFI      = "wifi"
)

func GetNetworkName(network int) string {
	nwMap := map[int]string{
		NETWORK_TYPE_UNKNOWN:   NETWORK_TYPE_NAME_UNKNOWN,
		NETWORK_TYPE_2G:        NETWORK_TYPE_NAME_2G,
		NETWORK_TYPE_3G:        NETWORK_TYPE_NAME_3G,
		NETWORK_TYPE_4G:        NETWORK_TYPE_NAME_4G,
		NETWORK_TYPE_5G:        NETWORK_TYPE_NAME_5G,
		NETWORK_TYPE_UNKNOWN_1: NETWORK_TYPE_NAME_UNKNOWN_1,
		NETWORK_TYPE_WIFI:      NETWORK_TYPE_NAME_WIFI,
	}
	name, ok := nwMap[network]
	if ok {
		return name
	}
	return NETWORK_TYPE_NAME_UNKNOWN
}

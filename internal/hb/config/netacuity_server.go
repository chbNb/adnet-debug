package config

type NetacuitySvrCfg struct {
	FeatureCode  int
	APIId        int
	NetServerIp  string
	TimeoutDelay int
	Expire       int64
	GRPCAddress  string
}

func parseNetacuitySvrCfg(fileName string) (*NetacuitySvrCfg, error) {
	viper, err := ReadConfig(fileName)
	if err != nil {
		return nil, err
	}
	netSvrCfg := &NetacuitySvrCfg{}
	netSvrCfg.FeatureCode = viper.GetInt("feature_code")
	netSvrCfg.APIId = viper.GetInt("api_id")
	netSvrCfg.NetServerIp = viper.GetString("net_server_ip")
	netSvrCfg.TimeoutDelay = viper.GetInt("timeout_delay")
	netSvrCfg.Expire = viper.GetInt64("expire")
	netSvrCfg.GRPCAddress = viper.GetString("grpc_address")
	return netSvrCfg, nil
}

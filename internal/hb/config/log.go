package config

type LogCfg struct {
	Bid                  string
	Load                 string
	Event                string
	Watch                string
	Creative             string
	DspCreative          string
	Runtime              string
	Request              string
	LRequest             string
	ReqMonitor           string
	TreasureBox          string
	ConsulAerospike      string
	ConsulAdx            string
	ConsulWatch          string
	ConsulAdnetAerospike string
	MappingServer        string
	ConsulMappingServer  string
	DeviceAerospike      string
}

func parseLogCfg(configDir, cloud, region, fileName string) (*LogCfg, error) {
	viper, err := ReadConfig(getConfigFile(configDir, cloud, region, fileName))
	if err != nil {
		return nil, err
	}
	logCfg := &LogCfg{}
	logCfg.Bid = getConfigFile(configDir, cloud, region, viper.GetString("bid_log"))
	logCfg.Load = getConfigFile(configDir, cloud, region, viper.GetString("load_log"))
	logCfg.Event = getConfigFile(configDir, cloud, region, viper.GetString("event_log"))
	logCfg.Watch = getConfigFile(configDir, cloud, region, viper.GetString("watch_log"))
	logCfg.Creative = getConfigFile(configDir, cloud, region, viper.GetString("creative_log"))
	logCfg.DspCreative = getConfigFile(configDir, cloud, region, viper.GetString("dsp_creative_log"))
	logCfg.Runtime = getConfigFile(configDir, cloud, region, viper.GetString("runtime_log"))
	logCfg.Request = getConfigFile(configDir, cloud, region, viper.GetString("req_log"))
	logCfg.LRequest = getConfigFile(configDir, cloud, region, viper.GetString("loss_req_log"))
	logCfg.ReqMonitor = getConfigFile(configDir, cloud, region, viper.GetString("req_monitor_log"))
	logCfg.TreasureBox = getConfigFile(configDir, cloud, region, viper.GetString("treasure_box_log"))
	logCfg.ConsulAerospike = getConfigFile(configDir, cloud, region, viper.GetString("consul_aerospike_log"))
	logCfg.ConsulAdx = getConfigFile(configDir, cloud, region, viper.GetString("consul_adx_log"))
	logCfg.ConsulWatch = getConfigFile(configDir, cloud, region, viper.GetString("consul_watch_log"))
	logCfg.ConsulAdnetAerospike = getConfigFile(configDir, cloud, region, viper.GetString("consul_adnet_aerospike_log"))
	logCfg.MappingServer = getConfigFile(configDir, cloud, region, viper.GetString("mapping_server_log"))
	logCfg.ConsulMappingServer = getConfigFile(configDir, cloud, region, viper.GetString("consul_mapping_server_log"))
	logCfg.DeviceAerospike = getConfigFile(configDir, cloud, region, viper.GetString("device_aerospike_log"))
	return logCfg, nil
}

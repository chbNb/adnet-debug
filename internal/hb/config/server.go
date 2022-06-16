package config

type ServerCfg struct {
	HTTPAddr           string
	ServerMetaUrl      string
	AerospikeMultiZone bool
}

func parseServerCfg(fileName string) (*ServerCfg, error) {
	viper, err := ReadConfig(fileName)
	if err != nil {
		return nil, err
	}
	serverCfg := &ServerCfg{}
	serverCfg.HTTPAddr = viper.GetString("http_addr")
	serverCfg.ServerMetaUrl = viper.GetString("server_meta_url")
	serverCfg.AerospikeMultiZone = viper.GetBool("aerospike_multi_zone")
	return serverCfg, nil
}

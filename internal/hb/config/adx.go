package config

import "os"

type AdxCfg struct {
	EndPoint       string
	ServiceName    string
	ForceUseConfig bool
	TimeMax        int
}

func parseAdxCfg(fileName string) (*AdxCfg, error) {
	viper, err := ReadConfig(fileName)
	if err != nil {
		return nil, err
	}
	adxCfg := &AdxCfg{}
	adxCfg.EndPoint = viper.GetString("end_point")
	adxCfg.TimeMax = viper.GetInt("time_max")
	adxCfg.ForceUseConfig = viper.GetBool("force_use_config")
	if os.Getenv("ADX_SERVICE_NAME") != "" {
		adxCfg.ServiceName = os.Getenv("ADX_SERVICE_NAME")
	}
	return adxCfg, nil
}

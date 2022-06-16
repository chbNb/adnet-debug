package config

import (
	"log"
	"path/filepath"

	"github.com/mae-pax/consul-loadbalancer/balancer"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
)

const (
	ConfigDirFlag = "config"
	HTTPAddrFlag  = "http_addr"
	CloudFlag     = "cloud"
	RegionFlag    = "region"
	ModeFlag      = "mode"
	StartModeFlag = "start_mode"
)

const (
	configDir           = "./config"
	serverConfig        = "server.yaml"
	logConfig           = "log.yaml"
	mgoConfig           = "mgo.yaml"
	netSvrConfig        = "netacuity_server.yaml"
	adxConfig           = "adx.yaml"
	aerospikeConfig     = "aerospike.yaml"
	aerospikeConfigAs   = "aerospike_as.yaml"
	redisConfig         = "redis.yaml"
	creativeCacheConfig = "redis_cache.yaml"
	uaParserConfig      = "hb_ua_regexes.yaml"
	tkConfig            = "tk.yaml"
	tbConfig            = "treasure_box.yaml"
	consulConfig        = "consul.yaml"
	geoConfig           = "geo.yaml"
)

type Config struct {
	ServerCfg                 *ServerCfg
	LogCfg                    *LogCfg
	ExtraCfg                  *mvutil.ExtraConfig
	NetSvrCfg                 *NetacuitySvrCfg
	AdxCfg                    *AdxCfg
	TkCfg                     *TkCfg
	ConsulCfg                 *mvutil.Consul
	AerospikeConsulBuild      *balancer.ConsulResolver
	AdxConsulBuild            *balancer.ConsulResolver
	AdnetAerospikeConsulBuild *balancer.ConsulResolver
	TBConfigPath              string
	GeoConfigPath             string
}

func DefaultConfigDir() string {
	return configDir
}

func DefaultSrvAddr() string {
	return ":9102"
}

func UaParserConfig(configDir, cloud, region string) string {
	return getConfigFile(configDir, cloud, region, uaParserConfig)
}

func CreativeCacheConfig(configDir, cloud, region string) string {
	return getConfigFile(configDir, cloud, region, creativeCacheConfig)
}

func RedisConfig(configDir, cloud, region string) string {
	return getConfigFile(configDir, cloud, region, redisConfig)
}

func AerospikeConfig(configDir, cloud, region string) string {
	return getConfigFile(configDir, cloud, region, aerospikeConfig)
}

func AerospikeConfigWithZone(configDir, cloud, region, zone string) string {
	fileName := "aerospike_" + zone + ".yaml"
	return getConfigFile(configDir, cloud, region, fileName)
}

func AerospikeConfigAs(configDir, cloud, region string) string {
	return getConfigFile(configDir, cloud, region, aerospikeConfigAs)
}

func getConfigFile(configDir, cloud, region, fileName string) string {
	cfgDir := DefaultConfigDir()
	if len(configDir) > 0 {
		cfgDir = configDir
	}
	cfgFile := filepath.Join(cfgDir, cloud, region, fileName)
	log.Printf("config file: %s", cfgFile)
	return cfgFile
}
func GetTBRegisterConfigPath(configDir string) string {
	configPath := DefaultConfigDir()
	if len(configDir) > 0 {
		configPath = configDir
	}
	log.Printf("GetTBRegisterConfigPath: %s", configPath)
	return configPath
}

func ParseConfig(configDir, cloud, region string) (*Config, error) {
	cfg := &Config{}

	srvCfg, err := parseServerCfg(getConfigFile(configDir, cloud, region, serverConfig))
	if err != nil {
		return nil, errors.Wrap(err, "parseServerCfg")
	}
	cfg.ServerCfg = srvCfg

	logCfg, err := parseLogCfg(configDir, cloud, region, logConfig)
	if err != nil {
		return nil, errors.Wrap(err, "parseLogCfg")
	}
	cfg.LogCfg = logCfg

	etrCfg, err := mvutil.ParseExtraConfig(getConfigFile(configDir, cloud, region, mgoConfig))
	if err != nil {
		return nil, errors.Wrap(err, "parseExtraConfig")
	}
	cfg.ExtraCfg = etrCfg

	// netSvrCfg, err := parseNetacuitySvrCfg(getConfigFile(configDir, cloud, region, netSvrConfig))
	// if err != nil {
	// 	return nil, errors.Wrap(err, "parseNetacuitySvrCfg")
	// }
	// cfg.NetSvrCfg = netSvrCfg

	adxCfg, err := parseAdxCfg(getConfigFile(configDir, cloud, region, adxConfig))
	if err != nil {
		return nil, errors.Wrap(err, "parseAdxCfg")
	}
	cfg.AdxCfg = adxCfg

	tkCfg, err := parseTkCfg(getConfigFile(configDir, cloud, region, tkConfig))
	if err != nil {
		return nil, errors.Wrap(err, "parseTkCfg")
	}
	cfg.TkCfg = tkCfg

	consulCfg, err := mvutil.ParseConsulCfg(getConfigFile(configDir, cloud, region, consulConfig))
	if err != nil {
		return nil, errors.Wrap(err, "parseConsulCfg")
	}
	cfg.ConsulCfg = consulCfg

	cfg.TBConfigPath = getConfigFile(configDir, cloud, region, tbConfig)
	cfg.GeoConfigPath = getConfigFile(configDir, cloud, region, geoConfig)

	return cfg, nil
}

var cfg *Config

func SetConfig(conf *Config) {
	cfg = conf
}

func GetConfig() *Config {
	return cfg
}

func ReadConfig(filePath string) (*viper.Viper, error) {
	config := viper.New()
	config.SetConfigFile(filePath)
	if err := config.ReadInConfig(); err != nil {
		return nil, err
	}
	return config, nil
}

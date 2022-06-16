package mvutil

import (
	"time"

	"github.com/mae-pax/consul-loadbalancer/balancer"
)

type Consul struct {
	Aerospike       *Aerospike
	AdnetAerospike  *Aerospike
	Adx             *Adx
	MappingServer   *MappingServer
	Cloud           string
	Address         string
	KeyPath         string
	CpuThreshold    string
	InstanceFactor  string
	OnlineLabFactor string
	Services        string
	ZoneCPU         string
}

type MappingServer struct {
	Enable      bool
	ServiceName string
	Timeout     time.Duration
	Interval    time.Duration
	HttpUrl     string
}

type Adx struct {
	Enable      bool
	ServiceName string
	Timeout     time.Duration
	Interval    time.Duration
}

type Aerospike struct {
	Enable              bool
	ServiceName         string
	Timeout             time.Duration
	WriteTimeout        time.Duration
	Interval            time.Duration
	Expiration          time.Duration
	Namespace           string
	SetName             string
	BinName             string
	Retries             int
	ConnectionQueueSize int
}

func NewAdxConsulResolver(cfg *Consul) (*balancer.ConsulResolver, error) {
	return balancer.NewConsulResolver(
		cfg.Cloud,
		cfg.Address,
		cfg.Adx.ServiceName,
		cfg.KeyPath+"/"+cfg.Adx.ServiceName+"/"+cfg.CpuThreshold,
		cfg.KeyPath+"/"+cfg.Adx.ServiceName+"/"+cfg.ZoneCPU,
		cfg.KeyPath+"/"+cfg.Adx.ServiceName+"/"+cfg.InstanceFactor,
		cfg.KeyPath+"/"+cfg.Adx.ServiceName+"/"+cfg.OnlineLabFactor,
		cfg.Adx.Interval,
		cfg.Adx.Timeout,
		cfg.KeyPath+"/"+cfg.Adx.ServiceName+"/"+cfg.Services,
	)
}

func NewHBConsulResolver(cfg *Consul) (*balancer.ConsulResolver, error) {
	return balancer.NewConsulResolver(
		cfg.Cloud,
		cfg.Address,
		cfg.Aerospike.ServiceName,
		cfg.KeyPath+"/"+cfg.Aerospike.ServiceName+"/"+cfg.CpuThreshold,
		cfg.KeyPath+"/"+cfg.Aerospike.ServiceName+"/"+cfg.ZoneCPU,
		cfg.KeyPath+"/"+cfg.Aerospike.ServiceName+"/"+cfg.InstanceFactor,
		cfg.KeyPath+"/"+cfg.Aerospike.ServiceName+"/"+cfg.OnlineLabFactor,
		cfg.Aerospike.Interval,
		cfg.Aerospike.Timeout,
	)
}

func NewAdnetConsulResolver(cfg *Consul) (*balancer.ConsulResolver, error) {
	return balancer.NewConsulResolver(
		cfg.Cloud,
		cfg.Address,
		cfg.AdnetAerospike.ServiceName,
		cfg.KeyPath+"/"+cfg.AdnetAerospike.ServiceName+"/"+cfg.CpuThreshold,
		cfg.KeyPath+"/"+cfg.AdnetAerospike.ServiceName+"/"+cfg.ZoneCPU,
		cfg.KeyPath+"/"+cfg.AdnetAerospike.ServiceName+"/"+cfg.InstanceFactor,
		cfg.KeyPath+"/"+cfg.AdnetAerospike.ServiceName+"/"+cfg.OnlineLabFactor,
		cfg.AdnetAerospike.Interval,
		cfg.AdnetAerospike.Timeout,
	)
}

func NewMappingServerConsulResolver(cfg *Consul) (*balancer.ConsulResolver, error) {
	return balancer.NewConsulResolver(
		cfg.Cloud,
		cfg.Address,
		cfg.MappingServer.ServiceName,
		cfg.KeyPath+"/mapping_service/"+cfg.CpuThreshold,
		cfg.KeyPath+"/mapping_service/"+cfg.ZoneCPU,
		cfg.KeyPath+"/mapping_service/"+cfg.InstanceFactor,
		cfg.KeyPath+"/mapping_service/"+cfg.OnlineLabFactor,
		cfg.MappingServer.Interval,
		cfg.MappingServer.Timeout,
	)
}

func ParseConsulCfg(fileName string) (*Consul, error) {
	viper, err := ReadConfig(fileName)
	if err != nil {
		return nil, err
	}
	aerospikeCfg := &Aerospike{}
	aerospikeCfg.Enable = viper.GetBool("aerospike.enable")
	aerospikeCfg.ServiceName = viper.GetString("aerospike.service_name")
	aerospikeCfg.Timeout = viper.GetDuration("aerospike.timeout")
	aerospikeCfg.WriteTimeout = viper.GetDuration("aerospike.write_timeout")
	aerospikeCfg.Interval = viper.GetDuration("aerospike.interval")
	aerospikeCfg.Expiration = viper.GetDuration("aerospike.expiration")
	aerospikeCfg.Namespace = viper.GetString("aerospike.namespace")
	aerospikeCfg.SetName = viper.GetString("aerospike.setname")
	aerospikeCfg.Retries = viper.GetInt("aerospike.retries")
	aerospikeCfg.ConnectionQueueSize = viper.GetInt("aerospike.connection_queue_size")

	adxCfg := &Adx{}
	adxCfg.Enable = viper.GetBool("adx.enable")
	adxCfg.ServiceName = viper.GetString("adx.service_name")
	adxCfg.Timeout = viper.GetDuration("adx.timeout")
	adxCfg.Interval = viper.GetDuration("adx.interval")

	adnetAerospikeCfg := &Aerospike{}
	adnetAerospikeCfg.Enable = viper.GetBool("adnet_aerospike.enable")
	adnetAerospikeCfg.ServiceName = viper.GetString("adnet_aerospike.service_name")
	adnetAerospikeCfg.Timeout = viper.GetDuration("adnet_aerospike.timeout")
	adnetAerospikeCfg.Interval = viper.GetDuration("adnet_aerospike.interval")
	adnetAerospikeCfg.Expiration = viper.GetDuration("adnet_aerospike.expiration")
	adnetAerospikeCfg.Namespace = viper.GetString("adnet_aerospike.namespace")
	adnetAerospikeCfg.SetName = viper.GetString("adnet_aerospike.setname")
	adnetAerospikeCfg.Retries = viper.GetInt("adnet_aerospike.retries")
	adnetAerospikeCfg.ConnectionQueueSize = viper.GetInt("adnet_aerospike.connection_queue_size")

	mappingServerCfg := &MappingServer{}
	mappingServerCfg.Enable = viper.GetBool("mapping_server.enable")
	mappingServerCfg.ServiceName = viper.GetString("mapping_server.service_name")
	mappingServerCfg.Timeout = viper.GetDuration("mapping_server.timeout")
	mappingServerCfg.Interval = viper.GetDuration("mapping_server.interval")
	mappingServerCfg.HttpUrl = viper.GetString("mapping_server.httpUrl")

	consulCfg := &Consul{}
	consulCfg.Aerospike = aerospikeCfg
	consulCfg.Adx = adxCfg
	consulCfg.AdnetAerospike = adnetAerospikeCfg
	consulCfg.MappingServer = mappingServerCfg
	consulCfg.Address = viper.GetString("address")
	consulCfg.KeyPath = viper.GetString("key_path")
	consulCfg.CpuThreshold = viper.GetString("cpu_threshold")
	consulCfg.InstanceFactor = viper.GetString("instance_factor")
	consulCfg.OnlineLabFactor = viper.GetString("onlinelab_factor")
	consulCfg.ZoneCPU = viper.GetString("zone_cpu")
	consulCfg.Cloud = viper.GetString("cloud")
	consulCfg.Services = viper.GetString("services")

	return consulCfg, nil
}

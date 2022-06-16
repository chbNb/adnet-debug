package req_context

import (
	"encoding/binary"
	"errors"
	"hash/crc32"
	"net"
	"strconv"
	"time"

	"github.com/easierway/concurrent_map"
	"github.com/mae-pax/consul-loadbalancer/balancer"
	"github.com/ua-parser/uap-go/uaparser"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/config"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/mlog"
	"gitlab.mobvista.com/ADN/adnet/internal/hb/storage"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	"gitlab.mobvista.com/ADN/exporter/metrics"
	"gitlab.mobvista.com/mtech/mkv/pkg/kvclient"
)

const (
	HongKong  = "hongkong"
	Virginia  = "virginia"
	Singapore = "singapore"
	Seoul     = "seoul"
	Frankfurt = "frankfurt"
	Ohio      = "ohio"
)

type ReqContext struct {
	Cfg                 *config.Config
	BidCacheClient      *storage.KVConfig
	AdCacheClient       *storage.KVConfig
	AsCacheClient       *storage.KVConfig
	CreativeCacheClient *storage.KVRedis
	UaParser            *uaparser.Parser
	ServerIp            string
	Region              string
	Cloud               string
	Zone                string
	MLogs               *mlog.MLogger
	AerospikeClientMap  *concurrent_map.ConcurrentMap
	AerospikeOnlineConf *mvutil.HBAerospikeConf
}

var ctx *ReqContext

func GetInstance() *ReqContext {
	if ctx == nil {
		ctx = &ReqContext{AerospikeClientMap: concurrent_map.CreateConcurrentMap(99)}
	}
	return ctx
}

func (reqCtx *ReqContext) Close() {
	reqCtx.BidCacheClient.Close()
	// reqCtx.AdCacheClient.Close()
	// reqCtx.AsCacheClient.Close()
	reqCtx.CreativeCacheClient.Conn.Close()
	if reqCtx.AdCacheClient != nil {
		reqCtx.AdCacheClient.Close()
	}
	if reqCtx.Cfg.AerospikeConsulBuild != nil {
		reqCtx.Cfg.AerospikeConsulBuild.Stop()
	}
	if reqCtx.Cfg.AdxConsulBuild != nil {
		reqCtx.Cfg.AdxConsulBuild.Stop()
	}
	if reqCtx.Cfg.AdnetAerospikeConsulBuild != nil {
		reqCtx.Cfg.AdnetAerospikeConsulBuild.Stop()
	}
}

func ip2int(ipAddr string) (uint32, error) {
	ip := net.ParseIP(ipAddr)
	if ip == nil {
		return 0, errors.New("wrong ipAddr format")
	}
	ip = ip.To4()
	return binary.BigEndian.Uint32(ip), nil
}

func (reqCtx *ReqContext) GetMkvAerospikeClient() (*storage.AerospikeClient, error) {
	node := reqCtx.Cfg.AerospikeConsulBuild.SelectNode()
	metrics.IncCounterWithLabelValues(22, "hb_aerospike", mvutil.Zone(), node.Zone)
	key, _ := ip2int(node.Host)
	if val, found := reqCtx.AerospikeClientMap.Get(concurrent_map.I64Key(int64(key))); found {
		if client, ok := val.(*storage.AerospikeClient); ok {
			return client, nil
		}
	}
	address := node.Host + ":" + strconv.Itoa(node.Port)
	aerospikeBuilder := kvclient.NewAerospikeBuilder()
	aerospikeBuilder.WithAddress(address)
	aerospikeBuilder.WithNamespace(reqCtx.Cfg.ConsulCfg.Aerospike.Namespace)
	aerospikeBuilder.WithSetName(reqCtx.Cfg.ConsulCfg.Aerospike.SetName)
	aerospikeBuilder.WithRetries(reqCtx.Cfg.ConsulCfg.Aerospike.Retries)
	aerospikeBuilder.WithTimeout(reqCtx.Cfg.ConsulCfg.Aerospike.Timeout)
	aerospikeBuilder.WithWriteTimeout(reqCtx.Cfg.ConsulCfg.Aerospike.WriteTimeout)
	aerospikeBuilder.WithExpiration(reqCtx.Cfg.ConsulCfg.Aerospike.Expiration)
	aerospikeBuilder.WithConnectionQueueSize(reqCtx.Cfg.ConsulCfg.Aerospike.ConnectionQueueSize)
	client, err := aerospikeBuilder.Build()
	if err != nil {
		return nil, err
	}
	ac := storage.NewAerospikeClient(client)
	ac.SetCompressor()
	ac.SetSerializer()
	reqCtx.AerospikeClientMap.Set(concurrent_map.I64Key(int64(key)), ac)
	return ac, nil
}

func (reqCtx *ReqContext) GetLoadAerospikeClient(conf *mvutil.HBAerospikeConf, migrate bool) (*storage.AerospikeClient, error) {
	return reqCtx.GetMkvAerospikeClientWithConf(conf, migrate)
}

func (reqCtx *ReqContext) GetBidAerospikeClient(conf *mvutil.HBAerospikeConf, migrate bool) (*storage.AerospikeClient, error) {
	return reqCtx.GetMkvAerospikeClientWithConf(conf, migrate)
}

func (reqCtx *ReqContext) GetMkvAerospikeClientWithConf(conf *mvutil.HBAerospikeConf, migrate bool) (*storage.AerospikeClient, error) {
	address := conf.Endpoint
	key := crc32.ChecksumIEEE([]byte(address))
	if migrate {
		address = conf.MigrateEndpoint
		key = crc32.ChecksumIEEE([]byte(address))
	}
	// 不更新配置不需要重新 build aerospike client
	if reqCtx.AerospikeOnlineConf != nil && reqCtx.AerospikeOnlineConf.Updated == conf.Updated {
		if val, found := reqCtx.AerospikeClientMap.Get(concurrent_map.I64Key(int64(key))); found {
			if client, ok := val.(*storage.AerospikeClient); ok {
				return client, nil
			}
		}
	}

	reqCtx.AerospikeOnlineConf = &mvutil.HBAerospikeConf{}
	reqCtx.AerospikeOnlineConf.Updated = conf.Updated

	aerospikeBuilder := kvclient.NewAerospikeBuilder()
	aerospikeBuilder.WithAddress(address)
	aerospikeBuilder.WithNamespace(conf.Namespace)
	aerospikeBuilder.WithSetName(conf.SetName)
	aerospikeBuilder.WithRetries(conf.ReadRetry)
	aerospikeBuilder.WithWriteRetries(conf.WriteRetry)
	aerospikeBuilder.WithTimeout(time.Duration(conf.ReadTimeout) * time.Millisecond)
	aerospikeBuilder.WithWriteTimeout(time.Duration(conf.WriteTimeout) * time.Millisecond)
	aerospikeBuilder.WithExpiration(time.Duration(conf.Expiration) * time.Minute)
	aerospikeBuilder.WithConnectionQueueSize(conf.ConnectionSize)
	client, err := aerospikeBuilder.Build()
	if err != nil {
		return nil, err
	}
	ac := storage.NewAerospikeClient(client)
	ac.SetCompressor()
	ac.SetSerializer()
	reqCtx.AerospikeClientMap.Set(concurrent_map.I64Key(int64(key)), ac)
	return ac, nil
}

func (reqCtx *ReqContext) SetAdCacheClient(engine storage.EngineType, cfg string) error {
	kvClient, err := storage.NewKVStorage(engine, cfg)
	if err != nil {
		return err
	}
	kvClient.SetCompressor(storage.NewSimpleCompressor())
	kvClient.SetSerializer(storage.NewSimpleSerializer())
	reqCtx.AdCacheClient = kvClient
	return nil
}

// func (reqCtx *ReqContext) SetAsCacheClient(engine storage.EngineType, cfg string) error {
// 	kvClient, err := storage.NewKVStorage(engine, cfg)
// 	if err != nil {
// 		return err
// 	}
// 	reqCtx.AsCacheClient = kvClient
// 	return nil
// }

func (reqCtx *ReqContext) SetAerospikeConsulBuild() error {
	asServiceName := reqCtx.Cfg.ConsulCfg.Aerospike.ServiceName
	if reqCtx.Cfg.ServerCfg.AerospikeMultiZone {
		asServiceName = asServiceName + "-" + mvutil.Zone()
	}
	consulResolver, err := balancer.NewConsulResolver(
		reqCtx.Cloud,
		reqCtx.Cfg.ConsulCfg.Address,
		asServiceName,
		reqCtx.Cfg.ConsulCfg.KeyPath+"/"+reqCtx.Cfg.ConsulCfg.Aerospike.ServiceName+"/"+reqCtx.Cfg.ConsulCfg.CpuThreshold,
		reqCtx.Cfg.ConsulCfg.KeyPath+"/"+reqCtx.Cfg.ConsulCfg.Aerospike.ServiceName+"/"+reqCtx.Cfg.ConsulCfg.ZoneCPU,
		reqCtx.Cfg.ConsulCfg.KeyPath+"/"+reqCtx.Cfg.ConsulCfg.Aerospike.ServiceName+"/"+reqCtx.Cfg.ConsulCfg.InstanceFactor,
		reqCtx.Cfg.ConsulCfg.KeyPath+"/"+reqCtx.Cfg.ConsulCfg.Aerospike.ServiceName+"/"+reqCtx.Cfg.ConsulCfg.OnlineLabFactor,
		reqCtx.Cfg.ConsulCfg.Aerospike.Interval,
		reqCtx.Cfg.ConsulCfg.Aerospike.Timeout,
	)
	if err != nil {
		return err
	}
	reqCtx.Cfg.AerospikeConsulBuild = consulResolver
	return nil
}

func (reqCtx *ReqContext) SetAdxConsulBuild() error {
	consulResolver, err := balancer.NewConsulResolver(
		reqCtx.Cfg.ConsulCfg.Cloud,
		reqCtx.Cfg.ConsulCfg.Address,
		reqCtx.Cfg.ConsulCfg.Adx.ServiceName,
		reqCtx.Cfg.ConsulCfg.KeyPath+"/"+reqCtx.Cfg.ConsulCfg.Adx.ServiceName+"/"+reqCtx.Cfg.ConsulCfg.CpuThreshold,
		reqCtx.Cfg.ConsulCfg.KeyPath+"/"+reqCtx.Cfg.ConsulCfg.Adx.ServiceName+"/"+reqCtx.Cfg.ConsulCfg.ZoneCPU,
		reqCtx.Cfg.ConsulCfg.KeyPath+"/"+reqCtx.Cfg.ConsulCfg.Adx.ServiceName+"/"+reqCtx.Cfg.ConsulCfg.InstanceFactor,
		reqCtx.Cfg.ConsulCfg.KeyPath+"/"+reqCtx.Cfg.ConsulCfg.Adx.ServiceName+"/"+reqCtx.Cfg.ConsulCfg.OnlineLabFactor,
		reqCtx.Cfg.ConsulCfg.Adx.Interval,
		reqCtx.Cfg.ConsulCfg.Adx.Timeout,
		reqCtx.Cfg.ConsulCfg.KeyPath+"/"+reqCtx.Cfg.ConsulCfg.Adx.ServiceName+"/"+reqCtx.Cfg.ConsulCfg.Services,
	)
	if err != nil {
		return err
	}
	reqCtx.Cfg.AdxConsulBuild = consulResolver
	return nil
}

func (reqCtx *ReqContext) SetAdnetAerospikeConsulBuild(consulResolver *balancer.ConsulResolver) {
	reqCtx.Cfg.AdnetAerospikeConsulBuild = consulResolver
}

func (reqCtx *ReqContext) SetCfg(cfg *config.Config) {
	reqCtx.Cfg = cfg
}

func (reqCtx *ReqContext) SetUaParser(parser *uaparser.Parser) {
	reqCtx.UaParser = parser
}

func (reqCtx *ReqContext) SetServerIp(serverIp string) {
	reqCtx.ServerIp = serverIp
}

func (reqCtx *ReqContext) SetRegion(region string) {
	reqCtx.Region = region
	mvutil.InitRegion(region)
}

func (reqCtx *ReqContext) SetCloud(cloud string) {
	reqCtx.Cloud = cloud
	mvutil.InitZone(cloud)
}

func (reqCtx *ReqContext) SetBidCacheClient(engine storage.EngineType, cfg string) error {
	kvClient, err := storage.NewKVStorage(engine, cfg)
	if err != nil {
		return err
	}
	kvClient.SetCompressor(storage.NewReqCtxCompressor())
	kvClient.SetSerializer(storage.NewReqCtxSerializer())
	reqCtx.BidCacheClient = kvClient
	return nil
}

func (reqCtx *ReqContext) SetCreativeCacheClient(cfg, ds string) error {
	client, err := storage.NewRedisClient(cfg, ds)
	if err != nil {
		return err
	}
	reqCtx.CreativeCacheClient = client
	return nil
}

func (reqCtx *ReqContext) SetLoggers(loger *mlog.MLogger) {
	reqCtx.MLogs = loger
}

func (reqCtx *ReqContext) GetRegionPrefix() string {
	region := reqCtx.Region
	var regionPrefix string
	switch region {
	case Virginia:
		regionPrefix = "vg"
	case Singapore:
		regionPrefix = "sg"
	case Seoul:
		regionPrefix = "se"
	case Frankfurt:
		regionPrefix = "fk"
	case Ohio:
		regionPrefix = "oh"
	default:
		regionPrefix = "sg"
	}
	return regionPrefix
}

package mkv

import (
	"encoding/binary"
	"errors"
	"math/rand"
	"net"
	"strconv"

	"github.com/easierway/concurrent_map"
	"github.com/mae-pax/consul-loadbalancer/balancer"
	"gitlab.mobvista.com/ADN/adnet/internal/extractor"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	"gitlab.mobvista.com/ADN/exporter/metrics"
	"gitlab.mobvista.com/mtech/mkv/pkg/kvcfg"
	"gitlab.mobvista.com/mtech/mkv/pkg/kvclient"
	"gitlab.mobvista.com/mtech/mkv/pkg/mykv"
)

var (
	client      kvclient.KVClient
	consulBuild *balancer.ConsulResolver
	clientMap   *concurrent_map.ConcurrentMap
	cfg         *mvutil.Aerospike
)

func ip2int(ipAddr string) (uint32, error) {
	ip := net.ParseIP(ipAddr)
	if ip == nil {
		return 0, errors.New("wrong ipAddr format")
	}
	ip = ip.To4()
	return binary.BigEndian.Uint32(ip), nil
}

func SetConsulBuild(r *balancer.ConsulResolver, c *mvutil.Aerospike) {
	consulBuild = r
	cfg = c
	clientMap = concurrent_map.CreateConcurrentMap(99)
}

func getClient() (kvclient.KVClient, error) {
	node := consulBuild.SelectNode()
	metrics.IncCounterWithLabelValues(22, "mkv", mvutil.Zone(), node.Zone)
	key, _ := ip2int(node.Host)
	if val, found := clientMap.Get(concurrent_map.I64Key(int64(key))); found {
		if client, ok := val.(kvclient.KVClient); ok {
			return client, nil
		}
	}
	address := node.Host + ":" + strconv.Itoa(node.Port)
	aerospikeBuilder := kvclient.NewAerospikeBuilder()
	aerospikeBuilder.WithAddress(address)
	aerospikeBuilder.WithNamespace(cfg.Namespace)
	aerospikeBuilder.WithSetName(cfg.SetName)
	aerospikeBuilder.WithRetries(cfg.Retries)
	aerospikeBuilder.WithTimeout(cfg.Timeout)
	aerospikeBuilder.WithExpiration(cfg.Expiration)
	aerospikeBuilder.WithConnectionQueueSize(cfg.ConnectionQueueSize)
	cache, err := aerospikeBuilder.Build()
	if err != nil {
		return nil, err
	}
	var caches []kvclient.Cache
	caches = append(caches, cache)
	client := kvclient.NewBuilder().WithCaches(caches).Build()
	client.SetCompressor(&mykv.Compressor{})
	client.SetSerializer(&mykv.Serializer{})
	clientMap.Set(concurrent_map.I64Key(int64(key)), client)
	return client, nil
}

func InitClient() error {
	var err error
	client, err = kvcfg.NewKVClientWithFile(mvutil.Config.AreaConfig.HttpConfig.MKVConf)
	if err == nil {
		client.SetCompressor(&mykv.Compressor{})
		client.SetSerializer(&mykv.Serializer{})
	}
	return err
}

func MKVSetFB(key string, value string) error {
	k2 := &mykv.Key{Message: key}
	v2 := &mykv.Val{Message: value}
	err := client.Set(k2, v2)
	return err
}

func MKVGetFB(key string) (string, error) {
	var v2 mykv.Val
	k2 := &mykv.Key{Message: key}
	ok, err := client.Get(k2, &v2)
	if err != nil {
		return "", err
	}
	if ok {
		return v2.Message, nil
	}
	return "", nil
}

func GetMKVFieldVal(key string) (map[string][]byte, error) {
	if cfg == nil || !cfg.Enable {
		return client.GetASFieldVal(key)
	}
	ratio := extractor.GetUseConsulServicesV2Ratio(mvutil.Cloud(), mvutil.Region(), "aerospike")
	if ratio <= 0 || ratio < rand.Float64() {
		return client.GetASFieldVal(key)
	}
	c, err := getClient()
	if err != nil {
		return nil, err
	}
	val, err := c.GetASFieldVal(key)
	if err != nil {
		return nil, err
	}
	return val, nil
}

func SetMKVKeyField(key string, val map[string][]byte) error {
	return client.SetASKeyField(key, val)
}

package redis

import (
	"net"
	"sync"
	"time"

	"github.com/go-redis/redis"
	mlogger "github.com/mae-pax/logger"
	"gitlab.mobvista.com/ADN/adnet/internal/consuls"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
)

const (
	AlgoRedisPort = "6382"
)

type RedisPool struct {
	lck  sync.RWMutex
	pool map[string]*redis.Client

	connectTime  int
	readTimeout  int
	writeTimeout int
	poolSize     int
	log          *mlogger.Log
}

func InitPoolFromConsul(consulConfig mvutil.ConsulConfig, connectTime, readTimeout, writeTimeout, poolSize int, log *mlogger.Log) error {
	subjectlPool = &RedisPool{
		pool:         make(map[string]*redis.Client),
		connectTime:  connectTime,
		readTimeout:  readTimeout,
		writeTimeout: writeTimeout,
		poolSize:     poolSize,
		log:          log,
	}

	if err := consuls.InitCreativeRedisConsul(consulConfig, log); err != nil {
		return err
	}

	subjectAlgoPool = &RedisPool{
		pool:         make(map[string]*redis.Client),
		connectTime:  connectTime,
		readTimeout:  readTimeout,
		writeTimeout: writeTimeout,
		poolSize:     poolSize,
		log:          log,
	}

	localZone := consuls.CreativeRedisResolver.GetLocalZone()
	for _, node := range localZone.Nodes {
		client, err := subjectlPool.newClient(node.Address)
		if err == nil {
			subjectlPool.addNode(node.Address, client)
		}
		client, err = subjectAlgoPool.newClient(net.JoinHostPort(node.Host, AlgoRedisPort))
		if err == nil {
			subjectAlgoPool.addNode(node.Address, client)
		}
	}
	otherZone := consuls.CreativeRedisResolver.GetOtherZone()
	for _, node := range otherZone.Nodes {
		client, err := subjectlPool.newClient(node.Address)
		if err != nil {
			log.Warnf("create redis from consul error:%s", err.Error())
		} else {
			subjectlPool.addNode(node.Address, client)
		}
		client, err = subjectAlgoPool.newClient(net.JoinHostPort(node.Host, AlgoRedisPort))
		if err == nil {
			subjectAlgoPool.addNode(node.Address, client)
		}
	}
	// 异常处理

	return nil
}

func (rp *RedisPool) GetNode() *redis.Client {

	tryMax := 3
TRY_AGAIN:
	node := consuls.CreativeRedisResolver.DiscoverNode()
	key := node.Address

	client := rp.getNode(key)
	if client == nil {
		rp.log.Infof("new redis node from consul:%s", key)
		var err error
		if client, err = rp.newClient(key); err == nil {
			rp.addNode(key, client)
		} else {
			tryMax--
			if tryMax >= 0 {
				goto TRY_AGAIN
			}
		}
	}

	return client
}

func (rp *RedisPool) GetLen() int {
	rp.lck.RLock()
	defer rp.lck.RUnlock()
	return len(rp.pool)
}

func (rp *RedisPool) getNode(key string) *redis.Client {
	rp.lck.RLock()
	defer rp.lck.RUnlock()
	if node, ok := rp.pool[key]; ok {
		return node
	} else {
		return nil
	}
}

func (rp *RedisPool) addNode(key string, redisClient *redis.Client) {
	rp.lck.Lock()
	defer rp.lck.Unlock()
	rp.pool[key] = redisClient
}

func (rp *RedisPool) delete(key string) {
	rp.lck.Lock()
	defer rp.lck.Unlock()
	delete(rp.pool, key)
}

func (rp *RedisPool) newClient(addr string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:               addr,
		DialTimeout:        time.Duration(rp.connectTime) * time.Millisecond,
		ReadTimeout:        time.Duration(rp.readTimeout) * time.Millisecond,
		WriteTimeout:       time.Duration(rp.writeTimeout) * time.Millisecond,
		IdleTimeout:        -1,
		IdleCheckFrequency: -1,
		PoolSize:           rp.poolSize,
	})
	_, err := client.Info().Result()
	if err != nil {
		rp.log.Warnf("create redis from consul error:%s", err.Error())
		return nil, err
	} else {
		return client, nil
	}
}

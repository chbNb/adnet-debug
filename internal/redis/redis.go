package redis

import (
	"errors"
	"time"

	"github.com/go-redis/redis"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	"gitlab.mobvista.com/ADN/adnet/internal/watcher"
)

var (
	// subject       *redis.ClusterClient
	subjectAlgo     *redis.Client
	subjectl        *redis.Client
	subjectlPool    *RedisPool
	subjectAlgoPool *RedisPool
	// subjectlDB1   *redis.Client
	ipInfoCluster *redis.ClusterClient
)

func SetSubjectl(c *redis.Client) {
	if c != nil {
		subjectl = c
	}
}

func InitLocalRedis(addr string, connectTime, readTimeout, writeTimeout, poolSize int) error {
	subjectl = redis.NewClient(&redis.Options{
		Addr:               addr,
		DialTimeout:        time.Duration(connectTime) * time.Millisecond,
		ReadTimeout:        time.Duration(readTimeout) * time.Millisecond,
		WriteTimeout:       time.Duration(writeTimeout) * time.Millisecond,
		IdleTimeout:        -1,
		IdleCheckFrequency: -1,
		PoolSize:           poolSize,
	})
	_, err := subjectl.Info().Result()
	return err
}

func InitIpInfoCluster(addrs []string, connectTime, readTimeout, writeTimeout, poolSize int) error {
	err := errors.New("InitIpInfoCluster start...")
	for i := 0; i < 10 && err != nil; i++ { // 增加10次重试
		ipInfoCluster = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:              addrs,
			DialTimeout:        time.Duration(connectTime) * time.Millisecond,
			ReadTimeout:        time.Duration(readTimeout) * time.Millisecond,
			WriteTimeout:       time.Duration(writeTimeout) * time.Millisecond,
			IdleTimeout:        -1,
			IdleCheckFrequency: -1,
			PoolSize:           poolSize,
		})
		_, err = ipInfoCluster.ClusterInfo().Result()
	}
	return err
}

func InitLocalAlgoRedis(addr string, connectTime, readTimeout, writeTimeout, poolSize int) error {
	subjectAlgo = redis.NewClient(&redis.Options{
		Addr:               addr,
		DialTimeout:        time.Duration(connectTime) * time.Millisecond,
		ReadTimeout:        time.Duration(readTimeout) * time.Millisecond,
		WriteTimeout:       time.Duration(writeTimeout) * time.Millisecond,
		IdleTimeout:        -1,
		IdleCheckFrequency: -1,
		PoolSize:           poolSize,
	})
	_, err := subjectAlgo.Info().Result()
	return err
}

// 算法使用的redis
func LocalRedisAlgoHGet(key string, field string) (val string, err error) {
	if subjectAlgoPool != nil {
		subjectAlgo = subjectAlgoPool.GetNode()
	}
	return subjectAlgo.HGet(key, field).Result()
}

func LocalRedisAlgoHExists(key, field string) (bool, error) {
	if subjectAlgoPool != nil {
		subjectAlgo = subjectAlgoPool.GetNode()
	}
	return subjectAlgo.HExists(key, field).Result()
}

func GetIpInfo(key string) (val string, err error) {
	return ipInfoCluster.Get(key).Result()
}

func SetIpInfo(key, val string, expireTime int64) error {
	scmd := ipInfoCluster.Set(key, val, time.Duration(expireTime)*time.Second)
	if scmd.Err() == redis.Nil {
		return redis.Nil
	}
	// close retry
	// if scmd.Err() == redis.Nil {
	//	scmd = ipInfoCluster.Set(key, val, time.Duration(expireTime)*time.Second)
	//	if scmd.Err() == redis.Nil {
	//		return redis.Nil
	//	}
	// }
	return scmd.Err()
}

func DelIpInfo(key string) (i int64, err error) {
	return ipInfoCluster.Del(key).Result()
}

// // set  key value 超时时间，失败则返回err
// func RedisSet(key string, value string, expireTime int64) error {
//	scmd := subject.Set(key, value, time.Duration(expireTime)*time.Second)
//	if scmd.Err() == redis.Nil {
//		scmd = subject.Set(key, value, time.Duration(expireTime)*time.Second)
//		if scmd.Err() == redis.Nil {
//			return redis.Nil
//		}
//	}
//
//	//_, err := subject.Set(key, value, time.Duration(expireTime*1000000000)).Result()
//	return scmd.Err()
// }
//
// // del 删除key，失败则返回err
// func RedisDel(key string) error {
//	return subject.Del(key).Err()
// }
//
// // get key，返回val，失败则返回err
// func RedisGet(key string) (val string, err error) {
//	return subject.Get(key).Result()
// }

func RedisSAdd(key string, members ...interface{}) error {
	return ipInfoCluster.SAdd(key, members...).Err()
}

func RedisExpire(key string, expiration time.Duration) (bool, error) {
	return ipInfoCluster.Expire(key, expiration).Result()
}

func RedisSMembers(key string) ([]string, error) {
	return ipInfoCluster.SMembers(key).Result()
}

// get key，返回val，失败则返回err
func LocalRedisGet(key string) (val string, err error) {
	now := time.Now().UnixNano()
	if mvutil.Config.AreaConfig.HttpConfig.UseCtRedisConsul {
		defer func() {
			watcher.AddAvgWatchValue("creative_consul_redis", float64((time.Now().UnixNano()-now)/1e3))
		}()
		subjectl = subjectlPool.GetNode()
	} else {
		defer func() {
			watcher.AddAvgWatchValue("creative_local_redis", float64((time.Now().UnixNano()-now)/1e3))
		}()
	}
	return subjectl.Get(key).Result()
}

// get DB1 key，返回val，失败则返回err
// func LocalRedisDB1HGet(key string, field string) (val string, err error) {
//	return subjectlDB1.HGet(key, field).Result()
// }

func LocalRedisHMGetPipeline(args map[string][]string) (map[string][]string, error) {
	now := time.Now().UnixNano()
	if mvutil.Config.AreaConfig.HttpConfig.UseCtRedisConsul {
		defer func() {
			watcher.AddAvgWatchValue("creative_consul_redis", float64((time.Now().UnixNano()-now)/1e3))
		}()
		subjectl = subjectlPool.GetNode()
	} else {
		defer func() {
			watcher.AddAvgWatchValue("creative_local_redis", float64((time.Now().UnixNano()-now)/1e3))
		}()
	}
	pipe := subjectl.Pipeline()
	result := make(map[string]*redis.SliceCmd, len(args))
	retLen := make(map[string]int, len(args))
	for k, v := range args {
		retLen[k] = len(v)
		result[k] = pipe.HMGet(k, v...)
	}
	_, err := pipe.Exec()
	if err != nil {
		return nil, err
	}
	pipe.Close()
	ret := make(map[string][]string, len(args))
	for k, v := range result {
		if v.Err() != nil {
			// ret[k] = nil
			continue
		}
		tmpRet := make([]string, retLen[k])
		for i, vv := range v.Val() {
			vstr := ""
			switch reply := vv.(type) {
			case []byte:
				vstr = string(reply)
			case string:
				vstr = reply
			case nil:
				vstr = ""
			}
			tmpRet[i] = vstr
		}
		ret[k] = tmpRet
	}
	return ret, nil
}

// hmget
func LocalRedisHMGet(key string, args []string) ([]string, error) {
	now := time.Now().UnixNano()
	if mvutil.Config.AreaConfig.HttpConfig.UseCtRedisConsul {
		defer func() {
			watcher.AddAvgWatchValue("creative_consul_redis", float64((time.Now().UnixNano()-now)/1e3))
		}()
		subjectl = subjectlPool.GetNode()
	} else {
		defer func() {
			watcher.AddAvgWatchValue("creative_local_redis", float64((time.Now().UnixNano()-now)/1e3))
		}()
	}
	pipe := subjectl.Pipeline()
	data := pipe.HMGet(key, args...)
	if data.Err() != nil {
		return nil, data.Err()
	}
	_, err := pipe.Exec()
	if err != nil {
		return nil, err
	}
	pipe.Close()
	result := make([]string, len(args))
	for i, v := range data.Val() {
		vstr := ""
		switch reply := v.(type) {
		case []byte:
			vstr = string(reply)
		case string:
			vstr = reply
		case nil:
			vstr = ""
		}
		result[i] = vstr
	}
	return result, nil
}

// // pipeLine
// func RedisPipelineGet(ipkey, device string) (string, int, error) {
//	pipe := subject.Pipeline()
//	ipData := pipe.Get(ipkey)
//	capData := pipe.Get(device)
//	_, errp := pipe.Exec()
//	pipe.Close()
//	if errp != nil {
//		return "", 0, errp
//	}
//	icap, _ := strconv.Atoi(capData.Val())
//	return ipData.Val(), icap, nil
// }

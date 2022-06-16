package storage

import (
	"time"

	"github.com/go-redis/redis"
	"github.com/spf13/viper"
)

type KVRedis struct {
	Conn *redis.Client
}

func NewRedisClient(configFile, datasource string) (*KVRedis, error) {
	client, err := factory(configFile, datasource)
	if err != nil {
		return nil, err
	}
	if _, err := client.Ping().Result(); err != nil {
		return nil, err
	}
	return &KVRedis{client}, nil
}

// readConfig resolve the config file.
// filePath: /a/b.toml, ./a/b.yaml, a.json
// fileName: local, support [toml, yaml, json]
func readConfig(filePath string) (*viper.Viper, error) {
	config := viper.New()
	config.SetConfigFile(filePath)
	if err := config.ReadInConfig(); err != nil {
		return nil, err
	}
	return config, nil
}

// redis client factory
func factory(configFile, datasource string) (*redis.Client, error) {
	config, err := readConfig(configFile)
	if err != nil {
		return nil, err
	}
	network := config.GetString("redis." + datasource + ".network")
	address := config.GetString("redis."+datasource+".host") + ":" + config.GetString("redis."+datasource+".port")
	auth := config.GetString("redis." + datasource + ".auth")

	readTimeout := config.GetDuration("redis.readTimeout") * time.Millisecond
	writeTimeout := config.GetDuration("redis.writeTimeout") * time.Millisecond
	dialTimeout := config.GetDuration("redis.dialTimeout") * time.Millisecond
	idleTimeout := config.GetDuration("redis.idleTimeout") * time.Minute
	database := config.GetInt("redis.database")
	poolSize := config.GetInt("redis.poolSize")

	client := redis.NewClient(&redis.Options{
		Network:      network,
		Addr:         address,
		Password:     auth,
		DB:           database,
		PoolSize:     poolSize,
		DialTimeout:  dialTimeout,
		IdleTimeout:  idleTimeout,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	})
	return client, nil
}

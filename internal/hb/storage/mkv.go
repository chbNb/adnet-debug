package storage

import (
	"errors"
	"time"

	"gitlab.mobvista.com/mtech/mkv/pkg/kvcfg"
	"gitlab.mobvista.com/mtech/mkv/pkg/kvclient"
)

type EngineType uint8

const (
	Redis     EngineType = iota + 1 // 1
	Aerospike                       // 2
)

type AerospikeClient struct {
	client     *kvclient.Aerospike
	compressor *ReqCtxCompressor
	serializer *ReqCtxSerializer
}

func NewAerospikeClient(client *kvclient.Aerospike) *AerospikeClient {
	c := AerospikeClient{client: client}
	return &c
}

type KVConfig struct {
	storageName EngineType
	configPath  string
	kvEngine    kvclient.KVClient
}

func NewKVStorage(storageName EngineType, configPath string) (*KVConfig, error) {
	cfg := &KVConfig{
		storageName: storageName,
		configPath:  configPath,
	}
	err := cfg.connector()
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

// Connector from kv config file
func (config *KVConfig) connector() error {
	client, err := kvcfg.NewKVClientWithFile(config.configPath)
	if err != nil {
		return err
	}
	config.kvEngine = client
	return nil
}

func (config *KVConfig) Close() error {
	return config.kvEngine.Close()
}

// MixSet is call Set or SetBidding
func (config *KVConfig) MixSet(key string, val interface{}) error {
	switch v := val.(type) {
	case []byte:
		return config.Set(key, v)
	case *BiddingVal:
		return config.SetBidding(key, v)
	case *ReqCtxVal:
		return config.SetReqCtx(key, v)
	default:
		return errors.New("value type undefined")
	}
}

// Set set string data to kv storage by the key
func (config *KVConfig) Set(key string, val []byte) error {
	var err error
	k := &SimpleKey{Message: key}
	v := &SimpleVal{Message: val}
	config.kvEngine.SetCompressor(&SimpleCompressor{})
	config.kvEngine.SetSerializer(&SimpleSerializer{})
	switch config.storageName {
	case Redis, Aerospike:
		err = config.kvEngine.Set(k, v)
	default:
		err = errors.New("kv engine unimplemented.")
	}
	if err != nil {
		return err
	}
	return nil
}

// SetBidding set BiddingVal data to kv storage by the key
func (config *KVConfig) SetBidding(key string, val *BiddingVal) error {
	var err error
	k := &BiddingKey{BiddingID: key}
	compressor := NewBiddingCompressor()
	config.kvEngine.SetCompressor(compressor)
	serializer := NewBiddingSerializer()
	config.kvEngine.SetSerializer(serializer)
	switch config.storageName {
	case Redis, Aerospike:
		err = config.kvEngine.Set(k, val)
	default:
		err = errors.New("kv engine unimplemented.")
	}
	if err != nil {
		return err
	}
	return nil
}

func (config *KVConfig) SetCompressor(c kvclient.Compressor) {
	config.kvEngine.SetCompressor(c)
}

func (config *KVConfig) SetSerializer(s kvclient.Serializer) {
	config.kvEngine.SetSerializer(s)
}

// SetReqCtx set ReqCtx data to kv storage by the key
func (config *KVConfig) SetReqCtx(key string, val *ReqCtxVal) error {
	var err error
	k := &ReqCtxKey{Token: key}
	//compressor := NewReqCtxCompressor()
	//config.kvEngine.SetCompressor(compressor)
	//serializer := NewReqCtxSerializer()
	//config.kvEngine.SetSerializer(serializer)
	switch config.storageName {
	case Redis, Aerospike:
		err = config.kvEngine.Set(k, val)
	default:
		err = errors.New("kv engine unimplemented.")
	}
	if err != nil {
		return err
	}
	return nil
}

// SetExReqCtx set ReqCtx data to kv storage by the key and expiration time
func (config *KVConfig) SetExReqCtx(key string, val *ReqCtxVal, exp time.Duration) error {
	var err error
	k := &ReqCtxKey{Token: key}
	compressor := NewReqCtxCompressor()
	config.kvEngine.SetCompressor(compressor)
	serializer := NewReqCtxSerializer()
	config.kvEngine.SetSerializer(serializer)
	switch config.storageName {
	case Redis, Aerospike:
		err = config.kvEngine.SetEx(k, val, exp)
	default:
		err = errors.New("kv engine unimplemented.")
	}
	if err != nil {
		return err
	}
	return nil
}

func (config *KVConfig) Delete(key string) error {
	var err error
	k := &ReqCtxKey{Token: key}
	compressor := NewReqCtxCompressor()
	config.kvEngine.SetCompressor(compressor)
	serializer := NewReqCtxSerializer()
	config.kvEngine.SetSerializer(serializer)
	switch config.storageName {
	case Redis, Aerospike:
		err = config.kvEngine.Del(k)
	default:
		err = errors.New("kv engine unimplemented.")
	}
	if err != nil {
		return err
	}
	return nil
}

func (c *AerospikeClient) SetCompressor() {
	c.compressor = NewReqCtxCompressor()
}

func (c *AerospikeClient) SetSerializer() {
	c.serializer = NewReqCtxSerializer()
}

func (c *AerospikeClient) Set(key, val interface{}) error {
	if c.client == nil {
		return errors.New("aerospike client did not complete initialization")
	}
	//compressor := NewReqCtxCompressor()
	k := c.compressor.Compress(key)
	//serializer := NewReqCtxSerializer()
	v, _ := c.serializer.Marshal(val)
	err := c.client.Set(k, v)
	if err != nil {
		return err
	}
	return nil
}

func (c *AerospikeClient) GetReqCtx(key interface{}) (*ReqCtxVal, error) {
	if c.client == nil {
		return nil, errors.New("aerospike client did not complete initialization")
	}
	compressor := NewReqCtxCompressor()
	k := compressor.Compress(key)
	serializer := NewReqCtxSerializer()
	b, err := c.client.Get(k)
	if err != nil {
		return nil, err
	}
	var res ReqCtxVal
	serializer.Unmarshal(b, &res)
	return &res, nil
}

func (c *AerospikeClient) SetExReqCtx(key, val interface{}, exp time.Duration) error {
	if c.client == nil {
		return errors.New("aerospike client did not complete initialization")
	}
	compressor := NewReqCtxCompressor()
	k := compressor.Compress(key)
	serializer := NewReqCtxSerializer()
	v, _ := serializer.Marshal(val)
	err := c.client.SetEx(k, v, exp)
	if err != nil {
		return err
	}
	return nil
}

func (c *AerospikeClient) Delete(key interface{}) error {
	if c.client == nil {
		return errors.New("aerospike client did not complete initialization")
	}
	compressor := NewReqCtxCompressor()
	k := compressor.Compress(key)
	return c.client.Del(k)
}

// SetEx set string data to kv storage by the key and expiration time
func (config *KVConfig) SetEx(key string, val []byte, exp time.Duration) error {
	var err error
	k := &SimpleKey{Message: key}
	v := &SimpleVal{Message: val}
	config.kvEngine.SetCompressor(&SimpleCompressor{})
	config.kvEngine.SetSerializer(&SimpleSerializer{})
	switch config.storageName {
	case Redis, Aerospike:
		err = config.kvEngine.SetEx(k, v, exp)
	default:
		err = errors.New("kv engine unimplemented.")
	}
	if err != nil {
		return err
	}
	return nil
}

// SetEx set BiddingVal data to kv storage by the key and expiration time
func (config *KVConfig) SetExBidding(key string, val *BiddingVal, exp time.Duration) error {
	var err error
	k := &BiddingKey{BiddingID: key}
	compressor := NewBiddingCompressor()
	config.kvEngine.SetCompressor(compressor)
	serializer := NewBiddingSerializer()
	config.kvEngine.SetSerializer(serializer)
	switch config.storageName {
	case Redis, Aerospike:
		err = config.kvEngine.SetEx(k, val, exp)
	default:
		err = errors.New("kv engine unimplemented.")
	}
	if err != nil {
		return err
	}
	return nil
}

// Get string data from kv storage by the key
func (config *KVConfig) Get(key string) ([]byte, error) {
	var (
		err error
		ok  bool
		v   SimpleVal
	)
	config.kvEngine.SetCompressor(&SimpleCompressor{})
	config.kvEngine.SetSerializer(&SimpleSerializer{})
	k := &SimpleKey{Message: key}
	switch config.storageName {
	case Redis, Aerospike:
		ok, err = config.kvEngine.Get(k, &v)
	default:
		err = errors.New("kv engine unimplemented.")
	}
	if err != nil {
		return nil, err
	}
	if ok {
		return v.Message, nil
	}
	return nil, errors.New("the key is missing in kv")
}

// GetBidding get string data from kv storage by the key
func (config *KVConfig) GetBidding(key string) (*BiddingVal, error) {
	var (
		err error
		ok  bool
		v   BiddingVal
	)
	compressor := NewBiddingCompressor()
	config.kvEngine.SetCompressor(compressor)
	serializer := NewBiddingSerializer()
	config.kvEngine.SetSerializer(serializer)
	k := &BiddingKey{BiddingID: key}
	switch config.storageName {
	case Redis, Aerospike:
		ok, err = config.kvEngine.Get(k, &v)
	default:
		err = errors.New("kv engine unimplemented.")
	}
	if err != nil {
		return nil, err
	}
	if ok {
		return &v, nil
	}
	return nil, errors.New("the key is missing in kv")
}

// GetReqCtx get string data from kv storage by the key
func (config *KVConfig) GetReqCtx(key string) (*ReqCtxVal, error) {
	var (
		err error
		ok  bool
		v   ReqCtxVal
	)
	compressor := NewReqCtxCompressor()
	config.kvEngine.SetCompressor(compressor)
	serializer := NewReqCtxSerializer()
	config.kvEngine.SetSerializer(serializer)
	k := &ReqCtxKey{Token: key}
	switch config.storageName {
	case Redis, Aerospike:
		ok, err = config.kvEngine.Get(k, &v)
	default:
		err = errors.New("kv engine unimplemented.")
	}
	if err != nil {
		return nil, err
	}
	if ok {
		return &v, nil
	}
	return nil, errors.New("the key is missing in kv")
}

func (config *KVConfig) GetASFieldVal(key string) (map[string][]byte, error) {
	val, err := config.kvEngine.GetASFieldVal(key)
	if err != nil {
		return nil, err
	}
	return val, nil
}

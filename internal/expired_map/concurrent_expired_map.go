package expired_map

import (
	"github.com/easierway/concurrent_map"
	"sync"
	"sync/atomic"
	"time"
)

type ConcurrentExpiredMap struct {
	partitions    []*ExpiredMap
	numOfBlockets int

	//过期删除的周期
	expiredDefaultTime int64 //默认通用时间，只用使用Set方法时使用

	tickerTime time.Duration
	stop       chan struct{}
	needStop   int32
}

func createExpiredMap(batchDeleteTime time.Duration) *ExpiredMap {

	e := ExpiredMap{
		dataMap:    make(map[interface{}]*dataStruct),
		lck:        new(sync.RWMutex),
		tickerTime: batchDeleteTime,
		timeMap:    make(map[int64]map[interface{}]bool),
		stop:       make(chan struct{}),
	}
	atomic.StoreInt32(&e.needStop, 1)
	//go e.run()  //不需要单独维护key的删除
	return &e
}

// CreateConcurrentExpiredMap is to create a ConcurrentExpiredMap with the setting number of the partitions
func CreateConcurrentExpiredMap(numOfPartitions int, batchDeleteTime time.Duration, expiredDefaultTime int64) *ConcurrentExpiredMap {
	var partitions []*ExpiredMap
	for i := 0; i < numOfPartitions; i++ {
		partitions = append(partitions, createExpiredMap(batchDeleteTime))
	}
	cem := ConcurrentExpiredMap{
		partitions:    partitions,
		numOfBlockets: numOfPartitions,

		expiredDefaultTime: expiredDefaultTime,

		tickerTime: batchDeleteTime,
		stop:       make(chan struct{}),
	}

	go cem.run()
	return &cem
}

func (cem *ConcurrentExpiredMap) getPartition(key concurrent_map.Partitionable) *ExpiredMap {
	partitionID := key.PartitionKey() % int64(cem.numOfBlockets)
	return (*ExpiredMap)(cem.partitions[partitionID])
}

// Get is to get the value by the key
func (cem *ConcurrentExpiredMap) Get(key concurrent_map.Partitionable) (interface{}, bool) {
	return cem.getPartition(key).Get(key.Value())
}

// Set is to store the KV entry to the map
func (cem *ConcurrentExpiredMap) Set(key concurrent_map.Partitionable, v interface{}) {
	im := cem.getPartition(key)
	im.Set(key.Value(), v, cem.expiredDefaultTime)
}

func (cem *ConcurrentExpiredMap) SetWithTime(key concurrent_map.Partitionable, v interface{}, expireSeconds int64) {
	im := cem.getPartition(key)
	im.Set(key.Value(), v, expireSeconds)
}

//缓存时间少于一半时刷新缓存
func (cem *ConcurrentExpiredMap) ReSetAfterHalfexpire(key concurrent_map.Partitionable, v interface{}, expireSeconds int64) {
	im := cem.getPartition(key)
	ttl := im.TTL(key.Value())
	if ttl < expireSeconds/2 {
		im.Set(key.Value(), v, expireSeconds)
	}
}

// Del is to delete the entries by the key
func (cem *ConcurrentExpiredMap) Del(key concurrent_map.Partitionable) {
	im := cem.getPartition(key)
	im.Delete(key.Value())
}

//

//background goroutine 主动删除过期的key
//数据实际删除时间比应该删除的时间稍晚一些，这个误差我们应该能接受。
func (cem *ConcurrentExpiredMap) run() {
	t := time.NewTicker(time.Second * cem.tickerTime * 1)
	delCh := make(chan *delMsg, delChannelCap)
	times := make([]int64, cem.tickerTime+3)
	now := time.Now().Unix()
	go func() {
		for v := range delCh {
			if atomic.LoadInt32(&cem.needStop) == 1 {
				//fmt.Println("---del stop---")
				return
			}
			cem.deleteByTimes(v.times)
		}
	}()
	for {
		select {
		case <-t.C:
			now = time.Now().Unix() - 1 //当前这一秒不删除
			for k := range times {
				times[k] = now
				now--
			}
			delCh <- &delMsg{times}
		case <-cem.stop:
			//fmt.Println("=== STOP ===")
			atomic.StoreInt32(&cem.needStop, 1)
			delCh <- &delMsg{times: []int64{}}
			return
		}
	}
}

//根据时间来删除
func (cem *ConcurrentExpiredMap) deleteByTimes(times []int64) {
	for _, e := range cem.partitions { //分片处理
		e.deleteByTimes(times)
	}
}

func (cem *ConcurrentExpiredMap) Close() { // todo 关闭后在使用怎么处理
	cem.stop <- struct{}{}
	cem.partitions = nil
	cem.numOfBlockets = 0
}

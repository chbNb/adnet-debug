package expired_map

import (
	"sync"
	"sync/atomic"
	"time"
)

const delChannelCap = 100

type dataStruct struct {
	realData    interface{} //真实的数据
	expiredTime int64       //过期时间，时间戳
}

type ExpiredMap struct {
	dataMap map[interface{}]*dataStruct
	timeMap map[int64]map[interface{}]bool //{时间戳: {key:  true}}

	tickerTime time.Duration //过期key的删除周期
	lck        *sync.RWMutex
	stop       chan struct{}
	needStop   int32
}

type delMsg struct {
	times []int64
}

func CreateExpiredMap(batchDeleteTime time.Duration) *ExpiredMap {
	e := ExpiredMap{
		dataMap:    make(map[interface{}]*dataStruct),
		lck:        new(sync.RWMutex),
		tickerTime: batchDeleteTime,
		timeMap:    make(map[int64]map[interface{}]bool),
		stop:       make(chan struct{}),
	}
	atomic.StoreInt32(&e.needStop, 0)
	go e.run()
	return &e
}

//获取key的删除周期
func (e *ExpiredMap) GetTickerTime() time.Duration {
	return e.tickerTime
}

//background goroutine 主动删除过期的key
//数据实际删除时间比应该删除的时间稍晚一些，这个误差我们应该能接受。
func (e *ExpiredMap) run() {
	t := time.NewTicker(time.Second * e.tickerTime * 1)
	delCh := make(chan *delMsg, delChannelCap)
	times := make([]int64, e.tickerTime+3)
	now := time.Now().Unix()
	go func() {
		for v := range delCh {
			if atomic.LoadInt32(&e.needStop) == 1 {
				return
			}
			e.deleteByTimes(v.times)
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
		case <-e.stop:
			atomic.StoreInt32(&e.needStop, 1)
			delCh <- &delMsg{times: []int64{}}
			return
		}
	}
}

func (e *ExpiredMap) Set(key, value interface{}, expireSeconds int64) {
	if expireSeconds <= 0 {
		return
	}
	expiredTime := time.Now().Unix() + expireSeconds
	e.lck.Lock()
	defer e.lck.Unlock()

	val, ok := e.dataMap[key]
	if ok { //可以找到历史数据，直接删除
		needToDeleteKey := val.expiredTime
		if _, ok := e.timeMap[needToDeleteKey]; ok {
			delete(e.timeMap[needToDeleteKey], key)
		}
	}

	e.dataMap[key] = &dataStruct{
		realData:    value,
		expiredTime: expiredTime,
	}

	if e.timeMap[expiredTime] == nil {
		e.timeMap[expiredTime] = make(map[interface{}]bool)
	}
	e.timeMap[expiredTime][key] = true //过期时间作为key，放在map中
}

func (e *ExpiredMap) Get(key interface{}) (value interface{}, found bool) {
	e.lck.RLock()
	defer e.lck.RUnlock()

	val, ok := e.dataMap[key]
	if ok {
		if !e.checkExpired(val) { //已过期
			ok = false
		}
	}
	if !ok {
		return
	}

	return val.realData, true
}

func (e *ExpiredMap) Delete(key interface{}) {
	e.lck.Lock()
	defer e.lck.Unlock()
	if old, ok := e.dataMap[key]; ok {
		if _, ok := e.timeMap[old.expiredTime]; ok {
			delete(e.timeMap[old.expiredTime], key)
		}
		delete(e.dataMap, key)
	}
}

func (e *ExpiredMap) Remove(key interface{}) {
	e.Delete(key)
}

//根据时间来删除
func (e *ExpiredMap) deleteByTimes(times []int64) {
	e.lck.Lock()
	defer e.lck.Unlock()
	for _, t := range times {
		if keys, ok := e.timeMap[t]; ok {
			for key := range keys {
				delete(e.dataMap, key) //删除实际数据
			}
			delete(e.timeMap, t) //删除时间中的key
		}
	}
}

func (e *ExpiredMap) Length() int { //结果是不准确的，因为有未删除的key，用于统计当前实时内存保存的条数
	e.lck.RLock()
	defer e.lck.RUnlock()
	return len(e.dataMap)
}

//返回key的剩余生存时间 key不存在返回负数
func (e *ExpiredMap) TTL(key interface{}) int64 {
	e.lck.RLock()
	defer e.lck.RUnlock()

	val, found := e.dataMap[key]
	if found {
		if !e.checkExpired(val) { //过期
			return -1
		}
	} else { //找不到
		return -1
	}

	return e.dataMap[key].expiredTime - time.Now().Unix()
}

//检测是否过期， 最大会有1s差值，也就是设置了1s，实际上是2s才过期
func (e *ExpiredMap) checkExpired(val *dataStruct) bool {
	if val.expiredTime < time.Now().Unix() {
		return false
	}
	return true
}

//====================意义不大的方法===============================
func (e *ExpiredMap) Clear() {
	e.lck.Lock()
	defer e.lck.Unlock()
	e.dataMap = make(map[interface{}]*dataStruct)
	e.timeMap = make(map[int64]map[interface{}]bool)
}

func (e *ExpiredMap) Close() { // todo 关闭后在使用怎么处理
	e.lck.Lock()
	defer e.lck.Unlock()
	e.stop <- struct{}{}
	e.dataMap = nil
	e.timeMap = nil
}

func (e *ExpiredMap) Stop() {
	e.Close()
}

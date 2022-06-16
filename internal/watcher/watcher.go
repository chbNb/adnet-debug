package watcher

import (
	"bytes"
	"runtime"
	"sort"
	"strconv"
	"sync"
	"time"

	mlogger "github.com/mae-pax/logger"
)

type AvgWatch struct {
	count float64
	total float64
}

type watch struct {
	stopChan    chan bool
	mutex       *sync.Mutex
	avg_mutex   *sync.Mutex
	content     map[string]float64
	avg_content map[string]*AvgWatch
}

var watcher watch
var logger *mlogger.Log

func AddWatchValue(key string, value float64) {
	watcher.mutex.Lock()
	defer watcher.mutex.Unlock()
	v, in := watcher.content[key]
	if in {
		v = v + value
	} else {
		v = value
	}

	watcher.content[key] = v
}

func SetWatchValue(key string, value float64) {
	watcher.mutex.Lock()
	defer watcher.mutex.Unlock()
	watcher.content[key] = value
}

func AddAvgWatchValue(key string, value float64) {
	watcher.avg_mutex.Lock()
	defer watcher.avg_mutex.Unlock()
	v, in := watcher.avg_content[key]
	if !in {
		v = &AvgWatch{}
		watcher.avg_content[key] = v
	}

	v.total += value
	v.count++
}

func RunWatch() {
	go func() {
		for {
			select {
			case <-time.Tick(1 * time.Minute):
				WriteDisk()
			case <-watcher.stopChan:
				logger.Info("Watcher will be stop.")
				return
			}
		}
	}()
}

func Stop() {
	watcher.stopChan <- true
	WriteDisk()
}

func WriteDisk() {
	internalStr := getInternal()
	content := getContent()
	avg_content := getAvgContent()
	logger.Info(internalStr + content + avg_content)
	// logger.Out.(*bufio.Writer).Flush()
}

func getInternal() string {
	var buffer bytes.Buffer
	// goroutines
	// buffer.WriteString("")
	buffer.WriteString("\tgoroutines:")
	// buffer.WriteString(":")
	buffer.WriteString(strconv.FormatFloat(float64(runtime.NumGoroutine()), 'f', -1, 64))
	// local redis pool
	// pool := redis.GetRedisPool()
	// if pool != nil {
	// 	poolStats := pool.Stats()
	// 	buffer.WriteString("\tlredis_active_count:")
	// 	buffer.WriteString(strconv.Itoa(poolStats.ActiveCount))
	// 	buffer.WriteString("\tlredis_idle_count:")
	// 	buffer.WriteString(strconv.Itoa(poolStats.IdleCount))
	// }
	return buffer.String()
}

func getContent() string {
	var buffer bytes.Buffer
	watcher.mutex.Lock()
	defer watcher.mutex.Unlock()

	keys := make([]string, len(watcher.content))
	i := 0
	for k := range watcher.content {
		keys[i] = k
		i++
	}
	sort.Strings(keys)

	for _, key := range keys {
		buffer.WriteString("\t")
		buffer.WriteString(key)
		buffer.WriteString(":")
		buffer.WriteString(strconv.FormatFloat(watcher.content[key], 'f', -1, 64))
		watcher.content[key] = 0.0
	}

	return buffer.String()
}

func getAvgContent() string {
	var buffer bytes.Buffer
	watcher.avg_mutex.Lock()
	defer watcher.avg_mutex.Unlock()
	// avg
	keys := make([]string, len(watcher.avg_content))
	i := 0
	for k := range watcher.avg_content {
		keys[i] = k
		i++
	}
	sort.Strings(keys)

	for _, key := range keys {
		avg := 0
		if watcher.avg_content[key].count > 0 {
			avg = int(watcher.avg_content[key].total / watcher.avg_content[key].count)
		}
		buffer.WriteString("\t")
		buffer.WriteString(key)
		buffer.WriteString(":")
		buffer.WriteString(strconv.Itoa(avg))
		buffer.WriteString("\t")
		buffer.WriteString(key + "_count")
		buffer.WriteString(":")
		buffer.WriteString(strconv.FormatFloat(watcher.avg_content[key].count, 'f', -1, 64))
		watcher.avg_content[key].count = 0.0
		watcher.avg_content[key].total = 0.0
	}

	return buffer.String()
}

func Init(log *mlogger.Log) {
	logger = log
	watcher.stopChan = make(chan bool)
	watcher.mutex = new(sync.Mutex)
	watcher.avg_mutex = new(sync.Mutex)
	watcher.content = make(map[string]float64)
	watcher.avg_content = make(map[string]*AvgWatch)
}

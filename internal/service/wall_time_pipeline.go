package service

import (
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	"github.com/easierway/pipefiter_framework/pipefilter"
	"gitlab.mobvista.com/ADN/exporter/metrics"
)

type AtomicInt int64

func (a *AtomicInt) Add(i int64) {
	atomic.AddInt64((*int64)(a), i)
}

func (a *AtomicInt) Val() int64 {
	return *(*int64)(a)
}

func (a *AtomicInt) Reset() int64 {
	return atomic.SwapInt64((*int64)(a), 0)
}

type WallTimePipeline struct {
	Name        string
	Filters     *[]pipefilter.Filter
	TimeElapsed []AtomicInt
	Labels      []string // pipeline metrics采集需要
}

func (wtp *WallTimePipeline) Process(data interface{}) (interface{}, error) {
	var ret interface{}
	var err error

	var isMetrics bool
	if len(wtp.Labels) != 0 {
		isMetrics = true
	}

	now := time.Now()
	for i, filter := range *wtp.Filters {
		ret, err = filter.Process(data)
		wtp.TimeElapsed[i].Add(time.Since(now).Nanoseconds())
		if isMetrics {
			costTime := (float64)(time.Since(now) / time.Millisecond)   // 毫秒
			metrics.SetGaugeWithLabelValues(costTime, 5, wtp.Labels[i]) // pipeline处理时间metrics采集
		}

		now = time.Now()
		if err != nil {
			return ret, err
		}
		data = ret
	}
	return ret, err
}

func (wtp *WallTimePipeline) Show(path string) {
	go func() {
		for range time.Tick(time.Second * time.Duration(60)) {
			ts := make([]string, len(wtp.TimeElapsed))
			for i := range wtp.TimeElapsed {
				ts[i] = fmt.Sprintf("%v", wtp.TimeElapsed[i].Reset()/(int64)(time.Millisecond))
			}
			fmt.Println(time.Now().Format("2006-01-02 15:04:05") + " | " + path + " | " + strings.Join(ts, "\t"))
		}
	}()
}

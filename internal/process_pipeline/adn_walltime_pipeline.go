package process_pipeline

import (
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	"github.com/easierway/pipefiter_framework/pipefilter"
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

type AdnWallTimePipeline struct {
	Name        string
	Filters     *[]pipefilter.Filter
	TimeElapsed []AtomicInt
}

func (awtp *AdnWallTimePipeline) Process(data interface{}) (interface{}, error) {
	var ret interface{}
	var err error

	now := time.Now()
	for i, filter := range *awtp.Filters {
		ret, err = filter.Process(data)
		awtp.TimeElapsed[i].Add(time.Since(now).Nanoseconds())
		now = time.Now()
		if err != nil {
			return ret, err
		}
		data = ret
	}
	return ret, err
}

func (awtp *AdnWallTimePipeline) Show() {
	go func() {
		for range time.Tick(time.Second * time.Duration(30)) {
			ts := make([]string, len(awtp.TimeElapsed))
			for i := range awtp.TimeElapsed {
				ts[i] = fmt.Sprintf("%v", awtp.TimeElapsed[i].Reset()/(int64)(time.Millisecond))
			}
			fmt.Println(time.Now().Format("2006-01-02 15:04:05") + " | " + strings.Join(ts, "\t"))
		}
	}()
}

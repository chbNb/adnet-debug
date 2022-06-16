package process_pipeline

import (
	"io"
	"net/http"
	"time"

	"github.com/easierway/pipefiter_framework/pipefilter"
	"gitlab.mobvista.com/ADN/adnet/internal/extractor"
	"gitlab.mobvista.com/ADN/adnet/internal/mvutil"
	"gitlab.mobvista.com/ADN/adnet/internal/utility"
	"gitlab.mobvista.com/ADN/adnet/internal/watcher"
)

//var Bucket chan time.Time
var RateLimitDecorator *utility.RateLimitDecorator
var bucketSpeed int
var preTokens int
var stopChan chan bool

type AdnHandler struct {
	Pipeline pipefilter.Filter
}

func UpdateBucket(bucketSpeed, preTokens int) {
	bucketSpeed = bucketSpeed
	preTokens = preTokens
	// decorator, err := service_decorators.CreateRateLimitDecorator(1*time.Second, bucketSpeed, preTokens)
	// if err != nil {
	// 	return err
	// }
	// RateLimitDecorator = decorator
	//Bucket = microservice_helper.CreateTokenBucket(preTokens, bucketSpeed, 1*time.Second)
	go func() {
		inc := time.After(2 * time.Minute)
		for {
			select {
			case <-inc:
				rateLimit := extractor.GetRateLimit()
				if rateLimit != nil && rateLimit.BucketSpeed > 0 && rateLimit.PreTokens > 0 &&
					(rateLimit.BucketSpeed != bucketSpeed || rateLimit.PreTokens != preTokens) {
					bucketSpeed = rateLimit.BucketSpeed
					preTokens = rateLimit.PreTokens
					decorator, err := utility.CreateRateLimitDecorator(1*time.Second, bucketSpeed, preTokens)
					if err != nil {
						mvutil.Logger.Runtime.Infof("new RateLimit config PreTokens:%d, BucketSpeed:%d error: %s", preTokens, bucketSpeed, err.Error())
					} else {
						//Bucket = microservice_helper.CreateTokenBucket(preTokens, bucketSpeed, 1*time.Second)
						RateLimitDecorator = decorator
						mvutil.Logger.Runtime.Infof("new RateLimit config PreTokens:%d, BucketSpeed:%d", preTokens, bucketSpeed)
					}
				}
				//getOtherConfig()
				inc = time.After(2 * time.Minute)
			case <-stopChan:
				mvutil.Logger.Runtime.Infof("UpdateBucket get stop signal, will stop")
				return
			}
		}
	}()
}

func init() {
	stopChan = make(chan bool)
}

func StopBucket() {
	stopChan <- true
}

func WriterData(now int64, c http.ResponseWriter, data string) {
	io.WriteString(c, data)
	watcher.AddAvgWatchValue("adnet_cost", float64((time.Now().UnixNano()-now)/1e6))
}

func WriterHeader(now int64, c http.ResponseWriter, code int) {
	c.WriteHeader(code)
	watcher.AddAvgWatchValue("adnet_cost", float64((time.Now().UnixNano()-now)/1e6))
}

func LimitError(now int64, c http.ResponseWriter, req *http.Request) {
	watcher.AddWatchValue("ratelimit", float64(1))
	WriterHeader(now, c, http.StatusTooManyRequests)
}

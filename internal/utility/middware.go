package utility

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	"gitlab.mobvista.com/ADN/exporter/metrics"
	"gitlab.mobvista.com/ADN/exporter/utility"
)

func HttpHandlerInterceptor(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timeStart := time.Now()
		httpMethod := r.Method

		// 防止转义合法的url被过滤
		var apiPath string
		u, _ := url.Parse(r.URL.String())
		paths := strings.Split(u.Path, "?")
		if len(paths) != 0 {
			apiPath = paths[0]
		}

		next.ServeHTTP(w, r)

		timeElapsed := (float64)(time.Since(timeStart) / time.Millisecond)
		if !utility.RegexpInvalidUrlPath(apiPath) && !utility.CheckBlackKeyPath(apiPath) {
			metrics.IncCounterWithLabelValues(0, httpMethod, apiPath)              // 自增qps
			metrics.SetGaugeWithLabelValues(timeElapsed, 1, httpMethod, apiPath)   // 响应时间
			metrics.AddSummaryWithLabelValues(timeElapsed, 2, httpMethod, apiPath) // 响应分位
		}
	})
}

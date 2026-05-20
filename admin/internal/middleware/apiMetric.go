package middleware

import (
	"net/http"
	"time"

	"myproject/admin/internal/metric"
)

// ApiMetricMiddleware 记录 API 请求指标
func ApiMetricMiddleware() func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// 创建响应记录器以捕获状态码
			recorder := &metricResponseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			next(recorder, r)

			// 记录请求耗时
			duration := time.Since(start)
			metric.ObserveApiRequest(r.Method, r.URL.Path, duration)

			// 记录请求计数
			status := "success"
			if recorder.statusCode >= 400 {
				status = "error"
			}
			metric.ApiRequestCounter.Inc(r.Method, r.URL.Path, status)
		}
	}
}

// metricResponseWriter 捕获响应状态码
type metricResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *metricResponseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

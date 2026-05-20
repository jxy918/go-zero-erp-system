package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"
)

// RateLimiter 简单的内存速率限制器
type RateLimiter struct {
	mu       sync.Mutex
	requests map[string][]time.Time
	maxReqs  int
	window   time.Duration
}

// NewRateLimiter 创建速率限制器
// maxReqs: 时间窗口内最大请求数
// window: 时间窗口大小
func NewRateLimiter(maxReqs int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]time.Time),
		maxReqs:  maxReqs,
		window:   window,
	}
}

// Allow 检查是否允许请求
func (rl *RateLimiter) Allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-rl.window)

	// 清理过期记录
	var valid []time.Time
	for _, t := range rl.requests[key] {
		if t.After(windowStart) {
			valid = append(valid, t)
		}
	}

	// 检查是否超过限制
	if len(valid) >= rl.maxReqs {
		rl.requests[key] = valid
		return false
	}

	// 添加新请求
	rl.requests[key] = append(valid, now)
	return true
}

// 全局登录速率限制器：5次/分钟
var loginRateLimiter = NewRateLimiter(5, time.Minute)

// LoginRateLimitMiddleware 登录接口速率限制中间件
func LoginRateLimitMiddleware() func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// 仅对登录接口生效
			if r.URL.Path != "/auth/login" {
				next(w, r)
				return
			}

			// 获取客户端 IP
			ip := getClientIP(r)

			// 检查速率限制
			if !loginRateLimiter.Allow(ip) {
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				w.WriteHeader(http.StatusTooManyRequests)
				w.Write([]byte(`{"code":429,"data":null,"message":"登录请求过于频繁，请稍后再试"}`))
				return
			}

			next(w, r)
		}
	}
}

// getClientIP 获取客户端真实 IP
func getClientIP(r *http.Request) string {
	// 优先从代理头获取
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return xff
	}
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

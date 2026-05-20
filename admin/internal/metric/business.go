package metric

import (
	"sync"
	"time"

	"github.com/zeromicro/go-zero/core/metric"
)

var (
	enabled   = true
	enabledMu sync.RWMutex
)

// SetEnabled 设置指标收集开关
func SetEnabled(on bool) {
	enabledMu.Lock()
	defer enabledMu.Unlock()
	enabled = on
}

// IsEnabled 检查指标收集是否启用
func IsEnabled() bool {
	enabledMu.RLock()
	defer enabledMu.RUnlock()
	return enabled
}

// 业务指标收集器
// 所有指标在 prometheus.Enabled() 为 true 时自动收集
// DevServer 启用后自动启用 prometheus

var (
	// 登录指标
	LoginCounter = metric.NewCounterVec(&metric.CounterVecOpts{
		Namespace: "erp",
		Subsystem: "auth",
		Name:      "login_total",
		Help:      "Total login attempts",
		Labels:    []string{"status"}, // success, failure
	})

	// API 请求指标
	ApiRequestCounter = metric.NewCounterVec(&metric.CounterVecOpts{
		Namespace: "erp",
		Subsystem: "api",
		Name:      "request_total",
		Help:      "Total API requests",
		Labels:    []string{"method", "path", "status"}, // GET/POST, /user/list, 200/400/500
	})

	// API 请求耗时（毫秒）
	ApiRequestDuration = metric.NewHistogramVec(&metric.HistogramVecOpts{
		Namespace: "erp",
		Subsystem: "api",
		Name:      "request_duration_ms",
		Help:      "API request duration in milliseconds",
		Labels:    []string{"method", "path"},
		Buckets:   []float64{10, 50, 100, 200, 500, 1000, 2000, 5000},
	})

	// 数据库查询耗时（毫秒）
	DbQueryDuration = metric.NewHistogramVec(&metric.HistogramVecOpts{
		Namespace: "erp",
		Subsystem: "db",
		Name:      "query_duration_ms",
		Help:      "Database query duration in milliseconds",
		Labels:    []string{"table", "operation"}, // users, select/update/insert
		Buckets:   []float64{5, 10, 20, 50, 100, 200, 500, 1000},
	})

	// 订单创建指标
	OrderCreateCounter = metric.NewCounterVec(&metric.CounterVecOpts{
		Namespace: "erp",
		Subsystem: "order",
		Name:      "create_total",
		Help:      "Total order creations",
		Labels:    []string{"type", "status"}, // purchase/sales, success/failure
	})

	// 库存调整指标
	InventoryAdjustCounter = metric.NewCounterVec(&metric.CounterVecOpts{
		Namespace: "erp",
		Subsystem: "inventory",
		Name:      "adjust_total",
		Help:      "Total inventory adjustments",
		Labels:    []string{"type"}, // inbound/outbound/adjust
	})

	// 活跃用户数（Gauge）
	ActiveUsersGauge = metric.NewGaugeVec(&metric.GaugeVecOpts{
		Namespace: "erp",
		Subsystem: "system",
		Name:      "active_users",
		Help:      "Number of active users",
		Labels:    []string{"status"}, // enabled/disabled
	})
)

// ObserveApiRequest 记录 API 请求耗时
func ObserveApiRequest(method, path string, duration time.Duration) {
	if !IsEnabled() {
		return
	}
	ms := duration.Milliseconds()
	ApiRequestDuration.Observe(ms, method, path)
}

// ObserveDbQuery 记录数据库查询耗时
func ObserveDbQuery(table, operation string, duration time.Duration) {
	if !IsEnabled() {
		return
	}
	ms := duration.Milliseconds()
	DbQueryDuration.Observe(ms, table, operation)
}

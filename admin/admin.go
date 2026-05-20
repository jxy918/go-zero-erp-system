// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"strings"

	"myproject/admin/internal/config"
	"myproject/admin/internal/handler"
	"myproject/admin/internal/metric"
	"myproject/admin/internal/middleware"
	"myproject/admin/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/admin-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)

	// 设置自定义业务指标开关
	metric.SetEnabled(c.Metrics.Enabled)

	// 统一响应格式中间件（需要在最外层，确保所有响应都被统一格式化）
	server.Use(middleware.ResponseMiddleware())

	// API 指标收集中间件（记录请求计数和耗时）- 可通过配置开关控制
	if c.Metrics.Enabled {
		server.Use(middleware.ApiMetricMiddleware())
	}

	// CORS中间件（处理跨域请求）
	server.Use(middleware.CorsMiddleware(c.CORS.AllowedOrigins))

	// 登录速率限制中间件（防止暴力破解）
	server.Use(middleware.LoginRateLimitMiddleware())

	// 认证中间件（跳过登录相关接口）
	server.Use(func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// 跳过认证的路径
			skipPaths := []string{"/auth/login", "/auth/refresh"}

			for _, path := range skipPaths {
				if strings.HasPrefix(r.URL.Path, path) {
					next.ServeHTTP(w, r)
					return
				}
			}

			// /system/init-data 需要特殊保护：仅允许本地访问
			if strings.HasPrefix(r.URL.Path, "/system/init-data") {
				clientIP := getClientIP(r)
				if !isLocalhost(clientIP) {
					w.Header().Set("Content-Type", "application/json; charset=utf-8")
					w.WriteHeader(http.StatusForbidden)
					w.Write([]byte(`{"code":403,"data":null,"message":"init-data 接口仅允许本地访问"}`))
					return
				}
				next.ServeHTTP(w, r)
				return
			}

			// 其他接口需要认证
			middleware.AuthMiddleware(ctx)(next).ServeHTTP(w, r)
		}
	})

	// 权限中间件（基于路径的权限检查）
	server.Use(func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// 跳过权限检查的路径
			skipPaths := []string{"/auth/login", "/auth/refresh", "/auth/logout", "/system/init-data"}

			for _, path := range skipPaths {
				if strings.HasPrefix(r.URL.Path, path) {
					next.ServeHTTP(w, r)
					return
				}
			}

			// 其他接口需要权限检查
			middleware.PermissionMiddlewareByPath(ctx)(next).ServeHTTP(w, r)
		}
	})

	// 调用routes.go中的路由注册函数（由goctl自动生成）
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}

// getClientIP 获取客户端真实 IP
func getClientIP(r *http.Request) string {
	// 优先从代理头获取
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		return strings.TrimSpace(parts[0])
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

// isLocalhost 检查是否为本地地址
func isLocalhost(ip string) bool {
	if ip == "127.0.0.1" || ip == "::1" || ip == "localhost" {
		return true
	}
	// 检查是否为私有地址
	if parsedIP := net.ParseIP(ip); parsedIP != nil {
		return parsedIP.IsLoopback() || parsedIP.IsPrivate()
	}
	return false
}

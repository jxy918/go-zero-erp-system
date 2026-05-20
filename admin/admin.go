// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package main

import (
	"flag"
	"fmt"
	"net/http"
	"strings"

	"myproject/admin/internal/config"
	"myproject/admin/internal/handler"
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

	// 统一响应格式中间件（需要在最外层，确保所有响应都被统一格式化）
	server.Use(middleware.ResponseMiddleware())

	// CORS中间件（处理跨域请求）
	server.Use(middleware.CorsMiddleware())

	// 认证中间件（跳过登录相关接口）
	server.Use(func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// 跳过认证的路径
			skipPaths := []string{"/auth/login", "/auth/refresh", "/system/init-data"}

			for _, path := range skipPaths {
				if strings.HasPrefix(r.URL.Path, path) {
					next.ServeHTTP(w, r)
					return
				}
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

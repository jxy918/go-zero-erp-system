package middleware

import (
	"net/http"
	"strings"
)

// CorsMiddleware 处理跨域请求的中间件
// allowedOrigins: 允许的域名列表，逗号分隔。空字符串表示允许所有域名（开发环境）
func CorsMiddleware(allowedOrigins string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			if allowedOrigins == "" || allowedOrigins == "*" {
				// 开发环境：允许所有域名
				w.Header().Set("Access-Control-Allow-Origin", "*")
			} else {
				// 生产环境：检查请求域名是否在允许列表中
				allowedList := strings.Split(allowedOrigins, ",")
				for _, allowed := range allowedList {
					if origin == strings.TrimSpace(allowed) {
						w.Header().Set("Access-Control-Allow-Origin", origin)
						w.Header().Set("Vary", "Origin")
						break
					}
				}
				// 如果请求域名不在允许列表中，不设置 Access-Control-Allow-Origin
				// 浏览器会阻止跨域请求
			}

			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, user_id")
			w.Header().Set("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
			w.Header().Set("Access-Control-Allow-Credentials", "true")

			// 处理OPTIONS请求
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			// 继续处理请求
			next(w, r)
		}
	}
}

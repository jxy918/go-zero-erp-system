package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"myproject/admin/internal/model"
	"myproject/admin/internal/svc"
	"myproject/admin/internal/util"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// contextKey 定义上下文键类型，用于类型安全的上下文值存储
type contextKey string

// 请求上下文中存储用户信息的键名常量
const (
	usernameKey contextKey = "username" // 用户名
)

// AuthMiddleware 认证中间件 - 核心认证逻辑
// 功能：验证JWT Token、检查用户状态、存储用户信息到上下文、记录活动日志
// 所有需要认证的API请求都会经过此中间件
func AuthMiddleware(svc *svc.ServiceContext) func(next http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// 步骤1: 从请求头获取Authorization Token
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				httpx.Error(w, errors.New("unauthorized"))
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			// 步骤2: 解析Bearer Token格式
			// 格式要求: "Bearer <token>"
			parts := strings.SplitN(authHeader, " ", 2)
			if !(len(parts) == 2 && parts[0] == "Bearer") {
				httpx.Error(w, errors.New("unauthorized"))
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			// 步骤3: 验证并解析JWT Token
			token := parts[1]
			claims, err := util.ParseToken(token, svc.Config.JWT.AccessSecret)
			if err != nil {
				httpx.Error(w, err)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

		// 步骤4: 检查用户是否存在且状态正常
		user, err := svc.UserModel.GetByID(claims.UserID)
		if err != nil {
			httpx.Error(w, errors.New("认证服务异常，请稍后重试"))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if user == nil {
			httpx.Error(w, errors.New("用户不存在"))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if user.Status == 0 {
			httpx.Error(w, errors.New("账号已被禁用，请联系管理员"))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// 步骤5: 获取用户的角色ID和管理员状态
		roleID := uint(0)
		isAdmin := false
		if len(user.Roles) > 0 {
			roleID = user.Roles[0].ID
			// 基于角色 code 判断是否为管理员，而非用户名
			isAdmin = user.Roles[0].Code == "admin"
		}

		// 步骤6: 将用户信息存储到请求上下文，供后续逻辑使用
		ctx := r.Context()
		ctx = context.WithValue(ctx, util.UserIDKey, claims.UserID)    // 用户ID
		ctx = context.WithValue(ctx, usernameKey, claims.Username)      // 用户名
		ctx = context.WithValue(ctx, util.IsAdminKey, isAdmin)          // 是否管理员（基于角色）
		ctx = context.WithValue(ctx, util.RoleIDKey, roleID)            // 角色ID
		r = r.WithContext(ctx)

		// 步骤7: 异步记录用户活动日志（不阻塞请求）
		activity := &model.Activity{
			UserID:   claims.UserID,
			Username: claims.Username,
			Action:   "访问API",
			URL:      r.URL.Path,
			IP:       util.GetClientIP(r), // 获取真实客户端IP（支持代理）
		}
		go func() {
			defer func() {
				if r := recover(); r != nil {
					// 记录 panic 但不影响主流程
					logx.Errorf("Activity log panic: %v", r)
				}
			}()
			svc.ActivityModel.Create(activity)
		}()

			// 步骤8: 继续处理下一个中间件或处理函数
			next.ServeHTTP(w, r)
		}
	}
}

// PermissionMiddlewareByPath 基于路径的权限中间件（从数据库读取权限配置）
func PermissionMiddlewareByPath(svc *svc.ServiceContext) func(next http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			// 从 context 获取用户ID
			userID, ok := ctx.Value(util.UserIDKey).(uint)
		if !ok || userID == 0 {
			httpx.Error(w, errors.New("认证信息无效"))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

			// 检查是否为管理员（管理员拥有所有权限）
			isAdmin, ok := ctx.Value(util.IsAdminKey).(bool)
			if ok && isAdmin {
				next.ServeHTTP(w, r)
				return
			}

			// 获取请求路径
			path := r.URL.Path

			// 从数据库查询该路径对应的权限
			requiredPermission, err := svc.PermissionModel.GetByPath(path)
			if err != nil {
				// 没找到对应的权限配置，直接通过
				next.ServeHTTP(w, r)
				return
			}

			// 根据用户ID获取用户的权限列表
			userPermissions, err := svc.PermissionModel.GetAllByUserID(userID)
		if err != nil {
			httpx.Error(w, errors.New("权限服务异常，请稍后重试"))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

			// 检查是否包含所需的权限
			hasPermission := false
			for _, perm := range userPermissions {
				if perm.Code == requiredPermission.Code {
					hasPermission = true
					break
				}
			}

			// 如果没有权限，返回403错误并给出友好提示
			if !hasPermission {
				jsonError(w, http.StatusForbidden, "您暂无权限访问此功能，请联系管理员")
				return
			}

			next.ServeHTTP(w, r)
		}
	}
}

// GetUsername 获取当前用户的用户名
func GetUsername(ctx context.Context) string {
	if username, ok := ctx.Value(usernameKey).(string); ok {
		return username
	}
	return ""
}

// jsonError 返回标准JSON格式的错误响应
func jsonError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	response := map[string]interface{}{
		"code":    code,
		"data":    nil,
		"message": message,
	}
	json.NewEncoder(w).Encode(response)
}

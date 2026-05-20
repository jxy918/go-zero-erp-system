package util

import (
	"context"
)

// contextKey 定义上下文键类型，用于类型安全的上下文值存储
type contextKey string

// 请求上下文中存储用户信息的键名常量
const (
	UserIDKey  contextKey = "user_id"  // 用户ID - 用于数据权限过滤
	RoleIDKey  contextKey = "role_id"  // 角色ID - 用于数据权限过滤
	IsAdminKey contextKey = "is_admin" // 是否管理员 - 用于判断是否跳过权限检查
)

// GetUserID 从请求上下文中获取当前用户ID
// 该函数在数据权限过滤逻辑中使用，用于限制普通用户只能访问自己的数据
// 返回值：用户ID，如果未登录或获取失败则返回0
func GetUserID(ctx context.Context) uint {
	if id, ok := ctx.Value(UserIDKey).(uint); ok {
		return id
	}
	return 0
}

// GetRoleID 从请求上下文中获取当前角色ID
// 该函数在数据权限过滤逻辑中使用，用于限制用户只能查看自己角色范围内的数据
// 返回值：角色ID，如果未登录或获取失败则返回0
func GetRoleID(ctx context.Context) uint {
	if id, ok := ctx.Value(RoleIDKey).(uint); ok {
		return id
	}
	return 0
}

// IsAdmin 判断当前用户是否为管理员
// 管理员拥有所有权限，可以访问所有数据，不受数据权限过滤限制
// 返回值：true-管理员，false-普通用户或未登录
func IsAdmin(ctx context.Context) bool {
	if isAdmin, ok := ctx.Value(IsAdminKey).(bool); ok {
		return isAdmin
	}
	return false
}

// PermissionFunc 定义数据权限检查函数的类型
type PermissionFunc func(ctx context.Context) bool

// CheckDataPermission 数据权限检查工具
// 这个函数用于在业务逻辑中进行数据级别的权限检查
// 例如：检查用户是否有权限修改某条数据
func CheckDataPermission(ctx context.Context, checks ...PermissionFunc) bool {
	// 管理员拥有所有数据权限
	if IsAdmin(ctx) {
		return true
	}

	// 执行所有检查函数，必须全部通过
	for _, check := range checks {
		if !check(ctx) {
			return false
		}
	}

	return true
}

// IsOwnerCheck 检查是否为数据的创建者/所有者
func IsOwnerCheck(ctx context.Context, ownerID uint) PermissionFunc {
	return func(ctx context.Context) bool {
		return GetUserID(ctx) == ownerID
	}
}

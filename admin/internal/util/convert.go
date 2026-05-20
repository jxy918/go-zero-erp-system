package util

import (
	"net"
	"net/http"
	"strings"

	"myproject/admin/internal/model"
	"myproject/admin/internal/types"
)

// GetClientIP 获取客户端真实IP地址
func GetClientIP(r *http.Request) string {
	// 按优先级检查常用的代理头
	headers := []string{
		"X-Forwarded-For",
		"X-Real-IP",
		"Proxy-Client-IP",
		"WL-Proxy-Client-IP",
		"HTTP_X_FORWARDED_FOR",
		"HTTP_X_FORWARDED",
		"HTTP_X_CLUSTER_CLIENT_IP",
		"HTTP_CLIENT_IP",
		"HTTP_FORWARDED_FOR",
		"HTTP_FORWARDED",
	}

	for _, header := range headers {
		ip := r.Header.Get(header)
		if ip != "" {
			// X-Forwarded-For 可能包含多个IP，取第一个
			ips := strings.Split(ip, ",")
			clientIP := strings.TrimSpace(ips[0])
			// 去掉可能的端口号
			clientIP = extractIP(clientIP)
			if isValidIP(clientIP) {
				return formatIP(clientIP)
			}
		}
	}

	// 如果没有代理头，使用 RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		// 如果无法解析，尝试直接提取IP
		ip = extractIP(r.RemoteAddr)
	}
	return formatIP(ip)
}

// extractIP 从字符串中提取IP地址（去掉端口）
func extractIP(ipPort string) string {
	// 处理 IPv6 格式 [::1]:port
	if strings.HasPrefix(ipPort, "[") {
		idx := strings.Index(ipPort, "]")
		if idx != -1 {
			return ipPort[1:idx]
		}
	}
	// 处理 IPv4 格式 ip:port
	ip, _, err := net.SplitHostPort(ipPort)
	if err == nil {
		return ip
	}
	return ipPort
}

// formatIP 格式化IP地址，将IPv6本地地址转换为IPv4格式
func formatIP(ip string) string {
	// 如果是IPv6本地回环地址，转换为IPv4格式
	if ip == "::1" {
		return "127.0.0.1"
	}
	// 如果是完整的IPv6地址（包含::），保持原样
	return ip
}

// isValidIP 验证IP地址是否有效
func isValidIP(ip string) bool {
	return net.ParseIP(ip) != nil
}

// ConvertRoles 转换角色列表
func ConvertRoles(roles []model.Role) []types.RoleInfo {
	var roleInfos []types.RoleInfo
	for _, role := range roles {
		roleInfos = append(roleInfos, types.RoleInfo{
			ID:          role.ID,
			Name:        role.Name,
			Code:        role.Code,
			Desc:        role.Desc,
			Status:      role.Status,
			Permissions: ConvertPermissions(role.Permissions),
		})
	}
	return roleInfos
}

// ConvertPermissions 转换权限列表
// 角色关联的权限都是按钮权限（type=2）
func ConvertPermissions(permissions []model.Permission) []types.PermissionInfo {
	var permissionInfos []types.PermissionInfo
	for _, permission := range permissions {
		permissionInfos = append(permissionInfos, types.PermissionInfo{
			ID:       permission.ID,
			Name:     permission.Name,
			Code:     permission.Code,
			Desc:     permission.Desc,
			Type:     2, // 按钮权限
			ParentID: 0,
			Sort:     permission.Sort,
			Status:   permission.Status,
		})
	}
	return permissionInfos
}

// ConvertMenus 转换菜单列表
func ConvertMenus(menus []model.Menu) []types.MenuInfo {
	var menuInfos []types.MenuInfo
	for _, menu := range menus {
		menuInfos = append(menuInfos, types.MenuInfo{
			ID:          menu.ID,
			Name:        menu.Name,
			Code:        menu.Code,
			Desc:        menu.Desc,
			ParentID:    menu.ParentID,
			Path:        menu.Path,
			Component:   menu.Component,
			Icon:        menu.Icon,
			Sort:        menu.Sort,
			Status:      menu.Status,
			Permissions: ConvertPermissions(menu.Permissions),
		})
	}
	return menuInfos
}

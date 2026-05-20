// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"myproject/admin/internal/model"
	"myproject/admin/internal/svc"
	"myproject/admin/internal/types"
	"myproject/admin/internal/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListRoleLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
	ctx    context.Context
}

func NewListRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListRoleLogic {
	return &ListRoleLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
		ctx:    ctx,
	}
}

func (l *ListRoleLogic) ListRole(req *types.ListRoleRequest) (resp *types.ListRoleResponse, err error) {
	// 1. 设置默认分页参数
	page := req.Page
	if page <= 0 {
		page = 1
	}

	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	// 2. 获取当前用户权限信息
	isAdmin := util.IsAdmin(l.ctx)
	userID := util.GetUserID(l.ctx)

	// 3. 根据权限获取角色列表
	var roles []model.Role
	var total int64
	if isAdmin {
		// 管理员可以查看所有角色
		roles, total, err = l.svcCtx.RoleModel.ListWithName(page, pageSize, req.Name)
	} else if userID > 0 {
		// 普通用户只能查看自己所属的角色
		roles, total, err = l.svcCtx.RoleModel.ListByUserID(userID, page, pageSize, req.Name)
	} else {
		// 如果用户ID为0，返回空列表
		roles = []model.Role{}
		total = 0
	}
	if err != nil {
		return nil, err
	}

	// 3. 获取所有权限与菜单的关联关系
	var menuPerms []model.MenuPermission
	if err := l.svcCtx.DB.Find(&menuPerms).Error; err != nil {
		return nil, err
	}
	
	// 构建权限ID到菜单ID的映射
	permToMenuMap := make(map[uint]uint)
	for _, mp := range menuPerms {
		permToMenuMap[mp.PermissionID] = mp.MenuID
	}
	
	// 4. 构建响应
	var roleInfos []types.RoleInfo
	for _, role := range roles {
		// 为角色的每个权限设置正确的父菜单ID
		var permissionInfos []types.PermissionInfo
		for _, perm := range role.Permissions {
			parentID := uint(0)
			if menuID, exists := permToMenuMap[perm.ID]; exists {
				parentID = menuID
			}
			permissionInfos = append(permissionInfos, types.PermissionInfo{
				ID:       perm.ID,
				Name:     perm.Name,
				Code:     perm.Code,
				Desc:     perm.Desc,
				Type:     2, // 按钮权限
				ParentID: parentID,
				Sort:     perm.Sort,
				Status:   perm.Status,
			})
		}
		
		roleInfos = append(roleInfos, types.RoleInfo{
			ID:          role.ID,
			Name:        role.Name,
			Code:        role.Code,
			Desc:        role.Desc,
			Status:      role.Status,
			Permissions: permissionInfos,
			Menus:       util.ConvertMenus(role.Menus),
		})
	}

	resp = &types.ListRoleResponse{
		Roles: roleInfos,
		Total: total,
	}

	return resp, nil
}

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

type ListPermissionLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
	ctx    context.Context
}

func NewListPermissionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListPermissionLogic {
	return &ListPermissionLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
		ctx:    ctx,
	}
}

func (l *ListPermissionLogic) ListPermission(req *types.ListPermissionRequest) (resp *types.ListPermissionResponse, err error) {
	// 1. 设置默认分页参数
	page := req.Page
	if page <= 0 {
		page = 1
	}

	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 500
	}

	// 2. 获取当前用户权限信息
	isAdmin := util.IsAdmin(l.ctx)
	userID := util.GetUserID(l.ctx)

	// 3. 构建权限列表（同时包含菜单和按钮权限）
	var permissionInfos []types.PermissionInfo

	if isAdmin || userID > 0 {
		// 获取所有菜单
		var menus []model.Menu
		menuQuery := l.svcCtx.DB.Model(&model.Menu{})

		if req.Name != "" {
			menuQuery = menuQuery.Where("name LIKE ? OR code LIKE ?", "%"+req.Name+"%", "%"+req.Name+"%")
		}

		if err := menuQuery.Order("sort asc").Find(&menus).Error; err != nil {
			return nil, err
		}

		// 获取所有按钮权限
		var buttonPermissions []model.Permission
		permQuery := l.svcCtx.DB.Model(&model.Permission{})

		if req.Name != "" {
			permQuery = permQuery.Where("name LIKE ? OR code LIKE ?", "%"+req.Name+"%", "%"+req.Name+"%")
		}

		if err := permQuery.Order("sort asc").Find(&buttonPermissions).Error; err != nil {
			return nil, err
		}

		// 构建按钮权限ID到权限对象的映射
		permMap := make(map[uint]model.Permission)
		for _, perm := range buttonPermissions {
			permMap[perm.ID] = perm
		}

		// 通过 menu_permissions 表获取所有按钮权限及其关联的菜单
		var menuPerms []model.MenuPermission
		if err := l.svcCtx.DB.Find(&menuPerms).Error; err != nil {
			return nil, err
		}

		// 构建菜单ID到按钮权限ID的映射
		menuToPermsMap := make(map[uint][]uint)
		for _, mp := range menuPerms {
			menuToPermsMap[mp.MenuID] = append(menuToPermsMap[mp.MenuID], mp.PermissionID)
		}

		// 构建按钮权限ID到菜单ID的映射
		permToMenuMap := make(map[uint]uint)
		for _, mp := range menuPerms {
			permToMenuMap[mp.PermissionID] = mp.MenuID
		}

		// 构建树形结构：先添加菜单，再添加按钮权限
		menuMap := make(map[uint]bool)
		for _, menu := range menus {
			menuMap[menu.ID] = true
			// 添加菜单
			permissionInfos = append(permissionInfos, types.PermissionInfo{
				ID:       menu.ID,
				Name:     menu.Name,
				Code:     menu.Code,
				Desc:     menu.Desc,
				Path:     menu.Path,
				Type:     1, // 1: 菜单
				ParentID: menu.ParentID,
				MenuID:   0,
				Sort:     menu.Sort,
				Status:   menu.Status,
			})

			// 添加该菜单下的按钮权限
			if permIDs, exists := menuToPermsMap[menu.ID]; exists {
				for _, permID := range permIDs {
					if perm, ok := permMap[permID]; ok {
						permissionInfos = append(permissionInfos, types.PermissionInfo{
							ID:       perm.ID,
							Name:     perm.Name,
							Code:     perm.Code,
							Desc:     perm.Desc,
							Path:     perm.Path,
							Type:     2,       // 2: 按钮权限
							ParentID: menu.ID, // 通过 menu_permissions 表关联
							MenuID:   menu.ID, // 记录关联的菜单ID
							Sort:     perm.Sort,
							Status:   perm.Status,
						})
					}
				}
			}
		}

		resp = &types.ListPermissionResponse{
			Permissions: permissionInfos,
			Total:       int64(len(permissionInfos)),
		}
	} else {
		resp = &types.ListPermissionResponse{
			Permissions: []types.PermissionInfo{},
			Total:       0,
		}
	}

	return resp, nil
}

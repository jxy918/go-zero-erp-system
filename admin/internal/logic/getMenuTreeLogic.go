package logic

import (
	"context"

	"myproject/admin/internal/model"
	"myproject/admin/internal/svc"
	"myproject/admin/internal/types"
	"myproject/admin/internal/util"

	"github.com/zeromicro/go-zero/core/logx"
)

// GetMenuTreeLogic 获取菜单树逻辑
func NewGetMenuTreeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMenuTreeLogic {
	return &GetMenuTreeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

type GetMenuTreeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// GetMenuTree 获取菜单树
func (l *GetMenuTreeLogic) GetMenuTree(req types.ListMenuRequest) (*types.ListMenuResponse, error) {
	// 从请求上下文中获取用户ID
	userID := util.GetUserID(l.ctx)
	isAdmin := util.IsAdmin(l.ctx)

	logx.Infof("GetMenuTree - userID: %d, isAdmin: %v", userID, isAdmin)

	if userID == 0 {
		logx.Error("user_id not found in context")
		return &types.ListMenuResponse{
			Menus: []types.MenuInfo{},
			Total: 0,
		}, nil
	}

	// 获取用户所有菜单
	menus, err := l.svcCtx.MenuModel.GetUserMenus(userID, isAdmin)
	if err != nil {
		logx.Errorf("GetUserMenus error: %v", err)
		return nil, err
	}

	logx.Infof("GetMenuTree - 从数据库获取到 %d 个菜单", len(menus))
	for _, menu := range menus {
		logx.Infof("菜单: id=%d, name=%s, parent_id=%d", menu.ID, menu.Name, menu.ParentID)
	}

	// 构建菜单树
	menuInfos := l.buildMenuTree(menus)

	return &types.ListMenuResponse{
		Menus: menuInfos,
		Total: int64(len(menuInfos)),
	}, nil
}

// buildMenuTree 构建菜单树
func (l *GetMenuTreeLogic) buildMenuTree(menus []model.Menu) []types.MenuInfo {
	// 创建所有菜单的Info并存储在map中（自动去重，因为map的key是唯一的）
	menuMap := make(map[uint]types.MenuInfo)
	for _, menu := range menus {
		menuMap[menu.ID] = types.MenuInfo{
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
			Permissions: []types.PermissionInfo{},
			Children:    []types.MenuInfo{},
		}
	}

	// 首先添加所有根菜单（ParentID=0），使用map去重
	var menuTree []types.MenuInfo
	addedRootIDs := make(map[uint]bool) // 用于记录已添加的根菜单ID
	for _, menu := range menus {
		if menu.ParentID == 0 && !addedRootIDs[menu.ID] {
			if info, ok := menuMap[menu.ID]; ok {
				menuTree = append(menuTree, info)
				addedRootIDs[menu.ID] = true // 标记已添加
			}
		}
	}

	// 然后为每个根菜单添加子菜单
	for i := range menuTree {
		rootID := menuTree[i].ID
		for _, menu := range menus {
			if menu.ParentID == rootID {
				if child, ok := menuMap[menu.ID]; ok {
					menuTree[i].Children = append(menuTree[i].Children, child)
				}
			}
		}
	}

	return menuTree
}

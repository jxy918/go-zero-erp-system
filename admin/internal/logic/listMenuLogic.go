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

type ListMenuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListMenuLogic {
	return &ListMenuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListMenuLogic) ListMenu(req *types.ListMenuRequest) (resp *types.ListMenuResponse, err error) {
	// 获取当前用户权限信息
	isAdmin := util.IsAdmin(l.ctx)
	userID := util.GetUserID(l.ctx)

	var menus []interface{}
	var total int64

	if isAdmin {
		// 管理员可以查看所有菜单
		adminMenus, cnt, err := l.svcCtx.MenuModel.List(req.Page, req.PageSize)
		if err != nil {
			return nil, err
		}
		total = cnt
		for _, menu := range adminMenus {
			menus = append(menus, menu)
		}
	} else if userID > 0 {
		// 普通用户只能查看自己有权限的菜单
		menus, total, err = l.svcCtx.MenuModel.ListByUserID(userID, req.Page, req.PageSize)
	} else {
		// 如果用户ID为0，返回空列表
		return &types.ListMenuResponse{
			Menus: []types.MenuInfo{},
			Total: 0,
		}, nil
	}
	if err != nil {
		return nil, err
	}

	menuInfos := make([]types.MenuInfo, 0, len(menus))
	for _, item := range menus {
		switch v := item.(type) {
		case model.Menu:
			menuInfos = append(menuInfos, types.MenuInfo{
				ID:          v.ID,
				Name:        v.Name,
				Code:        v.Code,
				Desc:        v.Desc,
				ParentID:    v.ParentID,
				Path:        v.Path,
				Component:   v.Component,
				Icon:        v.Icon,
				Sort:        v.Sort,
				Status:      v.Status,
				Permissions: []types.PermissionInfo{},
				Children:    []types.MenuInfo{},
			})
		case *model.Menu:
			menuInfos = append(menuInfos, types.MenuInfo{
				ID:          v.ID,
				Name:        v.Name,
				Code:        v.Code,
				Desc:        v.Desc,
				ParentID:    v.ParentID,
				Path:        v.Path,
				Component:   v.Component,
				Icon:        v.Icon,
				Sort:        v.Sort,
				Status:      v.Status,
				Permissions: []types.PermissionInfo{},
				Children:    []types.MenuInfo{},
			})
		}
	}

	return &types.ListMenuResponse{
		Menus: menuInfos,
		Total: total,
	}, nil
}

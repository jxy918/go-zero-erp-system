package logic

import (
	"context"

	"myproject/admin/internal/svc"
	"myproject/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateMenuLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateMenuLogic {
	return &UpdateMenuLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateMenuLogic) UpdateMenu(req types.UpdateMenuRequest) (resp *types.MenuInfo, err error) {
	menu, err := l.svcCtx.MenuModel.Get(req.ID)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		menu.Name = req.Name
	}
	if req.Code != "" {
		menu.Code = req.Code
	}
	if req.Desc != "" {
		menu.Desc = req.Desc
	}
	menu.ParentID = req.ParentID
	if req.Path != "" {
		menu.Path = req.Path
	}
	if req.Component != "" {
		menu.Component = req.Component
	}
	if req.Icon != "" {
		menu.Icon = req.Icon
	}
	if req.Sort != 0 {
		menu.Sort = req.Sort
	}
	// 始终更新状态字段（0 和 1 都是有效的状态值）
	menu.Status = req.Status

	if err := l.svcCtx.MenuModel.Update(menu); err != nil {
		return nil, err
	}

	return &types.MenuInfo{
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
	}, nil
}

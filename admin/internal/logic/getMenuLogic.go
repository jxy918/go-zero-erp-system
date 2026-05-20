package logic

import (
	"context"

	"myproject/admin/internal/svc"
	"myproject/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMenuLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMenuLogic {
	return &GetMenuLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetMenuLogic) GetMenu(req types.GetMenuRequest) (resp *types.MenuInfo, err error) {
	menu, err := l.svcCtx.MenuModel.Get(req.ID)
	if err != nil {
		return nil, err
	}

	// 将权限列表转换为响应格式
	var permissions []types.PermissionInfo
	for _, p := range menu.Permissions {
		permissions = append(permissions, types.PermissionInfo{
			ID:     p.ID,
			Name:   p.Name,
			Code:   p.Code,
			Desc:   p.Desc,
			Sort:   p.Sort,
			Status: p.Status,
		})
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
		Permissions: permissions,
		Children:    []types.MenuInfo{},
	}, nil
}

package logic

import (
	"context"

	"myproject/admin/internal/svc"
	"myproject/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AssignMenuPermissionsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAssignMenuPermissionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AssignMenuPermissionsLogic {
	return &AssignMenuPermissionsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AssignMenuPermissionsLogic) AssignMenuPermissions(req types.AssignMenuPermissionsRequest) (resp *types.MenuInfo, err error) {
	if err := l.svcCtx.MenuModel.AssignPermissions(req.MenuID, req.PermissionIDs); err != nil {
		return nil, err
	}

	menu, err := l.svcCtx.MenuModel.Get(req.MenuID)
	if err != nil {
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

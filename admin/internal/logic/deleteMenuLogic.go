package logic

import (
	"context"

	"myproject/admin/internal/svc"
	"myproject/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteMenuLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteMenuLogic {
	return &DeleteMenuLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteMenuLogic) DeleteMenu(req *types.DeleteMenuRequest) (resp *types.MenuInfo, err error) {
	menu, err := l.svcCtx.MenuModel.Get(req.ID)
	if err != nil {
		return nil, err
	}

	if err := l.svcCtx.MenuModel.Delete(req.ID); err != nil {
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

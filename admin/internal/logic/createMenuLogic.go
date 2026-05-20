package logic

import (
	"context"

	"myproject/admin/internal/model"
	"myproject/admin/internal/svc"
	"myproject/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateMenuLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateMenuLogic {
	return &CreateMenuLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateMenuLogic) CreateMenu(req types.CreateMenuRequest) (resp *types.MenuInfo, err error) {
	menu := &model.Menu{
		Name:      req.Name,
		Code:      req.Code,
		Desc:      req.Desc,
		ParentID:  req.ParentID,
		Path:      req.Path,
		Component: req.Component,
		Icon:      req.Icon,
		Sort:      req.Sort,
		Status:    req.Status,
	}

	if err := l.svcCtx.MenuModel.Create(menu); err != nil {
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

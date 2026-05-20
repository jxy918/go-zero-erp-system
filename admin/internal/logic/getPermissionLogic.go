// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"errors"

	"myproject/admin/internal/model"
	"myproject/admin/internal/svc"
	"myproject/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPermissionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPermissionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPermissionLogic {
	return &GetPermissionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPermissionLogic) GetPermission(req *types.GetPermissionRequest) (resp *types.PermissionInfo, err error) {
	// 1. 根据ID获取权限
	permission, err := l.svcCtx.PermissionModel.GetByID(req.ID)
	if err != nil {
		return nil, err
	}
	if permission == nil {
		return nil, errors.New("权限不存在")
	}

	// 2. 获取关联的菜单ID
	var menuID uint
	var menuPerm model.MenuPermission
	if err := l.svcCtx.DB.Model(&model.MenuPermission{}).
		Where("permission_id = ?", permission.ID).
		First(&menuPerm).Error; err == nil {
		menuID = menuPerm.MenuID
	}

	// 3. 构建响应
	resp = &types.PermissionInfo{
		ID:       permission.ID,
		Name:     permission.Name,
		Code:     permission.Code,
		Desc:     permission.Desc,
		Path:     permission.Path,
		Type:     2, // 按钮权限
		ParentID: menuID,
		Sort:     permission.Sort,
		Status:   permission.Status,
	}

	return resp, nil
}

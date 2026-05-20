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

type DeletePermissionLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewDeletePermissionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeletePermissionLogic {
	return &DeletePermissionLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *DeletePermissionLogic) DeletePermission(req *types.DeletePermissionRequest) (resp *types.PermissionInfo, err error) {
	// 1. 根据ID获取权限
	permission, err := l.svcCtx.PermissionModel.GetByID(req.ID)
	if err != nil {
		return nil, err
	}
	if permission == nil {
		return nil, errors.New("权限不存在")
	}

	// 2. 保存权限信息用于响应
	permissionInfo := &types.PermissionInfo{
		ID:     permission.ID,
		Name:   permission.Name,
		Code:   permission.Code,
		Desc:   permission.Desc,
		Sort:   permission.Sort,
		Status: permission.Status,
	}

	// 3. 硬删除 menu_permissions 关联
	if err := l.svcCtx.DB.Unscoped().Where("permission_id = ?", permission.ID).Delete(&model.MenuPermission{}).Error; err != nil {
		logx.Errorf("删除菜单权限关联失败: %v", err)
	}

	// 4. 删除权限
	err = l.svcCtx.PermissionModel.Delete(req.ID)
	if err != nil {
		return nil, err
	}

	return permissionInfo, nil
}

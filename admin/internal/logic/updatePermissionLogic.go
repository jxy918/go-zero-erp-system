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

type UpdatePermissionLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewUpdatePermissionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatePermissionLogic {
	return &UpdatePermissionLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *UpdatePermissionLogic) UpdatePermission(req *types.UpdatePermissionRequest) (resp *types.PermissionInfo, err error) {
	// 1. 根据ID获取权限
	permission, err := l.svcCtx.PermissionModel.GetByID(req.ID)
	if err != nil {
		return nil, err
	}
	if permission == nil {
		return nil, errors.New("权限不存在")
	}

	// 2. 更新权限信息
	if req.Name != "" {
		permission.Name = req.Name
	}

	if req.Code != "" {
		// 检查权限代码是否已存在
		existingPermission, err := l.svcCtx.PermissionModel.GetByCode(req.Code)
		if err != nil {
			return nil, err
		}
		if existingPermission != nil && existingPermission.ID != req.ID {
			return nil, errors.New("权限代码已存在")
		}
		permission.Code = req.Code
	}

	if req.Desc != "" {
		permission.Desc = req.Desc
	}

	if req.Path != "" {
		permission.Path = req.Path
	}

	if req.Sort != 0 {
		permission.Sort = req.Sort
	}

	// 状态字段：0 表示禁用，1 表示启用，需要始终更新
	permission.Status = req.Status

	// 3. 保存更新
	err = l.svcCtx.PermissionModel.Update(permission)
	if err != nil {
		return nil, err
	}

	// 4. 更新 menu_permissions 关联
	var currentMenuID uint
	err = l.svcCtx.DB.Model(&model.MenuPermission{}).
		Where("permission_id = ?", permission.ID).
		Select("menu_id").Scan(&currentMenuID).Error
	if err != nil {
		logx.Errorf("查询当前菜单关联失败: %v", err)
	}

	if currentMenuID != req.MenuID {
		if req.MenuID > 0 {
			if currentMenuID > 0 {
				if err := l.svcCtx.DB.Model(&model.MenuPermission{}).
					Where("permission_id = ?", permission.ID).
					Update("menu_id", req.MenuID).Error; err != nil {
					logx.Errorf("更新菜单权限关联失败: %v", err)
				}
			} else {
				menuPerm := &model.MenuPermission{
					MenuID:       req.MenuID,
					PermissionID: permission.ID,
				}
				if err := l.svcCtx.DB.Create(menuPerm).Error; err != nil {
					logx.Errorf("创建菜单权限关联失败: %v", err)
				}
			}
		} else {
			if currentMenuID > 0 {
				if err := l.svcCtx.DB.Where("permission_id = ?", permission.ID).Delete(&model.MenuPermission{}).Error; err != nil {
					logx.Errorf("删除菜单权限关联失败: %v", err)
				}
			}
		}
	}

	// 5. 构建响应
	resp = &types.PermissionInfo{
		ID:     permission.ID,
		Name:   permission.Name,
		Code:   permission.Code,
		Desc:   permission.Desc,
		Sort:   permission.Sort,
		Status: permission.Status,
	}

	return resp, nil
}

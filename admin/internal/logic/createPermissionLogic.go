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

type CreatePermissionLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewCreatePermissionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreatePermissionLogic {
	return &CreatePermissionLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *CreatePermissionLogic) CreatePermission(req *types.CreatePermissionRequest) (resp *types.PermissionInfo, err error) {
	// 调试日志
	logx.Infof("创建权限请求: Name=%s, Code=%s, Desc=%s, Path=%s, Sort=%d, Status=%d, MenuID=%d",
		req.Name, req.Code, req.Desc, req.Path, req.Sort, req.Status, req.MenuID)

	// 1. 检查权限代码是否已存在
	existingPermission, err := l.svcCtx.PermissionModel.GetByCode(req.Code)
	if err != nil {
		return nil, err
	}
	if existingPermission != nil {
		return nil, errors.New("权限代码已存在")
	}

	// 2. 创建权限
	permission := &model.Permission{
		Name:   req.Name,
		Code:   req.Code,
		Desc:   req.Desc,
		Path:   req.Path,
		Sort:   req.Sort,
		Status: req.Status,
	}

	err = l.svcCtx.PermissionModel.Create(permission)
	if err != nil {
		return nil, err
	}

	// 3. 如果关联了菜单，创建 menu_permissions 关联
	logx.Infof("创建权限: MenuID=%d, PermissionID=%d", req.MenuID, permission.ID)
	if req.MenuID > 0 {
		menuPerm := &model.MenuPermission{
			MenuID:       req.MenuID,
			PermissionID: permission.ID,
		}
		if err := l.svcCtx.DB.Create(menuPerm).Error; err != nil {
			logx.Errorf("创建菜单权限关联失败: %v", err)
		} else {
			logx.Infof("菜单权限关联创建成功: MenuID=%d, PermissionID=%d", req.MenuID, permission.ID)
		}
	} else {
		logx.Infof("MenuID=%d，跳过菜单权限关联创建", req.MenuID)
	}

	// 4. 构建响应
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

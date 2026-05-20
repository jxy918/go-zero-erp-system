// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"errors"

	"myproject/admin/internal/svc"
	"myproject/admin/internal/types"
	"myproject/admin/internal/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type AssignPermissionsLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAssignPermissionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AssignPermissionsLogic {
	return &AssignPermissionsLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AssignPermissionsLogic) AssignPermissions(req *types.AssignPermissionsRequest) (resp *types.RoleInfo, err error) {
	// 1. 检查角色是否存在
	role, err := l.svcCtx.RoleModel.GetByID(req.RoleID)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, errors.New("角色不存在")
	}

	// 2. 为角色分配权限
	err = l.svcCtx.RoleModel.AssignPermissions(req.RoleID, req.PermissionIDs)
	if err != nil {
		return nil, err
	}

	// 3. 重新获取角色信息（包含更新后的权限）
	updatedRole, err := l.svcCtx.RoleModel.GetByID(req.RoleID)
	if err != nil {
		return nil, err
	}

	// 4. 构建响应
	resp = &types.RoleInfo{
		ID:          updatedRole.ID,
		Name:        updatedRole.Name,
		Code:        updatedRole.Code,
		Desc:        updatedRole.Desc,
		Status:      updatedRole.Status,
		Permissions: util.ConvertPermissions(updatedRole.Permissions),
	}

	return resp, nil
}

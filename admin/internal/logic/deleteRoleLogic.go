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

type DeleteRoleLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewDeleteRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteRoleLogic {
	return &DeleteRoleLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *DeleteRoleLogic) DeleteRole(req *types.DeleteRoleRequest) (resp *types.RoleInfo, err error) {
	// 1. 根据ID获取角色
	role, err := l.svcCtx.RoleModel.GetByID(req.ID)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, errors.New("角色不存在")
	}

	// 2. 保存角色信息用于响应
	roleInfo := &types.RoleInfo{
		ID:          role.ID,
		Name:        role.Name,
		Code:        role.Code,
		Desc:        role.Desc,
		Status:      role.Status,
		Permissions: util.ConvertPermissions(role.Permissions),
	}

	// 3. 删除角色
	err = l.svcCtx.RoleModel.Delete(req.ID)
	if err != nil {
		return nil, err
	}

	return roleInfo, nil
}

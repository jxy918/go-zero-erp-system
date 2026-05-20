// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"myproject/admin/internal/svc"
	"myproject/admin/internal/types"
	"myproject/admin/internal/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRoleLogic {
	return &GetRoleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetRoleLogic) GetRole(req *types.GetRoleRequest) (resp *types.RoleInfo, err error) {
	// 1. 获取角色信息
	role, err := l.svcCtx.RoleModel.GetByID(req.ID)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, nil
	}

	// 2. 构建响应
	resp = &types.RoleInfo{
		ID:          role.ID,
		Name:        role.Name,
		Code:        role.Code,
		Desc:        role.Desc,
		Status:      role.Status,
		Permissions: util.ConvertPermissions(role.Permissions),
		Menus:       util.ConvertMenus(role.Menus),
	}

	return resp, nil
}

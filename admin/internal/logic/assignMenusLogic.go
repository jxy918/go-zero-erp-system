package logic

import (
	"context"
	"errors"

	"myproject/admin/internal/svc"
	"myproject/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AssignMenusLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAssignMenusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AssignMenusLogic {
	return &AssignMenusLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AssignMenusLogic) AssignMenus(req *types.AssignMenusRequest) (resp *types.RoleInfo, err error) {
	// 1. 检查角色是否存在
	role, err := l.svcCtx.RoleModel.GetByID(req.RoleID)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, errors.New("角色不存在")
	}

	// 2. 为角色分配菜单
	err = l.svcCtx.RoleModel.AssignMenus(req.RoleID, req.MenuIDs)
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
		Permissions: nil, // 这里可以根据需要添加菜单信息
	}

	return resp, nil
}

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

type AssignRolesLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAssignRolesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AssignRolesLogic {
	return &AssignRolesLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AssignRolesLogic) AssignRoles(req *types.AssignRolesRequest) (resp *types.UserInfo, err error) {
	// 1. 检查用户是否存在
	user, err := l.svcCtx.UserModel.GetByID(req.UserID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("用户不存在")
	}

	// 2. 为用户分配角色
	err = l.svcCtx.UserModel.AssignRoles(req.UserID, req.RoleIDs)
	if err != nil {
		return nil, err
	}

	// 3. 重新获取用户信息（包含更新后的角色）
	updatedUser, err := l.svcCtx.UserModel.GetByID(req.UserID)
	if err != nil {
		return nil, err
	}

	// 4. 构建响应
	resp = &types.UserInfo{
		ID:       updatedUser.ID,
		Username: updatedUser.Username,
		Nickname: updatedUser.Nickname,
		Email:    updatedUser.Email,
		Phone:    updatedUser.Phone,
		Status:   updatedUser.Status,
		Roles:    util.ConvertRoles(updatedUser.Roles),
	}

	return resp, nil
}

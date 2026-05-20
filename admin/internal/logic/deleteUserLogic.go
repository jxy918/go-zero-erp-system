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

type DeleteUserLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewDeleteUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteUserLogic {
	return &DeleteUserLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *DeleteUserLogic) DeleteUser(req *types.DeleteUserRequest) (resp *types.UserInfo, err error) {
	// 1. 根据ID获取用户
	user, err := l.svcCtx.UserModel.GetByID(req.ID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("用户不存在")
	}

	// 2. 保存用户信息用于响应
	userInfo := &types.UserInfo{
		ID:       user.ID,
		Username: user.Username,
		Nickname: user.Nickname,
		Email:    user.Email,
		Phone:    user.Phone,
		Status:   user.Status,
		Roles:    util.ConvertRoles(user.Roles),
	}

	// 3. 删除用户之前先清除角色关联
	if err := l.svcCtx.UserModel.AssignRoles(req.ID, []uint{}); err != nil {
		l.Logger.Info("Failed to clear roles for user ", req.ID, ": ", err)
	}

	// 4. 删除用户
	err = l.svcCtx.UserModel.Delete(req.ID)
	if err != nil {
		return nil, err
	}

	return userInfo, nil
}

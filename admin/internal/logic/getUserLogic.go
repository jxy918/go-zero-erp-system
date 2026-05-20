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

type GetUserLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewGetUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserLogic {
	return &GetUserLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *GetUserLogic) GetUser(req *types.GetUserRequest) (resp *types.UserInfo, err error) {
	// 1. 根据ID获取用户
	user, err := l.svcCtx.UserModel.GetByID(req.ID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("用户不存在")
	}

	// 2. 构建响应
	resp = &types.UserInfo{
		ID:       user.ID,
		Username: user.Username,
		Nickname: user.Nickname,
		Email:    user.Email,
		Phone:    user.Phone,
		Status:   user.Status,
		Roles:    util.ConvertRoles(user.Roles),
	}

	return resp, nil
}

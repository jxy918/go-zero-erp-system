// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"myproject/admin/internal/model"
	"myproject/admin/internal/svc"
	"myproject/admin/internal/types"
	"myproject/admin/internal/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListUserLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
	ctx    context.Context
}

func NewListUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListUserLogic {
	return &ListUserLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
		ctx:    ctx,
	}
}

func (l *ListUserLogic) ListUser(req *types.ListUserRequest) (resp *types.ListUserResponse, err error) {
	// 1. 设置默认分页参数
	page := req.Page
	if page <= 0 {
		page = 1
	}

	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	// 2. 获取当前用户权限信息
	isAdmin := util.IsAdmin(l.ctx)
	roleID := util.GetRoleID(l.ctx)

	// 3. 根据权限获取用户列表
	var users []model.User
	var total int64
	if isAdmin {
		// 管理员可以查看所有用户
		users, total, err = l.svcCtx.UserModel.ListWithUsername(page, pageSize, req.Username)
	} else if roleID > 0 {
		// 普通用户只能查看自己角色范围内的用户（数据权限）
		users, total, err = l.svcCtx.UserModel.ListByRoleID(roleID, page, pageSize, req.Username)
	} else {
		// 如果角色ID为0，返回空列表
		users = []model.User{}
		total = 0
	}
	if err != nil {
		return nil, err
	}

	// 4. 构建响应
	var userInfos []types.UserInfo
	for _, user := range users {
		userInfos = append(userInfos, types.UserInfo{
			ID:       user.ID,
			Username: user.Username,
			Nickname: user.Nickname,
			Email:    user.Email,
			Phone:    user.Phone,
			Status:   user.Status,
			Roles:    util.ConvertRoles(user.Roles),
		})
	}

	resp = &types.ListUserResponse{
		Users: userInfos,
		Total: total,
	}

	return resp, nil
}

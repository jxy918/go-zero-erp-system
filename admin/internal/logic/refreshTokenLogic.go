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

type RefreshTokenLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewRefreshTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefreshTokenLogic {
	return &RefreshTokenLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *RefreshTokenLogic) RefreshToken(req *types.RefreshTokenRequest) (resp *types.RefreshTokenResponse, err error) {
	// 1. 解析旧token
	claims, err := util.ParseToken(req.Token, l.svcCtx.Config.JWT.AccessSecret)
	if err != nil {
		return nil, err
	}

	// 2. 生成新token
	newToken, err := util.GenerateToken(claims.UserID, claims.Username, claims.RoleID, claims.IsAdmin, l.svcCtx.Config.JWT.AccessSecret, l.svcCtx.Config.JWT.AccessExpire)
	if err != nil {
		return nil, err
	}

	// 3. 构建响应
	resp = &types.RefreshTokenResponse{
		Token: newToken,
	}

	return resp, nil
}

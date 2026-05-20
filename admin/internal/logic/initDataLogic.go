// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"myproject/admin/internal/model"
	"myproject/admin/internal/svc"
	"myproject/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type InitDataLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewInitDataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InitDataLogic {
	return &InitDataLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *InitDataLogic) InitData() (resp *types.InitDataResponse, err error) {
	// 调用模型层的初始化函数
	if err := model.InitData(); err != nil {
		return &types.InitDataResponse{
			Success: false,
			Message: "初始化失败: " + err.Error(),
		}, nil
	}

	return &types.InitDataResponse{
		Success: true,
		Message: "初始化成功",
	}, nil
}

// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"myproject/admin/internal/svc"
	"myproject/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetInventoryCheckLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetInventoryCheckLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetInventoryCheckLogic {
	return &GetInventoryCheckLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetInventoryCheckLogic) GetInventoryCheck(req *types.GetInventoryCheckRequest) (resp *types.InventoryCheckResponse, err error) {
	check, err := l.svcCtx.InventoryCheckModel.GetByID(req.ID)
	if err != nil {
		return nil, err
	}
	if check == nil {
		return &types.InventoryCheckResponse{}, nil
	}

	return &types.InventoryCheckResponse{
		Check: *convertInventoryCheckToInfo(l.svcCtx.InventoryModel, check),
	}, nil
}

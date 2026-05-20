// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"myproject/admin/internal/svc"
	"myproject/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteInventoryCheckLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteInventoryCheckLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteInventoryCheckLogic {
	return &DeleteInventoryCheckLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteInventoryCheckLogic) DeleteInventoryCheck(req *types.DeleteInventoryCheckRequest) (resp *types.InventoryCheckResponse, err error) {
	check, err := l.svcCtx.InventoryCheckModel.GetByID(req.ID)
	if err != nil {
		return nil, err
	}
	if check == nil {
		return &types.InventoryCheckResponse{}, nil
	}
	if check.Status != 1 {
		return &types.InventoryCheckResponse{}, nil
	}

	err = l.svcCtx.InventoryCheckModel.Delete(req.ID)
	if err != nil {
		return nil, err
	}

	return &types.InventoryCheckResponse{
		Check: *convertInventoryCheckToInfo(l.svcCtx.InventoryModel, check),
	}, nil
}

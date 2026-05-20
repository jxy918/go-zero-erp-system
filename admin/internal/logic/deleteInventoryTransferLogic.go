// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"myproject/admin/internal/svc"
	"myproject/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteInventoryTransferLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteInventoryTransferLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteInventoryTransferLogic {
	return &DeleteInventoryTransferLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteInventoryTransferLogic) DeleteInventoryTransfer(req *types.DeleteInventoryTransferRequest) (resp *types.InventoryTransferResponse, err error) {
	transfer, err := l.svcCtx.InventoryTransferModel.GetByID(uint(req.ID))
	if err != nil {
		return nil, err
	}
	if transfer == nil {
		return &types.InventoryTransferResponse{}, nil
	}

	if transfer.Status != 1 && transfer.Status != 4 {
		return &types.InventoryTransferResponse{}, nil
	}

	if err := l.svcCtx.InventoryTransferModel.Delete(uint(req.ID)); err != nil {
		return nil, err
	}

	return &types.InventoryTransferResponse{}, nil
}

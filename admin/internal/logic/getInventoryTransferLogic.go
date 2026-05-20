// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"myproject/admin/internal/svc"
	"myproject/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetInventoryTransferLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetInventoryTransferLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetInventoryTransferLogic {
	return &GetInventoryTransferLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetInventoryTransferLogic) GetInventoryTransfer(req *types.GetInventoryTransferRequest) (resp *types.InventoryTransferResponse, err error) {
	transfer, err := l.svcCtx.InventoryTransferModel.GetByID(uint(req.ID))
	if err != nil {
		return nil, err
	}
	if transfer == nil {
		return &types.InventoryTransferResponse{}, nil
	}

	return &types.InventoryTransferResponse{
		Transfer: *convertInventoryTransferToInfo(l.svcCtx.InventoryModel, transfer),
	}, nil
}

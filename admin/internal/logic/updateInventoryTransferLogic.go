// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"myproject/admin/internal/svc"
	"myproject/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateInventoryTransferLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateInventoryTransferLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateInventoryTransferLogic {
	return &UpdateInventoryTransferLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateInventoryTransferLogic) UpdateInventoryTransfer(req *types.UpdateInventoryTransferRequest) (resp *types.InventoryTransferResponse, err error) {
	transfer, err := l.svcCtx.InventoryTransferModel.GetByID(uint(req.ID))
	if err != nil {
		return nil, err
	}
	if transfer == nil {
		return &types.InventoryTransferResponse{}, nil
	}

	if transfer.Status != 1 {
		return &types.InventoryTransferResponse{}, nil
	}

	if req.FromWarehouseID != 0 && req.ToWarehouseID != 0 && req.FromWarehouseID == req.ToWarehouseID {
		return &types.InventoryTransferResponse{}, nil
	}

	if req.Quantity < 0 {
		return &types.InventoryTransferResponse{}, nil
	}

	if req.FromWarehouseID != 0 {
		transfer.FromWarehouseID = req.FromWarehouseID
	}
	if req.ToWarehouseID != 0 {
		transfer.ToWarehouseID = req.ToWarehouseID
	}
	if req.ProductID != 0 {
		transfer.ProductID = req.ProductID
	}
	if req.Quantity > 0 {
		transfer.Quantity = req.Quantity
	}
	if req.Remark != "" {
		transfer.Remark = req.Remark
	}

	if err := l.svcCtx.InventoryTransferModel.Update(transfer); err != nil {
		return nil, err
	}

	transfer, err = l.svcCtx.InventoryTransferModel.GetByID(transfer.ID)
	if err != nil {
		return nil, err
	}

	return &types.InventoryTransferResponse{
		Transfer: *convertInventoryTransferToInfo(l.svcCtx.InventoryModel, transfer),
	}, nil
}

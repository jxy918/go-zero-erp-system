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

type CreateInventoryTransferLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateInventoryTransferLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateInventoryTransferLogic {
	return &CreateInventoryTransferLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateInventoryTransferLogic) CreateInventoryTransfer(req *types.CreateInventoryTransferRequest) (resp *types.InventoryTransferResponse, err error) {
	operatorID, _ := l.ctx.Value(util.UserIDKey).(uint)
	if req.FromWarehouseID == req.ToWarehouseID {
		return &types.InventoryTransferResponse{}, nil
	}

	if req.Quantity <= 0 {
		return &types.InventoryTransferResponse{}, nil
	}

	fromWarehouse, err := l.svcCtx.WarehouseModel.GetByID(req.FromWarehouseID)
	if err != nil {
		return nil, err
	}
	if fromWarehouse == nil {
		return &types.InventoryTransferResponse{}, nil
	}

	toWarehouse, err := l.svcCtx.WarehouseModel.GetByID(req.ToWarehouseID)
	if err != nil {
		return nil, err
	}
	if toWarehouse == nil {
		return &types.InventoryTransferResponse{}, nil
	}

	product, err := l.svcCtx.ProductModel.GetByID(req.ProductID)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return &types.InventoryTransferResponse{}, nil
	}

	transferNo, err := l.svcCtx.InventoryTransferModel.GenerateTransferNo()
	if err != nil {
		return nil, err
	}

	transfer := &model.InventoryTransfer{
		TransferNo:      transferNo,
		FromWarehouseID: req.FromWarehouseID,
		ToWarehouseID:   req.ToWarehouseID,
		ProductID:       req.ProductID,
		Quantity:        req.Quantity,
		Status:          1,
		Remark:          req.Remark,
		CreatedBy:       operatorID,
	}

	if err := l.svcCtx.InventoryTransferModel.Create(transfer); err != nil {
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

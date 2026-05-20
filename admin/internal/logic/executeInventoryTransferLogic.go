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

type ExecuteInventoryTransferLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewExecuteInventoryTransferLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExecuteInventoryTransferLogic {
	return &ExecuteInventoryTransferLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ExecuteInventoryTransferLogic) ExecuteInventoryTransfer(req *types.ExecuteInventoryTransferRequest) (resp *types.InventoryTransferResponse, err error) {
	operatorID, _ := l.ctx.Value(util.UserIDKey).(uint)
	transfer, err := l.svcCtx.InventoryTransferModel.GetByID(uint(req.ID))
	if err != nil {
		return nil, err
	}
	if transfer == nil {
		return &types.InventoryTransferResponse{}, nil
	}

	if transfer.Status != 2 {
		return &types.InventoryTransferResponse{}, nil
	}

	if err := l.svcCtx.InventoryTransferModel.Execute(uint(req.ID), operatorID); err != nil {
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

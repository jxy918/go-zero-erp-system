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

type AuditInventoryTransferLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAuditInventoryTransferLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AuditInventoryTransferLogic {
	return &AuditInventoryTransferLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AuditInventoryTransferLogic) AuditInventoryTransfer(req *types.AuditInventoryTransferRequest) (resp *types.InventoryTransferResponse, err error) {
	operatorID, _ := l.ctx.Value(util.UserIDKey).(uint)
	transfer, err := l.svcCtx.InventoryTransferModel.GetByID(uint(req.ID))
	if err != nil {
		return &types.InventoryTransferResponse{
			Transfer: types.InventoryTransferInfo{},
		}, nil
	}
	if transfer == nil {
		return &types.InventoryTransferResponse{
			Transfer: types.InventoryTransferInfo{},
		}, nil
	}

	if transfer.Status != 1 {
		return &types.InventoryTransferResponse{
			Transfer: types.InventoryTransferInfo{},
		}, nil
	}

	if req.Status != 2 && req.Status != 4 {
		return &types.InventoryTransferResponse{
			Transfer: types.InventoryTransferInfo{},
		}, nil
	}

	if err := l.svcCtx.InventoryTransferModel.Audit(uint(req.ID), req.Status, operatorID); err != nil {
		return &types.InventoryTransferResponse{
			Transfer: types.InventoryTransferInfo{},
		}, nil
	}

	transfer, err = l.svcCtx.InventoryTransferModel.GetByID(transfer.ID)
	if err != nil {
		return &types.InventoryTransferResponse{
			Transfer: types.InventoryTransferInfo{},
		}, nil
	}

	return &types.InventoryTransferResponse{
		Transfer: *convertInventoryTransferToInfo(l.svcCtx.InventoryModel, transfer),
	}, nil
}

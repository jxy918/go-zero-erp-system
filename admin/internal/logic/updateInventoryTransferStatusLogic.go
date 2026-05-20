// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"errors"

	"myproject/admin/internal/svc"
	"myproject/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateInventoryTransferStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateInventoryTransferStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateInventoryTransferStatusLogic {
	return &UpdateInventoryTransferStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateInventoryTransferStatusLogic) UpdateInventoryTransferStatus(req *types.InventoryTransferStatusRequest) (resp *types.EmptyResponse, err error) {
	transfer, err := l.svcCtx.InventoryTransferModel.GetByID(req.ID)
	if err != nil {
		return nil, err
	}
	if transfer == nil {
		return nil, errors.New("调拨单不存在")
	}

	transfer.Status = req.Status
	if err := l.svcCtx.InventoryTransferModel.Update(transfer); err != nil {
		return nil, err
	}

	return &types.EmptyResponse{}, nil
}

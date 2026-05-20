// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"myproject/admin/internal/svc"
	"myproject/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListInventoryTransferLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListInventoryTransferLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListInventoryTransferLogic {
	return &ListInventoryTransferLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListInventoryTransferLogic) ListInventoryTransfer(req *types.ListInventoryTransferRequest) (resp *types.ListInventoryTransferResponse, err error) {
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}

	transfers, total, err := l.svcCtx.InventoryTransferModel.List(
		page, pageSize,
		req.TransferNo,
		uint(req.FromWarehouseID),
		uint(req.ToWarehouseID),
		uint(req.ProductID),
		uint(req.Status),
	)
	if err != nil {
		return nil, err
	}

	transferInfos := make([]types.InventoryTransferInfo, 0, len(transfers))
	for _, transfer := range transfers {
		transferInfos = append(transferInfos, *convertInventoryTransferToInfo(l.svcCtx.InventoryModel, &transfer))
	}

	return &types.ListInventoryTransferResponse{
		Transfers: transferInfos,
		Total:     total,
	}, nil
}

package logic

import (
	"context"

	"myproject/admin/internal/svc"
	"myproject/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCurrentStockLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewGetCurrentStockLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCurrentStockLogic {
	return &GetCurrentStockLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *GetCurrentStockLogic) GetCurrentStock(req *types.GetCurrentStockRequest) (resp *types.GetCurrentStockResponse, err error) {
	quantity := 0

	if req.ProductID > 0 && req.WarehouseID > 0 {
		total, err := l.svcCtx.InventoryModel.GetStockByProductAndWarehouse(req.ProductID, req.WarehouseID)
		if err != nil {
			return nil, err
		}
		quantity = total
	}

	return &types.GetCurrentStockResponse{
		ProductID:   req.ProductID,
		WarehouseID: req.WarehouseID,
		Quantity:    quantity,
	}, nil
}

package logic

import (
	"context"

	"myproject/admin/internal/svc"
	"myproject/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListWarehouseInventoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListWarehouseInventoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListWarehouseInventoryLogic {
	return &ListWarehouseInventoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListWarehouseInventoryLogic) ListWarehouseInventory(req *types.ListWarehouseInventoryRequest) (*types.ListWarehouseInventoryResponse, error) {
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}

	inventories, total, err := l.svcCtx.InventoryModel.ListWarehouseInventory(page, pageSize, req.ProductID, req.WarehouseID)
	if err != nil {
		return nil, err
	}

	var inventoryInfos []types.WarehouseInventoryInfo
	for _, inv := range inventories {
		info := types.WarehouseInventoryInfo{
			ID:          inv.ID,
			ProductID:   inv.ProductID,
			WarehouseID: inv.WarehouseID,
			Quantity:    inv.Quantity,
			UpdatedAt:   inv.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
		if inv.Product.ID > 0 {
			info.Product = inv.Product.Name
			info.ProductCode = inv.Product.Code
			info.Spec = inv.Product.Spec
		}
		if inv.Warehouse.ID > 0 {
			info.Warehouse = inv.Warehouse.Name
		}
		inventoryInfos = append(inventoryInfos, info)
	}

	return &types.ListWarehouseInventoryResponse{
		Inventories: inventoryInfos,
		Total:       total,
	}, nil
}

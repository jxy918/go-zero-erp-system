package logic

import (
	"context"

	"myproject/admin/internal/svc"
	"myproject/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListInventoryLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewListInventoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListInventoryLogic {
	return &ListInventoryLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *ListInventoryLogic) ListInventory(req *types.ListInventoryRequest) (resp *types.ListInventoryResponse, err error) {
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	inventories, total, err := l.svcCtx.InventoryModel.List(page, pageSize, req.ProductName, req.ProductID, req.WarehouseID)
	if err != nil {
		return nil, err
	}

	resp = &types.ListInventoryResponse{
		Inventory: make([]types.InventoryInfo, 0, len(inventories)),
		Total:     total,
	}

	for _, inv := range inventories {
		item := types.InventoryInfo{
			ID:             inv.ID,
			ProductID:      inv.ProductID,
			WarehouseID:    inv.WarehouseID,
			Quantity:       inv.Quantity,
			Type:           inv.Type,
			OrderID:        inv.OrderID,
			OrderType:      inv.OrderType,
			Remark:         "",
			BeforeQuantity: 0,
			AfterQuantity:  0,
			CreatedAt:      inv.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:      inv.UpdatedAt.Format("2006-01-02 15:04:05"),
		}

		if inv.Product.ID > 0 {
			stock, _ := l.svcCtx.InventoryModel.GetProductStockByWarehouse(inv.ProductID, inv.WarehouseID)

			mainUnit := ""
			if len(inv.Product.Units) > 0 {
				for _, u := range inv.Product.Units {
					if u.IsMain == 1 {
						mainUnit = u.UnitName
						break
					}
				}
				if mainUnit == "" {
					mainUnit = inv.Product.Units[0].UnitName
				}
			}

			item.Product = types.ProductInfo{
				ID:        inv.Product.ID,
				Name:      inv.Product.Name,
				Code:      inv.Product.Code,
				Price:     inv.Product.Price,
				Stock:     stock,
				MainUnit:  mainUnit,
				Status:    inv.Product.Status,
				CreatedAt: inv.Product.CreatedAt.Format("2006-01-02 15:04:05"),
				UpdatedAt: inv.Product.UpdatedAt.Format("2006-01-02 15:04:05"),
			}

			if len(inv.Product.Units) > 0 {
				item.Product.Units = make([]types.ProductUnitInfo, 0, len(inv.Product.Units))
				for _, u := range inv.Product.Units {
					item.Product.Units = append(item.Product.Units, types.ProductUnitInfo{
						ID:        u.ID,
						ProductID: u.ProductID,
						UnitName:  u.UnitName,
						Ratio:     u.Ratio,
						IsMain:    u.IsMain,
						CreatedAt: u.CreatedAt.Format("2006-01-02 15:04:05"),
						UpdatedAt: u.UpdatedAt.Format("2006-01-02 15:04:05"),
					})
				}
			}
		}

		if inv.Warehouse.ID > 0 {
			item.Warehouse = types.WarehouseInfo{
				ID:        inv.Warehouse.ID,
				Name:      inv.Warehouse.Name,
				Code:      inv.Warehouse.Code,
				Contact:   inv.Warehouse.Contact,
				Desc:      inv.Warehouse.Desc,
				Status:    inv.Warehouse.Status,
				CreatedAt: inv.Warehouse.CreatedAt.Format("2006-01-02 15:04:05"),
				UpdatedAt: inv.Warehouse.UpdatedAt.Format("2006-01-02 15:04:05"),
			}
		}

		resp.Inventory = append(resp.Inventory, item)
	}

	return resp, nil
}

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

type GetInventoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetInventoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetInventoryLogic {
	return &GetInventoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetInventoryLogic) GetInventory(req *types.GetInventoryRequest) (resp *types.InventoryInfo, err error) {
	inv, err := l.svcCtx.InventoryModel.GetByID(req.ID)
	if err != nil {
		return nil, err
	}
	if inv == nil {
		return nil, errors.New("库存记录不存在")
	}

	resp = &types.InventoryInfo{
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
		resp.Product = types.ProductInfo{
			ID:        inv.Product.ID,
			Name:      inv.Product.Name,
			Code:      inv.Product.Code,
			Price:     inv.Product.Price,
			Stock:     stock,
			MainUnit:  "",
			Status:    inv.Product.Status,
			CreatedAt: inv.Product.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: inv.Product.UpdatedAt.Format("2006-01-02 15:04:05"),
		}

		if len(inv.Product.Units) > 0 {
			resp.Product.Units = make([]types.ProductUnitInfo, 0, len(inv.Product.Units))
			for _, u := range inv.Product.Units {
				resp.Product.Units = append(resp.Product.Units, types.ProductUnitInfo{
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
		resp.Warehouse = types.WarehouseInfo{
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

	return resp, nil
}

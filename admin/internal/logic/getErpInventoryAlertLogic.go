package logic

import (
	"context"
	"time"

	"myproject/admin/internal/model"
	"myproject/admin/internal/svc"
	"myproject/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetErpInventoryAlertLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewGetErpInventoryAlertLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetErpInventoryAlertLogic {
	return &GetErpInventoryAlertLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *GetErpInventoryAlertLogic) GetErpInventoryAlert() (resp []*types.ErpInventoryAlertResponse, err error) {
	type inventoryWithStock struct {
		ProductID     uint
		WarehouseID   uint
		ProductName   string
		ProductCode   string
		WarehouseName string
		CurrentStock  int
		MinStock      int
		MaxStock      int
		SafetyStock   int
	}

	var inventoryData []inventoryWithStock

	err = l.svcCtx.DB.Model(&model.InventoryRecord{}).
		Select("inventory_records.product_id, inventory_records.warehouse_id, p.name as product_name, p.code as product_code, w.name as warehouse_name, COALESCE(SUM(inventory_records.quantity), 0) as current_stock, p.min_stock, p.max_stock, p.safety_stock").
		Joins("LEFT JOIN products p ON inventory_records.product_id = p.id").
		Joins("LEFT JOIN warehouses w ON inventory_records.warehouse_id = w.id").
		Where("p.status = 1").
		Where("w.status = 1").
		Group("inventory_records.product_id, inventory_records.warehouse_id, p.name, p.code, w.name, p.min_stock, p.max_stock, p.safety_stock").
		Having("(COALESCE(SUM(inventory_records.quantity), 0) <= p.safety_stock OR COALESCE(SUM(inventory_records.quantity), 0) <= p.min_stock OR COALESCE(SUM(inventory_records.quantity), 0) >= p.max_stock)").
		Find(&inventoryData).Error

	if err != nil {
		return nil, err
	}

	resp = make([]*types.ErpInventoryAlertResponse, 0, len(inventoryData))
	for _, data := range inventoryData {
		alertType := 0
		alertLevel := 0

		if data.CurrentStock <= data.MinStock && data.MinStock > 0 {
			alertType = 2
			alertLevel = 3
		} else if data.CurrentStock <= data.SafetyStock && data.SafetyStock > 0 {
			alertType = 1
			alertLevel = 2
		} else if data.CurrentStock >= data.MaxStock && data.MaxStock > 0 {
			alertType = 3
			alertLevel = 1
		}

		resp = append(resp, &types.ErpInventoryAlertResponse{
			ID:          data.ProductID,
			ProductID:   data.ProductID,
			ProductName: data.ProductName,
			ProductCode: data.ProductCode,
			WarehouseID: data.WarehouseID,
			Warehouse:   data.WarehouseName,
			Quantity:    data.CurrentStock,
			MinStock:    data.MinStock,
			MaxStock:    data.MaxStock,
			SafetyStock: data.SafetyStock,
			AlertType:   alertType,
			AlertLevel:  alertLevel,
			CreatedAt:   time.Now().Format("2006-01-02 15:04:05"),
		})
	}

	return resp, nil
}

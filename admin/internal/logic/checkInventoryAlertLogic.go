package logic

import (
	"context"
	"time"

	"myproject/admin/internal/model"
	"myproject/admin/internal/svc"
	"myproject/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CheckInventoryAlertLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCheckInventoryAlertLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckInventoryAlertLogic {
	return &CheckInventoryAlertLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CheckInventoryAlertLogic) CheckInventoryAlert() (resp *types.ListInventoryAlertResponse, err error) {
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

	alerts := make([]types.InventoryAlertInfo, 0, len(inventoryData))
	for _, data := range inventoryData {
		alertInfo := types.InventoryAlertInfo{
			ID:          0,
			ProductID:   data.ProductID,
			ProductName: data.ProductName,
			ProductCode: data.ProductCode,
			WarehouseID: data.WarehouseID,
			Warehouse:   data.WarehouseName,
			Quantity:    data.CurrentStock,
			MinStock:    data.MinStock,
			MaxStock:    data.MaxStock,
			SafetyStock: data.SafetyStock,
			AlertType:   0,
			AlertLevel:  0,
			CreatedAt:   time.Now().Format("2006-01-02 15:04:05"),
		}

		alerts = append(alerts, alertInfo)
	}

	return &types.ListInventoryAlertResponse{
		Alerts: alerts,
		Total:  int64(len(alerts)),
	}, nil
}

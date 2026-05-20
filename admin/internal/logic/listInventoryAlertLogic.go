package logic

import (
	"context"
	"time"

	"myproject/admin/internal/model"
	"myproject/admin/internal/svc"
	"myproject/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListInventoryAlertLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListInventoryAlertLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListInventoryAlertLogic {
	return &ListInventoryAlertLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListInventoryAlertLogic) ListInventoryAlert(req *types.ListInventoryAlertRequest) (resp *types.ListInventoryAlertResponse, err error) {
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}

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
	var total int64

	query := l.svcCtx.DB.Model(&model.InventoryRecord{}).
		Select("inventory_records.product_id, inventory_records.warehouse_id, p.name as product_name, p.code as product_code, w.name as warehouse_name, COALESCE(SUM(inventory_records.quantity), 0) as current_stock, p.min_stock, p.max_stock, p.safety_stock").
		Joins("LEFT JOIN products p ON inventory_records.product_id = p.id").
		Joins("LEFT JOIN warehouses w ON inventory_records.warehouse_id = w.id").
		Where("p.status = 1").
		Where("w.status = 1").
		Group("inventory_records.product_id, inventory_records.warehouse_id, p.name, p.code, w.name, p.min_stock, p.max_stock, p.safety_stock").
		Having("(COALESCE(SUM(inventory_records.quantity), 0) <= p.safety_stock OR COALESCE(SUM(inventory_records.quantity), 0) <= p.min_stock OR COALESCE(SUM(inventory_records.quantity), 0) >= p.max_stock)")

	if req.AlertType > 0 {
		switch req.AlertType {
		case 1:
			query = query.Having("COALESCE(SUM(inventory_records.quantity), 0) <= p.safety_stock AND COALESCE(SUM(inventory_records.quantity), 0) > p.min_stock")
		case 2:
			query = query.Having("COALESCE(SUM(inventory_records.quantity), 0) <= p.min_stock")
		case 3:
			query = query.Having("COALESCE(SUM(inventory_records.quantity), 0) >= p.max_stock")
		}
	}

	err = query.Count(&total).Error
	if err != nil {
		return nil, err
	}

	err = query.Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&inventoryData).Error
	if err != nil {
		return nil, err
	}

	alerts := make([]types.InventoryAlertInfo, 0, len(inventoryData))
	for _, data := range inventoryData {
		alertInfo := types.InventoryAlertInfo{
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

		if data.CurrentStock <= data.MinStock && data.MinStock > 0 {
			alertInfo.AlertType = 2
			alertInfo.AlertLevel = 3
		} else if data.CurrentStock <= data.SafetyStock && data.SafetyStock > 0 {
			alertInfo.AlertType = 1
			alertInfo.AlertLevel = 2
		} else if data.CurrentStock >= data.MaxStock && data.MaxStock > 0 {
			alertInfo.AlertType = 3
			alertInfo.AlertLevel = 1
		}

		alerts = append(alerts, alertInfo)
	}

	return &types.ListInventoryAlertResponse{
		Alerts: alerts,
		Total:  total,
	}, nil
}

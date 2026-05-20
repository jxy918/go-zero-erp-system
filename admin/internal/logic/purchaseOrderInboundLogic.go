package logic

import (
	"context"
	"errors"

	"myproject/admin/internal/metric"
	"myproject/admin/internal/model"
	"myproject/admin/internal/svc"
	"myproject/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type PurchaseOrderInboundLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPurchaseOrderInboundLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PurchaseOrderInboundLogic {
	return &PurchaseOrderInboundLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PurchaseOrderInboundLogic) PurchaseOrderInbound(req *types.PurchaseInboundRequest) (resp *types.PurchaseOrderInfo, err error) {
	if req.OrderID == 0 {
		return nil, errors.New("采购订单不能为空")
	}

	if req.WarehouseID == 0 {
		return nil, errors.New("仓库不能为空")
	}

	order, err := l.svcCtx.PurchaseModel.GetByID(req.OrderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, errors.New("采购订单不存在")
	}

	if order.Status != 2 {
		return nil, errors.New("订单状态不正确，请先审核订单")
	}

	warehouse, err := l.svcCtx.WarehouseModel.GetByID(req.WarehouseID)
	if err != nil {
		return nil, err
	}
	if warehouse == nil {
		return nil, errors.New("仓库不存在")
	}

	err = l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		for _, item := range order.Items {
			baseQty := item.BaseQty
			if baseQty == 0 {
				baseQty = item.Quantity
			}

			var beforeQty int
			err := tx.Model(&model.WarehouseInventory{}).
				Select("quantity").
				Where("product_id = ? AND warehouse_id = ?", item.ProductID, req.WarehouseID).
				Scan(&beforeQty).Error
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}

			var inv model.WarehouseInventory
			result := tx.Where("product_id = ? AND warehouse_id = ?", item.ProductID, req.WarehouseID).First(&inv)
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				inv = model.WarehouseInventory{
					ProductID:   item.ProductID,
					WarehouseID: req.WarehouseID,
					Quantity:    baseQty,
				}
				if err := tx.Create(&inv).Error; err != nil {
					return err
				}
			} else {
				if err := tx.Model(&model.WarehouseInventory{}).
					Where("product_id = ? AND warehouse_id = ?", item.ProductID, req.WarehouseID).
					Update("quantity", gorm.Expr("quantity + ?", baseQty)).Error; err != nil {
					return err
				}
			}

			var currentStock int
			err = tx.Model(&model.InventoryRecord{}).
				Select("COALESCE(SUM(quantity), 0)").
				Where("product_id = ? AND warehouse_id = ?", item.ProductID, req.WarehouseID).
				Scan(&currentStock).Error
			if err != nil {
				return err
			}

			invRecord := model.InventoryRecord{
				ProductID:   item.ProductID,
				WarehouseID: req.WarehouseID,
				Quantity:    baseQty,
				Type:        1,
				OrderID:     req.OrderID,
				OrderType:   1,
			}
			if err := tx.Create(&invRecord).Error; err != nil {
				return err
			}

			change := &model.InventoryChange{
				ProductID:      item.ProductID,
				WarehouseID:    req.WarehouseID,
				BeforeQuantity: beforeQty,
				AfterQuantity:  beforeQty + baseQty,
				Quantity:       baseQty,
				Type:           1,
				OrderID:        req.OrderID,
				OrderType:      1,
				Remark:         "采购入库",
			}
			if err := tx.Create(change).Error; err != nil {
				return err
			}
		}

		order.Status = 3
		if err := tx.Save(order).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		if metric.IsEnabled() {
			metric.OrderCreateCounter.Inc("purchase", "failure")
		}
		return nil, err
	}

	// 记录采购入库成功指标
	if metric.IsEnabled() {
		metric.OrderCreateCounter.Inc("purchase", "success")
		metric.InventoryAdjustCounter.Inc("inbound")
	}

	// 重新查询更新后的订单信息
	updatedOrder, err := l.svcCtx.PurchaseModel.GetByID(req.OrderID)
	if err != nil {
		return nil, err
	}

	// 构建响应
	resp = &types.PurchaseOrderInfo{
		ID:          updatedOrder.ID,
		OrderNo:     updatedOrder.OrderNo,
		SupplierID:  updatedOrder.SupplierID,
		WarehouseID: updatedOrder.WarehouseID,
		TotalAmount: updatedOrder.TotalAmount,
		Status:      updatedOrder.Status,
		Remark:      updatedOrder.Remark,
		CreatedAt:   updatedOrder.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   updatedOrder.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	return resp, nil
}

package logic

import (
	"context"
	"errors"

	"myproject/admin/internal/middleware"
	"myproject/admin/internal/model"
	"myproject/admin/internal/svc"
	"myproject/admin/internal/types"
	"myproject/admin/internal/util"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type PurchaseInboundLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPurchaseInboundLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PurchaseInboundLogic {
	return &PurchaseInboundLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PurchaseInboundLogic) PurchaseInbound(req *types.PurchaseInboundRequest) (resp *types.PurchaseOrderInfo, err error) {
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

	err = model.DB.Transaction(func(tx *gorm.DB) error {
		for _, item := range order.Items {
			baseQty := item.BaseQty
			if baseQty == 0 {
				baseQty = item.Quantity
			}

			if err := tx.Model(&model.Product{}).Where("id = ?", item.ProductID).Update("stock", gorm.Expr("stock + ?", baseQty)).Error; err != nil {
				return err
			}

			var currentStock int
			err := tx.Model(&model.InventoryRecord{}).
				Select("COALESCE(SUM(quantity), 0)").
				Where("product_id = ? AND warehouse_id = ?", item.ProductID, req.WarehouseID).
				Scan(&currentStock).Error
			if err != nil {
				return err
			}

			inv := model.InventoryRecord{
				ProductID:   item.ProductID,
				WarehouseID: req.WarehouseID,
				Quantity:    baseQty,
				Type:        1,
				OrderID:     req.OrderID,
				OrderType:   1,
			}
			if err := tx.Create(&inv).Error; err != nil {
				return err
			}

			change := &model.InventoryChange{
				ProductID:      item.ProductID,
				WarehouseID:    req.WarehouseID,
				BeforeQuantity: currentStock,
				AfterQuantity:  currentStock + baseQty,
				Quantity:       baseQty,
				Type:           1,
				OrderID:        req.OrderID,
				OrderType:      1,
				Remark:         "采购入库",
			}
			if err := tx.Create(change).Error; err != nil {
				return err
			}

			var warehouseInventory struct {
				ID uint
			}
			err = tx.Table("warehouse_inventories").
				Select("id").
				Where("product_id = ? AND warehouse_id = ?", item.ProductID, req.WarehouseID).
				Take(&warehouseInventory).Error
			if err == gorm.ErrRecordNotFound {
				err = tx.Table("warehouse_inventories").Create(map[string]interface{}{
					"product_id":   item.ProductID,
					"warehouse_id": req.WarehouseID,
					"quantity":     baseQty,
				}).Error
				if err != nil {
					return err
				}
			} else if err == nil {
				err = tx.Table("warehouse_inventories").
					Where("id = ?", warehouseInventory.ID).
					Update("quantity", gorm.Expr("quantity + ?", baseQty)).Error
				if err != nil {
					return err
				}
			} else {
				return err
			}
		}

		return tx.Model(&model.PurchaseOrder{}).Where("id = ?", req.OrderID).Update("status", 3).Error
	})

	if err != nil {
		return nil, err
	}

	operatorID, _ := l.ctx.Value(util.UserIDKey).(uint)
	operatorName := middleware.GetUsername(l.ctx)

	log := &model.OrderLog{
		OrderID:      req.OrderID,
		OrderType:    1,
		BeforeStatus: 2,
		AfterStatus:  3,
		OperatorID:   operatorID,
		OperatorName: operatorName,
		Remark:       "采购订单入库",
	}
	if err := l.svcCtx.OrderLogModel.Insert(log); err != nil {
		return nil, err
	}

	order, _ = l.svcCtx.PurchaseModel.GetByID(req.OrderID)
	resp = &types.PurchaseOrderInfo{
		ID:          order.ID,
		OrderNo:     order.OrderNo,
		SupplierID:  order.SupplierID,
		TotalAmount: order.TotalAmount,
		Status:      order.Status,
		Remark:      order.Remark,
		CreatedAt:   order.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   order.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	if order.Supplier.ID > 0 {
		resp.Supplier = types.SupplierInfo{
			ID:        order.Supplier.ID,
			Name:      order.Supplier.Name,
			Contact:   order.Supplier.Contact,
			Phone:     order.Supplier.Phone,
			Email:     order.Supplier.Email,
			Address:   order.Supplier.Address,
			Status:    order.Supplier.Status,
			CreatedAt: order.Supplier.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: order.Supplier.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	return resp, nil
}

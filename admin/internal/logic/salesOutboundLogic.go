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

type SalesOutboundLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSalesOutboundLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SalesOutboundLogic {
	return &SalesOutboundLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SalesOutboundLogic) SalesOutbound(req *types.SalesOutboundRequest) (resp *types.SalesOrderInfo, err error) {
	if req.OrderID == 0 {
		return nil, errors.New("销售订单不能为空")
	}

	if req.WarehouseID == 0 {
		return nil, errors.New("仓库不能为空")
	}

	order, err := l.svcCtx.SalesModel.GetByID(req.OrderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, errors.New("销售订单不存在")
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

			var currentStock int
			err := tx.Model(&model.InventoryRecord{}).
				Select("COALESCE(SUM(quantity), 0)").
				Where("product_id = ? AND warehouse_id = ?", item.ProductID, req.WarehouseID).
				Scan(&currentStock).Error
			if err != nil {
				return err
			}

			if currentStock < baseQty {
				return errors.New("该仓库库存不足")
			}

			if err := tx.Model(&model.Product{}).Where("id = ?", item.ProductID).Update("stock", gorm.Expr("stock - ?", baseQty)).Error; err != nil {
				return err
			}

			inv := model.InventoryRecord{
				ProductID:   item.ProductID,
				WarehouseID: req.WarehouseID,
				Quantity:    -baseQty,
				Type:        2,
				OrderID:     req.OrderID,
				OrderType:   2,
			}
			if err := tx.Create(&inv).Error; err != nil {
				return err
			}

			change := &model.InventoryChange{
				ProductID:      item.ProductID,
				WarehouseID:    req.WarehouseID,
				BeforeQuantity: currentStock,
				AfterQuantity:  currentStock - baseQty,
				Quantity:       -baseQty,
				Type:           2,
				OrderID:        req.OrderID,
				OrderType:      2,
				Remark:         "销售出库",
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
			if err == nil {
				err = tx.Table("warehouse_inventories").
					Where("id = ?", warehouseInventory.ID).
					Update("quantity", gorm.Expr("quantity - ?", baseQty)).Error
				if err != nil {
					return err
				}
			} else if err != gorm.ErrRecordNotFound {
				return err
			}
		}

		return tx.Model(&model.SalesOrder{}).Where("id = ?", req.OrderID).Update("status", 3).Error
	})

	if err != nil {
		return nil, err
	}

	operatorID, _ := l.ctx.Value(util.UserIDKey).(uint)
	operatorName := middleware.GetUsername(l.ctx)

	log := &model.OrderLog{
		OrderID:      req.OrderID,
		OrderType:    2,
		BeforeStatus: 2,
		AfterStatus:  3,
		OperatorID:   operatorID,
		OperatorName: operatorName,
		Remark:       "销售订单出库",
	}
	if err := l.svcCtx.OrderLogModel.Insert(log); err != nil {
		return nil, err
	}

	order, _ = l.svcCtx.SalesModel.GetByID(req.OrderID)
	resp = &types.SalesOrderInfo{
		ID:          order.ID,
		OrderNo:     order.OrderNo,
		CustomerID:  order.CustomerID,
		TotalAmount: order.TotalAmount,
		Status:      order.Status,
		Remark:      order.Remark,
		CreatedAt:   order.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   order.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	if order.Customer.ID > 0 {
		resp.Customer = types.CustomerInfo{
			ID:        order.Customer.ID,
			Name:      order.Customer.Name,
			Contact:   order.Customer.Contact,
			Phone:     order.Customer.Phone,
			Email:     order.Customer.Email,
			Address:   order.Customer.Address,
			Status:    order.Customer.Status,
			CreatedAt: order.Customer.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: order.Customer.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	return resp, nil
}

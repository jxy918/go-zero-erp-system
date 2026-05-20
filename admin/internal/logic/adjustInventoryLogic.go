package logic

import (
	"context"
	"errors"

	"myproject/admin/internal/model"
	"myproject/admin/internal/svc"
	"myproject/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type AdjustInventoryLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAdjustInventoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdjustInventoryLogic {
	return &AdjustInventoryLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *AdjustInventoryLogic) AdjustInventory(req *types.AdjustInventoryRequest) (resp *types.InventoryInfo, err error) {
	if req.ProductID == 0 {
		return nil, errors.New("产品不能为空")
	}

	if req.WarehouseID == 0 {
		return nil, errors.New("仓库不能为空")
	}

	if req.Quantity == 0 {
		return nil, errors.New("调整数量不能为0")
	}

	product, err := l.svcCtx.ProductModel.GetByID(req.ProductID)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, errors.New("产品不存在")
	}

	warehouse, err := l.svcCtx.WarehouseModel.GetByID(req.WarehouseID)
	if err != nil {
		return nil, err
	}
	if warehouse == nil {
		return nil, errors.New("仓库不存在")
	}

	quantity := req.Quantity

	var adjustType int
	if quantity > 0 {
		adjustType = 1
	} else {
		adjustType = 2
	}

	var beforeStock, afterStock int

	err = model.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.Product{}).Where("id = ?", req.ProductID).Update("stock", gorm.Expr("stock + ?", quantity)).Error; err != nil {
			return err
		}

		var warehouseInventory struct {
			ID       uint
			Quantity int
		}
		err = tx.Table("warehouse_inventories").
			Select("id, quantity").
			Where("product_id = ? AND warehouse_id = ?", req.ProductID, req.WarehouseID).
			Take(&warehouseInventory).Error

		if err == gorm.ErrRecordNotFound {
			beforeStock = 0
			afterStock = quantity
			if afterStock < 0 {
				return errors.New("库存不足")
			}
			err = tx.Table("warehouse_inventories").Create(map[string]interface{}{
				"product_id":   req.ProductID,
				"warehouse_id": req.WarehouseID,
				"quantity":     quantity,
			}).Error
			if err != nil {
				return err
			}
		} else if err == nil {
			beforeStock = warehouseInventory.Quantity
			if quantity < 0 && beforeStock < -quantity {
				return errors.New("库存不足")
			}
			afterStock = beforeStock + quantity
			err = tx.Table("warehouse_inventories").
				Where("id = ?", warehouseInventory.ID).
				Update("quantity", afterStock).Error
			if err != nil {
				return err
			}
		} else {
			return err
		}

		inv := model.InventoryRecord{
			ProductID:   req.ProductID,
			WarehouseID: req.WarehouseID,
			Quantity:    quantity,
			Type:        3,
		}
		return tx.Create(&inv).Error
	})

	if err != nil {
		return nil, err
	}

	change := &model.InventoryChange{
		ProductID:      req.ProductID,
		WarehouseID:    req.WarehouseID,
		BeforeQuantity: beforeStock,
		AfterQuantity:  afterStock,
		Quantity:       quantity,
		Type:           adjustType,
		Remark:         req.Remark,
	}
	if err := l.svcCtx.InventoryChangeModel.Create(change); err != nil {
		return nil, err
	}

	resp = &types.InventoryInfo{
		ID:             0,
		ProductID:      req.ProductID,
		WarehouseID:    req.WarehouseID,
		Quantity:       quantity,
		Type:           adjustType,
		OrderType:      0,
		Remark:         req.Remark,
		BeforeQuantity: beforeStock,
		AfterQuantity:  afterStock,
	}

	resp.Product = types.ProductInfo{
		ID:        product.ID,
		Name:      product.Name,
		Code:      product.Code,
		CostPrice: product.CostPrice,
		Stock:     beforeStock + quantity,
		Status:    product.Status,
	}

	resp.Warehouse = types.WarehouseInfo{
		ID:      warehouse.ID,
		Name:    warehouse.Name,
		Code:    warehouse.Code,
		Contact: warehouse.Contact,
		Desc:    warehouse.Desc,
		Status:  warehouse.Status,
	}

	return resp, nil
}

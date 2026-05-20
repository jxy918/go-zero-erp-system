package model

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type PurchaseModel struct {
	db *gorm.DB
}

func NewPurchaseModel(db *gorm.DB) *PurchaseModel {
	return &PurchaseModel{db: db}
}

func generateOrderNo(prefix string) string {
	return fmt.Sprintf("%s%s%d", prefix, time.Now().Format("20060102"), time.Now().UnixNano()%100000000)
}

func (m *PurchaseModel) Create(order *PurchaseOrder, items []PurchaseOrderItem) error {
	return m.db.Transaction(func(tx *gorm.DB) error {
		order.OrderNo = generateOrderNo("PO")
		if err := tx.Create(order).Error; err != nil {
			return err
		}
		for i := range items {
			items[i].OrderID = order.ID
			if err := tx.Create(&items[i]).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (m *PurchaseModel) Update(order *PurchaseOrder) error {
	return m.db.Model(order).Updates(map[string]interface{}{
		"supplier_id":  order.SupplierID,
		"warehouse_id": order.WarehouseID,
		"status":       order.Status,
		"remark":       order.Remark,
	}).Error
}

func (m *PurchaseModel) Delete(id uint) error {
	return m.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Unscoped().Where("order_id = ?", id).Delete(&PurchaseOrderItem{}).Error; err != nil {
			return err
		}
		return tx.Unscoped().Delete(&PurchaseOrder{}, id).Error
	})
}

func (m *PurchaseModel) GetByID(id uint) (*PurchaseOrder, error) {
	var order PurchaseOrder
	err := m.db.Preload("Supplier").Preload("Warehouse").Preload("Items.Product").First(&order, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &order, nil
}

func (m *PurchaseModel) List(page, pageSize int, orderNo string, status int) ([]PurchaseOrder, int64, error) {
	var orders []PurchaseOrder
	var total int64

	offset := (page - 1) * pageSize

	query := m.db.Model(&PurchaseOrder{})

	if orderNo != "" {
		query = query.Where("order_no LIKE ?", "%"+orderNo+"%")
	}

	if status > 0 {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Preload("Supplier").Preload("Warehouse").Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func (m *PurchaseModel) UpdateStatus(id uint, status int) error {
	return m.db.Model(&PurchaseOrder{}).Where("id = ?", id).Update("status", status).Error
}

// PurchaseInStock 采购入库
func (m *PurchaseModel) PurchaseInStock(orderID uint) error {
	return m.db.Transaction(func(tx *gorm.DB) error {
		// 1. 获取采购订单信息
		var order PurchaseOrder
		if err := tx.Preload("Items").First(&order, orderID).Error; err != nil {
			return err
		}

		// 2. 检查订单状态
		if order.Status != 2 { // 2: 已审核
			return errors.New("订单状态不正确，只有已审核的订单才能入库")
		}

		// 3. 更新库存并记录变动
		for _, item := range order.Items {
			// 获取当前库存（从流水记录计算）
			var currentStock int
			err := tx.Model(&InventoryRecord{}).
				Select("COALESCE(SUM(quantity), 0)").
				Where("product_id = ? AND warehouse_id = ?", item.ProductID, order.WarehouseID).
				Scan(&currentStock).Error
			if err != nil {
				return err
			}

			// 记录库存变动历史
			change := InventoryChange{
				ProductID:      item.ProductID,
				WarehouseID:    order.WarehouseID,
				BeforeQuantity: currentStock,
				AfterQuantity:  currentStock + item.Quantity,
				Quantity:       item.Quantity,
				Type:           1, // 1: 入库
				OrderID:        orderID,
				OrderType:      1, // 1: 采购订单
				Remark:         "采购入库",
			}
			if err := tx.Create(&change).Error; err != nil {
				return err
			}

			// 记录库存流水
			record := InventoryRecord{
				ProductID:   item.ProductID,
				WarehouseID: order.WarehouseID,
				Quantity:    item.Quantity,
				Type:        1, // 1: 入库
				OrderID:     orderID,
				OrderType:   1, // 1: 采购订单
			}
			if err := tx.Create(&record).Error; err != nil {
				return err
			}
		}

		// 4. 更新订单状态为已入库
		return tx.Model(&PurchaseOrder{}).Where("id = ?", orderID).Update("status", 3).Error // 3: 已入库
	})
}

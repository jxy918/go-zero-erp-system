package model

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type SalesModel struct {
	db *gorm.DB
}

func NewSalesModel(db *gorm.DB) *SalesModel {
	return &SalesModel{db: db}
}

func generateSalesOrderNo() string {
	return fmt.Sprintf("%s%s%d", "SO", time.Now().Format("20060102"), time.Now().UnixNano()%100000000)
}

func (m *SalesModel) Create(order *SalesOrder, items []SalesOrderItem) error {
	return m.db.Transaction(func(tx *gorm.DB) error {
		order.OrderNo = generateSalesOrderNo()
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

func (m *SalesModel) Update(order *SalesOrder) error {
	return m.db.Model(order).Updates(map[string]interface{}{
		"customer_id":  order.CustomerID,
		"warehouse_id": order.WarehouseID,
		"status":       order.Status,
		"remark":       order.Remark,
	}).Error
}

func (m *SalesModel) UpdateFields(id uint, fields map[string]interface{}) error {
	return m.db.Model(&SalesOrder{}).Where("id = ?", id).Updates(fields).Error
}

func (m *SalesModel) Delete(id uint) error {
	return m.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Unscoped().Where("order_id = ?", id).Delete(&SalesOrderItem{}).Error; err != nil {
			return err
		}
		return tx.Unscoped().Delete(&SalesOrder{}, id).Error
	})
}

func (m *SalesModel) GetByID(id uint) (*SalesOrder, error) {
	var order SalesOrder
	err := m.db.Preload("Customer").Preload("Warehouse").Preload("Items.Product").First(&order, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &order, nil
}

func (m *SalesModel) List(page, pageSize int, orderNo string, status int) ([]SalesOrder, int64, error) {
	var orders []SalesOrder
	var total int64

	offset := (page - 1) * pageSize

	query := m.db.Model(&SalesOrder{})

	if orderNo != "" {
		query = query.Where("order_no LIKE ?", "%"+orderNo+"%")
	}

	if status > 0 {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Preload("Customer").Preload("Warehouse").Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func (m *SalesModel) UpdateStatus(id uint, status int) error {
	return m.db.Model(&SalesOrder{}).Where("id = ?", id).Update("status", status).Error
}

// SalesOutStock 销售出库
func (m *SalesModel) SalesOutStock(orderID uint) error {
	return m.db.Transaction(func(tx *gorm.DB) error {
		// 1. 获取销售订单信息
		var order SalesOrder
		if err := tx.Preload("Items").First(&order, orderID).Error; err != nil {
			return err
		}

		// 2. 检查订单状态
		if order.Status != 2 { // 2: 已审核
			return errors.New("订单状态不正确，只有已审核的订单才能出库")
		}

		// 3. 检查库存并扣减
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

			// 检查库存是否充足
			if currentStock < item.Quantity {
				return errors.New("库存不足，无法出库")
			}

			// 记录库存变动历史
			change := InventoryChange{
				ProductID:      item.ProductID,
				WarehouseID:    order.WarehouseID,
				BeforeQuantity: currentStock,
				AfterQuantity:  currentStock - item.Quantity,
				Quantity:       -item.Quantity,
				Type:           2, // 2: 出库
				OrderID:        orderID,
				OrderType:      2, // 2: 销售订单
				Remark:         "销售出库",
			}
			if err := tx.Create(&change).Error; err != nil {
				return err
			}

			// 记录库存流水
			record := InventoryRecord{
				ProductID:   item.ProductID,
				WarehouseID: order.WarehouseID,
				Quantity:    -item.Quantity,
				Type:        2, // 2: 出库
				OrderID:     orderID,
				OrderType:   2, // 2: 销售订单
			}
			if err := tx.Create(&record).Error; err != nil {
				return err
			}
		}

		// 4. 更新订单状态为已出库
		return tx.Model(&SalesOrder{}).Where("id = ?", orderID).Update("status", 3).Error // 3: 已出库
	})
}

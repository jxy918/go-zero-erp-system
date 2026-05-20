package model

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type InventoryModel struct {
	db *gorm.DB
}

func NewInventoryModel(db *gorm.DB) *InventoryModel {
	return &InventoryModel{db: db}
}

func (m *InventoryModel) Create(record *InventoryRecord) error {
	return m.db.Create(record).Error
}

func (m *InventoryModel) GetByID(id uint) (*InventoryRecord, error) {
	var record InventoryRecord
	err := m.db.Preload("Product").Preload("Warehouse").First(&record, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &record, nil
}

func (m *InventoryModel) List(page, pageSize int, productName string, productID, warehouseID uint) ([]InventoryRecord, int64, error) {
	var records []InventoryRecord
	var total int64

	offset := (page - 1) * pageSize

	query := m.db.Model(&InventoryRecord{})

	if productName != "" {
		query = query.Joins("JOIN products ON products.id = inventory_records.product_id").
			Where("products.name LIKE ?", "%"+productName+"%")
	}

	if productID > 0 {
		query = query.Where("product_id = ?", productID)
	}

	if warehouseID > 0 {
		query = query.Where("warehouse_id = ?", warehouseID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Preload("Product").Preload("Product.Units").Preload("Warehouse").Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&records).Error; err != nil {
		return nil, 0, err
	}

	return records, total, nil
}

func (m *InventoryModel) GetStockByProduct(productID uint) (int, error) {
	var total int
	err := m.db.Model(&InventoryRecord{}).
		Select("COALESCE(SUM(quantity), 0)").
		Where("product_id = ?", productID).
		Scan(&total).Error
	return total, err
}

func (m *InventoryModel) GetStockByProductAndWarehouse(productID, warehouseID uint) (int, error) {
	var total int
	err := m.db.Model(&InventoryRecord{}).
		Select("COALESCE(SUM(quantity), 0)").
		Where("product_id = ? AND warehouse_id = ?", productID, warehouseID).
		Scan(&total).Error
	return total, err
}

func (m *InventoryModel) GetHistory(productID, warehouseID uint) ([]InventoryRecord, error) {
	var records []InventoryRecord
	query := m.db.Preload("Product").Preload("Product.Units").Preload("Warehouse")

	if productID > 0 {
		query = query.Where("product_id = ?", productID)
	}
	if warehouseID > 0 {
		query = query.Where("warehouse_id = ?", warehouseID)
	}

	err := query.Order("created_at DESC").Find(&records).Error
	if err != nil {
		return nil, err
	}
	return records, nil
}

// AdjustStock 调整库存（增量调整）
func (m *InventoryModel) AdjustStock(productID, warehouseID uint, quantity int, orderID uint, orderType int, remark string) error {
	fmt.Println("==========================================")
	fmt.Println("=== AdjustStock START ===")
	fmt.Printf("输入参数:\n")
	fmt.Printf("  productID: %d\n", productID)
	fmt.Printf("  warehouseID: %d\n", warehouseID)
	fmt.Printf("  quantity: %d\n", quantity)
	fmt.Printf("  orderID: %d\n", orderID)
	fmt.Printf("  orderType: %d (1=入库, 2=出库, 3=调整)\n", orderType)
	fmt.Printf("  remark: %s\n", remark)

	return m.db.Transaction(func(tx *gorm.DB) error {
		// 获取当前库存（从流水记录计算）
		var currentStock int
		err := tx.Model(&InventoryRecord{}).
			Select("COALESCE(SUM(quantity), 0)").
			Where("product_id = ? AND warehouse_id = ?", productID, warehouseID).
			Scan(&currentStock).Error
		if err != nil {
			return err
		}

		// 检查调整后的库存是否为负数
		if currentStock+quantity < 0 {
			return errors.New("库存调整后不能为负数")
		}

		// 记录库存变动历史
		change := InventoryChange{
			ProductID:      productID,
			WarehouseID:    warehouseID,
			BeforeQuantity: currentStock,
			AfterQuantity:  currentStock + quantity,
			Quantity:       quantity,
			Type:           orderType, // 1:入库, 2:出库, 3:调整
			OrderID:        orderID,
			OrderType:      4, // 4: 库存调整申请
			Remark:         remark,
		}
		fmt.Println("  -> 准备写入 inventory_changes:")
		fmt.Printf("     Type: %d\n", change.Type)
		fmt.Printf("     Quantity: %d\n", change.Quantity)
		if err := tx.Create(&change).Error; err != nil {
			return err
		}

		// 记录库存流水
		record := InventoryRecord{
			ProductID:   productID,
			WarehouseID: warehouseID,
			Quantity:    quantity,
			Type:        orderType, // 1:入库, 2:出库, 3:调整
			OrderID:     orderID,
			OrderType:   4, // 4: 库存调整申请
		}
		fmt.Println("  -> 准备写入 inventory_records:")
		fmt.Printf("     Type: %d\n", record.Type)
		fmt.Printf("     Quantity: %d\n", record.Quantity)
		if err := tx.Create(&record).Error; err != nil {
			return err
		}

		// 更新 warehouse_inventories 表（实时库存）
		var warehouseInv struct {
			ID uint
		}
		err = tx.Table("warehouse_inventories").
			Select("id").
			Where("product_id = ? AND warehouse_id = ?", productID, warehouseID).
			Take(&warehouseInv).Error
		if err == gorm.ErrRecordNotFound {
			err = tx.Table("warehouse_inventories").Create(map[string]interface{}{
				"product_id":   productID,
				"warehouse_id": warehouseID,
				"quantity":     currentStock + quantity,
			}).Error
			if err != nil {
				return err
			}
		} else if err == nil {
			err = tx.Table("warehouse_inventories").
				Where("id = ?", warehouseInv.ID).
				Update("quantity", currentStock+quantity).Error
			if err != nil {
				return err
			}
		} else {
			return err
		}

		return nil
	})
}

// GetWarehouseInventory 获取指定仓库的所有产品库存
func (m *InventoryModel) GetWarehouseInventory(warehouseID uint) ([]struct {
	ProductID uint
	Quantity  int
}, error) {
	var inventories []struct {
		ProductID uint
		Quantity  int
	}
	err := m.db.Table("warehouse_inventories").
		Select("product_id, quantity").
		Where("warehouse_id = ?", warehouseID).
		Scan(&inventories).Error
	return inventories, err
}

// ListWarehouseInventory 获取仓库库存列表
func (m *InventoryModel) ListWarehouseInventory(page, pageSize int, productID, warehouseID uint) ([]WarehouseInventory, int64, error) {
	var inventories []WarehouseInventory
	var total int64

	offset := (page - 1) * pageSize

	query := m.db.Model(&WarehouseInventory{})

	if productID > 0 {
		query = query.Where("product_id = ?", productID)
	}

	if warehouseID > 0 {
		query = query.Where("warehouse_id = ?", warehouseID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Preload("Product").Preload("Warehouse").Offset(offset).Limit(pageSize).Order("updated_at DESC").Find(&inventories).Error; err != nil {
		return nil, 0, err
	}

	return inventories, total, nil
}

// GetProductTotalStock 获取产品的总库存（从 warehouse_inventories）
func (m *InventoryModel) GetProductTotalStock(productID uint) (int, error) {
	var total int64
	err := m.db.Model(&WarehouseInventory{}).
		Where("product_id = ?", productID).
		Select("COALESCE(SUM(quantity), 0)").
		Scan(&total).Error
	return int(total), err
}

// GetProductStockByWarehouse 获取产品在指定仓库的库存
func (m *InventoryModel) GetProductStockByWarehouse(productID, warehouseID uint) (int, error) {
	var inv WarehouseInventory
	err := m.db.Where("product_id = ? AND warehouse_id = ?", productID, warehouseID).First(&inv).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil
		}
		return 0, err
	}
	return inv.Quantity, nil
}

// SetStock 直接设置库存数量
func (m *InventoryModel) SetStock(productID, warehouseID uint, quantity int, remark string) error {
	return m.db.Transaction(func(tx *gorm.DB) error {
		// 检查库存数量是否为负数
		if quantity < 0 {
			return errors.New("库存数量不能为负数")
		}

		// 获取当前库存（从流水记录计算）
		var currentStock int
		err := tx.Model(&InventoryRecord{}).
			Select("COALESCE(SUM(quantity), 0)").
			Where("product_id = ? AND warehouse_id = ?", productID, warehouseID).
			Scan(&currentStock).Error
		if err != nil {
			return err
		}

		// 记录库存变动历史
		change := InventoryChange{
			ProductID:      productID,
			WarehouseID:    warehouseID,
			BeforeQuantity: currentStock,
			AfterQuantity:  quantity,
			Quantity:       quantity - currentStock,
			Type:           3, // 3: 调整
			Remark:         remark,
		}
		if err := tx.Create(&change).Error; err != nil {
			return err
		}

		// 记录库存流水
		record := InventoryRecord{
			ProductID:   productID,
			WarehouseID: warehouseID,
			Quantity:    quantity - currentStock,
			Type:        3, // 3: 调整
		}
		if err := tx.Create(&record).Error; err != nil {
			return err
		}

		return nil
	})
}

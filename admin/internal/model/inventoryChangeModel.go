package model

import (
	"time"

	"gorm.io/gorm"
)

type InventoryChange struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	ProductID      uint      `gorm:"not null;index" json:"product_id"`
	WarehouseID    uint      `gorm:"not null;index" json:"warehouse_id"`
	BeforeQuantity int       `gorm:"not null" json:"before_quantity"`
	AfterQuantity  int       `gorm:"not null" json:"after_quantity"`
	Quantity       int       `gorm:"not null" json:"quantity"`
	Type           int       `gorm:"not null;index" json:"type"`        // 1: 入库, 2: 出库, 3: 调整
	OrderID        uint      `gorm:"index" json:"order_id"`             // 关联订单ID
	OrderType      int       `gorm:"default:0;index" json:"order_type"` // 0: 无, 1: 采购订单, 2: 销售订单, 3: 库存调整申请
	Remark         string    `gorm:"size:500" json:"remark"`
	CreatedAt      time.Time `json:"created_at"`
}

type InventoryChangeModel struct {
	db *gorm.DB
}

func NewInventoryChangeModel(db *gorm.DB) *InventoryChangeModel {
	return &InventoryChangeModel{db: db}
}

func (m *InventoryChangeModel) Create(record *InventoryChange) error {
	return m.db.Create(record).Error
}

func (m *InventoryChangeModel) GetHistory(productID, warehouseID uint) ([]InventoryChange, error) {
	var records []InventoryChange
	err := m.db.Where("product_id = ? AND warehouse_id = ?", productID, warehouseID).
		Order("created_at DESC").Find(&records).Error
	return records, err
}

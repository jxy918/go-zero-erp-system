package model

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"gorm.io/gorm"
)

// InventoryTransfer 调拨单主表
type InventoryTransfer struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	TransferNo      string         `gorm:"size:50;not null;unique;index" json:"transfer_no"` // 调拨单号（格式：TF-YYYYMMDD-0001）
	FromWarehouseID uint           `gorm:"not null;index" json:"from_warehouse_id"`
	ToWarehouseID   uint           `gorm:"not null;index" json:"to_warehouse_id"`
	ProductID       uint           `gorm:"not null;index" json:"product_id"`
	Quantity        int            `gorm:"not null" json:"quantity"`
	Status          int            `gorm:"default:1;index" json:"status"` // 1:待审核, 2:已审核, 3:已完成, 4:已拒绝
	Remark          string         `gorm:"size:500" json:"remark"`
	CreatedBy       uint           `gorm:"index" json:"created_by"`
	AuditedBy       uint           `gorm:"index" json:"audited_by"`
	ExecutedBy      uint           `gorm:"index" json:"executed_by"`
	AuditedAt       *time.Time     `json:"audited_at"`
	ExecutedAt      *time.Time     `json:"executed_at"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`

	FromWarehouse Warehouse `gorm:"foreignKey:FromWarehouseID;references:ID" json:"from_warehouse"`
	ToWarehouse   Warehouse `gorm:"foreignKey:ToWarehouseID;references:ID" json:"to_warehouse"`
	Product       Product   `gorm:"foreignKey:ProductID;references:ID" json:"product"`
}

// InventoryTransferModel 调拨单数据模型
type InventoryTransferModel struct {
	db *gorm.DB
}

func NewInventoryTransferModel(db *gorm.DB) *InventoryTransferModel {
	return &InventoryTransferModel{db: db}
}

// Create 创建调拨单
func (m *InventoryTransferModel) Create(transfer *InventoryTransfer) error {
	return m.db.Create(transfer).Error
}

// GetByID 根据ID获取调拨单
func (m *InventoryTransferModel) GetByID(id uint) (*InventoryTransfer, error) {
	var transfer InventoryTransfer
	err := m.db.Preload("FromWarehouse").Preload("ToWarehouse").Preload("Product").Preload("Product.Units").First(&transfer, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &transfer, nil
}

// List 获取调拨单列表
func (m *InventoryTransferModel) List(page, pageSize int, transferNo string, fromWarehouseID, toWarehouseID, productID, status uint) ([]InventoryTransfer, int64, error) {
	var transfers []InventoryTransfer
	var total int64

	offset := (page - 1) * pageSize

	query := m.db.Model(&InventoryTransfer{}).Preload("FromWarehouse").Preload("ToWarehouse").Preload("Product").Preload("Product.Units")

	if transferNo != "" {
		query = query.Where("transfer_no LIKE ?", "%"+transferNo+"%")
	}

	if fromWarehouseID > 0 {
		query = query.Where("from_warehouse_id = ?", fromWarehouseID)
	}

	if toWarehouseID > 0 {
		query = query.Where("to_warehouse_id = ?", toWarehouseID)
	}

	if productID > 0 {
		query = query.Where("product_id = ?", productID)
	}

	if status > 0 {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&transfers).Error; err != nil {
		return nil, 0, err
	}

	return transfers, total, nil
}

// Update 更新调拨单
func (m *InventoryTransferModel) Update(transfer *InventoryTransfer) error {
	return m.db.Save(transfer).Error
}

// Delete 删除调拨单（硬删除）
func (m *InventoryTransferModel) Delete(id uint) error {
	return m.db.Unscoped().Delete(&InventoryTransfer{}, id).Error
}

// GenerateTransferNo 生成调拨单号
func (m *InventoryTransferModel) GenerateTransferNo() (string, error) {
	now := time.Now()
	dateStr := now.Format("20060102")
	timestamp := now.Format("150405")
	prefix := fmt.Sprintf("TF-%s-", dateStr)

	for i := 0; i < 10; i++ {
		random := rand.Intn(9000) + 1000
		transferNo := fmt.Sprintf("%s%s%04d", prefix, timestamp, random)

		var exists int64
		m.db.Unscoped().Model(&InventoryTransfer{}).
			Where("transfer_no = ?", transferNo).
			Count(&exists)

		if exists == 0 {
			return transferNo, nil
		}
	}

	return "", fmt.Errorf("failed to generate unique transfer no")
}

// Audit 审核调拨单
func (m *InventoryTransferModel) Audit(id uint, status int, auditedBy uint) error {
	now := time.Now()
	return m.db.Model(&InventoryTransfer{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     status,
			"audited_by": auditedBy,
			"audited_at": &now,
		}).Error
}

// Execute 执行调拨
func (m *InventoryTransferModel) Execute(id uint, operatorID uint) error {
	return m.db.Transaction(func(tx *gorm.DB) error {
		var transfer InventoryTransfer
		if err := tx.First(&transfer, id).Error; err != nil {
			return err
		}

		// 检查状态是否为已审核
		if transfer.Status != 2 {
			return errors.New("调拨单状态不正确，无法执行")
		}

		// 检查源仓库库存
		var stock int
		tx.Model(&InventoryRecord{}).
			Select("COALESCE(SUM(quantity), 0)").
			Where("product_id = ? AND warehouse_id = ?", transfer.ProductID, transfer.FromWarehouseID).
			Scan(&stock)

		if stock < transfer.Quantity {
			return errors.New("源仓库库存不足")
		}

		// 源仓库扣减库存（出库）
		inventoryModel := NewInventoryModel(tx)
		remark := fmt.Sprintf("调拨出库，调拨单号：%s", transfer.TransferNo)
		if err := inventoryModel.AdjustStock(transfer.ProductID, transfer.FromWarehouseID, -transfer.Quantity, transfer.ID, 2, remark); err != nil {
			return err
		}

		// 目标仓库增加库存（入库）
		remark = fmt.Sprintf("调拨入库，调拨单号：%s", transfer.TransferNo)
		if err := inventoryModel.AdjustStock(transfer.ProductID, transfer.ToWarehouseID, transfer.Quantity, transfer.ID, 1, remark); err != nil {
			return err
		}

		// 更新调拨单状态为已完成
		now := time.Now()
		return tx.Model(&InventoryTransfer{}).
			Where("id = ?", id).
			Updates(map[string]interface{}{
				"status":      3,
				"executed_at": &now,
				"executed_by": operatorID,
			}).Error
	})
}

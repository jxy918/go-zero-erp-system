package model

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

// InventoryCheck 盘点单主表
type InventoryCheck struct {
	ID          uint                 `gorm:"primaryKey" json:"id"`
	CheckNo     string               `gorm:"size:50;not null;unique;index" json:"check_no"` // 盘点单号（格式：CK-YYYYMMDD-0001）
	WarehouseID uint                 `gorm:"not null;index" json:"warehouse_id"`
	Status      int                  `gorm:"default:1;index" json:"status"` // 1:待盘点, 2:盘点中, 3:已完成, 4:已提交
	TotalDiff   int                  `gorm:"default:0" json:"total_diff"`   // 总差异数量
	Remark      string               `gorm:"size:500" json:"remark"`
	CreatedAt   time.Time            `json:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at"`
	DeletedAt   gorm.DeletedAt       `gorm:"index" json:"-"`
	Warehouse   Warehouse            `gorm:"foreignKey:WarehouseID;references:ID" json:"warehouse"`
	Items       []InventoryCheckItem `gorm:"foreignKey:CheckID;references:ID" json:"items"`
}

// InventoryCheckItem 盘点明细表
type InventoryCheckItem struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CheckID   uint           `gorm:"not null;index" json:"check_id"`
	ProductID uint           `gorm:"not null;index" json:"product_id"`
	SystemQty int            `gorm:"not null" json:"system_qty"`    // 系统库存数量
	ActualQty int            `gorm:"not null" json:"actual_qty"`    // 实际盘点数量
	DiffQty   int            `gorm:"not null" json:"diff_qty"`      // 差异数量（实际-系统）
	Status    int            `gorm:"default:1;index" json:"status"` // 1:待处理, 2:已处理
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Product   Product        `gorm:"foreignKey:ProductID;references:ID" json:"product"`
}

// InventoryCheckModel 盘点单数据模型
type InventoryCheckModel struct {
	db *gorm.DB
}

func NewInventoryCheckModel(db *gorm.DB) *InventoryCheckModel {
	return &InventoryCheckModel{db: db}
}

// Create 创建盘点单
func (m *InventoryCheckModel) Create(check *InventoryCheck) error {
	return m.db.Create(check).Error
}

// GetByID 根据ID获取盘点单
func (m *InventoryCheckModel) GetByID(id uint) (*InventoryCheck, error) {
	var check InventoryCheck
	err := m.db.Preload("Warehouse").Preload("Items").Preload("Items.Product").Preload("Items.Product.Units").First(&check, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &check, nil
}

// List 获取盘点单列表
func (m *InventoryCheckModel) List(page, pageSize int, checkNo string, warehouseID uint, status int) ([]InventoryCheck, int64, error) {
	var checks []InventoryCheck
	var total int64

	offset := (page - 1) * pageSize

	query := m.db.Model(&InventoryCheck{}).Preload("Warehouse").Preload("Items")

	if checkNo != "" {
		query = query.Where("check_no LIKE ?", "%"+checkNo+"%")
	}

	if warehouseID > 0 {
		query = query.Where("warehouse_id = ?", warehouseID)
	}

	if status > 0 {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&checks).Error; err != nil {
		return nil, 0, err
	}

	return checks, total, nil
}

// Update 更新盘点单
func (m *InventoryCheckModel) Update(check *InventoryCheck) error {
	return m.db.Save(check).Error
}

// Delete 删除盘点单（硬删除）
func (m *InventoryCheckModel) Delete(id uint) error {
	return m.db.Unscoped().Delete(&InventoryCheck{}, id).Error
}

// CreateItems 批量创建盘点明细
func (m *InventoryCheckModel) CreateItems(items []InventoryCheckItem) error {
	return m.db.Create(&items).Error
}

// UpdateItem 更新盘点明细
func (m *InventoryCheckModel) UpdateItem(item *InventoryCheckItem) error {
	return m.db.Save(item).Error
}

// UpdateItems 批量更新盘点明细
func (m *InventoryCheckModel) UpdateItems(items []InventoryCheckItem) error {
	for _, item := range items {
		if err := m.db.Save(&item).Error; err != nil {
			return err
		}
	}
	return nil
}

// DeleteItemsByCheckID 根据盘点单ID删除明细（硬删除）
func (m *InventoryCheckModel) DeleteItemsByCheckID(checkID uint) error {
	return m.db.Unscoped().Where("check_id = ?", checkID).Delete(&InventoryCheckItem{}).Error
}

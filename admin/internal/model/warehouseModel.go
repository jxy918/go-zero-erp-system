package model

import (
	"errors"

	"gorm.io/gorm"
)

type WarehouseModel struct {
	db *gorm.DB
}

func NewWarehouseModel(db *gorm.DB) *WarehouseModel {
	return &WarehouseModel{db: db}
}

func (m *WarehouseModel) Create(warehouse *Warehouse) error {
	return m.db.Create(warehouse).Error
}

func (m *WarehouseModel) Update(warehouse *Warehouse) error {
	return m.db.Save(warehouse).Error
}

func (m *WarehouseModel) Delete(id uint) error {
	return m.db.Unscoped().Delete(&Warehouse{}, id).Error
}

func (m *WarehouseModel) GetByID(id uint) (*Warehouse, error) {
	var warehouse Warehouse
	err := m.db.First(&warehouse, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &warehouse, nil
}

func (m *WarehouseModel) GetByName(name string) (*Warehouse, error) {
	var warehouse Warehouse
	err := m.db.Where("name = ?", name).First(&warehouse).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &warehouse, nil
}

func (m *WarehouseModel) List(page, pageSize int, name, code string) ([]Warehouse, int64, error) {
	var warehouses []Warehouse
	var total int64

	offset := (page - 1) * pageSize

	query := m.db.Model(&Warehouse{})

	if name != "" || code != "" {
		query = query.Where("name LIKE ? OR code LIKE ?", "%"+name+"%", "%"+code+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(pageSize).Find(&warehouses).Error; err != nil {
		return nil, 0, err
	}

	return warehouses, total, nil
}

func (m *WarehouseModel) ListActive(page, pageSize int, name, code string) ([]Warehouse, int64, error) {
	var warehouses []Warehouse
	var total int64

	offset := (page - 1) * pageSize

	query := m.db.Model(&Warehouse{}).Where("status = ?", 1)

	if name != "" || code != "" {
		query = query.Where("name LIKE ? OR code LIKE ?", "%"+name+"%", "%"+code+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(pageSize).Find(&warehouses).Error; err != nil {
		return nil, 0, err
	}

	return warehouses, total, nil
}

package model

import (
	"errors"

	"gorm.io/gorm"
)

type SupplierModel struct {
	db *gorm.DB
}

func NewSupplierModel(db *gorm.DB) *SupplierModel {
	return &SupplierModel{db: db}
}

func (m *SupplierModel) Create(supplier *Supplier) error {
	return m.db.Create(supplier).Error
}

func (m *SupplierModel) Update(supplier *Supplier) error {
	return m.db.Save(supplier).Error
}

func (m *SupplierModel) Delete(id uint) error {
	return m.db.Unscoped().Delete(&Supplier{}, id).Error
}

func (m *SupplierModel) GetByID(id uint) (*Supplier, error) {
	var supplier Supplier
	err := m.db.First(&supplier, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &supplier, nil
}

func (m *SupplierModel) GetByName(name string) (*Supplier, error) {
	var supplier Supplier
	err := m.db.Where("name = ?", name).First(&supplier).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &supplier, nil
}

func (m *SupplierModel) List(page, pageSize int, name, code string) ([]Supplier, int64, error) {
	var suppliers []Supplier
	var total int64

	offset := (page - 1) * pageSize

	query := m.db.Model(&Supplier{})

	if name != "" || code != "" {
		query = query.Where("name LIKE ? OR code LIKE ?", "%"+name+"%", "%"+code+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(pageSize).Find(&suppliers).Error; err != nil {
		return nil, 0, err
	}

	return suppliers, total, nil
}

func (m *SupplierModel) ListActive(page, pageSize int, name, code string) ([]Supplier, int64, error) {
	var suppliers []Supplier
	var total int64

	offset := (page - 1) * pageSize

	query := m.db.Model(&Supplier{}).Where("status = ?", 1)

	if name != "" || code != "" {
		query = query.Where("name LIKE ? OR code LIKE ?", "%"+name+"%", "%"+code+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(pageSize).Find(&suppliers).Error; err != nil {
		return nil, 0, err
	}

	return suppliers, total, nil
}

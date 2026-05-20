package model

import (
	"errors"

	"gorm.io/gorm"
)

type CustomerModel struct {
	db *gorm.DB
}

func NewCustomerModel(db *gorm.DB) *CustomerModel {
	return &CustomerModel{db: db}
}

func (m *CustomerModel) Create(customer *Customer) error {
	return m.db.Create(customer).Error
}

func (m *CustomerModel) Update(customer *Customer) error {
	return m.db.Save(customer).Error
}

func (m *CustomerModel) Delete(id uint) error {
	return m.db.Unscoped().Delete(&Customer{}, id).Error
}

func (m *CustomerModel) GetByID(id uint) (*Customer, error) {
	var customer Customer
	err := m.db.First(&customer, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &customer, nil
}

func (m *CustomerModel) GetByName(name string) (*Customer, error) {
	var customer Customer
	err := m.db.Where("name = ?", name).First(&customer).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &customer, nil
}

func (m *CustomerModel) List(page, pageSize int, name, code string) ([]Customer, int64, error) {
	var customers []Customer
	var total int64

	offset := (page - 1) * pageSize

	query := m.db.Model(&Customer{})

	if name != "" || code != "" {
		query = query.Where("name LIKE ? OR code LIKE ?", "%"+name+"%", "%"+code+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(pageSize).Find(&customers).Error; err != nil {
		return nil, 0, err
	}

	return customers, total, nil
}

func (m *CustomerModel) ListActive(page, pageSize int, name, code string) ([]Customer, int64, error) {
	var customers []Customer
	var total int64

	offset := (page - 1) * pageSize

	query := m.db.Model(&Customer{}).Where("status = ?", 1)

	if name != "" || code != "" {
		query = query.Where("name LIKE ? OR code LIKE ?", "%"+name+"%", "%"+code+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(pageSize).Find(&customers).Error; err != nil {
		return nil, 0, err
	}

	return customers, total, nil
}

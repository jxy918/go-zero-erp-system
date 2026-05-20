package model

import (
	"errors"

	"gorm.io/gorm"
)

type ProductUnitModel struct {
	db *gorm.DB
}

func NewProductUnitModel(db *gorm.DB) *ProductUnitModel {
	return &ProductUnitModel{db: db}
}

func (m *ProductUnitModel) GetMainUnit(productID uint) (*ProductUnit, error) {
	var unit ProductUnit
	err := m.db.Where("product_id = ? AND is_main = 1", productID).First(&unit).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &unit, nil
}

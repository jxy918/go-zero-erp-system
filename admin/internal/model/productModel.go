package model

import (
	"errors"

	"gorm.io/gorm"
)

type ProductModel struct {
	db *gorm.DB
}

func NewProductModel(db *gorm.DB) *ProductModel {
	return &ProductModel{db: db}
}

func (m *ProductModel) Create(product *Product) error {
	return m.db.Create(product).Error
}

func (m *ProductModel) Update(product *Product) error {
	return m.db.Model(&Product{}).Where("id = ?", product.ID).Updates(map[string]interface{}{
		"name":         product.Name,
		"code":         product.Code,
		"category_id":  product.CategoryID,
		"spec":         product.Spec,
		"price":        product.Price,
		"cost_price":   product.CostPrice,
		"min_stock":    product.MinStock,
		"max_stock":    product.MaxStock,
		"safety_stock": product.SafetyStock,
		"desc":         product.Desc,
		"status":       product.Status,
	}).Error
}

func (m *ProductModel) Delete(id uint) error {
	return m.db.Unscoped().Delete(&Product{}, id).Error
}

func (m *ProductModel) GetByID(id uint) (*Product, error) {
	var product Product
	err := m.db.Preload("Category").Preload("Units").First(&product, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &product, nil
}

func (m *ProductModel) GetByCode(code string) (*Product, error) {
	var product Product
	err := m.db.Where("code = ?", code).First(&product).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &product, nil
}

func (m *ProductModel) List(page, pageSize int, name, code string) ([]Product, int64, error) {
	var products []Product
	var total int64

	offset := (page - 1) * pageSize

	query := m.db.Model(&Product{})

	if name != "" && code != "" {
		query = query.Where("name LIKE ? OR code LIKE ?", "%"+name+"%", "%"+code+"%")
	} else if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	} else if code != "" {
		query = query.Where("code LIKE ?", "%"+code+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Preload("Category").Preload("Units").Offset(offset).Limit(pageSize).Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (m *ProductModel) ListActive(page, pageSize int, name, code string) ([]Product, int64, error) {
	var products []Product
	var total int64

	offset := (page - 1) * pageSize

	query := m.db.Model(&Product{}).Where("status = ?", 1)

	if name != "" && code != "" {
		query = query.Where("name LIKE ? OR code LIKE ?", "%"+name+"%", "%"+code+"%")
	} else if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	} else if code != "" {
		query = query.Where("code LIKE ?", "%"+code+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Preload("Category").Preload("Units").Offset(offset).Limit(pageSize).Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (m *ProductModel) UpdateStock(productID uint, quantity int) error {
	return m.db.Model(&Product{}).Where("id = ?", productID).Update("stock", gorm.Expr("stock + ?", quantity)).Error
}

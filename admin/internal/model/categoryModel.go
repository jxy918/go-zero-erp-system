package model

import (
	"errors"

	"gorm.io/gorm"
)

type CategoryModel struct {
	db *gorm.DB
}

func NewCategoryModel(db *gorm.DB) *CategoryModel {
	return &CategoryModel{db: db}
}

func (m *CategoryModel) Create(category *Category) error {
	return m.db.Create(category).Error
}

func (m *CategoryModel) Update(category *Category) error {
	return m.db.Save(category).Error
}

func (m *CategoryModel) Delete(id uint) error {
	return m.db.Unscoped().Delete(&Category{}, id).Error
}

func (m *CategoryModel) GetByID(id uint) (*Category, error) {
	var category Category
	err := m.db.First(&category, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &category, nil
}

func (m *CategoryModel) GetByName(name string) (*Category, error) {
	var category Category
	err := m.db.Where("name = ?", name).First(&category).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &category, nil
}

func (m *CategoryModel) GetByCode(code string) (*Category, error) {
	var category Category
	err := m.db.Where("code = ?", code).First(&category).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &category, nil
}

func (m *CategoryModel) List(page, pageSize int, name string) ([]Category, int64, error) {
	var categories []Category
	var total int64

	offset := (page - 1) * pageSize

	query := m.db.Model(&Category{})

	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Order("sort ASC").Offset(offset).Limit(pageSize).Find(&categories).Error; err != nil {
		return nil, 0, err
	}

	return categories, total, nil
}

func (m *CategoryModel) ListAll() ([]Category, error) {
	var categories []Category
	err := m.db.Order("sort ASC").Find(&categories).Error
	return categories, err
}

func (m *CategoryModel) BuildTree(categories []Category) []Category {
	parentMap := make(map[uint][]Category)

	for _, cat := range categories {
		parentMap[cat.ParentID] = append(parentMap[cat.ParentID], cat)
	}

	var build func(parentID uint) []Category
	build = func(parentID uint) []Category {
		children, ok := parentMap[parentID]
		if !ok {
			return nil
		}
		for i := range children {
			children[i].Children = build(children[i].ID)
		}
		return children
	}

	return build(0)
}

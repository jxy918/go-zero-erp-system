package model

import (
	"errors"

	"gorm.io/gorm"
)

// PermissionModel 权限模型
type PermissionModel struct {
	db *gorm.DB
}

// NewPermissionModel 创建权限模型实例
func NewPermissionModel(db *gorm.DB) *PermissionModel {
	return &PermissionModel{db: db}
}

// Create 创建权限
func (m *PermissionModel) Create(permission *Permission) error {
	return m.db.Create(permission).Error
}

// Update 更新权限
func (m *PermissionModel) Update(permission *Permission) error {
	return m.db.Save(permission).Error
}

// Delete 删除权限（硬删除）
func (m *PermissionModel) Delete(id uint) error {
	return m.db.Unscoped().Delete(&Permission{}, id).Error
}

// GetByID 根据ID获取权限
func (m *PermissionModel) GetByID(id uint) (*Permission, error) {
	var permission Permission
	err := m.db.First(&permission, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &permission, nil
}

// GetByCode 根据代码获取权限
func (m *PermissionModel) GetByCode(code string) (*Permission, error) {
	var permission Permission
	err := m.db.Where("code = ?", code).First(&permission).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &permission, nil
}

// List 获取权限列表
func (m *PermissionModel) List(page, pageSize int) ([]Permission, int64, error) {
	return m.ListWithName(page, pageSize, "")
}

// ListWithName 根据名称获取权限列表
func (m *PermissionModel) ListWithName(page, pageSize int, name string) ([]Permission, int64, error) {
	var permissions []Permission
	var total int64

	offset := (page - 1) * pageSize

	query := m.db.Model(&Permission{})

	if name != "" {
		query = query.Where("name LIKE ? OR code LIKE ?", "%"+name+"%", "%"+name+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query = query.Offset(offset).Limit(pageSize).Order("sort asc")

	if err := query.Find(&permissions).Error; err != nil {
		return nil, 0, err
	}

	return permissions, total, nil
}

// ListByRoleID 根据角色ID获取权限列表
func (m *PermissionModel) ListByRoleID(roleID uint) ([]Permission, error) {
	var permissions []Permission
	err := m.db.Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").Where("role_permissions.role_id = ?", roleID).Find(&permissions).Error
	if err != nil {
		return nil, err
	}
	return permissions, nil
}

// GetAll 获取所有权限
func (m *PermissionModel) GetAll() ([]Permission, error) {
	var permissions []Permission
	err := m.db.Order("sort asc").Find(&permissions).Error
	if err != nil {
		return nil, err
	}
	return permissions, nil
}

// GetAllByUserID 根据用户ID获取所有权限
func (m *PermissionModel) GetAllByUserID(userID uint) ([]Permission, error) {
	var permissions []Permission
	err := m.db.Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").Joins("JOIN user_roles ON user_roles.role_id = role_permissions.role_id").Where("user_roles.user_id = ?", userID).Distinct().Order("sort asc").Find(&permissions).Error
	if err != nil {
		return nil, err
	}
	return permissions, nil
}

// GetByPath 根据后端接口路径获取权限
func (m *PermissionModel) GetByPath(apiPath string) (*Permission, error) {
	var permission Permission
	err := m.db.Where("path = ? AND status = 1", apiPath).First(&permission).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

// ListByUserID 根据用户ID获取用户有权限的权限列表（分页，支持名称筛选）
func (m *PermissionModel) ListByUserID(userID uint, page, pageSize int, name string) ([]interface{}, int64, error) {
	var permissions []*Permission
	var total int64

	offset := (page - 1) * pageSize

	query := m.db.Model(&Permission{}).
		Distinct("permissions.id", "permissions.name", "permissions.code", "permissions.desc", "permissions.sort", "permissions.status").
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Joins("JOIN user_roles ON user_roles.role_id = role_permissions.role_id").
		Where("user_roles.user_id = ?", userID)

	if name != "" {
		query = query.Where("permissions.name LIKE ? OR permissions.code LIKE ?", "%"+name+"%", "%"+name+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query = query.Order("sort asc").Offset(offset).Limit(pageSize)

	if err := query.Find(&permissions).Error; err != nil {
		return nil, 0, err
	}

	// 转换为 []interface{}
	result := make([]interface{}, len(permissions))
	for i, p := range permissions {
		result[i] = p
	}

	return result, total, nil
}

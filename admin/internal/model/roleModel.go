package model

import (
	"errors"

	"gorm.io/gorm"
)

// RoleModel 角色模型
type RoleModel struct {
	db *gorm.DB
}

// NewRoleModel 创建角色模型实例
func NewRoleModel(db *gorm.DB) *RoleModel {
	return &RoleModel{db: db}
}

// Create 创建角色
func (m *RoleModel) Create(role *Role) error {
	return m.db.Create(role).Error
}

// Update 更新角色
func (m *RoleModel) Update(role *Role) error {
	return m.db.Save(role).Error
}

// Delete 删除角色（硬删除）
func (m *RoleModel) Delete(id uint) error {
	return m.db.Unscoped().Delete(&Role{}, id).Error
}

// GetByID 根据ID获取角色
func (m *RoleModel) GetByID(id uint) (*Role, error) {
	var role Role
	err := m.db.Preload("Permissions").Preload("Menus").First(&role, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &role, nil
}

// GetByCode 根据代码获取角色
func (m *RoleModel) GetByCode(code string) (*Role, error) {
	var role Role
	err := m.db.Where("code = ?", code).First(&role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &role, nil
}

// List 获取角色列表
func (m *RoleModel) List(page, pageSize int) ([]Role, int64, error) {
	return m.ListWithName(page, pageSize, "")
}

// ListWithName 获取角色列表（支持按名称或编码筛选）
func (m *RoleModel) ListWithName(page, pageSize int, name string) ([]Role, int64, error) {
	var roles []Role
	var total int64

	offset := (page - 1) * pageSize

	query := m.db.Model(&Role{})

	// 如果有名称参数，添加筛选条件（同时匹配名称和编码）
	if name != "" {
		query = query.Where("name LIKE ? OR code LIKE ?", "%"+name+"%", "%"+name+"%")
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	query = query.Preload("Permissions").Preload("Menus").Offset(offset).Limit(pageSize)

	if err := query.Find(&roles).Error; err != nil {
		return nil, 0, err
	}

	return roles, total, nil
}

// ListByRoleID 根据角色ID获取角色列表
func (m *RoleModel) ListByRoleID(roleID uint, page, pageSize int, name string) ([]Role, int64, error) {
	var roles []Role
	var total int64

	offset := (page - 1) * pageSize

	query := m.db.Model(&Role{}).Where("id = ?", roleID)

	if name != "" {
		query = query.Where("name LIKE ? OR code LIKE ?", "%"+name+"%", "%"+name+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query = query.Preload("Permissions").Preload("Menus").Offset(offset).Limit(pageSize)

	if err := query.Find(&roles).Error; err != nil {
		return nil, 0, err
	}

	return roles, total, nil
}

// ListByUserID 根据用户ID获取用户所属的所有角色
func (m *RoleModel) ListByUserID(userID uint, page, pageSize int, name string) ([]Role, int64, error) {
	var roles []Role
	var total int64

	offset := (page - 1) * pageSize

	query := m.db.Model(&Role{}).
		Joins("JOIN user_roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ?", userID)

	if name != "" {
		query = query.Where("roles.name LIKE ? OR roles.code LIKE ?", "%"+name+"%", "%"+name+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query = query.Preload("Permissions").Preload("Menus").Offset(offset).Limit(pageSize)

	if err := query.Find(&roles).Error; err != nil {
		return nil, 0, err
	}

	return roles, total, nil
}

// AssignPermissions 为角色分配权限（只保存按钮权限，不保存菜单）
func (m *RoleModel) AssignPermissions(roleID uint, permissionIDs []uint) error {
	// 先清除角色的所有权限（使用原始SQL直接物理删除，最可靠）
	if err := m.db.Exec("DELETE FROM role_permissions WHERE role_id = ?", roleID).Error; err != nil {
		return err
	}

	// 如果权限ID为空，直接返回
	if len(permissionIDs) == 0 {
		return nil
	}

	// 去重并收集所有权限ID
	seen := make(map[uint]bool)
	for _, permID := range permissionIDs {
		if !seen[permID] {
			seen[permID] = true
		}
	}

	// 使用批量插入提高性能（只保存按钮权限，不保存菜单）
	rolePermissions := make([]RolePermission, 0, len(seen))
	for permID := range seen {
		rolePermissions = append(rolePermissions, RolePermission{RoleID: roleID, PermissionID: permID})
	}
	if err := m.db.Create(&rolePermissions).Error; err != nil {
		return err
	}

	return nil
}

// AssignMenus 为角色分配菜单
func (m *RoleModel) AssignMenus(roleID uint, menuIDs []uint) error {
	// 先清除角色的所有菜单（使用原始SQL直接物理删除，最可靠）
	if err := m.db.Exec("DELETE FROM role_menus WHERE role_id = ?", roleID).Error; err != nil {
		return err
	}

	// 如果菜单ID为空，直接返回
	if len(menuIDs) == 0 {
		return nil
	}

	// 去重
	seen := make(map[uint]bool)
	for _, menuID := range menuIDs {
		if !seen[menuID] {
			seen[menuID] = true
		}
	}

	// 使用批量插入提高性能
	roleMenus := make([]RoleMenu, 0, len(seen))
	for menuID := range seen {
		roleMenus = append(roleMenus, RoleMenu{RoleID: roleID, MenuID: menuID})
	}
	if err := m.db.Create(&roleMenus).Error; err != nil {
		return err
	}

	return nil
}

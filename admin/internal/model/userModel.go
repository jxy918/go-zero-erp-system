package model

import (
	"errors"

	"gorm.io/gorm"
)

// UserModel 用户模型
type UserModel struct {
	db *gorm.DB
}

// NewUserModel 创建用户模型实例
func NewUserModel(db *gorm.DB) *UserModel {
	return &UserModel{db: db}
}

// Create 创建用户
func (m *UserModel) Create(user *User) error {
	return m.db.Create(user).Error
}

// Update 更新用户
func (m *UserModel) Update(user *User) error {
	return m.db.Save(user).Error
}

// Delete 删除用户（硬删除）
func (m *UserModel) Delete(id uint) error {
	return m.db.Unscoped().Delete(&User{}, id).Error
}

// GetByID 根据ID获取用户
func (m *UserModel) GetByID(id uint) (*User, error) {
	var user User
	err := m.db.Preload("Roles").First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// GetByUsername 根据用户名获取用户（包含角色、权限和菜单）
func (m *UserModel) GetByUsername(username string) (*User, error) {
	var user User
	// 预加载角色及其权限和菜单
	err := m.db.Preload("Roles").Preload("Roles.Permissions").Preload("Roles.Menus").Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

// List 获取用户列表
func (m *UserModel) List(page, pageSize int) ([]User, int64, error) {
	return m.ListWithUsername(page, pageSize, "")
}

// ListWithUsername 获取用户列表（支持按用户名或昵称筛选）
func (m *UserModel) ListWithUsername(page, pageSize int, keyword string) ([]User, int64, error) {
	var users []User
	var total int64

	offset := (page - 1) * pageSize

	query := m.db.Model(&User{})

	// 如果有关键字参数，同时按用户名和昵称搜索
	if keyword != "" {
		query = query.Where("username LIKE ? OR nickname LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	query = query.Preload("Roles").Preload("Roles.Permissions").Offset(offset).Limit(pageSize)

	if err := query.Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// ListByRoleID 根据角色ID获取用户列表
func (m *UserModel) ListByRoleID(roleID uint, page, pageSize int, username string) ([]User, int64, error) {
	var users []User
	var total int64

	offset := (page - 1) * pageSize

	query := m.db.Model(&User{}).Joins("JOIN user_roles ON user_roles.user_id = users.id").
		Where("user_roles.role_id = ?", roleID)

	if username != "" {
		query = query.Where("users.username LIKE ?", "%"+username+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query = query.Preload("Roles").Preload("Roles.Permissions").Offset(offset).Limit(pageSize)

	if err := query.Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// AssignRoles 为用户分配角色
func (m *UserModel) AssignRoles(userID uint, roleIDs []uint) error {
	// 先清除用户的所有角色（使用 Unscoped 强制物理删除，避免软删除导致主键冲突）
	if err := m.db.Unscoped().Where("user_id = ?", userID).Delete(&UserRole{}).Error; err != nil {
		return err
	}

	// 如果角色ID为空，直接返回
	if len(roleIDs) == 0 {
		return nil
	}

	// 去重
	seen := make(map[uint]bool)
	for _, roleID := range roleIDs {
		if !seen[roleID] {
			seen[roleID] = true
		}
	}

	// 使用批量插入提高性能
	userRoles := make([]UserRole, 0, len(seen))
	for roleID := range seen {
		userRoles = append(userRoles, UserRole{UserID: userID, RoleID: roleID})
	}
	if err := m.db.Create(&userRoles).Error; err != nil {
		return err
	}

	return nil
}

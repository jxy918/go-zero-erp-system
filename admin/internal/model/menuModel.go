package model

import (
	"gorm.io/gorm"
)

// MenuModel 菜单模型接口
type MenuModel interface {
	Create(menu *Menu) error
	Update(menu *Menu) error
	Delete(id uint) error
	Get(id uint) (*Menu, error)
	GetAll() ([]Menu, error)
	GetTree() ([]Menu, error)
	GetUserMenus(userID uint, isAdmin bool) ([]Menu, error)
	AssignPermissions(menuID uint, permissionIDs []uint) error
	List(page, pageSize int) ([]Menu, int64, error)
	ListByUserID(userID uint, page, pageSize int) ([]interface{}, int64, error)
}

// menuModel 菜单模型实现
type menuModel struct {
	db *gorm.DB
}

// NewMenuModel 创建菜单模型实例
func NewMenuModel(db *gorm.DB) MenuModel {
	return &menuModel{db: db}
}

// Create 创建菜单
func (m *menuModel) Create(menu *Menu) error {
	return m.db.Create(menu).Error
}

// Update 更新菜单
func (m *menuModel) Update(menu *Menu) error {
	return m.db.Save(menu).Error
}

// Delete 删除菜单（硬删除）
func (m *menuModel) Delete(id uint) error {
	return m.db.Unscoped().Delete(&Menu{}, id).Error
}

// Get 根据ID获取菜单
func (m *menuModel) Get(id uint) (*Menu, error) {
	var menu Menu
	err := m.db.Preload("Permissions").First(&menu, id).Error
	if err != nil {
		return nil, err
	}
	return &menu, nil
}

// GetAll 获取所有菜单（管理员使用，包含禁用状态）
func (m *menuModel) GetAll() ([]Menu, error) {
	var menus []Menu
	err := m.db.Order("sort asc").Find(&menus).Error
	return menus, err
}

// GetTree 获取菜单树
func (m *menuModel) GetTree() ([]Menu, error) {
	var menus []Menu
	err := m.db.Order("sort asc").Find(&menus).Error
	if err != nil {
		return nil, err
	}

	// 构建菜单树
	menuMap := make(map[uint]*Menu)
	var rootMenus []Menu

	// 首先将所有菜单放入map
	for i := range menus {
		menuMap[menus[i].ID] = &menus[i]
	}

	// 构建树结构
	for i := range menus {
		if menus[i].ParentID == 0 {
			// 根菜单
			rootMenus = append(rootMenus, menus[i])
		} else {
			// 子菜单，添加到父菜单的子列表
			if parent, ok := menuMap[menus[i].ParentID]; ok {
				// 这里需要注意，由于Go的切片引用特性，我们需要使用指针来修改原始数据
				parent.Children = append(parent.Children, menus[i])
			}
		}
	}

	return rootMenus, nil
}

// GetUserMenus 根据用户ID获取用户的菜单
func (m *menuModel) GetUserMenus(userID uint, isAdmin bool) ([]Menu, error) {
	// 如果是admin用户，返回所有菜单
	if isAdmin {
		return m.GetAll()
	}

	// 根据用户的角色获取菜单（使用DISTINCT去重，避免多个角色分配相同菜单导致重复）
	var menus []Menu
	err := m.db.Distinct("menus.id", "menus.name", "menus.code", "menus.desc", "menus.parent_id", "menus.path", "menus.component", "menus.icon", "menus.sort", "menus.status", "menus.created_at", "menus.updated_at", "menus.deleted_at").
		Table("menus").
		Joins("JOIN role_menus ON role_menus.menu_id = menus.id").
		Joins("JOIN user_roles ON user_roles.role_id = role_menus.role_id").
		Where("user_roles.user_id = ? AND menus.status = 1", userID).
		Order("menus.sort asc").
		Find(&menus).Error
	return menus, err
}

// AssignPermissions 为菜单分配权限
func (m *menuModel) AssignPermissions(menuID uint, permissionIDs []uint) error {
	// 先清除现有的权限关联（使用 Unscoped 强制物理删除，避免软删除导致主键冲突）
	if err := m.db.Unscoped().Where("menu_id = ?", menuID).Delete(&MenuPermission{}).Error; err != nil {
		return err
	}

	// 如果有权限ID，则添加新的关联
	if len(permissionIDs) > 0 {
		// 去重
		seen := make(map[uint]bool)
		for _, permID := range permissionIDs {
			if !seen[permID] {
				seen[permID] = true
			}
		}

		// 使用批量插入提高性能
		menuPermissions := make([]MenuPermission, 0, len(seen))
		for permID := range seen {
			menuPermissions = append(menuPermissions, MenuPermission{MenuID: menuID, PermissionID: permID})
		}
		if err := m.db.Create(&menuPermissions).Error; err != nil {
			return err
		}
	}

	return nil
}

// List 获取菜单列表（分页）
func (m *menuModel) List(page, pageSize int) ([]Menu, int64, error) {
	var menus []Menu
	var total int64

	// 计算总数
	if err := m.db.Model(&Menu{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := m.db.Order("sort asc").Offset(offset).Limit(pageSize).Find(&menus).Error; err != nil {
		return nil, 0, err
	}

	return menus, total, nil
}

// ListByUserID 根据用户ID获取用户有权限的菜单列表（分页）
func (m *menuModel) ListByUserID(userID uint, page, pageSize int) ([]interface{}, int64, error) {
	var menus []*Menu
	var total int64

	offset := (page - 1) * pageSize

	// 查询用户有权限的菜单
	query := m.db.Model(&Menu{}).
		Distinct("menus.id", "menus.name", "menus.code", "menus.desc", "menus.parent_id", "menus.path", "menus.component", "menus.icon", "menus.sort", "menus.status").
		Joins("JOIN role_menus ON role_menus.menu_id = menus.id").
		Joins("JOIN user_roles ON user_roles.role_id = role_menus.role_id").
		Where("user_roles.user_id = ? AND menus.status = 1", userID)

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	if err := query.Order("sort asc").Offset(offset).Limit(pageSize).Find(&menus).Error; err != nil {
		return nil, 0, err
	}

	// 转换为 []interface{}
	result := make([]interface{}, len(menus))
	for i, menu := range menus {
		result[i] = menu
	}

	return result, total, nil
}

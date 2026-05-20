-- 库存盘点功能菜单和权限 SQL
-- 执行前请确保数据库中有 warehouses 和 products 表

-- =============================================
-- 1. 添加菜单
-- =============================================

-- 查找库存管理父菜单ID
SET @parent_id = (SELECT id FROM menus WHERE code = 'menu_inventory' LIMIT 1);

-- 如果没有库存管理菜单，则创建
IF @parent_id IS NULL THEN
    INSERT INTO menus (name, code, parent_id, path, component, icon, sort, type, status, created_at, updated_at)
    VALUES ('库存管理', 'menu_inventory', 0, '/inventory', 'Layout', 'Box', 5, 1, 1, NOW(), NOW());
    SET @parent_id = LAST_INSERT_ID();
END IF;

-- 创建库存盘点子菜单
INSERT INTO menus (name, code, parent_id, path, component, icon, sort, type, status, created_at, updated_at)
VALUES ('库存盘点', 'menu_inventory_check', @parent_id, '/inventory/check', 'inventory/InventoryCheck', 'DocumentChecked', 1, 1, 1, NOW(), NOW());

-- 获取盘点菜单ID
SET @check_menu_id = LAST_INSERT_ID();

-- =============================================
-- 2. 添加权限
-- =============================================

-- 盘点菜单查看权限
INSERT INTO permissions (name, code, type, menu_id, status, created_at, updated_at)
VALUES ('查看盘点列表', 'inventory_check:list', 2, @check_menu_id, 1, NOW(), NOW());

-- 盘点创建权限
INSERT INTO permissions (name, code, type, menu_id, status, created_at, updated_at)
VALUES ('创建盘点单', 'inventory_check:create', 2, @check_menu_id, 1, NOW(), NOW());

-- 盘点更新权限
INSERT INTO permissions (name, code, type, menu_id, status, created_at, updated_at)
VALUES ('更新盘点单', 'inventory_check:update', 2, @check_menu_id, 1, NOW(), NOW());

-- 盘点删除权限
INSERT INTO permissions (name, code, type, menu_id, status, created_at, updated_at)
VALUES ('删除盘点单', 'inventory_check:delete', 2, @check_menu_id, 1, NOW(), NOW());

-- 盘点提交权限
INSERT INTO permissions (name, code, type, menu_id, status, created_at, updated_at)
VALUES ('提交盘点单', 'inventory_check:submit', 2, @check_menu_id, 1, NOW(), NOW());

-- 盘点生成调整权限
INSERT INTO permissions (name, code, type, menu_id, status, created_at, updated_at)
VALUES ('生成调整申请', 'inventory_check:generate', 2, @check_menu_id, 1, NOW(), NOW());

-- =============================================
-- 3. 分配权限给超级管理员角色
-- =============================================

-- 查找超级管理员角色ID
SET @admin_role_id = (SELECT id FROM roles WHERE code = 'admin' LIMIT 1);

-- 如果没有admin角色，查找第一个角色
IF @admin_role_id IS NULL THEN
    SET @admin_role_id = (SELECT id FROM roles LIMIT 1);
END IF;

-- 分配所有盘点权限给管理员
INSERT INTO role_permissions (role_id, permission_id)
SELECT @admin_role_id, id FROM permissions WHERE code LIKE 'inventory_check:%';

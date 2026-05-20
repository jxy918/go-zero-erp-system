SET NAMES utf8mb4;

-- 获取查看报表菜单ID
SELECT @menu_id := id FROM menus WHERE code = 'menu_erp_report';

-- 获取查看ERP报表权限ID
SELECT @permission_id := id FROM permissions WHERE code = 'btn_erp_view';

-- 获取管理员角色ID
SELECT @admin_role_id := id FROM roles WHERE code = 'admin';

-- 插入菜单权限关联（如果不存在）
INSERT IGNORE INTO menu_permissions (menu_id, permission_id, created_at, updated_at)
VALUES (@menu_id, @permission_id, NOW(), NOW());

-- 插入角色菜单关联（如果不存在）
INSERT IGNORE INTO role_menus (role_id, menu_id, created_at, updated_at)
VALUES (@admin_role_id, @menu_id, NOW(), NOW());

-- 验证插入结果
SELECT '菜单权限关联' as type, * FROM menu_permissions WHERE menu_id = @menu_id;
SELECT '角色菜单关联' as type, * FROM role_menus WHERE menu_id = @menu_id;

SELECT '初始化完成' as result;
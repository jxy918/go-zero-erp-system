
-- 1. 修复菜单权限关联：把按钮权限关联到对应的菜单上
-- 菜单ID和按钮权限的对应关系：
-- 菜单ID: 10 (menu_product) -> 权限: 36, 37, 38
-- 菜单ID: 11 (menu_purchase) -> 权限: 39, 40, 41, 42, 75
-- 菜单ID: 12 (menu_sales) -> 权限: 43, 44, 45, 46, 76
-- 菜单ID: 13 (menu_inventory) -> 权限: 47
-- 菜单ID: 18 (menu_category) -> 权限: 63, 64, 65
-- 菜单ID: 19 (menu_supplier) -> 权限: 66, 67, 68
-- 菜单ID: 20 (menu_customer) -> 权限: 69, 70, 71
-- 菜单ID: 21 (menu_warehouse) -> 权限: 72, 73, 74
-- 菜单ID: 9 (menu_erp) -> 权限: 48, 77

INSERT IGNORE INTO menu_permissions (menu_id, permission_id) VALUES
(10, 36), (10, 37), (10, 38),
(11, 39), (11, 40), (11, 41), (11, 42), (11, 75),
(12, 43), (12, 44), (12, 45), (12, 46), (12, 76),
(13, 47),
(18, 63), (18, 64), (18, 65),
(19, 66), (19, 67), (19, 68),
(20, 69), (20, 70), (20, 71),
(21, 72), (21, 73), (21, 74),
(9, 48), (9, 77);

-- 2. 为管理员角色分配所有按钮权限
INSERT IGNORE INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.code = 'admin'
AND p.type = 2;

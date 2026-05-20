INSERT IGNORE INTO permissions (name, code, type, `desc`) VALUES 
('取消采购订单', 'btn_purchase_cancel', 2, '取消采购订单按钮'),
('取消销售订单', 'btn_sales_cancel', 2, '取消销售订单按钮'),
('查看订单日志', 'btn_order_log_view', 2, '查看订单日志按钮');

INSERT IGNORE INTO menu_permissions (menu_id, permission_id)
SELECT m.id, p.id FROM menus m, permissions p 
WHERE m.code = 'menu_purchase' AND p.code IN ('btn_purchase_cancel');

INSERT IGNORE INTO menu_permissions (menu_id, permission_id)
SELECT m.id, p.id FROM menus m, permissions p 
WHERE m.code = 'menu_sales' AND p.code IN ('btn_sales_cancel');

INSERT IGNORE INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r, permissions p 
WHERE r.code = 'admin' AND p.code IN ('btn_purchase_cancel', 'btn_sales_cancel', 'btn_order_log_view');

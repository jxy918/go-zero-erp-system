-- =============================================
-- ERP 进销存系统测试数据脚本
-- 版本: v1.0
-- 日期: 2026-04-28
-- =============================================

-- 注意：密码使用 bcrypt 加密，明文为 "admin123"
SET @ADMIN_PWD = '$2a$10$N9qo8uLOickgx2ZMRZoMye.IjzqAKL9xL5jvMFVdNJHvGCgTq/VEq';

-- =============================================
-- 1. 用户数据
-- =============================================
INSERT INTO users (username, password, nickname, email, phone, status, created_at, updated_at) VALUES
('admin', @ADMIN_PWD, '系统管理员', 'admin@example.com', '13800138000', 1, NOW(), NOW()),
('purchase', @ADMIN_PWD, '采购专员', 'purchase@example.com', '13800138001', 1, NOW(), NOW()),
('sales', @ADMIN_PWD, '销售专员', 'sales@example.com', '13800138002', 1, NOW(), NOW()),
('warehouse', @ADMIN_PWD, '仓库管理员', 'warehouse@example.com', '13800138003', 1, NOW(), NOW()),
('disabled', @ADMIN_PWD, '禁用用户', 'disabled@example.com', '13800138004', 0, NOW(), NOW());

-- =============================================
-- 2. 角色数据
-- =============================================
INSERT INTO roles (name, code, desc, status, created_at, updated_at) VALUES
('系统管理员', 'admin', '拥有系统所有权限', 1, NOW(), NOW()),
('采购管理员', 'purchase_admin', '负责采购订单管理', 1, NOW(), NOW()),
('销售管理员', 'sales_admin', '负责销售订单管理', 1, NOW(), NOW()),
('仓库管理员', 'warehouse_admin', '负责库存管理', 1, NOW(), NOW());

-- =============================================
-- 3. 用户角色关联
-- =============================================
INSERT INTO user_roles (user_id, role_id) VALUES
(1, 1), -- admin -> 系统管理员
(2, 2), -- purchase -> 采购管理员
(3, 3), -- sales -> 销售管理员
(4, 4); -- warehouse -> 仓库管理员

-- =============================================
-- 4. 权限数据
-- =============================================
INSERT INTO permissions (name, code, desc, type, parent_id, path, component, icon, sort, status, created_at, updated_at) VALUES
-- 系统管理菜单
('系统管理', 'system', '系统管理模块', 1, 0, '/system', 'System', 'Setting', 1, 1, NOW(), NOW()),
('用户管理', 'user', '用户管理', 1, 1, '/user', 'User', 'User', 1, 1, NOW(), NOW()),
('角色管理', 'role', '角色管理', 1, 1, '/role', 'Role', 'UserGroup', 2, 1, NOW(), NOW()),
('权限管理', 'permission', '权限管理', 1, 1, '/permission', 'Permission', 'Key', 3, 1, NOW(), NOW()),
('菜单管理', 'menu', '菜单管理', 1, 1, '/menu', 'Menu', 'Menu', 4, 1, NOW(), NOW()),
-- 用户管理按钮权限
('查看用户列表', 'user:list', '查看用户列表', 2, 2, '', '', '', 1, 1, NOW(), NOW()),
('创建用户', 'user:create', '创建用户', 2, 2, '', '', '', 2, 1, NOW(), NOW()),
('更新用户', 'user:update', '更新用户', 2, 2, '', '', '', 3, 1, NOW(), NOW()),
('删除用户', 'user:delete', '删除用户', 2, 2, '', '', '', 4, 1, NOW(), NOW()),
('分配角色', 'user:assign_roles', '分配角色', 2, 2, '', '', '', 5, 1, NOW(), NOW()),

-- 采购管理菜单
('采购管理', 'purchase', '采购管理模块', 1, 0, '/purchase', 'Purchase', 'ShoppingCart', 2, 1, NOW(), NOW()),
('采购订单', 'purchase_order', '采购订单管理', 1, 11, '/purchase/order', 'PurchaseOrder', 'FileText', 1, 1, NOW(), NOW()),
-- 采购订单按钮权限
('查看采购订单', 'purchase_order:list', '查看采购订单', 2, 12, '', '', '', 1, 1, NOW(), NOW()),
('创建采购订单', 'purchase_order:create', '创建采购订单', 2, 12, '', '', '', 2, 1, NOW(), NOW()),
('审核采购订单', 'purchase_order:audit', '审核采购订单', 2, 12, '', '', '', 3, 1, NOW(), NOW()),
('采购入库', 'purchase_order:inbound', '采购入库', 2, 12, '', '', '', 4, 1, NOW(), NOW()),
('取消采购订单', 'purchase_order:cancel', '取消采购订单', 2, 12, '', '', '', 5, 1, NOW(), NOW()),
('删除采购订单', 'purchase_order:delete', '删除采购订单', 2, 12, '', '', '', 6, 1, NOW(), NOW()),

-- 销售管理菜单
('销售管理', 'sales', '销售管理模块', 1, 0, '/sales', 'Sales', 'TrendingUp', 3, 1, NOW(), NOW()),
('销售订单', 'sales_order', '销售订单管理', 1, 18, '/sales/order', 'SalesOrder', 'FileText', 1, 1, NOW(), NOW()),
-- 销售订单按钮权限
('查看销售订单', 'sales_order:list', '查看销售订单', 2, 19, '', '', '', 1, 1, NOW(), NOW()),
('创建销售订单', 'sales_order:create', '创建销售订单', 2, 19, '', '', '', 2, 1, NOW(), NOW()),
('审核销售订单', 'sales_order:audit', '审核销售订单', 2, 19, '', '', '', 3, 1, NOW(), NOW()),
('销售出库', 'sales_order:outbound', '销售出库', 2, 19, '', '', '', 4, 1, NOW(), NOW()),
('取消销售订单', 'sales_order:cancel', '取消销售订单', 2, 19, '', '', '', 5, 1, NOW(), NOW()),
('删除销售订单', 'sales_order:delete', '删除销售订单', 2, 19, '', '', '', 6, 1, NOW(), NOW()),

-- 库存管理菜单
('库存管理', 'inventory', '库存管理模块', 1, 0, '/inventory', 'Inventory', 'Package', 4, 1, NOW(), NOW()),
-- 库存管理按钮权限
('查看库存', 'inventory:list', '查看库存', 2, 26, '', '', '', 1, 1, NOW(), NOW()),
('库存调整', 'inventory:adjust', '库存调整', 2, 26, '', '', '', 2, 1, NOW(), NOW()),
('库存变动记录', 'inventory:history', '库存变动记录', 2, 26, '', '', '', 3, 1, NOW(), NOW()),

-- 产品管理菜单
('产品管理', 'product', '产品管理模块', 1, 0, '/product', 'Product', 'Box', 5, 1, NOW(), NOW()),
('产品列表', 'product_list', '产品列表', 1, 30, '/product/list', 'Product', 'Box', 1, 1, NOW(), NOW()),
('产品分类', 'category', '产品分类', 1, 30, '/product/category', 'Category', 'FolderOpen', 2, 1, NOW(), NOW()),
-- 产品管理按钮权限
('查看产品', 'product:list', '查看产品', 2, 31, '', '', '', 1, 1, NOW(), NOW()),
('创建产品', 'product:create', '创建产品', 2, 31, '', '', '', 2, 1, NOW(), NOW()),
('更新产品', 'product:update', '更新产品', 2, 31, '', '', '', 3, 1, NOW(), NOW()),
('删除产品', 'product:delete', '删除产品', 2, 31, '', '', '', 4, 1, NOW(), NOW()),

-- 供应商管理菜单
('供应商管理', 'supplier', '供应商管理模块', 1, 0, '/supplier', 'Supplier', 'Building', 6, 1, NOW(), NOW()),
-- 供应商管理按钮权限
('查看供应商', 'supplier:list', '查看供应商', 2, 38, '', '', '', 1, 1, NOW(), NOW()),
('创建供应商', 'supplier:create', '创建供应商', 2, 38, '', '', '', 2, 1, NOW(), NOW()),
('更新供应商', 'supplier:update', '更新供应商', 2, 38, '', '', '', 3, 1, NOW(), NOW()),
('删除供应商', 'supplier:delete', '删除供应商', 2, 38, '', '', '', 4, 1, NOW(), NOW()),

-- 客户管理菜单
('客户管理', 'customer', '客户管理模块', 1, 0, '/customer', 'Customer', 'Users', 7, 1, NOW(), NOW()),
-- 客户管理按钮权限
('查看客户', 'customer:list', '查看客户', 2, 44, '', '', '', 1, 1, NOW(), NOW()),
('创建客户', 'customer:create', '创建客户', 2, 44, '', '', '', 2, 1, NOW(), NOW()),
('更新客户', 'customer:update', '更新客户', 2, 44, '', '', '', 3, 1, NOW(), NOW()),
('删除客户', 'customer:delete', '删除客户', 2, 44, '', '', '', 4, 1, NOW(), NOW()),

-- 活动日志菜单
('活动日志', 'activity', '操作日志', 1, 0, '/activity', 'Activity', 'Clock', 8, 1, NOW(), NOW()),
('查看活动日志', 'activity:list', '查看活动日志', 2, 50, '', '', '', 1, 1, NOW(), NOW());

-- =============================================
-- 5. 角色权限关联
-- =============================================
-- 系统管理员拥有所有权限
INSERT INTO role_permissions (role_id, permission_id)
SELECT 1, id FROM permissions;

-- 采购管理员权限
INSERT INTO role_permissions (role_id, permission_id) VALUES
(2, 11), (2, 12), (2, 13), (2, 14), (2, 15), (2, 16), (2, 17), -- 采购管理相关
(2, 30), (2, 31), (2, 32), (2, 33), -- 产品管理相关
(2, 38), (2, 39), (2, 40), (2, 41), (2, 42), -- 供应商管理相关
(2, 26), (2, 27); -- 库存查看权限

-- 销售管理员权限
INSERT INTO role_permissions (role_id, permission_id) VALUES
(3, 18), (3, 19), (3, 20), (3, 21), (3, 22), (3, 23), (3, 24), (3, 25), -- 销售管理相关
(3, 30), (3, 31), (3, 32), -- 产品查看权限
(3, 44), (3, 45), (3, 46), (3, 47), (3, 48), -- 客户管理相关
(3, 26), (3, 27); -- 库存查看权限

-- 仓库管理员权限
INSERT INTO role_permissions (role_id, permission_id) VALUES
(4, 26), (4, 27), (4, 28), (4, 29), -- 库存管理相关
(4, 30), (3, 31), (3, 32), -- 产品查看权限
(4, 11), (4, 12), (4, 13), (4, 16), -- 采购订单入库权限
(4, 18), (4, 19), (4, 20), (4, 23); -- 销售订单出库权限

-- =============================================
-- 6. 产品分类数据
-- =============================================
INSERT INTO categories (name, code, parent_id, sort, status, desc, created_at, updated_at) VALUES
('电子产品', 'electronics', 0, 1, 1, '电子类产品', NOW(), NOW()),
('食品饮料', 'food', 0, 2, 1, '食品饮料类产品', NOW(), NOW()),
('日用百货', 'daily', 0, 3, 1, '日用百货类产品', NOW(), NOW()),
('办公用品', 'office', 0, 4, 1, '办公用品类产品', NOW(), NOW()),
('手机', 'phone', 1, 1, 1, '手机产品', NOW(), NOW()),
('电脑', 'computer', 1, 2, 1, '电脑产品', NOW(), NOW()),
('零食', 'snack', 2, 1, 1, '零食类', NOW(), NOW()),
('饮料', 'beverage', 2, 2, 1, '饮料类', NOW(), NOW());

-- =============================================
-- 7. 产品数据
-- =============================================
INSERT INTO products (name, code, category_id, spec, unit, price, cost_price, stock, desc, status, created_at, updated_at) VALUES
('iPhone 15 Pro', 'P001', 5, '256GB', '台', 8999.00, 7500.00, 100, '苹果手机', 1, NOW(), NOW()),
('MacBook Pro 14', 'P002', 6, 'M3 Pro/18GB/512GB', '台', 16999.00, 14000.00, 50, '苹果笔记本', 1, NOW(), NOW()),
('华为 Mate60 Pro', 'P003', 5, '512GB', '台', 6999.00, 5800.00, 80, '华为手机', 1, NOW(), NOW()),
('联想 ThinkPad', 'P004', 6, 'i7/16GB/512GB', '台', 8999.00, 7200.00, 60, '联想笔记本', 1, NOW(), NOW()),
('乐事薯片', 'P005', 7, '原味/104g', '袋', 8.50, 5.00, 500, '乐事原味薯片', 1, NOW(), NOW()),
('奥利奥饼干', 'P006', 7, '夹心/388g', '盒', 18.90, 12.00, 300, '奥利奥夹心饼干', 1, NOW(), NOW()),
('可口可乐', 'P007', 8, '330ml*24', '箱', 58.00, 45.00, 200, '可口可乐整箱', 1, NOW(), NOW()),
('农夫山泉', 'P008', 8, '550ml*24', '箱', 28.00, 20.00, 300, '农夫山泉矿泉水', 1, NOW(), NOW()),
('得力笔记本', 'P009', 4, 'A5/80页', '本', 5.50, 3.50, 1000, '得力办公笔记本', 1, NOW(), NOW()),
('晨光中性笔', 'P010', 4, '0.5mm/12支', '盒', 12.00, 7.00, 800, '晨光中性笔', 1, NOW(), NOW()),
('纸巾', 'P011', 3, '抽纸/3层/100抽', '包', 6.50, 4.00, 600, '家用抽纸', 1, NOW(), NOW()),
('洗衣液', 'P012', 3, '2kg', '瓶', 28.00, 18.00, 200, '洗衣液', 1, NOW(), NOW());

-- =============================================
-- 8. 供应商数据
-- =============================================
INSERT INTO suppliers (name, code, contact, phone, desc, email, address, status, created_at, updated_at) VALUES
('苹果公司', 'S001', '张三', '13800138101', '苹果产品供应商', 'apple@example.com', '北京市朝阳区', 1, NOW(), NOW()),
('华为技术', 'S002', '李四', '13800138102', '华为产品供应商', 'huawei@example.com', '深圳市南山区', 1, NOW(), NOW()),
('百事食品', 'S003', '王五', '13800138103', '休闲食品供应商', 'pepsi@example.com', '上海市浦东新区', 1, NOW(), NOW()),
('可口可乐中国', 'S004', '赵六', '13800138104', '饮料供应商', 'coke@example.com', '广州市天河区', 1, NOW(), NOW()),
('得力集团', 'S005', '孙七', '13800138105', '办公用品供应商', 'deli@example.com', '宁波市江北区', 1, NOW(), NOW()),
('维达纸业', 'S006', '周八', '13800138106', '日用品供应商', 'vinda@example.com', '佛山市顺德区', 1, NOW(), NOW());

-- =============================================
-- 9. 客户数据
-- =============================================
INSERT INTO customers (name, code, contact, phone, desc, email, address, status, created_at, updated_at) VALUES
('京东商城', 'C001', '刘强', '13800138201', '电商平台大客户', 'jd@example.com', '北京市大兴区', 1, NOW(), NOW()),
('天猫超市', 'C002', '马云', '13800138202', '电商平台大客户', 'tmall@example.com', '杭州市余杭区', 1, NOW(), NOW()),
('沃尔玛', 'C003', '山姆', '13800138203', '超市连锁客户', 'walmart@example.com', '深圳市福田区', 1, NOW(), NOW()),
('华润万家', 'C004', '陈明', '13800138204', '超市连锁客户', 'crv@example.com', '广州市越秀区', 1, NOW(), NOW()),
('苏宁易购', 'C005', '张近', '13800138205', '电商平台客户', 'suning@example.com', '南京市鼓楼区', 1, NOW(), NOW()),
('步步高', 'C006', '王填', '13800138206', '超市连锁客户', 'bbg@example.com', '长沙市芙蓉区', 1, NOW(), NOW());

-- =============================================
-- 10. 仓库数据
-- =============================================
INSERT INTO warehouses (name, code, address, contact, phone, desc, status, created_at, updated_at) VALUES
('北京仓库', 'W001', '北京市朝阳区仓库路1号', '李经理', '13800138301', '北京地区主仓库', 1, NOW(), NOW()),
('上海仓库', 'W002', '上海市浦东新区仓库路2号', '王经理', '13800138302', '上海地区主仓库', 1, NOW(), NOW()),
('广州仓库', 'W003', '广州市天河区仓库路3号', '张经理', '13800138303', '广州地区主仓库', 1, NOW(), NOW()),
('深圳仓库', 'W004', '深圳市南山区仓库路4号', '刘经理', '13800138304', '深圳地区主仓库', 1, NOW(), NOW());

-- =============================================
-- 11. 库存初始数据
-- =============================================
INSERT INTO inventory_records (product_id, warehouse_id, quantity, type, created_at, updated_at) VALUES
(1, 1, 30, 1, NOW(), NOW()), -- iPhone 15 Pro 北京仓 30台
(1, 2, 40, 1, NOW(), NOW()), -- iPhone 15 Pro 上海仓 40台
(1, 3, 30, 1, NOW(), NOW()), -- iPhone 15 Pro 广州仓 30台
(2, 1, 15, 1, NOW(), NOW()), -- MacBook Pro 北京仓 15台
(2, 2, 20, 1, NOW(), NOW()), -- MacBook Pro 上海仓 20台
(3, 1, 25, 1, NOW(), NOW()), -- 华为 Mate60 北京仓 25台
(3, 4, 35, 1, NOW(), NOW()), -- 华为 Mate60 深圳仓 35台
(5, 1, 150, 1, NOW(), NOW()), -- 乐事薯片 北京仓 150袋
(5, 2, 200, 1, NOW(), NOW()), -- 乐事薯片 上海仓 200袋
(7, 3, 80, 1, NOW(), NOW()), -- 可口可乐 广州仓 80箱
(7, 4, 70, 1, NOW(), NOW()), -- 可口可乐 深圳仓 70箱
(10, 1, 200, 1, NOW(), NOW()), -- 晨光中性笔 北京仓 200盒
(11, 2, 150, 1, NOW(), NOW()), -- 纸巾 上海仓 150包
(12, 3, 60, 1, NOW(), NOW()); -- 洗衣液 广州仓 60瓶

-- =============================================
-- 12. 测试采购订单数据
-- =============================================
-- 待审核采购订单
INSERT INTO purchase_orders (order_no, supplier_id, warehouse_id, total_amount, status, remark, created_at, updated_at) VALUES
('PO202604280001', 1, 1, 89990.00, 1, '采购iPhone 15 Pro 10台', NOW(), NOW()),
('PO202604280002', 3, 2, 945.00, 1, '采购奥利奥饼干50盒', NOW(), NOW());

INSERT INTO purchase_order_items (order_id, product_id, quantity, unit_price, amount, created_at, updated_at) VALUES
(1, 1, 10, 8999.00, 89990.00, NOW(), NOW()),
(2, 6, 50, 18.90, 945.00, NOW(), NOW());

-- 已审核采购订单（可入库）
INSERT INTO purchase_orders (order_no, supplier_id, warehouse_id, total_amount, status, remark, created_at, updated_at) VALUES
('PO202604270001', 2, 4, 279960.00, 2, '采购华为Mate60 Pro 40台', NOW(), NOW());

INSERT INTO purchase_order_items (order_id, product_id, quantity, unit_price, amount, created_at, updated_at) VALUES
(3, 3, 40, 6999.00, 279960.00, NOW(), NOW());

-- 已完成采购订单
INSERT INTO purchase_orders (order_no, supplier_id, warehouse_id, total_amount, status, remark, created_at, updated_at) VALUES
('PO202604260001', 5, 1, 5500.00, 3, '采购得力笔记本1000本', NOW(), NOW());

INSERT INTO purchase_order_items (order_id, product_id, quantity, unit_price, amount, created_at, updated_at) VALUES
(4, 9, 1000, 5.50, 5500.00, NOW(), NOW());

-- =============================================
-- 13. 测试销售订单数据
-- =============================================
-- 待审核销售订单
INSERT INTO sales_orders (order_no, customer_id, warehouse_id, total_amount, status, remark, created_at, updated_at) VALUES
('SO202604280001', 1, 1, 44995.00, 1, '京东订单-5台iPhone', NOW(), NOW()),
('SO202604280002', 3, 2, 4250.00, 1, '沃尔玛订单-500袋薯片', NOW(), NOW());

INSERT INTO sales_order_items (order_id, product_id, quantity, unit_price, amount, created_at, updated_at) VALUES
(1, 1, 5, 8999.00, 44995.00, NOW(), NOW()),
(2, 5, 500, 8.50, 4250.00, NOW(), NOW());

-- 已审核销售订单（可出库）
INSERT INTO sales_orders (order_no, customer_id, warehouse_id, total_amount, status, remark, created_at, updated_at) VALUES
('SO202604270001', 2, 3, 23400.00, 2, '天猫超市订单-400箱可乐', NOW(), NOW());

INSERT INTO sales_order_items (order_id, product_id, quantity, unit_price, amount, created_at, updated_at) VALUES
(3, 7, 400, 58.50, 23400.00, NOW(), NOW());

-- 已完成销售订单
INSERT INTO sales_orders (order_no, customer_id, warehouse_id, total_amount, status, remark, created_at, updated_at) VALUES
('SO202604260001', 4, 1, 1200.00, 3, '华润万家订单-100盒中性笔', NOW(), NOW());

INSERT INTO sales_order_items (order_id, product_id, quantity, unit_price, amount, created_at, updated_at) VALUES
(4, 10, 100, 12.00, 1200.00, NOW(), NOW());

-- =============================================
-- 14. 库存变动记录
-- =============================================
INSERT INTO inventory_changes (product_id, warehouse_id, before_quantity, after_quantity, quantity, type, order_id, remark, created_at) VALUES
(1, 1, 25, 30, 5, 1, 4, '采购入库', NOW()),
(5, 2, 180, 200, 20, 1, 4, '采购入库', NOW()),
(10, 1, 210, 200, -10, 2, 4, '销售出库', NOW()),
(11, 3, 50, 55, 5, 3, 0, '库存调整-报损补库', NOW());

-- =============================================
-- 15. 菜单数据
-- =============================================
INSERT INTO menus (name, code, desc, parent_id, path, component, icon, sort, status, created_at, updated_at) VALUES
('系统管理', 'menu_system', '系统管理模块', 0, '/system', 'System', 'Setting', 1, 1, NOW(), NOW()),
('用户管理', 'menu_user', '用户管理', 1, '/user', 'User', 'User', 1, 1, NOW(), NOW()),
('角色管理', 'menu_role', '角色管理', 1, '/role', 'Role', 'UserGroup', 2, 1, NOW(), NOW()),
('权限管理', 'menu_permission', '权限管理', 1, '/permission', 'Permission', 'Key', 3, 1, NOW(), NOW()),
('菜单管理', 'menu_menu', '菜单管理', 1, '/menu', 'Menu', 'Menu', 4, 1, NOW(), NOW()),
('采购管理', 'menu_purchase', '采购管理模块', 0, '/purchase', 'Purchase', 'ShoppingCart', 2, 1, NOW(), NOW()),
('采购订单', 'menu_purchase_order', '采购订单管理', 6, '/purchase/order', 'PurchaseOrder', 'FileText', 1, 1, NOW(), NOW()),
('销售管理', 'menu_sales', '销售管理模块', 0, '/sales', 'Sales', 'TrendingUp', 3, 1, NOW(), NOW()),
('销售订单', 'menu_sales_order', '销售订单管理', 8, '/sales/order', 'SalesOrder', 'FileText', 1, 1, NOW(), NOW()),
('库存管理', 'menu_inventory', '库存管理模块', 0, '/inventory', 'Inventory', 'Package', 4, 1, NOW(), NOW()),
('产品管理', 'menu_product', '产品管理模块', 0, '/product', 'Product', 'Box', 5, 1, NOW(), NOW()),
('产品列表', 'menu_product_list', '产品列表', 11, '/product/list', 'Product', 'Box', 1, 1, NOW(), NOW()),
('产品分类', 'menu_category', '产品分类', 11, '/product/category', 'Category', 'FolderOpen', 2, 1, NOW(), NOW()),
('供应商管理', 'menu_supplier', '供应商管理模块', 0, '/supplier', 'Supplier', 'Building', 6, 1, NOW(), NOW()),
('客户管理', 'menu_customer', '客户管理模块', 0, '/customer', 'Customer', 'Users', 7, 1, NOW(), NOW()),
('活动日志', 'menu_activity', '操作日志', 0, '/activity', 'Activity', 'Clock', 8, 1, NOW(), NOW());

-- =============================================
-- 16. 角色菜单关联
-- =============================================
-- 系统管理员拥有所有菜单
INSERT INTO role_menus (role_id, menu_id)
SELECT 1, id FROM menus;

-- 采购管理员菜单
INSERT INTO role_menus (role_id, menu_id) VALUES
(2, 6), (2, 7), (2, 11), (2, 12), (2, 13), (2, 14), (2, 10);

-- 销售管理员菜单
INSERT INTO role_menus (role_id, menu_id) VALUES
(3, 8), (3, 9), (3, 11), (3, 12), (3, 15), (3, 10);

-- 仓库管理员菜单
INSERT INTO role_menus (role_id, menu_id) VALUES
(4, 10), (4, 11), (4, 12), (4, 6), (4, 7), (4, 8), (4, 9);

-- =============================================
-- 17. 活动日志数据
-- =============================================
INSERT INTO activities (user_id, username, action, module, ip, desc, created_at) VALUES
(1, 'admin', 'login', 'system', '127.0.0.1', '用户登录系统', NOW()),
(1, 'admin', 'create', 'product', '127.0.0.1', '创建产品: iPhone 15 Pro', NOW()),
(2, 'purchase', 'create', 'purchase_order', '127.0.0.1', '创建采购订单: PO202604280001', NOW()),
(3, 'sales', 'create', 'sales_order', '127.0.0.1', '创建销售订单: SO202604280001', NOW()),
(4, 'warehouse', 'update', 'inventory', '127.0.0.1', '库存调整: 纸巾 +5', NOW());

-- =============================================
-- 数据插入完成
-- =============================================
SELECT '测试数据插入完成' AS result;
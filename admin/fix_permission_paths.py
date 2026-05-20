import re

# 从 routes.go 提取所有API路径
routes_content = """
/product/active-list
/product/category/create
/product/category/delete
/product/category/list
/product/category/tree
/product/category/update
/product/create
/product/delete
/product/get/:id
/product/list
/product/unit/create
/product/unit/delete
/product/unit/get/:id
/product/unit/list
/product/unit/update
/product/update
/warehouse/active-list
/warehouse/create
/warehouse/delete
/warehouse/get/:id
/warehouse/list
/warehouse/update
/supplier/active-list
/supplier/create
/supplier/delete
/supplier/get/:id
/supplier/list
/supplier/update
/customer/active-list
/customer/create
/customer/delete
/customer/get/:id
/customer/list
/customer/update
/purchase/order/create
/purchase/order/delete
/purchase/order/get/:id
/purchase/order/inbound
/purchase/order/list
/purchase/order/status
/purchase/order/update
/sales/order/create
/sales/order/delete
/sales/order/get/:id
/sales/order/list
/sales/order/status
/sales/order/update
/inventory/adjust/create
/inventory/adjust/approve
/inventory/adjust/list
/inventory/adjust/reject
/inventory/alert/check
/inventory/alert/list
/inventory/check/create
/inventory/check/delete
/inventory/check/get/:id
/inventory/check/list
/inventory/check/submit
/inventory/check/update
/inventory/generate-adjust-from-check
/inventory/get/:id
/inventory/history
/inventory/list
/inventory/order-log/list
/inventory/sales-outbound
/inventory/transfer/audit
/inventory/transfer/create
/inventory/transfer/delete
/inventory/transfer/execute
/inventory/transfer/get/:id
/inventory/transfer/list
/inventory/transfer/update
/user/assign-roles
/user/create
/user/delete
/user/get/:id
/user/list
/user/update
/role/assign-menus
/role/assign-permissions
/role/create
/role/delete
/role/get/:id
/role/list
/role/update
/permission/create
/permission/delete
/permission/get/:id
/permission/list
/permission/update
/menu/assign-permissions
/menu/create
/menu/delete
/menu/get/:id
/menu/list
/menu/tree
/menu/update
/activity/list
/erp/statistics/business
/erp/statistics/dashboard
/erp/statistics/inventory-alert
/erp/statistics/order-status
/erp/statistics/top-products
/erp/statistics/trend
/inventory/adjust/get/:id
/inventory/check/status
/inventory/transfer/status
/purchase/order/inbound/list
/auth/login
/auth/logout
/auth/refresh
/erp/init-data
"""

# 解析路由路径集合
route_paths = set()
for line in routes_content.strip().split('\n'):
    path = line.strip()
    if path:
        # 移除动态参数如 :id
        clean_path = re.sub(r'/:id', '/{id}', path)
        route_paths.add(path)
        route_paths.add(clean_path)

# 权限配置路径（从 initData.go 提取）
permission_paths = {
    # 用户管理
    "/user/list": "user:list",
    "/user/create": "btn_user_create",
    "/user/update": "btn_user_update",
    "/user/delete": "btn_user_delete",
    "/user/assign-roles": "btn_user_assign",
    # 角色管理
    "/role/list": "role:list",
    "/role/create": "btn_role_create",
    "/role/update": "btn_role_update",
    "/role/delete": "btn_role_delete",
    "/role/assign-permissions": "btn_role_assign",
    "/role/assign-menus": "btn_role_assign_menus",
    # 权限管理
    "/permission/list": "permission:list",
    "/permission/create": "btn_permission_create",
    "/permission/update": "btn_permission_update",
    "/permission/delete": "btn_permission_delete",
    # 菜单管理
    "/menu/list": "menu:list",
    "/menu/tree": "menu:tree",
    "/menu/create": "btn_menu_create",
    "/menu/update": "btn_menu_update",
    "/menu/delete": "btn_menu_delete",
    "/menu/assign-permissions": "btn_menu_assign_permissions",
    # 活动日志
    "/activity/list": "activity:list",
    # 产品管理
    "/product/list": "product:list",
    "/product/create": "btn_product_create",
    "/product/update": "btn_product_update",
    "/product/delete": "btn_product_delete",
    "/product/active-list": "product:active_list",
    # 产品分类
    "/product/category/list": "category:list",
    "/product/category/create": "btn_category_create",
    "/product/category/update": "btn_category_update",
    "/product/category/delete": "btn_category_delete",
    # 供应商
    "/supplier/list": "supplier:list",
    "/supplier/create": "btn_supplier_create",
    "/supplier/update": "btn_supplier_update",
    "/supplier/delete": "btn_supplier_delete",
    "/supplier/active-list": "supplier:active_list",
    # 客户
    "/customer/list": "customer:list",
    "/customer/create": "btn_customer_create",
    "/customer/update": "btn_customer_update",
    "/customer/delete": "btn_customer_delete",
    "/customer/active-list": "customer:active_list",
    # 仓库
    "/warehouse/list": "warehouse:list",
    "/warehouse/create": "btn_warehouse_create",
    "/warehouse/update": "btn_warehouse_update",
    "/warehouse/delete": "btn_warehouse_delete",
    "/warehouse/active-list": "warehouse:active_list",
    # 采购订单
    "/purchase/order/list": "purchase:list",
    "/purchase/order/create": "btn_purchase_create",
    "/purchase/order/status": "btn_purchase_approve",
    "/purchase/order/inbound": "btn_purchase_inbound",
    "/purchase/order/status": "btn_purchase_cancel",
    "/purchase/order/delete": "btn_purchase_delete",
    # 销售订单
    "/sales/order/list": "sales:list",
    "/sales/order/create": "btn_sales_create",
    "/sales/order/status": "btn_sales_approve",
    "/inventory/sales-outbound": "btn_sales_outbound",
    "/sales/order/status": "btn_sales_cancel",
    "/sales/order/delete": "btn_sales_delete",
    # 订单日志
    "/inventory/order-log/list": "order-log:list",
    # 库存管理
    "/inventory/list": "inventory:list",
    "/inventory/history": "inventory:history",
    # 库存调整
    "/inventory/adjust/list": "inventory_adjust:list",
    "/inventory/adjust/create": "inventory_adjust:create",
    "/inventory/adjust/approve": "inventory_adjust:approve",
    "/inventory/adjust/reject": "inventory_adjust:reject",
    # 库存调拨
    "/inventory/transfer/list": "inventory_transfer:list",
    "/inventory/transfer/create": "inventory_transfer:create",
    "/inventory/transfer/update": "inventory_transfer:update",
    "/inventory/transfer/delete": "inventory_transfer:delete",
    "/inventory/transfer/audit": "inventory_transfer:audit",
    "/inventory/transfer/execute": "inventory_transfer:execute",
    # 库存盘点
    "/inventory/check/list": "inventory_check:list",
    "/inventory/check/create": "inventory_check:create",
    "/inventory/check/update": "inventory_check:update",
    "/inventory/check/delete": "inventory_check:delete",
    "/inventory/check/submit": "inventory_check:submit",
    "/inventory/generate-adjust-from-check": "inventory_check:generate",
    # 库存预警
    "/inventory/alert/list": "inventory_alert:list",
    "/inventory/alert/check": "inventory_alert:check",
    # 产品计量单位
    "/product/unit/list": "product_unit:list",
    "/product/unit/create": "product_unit:create",
    "/product/unit/update": "product_unit:update",
    "/product/unit/delete": "product_unit:delete",
    # ERP统计报表
    "/erp/statistics/dashboard": "erp:view",
    "/erp/statistics/trend": "erp:trend",
    "/erp/statistics/inventory-alert": "erp:alert",
    "/erp/statistics/top-products": "erp:top",
    "/erp/statistics/order-status": "erp:status",
    "/erp/statistics/business": "erp:business",
}

print("=== 权限路径与路由路径对比结果 ===\n")
print("【不匹配的权限路径】")
mismatched = []
for perm_path, perm_code in permission_paths.items():
    if perm_path not in route_paths:
        # 检查是否是动态路径
        dynamic_path = re.sub(r'/{id}', '/:id', perm_path)
        if dynamic_path not in route_paths:
            mismatched.append((perm_path, perm_code))
            print(f"  ❌ {perm_path} ({perm_code})")

print(f"\n【匹配的权限路径】: {len(permission_paths) - len(mismatched)} 个")

if mismatched:
    print(f"\n【需要修复的路径】共 {len(mismatched)} 个:")
    for path, code in mismatched:
        print(f"  - {code}: {path}")
else:
    print("\n✅ 所有权限路径都与路由路径匹配！")

# 生成SQL修复语句
if mismatched:
    print("\n【生成的SQL修复语句】")
    for path, code in mismatched:
        print(f"UPDATE permissions SET path = '{path}' WHERE code = '{code}';")

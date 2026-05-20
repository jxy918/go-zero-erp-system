// 权限码常量定义
// 统一管理所有权限码，避免散落在各个组件中
// 格式: {模块}:{操作}

export const PermissionCodes = {
  // 认证模块
  AUTH_LOGIN: 'auth:login',
  AUTH_LOGOUT: 'auth:logout',
  AUTH_REFRESH: 'auth:refresh',

  // 用户管理
  USER_LIST: 'user:list',
  USER_CREATE: 'user:create',
  USER_UPDATE: 'user:update',
  USER_DELETE: 'user:delete',
  USER_ASSIGN_ROLES: 'user:assign_roles',

  // 角色管理
  ROLE_LIST: 'role:list',
  ROLE_CREATE: 'role:create',
  ROLE_UPDATE: 'role:update',
  ROLE_DELETE: 'role:delete',
  ROLE_ASSIGN_PERMISSIONS: 'role:assign_permissions',
  ROLE_ASSIGN_MENUS: 'role:assign_menus',

  // 权限管理
  PERMISSION_LIST: 'permission:list',
  PERMISSION_CREATE: 'permission:create',
  PERMISSION_UPDATE: 'permission:update',
  PERMISSION_DELETE: 'permission:delete',

  // 菜单管理
  MENU_LIST: 'menu:list',
  MENU_CREATE: 'menu:create',
  MENU_UPDATE: 'menu:update',
  MENU_DELETE: 'menu:delete',
  MENU_ASSIGN_PERMISSIONS: 'menu:assign_permissions',

  // 产品管理
  PRODUCT_LIST: 'product:list',
  PRODUCT_CREATE: 'product:create',
  PRODUCT_UPDATE: 'product:update',
  PRODUCT_DELETE: 'product:delete',

  // 产品分类
  CATEGORY_LIST: 'category:list',
  CATEGORY_CREATE: 'category:create',
  CATEGORY_UPDATE: 'category:update',
  CATEGORY_DELETE: 'category:delete',

  // 产品单位
  PRODUCT_UNIT_LIST: 'product_unit:list',
  PRODUCT_UNIT_CREATE: 'product_unit:create',
  PRODUCT_UNIT_UPDATE: 'product_unit:update',
  PRODUCT_UNIT_DELETE: 'product_unit:delete',

  // 供应商管理
  SUPPLIER_LIST: 'supplier:list',
  SUPPLIER_CREATE: 'supplier:create',
  SUPPLIER_UPDATE: 'supplier:update',
  SUPPLIER_DELETE: 'supplier:delete',

  // 客户管理
  CUSTOMER_LIST: 'customer:list',
  CUSTOMER_CREATE: 'customer:create',
  CUSTOMER_UPDATE: 'customer:update',
  CUSTOMER_DELETE: 'customer:delete',

  // 仓库管理
  WAREHOUSE_LIST: 'warehouse:list',
  WAREHOUSE_CREATE: 'warehouse:create',
  WAREHOUSE_UPDATE: 'warehouse:update',
  WAREHOUSE_DELETE: 'warehouse:delete',

  // 采购订单
  PURCHASE_ORDER_LIST: 'purchase_order:list',
  PURCHASE_ORDER_CREATE: 'purchase_order:create',
  PURCHASE_ORDER_UPDATE: 'purchase_order:update',
  PURCHASE_ORDER_DELETE: 'purchase_order:delete',
  PURCHASE_ORDER_STATUS: 'purchase_order:status',

  // 销售订单
  SALES_ORDER_LIST: 'sales_order:list',
  SALES_ORDER_CREATE: 'sales_order:create',
  SALES_ORDER_UPDATE: 'sales_order:update',
  SALES_ORDER_DELETE: 'sales_order:delete',
  SALES_ORDER_STATUS: 'sales_order:status',

  // 库存管理
  INVENTORY_RECORD_LIST: 'inventory_record:list',
  INVENTORY_ADJUST_CREATE: 'inventory_adjust:create',
  INVENTORY_ADJUST_LIST: 'inventory_adjust:list',
  INVENTORY_CHECK_CREATE: 'inventory_check:create',
  INVENTORY_CHECK_LIST: 'inventory_check:list',
  INVENTORY_CHECK_STATUS: 'inventory_check:status',
  INVENTORY_TRANSFER_CREATE: 'inventory_transfer:create',
  INVENTORY_TRANSFER_LIST: 'inventory_transfer:list',
  INVENTORY_TRANSFER_STATUS: 'inventory_transfer:status',
  WAREHOUSE_INVENTORY_LIST: 'warehouse_inventory:list',

  // ERP报表
  ERP_DASHBOARD: 'erp:dashboard',
  ERP_TREND: 'erp:trend',
  ERP_ORDER_STATUS: 'erp:order_status',
  ERP_INVENTORY_ALERT: 'erp:inventory_alert',
  ERP_TODO_DATA: 'erp:todo_data',

  // 活动日志
  ACTIVITY_LIST: 'activity:list'
}

// 检查用户是否有指定权限
export const hasPermission = (permissionCode, userPermissions) => {
  if (!permissionCode) {
    return false
  }
  if (!userPermissions || !Array.isArray(userPermissions)) {
    return false
  }
  const codes = userPermissions.map(p => p.code || p.Code)
  return codes.includes(permissionCode)
}

// 检查用户是否有多个权限中的任意一个
export const hasAnyPermission = (permissionCodes, userPermissions) => {
  if (!permissionCodes || !Array.isArray(permissionCodes)) {
    return false
  }
  return permissionCodes.some(code => hasPermission(code, userPermissions))
}

// 检查用户是否有所有指定权限
export const hasAllPermissions = (permissionCodes, userPermissions) => {
  if (!permissionCodes || !Array.isArray(permissionCodes)) {
    return false
  }
  return permissionCodes.every(code => hasPermission(code, userPermissions))
}
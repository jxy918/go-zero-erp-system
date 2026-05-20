package model

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

// InitData 初始化默认数据（手动调用接口触发）
func InitData() error {
	log.Println("🔴 InitData() 被调用")
	return doInitData()
}

// doInitData 执行实际的初始化逻辑
func doInitData() error {
	log.Println("🔴 doInitData() 被调用")
	// 1. 快速检查是否已初始化（通过检查管理员用户）
	var userCount int64
	DB.Model(&User{}).Count(&userCount)

	if userCount > 0 {
		// 快速路径：数据已存在，只做必要的后台维护
		log.Println("📦 检测到数据已初始化，执行后台维护...")

		// 快速修复权限路径
		if err := quickFixPermissionPaths(); err != nil {
			log.Printf("⚠️ 权限路径修复: %v", err)
		}

		// 快速修复权限code（只检查不存在的）
		if err := quickFixPermissionCodes(); err != nil {
			log.Printf("⚠️ 权限修复: %v", err)
		}

		// 快速同步菜单权限关联
		if err := quickSyncMenuPermissions(); err != nil {
			log.Printf("⚠️ 菜单权限同步: %v", err)
		}

		// 确保管理员角色有所有菜单
		if err := ensureAdminMenus(); err != nil {
			log.Printf("⚠️ 管理员菜单同步: %v", err)
		}

		// 确保管理员角色有所有权限
		if err := ensureAdminPermissions(); err != nil {
			log.Printf("⚠️ 管理员权限同步: %v", err)
		}

		// 确保缺失的菜单被创建
		if err := ensureMissingMenus(); err != nil {
			log.Printf("⚠️ 缺失菜单创建: %v", err)
		}

		// 确保测试角色存在
		if err := ensureTestRole(); err != nil {
			log.Printf("⚠️ 测试角色: %v", err)
		}

		return nil
	}

	// 完整初始化路径：第一次启动
	log.Println("🔧 首次启动，执行完整初始化...")

	// 1. 创建根菜单
	// 控制台菜单
	dashboardMenu := Menu{
		Name:      "控制台",
		Code:      "menu_dashboard",
		Desc:      "控制台模块",
		ParentID:  0,
		Path:      "/",
		Component: "Dashboard",
		Icon:      "House",
		Sort:      1,
		Status:    1,
	}
	if err := DB.Create(&dashboardMenu).Error; err != nil {
		log.Printf("创建控制台菜单失败: %v", err)
	} else {
		log.Printf("创建控制台菜单: %s (id: %d)", dashboardMenu.Name, dashboardMenu.ID)
	}

	// 系统管理菜单
	systemMenu := Menu{
		Name:      "系统管理",
		Code:      "menu_system",
		Desc:      "系统管理模块",
		ParentID:  0,
		Path:      "/system",
		Component: "",
		Icon:      "Setting",
		Sort:      2,
		Status:    1,
	}
	if err := DB.Create(&systemMenu).Error; err != nil {
		log.Printf("创建系统管理菜单失败: %v", err)
	} else {
		log.Printf("创建系统管理菜单: %s (id: %d)", systemMenu.Name, systemMenu.ID)
	}

	// ERP管理菜单
	erpMenu := Menu{
		Name:      "ERP管理",
		Code:      "menu_erp",
		Desc:      "ERP管理模块",
		ParentID:  0,
		Path:      "/erp",
		Component: "",
		Icon:      "ShoppingCart",
		Sort:      3,
		Status:    1,
	}
	if err := DB.Create(&erpMenu).Error; err != nil {
		log.Printf("创建ERP管理菜单失败: %v", err)
	} else {
		log.Printf("创建ERP管理菜单: %s (id: %d)", erpMenu.Name, erpMenu.ID)
	}

	// 2. 创建子菜单
	subMenus := []Menu{
		{Name: "用户管理", Code: "menu_user", Desc: "用户管理模块", ParentID: systemMenu.ID, Path: "/user", Component: "User", Icon: "User", Sort: 1, Status: 1},
		{Name: "角色管理", Code: "menu_role", Desc: "角色管理模块", ParentID: systemMenu.ID, Path: "/role", Component: "Role", Icon: "Avatar", Sort: 2, Status: 1},
		{Name: "权限管理", Code: "menu_permission", Desc: "权限管理模块", ParentID: systemMenu.ID, Path: "/permission", Component: "Permission", Icon: "Lock", Sort: 3, Status: 1},
		{Name: "菜单管理", Code: "menu_menu", Desc: "菜单管理模块", ParentID: systemMenu.ID, Path: "/menu", Component: "Menu", Icon: "Menu", Sort: 4, Status: 1},
		{Name: "活动日志", Code: "menu_activity", Desc: "活动日志模块", ParentID: systemMenu.ID, Path: "/activity", Component: "Activity", Icon: "Document", Sort: 5, Status: 1},
		{Name: "产品管理", Code: "menu_product", Desc: "产品管理模块", ParentID: erpMenu.ID, Path: "/product", Component: "Product", Icon: "Package", Sort: 1, Status: 1},
		{Name: "产品分类", Code: "menu_category", Desc: "产品分类模块", ParentID: erpMenu.ID, Path: "/product/category", Component: "product/Category", Icon: "Folder", Sort: 2, Status: 1},
		{Name: "供应商管理", Code: "menu_supplier", Desc: "供应商管理模块", ParentID: erpMenu.ID, Path: "/supplier", Component: "Supplier", Icon: "Shop", Sort: 3, Status: 1},
		{Name: "客户管理", Code: "menu_customer", Desc: "客户管理模块", ParentID: erpMenu.ID, Path: "/customer", Component: "Customer", Icon: "UserFilled", Sort: 4, Status: 1},
		{Name: "仓库管理", Code: "menu_warehouse", Desc: "仓库管理模块", ParentID: erpMenu.ID, Path: "/warehouse", Component: "Warehouse", Icon: "OfficeBuilding", Sort: 5, Status: 1},
		{Name: "采购管理", Code: "menu_purchase", Desc: "采购管理模块", ParentID: erpMenu.ID, Path: "/purchase", Component: "Purchase", Icon: "Shopping", Sort: 6, Status: 1},
		{Name: "销售管理", Code: "menu_sales", Desc: "销售管理模块", ParentID: erpMenu.ID, Path: "/sales", Component: "Sales", Icon: "TrendingUp", Sort: 7, Status: 1},
		{Name: "库存管理", Code: "menu_inventory", Desc: "库存管理模块", ParentID: erpMenu.ID, Path: "/inventory", Component: "Inventory", Icon: "Warehouse", Sort: 8, Status: 1},
		{Name: "库存调整", Code: "menu_inventory_adjust", Desc: "库存调整模块", ParentID: erpMenu.ID, Path: "/inventory/adjust-request", Component: "inventory/InventoryAdjustRequest", Icon: "RefreshCw", Sort: 9, Status: 1},
		{Name: "库存调拨", Code: "menu_inventory_transfer", Desc: "库存调拨模块", ParentID: erpMenu.ID, Path: "/inventory/transfer", Component: "inventory/InventoryTransfer", Icon: "ArrowRightCircle", Sort: 10, Status: 1},
		{Name: "库存盘点", Code: "menu_inventory_check", Desc: "库存盘点模块", ParentID: erpMenu.ID, Path: "/inventory/check", Component: "inventory/InventoryCheck", Icon: "DocumentChecked", Sort: 11, Status: 1},
		{Name: "库存预警", Code: "menu_inventory_alert", Desc: "库存预警模块", ParentID: erpMenu.ID, Path: "/inventory/alert", Component: "inventory/InventoryAlert", Icon: "Bell", Sort: 12, Status: 1},
		{Name: "计量单位", Code: "menu_product_unit", Desc: "计量单位模块", ParentID: erpMenu.ID, Path: "/product/unit", Component: "product/ProductUnit", Icon: "Calculator", Sort: 13, Status: 1},
		{Name: "预览报表", Code: "menu_erp_report", Desc: "ERP统计报表", ParentID: erpMenu.ID, Path: "/erp", Component: "Erp", Icon: "Eye", Sort: 14, Status: 1},
	}

	for i := range subMenus {
		if err := DB.Create(&subMenus[i]).Error; err != nil {
			log.Printf("创建子菜单失败: %s, err: %v", subMenus[i].Name, err)
		} else {
			log.Printf("创建子菜单: %s (parent_id: %d)", subMenus[i].Name, systemMenu.ID)
		}
	}

	log.Println("菜单数据初始化完成")

	// 4. 创建默认权限（按钮权限，与菜单关联）
	// 定义权限与菜单的关联关系
	type permItem struct {
		Name     string
		Code     string
		Type     int
		Desc     string
		Path     string
		MenuCode string // 关联的菜单code
	}

	permissions := []permItem{
		// 用户管理权限
		{Name: "查看用户列表", Code: "user:list", Type: 2, Desc: "查看用户列表", Path: "/user/list", MenuCode: "menu_user"},
		{Name: "创建用户", Code: "btn_user_create", Type: 2, Desc: "创建用户按钮", Path: "/user/create", MenuCode: "menu_user"},
		{Name: "编辑用户", Code: "btn_user_update", Type: 2, Desc: "编辑用户按钮", Path: "/user/update", MenuCode: "menu_user"},
		{Name: "删除用户", Code: "btn_user_delete", Type: 2, Desc: "删除用户按钮", Path: "/user/delete", MenuCode: "menu_user"},
		{Name: "分配角色", Code: "btn_user_assign", Type: 2, Desc: "分配角色按钮", Path: "/user/assign-roles", MenuCode: "menu_user"},
		// 角色管理权限
		{Name: "查看角色列表", Code: "role:list", Type: 2, Desc: "查看角色列表", Path: "/role/list", MenuCode: "menu_role"},
		{Name: "创建角色", Code: "btn_role_create", Type: 2, Desc: "创建角色按钮", Path: "/role/create", MenuCode: "menu_role"},
		{Name: "编辑角色", Code: "btn_role_update", Type: 2, Desc: "编辑角色按钮", Path: "/role/update", MenuCode: "menu_role"},
		{Name: "删除角色", Code: "btn_role_delete", Type: 2, Desc: "删除角色按钮", Path: "/role/delete", MenuCode: "menu_role"},
		{Name: "分配权限", Code: "btn_role_assign", Type: 2, Desc: "分配权限按钮", Path: "/role/assign-permissions", MenuCode: "menu_role"},
		{Name: "分配菜单", Code: "btn_role_assign_menus", Type: 2, Desc: "分配菜单按钮", Path: "/role/assign-menus", MenuCode: "menu_role"},
		// 权限管理权限
		{Name: "查看权限列表", Code: "permission:list", Type: 2, Desc: "查看权限列表", Path: "/permission/list", MenuCode: "menu_permission"},
		{Name: "创建权限", Code: "btn_permission_create", Type: 2, Desc: "创建权限按钮", Path: "/permission/create", MenuCode: "menu_permission"},
		{Name: "编辑权限", Code: "btn_permission_update", Type: 2, Desc: "编辑权限按钮", Path: "/permission/update", MenuCode: "menu_permission"},
		{Name: "删除权限", Code: "btn_permission_delete", Type: 2, Desc: "删除权限按钮", Path: "/permission/delete", MenuCode: "menu_permission"},
		// 菜单管理权限
		{Name: "查看菜单", Code: "menu:list", Type: 2, Desc: "查看菜单", Path: "/menu/list", MenuCode: "menu_menu"},
		{Name: "查看菜单树", Code: "menu:tree", Type: 2, Desc: "查看菜单树", Path: "/menu/tree", MenuCode: "menu_menu"},
		{Name: "创建菜单", Code: "btn_menu_create", Type: 2, Desc: "创建菜单按钮", Path: "/menu/create", MenuCode: "menu_menu"},
		{Name: "编辑菜单", Code: "btn_menu_update", Type: 2, Desc: "编辑菜单按钮", Path: "/menu/update", MenuCode: "menu_menu"},
		{Name: "删除菜单", Code: "btn_menu_delete", Type: 2, Desc: "删除菜单按钮", Path: "/menu/delete", MenuCode: "menu_menu"},
		{Name: "分配菜单权限", Code: "btn_menu_assign_permissions", Type: 2, Desc: "分配菜单权限", Path: "/menu/assign-permissions", MenuCode: "menu_menu"},
		// 活动日志权限
		{Name: "查看活动日志", Code: "activity:list", Type: 2, Desc: "查看活动日志", Path: "/activity/list", MenuCode: "menu_activity"},
		// ERP产品管理权限
		{Name: "查看产品列表", Code: "product:list", Type: 2, Desc: "查看产品列表", Path: "/product/list", MenuCode: "menu_product"},
		{Name: "创建产品", Code: "btn_product_create", Type: 2, Desc: "创建产品按钮", Path: "/product/create", MenuCode: "menu_product"},
		{Name: "编辑产品", Code: "btn_product_update", Type: 2, Desc: "编辑产品按钮", Path: "/product/update", MenuCode: "menu_product"},
		{Name: "删除产品", Code: "btn_product_delete", Type: 2, Desc: "删除产品按钮", Path: "/product/delete", MenuCode: "menu_product"},
		{Name: "查看激活产品", Code: "product:active_list", Type: 2, Desc: "查看激活产品列表", Path: "/product/active-list", MenuCode: "menu_product"},
		// ERP产品分类权限
		{Name: "查看分类列表", Code: "category:list", Type: 2, Desc: "查看产品分类列表", Path: "/product/category/list", MenuCode: "menu_category"},
		{Name: "创建分类", Code: "btn_category_create", Type: 2, Desc: "创建产品分类", Path: "/product/category/create", MenuCode: "menu_category"},
		{Name: "编辑分类", Code: "btn_category_update", Type: 2, Desc: "编辑产品分类", Path: "/product/category/update", MenuCode: "menu_category"},
		{Name: "删除分类", Code: "btn_category_delete", Type: 2, Desc: "删除产品分类", Path: "/product/category/delete", MenuCode: "menu_category"},
		// ERP供应商权限
		{Name: "查看供应商列表", Code: "supplier:list", Type: 2, Desc: "查看供应商列表", Path: "/supplier/list", MenuCode: "menu_supplier"},
		{Name: "创建供应商", Code: "btn_supplier_create", Type: 2, Desc: "创建供应商", Path: "/supplier/create", MenuCode: "menu_supplier"},
		{Name: "编辑供应商", Code: "btn_supplier_update", Type: 2, Desc: "编辑供应商", Path: "/supplier/update", MenuCode: "menu_supplier"},
		{Name: "删除供应商", Code: "btn_supplier_delete", Type: 2, Desc: "删除供应商", Path: "/supplier/delete", MenuCode: "menu_supplier"},
		{Name: "查看激活供应商", Code: "supplier:active_list", Type: 2, Desc: "查看激活供应商列表", Path: "/supplier/active-list", MenuCode: "menu_supplier"},
		// ERP客户权限
		{Name: "查看客户列表", Code: "customer:list", Type: 2, Desc: "查看客户列表", Path: "/customer/list", MenuCode: "menu_customer"},
		{Name: "创建客户", Code: "btn_customer_create", Type: 2, Desc: "创建客户", Path: "/customer/create", MenuCode: "menu_customer"},
		{Name: "编辑客户", Code: "btn_customer_update", Type: 2, Desc: "编辑客户", Path: "/customer/update", MenuCode: "menu_customer"},
		{Name: "删除客户", Code: "btn_customer_delete", Type: 2, Desc: "删除客户", Path: "/customer/delete", MenuCode: "menu_customer"},
		{Name: "查看激活客户", Code: "customer:active_list", Type: 2, Desc: "查看激活客户列表", Path: "/customer/active-list", MenuCode: "menu_customer"},
		// ERP仓库权限
		{Name: "查看仓库列表", Code: "warehouse:list", Type: 2, Desc: "查看仓库列表", Path: "/warehouse/list", MenuCode: "menu_warehouse"},
		{Name: "创建仓库", Code: "btn_warehouse_create", Type: 2, Desc: "创建仓库", Path: "/warehouse/create", MenuCode: "menu_warehouse"},
		{Name: "编辑仓库", Code: "btn_warehouse_update", Type: 2, Desc: "编辑仓库", Path: "/warehouse/update", MenuCode: "menu_warehouse"},
		{Name: "删除仓库", Code: "btn_warehouse_delete", Type: 2, Desc: "删除仓库", Path: "/warehouse/delete", MenuCode: "menu_warehouse"},
		{Name: "查看激活仓库", Code: "warehouse:active_list", Type: 2, Desc: "查看激活仓库列表", Path: "/warehouse/active-list", MenuCode: "menu_warehouse"},
		// ERP采购管理权限
		{Name: "查看采购订单", Code: "purchase:list", Type: 2, Desc: "查看采购订单列表", Path: "/purchase/order/list", MenuCode: "menu_purchase"},
		{Name: "创建采购订单", Code: "btn_purchase_create", Type: 2, Desc: "创建采购订单按钮", Path: "/purchase/order/create", MenuCode: "menu_purchase"},
		{Name: "审核采购订单", Code: "btn_purchase_approve", Type: 2, Desc: "审核采购订单按钮", Path: "/purchase/order/status", MenuCode: "menu_purchase"},
		{Name: "采购入库", Code: "btn_purchase_inbound", Type: 2, Desc: "采购入库按钮", Path: "/purchase/order/inbound", MenuCode: "menu_purchase"},
		{Name: "取消采购订单", Code: "btn_purchase_cancel", Type: 2, Desc: "取消采购订单按钮", Path: "/purchase/order/status", MenuCode: "menu_purchase"},
		{Name: "删除采购订单", Code: "btn_purchase_delete", Type: 2, Desc: "删除采购订单按钮", Path: "/purchase/order/delete", MenuCode: "menu_purchase"},
		// ERP销售管理权限
		{Name: "查看销售订单", Code: "sales:list", Type: 2, Desc: "查看销售订单列表", Path: "/sales/order/list", MenuCode: "menu_sales"},
		{Name: "创建销售订单", Code: "btn_sales_create", Type: 2, Desc: "创建销售订单按钮", Path: "/sales/order/create", MenuCode: "menu_sales"},
		{Name: "审核销售订单", Code: "btn_sales_approve", Type: 2, Desc: "审核销售订单按钮", Path: "/sales/order/status", MenuCode: "menu_sales"},
		{Name: "销售出库", Code: "btn_sales_outbound", Type: 2, Desc: "销售出库按钮", Path: "/inventory/sales-outbound", MenuCode: "menu_sales"},
		{Name: "取消销售订单", Code: "btn_sales_cancel", Type: 2, Desc: "取消销售订单按钮", Path: "/sales/order/status", MenuCode: "menu_sales"},
		{Name: "删除销售订单", Code: "btn_sales_delete", Type: 2, Desc: "删除销售订单按钮", Path: "/sales/order/delete", MenuCode: "menu_sales"},
		// ERP订单日志权限
		{Name: "查看订单日志", Code: "order-log:list", Type: 2, Desc: "查看订单日志按钮", Path: "/inventory/order-log/list", MenuCode: "menu_order_log"},
		// ERP库存管理权限
		{Name: "查看库存", Code: "inventory:list", Type: 2, Desc: "查看库存列表", Path: "/inventory/list", MenuCode: "menu_inventory"},
		{Name: "查看库存历史", Code: "inventory:history", Type: 2, Desc: "查看库存历史", Path: "/inventory/history", MenuCode: "menu_inventory"},
		// 库存调整权限
		{Name: "查看调整申请", Code: "inventory_adjust:list", Type: 2, Desc: "查看库存调整申请列表", Path: "/inventory/adjust/list", MenuCode: "menu_inventory_adjust"},
		{Name: "创建调整申请", Code: "inventory_adjust:create", Type: 2, Desc: "创建库存调整申请", Path: "/inventory/adjust/create", MenuCode: "menu_inventory_adjust"},
		{Name: "审核调整申请", Code: "inventory_adjust:approve", Type: 2, Desc: "审核库存调整申请", Path: "/inventory/adjust/approve", MenuCode: "menu_inventory_adjust"},
		{Name: "拒绝调整申请", Code: "inventory_adjust:reject", Type: 2, Desc: "拒绝库存调整申请", Path: "/inventory/adjust/reject", MenuCode: "menu_inventory_adjust"},
		// 库存调拨权限
		{Name: "查看调拨列表", Code: "inventory_transfer:list", Type: 2, Desc: "查看库存调拨列表", Path: "/inventory/transfer/list", MenuCode: "menu_inventory_transfer"},
		{Name: "创建调拨单", Code: "inventory_transfer:create", Type: 2, Desc: "创建库存调拨单", Path: "/inventory/transfer/create", MenuCode: "menu_inventory_transfer"},
		{Name: "编辑调拨单", Code: "inventory_transfer:update", Type: 2, Desc: "编辑库存调拨单", Path: "/inventory/transfer/update", MenuCode: "menu_inventory_transfer"},
		{Name: "删除调拨单", Code: "inventory_transfer:delete", Type: 2, Desc: "删除库存调拨单", Path: "/inventory/transfer/delete", MenuCode: "menu_inventory_transfer"},
		{Name: "审核调拨单", Code: "inventory_transfer:audit", Type: 2, Desc: "审核库存调拨单", Path: "/inventory/transfer/audit", MenuCode: "menu_inventory_transfer"},
		{Name: "执行调拨", Code: "inventory_transfer:execute", Type: 2, Desc: "执行库存调拨", Path: "/inventory/transfer/execute", MenuCode: "menu_inventory_transfer"},
		// 库存盘点权限
		{Name: "查看盘点列表", Code: "inventory_check:list", Type: 2, Desc: "查看盘点列表", Path: "/inventory/check/list", MenuCode: "menu_inventory_check"},
		{Name: "创建盘点单", Code: "inventory_check:create", Type: 2, Desc: "创建盘点单按钮", Path: "/inventory/check/create", MenuCode: "menu_inventory_check"},
		{Name: "更新盘点单", Code: "inventory_check:update", Type: 2, Desc: "更新盘点单按钮", Path: "/inventory/check/update", MenuCode: "menu_inventory_check"},
		{Name: "删除盘点单", Code: "inventory_check:delete", Type: 2, Desc: "删除盘点单按钮", Path: "/inventory/check/delete", MenuCode: "menu_inventory_check"},
		{Name: "提交盘点单", Code: "inventory_check:submit", Type: 2, Desc: "提交盘点单按钮", Path: "/inventory/check/submit", MenuCode: "menu_inventory_check"},
		{Name: "生成调整申请", Code: "inventory_check:generate", Type: 2, Desc: "生成调整申请按钮", Path: "/inventory/generate-adjust-from-check", MenuCode: "menu_inventory_check"},
		// 库存预警权限
		{Name: "查看库存预警", Code: "inventory_alert:list", Type: 2, Desc: "查看库存预警列表", Path: "/inventory/alert/list", MenuCode: "menu_inventory_alert"},
		{Name: "手动检查预警", Code: "inventory_alert:check", Type: 2, Desc: "手动检查库存预警", Path: "/inventory/alert/check", MenuCode: "menu_inventory_alert"},
		// 产品计量单位权限
		{Name: "查看计量单位", Code: "product_unit:list", Type: 2, Desc: "查看计量单位列表", Path: "/product/unit/list", MenuCode: "menu_product_unit"},
		{Name: "创建计量单位", Code: "product_unit:create", Type: 2, Desc: "创建计量单位", Path: "/product/unit/create", MenuCode: "menu_product_unit"},
		{Name: "编辑计量单位", Code: "product_unit:update", Type: 2, Desc: "编辑计量单位", Path: "/product/unit/update", MenuCode: "menu_product_unit"},
		{Name: "删除计量单位", Code: "product_unit:delete", Type: 2, Desc: "删除计量单位", Path: "/product/unit/delete", MenuCode: "menu_product_unit"},
		// ERP统计报表权限
		{Name: "查看ERP统计", Code: "erp:view", Type: 2, Desc: "查看ERP管理模块", Path: "/erp/statistics/dashboard", MenuCode: "menu_erp_report"},
		{Name: "查看ERP趋势", Code: "erp:trend", Type: 2, Desc: "查看ERP趋势", Path: "/erp/statistics/trend", MenuCode: "menu_erp_report"},
		{Name: "查看库存预警", Code: "erp:alert", Type: 2, Desc: "查看库存预警", Path: "/erp/statistics/inventory-alert", MenuCode: "menu_erp_report"},
		{Name: "查看热门产品", Code: "erp:top", Type: 2, Desc: "查看热门产品", Path: "/erp/statistics/top-products", MenuCode: "menu_erp_report"},
		{Name: "查看订单状态", Code: "erp:status", Type: 2, Desc: "查看订单状态", Path: "/erp/statistics/order-status", MenuCode: "menu_erp_report"},
		{Name: "查看业务统计", Code: "erp:business", Type: 2, Desc: "查看业务统计", Path: "/erp/statistics/business", MenuCode: "menu_erp_report"},
	}

	for _, p := range permissions {
		// 先检查权限是否已存在
		var existing Permission
		err := DB.Where("code = ?", p.Code).First(&existing).Error
		if err == nil {
			// 权限已存在，更新path和其他字段
			log.Printf("更新权限path: %s, path: %s", p.Code, p.Path)
			DB.Model(&existing).Updates(map[string]interface{}{
				"name": p.Name,
				"desc": p.Desc,
				"path": p.Path,
			})
			continue
		}

		// 查找关联的菜单
		var menu Menu
		menuID := uint(0)
		if p.MenuCode != "" {
			if err := DB.Where("code = ?", p.MenuCode).First(&menu).Error; err == nil {
				menuID = menu.ID
			}
		}

		// 创建权限
		newPerm := Permission{
			Name:   p.Name,
			Code:   p.Code,
			Desc:   p.Desc,
			Path:   p.Path,
			Status: 1,
		}

		if err := DB.Create(&newPerm).Error; err != nil {
			log.Printf("创建权限失败: %s, err: %v", p.Code, err)
		} else {
			log.Printf("创建权限: %s", p.Code)
			if menuID > 0 {
				mp := MenuPermission{MenuID: menuID, PermissionID: newPerm.ID}
				DB.Create(&mp)
			}
		}
	}

	log.Println("权限数据初始化完成")

	// 5. 为菜单分配对应权限
	if err := assignMenuPermissions(); err != nil {
		log.Printf("菜单权限关联失败: %v", err)
	}
	log.Println("菜单权限关联初始化完成")

	// 6. 创建默认管理员角色
	adminRole := Role{Name: "管理员", Code: "admin", Desc: "系统管理员", Status: 1}
	if err := DB.Create(&adminRole).Error; err != nil {
		return err
	}

	// 7. 为管理员角色分配所有权限和菜单
	var allPermissions []Permission
	var allMenus []Menu
	DB.Find(&allPermissions)
	DB.Find(&allMenus)
	DB.Model(&adminRole).Association("Permissions").Append(allPermissions)
	DB.Model(&adminRole).Association("Menus").Append(allMenus)

	// 8. 创建默认管理员用户
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	adminUser := User{
		Username: "admin",
		Password: string(hashedPassword),
		Nickname: "系统管理员",
		Email:    "admin@example.com",
		Phone:    "13800138000",
		Status:   1,
	}

	if err := DB.Create(&adminUser).Error; err != nil {
		return err
	}

	// 9. 为管理员用户分配角色
	DB.Model(&adminUser).Association("Roles").Append(&adminRole)

	log.Println("用户数据初始化成功")

	// 10. 创建测试角色
	if err := createTestRole(); err != nil {
		log.Printf("测试角色创建: %v", err)
	}

	return nil
}

// quickFixPermissionPaths 快速修复权限路径（确保与路由一致）
func quickFixPermissionPaths() error {
	log.Println("🔴 quickFixPermissionPaths() 被调用")
	pathFixes := []struct {
		oldPath string
		newPath string
	}{
		// 产品分类权限路径修复
		{"/category/list", "/product/category/list"},
		{"/category/create", "/product/category/create"},
		{"/category/update", "/product/category/update"},
		{"/category/delete", "/product/category/delete"},
	}

	for _, fix := range pathFixes {
		result := DB.Model(&Permission{}).Where("path = ?", fix.oldPath).Update("path", fix.newPath)
		if result.Error != nil {
			log.Printf("⚠️ 修复路径失败 %s -> %s: %v", fix.oldPath, fix.newPath, result.Error)
		} else if result.RowsAffected > 0 {
			log.Printf("✅ 修复路径: %s -> %s (影响 %d 条)", fix.oldPath, fix.newPath, result.RowsAffected)
		}
	}
	return nil
}

// quickFixPermissionCodes 快速修复权限code
func quickFixPermissionCodes() error {
	// 只检查缺失的权限code并修复
	permCodes := []struct {
		name string
		code string
	}{
		{"创建用户", "btn_user_create"},
		{"编辑用户", "btn_user_update"},
		{"删除用户", "btn_user_delete"},
		{"分配角色", "btn_user_assign"},
		{"创建角色", "btn_role_create"},
		{"编辑角色", "btn_role_update"},
		{"删除角色", "btn_role_delete"},
		{"分配权限", "btn_role_assign"},
		{"分配菜单", "btn_role_assign_menus"},
		{"创建权限", "btn_permission_create"},
		{"编辑权限", "btn_permission_update"},
		{"删除权限", "btn_permission_delete"},
		{"创建菜单", "btn_menu_create"},
		{"编辑菜单", "btn_menu_update"},
		{"删除菜单", "btn_menu_delete"},
		// ERP产品管理权限
		{"创建产品", "btn_product_create"},
		{"编辑产品", "btn_product_update"},
		{"删除产品", "btn_product_delete"},
		// ERP采购管理权限
		{"创建采购订单", "btn_purchase_create"},
		{"审核采购订单", "btn_purchase_approve"},
		{"采购入库", "btn_purchase_inbound"},
		{"取消采购订单", "btn_purchase_cancel"},
		{"删除采购订单", "btn_purchase_delete"},
		// ERP销售管理权限
		{"创建销售订单", "btn_sales_create"},
		{"审核销售订单", "btn_sales_approve"},
		{"销售出库", "btn_sales_outbound"},
		{"取消销售订单", "btn_sales_cancel"},
		{"删除销售订单", "btn_sales_delete"},
		// ERP订单日志权限
		{"查看订单日志", "btn_order_log_view"},
		// ERP库存管理权限
		{"库存调整", "btn_inventory_adjust"},
		// 库存调整权限
		{"查看调整申请", "inventory_adjust:list"},
		{"创建调整申请", "inventory_adjust:create"},
		{"审核调整申请", "inventory_adjust:approve"},
		{"拒绝调整申请", "inventory_adjust:reject"},
		// ERP管理父菜单权限
		{"查看ERP", "btn_erp_view"},
	}

	for _, p := range permCodes {
		var count int64
		DB.Model(&Permission{}).Where("name = ? AND code = ?", p.name, p.code).Count(&count)
		if count == 0 {
			// 权限缺失，尝试修复
			DB.Exec("UPDATE permissions SET code = ? WHERE name = ? AND code != ?", p.code, p.name, p.code)
		}
	}
	return nil
}

// quickSyncMenuPermissions 快速同步菜单权限关联
func quickSyncMenuPermissions() error {
	// 只在权限缺失时同步
	return assignMenuPermissions()
}

// ensureAdminMenus 确保管理员角色有所有菜单
func ensureAdminMenus() error {
	var adminRole Role
	if err := DB.Where("code = ?", "admin").First(&adminRole).Error; err != nil {
		return err
	}

	var allMenus []Menu
	if err := DB.Find(&allMenus).Error; err != nil {
		return err
	}

	return DB.Model(&adminRole).Association("Menus").Append(allMenus)
}

// ensureAdminPermissions 确保管理员角色有所有权限
func ensureAdminPermissions() error {
	log.Println("🔴 ensureAdminPermissions() 被调用")
	var adminRole Role
	if err := DB.Where("code = ?", "admin").First(&adminRole).Error; err != nil {
		log.Printf("❌ 未找到管理员角色: %v", err)
		return err
	}

	var allPermissions []Permission
	if err := DB.Find(&allPermissions).Error; err != nil {
		log.Printf("❌ 查询权限失败: %v", err)
		return err
	}

	log.Printf("📦 管理员角色: %s (id: %d)，准备分配 %d 个权限", adminRole.Name, adminRole.ID, len(allPermissions))
	return DB.Model(&adminRole).Association("Permissions").Append(allPermissions)
}

// ensureMissingMenus 确保数据库中不存在新功能菜单时创建它们
func ensureMissingMenus() error {
	menusToCheck := []struct {
		Code       string
		Name       string
		ParentCode string
		Path       string
		Component  string
		Icon       string
		Sort       int
	}{
		{Code: "menu_category", Name: "产品分类", ParentCode: "menu_erp", Path: "/product/category", Component: "product/Category", Icon: "Folder", Sort: 2},
		{Code: "menu_supplier", Name: "供应商管理", ParentCode: "menu_erp", Path: "/supplier", Component: "Supplier", Icon: "Shop", Sort: 3},
		{Code: "menu_customer", Name: "客户管理", ParentCode: "menu_erp", Path: "/customer", Component: "Customer", Icon: "UserFilled", Sort: 4},
		{Code: "menu_warehouse", Name: "仓库管理", ParentCode: "menu_erp", Path: "/warehouse", Component: "Warehouse", Icon: "OfficeBuilding", Sort: 5},
		{Code: "menu_inventory_adjust", Name: "库存调整", ParentCode: "menu_erp", Path: "/inventory/adjust-request", Component: "inventory/InventoryAdjustRequest", Icon: "RefreshCw", Sort: 9},
		{Code: "menu_inventory_transfer", Name: "库存调拨", ParentCode: "menu_erp", Path: "/inventory/transfer", Component: "inventory/InventoryTransfer", Icon: "ArrowRightCircle", Sort: 10},
		{Code: "menu_inventory_check", Name: "库存盘点", ParentCode: "menu_erp", Path: "/inventory/check", Component: "inventory/InventoryCheck", Icon: "DocumentChecked", Sort: 11},
		{Code: "menu_inventory_alert", Name: "库存预警", ParentCode: "menu_erp", Path: "/inventory/alert", Component: "inventory/InventoryAlert", Icon: "Bell", Sort: 12},
		{Code: "menu_product_unit", Name: "计量单位", ParentCode: "menu_erp", Path: "/product/unit", Component: "product/ProductUnit", Icon: "Calculator", Sort: 13},
		{Code: "menu_erp_report", Name: "预览报表", ParentCode: "menu_erp", Path: "/erp", Component: "Erp", Icon: "Eye", Sort: 14},
	}

	for _, m := range menusToCheck {
		var existing Menu
		err := DB.Where("code = ?", m.Code).First(&existing).Error
		if err != nil {
			// 菜单不存在，需要创建
			var parentMenu Menu
			if err := DB.Where("code = ?", m.ParentCode).First(&parentMenu).Error; err != nil {
				log.Printf("⚠️ 未找到父菜单 %s，跳过创建 %s", m.ParentCode, m.Name)
				continue
			}

			newMenu := Menu{
				Name:      m.Name,
				Code:      m.Code,
				Desc:      m.Name + "模块",
				ParentID:  parentMenu.ID,
				Path:      m.Path,
				Component: m.Component,
				Icon:      m.Icon,
				Sort:      m.Sort,
				Status:    1,
			}
			if err := DB.Create(&newMenu).Error; err != nil {
				log.Printf("⚠️ 创建菜单 %s 失败: %v", m.Name, err)
				continue
			}
			log.Printf("✅ 新增菜单: %s (id: %d)", newMenu.Name, newMenu.ID)
		}
	}

	// 确保所有权限存在
	permissionsToCheck := []struct {
		Code     string
		Name     string
		MenuCode string
	}{
		// 产品分类权限
		{Code: "category:list", Name: "查看分类列表", MenuCode: "menu_category"},
		{Code: "btn_category_create", Name: "创建分类", MenuCode: "menu_category"},
		{Code: "btn_category_update", Name: "编辑分类", MenuCode: "menu_category"},
		{Code: "btn_category_delete", Name: "删除分类", MenuCode: "menu_category"},
		// 供应商权限
		{Code: "supplier:list", Name: "查看供应商列表", MenuCode: "menu_supplier"},
		{Code: "btn_supplier_create", Name: "创建供应商", MenuCode: "menu_supplier"},
		{Code: "btn_supplier_update", Name: "编辑供应商", MenuCode: "menu_supplier"},
		{Code: "btn_supplier_delete", Name: "删除供应商", MenuCode: "menu_supplier"},
		// 客户权限
		{Code: "customer:list", Name: "查看客户列表", MenuCode: "menu_customer"},
		{Code: "btn_customer_create", Name: "创建客户", MenuCode: "menu_customer"},
		{Code: "btn_customer_update", Name: "编辑客户", MenuCode: "menu_customer"},
		{Code: "btn_customer_delete", Name: "删除客户", MenuCode: "menu_customer"},
		// 仓库权限
		{Code: "warehouse:list", Name: "查看仓库列表", MenuCode: "menu_warehouse"},
		{Code: "btn_warehouse_create", Name: "创建仓库", MenuCode: "menu_warehouse"},
		{Code: "btn_warehouse_update", Name: "编辑仓库", MenuCode: "menu_warehouse"},
		{Code: "btn_warehouse_delete", Name: "删除仓库", MenuCode: "menu_warehouse"},
		// 库存调整权限
		{Code: "inventory_adjust:list", Name: "查看调整申请", MenuCode: "menu_inventory_adjust"},
		{Code: "inventory_adjust:create", Name: "创建调整申请", MenuCode: "menu_inventory_adjust"},
		{Code: "inventory_adjust:approve", Name: "审核调整申请", MenuCode: "menu_inventory_adjust"},
		{Code: "inventory_adjust:reject", Name: "拒绝调整申请", MenuCode: "menu_inventory_adjust"},
		// 库存调拨权限
		{Code: "inventory_transfer:list", Name: "查看调拨列表", MenuCode: "menu_inventory_transfer"},
		{Code: "inventory_transfer:create", Name: "创建调拨单", MenuCode: "menu_inventory_transfer"},
		{Code: "inventory_transfer:update", Name: "编辑调拨单", MenuCode: "menu_inventory_transfer"},
		{Code: "inventory_transfer:delete", Name: "删除调拨单", MenuCode: "menu_inventory_transfer"},
		{Code: "inventory_transfer:audit", Name: "审核调拨单", MenuCode: "menu_inventory_transfer"},
		{Code: "inventory_transfer:execute", Name: "执行调拨", MenuCode: "menu_inventory_transfer"},
		// 库存盘点权限
		{Code: "inventory_check:list", Name: "查看盘点列表", MenuCode: "menu_inventory_check"},
		{Code: "inventory_check:create", Name: "创建盘点单", MenuCode: "menu_inventory_check"},
		{Code: "inventory_check:update", Name: "更新盘点单", MenuCode: "menu_inventory_check"},
		{Code: "inventory_check:delete", Name: "删除盘点单", MenuCode: "menu_inventory_check"},
		{Code: "inventory_check:submit", Name: "提交盘点单", MenuCode: "menu_inventory_check"},
		{Code: "inventory_check:generate", Name: "生成调整申请", MenuCode: "menu_inventory_check"},
		// 库存预警权限
		{Code: "inventory_alert:list", Name: "查看库存预警", MenuCode: "menu_inventory_alert"},
		{Code: "inventory_alert:check", Name: "手动检查预警", MenuCode: "menu_inventory_alert"},
		// 计量单位权限
		{Code: "product_unit:list", Name: "查看计量单位", MenuCode: "menu_product_unit"},
		{Code: "product_unit:create", Name: "创建计量单位", MenuCode: "menu_product_unit"},
		{Code: "product_unit:update", Name: "编辑计量单位", MenuCode: "menu_product_unit"},
		{Code: "product_unit:delete", Name: "删除计量单位", MenuCode: "menu_product_unit"},
		// ERP统计报表权限
		{Code: "erp:view", Name: "查看ERP统计", MenuCode: "menu_erp_report"},
		{Code: "erp:trend", Name: "查看ERP趋势", MenuCode: "menu_erp_report"},
		{Code: "erp:alert", Name: "查看库存预警", MenuCode: "menu_erp_report"},
		{Code: "erp:top", Name: "查看热门产品", MenuCode: "menu_erp_report"},
		{Code: "erp:status", Name: "查看订单状态", MenuCode: "menu_erp_report"},
		{Code: "erp:business", Name: "查看业务统计", MenuCode: "menu_erp_report"},
	}

	for _, p := range permissionsToCheck {
		var existing Permission
		err := DB.Where("code = ?", p.Code).First(&existing).Error
		if err != nil {
			// 权限不存在，需要创建
			var menu Menu
			if err := DB.Where("code = ?", p.MenuCode).First(&menu).Error; err != nil {
				log.Printf("⚠️ 未找到菜单 %s，跳过创建权限 %s", p.MenuCode, p.Name)
				continue
			}

			newPerm := Permission{
				Name:   p.Name,
				Code:   p.Code,
				Desc:   p.Name,
				Status: 1,
			}
			if err := DB.Create(&newPerm).Error; err != nil {
				log.Printf("⚠️ 创建权限 %s 失败: %v", p.Name, err)
				continue
			}
			mp := MenuPermission{MenuID: menu.ID, PermissionID: newPerm.ID}
			DB.Create(&mp)
			log.Printf("✅ 新增权限: %s (id: %d)", newPerm.Name, newPerm.ID)
		}
	}

	// 确保菜单-权限关联
	var invCheckMenu Menu
	if err := DB.Where("code = ?", "menu_inventory_check").First(&invCheckMenu).Error; err == nil {
		var perms []Permission
		DB.Where("code LIKE ?", "inventory_check:%").Find(&perms)
		for _, perm := range perms {
			var mp MenuPermission
			if err := DB.Where("menu_id = ? AND permission_id = ?", invCheckMenu.ID, perm.ID).First(&mp).Error; err != nil {
				mp = MenuPermission{MenuID: invCheckMenu.ID, PermissionID: perm.ID}
				DB.Create(&mp)
			}
		}
	}

	return nil
}

// ensureTestRole 确保测试角色存在（不修改已有权限）
func ensureTestRole() error {
	// 只确保角色存在，不修改已有权限
	var testRole Role
	err := DB.Where("code = ?", "test-role").First(&testRole).Error
	if err != nil {
		// 角色不存在，创建角色并分配初始权限
		log.Println("📦 测试角色不存在，创建并分配初始权限")
		return createTestRole()
	}
	log.Printf("📦 测试角色已存在: %s (id: %d)，跳过任何修改", testRole.Name, testRole.ID)
	return nil
}

// createTestRole 创建测试角色（只在首次创建时分配权限）
func createTestRole() error {
	// 查找测试角色
	var testRole Role
	err := DB.Where("code = ?", "test-role").First(&testRole).Error
	if err != nil {
		// 测试角色不存在，创建新角色
		testRole = Role{Name: "测试角色", Code: "test-role", Desc: "测试角色，用于权限测试", Status: 1}
		if err := DB.Create(&testRole).Error; err != nil {
			return err
		}
		log.Printf("✅ 创建测试角色: %s (id: %d)", testRole.Name, testRole.ID)

		// 仅在角色首次创建时分配菜单和权限
		if err := assignInitialPermissionsToTestRole(&testRole); err != nil {
			log.Printf("⚠️ 为测试角色分配初始权限失败: %v", err)
		}
	} else {
		log.Printf("📦 测试角色已存在: %s (id: %d)，跳过权限分配", testRole.Name, testRole.ID)
	}

	return nil
}

// assignInitialPermissionsToTestRole 仅在角色首次创建时为测试角色分配初始权限
func assignInitialPermissionsToTestRole(testRole *Role) error {
	log.Println("🔴 assignInitialPermissionsToTestRole() 被调用，角色ID:", testRole.ID)
	// 为测试角色分配菜单权限（所有主要菜单）
	var menus []Menu
	if err := DB.Where("code IN ?", []string{"menu_dashboard", "menu_user", "menu_role", "menu_permission", "menu_menu"}).Find(&menus).Error; err != nil {
		return err
	}
	if len(menus) > 0 {
		if err := DB.Model(testRole).Association("Menus").Append(menus); err != nil {
			log.Printf("⚠️ 为角色分配菜单失败: %v", err)
		}
		log.Printf("✅ 为测试角色分配了 %d 个菜单", len(menus))
	} else {
		log.Printf("⚠️ 未找到菜单，跳过分配")
	}

	// 批量为测试角色分配权限
	permCodes := []string{
		"btn_user_create", "btn_user_update", "btn_user_delete", "btn_user_assign",
		"btn_role_create", "btn_role_update", "btn_role_delete", "btn_role_assign",
		"btn_permission_update",
		"btn_menu_create", "btn_menu_update", "btn_menu_delete",
	}
	var perms []Permission
	if err := DB.Where("code IN ?", permCodes).Find(&perms).Error; err != nil {
		return err
	}
	DB.Model(testRole).Association("Permissions").Append(perms)
	log.Printf("✅ 为测试角色分配了 %d 个初始权限", len(perms))

	return nil
}

// assignMenuPermissions 为菜单分配对应权限
func assignMenuPermissions() error {
	// 用户管理菜单权限
	if err := assignPermissionsToMenu("menu_user", []string{
		"user:list",
		"btn_user_create",
		"btn_user_update",
		"btn_user_delete",
		"btn_user_assign",
	}); err != nil {
		return err
	}

	// 角色管理菜单权限
	if err := assignPermissionsToMenu("menu_role", []string{
		"role:list",
		"btn_role_create",
		"btn_role_update",
		"btn_role_delete",
		"btn_role_assign",
		"btn_role_assign_menus",
	}); err != nil {
		return err
	}

	// 权限管理菜单权限
	if err := assignPermissionsToMenu("menu_permission", []string{
		"permission:list",
		"btn_permission_create",
		"btn_permission_update",
		"btn_permission_delete",
	}); err != nil {
		return err
	}

	// 菜单管理菜单权限
	if err := assignPermissionsToMenu("menu_menu", []string{
		"menu:list",
		"menu:tree",
		"btn_menu_create",
		"btn_menu_update",
		"btn_menu_delete",
		"btn_menu_assign_permissions",
	}); err != nil {
		return err
	}

	// ERP产品管理菜单权限
	if err := assignPermissionsToMenu("menu_product", []string{
		"product:list",
		"btn_product_create",
		"btn_product_update",
		"btn_product_delete",
	}); err != nil {
		return err
	}

	// ERP产品分类菜单权限
	if err := assignPermissionsToMenu("menu_category", []string{
		"category:list",
		"btn_category_create",
		"btn_category_update",
		"btn_category_delete",
	}); err != nil {
		return err
	}

	// ERP供应商菜单权限
	if err := assignPermissionsToMenu("menu_supplier", []string{
		"supplier:list",
		"btn_supplier_create",
		"btn_supplier_update",
		"btn_supplier_delete",
	}); err != nil {
		return err
	}

	// ERP客户菜单权限
	if err := assignPermissionsToMenu("menu_customer", []string{
		"customer:list",
		"btn_customer_create",
		"btn_customer_update",
		"btn_customer_delete",
	}); err != nil {
		return err
	}

	// ERP仓库菜单权限
	if err := assignPermissionsToMenu("menu_warehouse", []string{
		"warehouse:list",
		"btn_warehouse_create",
		"btn_warehouse_update",
		"btn_warehouse_delete",
	}); err != nil {
		return err
	}

	// ERP采购管理菜单权限
	if err := assignPermissionsToMenu("menu_purchase", []string{
		"purchase:list",
		"btn_purchase_create",
		"btn_purchase_approve",
		"btn_purchase_inbound",
		"btn_purchase_cancel",
		"btn_purchase_delete",
	}); err != nil {
		return err
	}

	// ERP销售管理菜单权限
	if err := assignPermissionsToMenu("menu_sales", []string{
		"sales:list",
		"btn_sales_create",
		"btn_sales_approve",
		"btn_sales_outbound",
		"btn_sales_cancel",
		"btn_sales_delete",
	}); err != nil {
		return err
	}

	// ERP库存管理菜单权限
	if err := assignPermissionsToMenu("menu_inventory", []string{
		"inventory:list",
		"inventory:history",
	}); err != nil {
		return err
	}

	// 库存调整菜单权限
	if err := assignPermissionsToMenu("menu_inventory_adjust", []string{
		"inventory_adjust:list",
		"inventory_adjust:create",
		"inventory_adjust:approve",
		"inventory_adjust:reject",
	}); err != nil {
		return err
	}

	// 库存调拨菜单权限
	if err := assignPermissionsToMenu("menu_inventory_transfer", []string{
		"inventory_transfer:list",
		"inventory_transfer:create",
		"inventory_transfer:update",
		"inventory_transfer:delete",
		"inventory_transfer:audit",
		"inventory_transfer:execute",
	}); err != nil {
		return err
	}

	// 库存盘点菜单权限
	if err := assignPermissionsToMenu("menu_inventory_check", []string{
		"inventory_check:list",
		"inventory_check:create",
		"inventory_check:update",
		"inventory_check:delete",
		"inventory_check:submit",
		"inventory_check:generate",
	}); err != nil {
		return err
	}

	// 库存预警菜单权限
	if err := assignPermissionsToMenu("menu_inventory_alert", []string{
		"inventory_alert:list",
		"inventory_alert:check",
	}); err != nil {
		return err
	}

	// 计量单位菜单权限
	if err := assignPermissionsToMenu("menu_product_unit", []string{
		"product_unit:list",
		"product_unit:create",
		"product_unit:update",
		"product_unit:delete",
	}); err != nil {
		return err
	}

	// ERP统计报表菜单权限
	if err := assignPermissionsToMenu("menu_erp_report", []string{
		"erp:view",
		"erp:trend",
		"erp:alert",
		"erp:top",
		"erp:status",
		"erp:business",
	}); err != nil {
		return err
	}

	return nil
}

// assignPermissionsToMenu 根据菜单code和权限code列表，为菜单分配权限
func assignPermissionsToMenu(menuCode string, permissionCodes []string) error {
	var menu Menu
	if err := DB.Where("code = ?", menuCode).First(&menu).Error; err != nil {
		log.Printf("查找菜单失败: %s, err: %v", menuCode, err)
		return err
	}

	var permissions []Permission
	if err := DB.Where("code IN ?", permissionCodes).Find(&permissions).Error; err != nil {
		log.Printf("查找权限失败: %v, err: %v", permissionCodes, err)
		return err
	}

	if len(permissions) == 0 {
		log.Printf("未找到权限: %v", permissionCodes)
		return nil
	}

	if err := DB.Model(&menu).Association("Permissions").Append(permissions); err != nil {
		log.Printf("为菜单分配权限失败: %s, err: %v", menuCode, err)
		return err
	}

	log.Printf("为菜单 %s 分配了 %d 个权限", menuCode, len(permissions))
	return nil
}

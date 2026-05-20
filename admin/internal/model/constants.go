package model

// 用户状态常量
const (
	UserStatusDisabled = 0 // 禁用
	UserStatusEnabled  = 1 // 启用
)

// 采购订单状态常量
const (
	PurchaseStatusPending   = 1 // 待审核
	PurchaseStatusApproved  = 2 // 已审核
	PurchaseStatusInbound   = 3 // 已入库
	PurchaseStatusCancelled = 4 // 已取消
)

// 销售订单状态常量
const (
	SalesStatusPending  = 1 // 待审核
	SalesStatusApproved = 2 // 已审核
	SalesStatusOutbound = 3 // 已出库
	SalesStatusCancelled = 4 // 已取消
)

// 库存记录类型常量
const (
	InventoryTypeInbound  = 1 // 入库
	InventoryTypeOutbound = 2 // 出库
	InventoryTypeAdjust   = 3 // 调整
)

// 库存调整申请状态常量
const (
	AdjustStatusPending  = 1 // 待审核
	AdjustStatusApproved = 2 // 已审核
	AdjustStatusRejected = 4 // 已拒绝
)

// 库存调拨单状态常量
const (
	TransferStatusPending   = 1 // 待审核
	TransferStatusApproved  = 2 // 已审核
	TransferStatusCompleted = 3 // 已完成
	TransferStatusRejected  = 4 // 已拒绝
)

// 库存盘点单状态常量
const (
	CheckStatusPending   = 1 // 待盘点
	CheckStatusCompleted = 2 // 已完成
)

// 权限类型常量
const (
	PermissionTypeMenu   = 1 // 菜单权限
	PermissionTypeButton = 2 // 按钮权限
)

// 订单类型常量（用于库存记录）
const (
	OrderTypePurchase = 1 // 采购订单
	OrderTypeSales    = 2 // 销售订单
)

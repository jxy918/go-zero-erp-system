package model

import (
	"time"

	"gorm.io/gorm"
)

// User 用户表
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Username  string         `gorm:"size:50;not null;unique;index" json:"username"`
	Password  string         `gorm:"size:100;not null" json:"-"`
	Nickname  string         `gorm:"size:50" json:"nickname"`
	Email     string         `gorm:"size:100" json:"email"`
	Phone     string         `gorm:"size:20" json:"phone"`
	Status    int            `gorm:"default:1;index" json:"status"` // 1: 启用, 0: 禁用
	Roles     []Role         `gorm:"many2many:user_roles;" json:"roles"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Role 角色表
type Role struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:50;not null;unique;index" json:"name"`
	Code        string         `gorm:"size:50;not null;unique" json:"code"`
	Desc        string         `gorm:"size:200" json:"desc"`
	Status      int            `gorm:"default:1;index" json:"status"` // 1: 启用, 0: 禁用
	Permissions []Permission   `gorm:"many2many:role_permissions;" json:"permissions"`
	Menus       []Menu         `gorm:"many2many:role_menus;" json:"menus"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// Permission 权限表
type Permission struct {
	ID        uint           `gorm:"primaryKey;column:id" json:"id"`
	Name      string         `gorm:"column:name;size:100;not null" json:"name"`
	Code      string         `gorm:"column:code;size:100;not null;unique;index" json:"code"`
	Status    int            `gorm:"column:status;default:1;index" json:"status"`
	Desc      string         `gorm:"column:desc;size:200" json:"desc"`
	Sort      int            `gorm:"column:sort;default:0" json:"sort"`
	Path      string         `gorm:"column:path;size:200" json:"path"`
	CreatedAt time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
}

// Menu 菜单表
type Menu struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:50;not null" json:"name"`
	Code        string         `gorm:"size:50;not null;unique;index" json:"code"`
	ParentID    uint           `gorm:"default:0;index" json:"parent_id"`
	Path        string         `gorm:"size:200" json:"path"`
	Component   string         `gorm:"size:200" json:"component"`
	Icon        string         `gorm:"size:50" json:"icon"`
	Sort        int            `gorm:"default:0" json:"sort"`
	Status      int            `gorm:"default:1;index" json:"status"`
	Desc        string         `gorm:"size:200" json:"desc"`
	Children    []Menu         `gorm:"foreignKey:ParentID;references:ID" json:"children"`
	Permissions []Permission   `gorm:"many2many:menu_permissions;" json:"permissions"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// Category 产品分类表
type Category struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:100;not null;unique;index" json:"name"`
	Code      string         `gorm:"size:50;not null;unique" json:"code"`
	ParentID  uint           `gorm:"default:0;index" json:"parent_id"`
	Desc      string         `gorm:"size:200" json:"desc"`
	Sort      int            `gorm:"default:0" json:"sort"`
	Status    int            `gorm:"default:1;index" json:"status"`
	Children  []Category     `gorm:"foreignKey:ParentID;references:ID" json:"children"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Product 产品表
type Product struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:100;not null;index" json:"name"`
	Code        string         `gorm:"size:50;not null;unique;index" json:"code"`
	CategoryID  uint           `gorm:"not null;index" json:"category_id"`
	Spec        string         `gorm:"size:100" json:"spec"`
	Price       float64        `gorm:"type:decimal(10,2);not null" json:"price"`
	CostPrice   float64        `gorm:"type:decimal(10,2)" json:"cost_price"`
	MinStock    int            `gorm:"default:0" json:"min_stock"`
	MaxStock    int            `gorm:"default:99999" json:"max_stock"`
	SafetyStock int            `gorm:"default:0" json:"safety_stock"`
	Desc        string         `gorm:"size:200" json:"desc"`
	Status      int            `gorm:"default:1;index" json:"status"`
	Category    Category       `gorm:"foreignKey:CategoryID" json:"category"`
	Units       []ProductUnit  `gorm:"foreignKey:ProductID" json:"units"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// ProductUnit 产品计量单位表
type ProductUnit struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	ProductID uint           `gorm:"not null;index" json:"product_id"`
	UnitName  string         `gorm:"size:20;not null" json:"unit_name"`
	Ratio     float64        `gorm:"type:decimal(10,4);default:1" json:"ratio"`
	IsMain    int            `gorm:"default:0" json:"is_main"`
	Product   *Product       `gorm:"foreignKey:ProductID" json:"product"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Supplier 供应商表
type Supplier struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:100;not null;index" json:"name"`
	Code      string         `gorm:"size:50;not null;unique" json:"code"`
	Contact   string         `gorm:"size:50" json:"contact"`
	Phone     string         `gorm:"size:20" json:"phone"`
	Email     string         `gorm:"size:100" json:"email"`
	Address   string         `gorm:"size:200" json:"address"`
	Desc      string         `gorm:"size:200" json:"desc"`
	Status    int            `gorm:"default:1;index" json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Customer 客户表
type Customer struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:100;not null;index" json:"name"`
	Code      string         `gorm:"size:50;not null;unique" json:"code"`
	Contact   string         `gorm:"size:50" json:"contact"`
	Phone     string         `gorm:"size:20" json:"phone"`
	Email     string         `gorm:"size:100" json:"email"`
	Address   string         `gorm:"size:200" json:"address"`
	Desc      string         `gorm:"size:200" json:"desc"`
	Status    int            `gorm:"default:1;index" json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Warehouse 仓库表
type Warehouse struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:100;not null;unique;index" json:"name"`
	Code      string         `gorm:"size:50;not null;unique" json:"code"`
	Contact   string         `gorm:"size:50" json:"contact"`
	Phone     string         `gorm:"size:20" json:"phone"`
	Address   string         `gorm:"size:200" json:"address"`
	Desc      string         `gorm:"size:200" json:"desc"`
	Status    int            `gorm:"default:1;index" json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// PurchaseOrder 采购订单表
type PurchaseOrder struct {
	ID          uint                `gorm:"primaryKey" json:"id"`
	OrderNo     string              `gorm:"size:50;not null;unique;index" json:"order_no"`
	SupplierID  uint                `gorm:"not null;index" json:"supplier_id"`
	WarehouseID uint                `gorm:"index" json:"warehouse_id"`
	TotalAmount float64             `gorm:"type:decimal(12,2);not null" json:"total_amount"`
	Status      int                 `gorm:"default:1;index" json:"status"` // 1: 待审核, 2: 已审核, 3: 已入库, 4: 已取消
	Remark      string              `gorm:"size:500" json:"remark"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
	DeletedAt   gorm.DeletedAt      `gorm:"index" json:"-"`
	Supplier    Supplier            `gorm:"foreignKey:SupplierID;references:ID" json:"supplier"`
	Warehouse   Warehouse           `gorm:"foreignKey:WarehouseID;references:ID" json:"warehouse"`
	Items       []PurchaseOrderItem `gorm:"foreignKey:OrderID;references:ID" json:"items"`
}

// PurchaseOrderItem 采购订单项
type PurchaseOrderItem struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	OrderID   uint           `gorm:"not null;index" json:"order_id"`
	ProductID uint           `gorm:"not null;index" json:"product_id"`
	UnitID    uint           `json:"unit_id"`
	UnitName  string         `gorm:"size:50" json:"unit_name"`
	Ratio     int            `gorm:"default:1" json:"ratio"`
	Quantity  int            `gorm:"not null" json:"quantity"`
	BaseQty   int            `gorm:"not null" json:"base_qty"`
	UnitPrice float64        `gorm:"type:decimal(10,2);not null" json:"unit_price"`
	Amount    float64        `gorm:"type:decimal(12,2);not null" json:"amount"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Product   Product        `gorm:"foreignKey:ProductID;references:ID" json:"product"`
}

// SalesOrder 销售订单表
type SalesOrder struct {
	ID          uint             `gorm:"primaryKey" json:"id"`
	OrderNo     string           `gorm:"size:50;not null;unique;index" json:"order_no"`
	CustomerID  uint             `gorm:"not null;index" json:"customer_id"`
	WarehouseID uint             `gorm:"not null;index" json:"warehouse_id"`
	TotalAmount float64          `gorm:"type:decimal(12,2);not null" json:"total_amount"`
	Status      int              `gorm:"default:1;index" json:"status"` // 1: 待审核, 2: 已审核, 3: 已出库, 4: 已取消
	Remark      string           `gorm:"size:500" json:"remark"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
	DeletedAt   gorm.DeletedAt   `gorm:"index" json:"-"`
	Customer    Customer         `gorm:"foreignKey:CustomerID;references:ID" json:"customer"`
	Warehouse   Warehouse        `gorm:"foreignKey:WarehouseID;references:ID" json:"warehouse"`
	Items       []SalesOrderItem `gorm:"foreignKey:OrderID;references:ID" json:"items"`
}

// SalesOrderItem 销售订单项
type SalesOrderItem struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	OrderID   uint           `gorm:"not null;index" json:"order_id"`
	ProductID uint           `gorm:"not null;index" json:"product_id"`
	UnitID    uint           `json:"unit_id"`
	UnitName  string         `gorm:"size:50" json:"unit_name"`
	Ratio     int            `gorm:"default:1" json:"ratio"`
	Quantity  int            `gorm:"not null" json:"quantity"`
	BaseQty   int            `gorm:"not null" json:"base_qty"`
	UnitPrice float64        `gorm:"type:decimal(10,2);not null" json:"unit_price"`
	Amount    float64        `gorm:"type:decimal(12,2);not null" json:"amount"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Product   Product        `gorm:"foreignKey:ProductID;references:ID" json:"product"`
}

// InventoryRecord 库存记录表
type InventoryRecord struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	ProductID   uint           `gorm:"not null;index" json:"product_id"`
	WarehouseID uint           `gorm:"not null;index" json:"warehouse_id"`
	Quantity    int            `gorm:"default:0" json:"quantity"`
	Type        int            `gorm:"default:1;index" json:"type"`       // 1: 入库, 2: 出库, 3: 调整
	OrderID     uint           `gorm:"index" json:"order_id"`             // 关联订单ID
	OrderType   int            `gorm:"default:0;index" json:"order_type"` // 1: 采购订单, 2: 销售订单
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Product     Product        `gorm:"foreignKey:ProductID;references:ID" json:"product"`
	Warehouse   Warehouse      `gorm:"foreignKey:WarehouseID;references:ID" json:"warehouse"`
}

// WarehouseInventory 仓库库存缓存表
type WarehouseInventory struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	ProductID   uint      `gorm:"not null;uniqueIndex:idx_product_warehouse" json:"product_id"`
	WarehouseID uint      `gorm:"not null;uniqueIndex:idx_product_warehouse" json:"warehouse_id"`
	Quantity    int       `gorm:"default:0" json:"quantity"`
	UpdatedAt   time.Time `json:"updated_at"`
	Product     Product   `gorm:"foreignKey:ProductID;references:ID" json:"product"`
	Warehouse   Warehouse `gorm:"foreignKey:WarehouseID;references:ID" json:"warehouse"`
}

// InventoryAdjustRequest 库存调整申请表
type InventoryAdjustRequest struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	RequestNo   string         `gorm:"size:50;not null;unique;index" json:"request_no"`
	ProductID   uint           `gorm:"not null;index" json:"product_id"`
	WarehouseID uint           `gorm:"not null;index" json:"warehouse_id"`
	BeforeQty   int            `gorm:"not null" json:"before_qty"`
	Quantity    int            `gorm:"not null" json:"quantity"`
	AfterQty    int            `gorm:"not null" json:"after_qty"`
	Type        int            `gorm:"not null" json:"type"` // 1: 增加, 2: 减少
	Reason      string         `gorm:"size:500" json:"reason"`
	Status      int            `gorm:"default:1;index" json:"status"` // 1: 待审核, 2: 已审核, 4: 已拒绝
	ApplicantID uint           `gorm:"not null;index" json:"applicant_id"`
	ApproverID  uint           `json:"approver_id"`
	ApproveTime *time.Time     `json:"approve_time"`
	ApproveNote string         `gorm:"size:500" json:"approve_note"`
	Product     Product        `gorm:"foreignKey:ProductID;references:ID" json:"product"`
	Warehouse   Warehouse      `gorm:"foreignKey:WarehouseID;references:ID" json:"warehouse"`
	Applicant   User           `gorm:"foreignKey:ApplicantID;references:ID" json:"applicant"`
	Approver    User           `gorm:"foreignKey:ApproverID;references:ID" json:"approver"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// MenuPermission 菜单权限关联表
type MenuPermission struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	MenuID       uint           `gorm:"not null;index" json:"menu_id"`
	PermissionID uint           `gorm:"not null;index" json:"permission_id"`
	CreatedAt    time.Time      `json:"created_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// RolePermission 角色权限关联表
type RolePermission struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	RoleID       uint           `gorm:"not null;index" json:"role_id"`
	PermissionID uint           `gorm:"not null;index" json:"permission_id"`
	CreatedAt    time.Time      `json:"created_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// RoleMenu 角色菜单关联表
type RoleMenu struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	RoleID    uint           `gorm:"not null;index" json:"role_id"`
	MenuID    uint           `gorm:"not null;index" json:"menu_id"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// UserRole 用户角色关联表
type UserRole struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"not null;index" json:"user_id"`
	RoleID    uint           `gorm:"not null;index" json:"role_id"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

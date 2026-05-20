// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"myproject/admin/internal/config"
	"myproject/admin/internal/model"

	"gorm.io/gorm"
)

type ServiceContext struct {
	Config                  config.Config
	DB                      *gorm.DB
	UserModel               *model.UserModel
	RoleModel               *model.RoleModel
	PermissionModel         *model.PermissionModel
	ActivityModel           model.ActivityModel
	MenuModel               model.MenuModel
	ProductModel            *model.ProductModel
	ProductUnitModel        *model.ProductUnitModel
	CategoryModel           *model.CategoryModel
	SupplierModel           *model.SupplierModel
	CustomerModel           *model.CustomerModel
	WarehouseModel          *model.WarehouseModel
	PurchaseModel           *model.PurchaseModel
	SalesModel              *model.SalesModel
	InventoryModel          *model.InventoryModel
	InventoryChangeModel    *model.InventoryChangeModel
	InventoryAdjustReqModel *model.InventoryAdjustRequestModel
	InventoryCheckModel     *model.InventoryCheckModel
	InventoryTransferModel  *model.InventoryTransferModel
	OrderLogModel           *model.OrderLogModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 初始化数据库连接
	if err := model.InitDB(c.DB.DataSource); err != nil {
		panic(err)
	}

	return &ServiceContext{
		Config:                  c,
		DB:                      model.DB,
		UserModel:               model.NewUserModel(model.DB),
		RoleModel:               model.NewRoleModel(model.DB),
		PermissionModel:         model.NewPermissionModel(model.DB),
		ActivityModel:           model.NewActivityModel(model.DB),
		MenuModel:               model.NewMenuModel(model.DB),
		ProductModel:            model.NewProductModel(model.DB),
		ProductUnitModel:        model.NewProductUnitModel(model.DB),
		CategoryModel:           model.NewCategoryModel(model.DB),
		SupplierModel:           model.NewSupplierModel(model.DB),
		CustomerModel:           model.NewCustomerModel(model.DB),
		WarehouseModel:          model.NewWarehouseModel(model.DB),
		PurchaseModel:           model.NewPurchaseModel(model.DB),
		SalesModel:              model.NewSalesModel(model.DB),
		InventoryModel:          model.NewInventoryModel(model.DB),
		InventoryChangeModel:    model.NewInventoryChangeModel(model.DB),
		InventoryAdjustReqModel: model.NewInventoryAdjustRequestModel(model.DB),
		InventoryCheckModel:     model.NewInventoryCheckModel(model.DB),
		InventoryTransferModel:  model.NewInventoryTransferModel(model.DB),
		OrderLogModel:           model.NewOrderLogModel(model.DB),
	}
}

package logic

import (
	"context"
	"time"

	"myproject/admin/internal/model"
	"myproject/admin/internal/svc"
	"myproject/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetErpBusinessLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewGetErpBusinessLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetErpBusinessLogic {
	return &GetErpBusinessLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *GetErpBusinessLogic) GetErpBusiness() (resp *types.ErpBusinessResponse, err error) {
	var supplierCount int64
	err = model.DB.Model(&model.Supplier{}).Where("status = ?", 1).Count(&supplierCount).Error
	if err != nil {
		return nil, err
	}

	var customerCount int64
	err = model.DB.Model(&model.Customer{}).Where("status = ?", 1).Count(&customerCount).Error
	if err != nil {
		return nil, err
	}

	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	var activeSupplierCount int64
	err = model.DB.Table("purchase_orders po").
		Distinct("po.supplier_id").
		Where("po.created_at >= ? AND po.status != ?", startOfMonth, 4).
		Count(&activeSupplierCount).Error
	if err != nil {
		return nil, err
	}

	var activeCustomerCount int64
	err = model.DB.Table("sales_orders so").
		Distinct("so.customer_id").
		Where("so.created_at >= ? AND so.status != ?", startOfMonth, 4).
		Count(&activeCustomerCount).Error
	if err != nil {
		return nil, err
	}

	var productCount int64
	err = model.DB.Model(&model.Product{}).Where("status = ?", 1).Count(&productCount).Error
	if err != nil {
		return nil, err
	}

	var warehouseCount int64
	err = model.DB.Model(&model.Warehouse{}).Where("status = ?", 1).Count(&warehouseCount).Error
	if err != nil {
		return nil, err
	}

	var lowStockProductCount int64
	err = model.DB.Table("warehouse_inventories wi").
		Joins("LEFT JOIN products p ON wi.product_id = p.id").
		Where("wi.quantity < p.safety_stock").
		Distinct("wi.product_id").
		Count(&lowStockProductCount).Error
	if err != nil {
		return nil, err
	}

	resp = &types.ErpBusinessResponse{
		SupplierCount:        int(supplierCount),
		CustomerCount:        int(customerCount),
		ActiveSupplierCount:  int(activeSupplierCount),
		ActiveCustomerCount:  int(activeCustomerCount),
		ProductCount:         int(productCount),
		WarehouseCount:       int(warehouseCount),
		LowStockProductCount: int(lowStockProductCount),
	}

	return resp, nil
}

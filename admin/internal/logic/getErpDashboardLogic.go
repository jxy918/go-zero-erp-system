// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"myproject/admin/internal/model"
	"myproject/admin/internal/svc"
	"myproject/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetErpDashboardLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetErpDashboardLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetErpDashboardLogic {
	return &GetErpDashboardLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetErpDashboardLogic) GetErpDashboard() (resp *types.ErpDashboardResponse, err error) {
	var totalProducts int64
	err = l.svcCtx.DB.Model(&model.Product{}).Where("status = ?", 1).Count(&totalProducts).Error
	if err != nil {
		return nil, err
	}

	var totalSuppliers int64
	err = l.svcCtx.DB.Model(&model.Supplier{}).Where("status = ?", 1).Count(&totalSuppliers).Error
	if err != nil {
		return nil, err
	}

	var totalCustomers int64
	err = l.svcCtx.DB.Model(&model.Customer{}).Where("status = ?", 1).Count(&totalCustomers).Error
	if err != nil {
		return nil, err
	}

	var purchaseOrderCount int64
	err = l.svcCtx.DB.Model(&model.PurchaseOrder{}).Count(&purchaseOrderCount).Error
	if err != nil {
		return nil, err
	}

	var salesOrderCount int64
	err = l.svcCtx.DB.Model(&model.SalesOrder{}).Count(&salesOrderCount).Error
	if err != nil {
		return nil, err
	}

	return &types.ErpDashboardResponse{
		TotalProducts:  int(totalProducts),
		TotalSuppliers: int(totalSuppliers),
		TotalCustomers: int(totalCustomers),
		TotalOrders:    int(purchaseOrderCount + salesOrderCount),
	}, nil
}

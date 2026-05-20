package logic

import (
	"context"
	"time"

	"myproject/admin/internal/model"
	"myproject/admin/internal/svc"
	"myproject/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetErpOverviewLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewGetErpOverviewLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetErpOverviewLogic {
	return &GetErpOverviewLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *GetErpOverviewLogic) GetErpOverview() (resp *types.ErpOverviewResponse, err error) {
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	var purchaseAmount float64
	err = model.DB.Model(&model.PurchaseOrder{}).
		Where("created_at >= ?", startOfMonth).
		Select("COALESCE(SUM(total_amount), 0)").
		Find(&purchaseAmount).Error
	if err != nil {
		return nil, err
	}

	var salesAmount float64
	err = model.DB.Model(&model.SalesOrder{}).
		Where("created_at >= ?", startOfMonth).
		Select("COALESCE(SUM(total_amount), 0)").
		Find(&salesAmount).Error
	if err != nil {
		return nil, err
	}

	var totalInventory int
	err = model.DB.Model(&model.InventoryRecord{}).
		Select("COALESCE(SUM(quantity), 0)").
		Find(&totalInventory).Error
	if err != nil {
		return nil, err
	}

	var totalOrders int64
	err = model.DB.Model(&model.PurchaseOrder{}).
		Where("created_at >= ?", startOfMonth).
		Count(&totalOrders).Error
	if err != nil {
		return nil, err
	}

	var salesOrderCount int64
	err = model.DB.Model(&model.SalesOrder{}).
		Where("created_at >= ?", startOfMonth).
		Count(&salesOrderCount).Error
	if err != nil {
		return nil, err
	}

	var userCount int64
	err = model.DB.Model(&model.User{}).Count(&userCount).Error
	if err != nil {
		return nil, err
	}

	var productCount int64
	err = model.DB.Model(&model.Product{}).Count(&productCount).Error
	if err != nil {
		return nil, err
	}

	resp = &types.ErpOverviewResponse{
		PurchaseAmount: purchaseAmount,
		SalesAmount:    salesAmount,
		TotalInventory: totalInventory,
		TotalOrders:    int(totalOrders + salesOrderCount),
		UserCount:      userCount,
		ProductCount:   productCount,
	}

	return resp, nil
}

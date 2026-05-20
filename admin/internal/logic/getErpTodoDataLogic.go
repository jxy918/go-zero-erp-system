package logic

import (
	"context"

	"myproject/admin/internal/model"
	"myproject/admin/internal/svc"
	"myproject/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetErpTodoDataLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewGetErpTodoDataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetErpTodoDataLogic {
	return &GetErpTodoDataLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *GetErpTodoDataLogic) GetErpTodoData() (resp *types.ErpTodoData, err error) {
	var pendingCheckOrders int64
	err = model.DB.Model(&model.InventoryCheck{}).
		Where("status IN (1, 2)").
		Count(&pendingCheckOrders).Error
	if err != nil {
		l.Logger.Error("查询待盘点单失败:", err)
	}

	var pendingAdjustRequests int64
	err = model.DB.Model(&model.InventoryAdjustRequest{}).
		Where("status = 1").
		Count(&pendingAdjustRequests).Error
	if err != nil {
		l.Logger.Error("查询待处理调整失败:", err)
	}

	var pendingTransfers int64
	err = model.DB.Model(&model.InventoryTransfer{}).
		Where("status = 1").
		Count(&pendingTransfers).Error
	if err != nil {
		l.Logger.Error("查询待审核调拨失败:", err)
	}

	var pendingPurchases int64
	err = model.DB.Model(&model.PurchaseOrder{}).
		Where("status = 1").
		Count(&pendingPurchases).Error
	if err != nil {
		l.Logger.Error("查询待审核采购失败:", err)
	}

	var pendingSales int64
	err = model.DB.Model(&model.SalesOrder{}).
		Where("status = 1").
		Count(&pendingSales).Error
	if err != nil {
		l.Logger.Error("查询待审核销售失败:", err)
	}

	resp = &types.ErpTodoData{
		PendingCheckOrders:    int(pendingCheckOrders),
		PendingAdjustRequests: int(pendingAdjustRequests),
		PendingTransfers:      int(pendingTransfers),
		PendingPurchases:      int(pendingPurchases),
		PendingSales:          int(pendingSales),
	}

	return resp, nil
}

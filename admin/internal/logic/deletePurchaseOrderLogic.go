package logic

import (
	"context"
	"errors"

	"myproject/admin/internal/svc"
	"myproject/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeletePurchaseOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeletePurchaseOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeletePurchaseOrderLogic {
	return &DeletePurchaseOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeletePurchaseOrderLogic) DeletePurchaseOrder(req *types.DeletePurchaseOrderRequest) (resp *types.PurchaseOrderInfo, err error) {
	order, err := l.svcCtx.PurchaseModel.GetByID(req.ID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, errors.New("订单不存在")
	}

	if err := l.svcCtx.PurchaseModel.Delete(req.ID); err != nil {
		return nil, err
	}

	resp = &types.PurchaseOrderInfo{
		ID:          order.ID,
		OrderNo:     order.OrderNo,
		SupplierID:  order.SupplierID,
		WarehouseID: order.WarehouseID,
		TotalAmount: order.TotalAmount,
		Status:      order.Status,
		Remark:      order.Remark,
		CreatedAt:   order.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   order.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	return resp, nil
}

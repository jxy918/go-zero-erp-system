package logic

import (
	"context"
	"errors"

	"myproject/admin/internal/svc"
	"myproject/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteSalesOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteSalesOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteSalesOrderLogic {
	return &DeleteSalesOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteSalesOrderLogic) DeleteSalesOrder(req *types.DeleteSalesOrderRequest) (resp *types.SalesOrderInfo, err error) {
	order, err := l.svcCtx.SalesModel.GetByID(req.ID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, errors.New("订单不存在")
	}

	if order.Status != 1 {
		return nil, errors.New("只有待审核的订单才能删除")
	}

	if err := l.svcCtx.SalesModel.Delete(req.ID); err != nil {
		return nil, err
	}

	return nil, nil
}

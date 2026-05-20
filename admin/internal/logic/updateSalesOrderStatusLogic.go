// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"errors"

	"myproject/admin/internal/middleware"
	"myproject/admin/internal/model"
	"myproject/admin/internal/svc"
	"myproject/admin/internal/types"
	"myproject/admin/internal/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateSalesOrderStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateSalesOrderStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateSalesOrderStatusLogic {
	return &UpdateSalesOrderStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateSalesOrderStatusLogic) UpdateSalesOrderStatus(req *types.SalesOrderStatusRequest) (resp *types.EmptyResponse, err error) {
	order, err := l.svcCtx.SalesModel.GetByID(req.ID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, errors.New("订单不存在")
	}

	beforeStatus := order.Status
	if err := l.svcCtx.SalesModel.UpdateStatus(req.ID, req.Status); err != nil {
		return nil, err
	}

	operatorID, _ := l.ctx.Value(util.UserIDKey).(uint)
	operatorName := middleware.GetUsername(l.ctx)

	log := &model.OrderLog{
		OrderID:      req.ID,
		OrderType:    2,
		BeforeStatus: beforeStatus,
		AfterStatus:  req.Status,
		OperatorID:   operatorID,
		OperatorName: operatorName,
		Remark:       "销售订单状态变更",
	}
	if err := l.svcCtx.OrderLogModel.Insert(log); err != nil {
		return nil, err
	}

	return &types.EmptyResponse{}, nil
}

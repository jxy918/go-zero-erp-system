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

type UpdateSalesOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateSalesOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateSalesOrderLogic {
	return &UpdateSalesOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateSalesOrderLogic) UpdateSalesOrder(req *types.UpdateSalesOrderRequest) (resp *types.SalesOrderInfo, err error) {
	order, err := l.svcCtx.SalesModel.GetByID(req.ID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, errors.New("订单不存在")
	}

	updateData := make(map[string]interface{})

	if req.Status != 0 && req.Status != order.Status {
		updateData["status"] = req.Status

		operatorID, _ := l.ctx.Value(util.UserIDKey).(uint)
		operatorName := middleware.GetUsername(l.ctx)

		log := &model.OrderLog{
			OrderID:      req.ID,
			OrderType:    2,
			BeforeStatus: order.Status,
			AfterStatus:  req.Status,
			OperatorID:   operatorID,
			OperatorName: operatorName,
			Remark:       "销售订单状态变更",
		}
		if err := l.svcCtx.OrderLogModel.Insert(log); err != nil {
			return nil, err
		}
	}
	if req.Remark != "" {
		updateData["remark"] = req.Remark
	}

	if len(updateData) == 0 {
		return nil, errors.New("没有需要更新的字段")
	}

	if err := l.svcCtx.SalesModel.UpdateFields(req.ID, updateData); err != nil {
		return nil, err
	}

	return nil, nil
}

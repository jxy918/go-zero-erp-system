package logic

import (
	"context"

	"myproject/admin/internal/svc"
	"myproject/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListOrderLogLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListOrderLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListOrderLogLogic {
	return &ListOrderLogLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListOrderLogLogic) ListOrderLog(req *types.ListOrderLogRequest) (resp *types.ListOrderLogResponse, err error) {
	logs, err := l.svcCtx.OrderLogModel.List(req.OrderID, uint(req.OrderType), req.OperatorID, req.StartTime, req.EndTime)
	if err != nil {
		return nil, err
	}

	resp = &types.ListOrderLogResponse{
		Logs:  make([]types.OrderLogInfo, 0, len(logs)),
		Total: int64(len(logs)),
	}

	for _, log := range logs {
		resp.Logs = append(resp.Logs, types.OrderLogInfo{
			ID:               log.ID,
			OrderID:          log.OrderID,
			OrderType:        log.OrderType,
			OrderTypeDesc:    l.getOrderTypeDesc(log.OrderType),
			BeforeStatus:     log.BeforeStatus,
			BeforeStatusDesc: l.getStatusDesc(log.OrderType, log.BeforeStatus),
			AfterStatus:      log.AfterStatus,
			AfterStatusDesc:  l.getStatusDesc(log.OrderType, log.AfterStatus),
			OperatorID:       log.OperatorID,
			OperatorName:     log.OperatorName,
			Remark:           log.Remark,
			CreatedAt:        log.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return resp, nil
}

func (l *ListOrderLogLogic) getOrderTypeDesc(orderType int) string {
	switch orderType {
	case 1:
		return "采购订单"
	case 2:
		return "销售订单"
	default:
		return "未知"
	}
}

func (l *ListOrderLogLogic) getStatusDesc(orderType, status int) string {
	switch status {
	case 1:
		return "待审核"
	case 2:
		return "已审核"
	case 3:
		if orderType == 1 {
			return "已入库"
		}
		return "已出库"
	case 4:
		return "已取消"
	default:
		return "未知"
	}
}

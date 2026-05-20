package logic

import (
	"context"
	"errors"

	"myproject/admin/internal/model"
	"myproject/admin/internal/svc"
	"myproject/admin/internal/types"
	"myproject/admin/internal/util"
)

type RejectInventoryAdjustRequestLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRejectInventoryAdjustRequestLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RejectInventoryAdjustRequestLogic {
	return &RejectInventoryAdjustRequestLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RejectInventoryAdjustRequestLogic) RejectInventoryAdjustRequest(req *types.RejectInventoryAdjustRequest) (*types.InventoryAdjustRequestInfo, error) {
	approverID := util.GetUserID(l.ctx)

	adjustReq, err := l.svcCtx.InventoryAdjustReqModel.GetByID(req.ID)
	if err != nil {
		return nil, errors.New("申请不存在")
	}

	if adjustReq.Status != 1 {
		return nil, errors.New("只能拒绝待处理的申请")
	}

	if err := l.svcCtx.InventoryAdjustReqModel.Reject(req.ID, approverID, req.Note); err != nil {
		return nil, err
	}

	adjustReq, _ = l.svcCtx.InventoryAdjustReqModel.GetByID(req.ID)
	return l.convertToResponse(adjustReq), nil
}

func (l *RejectInventoryAdjustRequestLogic) convertToResponse(req *model.InventoryAdjustRequest) *types.InventoryAdjustRequestInfo {
	return &types.InventoryAdjustRequestInfo{
		ID:          req.ID,
		RequestNo:   req.RequestNo,
		ProductID:   req.ProductID,
		WarehouseID: req.WarehouseID,
		BeforeQty:   req.BeforeQty,
		Quantity:    req.Quantity,
		AfterQty:    req.AfterQty,
		Type:        req.Type,
		TypeDesc:    l.getTypeDesc(req.Type),
		Reason:      req.Reason,
		Status:      req.Status,
		StatusDesc:  l.getStatusDesc(req.Status),
		ApplicantID: req.ApplicantID,
		ApproverID:  req.ApproverID,
	}
}

func (l *RejectInventoryAdjustRequestLogic) getTypeDesc(t int) string {
	switch t {
	case 1:
		return "盘盈"
	case 2:
		return "盘亏"
	case 3:
		return "调拨"
	default:
		return "其他"
	}
}

func (l *RejectInventoryAdjustRequestLogic) getStatusDesc(s int) string {
	switch s {
	case 1:
		return "待审核"
	case 2:
		return "已审核"
	case 3:
		return "已拒绝"
	default:
		return "未知"
	}
}

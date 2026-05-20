package logic

import (
	"context"
	"errors"

	"myproject/admin/internal/model"
	"myproject/admin/internal/svc"
	"myproject/admin/internal/types"
	"myproject/admin/internal/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateInventoryAdjustLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateInventoryAdjustLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateInventoryAdjustLogic {
	return &CreateInventoryAdjustLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateInventoryAdjustLogic) CreateInventoryAdjust(req *types.CreateInventoryAdjustRequest) (*types.InventoryAdjustRequestInfo, error) {
	userID := util.GetUserID(l.ctx)

	beforeQty, err := l.svcCtx.InventoryModel.GetStockByProductAndWarehouse(req.ProductID, req.WarehouseID)
	if err != nil {
		beforeQty = 0
	}

	var quantity = req.Quantity
	switch req.Type {
	case 1:
		if quantity < 0 {
			quantity = -quantity
		}
	case 2:
		if quantity > 0 {
			quantity = -quantity
		}
	case 4:
	}

	if quantity == 0 {
		return nil, errors.New("调整数量不能为0")
	}

	afterQty := beforeQty + quantity
	if afterQty < 0 {
		return nil, errors.New("调整后库存不能为负数")
	}

	adjustReq := &model.InventoryAdjustRequest{
		ProductID:   req.ProductID,
		WarehouseID: req.WarehouseID,
		BeforeQty:   beforeQty,
		Quantity:    quantity,
		AfterQty:    afterQty,
		Type:        req.Type,
		Reason:      req.Reason,
		Status:      1,
		ApplicantID: userID,
	}

	if err := l.svcCtx.InventoryAdjustReqModel.Create(adjustReq); err != nil {
		return nil, err
	}

	return l.convertToResponseWithQuantity(adjustReq, quantity)
}

func (l *CreateInventoryAdjustLogic) convertToResponse(req *model.InventoryAdjustRequest) (*types.InventoryAdjustRequestInfo, error) {
	return l.convertToResponseWithQuantity(req, req.Quantity)
}

func (l *CreateInventoryAdjustLogic) convertToResponseWithQuantity(req *model.InventoryAdjustRequest, quantity int) (*types.InventoryAdjustRequestInfo, error) {
	return &types.InventoryAdjustRequestInfo{
		ID:          req.ID,
		RequestNo:   req.RequestNo,
		ProductID:   req.ProductID,
		WarehouseID: req.WarehouseID,
		BeforeQty:   req.BeforeQty,
		Quantity:    quantity,
		AfterQty:    req.AfterQty,
		Type:        req.Type,
		TypeDesc:    l.getTypeDesc(req.Type),
		Reason:      req.Reason,
		Status:      req.Status,
		StatusDesc:  l.getStatusDesc(req.Status),
		ApplicantID: req.ApplicantID,
		ApproverID:  req.ApproverID,
	}, nil
}

func (l *CreateInventoryAdjustLogic) getTypeDesc(t int) string {
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

func (l *CreateInventoryAdjustLogic) getStatusDesc(s int) string {
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

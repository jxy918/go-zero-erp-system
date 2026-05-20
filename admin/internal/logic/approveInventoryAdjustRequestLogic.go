package logic

import (
	"context"
	"errors"
	"fmt"

	"myproject/admin/internal/model"
	"myproject/admin/internal/svc"
	"myproject/admin/internal/types"
	"myproject/admin/internal/util"
)

type ApproveInventoryAdjustRequestLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApproveInventoryAdjustRequestLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApproveInventoryAdjustRequestLogic {
	return &ApproveInventoryAdjustRequestLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApproveInventoryAdjustRequestLogic) ApproveInventoryAdjustRequest(req *types.ApproveInventoryAdjustRequest) (*types.InventoryAdjustRequestInfo, error) {
	approverID := util.GetUserID(l.ctx)

	adjustReq, err := l.svcCtx.InventoryAdjustReqModel.GetByID(req.ID)
	if err != nil {
		return nil, errors.New("申请不存在")
	}

	if adjustReq.Status != 1 {
		return nil, errors.New("只能审核待处理的申请")
	}

	if err := l.svcCtx.InventoryAdjustReqModel.Approve(req.ID, approverID, req.Note); err != nil {
		return nil, err
	}

	fmt.Println("==========================================")
	fmt.Println("=== 审核库存调整申请 START ===")
	fmt.Printf("请求ID: %d\n", req.ID)
	fmt.Printf("数据库读取的调整申请:\n")
	fmt.Printf("  ID: %d\n", adjustReq.ID)
	fmt.Printf("  RequestNo: %s\n", adjustReq.RequestNo)
	fmt.Printf("  Type: %d (1=盘盈, 2=盘亏, 4=其他)\n", adjustReq.Type)
	fmt.Printf("  Quantity: %d\n", adjustReq.Quantity)
	fmt.Printf("  BeforeQty: %d\n", adjustReq.BeforeQty)
	fmt.Printf("  AfterQty: %d\n", adjustReq.AfterQty)
	fmt.Printf("  ProductID: %d\n", adjustReq.ProductID)
	fmt.Printf("  WarehouseID: %d\n", adjustReq.WarehouseID)
	fmt.Println("=== 审核库存调整申请 END ===")

	if err := l.executeAdjust(adjustReq); err != nil {
		return nil, err
	}

	adjustReq, _ = l.svcCtx.InventoryAdjustReqModel.GetByID(req.ID)
	return l.convertToResponse(adjustReq), nil
}

func (l *ApproveInventoryAdjustRequestLogic) executeAdjust(req *model.InventoryAdjustRequest) error {
	var changeType int
	var quantity = req.Quantity
	var reason string

	fmt.Println("==========================================")
	fmt.Println("=== executeAdjust START ===")
	fmt.Printf("输入参数:\n")
	fmt.Printf("  Type: %d\n", req.Type)
	fmt.Printf("  Quantity: %d\n", req.Quantity)
	fmt.Printf("  RequestNo: %s\n", req.RequestNo)
	fmt.Printf("  BeforeQty: %d\n", req.BeforeQty)
	fmt.Printf("  AfterQty: %d\n", req.AfterQty)

	switch req.Type {
	case 1: // 盘盈
		fmt.Println("  -> 匹配到 case 1 (盘盈)")
		changeType = 1 // 入库
		if quantity < 0 {
			fmt.Printf("  -> quantity=%d < 0，取反\n", quantity)
			quantity = -quantity
		}
		reason = fmt.Sprintf("[盘盈] 库存调整申请#%s，调整前:%d，调整数量:+%d，调整后:%d，原因:%s",
			req.RequestNo, req.BeforeQty, quantity, req.AfterQty, req.Reason)
	case 2: // 盘亏
		fmt.Println("  -> 匹配到 case 2 (盘亏)")
		changeType = 2 // 出库
		if quantity > 0 {
			fmt.Printf("  -> quantity=%d > 0，取反\n", quantity)
			quantity = -quantity
		}
		reason = fmt.Sprintf("[盘亏] 库存调整申请#%s，调整前:%d，调整数量:%d，调整后:%d，原因:%s",
			req.RequestNo, req.BeforeQty, quantity, req.AfterQty, req.Reason)
	case 4: // 其他
		fmt.Println("  -> 匹配到 case 4 (其他)")
		if quantity >= 0 {
			fmt.Printf("  -> quantity=%d >= 0，入库类型\n", quantity)
			changeType = 1 // 入库
			reason = fmt.Sprintf("[其他-入库] 库存调整申请#%s，调整前:%d，调整数量:+%d，调整后:%d，原因:%s",
				req.RequestNo, req.BeforeQty, quantity, req.AfterQty, req.Reason)
		} else {
			fmt.Printf("  -> quantity=%d < 0，出库类型\n", quantity)
			changeType = 2 // 出库
			reason = fmt.Sprintf("[其他-出库] 库存调整申请#%s，调整前:%d，调整数量:%d，调整后:%d，原因:%s",
				req.RequestNo, req.BeforeQty, quantity, req.AfterQty, req.Reason)
		}
	default:
		fmt.Printf("  -> 匹配到 default case (Type=%d)\n", req.Type)
		if quantity >= 0 {
			changeType = 1 // 入库
			reason = fmt.Sprintf("[库存调整-入库] 申请#%s，调整前:%d，调整数量:+%d，调整后:%d，原因:%s",
				req.RequestNo, req.BeforeQty, quantity, req.AfterQty, req.Reason)
		} else {
			changeType = 2 // 出库
			reason = fmt.Sprintf("[库存调整-出库] 申请#%s，调整前:%d，调整数量:%d，调整后:%d，原因:%s",
				req.RequestNo, req.BeforeQty, quantity, req.AfterQty, req.Reason)
		}
	}

	fmt.Println("=== executeAdjust 计算结果 ===")
	fmt.Printf("  changeType: %d (1=入库, 2=出库)\n", changeType)
	fmt.Printf("  quantity: %d\n", quantity)
	fmt.Printf("  reason: %s\n", reason)
	fmt.Println("=== executeAdjust END ===")

	return l.svcCtx.InventoryModel.AdjustStock(req.ProductID, req.WarehouseID, quantity, req.ID, changeType, reason)
}

func (l *ApproveInventoryAdjustRequestLogic) convertToResponse(req *model.InventoryAdjustRequest) *types.InventoryAdjustRequestInfo {
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

func (l *ApproveInventoryAdjustRequestLogic) getTypeDesc(t int) string {
	switch t {
	case 1:
		return "盘盈"
	case 2:
		return "盘亏"
	case 4:
		return "其他"
	default:
		return "其他"
	}
}

func (l *ApproveInventoryAdjustRequestLogic) getStatusDesc(s int) string {
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

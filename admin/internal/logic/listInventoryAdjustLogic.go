package logic

import (
	"context"

	"myproject/admin/internal/svc"
	"myproject/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListInventoryAdjustLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListInventoryAdjustLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListInventoryAdjustLogic {
	return &ListInventoryAdjustLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListInventoryAdjustLogic) ListInventoryAdjust(req *types.ListInventoryAdjustRequest) (*types.ListInventoryAdjustResponse, error) {
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}

	requests, total, err := l.svcCtx.InventoryAdjustReqModel.List(page, pageSize, req.Status)
	if err != nil {
		return nil, err
	}

	var requestInfos []types.InventoryAdjustRequestInfo
	for _, r := range requests {
		info := types.InventoryAdjustRequestInfo{
			ID:          r.ID,
			RequestNo:   r.RequestNo,
			ProductID:   r.ProductID,
			WarehouseID: r.WarehouseID,
			BeforeQty:   r.BeforeQty,
			Quantity:    r.Quantity,
			AfterQty:    r.AfterQty,
			Type:        r.Type,
			TypeDesc:    l.getTypeDesc(r.Type),
			Reason:      r.Reason,
			Status:      r.Status,
			StatusDesc:  l.getStatusDesc(r.Status),
			ApplicantID: r.ApplicantID,
			ApproverID:  r.ApproverID,
			ApproveNote: r.ApproveNote,
			CreatedAt:   r.CreatedAt.Format("2006-01-02 15:04:05"),
		}

		if r.Product.ID > 0 {
			info.Product = types.ProductInfo{
				ID:       r.Product.ID,
				Name:     r.Product.Name,
				Code:     r.Product.Code,
				Spec:     r.Product.Spec,
				Price:    r.Product.Price,
				MainUnit: "",
			}
		}
		if r.Warehouse.ID > 0 {
			info.Warehouse = types.WarehouseInfo{
				ID:      r.Warehouse.ID,
				Name:    r.Warehouse.Name,
				Address: r.Warehouse.Address,
			}
		}
		if r.Applicant.ID > 0 {
			info.Applicant = types.UserInfo{
				ID:       r.Applicant.ID,
				Username: r.Applicant.Username,
				Nickname: r.Applicant.Nickname,
			}
		}
		if r.Approver.ID > 0 {
			info.Approver = types.UserInfo{
				ID:       r.Approver.ID,
				Username: r.Approver.Username,
				Nickname: r.Approver.Nickname,
			}
		}
		if r.ApproveTime != nil {
			info.ApproveTime = r.ApproveTime.Format("2006-01-02 15:04:05")
		}

		requestInfos = append(requestInfos, info)
	}

	return &types.ListInventoryAdjustResponse{
		Requests: requestInfos,
		Total:    total,
	}, nil
}

func (l *ListInventoryAdjustLogic) getTypeDesc(t int) string {
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

func (l *ListInventoryAdjustLogic) getStatusDesc(s int) string {
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

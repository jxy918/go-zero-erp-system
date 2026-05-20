// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"myproject/admin/internal/svc"
	"myproject/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListInventoryCheckLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListInventoryCheckLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListInventoryCheckLogic {
	return &ListInventoryCheckLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListInventoryCheckLogic) ListInventoryCheck(req *types.ListInventoryCheckRequest) (resp *types.ListInventoryCheckResponse, err error) {
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	checks, total, err := l.svcCtx.InventoryCheckModel.List(page, pageSize, req.CheckNo, req.WarehouseID, req.Status)
	if err != nil {
		return nil, err
	}

	items := make([]types.InventoryCheckInfo, 0)
	for _, check := range checks {
		items = append(items, *convertInventoryCheckToInfo(l.svcCtx.InventoryModel, &check))
	}

	return &types.ListInventoryCheckResponse{
		Checks: items,
		Total:  total,
	}, nil
}

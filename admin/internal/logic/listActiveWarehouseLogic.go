// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"myproject/admin/internal/svc"
	"myproject/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListActiveWarehouseLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListActiveWarehouseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListActiveWarehouseLogic {
	return &ListActiveWarehouseLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListActiveWarehouseLogic) ListActiveWarehouse(req *types.ListWarehouseRequest) (resp *types.ListWarehouseResponse, err error) {
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	warehouses, total, err := l.svcCtx.WarehouseModel.ListActive(page, pageSize, req.Name, req.Name)
	if err != nil {
		return nil, err
	}

	resp = &types.ListWarehouseResponse{
		Warehouses: make([]types.WarehouseInfo, 0, len(warehouses)),
		Total:      total,
	}

	for _, warehouse := range warehouses {
		resp.Warehouses = append(resp.Warehouses, types.WarehouseInfo{
			ID:        warehouse.ID,
			Name:      warehouse.Name,
			Code:      warehouse.Code,
			Address:   warehouse.Address,
			Contact:   warehouse.Contact,
			Phone:     warehouse.Phone,
			Desc:      warehouse.Desc,
			Status:    warehouse.Status,
			CreatedAt: warehouse.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: warehouse.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return resp, nil
}

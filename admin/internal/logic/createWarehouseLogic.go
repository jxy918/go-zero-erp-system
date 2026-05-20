// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"myproject/admin/internal/model"
	"myproject/admin/internal/svc"
	"myproject/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateWarehouseLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateWarehouseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateWarehouseLogic {
	return &CreateWarehouseLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateWarehouseLogic) CreateWarehouse(req *types.CreateWarehouseRequest) (resp *types.WarehouseInfo, err error) {
	warehouse := &model.Warehouse{
		Name:    req.Name,
		Code:    req.Code,
		Address: req.Address,
		Contact: req.Contact,
		Phone:   req.Phone,
		Desc:    req.Desc,
		Status:  req.Status,
	}

	if err := l.svcCtx.WarehouseModel.Create(warehouse); err != nil {
		return nil, err
	}

	resp = &types.WarehouseInfo{
		ID:          warehouse.ID,
		Name:        warehouse.Name,
		Code:        warehouse.Code,
		Address:     warehouse.Address,
		Contact:     warehouse.Contact,
		Phone:       warehouse.Phone,
		Desc:        warehouse.Desc,
		Status:      warehouse.Status,
		CreatedAt:   warehouse.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   warehouse.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	return resp, nil
}

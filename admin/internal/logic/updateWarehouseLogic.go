// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"errors"

	"myproject/admin/internal/svc"
	"myproject/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateWarehouseLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateWarehouseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateWarehouseLogic {
	return &UpdateWarehouseLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateWarehouseLogic) UpdateWarehouse(req *types.UpdateWarehouseRequest) (resp *types.WarehouseInfo, err error) {
	warehouse, err := l.svcCtx.WarehouseModel.GetByID(req.ID)
	if err != nil {
		return nil, err
	}
	if warehouse == nil {
		return nil, errors.New("仓库不存在")
	}

	if req.Name != "" {
		warehouse.Name = req.Name
	}
	if req.Code != "" {
		warehouse.Code = req.Code
	}
	if req.Address != "" {
		warehouse.Address = req.Address
	}
	if req.Contact != "" {
		warehouse.Contact = req.Contact
	}
	if req.Phone != "" {
		warehouse.Phone = req.Phone
	}
	if req.Desc != "" {
		warehouse.Desc = req.Desc
	}
	if req.Status >= 0 {
		warehouse.Status = req.Status
	}

	if err := l.svcCtx.WarehouseModel.Update(warehouse); err != nil {
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

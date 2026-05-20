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

type GetWarehouseLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetWarehouseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetWarehouseLogic {
	return &GetWarehouseLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetWarehouseLogic) GetWarehouse(req *types.GetWarehouseRequest) (resp *types.WarehouseInfo, err error) {
	warehouse, err := l.svcCtx.WarehouseModel.GetByID(req.ID)
	if err != nil {
		return nil, err
	}
	if warehouse == nil {
		return nil, errors.New("仓库不存在")
	}

	resp = &types.WarehouseInfo{
		ID:        warehouse.ID,
		Name:      warehouse.Name,
		Code:      warehouse.Code,
		Address:   warehouse.Address,
		Contact:   warehouse.Contact,
		Phone:     warehouse.Phone,
		Status:    warehouse.Status,
		Desc:      warehouse.Desc,
		CreatedAt: warehouse.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: warehouse.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	return resp, nil
}

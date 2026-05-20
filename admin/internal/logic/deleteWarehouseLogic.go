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

type DeleteWarehouseLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteWarehouseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteWarehouseLogic {
	return &DeleteWarehouseLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteWarehouseLogic) DeleteWarehouse(req *types.DeleteWarehouseRequest) (resp *types.WarehouseInfo, err error) {
	warehouse, err := l.svcCtx.WarehouseModel.GetByID(req.ID)
	if err != nil {
		return nil, err
	}
	if warehouse == nil {
		return nil, errors.New("仓库不存在")
	}

	if err := l.svcCtx.WarehouseModel.Delete(req.ID); err != nil {
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

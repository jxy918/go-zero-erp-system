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

type UpdateInventoryCheckStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateInventoryCheckStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateInventoryCheckStatusLogic {
	return &UpdateInventoryCheckStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateInventoryCheckStatusLogic) UpdateInventoryCheckStatus(req *types.InventoryCheckStatusRequest) (resp *types.EmptyResponse, err error) {
	check, err := l.svcCtx.InventoryCheckModel.GetByID(req.ID)
	if err != nil {
		return nil, err
	}
	if check == nil {
		return nil, errors.New("盘点单不存在")
	}

	check.Status = req.Status
	if err := l.svcCtx.InventoryCheckModel.Update(check); err != nil {
		return nil, err
	}

	return &types.EmptyResponse{}, nil
}

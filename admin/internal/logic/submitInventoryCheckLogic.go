// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"myproject/admin/internal/model"
	"myproject/admin/internal/svc"
	"myproject/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type SubmitInventoryCheckLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSubmitInventoryCheckLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SubmitInventoryCheckLogic {
	return &SubmitInventoryCheckLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SubmitInventoryCheckLogic) SubmitInventoryCheck(req *types.SubmitInventoryCheckRequest) (resp *types.InventoryCheckResponse, err error) {
	check, err := l.svcCtx.InventoryCheckModel.GetByID(req.ID)
	if err != nil {
		return nil, err
	}
	if check == nil {
		return &types.InventoryCheckResponse{}, nil
	}
	if check.Status != 1 && check.Status != 2 {
		return &types.InventoryCheckResponse{}, nil
	}

	totalDiff := 0
	for _, item := range check.Items {
		totalDiff += item.DiffQty
	}

	check.Status = 4
	check.TotalDiff = totalDiff

	err = model.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(check).Error; err != nil {
			return err
		}

		return tx.Model(&model.InventoryCheckItem{}).
			Where("check_id = ?", req.ID).
			Update("status", 2).Error
	})

	if err != nil {
		return nil, err
	}

	check, err = l.svcCtx.InventoryCheckModel.GetByID(req.ID)
	if err != nil {
		return nil, err
	}

	return &types.InventoryCheckResponse{
		Check: *convertInventoryCheckToInfo(l.svcCtx.InventoryModel, check),
	}, nil
}

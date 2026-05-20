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

type UpdateInventoryCheckLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateInventoryCheckLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateInventoryCheckLogic {
	return &UpdateInventoryCheckLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateInventoryCheckLogic) UpdateInventoryCheck(req *types.UpdateInventoryCheckRequest) (resp *types.InventoryCheckResponse, err error) {
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

	check.Remark = req.Remark
	if check.Status == 1 && len(req.Items) > 0 {
		check.Status = 2
	}

	err = model.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(check).Error; err != nil {
			return err
		}

		for _, itemUpdate := range req.Items {
			var item model.InventoryCheckItem
			if err := tx.Where("id = ? AND check_id = ?", itemUpdate.ID, req.ID).First(&item).Error; err != nil {
				return err
			}

			diffQty := itemUpdate.ActualQty - item.SystemQty
			if err := tx.Model(&model.InventoryCheckItem{}).
				Where("id = ? AND check_id = ?", itemUpdate.ID, req.ID).
				Updates(map[string]interface{}{
					"actual_qty": itemUpdate.ActualQty,
					"diff_qty":   diffQty,
				}).Error; err != nil {
				return err
			}
		}

		return nil
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

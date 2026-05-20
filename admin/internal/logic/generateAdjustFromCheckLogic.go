// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"fmt"

	"myproject/admin/internal/model"
	"myproject/admin/internal/svc"
	"myproject/admin/internal/types"
	"myproject/admin/internal/util"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type GenerateAdjustFromCheckLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGenerateAdjustFromCheckLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GenerateAdjustFromCheckLogic {
	return &GenerateAdjustFromCheckLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GenerateAdjustFromCheckLogic) GenerateAdjustFromCheck(req *types.GenerateAdjustFromCheckRequest) (resp *types.GenerateAdjustFromCheckResponse, err error) {
	check, err := l.svcCtx.InventoryCheckModel.GetByID(req.ID)
	if err != nil {
		return nil, err
	}
	if check == nil {
		return nil, fmt.Errorf("盘点单不存在")
	}
	if check.Status != 4 {
		return nil, fmt.Errorf("只能对已提交的盘点单生成调整单")
	}

	applicantID := util.GetUserID(l.ctx)

	var adjustIds []uint
	adjustCount := 0

	err = model.DB.Transaction(func(tx *gorm.DB) error {
		for _, item := range check.Items {
			if item.DiffQty != 0 {
				adjustType := 1 // 1: 盘盈
				if item.DiffQty < 0 {
					adjustType = 2 // 2: 盘亏
				}

				// 获取调整前的库存
				var beforeQty int
				err := tx.Model(&model.InventoryRecord{}).
					Select("COALESCE(SUM(quantity), 0)").
					Where("product_id = ? AND warehouse_id = ?", item.ProductID, check.WarehouseID).
					Scan(&beforeQty).Error
				if err != nil {
					return err
				}

				adjustQty := item.DiffQty

				adjust := &model.InventoryAdjustRequest{
					WarehouseID: check.WarehouseID,
					ProductID:   item.ProductID,
					BeforeQty:   beforeQty,
					Quantity:    adjustQty,
					AfterQty:    beforeQty + item.DiffQty,
					Type:        adjustType,
					Reason:      fmt.Sprintf("盘点差异调整，盘点单号：%s", check.CheckNo),
					ApplicantID: applicantID,
					Status:      1,
				}

				adjustModel := model.NewInventoryAdjustRequestModel(tx)
				if err := adjustModel.Create(adjust); err != nil {
					return err
				}

				adjustIds = append(adjustIds, adjust.ID)
				adjustCount++
			}
		}

		check.Status = 3
		if err := tx.Save(check).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &types.GenerateAdjustFromCheckResponse{
		Success:     true,
		AdjustCount: adjustCount,
		AdjustIds:   adjustIds,
		Message:     fmt.Sprintf("成功生成%d个库存调整申请", adjustCount),
	}, nil
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

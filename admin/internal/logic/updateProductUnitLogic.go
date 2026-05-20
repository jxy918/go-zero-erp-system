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

type UpdateProductUnitLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateProductUnitLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateProductUnitLogic {
	return &UpdateProductUnitLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateProductUnitLogic) UpdateProductUnit(req *types.UpdateProductUnitRequest) (resp *types.ProductUnitInfo, err error) {
	var unit model.ProductUnit
	err = l.svcCtx.DB.Where("id = ?", req.ID).First(&unit).Error
	if err != nil {
		return nil, err
	}

	if req.IsMain == 1 && unit.IsMain != 1 {
		err = l.svcCtx.DB.Model(&model.ProductUnit{}).
			Where("product_id = ? AND is_main = 1", unit.ProductID).
			Update("is_main", 0).Error
		if err != nil {
			return nil, err
		}
	}

	if req.UnitName != "" {
		unit.UnitName = req.UnitName
	}
	if req.Ratio > 0 {
		unit.Ratio = req.Ratio
	}
	if req.IsMain >= 0 {
		unit.IsMain = req.IsMain
	}

	err = l.svcCtx.DB.Save(&unit).Error
	if err != nil {
		return nil, err
	}

	var product model.Product
	err = l.svcCtx.DB.Where("id = ?", unit.ProductID).First(&product).Error
	if err != nil {
		return nil, err
	}

	return &types.ProductUnitInfo{
		ID:          unit.ID,
		ProductID:   unit.ProductID,
		ProductName: product.Name,
		UnitName:    unit.UnitName,
		Ratio:       unit.Ratio,
		IsMain:      unit.IsMain,
		CreatedAt:   unit.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   unit.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

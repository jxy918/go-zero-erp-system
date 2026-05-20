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

type DeleteProductUnitLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteProductUnitLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteProductUnitLogic {
	return &DeleteProductUnitLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteProductUnitLogic) DeleteProductUnit(req *types.DeleteProductUnitRequest) (resp *types.ProductUnitInfo, err error) {
	var unit model.ProductUnit
	err = l.svcCtx.DB.Where("id = ?", req.ID).First(&unit).Error
	if err != nil {
		return nil, err
	}

	var product model.Product
	err = l.svcCtx.DB.Where("id = ?", unit.ProductID).First(&product).Error
	if err != nil {
		return nil, err
	}

	err = l.svcCtx.DB.Delete(&unit).Error
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

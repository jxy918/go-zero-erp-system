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

type CreateProductUnitLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateProductUnitLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateProductUnitLogic {
	return &CreateProductUnitLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateProductUnitLogic) CreateProductUnit(req *types.CreateProductUnitRequest) (resp *types.ProductUnitInfo, err error) {
	var product model.Product
	err = l.svcCtx.DB.Where("id = ? AND status = 1", req.ProductID).First(&product).Error
	if err != nil {
		return nil, err
	}

	if req.IsMain == 1 {
		err = l.svcCtx.DB.Model(&model.ProductUnit{}).
			Where("product_id = ? AND is_main = 1", req.ProductID).
			Update("is_main", 0).Error
		if err != nil {
			return nil, err
		}
	}

	ratio := req.Ratio
	if ratio <= 0 {
		ratio = 1
	}

	unit := model.ProductUnit{
		ProductID: req.ProductID,
		UnitName:  req.UnitName,
		Ratio:     ratio,
		IsMain:    req.IsMain,
	}

	err = l.svcCtx.DB.Create(&unit).Error
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

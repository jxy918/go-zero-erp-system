package logic

import (
	"context"
	"errors"

	"myproject/admin/internal/svc"
	"myproject/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetProductLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewGetProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProductLogic {
	return &GetProductLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *GetProductLogic) GetProduct(req *types.GetProductRequest) (resp *types.ProductInfo, err error) {
	product, err := l.svcCtx.ProductModel.GetByID(req.ID)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, errors.New("产品不存在")
	}

	var mainUnit string
	for _, unit := range product.Units {
		if unit.IsMain == 1 {
			mainUnit = unit.UnitName
			break
		}
	}

	stock, _ := l.svcCtx.InventoryModel.GetProductTotalStock(product.ID)

	resp = &types.ProductInfo{
		ID:          product.ID,
		Name:        product.Name,
		Code:        product.Code,
		Spec:        product.Spec,
		CategoryID:  product.CategoryID,
		Unit:        mainUnit,
		MinStock:    product.MinStock,
		MaxStock:    product.MaxStock,
		SafetyStock: product.SafetyStock,
		CostPrice:   product.CostPrice,
		SalePrice:   product.Price,
		Price:       product.Price,
		Status:      product.Status,
		Desc:        product.Desc,
		Stock:       stock,
		MainUnit:    mainUnit,
		Units:       []types.ProductUnitInfo{},
		CreatedAt:   product.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   product.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	if product.Category.ID > 0 {
		resp.Category = types.ProductCategoryInfo{
			ID:        product.Category.ID,
			Name:      product.Category.Name,
			Code:      product.Category.Code,
			ParentID:  product.Category.ParentID,
			Sort:      product.Category.Sort,
			Status:    product.Category.Status,
			Desc:      product.Category.Desc,
			Children:  []types.ProductCategoryInfo{},
			CreatedAt: product.Category.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: product.Category.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	return resp, nil
}

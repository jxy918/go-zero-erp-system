package logic

import (
	"context"

	"myproject/admin/internal/svc"
	"myproject/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListProductLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewListProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListProductLogic {
	return &ListProductLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *ListProductLogic) ListProduct(req *types.ListProductRequest) (resp *types.ListProductResponse, err error) {
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	products, total, err := l.svcCtx.ProductModel.List(page, pageSize, req.Name, req.Code)
	if err != nil {
		return nil, err
	}

	resp = &types.ListProductResponse{
		Products: make([]types.ProductInfo, 0, len(products)),
		Total:    total,
	}

	for _, product := range products {
		var mainUnit string
		for _, unit := range product.Units {
			if unit.IsMain == 1 {
				mainUnit = unit.UnitName
				break
			}
		}

		stock, _ := l.svcCtx.InventoryModel.GetProductTotalStock(product.ID)

		item := types.ProductInfo{
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
			item.Category = types.ProductCategoryInfo{
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

		resp.Products = append(resp.Products, item)
	}

	return resp, nil
}

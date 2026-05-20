// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"myproject/admin/internal/svc"
	"myproject/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListActiveProductLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListActiveProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListActiveProductLogic {
	return &ListActiveProductLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListActiveProductLogic) ListActiveProduct(req *types.ListProductRequest) (resp *types.ListProductResponse, err error) {
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 100
	}

	name := req.Name
	code := req.Code

	products, total, err := l.svcCtx.ProductModel.ListActive(page, pageSize, name, code)
	if err != nil {
		return nil, err
	}

	resp = &types.ListProductResponse{
		Products: make([]types.ProductInfo, 0, len(products)),
		Total:    total,
	}

	for _, product := range products {
		units := make([]types.ProductUnitInfo, 0, len(product.Units))
		for _, unit := range product.Units {
			units = append(units, types.ProductUnitInfo{
				ID:        unit.ID,
				ProductID: unit.ProductID,
				UnitName:  unit.UnitName,
				Ratio:     unit.Ratio,
				IsMain:    unit.IsMain,
				CreatedAt: unit.CreatedAt.Format("2006-01-02 15:04:05"),
				UpdatedAt: unit.UpdatedAt.Format("2006-01-02 15:04:05"),
			})
		}

		category := types.ProductCategoryInfo{}
		if product.Category.ID > 0 {
			category = types.ProductCategoryInfo{
				ID:        product.Category.ID,
				Name:      product.Category.Name,
				Code:      product.Category.Code,
				ParentID:  product.Category.ParentID,
				Sort:      product.Category.Sort,
				Status:    product.Category.Status,
				Desc:      product.Category.Desc,
				CreatedAt: product.Category.CreatedAt.Format("2006-01-02 15:04:05"),
				UpdatedAt: product.Category.UpdatedAt.Format("2006-01-02 15:04:05"),
			}
		}

		mainUnit := ""
		if len(units) > 0 {
			for _, unit := range units {
				if unit.IsMain == 1 {
					mainUnit = unit.UnitName
					break
				}
			}
			if mainUnit == "" {
				mainUnit = units[0].UnitName
			}
		}

		resp.Products = append(resp.Products, types.ProductInfo{
			ID:          product.ID,
			Name:        product.Name,
			Code:        product.Code,
			Spec:        product.Spec,
			CategoryID:  product.CategoryID,
			CostPrice:   product.CostPrice,
			Price:       product.Price,
			MinStock:    product.MinStock,
			MaxStock:    product.MaxStock,
			SafetyStock: product.SafetyStock,
			Status:      product.Status,
			Desc:        product.Desc,
			MainUnit:    mainUnit,
			Units:       units,
			Category:    category,
			CreatedAt:   product.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   product.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return resp, nil
}

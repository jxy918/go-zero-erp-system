package logic

import (
	"context"
	"time"

	"myproject/admin/internal/model"
	"myproject/admin/internal/svc"
	"myproject/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetErpTopProductsLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewGetErpTopProductsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetErpTopProductsLogic {
	return &GetErpTopProductsLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *GetErpTopProductsLogic) GetErpTopProducts() (resp []*types.ErpTopProductResponse, err error) {
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	type topProduct struct {
		ProductName string  `json:"product_name"`
		Quantity    int     `json:"quantity"`
		Amount      float64 `json:"amount"`
	}

	var products []topProduct
	err = model.DB.Table("sales_order_items soi").
		Select("p.name as product_name, SUM(soi.quantity) as quantity, SUM(soi.quantity * soi.unit_price) as amount").
		Joins("LEFT JOIN sales_orders so ON soi.order_id = so.id").
		Joins("LEFT JOIN products p ON soi.product_id = p.id").
		Where("so.created_at >= ?", startOfMonth).
		Group("soi.product_id, p.name").
		Order("quantity DESC").
		Limit(10).
		Find(&products).Error
	if err != nil {
		return nil, err
	}

	resp = make([]*types.ErpTopProductResponse, 0, len(products))
	for i, product := range products {
		resp = append(resp, &types.ErpTopProductResponse{
			Rank:        i + 1,
			ProductName: product.ProductName,
			Quantity:    product.Quantity,
			Amount:      product.Amount,
		})
	}

	return resp, nil
}
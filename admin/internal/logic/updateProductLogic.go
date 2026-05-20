package logic

import (
	"context"
	"errors"

	"myproject/admin/internal/svc"
	"myproject/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateProductLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewUpdateProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateProductLogic {
	return &UpdateProductLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *UpdateProductLogic) UpdateProduct(req *types.UpdateProductRequest) (resp *types.ProductInfo, err error) {
	product, err := l.svcCtx.ProductModel.GetByID(req.ID)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, errors.New("产品不存在")
	}

	if req.Code != "" && req.Code != product.Code {
		existing, err := l.svcCtx.ProductModel.GetByCode(req.Code)
		if err != nil {
			return nil, err
		}
		if existing != nil && existing.ID != req.ID {
			return nil, errors.New("产品编码已存在")
		}
		product.Code = req.Code
	}

	if req.Name != "" {
		product.Name = req.Name
	}
	product.CategoryID = req.CategoryID
	if req.Spec != "" {
		product.Spec = req.Spec
	}
	if req.Price > 0 {
		product.Price = req.Price
	}
	if req.CostPrice >= 0 {
		product.CostPrice = req.CostPrice
	}
	if req.MinStock >= 0 {
		product.MinStock = req.MinStock
	}
	if req.SafetyStock >= 0 {
		product.SafetyStock = req.SafetyStock
	}
	if req.MaxStock >= 0 {
		product.MaxStock = req.MaxStock
	}
	if req.Desc != "" {
		product.Desc = req.Desc
	}
	if req.Status >= 0 {
		product.Status = req.Status
	}

	err = l.svcCtx.ProductModel.Update(product)
	if err != nil {
		return nil, err
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

	return resp, nil
}

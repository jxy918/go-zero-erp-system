package logic

import (
	"context"
	"errors"

	"myproject/admin/internal/model"
	"myproject/admin/internal/svc"
	"myproject/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateProductLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewCreateProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateProductLogic {
	return &CreateProductLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *CreateProductLogic) CreateProduct(req *types.CreateProductRequest) (resp *types.ProductInfo, err error) {
	if req.Name == "" {
		return nil, errors.New("产品名称不能为空")
	}

	if req.Code == "" {
		return nil, errors.New("产品编码不能为空")
	}

	if req.Price <= 0 {
		return nil, errors.New("产品价格必须大于0")
	}

	existing, err := l.svcCtx.ProductModel.GetByCode(req.Code)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("产品编码已存在")
	}

	product := &model.Product{
		Name:        req.Name,
		Code:        req.Code,
		CategoryID:  req.CategoryID,
		Spec:        req.Spec,
		Price:       req.Price,
		CostPrice:   req.CostPrice,
		MinStock:    req.MinStock,
		SafetyStock: req.SafetyStock,
		MaxStock:    req.MaxStock,
		Desc:        req.Desc,
		Status:      req.Status,
	}

	if product.Status == 0 {
		product.Status = 1
	}

	err = l.svcCtx.ProductModel.Create(product)
	if err != nil {
		return nil, err
	}

	var mainUnit string
	if req.Unit != "" {
		productUnit := &model.ProductUnit{
			ProductID: product.ID,
			UnitName:  req.Unit,
			Ratio:     1,
			IsMain:    1,
		}
		err = l.svcCtx.DB.Create(productUnit).Error
		if err != nil {
			l.Error("创建产品单位失败:", err)
		}
		mainUnit = req.Unit
	}

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
		Stock:       0,
		MainUnit:    mainUnit,
		Units:       []types.ProductUnitInfo{},
		CreatedAt:   product.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   product.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	return resp, nil
}

package logic

import (
	"context"
	"errors"

	"myproject/admin/internal/svc"
	"myproject/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteProductLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewDeleteProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteProductLogic {
	return &DeleteProductLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *DeleteProductLogic) DeleteProduct(req *types.DeleteProductRequest) (resp *types.ProductInfo, err error) {
	product, err := l.svcCtx.ProductModel.GetByID(req.ID)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, errors.New("产品不存在")
	}

	err = l.svcCtx.ProductModel.Delete(req.ID)
	if err != nil {
		return nil, err
	}

	resp = &types.ProductInfo{
		ID:         product.ID,
		Name:       product.Name,
		Code:       product.Code,
		CategoryID: product.CategoryID,
		Price:      product.Price,
		Stock:      0,
		MainUnit:   "",
		Desc:       product.Desc,
		Status:     product.Status,
		CreatedAt:  product.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:  product.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	return resp, nil
}

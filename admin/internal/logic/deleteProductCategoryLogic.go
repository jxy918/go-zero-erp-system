package logic

import (
	"context"
	"errors"

	"myproject/admin/internal/svc"
	"myproject/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteProductCategoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteProductCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteProductCategoryLogic {
	return &DeleteProductCategoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteProductCategoryLogic) DeleteProductCategory(req *types.DeleteProductCategoryRequest) (resp *types.EmptyResponse, err error) {
	category, err := l.svcCtx.CategoryModel.GetByID(req.ID)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, errors.New("分类不存在")
	}

	err = l.svcCtx.CategoryModel.Delete(req.ID)
	if err != nil {
		return nil, err
	}

	return &types.EmptyResponse{}, nil
}
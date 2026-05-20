package logic

import (
	"context"
	"errors"

	"myproject/admin/internal/model"
	"myproject/admin/internal/svc"
	"myproject/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateProductCategoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateProductCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateProductCategoryLogic {
	return &CreateProductCategoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateProductCategoryLogic) CreateProductCategory(req *types.CreateProductCategoryRequest) (resp *types.ProductCategoryInfo, err error) {
	if req.Name == "" {
		return nil, errors.New("分类名称不能为空")
	}
	if req.Code == "" {
		return nil, errors.New("分类编码不能为空")
	}

	existing, err := l.svcCtx.CategoryModel.GetByName(req.Name)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("分类名称已存在")
	}

	existingCode, err := l.svcCtx.CategoryModel.GetByCode(req.Code)
	if err != nil {
		return nil, err
	}
	if existingCode != nil {
		return nil, errors.New("分类编码已存在")
	}

	category := &model.Category{
		Name:     req.Name,
		Code:     req.Code,
		ParentID: req.ParentID,
		Sort:     req.Sort,
		Status:   req.Status,
		Desc:     req.Desc,
	}

	if category.Status == 0 {
		category.Status = 1
	}

	err = l.svcCtx.CategoryModel.Create(category)
	if err != nil {
		return nil, err
	}

	resp = &types.ProductCategoryInfo{
		ID:        category.ID,
		Name:      category.Name,
		Code:      category.Code,
		ParentID:  category.ParentID,
		Sort:      category.Sort,
		Status:    category.Status,
		Desc:      category.Desc,
		Children:  []types.ProductCategoryInfo{},
		CreatedAt: category.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: category.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	return resp, nil
}

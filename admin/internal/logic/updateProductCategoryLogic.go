package logic

import (
	"context"
	"errors"

	"myproject/admin/internal/svc"
	"myproject/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateProductCategoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateProductCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateProductCategoryLogic {
	return &UpdateProductCategoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateProductCategoryLogic) UpdateProductCategory(req *types.UpdateProductCategoryRequest) (resp *types.ProductCategoryInfo, err error) {
	category, err := l.svcCtx.CategoryModel.GetByID(req.ID)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, errors.New("分类不存在")
	}

	if req.Name != "" && req.Name != category.Name {
		existing, err := l.svcCtx.CategoryModel.GetByName(req.Name)
		if err != nil {
			return nil, err
		}
		if existing != nil && existing.ID != req.ID {
			return nil, errors.New("分类名称已存在")
		}
		category.Name = req.Name
	}

	if req.Code != category.Code {
		if req.Code != "" {
			existing, err := l.svcCtx.CategoryModel.GetByCode(req.Code)
			if err != nil {
				return nil, err
			}
			if existing != nil && existing.ID != req.ID {
				return nil, errors.New("分类编码已存在")
			}
		}
		category.Code = req.Code
	}

	if req.ParentID >= 0 {
		category.ParentID = req.ParentID
	}
	if req.Sort >= 0 {
		category.Sort = req.Sort
	}
	if req.Status >= 0 {
		category.Status = req.Status
	}
	category.Desc = req.Desc

	err = l.svcCtx.CategoryModel.Update(category)
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

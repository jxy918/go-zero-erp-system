package logic

import (
	"context"

	"myproject/admin/internal/model"
	"myproject/admin/internal/svc"
	"myproject/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetProductCategoryTreeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetProductCategoryTreeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProductCategoryTreeLogic {
	return &GetProductCategoryTreeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetProductCategoryTreeLogic) GetProductCategoryTree() (resp *types.ListProductCategoryResponse, err error) {
	categories, err := l.svcCtx.CategoryModel.ListAll()
	if err != nil {
		return nil, err
	}

	var categoryInfos []types.ProductCategoryInfo
	for _, cat := range categories {
		if cat.ParentID == 0 {
			info := convertToCategoryInfo(&cat)
			info.Children = buildCategoryChildren(categories, cat.ID)
			categoryInfos = append(categoryInfos, *info)
		}
	}

	resp = &types.ListProductCategoryResponse{
		Categories: categoryInfos,
		Total:      int64(len(categoryInfos)),
	}

	return resp, nil
}

func buildCategoryChildren(categories []model.Category, parentID uint) []types.ProductCategoryInfo {
	var children []types.ProductCategoryInfo
	for _, cat := range categories {
		if cat.ParentID == parentID {
			info := convertToCategoryInfo(&cat)
			info.Children = buildCategoryChildren(categories, cat.ID)
			children = append(children, *info)
		}
	}
	return children
}

func convertToCategoryInfo(cat *model.Category) *types.ProductCategoryInfo {
	return &types.ProductCategoryInfo{
		ID:        cat.ID,
		Name:      cat.Name,
		Code:      cat.Code,
		ParentID:  cat.ParentID,
		Sort:      cat.Sort,
		Status:    cat.Status,
		Desc:      cat.Desc,
		Children:  []types.ProductCategoryInfo{},
		CreatedAt: cat.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: cat.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

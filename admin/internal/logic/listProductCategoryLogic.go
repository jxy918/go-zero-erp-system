package logic

import (
	"context"

	"myproject/admin/internal/svc"
	"myproject/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListProductCategoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListProductCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListProductCategoryLogic {
	return &ListProductCategoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListProductCategoryLogic) ListProductCategory(req *types.ListProductCategoryRequest) (resp *types.ListProductCategoryResponse, err error) {
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 100
	}

	name := ""
	if req.Name != "" {
		name = req.Name
	}

	categories, total, err := l.svcCtx.CategoryModel.List(page, pageSize, name)
	if err != nil {
		return nil, err
	}

	var categoryInfos []types.ProductCategoryInfo
	for _, cat := range categories {
		categoryInfos = append(categoryInfos, types.ProductCategoryInfo{
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
		})
	}

	resp = &types.ListProductCategoryResponse{
		Categories: categoryInfos,
		Total:      total,
	}

	return resp, nil
}

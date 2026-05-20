// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"myproject/admin/internal/svc"
	"myproject/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListActiveSupplierLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListActiveSupplierLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListActiveSupplierLogic {
	return &ListActiveSupplierLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListActiveSupplierLogic) ListActiveSupplier(req *types.ListSupplierRequest) (resp *types.ListSupplierResponse, err error) {
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	suppliers, total, err := l.svcCtx.SupplierModel.ListActive(page, pageSize, req.Name, req.Name)
	if err != nil {
		return nil, err
	}

	resp = &types.ListSupplierResponse{
		Suppliers: make([]types.SupplierInfo, 0, len(suppliers)),
		Total:     total,
	}

	for _, supplier := range suppliers {
		resp.Suppliers = append(resp.Suppliers, types.SupplierInfo{
			ID:        supplier.ID,
			Name:      supplier.Name,
			Code:      supplier.Code,
			Contact:   supplier.Contact,
			Phone:     supplier.Phone,
			Desc:      supplier.Desc,
			Email:     supplier.Email,
			Address:   supplier.Address,
			Status:    supplier.Status,
			CreatedAt: supplier.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: supplier.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return resp, nil
}

// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"myproject/admin/internal/svc"
	"myproject/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListActiveCustomerLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListActiveCustomerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListActiveCustomerLogic {
	return &ListActiveCustomerLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListActiveCustomerLogic) ListActiveCustomer(req *types.ListCustomerRequest) (resp *types.ListCustomerResponse, err error) {
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	customers, total, err := l.svcCtx.CustomerModel.ListActive(page, pageSize, req.Name, req.Name)
	if err != nil {
		return nil, err
	}

	resp = &types.ListCustomerResponse{
		Customers: make([]types.CustomerInfo, 0, len(customers)),
		Total:     total,
	}

	for _, customer := range customers {
		resp.Customers = append(resp.Customers, types.CustomerInfo{
			ID:        customer.ID,
			Name:      customer.Name,
			Code:      customer.Code,
			Contact:   customer.Contact,
			Phone:     customer.Phone,
			Email:     customer.Email,
			Address:   customer.Address,
			Status:    customer.Status,
			Desc:      customer.Desc,
			CreatedAt: customer.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: customer.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return resp, nil
}

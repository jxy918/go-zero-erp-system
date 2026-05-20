package logic

import (
	"context"

	"myproject/admin/internal/svc"
	"myproject/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListCustomerLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListCustomerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListCustomerLogic {
	return &ListCustomerLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListCustomerLogic) ListCustomer(req *types.ListCustomerRequest) (resp *types.ListCustomerResponse, err error) {
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	customers, total, err := l.svcCtx.CustomerModel.List(page, pageSize, req.Name, req.Name)
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
			Desc:      customer.Desc,
			Email:     customer.Email,
			Address:   customer.Address,
			Status:    customer.Status,
			CreatedAt: customer.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: customer.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return resp, nil
}

package logic

import (
	"context"
	"errors"

	"myproject/admin/internal/svc"
	"myproject/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateCustomerLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateCustomerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateCustomerLogic {
	return &UpdateCustomerLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateCustomerLogic) UpdateCustomer(req *types.UpdateCustomerRequest) (resp *types.CustomerInfo, err error) {
	customer, err := l.svcCtx.CustomerModel.GetByID(req.ID)
	if err != nil {
		return nil, err
	}
	if customer == nil {
		return nil, errors.New("客户不存在")
	}

	if req.Name != "" {
		customer.Name = req.Name
	}
	if req.Code != "" {
		customer.Code = req.Code
	}
	if req.Contact != "" {
		customer.Contact = req.Contact
	}
	if req.Phone != "" {
		customer.Phone = req.Phone
	}
	if req.Desc != "" {
		customer.Desc = req.Desc
	}
	if req.Email != "" {
		customer.Email = req.Email
	}
	if req.Address != "" {
		customer.Address = req.Address
	}
	if req.Status >= 0 {
		customer.Status = req.Status
	}

	if err := l.svcCtx.CustomerModel.Update(customer); err != nil {
		return nil, err
	}

	resp = &types.CustomerInfo{
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
	}

	return resp, nil
}

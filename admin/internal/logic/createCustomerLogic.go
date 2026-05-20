package logic

import (
	"context"
	"errors"

	"myproject/admin/internal/model"
	"myproject/admin/internal/svc"
	"myproject/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateCustomerLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewCreateCustomerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateCustomerLogic {
	return &CreateCustomerLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *CreateCustomerLogic) CreateCustomer(req *types.CreateCustomerRequest) (resp *types.CustomerInfo, err error) {
	if req.Name == "" {
		return nil, errors.New("客户名称不能为空")
	}

	existing, err := l.svcCtx.CustomerModel.GetByName(req.Name)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("客户名称已存在")
	}

	customer := &model.Customer{
		Name:    req.Name,
		Code:    req.Code,
		Contact: req.Contact,
		Phone:   req.Phone,
		Desc:    req.Desc,
		Email:   req.Email,
		Address: req.Address,
		Status:  req.Status,
	}

	if customer.Status == 0 {
		customer.Status = 1
	}

	err = l.svcCtx.CustomerModel.Create(customer)
	if err != nil {
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

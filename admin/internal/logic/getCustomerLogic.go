// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"errors"

	"myproject/admin/internal/svc"
	"myproject/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCustomerLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCustomerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCustomerLogic {
	return &GetCustomerLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCustomerLogic) GetCustomer(req *types.GetCustomerRequest) (resp *types.CustomerInfo, err error) {
	customer, err := l.svcCtx.CustomerModel.GetByID(req.ID)
	if err != nil {
		return nil, err
	}
	if customer == nil {
		return nil, errors.New("客户不存在")
	}

	resp = &types.CustomerInfo{
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
	}

	return resp, nil
}

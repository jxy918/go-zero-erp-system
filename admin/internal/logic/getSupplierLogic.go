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

type GetSupplierLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetSupplierLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSupplierLogic {
	return &GetSupplierLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSupplierLogic) GetSupplier(req *types.GetSupplierRequest) (resp *types.SupplierInfo, err error) {
	supplier, err := l.svcCtx.SupplierModel.GetByID(req.ID)
	if err != nil {
		return nil, err
	}
	if supplier == nil {
		return nil, errors.New("供应商不存在")
	}

	resp = &types.SupplierInfo{
		ID:        supplier.ID,
		Name:      supplier.Name,
		Code:      supplier.Code,
		Contact:   supplier.Contact,
		Phone:     supplier.Phone,
		Email:     supplier.Email,
		Address:   supplier.Address,
		Status:    supplier.Status,
		Desc:      supplier.Desc,
		CreatedAt: supplier.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: supplier.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	return resp, nil
}

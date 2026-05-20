package logic

import (
	"context"
	"errors"

	"myproject/admin/internal/svc"
	"myproject/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateSupplierLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateSupplierLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateSupplierLogic {
	return &UpdateSupplierLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateSupplierLogic) UpdateSupplier(req *types.UpdateSupplierRequest) (resp *types.SupplierInfo, err error) {
	supplier, err := l.svcCtx.SupplierModel.GetByID(req.ID)
	if err != nil {
		return nil, err
	}
	if supplier == nil {
		return nil, errors.New("供应商不存在")
	}

	if req.Name != "" {
		supplier.Name = req.Name
	}
	if req.Code != "" {
		supplier.Code = req.Code
	}
	if req.Contact != "" {
		supplier.Contact = req.Contact
	}
	if req.Phone != "" {
		supplier.Phone = req.Phone
	}
	if req.Desc != "" {
		supplier.Desc = req.Desc
	}
	if req.Email != "" {
		supplier.Email = req.Email
	}
	if req.Address != "" {
		supplier.Address = req.Address
	}
	if req.Status >= 0 {
		supplier.Status = req.Status
	}

	if err := l.svcCtx.SupplierModel.Update(supplier); err != nil {
		return nil, err
	}

	resp = &types.SupplierInfo{
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
	}

	return resp, nil
}

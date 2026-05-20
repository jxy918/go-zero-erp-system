package logic

import (
	"context"
	"errors"

	"myproject/admin/internal/model"
	"myproject/admin/internal/svc"
	"myproject/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateSupplierLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewCreateSupplierLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateSupplierLogic {
	return &CreateSupplierLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *CreateSupplierLogic) CreateSupplier(req *types.CreateSupplierRequest) (resp *types.SupplierInfo, err error) {
	if req.Name == "" {
		return nil, errors.New("供应商名称不能为空")
	}

	existing, err := l.svcCtx.SupplierModel.GetByName(req.Name)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("供应商名称已存在")
	}

	supplier := &model.Supplier{
		Name:    req.Name,
		Code:    req.Code,
		Contact: req.Contact,
		Phone:   req.Phone,
		Desc:    req.Desc,
		Email:   req.Email,
		Address: req.Address,
		Status:  req.Status,
	}

	if supplier.Status == 0 {
		supplier.Status = 1
	}

	err = l.svcCtx.SupplierModel.Create(supplier)
	if err != nil {
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
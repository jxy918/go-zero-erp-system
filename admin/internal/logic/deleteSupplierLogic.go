package logic

import (
	"context"
	"errors"

	"myproject/admin/internal/svc"
	"myproject/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteSupplierLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteSupplierLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteSupplierLogic {
	return &DeleteSupplierLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteSupplierLogic) DeleteSupplier(req *types.DeleteSupplierRequest) (resp *types.SupplierInfo, err error) {
	supplier, err := l.svcCtx.SupplierModel.GetByID(req.ID)
	if err != nil {
		return nil, err
	}
	if supplier == nil {
		return nil, errors.New("供应商不存在")
	}

	if err := l.svcCtx.SupplierModel.Delete(req.ID); err != nil {
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

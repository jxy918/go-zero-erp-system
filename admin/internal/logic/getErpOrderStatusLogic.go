package logic

import (
	"context"

	"myproject/admin/internal/model"
	"myproject/admin/internal/svc"
	"myproject/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetErpOrderStatusLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewGetErpOrderStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetErpOrderStatusLogic {
	return &GetErpOrderStatusLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *GetErpOrderStatusLogic) GetErpOrderStatus(req *types.ErpOrderStatusRequest) (resp *types.ErpOrderStatusResponse, err error) {
	orderType := req.Type
	if orderType == "" {
		orderType = "purchase"
	}

	resp = &types.ErpOrderStatusResponse{}

	switch orderType {
	case "purchase":
		var pending, approved, completed, cancelled int64
		err = model.DB.Model(&model.PurchaseOrder{}).Where("status = ?", 1).Count(&pending).Error
		if err != nil {
			return nil, err
		}
		err = model.DB.Model(&model.PurchaseOrder{}).Where("status = ?", 2).Count(&approved).Error
		if err != nil {
			return nil, err
		}
		err = model.DB.Model(&model.PurchaseOrder{}).Where("status = ?", 3).Count(&completed).Error
		if err != nil {
			return nil, err
		}
		err = model.DB.Model(&model.PurchaseOrder{}).Where("status = ?", 4).Count(&cancelled).Error
		if err != nil {
			return nil, err
		}
		resp.Pending = int(pending)
		resp.Approved = int(approved)
		resp.Completed = int(completed)
		resp.Cancelled = int(cancelled)
	case "sales":
		var pending, approved, completed, cancelled int64
		err = model.DB.Model(&model.SalesOrder{}).Where("status = ?", 1).Count(&pending).Error
		if err != nil {
			return nil, err
		}
		err = model.DB.Model(&model.SalesOrder{}).Where("status = ?", 2).Count(&approved).Error
		if err != nil {
			return nil, err
		}
		err = model.DB.Model(&model.SalesOrder{}).Where("status = ?", 3).Count(&completed).Error
		if err != nil {
			return nil, err
		}
		err = model.DB.Model(&model.SalesOrder{}).Where("status = ?", 4).Count(&cancelled).Error
		if err != nil {
			return nil, err
		}
		resp.Pending = int(pending)
		resp.Approved = int(approved)
		resp.Completed = int(completed)
		resp.Cancelled = int(cancelled)
	default:
		var pending, approved, completed, cancelled int64
		err = model.DB.Model(&model.PurchaseOrder{}).Where("status = ?", 1).Count(&pending).Error
		if err != nil {
			return nil, err
		}
		err = model.DB.Model(&model.PurchaseOrder{}).Where("status = ?", 2).Count(&approved).Error
		if err != nil {
			return nil, err
		}
		err = model.DB.Model(&model.PurchaseOrder{}).Where("status = ?", 3).Count(&completed).Error
		if err != nil {
			return nil, err
		}
		err = model.DB.Model(&model.PurchaseOrder{}).Where("status = ?", 4).Count(&cancelled).Error
		if err != nil {
			return nil, err
		}
		resp.Pending = int(pending)
		resp.Approved = int(approved)
		resp.Completed = int(completed)
		resp.Cancelled = int(cancelled)
	}

	return resp, nil
}

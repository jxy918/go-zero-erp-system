package logic

import (
	"context"

	"myproject/admin/internal/svc"
	"myproject/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListPurchaseOrderLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewListPurchaseOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListPurchaseOrderLogic {
	return &ListPurchaseOrderLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *ListPurchaseOrderLogic) ListPurchaseOrder(req *types.ListPurchaseOrderRequest) (resp *types.ListPurchaseOrderResponse, err error) {
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	orders, total, err := l.svcCtx.PurchaseModel.List(page, pageSize, req.OrderNo, req.Status)
	if err != nil {
		return nil, err
	}

	resp = &types.ListPurchaseOrderResponse{
		Orders: make([]types.PurchaseOrderInfo, 0, len(orders)),
		Total:  total,
	}

	for _, order := range orders {
		item := types.PurchaseOrderInfo{
			ID:          order.ID,
			OrderNo:     order.OrderNo,
			SupplierID:  order.SupplierID,
			WarehouseID: order.WarehouseID,
			TotalAmount: order.TotalAmount,
			Status:      order.Status,
			Remark:      order.Remark,
			CreatedAt:   order.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   order.UpdatedAt.Format("2006-01-02 15:04:05"),
		}

		if order.Supplier.ID > 0 {
			item.Supplier = types.SupplierInfo{
				ID:        order.Supplier.ID,
				Name:      order.Supplier.Name,
				Contact:   order.Supplier.Contact,
				Phone:     order.Supplier.Phone,
				Email:     order.Supplier.Email,
				Address:   order.Supplier.Address,
				Status:    order.Supplier.Status,
				CreatedAt: order.Supplier.CreatedAt.Format("2006-01-02 15:04:05"),
				UpdatedAt: order.Supplier.UpdatedAt.Format("2006-01-02 15:04:05"),
			}
		}

		if order.Warehouse.ID > 0 {
			item.Warehouse = types.WarehouseInfo{
				ID:        order.Warehouse.ID,
				Name:      order.Warehouse.Name,
				Code:      order.Warehouse.Code,
				Contact:   order.Warehouse.Contact,
				Desc:      order.Warehouse.Desc,
				Status:    order.Warehouse.Status,
				CreatedAt: order.Warehouse.CreatedAt.Format("2006-01-02 15:04:05"),
				UpdatedAt: order.Warehouse.UpdatedAt.Format("2006-01-02 15:04:05"),
			}
		}

		resp.Orders = append(resp.Orders, item)
	}

	return resp, nil
}

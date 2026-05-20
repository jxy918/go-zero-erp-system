package logic

import (
	"context"

	"myproject/admin/internal/svc"
	"myproject/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListSalesOrderLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewListSalesOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListSalesOrderLogic {
	return &ListSalesOrderLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *ListSalesOrderLogic) ListSalesOrder(req *types.ListSalesOrderRequest) (resp *types.ListSalesOrderResponse, err error) {
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	orders, total, err := l.svcCtx.SalesModel.List(page, pageSize, req.OrderNo, req.Status)
	if err != nil {
		return nil, err
	}

	resp = &types.ListSalesOrderResponse{
		Orders: make([]types.SalesOrderInfo, 0, len(orders)),
		Total:  total,
	}

	for _, order := range orders {
		item := types.SalesOrderInfo{
			ID:          order.ID,
			OrderNo:     order.OrderNo,
			CustomerID:  order.CustomerID,
			WarehouseID: order.WarehouseID,
			TotalAmount: order.TotalAmount,
			Status:      order.Status,
			Remark:      order.Remark,
			CreatedAt:   order.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   order.UpdatedAt.Format("2006-01-02 15:04:05"),
		}

		if order.Customer.ID > 0 {
			item.Customer = types.CustomerInfo{
				ID:        order.Customer.ID,
				Name:      order.Customer.Name,
				Contact:   order.Customer.Contact,
				Phone:     order.Customer.Phone,
				Email:     order.Customer.Email,
				Address:   order.Customer.Address,
				Status:    order.Customer.Status,
				CreatedAt: order.Customer.CreatedAt.Format("2006-01-02 15:04:05"),
				UpdatedAt: order.Customer.UpdatedAt.Format("2006-01-02 15:04:05"),
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

package logic

import (
	"context"
	"errors"

	"myproject/admin/internal/svc"
	"myproject/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSalesOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetSalesOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSalesOrderLogic {
	return &GetSalesOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSalesOrderLogic) GetSalesOrder(req *types.GetSalesOrderRequest) (resp *types.SalesOrderInfo, err error) {
	order, err := l.svcCtx.SalesModel.GetByID(req.ID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, errors.New("订单不存在")
	}

	resp = &types.SalesOrderInfo{
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
		resp.Customer = types.CustomerInfo{
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
		resp.Warehouse = types.WarehouseInfo{
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

	resp.Items = make([]types.SalesOrderItemInfo, 0, len(order.Items))
	for _, item := range order.Items {
		productInfo := types.ProductInfo{}
		if item.Product.ID > 0 {
			stock, _ := l.svcCtx.InventoryModel.GetProductStockByWarehouse(item.ProductID, order.WarehouseID)
			productInfo = types.ProductInfo{
				ID:        item.Product.ID,
				Name:      item.Product.Name,
				Code:      item.Product.Code,
				Price:     item.Product.Price,
				Stock:     stock,
				MainUnit:  "",
				Status:    item.Product.Status,
				CreatedAt: item.Product.CreatedAt.Format("2006-01-02 15:04:05"),
				UpdatedAt: item.Product.UpdatedAt.Format("2006-01-02 15:04:05"),
			}
		}

		resp.Items = append(resp.Items, types.SalesOrderItemInfo{
			ID:        item.ID,
			OrderID:   item.OrderID,
			ProductID: item.ProductID,
			Product:   productInfo,
			UnitID:    item.UnitID,
			UnitName:  item.UnitName,
			Ratio:     item.Ratio,
			Quantity:  item.Quantity,
			BaseQty:   item.BaseQty,
			UnitPrice: item.UnitPrice,
			Amount:    item.Amount,
		})
	}

	return resp, nil
}

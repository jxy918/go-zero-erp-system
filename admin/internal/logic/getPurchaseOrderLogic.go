package logic

import (
	"context"
	"errors"

	"myproject/admin/internal/svc"
	"myproject/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPurchaseOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPurchaseOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPurchaseOrderLogic {
	return &GetPurchaseOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPurchaseOrderLogic) GetPurchaseOrder(req *types.GetPurchaseOrderRequest) (resp *types.PurchaseOrderInfo, err error) {
	order, err := l.svcCtx.PurchaseModel.GetByID(req.ID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, errors.New("订单不存在")
	}

	resp = &types.PurchaseOrderInfo{
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
		resp.Supplier = types.SupplierInfo{
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

	resp.Items = make([]types.PurchaseOrderItemInfo, 0, len(order.Items))
	for _, item := range order.Items {
		productInfo := types.ProductInfo{}
		if item.Product.ID > 0 {
			stock, _ := l.svcCtx.InventoryModel.GetProductStockByWarehouse(item.ProductID, order.WarehouseID)
			productInfo = types.ProductInfo{
				ID:         item.Product.ID,
				Name:       item.Product.Name,
				Code:       item.Product.Code,
				CategoryID: item.Product.CategoryID,
				Spec:       item.Product.Spec,
				CostPrice:  item.Product.CostPrice,
				Stock:      stock,
				Status:     item.Product.Status,
				CreatedAt:  item.Product.CreatedAt.Format("2006-01-02 15:04:05"),
				UpdatedAt:  item.Product.UpdatedAt.Format("2006-01-02 15:04:05"),
			}
		}

		resp.Items = append(resp.Items, types.PurchaseOrderItemInfo{
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

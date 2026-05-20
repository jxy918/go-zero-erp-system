package logic

import (
	"context"
	"errors"

	"myproject/admin/internal/middleware"
	"myproject/admin/internal/model"
	"myproject/admin/internal/svc"
	"myproject/admin/internal/types"
	"myproject/admin/internal/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdatePurchaseOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdatePurchaseOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatePurchaseOrderLogic {
	return &UpdatePurchaseOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdatePurchaseOrderLogic) UpdatePurchaseOrder(req *types.UpdatePurchaseOrderRequest) (resp *types.PurchaseOrderInfo, err error) {
	order, err := l.svcCtx.PurchaseModel.GetByID(req.ID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, errors.New("订单不存在")
	}

	if req.Status != 0 && req.Status != order.Status {
		beforeStatus := order.Status
		if err := l.svcCtx.PurchaseModel.UpdateStatus(req.ID, req.Status); err != nil {
			return nil, err
		}

		operatorID, _ := l.ctx.Value(util.UserIDKey).(uint)
		operatorName := middleware.GetUsername(l.ctx)

		log := &model.OrderLog{
			OrderID:      req.ID,
			OrderType:    1,
			BeforeStatus: beforeStatus,
			AfterStatus:  req.Status,
			OperatorID:   operatorID,
			OperatorName: operatorName,
			Remark:       "采购订单状态变更",
		}
		if err := l.svcCtx.OrderLogModel.Insert(log); err != nil {
			return nil, err
		}
	}

	updatedOrder, err := l.svcCtx.PurchaseModel.GetByID(req.ID)
	if err != nil {
		return nil, err
	}

	resp = &types.PurchaseOrderInfo{
		ID:          updatedOrder.ID,
		OrderNo:     updatedOrder.OrderNo,
		SupplierID:  updatedOrder.SupplierID,
		WarehouseID: updatedOrder.WarehouseID,
		TotalAmount: updatedOrder.TotalAmount,
		Status:      updatedOrder.Status,
		Remark:      updatedOrder.Remark,
		CreatedAt:   updatedOrder.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   updatedOrder.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	if updatedOrder.Supplier.ID > 0 {
		resp.Supplier = types.SupplierInfo{
			ID:        updatedOrder.Supplier.ID,
			Name:      updatedOrder.Supplier.Name,
			Contact:   updatedOrder.Supplier.Contact,
			Phone:     updatedOrder.Supplier.Phone,
			Email:     updatedOrder.Supplier.Email,
			Address:   updatedOrder.Supplier.Address,
			Status:    updatedOrder.Supplier.Status,
			CreatedAt: updatedOrder.Supplier.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: updatedOrder.Supplier.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	if updatedOrder.Warehouse.ID > 0 {
		resp.Warehouse = types.WarehouseInfo{
			ID:        updatedOrder.Warehouse.ID,
			Name:      updatedOrder.Warehouse.Name,
			Code:      updatedOrder.Warehouse.Code,
			Contact:   updatedOrder.Warehouse.Contact,
			Desc:      updatedOrder.Warehouse.Desc,
			Status:    updatedOrder.Warehouse.Status,
			CreatedAt: updatedOrder.Warehouse.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: updatedOrder.Warehouse.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	resp.Items = make([]types.PurchaseOrderItemInfo, 0, len(updatedOrder.Items))
	for _, item := range updatedOrder.Items {
		productInfo := types.ProductInfo{}
		if item.Product.ID > 0 {
			stock, _ := l.svcCtx.InventoryModel.GetProductStockByWarehouse(item.ProductID, updatedOrder.WarehouseID)
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
			Quantity:  item.Quantity,
			UnitPrice: item.UnitPrice,
			Amount:    item.Amount,
		})
	}

	return resp, nil
}

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

type CreateSalesOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateSalesOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateSalesOrderLogic {
	return &CreateSalesOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateSalesOrderLogic) CreateSalesOrder(req *types.CreateSalesOrderRequest) (resp *types.SalesOrderInfo, err error) {
	if req.CustomerID == 0 {
		return nil, errors.New("客户不能为空")
	}

	if len(req.Items) == 0 {
		return nil, errors.New("销售订单不能为空")
	}

	customer, err := l.svcCtx.CustomerModel.GetByID(req.CustomerID)
	if err != nil {
		return nil, err
	}
	if customer == nil {
		return nil, errors.New("客户不存在")
	}

	order := &model.SalesOrder{
		CustomerID:  req.CustomerID,
		WarehouseID: req.WarehouseID,
		Remark:      req.Remark,
		Status:      1,
	}

	if req.WarehouseID == 0 {
		return nil, errors.New("出库仓库不能为空")
	}

	warehouse, err := l.svcCtx.WarehouseModel.GetByID(req.WarehouseID)
	if err != nil {
		return nil, err
	}
	if warehouse == nil {
		return nil, errors.New("仓库不存在")
	}

	var totalAmount float64
	items := make([]model.SalesOrderItem, 0, len(req.Items))
	for _, item := range req.Items {
		product, err := l.svcCtx.ProductModel.GetByID(item.ProductID)
		if err != nil {
			return nil, err
		}
		if product == nil {
			return nil, errors.New("产品不存在")
		}

		ratio := int(item.Ratio)
		if ratio <= 0 {
			ratio = 1
		}

		unitID := item.UnitID
		unitName := item.UnitName
		if unitID == 0 {
			var units []model.ProductUnit
			err := l.svcCtx.DB.Where("product_id = ? AND is_main = 1", item.ProductID).First(&units).Error
			if err != nil {
				err = l.svcCtx.DB.Where("product_id = ?", item.ProductID).First(&units).Error
			}
			if err == nil && len(units) > 0 {
				unitID = units[0].ID
				unitName = units[0].UnitName
				ratio = int(units[0].Ratio)
			}
		}

		baseQty := item.Quantity * ratio

		stock, err := l.svcCtx.InventoryModel.GetStockByProductAndWarehouse(item.ProductID, req.WarehouseID)
		if err != nil {
			return nil, err
		}
		if stock < baseQty {
			return nil, errors.New("仓库库存不足")
		}

		amount := float64(item.Quantity) * item.UnitPrice
		totalAmount += amount

		items = append(items, model.SalesOrderItem{
			ProductID: item.ProductID,
			UnitID:    unitID,
			UnitName:  unitName,
			Ratio:     ratio,
			Quantity:  item.Quantity,
			BaseQty:   baseQty,
			UnitPrice: item.UnitPrice,
			Amount:    amount,
		})
	}

	order.TotalAmount = totalAmount

	err = l.svcCtx.SalesModel.Create(order, items)
	if err != nil {
		return nil, err
	}

	operatorID, _ := l.ctx.Value(util.UserIDKey).(uint)
	operatorName := middleware.GetUsername(l.ctx)

	log := &model.OrderLog{
		OrderID:      order.ID,
		OrderType:    2,
		BeforeStatus: 0,
		AfterStatus:  1,
		OperatorID:   operatorID,
		OperatorName: operatorName,
		Remark:       "创建销售订单",
	}
	if err := l.svcCtx.OrderLogModel.Insert(log); err != nil {
		return nil, err
	}

	resp = buildSalesOrderInfo(l.svcCtx.InventoryModel, order, customer, warehouse, items)

	return resp, nil
}

func buildSalesOrderInfo(inventoryModel *model.InventoryModel, order *model.SalesOrder, customer *model.Customer, warehouse *model.Warehouse, items []model.SalesOrderItem) *types.SalesOrderInfo {
	resp := &types.SalesOrderInfo{
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

	if customer != nil {
		resp.Customer = types.CustomerInfo{
			ID:        customer.ID,
			Name:      customer.Name,
			Contact:   customer.Contact,
			Phone:     customer.Phone,
			Email:     customer.Email,
			Address:   customer.Address,
			Status:    customer.Status,
			CreatedAt: customer.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: customer.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	if warehouse != nil {
		resp.Warehouse = types.WarehouseInfo{
			ID:        warehouse.ID,
			Name:      warehouse.Name,
			Code:      warehouse.Code,
			Contact:   warehouse.Contact,
			Desc:      warehouse.Desc,
			Status:    warehouse.Status,
			CreatedAt: warehouse.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: warehouse.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	resp.Items = make([]types.SalesOrderItemInfo, 0, len(items))
	for _, item := range items {
		var product types.ProductInfo
		if item.Product.ID > 0 {
			stock, _ := inventoryModel.GetProductStockByWarehouse(item.ProductID, order.WarehouseID)
			product = types.ProductInfo{
				ID:         item.Product.ID,
				Name:       item.Product.Name,
				Code:       item.Product.Code,
				CategoryID: item.Product.CategoryID,
				Price:      item.Product.Price,
				Stock:      stock,
				Status:     item.Product.Status,
			}
		}

		resp.Items = append(resp.Items, types.SalesOrderItemInfo{
			ID:        item.ID,
			OrderID:   item.OrderID,
			ProductID: item.ProductID,
			Product:   product,
			UnitID:    item.UnitID,
			UnitName:  item.UnitName,
			Ratio:     item.Ratio,
			Quantity:  item.Quantity,
			BaseQty:   item.BaseQty,
			UnitPrice: item.UnitPrice,
			Amount:    item.Amount,
		})
	}

	return resp
}

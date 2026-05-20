package logic

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/logx"

	"myproject/admin/internal/model"
	"myproject/admin/internal/types"
)

// TestCreatePurchaseOrder tests the CreatePurchaseOrder function
func TestCreatePurchaseOrder(t *testing.T) {
	logx.Disable() // Disable logging for tests

	tests := []struct {
		name    string
		req     *types.CreatePurchaseOrderRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "供应商为空",
			req: &types.CreatePurchaseOrderRequest{
				SupplierID: 0,
				Items: []types.CreatePurchaseOrderItem{
					{ProductID: 1, Quantity: 10, UnitPrice: 100.00},
				},
			},
			wantErr: true,
			errMsg:  "供应商不能为空",
		},
		{
			name: "订单商品为空",
			req: &types.CreatePurchaseOrderRequest{
				SupplierID: 1,
				Items:      []types.CreatePurchaseOrderItem{},
			},
			wantErr: true,
			errMsg:  "采购订单不能为空",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create logic with nil svcCtx - tests basic validation
			l := &CreatePurchaseOrderLogic{
				Logger: logx.WithContext(context.Background()),
				svcCtx: nil, // We're only testing validation logic that doesn't need DB
			}

			_, err := l.CreatePurchaseOrder(tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestCreateSalesOrder tests the CreateSalesOrder function
func TestCreateSalesOrder(t *testing.T) {
	logx.Disable() // Disable logging for tests

	tests := []struct {
		name    string
		req     *types.CreateSalesOrderRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "客户为空",
			req: &types.CreateSalesOrderRequest{
				CustomerID: 0,
				Items: []types.SalesOrderItem{
					{ProductID: 1, Quantity: 10, UnitPrice: 100.00},
				},
			},
			wantErr: true,
			errMsg:  "客户不能为空",
		},
		{
			name: "订单商品为空",
			req: &types.CreateSalesOrderRequest{
				CustomerID: 1,
				Items:      []types.SalesOrderItem{},
			},
			wantErr: true,
			errMsg:  "销售订单不能为空",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &CreateSalesOrderLogic{
				Logger: logx.WithContext(context.Background()),
				svcCtx: nil, // We're only testing validation logic that doesn't need DB
			}

			_, err := l.CreateSalesOrder(tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestAdjustInventory tests the AdjustInventory function
func TestAdjustInventory(t *testing.T) {
	logx.Disable() // Disable logging for tests

	tests := []struct {
		name    string
		req     *types.AdjustInventoryRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "产品为空",
			req: &types.AdjustInventoryRequest{
				ProductID:   0,
				WarehouseID: 1,
				Quantity:    50,
			},
			wantErr: true,
			errMsg:  "产品不能为空",
		},
		{
			name: "仓库为空",
			req: &types.AdjustInventoryRequest{
				ProductID:   1,
				WarehouseID: 0,
				Quantity:    50,
			},
			wantErr: true,
			errMsg:  "仓库不能为空",
		},
		{
			name: "调整数量为0",
			req: &types.AdjustInventoryRequest{
				ProductID:   1,
				WarehouseID: 1,
				Quantity:    0,
			},
			wantErr: true,
			errMsg:  "调整数量不能为0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &AdjustInventoryLogic{
				Logger: logx.WithContext(context.Background()),
				svcCtx: nil, // We're only testing validation logic that doesn't need DB
			}

			_, err := l.AdjustInventory(tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestBuildPurchaseOrderInfo tests the helper function
func TestBuildPurchaseOrderInfo(t *testing.T) {
	order := &model.PurchaseOrder{
		ID:          1,
		OrderNo:     "PO202604280001",
		SupplierID:  1,
		WarehouseID: 1,
		TotalAmount: 1000.00,
		Status:      1,
		Remark:      "Test Order",
	}

	supplier := &model.Supplier{
		ID:      1,
		Name:    "Test Supplier",
		Contact: "Test Contact",
		Phone:   "13800138000",
	}

	warehouse := &model.Warehouse{
		ID:   1,
		Name: "Test Warehouse",
		Code: "WH001",
	}

	items := []model.PurchaseOrderItem{
		{
			ID:        1,
			OrderID:   1,
			ProductID: 1,
			Quantity:  10,
			UnitPrice: 100.00,
			Amount:    1000.00,
		},
	}

	result := buildPurchaseOrderInfo(nil, order, supplier, warehouse, items)

	assert.Equal(t, uint(1), result.ID)
	assert.Equal(t, "PO202604280001", result.OrderNo)
	assert.Equal(t, float64(1000.00), result.TotalAmount)
	assert.Equal(t, "Test Supplier", result.Supplier.Name)
	assert.Equal(t, "Test Warehouse", result.Warehouse.Name)
	assert.Len(t, result.Items, 1)
}

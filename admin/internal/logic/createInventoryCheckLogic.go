// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"fmt"
	"time"

	"myproject/admin/internal/model"
	"myproject/admin/internal/svc"
	"myproject/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type CreateInventoryCheckLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateInventoryCheckLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateInventoryCheckLogic {
	return &CreateInventoryCheckLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func generateCheckNo() string {
	return fmt.Sprintf("CK-%s-%04d", time.Now().Format("20060102"), time.Now().UnixMilli()%10000)
}

func (l *CreateInventoryCheckLogic) CreateInventoryCheck(req *types.CreateInventoryCheckRequest) (resp *types.InventoryCheckResponse, err error) {
	warehouse, err := l.svcCtx.WarehouseModel.GetByID(req.WarehouseID)
	if err != nil {
		return nil, err
	}
	if warehouse == nil {
		return &types.InventoryCheckResponse{}, nil
	}

	warehouseInventories, err := l.svcCtx.InventoryModel.GetWarehouseInventory(req.WarehouseID)
	if err != nil {
		return nil, err
	}

	if len(warehouseInventories) == 0 {
		return &types.InventoryCheckResponse{}, nil
	}

	check := &model.InventoryCheck{
		CheckNo:     generateCheckNo(),
		WarehouseID: req.WarehouseID,
		Status:      1,
		Remark:      req.Remark,
	}

	err = model.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(check).Error; err != nil {
			return err
		}

		var items []model.InventoryCheckItem
		for _, inv := range warehouseInventories {
			items = append(items, model.InventoryCheckItem{
				CheckID:   check.ID,
				ProductID: inv.ProductID,
				SystemQty: inv.Quantity,
				ActualQty: 0,
				DiffQty:   -inv.Quantity,
				Status:    1,
			})
		}

		if len(items) > 0 {
			if err := tx.Create(&items).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	check, err = l.svcCtx.InventoryCheckModel.GetByID(check.ID)
	if err != nil {
		return nil, err
	}

	return &types.InventoryCheckResponse{
		Check: *convertInventoryCheckToInfo(l.svcCtx.InventoryModel, check),
	}, nil
}

func convertInventoryCheckToInfo(inventoryModel *model.InventoryModel, check *model.InventoryCheck) *types.InventoryCheckInfo {
	if check == nil {
		return nil
	}

	items := make([]types.InventoryCheckItemInfo, 0)
	totalDiff := 0
	for _, item := range check.Items {
		totalDiff += item.DiffQty
		stock, _ := inventoryModel.GetProductStockByWarehouse(item.ProductID, check.WarehouseID)
		items = append(items, types.InventoryCheckItemInfo{
			ID:        item.ID,
			CheckID:   item.CheckID,
			ProductID: item.ProductID,
			Product:   convertProductInfo(&item.Product, stock),
			SystemQty: item.SystemQty,
			ActualQty: item.ActualQty,
			DiffQty:   item.DiffQty,
			Status:    item.Status,
		})
	}

	return &types.InventoryCheckInfo{
		ID:          check.ID,
		CheckNo:     check.CheckNo,
		WarehouseID: check.WarehouseID,
		Warehouse:   convertWarehouseInfo(&check.Warehouse),
		Status:      check.Status,
		TotalDiff:   totalDiff,
		Remark:      check.Remark,
		CreatedAt:   check.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   check.UpdatedAt.Format("2006-01-02 15:04:05"),
		Items:       items,
	}
}

func convertProductInfo(product *model.Product, stock int) types.ProductInfo {
	if product == nil {
		return types.ProductInfo{}
	}
	var mainUnit string
	for _, unit := range product.Units {
		if unit.IsMain == 1 {
			mainUnit = unit.UnitName
			break
		}
	}
	return types.ProductInfo{
		ID:         product.ID,
		Name:       product.Name,
		Code:       product.Code,
		CategoryID: product.CategoryID,
		Spec:       product.Spec,
		MainUnit:   mainUnit,
		Price:      product.Price,
		CostPrice:  product.CostPrice,
		Stock:      stock,
		Desc:       product.Desc,
		Status:     product.Status,
		CreatedAt:  product.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:  product.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func convertWarehouseInfo(warehouse *model.Warehouse) types.WarehouseInfo {
	if warehouse == nil {
		return types.WarehouseInfo{}
	}
	return types.WarehouseInfo{
		ID:      warehouse.ID,
		Name:    warehouse.Name,
		Code:    warehouse.Code,
		Address: warehouse.Address,
		Status:  warehouse.Status,
	}
}

func convertInventoryTransferToInfo(inventoryModel *model.InventoryModel, transfer *model.InventoryTransfer) *types.InventoryTransferInfo {
	if transfer == nil {
		return nil
	}
	stock, _ := inventoryModel.GetProductStockByWarehouse(transfer.ProductID, transfer.ToWarehouseID)
	info := &types.InventoryTransferInfo{
		ID:              transfer.ID,
		TransferNo:      transfer.TransferNo,
		FromWarehouseID: transfer.FromWarehouseID,
		FromWarehouse:   convertWarehouseInfo(&transfer.FromWarehouse),
		ToWarehouseID:   transfer.ToWarehouseID,
		ToWarehouse:     convertWarehouseInfo(&transfer.ToWarehouse),
		ProductID:       transfer.ProductID,
		Product:         convertProductInfo(&transfer.Product, stock),
		Quantity:        transfer.Quantity,
		Status:          transfer.Status,
		Remark:          transfer.Remark,
		CreatedBy:       transfer.CreatedBy,
		CreatedAt:       transfer.CreatedAt.Format("2006-01-02 15:04:05"),
		AuditedBy:       transfer.AuditedBy,
		ExecutedBy:      transfer.ExecutedBy,
		UpdatedAt:       transfer.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
	if transfer.AuditedAt != nil {
		info.AuditedAt = transfer.AuditedAt.Format("2006-01-02 15:04:05")
	}
	if transfer.ExecutedAt != nil {
		info.ExecutedAt = transfer.ExecutedAt.Format("2006-01-02 15:04:05")
	}
	if transfer.CreatedBy > 0 {
		var creator model.User
		if err := model.DB.First(&creator, transfer.CreatedBy).Error; err == nil {
			info.CreatedByName = creator.Nickname
			if info.CreatedByName == "" {
				info.CreatedByName = creator.Username
			}
		}
	}
	if transfer.AuditedBy > 0 {
		var auditor model.User
		if err := model.DB.First(&auditor, transfer.AuditedBy).Error; err == nil {
			info.AuditedByName = auditor.Nickname
			if info.AuditedByName == "" {
				info.AuditedByName = auditor.Username
			}
		}
	}
	if transfer.ExecutedBy > 0 {
		var executor model.User
		if err := model.DB.First(&executor, transfer.ExecutedBy).Error; err == nil {
			info.ExecutedByName = executor.Nickname
			if info.ExecutedByName == "" {
				info.ExecutedByName = executor.Username
			}
		}
	}
	return info
}

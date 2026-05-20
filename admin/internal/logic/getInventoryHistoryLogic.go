package logic

import (
	"context"

	"myproject/admin/internal/svc"
	"myproject/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetInventoryHistoryLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewGetInventoryHistoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetInventoryHistoryLogic {
	return &GetInventoryHistoryLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *GetInventoryHistoryLogic) GetInventoryHistory(req *types.GetInventoryHistoryRequest) (resp *types.ListInventoryRecordResponse, err error) {
	records, err := l.svcCtx.InventoryChangeModel.GetHistory(req.ProductID, req.WarehouseID)
	if err != nil {
		return nil, err
	}

	var recordInfos []types.InventoryRecordInfo
	for _, record := range records {
		info := types.InventoryRecordInfo{
			ID:             record.ID,
			ProductID:      record.ProductID,
			WarehouseID:    record.WarehouseID,
			Type:           record.Type,
			Quantity:       record.Quantity,
			Balance:        record.AfterQuantity,
			OrderNo:        "",
			Remark:         record.Remark,
			CreatedAt:      record.CreatedAt.Format("2006-01-02 15:04:05"),
			BeforeQuantity: record.BeforeQuantity,
			AfterQuantity:  record.AfterQuantity,
		}

		product, err := l.svcCtx.ProductModel.GetByID(record.ProductID)
		if err == nil && product != nil {
			info.Product = product.Name
			info.ProductCode = product.Code
			unit, err := l.svcCtx.ProductUnitModel.GetMainUnit(record.ProductID)
			if err == nil && unit != nil {
				info.Unit = unit.UnitName
			}
		}

		warehouse, _ := l.svcCtx.WarehouseModel.GetByID(record.WarehouseID)
		if warehouse != nil {
			info.Warehouse = warehouse.Name
		}

		recordInfos = append(recordInfos, info)
	}

	return &types.ListInventoryRecordResponse{
		Records: recordInfos,
		Total:   int64(len(recordInfos)),
	}, nil
}

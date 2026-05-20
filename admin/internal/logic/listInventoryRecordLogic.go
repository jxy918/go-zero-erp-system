package logic

import (
	"context"

	"myproject/admin/internal/svc"
	"myproject/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListInventoryRecordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListInventoryRecordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListInventoryRecordLogic {
	return &ListInventoryRecordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListInventoryRecordLogic) ListInventoryRecord(req *types.ListInventoryRecordRequest) (*types.ListInventoryRecordResponse, error) {
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}

	records, total, err := l.svcCtx.InventoryModel.List(page, pageSize, "", req.ProductID, req.WarehouseID)
	if err != nil {
		return nil, err
	}

	var recordInfos []types.InventoryRecordInfo
	for _, r := range records {
		info := types.InventoryRecordInfo{
			ID:          r.ID,
			ProductID:   r.ProductID,
			WarehouseID: r.WarehouseID,
			Type:        r.Type,
			Quantity:    r.Quantity,
			CreatedAt:   r.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		if r.Product.ID > 0 {
			info.Product = r.Product.Name
			info.ProductCode = r.Product.Code
			unit, err := l.svcCtx.ProductUnitModel.GetMainUnit(r.ProductID)
			if err == nil && unit != nil {
				info.Unit = unit.UnitName
			}
		}
		if r.Warehouse.ID > 0 {
			info.Warehouse = r.Warehouse.Name
		}
		recordInfos = append(recordInfos, info)
	}

	return &types.ListInventoryRecordResponse{
		Records: recordInfos,
		Total:   total,
	}, nil
}

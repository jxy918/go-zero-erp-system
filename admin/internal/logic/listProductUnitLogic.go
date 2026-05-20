package logic

import (
	"context"

	"myproject/admin/internal/model"
	"myproject/admin/internal/svc"
	"myproject/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListProductUnitLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListProductUnitLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListProductUnitLogic {
	return &ListProductUnitLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListProductUnitLogic) ListProductUnit(req *types.ListProductUnitRequest) (resp *types.ListProductUnitResponse, err error) {
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}

	var productUnits []model.ProductUnit
	var total int64

	query := l.svcCtx.DB.Model(&model.ProductUnit{}).
		Joins(" LEFT JOIN products p ON product_units.product_id = p.id").
		Where("p.status = 1 OR p.status IS NULL")

	if req.ProductID > 0 {
		query = query.Where("product_units.product_id = ?", req.ProductID)
	}

	if req.ProductName != "" {
		query = query.Where("p.name LIKE ?", "%"+req.ProductName+"%")
	}

	err = query.Count(&total).Error
	if err != nil {
		return nil, err
	}

	err = query.Preload("Product").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&productUnits).Error
	if err != nil {
		return nil, err
	}

	units := make([]types.ProductUnitInfo, 0, len(productUnits))
	for _, pu := range productUnits {
		productName := ""
		if pu.Product.ID > 0 {
			productName = pu.Product.Name
		}
		units = append(units, types.ProductUnitInfo{
			ID:          pu.ID,
			ProductID:   pu.ProductID,
			ProductName: productName,
			UnitName:    pu.UnitName,
			Ratio:       pu.Ratio,
			IsMain:      pu.IsMain,
			CreatedAt:   pu.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   pu.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &types.ListProductUnitResponse{
		Units: units,
		Total: total,
	}, nil
}

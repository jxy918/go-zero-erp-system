// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"errors"

	"myproject/admin/internal/model"
	"myproject/admin/internal/svc"
	"myproject/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type GetProductUnitLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetProductUnitLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProductUnitLogic {
	return &GetProductUnitLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetProductUnitLogic) GetProductUnit(req *types.GetProductUnitRequest) (resp *types.ProductUnitInfo, err error) {
	var pu model.ProductUnit
	err = l.svcCtx.DB.Preload("Product").First(&pu, req.ID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("产品单位不存在")
		}
		return nil, err
	}

	productName := ""
	if pu.Product.ID > 0 {
		productName = pu.Product.Name
	}

	resp = &types.ProductUnitInfo{
		ID:          pu.ID,
		ProductID:   pu.ProductID,
		ProductName: productName,
		UnitName:    pu.UnitName,
		Ratio:       pu.Ratio,
		IsMain:      pu.IsMain,
		CreatedAt:   pu.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   pu.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	return resp, nil
}

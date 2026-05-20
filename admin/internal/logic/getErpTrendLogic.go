package logic

import (
	"context"
	"time"

	"myproject/admin/internal/model"
	"myproject/admin/internal/svc"
	"myproject/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetErpTrendLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewGetErpTrendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetErpTrendLogic {
	return &GetErpTrendLogic{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *GetErpTrendLogic) GetErpTrend(req *types.ErpTrendRequest) (resp *types.ErpTrendResponse, err error) {
	days := req.Days
	if days <= 0 {
		days = 30
	}

	now := time.Now()
	dates := make([]string, days)
	purchaseData := make([]float64, days)
	salesData := make([]float64, days)

	for i := 0; i < days; i++ {
		date := now.AddDate(0, 0, -(days - 1 - i))
		dateStr := date.Format("2006-01-02")
		dates[i] = dateStr

		var purchaseAmount float64
		err = model.DB.Model(&model.PurchaseOrder{}).
			Where("DATE(created_at) = ?", dateStr).
			Select("COALESCE(SUM(total_amount), 0)").
			Find(&purchaseAmount).Error
		if err != nil {
			return nil, err
		}
		purchaseData[i] = purchaseAmount

		var salesAmount float64
		err = model.DB.Model(&model.SalesOrder{}).
			Where("DATE(created_at) = ?", dateStr).
			Select("COALESCE(SUM(total_amount), 0)").
			Find(&salesAmount).Error
		if err != nil {
			return nil, err
		}
		salesData[i] = salesAmount
	}

	resp = &types.ErpTrendResponse{
		Dates:        dates,
		PurchaseData: purchaseData,
		SalesData:    salesData,
	}

	return resp, nil
}

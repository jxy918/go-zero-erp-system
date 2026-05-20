package logic

import (
	"context"

	"myproject/admin/internal/model"
	"myproject/admin/internal/svc"
	"myproject/admin/internal/types"
	"myproject/admin/internal/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListActivityLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListActivityLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListActivityLogic {
	return &ListActivityLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListActivityLogic) ListActivity(req *types.ListActivityRequest) (resp *types.ListActivityResponse, err error) {
	// 获取当前用户权限信息
	isAdmin := util.IsAdmin(l.ctx)
	userID := util.GetUserID(l.ctx)

	var activities []model.Activity
	var total int64

	// 根据权限获取活动日志列表
	if isAdmin {
		// 超级管理员：查看所有日志（支持按用户名搜索）
		username := req.Username
		if username != "" {
			activities, total, err = l.svcCtx.ActivityModel.ListByUsername(username, req.Page, req.PageSize)
		} else {
			activities, total, err = l.svcCtx.ActivityModel.List(req.Page, req.PageSize)
		}
	} else {
		// 普通用户：只能查看自己的日志（从会话获取用户ID，不接受前端参数）
		activities, total, err = l.svcCtx.ActivityModel.ListByUserID(userID, req.Page, req.PageSize)
	}

	if err != nil {
		return nil, err
	}

	activityInfos := make([]types.ActivityInfo, len(activities))
	for i, activity := range activities {
		activityInfos[i] = types.ActivityInfo{
			ID:        activity.ID,
			UserID:    activity.UserID,
			Username:  activity.Username,
			Action:    activity.Action,
			URL:       activity.URL,
			IP:        activity.IP,
			CreatedAt: activity.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	return &types.ListActivityResponse{
		Activities: activityInfos,
		Total:      total,
	}, nil
}

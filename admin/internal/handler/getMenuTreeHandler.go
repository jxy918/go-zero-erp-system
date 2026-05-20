// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package handler

import (
	"net/http"

	"myproject/admin/internal/logic"
	"myproject/admin/internal/svc"
	"myproject/admin/internal/types"
	"myproject/admin/internal/util"

	"github.com/zeromicro/go-zero/core/logx"
)

func GetMenuTreeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		logx.Info("GetMenuTreeHandler - 请求到达")

		req := types.ListMenuRequest{
			Page:     1,
			PageSize: 100,
		}

		l := logic.NewGetMenuTreeLogic(r.Context(), svcCtx)
		resp, err := l.GetMenuTree(req)
		if err != nil {
			logx.Errorf("GetMenuTreeHandler - 处理失败: %v", err)
			util.ErrorResponse(w, r, 500, err.Error())
		} else {
			logx.Infof("GetMenuTreeHandler - 处理成功，返回 %d 个菜单", len(resp.Menus))
			util.SuccessResponse(w, r, resp)
		}
	}
}

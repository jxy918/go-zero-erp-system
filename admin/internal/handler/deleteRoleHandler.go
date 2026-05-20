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
	"github.com/zeromicro/go-zero/rest/httpx"
)

func DeleteRoleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		logx.Info("DeleteRoleHandler: Request received")
		logx.Info("DeleteRoleHandler: Request URL:", r.URL.String())
		logx.Info("DeleteRoleHandler: Request method:", r.Method)

		var req types.DeleteRoleRequest
		if err := httpx.Parse(r, &req); err != nil {
			logx.Error("DeleteRoleHandler: Parse error:", err)
			util.ErrorResponse(w, r, 400, "field \"id\" is not set")
			return
		}

		logx.Info("DeleteRoleHandler: Parsed request:", req)

		l := logic.NewDeleteRoleLogic(r.Context(), svcCtx)
		resp, err := l.DeleteRole(&req)
		if err != nil {
			logx.Error("DeleteRoleHandler: Logic error:", err)
			util.ErrorResponse(w, r, 500, err.Error())
		} else {
			logx.Info("DeleteRoleHandler: Success")
			util.SuccessResponse(w, r, resp)
		}
	}
}

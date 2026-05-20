// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"myproject/admin/internal/logic"
	"myproject/admin/internal/svc"
	"myproject/admin/internal/types"
	"myproject/admin/internal/util"
)

func AssignMenuPermissionsHandler(serverCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		
		var req types.AssignMenuPermissionsRequest
		if err := httpx.Parse(r, &req); err != nil {
			util.ErrorResponse(w, r, 400, err.Error())
			return
		}

		l := logic.NewAssignMenuPermissionsLogic(r.Context(), serverCtx)
		resp, err := l.AssignMenuPermissions(req)
		if err != nil {
			util.ErrorResponse(w, r, 500, err.Error())
		} else {
			util.SuccessResponse(w, r, resp)
		}
	}
}
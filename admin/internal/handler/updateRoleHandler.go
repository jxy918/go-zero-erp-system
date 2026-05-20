// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package handler

import (
	"encoding/json"
	"net/http"

	"myproject/admin/internal/logic"
	"myproject/admin/internal/svc"
	"myproject/admin/internal/types"
	"myproject/admin/internal/util"
)

func UpdateRoleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		
		var req types.UpdateRoleRequest
		
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			util.ErrorResponse(w, r, 400, "解析请求失败: "+err.Error())
			return
		}

		l := logic.NewUpdateRoleLogic(r.Context(), svcCtx)
		resp, err := l.UpdateRole(&req)
		if err != nil {
			util.ErrorResponse(w, r, 500, err.Error())
		} else {
			util.SuccessResponse(w, r, resp)
		}
	}
}
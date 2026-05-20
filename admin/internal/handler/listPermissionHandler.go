// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package handler

import (
	"net/http"
	"strconv"

	"myproject/admin/internal/logic"
	"myproject/admin/internal/svc"
	"myproject/admin/internal/types"
	"myproject/admin/internal/util"
)

func ListPermissionHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		
		req := types.ListPermissionRequest{
			Page:     1,
			PageSize: 10,
			Name:     r.URL.Query().Get("name"),
		}

		if err := r.ParseForm(); err == nil {
			if page := r.Form.Get("page"); page != "" {
				req.Page, _ = strconv.Atoi(page)
			}
			if pageSize := r.Form.Get("page_size"); pageSize != "" {
				req.PageSize, _ = strconv.Atoi(pageSize)
			}
		}

		l := logic.NewListPermissionLogic(r.Context(), svcCtx)
		resp, err := l.ListPermission(&req)
		if err != nil {
			util.ErrorResponse(w, r, 500, err.Error())
		} else {
			util.SuccessResponse(w, r, resp)
		}
	}
}
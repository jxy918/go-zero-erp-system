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

func ListUserHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page := 1
		pageSize := 10
		username := ""

		if pageStr := r.URL.Query().Get("page"); pageStr != "" {
			if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
				page = p
			}
		}

		if pageSizeStr := r.URL.Query().Get("page_size"); pageSizeStr != "" {
			if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 {
				pageSize = ps
			}
		}

		username = r.URL.Query().Get("username")

		req := types.ListUserRequest{
			Page:     page,
			PageSize: pageSize,
			Username: username,
		}

		l := logic.NewListUserLogic(r.Context(), svcCtx)
		resp, err := l.ListUser(&req)
		if err != nil {
			util.ErrorResponse(w, r, 500, err.Error())
		} else {
			util.SuccessResponse(w, r, resp)
		}
	}
}

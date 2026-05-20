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
	"github.com/zeromicro/go-zero/core/logx"
)

func ListActivityHandler(serverCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		
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

		if usernameStr := r.URL.Query().Get("username"); usernameStr != "" {
			username = usernameStr
		}

		req := &types.ListActivityRequest{
			Page:     page,
			PageSize: pageSize,
			Username: username,
		}

		logx.Infof("ListActivityRequest: Page=%d, PageSize=%d, Username=%s", req.Page, req.PageSize, req.Username)

		l := logic.NewListActivityLogic(r.Context(), serverCtx)
		resp, err := l.ListActivity(req)
		if err != nil {
			util.ErrorResponse(w, r, 500, err.Error())
		} else {
			util.SuccessResponse(w, r, resp)
		}
	}
}
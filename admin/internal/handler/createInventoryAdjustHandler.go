// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package handler

import (
	"net/http"

	"myproject/admin/internal/logic"
	"myproject/admin/internal/svc"
	"myproject/admin/internal/types"
	"myproject/admin/internal/util"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func CreateInventoryAdjustHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateInventoryAdjustRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewCreateInventoryAdjustLogic(r.Context(), svcCtx)
		resp, err := l.CreateInventoryAdjust(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			util.SuccessResponse(w, r, resp)
		}
	}
}

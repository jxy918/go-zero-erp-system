// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package handler

import (
	"encoding/json"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"myproject/admin/internal/logic"
	"myproject/admin/internal/svc"
	"myproject/admin/internal/types"
	"myproject/admin/internal/util"
)

func AuditInventoryTransferHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AuditInventoryTransferRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewAuditInventoryTransferLogic(r.Context(), svcCtx)
		resp, err := l.AuditInventoryTransfer(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			util.SuccessResponse(w, r, resp)
		}
	}
}

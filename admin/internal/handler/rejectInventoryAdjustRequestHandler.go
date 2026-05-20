package handler

import (
	"net/http"

	"myproject/admin/internal/logic"
	"myproject/admin/internal/svc"
	"myproject/admin/internal/types"
	"myproject/admin/internal/util"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func RejectInventoryAdjustRequestHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RejectInventoryAdjustRequest
		if err := httpx.Parse(r, &req); err != nil {
			util.ErrorResponse(w, r, 400, err.Error())
			return
		}

		l := logic.NewRejectInventoryAdjustRequestLogic(r.Context(), svcCtx)
		resp, err := l.RejectInventoryAdjustRequest(&req)
		if err != nil {
			util.ErrorResponse(w, r, 500, err.Error())
		} else {
			util.SuccessResponse(w, r, resp)
		}
	}
}

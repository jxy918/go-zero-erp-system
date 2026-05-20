package handler

import (
	"net/http"

	"myproject/admin/internal/logic"
	"myproject/admin/internal/svc"
	"myproject/admin/internal/util"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetErpTodoDataHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewGetErpTodoDataLogic(r.Context(), svcCtx)
		resp, err := l.GetErpTodoData()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			util.SuccessResponse(w, r, resp)
		}
	}
}

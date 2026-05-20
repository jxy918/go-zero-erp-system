// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package handler

import (
	"net/http"

	"myproject/admin/internal/logic"
	"myproject/admin/internal/svc"
	"myproject/admin/internal/util"
)

func InitDataHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewInitDataLogic(r.Context(), svcCtx)
		resp, err := l.InitData()
		if err != nil {
			util.ErrorResponse(w, r, 500, err.Error())
		} else {
			util.SuccessResponse(w, r, resp)
		}
	}
}

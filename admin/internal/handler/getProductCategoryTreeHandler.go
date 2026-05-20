// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package handler

import (
	"net/http"

	"myproject/admin/internal/logic"
	"myproject/admin/internal/svc"
	"myproject/admin/internal/util"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetProductCategoryTreeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewGetProductCategoryTreeLogic(r.Context(), svcCtx)
		resp, err := l.GetProductCategoryTree()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			util.SuccessResponse(w, r, resp)
		}
	}
}

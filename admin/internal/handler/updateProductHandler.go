// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package handler

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"myproject/admin/internal/logic"
	"myproject/admin/internal/svc"
	"myproject/admin/internal/types"
	"myproject/admin/internal/util"
)

func UpdateProductHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			util.ErrorResponse(w, r, 400, "Failed to read request body")
			return
		}

		r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		var req types.UpdateProductRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			util.ErrorResponse(w, r, 400, err.Error())
			return
		}

		l := logic.NewUpdateProductLogic(r.Context(), svcCtx)
		resp, err := l.UpdateProduct(&req)
		if err != nil {
			util.ErrorResponse(w, r, 500, err.Error())
		} else {
			util.SuccessResponse(w, r, resp)
		}
	}
}

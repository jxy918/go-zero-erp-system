// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package handler

import (
	"fmt"
	"net/http"

	"myproject/admin/internal/logic"
	"myproject/admin/internal/svc"
	"myproject/admin/internal/types"
	"myproject/admin/internal/util"
	"github.com/zeromicro/go-zero/core/logx"
)

func DeleteUserHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		logx.Info("DeleteUserHandler: Request received")
		logx.Info("DeleteUserHandler: Request URL:", r.URL.String())
		logx.Info("DeleteUserHandler: Request method:", r.Method)
		logx.Info("DeleteUserHandler: Query parameters:", r.URL.Query())

		var req types.DeleteUserRequest
		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			logx.Error("DeleteUserHandler: id parameter is required")
			util.ErrorResponse(w, r, 400, "field \"id\" is not set")
			return
		}

		var id uint
		_, err := fmt.Sscanf(idStr, "%d", &id)
		if err != nil {
			logx.Error("DeleteUserHandler: Parse error:", err)
			util.ErrorResponse(w, r, 400, "invalid id parameter")
			return
		}

		req.ID = id

		logx.Info("DeleteUserHandler: Parsed request:", req)

		l := logic.NewDeleteUserLogic(r.Context(), svcCtx)
		resp, err := l.DeleteUser(&req)
		if err != nil {
			logx.Error("DeleteUserHandler: Logic error:", err)
			util.ErrorResponse(w, r, 500, err.Error())
		} else {
			logx.Info("DeleteUserHandler: Success")
			util.SuccessResponse(w, r, resp)
		}
	}
}
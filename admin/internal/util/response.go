package util

import (
	"net/http"

	"myproject/admin/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// SuccessResponse 成功响应
func SuccessResponse(w http.ResponseWriter, r *http.Request, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	response := types.Response{
		Code:    0,
		Data:    data,
		Message: "success",
	}
	httpx.OkJsonCtx(r.Context(), w, response)
}

// ErrorResponse 错误响应
func ErrorResponse(w http.ResponseWriter, r *http.Request, code int, message string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	response := types.Response{
		Code:    code,
		Data:    nil,
		Message: message,
	}
	httpx.OkJsonCtx(r.Context(), w, response)
}

// ParseRequest 解析请求，返回错误时直接写错误响应
func ParseRequest(w http.ResponseWriter, r *http.Request, v interface{}) bool {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := httpx.Parse(r, v); err != nil {
		ErrorResponse(w, r, 400, err.Error())
		return false
	}
	return true
}

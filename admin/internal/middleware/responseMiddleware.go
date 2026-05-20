package middleware

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// responseRecorder 用于捕获响应
type responseRecorder struct {
	http.ResponseWriter
	body       *bytes.Buffer
	statusCode int
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	return r.body.Write(b)
}

func (r *responseRecorder) WriteHeader(code int) {
	r.statusCode = code
}

// ResponseMiddleware 统一响应格式中间件
// 该中间件会统一所有接口的响应格式为 {"code":0,"data":{},"message":"success"}
func ResponseMiddleware() func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// 创建响应记录器
			recorder := &responseRecorder{
				ResponseWriter: w,
				body:           bytes.NewBufferString(""),
				statusCode:     http.StatusOK,
			}

			// 执行下一个中间件/handler
			next(recorder, r)

			// 获取响应内容
			responseBody := recorder.body.String()

			// 检查是否是有效的JSON
			var responseMap map[string]interface{}
			if err := json.Unmarshal([]byte(responseBody), &responseMap); err == nil {
				// 检查是否已经是统一格式的响应（包含 code、data、message 字段）
				if _, hasCode := responseMap["code"]; hasCode {
					if _, hasData := responseMap["data"]; hasData {
						if _, hasMessage := responseMap["message"]; hasMessage {
							// 已经是统一格式，保留原始状态码
							w.Header().Set("Content-Type", "application/json; charset=utf-8")
							w.WriteHeader(recorder.statusCode)
							w.Write([]byte(responseBody))
							return
						}
					}
				}

				// 检查是否是 go-zero 的错误响应格式 {"error":"xxx"}
				if _, exists := responseMap["error"]; exists {
					// 转换为统一响应格式
					w.Header().Set("Content-Type", "application/json; charset=utf-8")
					w.WriteHeader(recorder.statusCode)

					// 获取错误消息
					errorMsg := ""
					if errVal, ok := responseMap["error"]; ok {
						errorMsg = errVal.(string)
					}

					// 返回统一格式
					json.NewEncoder(w).Encode(map[string]interface{}{
						"code":    recorder.statusCode,
						"data":    nil,
						"message": errorMsg,
					})
					return
				}

				// 如果是普通业务数据（不是统一格式也不是错误格式），包装成统一格式
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				w.WriteHeader(recorder.statusCode)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"code":    0,
					"data":    responseMap,
					"message": "success",
				})
				return
			}

			// 如果不是有效的JSON，直接返回原始响应
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(recorder.statusCode)
			w.Write([]byte(responseBody))
		}
	}
}

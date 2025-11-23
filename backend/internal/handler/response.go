// backend/internal/handler/response.go
package handler

import (
	"encoding/json"
	"net/http"
)

// APIResponse は統一されたAPIレスポンス形式
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
}

// APIError はエラー情報
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// WriteJSON はJSONレスポンスを書き込みます
func WriteJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := APIResponse{
		Success: statusCode >= 200 && statusCode < 300,
		Data:    data,
	}

	json.NewEncoder(w).Encode(response)
}

// WriteError はエラーレスポンスを書き込みます
func WriteError(w http.ResponseWriter, statusCode int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := APIResponse{
		Success: false,
		Error: &APIError{
			Code:    code,
			Message: message,
		},
	}

	json.NewEncoder(w).Encode(response)
}

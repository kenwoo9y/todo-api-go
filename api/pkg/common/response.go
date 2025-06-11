package common

import (
	"encoding/json"
	"net/http"
)

// エラーレスポンスの共通構造体
type ErrorResponse struct {
	Message string `json:"message"`
}

// JSONレスポンスを送信する共通関数
func JSONResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// エラーレスポンスを送信する共通関数
func ErrorJSONResponse(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{
		Message: message,
	})
}

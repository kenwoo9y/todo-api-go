package common

import (
	"encoding/json"
	"net/http"
)

// Common error response structure
type ErrorResponse struct {
	Message string `json:"message"`
}

// Common function to send JSON responses
func JSONResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// Common function to send error responses
func ErrorJSONResponse(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{
		Message: message,
	})
}

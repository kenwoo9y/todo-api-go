package common

import (
	"net/http"
	"strconv"
	"strings"
)

// Common function to extract ID from URL path
func ExtractIDFromPath(path string, prefix string) (int64, error) {
	idStr := strings.TrimPrefix(path, prefix)
	return strconv.ParseInt(idStr, 10, 64)
}

// Common function to extract ownerID from /users/{id}/tasks path
func ExtractOwnerIDFromPath(path string) (int64, error) {
	pathParts := strings.Split(path, "/")
	if len(pathParts) != 4 {
		return 0, ErrInvalidPathFormat
	}
	return strconv.ParseInt(pathParts[2], 10, 64)
}

// Common function to validate HTTP methods
func ValidateRequestMethod(w http.ResponseWriter, r *http.Request, allowedMethods ...string) bool {
	for _, method := range allowedMethods {
		if r.Method == method {
			return true
		}
	}
	ErrorJSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
	return false
}

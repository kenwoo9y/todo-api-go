package common

import (
	"net/http"
	"strconv"
	"strings"
)

// URLパスからIDを抽出する共通関数
func ExtractIDFromPath(path string, prefix string) (int64, error) {
	idStr := strings.TrimPrefix(path, prefix)
	return strconv.ParseInt(idStr, 10, 64)
}

// /users/{id}/tasks からownerIDを抽出する共通関数
func ExtractOwnerIDFromPath(path string) (int64, error) {
	pathParts := strings.Split(path, "/")
	if len(pathParts) != 4 {
		return 0, ErrInvalidPathFormat
	}
	return strconv.ParseInt(pathParts[2], 10, 64)
}

// HTTPメソッドを検証する共通関数
func ValidateRequestMethod(w http.ResponseWriter, r *http.Request, allowedMethods ...string) bool {
	for _, method := range allowedMethods {
		if r.Method == method {
			return true
		}
	}
	ErrorJSONResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
	return false
}

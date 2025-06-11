package common

import (
	"errors"
	"net/http"
)

// 共通エラー定義
var (
	ErrInvalidPathFormat = errors.New("invalid path format")
	ErrInvalidID         = errors.New("invalid id")
	ErrInvalidOwnerID    = errors.New("invalid owner id")
	ErrNotFound          = errors.New("not found")
	ErrInternalServer    = errors.New("internal server error")
)

// エラーを適切なHTTPレスポンスに変換する共通関数
func HandleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrInvalidPathFormat),
		errors.Is(err, ErrInvalidID),
		errors.Is(err, ErrInvalidOwnerID):
		ErrorJSONResponse(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, ErrNotFound):
		ErrorJSONResponse(w, http.StatusNotFound, err.Error())
	default:
		ErrorJSONResponse(w, http.StatusInternalServerError, ErrInternalServer.Error())
	}
}

package common

import (
	"errors"
	"net/http"
)

// Common error definitions
var (
	ErrInvalidPathFormat = errors.New("invalid path format")
	ErrInvalidID         = errors.New("invalid id")
	ErrInvalidOwnerID    = errors.New("invalid owner id")
	ErrNotFound          = errors.New("not found")
	ErrInternalServer    = errors.New("internal server error")
)

// Common function to convert errors to appropriate HTTP responses
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

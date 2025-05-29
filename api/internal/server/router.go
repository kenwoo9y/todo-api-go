package server

import (
	"net/http"

	"github.com/kenwoo9y/todo-api-go/api/internal/handler"
)

func SetupServer(userHandler *handler.UserHandler) *http.Server {
	return &http.Server{
		Addr:    ":8080",
		Handler: userHandler,
	}
}

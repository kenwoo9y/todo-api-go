package server

import (
	"net/http"

	"github.com/kenwoo9y/todo-api-go/api/internal/handler"
)

func SetupServer(userHandler *handler.UserHandler, taskHandler *handler.TaskHandler) *http.Server {
	mux := http.NewServeMux()
	mux.Handle("/users", userHandler)
	mux.Handle("/users/", userHandler)
	mux.Handle("/tasks", taskHandler)
	mux.Handle("/tasks/", taskHandler)

	return &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
}

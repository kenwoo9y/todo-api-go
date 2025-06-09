package server

import (
	"net/http"
	"strings"

	"github.com/kenwoo9y/todo-api-go/api/internal/handler"
)

type customRouter struct {
	userHandler *handler.UserHandler
	taskHandler *handler.TaskHandler
}

func (r *customRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path

	switch {
	case path == "/users" || path == "/users/":
		r.userHandler.ServeHTTP(w, req)
	case strings.HasPrefix(path, "/users/username/"):
		r.userHandler.ServeHTTP(w, req)
	case strings.HasPrefix(path, "/users/") && strings.HasSuffix(path, "/tasks"):
		r.taskHandler.ServeHTTP(w, req)
	case strings.HasPrefix(path, "/users/"):
		r.userHandler.ServeHTTP(w, req)
	case path == "/tasks" || path == "/tasks/":
		r.taskHandler.ServeHTTP(w, req)
	case strings.HasPrefix(path, "/tasks/"):
		r.taskHandler.ServeHTTP(w, req)
	default:
		http.Error(w, "Not Found", http.StatusNotFound)
	}
}

func SetupServer(userHandler *handler.UserHandler, taskHandler *handler.TaskHandler) *http.Server {
	router := &customRouter{
		userHandler: userHandler,
		taskHandler: taskHandler,
	}

	return &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
}

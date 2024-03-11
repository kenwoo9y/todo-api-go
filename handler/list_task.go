package handler

import (
	"net/http"
	"time"

	"github.com/kenwoo9y/todo-api-go/entity"
	"github.com/kenwoo9y/todo-api-go/store"
)

type ListTask struct {
	Store *store.TaskStore
}

type task struct {
	ID          entity.TaskID     `json:"id"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Status      entity.TaskStatus `json:"status"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

func (lt *ListTask) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tasks := lt.Store.All()
	rsp := []task{}
	for _, t := range tasks {
		rsp = append(rsp, task{
			ID:          t.ID,
			Title:       t.Title,
			Description: t.Description,
			Status:      t.Status,
			CreatedAt:   t.CreatedAt,
			UpdatedAt:   t.UpdatedAt,
		})
	}
	RespondJSON(ctx, w, rsp, http.StatusOK)
}

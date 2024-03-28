package handler

import (
	"context"

	"github.com/kenwoo9y/todo-api-go/entity"
)

type ListTasksService interface {
	ListTasks(ctx context.Context) (entity.Tasks, error)
}
type AddTaskService interface {
	AddTask(ctx context.Context, title string) (*entity.Task, error)
}

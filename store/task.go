package store

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/kenwoo9y/todo-api-go/entity"
)

func (r *Repository) ListTasks(
	ctx context.Context, db *sqlx.DB,
) (entity.Tasks, error) {
	tasks := entity.Tasks{}
	sql := `SELECT id, title, description, status, created_at, updated_at FROM tasks;`
	if err := db.SelectContext(ctx, &tasks, sql); err != nil {
		return nil, err
	}
	return tasks, nil
}

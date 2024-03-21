package store

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/kenwoo9y/todo-api-go/entity"
)

func (r *Repository) AddTask(
	ctx context.Context, db Execer, t *entity.Task,
) error {
	t.CreatedAt = r.Clocker.Now()
	t.UpdatedAt = r.Clocker.Now()
	sql := `INSERT INTO task
		(title, description, status, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`
	result, err := db.ExecContext(
		ctx, sql, t.Title, t.Status,
		t.CreatedAt, t.UpdatedAt,
	)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	t.ID = entity.TaskID(id)
	return nil
}

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

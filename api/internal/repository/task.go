package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/kenwoo9y/todo-api-go/api/internal/config"
	"github.com/kenwoo9y/todo-api-go/api/internal/entity"
)

type TaskRepository interface {
	Create(ctx context.Context, task *entity.Task) error
	GetAll(ctx context.Context) ([]entity.Task, error)
	GetByID(ctx context.Context, id int64) (*entity.Task, error)
	GetByOwnerID(ctx context.Context, ownerID int64) ([]entity.Task, error)
	Update(ctx context.Context, task *entity.Task) error
	Delete(ctx context.Context, id int64) error
}

type taskRepository struct {
	db     *sql.DB
	dbType string
}

func NewTaskRepository(db *sql.DB, cfg *config.Config) TaskRepository {
	return &taskRepository{
		db:     db,
		dbType: cfg.DBType,
	}
}

func (r *taskRepository) Create(ctx context.Context, task *entity.Task) error {
	var query string
	if r.dbType == "mysql" {
		query = `
			INSERT INTO tasks (title, description, due_date, status, owner_id, created_at, updated_at)
			VALUES (?, ?, STR_TO_DATE(?, '%Y-%m-%d'), ?, ?, ?, ?)`
	} else {
		query = `
			INSERT INTO tasks (title, description, due_date, status, owner_id, created_at, updated_at)
			VALUES ($1, $2, $3::date, $4, $5, $6, $7)
			RETURNING id`
	}

	now := time.Now()
	if r.dbType == "mysql" {
		result, err := r.db.ExecContext(ctx,
			query,
			task.Title,
			task.Description,
			task.DueDate,
			task.Status,
			task.OwnerID,
			now,
			now,
		)
		if err != nil {
			return err
		}

		id, err := result.LastInsertId()
		if err != nil {
			return err
		}
		task.ID = id
		return nil
	} else {
		return r.db.QueryRowContext(ctx,
			query,
			task.Title,
			task.Description,
			task.DueDate,
			task.Status,
			task.OwnerID,
			now,
			now,
		).Scan(&task.ID)
	}
}

func (r *taskRepository) GetAll(ctx context.Context) ([]entity.Task, error) {
	query := `SELECT id, title, description, DATE_FORMAT(due_date, '%Y-%m-%d') as due_date, status, owner_id, created_at, updated_at FROM tasks
		ORDER BY
			CASE status
				WHEN 'Done' THEN 1
				ELSE 0
			END ASC,
			due_date ASC,
			created_at DESC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []entity.Task
	for rows.Next() {
		var task entity.Task
		if err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.DueDate,
			&task.Status,
			&task.OwnerID,
			&task.CreatedAt,
			&task.UpdatedAt,
		); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, rows.Err()
}

func (r *taskRepository) GetByID(ctx context.Context, id int64) (*entity.Task, error) {
	var task entity.Task
	var query string
	if r.dbType == "mysql" {
		query = `SELECT id, title, description, DATE_FORMAT(due_date, '%Y-%m-%d') as due_date, status, owner_id, created_at, updated_at FROM tasks WHERE id = ?`
	} else {
		query = `SELECT id, title, description, TO_CHAR(due_date, 'YYYY-MM-DD') as due_date, status, owner_id, created_at, updated_at FROM tasks WHERE id = $1`
	}

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.DueDate,
		&task.Status,
		&task.OwnerID,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &task, err
}

func (r *taskRepository) GetByOwnerID(ctx context.Context, ownerID int64) ([]entity.Task, error) {
	baseQuery := `
		SELECT id, title, description, DATE_FORMAT(due_date, '%%Y-%%m-%%d') as due_date, status, owner_id, created_at, updated_at FROM tasks 
		WHERE owner_id = %s
		ORDER BY
			CASE status
				WHEN 'Done' THEN 1
				ELSE 0
			END ASC,
			due_date ASC,
			created_at DESC`

	var query string
	if r.dbType == "mysql" {
		query = fmt.Sprintf(baseQuery, "?")
	} else {
		query = fmt.Sprintf(`
			SELECT id, title, description, TO_CHAR(due_date, 'YYYY-MM-DD') as due_date, status, owner_id, created_at, updated_at FROM tasks 
			WHERE owner_id = $1
			ORDER BY
				CASE status
					WHEN 'Done' THEN 1
					ELSE 0
				END ASC,
				due_date ASC,
				created_at DESC`)
	}

	rows, err := r.db.QueryContext(ctx, query, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []entity.Task
	for rows.Next() {
		var task entity.Task
		if err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.DueDate,
			&task.Status,
			&task.OwnerID,
			&task.CreatedAt,
			&task.UpdatedAt,
		); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, rows.Err()
}

func (r *taskRepository) Update(ctx context.Context, task *entity.Task) error {
	var query string
	if r.dbType == "mysql" {
		query = `
			UPDATE tasks
			SET title = ?, description = ?, due_date = STR_TO_DATE(?, '%Y-%m-%d'), status = ?, owner_id = ?, updated_at = ?
			WHERE id = ?`
	} else {
		query = `
			UPDATE tasks
			SET title = $1, description = $2, due_date = $3::date, status = $4, owner_id = $5, updated_at = $6
			WHERE id = $7`
	}

	_, err := r.db.ExecContext(ctx,
		query,
		task.Title,
		task.Description,
		task.DueDate,
		task.Status,
		task.OwnerID,
		time.Now(),
		task.ID,
	)
	return err
}

func (r *taskRepository) Delete(ctx context.Context, id int64) error {
	var query string
	if r.dbType == "mysql" {
		query = `DELETE FROM tasks WHERE id = ?`
	} else {
		query = `DELETE FROM tasks WHERE id = $1`
	}
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

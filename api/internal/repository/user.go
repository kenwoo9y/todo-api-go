package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/kenwoo9y/todo-api-go/api/internal/config"
	"github.com/kenwoo9y/todo-api-go/api/internal/entity"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	GetAll(ctx context.Context) ([]entity.User, error)
	GetByID(ctx context.Context, id int64) (*entity.User, error)
	GetByUsername(ctx context.Context, username string) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id int64) error
}

type userRepository struct {
	db     *sql.DB
	dbType string
}

func NewUserRepository(db *sql.DB, cfg *config.Config) UserRepository {
	return &userRepository{
		db:     db,
		dbType: cfg.DBType,
	}
}

func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	var query string
	if r.dbType == "mysql" {
		query = `
			INSERT INTO users (username, email, first_name, last_name, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?)`
	} else {
		query = `
			INSERT INTO users (username, email, first_name, last_name, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id`
	}

	now := time.Now()
	if r.dbType == "mysql" {
		result, err := r.db.ExecContext(ctx,
			query,
			user.Username,
			user.Email,
			user.FirstName,
			user.LastName,
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
		user.ID = id
		return nil
	} else {
		return r.db.QueryRowContext(ctx,
			query,
			user.Username,
			user.Email,
			user.FirstName,
			user.LastName,
			now,
			now,
		).Scan(&user.ID)
	}
}

func (r *userRepository) GetAll(ctx context.Context) ([]entity.User, error) {
	query := `SELECT * FROM users`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []entity.User
	for rows.Next() {
		var user entity.User
		if err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.FirstName,
			&user.LastName,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, rows.Err()
}

func (r *userRepository) GetByID(ctx context.Context, id int64) (*entity.User, error) {
	var user entity.User
	var query string
	if r.dbType == "mysql" {
		query = `SELECT * FROM users WHERE id = ?`
	} else {
		query = `SELECT * FROM users WHERE id = $1`
	}

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &user, err
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	var user entity.User
	var query string
	if r.dbType == "mysql" {
		query = `SELECT * FROM users WHERE username = ?`
	} else {
		query = `SELECT * FROM users WHERE username = $1`
	}

	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &user, err
}

func (r *userRepository) Update(ctx context.Context, user *entity.User) error {
	var query string
	if r.dbType == "mysql" {
		query = `
			UPDATE users
			SET username = ?, email = ?, first_name = ?, last_name = ?, updated_at = ?
			WHERE id = ?`
	} else {
		query = `
			UPDATE users
			SET username = $1, email = $2, first_name = $3, last_name = $4, updated_at = $5
			WHERE id = $6`
	}

	_, err := r.db.ExecContext(ctx,
		query,
		user.Username,
		user.Email,
		user.FirstName,
		user.LastName,
		time.Now(),
		user.ID,
	)
	return err
}

func (r *userRepository) Delete(ctx context.Context, id int64) error {
	var query string
	if r.dbType == "mysql" {
		query = `DELETE FROM users WHERE id = ?`
	} else {
		query = `DELETE FROM users WHERE id = $1`
	}
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

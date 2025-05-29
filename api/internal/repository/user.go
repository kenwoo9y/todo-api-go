package repository

import (
	"database/sql"
	"time"

	"github.com/kenwoo9y/todo-api-go/api/internal/entity"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *entity.User) error {
	query := `
		INSERT INTO users (username, email, first_name, last_name, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`

	return r.db.QueryRow(
		query,
		user.Username,
		user.Email,
		user.FirstName,
		user.LastName,
		time.Now(),
		time.Now(),
	).Scan(&user.ID)
}

func (r *UserRepository) GetAll() ([]entity.User, error) {
	query := `SELECT * FROM users`
	rows, err := r.db.Query(query)
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

func (r *UserRepository) GetByID(id int64) (*entity.User, error) {
	var user entity.User
	query := `SELECT * FROM users WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(
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

func (r *UserRepository) Update(user *entity.User) error {
	query := `
		UPDATE users
		SET username = $1, email = $2, first_name = $3, last_name = $4, updated_at = $5
		WHERE id = $6`

	_, err := r.db.Exec(
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

func (r *UserRepository) Delete(id int64) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

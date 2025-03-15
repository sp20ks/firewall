package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	"auth/internal/entity"

	"auth/internal/repository"

	"github.com/google/uuid"
)

type PostgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) repository.UserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) GetUser(username string) (*entity.User, error) {
	query := `SELECT id, username, password_hash, created_at FROM users WHERE username = $1`

	user := &entity.User{}
	err := r.db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Password, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

func (r *PostgresUserRepository) CreateUser(user *entity.User) error {
	query := `INSERT INTO users (id, username, password_hash, created_at) VALUES ($1, $2, $3, $4)`

	user.ID = uuid.New().String()
	_, err := r.db.Exec(query, user.ID, user.Username, user.Password, user.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

package repository

import (
	"database/sql"
	//"time"
	"fmt"

	"github.com/Paukku/ajanvarausjarjestelma/backend/internal/model"
	"github.com/google/uuid"
)

type PostgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) ExistsByEmail(email string) (bool, error) {
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email=$1)", email).Scan(&exists)
	return exists, err
}

func (r *PostgresUserRepository) CreateUser(user *model.User) (*model.User, error) {
	if user.UUID == uuid.Nil {
		user.UUID = uuid.New()
	}

	err := r.db.QueryRow(
		`INSERT INTO users(uuid, name, email, password_hash, role) 
		 VALUES($1, $2, $3, $4, $5) 
		 RETURNING created_at, updated_at`,
		user.UUID, user.Name, user.Email, user.PasswordHash, user.Role,
	).Scan(&user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *PostgresUserRepository) GetUsers(limit, offset int32) ([]*model.User, error) {
	// varmista, ettei limit ole negatiivinen tai liian suuri
	if limit <= 0 {
		limit = 50
	} else if limit > 100 {
		limit = 100
	}

	if offset < 0 {
		offset = 0
	}

	query := `
		SELECT uuid, name, email, role, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query users with limit=%d offset=%d: %w", limit, offset, err)
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		user := &model.User{}
		err := rows.Scan(&user.UUID, &user.Name, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user row: %w", err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return users, nil
}

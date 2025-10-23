package repository

import (
	"database/sql"
	//"time"

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

func (r *PostgresUserRepository) Create(user *model.User) (*model.User, error) {
	if user.UUID == uuid.Nil {
		user.UUID = uuid.New()
	}

	err := r.db.QueryRow(
		`INSERT INTO users(uuid, name, email, password_hash) 
		 VALUES($1, $2, $3, $4) 
		 RETURNING created_at, updated_at`,
		user.UUID, user.Name, user.Email, user.PasswordHash,
	).Scan(&user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return user, nil
}

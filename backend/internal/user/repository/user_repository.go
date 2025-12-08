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

func (r *PostgresUserRepository) GetUsers() ([]*model.User, error) {
	rows, err := r.db.Query("SELECT uuid, name, email, role, created_at, updated_at FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		user := &model.User{}
		err := rows.Scan(&user.UUID, &user.Name, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

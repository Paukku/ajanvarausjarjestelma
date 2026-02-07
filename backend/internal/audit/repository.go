package audit

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

/*
type Repository interface {
	Insert(
		ctx context.Context,
		action string,
		entity string,
		entityID *uuid.UUID,
		actorID *uuid.UUID,
	) error

	Find(
		ctx context.Context,
		limit, offset int,
	) ([]Log, error)
}*/

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Insert(
	ctx context.Context,
	action string,
	entity string,
	entityID *uuid.UUID,
	actorID *uuid.UUID,
	ipAddress *string,
	userAgent *string,
) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO audit_logs (
			id, action, entity, entity_id, actor_id, ip_address, user_agent
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7
		)
	`, uuid.New(), action, entity, entityID, actorID, ipAddress, userAgent)
	return err
}

func (r *PostgresRepository) Find(
	ctx context.Context,
	limit, offset int,
) ([]Log, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT
			id,
			action,
			entity,
			entity_id,
			actor_id,
			ip_address,
			user_agent,
			created_at
		FROM audit_logs
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []Log

	for rows.Next() {
		var log Log
		err := rows.Scan(
			&log.ID,
			&log.Action,
			&log.Entity,
			&log.EntityID,
			&log.ActorID,
			&log.IPAddress,
			&log.UserAgent,
			&log.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	return logs, nil
}

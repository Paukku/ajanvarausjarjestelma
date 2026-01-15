package audit

import (
	"context"

	"github.com/Paukku/ajanvarausjarjestelma/backend/internal/auth/actorctx"
	"github.com/google/uuid"
)

type Service struct {
	repo *PostgresRepository
}

func (s *Service) Log(
	ctx context.Context,
	action string,
	entity string,
	entityID *uuid.UUID,
) {

	var actorID *uuid.UUID

	if id, ok := actorctx.ActorIDFromContext(ctx); ok {
		actorID = id
	}
	// Ignoring error for logging
	_ = s.repo.Insert(ctx, action, entity, entityID, actorID, nil, nil)
}

func (s *Service) GetLogs(
	ctx context.Context,
	limit, offset int,
) ([]Log, error) {
	return s.repo.Find(ctx, limit, offset)
}

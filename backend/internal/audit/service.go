package audit

import (
	"context"

	"github.com/Paukku/ajanvarausjarjestelma/backend/internal/auth/actorctx"
	"github.com/google/uuid"
)

type Repository interface {
	Insert(ctx context.Context, action string, entity string, entityID *uuid.UUID, actorID *uuid.UUID) error
}

type Service struct {
	repo Repository
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
	_ = s.repo.Insert(ctx, action, entity, entityID, actorID)
}

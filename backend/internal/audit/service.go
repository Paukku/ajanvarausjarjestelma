package audit

import (
	"context"

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
	actorID *uuid.UUID,
) {
	// EI saa kaataa requestia
	_ = s.repo.Insert(ctx, action, entity, entityID, actorID)
}

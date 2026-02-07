package actorctx

import (
	"context"

	"github.com/google/uuid"
)

// unexported key type variable for storing actorID in context
type actorIDKeyType struct{}

var actorIDKey = actorIDKeyType{}

// ActorIDFromContext hakee actorID:n contextista
func ActorIDFromContext(ctx context.Context) (*uuid.UUID, bool) {
	actorID, ok := ctx.Value(actorIDKey).(uuid.UUID)
	if !ok {
		return nil, false
	}
	return &actorID, true
}

// ContextWithActorID lisää actorID:n contextiin
func ContextWithActorID(ctx context.Context, actorID uuid.UUID) context.Context {
	return context.WithValue(ctx, actorIDKey, actorID)
}

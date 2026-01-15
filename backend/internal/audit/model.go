package audit

import (
	"time"

	"github.com/google/uuid"
)

type Log struct {
	ID        uuid.UUID
	Action    string
	Entity    string
	EntityID  *uuid.UUID
	ActorID   *uuid.UUID
	IPAddress *string
	UserAgent *string
	CreatedAt time.Time
}

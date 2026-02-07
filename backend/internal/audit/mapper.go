package audit

import (
	"time"

	pb "github.com/Paukku/ajanvarausjarjestelma/backend/pb/common"
	"github.com/google/uuid"
)

func mapLogToProto(log Log) *pb.AuditLog {
	return &pb.AuditLog{
		Id:        log.ID.String(),
		Action:    log.Action,
		Entity:    log.Entity,
		EntityId:  uuidPtrToString(log.EntityID),
		ActorId:   uuidPtrToString(log.ActorID),
		IpAddress: stringPtr(log.IPAddress),
		UserAgent: stringPtr(log.UserAgent),
		CreatedAt: log.CreatedAt.Format(time.RFC3339),
	}
}

func uuidPtrToString(id *uuid.UUID) string {
	if id == nil {
		return ""
	}
	return id.String()
}

func stringPtr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

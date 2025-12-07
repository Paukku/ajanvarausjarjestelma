package model

import (
	"strings"
	"time"

	pb "github.com/Paukku/ajanvarausjarjestelma/backend/pb/common"
	"github.com/google/uuid"
)

type User struct {
	UUID         uuid.UUID
	Name         string
	Email        string
	PasswordHash string
	Role         string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
}

func ConvertUserListToPB(users []*User) []*pb.User {
	result := make([]*pb.User, 0, len(users))

	for _, u := range users {
		result = append(result, &pb.User{
			Uuid:  u.UUID.String(),
			Name:  u.Name,
			Email: u.Email,
			Role:  convertRoleToPB(u.Role),
		})
	}

	return result
}

func convertRoleToPB(role string) pb.UserRole {
	switch strings.ToUpper(strings.TrimSpace(role)) {
	case "ADMIN":
		return pb.UserRole_ADMIN
	case "OWNER":
		return pb.UserRole_OWNER
	case "EMPLOYEE":
		return pb.UserRole_EMPLOYEE
	case "UNAUTHORIZED":
		return pb.UserRole_UNAUTHORIZED
	default:
		return pb.UserRole_UNAUTHORIZED
	}
}

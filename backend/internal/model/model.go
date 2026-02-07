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
		result = append(result, convertUserToPB(u))
	}

	return result
}

func ConvertUserToPB(user *User) *pb.User {
	return convertUserToPB(user)
}

func convertUserToPB(user *User) *pb.User {
	return &pb.User{
		Uuid:  user.UUID.String(),
		Name:  user.Name,
		Email: user.Email,
		Role:  convertRoleToPB(user.Role),
	}
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

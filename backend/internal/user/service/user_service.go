package service

import (
	"context"
	"fmt"

	"github.com/Paukku/ajanvarausjarjestelma/backend/internal/audit"
	"github.com/Paukku/ajanvarausjarjestelma/backend/internal/model"
	"github.com/Paukku/ajanvarausjarjestelma/backend/internal/user/repository"
	pb "github.com/Paukku/ajanvarausjarjestelma/backend/pb/common"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceServer struct {
	audit *audit.Service
	Repo  *repository.PostgresUserRepository
}

func NewUserServiceServer(repo *repository.PostgresUserRepository) *UserServiceServer {
	return &UserServiceServer{Repo: repo}
}

func (s *UserServiceServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.GeneralResponse, error) {
	exists, err := s.Repo.ExistsByEmail(req.Email)
	if err != nil {
		return nil, err
	}

	if exists {
		return &pb.GeneralResponse{Success: false, Message: "Email already exists"}, nil
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.GetPassword()), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &model.User{
		Name:         req.GetName(),
		Email:        req.GetEmail(),
		PasswordHash: string(hashed),
		Role:         pb.UserRole_UNAUTHORIZED.String(),
	}

	createdUser, err := s.Repo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	s.audit.Log(
		ctx,
		"USER_CREATED",
		"user",
		&createdUser.UUID,
	)

	return &pb.GeneralResponse{Success: true, Message: "User created!"}, nil
}

func (s *UserServiceServer) GetUsers(ctx context.Context, limit, offset int32) (*pb.GetUsersResponse, error) {
	users, err := s.Repo.GetUsers(limit, offset)
	if err != nil {
		return nil, err
	}

	return &pb.GetUsersResponse{Users: model.ConvertUserListToPB(users)}, nil
}

func (s *UserServiceServer) GetUserById(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
	return &pb.User{Uuid: "1", Name: "Test User"}, nil
}

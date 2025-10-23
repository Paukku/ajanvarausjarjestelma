package service

import (
	"context"

	"github.com/Paukku/ajanvarausjarjestelma/backend/internal/model"
	"github.com/Paukku/ajanvarausjarjestelma/backend/internal/user/repository"
	pb "github.com/Paukku/ajanvarausjarjestelma/backend/pb/common"
)

type UserServiceServer struct {
	Repo *repository.PostgresUserRepository
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

	user := &model.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: req.Password, // myöhemmin hashaus tähän
	}

	_, err = s.Repo.Create(user)
	if err != nil {
		return nil, err
	}

	return &pb.GeneralResponse{Success: true, Message: "User created!"}, nil
}

func (s *UserServiceServer) GetUser(ctx context.Context, req *pb.EmptyRequest) (*pb.UserList, error) {
	return &pb.UserList{Users: []*pb.User{}}, nil
}

func (s *UserServiceServer) GetUserById(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
	return &pb.User{Uuid: "1", Name: "Test User"}, nil
}

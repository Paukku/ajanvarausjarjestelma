package handler

import (
	"context"

	"github.com/Paukku/ajanvarausjarjestelma/backend/internal/user/service"
	pb "github.com/Paukku/ajanvarausjarjestelma/backend/pb/common"
)

type UserHandler struct {
	userService *service.UserServiceServer
}

func NewUserHandler(userService *service.UserServiceServer) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.GeneralResponse, error) {
	return h.userService.CreateUser(ctx, req)
}

func (h *UserHandler) GetUser(ctx context.Context, req *pb.EmptyRequest) (*pb.UserList, error) {
	return h.userService.GetUser(ctx, req)
}

func (h *UserHandler) GetUserById(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
	return h.userService.GetUserById(ctx, req)
}

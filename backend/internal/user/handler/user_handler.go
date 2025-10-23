package handler

import (
	"github.com/Paukku/ajanvarausjarjestelma/backend/internal/user/service"
)

type UserHandler struct {
	userService *service.UserServiceServer
}

func NewUserHandler(userService *service.UserServiceServer) *UserHandler {
	return &UserHandler{userService: userService}
}

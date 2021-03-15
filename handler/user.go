package handler

import (
	"context"
	"errors"

	"github.com/booking-man-be/lib/logger"
	userPb "github.com/booking-man-be/proto/user"
	"github.com/golang/protobuf/ptypes/empty"
)

type userHandler struct {
}

func NewUserHandler() userPb.UserServer {
	return &userHandler{}
}

func (h *userHandler) RegisterUser(ctx context.Context, req *userPb.RegisterUserRequest) (*empty.Empty, error) {
	logger.Info(req)
	return nil, errors.New("Not implemented yet")
}

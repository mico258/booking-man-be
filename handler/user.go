package handler

import (
	"context"
	"errors"

	"github.com/booking-man-be/lib/logger"
	userPb "github.com/booking-man-be/proto/user"
	"github.com/booking-man-be/user"
	"github.com/golang/protobuf/ptypes/empty"
)

type userHandler struct {
	service user.Service
}

func NewUserHandler(service user.Service) userPb.UserServer {
	return &userHandler{
		service: service,
	}
}

func (h *userHandler) RegisterUser(ctx context.Context, req *userPb.RegisterUserRequest) (*empty.Empty, error) {
	logger.Info(req)
	return nil, errors.New("Not implemented yet")
}

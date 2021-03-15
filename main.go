package main

import (
	"github.com/booking-man-be/handler"
	"github.com/booking-man-be/lib/logger"
	"github.com/booking-man-be/lib/server"
	userPb "github.com/booking-man-be/proto/user"
)

func main() {

	// TODO change port to config
	svc := server.NewService(
		server.GRPCPort("9090"),
		server.RESTPort("80"),
	)

	svc.Init()

	// init handler
	userHandler := handler.NewUserHandler()

	// register handler to grpc and rest
	userPb.RegisterUserServer(svc.Server(), userHandler)
	svc.RegisterRESTHandler(userPb.RegisterUserHandler)

	if err := <-svc.RunServers(); err != nil {
		logger.Fatal(err)
	}

}

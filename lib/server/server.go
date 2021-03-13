package server

import (
	"context"
	"errors"
	"net"
	"sync"

	gwruntime "github.com/grpc-ecosystem/grpc-gateway/runtime"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type restHandler func(ctx context.Context, mux *gwruntime.ServeMux, conn *grpc.ClientConn) error

type Service interface {
	// GRPC method
	Server() *grpc.Server
	UseServerUnaryInterceptor(interceptors ...grpc.UnaryServerInterceptor)
	Init()
	ListenAndServeGRPC() error

	// Rest gateway method
	ListenAndServeGateway(ctx context.Context, enableCORS bool) error
	RunServers() <-chan error
	RegisterRESTHandler(restHandler ...restHandler)
	Shutdown(context.Context) error
}

type service struct {
	server       *grpc.Server
	serverMux    sync.Mutex
	options      Options
	interceptors Interceptors
	restHandlers []restHandler
}

func (s *service) initServer() {
	s.serverMux.Lock()
	defer s.serverMux.Unlock()

	srv := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.ChainUnaryServer(s.interceptors.serverUnary...)),
		// TODO: need to tidy up this code
		grpc.MaxRecvMsgSize(rgDefaultServerMaxReceiveMessageSize),
		grpc.MaxSendMsgSize(rgDefaultServerMaxSendMessageSize),
		grpc.MaxMsgSize(rgDefaultServerMaxSizeMessageSize),
		// TODO: create StreamInterceptor
	)
	s.server = srv
	if s.options.gRPCReflection {
		reflection.Register(srv)
	}
}

func (s *service) Shutdown(ctx context.Context) error {
	<-ctx.Done()
	return nil
}

func (s *service) initServer() {
	s.serverMux.Lock()
	defer s.serverMux.Unlock()

	srv := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.ChainUnaryServer(s.interceptors.serverUnary...)),
		// TODO: need to tidy up this code
		grpc.MaxRecvMsgSize(rgDefaultServerMaxReceiveMessageSize),
		grpc.MaxSendMsgSize(rgDefaultServerMaxSendMessageSize),
		grpc.MaxMsgSize(rgDefaultServerMaxSizeMessageSize),
		// TODO: create StreamInterceptor
	)
	s.server = srv
	if s.options.gRPCReflection {
		reflection.Register(srv)
	}
}

// ListenAndServeGRPC function to run and listen GRPC service
func (s *service) ListenAndServeGRPC() error {
	if s.server == nil {
		return errors.New("err: please init server first")
	}
	return gRPCServe(s.Server(), s.options.gRPCPort)
}

func gRPCServe(srv *grpc.Server, port string) error {
	logger.Infof("starting gRPC server at :%s...", port)

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	return srv.Serve(lis)
}

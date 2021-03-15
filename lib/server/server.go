package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/booking-man-be/lib/logger"
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

type Interceptors struct {
	serverStream []grpc.StreamServerInterceptor
	serverUnary  []grpc.UnaryServerInterceptor
	clientStream []grpc.StreamClientInterceptor
	clientUnary  []grpc.UnaryClientInterceptor
}

type service struct {
	server       *grpc.Server
	serverMux    sync.Mutex
	options      Options
	interceptors Interceptors
	restHandlers []restHandler
}

// OptionFunc function type of Option struct
type OptionFunc func(o *Options)

func generateOptions(fs ...OptionFunc) Options {
	var marshalers []gwruntime.ServeMuxOption

	opts := Options{
		gRPCPort: DefaultGRPCPort,
		restPort: DefaultRESTPort,
		network:  DefaultNetwork,
		gRPCHost: DefaultGRPHost,
	}
	for _, f := range fs {
		f(&opts)
	}

	opts.restServeMuxOpts = append(marshalers, opts.restServeMuxOpts...)
	return opts
}

// GRPCPort to set gRPC port option
func GRPCPort(p string) OptionFunc {
	return func(o *Options) {
		o.gRPCPort = p
	}
}

// RESTPort to set rest gateway port option
func RESTPort(p string) OptionFunc {
	return func(o *Options) {
		o.restPort = p
	}
}

func NewService(fs ...OptionFunc) Service {
	opts := generateOptions(fs...)

	return &service{
		options: opts,
	}
}

// UseServerUnaryInterceptor method to collect define interceptor
func (s *service) UseServerUnaryInterceptor(interceptors ...grpc.UnaryServerInterceptor) {

	s.interceptors.serverUnary = append(s.interceptors.serverUnary, interceptors...)

}

func (s *service) Init() {
	s.initServer()
}

// RunServers starts the grpc and rest gateway server
func (s *service) RunServers() <-chan error {
	ch := make(chan error, 2)
	go func() {
		logger.Infof("Initializing gRPC connection in port %s", s.options.gRPCPort)
		if err := s.ListenAndServeGRPC(); err != nil {
			ch <- fmt.Errorf("cannot run grpc service: %v", err)
		}
	}()

	go func() {
		logger.Infof("Initializing rest gateway connection in port %s", s.options.restPort)
		if err := s.ListenAndServeGateway(context.Background(), true); err != nil {
			ch <- fmt.Errorf("cannot run gateway service: %v", err)
		}
	}()

	return ch
}

func (s *service) RegisterRESTHandler(restHandler ...restHandler) {
	s.restHandlers = append(s.restHandlers, restHandler...)
}

func (s *service) Shutdown(ctx context.Context) error {
	<-ctx.Done()
	return nil
}

// ListenAndServeGRPC function to run and listen GRPC service
func (s *service) ListenAndServeGRPC() error {
	if s.server == nil {
		return errors.New("err: please init server first")
	}
	return gRPCServe(s.Server(), s.options.gRPCPort)
}

func (s *service) Server() *grpc.Server {
	if s.server == nil {
		s.initServer()
	}
	return s.server
}

func (s *service) initServer() {
	s.serverMux.Lock()
	defer s.serverMux.Unlock()

	srv := grpc.NewServer(
		grpc.UnaryInterceptor(ChainUnaryServer(s.interceptors.serverUnary...)),
		grpc.MaxRecvMsgSize(defaultMessageSize),
		grpc.MaxSendMsgSize(defaultMessageSize),
		grpc.MaxMsgSize(defaultMessageSize),
	)
	s.server = srv
	if s.options.gRPCReflection {
		reflection.Register(srv)
	}
}

func gRPCServe(srv *grpc.Server, port string) error {
	logger.Infof("starting gRPC server at :%s...", port)

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	return srv.Serve(lis)
}

func (s *service) initRESTHandler(ctx context.Context) (http.Handler, error) {
	// prepare mux and connection
	mux := gwruntime.NewServeMux(s.options.restServeMuxOpts...)
	conn, err := s.dialSelf()
	if err != nil {
		return nil, err
	}

	// init rest handler
	for _, h := range s.restHandlers {
		if err := h(ctx, mux, conn); err != nil {
			return nil, err
		}
	}
	return mux, err
}

func (s *service) dialSelf() (*grpc.ClientConn, error) {
	return dial(s.options.network, fmt.Sprintf("%s:%s", s.options.gRPCHost, s.options.gRPCPort))
}

func dial(network, addr string) (*grpc.ClientConn, error) {
	return dialContext(context.Background(), network, addr)
}

func dialContext(ctx context.Context, network, addr string) (*grpc.ClientConn, error) {
	switch network {
	case "tcp":
		return grpc.DialContext(ctx, addr, grpc.WithInsecure(), grpc.WithMaxMsgSize(defaultMessageSize))
	case "unix":
		d := func(addr string, timeout time.Duration) (net.Conn, error) {
			return net.DialTimeout("unix", addr, timeout)
		}
		return grpc.DialContext(ctx, addr, grpc.WithInsecure(), grpc.WithDialer(d), grpc.WithMaxMsgSize(defaultMessageSize))

	default:
		return nil, fmt.Errorf("unsupported network type %q", network)
	}
}

func (s *service) ListenAndServeGateway(ctx context.Context, enableCORS bool) error {
	mux := http.NewServeMux()
	args := &http.Server{
		Addr: fmt.Sprintf(":%s", s.options.restPort),
	}
	if enableCORS {
		args.Handler = allowCORS(mux)
	}

	// register gateway handlers
	handler, err := s.initRESTHandler(ctx)
	if err != nil {
		return err
	}

	mux.Handle("/", handler)
	go func() {
		<-ctx.Done()
		logger.Infof("Shutting down the http server")
		if err := args.Shutdown(context.Background()); err != nil {
			logger.Errorf("Failed to shutdown http server: %v\n", err)
		}
	}()
	logger.Infof("starting rest gateway server at :%s...", s.options.restPort)
	if err := args.ListenAndServe(); err != http.ErrServerClosed {
		logger.Errorf("Failed to listen and serve: %v\n", err)
		return err
	}
	return nil
}

func allowCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			if r.Method == "OPTIONS" && r.Header.Get("Access-Control-Request-Method") != "" {
				preflightHandler(w, r)
				return
			}
		}
		h.ServeHTTP(w, r)
	})
}

func preflightHandler(w http.ResponseWriter, r *http.Request) {
	headers := []string{HeaderContentType, HeaderAccept, HeaderAuthorization}
	w.Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ","))
	methods := []string{"GET", "HEAD", "POST", "PUT", "DELETE"}
	w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ","))
	logger.Infof("preflight request for %s\n", r.URL.Path)
}

func ChainUnaryServer(interceptors ...grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
	n := len(interceptors)
	if n > 1 {
		lastI := n - 1
		return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
			var (
				chainHandler grpc.UnaryHandler
				curI         int
			)

			chainHandler = func(currentCtx context.Context, currentReq interface{}) (interface{}, error) {
				if curI == lastI {
					return handler(currentCtx, currentReq)
				}
				curI++
				resp, err := interceptors[curI](currentCtx, currentReq, info, chainHandler)
				curI--
				return resp, err
			}

			return interceptors[0](ctx, req, info, chainHandler)
		}
	}

	if n == 1 {
		return interceptors[0]
	}

	// n == 0; Dummy interceptor maintained for backward compatibility to avoid returning nil.
	return func(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		return handler(ctx, req)
	}
}

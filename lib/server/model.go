package server

import (
	gwruntime "github.com/grpc-ecosystem/grpc-gateway/runtime"
)

const (
	DefaultGRPCPort      = "9090"
	DefaultRESTPort      = "8080"
	DefaultNetwork       = "tcp"
	DefaultGRPHost       = "localhost"
	formURLEncodedHeader = "application/x-www-form-urlencoded"
)

const (
	HeaderContentType   = "Content-Type"
	HeaderAccept        = "Accept"
	HeaderAuthorization = "authorization"
	HeaderRequestID     = "requestid"
	HeaderDeviceID      = "deviceid"
	HeaderClientID      = "clientid"
)

const defaultMessageSize = 1024 * 1024 * 20

// Options is model to set option configuration
// of gRPC and REST gateway
type Options struct {
	// gRPCHost is host of gRPC service
	// by default is http://localhost
	gRPCHost string
	// gRPCPort is port of gRPC service
	// by default is 9090
	gRPCPort string
	// restPort is port of rest gateway service
	// by default is 8080
	restPort string
	// network is type of network protocol
	// by default using tcp connection
	network string

	// restServerMuxOpts is property to set option of rest service
	restServeMuxOpts []gwruntime.ServeMuxOption

	// gRPCReflection
	gRPCReflection bool

	// to enable request using form data
	enableFormData bool
}

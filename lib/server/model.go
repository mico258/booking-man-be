package server

import (
	"gitlab.com/ruangguru/source/shared-lib/go/morse/gateway"

	gwruntime "github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	HeaderContentType         = "Content-Type"
	HeaderAccept              = "Accept"
	HeaderAuthorization       = "authorization"
	HeaderInternalAPIPassword = "internalapipassword"
	HeaderUserCode            = "usercode"
	HeaderRole                = "role"
	HeaderRequestID           = "requestid"
	HeaderCoutry              = "country"
	HeaderDeviceID            = "deviceid"
	HeaderClientID            = "clientid"
	HeaderOAuthClientID       = "oauthclientid" // will be used in login endpoint header
	HeaderAppName             = "appname"
	HeaderAppVersion          = "appversion"
	HeaderAppVersionCode      = "appversioncode"
	HeaderManufacture         = "manufacture"
	HeaderModel               = "model"
	HeaderPlatform            = "platform"
	HeaderPlatformVersion     = "platformversion"
	HeaderSDKVersion          = "sdkversion"
	HeaderUserAgent           = "user-agent"
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
	// disableRest is property to diable rest service
	// in case we only need gRPC service
	disableRest bool
	// restServerMuxOpts is property to set option of rest service
	restServeMuxOpts []gwruntime.ServeMuxOption

	// gRPCReflection
	gRPCReflection bool

	// to enable request using form data
	enableFormData bool
	// prometheusCollectors are a list of prometheus collector
	prometheusCollectors []prometheus.Collector
	// enableLegacyError to enable backward compatible with previous api (GIN)
	// added field 'code' & 'error_description' on API response
	// note: you must add 'option (morse.api.response.data) = Array;' on your proto file to make the response data return array
	errorResponseType gateway.ErrorResponseType
}

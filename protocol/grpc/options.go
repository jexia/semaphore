package grpc

import (
	"github.com/jexia/maestro/specs"
)

const (
	// ServiceOption represents the gRPC service option key
	ServiceOption = "service"
	// MethodOption represents the gRPC method option key
	MethodOption = "method"
)

// ListenerOptions represents the available HTTP options
type ListenerOptions struct {
	Service string
	Method  string
}

// ParseListenerOptions parses the given specs options into HTTP options
func ParseListenerOptions(options specs.Options) (*ListenerOptions, error) {
	result := &ListenerOptions{}

	service, has := options[ServiceOption]
	if has {
		result.Service = service
	}

	method, has := options[MethodOption]
	if has {
		result.Method = method
	}

	return result, nil
}

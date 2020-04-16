package grpc

import (
	"github.com/jexia/maestro/pkg/specs"
	"github.com/jexia/maestro/pkg/transport"
)

const (
	// ServiceOption represents the service name option key
	ServiceOption = "service"
	// MethodOption represents the method name option key
	MethodOption = "method"
	// PackageOption represents the package name option key
	PackageOption = "package"
)

// ListenerOptions represents the available HTTP options
type ListenerOptions struct {
}

// ParseListenerOptions parses the given specs options into HTTP options
func ParseListenerOptions(options specs.Options) (*ListenerOptions, error) {
	result := &ListenerOptions{}

	return result, nil
}

// EndpointOptions represents the available HTTP options
type EndpointOptions struct {
	Package string
	Service string
	Method  string
}

// ParseEndpointOptions parses the given specs options into HTTP options
func ParseEndpointOptions(endpoint *transport.Endpoint) (*EndpointOptions, error) {
	result := &EndpointOptions{
		Package: "maestro",
		Service: "service",
		Method:  endpoint.Flow.GetName(),
	}

	pkg, has := endpoint.Options[PackageOption]
	if has {
		result.Package = pkg
	}

	service, has := endpoint.Options[ServiceOption]
	if has {
		result.Service = service
	}

	method, has := endpoint.Options[MethodOption]
	if has {
		result.Method = method
	}

	return result, nil
}

// CallerOptions represents the available HTTP options
type CallerOptions struct {
}

// ParseCallerOptions parses the given specs options into HTTP options
func ParseCallerOptions(options specs.Options) (*CallerOptions, error) {
	result := &CallerOptions{}

	return result, nil
}

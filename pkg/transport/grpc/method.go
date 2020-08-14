package grpc

import (
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/transport"
)

// Method represents a gRPC endpoint
type Method struct {
	*transport.Endpoint
	fqn  string
	name string
	flow transport.Flow
}

// GetName returns the method name
func (method *Method) GetName() string { return method.name }

// References returns the available method references
func (method *Method) References() []*specs.Property {
	return make([]*specs.Property, 0)
}

// GetRequest returns the request input parameter map
func (method *Method) GetRequest() map[string]*specs.Property {
	if method.Request == nil {
		return make(map[string]*specs.Property, 0)
	}

	return method.Request.Definition.Property.Nested
}

// GetResponse returns the request output parameter map
func (method *Method) GetResponse() map[string]*specs.Property {
	if method.Response == nil {
		return make(map[string]*specs.Property, 0)
	}

	return method.Response.Definition.Property.Nested
}

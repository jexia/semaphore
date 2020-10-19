package grpc

import (
	"fmt"

	"github.com/jexia/semaphore/pkg/codec/proto"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/transport"
)

// Method represents a gRPC endpoint
type Method struct {
	*transport.Endpoint
	Service *proto.Service
	Name    string
	Flow    transport.Flow
}

func (method *Method) String() string {
	return fmt.Sprintf("%s/%s", method.Service, method.Name)
}

// GetName returns the method name
func (method *Method) GetName() string { return method.Name }

// References returns the available method references
func (method *Method) References() []*specs.Property {
	return make([]*specs.Property, 0)
}

// GetRequest returns the request input parameter map
func (method *Method) GetRequest() specs.Message {
	if method.Request == nil {
		return make(specs.Message, 0)
	}

	return method.Request.Definition.Property.Message
}

// GetResponse returns the request output parameter map
func (method *Method) GetResponse() specs.Message {
	if method.Response == nil {
		return make(specs.Message, 0)
	}

	return method.Response.Definition.Property.Message
}

package grpc

import (
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/jexia/maestro/pkg/codec"
	"github.com/jexia/maestro/pkg/metadata"
	"github.com/jexia/maestro/pkg/specs"
	"github.com/jexia/maestro/pkg/transport"
)

// Request method request
type Request struct {
	param  *specs.ParameterMap
	header *metadata.Manager
	codec  codec.Manager
}

// Method represents a gRPC endpoint
type Method struct {
	*transport.Endpoint
	fqn        string
	name       string
	flow       transport.Flow
	descriptor []byte
}

// GetName returns the method name
func (method *Method) GetName() string {
	return method.name
}

// References returns the available method references
func (method *Method) References() []*specs.Property {
	return make([]*specs.Property, 0)
}

// GetRequest returns the request input parameter map
func (method *Method) GetRequest() map[string]*specs.Property {
	if method.Request == nil {
		return make(map[string]*specs.Property, 0)
	}

	return method.Request.Schema.Property.Nested
}

// GetResponse returns the request output parameter map
func (method *Method) GetResponse() map[string]*specs.Property {
	if method.Response == nil {
		return make(map[string]*specs.Property, 0)
	}

	return method.Response.Schema.Property.Nested
}

// Service represents a gRPC service
type Service struct {
	pkg     string
	name    string
	methods map[string]*Method
	proto   *descriptor.FileDescriptorProto
}

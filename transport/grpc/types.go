package grpc

import (
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/jexia/maestro/codec"
	"github.com/jexia/maestro/specs"
	"github.com/jexia/maestro/transport"
)

// Method represents a gRPC endpoint
type Method struct {
	fqn        string
	name       string
	flow       transport.Flow
	in         *specs.ParameterMap
	req        codec.Manager
	out        *specs.ParameterMap
	res        codec.Manager
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
	return method.in.Property.Nested
}

// GetResponse returns the request output parameter map
func (method *Method) GetResponse() map[string]*specs.Property {
	return method.out.Property.Nested
}

// Service represents a gRPC service
type Service struct {
	pkg        string
	name       string
	methods    map[string]*Method
	descriptor []byte
	file       *descriptor.FileDescriptorProto
}

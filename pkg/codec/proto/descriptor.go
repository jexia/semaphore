package proto

import (
	"strings"

	"github.com/jexia/semaphore/pkg/broker/trace"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/types"
	"github.com/jhump/protoreflect/desc/builder"
)

// Method represents a service mthod
type Method interface {
	GetName() string
	GetRequest() map[string]*specs.Property
	GetResponse() map[string]*specs.Property
}

// Methods represents a collection of methods
type Methods map[string]Method

// NewFile constructs and returns a new file builder
func NewFile(name string) *builder.FileBuilder {
	return builder.NewFile(name)
}

// NewMessageDescriptor constructs a new protobuf message descriptor for the given parameter map
func NewMessageDescriptor(path string, specs *specs.ParameterMap) ([]byte, error) {
	prop := specs.Property
	if prop.Type != types.Message {
		return nil, trace.New(trace.WithMessage("a proto message always requires a root message"))
	}

	desc, err := NewMessage(path, prop.Nested)
	if err != nil {
		return nil, err
	}

	bb, _ := desc.AsDescriptorProto().Descriptor()
	return bb, nil
}

// NewServiceDescriptor constructs a new protobuf service descriptor for the given parameter map
func NewServiceDescriptor(file *builder.FileBuilder, name string, methods Methods) error {
	service := builder.NewService(name)

	for _, method := range methods {
		name := strings.Title(name) + strings.Title(method.GetName())
		req := builder.NewMessage(name + "Request")
		err := ConstructMessage(req, method.GetRequest())
		if err != nil {
			return err
		}

		resp := builder.NewMessage(name + "Response")
		err = ConstructMessage(resp, method.GetResponse())
		if err != nil {
			return err
		}

		err = file.TryAddMessage(req)
		if err != nil {
			return err
		}

		err = file.TryAddMessage(resp)
		if err != nil {
			return err
		}

		method := builder.NewMethod(method.GetName(), builder.RpcTypeMessage(req, false), builder.RpcTypeMessage(resp, false))
		err = service.TryAddMethod(method)
		if err != nil {
			return err
		}
	}

	err := file.TryAddService(service)
	if err != nil {
		return err
	}

	return nil
}

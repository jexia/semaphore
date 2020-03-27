package proto

import (
	"github.com/jexia/maestro/specs"
	"github.com/jexia/maestro/specs/trace"
	"github.com/jexia/maestro/specs/types"
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
	file.AddService(service)

	for _, method := range methods {
		in := builder.NewMessage("something")
		err := ConstructMessage(in, method.GetRequest())
		if err != nil {
			return err
		}

		file.AddMessage(in)

		out := builder.NewMessage("else")
		err = ConstructMessage(out, method.GetResponse())
		if err != nil {
			return err
		}

		file.AddMessage(out)

		method := builder.NewMethod(method.GetName(), builder.RpcTypeMessage(in, false), builder.RpcTypeMessage(out, false))
		service.AddMethod(method)
	}

	return nil
}

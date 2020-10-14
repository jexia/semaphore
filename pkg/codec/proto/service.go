package proto

import (
	"fmt"
	"strings"

	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/builder"
)

// Service represents a gRPC service.
type Service struct {
	Package string
	Name    string
	Methods Methods
}

// String returns service full name.
func (service *Service) String() string {
	return fmt.Sprintf("%s.%s", service.Package, service.Name)
}

// FileDescriptor generates protobuf file descriptor.
func (service *Service) FileDescriptor() (*desc.FileDescriptor, error) {
	var (
		file         = builder.NewFile(service.String())
		protoService = builder.NewService(service.Name)
	)

	file.IsProto3 = true

	for _, method := range service.Methods {
		var (
			name = strings.Title(service.Name) + strings.Title(method.GetName())
			req  = builder.NewMessage(name + "Request")
		)

		if err := ConstructMessage(
			make(map[string]*builder.MessageBuilder),
			make(map[string]*builder.FieldType),
			req,
			method.GetRequest(),
		); err != nil {
			return nil, err
		}

		resp := builder.NewMessage(name + "Response")
		if err := ConstructMessage(
			make(map[string]*builder.MessageBuilder),
			make(map[string]*builder.FieldType),
			resp,
			method.GetResponse(),
		); err != nil {
			return nil, err
		}

		if err := file.TryAddMessage(req); err != nil {
			return nil, err
		}

		if err := file.TryAddMessage(resp); err != nil {
			return nil, err
		}

		protoMethod := builder.NewMethod(method.GetName(), builder.RpcTypeMessage(req, false), builder.RpcTypeMessage(resp, false))
		if err := protoService.TryAddMethod(protoMethod); err != nil {
			return nil, err
		}
	}

	file.Package = service.Package

	if err := file.TryAddService(protoService); err != nil {
		return nil, err
	}

	return file.Build()
}

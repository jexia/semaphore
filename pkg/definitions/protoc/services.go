package protoc

import (
	"github.com/golang/protobuf/proto"
	"github.com/jexia/maestro/annotations"
	"github.com/jexia/maestro/pkg/specs"
	"github.com/jexia/maestro/pkg/transport/http"
	"github.com/jhump/protoreflect/desc"
)

const (
	// NameOption represents the Service name option key
	NameOption = "service_name"
	// HostOption represents the Service host option key
	HostOption = "service_host"
	// TransportOption represents the Service transport option key
	TransportOption = "service_transport"
	// CodecOption represents the Service codec option key
	CodecOption = "service_codec"
)

// NewServices constructs a new service(s) manifest from the given file descriptors
func NewServices(descriptors []*desc.FileDescriptor) []*specs.ServicesManifest {
	result := &specs.ServicesManifest{
		Services: make([]*specs.Service, 0),
	}

	for _, descriptor := range descriptors {
		for _, service := range descriptor.GetServices() {
			result.Services = append(result.Services, NewService(service))
		}
	}

	return []*specs.ServicesManifest{result}
}

// NewService constructs a new service with the given descriptor
func NewService(descriptor *desc.ServiceDescriptor) *specs.Service {
	options := specs.Options{}

	ext, err := proto.GetExtension(descriptor.GetOptions(), annotations.E_Service)
	if err == nil {
		ext := ext.(*annotations.Service)
		options[HostOption] = ext.GetHost()
		options[TransportOption] = ext.GetTransport()
		options[CodecOption] = ext.GetCodec()
	}

	result := &specs.Service{
		FullyQualifiedName: descriptor.GetFullyQualifiedName(),
		Name:               descriptor.GetName(),
		Package:            descriptor.GetFile().GetPackage(),
		Comment:            descriptor.GetSourceInfo().GetLeadingComments(),
		Host:               options[HostOption],
		Transport:          options[TransportOption],
		Codec:              options[CodecOption],
		Options:            options,
	}

	result.Methods = make([]*specs.Method, len(descriptor.GetMethods()))
	for index, method := range descriptor.GetMethods() {
		result.Methods[index] = NewMethod(method)
	}

	return result
}

// NewMethod constructs a new method with the given descriptor
func NewMethod(descriptor *desc.MethodDescriptor) *specs.Method {
	options := make(specs.Options)

	ext, err := proto.GetExtension(descriptor.GetOptions(), annotations.E_Http)
	if err == nil {
		ext := ext.(*annotations.HTTP)
		options[http.EndpointOption] = ext.GetEndpoint()
		options[http.MethodOption] = ext.GetMethod()
	}

	result := &specs.Method{
		Name:    descriptor.GetName(),
		Comment: descriptor.GetSourceInfo().GetLeadingComments(),
		Input:   descriptor.GetInputType().GetFullyQualifiedName(),
		Output:  descriptor.GetOutputType().GetFullyQualifiedName(),
		Options: options,
	}

	return result
}

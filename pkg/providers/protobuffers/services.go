package protobuffers

import (
	"github.com/golang/protobuf/proto"
	annotations "github.com/jexia/semaphore/api"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/transport/http"
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
	// RequestCodecOption represents the Service request codec option key
	RequestCodecOption = "service_request_codec"
	// ResponseCodecOption represents the Service response codec option key
	ResponseCodecOption = "service_response_codec"
)

// NewServices constructs a new service(s) manifest from the given file descriptors
func NewServices(descriptors []*desc.FileDescriptor) specs.ServiceList {
	result := make(specs.ServiceList, 0)

	for _, descriptor := range descriptors {
		for _, service := range descriptor.GetServices() {
			result = append(result, NewService(service))
		}
	}

	return result
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
		options[RequestCodecOption] = ext.GetRequestCodec()
		options[ResponseCodecOption] = ext.GetResponseCodec()
	}

	if options[TransportOption] == "" {
		options[TransportOption] = "grpc"
	}

	if options[CodecOption] == "" {
		options[CodecOption] = "proto"
	}

	result := &specs.Service{
		FullyQualifiedName: descriptor.GetFullyQualifiedName(),
		Name:               descriptor.GetName(),
		Package:            descriptor.GetFile().GetPackage(),
		Comment:            descriptor.GetSourceInfo().GetLeadingComments(),
		Host:               options[HostOption],
		Transport:          options[TransportOption],
		RequestCodec:       options[CodecOption],
		ResponseCodec:      options[CodecOption],
		Options:            options,
	}

	req := options[RequestCodecOption]
	if req != "" {
		result.RequestCodec = req
	}

	res := options[ResponseCodecOption]
	if res != "" {
		result.ResponseCodec = res
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

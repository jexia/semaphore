package protoc

import (
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/jexia/maestro/annotations"
	"github.com/jexia/maestro/protocol/http"
	"github.com/jexia/maestro/schema"
	"github.com/jexia/maestro/specs/types"
	"github.com/jhump/protoreflect/desc"
)

const (
	// HostOption represents the Service host option key
	HostOption = "service_host"
	// ProtocolOption represents the Service protocol option key
	ProtocolOption = "service_protocol"
	// CodecOption represents the Service codec option key
	CodecOption = "service_codec"
)

// NewCollection constructs a new schema collection from the given descriptors
func NewCollection(descriptors []*desc.FileDescriptor) schema.Collection {
	return &collection{
		descriptors: descriptors,
	}
}

type collection struct {
	descriptors []*desc.FileDescriptor
}

func (collection *collection) GetService(service string) schema.Service {
	for _, descriptor := range collection.descriptors {
		service := descriptor.FindService(service)
		if service == nil {
			continue
		}

		return NewService(service)
	}

	return nil
}

func (collection *collection) GetServices() []schema.Service {
	result := make([]schema.Service, 0)
	for _, descriptor := range collection.descriptors {
		for _, service := range descriptor.GetServices() {
			result = append(result, NewService(service))
		}
	}

	return result
}

func (collection *collection) GetMessage(message string) schema.Property {
	for _, descriptor := range collection.descriptors {
		message := descriptor.FindMessage(message)
		if message == nil {
			continue
		}

		return NewMessage(message)
	}

	return nil
}

func (collection *collection) GetMessages() []schema.Property {
	result := make([]schema.Property, 0)
	for _, descriptor := range collection.descriptors {
		for _, message := range descriptor.GetMessageTypes() {
			result = append(result, NewMessage(message))
		}
	}

	return result
}

// NewService constructs a new service with the given descriptor
func NewService(descriptor *desc.ServiceDescriptor) schema.Service {
	options := schema.Options{}

	ext, err := proto.GetExtension(descriptor.GetOptions(), annotations.E_Service)
	if err == nil {
		ext := ext.(*annotations.Service)
		options[HostOption] = ext.GetHost()
		options[ProtocolOption] = ext.GetProtocol()
		options[CodecOption] = ext.GetCodec()
	}

	return &service{
		descriptor: descriptor,
		options:    options,
	}
}

type service struct {
	descriptor *desc.ServiceDescriptor
	options    schema.Options
}

func (service *service) GetName() string {
	return service.descriptor.GetFullyQualifiedName()
}

func (service *service) GetHost() string {
	return service.options[HostOption]
}

func (service *service) GetProtocol() string {
	return service.options[ProtocolOption]
}

func (service *service) GetCodec() string {
	return service.options[CodecOption]
}

func (service *service) GetMethod(name string) schema.Method {
	for _, method := range service.descriptor.GetMethods() {
		if method.GetName() != name {
			continue
		}

		return NewMethod(method)
	}

	return nil
}

func (service *service) GetMethods() schema.Methods {
	result := make([]schema.Method, len(service.descriptor.GetMethods()))
	for index, method := range service.descriptor.GetMethods() {
		result[index] = NewMethod(method)
	}

	return result
}

func (service *service) GetOptions() schema.Options {
	return service.options
}

// NewMethod constructs a new method with the given descriptor
func NewMethod(descriptor *desc.MethodDescriptor) schema.Method {
	options := make(schema.Options)

	ext, err := proto.GetExtension(descriptor.GetOptions(), annotations.E_Http)
	if err == nil {
		ext := ext.(*annotations.HTTP)
		options[http.EndpointOption] = ext.GetEndpoint()
		options[http.MethodOption] = ext.GetMethod()
	}

	return &method{
		descriptor: descriptor,
		options:    options,
	}
}

type method struct {
	descriptor *desc.MethodDescriptor
	options    schema.Options
}

func (method *method) GetName() string {
	return method.descriptor.GetName()
}

func (method *method) GetInput() schema.Property {
	return NewMessage(method.descriptor.GetInputType())
}

func (method *method) GetOutput() schema.Property {
	method.descriptor.GetOutputType().AsDescriptorProto()
	return NewMessage(method.descriptor.GetOutputType())
}

func (method *method) GetOptions() schema.Options {
	return method.options
}

// NewMessage constructs a schema Property with the given message descriptor
func NewMessage(descriptor *desc.MessageDescriptor) schema.Property {
	return &message{
		desc:    descriptor,
		options: make(schema.Options),
	}
}

type message struct {
	desc    *desc.MessageDescriptor
	options schema.Options
}

// GetName returns the message name
func (message *message) GetName() string {
	return message.desc.GetFullyQualifiedName()
}

// GetPosition returns the property position inside a message
func (message *message) GetPosition() int32 {
	return 1
}

// GetType returns the message type
func (message *message) GetType() types.Type {
	return types.TypeMessage
}

// GetLabel returns the message label
func (message *message) GetLabel() types.Label {
	return types.LabelOptional
}

// GetNested attempts to return a all the nested properties
func (message *message) GetNested() map[string]schema.Property {
	fields := message.desc.GetFields()
	result := make(map[string]schema.Property, len(fields))
	for _, field := range fields {
		result[field.GetName()] = NewProperty(field)
	}

	return result
}

func (message *message) GetOptions() schema.Options {
	return message.options
}

// NewProperty constructs a schema Property with the given field descriptor
func NewProperty(descriptor *desc.FieldDescriptor) schema.Property {
	return &property{
		desc:    descriptor,
		options: make(schema.Options),
	}
}

type property struct {
	desc    *desc.FieldDescriptor
	options schema.Options
}

// GetName returns the property name
func (property *property) GetName() string {
	return property.desc.GetName()
}

// GetPosition returns the property position inside a message
func (property *property) GetPosition() int32 {
	return property.desc.GetNumber()
}

// GetType returns the property type
func (property *property) GetType() types.Type {
	return Types[property.desc.GetType()]
}

// GetLabel returns the property label
func (property *property) GetLabel() types.Label {
	return Labels[property.desc.GetLabel()]
}

// GetNested attempts to return a all the nested properties
func (property *property) GetNested() map[string]schema.Property {
	if property.desc.GetType() != descriptor.FieldDescriptorProto_TYPE_MESSAGE {
		return make(map[string]schema.Property)
	}

	fields := property.desc.GetMessageType().GetFields()
	result := make(map[string]schema.Property, len(fields))
	for _, field := range fields {
		result[field.GetName()] = NewProperty(field)
	}

	return result
}

func (property *property) GetOptions() schema.Options {
	return property.options
}

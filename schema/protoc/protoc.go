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

// Collection represents a collection of proto schemas
type Collection interface {
	schema.Collection
	GetDescriptors() []*desc.FileDescriptor
}

// NewCollection constructs a new schema collection from the given descriptors
func NewCollection(descriptors []*desc.FileDescriptor) Collection {
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

func (collection *collection) GetProperty(message string) schema.Property {
	for _, descriptor := range collection.descriptors {
		message := descriptor.FindMessage(message)
		if message == nil {
			continue
		}

		return NewMessageProperty(message)
	}

	return nil
}

func (collection *collection) GetDescriptors() []*desc.FileDescriptor {
	return collection.descriptors
}

// Service represents a proto service
type Service interface {
	schema.Service
	GetDescriptor() *desc.ServiceDescriptor
}

// NewService constructs a new service with the given descriptor
func NewService(descriptor *desc.ServiceDescriptor) Service {
	return &service{
		descriptor: descriptor,
		options:    make(schema.Options),
	}
}

type service struct {
	descriptor *desc.ServiceDescriptor
	options    schema.Options
}

func (service *service) GetName() string {
	return service.descriptor.GetName()
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

func (service *service) GetMethods() []schema.Method {
	result := make([]schema.Method, len(service.descriptor.GetMethods()))
	for index, method := range service.descriptor.GetMethods() {
		result[index] = NewMethod(method)
	}

	return result
}

func (service *service) GetDescriptor() *desc.ServiceDescriptor {
	return service.descriptor
}

func (service *service) GetOptions() schema.Options {
	return service.options
}

// Method represents a proto service method
type Method interface {
	schema.Method
	GetDescriptor() *desc.MethodDescriptor
}

// NewMethod constructs a new method with the given descriptor
func NewMethod(descriptor *desc.MethodDescriptor) Method {
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
	return NewMessageProperty(method.descriptor.GetInputType())
}

func (method *method) GetOutput() schema.Property {
	return NewMessageProperty(method.descriptor.GetOutputType())
}

func (method *method) GetDescriptor() *desc.MethodDescriptor {
	return method.descriptor
}

func (method *method) GetOptions() schema.Options {
	return method.options
}

// NewMessageProperty constructs a schema Property with the given message descriptor
func NewMessageProperty(descriptor *desc.MessageDescriptor) Property {
	return &property{
		message: descriptor,
		options: make(schema.Options),
	}
}

// NewFieldProperty constructs a schema Property with the given field descriptor
func NewFieldProperty(desc *desc.FieldDescriptor) Property {
	result := &property{
		field:   desc,
		options: make(schema.Options),
	}

	if desc.GetType() == descriptor.FieldDescriptorProto_TYPE_MESSAGE {
		result.field = nil
		result.message = desc.GetMessageType()
	}

	return result
}

// Property represents a proto property
type Property interface {
	schema.Property
	GetProtoField(string) Property
	GetFieldDescriptor() *desc.FieldDescriptor
	GetMessageDescriptor() *desc.MessageDescriptor
}

type property struct {
	message *desc.MessageDescriptor
	field   *desc.FieldDescriptor
	options schema.Options
}

// GetName returns the property name
func (property *property) GetName() string {
	if property.message != nil {
		return property.message.GetName()
	}

	if property.field != nil {
		return property.field.GetName()
	}

	return ""
}

// GetType returns the property type
func (property *property) GetType() types.Type {
	if property.message != nil {
		return types.TypeMessage
	}

	if property.field != nil {
		return Types[property.field.GetType()]
	}

	return ""
}

// GetProtoField attempts to return a proto field matching the given name
func (property *property) GetProtoField(name string) Property {
	if property.message == nil {
		return nil
	}

	for _, field := range property.message.GetFields() {
		if field.GetName() == name {
			return NewFieldProperty(field)
		}
	}

	return nil
}

// GetLabel returns the property label
func (property *property) GetLabel() types.Label {
	if property.message != nil {
		return types.LabelOptional
	}

	if property.field != nil {
		return Labels[property.field.GetLabel()]
	}

	return ""
}

// GetNested attempts to return a all the nested properties
func (property *property) GetNested() map[string]schema.Property {
	if property.message == nil {
		return make(map[string]schema.Property)
	}

	result := make(map[string]schema.Property, len(property.message.GetFields()))
	for _, field := range property.message.GetFields() {
		result[field.GetName()] = NewFieldProperty(field)
	}

	return nil
}

func (property *property) GetMessageDescriptor() *desc.MessageDescriptor {
	return property.message
}

func (property *property) GetFieldDescriptor() *desc.FieldDescriptor {
	return property.field
}

func (property *property) GetOptions() schema.Options {
	return property.options
}

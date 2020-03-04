package protoc

import (
	"github.com/golang/protobuf/proto"
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

func (method *method) GetInput() schema.Object {
	return NewObject(method.descriptor.GetInputType())
}

func (method *method) GetOutput() schema.Object {
	return NewObject(method.descriptor.GetOutputType())
}

func (method *method) GetDescriptor() *desc.MethodDescriptor {
	return method.descriptor
}

func (method *method) GetOptions() schema.Options {
	return method.options
}

// NewObject constructs a schema Object with the given descriptor
func NewObject(descriptor *desc.MessageDescriptor) Object {
	return &object{
		descriptor: descriptor,
		options:    make(schema.Options),
	}
}

// Object represents a proto message
type Object interface {
	schema.Object
	GetProtoField(name string) Field
	GetDescriptor() *desc.MessageDescriptor
}

type object struct {
	descriptor *desc.MessageDescriptor
	options    schema.Options
}

// GetField attempts to return a field matching the given name
func (object *object) GetField(name string) schema.Field {
	for _, field := range object.descriptor.GetFields() {
		if field.GetName() == name {
			return NewField(field)
		}
	}

	return nil
}

// GetProtoField attempts to return a proto field matching the given name
func (object *object) GetProtoField(name string) Field {
	for _, field := range object.descriptor.GetFields() {
		if field.GetName() == name {
			return NewField(field)
		}
	}

	return nil
}

// GetFields returns all available fields inside the given object
func (object *object) GetFields() []schema.Field {
	result := make([]schema.Field, len(object.descriptor.GetFields()))

	for index, field := range object.descriptor.GetFields() {
		result[index] = NewField(field)
	}

	return result
}

func (object *object) GetDescriptor() *desc.MessageDescriptor {
	return object.descriptor
}

func (object *object) GetOptions() schema.Options {
	return object.options
}

// NewField constructs a new object field with the given descriptor
func NewField(descriptor *desc.FieldDescriptor) Field {
	return &field{
		descriptor: descriptor,
		options:    make(schema.Options),
	}
}

// Field represents a proto message field
type Field interface {
	schema.Field
	GetDescriptor() *desc.FieldDescriptor
}

type field struct {
	descriptor *desc.FieldDescriptor
	options    schema.Options
}

func (field *field) GetName() string {
	return field.descriptor.GetName()
}

func (field *field) GetType() types.Type {
	return Types[field.descriptor.GetType()]
}

func (field *field) GetLabel() types.Label {
	return Labels[field.descriptor.GetLabel()]
}

func (field *field) GetObject() schema.Object {
	return NewObject(field.descriptor.GetMessageType())
}

func (field *field) GetDescriptor() *desc.FieldDescriptor {
	return field.descriptor
}

func (field *field) GetOptions() schema.Options {
	return field.options
}

package mock

import (
	"github.com/jexia/maestro/schema"
	"github.com/jexia/maestro/specs/labels"
	"github.com/jexia/maestro/specs/types"
)

// NewCollection constructs a new schema collection from the given descriptors
func NewCollection(descriptor Collection) *Collection {
	return &descriptor
}

// Exception represents a exception thrown during runtime
type Exception struct {
	File    string `yaml:"file"`
	Line    int    `yaml:"line"`
	Message string `yaml:"message"`
}

// Collection represents a mock YAML file
type Collection struct {
	Exception Exception            `yaml:"exception"`
	Services  map[string]*Service  `yaml:"services"`
	Objects   map[string]*Property `yaml:"objects"`
}

// GetService attempts to find the given service
func (collection *Collection) GetService(name string) schema.Service {
	for key, service := range collection.Services {
		if key != name {
			continue
		}

		return NewService(name, service)
	}

	return nil
}

// GetServices returns all available services inside the given collection
func (collection *Collection) GetServices() []schema.Service {
	result := make([]schema.Service, len(collection.Services))

	for name, service := range collection.Services {
		result = append(result, NewService(name, service))
	}

	return result
}

// GetMessage attempts to find and return the given schema property
func (collection *Collection) GetMessage(name string) schema.Property {
	for key, object := range collection.Objects {
		if key != name {
			continue
		}

		return NewProperty(key, object)
	}

	return nil
}

// GetMessages returns all available messages inside the given collection
func (collection *Collection) GetMessages() []schema.Property {
	result := make([]schema.Property, len(collection.Objects))

	for key, object := range collection.Objects {
		result = append(result, NewProperty(key, object))
	}

	return result
}

// NewService constructs a new service with the given descriptor
func NewService(name string, service *Service) *Service {
	service.Name = name
	return service
}

// Service represents a mocking service
type Service struct {
	Name      string
	Package   string             `yaml:"package"`
	Comment   string             `yaml:"comment"`
	Host      string             `yaml:"host"`
	Transport string             `yaml:"transport"`
	Codec     string             `yaml:"codec"`
	Methods   map[string]*Method `yaml:"methods"`
	Options   schema.Options     `yaml:"options"`
}

// GetPackage returns the service package
func (service *Service) GetPackage() string {
	return service.Package
}

// GetFullyQualifiedName returns the fully qualified service name
func (service *Service) GetFullyQualifiedName() string {
	return service.Name
}

// GetName returns the service name
func (service *Service) GetName() string {
	return service.Name
}

// GetComment returns the service comment
func (service *Service) GetComment() string {
	return service.Comment
}

// GetHost returns the service host
func (service *Service) GetHost() string {
	return service.Host
}

// GetTransport returns the service transport
func (service *Service) GetTransport() string {
	return service.Transport
}

// GetCodec returns the service codec
func (service *Service) GetCodec() string {
	return service.Codec
}

// GetMethod attempts to return the given service method
func (service *Service) GetMethod(name string) schema.Method {
	for key, method := range service.Methods {
		if key != name {
			continue
		}

		return NewMethod(key, method)
	}

	return nil
}

// GetOptions returns the service options
func (service *Service) GetOptions() schema.Options {
	return service.Options
}

// GetMethods attempts to return the given service methods
func (service *Service) GetMethods() schema.Methods {
	result := make([]schema.Method, len(service.Methods))

	index := 0
	for key, method := range service.Methods {
		result[index] = NewMethod(key, method)
		index++
	}

	return result
}

// Method represents a mock YAML service method
type Method struct {
	Name    string
	Comment string         `yaml:"comment"`
	Input   *Property      `yaml:"input"`
	Output  *Property      `yaml:"output"`
	Options schema.Options `yaml:"options"`
}

// NewMethod constructs a new method with the given descriptor
func NewMethod(name string, method *Method) *Method {
	method.Name = name
	return method
}

// GetName returns the method name
func (method *Method) GetName() string {
	return method.Name
}

// GetComment returns the method comment
func (method *Method) GetComment() string {
	return method.Comment
}

// GetInput returns the method input
func (method *Method) GetInput() schema.Property {
	return method.Input
}

// GetOutput returns the method output
func (method *Method) GetOutput() schema.Property {
	return method.Output
}

// GetOptions returns the method options
func (method *Method) GetOptions() schema.Options {
	return method.Options
}

// NewProperty appends the name to the given property
func NewProperty(name string, property *Property) *Property {
	property.Name = name

	for key, prop := range property.Nested {
		property.Nested[key] = NewProperty(key, prop)
	}

	return property
}

// Property represents a proto message property
type Property struct {
	Name     string
	Comment  string               `yaml:"comment"`
	Type     types.Type           `yaml:"type"`
	Label    labels.Label         `yaml:"label"`
	Position int32                `yaml:"position"`
	Nested   map[string]*Property `yaml:"nested"`
	Options  schema.Options       `yaml:"options"`
}

// GetName returns the field name
func (property *Property) GetName() string {
	return property.Name
}

// GetComment returns the field comment
func (property *Property) GetComment() string {
	return property.Comment
}

// GetPosition returns the field position
func (property *Property) GetPosition() int32 {
	return property.Position
}

// GetType returns tye field type
func (property *Property) GetType() types.Type {
	return property.Type
}

// GetLabel returns the field label
func (property *Property) GetLabel() labels.Label {
	return property.Label
}

// GetNested returns the field nested object
func (property *Property) GetNested() map[string]schema.Property {
	result := make(map[string]schema.Property, len(property.Nested))
	for key, nested := range property.Nested {
		result[key] = NewProperty(key, nested)
	}

	return result
}

// GetOptions returns the field options
func (property *Property) GetOptions() schema.Options {
	return property.Options
}

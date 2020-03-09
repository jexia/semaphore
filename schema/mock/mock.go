package mock

import (
	"github.com/jexia/maestro/schema"
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

// GetProperty attempts to find and return the given schema property
func (collection *Collection) GetProperty(name string) schema.Property {
	for key, object := range collection.Objects {
		if key != name {
			continue
		}

		return object
	}

	return nil
}

// NewService constructs a new service with the given descriptor
func NewService(name string, service *Service) *Service {
	service.Name = name
	return service
}

// Service represents a mocking service
type Service struct {
	Name    string
	Methods map[string]*Method `yaml:"methods"`
	Options schema.Options     `yaml:"options"`
}

// GetName returns the service name
func (service *Service) GetName() string {
	return service.Name
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
func (service *Service) GetMethods() []schema.Method {
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

// Property represents a proto message property
type Property struct {
	Name     string               `yaml:"name"`
	Type     types.Type           `yaml:"type"`
	Label    types.Label          `yaml:"label"`
	Position int32                `yaml:"position"`
	Nested   map[string]*Property `yaml:"nested"`
	Options  schema.Options       `yaml:"options"`
}

// GetName returns the field name
func (property *Property) GetName() string {
	return property.Name
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
func (property *Property) GetLabel() types.Label {
	return property.Label
}

// GetNested returns the field nested object
func (property *Property) GetNested() map[string]schema.Property {
	result := make(map[string]schema.Property, len(property.Nested))
	for key, nested := range property.Nested {
		result[key] = nested
	}

	return result
}

// GetOptions returns the field options
func (property *Property) GetOptions() schema.Options {
	return property.Options
}

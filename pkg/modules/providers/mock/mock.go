package mock

import "github.com/jexia/semaphore/pkg/specs"

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
	Exception  Exception                  `yaml:"exception"`
	Services   map[string]*Service        `yaml:"services"`
	Properties map[string]*specs.Property `yaml:"properties"`
}

// GetService attempts to find the given service
func (collection *Collection) GetService(name string) *Service {
	for key, service := range collection.Services {
		if key != name {
			continue
		}

		return NewService(name, service)
	}

	return nil
}

// GetServices returns all available services inside the given collection
func (collection *Collection) GetServices() []*Service {
	result := make([]*Service, 0, len(collection.Services))

	for name, service := range collection.Services {
		result = append(result, NewService(name, service))
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
	Options   specs.Options      `yaml:"options"`
}

// GetMethod attempts to return the given service method
func (service *Service) GetMethod(name string) *Method {
	for key, method := range service.Methods {
		if key != name {
			continue
		}

		return NewMethod(key, method)
	}

	return nil
}

// GetMethods attempts to return the given service methods
func (service *Service) GetMethods() []*Method {
	result := make([]*Method, len(service.Methods))

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
	Comment string        `yaml:"comment"`
	Input   string        `yaml:"input"`
	Output  string        `yaml:"output"`
	Options specs.Options `yaml:"options"`
}

// NewMethod constructs a new method with the given descriptor
func NewMethod(name string, method *Method) *Method {
	method.Name = name
	return method
}

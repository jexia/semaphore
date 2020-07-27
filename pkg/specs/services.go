package specs

import "github.com/jexia/semaphore/pkg/specs/metadata"

// ServiceList represents a collection of services
type ServiceList []*Service

// Append appends the given services to the current services collection
func (services *ServiceList) Append(list ServiceList) {
	*services = append(*services, list...)
}

// Get attempts to find and return a flow with the given name
func (services ServiceList) Get(name string) *Service {
	for _, service := range services {
		if service.FullyQualifiedName == name {
			return service
		}
	}

	return nil
}

// Service represents a service which exposes a set of methods
type Service struct {
	*metadata.Meta
	Comment            string    `json:"comment,omitempty"`
	Package            string    `json:"package,omitempty"`
	FullyQualifiedName string    `json:"fully_qualified_name,omitempty"`
	Name               string    `json:"name,omitempty"`
	Transport          string    `json:"transport,omitempty"`
	Codec              string    `json:"codec,omitempty"`
	Host               string    `json:"host,omitempty"`
	Methods            []*Method `json:"methods,omitempty"`
	Options            Options   `json:"options,omitempty"`
}

// GetMethod attempts to find and return a method matching the given name
func (service *Service) GetMethod(name string) *Method {
	for _, method := range service.Methods {
		if method.Name == name {
			return method
		}
	}

	return nil
}

// Method represents a service method
type Method struct {
	*metadata.Meta
	Comment string  `json:"comment,omitempty"`
	Name    string  `json:"name,omitempty"`
	Input   string  `json:"input,omitempty"`
	Output  string  `json:"output,omitempty"`
	Options Options `json:"options,omitempty"`
}

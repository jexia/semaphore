package specs

// NewServicesManifest constructs a new service(s) manifest
func NewServicesManifest() *ServicesManifest {
	return &ServicesManifest{
		Services: make([]*Service, 0),
	}
}

// ServicesManifest holds a collection of services
type ServicesManifest struct {
	Services []*Service `json:"services"`
}

// Merge merges the incoming services into the given service manifest
func (services *ServicesManifest) Merge(incoming *ServicesManifest) {
	services.Services = append(services.Services, incoming.Services...)
}

// GetService attempts to find and return a flow with the given name
func (services *ServicesManifest) GetService(name string) *Service {
	for _, service := range services.Services {
		if service.FullyQualifiedName == name {
			return service
		}
	}

	return nil
}

// Service represents a service which exposes a set of methods
type Service struct {
	Comment            string    `json:"comment"`
	Package            string    `json:"package"`
	FullyQualifiedName string    `json:"fully_qualified_name"`
	Name               string    `json:"name"`
	Transport          string    `json:"transport"`
	Codec              string    `json:"codec"`
	Host               string    `json:"host"`
	Methods            []*Method `json:"methods"`
	Options            Options   `json:"options"`
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
	Comment string  `json:"comment"`
	Name    string  `json:"name"`
	Input   string  `json:"input"`
	Output  string  `json:"output"`
	Options Options `json:"options"`
}

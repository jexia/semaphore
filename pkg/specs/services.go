package specs

// ServicesManifest holds a collection of services
type ServicesManifest struct {
	Services []*Service `json:"services,omitempty"`
}

// Append merges the incoming services into the given service manifest
func (manifest *ServicesManifest) Append(incoming ...*ServicesManifest) {
	if manifest == nil {
		return
	}

	for _, right := range incoming {
		manifest.Services = append(manifest.Services, right.Services...)
	}
}

// GetService attempts to find and return a flow with the given name
func (manifest *ServicesManifest) GetService(name string) *Service {
	for _, service := range manifest.Services {
		if service.FullyQualifiedName == name {
			return service
		}
	}

	return nil
}

// Service represents a service which exposes a set of methods
type Service struct {
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
	Comment string  `json:"comment,omitempty"`
	Name    string  `json:"name,omitempty"`
	Input   string  `json:"input,omitempty"`
	Output  string  `json:"output,omitempty"`
	Options Options `json:"options,omitempty"`
}

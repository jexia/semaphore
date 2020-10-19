package mock

import "github.com/jexia/semaphore/pkg/specs"

// ServiceManifest formats the given mock collection to a specs service(s) manifest
func ServiceManifest(collection *Collection) specs.ServiceList {
	result := make(specs.ServiceList, len(collection.GetServices()))

	for index, service := range collection.GetServices() {
		result[index] = SpecsService(service)
	}

	return result
}

// SpecsService formats the given mock service to a specs service
func SpecsService(service *Service) *specs.Service {
	result := &specs.Service{
		Package:            service.Package,
		FullyQualifiedName: service.Name,
		Name:               service.Name,
		Comment:            service.Comment,
		RequestCodec:       service.Codec,
		ResponseCodec:      service.Codec,
		Host:               service.Host,
		Options:            service.Options,
		Methods:            make([]*specs.Method, len(service.GetMethods())),
	}

	for index, method := range service.GetMethods() {
		result.Methods[index] = SpecsMethod(method)
	}

	return result
}

// SpecsMethod formats the given mock method to a specs method
func SpecsMethod(method *Method) *specs.Method {
	result := &specs.Method{
		Name:    method.Name,
		Comment: method.Comment,
		Input:   method.Input,
		Output:  method.Output,
		Options: method.Options,
	}

	return result
}

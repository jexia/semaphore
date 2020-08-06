package mock

import (
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/template"
)

// SchemaManifest formats the given mock collection to a specs schema manifest
func SchemaManifest(collection *Collection) specs.Schemas {
	result := make(specs.Schemas, 0)

	for _, prop := range collection.GetMessages() {
		result[prop.Name] = SpecsProperty("", prop)
	}

	return result
}

// SpecsProperty formats the given mock property to a specs property
func SpecsProperty(path string, property *Property) *specs.Property {
	result := &specs.Property{
		Name:     property.Name,
		Path:     path,
		Comment:  property.Comment,
		Default:  property.Default,
		Type:     property.Type,
		Label:    property.Label,
		Position: property.Position,
		Options:  property.Options,
	}

	if property.Enum != nil {
		result.Enum = &specs.Enum{
			Name:      property.Name,
			Keys:      make(map[string]*specs.EnumValue, len(property.Enum)),
			Positions: make(map[int32]*specs.EnumValue, len(property.Enum)),
		}

		for key, value := range property.Enum {
			value := &specs.EnumValue{
				Key:         key,
				Position:    value.Position,
				Description: value.Description,
			}

			result.Enum.Keys[value.Key] = value
			result.Enum.Positions[value.Position] = value
		}
	}

	if property.Nested != nil {
		result.Nested = make(map[string]*specs.Property, len(property.Nested))

		for key, nested := range property.GetNested() {
			result.Nested[key] = SpecsProperty(template.JoinPath(path, key), nested)
		}
	}

	return result
}

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
		Codec:              service.Codec,
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

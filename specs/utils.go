package specs

import (
	"github.com/jexia/maestro/schema"
)

// ToParameterMap maps the given schema object to a parameter map
func ToParameterMap(origin *ParameterMap, path string, prop schema.Property) *ParameterMap {
	result := &ParameterMap{}

	if origin != nil {
		result.Options = origin.Options
		result.Header = origin.Header
		result.Schema = origin.Schema
	}

	if prop != nil {
		result.Property = ToProperty("", "", prop)
	}

	return result
}

// ToProperty transforms the given schema property to a specs property
func ToProperty(path string, name string, prop schema.Property) *Property {
	result := &Property{
		Path:  path,
		Name:  name,
		Type:  prop.GetType(),
		Label: prop.GetLabel(),
	}

	if prop.GetNested() != nil {
		result.Nested = make(map[string]*Property, len(prop.GetNested()))
		for key, object := range prop.GetNested() {
			result.Nested[key] = ToProperty(JoinPath(path, key), key, object)
		}
	}

	return result
}

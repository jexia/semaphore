package specs

import (
	"github.com/jexia/maestro/schema"
	"github.com/jexia/maestro/specs/types"
)

// ToParameterMap maps the given schema object to a parameter map
func ToParameterMap(origin *ParameterMap, path string, object schema.Object) *ParameterMap {
	result := &ParameterMap{
		Properties: make(map[string]*Property),
		Nested:     make(map[string]*NestedParameterMap),
		Repeated:   make(map[string]*RepeatedParameterMap),
	}

	if origin != nil {
		result.Options = origin.Options
		result.Header = origin.Header
		result.Schema = origin.Schema
	}

	if object == nil {
		return result
	}

	for _, field := range object.GetFields() {
		path := JoinPath(path, field.GetName())

		if field.GetLabel() == types.LabelRepeated {
			param := ToParameterMap(origin, path, field.GetObject())
			result.Repeated[field.GetName()] = &RepeatedParameterMap{
				Path:       path,
				Name:       field.GetName(),
				Properties: param.Properties,
				Nested:     param.Nested,
				Repeated:   param.Repeated,
			}
			continue
		}

		if field.GetType() == types.TypeMessage {
			param := ToParameterMap(origin, path, field.GetObject())
			result.Nested[field.GetName()] = &NestedParameterMap{
				Path:       path,
				Name:       field.GetName(),
				Properties: param.Properties,
				Nested:     param.Nested,
				Repeated:   param.Repeated,
			}
			continue
		}

		result.Properties[field.GetName()] = &Property{
			Path: path,
			Type: field.GetType(),
		}
	}

	return result
}

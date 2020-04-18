package references

import (
	"github.com/jexia/maestro/pkg/specs"
)

// ToParameterMap maps the given schema object to a parameter map
func ToParameterMap(origin *specs.ParameterMap, path string, prop *specs.Property) *specs.ParameterMap {
	result := &specs.ParameterMap{}

	if origin != nil {
		result.Options = origin.Options
		result.Header = origin.Header
		result.Schema = origin.Schema
	}

	if prop != nil {
		result.Property = prop
	}

	return result
}

package refs

import (
	"github.com/jexia/maestro/specs"
)

// References returns all the available references inside the given object
func References(params *specs.ParameterMap) map[string]*specs.PropertyReference {
	result := make(map[string]*specs.PropertyReference)
	for _, prop := range params.Header {
		if prop.Reference != nil {
			result[prop.Reference.String()] = prop.Reference
		}
	}

	for key, prop := range PropertyReferences(params.Property) {
		result[key] = prop
	}

	return result
}

// PropertyReferences returns the available references within the given property
func PropertyReferences(property *specs.Property) map[string]*specs.PropertyReference {
	result := make(map[string]*specs.PropertyReference)

	if property.Reference != nil {
		result[property.Reference.String()] = property.Reference
	}

	if property.Nested == nil {
		return result
	}

	for _, nested := range property.Nested {
		for key, ref := range PropertyReferences(nested) {
			result[key] = ref
		}
	}

	return result
}

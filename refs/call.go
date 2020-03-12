package refs

import (
	"github.com/jexia/maestro/specs"
)

// References represents a map of property references
type References map[string]*specs.PropertyReference

// MergeLeft merges the references into the given reference
func (references References) MergeLeft(incoming ...References) {
	for _, refs := range incoming {
		for key, val := range refs {
			references[key] = val
		}
	}
}

// ParameterReferences returns all the available references inside the given parameter map
func ParameterReferences(params *specs.ParameterMap) References {
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
func PropertyReferences(property *specs.Property) References {
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

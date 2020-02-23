package refs

import "github.com/jexia/maestro/specs"

// References returns all the available references inside the given object
func References(object specs.Object) (result map[string]*specs.PropertyReference) {
	for _, prop := range object.GetProperties() {
		if prop.Reference != nil {
			result[prop.Reference.String()] = prop.Reference
		}
	}

	for _, nested := range object.GetNestedProperties() {
		for key, val := range References(nested) {
			result[key] = val
		}
	}

	for _, repeated := range object.GetRepeatedProperties() {
		for key, val := range References(repeated) {
			result[key] = val
		}
	}

	return result
}

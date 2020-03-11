package graphql

import (
	"github.com/jexia/maestro/refs"
	"github.com/jexia/maestro/specs"
	"github.com/jexia/maestro/specs/types"
)

// ResponseValue constructs the response value send back to the client
func ResponseValue(specs *specs.Property, refs *refs.Store) (interface{}, error) {
	if specs.Type != types.TypeMessage {
		return nil, ErrInvalidObject
	}

	result := make(map[string]interface{}, len(specs.Nested))
	for _, nested := range specs.Nested {
		if nested.Type == types.TypeMessage {
			value, err := ResponseValue(nested, refs)
			if err != nil {
				return nil, err
			}

			result[nested.Name] = value
			continue
		}

		ref := refs.Load(nested.Reference.Resource, nested.Reference.Path)
		if ref == nil {
			continue
		}

		result[nested.Name] = ref.Value
	}

	return result, nil
}

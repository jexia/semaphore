package graphql

import (
	"github.com/jexia/maestro/pkg/refs"
	"github.com/jexia/maestro/pkg/specs"
	"github.com/jexia/maestro/pkg/specs/labels"
	"github.com/jexia/maestro/pkg/specs/types"
)

// ResponseValue constructs the response value send back to the client
func ResponseValue(specs *specs.Property, refs refs.Store) (interface{}, error) {
	if specs.Type != types.Message {
		return nil, ErrInvalidObject
	}

	result := make(map[string]interface{}, len(specs.Nested))
	for _, nested := range specs.Nested {
		if nested.Label == labels.Repeated {
			store := refs.Load(nested.Reference.Resource, nested.Reference.Path)
			repeating := make([]interface{}, len(store.Repeated))

			for index, store := range store.Repeated {
				if nested.Type == types.Message {
					value, err := ResponseValue(nested, store)
					if err != nil {
						return nil, err
					}

					repeating[index] = value
					continue
				}

				ref := store.Load("", "")
				if ref == nil {
					continue
				}

				repeating[index] = ref.Value
			}

			result[nested.Name] = repeating
			continue
		}

		if nested.Type == types.Message {
			value, err := ResponseValue(nested, refs)
			if err != nil {
				return nil, err
			}

			result[nested.Name] = value
			continue
		}

		if nested.Reference == nil {
			continue
		}

		ref := refs.Load(nested.Reference.Resource, nested.Reference.Path)
		if ref == nil {
			continue
		}

		if ref.Enum != nil && nested.Enum != nil {
			// Enum should exist and a additional nil check should not be necessary
			result[nested.Name] = nested.Enum.Positions[*ref.Enum].Key
			continue
		}

		result[nested.Name] = ref.Value
	}

	return result, nil
}

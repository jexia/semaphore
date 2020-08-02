package graphql

import (
	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/labels"
	"github.com/jexia/semaphore/pkg/specs/types"
)

// ResponseValue constructs the response value send back to the client
func ResponseValue(specs *specs.Property, refs references.Store) (interface{}, error) {
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

				val := nested.Default
				ref := store.Load("", "")
				if ref != nil {
					if ref.Enum != nil {
						val = nested.Enum.Positions[*ref.Enum].Key
					}

					if ref.Value != nil && val == nil {
						val = ref.Value
					}
				}

				repeating[index] = val
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

		val := nested.Default

		if nested.Reference != nil {
			ref := refs.Load(nested.Reference.Resource, nested.Reference.Path)
			if ref != nil {
				if ref.Enum != nil {
					// Enum should exist and a additional nil check should not be necessary
					val = nested.Enum.Positions[*ref.Enum].Key
				}

				if ref.Value != nil && val == nil {
					val = ref.Value
				}
			}
		}

		result[nested.Name] = val
	}

	return result, nil
}

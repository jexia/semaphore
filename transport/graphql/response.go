package graphql

import (
	"github.com/jexia/maestro/specs"
	"github.com/jexia/maestro/specs/labels"
	"github.com/jexia/maestro/specs/types"
)

// ResponseValue constructs the response value send back to the client
func ResponseValue(specs *specs.Property, refs *specs.Store) (interface{}, error) {
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

				// TODO: support repeated types
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

		result[nested.Name] = ref.Value
	}

	return result, nil
}

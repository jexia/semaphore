package graphql

import (
	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/types"
)

// ResponseObject constructs the response value send back to the client
func ResponseObject(specs *specs.Property, refs references.Store) (map[string]interface{}, error) {
	if specs == nil || refs == nil {
		return nil, ErrInvalidObject
	}

	if specs.Type() != types.Message {
		return nil, ErrInvalidObject
	}

	value, err := ResponseValue(specs.Template, refs)
	if err != nil {
		return nil, err
	}

	// could safely assume that the given type is a map[string]interface{}
	object := value.(map[string]interface{})
	return object, nil
}

// ResponseValue constructs a new response used inside a response object
func ResponseValue(template specs.Template, refs references.Store) (interface{}, error) {
	switch {
	case template.Message != nil:
		result := make(map[string]interface{}, len(template.Message))
		for _, nested := range template.Message {
			value, err := ResponseValue(nested.Template, refs)
			if err != nil {
				return nil, err
			}

			if value == nil {
				continue
			}

			result[nested.Name] = value
		}

		return result, nil
	case template.Repeated != nil:
		// TODO: implement default types
		store := refs.Load(template.Reference.Resource, template.Reference.Path)
		if store == nil {
			return nil, nil
		}

		result := make([]interface{}, 0, len(store.Repeated))

		if len(store.Repeated) == 0 {
			return result, nil
		}

		// TODO: cache repeated template
		template, err := template.Repeated.Template()
		if err != nil {
			return nil, err
		}

		for _, store := range store.Repeated {
			value, err := ResponseValue(template, store)
			if err != nil {
				return nil, err
			}

			if value == nil {
				continue
			}

			result = append(result, value)
		}

		return result, nil
	case template.Enum != nil:
		if template.Reference == nil {
			return nil, nil
		}

		ref := refs.Load(template.Reference.Resource, template.Reference.Path)
		if ref == nil {
			return nil, nil
		}

		return template.Enum.Positions[*ref.Enum].Key, nil
	case template.Scalar != nil:
		value := template.Scalar.Default

		if template.Reference != nil {
			ref := refs.Load(template.Reference.Resource, template.Reference.Path)
			if ref != nil {
				value = ref.Value
			}
		}

		return value, nil
	}

	return nil, nil
}

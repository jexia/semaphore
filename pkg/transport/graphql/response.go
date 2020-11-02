package graphql

import (
	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/types"
)

// ResponseObject constructs the response value send back to the client
func ResponseObject(specs *specs.Property, store references.Store, tracker references.Tracker) (map[string]interface{}, error) {
	if specs == nil || store == nil {
		return nil, ErrInvalidObject
	}

	if specs.Type() != types.Message {
		return nil, ErrInvalidObject
	}

	value, err := ResponseValue(specs.Template, store, tracker)
	if err != nil {
		return nil, err
	}

	// could safely assume that the given type is a map[string]interface{}
	object := value.(map[string]interface{})
	return object, nil
}

// ResponseValue constructs a new response used inside a response object
func ResponseValue(tmpl *specs.Template, store references.Store, tracker references.Tracker) (interface{}, error) {
	switch {
	case tmpl.Message != nil:
		result := make(map[string]interface{}, len(tmpl.Message))
		for _, nested := range tmpl.Message {
			value, err := ResponseValue(nested.Template, store, tracker)
			if err != nil {
				return nil, err
			}

			if value == nil {
				continue
			}

			result[nested.Name] = value
		}

		return result, nil
	case tmpl.Repeated != nil:
		// TODO: implement default types

		rpath := tmpl.Reference.String()
		length := store.Length(tracker.Resolve(rpath))

		if length == 0 {
			return nil, nil
		}

		result := make([]interface{}, 0, length)
		template, err := tmpl.Repeated.Template()
		if err != nil {
			return nil, err
		}

		tracker.Track(rpath, 0)

		for index := 0; index < length; index++ {
			value, err := ResponseValue(template, store, tracker)
			if err != nil {
				return nil, err
			}

			if value == nil {
				continue
			}

			result = append(result, value)
			tracker.Next(rpath)
		}

		return result, nil
	case tmpl.Enum != nil:
		if tmpl.Reference == nil {
			return nil, nil
		}

		ref := store.Load(tracker.Resolve(tmpl.Reference.String()))
		if ref == nil {
			return nil, nil
		}

		return tmpl.Enum.Positions[*ref.Enum].Key, nil
	case tmpl.Scalar != nil:
		value := tmpl.Scalar.Default

		if tmpl.Reference != nil {
			ref := store.Load(tracker.Resolve(tmpl.Reference.String()))
			if ref != nil {
				value = ref.Value
			}
		}

		return value, nil
	}

	return nil, nil
}

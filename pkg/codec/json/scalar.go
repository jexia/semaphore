package json

import (
	"github.com/francoispqt/gojay"
	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
)

// Scalar represents a scalar type such as a string, int or float. The scalar
// type is used to encode or decode values inside gojay. These values are used
// to construct messages to a service or user.
type Scalar specs.Template

func (template Scalar) value(store references.Store, tracker references.Tracker) interface{} {
	var value = template.Scalar.Default

	if template.Reference != nil {
		var reference = store.Load(tracker.Resolve(template.Reference.String()))
		if reference != nil && reference.Value != nil {
			value = reference.Value
		}
	}

	return value
}

// Marshal marshals the scalar template as a JSON value.
func (template Scalar) Marshal(encoder *gojay.Encoder, store references.Store, tracker references.Tracker) {
	if template.Scalar == nil {
		return
	}

	AddType(encoder, template.Scalar.Type, template.value(store, tracker))
}

// MarshalKey marshals the scalar template as an object field using the given key.
// The key is not set if the value is `null`.
func (template Scalar) MarshalKey(encoder *gojay.Encoder, key string, store references.Store, tracker references.Tracker) {
	if template.Scalar == nil {
		return
	}

	value := template.value(store, tracker)
	if value == nil {
		return
	}

	AddTypeKey(encoder, key, template.Scalar.Type, value)
}

// Unmarshal attempts to unmarshal the value from the decoder as a scalar and
// stores it inside the reference store.
func (template Scalar) Unmarshal(decoder *gojay.Decoder, path string, store references.Store, tracker references.Tracker) error {
	if template.Scalar == nil {
		return nil
	}

	value, err := DecodeType(decoder, template.Scalar.Type)
	if err != nil {
		return err
	}

	store.Store(tracker.Resolve(path), &references.Reference{
		Value: value,
	})

	return nil
}

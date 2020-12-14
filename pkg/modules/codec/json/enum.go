package json

import (
	"github.com/francoispqt/gojay"
	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/types"
)

// Enum represents a enum type template. This type is used to encode or decode
// values inside gojay. These values are used to construct messages to a service
// or a user.
type Enum specs.Template

func (template *Enum) value(store references.Store, tracker references.Tracker) *specs.EnumValue {
	if template.Reference == nil {
		return nil
	}

	var reference = store.Load(tracker.Resolve(template.Reference.String()))
	if reference == nil || reference.Enum == nil {
		return nil
	}

	if position := reference.Enum; position != nil {
		return template.Enum.Positions[*reference.Enum]
	}

	return nil
}

// Marshal marshals the enum template as a JSON value. If no enum has been defined
// inside the given reference store is the type ignored.
func (template Enum) Marshal(encoder *gojay.Encoder, store references.Store, tracker references.Tracker) {
	value := template.value(store, tracker)
	if value == nil {
		return
	}

	AddType(encoder, types.String, value.Key)
}

// MarshalKey marshals the enum template as an object field using the given key.
// The key is not set if the enum value is `null`.
func (template Enum) MarshalKey(encoder *gojay.Encoder, key string, store references.Store, tracker references.Tracker) {
	var value = template.value(store, tracker)
	if value == nil {
		return
	}

	AddTypeKey(encoder, key, types.String, value.Key)
}

// Unmarshal attempts to unmarshal the value from the decoder as a enum and
// stores it inside the reference store.
func (template Enum) Unmarshal(decoder *gojay.Decoder, path string, store references.Store, tracker references.Tracker) error {
	var key string
	if err := decoder.AddString(&key); err != nil {
		return err
	}

	value := template.Enum.Keys[key]
	if value == nil {
		return nil
	}

	store.Store(tracker.Resolve(path), &references.Reference{
		Enum: &value.Position,
	})

	return nil
}

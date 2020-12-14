package json

import (
	"github.com/francoispqt/gojay"
	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/types"
)

func encode(encoder *gojay.Encoder, path string, template *specs.Template, store references.Store, tracker references.Tracker) error {
	typed := template.Type()

	switch {
	case typed == types.Message:
		encoder.Object(NewObject(path, template, store, tracker))
	case typed == types.Array:
		encoder.Array(NewArray(path, template, store, tracker))
	case template.Enum != nil:
		Enum(*template).Marshal(encoder, store, tracker)
	default:
		Scalar(*template).Marshal(encoder, store, tracker)
	}

	return nil
}

func encodeKey(encoder *gojay.Encoder, path, key string, template *specs.Template, store references.Store, tracker references.Tracker) {
	typed := template.Type()
	if (typed == types.Message || typed == types.Array) && template.Reference != nil {
		length := store.Length(tracker.Resolve(template.Reference.String()))
		if length == 0 {
			return
		}
	}

	switch {
	case typed == types.Message:
		encoder.AddObjectKeyOmitEmpty(key, NewObject(path, template, store, tracker))
	case template.Repeated != nil:
		encoder.AddArrayKey(key, NewArray(path, template, store, tracker))
	case template.Enum != nil:
		Enum(*template).MarshalKey(encoder, key, store, tracker)
	default:
		Scalar(*template).MarshalKey(encoder, key, store, tracker)
	}
}

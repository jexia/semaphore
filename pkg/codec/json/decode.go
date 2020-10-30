package json

import (
	"github.com/francoispqt/gojay"
	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
)

func decode(decoder *gojay.Decoder, path string, template specs.Template, store references.Store, tracker references.Tracker) error {
	switch {
	case template.Message != nil:
		return decoder.Object(NewObject(path, template, store, tracker))
	case template.Repeated != nil:
		return decoder.Array(NewArray(path, template, store, tracker))
	case template.Enum != nil:
		return Enum(template).Unmarshal(decoder, path, store, tracker)
	case template.Scalar != nil:
		return Scalar(template).Unmarshal(decoder, path, store, tracker)
	}

	return nil
}

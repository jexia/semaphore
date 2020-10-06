package json

import (
	"github.com/francoispqt/gojay"
	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
)

func decodeElement(decoder *gojay.Decoder, resource, path string, template specs.Template, store references.Store) error {
	var reference = &specs.PropertyReference{
		Resource: resource,
		Path:     path,
	}

	switch {
	case template.Message != nil:
		object := NewObject(resource, template.Message, store)
		if object == nil {
			break
		}

		return decoder.AddObject(object)
	case template.Repeated != nil:
		defer store.StoreReference(resource, &references.Reference{
			Path: path,
		})

		array := NewArray(resource, template.Repeated, template.Reference, store)
		if array == nil {
			break
		}

		return decoder.AddArray(array)
	case template.Enum != nil:
		return NewEnum("", template.Enum, reference, store).UnmarshalJSONEnum(decoder)
	case template.Scalar != nil:
		return NewScalar("", template.Scalar, reference, store).UnmarshalJSONScalar(decoder)
	}

	return nil
}

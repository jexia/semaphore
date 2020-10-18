package json

import (
	"github.com/francoispqt/gojay"
	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
)

func encodeElement(encoder *gojay.Encoder, resource string, template *specs.Template, store references.Store) {
	switch {
	case template.Message != nil:
		encoder.Object(
			NewObject(resource, template.Message, store),
		)
	case template.Repeated != nil:
		encoder.Array(
			NewArray(resource, template.Repeated, template.Reference, store),
		)
	case template.Enum != nil:
		NewEnum("", template.Enum, template.Reference, store).MarshalJSONEnum(encoder)
	case template.Scalar != nil:
		NewScalar("", template.Scalar, template.Reference, store).MarshalJSONScalar(encoder)
	}
}

func encodeElementKey(encoder *gojay.Encoder, resource, key string, template *specs.Template, store references.Store) {
	switch {
	case template.Message != nil:
		encoder.ObjectKey(
			key,
			NewObject(resource, template.Message, store),
		)
	case template.Repeated != nil:
		encoder.ArrayKey(
			key,
			NewArray(resource, template.Repeated, template.Reference, store),
		)
	case template.Enum != nil:
		NewEnum(key, template.Enum, template.Reference, store).MarshalJSONEnumKey(encoder)
	case template.Scalar != nil:
		NewScalar(key, template.Scalar, template.Reference, store).MarshalJSONScalarKey(encoder)
	}
}

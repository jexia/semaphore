package json

import (
	"github.com/francoispqt/gojay"
	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
)

// Array represents a JSON array
type Array struct {
	resource string
	template specs.Template
	ref      *references.Reference
	keys     int
}

// NewArray constructs a new JSON array encoder/decoder
func NewArray(resource string, repeated specs.Repeated, ref *specs.PropertyReference, refs references.Store) *Array {
	// skip arrays which have no elements
	if len(repeated) == 0 && ref == nil {
		return nil
	}

	var reference *references.Reference
	if ref != nil {
		reference = refs.Load(ref.Resource, ref.Path)
	}

	if ref != nil && reference == nil {
		return nil
	}

	template, err := repeated.Template()
	if err != nil {
		panic(err)
	}

	generator := &Array{
		resource: resource,
		template: template,
		ref:      reference,
	}

	if template.Repeated != nil {
		generator.keys = len(template.Repeated)
	}

	return generator
}

// MarshalJSONArray encodes the array into the given gojay encoder
func (array *Array) MarshalJSONArray(encoder *gojay.Encoder) {
	if array == nil || array.ref == nil {
		return
	}

	for _, store := range array.ref.Repeated {
		array.template.Reference = new(specs.PropertyReference)

		encodeElement(encoder, array.resource, array.template, store)
	}
}

// UnmarshalJSONArray unmarshals the given specs into the configured reference store
func (array *Array) UnmarshalJSONArray(decoder *gojay.Decoder) error {
	if array == nil {
		return nil
	}

	// FIXME: array.keys is derrived from a wrong value
	store := references.NewReferenceStore(array.keys)

	return decodeElement(decoder, "", "", array.template, store)
}

// IsNil returns whether the given array is null or not
func (array *Array) IsNil() bool {
	return array == nil
}

package json

import (
	"github.com/francoispqt/gojay"
	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
)

// Array represents a JSON array.
type Array struct {
	resource  string
	template  specs.Template
	repeated  specs.Repeated
	reference *specs.PropertyReference
	store     references.Store
}

// NewArray creates a new array to be JSON encoded/decoded.
func NewArray(resource string, repeated specs.Repeated, reference *specs.PropertyReference, store references.Store) *Array {
	template, err := repeated.Template()
	if err != nil {
		panic(err)
	}

	return &Array{
		resource:  resource,
		template:  template,
		repeated:  repeated,
		reference: reference,
		store:     store,
	}
}

// MarshalJSONArray encodes the array into the given gojay encoder.
func (array *Array) MarshalJSONArray(encoder *gojay.Encoder) {
	if array.reference == nil {
		for _, template := range array.repeated {
			encodeElement(encoder, "", template, array.store)
		}

		return
	}

	var reference = array.store.Load(array.reference.Resource, array.reference.Path)

	if reference == nil {
		return
	}

	for _, store := range reference.Repeated {
		array.template.Reference = new(specs.PropertyReference)

		encodeElement(encoder, array.resource, array.template, store)
	}
}

// UnmarshalJSONArray unmarshals the given specs into the configured reference store.
func (array *Array) UnmarshalJSONArray(decoder *gojay.Decoder) error {
	store := references.NewReferenceStore(0)

	if array.reference != nil {
		var reference = array.store.Load(array.reference.Resource, array.reference.Path)
		if reference == nil {
			return nil
		}

		reference.Append(store)
	}

	// NOTE: always consume an array even if the reference is not set
	return decodeElement(decoder, "", "", array.template, store)
}

// IsNil returns whether the given array is null or not.
func (array *Array) IsNil() bool {
	return array == nil
}

package json

import (
	"github.com/francoispqt/gojay"
	"github.com/jexia/maestro/codec/json/types"
	"github.com/jexia/maestro/refs"
	"github.com/jexia/maestro/specs"
)

// NewObject constructs a new object encoder/decoder for the given specs
func NewObject(resource string, specs specs.Object, refs *refs.Store) *Object {
	keys := len(specs.GetProperties()) + len(specs.GetNestedProperties()) + len(specs.GetRepeatedProperties())
	return &Object{
		resource: resource,
		keys:     keys,
		refs:     refs,
		specs:    specs,
	}
}

// Object represents a JSON object
type Object struct {
	resource string
	specs    specs.Object
	refs     *refs.Store
	keys     int
}

// MarshalJSONObject encodes the given specs object into the given gojay encoder
func (object *Object) MarshalJSONObject(encoder *gojay.Encoder) {
	for key, prop := range object.specs.GetProperties() {
		val := prop.Default

		if prop.Reference != nil {
			ref := object.refs.Load(prop.Reference.Resource, prop.Reference.Path)
			if ref != nil {
				val = ref.Value
			}
		}

		if val == nil {
			continue
		}

		types.Add(encoder, key, prop.Type, val)
	}

	for key, nested := range object.specs.GetNestedProperties() {
		result := NewObject(object.resource, nested, object.refs)
		encoder.AddObjectKey(key, result)
	}

	for key, repeated := range object.specs.GetRepeatedProperties() {
		ref := object.refs.Load(repeated.Template.Resource, repeated.Template.Path)
		if ref == nil {
			continue
		}

		array := NewArray(object.resource, repeated, ref.Repeated)
		encoder.AddArrayKey(key, array)
	}
}

// UnmarshalJSONObject unmarshals the given specs into the configured reference store
func (object *Object) UnmarshalJSONObject(dec *gojay.Decoder, k string) error {
	// TODO: unmarshal JSON

	return nil
}

// NKeys returns the ammount of available keys inside the given object
func (object *Object) NKeys() int {
	return object.keys
}

// IsNil returns whether the given object is null or not
func (object *Object) IsNil() bool {
	return false
}

// NewArray constructs a new JSON array encoder/decoder
func NewArray(resource string, object specs.Object, refs []*refs.Store) *Array {
	return &Array{
		resource: resource,
		specs:    object,
		items:    refs,
	}
}

// Array represents a JSON array
type Array struct {
	resource string
	specs    specs.Object
	items    []*refs.Store
}

// MarshalJSONArray encodes the array into the given gojay encoder
func (array *Array) MarshalJSONArray(enc *gojay.Encoder) {
	for _, store := range array.items {
		object := NewObject(array.resource, array.specs, store)
		enc.AddObject(object)
	}
}

// IsNil returns whether the given array is null or not
func (array *Array) IsNil() bool {
	return false
}

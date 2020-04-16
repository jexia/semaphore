package json

import (
	"github.com/francoispqt/gojay"
	"github.com/jexia/maestro/pkg/refs"
	"github.com/jexia/maestro/pkg/specs"
	"github.com/jexia/maestro/pkg/specs/labels"
	"github.com/jexia/maestro/pkg/specs/types"
)

// NewObject constructs a new object encoder/decoder for the given specs
func NewObject(resource string, specs map[string]*specs.Property, refs refs.Store) *Object {
	keys := len(specs)

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
	specs    map[string]*specs.Property
	refs     refs.Store
	keys     int
}

// MarshalJSONObject encodes the given specs object into the given gojay encoder
func (object *Object) MarshalJSONObject(encoder *gojay.Encoder) {
	for _, prop := range object.specs {
		if prop.Label == labels.Repeated {
			if prop.Reference == nil {
				continue
			}

			ref := object.refs.Load(prop.Reference.Resource, prop.Reference.Path)
			if ref == nil {
				continue
			}

			array := NewArray(object.resource, prop, ref, ref.Repeated)
			encoder.AddArrayKey(prop.Name, array)
			continue
		}

		if prop.Type == types.Message {
			result := NewObject(object.resource, prop.Nested, object.refs)
			encoder.AddObjectKey(prop.Name, result)
			continue
		}

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

		AddTypeKey(encoder, prop.Name, prop.Type, val)
	}
}

// UnmarshalJSONObject unmarshals the given specs into the configured reference store
func (object *Object) UnmarshalJSONObject(dec *gojay.Decoder, key string) error {
	prop, has := object.specs[key]
	if !has {
		return nil
	}

	if prop.Label == labels.Repeated {
		ref := refs.NewReference(prop.Path)
		array := NewArray(object.resource, prop, ref, nil)
		err := dec.AddArray(array)
		if err != nil {
			return err
		}

		object.refs.StoreReference(object.resource, ref)
		return nil
	}

	if prop.Type == types.Message {
		dynamic := NewObject(object.resource, prop.Nested, object.refs)
		err := dec.AddObject(dynamic)
		return err
	}

	ref := refs.NewReference(prop.Path)
	ref.Value = DecodeType(dec, prop.Type)
	object.refs.StoreReference(object.resource, ref)
	return nil
}

// NKeys returns the amount of available keys inside the given object
func (object *Object) NKeys() int {
	return object.keys
}

// IsNil returns whether the given object is null or not
func (object *Object) IsNil() bool {
	return false
}

// NewArray constructs a new JSON array encoder/decoder
func NewArray(resource string, object *specs.Property, ref *refs.Reference, refs []refs.Store) *Array {
	keys := 0

	if object.Nested != nil {
		keys = len(object.Nested)
	}

	return &Array{
		resource: resource,
		specs:    object,
		items:    refs,
		ref:      ref,
		keys:     keys,
	}
}

// Array represents a JSON array
type Array struct {
	resource string
	specs    *specs.Property
	items    []refs.Store
	ref      *refs.Reference
	keys     int
}

// MarshalJSONArray encodes the array into the given gojay encoder
func (array *Array) MarshalJSONArray(enc *gojay.Encoder) {
	for _, store := range array.items {
		if array.specs.Type == types.Message {
			object := NewObject(array.resource, array.specs.Nested, store)
			enc.AddObject(object)
			continue
		}

		val := array.specs.Default

		if array.specs.Reference != nil {
			ref := store.Load("", "")
			if ref != nil {
				val = ref.Value
			}
		}

		if val == nil {
			continue
		}

		AddType(enc, array.specs.Type, val)
	}
}

// UnmarshalJSONArray unmarshals the given specs into the configured reference store
func (array *Array) UnmarshalJSONArray(dec *gojay.Decoder) error {
	store := refs.NewReferenceStore(array.keys)

	if array.specs.Type == types.Message {
		object := NewObject(array.resource, array.specs.Nested, store)
		dec.AddObject(object)
		array.ref.Append(store)
		return nil
	}

	val := DecodeType(dec, array.specs.Type)
	store.StoreValue("", "", val)
	array.ref.Append(store)

	return nil
}

// IsNil returns whether the given array is null or not
func (array *Array) IsNil() bool {
	return false
}

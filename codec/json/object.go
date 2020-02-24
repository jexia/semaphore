package json

import (
	"github.com/francoispqt/gojay"
	"github.com/jexia/maestro/codec/json/types"
	"github.com/jexia/maestro/refs"
	"github.com/jexia/maestro/specs"
)

// TODO: configure whether to set default value, nil or ignore

func NewObject(object specs.Object, refs *refs.Store) *Object {
	keys := len(object.GetProperties()) + len(object.GetNestedProperties()) + len(object.GetRepeatedProperties())

	return &Object{
		specs: object,
		keys:  keys,
		refs:  refs,
	}
}

type Object struct {
	specs specs.Object
	refs  *refs.Store
	keys  int
}

func (object *Object) UnmarshalJSONObject(dec *gojay.Decoder, k string) error {
	// TODO: unmarshal JSON

	return nil
}

func (object *Object) NKeys() int {
	return object.keys
}

func (object *Object) MarshalJSONObject(encoder *gojay.Encoder) {
	for key, prop := range object.specs.GetProperties() {
		val := prop.Default

		if prop.Reference != nil {
			ref := object.refs.Load(prop.Reference.Resource, prop.Reference.Path)
			if ref != nil {
				val = ref.Value
			}
		}

		types.Add(encoder, key, prop.Type, val)
	}

	for key, nested := range object.specs.GetNestedProperties() {
		object := NewObject(nested, object.refs)
		encoder.AddObjectKey(key, object)
	}

	for key, repeated := range object.specs.GetRepeatedProperties() {
		ref := object.refs.Load(repeated.Template.Resource, repeated.Template.Path)
		if ref == nil {
			continue
		}

		array := NewArray(repeated, ref.Repeated)
		encoder.AddArrayKey(key, array)
	}
}

func (object *Object) IsNil() bool {
	return false
}

func NewArray(object specs.Object, refs []*refs.Store) *Array {
	return &Array{
		specs: object,
		refs:  refs,
	}
}

type Array struct {
	specs specs.Object
	refs  []*refs.Store
}

func (array *Array) MarshalJSONArray(enc *gojay.Encoder) {
	for _, store := range array.refs {
		object := NewObject(array.specs, store)
		enc.AddObject(object)
	}
}

func (array *Array) IsNil() bool {
	return false
}

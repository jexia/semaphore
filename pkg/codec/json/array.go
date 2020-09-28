package json

import (
	"github.com/francoispqt/gojay"
	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/types"
)

// NewArray constructs a new JSON array encoder/decoder
func NewArray(resource string, template specs.Template, refs references.Store) *Array {
	generator := &Array{
		resource: resource,
		template: template,
	}

	if template.Reference != nil {
		generator.ref = refs.Load(template.Reference.Resource, template.Reference.Path)
	}

	if template.Repeated != nil {
		generator.keys = len(template.Repeated)
	}

	return generator
}

// Array represents a JSON array
type Array struct {
	resource string
	template specs.Template
	ref      *references.Reference
	keys     int
}

// MarshalJSONArray encodes the array into the given gojay encoder
func (array *Array) MarshalJSONArray(enc *gojay.Encoder) {
	if array.ref == nil {
		return
	}

	for _, store := range array.ref.Repeated {
		switch {
		case array.template.Message != nil:
			enc.AddObject(NewObject(array.resource, array.template.Message, store))
			break
		case array.template.Repeated != nil:
			enc.AddArray(NewArray(array.resource, array.template, store))
			break
		case array.template.Enum != nil:
			// TODO: check if enums in arrays work
			if array.template.Reference == nil {
				break
			}

			ref := store.Load("", "")
			if ref == nil {
				break
			}

			key := array.template.Enum.Positions[*ref.Enum].Key
			AddType(enc, types.String, key)
			break
		case array.template.Scalar != nil:
			val := array.template.Scalar.Default

			ref := store.Load("", "")
			if ref != nil {
				val = ref.Value
			}

			AddType(enc, array.template.Scalar.Type, val)
			break
		}
	}
}

// UnmarshalJSONArray unmarshals the given specs into the configured reference store
func (array *Array) UnmarshalJSONArray(dec *gojay.Decoder) error {
	store := references.NewReferenceStore(array.keys)

	switch {
	case array.template.Message != nil:
		object := NewObject(array.resource, array.template.Message, store)
		err := dec.AddObject(object)
		if err != nil {
			return err
		}

		array.ref.Append(store)
		break
	case array.template.Repeated != nil:
		array := NewArray(array.resource, array.template, store)
		err := dec.AddArray(array)
		if err != nil {
			return err
		}

		array.ref.Append(store)
		break
	case array.template.Enum != nil:
		var key string
		err := dec.AddString(&key)
		if err != nil {
			return err
		}

		enum := array.template.Enum.Keys[key]
		if enum != nil {
			store.StoreEnum("", "", enum.Position)
			array.ref.Append(store)
		}
		break
	case array.template.Scalar != nil:
		val, err := DecodeType(dec, array.template.Scalar.Type)
		if err != nil {
			return err
		}

		store.StoreValue("", "", val)
		array.ref.Append(store)
		break
	}

	return nil
}

// IsNil returns whether the given array is null or not
func (array *Array) IsNil() bool {
	return false
}

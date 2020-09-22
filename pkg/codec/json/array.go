package json

import (
	"github.com/francoispqt/gojay"
	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/types"
)

// NewArray constructs a new JSON array encoder/decoder
func NewArray(resource string, repeated *specs.Repeated, refs references.Store) *Array {
	generator := &Array{
		resource: resource,
		specs:    repeated,
	}

	if repeated.Reference != nil {
		generator.ref = refs.Load(repeated.Reference.Resource, repeated.Reference.Path)
	}

	if repeated.Default != nil {
		generator.keys = len(repeated.Default)
	}

	return generator
}

// Array represents a JSON array
type Array struct {
	resource string
	specs    *specs.Repeated
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
		case array.specs.Message != nil:
			enc.AddObject(NewObject(array.resource, array.specs.Message, store))
			break
		case array.specs.Repeated != nil:
			enc.AddArray(NewArray(array.resource, array.specs.Repeated, store))
			break
		case array.specs.Enum != nil:
			// TODO: check if enums in arrays work
			if array.specs.Reference == nil {
				break
			}

			ref := store.Load("", "")
			if ref == nil {
				break
			}

			key := array.specs.Enum.Positions[*ref.Enum].Key
			AddType(enc, types.String, key)
			break
		case array.specs.Scalar != nil:
			val := array.specs.Scalar.Default

			ref := store.Load("", "")
			if ref != nil {
				val = ref.Value
			}

			AddType(enc, array.specs.Scalar.Type, val)
			break
		}
	}
}

// UnmarshalJSONArray unmarshals the given specs into the configured reference store
func (array *Array) UnmarshalJSONArray(dec *gojay.Decoder) error {
	store := references.NewReferenceStore(array.keys)

	switch {
	case array.specs.Message != nil:
		object := NewObject(array.resource, array.specs.Message, store)
		err := dec.AddObject(object)
		if err != nil {
			return err
		}

		array.ref.Append(store)
		break
	case array.specs.Repeated != nil:
		array := NewArray(array.resource, array.specs.Repeated, store)
		err := dec.AddArray(array)
		if err != nil {
			return err
		}

		array.ref.Append(store)
		break
	case array.specs.Enum != nil:
		var key string
		err := dec.AddString(&key)
		if err != nil {
			return err
		}

		enum := array.specs.Enum.Keys[key]
		if enum != nil {
			store.StoreEnum("", "", enum.Position)
			array.ref.Append(store)
		}
		break
	case array.specs.Scalar != nil:
		val, err := DecodeType(dec, array.specs.Scalar.Type)
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

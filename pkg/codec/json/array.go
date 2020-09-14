package json

import (
	"github.com/francoispqt/gojay"
	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/types"
)

// NewArray constructs a new JSON array encoder/decoder
func NewArray(resource string, object *specs.Property, ref *references.Reference, refs []references.Store) *Array {
	keys := 0

	if object.Repeated != nil {
		keys = len(object.Repeated)
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
	items    []references.Store
	ref      *references.Reference
	keys     int
}

// MarshalJSONArray encodes the array into the given gojay encoder
func (array *Array) MarshalJSONArray(enc *gojay.Encoder) {
	for _, store := range array.items {
		if array.specs.Type == types.Message {
			object := NewObject(array.resource, array.specs.Repeated, store)
			enc.AddObject(object)
			continue
		}

		val := array.specs.Default

		if array.specs.Reference != nil {
			ref := store.Load("", "")
			if ref != nil {
				if ref.Enum != nil && array.specs.Enum != nil {
					val = array.specs.Enum.Positions[*ref.Enum].Key
				}

				if val == nil {
					val = ref.Value
				}
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
	store := references.NewReferenceStore(array.keys)

	if array.specs.Type == types.Message {
		object := NewObject(array.resource, array.specs.Repeated, store)
		err := dec.AddObject(object)
		if err != nil {
			return err
		}

		array.ref.Append(store)
		return nil
	}

	if array.specs.Type == types.Enum && array.specs.Enum != nil {
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

		return nil
	}

	val, err := DecodeType(dec, array.specs.Type)
	if err != nil {
		return err
	}

	store.StoreValue("", "", val)
	array.ref.Append(store)

	return nil
}

// IsNil returns whether the given array is null or not
func (array *Array) IsNil() bool {
	return false
}

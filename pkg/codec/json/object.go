package json

import (
	"github.com/francoispqt/gojay"
	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/labels"
	"github.com/jexia/semaphore/pkg/specs/types"
)

// NewObject constructs a new object encoder/decoder for the given specs
func NewObject(resource string, items []*specs.Property, refs references.Store) *Object {
	// TODO: check overhead of initialization
	keys := make(map[string]*specs.Property, len(items))
	for _, spec := range items {
		keys[spec.Name] = spec
	}

	return &Object{
		resource: resource,
		length:   len(items),
		refs:     refs,
		specs:    items,
		keys:     keys,
	}
}

// Object represents a JSON object
type Object struct {
	resource string
	specs    []*specs.Property
	refs     references.Store
	length   int
	keys     map[string]*specs.Property
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
				if prop.Type == types.Enum && ref.Enum != nil {
					enum := prop.Enum.Positions[*ref.Enum]
					if enum != nil {
						val = enum.Key
					}
				} else if ref.Value != nil {
					val = ref.Value
				}
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
	prop, has := object.keys[key]
	if !has {
		return nil
	}

	if prop.Label == labels.Repeated {
		ref := &references.Reference{
			Path: prop.Path,
		}

		array := NewArray(object.resource, prop, ref, nil)
		err := dec.AddArray(array)
		if err != nil {
			return err
		}

		object.refs.StoreReference(object.resource, ref)
		return nil
	}

	if prop.Type == types.Message {
		return dec.AddObject(
			NewObject(object.resource, prop.Nested, object.refs),
		)
	}

	ref := &references.Reference{
		Path: prop.Path,
	}

	if prop.Type == types.Enum {
		var key string
		dec.AddString(&key)

		enum := prop.Enum.Keys[key]
		if enum != nil {
			ref.Enum = &enum.Position
		}
	} else {
		value, err := DecodeType(dec, prop.Type)
		if err != nil {
			return err
		}

		ref.Value = value
	}

	object.refs.StoreReference(object.resource, ref)
	return nil
}

// NKeys returns the amount of available keys inside the given object
func (object *Object) NKeys() int {
	return object.length
}

// IsNil returns whether the given object is null or not
func (object *Object) IsNil() bool {
	return false
}

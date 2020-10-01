package json

import (
	"github.com/francoispqt/gojay"
	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/types"
)

// NewObject constructs a new object encoder/decoder for the given specs
func NewObject(resource string, items specs.Message, refs references.Store) *Object {
	return &Object{
		resource: resource,
		length:   len(items),
		refs:     refs,
		specs:    items,
	}
}

// Object represents a JSON object
type Object struct {
	resource string
	specs    specs.Message
	refs     references.Store
	length   int
}

// MarshalJSONObject encodes the given specs object into the given gojay encoder
func (object *Object) MarshalJSONObject(encoder *gojay.Encoder) {
	if object == nil {
		return
	}

	for _, prop := range object.specs {
		switch {
		case prop.Repeated != nil:
			array := NewArray(object.resource, prop.Template, object.refs)
			if array == nil {
				break
			}

			encoder.AddArrayKey(prop.Name, array)
			break
		case prop.Message != nil:
			result := NewObject(object.resource, prop.Message, object.refs)
			if result == nil {
				break
			}

			encoder.AddObjectKey(prop.Name, result)
			break
		case prop.Enum != nil:
			if prop.Reference == nil {
				break
			}

			ref := object.refs.Load(prop.Reference.Resource, prop.Reference.Path)
			if ref == nil || ref.Enum == nil {
				break
			}

			enum := prop.Enum.Positions[*ref.Enum]
			if enum == nil {
				break
			}

			AddTypeKey(encoder, prop.Name, types.String, enum.Key)
			break
		case prop.Scalar != nil:
			val := prop.Scalar.Default

			if prop.Reference != nil {
				ref := object.refs.Load(prop.Reference.Resource, prop.Reference.Path)
				if ref != nil && ref.Value != nil {
					val = ref.Value
				}
			}

			if val == nil {
				break
			}

			AddTypeKey(encoder, prop.Name, prop.Scalar.Type, val)
			break
		}

	}
}

// UnmarshalJSONObject unmarshals the given specs into the configured reference store
func (object *Object) UnmarshalJSONObject(dec *gojay.Decoder, key string) error {
	if object == nil {
		return nil
	}

	property, has := object.specs[key]
	if !has {
		return nil
	}

	switch {
	case property.Message != nil:
		object := NewObject(object.resource, property.Message, object.refs)
		if object == nil {
			break
		}

		return dec.AddObject(object)
	case property.Repeated != nil:
		ref := &references.Reference{
			Path: property.Path,
		}

		array := NewArray(object.resource, property.Template, object.refs)
		if array == nil {
			break
		}

		err := dec.AddArray(array)
		if err != nil {
			return err
		}

		object.refs.StoreReference(object.resource, ref)
		return nil
	}

	ref := &references.Reference{
		Path: property.Path,
	}

	switch {
	case property.Enum != nil:
		var key string
		dec.AddString(&key)

		enum := property.Enum.Keys[key]
		if enum != nil {
			ref.Enum = &enum.Position
		}

		break
	case property.Scalar != nil:
		value, err := DecodeType(dec, property.Scalar.Type)
		if err != nil {
			return err
		}

		ref.Value = value
		break
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

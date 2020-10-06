package json

import (
	"github.com/francoispqt/gojay"
	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
)

type Scalar struct {
	name      string
	scalar    *specs.Scalar
	reference *specs.PropertyReference
	store     references.Store
}

func NewScalar(name string, scalar *specs.Scalar, reference *specs.PropertyReference, store references.Store) *Scalar {
	return &Scalar{
		name:      name,
		scalar:    scalar,
		reference: reference,
		store:     store,
	}
}

func (scalar Scalar) Value() interface{} {
	var value = scalar.scalar.Default

	if scalar.reference != nil {
		var reference = scalar.store.Load(scalar.reference.Resource, scalar.reference.Path)
		if reference != nil && reference.Value != nil {
			value = reference.Value
		}
	}

	return value
}

func (scalar Scalar) MarshalJSONScalar(encoder *gojay.Encoder) {
	AddType(encoder, scalar.scalar.Type, scalar.Value())
}

func (scalar Scalar) MarshalJSONScalarKey(encoder *gojay.Encoder) {
	var value = scalar.Value()
	if value == nil {
		return
	}

	AddTypeKey(encoder, scalar.name, scalar.scalar.Type, value)
}

func (scalar Scalar) UnmarshalJSONScalar(decoder *gojay.Decoder) error {
	value, err := DecodeType(decoder, scalar.scalar.Type)
	if err != nil {
		return err
	}

	var reference = &references.Reference{
		Path:  scalar.reference.Path,
		Value: value,
	}

	scalar.store.StoreReference(scalar.reference.Resource, reference)

	return nil
}

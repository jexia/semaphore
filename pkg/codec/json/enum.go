package json

import (
	"github.com/francoispqt/gojay"
	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/types"
)

// Enum is a vrapper over specs.Enum providing JSON encoding/decoding.
type Enum struct {
	name      string
	enum      *specs.Enum
	reference *specs.PropertyReference
	store     references.Store
}

// NewEnum creates a new enum by wrapping provided specs.Enum.
func NewEnum(name string, enum *specs.Enum, reference *specs.PropertyReference, store references.Store) *Enum {
	return &Enum{
		name:      name,
		enum:      enum,
		reference: reference,
		store:     store,
	}
}

func (enum *Enum) value() *specs.EnumValue {
	if enum.reference == nil {
		return nil
	}

	var reference = enum.store.Load(enum.reference.Resource, enum.reference.Path)
	if reference == nil || reference.Enum == nil {
		return nil
	}

	if position := reference.Enum; position != nil {
		return enum.enum.Positions[*reference.Enum]
	}

	return nil
}

// MarshalJSONEnum marshals enum which is not an object field (array element or
// standalone enum).
func (enum Enum) MarshalJSONEnum(encoder *gojay.Encoder) {
	var value interface{}
	if enumValue := enum.value(); enumValue != nil {
		value = enumValue.Key
	}

	AddType(encoder, types.String, value)
}

// MarshalJSONEnumKey marshals enum which is an object field.
func (enum Enum) MarshalJSONEnumKey(encoder *gojay.Encoder) {
	var enumValue = enum.value()
	if enumValue == nil {
		return
	}

	AddTypeKey(encoder, enum.name, types.String, enumValue.Key)
}

// UnmarshalJSONEnum unmarshals enum value from decoder to the reference store.
func (enum Enum) UnmarshalJSONEnum(decoder *gojay.Decoder) error {
	var key string
	if err := decoder.AddString(&key); err != nil {
		return err
	}

	var (
		reference = &references.Reference{
			Path: enum.reference.Path,
		}
		enumValue = enum.enum.Keys[key]
	)

	if enumValue != nil {
		reference.Enum = &enumValue.Position
	}

	enum.store.StoreReference(enum.reference.Resource, reference)

	return nil
}

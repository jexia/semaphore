package xml

import (
	"encoding/xml"

	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
)

// Enum is a vrapper over specs.Enum providing XML encoding/decoding.
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

// MarshalXML marshals given enum to XML.
func (enum *Enum) MarshalXML(encoder *xml.Encoder, _ xml.StartElement) error {
	var (
		value string
		start = xml.StartElement{
			Name: xml.Name{
				Local: enum.name,
			},
		}
	)

	if enum.reference != nil {
		var reference = enum.store.Load(enum.reference.Resource, enum.reference.Path)

		if reference == nil || reference.Enum == nil {
			return nil
		}

		var enumValue = enum.enum.Positions[*reference.Enum]
		if enumValue != nil {
			value = enumValue.Key
		}
	}

	return encoder.EncodeElement(value, start)
}

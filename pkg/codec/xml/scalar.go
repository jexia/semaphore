package xml

import (
	"encoding/xml"

	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
)

// Scalar is a wrapper for specs.Scalar providing XML encoding/decoding.
type Scalar struct {
	name      string
	scalar    *specs.Scalar
	reference *specs.PropertyReference
	store     references.Store
}

// NewScalar creates a wrapper for specs.Scalar to be XML encoded/decoded.
func NewScalar(name string, scalar *specs.Scalar, reference *specs.PropertyReference, store references.Store) *Scalar {
	return &Scalar{
		name:      name,
		scalar:    scalar,
		reference: reference,
		store:     store,
	}
}

// MarshalXML marshals scalar value to XML.
func (s Scalar) MarshalXML(encoder *xml.Encoder, _ xml.StartElement) error {
	var (
		value = s.scalar.Default
		start = xml.StartElement{
			Name: xml.Name{
				Local: s.name,
			},
		}
	)

	if s.reference != nil {
		var reference = s.store.Load(s.reference.Resource, s.reference.Path)
		if reference == nil {
			return nil
		}

		if reference.Value != nil {
			value = reference.Value
		}
	}

	return encoder.EncodeElement(value, start)
}

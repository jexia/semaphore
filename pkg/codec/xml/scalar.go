package xml

import (
	"encoding/xml"

	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/types"
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
func (scalar Scalar) MarshalXML(encoder *xml.Encoder, _ xml.StartElement) error {
	var (
		value = scalar.scalar.Default
		start = xml.StartElement{
			Name: xml.Name{
				Local: scalar.name,
			},
		}
	)

	if scalar.reference != nil {
		var reference = scalar.store.Load(scalar.reference.Resource, scalar.reference.Path)
		if reference != nil && reference.Value != nil {
			value = reference.Value
		}
	}

	return encoder.EncodeElement(value, start)
}

// UnmarshalXML unmarshals scalar value from XML stream.
func (scalar *Scalar) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) (err error) {
	const (
		waitForValue int = iota
		waitForClose
	)

	var state int

	for {
		tok, err := decoder.Token()
		if err != nil {
			return err
		}

		switch state {
		case waitForValue:
			var reference = &references.Reference{
				Path: buildPath(scalar.reference.Path, scalar.name),
			}

			switch t := tok.(type) {
			case xml.CharData:
				if reference.Value, err = types.DecodeFromString(string(t), scalar.scalar.Type); err != nil {
					return err
				}

				scalar.store.StoreReference(scalar.reference.Resource, reference)

				state = waitForClose
			case xml.EndElement:
				scalar.store.StoreReference(scalar.reference.Resource, reference)
				// scalar is closed with nil value
				return nil
			default:
				return errUnexpectedToken{
					actual: t,
					expected: []xml.Token{
						xml.CharData{},
						xml.EndElement{},
					},
				}
			}
		case waitForClose:
			switch t := tok.(type) {
			case xml.EndElement:
				// element is closed
				return nil
			default:
				return errUnexpectedToken{
					actual: t,
					expected: []xml.Token{
						xml.EndElement{},
					},
				}
			}
		}
	}
}

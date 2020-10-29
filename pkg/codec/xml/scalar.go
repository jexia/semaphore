package xml

import (
	"encoding/xml"

	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/types"
)

// Scalar is a wrapper for specs.Scalar providing XML encoding/decoding.
type Scalar struct {
	name, path string
	template   specs.Template
	store      references.Store
	tracker    references.Tracker
}

// NewScalar creates a wrapper for specs.Scalar to be XML encoded/decoded.
func NewScalar(name, path string, template specs.Template, store references.Store, tracker references.Tracker) *Scalar {
	return &Scalar{
		name:     name,
		path:     path,
		template: template,
		store:    store,
		tracker:  tracker,
	}
}

// MarshalXML marshals scalar value to XML.
func (scalar Scalar) MarshalXML(encoder *xml.Encoder, _ xml.StartElement) error {
	var (
		value = scalar.template.Scalar.Default
		start = xml.StartElement{
			Name: xml.Name{
				Local: scalar.name,
			},
		}
	)

	if scalar.template.Reference != nil {
		var reference = scalar.store.Load(scalar.tracker.Resolve(scalar.template.Reference.String()))
		if reference != nil && reference.Value != nil {
			value = reference.Value
		}
	}

	return encoder.EncodeElement(value, start)
}

// UnmarshalXML unmarshals scalar value from XML stream.
func (scalar *Scalar) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
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
			switch t := tok.(type) {
			case xml.CharData:
				value, err := types.DecodeFromString(string(t), scalar.template.Scalar.Type)
				if err != nil {
					return err
				}

				scalar.store.Store(scalar.tracker.Resolve(scalar.path), &references.Reference{
					Value: value,
				})

				state = waitForClose
			case xml.EndElement:
				scalar.store.Store(scalar.tracker.Resolve(scalar.path), new(references.Reference))
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

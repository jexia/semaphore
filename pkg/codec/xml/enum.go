package xml

import (
	"encoding/xml"

	"github.com/jexia/semaphore/v2/pkg/references"
	"github.com/jexia/semaphore/v2/pkg/specs"
	"github.com/jexia/semaphore/v2/pkg/specs/template"
)

// Enum is a vrapper over specs.Enum providing XML encoding/decoding.
type Enum struct {
	name, path string
	template   specs.Template
	store      references.Store
	tracker    references.Tracker
}

// NewEnum creates a new enum by wrapping provided specs.Enum.
func NewEnum(name, path string, template specs.Template, store references.Store, tracker references.Tracker) *Enum {
	return &Enum{
		name:     name,
		path:     path,
		template: template,
		store:    store,
		tracker:  tracker,
	}
}

// MarshalXML marshals given enum to XML.
func (enum *Enum) MarshalXML(encoder *xml.Encoder, _ xml.StartElement) error {
	var (
		start = xml.StartElement{
			Name: xml.Name{
				Local: enum.name,
			},
		}
		value = enum.value(enum.store, enum.tracker)
	)

	if value == nil {
		return nil
	}

	return encoder.EncodeElement(value.Key, start)
}

func (enum *Enum) value(store references.Store, tracker references.Tracker) *specs.EnumValue {
	if enum.template.Reference == nil {
		return nil
	}

	reference := store.Load(tracker.Resolve(enum.template.Reference.String()))
	if reference == nil || reference.Enum == nil {
		return nil
	}

	if position := reference.Enum; position != nil {
		return enum.template.Enum.Positions[*reference.Enum]
	}

	return nil
}

// UnmarshalXML unmarshals enum value from XML stream.
func (enum *Enum) UnmarshalXML(decoder *xml.Decoder, _ xml.StartElement) error {
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
				enumValue, ok := enum.template.Enum.Keys[string(t)]
				if !ok {
					return errUnknownEnum(t)
				}

				enum.store.Store(
					enum.tracker.Resolve(
						template.JoinPath(enum.path, enum.name),
					),
					&references.Reference{
						Enum: &enumValue.Position,
					},
				)

				state = waitForClose
			case xml.EndElement:
				enum.store.Store(
					enum.tracker.Resolve(
						template.JoinPath(enum.path, enum.name),
					),
					new(references.Reference),
				)

				// enum is closed with nil value
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
				// enum is closed
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

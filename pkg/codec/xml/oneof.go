package xml

import (
	"encoding/xml"

	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/template"
)

// OneOf represents an XML object with a single field allowed to be set.
// TODO: implement single element protection while decoding oneof.
type OneOf struct {
	name, path string
	template   specs.Template
	store      references.Store
	tracker    references.Tracker
}

// NewOneOf creates a new oneof.
func NewOneOf(name, path string, template specs.Template, store references.Store, tracker references.Tracker) *OneOf {
	return &OneOf{
		name:     name,
		path:     path,
		template: template,
		store:    store,
		tracker:  tracker,
	}
}

// MarshalXML encodes the given oneof into the given XML encoder.
func (oneof *OneOf) MarshalXML(encoder *xml.Encoder, _ xml.StartElement) error {
	var start = xml.StartElement{Name: xml.Name{Local: oneof.name}}

	if err := encoder.EncodeToken(start); err != nil {
		return err
	}

	// TODO: properties are now sorted during runtime. This process should be
	// moved to be prepared before MarshalXML is called.
	for _, property := range oneof.template.OneOf {
		if err := encodeElement(
			encoder,
			property.Name,
			template.JoinPath(oneof.path, property.Name),
			property.Template,
			oneof.store,
			oneof.tracker,
		); err != nil {
			return err
		}
	}

	return encoder.EncodeToken(xml.EndElement{Name: start.Name})
}

// UnmarshalXML decodes XML input into the receiver of type specs.OneOf.
func (oneof *OneOf) UnmarshalXML(decoder *xml.Decoder, _ xml.StartElement) error {
	for {
		tok, err := decoder.Token()
		if err != nil {
			return err
		}

		switch t := tok.(type) {
		case xml.StartElement:
			property, ok := oneof.template.OneOf[t.Name.Local]
			if !ok {
				// TODO: allow unknown fields
				return errUndefinedProperty(t.Name.Local)
			}

			oneof.store.Define(oneof.path, len(oneof.template.OneOf))

			if err := decodeElement(
				decoder,
				t,             // start element
				property.Name, // name
				template.JoinPath(oneof.path, oneof.name), // path
				property.Template,
				oneof.store,
				oneof.tracker,
			); err != nil {
				return err
			}
		case xml.EndElement:
			// oneof is closed
			return nil
		}
	}
}

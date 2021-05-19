package xml

import (
	"encoding/xml"

	"github.com/jexia/semaphore/v2/pkg/references"
	"github.com/jexia/semaphore/v2/pkg/specs"
	"github.com/jexia/semaphore/v2/pkg/specs/template"
)

// Object represents an XML object.
type Object struct {
	name, path string
	template   specs.Template
	store      references.Store
	tracker    references.Tracker
}

// NewObject creates a new object by wrapping provided specs.Message.
func NewObject(name, path string, template specs.Template, store references.Store, tracker references.Tracker) *Object {
	return &Object{
		name:     name,
		path:     path,
		template: template,
		store:    store,
		tracker:  tracker,
	}
}

// MarshalXML encodes the given specs object into the given XML encoder.
func (object *Object) MarshalXML(encoder *xml.Encoder, _ xml.StartElement) error {
	start := xml.StartElement{Name: xml.Name{Local: object.name}}

	if err := encoder.EncodeToken(start); err != nil {
		return err
	}

	// TODO: properties are now sorted during runtime. This process should be
	// moved to be prepared before MarshalXML is called.
	for _, property := range object.template.Message.SortedProperties() {
		if err := encodeElement(
			encoder,
			property.Name,
			template.JoinPath(object.path, property.Name),
			property.Template,
			object.store,
			object.tracker,
		); err != nil {
			return err
		}
	}

	return encoder.EncodeToken(xml.EndElement{Name: start.Name})
}

// UnmarshalXML decodes XML input into the receiver of type specs.Message.
func (object *Object) UnmarshalXML(decoder *xml.Decoder, _ xml.StartElement) error {
	for {
		tok, err := decoder.Token()
		if err != nil {
			return err
		}

		switch t := tok.(type) {
		case xml.StartElement:
			property, ok := object.template.Message[t.Name.Local]
			if !ok {
				// TODO: allow unknown fields
				return errUndefinedProperty(t.Name.Local)
			}

			object.store.Define(object.path, len(object.template.Message))

			if err := decodeElement(
				decoder,
				t,             // start element
				property.Name, // name
				template.JoinPath(object.path, object.name), // path
				property.Template,
				object.store,
				object.tracker,
			); err != nil {
				return err
			}
		case xml.EndElement:
			// object is closed
			return nil
		}
	}
}

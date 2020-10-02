package xml

import (
	"encoding/xml"

	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
)

// Object represents an XML object.
type Object struct {
	name      string
	message   specs.Message
	reference *specs.PropertyReference
	store     references.Store
}

// NewObject creates a new object by wrapping provided specs.Message.
func NewObject(name string, message specs.Message, reference *specs.PropertyReference, store references.Store) *Object {
	return &Object{
		name:      name,
		message:   message,
		reference: reference,
		store:     store,
	}
}

// MarshalXML encodes the given specs object into the given XML encoder.
func (object *Object) MarshalXML(encoder *xml.Encoder, _ xml.StartElement) error {
	var start = xml.StartElement{Name: xml.Name{Local: object.name}}

	if err := encoder.EncodeToken(start); err != nil {
		return err
	}

	// TODO: properties are now sorted during runtime. This process should be
	// moved to be prepared before MarshalXML is called.
	for _, property := range object.message.SortedProperties() {
		if err := encodeElement(encoder, property.Name, property.Template, object.store); err != nil {
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
			property, ok := object.message[t.Name.Local]
			if !ok {
				return errUndefinedProperty(t.Name.Local)
			}

			if err := decodeElement(
				decoder,
				t,                         // start element
				object.reference.Resource, // resource name
				buildPath(object.reference.Path, object.name), // path
				property.Name,
				property.Template,
				object.store,
			); err != nil {
				return err
			}
		case xml.EndElement:
			// object is closed
			return nil
		}
	}
}

package xml

import (
	"encoding/xml"
	"io"

	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/types"
)

// Object represents an XML object.
type Object struct {
	resource string
	specs    specs.Message
	refs     references.Store
}

// NewObject constructs a new object encoder/decoder for the given specs.
func NewObject(resource string, specs specs.Message, refs references.Store) *Object {
	return &Object{
		resource: resource,
		refs:     refs,
		specs:    specs,
	}
}

// MarshalXML encodes the given specs object into the given XML encoder.
func (object *Object) MarshalXML(encoder *xml.Encoder, _ xml.StartElement) error {
	var start = xml.StartElement{Name: xml.Name{Local: object.resource}}

	if err := encoder.EncodeToken(start); err != nil {
		return err
	}

	for _, spec := range object.specs {
		if err := object.encodeElement(encoder, spec); err != nil {
			return err
		}
	}

	return encoder.EncodeToken(xml.EndElement{Name: start.Name})
}

func (object *Object) encodeElement(encoder *xml.Encoder, prop *specs.Property) error {
	if prop.Repeated != nil {
		return encodeRepeated(encoder, object.resource, prop, object.refs)
	}

	// TODO: hide empty nested objects
	if prop.Type() == types.Message {
		return encodeNested(encoder, prop, object.refs)
	}

	return encodeValue(encoder, prop, object.refs, true)
}

// UnmarshalXML decodes XML input into the reference store.
func (object *Object) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	var refs = make(map[string]*references.Reference)

	defer func() {
		for _, reference := range refs {
			object.refs.StoreReference(object.resource, reference)
		}
	}()

	return object.unmarshalXML(decoder, refs)
}

func (object *Object) unmarshalXML(decoder *xml.Decoder, refs map[string]*references.Reference) error {
	tok, err := decoder.Token()
	if err == io.EOF {
		return nil
	}

	if err != nil {
		return err
	}

	return object.startElement(decoder, tok, refs)
}

func (object *Object) startElement(decoder *xml.Decoder, tok xml.Token, refs map[string]*references.Reference) error {
	switch t := tok.(type) {
	case xml.StartElement:
		var prop = object.specs[t.Name.Local]
		if prop == nil {
			return errUndefinedProperty(t.Name.Local)
		}

		if prop.Repeated != nil {
			return object.repeated(decoder, prop, refs)
		}

		return object.propertyValue(decoder, prop, refs)
	case xml.EndElement:
		// object is closed
		return nil
	default:
		return errUnexpectedToken{
			actual: t,
			expected: []xml.Token{
				xml.StartElement{},
				xml.EndElement{},
			},
		}
	}
}

func (object *Object) repeated(decoder *xml.Decoder, prop *specs.Property, refs map[string]*references.Reference) error {
	tok, err := decoder.Token()
	if err != nil {
		return err
	}

	switch t := tok.(type) {
	case xml.CharData:
		if err := decodeRepeatedValue(t, prop, refs); err != nil {
			return err
		}

		return object.closeElement(decoder, prop, refs)
	case xml.StartElement:
		if err := decodeRepeatedNested(decoder, t, prop, object.resource, refs); err != nil {
			return err
		}

		return object.unmarshalXML(decoder, refs)
	case xml.EndElement:
		return object.unmarshalXML(decoder, refs)
	default:
		return errUnexpectedToken{
			actual: t,
			expected: []xml.Token{
				xml.StartElement{},
				xml.EndElement{},
			},
		}
	}
}

func (object *Object) propertyValue(decoder *xml.Decoder, prop *specs.Property, refs map[string]*references.Reference) error {
	tok, err := decoder.Token()
	if err != nil {
		return err
	}

	switch t := tok.(type) {
	case xml.StartElement:
		if err := decodeNested(decoder, t, prop, object.resource, object.refs, refs); err != nil {
			return err
		}

		return object.unmarshalXML(decoder, refs)
	case xml.EndElement:
		return object.unmarshalXML(decoder, refs)
	case xml.CharData:
		if err := decodeValue(t, prop, object.resource, object.refs); err != nil {
			return err
		}

		return object.closeElement(decoder, prop, refs)
	default:
		return errUnexpectedToken{
			actual: t,
			expected: []xml.Token{
				xml.StartElement{},
				xml.CharData{},
				xml.EndElement{},
			},
		}
	}
}

func (object *Object) closeElement(decoder *xml.Decoder, prop *specs.Property, refs map[string]*references.Reference) error {
	tok, err := decoder.Token()
	if err != nil {
		return err
	}

	switch t := tok.(type) {
	case xml.EndElement:
		return object.unmarshalXML(decoder, refs)
	default:
		return errUnexpectedToken{
			actual: t,
			expected: []xml.Token{
				xml.EndElement{},
			},
		}
	}
}

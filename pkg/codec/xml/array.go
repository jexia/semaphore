package xml

import (
	"encoding/xml"
	"errors"
	"io"

	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/types"
)

// Array represents an array of values/references.
type Array struct {
	resource string
	property *specs.Property
	items    []references.Store
	ref      *references.Reference
}

// NewArray constructs a new JSON array encoder/decoder
func NewArray(resource string, property *specs.Property, refs []references.Store) *Array {
	return &Array{
		resource: resource,
		property: property,
		items:    refs,
		ref: &references.Reference{
			Path: property.Path,
		},
	}
}

// MarshalXML encodes the given specs object into the provided XML encoder.
func (array *Array) MarshalXML(encoder *xml.Encoder, _ xml.StartElement) error {
	for _, store := range array.items {
		if err := array.encodeElement(encoder, store); err != nil {
			return err
		}
	}

	return nil
}

func (array *Array) encodeElement(encoder *xml.Encoder, store references.Store) error {
	if array.property.Type == types.Message {
		return encodeNested(encoder, array.property, store)
	}

	return encodeValue(encoder, array.property, store, false)
}

func (array *Array) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	for {
		tok, err := decoder.Token()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}

		switch t := tok.(type) {
		case xml.StartElement:
			if err := array.decodeNested(decoder, t); err != nil {
				return err
			}
		default:
			errors.New("Boo!")
		}
	}

	return nil
}

func (array *Array) decodeNested(decoder *xml.Decoder, start xml.StartElement) error {
	if array.property.Type == types.Message {
		return array.decodeMessage(decoder, start)
	}

	return nil
}

func (array *Array) decodeMessage(decoder *xml.Decoder, start xml.StartElement) error {
	var (
		refs   = make(map[string]*references.Reference)
		store  = references.NewReferenceStore(1)
		object = NewObject(array.resource, array.property, store)
		err    = object.startElement(decoder, start, refs)
	)

	if err != nil && err != errEOS {
		return err
	}

	for _, reference := range refs {
		object.refs.StoreReference(object.resource, reference)
	}

	array.ref.Append(store)

	return nil
}

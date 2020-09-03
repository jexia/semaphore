package xml

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"sort"

	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/labels"
	"github.com/jexia/semaphore/pkg/specs/types"
)

// Object represents a JSON object
type Object struct {
	resource string
	specs    map[string]*specs.Property
	refs     references.Store
}

// NewObject constructs a new object encoder/decoder for the given specs
func NewObject(resource string, specs map[string]*specs.Property, refs references.Store) *Object {
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

	keys := make([]string, 0, len(object.specs))
	for key := range object.specs {
		keys = append(keys, key)
	}

	// sort properties by name
	sort.Strings(keys)

	for _, key := range keys {
		if err := object.encodeElement(encoder, object.specs[key]); err != nil {
			return err
		}
	}

	return encoder.EncodeToken(xml.EndElement{Name: start.Name})
}

func (object *Object) encodeElement(encoder *xml.Encoder, prop *specs.Property) error {
	if prop.Label == labels.Repeated {
		return encodeRepeated(encoder, object.resource, prop, object.refs)
	}

	// TODO: hide empty nested objects
	if prop.Type == types.Message {
		return encodeNested(encoder, prop, object.refs)
	}

	return encodeValue(encoder, prop, object.refs, true)
}

type decodeState int

const (
	waitForStart decodeState = iota
	waitForValue
	waitForClose
)

func spaces(i int) string {
	var str string

	for j := 0; j < i; j++ {
		str += "   "
	}

	return str
}

func (object *Object) UnmarshalXML(decoder *xml.Decoder, _ xml.StartElement) error {
	return object.unmarshalXML(decoder, nil, waitForStart, 0)
}

// <mock><repeating><value>repeating one</value></repeating><repeating><value>repeating two</value></repeating></mock>

func (object *Object) unmarshalXML(decoder *xml.Decoder, prop *specs.Property, state decodeState, level int) error {
	for {
		tok, err := decoder.Token()
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		log.Printf("%sOBJ[%d]: %T: %s", spaces(level), state, tok, tok)

		switch state {
		case waitForStart:
			switch t := tok.(type) {
			case xml.StartElement:
				var ok bool
				prop, ok = object.specs[t.Name.Local]
				if !ok {
					return fmt.Errorf("unknown property %q", t.Name.Local)
				}

				if prop.Label == labels.Repeated {
					var ref = &references.Reference{ // TODO
						Path: prop.Path,
					}

					var array = NewArray(object.resource, prop, ref, nil)
					if err := array.unmarshalXML(decoder, waitForValue, level+1); err != nil {
						return err
					}

					object.refs.StoreReference(object.resource, ref)

					continue
				}

				if prop.Type == types.Message {
					var nested = NewObject(object.resource, prop.Nested, object.refs)

					if err := nested.unmarshalXML(decoder, nil, waitForStart, level+1); err != nil {
						return err
					}

					continue
				}

				state = waitForValue
			case xml.EndElement:
				// object is closed
				return nil
			default:
				return fmt.Errorf("unexpected token type %T", t)
			}
		case waitForValue:
			switch t := tok.(type) {
			case xml.CharData:
				if err := decodeValue(prop, object.resource, t, object.refs); err != nil {
					return err
				}

				state = waitForClose
			default:
				return fmt.Errorf("unexpected token type %T", t)
			}
		case waitForClose:
			switch t := tok.(type) {
			case xml.EndElement:
				state = waitForStart
			default:
				return fmt.Errorf("unexpected token type %T", t)
			}
		}

	}

	return nil
}

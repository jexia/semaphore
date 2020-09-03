package xml

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log"

	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/types"
)

// Array represents an array of values/references.
type Array struct {
	resource string
	specs    *specs.Property
	items    []references.Store
	ref      *references.Reference
}

// NewArray constructs a new JSON array encoder/decoder
func NewArray(resource string, object *specs.Property, ref *references.Reference, refs []references.Store) *Array {
	return &Array{
		resource: resource,
		specs:    object,
		items:    refs,
		ref:      ref,
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
	if array.specs.Type == types.Message {
		return encodeNested(encoder, array.specs, store)
	}

	return encodeValue(encoder, array.specs, store, false)
}

// UnmarshalXML unmarshals the given specs into the configured reference store.
func (array *Array) UnmarshalXML(decoder *xml.Decoder, _ xml.StartElement) error {
	return array.unmarshalXML(decoder, waitForStart, 0)
}

func (array *Array) unmarshalXML(decoder *xml.Decoder, state decodeState, level int) error {
	for {
		tok, err := decoder.Token()
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		log.Printf("%sARR[%d]: %T: %s", spaces(level), state, tok, tok)

		switch state {
		case waitForStart:
			switch t := tok.(type) {
			case xml.StartElement:
				// TODO: check if name is still the same
				state = waitForValue
			case xml.EndElement:
				// no more elements
				return nil
			default:
				return fmt.Errorf("unexpected token type %T", t)
			}

		case waitForValue:
			switch t := tok.(type) {
			case xml.StartElement:
				if array.specs.Type != types.Message {
					return errors.New("TODO: not an object")
				}

				// log.Println()
				// log.Println()
				// log.Println(array.specs.Nested)
				// log.Println()
				// log.Println()

				var (
					store  = references.NewReferenceStore(0)
					nested = NewObject(array.resource, array.specs.Nested, store)
				)

				if err := nested.unmarshalXML(decoder, array.specs.Nested[t.Name.Local], waitForValue, level+1); err != nil {
					return err
				}

				array.ref.Append(store)

				state = waitForValue

			case xml.CharData:
				var store = references.NewReferenceStore(0)

				if array.specs.Type == types.Enum {
					log.Println("enum", string(t), array.specs.Enum.Keys)

					enum, ok := array.specs.Enum.Keys[string(t)]
					if !ok {
						return fmt.Errorf("unknown enum %s", t)
					}

					log.Println("position", enum.Position)

					store.StoreEnum("", "", enum.Position)
					array.ref.Append(store)

					state = waitForClose

					continue
				}

				value, err := DecodeType(string(t), array.specs.Type)
				if err != nil {
					return err
				}

				store.StoreValue("", "", value)
				array.ref.Append(store)

				state = waitForClose

				continue
			case xml.EndElement:
				continue
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

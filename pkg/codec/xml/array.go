package xml

import (
	"encoding/xml"

	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/types"
)

// Array represents an array of values/references.
type Array struct {
	// parent Entry

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

// // UnmarshalXML unmarshals the given specs into the configured reference store.
// func (array *Array) UnmarshalXML(decoder *xml.Decoder, _ xml.StartElement) error {
// 	return array.unmarshalXML(decoder)
// }
//
// func (array *Array) unmarshalXML(decoder *xml.Decoder) error {
// 	for {
// 		_, err := decoder.Token()
// 		if err == io.EOF {
// 			break
// 		}
//
// 		if err != nil {
// 			return err
// 		}
//
// 	}
//
// 	return nil
// }

package xml

import (
	"encoding/xml"
	"sort"

	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/labels"
	"github.com/jexia/semaphore/pkg/specs/types"
)

// NewObject constructs a new object encoder/decoder for the given specs
func NewObject(resource string, specs map[string]*specs.Property, refs references.Store) *Object {
	return &Object{
		resource: resource,
		refs:     refs,
		specs:    specs,
	}
}

// Object represents a JSON object
type Object struct {
	resource string
	specs    map[string]*specs.Property
	refs     references.Store
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

package xml

import (
	"encoding/xml"

	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/types"
)

func encodeRepeated(encoder *xml.Encoder, resource string, prop *specs.Property, store references.Store) error {
	if prop.Reference == nil {
		return nil
	}

	var ref = store.Load(prop.Reference.Resource, prop.Reference.Path)
	if ref == nil {
		return nil
	}

	var array = NewArray(resource, prop, ref.Repeated)

	return array.MarshalXML(encoder, xml.StartElement{})
}

func encodeNested(encoder *xml.Encoder, prop *specs.Property, store references.Store) error {
	// ignore malformed objects
	if prop.Nested == nil {
		return nil
	}

	var (
		nested = NewObject(prop.Name, prop, store)
		start  = xml.StartElement{Name: xml.Name{Local: prop.Name}}
	)

	return nested.MarshalXML(encoder, start)
}

func encodeValue(encoder *xml.Encoder, prop *specs.Property, store references.Store, loadByPath bool) error {
	var val = prop.Default

	if prop.Reference != nil {
		var ref *references.Reference

		if loadByPath {
			ref = store.Load(prop.Reference.Resource, prop.Reference.Path)
		} else {
			ref = store.Load("", "")
		}

		if ref != nil {
			if prop.Type == types.Enum && ref.Enum != nil {
				var enum = prop.Enum.Positions[*ref.Enum]
				if enum != nil {
					val = enum.Key
				}
			} else if ref.Value != nil {
				val = ref.Value
			}
		}
	}

	if val == nil {
		return nil
	}

	var start = xml.StartElement{Name: xml.Name{Local: prop.Name}}

	return encoder.EncodeElement(val, start)
}

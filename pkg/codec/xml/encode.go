package xml

import (
	"encoding/xml"
	"fmt"

	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
)

func encodeElement(encoder *xml.Encoder, name string, template specs.Template, store references.Store) error {
	var marshaler xml.Marshaler

	switch {
	case template.Message != nil:
		marshaler = NewObject(name, template.Message, store)
	case template.Repeated != nil:
		schema, err := template.Repeated.Template()
		if err != nil {
			return err
		}

		marshaler = NewArray(name, schema, template.Repeated, template.Reference, store)
	case template.Enum != nil:
		marshaler = NewEnum(name, template.Enum, template.Reference, store)
	case template.Scalar != nil:
		marshaler = NewScalar(name, template.Scalar, template.Reference, store)
	default:
		return fmt.Errorf("property '%s' has unknown type", name)
	}

	return marshaler.MarshalXML(encoder, xml.StartElement{})
}

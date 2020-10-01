package xml

import (
	"encoding/xml"
	"fmt"

	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
)

func decodeElement(decoder *xml.Decoder, start xml.StartElement, resource, prefix, name string, template specs.Template, store references.Store) (err error) {
	defer func() {
		if err != nil {
			err = errFailedToDecodeProperty{
				property: name,
				inner:    err,
			}
		}
	}()

	var unmarshaler xml.Unmarshaler

	switch {
	case template.Message != nil:
		unmarshaler = NewObject(resource, buildPath(prefix, name), name, template.Message, store)
	case template.Repeated != nil:
		schema, err := template.Repeated.Template()
		if err != nil {
			return err
		}

		unmarshaler = NewArray(resource, prefix, name, schema, template.Repeated, template.Reference, store)
	case template.Enum != nil:
		unmarshaler = NewEnum(resource, prefix, name, template.Enum, template.Reference, store)
	case template.Scalar != nil:
		unmarshaler = NewScalar(resource, prefix, name, template.Scalar, template.Reference, store)
	default:
		return fmt.Errorf("property '%s' has unknown type", name)
	}

	return unmarshaler.UnmarshalXML(decoder, start)
}

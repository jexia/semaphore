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
			err = errFailedToDecode{
				errStack{
					property: name,
					inner:    err,
				},
			}
		}
	}()

	var (
		unmarshaler xml.Unmarshaler
		reference   = specs.PropertyReference{
			Resource: resource,
			Path:     prefix,
		}
	)

	switch {
	case template.Message != nil:
		unmarshaler = NewObject(name, template.Message, &reference, store)
	case template.Repeated != nil:
		schema, err := template.Repeated.Template()
		if err != nil {
			return err
		}

		unmarshaler = NewArray(name, schema, template.Repeated, &reference, store)
	case template.Enum != nil:
		unmarshaler = NewEnum(name, template.Enum, &reference, store)
	case template.Scalar != nil:
		unmarshaler = NewScalar(name, template.Scalar, &reference, store)
	default:
		return fmt.Errorf("property '%s' has unknown type", name)
	}

	return unmarshaler.UnmarshalXML(decoder, start)
}

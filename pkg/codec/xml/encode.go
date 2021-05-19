package xml

import (
	"encoding/xml"
	"fmt"

	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
)

func encodeElement(encoder *xml.Encoder, name, path string, template specs.Template, store references.Store, tracker references.Tracker) (err error) {
	var marshaler xml.Marshaler

	switch {
	case template.Message != nil:
		marshaler = NewObject(name, path, template, store, tracker)
	case template.Repeated != nil:
		marshaler = NewArray(name, path, template, store, tracker)
	case template.OneOf != nil:
		marshaler = NewOneOf(name, path, template, store, tracker)
	case template.Enum != nil:
		marshaler = NewEnum(name, path, template, store, tracker)
	case template.Scalar != nil:
		marshaler = NewScalar(name, path, template, store, tracker)
	default:
		return fmt.Errorf("property '%s' has unknown type", path)
	}

	return marshaler.MarshalXML(encoder, xml.StartElement{})
}

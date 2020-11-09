package xml

import (
	"encoding/xml"
	"fmt"

	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
)

func decodeElement(decoder *xml.Decoder, start xml.StartElement, name, path string, template specs.Template, store references.Store, tracker references.Tracker) (err error) {
	switch {
	case template.Message != nil:
		return NewObject(name, path, template, store, tracker).UnmarshalXML(decoder, start)
	case template.Repeated != nil:
		return NewArray(name, path, template, store, tracker).UnmarshalXML(decoder, start)
	case template.Enum != nil:
		return NewEnum(name, path, template, store, tracker).UnmarshalXML(decoder, start)
	case template.Scalar != nil:
		return NewScalar(name, path, template, store, tracker).UnmarshalXML(decoder, start)
	default:
		return fmt.Errorf("property '%s' has unknown type", name)
	}
}

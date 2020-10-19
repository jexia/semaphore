package xml

import (
	"encoding/xml"
	"io"

	"github.com/jexia/semaphore/pkg/codec"
	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
)

// NewConstructor constructs a new XML constructor.
func NewConstructor() *Constructor { return &Constructor{} }

// Constructor is capable of constructing new codec managers for the given resource
// and specs.
type Constructor struct{}

// Name returns the name of the XML codec constructor.
func (constructor *Constructor) Name() string { return "xml" }

// New constructs a new XML codec manager
func (constructor *Constructor) New(resource string, specs *specs.ParameterMap) (codec.Manager, error) {
	if specs == nil {
		return nil, errNoSchema
	}

	manager := &Manager{
		resource: resource,
		property: specs.Property,
	}

	return manager, nil
}

// Manager manages a specs object and allows to encode/decode messages.
type Manager struct {
	resource string
	property *specs.Property
}

// Name returns the codec name.
func (manager *Manager) Name() string { return "xml" }

// Property returns the manager property which is used to marshal and unmarshal data.
func (manager *Manager) Property() *specs.Property { return manager.property }

// Marshal marshals the given reference store into a XML message.
// This method is called during runtime to encode a new message with the values
// stored inside the given reference store.
func (manager *Manager) Marshal(refs references.Store) (io.Reader, error) {
	if manager.property == nil {
		return nil, nil
	}

	var (
		reader, writer = io.Pipe()
		encoder        = xml.NewEncoder(writer)
	)

	go func() {
		if err := encodeElement(
			encoder,
			manager.property.Name,
			manager.property.Template,
			refs,
		); err != nil {
			_ = writer.CloseWithError(err)

			return
		}

		if err := encoder.Flush(); err != nil {
			_ = writer.CloseWithError(err)

			return
		}

		writer.Close()
	}()

	return reader, nil
}

// Unmarshal unmarshals the given XML io.Reader into the given reference store.
// This method is called during runtime to decode a new message and store it inside
// the given reference store.
func (manager *Manager) Unmarshal(reader io.Reader, refs references.Store) error {
	if manager.property == nil {
		return nil
	}

	var decoder = xml.NewDecoder(reader)

	for {
		tok, err := decoder.Token()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}

		switch t := tok.(type) {
		case xml.StartElement:
			if err := decodeElement(
				decoder,
				t,
				manager.resource,
				"", // prefix
				manager.property.Name,
				manager.property.Template,
				refs,
			); err != nil {
				return err
			}

			continue
		case xml.CharData:
			// ignore "\n", "\t" ...
			continue
		case xml.EndElement:
			// stream is closed
			return nil
		default:
			return errUnexpectedToken{
				actual: t,
				expected: []xml.Token{
					xml.StartElement{},
				},
			}
		}
	}
}

func buildPath(prefix, property string) string {
	if prefix == "" {
		return property
	}

	return prefix + "." + property
}

package xml

import (
	"bytes"
	"encoding/xml"
	"io"
	"io/ioutil"

	"github.com/jexia/semaphore/pkg/broker/trace"
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
		return nil, trace.New(trace.WithMessage("no object specs defined"))
	}

	return &Manager{
		resource: resource,
		specs:    specs.Property,
	}, nil
}

// Manager manages a specs object and allows to encode/decode messages.
type Manager struct {
	resource string
	specs    *specs.Property
}

// Name returns the codec name.
func (manager *Manager) Name() string { return "xml" }

// Property returns the manager property which is used to marshal and unmarshal data.
func (manager *Manager) Property() *specs.Property { return manager.specs }

// Marshal marshals the given reference store into a XML message.
// This method is called during runtime to encode a new message with the values
// stored inside the given reference store.
func (manager *Manager) Marshal(refs references.Store) (io.Reader, error) {
	if manager.specs == nil {
		return nil, nil
	}

	var object = NewObject(manager.resource, manager.specs.Repeated, refs)

	bb, err := xml.Marshal(object)
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(bb), nil
}

// Unmarshal unmarshals the given XML io.Reader into the given reference store.
// This method is called during runtime to decode a new message and store it inside
// the given reference store.
func (manager *Manager) Unmarshal(reader io.Reader, refs references.Store) error {
	if manager.specs == nil {
		return nil
	}

	bb, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}

	if len(bb) == 0 {
		return nil
	}

	var object = NewObject(manager.resource, manager.specs.Repeated, refs)

	return xml.Unmarshal(bb, object)
}

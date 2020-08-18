package formencoded

import (
	"io"

	"github.com/jexia/semaphore/pkg/broker/trace"
	"github.com/jexia/semaphore/pkg/codec"
	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
)

// Name represents the codec
const Name = "form-urlencoded"

// NewConstructor constructs a new www-form-urlencoded constructor
func NewConstructor() *Constructor {
	return &Constructor{}
}

// Constructor is capable of constructing new codec managers for the given resource and specs
type Constructor struct {
}

// Name returns the name of the www-form-urlencoded codec constructor
func (constructor *Constructor) Name() string {
	return Name
}

// New constructs a new www-form-urlencoded codec manager
func (constructor *Constructor) New(resource string, specs *specs.ParameterMap) (codec.Manager, error) {
	if specs == nil {
		return nil, trace.New(trace.WithMessage("no object specs defined"))
	}

	return &Manager{
		resource: resource,
		specs:    specs.Property,
	}, nil
}

// Manager manages a specs object and allows to encode/decode messages
type Manager struct {
	resource string
	specs    *specs.Property
}

// Name returns the proto codec name
func (manager *Manager) Name() string {
	return Name
}

// Property returns the manager property which is used to marshal and unmarshal data
func (manager *Manager) Property() *specs.Property {
	return manager.specs
}

// Marshal marshals the given reference store into a www-form-urlencoded message.
// This method is called during runtime to encode a new message with the values stored inside the given reference store
func (manager *Manager) Marshal(refs references.Store) (io.Reader, error) {
	return nil, nil
}

// Unmarshal unmarshals the given www-form-urlencoded io reader into the given reference store.
// This method is called during runtime to decode a new message and store it inside the given reference store
func (manager *Manager) Unmarshal(reader io.Reader, refs references.Store) error {
	return nil
}

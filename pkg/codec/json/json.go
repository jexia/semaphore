package json

import (
	"bytes"
	"io"
	"io/ioutil"

	"github.com/francoispqt/gojay"
	"github.com/jexia/maestro/pkg/codec"
	"github.com/jexia/maestro/pkg/refs"
	"github.com/jexia/maestro/pkg/specs"
	"github.com/jexia/maestro/pkg/specs/trace"
)

// NewConstructor constructs a new JSON constructor
func NewConstructor() *Constructor {
	return &Constructor{}
}

// Constructor is capable of constructing new codec managers for the given resource and specs
type Constructor struct {
}

// Name returns the name of the JSON codec constructor
func (constructor *Constructor) Name() string {
	return "json"
}

// New constructs a new JSON codec manager
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
	keys     int
}

// Property returns the manager property which is used to marshal and unmarshal data
func (manager *Manager) Property() *specs.Property {
	return manager.specs
}

// Marshal marshals the given reference store into a JSON message.
// This method is called during runtime to encode a new message with the values stored inside the given reference store
func (manager *Manager) Marshal(refs refs.Store) (io.Reader, error) {
	if manager.specs == nil {
		return nil, nil
	}

	object := NewObject(manager.resource, manager.specs.Nested, refs)
	bb, err := gojay.MarshalJSONObject(object)
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(bb), nil
}

// Unmarshal unmarshals the given JSON io reader into the given reference store.
// This method is called during runtime to decode a new message and store it inside the given reference store
func (manager *Manager) Unmarshal(reader io.Reader, refs refs.Store) error {
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

	object := NewObject(manager.resource, manager.specs.Nested, refs)
	err = gojay.UnmarshalJSONObject(bb, object)
	if err != nil {
		return err
	}

	return nil
}

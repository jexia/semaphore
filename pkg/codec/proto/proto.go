package proto

import (
	"bytes"
	"io"
	"io/ioutil"

	"github.com/jexia/semaphore/pkg/codec"
	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/template"
	"github.com/jexia/semaphore/pkg/specs/types"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"
)

// NewConstructor constructs a new Proto constructor
func NewConstructor() *Constructor {
	return &Constructor{}
}

// Constructor is capable of constructing new codec managers for the given resource and specs
type Constructor struct {
}

// Name returns the proto codec constructor name
func (constructor *Constructor) Name() string {
	return "proto"
}

// New constructs a new proto codec manager
func (constructor *Constructor) New(resource string, specs *specs.ParameterMap) (codec.Manager, error) {
	if specs == nil {
		return nil, ErrUndefinedSpecs{}
	}

	property := specs.Property
	if property == nil {
		return nil, nil
	}

	if property.Type() != types.Message {
		return nil, ErrNonRootMessage{}
	}

	desc, err := NewMessage(resource, property.Message)
	if err != nil {
		return nil, err
	}

	return &Manager{
		resource: resource,
		specs:    specs.Property,
		desc:     desc,
	}, nil
}

// Manager represents a proto message encoder/decoder
type Manager struct {
	resource string
	specs    *specs.Property
	desc     *desc.MessageDescriptor
}

// Name returns the proto codec name
func (manager *Manager) Name() string {
	return "proto"
}

// Property returns the property used to marshal and unmarshal data
func (manager *Manager) Property() *specs.Property {
	return manager.specs
}

// Marshal marshals the given reference store into a proto message.
// This method is called during runtime to encode a new message with the values stored inside the given reference store.
func (manager *Manager) Marshal(store references.Store) (io.Reader, error) {
	if manager.specs == nil {
		return nil, nil
	}

	tracker := references.NewTracker()
	result := dynamic.NewMessage(manager.desc)
	err := Message(manager.specs.Template).Marshal(result, manager.desc, template.ResourcePath(manager.resource), store, tracker)
	if err != nil {
		return nil, err
	}

	bb, err := result.Marshal()
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(bb), nil
}

// Unmarshal unmarshals the given io reader into the given reference store.
// This method is called during runtime to decode a new message and store it inside the given reference store
func (manager *Manager) Unmarshal(reader io.Reader, refs references.Store) error {
	if manager.specs == nil {
		return nil
	}

	bb, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}

	result := dynamic.NewMessage(manager.desc)
	err = result.Unmarshal(bb)
	if err != nil {
		return err
	}

	tracker := references.NewTracker()
	Message(manager.specs.Template).Unmarshal(result, template.ResourcePath(manager.resource), refs, tracker)

	return nil
}

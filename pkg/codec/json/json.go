package json

import (
	"bufio"
	"bytes"
	"fmt"
	"io"

	"github.com/francoispqt/gojay"
	"github.com/jexia/semaphore/pkg/codec"
	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
)

// Constructor is capable of constructing new codec managers for the given resource and specs.
type Constructor struct{}

// NewConstructor constructs a new JSON constructor.
func NewConstructor() *Constructor { return &Constructor{} }

// Name returns the name of the JSON codec constructor.
func (constructor *Constructor) Name() string { return "json" }

// Manager manages a specs object and allows to encode/decode messages.
type Manager struct {
	resource string
	property *specs.Property
	root     string
}

// New constructs a new JSON codec manager
func (constructor *Constructor) New(resource string, specs *specs.ParameterMap) (codec.Manager, error) {
	if specs == nil {
		return nil, ErrUndefinedSpecs{}
	}

	manager := &Manager{
		resource: resource,
		property: specs.Property,
		root:     fmt.Sprintf("%s:", resource),
	}

	return manager, nil
}

// Name returns the proto codec name
func (manager *Manager) Name() string { return "json" }

// Property returns the manager property which is used to marshal and unmarshal data
func (manager *Manager) Property() *specs.Property { return manager.property }

// Marshal marshals the given reference store into a JSON message.
// This method is called during runtime to encode a new message with the values stored inside the given reference store
func (manager *Manager) Marshal(store references.Store) (io.Reader, error) {
	if manager.property == nil {
		return nil, nil
	}

	var (
		bb      = &bytes.Buffer{}
		encoder = gojay.BorrowEncoder(bb)
		tracker = references.NewTracker()
	)

	encode(encoder, "", manager.property.Template, store, tracker)

	defer encoder.Release()
	_, err := encoder.Write()
	if err != nil {
		return nil, err
	}

	return bb, nil
}

// Unmarshal unmarshals the given JSON io reader into the given reference store.
// This method is called during runtime to decode a new message and store it inside the given reference store
func (manager *Manager) Unmarshal(reader io.Reader, store references.Store) error {
	if manager.property == nil {
		return nil
	}

	var (
		buff    = bufio.NewReader(reader)
		tracker = references.NewTracker()
		_, err  = buff.ReadByte()
	)

	if err == io.EOF {
		return nil
	}

	_ = buff.UnreadByte()

	var decoder = gojay.NewDecoder(buff)
	defer decoder.Release()

	return decode(
		decoder,
		manager.root,
		manager.property.Template,
		store,
		tracker,
	)
}

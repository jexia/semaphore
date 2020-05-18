package proto

import (
	"bytes"
	"io"
	"io/ioutil"

	"github.com/jexia/maestro/pkg/codec"
	"github.com/jexia/maestro/pkg/refs"
	"github.com/jexia/maestro/pkg/specs"
	"github.com/jexia/maestro/pkg/specs/trace"
	"github.com/jexia/maestro/pkg/specs/types"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"
)

// NewConstructor constructs a new JSON constructor
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
		return nil, trace.New(trace.WithMessage("no object specs defined"))
	}

	prop := specs.Property
	if prop.Type != types.Message {
		return nil, trace.New(trace.WithMessage("a proto message always requires a root message"))
	}

	desc, err := NewMessage(resource, prop.Nested)
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
func (manager *Manager) Marshal(refs refs.Store) (io.Reader, error) {
	if manager.specs == nil {
		return nil, nil
	}

	result := dynamic.NewMessage(manager.desc)
	err := manager.Encode(result, manager.desc, manager.specs.Nested, refs)
	if err != nil {
		return nil, err
	}

	bb, err := result.Marshal()
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(bb), nil
}

// Encode encodes the given specs object into the given dynamic proto message.
// References inside the specs are attempted to be fetched from the reference store.
func (manager *Manager) Encode(proto *dynamic.Message, desc *desc.MessageDescriptor, specs map[string]*specs.Property, store refs.Store) (err error) {
	for _, field := range desc.GetFields() {
		prop, has := specs[field.GetName()]
		if !has {
			continue
		}

		if field.IsRepeated() {
			if prop.Reference == nil {
				continue
			}

			ref := store.Load(prop.Reference.Resource, prop.Reference.Path)
			if ref == nil {
				continue
			}

			for _, store := range ref.Repeated {
				var value interface{}

				switch prop.Type {
				case types.Message:
					item := dynamic.NewMessage(field.GetMessageType())
					err = manager.Encode(item, field.GetMessageType(), prop.Nested, store)
					if err != nil {
						return err
					}

					value = item
				case types.Enum:
					ref := store.Load("", "")
					if ref == nil || ref.Enum == nil {
						continue
					}

					value = *ref.Enum
				default:
					ref := store.Load("", "")
					value = ref.Value
				}

				err = proto.TryAddRepeatedField(field, value)
				if err != nil {
					return err
				}
			}

			continue
		}

		if prop.Type == types.Message {
			dynamic := dynamic.NewMessage(field.GetMessageType())
			err = manager.Encode(dynamic, field.GetMessageType(), prop.Nested, store)
			if err != nil {
				return err
			}

			err = proto.TrySetField(field, dynamic)
			if err != nil {
				return err
			}

			continue
		}

		value := prop.Default

		if prop.Reference != nil {
			ref := store.Load(prop.Reference.Resource, prop.Reference.Path)
			if ref != nil {
				if prop.Type == types.Enum && ref.Enum != nil {
					value = ref.Enum
				}

				if value == nil {
					value = ref.Value
				}
			}
		}

		if value == nil {
			continue
		}

		err = proto.TrySetField(field, value)
		if err != nil {
			return err
		}
	}

	return nil
}

// Unmarshal unmarshals the given io reader into the given reference store.
// This method is called during runtime to decode a new message and store it inside the given reference store
func (manager *Manager) Unmarshal(reader io.Reader, refs refs.Store) error {
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

	manager.Decode(result, manager.specs.Nested, refs)
	return nil
}

// Decode decodes the given proto message into the given reference store.
func (manager *Manager) Decode(proto *dynamic.Message, properties map[string]*specs.Property, store refs.Store) {
	for _, field := range proto.GetKnownFields() {
		prop := properties[field.GetName()]

		if field.IsRepeated() {
			length := proto.FieldLength(field)

			ref := refs.NewReference(prop.Path)
			ref.Repeating(length)

			for index := 0; index < length; index++ {
				value := proto.GetRepeatedField(field, index)

				if prop.Type == types.Message {
					message := value.(*dynamic.Message)
					store := refs.NewReferenceStore(len(message.GetKnownFields()))
					manager.Decode(message, prop.Nested, store)
					ref.Set(index, store)
					continue
				}

				store := refs.NewReferenceStore(1)

				if prop.Type == types.Enum {
					enum, is := value.(int32)
					if !is {
						continue
					}

					store.StoreEnum("", "", enum)
					ref.Set(index, store)
					continue
				}

				store.StoreValue("", "", value)
				ref.Set(index, store)
			}

			store.StoreReference(manager.resource, ref)
			continue
		}

		if prop.Type == types.Message {
			nested := proto.GetField(field).(*dynamic.Message)
			manager.Decode(nested, prop.Nested, store)
			continue
		}

		value := proto.GetField(field)
		ref := refs.NewReference(prop.Path)

		if prop.Type == types.Enum {
			enum, is := value.(int32)
			if !is {
				continue
			}

			ref.Enum = &enum
			store.StoreReference(manager.resource, ref)
			continue
		}

		ref.Value = value
		store.StoreReference(manager.resource, ref)
	}
}

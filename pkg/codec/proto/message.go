package proto

import (
	"bytes"
	"io"
	"io/ioutil"

	"github.com/jexia/semaphore/pkg/codec"
	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
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
func (manager *Manager) Marshal(refs references.Store) (io.Reader, error) {
	if manager.specs == nil {
		return nil, nil
	}

	result := dynamic.NewMessage(manager.desc)
	err := manager.Encode(result, manager.desc, manager.specs.Message, refs)
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
func (manager *Manager) Encode(proto *dynamic.Message, desc *desc.MessageDescriptor, specs specs.Message, store references.Store) (err error) {
	if specs == nil {
		return
	}

	for _, field := range desc.GetFields() {
		property := specs[field.GetName()]
		if property == nil {
			continue
		}

		switch {
		case field.IsRepeated():
			err = manager.setRepeating(proto, property.Template, field, store)
			if err != nil {
				return err
			}

			break
		default:
			err = manager.setField(proto.TrySetField, property.Template, field, store)
			if err != nil {
				return err
			}
		}

	}

	return nil
}

func (manager *Manager) setRepeating(message *dynamic.Message, template specs.Template, field *desc.FieldDescriptor, store references.Store) error {
	// TODO: implement static values
	// 	if prop.Reference == nil {
	// 		for _, repeated := range prop.Repeated {
	// 			err = manager.setField(proto.TryAddRepeatedField, repeated, field, store)
	// 			if err != nil {
	// 				return err
	// 			}
	// 		}

	// 		continue
	// 	}

	if template.Reference == nil {
		return nil
	}

	ref := store.Load(template.Reference.Resource, template.Reference.Path)
	if ref == nil {
		return nil
	}

	// TODO: generate and store repeated template upfront
	template, err := template.Repeated.Template()
	if err != nil {
		return err
	}

	for _, store := range ref.Repeated {
		var value interface{}

		switch template.Type() {
		case types.Message:
			item := dynamic.NewMessage(field.GetMessageType())
			err = manager.Encode(item, field.GetMessageType(), template.Message, store)
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
			if ref == nil {
				continue
			}

			value = ref.Value
		}

		err = message.TryAddRepeatedField(field, value)
		if err != nil {
			return err
		}
	}

	return nil
}

type trySetProto func(fd *desc.FieldDescriptor, val interface{}) error

func (manager *Manager) setField(setter trySetProto, template specs.Template, field *desc.FieldDescriptor, store references.Store) error {
	switch {
	case template.Message != nil:
		dynamic := dynamic.NewMessage(field.GetMessageType())
		err := manager.Encode(dynamic, field.GetMessageType(), template.Message, store)
		if err != nil {
			return err
		}

		return setter(field, dynamic)
	case template.Enum != nil:
		if template.Reference == nil {
			break
		}

		ref := store.Load(template.Reference.Resource, template.Reference.Path)
		if ref == nil || ref.Enum == nil {
			break
		}

		return setter(field, ref.Enum)
	case template.Scalar != nil:
		value := template.Scalar.Default

		if template.Reference != nil {
			ref := store.Load(template.Reference.Resource, template.Reference.Path)
			if ref != nil && ref.Value != nil {
				value = ref.Value
			}
		}

		if value == nil {
			break
		}

		return setter(field, value)
	}

	return nil
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

	manager.Decode(result, manager.specs.Message, refs)
	return nil
}

// Decode decodes the given proto message into the given reference store.
func (manager *Manager) Decode(protobuf *dynamic.Message, message specs.Message, store references.Store) {
	if message == nil {
		return
	}

	for _, field := range protobuf.GetKnownFields() {
		property := message[field.GetName()]
		if property == nil {
			continue
		}

		switch property.Type() {
		case types.Array:
			length := protobuf.FieldLength(field)

			ref := &references.Reference{
				Path: property.Path,
			}

			ref.Repeating(length)

			tmpl, err := property.Repeated.Template()
			if err != nil {
				panic(err)
			}

			for index := 0; index < length; index++ {
				value := protobuf.GetRepeatedField(field, index)

				switch tmpl.Type() {
				case types.Message:
					message := value.(*dynamic.Message)
					store := references.NewReferenceStore(len(message.GetKnownFields()))
					manager.Decode(message, tmpl.Message, store)
					ref.Set(index, store)
				case types.Enum:
					store := references.NewReferenceStore(1)
					enum, is := value.(int32)
					if !is {
						continue
					}

					store.StoreEnum("", "", enum)
					ref.Set(index, store)
				default:
					store := references.NewReferenceStore(1)

					store.StoreValue("", "", value)
					ref.Set(index, store)
				}
			}

			store.StoreReference(manager.resource, ref)
		case types.Message:
			nested := protobuf.GetField(field).(*dynamic.Message)
			manager.Decode(nested, property.Message, store)
		case types.Enum:
			value := protobuf.GetField(field)
			enum, is := value.(int32)
			if !is {
				break
			}

			store.StoreEnum(manager.resource, property.Path, enum)
		default:
			value := protobuf.GetField(field)
			store.StoreValue(manager.resource, property.Path, value)
		}
	}
}

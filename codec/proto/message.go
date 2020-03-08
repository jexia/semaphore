package proto

import (
	"bytes"
	"io"
	"io/ioutil"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/jexia/maestro/codec"
	"github.com/jexia/maestro/refs"
	"github.com/jexia/maestro/specs"
	"github.com/jexia/maestro/specs/trace"
	"github.com/jexia/maestro/specs/types"
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
func (constructor *Constructor) New(resource string, specs *specs.Property) (codec.Manager, error) {
	if specs == nil {
		return nil, trace.New(trace.WithMessage("no object specs defined"))
	}

	if specs.Type != types.TypeMessage {
		return nil, trace.New(trace.WithMessage("a proto message always requires a root message"))
	}

	desc, err := NewMessage(resource, specs.Nested)
	if err != nil {
		return nil, err
	}

	return &Message{
		resource: resource,
		specs:    specs,
		desc:     desc,
	}, nil
}

// Message represents a proto message encoder/decoder
type Message struct {
	resource string
	specs    *specs.Property
	desc     *desc.MessageDescriptor
}

// Marshal marshals the given reference store into a proto message.
// This method is called during runtime to encode a new message with the values stored inside the given reference store.
func (message *Message) Marshal(refs *refs.Store) (io.Reader, error) {
	result := dynamic.NewMessage(message.desc)
	message.Encode(result, message.desc, message.specs.Nested, refs)
	bb, err := result.Marshal()
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(bb), nil
}

// Encode encodes the given specs object into the given dynamic proto message.
// References inside the specs are attempted to be fetched from the reference store.
func (message *Message) Encode(proto *dynamic.Message, desc *desc.MessageDescriptor, specs map[string]*specs.Property, store *refs.Store) (err error) {
	for _, field := range desc.GetFields() {
		prop, has := specs[field.GetName()]
		if !has {
			continue
		}

		val := prop.Default

		if prop.Reference != nil {
			ref := store.Load(prop.Reference.Resource, prop.Reference.Path)
			if ref != nil {
				val = ref.Value
			}
		}

		if prop.Label == types.LabelRepeated {
			if prop.Reference == nil {
				continue
			}

			if prop.Type == types.TypeMessage {
				dynamic := dynamic.NewMessage(field.GetMessageType())

				err = message.Encode(dynamic, field.GetMessageType(), prop.Nested, store)
				if err != nil {
					return err
				}

				val = dynamic
			}

			err = proto.TryAddRepeatedField(field, val)
			if err != nil {
				return err
			}
		}

		if prop.Type == types.TypeMessage {
			dynamic := dynamic.NewMessage(field.GetMessageType())
			err = message.Encode(dynamic, field.GetMessageType(), prop.Nested, store)
			if err != nil {
				return err
			}

			err = proto.TrySetField(field, dynamic)
			if err != nil {
				return err
			}
		}

		err = proto.TrySetField(field, val)
		if err != nil {
			return err
		}
	}

	return nil
}

// Unmarshal unmarshals the given io reader into the given reference store.
// This method is called during runtime to decode a new message and store it inside the given reference store
func (message *Message) Unmarshal(reader io.Reader, refs *refs.Store) error {
	bb, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}

	result := dynamic.NewMessage(message.desc)
	err = result.Unmarshal(bb)
	if err != nil {
		return err
	}

	message.Decode(result, "", refs)
	return nil
}

// Decode decodes the given proto message into the given reference store.
func (message *Message) Decode(proto *dynamic.Message, origin string, store *refs.Store) {
	for _, field := range proto.GetKnownFields() {
		key := field.GetName()
		path := specs.JoinPath(origin, key)

		if field.GetType() == descriptor.FieldDescriptorProto_TYPE_MESSAGE {
			if field.IsRepeated() {
				length := proto.FieldLength(field)

				ref := refs.New(path)
				ref.Repeating(length)

				for index := 0; index < length; index++ {
					repeated := proto.GetRepeatedField(field, index).(*dynamic.Message)
					store := refs.NewStore(len(repeated.GetKnownFields()))
					message.Decode(repeated, path, store)
					ref.Set(index, store)
				}

				store.StoreReference(message.resource, ref)
				continue
			}

			nested := proto.GetField(field).(*dynamic.Message)
			message.Decode(nested, path, store)
			continue
		}

		if field.IsRepeated() {
			continue
		}

		value := proto.GetField(field)

		ref := refs.New(path)
		ref.Value = value

		store.StoreReference(message.resource, ref)
	}
}

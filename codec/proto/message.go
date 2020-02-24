package proto

import (
	"bytes"
	"io"
	"io/ioutil"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/jexia/maestro/codec"
	"github.com/jexia/maestro/refs"
	"github.com/jexia/maestro/schema"
	"github.com/jexia/maestro/schema/protoc"
	"github.com/jexia/maestro/specs"
	"github.com/jexia/maestro/specs/trace"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"
)

// ErrUnkownSchema is thrown when the schema is not defined or then it is not a protoc object
var ErrUnkownSchema = trace.New(trace.WithMessage("unexpected schema type, a proto schema collection required for protobuf encoding/decoding"))

// New constructs a new proto message encoder/decoder for the given schema object and specifications
func New(resource string, schema schema.Object, specs specs.Object) (codec.Manager, error) {
	if schema == nil {
		return nil, ErrUnkownSchema
	}

	object, is := schema.(protoc.Object)
	if !is {
		return nil, ErrUnkownSchema
	}

	message := &Message{
		resource:   resource,
		specs:      specs,
		schema:     object,
		descriptor: object.GetDescriptor(),
	}

	return message, nil
}

// Message represents a proto message encoder/decoder
type Message struct {
	resource   string
	specs      specs.Object
	schema     protoc.Object
	descriptor *desc.MessageDescriptor
}

// Marshal marshals the given reference store into a proto message.
// This method is called during runtime to encode a new message with the values stored inside the given reference store.
func (message *Message) Marshal(refs *refs.Store) (io.Reader, error) {
	result := dynamic.NewMessage(message.descriptor)
	message.Encode(result, message.schema, message.specs, refs)
	bb, err := result.Marshal()
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(bb), nil
}

// Encode encodes the given specs object into the given dynamic proto message.
// References inside the specs are attempted to be fetched from the reference store.
func (message *Message) Encode(proto *dynamic.Message, schema protoc.Object, specs specs.Object, store *refs.Store) (err error) {
	for key, prop := range specs.GetProperties() {
		field := schema.GetField(key).(protoc.Field)
		val := prop.Default

		if prop.Reference != nil {
			ref := store.Load(prop.Reference.Resource, prop.Reference.Path)
			if ref != nil {
				val = ref.Value
			}
		}

		if val == nil {
			continue
		}

		err = proto.TrySetField(field.GetDescriptor(), val)
		if err != nil {
			return err
		}
	}

	for key, nested := range specs.GetNestedProperties() {
		field := schema.GetField(key).(protoc.Field)
		dynamic := dynamic.NewMessage(field.GetDescriptor().GetMessageType())
		err = message.Encode(dynamic, field.GetObject().(protoc.Object), nested.GetObject(), store)
		if err != nil {
			return err
		}

		err = proto.TrySetField(field.GetDescriptor(), dynamic)
		if err != nil {
			return err
		}
	}

	for key, repeated := range specs.GetRepeatedProperties() {
		ref := store.Load(repeated.Template.Resource, repeated.Template.Path)
		if ref == nil {
			continue
		}

		field := schema.GetField(key).(protoc.Field)

		for _, store := range ref.Repeated {
			dynamic := dynamic.NewMessage(field.GetDescriptor().GetMessageType())
			err = message.Encode(dynamic, field.GetObject().(protoc.Object), repeated.GetObject(), store)
			if err != nil {
				return err
			}

			err = proto.TryAddRepeatedField(field.GetDescriptor(), dynamic)
			if err != nil {
				return err
			}
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

	result := dynamic.NewMessage(message.descriptor)
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

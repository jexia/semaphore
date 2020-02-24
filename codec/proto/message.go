package proto

import (
	"github.com/jexia/maestro/refs"
	"github.com/jexia/maestro/schema/protoc"
	"github.com/jexia/maestro/specs"
	"github.com/jhump/protoreflect/dynamic"
)

// MarshalMessage ...
func MarshalMessage(proto *dynamic.Message, schema protoc.Object, specs specs.Object, store *refs.Store) (err error) {
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
			val = field.GetDescriptor().GetDefaultValue()
		}

		err = proto.TrySetField(field.GetDescriptor(), val)
		if err != nil {
			return err
		}
	}

	for key, nested := range specs.GetNestedProperties() {
		field := schema.GetField(key).(protoc.Field)
		message := dynamic.NewMessage(field.GetDescriptor().GetMessageType())
		err = MarshalMessage(message, field.GetObject().(protoc.Object), nested.GetObject(), store)
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
			message := dynamic.NewMessage(field.GetDescriptor().GetMessageType())

			err = MarshalMessage(message, field.GetObject().(protoc.Object), repeated.GetObject(), store)
			if err != nil {
				return err
			}

			err = proto.TryAddRepeatedField(field.GetDescriptor(), message)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

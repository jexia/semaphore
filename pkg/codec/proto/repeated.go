package proto

import (
	"github.com/jexia/semaphore/v2/pkg/references"
	"github.com/jexia/semaphore/v2/pkg/specs"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"
)

// Repeated represents a repeated field with the given template
type Repeated specs.Template

// Marshal attempts to marshal the given template as a protobuffer repeated value
func (tmpl Repeated) Marshal(message *dynamic.Message, field *desc.FieldDescriptor, path string, store references.Store, tracker references.Tracker) error {
	if tmpl.Repeated == nil {
		return nil
	}

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

	if tmpl.Reference == nil {
		return nil
	}

	repeated, err := tmpl.Repeated.Template()
	if err != nil {
		panic(err)
	}

	rpath := tracker.Resolve(tmpl.Reference.String())
	tracker.Track(rpath, 0)

	length := store.Length(rpath)

	for index := 0; index < length; index++ {
		switch {
		case repeated.Message != nil:
			desc := field.GetMessageType()
			nested := dynamic.NewMessage(desc)

			err := Message(repeated).Marshal(nested, desc, path, store, tracker)
			if err != nil {
				return err
			}

			message.TryAddRepeatedField(field, nested)
		default:
			err := Field(repeated).Marshal(message.TryAddRepeatedField, field, store, tracker)
			if err != nil {
				return err
			}
		}

		tracker.Next(rpath)
	}

	return nil
}

// Unmarshal unmarshals the given repeated field into the given reference store.
func (tmpl Repeated) Unmarshal(protobuf *dynamic.Message, field *desc.FieldDescriptor, path string, store references.Store, tracker references.Tracker) {
	tpath := tracker.Resolve(path)

	length := protobuf.FieldLength(field)
	store.Define(tpath, length)
	tracker.Track(tpath, 0)

	for index := 0; index < length; index++ {
		value := protobuf.GetRepeatedField(field, index)

		switch {
		case tmpl.Message != nil:
			message := value.(*dynamic.Message)
			Message(tmpl).Unmarshal(message, path, store, tracker)
		default:
			Field(tmpl).Unmarshal(value, path, store, tracker)
		}

		tracker.Next(tpath)
	}
}

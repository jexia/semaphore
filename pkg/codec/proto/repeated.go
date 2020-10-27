package proto

import (
	"log"

	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
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

	length := store.Length(tracker.Resolve(path))
	tracker.Track(path, 0)

	for index := 0; index < length; index++ {
		log.Println(repeated.Type())
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
			log.Println("--", tmpl.Reference, err, length)
			err := Field(repeated).Marshal(message.TryAddRepeatedField, field, store, tracker)
			if err != nil {
				return err
			}
		}

		tracker.Next(path)
	}

	return nil
}

// Unmarshal unmarshals the given repeated field into the given reference store.
func (tmpl Repeated) Unmarshal(protobuf *dynamic.Message, field *desc.FieldDescriptor, path string, store references.Store, tracker references.Tracker) {
	length := protobuf.FieldLength(field)
	store.Define(tracker.Resolve(path), length)

	for index := 0; index < length; index++ {
		value := protobuf.GetRepeatedField(field, index)

		switch {
		case tmpl.Message != nil:
			message := value.(*dynamic.Message)
			Message(tmpl).Unmarshal(message, path, store, tracker)
		default:
			Field(tmpl).Unmarshal(protobuf, field, path, store, tracker)
		}
	}
}

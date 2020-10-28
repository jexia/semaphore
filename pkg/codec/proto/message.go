package proto

import (
	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/template"
	"github.com/jexia/semaphore/pkg/specs/types"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"
)

// Message represents a protobuffer message
type Message specs.Template

// Marshal attempts to encode the given template as a protobuf message using the
// given store.
func (tmpl Message) Marshal(result *dynamic.Message, message *desc.MessageDescriptor, root string, store references.Store, tracker references.Tracker) error {
	if tmpl.Message == nil {
		return nil
	}

	for _, field := range message.GetFields() {
		name := field.GetName()
		path := template.JoinPath(root, name)
		specs, has := tmpl.Message[name]
		if !has {
			continue
		}

		switch {
		case specs.Repeated != nil:
			err := Repeated(specs.Template).Marshal(result, field, path, store, tracker)
			if err != nil {
				return err
			}
		case specs.Message != nil:
			desc := field.GetMessageType()
			nested := dynamic.NewMessage(desc)

			length := store.Length(tracker.Resolve(path))
			if length == 0 {
				break
			}

			err := Message(specs.Template).Marshal(nested, desc, path, store, tracker)
			if err != nil {
				return err
			}

			err = result.TrySetField(field, nested)
			if err != nil {
				return err
			}
		default:
			err := Field(specs.Template).Marshal(result.TrySetField, field, store, tracker)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Unmarshal unmarshals the given protobuffer message into the given reference store.
func (tmpl Message) Unmarshal(protobuf *dynamic.Message, path string, store references.Store, tracker references.Tracker) {
	if tmpl.Message == nil {
		return
	}

	fields := protobuf.GetKnownFields()
	store.Define(tracker.Resolve(path), len(fields))

	for _, field := range fields {
		property := tmpl.Message[field.GetName()]
		if property == nil {
			continue
		}

		path := template.JoinPath(path, field.GetName())

		switch property.Type() {
		case types.Array:
			tmpl, err := property.Repeated.Template()
			if err != nil {
				panic(err)
			}

			Repeated(tmpl).Unmarshal(protobuf, field, path, store, tracker)
		case types.Message:
			nested := protobuf.GetField(field).(*dynamic.Message)
			Message(property.Template).Unmarshal(nested, path, store, tracker)
		default:
			Field(tmpl).Unmarshal(protobuf, field, path, store, tracker)
		}
	}
}

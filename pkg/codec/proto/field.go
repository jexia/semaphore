package proto

import (
	"log"

	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"
)

// TrySetter represents a protoreflect setter used to define various values
type TrySetter func(fd *desc.FieldDescriptor, val interface{}) error

// Field represents a protobuf field
type Field specs.Template

// Marshal attempts to encode the given field as a protobuf field using the given setter
func (tmpl Field) Marshal(setter TrySetter, field *desc.FieldDescriptor, store references.Store, tracker references.Tracker) error {
	switch {
	case tmpl.Enum != nil:
		if tmpl.Reference == nil {
			break
		}

		ref := store.Load(tmpl.Reference.String())
		log.Println(ref, store)
		if ref == nil || ref.Enum == nil {
			break
		}

		return setter(field, ref.Enum)
	case tmpl.Scalar != nil:
		value := tmpl.Scalar.Default

		if tmpl.Reference != nil {
			ref := store.Load(tmpl.Reference.String())
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

// Unmarshal unmarshals the given protobuffer field into the given reference store.
func (tmpl Field) Unmarshal(protobuf *dynamic.Message, field *desc.FieldDescriptor, path string, store references.Store, tracker references.Tracker) {
	switch {
	case tmpl.Enum != nil:
		value := protobuf.GetField(field)
		enum, is := value.(int32)
		if !is {
			break
		}

		store.Store(tracker.Resolve(path), &references.Reference{Enum: &enum})
	case tmpl.Scalar != nil:
		value := protobuf.GetField(field)
		store.Store(tracker.Resolve(path), &references.Reference{Value: value})
	}
}

package proto

import (
	"log"

	"github.com/jexia/semaphore/pkg/providers/protobuffers"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/types"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/builder"
)

// NewMessage attempts to construct a new proto message descriptor for the given specs property
func NewMessage(resource string, message specs.Message) (*desc.MessageDescriptor, error) {
	log.Println("BUILD THE MESSAGE")

	msg := builder.NewMessage(resource)
	err := ConstructMessage(msg, message)
	if err != nil {
		return nil, err
	}

	return msg.Build()
}

// ConstructMessage constructs a proto message of the given specs into the given message builders
func ConstructMessage(message *builder.MessageBuilder, spec specs.Message) (err error) {
	if spec == nil {
		return nil
	}

	// types := make(map[string]*builder.FieldType)

	// TODO: implement registry !!!

	for _, property := range spec {
		typed, err := ConstructFieldType(message, property.Name, property.Template)
		if err != nil {
			return err
		}

		label := protobuffers.ProtoLabels[property.Label]
		if property.Type() == types.Array {
			label = protobuffers.Repeated
		}

		field := builder.NewField(property.Name, typed)
		field.SetJsonName(property.Name)
		field.SetLabel(label)
		field.SetComments(builder.Comments{
			LeadingComment: property.Description,
		})

		err = message.TryAddField(field.SetNumber(property.Position))
		if err != nil {
			return err
		}
	}

	return nil
}

// ConstructFieldType constructs a field constructor from the given property
func ConstructFieldType(message *builder.MessageBuilder, key string, template specs.Template) (*builder.FieldType, error) {
	// builder.FromOneOf()
	// OneOfDescriptor

	log.Println("<<<", template.Identifier)

	switch {
	case template.OneOf != nil:
		oneof := builder.NewOneOf("test")

		_ = oneof

	case template.Message != nil:
		// TODO:
		// - reuse global types
		// - appending a fixed prefix is probably not a good idea
		// - do not create inline messages?
		nested := builder.NewMessage(key + "Nested")
		err := ConstructMessage(nested, template.Message)
		if err != nil {
			return nil, err
		}

		err = message.TryAddNestedMessage(nested)
		if err != nil {
			return nil, err
		}

		return builder.FieldTypeMessage(nested), nil
	case template.Repeated != nil:
		field, err := template.Repeated.Template()
		if err != nil {
			return nil, err
		}

		// TODO: throw an error when attempting to construct a nested array
		return ConstructFieldType(message, key, field)
	case template.Enum != nil:
		enum := builder.NewEnum(key + "Enum")

		for _, value := range template.Enum.Keys {
			eval := builder.NewEnumValue(value.Key)

			eval.SetNumber(value.Position)
			eval.SetComments(builder.Comments{
				LeadingComment: value.Description,
			})

			err := enum.TryAddValue(eval)
			if err != nil {
				return nil, err
			}
		}

		err := message.TryAddNestedEnum(enum)
		if err != nil {
			return nil, err
		}

		return builder.FieldTypeEnum(enum), nil
	case template.Scalar != nil:
		return builder.FieldTypeScalar(protobuffers.ProtoTypes[template.Scalar.Type]), nil
	}

	return nil, ErrInvalidFieldType{template}
}

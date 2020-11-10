package proto

import (
	"github.com/jexia/semaphore/pkg/providers/protobuffers"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/types"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/builder"
)

// NewMessage attempts to construct a new proto message descriptor for the given specs property
func NewMessage(resource string, message specs.Message) (*desc.MessageDescriptor, error) {
	var (
		msg = builder.NewMessage(resource)
		err = ConstructMessage(msg, message)
	)

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

	for _, property := range spec {
		if property.OneOf != nil {
			if err := ConstructOneOf(message, property.Name, property.Template); err != nil {
				return err
			}

			continue
		}

		field, err := newFieldBuilder(message, property)
		if err != nil {
			return err
		}

		if err := message.TryAddField(field); err != nil {
			return err
		}
	}

	return nil
}

func newFieldBuilder(message *builder.MessageBuilder, property *specs.Property) (*builder.FieldBuilder, error) {
	fieldType, err := ConstructFieldType(message, property.Name, property.Template)
	if err != nil {
		return nil, err
	}

	field := builder.NewField(property.Name, fieldType)

	label := protobuffers.ProtoLabels[property.Label]
	if property.Type() == types.Array {
		label = protobuffers.Repeated
	}

	field.SetJsonName(property.Name)
	field.SetLabel(label)
	field.SetComments(builder.Comments{
		LeadingComment: property.Description,
	})

	return field.SetNumber(property.Position), nil
}

// ConstructFieldType constructs a field constructor from the given property
func ConstructFieldType(message *builder.MessageBuilder, key string, template specs.Template) (*builder.FieldType, error) {
	switch {
	case template.Message != nil:
		// TODO: appending a fixed prefix is probably not a good idea.
		nested := builder.NewMessage(key + "Type")
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

		// TODO: thrown a error when attempting to construct a nested array
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

func ConstructOneOf(message *builder.MessageBuilder, key string, template specs.Template) error {
	var oneOf = builder.NewOneOf(key)

	for _, property := range template.OneOf {
		field, err := newFieldBuilder(message, property)
		if err != nil {
			return err
		}

		if err := oneOf.TryAddChoice(field); err != nil {
			return err
		}
	}

	return message.TryAddOneOf(oneOf)
}

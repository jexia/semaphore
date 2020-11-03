package proto

import (
	"github.com/jexia/semaphore/pkg/providers/protobuffers"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/labels"
	"github.com/jexia/semaphore/pkg/specs/types"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/builder"
)

// NewMessage attempts to construct a new proto descriptor for the given property.
func NewMessage(property *specs.Property) (*desc.MessageDescriptor, error) {
	builder, _, err := newMessage(
		make(map[string]*builder.MessageBuilder),
		property.Name,
		property,
	)
	if err != nil {
		return nil, err
	}

	if err := builder.TrySetName(property.Name); err != nil {
		return nil, err
	}

	return builder.Build()
}

func newMessage(builders map[string]*builder.MessageBuilder, name string, property *specs.Property) (*builder.MessageBuilder, bool, error) {
	if property.Template.Identifier != "" {
		existing, ok := builders[property.Template.Identifier]
		if ok {
			return existing, true, nil
		}
	}

	builder := builder.NewMessage(name)
	builders[property.Template.Identifier] = builder

	if err := constructMessage(builders, builder, property.Message); err != nil {
		return nil, false, err
	}

	return builder, false, nil
}

// ConstructMessage constructs a proto message of the given specs into the given message builders
func ConstructMessage(messageBuilder *builder.MessageBuilder, message specs.Message) (err error) {
	return constructMessage(make(map[string]*builder.MessageBuilder), messageBuilder, message)
}

func constructMessage(builders map[string]*builder.MessageBuilder, messageBuilder *builder.MessageBuilder, message specs.Message) (err error) {
	for _, property := range message {
		typed, err := ConstructFieldType(builders, property.Name+"Type", messageBuilder, property)
		if err != nil {
			return err
		}

		label := protobuffers.ProtoLabels[labels.Optional]
		if property.Type() == types.Array {
			label = protobuffers.Repeated
		}

		field := builder.NewField(property.Name, typed)
		field.SetJsonName(property.Name)
		field.SetLabel(label)
		field.SetComments(builder.Comments{
			LeadingComment: property.Description,
		})

		if err = messageBuilder.TryAddField(field.SetNumber(property.Position)); err != nil {
			return err
		}
	}

	return nil
}

// ConstructFieldType constructs a field constructor from the given property
func ConstructFieldType(builders map[string]*builder.MessageBuilder, name string, message *builder.MessageBuilder, property *specs.Property) (ft *builder.FieldType, err error) {
	switch {
	case property.Message != nil:
		nested, seen, err := newMessage(builders, name, property)
		if err != nil {
			return nil, err
		}

		if !seen {
			if err := message.TryAddNestedMessage(nested); err != nil {
				return nil, err
			}
		}

		return builder.FieldTypeMessage(nested), nil
	case property.Repeated != nil:
		template, err := property.Repeated.Template()
		if err != nil {
			return nil, err
		}

		// TODO: thrown a error when attempting to construct a nested array
		field := &specs.Property{
			Name:     property.Name,
			Template: template,
		}

		return ConstructFieldType(builders, name, message, field)
	case property.Enum != nil:
		enum := builder.NewEnum(property.Name + "Enum")

		for _, value := range property.Enum.Keys {
			eval := builder.NewEnumValue(value.Key)

			eval.SetNumber(value.Position)
			eval.SetComments(builder.Comments{
				LeadingComment: value.Description,
			})

			if err := enum.TryAddValue(eval); err != nil {
				return nil, err
			}
		}

		if err := message.TryAddNestedEnum(enum); err != nil {
			return nil, err
		}

		return builder.FieldTypeEnum(enum), nil
	case property.Scalar != nil:
		return builder.FieldTypeScalar(protobuffers.ProtoTypes[property.Scalar.Type]), nil
	}

	return nil, ErrInvalidFieldType{*property.Template}
}
